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

func TestAccAlicloudSagDnatEntry_basic(t *testing.T) {
	var dnat smartag.DnatEntry
	resourceId := "alicloud_sag_dnat_entry.default"
	ra := resourceAttrInit(resourceId, nil)
	serviceFunc := func() interface{} {
		return &SagService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &dnat, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testSagDnatEntryName-%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceSagDnatEntryBasicDependence)
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
					"sag_id":        os.Getenv("SAG_INSTANCE_ID"),
					"type":          "Internet",
					"ip_protocol":   "tcp",
					"external_port": "65535",
					"internal_ip":   "192.168.0.2",
					"internal_port": "20",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"sag_id":        os.Getenv("SAG_INSTANCE_ID"),
						"type":          "Internet",
						"ip_protocol":   "tcp",
						"external_port": "65535",
						"internal_ip":   "192.168.0.2",
						"internal_port": "20",
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
					"sag_id":        os.Getenv("SAG_INSTANCE_ID"),
					"type":          "Internet",
					"ip_protocol":   "udp",
					"external_port": "65535",
					"internal_ip":   "192.168.0.4",
					"internal_port": "20",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"sag_id":        os.Getenv("SAG_INSTANCE_ID"),
						"type":          "Internet",
						"ip_protocol":   "udp",
						"external_port": "65535",
						"internal_ip":   "192.168.0.4",
						"internal_port": "20",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"sag_id":        os.Getenv("SAG_INSTANCE_ID"),
					"type":          "Internet",
					"ip_protocol":   "any",
					"external_port": "any",
					"internal_ip":   "172.16.0.4",
					"internal_port": "any",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"sag_id":        os.Getenv("SAG_INSTANCE_ID"),
						"type":          "Internet",
						"ip_protocol":   "any",
						"external_port": "any",
						"internal_ip":   "172.16.0.4",
						"internal_port": "any",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"sag_id":        os.Getenv("SAG_INSTANCE_ID"),
					"type":          "Intranet",
					"ip_protocol":   "tcp",
					"external_ip":   "1.0.0.2",
					"external_port": "1",
					"internal_ip":   "10.0.0.2",
					"internal_port": "20",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"sag_id":        os.Getenv("SAG_INSTANCE_ID"),
						"type":          "Intranet",
						"ip_protocol":   "tcp",
						"external_ip":   "1.0.0.2",
						"external_port": "1",
						"internal_ip":   "10.0.0.2",
						"internal_port": "20",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"sag_id":        os.Getenv("SAG_INSTANCE_ID"),
					"type":          "Intranet",
					"ip_protocol":   "udp",
					"external_ip":   "11.0.0.2",
					"external_port": "1",
					"internal_ip":   "10.0.0.4",
					"internal_port": "20",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"sag_id":        os.Getenv("SAG_INSTANCE_ID"),
						"type":          "Intranet",
						"ip_protocol":   "udp",
						"external_ip":   "11.0.0.2",
						"external_port": "1",
						"internal_ip":   "10.0.0.4",
						"internal_port": "20",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"sag_id":        os.Getenv("SAG_INSTANCE_ID"),
					"type":          "Intranet",
					"ip_protocol":   "any",
					"external_ip":   "172.32.0.2",
					"external_port": "any",
					"internal_ip":   "172.16.0.4",
					"internal_port": "any",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"sag_id":        os.Getenv("SAG_INSTANCE_ID"),
						"type":          "Intranet",
						"ip_protocol":   "any",
						"external_ip":   "172.32.0.2",
						"external_port": "any",
						"internal_ip":   "172.16.0.4",
						"internal_port": "any",
					}),
				),
			},
		},
	})
}

func TestAccAlicloudSagDnatEntry_multi(t *testing.T) {
	var dnat smartag.DnatEntry
	resourceId := "alicloud_sag_dnat_entry.default.4"
	ra := resourceAttrInit(resourceId, nil)
	serviceFunc := func() interface{} {
		return &SagService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &dnat, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testSagDnatEntryName-%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceSagDnatEntryBasicDependence)
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
					"sag_id":        os.Getenv("SAG_INSTANCE_ID"),
					"count":         "5",
					"type":          "Intranet",
					"ip_protocol":   "tcp",
					"external_ip":   fmt.Sprintf("1.0.0.%s", "${count.index}"),
					"external_port": "1",
					"internal_ip":   fmt.Sprintf("10.0.0.%s", "${count.index}"),
					"internal_port": "20",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"sag_id":        os.Getenv("SAG_INSTANCE_ID"),
						"type":          "Intranet",
						"ip_protocol":   "tcp",
						"external_ip":   "1.0.0.4",
						"external_port": "1",
						"internal_ip":   "10.0.0.4",
						"internal_port": "20",
					}),
				),
			},
		},
	})
}

func resourceSagDnatEntryBasicDependence(name string) string {
	return ""
}
