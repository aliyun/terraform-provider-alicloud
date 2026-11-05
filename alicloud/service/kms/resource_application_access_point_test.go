package kms_test

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/kms"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// collection-order: bypass policies -- pre-existing stock TypeList; schema and member-order coverage unchanged by this pure directory migration (route provider_handwritten)
// Test Kms ApplicationAccessPoint. >>> Resource test cases, automatically generated.
// Case 4108
func TestAccAliCloudKmsApplicationAccessPoint_basic4108(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_kms_application_access_point.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudKmsApplicationAccessPointMap4108)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return kms.NewKmsServiceV2(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeKmsApplicationAccessPoint")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%skmsapplicationaccesspoint%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudKmsApplicationAccessPointBasicDependence4108)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
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
	ra := acctest.ResourceAttrInit(resourceId, AlicloudKmsApplicationAccessPointMap4108)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return kms.NewKmsServiceV2(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeKmsApplicationAccessPoint")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%skmsapplicationaccesspoint%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudKmsApplicationAccessPointBasicDependence4108)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
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
