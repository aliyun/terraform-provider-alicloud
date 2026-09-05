package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAliCloudThreatDetectionDataSet_basic(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_threat_detection_data_set.default"
	ra := resourceAttrInit(resourceId, AlicloudThreatDetectionDataSetMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ThreatDetectionServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeThreatDetectionDataSet")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%s%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudThreatDetectionDataSetBasicDependence)
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
					"data_set_name":           "${var.name}",
					"data_set_field_key_name": "ip",
					"data_set_file_name":      "example.csv",
					"data_set_description":    "description for test",
					"data_set_type":           "userDefined",
					"data_set_status":         1,
					"role_for":                0,
					"lang":                    "zh",
					"region_id":               "cn-hangzhou",
					"ip_whitelist_recognizers": []map[string]interface{}{
						{
							"auto_recognize_status":        "on",
							"recognize_scope":              "all",
							"ip_whitelist_recognizer_type": "ip",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"data_set_name":           "CHECKSET",
						"data_set_field_key_name": "CHECKSET",
						"data_set_file_name":      "CHECKSET",
						"data_set_description":    "CHECKSET",
						"data_set_type":           "CHECKSET",
						"data_set_status":         "CHECKSET",
						"role_for":                "CHECKSET",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"data_set_name":           "${var.name}-u",
					"data_set_field_key_name": "ip",
					"data_set_file_name":      "test_file_u.csv",
					"data_set_description":    "description updated",
					"data_set_type":           "userDefined",
					"data_set_status":         2,
					"role_for":                0,
					"ip_whitelist_recognizers": []map[string]interface{}{
						{
							"auto_recognize_status":        "off",
							"recognize_scope":              "custom",
							"ip_whitelist_recognizer_type": "domain",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"data_set_name":        "CHECKSET",
						"data_set_description": "CHECKSET",
						"data_set_status":      "CHECKSET",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"data_set_name":           "${var.name}",
					"data_set_field_key_name": "ip",
					"data_set_file_name":      "test_file.csv",
					"data_set_description":    "description for test",
					"data_set_type":           "userDefined",
					"data_set_status":         1,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"data_set_name":        "CHECKSET",
						"data_set_description": "CHECKSET",
						"data_set_status":      "CHECKSET",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"lang", "region_id", "role_for"},
			},
		},
	})
}

var AlicloudThreatDetectionDataSetMap = map[string]string{}

func AlicloudThreatDetectionDataSetBasicDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}
`, name)
}
