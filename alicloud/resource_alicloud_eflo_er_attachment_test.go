// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// Test Eflo ErAttachment. >>> Resource test cases, automatically generated.
// Case er_attachment_basic 12485
func TestAccAliCloudEfloErAttachment_basic12485(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_eflo_er_attachment.default"
	ra := resourceAttrInit(resourceId, AlicloudEfloErAttachmentMap12485)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EfloServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeEfloErAttachment")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfacceflo%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudEfloErAttachmentBasicDependence12485)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-wulanchabu"})
			testAccPreCheck(t)
		},
		IDRefreshName:     resourceId,
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"er_id":                  "${alicloud_eflo_er.default.id}",
					"er_attachment_name":     "${var.name}",
					"instance_id":            "${alicloud_eflo_vpd.default.id}",
					"instance_type":          "VPD",
					"auto_receive_all_route": false,
					"resource_tenant_id":     "${data.alicloud_account.default.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"er_id":                  CHECKSET,
						"er_attachment_name":     name,
						"instance_id":            CHECKSET,
						"instance_type":          "VPD",
						"auto_receive_all_route": "false",
						"resource_tenant_id":     CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"er_attachment_name": "${var.name}-update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"er_attachment_name": name + "-update",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{},
			},
		},
	})
}

var AlicloudEfloErAttachmentMap12485 = map[string]string{
	"status":           CHECKSET,
	"create_time":      CHECKSET,
	"er_attachment_id": CHECKSET,
	"region_id":        CHECKSET,
}

func AlicloudEfloErAttachmentBasicDependence12485(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

variable "region_id" {
  default = "cn-wulanchabu"
}

variable "zone_id" {
  default = "cn-wulanchabu-b"
}

data "alicloud_account" "default" {}

resource "alicloud_eflo_vpd" "default" {
  cidr     = "10.0.0.0/8"
  vpd_name = "tf-attachment-vpd-${var.name}"
}

resource "alicloud_eflo_er" "default" {
  er_name        = "tf-attachment-er-${var.name}"
  master_zone_id = var.zone_id
}

`, name)
}

// Test Eflo ErAttachment. <<< Resource test cases, automatically generated.
