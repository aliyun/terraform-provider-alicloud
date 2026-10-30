package smartag_test

import (
	"fmt"
	"os"
	"testing"

	smartagsdk "github.com/aliyun/alibaba-cloud-sdk-go/services/smartag"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/smartag"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccAliCloudSagDnatEntry_basic(t *testing.T) {
	var dnat smartagsdk.DnatEntry
	resourceId := "alicloud_sag_dnat_entry.default"
	ra := acctest.ResourceAttrInit(resourceId, nil)
	serviceFunc := func() interface{} {
		return smartag.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}
	rc := acctest.ResourceCheckInit(resourceId, &dnat, serviceFunc)
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testSagDnatEntryName-%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, resourceSagDnatEntryBasicDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckWithRegions(t, true, connectivity.SmartagSupportedRegions)
			acctest.PreCheckWithEnvVariable(t, "SAG_INSTANCE_ID")
		},

		// module name
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
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

func TestAccAliCloudSagDnatEntry_multi(t *testing.T) {
	var dnat smartagsdk.DnatEntry
	resourceId := "alicloud_sag_dnat_entry.default.4"
	ra := acctest.ResourceAttrInit(resourceId, nil)
	serviceFunc := func() interface{} {
		return smartag.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}
	rc := acctest.ResourceCheckInit(resourceId, &dnat, serviceFunc)
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testSagDnatEntryName-%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, resourceSagDnatEntryBasicDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckWithRegions(t, true, connectivity.SmartagSupportedRegions)
			acctest.PreCheckWithEnvVariable(t, "SAG_INSTANCE_ID")
		},

		// module name
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
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
