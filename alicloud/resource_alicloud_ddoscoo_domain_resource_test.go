package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAlicloudDdoscooDomainResource_basic(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ddoscoo_domain_resource.default"
	ra := resourceAttrInit(resourceId, AlicloudDdoscooDomainResourceMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &DdoscooService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeDdoscooDomainResource")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	testAccConfig := resourceTestAccConfigFunc(resourceId, "", AlicloudDdoscooDomainResourceBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, connectivity.DdoscooSupportedRegions)
		},

		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"domain": "liduotesttf.qq.com",
					"proxy_types": []map[string]interface{}{
						{
							"proxy_ports": []string{"443"},
							"proxy_type":  "https",
						},
					},
					"instance_ids": []string{"${data.alicloud_ddoscoo_instances.default.ids.0}"},
					"real_servers": []string{"177.167.32.11"},
					"rs_type":      `0`,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_ids.#": "1",
						"proxy_types.#":  "1",
						"real_servers.#": "1",
						"rs_type":        "0",
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
					"https_ext": `{\"Http2\":1,\"Http2https\":0,\"Https2http\":0}`,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"https_ext": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_ids": []string{"${data.alicloud_ddoscoo_instances.default.ids.0}", "${data.alicloud_ddoscoo_instances.default.ids.1}"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_ids.#": "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"real_servers": []string{"177.167.32.11", "177.167.32.12"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"real_servers.#": "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_ids": []string{"${data.alicloud_ddoscoo_instances.default.ids.0}"},
					"real_servers": []string{"177.167.32.11"},
					"https_ext":    `{\"Http2\":0,\"Http2https\":0,\"Https2http\":0}`,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"real_servers.#": "1",
						"instance_ids.#": "1",
						"https_ext":      CHECKSET,
					}),
				),
			},
		},
	})
}

var AlicloudDdoscooDomainResourceMap0 = map[string]string{
	"domain":         CHECKSET,
	"https_ext":      CHECKSET,
	"instance_ids.#": "1",
	"proxy_types.#":  "1",
	"real_servers.#": "1",
	"rs_type":        "0",
}

func AlicloudDdoscooDomainResourceBasicDependence0(name string) string {
	return fmt.Sprintf(`
data "alicloud_ddoscoo_instances" "default" {}
`)
}
