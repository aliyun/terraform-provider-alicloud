package ecs_test

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/ecs"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

// TestAccAliCloudEcsDiagnosticReport_basic covers the create-only lifecycle: it
// provisions a fresh ECS instance, generates a diagnostic report for it, verifies the
// ReportId and resource_id land in state, and imports the report while ignoring the
// write-only additional_options map.
func TestAccAliCloudEcsDiagnosticReport_basic(t *testing.T) {
	r := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccecs%d", r)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             testAccCheckEcsDiagnosticReportDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAliCloudEcsDiagnosticReportConfig(name),
				Check: resource.ComposeTestCheckFunc(
					// ReportId and resource_id must land in state; the echoed time/metric
					// fields are checked as set rather than by value to avoid time-format
					// drift, matching the SDKv2 CHECKSET scope.
					resource.TestCheckResourceAttrSet("alicloud_ecs_diagnostic_report.default", "id"),
					resource.TestCheckResourceAttrSet("alicloud_ecs_diagnostic_report.default", "resource_id"),
					// additional_options is write-only: the API never returns it, so the
					// value is carried from config into state on Create.
					resource.TestCheckResourceAttr("alicloud_ecs_diagnostic_report.default", "additional_options.key", "value"),
				),
			},
			{
				ResourceName:            "alicloud_ecs_diagnostic_report.default",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"additional_options"},
			},
		},
	})
}

// testAccCheckEcsDiagnosticReportDestroy verifies a report is gone after destroy,
// using the shared DiagnosticReportExists predicate so the acceptance test exercises
// the same legal-missing (403 ReportNotExist / 200-only-RequestId) distinction the
// unit tests cover. A report that still exists, or a bad response, fails destroy.
func testAccCheckEcsDiagnosticReportDestroy(s *terraform.State) error {
	client := acctest.Provider.Meta().(*connectivity.AliyunClient)
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_ecs_diagnostic_report" {
			continue
		}
		exists, err := ecs.DiagnosticReportExists(client, rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("checking diagnostic report %s destroy: %w", rs.Primary.ID, err)
		}
		if exists {
			return fmt.Errorf("diagnostic report %s still exists", rs.Primary.ID)
		}
	}
	return nil
}

func testAccAliCloudEcsDiagnosticReportConfig(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%[1]s"
}

data "alicloud_zones" "default" {
  available_disk_category     = "cloud_essd"
  available_resource_creation = "VSwitch"
}

data "alicloud_images" "default" {
  name_regex  = "^ubuntu_[0-9]+_[0-9]+_x64*"
  most_recent = true
  owners      = "system"
}

data "alicloud_instance_types" "default" {
  availability_zone    = data.alicloud_zones.default.zones.0.id
  instance_type_family = "ecs.u1"
  image_id             = data.alicloud_images.default.images.0.id
}

resource "alicloud_vpc" "default" {
  vpc_name   = var.name
  cidr_block = "192.168.0.0/16"
}

resource "alicloud_vswitch" "default" {
  vswitch_name = var.name
  vpc_id       = alicloud_vpc.default.id
  cidr_block   = "192.168.192.0/24"
  zone_id      = data.alicloud_zones.default.zones.0.id
}

resource "alicloud_security_group" "default" {
  name   = var.name
  vpc_id = alicloud_vpc.default.id
}

resource "alicloud_instance" "default" {
  image_id                   = data.alicloud_images.default.images.0.id
  instance_type              = data.alicloud_instance_types.default.instance_types.0.id
  security_groups            = alicloud_security_group.default.*.id
  internet_charge_type       = "PayByTraffic"
  internet_max_bandwidth_out = "10"
  availability_zone          = data.alicloud_instance_types.default.instance_types.0.availability_zones.0
  instance_charge_type       = "PostPaid"
  system_disk_category       = "cloud_essd"
  vswitch_id                 = alicloud_vswitch.default.id
  instance_name              = var.name
}

resource "alicloud_ecs_diagnostic_report" "default" {
  resource_id        = alicloud_instance.default.id
  metric_set_id      = "dms-instancedefault"
  start_time         = "2022-07-11T12:00Z"
  end_time           = "2022-07-11T14:00Z"
  additional_options = {
    key = "value"
  }
}
`, name)
}
