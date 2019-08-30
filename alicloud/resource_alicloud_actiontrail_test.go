package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/actiontrail"

	"github.com/hashicorp/terraform/helper/acctest"

	"log"
	"strings"
	"time"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func init() {
	resource.AddTestSweepers(
		"alicloud_actiontrail",
		&resource.Sweeper{
			Name: "alicloud_actiontrail",
			F:    testSweepActiontrail,
		})
}

func testSweepActiontrail(region string) error {
	if testSweepPreCheckWithRegions(region, false, connectivity.ApiGatewayNoSupportedRegions) {
		log.Printf("[INFO] Skipping API Gateway unsupported region: %s", region)
		return nil
	}
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting Alicloud client: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)

	prefixes := []string{
		"tf-testacc",
		"tf_testacc",
	}

	request := actiontrail.CreateDescribeTrailsRequest()
	raw, err := client.WithActionTrailClient(func(actiontrailClient *actiontrail.Client) (interface{}, error) {
		return actiontrailClient.DescribeTrails(request)
	})
	if err != nil {
		return fmt.Errorf("Error Describe Apis: %s", err)
	}
	response := raw.(*actiontrail.DescribeTrailsResponse)

	swept := false

	for _, v := range response.TrailList {
		name := v.Name
		skip := true
		for _, prefix := range prefixes {
			if strings.HasPrefix(strings.ToLower(name), strings.ToLower(prefix)) {
				skip = false
				break
			}
		}
		if skip {
			log.Printf("[INFO] Skipping api: %s", name)
			continue
		}
		swept = true

		log.Printf("[INFO] Deleting Api: %s", name)

		request := actiontrail.CreateDeleteTrailRequest()
		request.Name = name

		_, err := client.WithActionTrailClient(func(actiontrailClient *actiontrail.Client) (interface{}, error) {
			return actiontrailClient.DeleteTrail(request)
		})

		if err != nil {
			log.Printf("[ERROR] Failed to delete Api (%s): %s", name, err)
		}
	}
	if swept {
		time.Sleep(5 * time.Second)
	}
	return nil
}

func TestAccAlicloudActionTrailUpdate(t *testing.T) {

	num := acctest.RandInt()
	num_role := acctest.RandInt()
	num_bucket := acctest.RandInt()

	var v actiontrail.TrailListItem
	resourceId := "alicloud_actiontrail.default"

	basicMap := map[string]string{
		"name":               "tf-testacc-actiontrail",
		"event_rw":           "Write",
		"oss_bucket_name":    fmt.Sprintf("tf-testacc-actiontrail-%v", num),
		"role_name":          fmt.Sprintf("tf-testacc-actiontrail-%v", num),
		"oss_key_prefix":     "",
		"sls_project_arn":    "",
		"sls_write_role_arn": "",
	}

	ra := resourceAttrInit(resourceId, basicMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ActionTrailService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeActionTrail")
	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, connectivity.ActiontrailNoSkipRegions)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  testActionTrailDestroy,
		Steps: []resource.TestStep{
			{
				Config: testActionTrailBasicConfig(num),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(nil),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testActionTrailUpdateConfig_event_rw(num),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"event_rw": "All",
					}),
				),
			},
			{
				Config: testActionTrailUpdateConfig_oss_key_prefix(num),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"oss_key_prefix": "at-product-account-audit-B",
					}),
				),
			},
			{
				Config: testActionTrailUpdateConfig_role(num, num_role),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"role_name": fmt.Sprintf("tf-testacc-actiontrail-%v", num_role),
					}),
				),
			},
			{
				Config: testActionTrailUpdateConfig_bucket(num_bucket, num_role),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"oss_bucket_name": fmt.Sprintf("tf-testacc-actiontrail-%v", num_bucket),
					}),
				),
			},
			{
				Config: testActionTrailUpdateConfig_sls_project_arn(num_bucket, num_role),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"sls_project_arn":    CHECKSET,
						"sls_write_role_arn": CHECKSET,
					}),
				),
			},
			{
				Config: testActionTrailBasicConfig(num),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(basicMap),
				),
			},
		},
	})
}

func testActionTrailDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*connectivity.AliyunClient)
	trailService := ActionTrailService{client}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_actiontrail" {
			continue
		}
		_, err := trailService.DescribeActionTrail(rs.Primary.ID)
		if err != nil {
			if NotFoundError(err) {
				continue
			}
			return WrapError(fmt.Errorf("Describe Action Trail error %#v", err))
		}
	}
	return nil
}

func testActionTrailBasicConfig(randInt int) string {
	return fmt.Sprintf(`
variable "name" {
    default = "tf-testacc-actiontrail-%v"
}


resource "alicloud_ram_role" "default" {
	  name = "${var.name}"
	  document = <<EOF
			{
			  "Statement": [
				{
				  "Action": "sts:AssumeRole",
				  "Effect": "Allow",
				  "Principal": {
					"Service": [
					  "actiontrail.aliyuncs.com",
					  "oss.aliyuncs.com"
					]
				  }
				}
			  ],
			  "Version": "1"
			}
		  EOF
	  description = "this is a test"
	  force = "true"
}

resource "alicloud_oss_bucket" "default" {
    bucket  = "${var.name}"
}

resource "alicloud_ram_policy" "default" {
	  name = "${var.name}"
	  document = <<EOF
		{
		  "Statement": [
			{
			  "Action": [
				"*"
			  ],
			  "Effect": "Allow",
			  "Resource": [
				"acs:oss:*:*:${alicloud_oss_bucket.default.id}",
				"acs:oss:*:*:${alicloud_oss_bucket.default.id}"
			  ]
			}
		  ],
			"Version": "1"
		}
	  EOF
	  description = "this is a policy test"
	  force = true
	}

	resource "alicloud_ram_role_policy_attachment" "default" {
    	  policy_name = "${alicloud_ram_policy.default.name}"
    	  role_name = "${alicloud_ram_role.default.name}"
    	  policy_type = "${alicloud_ram_policy.default.type}"
	}
	
resource "alicloud_actiontrail" "default" {
	name = "tf-testacc-actiontrail"
	oss_bucket_name = "${alicloud_oss_bucket.default.id}"
	role_name = "${alicloud_ram_role_policy_attachment.default.role_name}"
}
`, randInt)
}

func testActionTrailUpdateConfig_event_rw(randInt int) string {

	return fmt.Sprintf(`
variable "name" {
    default = "tf-testacc-actiontrail-%v"
}

	resource "alicloud_ram_role" "default" {
	  name = "${var.name}"
	  document = <<EOF
			{
			  "Statement": [
				{
				  "Action": "sts:AssumeRole",
				  "Effect": "Allow",
				  "Principal": {
					"Service": [
					  "actiontrail.aliyuncs.com",
					  "oss.aliyuncs.com"
					]
				  }
				}
			  ],
			  "Version": "1"
			}
		  EOF
	  description = "this is a test"
	  force = true
}

resource "alicloud_oss_bucket" "default" {
    bucket  = "${var.name}"
}

resource "alicloud_ram_policy" "default" {
	  name = "${var.name}"
	  document = <<EOF
		{
		  "Statement": [
			{
			  "Action": [
				"*"
			  ],
			  "Effect": "Allow",
			  "Resource": [
				"acs:oss:*:*:${alicloud_oss_bucket.default.id}",
				"acs:oss:*:*:${alicloud_oss_bucket.default.id}"
			  ]
			}
		  ],
			"Version": "1"
		}
	  EOF
	  description = "this is a policy test"
	  force = true
	}

	resource "alicloud_ram_role_policy_attachment" "default" {
    	  policy_name = "${alicloud_ram_policy.default.name}"
    	  role_name = "${alicloud_ram_role.default.name}"
    	  policy_type = "${alicloud_ram_policy.default.type}"
    	}

	resource "alicloud_actiontrail" "default" {
		name = "tf-testacc-actiontrail"
		event_rw = "All"
		oss_bucket_name = "${alicloud_oss_bucket.default.id}"
		role_name = "${alicloud_ram_role_policy_attachment.default.role_name}"
	}
`, randInt)

}

