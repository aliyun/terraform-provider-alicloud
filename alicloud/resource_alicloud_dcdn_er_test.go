package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAliCloudDcdnEr_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_dcdn_er.default"
	ra := resourceAttrInit(resourceId, AlicloudDcdnErMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &DcdnService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeDcdnEr")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sdcdn-er%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudDcdnErBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"er_name":     name,
					"description": name,
					"env_conf": []map[string]interface{}{
						{
							"staging": []map[string]interface{}{
								{
									"spec_name":     "5ms",
									"allowed_hosts": []string{"example.com"},
								},
							},
							"production": []map[string]interface{}{
								{
									"spec_name":     "5ms",
									"allowed_hosts": []string{"example.com"},
								},
							},
							"preset_canary_anhui": []map[string]interface{}{
								{
									"spec_name":     "5ms",
									"allowed_hosts": []string{"example.com"},
								},
							},
							"preset_canary_beijing": []map[string]interface{}{
								{
									"spec_name":     "5ms",
									"allowed_hosts": []string{"example.com"},
								},
							},
							"preset_canary_chongqing": []map[string]interface{}{
								{
									"spec_name":     "5ms",
									"allowed_hosts": []string{"example.com"},
								},
							},
							"preset_canary_fujian": []map[string]interface{}{
								{
									"spec_name":     "5ms",
									"allowed_hosts": []string{"example.com"},
								},
							},
							"preset_canary_gansu": []map[string]interface{}{
								{
									"spec_name":     "5ms",
									"allowed_hosts": []string{"example.com"},
								},
							},
							"preset_canary_guangdong": []map[string]interface{}{
								{
									"spec_name":     "5ms",
									"allowed_hosts": []string{"example.com"},
								},
							},
							"preset_canary_guangxi": []map[string]interface{}{
								{
									"spec_name":     "5ms",
									"allowed_hosts": []string{"example.com"},
								},
							},
							"preset_canary_guizhou": []map[string]interface{}{
								{
									"spec_name":     "5ms",
									"allowed_hosts": []string{"example.com"},
								},
							},
							"preset_canary_hainan": []map[string]interface{}{
								{
									"spec_name":     "5ms",
									"allowed_hosts": []string{"example.com"},
								},
							},
							"preset_canary_hebei": []map[string]interface{}{
								{
									"spec_name":     "5ms",
									"allowed_hosts": []string{"example.com"},
								},
							},
							"preset_canary_heilongjiang": []map[string]interface{}{
								{
									"spec_name":     "5ms",
									"allowed_hosts": []string{"example.com"},
								},
							},
							"preset_canary_henan": []map[string]interface{}{
								{
									"spec_name":     "5ms",
									"allowed_hosts": []string{"example.com"},
								},
							},
							"preset_canary_hong_kong": []map[string]interface{}{
								{
									"spec_name":     "5ms",
									"allowed_hosts": []string{"example.com"},
								},
							},
							"preset_canary_hubei": []map[string]interface{}{
								{
									"spec_name":     "5ms",
									"allowed_hosts": []string{"example.com"},
								},
							},
							"preset_canary_hunan": []map[string]interface{}{
								{
									"spec_name":     "5ms",
									"allowed_hosts": []string{"example.com"},
								},
							},
							"preset_canary_jiangsu": []map[string]interface{}{
								{
									"spec_name":     "5ms",
									"allowed_hosts": []string{"example.com"},
								},
							},
							"preset_canary_jiangxi": []map[string]interface{}{
								{
									"spec_name":     "5ms",
									"allowed_hosts": []string{"example.com"},
								},
							},
							"preset_canary_jilin": []map[string]interface{}{
								{
									"spec_name":     "5ms",
									"allowed_hosts": []string{"example.com"},
								},
							},
							"preset_canary_liaoning": []map[string]interface{}{
								{
									"spec_name":     "5ms",
									"allowed_hosts": []string{"example.com"},
								},
							},
							"preset_canary_macau": []map[string]interface{}{
								{
									"spec_name":     "5ms",
									"allowed_hosts": []string{"example.com"},
								},
							},
							"preset_canary_neimenggu": []map[string]interface{}{
								{
									"spec_name":     "5ms",
									"allowed_hosts": []string{"example.com"},
								},
							},
							"preset_canary_ningxia": []map[string]interface{}{
								{
									"spec_name":     "5ms",
									"allowed_hosts": []string{"example.com"},
								},
							},
							"preset_canary_qinghai": []map[string]interface{}{
								{
									"spec_name":     "5ms",
									"allowed_hosts": []string{"example.com"},
								},
							},
							"preset_canary_shaanxi": []map[string]interface{}{
								{
									"spec_name":     "5ms",
									"allowed_hosts": []string{"example.com"},
								},
							},
							"preset_canary_shandong": []map[string]interface{}{
								{
									"spec_name":     "5ms",
									"allowed_hosts": []string{"example.com"},
								},
							},
							"preset_canary_shanghai": []map[string]interface{}{
								{
									"spec_name":     "5ms",
									"allowed_hosts": []string{"example.com"},
								},
							},
							"preset_canary_shanxi": []map[string]interface{}{
								{
									"spec_name":     "5ms",
									"allowed_hosts": []string{"example.com"},
								},
							},
							"preset_canary_sichuan": []map[string]interface{}{
								{
									"spec_name":     "5ms",
									"allowed_hosts": []string{"example.com"},
								},
							},
							"preset_canary_taiwan": []map[string]interface{}{
								{
									"spec_name":     "5ms",
									"allowed_hosts": []string{"example.com"},
								},
							},
							"preset_canary_tianjin": []map[string]interface{}{
								{
									"spec_name":     "5ms",
									"allowed_hosts": []string{"example.com"},
								},
							},
							"preset_canary_xinjiang": []map[string]interface{}{
								{
									"spec_name":     "5ms",
									"allowed_hosts": []string{"example.com"},
								},
							},
							"preset_canary_xizang": []map[string]interface{}{
								{
									"spec_name":     "5ms",
									"allowed_hosts": []string{"example.com"},
								},
							},
							"preset_canary_yunnan": []map[string]interface{}{
								{
									"spec_name":     "5ms",
									"allowed_hosts": []string{"example.com"},
								},
							},
							"preset_canary_zhejiang": []map[string]interface{}{
								{
									"spec_name":     "5ms",
									"allowed_hosts": []string{"example.com"},
								},
							},
							"preset_canary_overseas": []map[string]interface{}{
								{
									"spec_name":     "5ms",
									"allowed_hosts": []string{"example.com"},
								},
							},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"er_name":     name,
						"description": name,
						"env_conf.#":  "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": name + "-update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": name + "-update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"env_conf": []map[string]interface{}{
						{
							"staging": []map[string]interface{}{
								{
									"spec_name":     "50ms",
									"allowed_hosts": []string{"example.com"},
								},
							},
							"production": []map[string]interface{}{
								{
									"spec_name":     "50ms",
									"allowed_hosts": []string{"example.com"},
								},
							},
							"preset_canary_anhui": []map[string]interface{}{
								{
									"spec_name":     "50ms",
									"allowed_hosts": []string{"example.com"},
								},
							},
							"preset_canary_beijing": []map[string]interface{}{
								{
									"spec_name":     "50ms",
									"allowed_hosts": []string{"example.com"},
								},
							},
							"preset_canary_chongqing": []map[string]interface{}{
								{
									"spec_name":     "50ms",
									"allowed_hosts": []string{"example.com"},
								},
							},
							"preset_canary_fujian": []map[string]interface{}{
								{
									"spec_name":     "50ms",
									"allowed_hosts": []string{"example.com"},
								},
							},
							"preset_canary_gansu": []map[string]interface{}{
								{
									"spec_name":     "50ms",
									"allowed_hosts": []string{"example.com"},
								},
							},
							"preset_canary_guangdong": []map[string]interface{}{
								{
									"spec_name":     "50ms",
									"allowed_hosts": []string{"example.com"},
								},
							},
							"preset_canary_guangxi": []map[string]interface{}{
								{
									"spec_name":     "50ms",
									"allowed_hosts": []string{"example.com"},
								},
							},
							"preset_canary_guizhou": []map[string]interface{}{
								{
									"spec_name":     "50ms",
									"allowed_hosts": []string{"example.com"},
								},
							},
							"preset_canary_hainan": []map[string]interface{}{
								{
									"spec_name":     "50ms",
									"allowed_hosts": []string{"example.com"},
								},
							},
							"preset_canary_hebei": []map[string]interface{}{
								{
									"spec_name":     "50ms",
									"allowed_hosts": []string{"example.com"},
								},
							},
							"preset_canary_heilongjiang": []map[string]interface{}{
								{
									"spec_name":     "50ms",
									"allowed_hosts": []string{"example.com"},
								},
							},
							"preset_canary_henan": []map[string]interface{}{
								{
									"spec_name":     "50ms",
									"allowed_hosts": []string{"example.com"},
								},
							},
							"preset_canary_hong_kong": []map[string]interface{}{
								{
									"spec_name":     "50ms",
									"allowed_hosts": []string{"example.com"},
								},
							},
							"preset_canary_hubei": []map[string]interface{}{
								{
									"spec_name":     "50ms",
									"allowed_hosts": []string{"example.com"},
								},
							},
							"preset_canary_hunan": []map[string]interface{}{
								{
									"spec_name":     "50ms",
									"allowed_hosts": []string{"example.com"},
								},
							},
							"preset_canary_jiangsu": []map[string]interface{}{
								{
									"spec_name":     "50ms",
									"allowed_hosts": []string{"example.com"},
								},
							},
							"preset_canary_jiangxi": []map[string]interface{}{
								{
									"spec_name":     "50ms",
									"allowed_hosts": []string{"example.com"},
								},
							},
							"preset_canary_jilin": []map[string]interface{}{
								{
									"spec_name":     "50ms",
									"allowed_hosts": []string{"example.com"},
								},
							},
							"preset_canary_liaoning": []map[string]interface{}{
								{
									"spec_name":     "50ms",
									"allowed_hosts": []string{"example.com"},
								},
							},
							"preset_canary_macau": []map[string]interface{}{
								{
									"spec_name":     "50ms",
									"allowed_hosts": []string{"example.com"},
								},
							},
							"preset_canary_neimenggu": []map[string]interface{}{
								{
									"spec_name":     "50ms",
									"allowed_hosts": []string{"example.com"},
								},
							},
							"preset_canary_ningxia": []map[string]interface{}{
								{
									"spec_name":     "50ms",
									"allowed_hosts": []string{"example.com"},
								},
							},
							"preset_canary_qinghai": []map[string]interface{}{
								{
									"spec_name":     "50ms",
									"allowed_hosts": []string{"example.com"},
								},
							},
							"preset_canary_shaanxi": []map[string]interface{}{
								{
									"spec_name":     "50ms",
									"allowed_hosts": []string{"example.com"},
								},
							},
							"preset_canary_shandong": []map[string]interface{}{
								{
									"spec_name":     "50ms",
									"allowed_hosts": []string{"example.com"},
								},
							},
							"preset_canary_shanghai": []map[string]interface{}{
								{
									"spec_name":     "50ms",
									"allowed_hosts": []string{"example.com"},
								},
							},
							"preset_canary_shanxi": []map[string]interface{}{
								{
									"spec_name":     "50ms",
									"allowed_hosts": []string{"example.com"},
								},
							},
							"preset_canary_sichuan": []map[string]interface{}{
								{
									"spec_name":     "50ms",
									"allowed_hosts": []string{"example.com"},
								},
							},
							"preset_canary_taiwan": []map[string]interface{}{
								{
									"spec_name":     "50ms",
									"allowed_hosts": []string{"example.com"},
								},
							},
							"preset_canary_tianjin": []map[string]interface{}{
								{
									"spec_name":     "50ms",
									"allowed_hosts": []string{"example.com"},
								},
							},
							"preset_canary_xinjiang": []map[string]interface{}{
								{
									"spec_name":     "50ms",
									"allowed_hosts": []string{"example.com"},
								},
							},
							"preset_canary_xizang": []map[string]interface{}{
								{
									"spec_name":     "50ms",
									"allowed_hosts": []string{"example.com"},
								},
							},
							"preset_canary_yunnan": []map[string]interface{}{
								{
									"spec_name":     "50ms",
									"allowed_hosts": []string{"example.com"},
								},
							},
							"preset_canary_zhejiang": []map[string]interface{}{
								{
									"spec_name":     "50ms",
									"allowed_hosts": []string{"example.com"},
								},
							},
							"preset_canary_overseas": []map[string]interface{}{
								{
									"spec_name":     "50ms",
									"allowed_hosts": []string{"example.com"},
								},
							},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"env_conf.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"env_conf": []map[string]interface{}{
						{
							"staging": []map[string]interface{}{
								{
									"spec_name":     "100ms",
									"allowed_hosts": []string{"example.com"},
								},
							},
							"production": []map[string]interface{}{
								{
									"spec_name":     "100ms",
									"allowed_hosts": []string{"example.com"},
								},
							},
							"preset_canary_anhui": []map[string]interface{}{
								{
									"spec_name":     "100ms",
									"allowed_hosts": []string{"example.com"},
								},
							},
							"preset_canary_beijing": []map[string]interface{}{
								{
									"spec_name":     "100ms",
									"allowed_hosts": []string{"example.com"},
								},
							},
							"preset_canary_chongqing": []map[string]interface{}{
								{
									"spec_name":     "100ms",
									"allowed_hosts": []string{"example.com"},
								},
							},
							"preset_canary_fujian": []map[string]interface{}{
								{
									"spec_name":     "100ms",
									"allowed_hosts": []string{"example.com"},
								},
							},
							"preset_canary_gansu": []map[string]interface{}{
								{
									"spec_name":     "100ms",
									"allowed_hosts": []string{"example.com"},
								},
							},
							"preset_canary_guangdong": []map[string]interface{}{
								{
									"spec_name":     "100ms",
									"allowed_hosts": []string{"example.com"},
								},
							},
							"preset_canary_guangxi": []map[string]interface{}{
								{
									"spec_name":     "100ms",
									"allowed_hosts": []string{"example.com"},
								},
							},
							"preset_canary_guizhou": []map[string]interface{}{
								{
									"spec_name":     "100ms",
									"allowed_hosts": []string{"example.com"},
								},
							},
							"preset_canary_hainan": []map[string]interface{}{
								{
									"spec_name":     "100ms",
									"allowed_hosts": []string{"example.com"},
								},
							},
							"preset_canary_hebei": []map[string]interface{}{
								{
									"spec_name":     "100ms",
									"allowed_hosts": []string{"example.com"},
								},
							},
							"preset_canary_heilongjiang": []map[string]interface{}{
								{
									"spec_name":     "100ms",
									"allowed_hosts": []string{"example.com"},
								},
							},
							"preset_canary_henan": []map[string]interface{}{
								{
									"spec_name":     "100ms",
									"allowed_hosts": []string{"example.com"},
								},
							},
							"preset_canary_hong_kong": []map[string]interface{}{
								{
									"spec_name":     "100ms",
									"allowed_hosts": []string{"example.com"},
								},
							},
							"preset_canary_hubei": []map[string]interface{}{
								{
									"spec_name":     "100ms",
									"allowed_hosts": []string{"example.com"},
								},
							},
							"preset_canary_hunan": []map[string]interface{}{
								{
									"spec_name":     "100ms",
									"allowed_hosts": []string{"example.com"},
								},
							},
							"preset_canary_jiangsu": []map[string]interface{}{
								{
									"spec_name":     "100ms",
									"allowed_hosts": []string{"example.com"},
								},
							},
							"preset_canary_jiangxi": []map[string]interface{}{
								{
									"spec_name":     "100ms",
									"allowed_hosts": []string{"example.com"},
								},
							},
							"preset_canary_jilin": []map[string]interface{}{
								{
									"spec_name":     "100ms",
									"allowed_hosts": []string{"example.com"},
								},
							},
							"preset_canary_liaoning": []map[string]interface{}{
								{
									"spec_name":     "100ms",
									"allowed_hosts": []string{"example.com"},
								},
							},
							"preset_canary_macau": []map[string]interface{}{
								{
									"spec_name":     "100ms",
									"allowed_hosts": []string{"example.com"},
								},
							},
							"preset_canary_neimenggu": []map[string]interface{}{
								{
									"spec_name":     "100ms",
									"allowed_hosts": []string{"example.com"},
								},
							},
							"preset_canary_ningxia": []map[string]interface{}{
								{
									"spec_name":     "100ms",
									"allowed_hosts": []string{"example.com"},
								},
							},
							"preset_canary_qinghai": []map[string]interface{}{
								{
									"spec_name":     "100ms",
									"allowed_hosts": []string{"example.com"},
								},
							},
							"preset_canary_shaanxi": []map[string]interface{}{
								{
									"spec_name":     "100ms",
									"allowed_hosts": []string{"example.com"},
								},
							},
							"preset_canary_shandong": []map[string]interface{}{
								{
									"spec_name":     "100ms",
									"allowed_hosts": []string{"example.com"},
								},
							},
							"preset_canary_shanghai": []map[string]interface{}{
								{
									"spec_name":     "100ms",
									"allowed_hosts": []string{"example.com"},
								},
							},
							"preset_canary_shanxi": []map[string]interface{}{
								{
									"spec_name":     "100ms",
									"allowed_hosts": []string{"example.com"},
								},
							},
							"preset_canary_sichuan": []map[string]interface{}{
								{
									"spec_name":     "100ms",
									"allowed_hosts": []string{"example.com"},
								},
							},
							"preset_canary_taiwan": []map[string]interface{}{
								{
									"spec_name":     "100ms",
									"allowed_hosts": []string{"example.com"},
								},
							},
							"preset_canary_tianjin": []map[string]interface{}{
								{
									"spec_name":     "100ms",
									"allowed_hosts": []string{"example.com"},
								},
							},
							"preset_canary_xinjiang": []map[string]interface{}{
								{
									"spec_name":     "100ms",
									"allowed_hosts": []string{"example.com"},
								},
							},
							"preset_canary_xizang": []map[string]interface{}{
								{
									"spec_name":     "100ms",
									"allowed_hosts": []string{"example.com"},
								},
							},
							"preset_canary_yunnan": []map[string]interface{}{
								{
									"spec_name":     "100ms",
									"allowed_hosts": []string{"example.com"},
								},
							},
							"preset_canary_zhejiang": []map[string]interface{}{
								{
									"spec_name":     "100ms",
									"allowed_hosts": []string{"example.com"},
								},
							},
							"preset_canary_overseas": []map[string]interface{}{
								{
									"spec_name":     "100ms",
									"allowed_hosts": []string{"example.com"},
								},
							},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"env_conf.#": "1",
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

var AlicloudDcdnErMap = map[string]string{}

func AlicloudDcdnErBasicDependence0(name string) string {
	return fmt.Sprintf(`
	variable "name" {
  		default = "%s"
	}
`, name)
}
