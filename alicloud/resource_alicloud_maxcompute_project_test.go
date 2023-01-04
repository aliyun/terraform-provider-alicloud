package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// Case 1
func TestAccAlicloudMaxcomputeProject_basic1968(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_maxcompute_project.default"
	ra := resourceAttrInit(resourceId, AlicloudMaxcomputeProjectMap1968)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &MaxComputeService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeMaxcomputeProject")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	checkoutSupportedRegions(t, true, connectivity.MaxComputeProjectSupportRegions)
	name := fmt.Sprintf("tf_testaccmp%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudMaxcomputeProjectBasicDependence1968)
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
					"default_quota": "默认后付费Quota",
					"project_name":  "${var.name}",
					"comment":       "${var.name}",
					"product_type":  "PayAsYouGo",
					"ip_white_list": []map[string]interface{}{
						{
							"ip_list":     "1.1.1.1,2.2.2.2",
							"vpc_ip_list": "10.10.10.10,11.11.11.11",
						},
					},
					"properties": []map[string]interface{}{
						{
							"allow_full_scan":  "false",
							"enable_decimal2":  "true",
							"retention_days":   "1",
							"sql_metering_max": "0",
							"timezone":         "Asia/Shanghai",
							"type_system":      "2",
							"encryption": []map[string]interface{}{
								{
									"enable":    "true",
									"algorithm": "AESCTR",
									"key":       "f58d854d-7bc0-4a6e-9205-160e10ffedec",
								},
							},
							"table_lifecycle": []map[string]interface{}{
								{
									"type":  "optional",
									"value": "37231",
								},
							},
						},
					},
					"security_properties": []map[string]interface{}{
						{
							"enable_download_privilege":            "false",
							"label_security":                       "true",
							"object_creator_has_access_permission": "true",
							"object_creator_has_grant_permission":  "true",
							"project_protection": []map[string]interface{}{
								{
									"protected": "false",
								},
							},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"default_quota":         CHECKSET,
						"project_name":          CHECKSET,
						"comment":               CHECKSET,
						"product_type":          CHECKSET,
						"ip_white_list.#":       "1",
						"properties.#":          "1",
						"security_properties.#": "1",
					}),
				),
			}, {
				Config: testAccConfig(map[string]interface{}{
					"comment":      "${var.name}_u",
					"project_name": "${var.name}_u",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"comment":      CHECKSET,
						"project_name": CHECKSET,
					}),
				),
			},
		},
	})
}

var AlicloudMaxcomputeProjectMap1968 = map[string]string{}

func AlicloudMaxcomputeProjectBasicDependence1968(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}
`, name)
}
