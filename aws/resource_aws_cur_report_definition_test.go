package aws

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/costandusagereportservice"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccAwsCurReportDefinition_basic(t *testing.T) {
	oldvar := os.Getenv("AWS_DEFAULT_REGION")
	os.Setenv("AWS_DEFAULT_REGION", "us-east-1")
	defer os.Setenv("AWS_DEFAULT_REGION", oldvar)

	resourceName := "aws_cur_report_definition.test"

	reportName := acctest.RandomWithPrefix("tf_acc_test")
	bucketName := fmt.Sprintf("tf-test-bucket-%d", acctest.RandInt())
	bucketRegion := "us-east-1"
	bucketPrefix := ""
	format := "textORcsv"
	compression := "GZIP"
	additionalArtifacts := []string{"REDSHIFT", "QUICKSIGHT"}
	refreshClosedReports := false
	reportVersioning := "CREATE_NEW_REPORT"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t); testAccPreCheckAWSCur(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAwsCurReportDefinitionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAwsCurReportDefinitionConfig_bucket(bucketName, bucketRegion),
			},
			{
				Config: testAccAwsCurReportDefinitionConfig_basic(reportName, bucketName, bucketPrefix, bucketRegion, format, compression, additionalArtifacts, refreshClosedReports, reportVersioning),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAwsCurReportDefinitionExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "report_name", reportName),
					resource.TestCheckResourceAttr(resourceName, "time_unit", "DAILY"),
					resource.TestCheckResourceAttr(resourceName, "format", format),
					resource.TestCheckResourceAttr(resourceName, "compression", compression),
					resource.TestCheckResourceAttr(resourceName, "additional_schema_elements.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "s3_bucket", bucketName),
					resource.TestCheckResourceAttr(resourceName, "s3_prefix", bucketPrefix),
					resource.TestCheckResourceAttr(resourceName, "s3_region", bucketRegion),
					resource.TestCheckResourceAttr(resourceName, "additional_artifacts.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "refresh_closed_reports", "false"),
					resource.TestCheckResourceAttr(resourceName, "report_versioning", reportVersioning),
				),
			},
		},
	})
}

func TestAccAwsCurReportDefinition_parquet(t *testing.T) {
	oldvar := os.Getenv("AWS_DEFAULT_REGION")
	os.Setenv("AWS_DEFAULT_REGION", "us-east-1")
	defer os.Setenv("AWS_DEFAULT_REGION", oldvar)

	resourceName := "aws_cur_report_definition.test"

	reportName := acctest.RandomWithPrefix("tf_acc_test")
	bucketName := fmt.Sprintf("tf-test-bucket-%d", acctest.RandInt())
	bucketRegion := "us-east-1"
	bucketPrefix := ""
	format := "Parquet"
	compression := "Parquet"
	additionalArtifacts := []string{}
	refreshClosedReports := false
	reportVersioning := "CREATE_NEW_REPORT"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t); testAccPreCheckAWSCur(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAwsCurReportDefinitionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAwsCurReportDefinitionConfig_bucket(bucketName, bucketRegion),
			},
			{
				Config: testAccAwsCurReportDefinitionConfig_basic(reportName, bucketName, bucketPrefix, bucketRegion, format, compression, additionalArtifacts, refreshClosedReports, reportVersioning),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAwsCurReportDefinitionExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "report_name", reportName),
					resource.TestCheckResourceAttr(resourceName, "time_unit", "DAILY"),
					resource.TestCheckResourceAttr(resourceName, "format", format),
					resource.TestCheckResourceAttr(resourceName, "compression", compression),
					resource.TestCheckResourceAttr(resourceName, "additional_schema_elements.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "s3_bucket", bucketName),
					resource.TestCheckResourceAttr(resourceName, "s3_prefix", bucketPrefix),
					resource.TestCheckResourceAttr(resourceName, "s3_region", bucketRegion),
					resource.TestCheckResourceAttr(resourceName, "refresh_closed_reports", "false"),
					resource.TestCheckResourceAttr(resourceName, "report_versioning", reportVersioning),
				),
			},
		},
	})
}

