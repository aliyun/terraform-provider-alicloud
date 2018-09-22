package alicloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccAlicloudFcServicesDataSource_basic(t *testing.T) {
	randInt := acctest.RandInt()
	serviceName := fmt.Sprintf("tf-testacc-fc-service-ds-basic-%d", randInt)
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlicloudFcServicesDataSourceBasic(randInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_fc_services.services"),
					resource.TestCheckResourceAttr("data.alicloud_fc_services.services", "services.#", "1"),
					resource.TestCheckResourceAttrSet("data.alicloud_fc_services.services", "services.0.id"),
					resource.TestCheckResourceAttr("data.alicloud_fc_services.services", "services.0.name", serviceName),
					resource.TestCheckResourceAttr("data.alicloud_fc_services.services", "services.0.description", serviceName+"-description"),
					resource.TestCheckResourceAttrSet("data.alicloud_fc_services.services", "services.0.role"),
					resource.TestCheckResourceAttr("data.alicloud_fc_services.services", "services.0.internet_access", "true"),
					resource.TestCheckResourceAttrSet("data.alicloud_fc_services.services", "services.0.creation_time"),
					resource.TestCheckResourceAttrSet("data.alicloud_fc_services.services", "services.0.last_modification_time"),
					resource.TestCheckResourceAttr("data.alicloud_fc_services.services", "services.0.log_config.#", "1"),
					resource.TestCheckResourceAttr("data.alicloud_fc_services.services", "services.0.log_config.0.project", serviceName+"-project"),
					resource.TestCheckResourceAttr("data.alicloud_fc_services.services", "services.0.log_config.0.logstore", serviceName+"-store"),
					resource.TestCheckResourceAttr("data.alicloud_fc_services.services", "services.0.vpc_config.#", "1"),
					resource.TestCheckResourceAttrSet("data.alicloud_fc_services.services", "services.0.vpc_config.0.vpc_id"),
					resource.TestCheckResourceAttr("data.alicloud_fc_services.services", "services.0.vpc_config.0.vswitch_ids.#", "1"),
					resource.TestCheckResourceAttrSet("data.alicloud_fc_services.services", "services.0.vpc_config.0.security_group_id"),
				),
			},
		},
	})
}

func testAccCheckAlicloudFcServicesDataSourceBasic(randInt int) string {
	return fmt.Sprintf(`
variable "name" {
	default = "tf-testacc-fc-service-ds-basic-%d"
}

data "alicloud_zones" "zones" {
    available_resource_creation = "VSwitch"
}

resource "alicloud_vpc" "sample_vpc" {
	name = "${var.name}"
	cidr_block = "172.16.0.0/16"
}

resource "alicloud_vswitch" "sample_vswitch" {
	name = "${var.name}"
	cidr_block = "172.16.0.0/24"
	availability_zone = "${data.alicloud_zones.zones.zones.0.id}"
	vpc_id = "${alicloud_vpc.sample_vpc.id}"
}

resource "alicloud_security_group" "sample_group" {
	name = "${var.name}"
	vpc_id = "${alicloud_vpc.sample_vpc.id}"
}

resource "alicloud_log_project" "sample_log_project" {
    name = "${var.name}-project"
}

resource "alicloud_log_store" "sample_log_store" {
    project = "${alicloud_log_project.sample_log_project.name}"
    name = "${var.name}-store"
}

resource "alicloud_ram_role" "sample_role" {
    name = "${var.name}"
    document = <<DEFINITION
    {
        "Statement": [
            {
                "Action": "sts:AssumeRole",
                "Effect": "Allow",
                "Principal": {
                    "Service": [
                        "fc.aliyuncs.com"
                    ]
                }
            }
        ],
        "Version": "1"
    }
    DEFINITION
    description = "this is a test"
    force = true
}

resource "alicloud_ram_policy" "sample_vpc_policy" {
    name = "${var.name}"
    document = <<DEFINITION
    {
  "Version": "1",
  "Statement": [
    {
      "Action": "vpc:*",
      "Resource": "*",
      "Effect": "Allow"
    },
    {
      "Action": [
        "ecs:*NetworkInterface*"
      ],
      "Resource": "*",
      "Effect": "Allow"
    }
  ]
}
    DEFINITION
}

resource "alicloud_ram_role_policy_attachment" "sample_attachment" {
    role_name = "${alicloud_ram_role.sample_role.name}"
    policy_name = "AliyunLogFullAccess"
    policy_type = "System"
}

resource "alicloud_ram_role_policy_attachment" "sample_attachment_vpc" {
    role_name = "${alicloud_ram_role.sample_role.name}"
    policy_name = "${alicloud_ram_policy.sample_vpc_policy.name}"
    policy_type = "Custom"
}

resource "alicloud_fc_service" "sample_service" {
    name = "${var.name}"
    description = "${var.name}-description"
    log_config {
	    project = "${alicloud_log_project.sample_log_project.name}"
	    logstore = "${alicloud_log_store.sample_log_store.name}"
    }
    vpc_config {
        vswitch_ids = ["${alicloud_vswitch.sample_vswitch.id}"]
        security_group_id = "${alicloud_security_group.sample_group.id}"
    }
    role = "${alicloud_ram_role.sample_role.arn}"
    depends_on = ["alicloud_ram_role_policy_attachment.sample_attachment", "alicloud_ram_role_policy_attachment.sample_attachment_vpc"]
    internet_access = true
}

data "alicloud_fc_services" "services" {
    name_regex = "${alicloud_fc_service.sample_service.name}"
}
`, randInt)
}
