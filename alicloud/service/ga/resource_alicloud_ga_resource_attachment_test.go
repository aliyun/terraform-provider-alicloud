package ga_test

import (
	"fmt"
	"regexp"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccAliCloudGaResourceAttachment_basic0(t *testing.T) {
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAccGaResourceAttachment%d", rand)
	resourceId := "alicloud_ga_resource_attachment.default"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             testAccCheckAliCloudGaResourceAttachmentDestroy,
		Steps: []resource.TestStep{
			{
				Config:      testAccAliCloudGaResourceAttachmentConfigInvalidType(name),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`associated_resource_type`),
			},
			{
				Config: testAccAliCloudGaResourceAttachmentConfig(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceId, "id"),
					resource.TestCheckResourceAttrPair(resourceId, "accelerator_id", "alicloud_ga_accelerator.default", "id"),
					resource.TestCheckResourceAttr(resourceId, "associated_resource_type", "WAF"),
					resource.TestCheckResourceAttr(resourceId, "associated_resource_region_id", "cn-hangzhou"),
					resource.TestCheckResourceAttrPair(resourceId, "associated_resource_id", "data.alicloud_wafv3_instances.default", "ids.0"),
					resource.TestCheckResourceAttr(resourceId, "associated_mode", "Associated"),
					resource.TestCheckResourceAttr(resourceId, "state", "active"),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckAliCloudGaResourceAttachmentDestroy(s *terraform.State) error {
	meta := acctest.Provider.Meta()
	if meta == nil {
		return fmt.Errorf("provider metadata is not available for CheckDestroy")
	}
	client := meta.(*connectivity.AliyunClient)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_ga_resource_attachment" {
			continue
		}
		parts := strings.Split(rs.Primary.ID, ":")
		if len(parts) != 4 {
			return fmt.Errorf("unexpected resource attachment ID %q", rs.Primary.ID)
		}
		request := map[string]interface{}{
			"AcceleratorId":              parts[0],
			"RegionId":                   client.RegionId,
			"AssociatedResourceType":     parts[1],
			"AssociatedResourceRegionId": parts[2],
			"AssociatedResourceId":       parts[3],
		}
		response, err := client.RpcPost("Ga", "2019-11-20", "GetGlobalAcceleratorResources", nil, request, true)
		if err != nil {
			if strings.Contains(err.Error(), "NotExist.Accelerator") {
				continue
			}
			return fmt.Errorf("checking resource attachment %s: %s", rs.Primary.ID, err)
		}
		records, _ := response["AssociatedResources"].([]interface{})
		for _, item := range records {
			record, ok := item.(map[string]interface{})
			if !ok {
				continue
			}
			if recordAcceleratorId, ok := record["AcceleratorId"].(string); ok && recordAcceleratorId != "" && recordAcceleratorId != parts[0] {
				continue
			}
			if fmt.Sprint(record["AssociatedResourceType"]) == parts[1] &&
				fmt.Sprint(record["AssociatedResourceRegionId"]) == parts[2] &&
				fmt.Sprint(record["AssociatedResourceId"]) == parts[3] {
				return fmt.Errorf("resource attachment %s still exists", rs.Primary.ID)
			}
		}
	}
	return nil
}

func testAccAliCloudGaResourceAttachmentConfig(name string) string {
	return fmt.Sprintf(`
provider "alicloud" {
  region = "cn-hangzhou"
}

variable "name" {
  default = "%s"
}

data "alicloud_wafv3_instances" "default" {}

resource "alicloud_ga_accelerator" "default" {
  accelerator_name       = var.name
  description            = var.name
  bandwidth_billing_type = "CDT"
  payment_type           = "PayAsYouGo"
}

resource "alicloud_ga_resource_attachment" "default" {
  accelerator_id                = alicloud_ga_accelerator.default.id
  associated_resource_type      = "WAF"
  associated_resource_region_id = "cn-hangzhou"
  associated_resource_id        = data.alicloud_wafv3_instances.default.ids.0
}
`, name)
}

func testAccAliCloudGaResourceAttachmentConfigInvalidType(name string) string {
	return fmt.Sprintf(`
provider "alicloud" {
  region = "cn-hangzhou"
}

variable "name" {
  default = "%s"
}

data "alicloud_wafv3_instances" "default" {}

resource "alicloud_ga_accelerator" "default" {
  accelerator_name       = var.name
  description            = var.name
  bandwidth_billing_type = "CDT"
  payment_type           = "PayAsYouGo"
}

resource "alicloud_ga_resource_attachment" "default" {
  accelerator_id                = alicloud_ga_accelerator.default.id
  associated_resource_type      = "ECS"
  associated_resource_region_id = "cn-hangzhou"
  associated_resource_id        = data.alicloud_wafv3_instances.default.ids.0
}
`, name)
}
