package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// Test Kms ClientKey. >>> Resource test cases, automatically generated.
// Case 4119
func TestAccAliCloudKmsClientKey_basic4119(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_kms_client_key.default"
	ra := resourceAttrInit(resourceId, AlicloudKmsClientKeyMap4119)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &KmsServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeKmsClientKey")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%skmsclientkey%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudKmsClientKeyBasicDependence4119)
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
					"aap_name":              "${alicloud_kms_application_access_point.AAP0.application_access_point_name}",
					"password":              "YouPassword123!",
					"not_before":            "2023-09-01T14:11:22Z",
					"not_after":             "2028-09-01T14:11:22Z",
					"private_key_data_file": "./hello.txt",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"aap_name":   CHECKSET,
						"password":   "YouPassword123!",
						"not_before": "2023-09-01T14:11:22Z",
						"not_after":  "2028-09-01T14:11:22Z",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"password", "private_key_data_file"},
			},
		},
	})
}

var AlicloudKmsClientKeyMap4119 = map[string]string{
	"create_time": CHECKSET,
}

func AlicloudKmsClientKeyBasicDependence4119(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

resource "alicloud_kms_application_access_point" "AAP0" {
  application_access_point_name = var.name
  description = var.name
  policies = ["aaa"]
}

`, name)
}

// Test Kms ClientKey. <<< Resource test cases, automatically generated.