func TestAccAwsCurReportDefinition_athena(t *testing.T) {
	oldvar := os.Getenv("AWS_DEFAULT_REGION")
	os.Setenv("AWS_DEFAULT_REGION", "us-east-1")
	defer os.Setenv("AWS_DEFAULT_REGION", oldvar)

	resourceName := "aws_cur_report_definition.test"

	reportName := acctest.RandomWithPrefix("tf_acc_test")
	bucketName := fmt.Sprintf("tf-test-bucket-%d", acctest.RandInt())
	bucketRegion := "us-east-1"
	bucketPrefix := "data"
	format := "Parquet"
	compression := "Parquet"
	additionalArtifacts := []string{"ATHENA"}
	refreshClosedReports := false
	reportVersioning := "OVERWRITE_REPORT"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t); testAccPreCheckAWSCur(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAwsCurReportDefinitionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAwsCurReportDefinitionConfig_bucket(bucketName, bucketRegion),
			},
			{
				Config: testAccAwsCurReportDefinitionConfig_basic(reportName, bucketName, bucketPrefix, bucketRegion, format, compression, additionalArtifacts, refreshClosedReports, reportVersioning),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAwsCurReportDefinitionExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "report_name", reportName),
					resource.TestCheckResourceAttr(resourceName, "time_unit", "DAILY"),
					resource.TestCheckResourceAttr(resourceName, "format", format),
					resource.TestCheckResourceAttr(resourceName, "compression", compression),
					resource.TestCheckResourceAttr(resourceName, "additional_schema_elements.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "s3_bucket", bucketName),
					resource.TestCheckResourceAttr(resourceName, "s3_prefix", bucketPrefix),
					resource.TestCheckResourceAttr(resourceName, "s3_region", bucketRegion),
					resource.TestCheckResourceAttr(resourceName, "additional_artifacts.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "refresh_closed_reports", "false"),
					resource.TestCheckResourceAttr(resourceName, "report_versioning", reportVersioning),
				),
			},
		},
	})
}

func TestAccAwsCurReportDefinition_refresh(t *testing.T) {
	oldvar := os.Getenv("AWS_DEFAULT_REGION")
	os.Setenv("AWS_DEFAULT_REGION", "us-east-1")
	defer os.Setenv("AWS_DEFAULT_REGION", oldvar)

	resourceName := "aws_cur_report_definition.test"

	reportName := acctest.RandomWithPrefix("tf_acc_test")
	bucketName := fmt.Sprintf("tf-test-bucket-%d", acctest.RandInt())
	bucketRegion := "us-east-1"
	bucketPrefix := ""
	format := "textORcsv"
	compression := "GZIP"
	additionalArtifacts := []string{"REDSHIFT", "QUICKSIGHT"}
	refreshClosedReports := true
	reportVersioning := "CREATE_NEW_REPORT"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t); testAccPreCheckAWSCur(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAwsCurReportDefinitionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAwsCurReportDefinitionConfig_bucket(bucketName, bucketRegion),
			},
			{
				Config: testAccAwsCurReportDefinitionConfig_basic(reportName, bucketName, bucketPrefix, bucketRegion, format, compression, additionalArtifacts, refreshClosedReports, reportVersioning),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAwsCurReportDefinitionExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "report_name", reportName),
					resource.TestCheckResourceAttr(resourceName, "time_unit", "DAILY"),
					resource.TestCheckResourceAttr(resourceName, "format", format),
					resource.TestCheckResourceAttr(resourceName, "compression", compression),
					resource.TestCheckResourceAttr(resourceName, "additional_schema_elements.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "s3_bucket", bucketName),
					resource.TestCheckResourceAttr(resourceName, "s3_prefix", bucketPrefix),
					resource.TestCheckResourceAttr(resourceName, "s3_region", bucketRegion),
					resource.TestCheckResourceAttr(resourceName, "additional_artifacts.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "refresh_closed_reports", "true"),
					resource.TestCheckResourceAttr(resourceName, "report_versioning", reportVersioning),
				),
			},
		},
	})
}

func TestAccAwsCurReportDefinition_overwrite(t *testing.T) {
	oldvar := os.Getenv("AWS_DEFAULT_REGION")
	os.Setenv("AWS_DEFAULT_REGION", "us-east-1")
	defer os.Setenv("AWS_DEFAULT_REGION", oldvar)

	resourceName := "aws_cur_report_definition.test"

	reportName := acctest.RandomWithPrefix("tf_acc_test")
	bucketName := fmt.Sprintf("tf-test-bucket-%d", acctest.RandInt())
	bucketRegion := "us-east-1"
	bucketPrefix := ""
	format := "textORcsv"
	compression := "GZIP"
	additionalArtifacts := []string{"REDSHIFT", "QUICKSIGHT"}
	refreshClosedReports := false
	reportVersioning := "OVERWRITE_REPORT"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t); testAccPreCheckAWSCur(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAwsCurReportDefinitionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAwsCurReportDefinitionConfig_bucket(bucketName, bucketRegion),
			},
			{
				Config: testAccAwsCurReportDefinitionConfig_basic(reportName, bucketName, bucketPrefix, bucketRegion, format, compression, additionalArtifacts, refreshClosedReports, reportVersioning),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAwsCurReportDefinitionExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "report_name", reportName),
					resource.TestCheckResourceAttr(resourceName, "time_unit", "DAILY"),
					resource.TestCheckResourceAttr(resourceName, "format", format),
					resource.TestCheckResourceAttr(resourceName, "compression", compression),
					resource.TestCheckResourceAttr(resourceName, "additional_schema_elements.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "s3_bucket", bucketName),
					resource.TestCheckResourceAttr(resourceName, "s3_prefix", bucketPrefix),
					resource.TestCheckResourceAttr(resourceName, "s3_region", bucketRegion),
					resource.TestCheckResourceAttr(resourceName, "additional_artifacts.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "refresh_closed_reports", "false"),
					resource.TestCheckResourceAttr(resourceName, "report_versioning", reportVersioning),
				),
			},
		},
	})
}