func testActionTrailUpdateConfig_oss_key_prefix(randInt int) string {
	return fmt.Sprintf(`
variable "name" {
    default = "tf-testacc-actiontrail-%v"
}

	resource "alicloud_ram_role" "default" {
	  name = "${var.name}"
	  document = <<EOF
			{
			  "Statement": [
				{
				  "Action": "sts:AssumeRole",
				  "Effect": "Allow",
				  "Principal": {
					"Service": [
					  "actiontrail.aliyuncs.com",
					  "oss.aliyuncs.com"
					]
				  }
				}
			  ],
			  "Version": "1"
			}
		  EOF
	  description = "this is a test"
	  force = true
}

resource "alicloud_oss_bucket" "default" {
    bucket  = "${var.name}"
}

resource "alicloud_ram_policy" "default" {
	  name = "${var.name}"
	  document = <<EOF
		{
		  "Statement": [
			{
			  "Action": [
				"*"
			  ],
			  "Effect": "Allow",
			  "Resource": [
				"acs:oss:*:*:${alicloud_oss_bucket.default.id}",
				"acs:oss:*:*:${alicloud_oss_bucket.default.id}"
			  ]
			}
		  ],
			"Version": "1"
		}
	  EOF
	  description = "this is a policy test"
	  force = true
	}

	resource "alicloud_ram_role_policy_attachment" "default" {
    	  policy_name = "${alicloud_ram_policy.default.name}"
    	  role_name = "${alicloud_ram_role.default.name}"
    	  policy_type = "${alicloud_ram_policy.default.type}"
    	}

	resource "alicloud_actiontrail" "default" {
		name = "tf-testacc-actiontrail"
		event_rw = "All"
		oss_bucket_name = "${alicloud_oss_bucket.default.id}"
		role_name = "${alicloud_ram_role_policy_attachment.default.role_name}"
		oss_key_prefix = "at-product-account-audit-B"
	}
`, randInt)
}

func testActionTrailUpdateConfig_role(rand_1, rand_2 int) string {
	return fmt.Sprintf(`
variable "name" {
    default = "tf-testacc-actiontrail-%v"
}

	resource "alicloud_ram_role" "default" {
	  name = "tf-testacc-actiontrail-%v"
	  document = <<EOF
			{
			  "Statement": [
				{
				  "Action": "sts:AssumeRole",
				  "Effect": "Allow",
				  "Principal": {
					"Service": [
					  "actiontrail.aliyuncs.com",
					  "oss.aliyuncs.com"
					]
				  }
				}
			  ],
			  "Version": "1"
			}
		  EOF
	  description = "this is a test"
	  force = true
}

resource "alicloud_oss_bucket" "default" {
    bucket  = "${var.name}"
}

resource "alicloud_ram_policy" "default" {
	  name = "${var.name}"
	  document = <<EOF
		{
		  "Statement": [
			{
			  "Action": [
				"*"
			  ],
			  "Effect": "Allow",
			  "Resource": [
				"acs:oss:*:*:${alicloud_oss_bucket.default.id}",
				"acs:oss:*:*:${alicloud_oss_bucket.default.id}"
			  ]
			}
		  ],
			"Version": "1"
		}
	  EOF
	  description = "this is a policy test"
	  force = true
	}

	resource "alicloud_ram_role_policy_attachment" "default" {
    	  policy_name = "${alicloud_ram_policy.default.name}"
    	  role_name = "${alicloud_ram_role.default.name}"
    	  policy_type = "${alicloud_ram_policy.default.type}"
    	}

	resource "alicloud_actiontrail" "default" {
		name = "tf-testacc-actiontrail"
		event_rw = "All"
		oss_bucket_name = "${alicloud_oss_bucket.default.id}"
		role_name = "${alicloud_ram_role_policy_attachment.default.role_name}"
		oss_key_prefix = "at-product-account-audit-B"
	}
`, rand_1, rand_2)
}

