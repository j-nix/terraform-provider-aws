package aws

import (
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/route53resolver"
	"github.com/hashicorp/terraform/helper/schema"
	"log"
)

func dataSourceAwsRoute53ResolverEndpoint() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAwsRoute53ResolverEndpointRead,

		Schema: map[string]*schema.Schema{
			"filter": dataSourceFiltersSchema(),
			"id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"arn": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"host_vpc_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"ip_addresses": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"security_group_ids": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func dataSourceAwsRoute53ResolverEndpointRead(d *schema.ResourceData, meta interface{}) error {
	var resolverId string
	conn := meta.(*AWSClient).route53resolverconn
	input := route53resolver.ListResolverEndpointsInput{}
	//conn.ListResolverEndpoints()

	if v, ok := d.GetOk("filter"); ok {
		input.Filters = buildAwsRoute53ResolverDataSourceFilters(v.(*schema.Set))
	}

	if id, ok := d.GetOk("id"); ok {
		resolverId = fmt.Sprintf("%v", id)
	}

	log.Printf("[DEBUG] Reading Route53 Resolver Endpoints: %s", input)
	output, err := conn.ListResolverEndpoints(&input)

	if err != nil {
		return fmt.Errorf("error reading Route53 Resolver Endpoints: %s", err)
	}

	if output == nil || len(output.ResolverEndpoints) == 0 {
		return errors.New("error reading Route53 Resolver Endpoints: no results found")
	}

	if len(output.ResolverEndpoints) > 1 {
		return errors.New("error reading Route53 Resolver Endpoints: multiple results found, try adjusting search criteria")
	}

	re := output.ResolverEndpoints[0]
	if re == nil {
		return errors.New("error reading Route53 Resolver Endpoint: empty result")
	}

	log.Printf("[DEBUG] Reading Route53 Resolver Endpoint IP Addresses: %s", resolverId)
	ips, err := endpointIps(&resolverId, conn)

	if err != nil {
		return fmt.Errorf("error reading Route53 Resolver Endpoint IP Addresses: %s", err)
	}

	if ips == nil {
		return errors.New("error reading Route53 Resolver Endpoint IP Addresses: no results found")
	}

	d.Set("arn", re.Arn)
	d.Set("host_vpc_id", re.HostVPCId)
	d.Set("ip_addresses", ips)
	d.Set("name", re.Name)
	d.Set("security_group_ids", re.SecurityGroupIds)
	d.SetId(aws.StringValue(re.Id))


	return nil
}

func endpointIps(id *string, conn *route53resolver.Route53Resolver) ([]string, error) {
	req := &route53resolver.ListResolverEndpointIpAddressesInput{
		ResolverEndpointId: id,
	}

	log.Printf("[DEBUG] Reading IPAddresses: %s", req)
	resp, err := conn.ListResolverEndpointIpAddresses(req)
	if err != nil {
		return nil, err
	}
	if resp == nil || len(resp.IpAddresses) == 0 {
		return nil, fmt.Errorf("no matching IP address found: %s", req)
	}

	ipResp := resp.IpAddresses
	log.Printf("[DEBUG] IP address response: %s", ipResp)

	list := make([]string, len(ipResp))
	for i, address := range ipResp {
		list[i] = *address.Ip
	}

	return list, nil
}
