package alicloud

import (
	"fmt"
	"os"
	"testing"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/smartag"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func TestAccAlicloudSagSnatEntry_basic(t *testing.T) {
	var snat smartag.SnatEntry
	resourceId := "alicloud_sag_snat_entry.default"
	ra := resourceAttrInit(resourceId, nil)
	serviceFunc := func() interface{} {
		return &SagService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &snat, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testSagSnatEntryName-%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceSagSnatEntryDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, connectivity.SmartagSupportedRegions)
			testAccPreCheckWithSmartAccessGatewaySetting(t)
		},

		// module name
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"sag_id":     os.Getenv("SAG_INSTANCE_ID"),
					"cidr_block": "192.168.7.0/24",
					"snat_ip":    "192.0.0.2",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"sag_id":     os.Getenv("SAG_INSTANCE_ID"),
						"cidr_block": "192.168.7.0/24",
						"snat_ip":    "192.0.0.2",
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"sag_id":     os.Getenv("SAG_INSTANCE_ID"),
					"cidr_block": "192.168.10.0/24",
					"snat_ip":    "192.169.0.2",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"sag_id":     os.Getenv("SAG_INSTANCE_ID"),
						"cidr_block": "192.168.10.0/24",
						"snat_ip":    "192.169.0.2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"sag_id":     os.Getenv("SAG_INSTANCE_ID"),
					"cidr_block": "172.16.7.0/24",
					"snat_ip":    "128.0.0.2",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"sag_id":     os.Getenv("SAG_INSTANCE_ID"),
						"cidr_block": "172.16.7.0/24",
						"snat_ip":    "128.0.0.2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"sag_id":     os.Getenv("SAG_INSTANCE_ID"),
					"cidr_block": "172.16.10.0/24",
					"snat_ip":    "172.32.0.2",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"sag_id":     os.Getenv("SAG_INSTANCE_ID"),
						"cidr_block": "172.16.10.0/24",
						"snat_ip":    "172.32.0.2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"sag_id":     os.Getenv("SAG_INSTANCE_ID"),
					"cidr_block": "10.7.0.0/24",
					"snat_ip":    "1.0.0.2",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"sag_id":     os.Getenv("SAG_INSTANCE_ID"),
						"cidr_block": "10.7.0.0/24",
						"snat_ip":    "1.0.0.2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"sag_id":     os.Getenv("SAG_INSTANCE_ID"),
					"cidr_block": "10.10.0.0/24",
					"snat_ip":    "11.0.0.2",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"sag_id":     os.Getenv("SAG_INSTANCE_ID"),
						"cidr_block": "10.10.0.0/24",
						"snat_ip":    "11.0.0.2",
					}),
				),
			},
		},
	})
}

func TestAccAlicloudSagSnatEntry_multi(t *testing.T) {
	var snat smartag.SnatEntry
	resourceId := "alicloud_sag_snat_entry.default.4"
	ra := resourceAttrInit(resourceId, nil)
	serviceFunc := func() interface{} {
		return &SagService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &snat, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testSagSnatEntryName-%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceSagSnatEntryDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, connectivity.SmartagSupportedRegions)
			testAccPreCheckWithSmartAccessGatewaySetting(t)
		},

		// module name
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"sag_id":     os.Getenv("SAG_INSTANCE_ID"),
					"count":      "5",
					"cidr_block": fmt.Sprintf("192.168.%s.0/24", "${count.index}"),
					"snat_ip":    fmt.Sprintf("192.0.0.%s", "${count.index}"),
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"sag_id":     os.Getenv("SAG_INSTANCE_ID"),
						"cidr_block": "192.168.4.0/24",
						"snat_ip":    "192.0.0.4",
					}),
				),
			},
		},
	})
}

func resourceSagSnatEntryDependence(name string) string {
	return ""
}