func testActionTrailUpdateConfig_bucket(rand_1, rand_2 int) string {
	return fmt.Sprintf(`
variable "name" {
    default = "tf-testacc-actiontrail-%v"
}

	resource "alicloud_ram_role" "default" {
	  name = "tf-testacc-actiontrail-%v"
	  document = <<EOF
			{
			  "Statement": [
				{
				  "Action": "sts:AssumeRole",
				  "Effect": "Allow",
				  "Principal": {
					"Service": [
					  "actiontrail.aliyuncs.com",
					  "oss.aliyuncs.com"
					]
				  }
				}
			  ],
			  "Version": "1"
			}
		  EOF
	  description = "this is a test"
	  force = true
}

resource "alicloud_oss_bucket" "default" {
    bucket  = "${var.name}"
}

resource "alicloud_ram_policy" "default" {
	  name = "${var.name}"
	  document = <<EOF
		{
		  "Statement": [
			{
			  "Action": [
				"*"
			  ],
			  "Effect": "Allow",
			  "Resource": [
				"acs:oss:*:*:${alicloud_oss_bucket.default.id}",
				"acs:oss:*:*:${alicloud_oss_bucket.default.id}"
			  ]
			}
		  ],
			"Version": "1"
		}
	  EOF
	  description = "this is a policy test"
	  force = true
	}

	resource "alicloud_ram_role_policy_attachment" "default" {
    	  policy_name = "${alicloud_ram_policy.default.name}"
    	  role_name = "${alicloud_ram_role.default.name}"
    	  policy_type = "${alicloud_ram_policy.default.type}"
    	}

	resource "alicloud_actiontrail" "default" {
		name = "tf-testacc-actiontrail"
		event_rw = "All"
		oss_bucket_name = "${alicloud_oss_bucket.default.id}"
		role_name = "${alicloud_ram_role_policy_attachment.default.role_name}"
		oss_key_prefix = "at-product-account-audit-B"
	}
`, rand_1, rand_2)
}

func testActionTrailUpdateConfig_sls_project_arn(rand_1, rand_2 int) string {
	return fmt.Sprintf(`
variable "name" {
    default = "tf-testacc-actiontrail-%v"
}

data "alicloud_regions" "current_region" {
  current = true
}
data "alicloud_account" "current" {
}

resource "alicloud_log_project" "foo" {
  name = "${var.name}"
  description = "tf unit test"
}

resource "alicloud_ram_role" "default" {
	  name = "tf-testacc-actiontrail-%v"
	  document = <<EOF
			{
			  "Statement": [
				{
				  "Action": "sts:AssumeRole",
				  "Effect": "Allow",
				  "Principal": {
					"Service": [
					  "actiontrail.aliyuncs.com",
					  "oss.aliyuncs.com"
					]
				  }
				}
			  ],
			  "Version": "1"
			}
		  EOF
	  description = "this is a test"
	  force = true
}

resource "alicloud_oss_bucket" "default" {
    bucket  = "${var.name}"
}

resource "alicloud_ram_policy" "default" {
	  name = "${var.name}"
	  document = <<EOF
		{
		  "Statement": [
			{
			  "Action": [
				"*"
			  ],
			  "Effect": "Allow",
			  "Resource": [
				"acs:oss:*:*:${alicloud_oss_bucket.default.id}",
				"acs:oss:*:*:${alicloud_oss_bucket.default.id}"
			  ]
			}
		  ],
			"Version": "1"
		}
	  EOF
	  description = "this is a policy test"
	  force = true
	}

	resource "alicloud_ram_role_policy_attachment" "default" {
    	  policy_name = "${alicloud_ram_policy.default.name}"
    	  role_name = "${alicloud_ram_role.default.name}"
    	  policy_type = "${alicloud_ram_policy.default.type}"
    	}

	resource "alicloud_actiontrail" "default" {
		name = "tf-testacc-actiontrail"
		event_rw = "All"
		oss_bucket_name = "${alicloud_oss_bucket.default.id}"
		role_name = "${alicloud_ram_role_policy_attachment.default.role_name}"
		oss_key_prefix = "at-product-account-audit-B"
	    sls_project_arn = "acs:log:${data.alicloud_regions.current_region.regions.0.id}:${data.alicloud_account.current.id}:project/${alicloud_log_project.foo.name}"
		sls_write_role_arn = "${alicloud_ram_role_policy_attachment.default.role_name}"
	}
`, rand_1, rand_2)
}
