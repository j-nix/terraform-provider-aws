package aws

import (
	"fmt"
	"github.com/aws/aws-sdk-go/service/route53resolver"
	"github.com/hashicorp/terraform/helper/schema"
	"log"
	"time"
)

func dataSourceAwsRoute53Resolver() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAwsRoute53ResolverRead,

		Schema: map[string]*schema.Schema{
			"resolver_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"ips": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func dataSourceAwsRoute53ResolverRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*AWSClient).route53resolverconn

	var resolverId string
	id, idExists := d.GetOk("resolver_id")
	if !idExists {
		return fmt.Errorf("Resolver ID doesn't exist")
	} else {
		resolverId = fmt.Sprintf("%v", id)
	}

	ips, err := endpointIps(&resolverId, conn)
	if err != nil {
		return fmt.Errorf("Error returning IPs for Route53Resolver endpoint %v", err)
	}

	d.SetId(time.Now().UTC().String())
	err = d.Set("ips", ips)
	if err != nil {
		return err
	}

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