func testAccCheckAwsCurReportDefinitionDestroy(s *terraform.State) error {
	conn := testAccProvider.Meta().(*AWSClient).costandusagereportconn

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "aws_cur_report_definition" {
			continue
		}
		reportName := rs.Primary.ID
		matchingReportDefinition, err := describeCurReportDefinition(conn, reportName)
		if err != nil {
			return err
		}
		if matchingReportDefinition != nil {
			return fmt.Errorf("Report Definition still exists: %q", rs.Primary.ID)
		}
	}
	return nil

}

func testAccCheckAwsCurReportDefinitionExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		conn := testAccProvider.Meta().(*AWSClient).costandusagereportconn

		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Resource not found: %s", resourceName)
		}
		reportName := rs.Primary.ID
		matchingReportDefinition, err := describeCurReportDefinition(conn, reportName)
		if err != nil {
			return err
		}
		if matchingReportDefinition == nil {
			return fmt.Errorf("Report Definition does not exist: %q", rs.Primary.ID)
		}
		return nil
	}
}

func testAccPreCheckAWSCur(t *testing.T) {
	conn := testAccProvider.Meta().(*AWSClient).costandusagereportconn

	input := &costandusagereportservice.DescribeReportDefinitionsInput{
		MaxResults: aws.Int64(5),
	}

	_, err := conn.DescribeReportDefinitions(input)

	if testAccPreCheckSkipError(err) || isAWSErr(err, "AccessDeniedException", "linked account is not allowed to modify report preference") {
		t.Skipf("skipping acceptance testing: %s", err)
	}

	if err != nil {
		t.Fatalf("unexpected PreCheck error: %s", err)
	}
}

// note: bucket needs to exist first due to ValidationException thrown if policy not yet applied
func testAccAwsCurReportDefinitionConfig_bucket(bucketName string, bucketRegion string) string {
	return fmt.Sprintf(`
provider "aws" {
  region = "us-east-1"
}

resource "aws_s3_bucket" "test" {
  bucket        = "%[1]s"
  acl           = "private"
  force_destroy = true
	region        = "%[2]s"
}

resource "aws_s3_bucket_policy" "test" {
	bucket = "${aws_s3_bucket.test.id}"

	policy = <<POLICY
{
    "Version": "2008-10-17",
    "Id": "s3policy",
    "Statement": [
        {
            "Sid": "AllowCURBillingACLPolicy",
            "Effect": "Allow",
            "Principal": {
                "AWS": "arn:aws:iam::386209384616:root"
            },
            "Action": [
                "s3:GetBucketAcl",
                "s3:GetBucketPolicy"
            ],
            "Resource": "${aws_s3_bucket.test.arn}"
        },
        {
            "Sid": "AllowCURPutObject",
            "Effect": "Allow",
            "Principal": {
                "AWS": "arn:aws:iam::386209384616:root"
            },
            "Action": "s3:PutObject",
            "Resource": "${aws_s3_bucket.test.arn}/*"
        }
    ]
}
POLICY
}
`, bucketName, bucketRegion)
}

// note: cur report definitions are currently only supported in us-east-1
func testAccAwsCurReportDefinitionConfig_basic(reportName string, bucketName string, bucketPrefix string, bucketRegion string, format string, compression string, additionalArtifacts []string, refreshClosedReports bool, reportVersioning string) string {
	artifactsStr := strings.Join(additionalArtifacts, "\", \"")

	if len(additionalArtifacts) > 0 {
		artifactsStr = fmt.Sprintf("additional_artifacts       = [\"%s\"]", artifactsStr)
	} else {
		artifactsStr = ""
	}

	return fmt.Sprintf(`
%[2]s

resource "aws_cur_report_definition" "test" {
  report_name                = "%[1]s"
  time_unit                  = "DAILY"
  format                     = "%[4]s"
  compression                = "%[5]s"
  additional_schema_elements = ["RESOURCES"]
  s3_bucket                  = "${aws_s3_bucket.test.id}"
  s3_prefix                  = "%[3]s"
  s3_region                  = "${aws_s3_bucket.test.region}"
	%[6]s
	refresh_closed_reports     = %[7]t
	report_versioning          = "%[8]s"
}
`, reportName, testAccAwsCurReportDefinitionConfig_bucket(bucketName, bucketRegion), bucketPrefix, format, compression, artifactsStr, refreshClosedReports, reportVersioning)
}
