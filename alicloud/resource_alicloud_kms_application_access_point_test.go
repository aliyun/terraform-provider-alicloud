package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// Test Kms ApplicationAccessPoint. >>> Resource test cases, automatically generated.
// Case 4108
func TestAccAliCloudKmsApplicationAccessPoint_basic4108(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_kms_application_access_point.default"
	ra := resourceAttrInit(resourceId, AlicloudKmsApplicationAccessPointMap4108)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &KmsServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeKmsApplicationAccessPoint")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%skmsapplicationaccesspoint%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudKmsApplicationAccessPointBasicDependence4108)
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
					"application_access_point_name": name,
					"policies": []string{
						"abc", "efg", "hfc"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"application_access_point_name": name,
						"policies.#":                    "3",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": "test aap",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": "test aap",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"policies": []string{
						"abc", "efg", "hfc"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"policies.#": "3",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": "asfdsfads",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": "asfdsfads",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"policies": []string{
						"guaguagua"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"policies.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description":                   "test aap",
					"application_access_point_name": name + "_update",
					"policies": []string{
						"abc", "efg", "hfc"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description":                   "test aap",
						"application_access_point_name": name + "_update",
						"policies.#":                    "3",
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

var AlicloudKmsApplicationAccessPointMap4108 = map[string]string{}

func AlicloudKmsApplicationAccessPointBasicDependence4108(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}


`, name)
}

// Case 4108  twin
func TestAccAliCloudKmsApplicationAccessPoint_basic4108_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_kms_application_access_point.default"
	ra := resourceAttrInit(resourceId, AlicloudKmsApplicationAccessPointMap4108)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &KmsServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeKmsApplicationAccessPoint")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%skmsapplicationaccesspoint%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudKmsApplicationAccessPointBasicDependence4108)
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
					"description":                   "asfdsfads",
					"application_access_point_name": name,
					"policies": []string{
						"guaguagua", "efg", "hfc"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description":                   "asfdsfads",
						"application_access_point_name": name,
						"policies.#":                    "3",
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

// Test Kms ApplicationAccessPoint. <<< Resource test cases, automatically generated.
