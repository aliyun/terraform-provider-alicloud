package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAliCloudGaAccessLog_basic0(t *testing.T) {
	var v map[string]interface{}
	checkoutSupportedRegions(t, true, connectivity.GaSupportRegions)
	resourceId := "alicloud_ga_access_log.default"
	ra := resourceAttrInit(resourceId, resourceAlicloudGaAccessLogMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &GaService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeGaAccessLog")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(100, 999)
	name := fmt.Sprintf("tf-testacc%s-name%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudGaAccessLogBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName:     resourceId,
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"accelerator_id":     "${data.alicloud_ga_accelerators.default.accelerators.0.id}",
					"listener_id":        "${alicloud_ga_listener.default.id}",
					"endpoint_group_id":  "${alicloud_ga_endpoint_group.default.id}",
					"sls_project_name":   "${alicloud_log_project.default.name}",
					"sls_log_store_name": "${alicloud_log_store.default.name}",
					"sls_region_id":      defaultRegionToTest,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"accelerator_id":     CHECKSET,
						"listener_id":        CHECKSET,
						"endpoint_group_id":  CHECKSET,
						"sls_project_name":   name,
						"sls_log_store_name": name,
						"sls_region_id":      defaultRegionToTest,
					}),
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

var resourceAlicloudGaAccessLogMap = map[string]string{
	"status":             CHECKSET,
	"accelerator_id":     CHECKSET,
	"listener_id":        CHECKSET,
	"endpoint_group_id":  CHECKSET,
	"sls_project_name":   CHECKSET,
	"sls_log_store_name": CHECKSET,
	"sls_region_id":      CHECKSET,
}

func AlicloudGaAccessLogBasicDependence0(name string) string {
	return fmt.Sprintf(`
	variable "name" {
  		default = "%s"
	}

	data "alicloud_ga_accelerators" "default" {
  		status = "active"
		bandwidth_billing_type = "BandwidthPackage"
	}

	resource "alicloud_log_project" "default" {
  		name = var.name
	}

	resource "alicloud_log_store" "default" {
  		project = alicloud_log_project.default.name
  		name    = var.name
	}

	resource "alicloud_ga_bandwidth_package" "default" {
		bandwidth      = 100
  		type           = "Basic"
  		bandwidth_type = "Basic"
  		payment_type   = "PayAsYouGo"
  		billing_type   = "PayBy95"
  		ratio          = 30
	}

	resource "alicloud_ga_bandwidth_package_attachment" "default" {
  		accelerator_id       = data.alicloud_ga_accelerators.default.accelerators.0.id
  		bandwidth_package_id = alicloud_ga_bandwidth_package.default.id
	}

	resource "alicloud_ga_listener" "default" {
  		accelerator_id = alicloud_ga_bandwidth_package_attachment.default.accelerator_id
  		port_ranges {
    		from_port = 80
    		to_port   = 80
  		}
	}

	resource "alicloud_eip_address" "default" {
  		payment_type = "PayAsYouGo"
	}

	resource "alicloud_ga_endpoint_group" "default" {
  		accelerator_id = alicloud_ga_listener.default.accelerator_id
		endpoint_configurations {
			endpoint = alicloud_eip_address.default.ip_address
			type     = "PublicIp"
			weight   = 20
		}
  		endpoint_group_region = "%s"
  		listener_id           = alicloud_ga_listener.default.id
}
`, name, defaultRegionToTest)
}

// lintignore: R001
