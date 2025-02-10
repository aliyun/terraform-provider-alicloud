package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// Test CloudPhone Policy. >>> Resource test cases, automatically generated.
// Case chuyuan_createPolicy_prod_all 10062
func TestAccAliCloudCloudPhonePolicy_basic10062(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cloud_phone_policy.default"
	ra := resourceAttrInit(resourceId, AlicloudCloudPhonePolicyMap10062)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CloudPhoneServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeCloudPhonePolicy")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%scloudphonepolicy%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudCloudPhonePolicyBasicDependence10062)
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
					"policy_group_name": "NewPolicyName",
					"resolution_width":  "720",
					"lock_resolution":   "on",
					"camera_redirect":   "on",
					"resolution_height": "1280",
					"clipboard":         "read",
					"net_redirect_policy": []map[string]interface{}{
						{
							"net_redirect":    "on",
							"custom_proxy":    "on",
							"proxy_type":      "socks5",
							"host_addr":       "192.168.12.13",
							"port":            "8888",
							"proxy_user_name": "user1",
							"proxy_password":  "123456",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"policy_group_name": "NewPolicyName",
						"resolution_width":  "720",
						"lock_resolution":   "on",
						"camera_redirect":   "on",
						"resolution_height": "1280",
						"clipboard":         "read",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"policy_group_name": "ModifyPolicyName",
					"resolution_width":  "1080",
					"resolution_height": "1920",
					"clipboard":         "write",
					"net_redirect_policy": []map[string]interface{}{
						{
							"net_redirect":    "on",
							"custom_proxy":    "on",
							"proxy_type":      "socks5",
							"host_addr":       "192.168.16.13",
							"port":            "9999",
							"proxy_user_name": "user2",
							"proxy_password":  "1234567",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"policy_group_name": "ModifyPolicyName",
						"resolution_width":  "1080",
						"resolution_height": "1920",
						"clipboard":         "write",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"policy_group_name": "defaultPolicyGroup",
					"resolution_width":  "720",
					"lock_resolution":   "off",
					"camera_redirect":   "off",
					"resolution_height": "1280",
					"clipboard":         "read",
					"net_redirect_policy": []map[string]interface{}{
						{
							"net_redirect": "off",
							"custom_proxy": "off",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"policy_group_name": "defaultPolicyGroup",
						"resolution_width":  "720",
						"lock_resolution":   "off",
						"camera_redirect":   "off",
						"resolution_height": "1280",
						"clipboard":         "read",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"resolution_width":  "1080",
					"lock_resolution":   "on",
					"camera_redirect":   "on",
					"resolution_height": "1920",
					"clipboard":         "write",
					"net_redirect_policy": []map[string]interface{}{
						{
							"net_redirect":    "on",
							"custom_proxy":    "on",
							"proxy_type":      "socks5",
							"host_addr":       "195.2.3.2",
							"port":            "1234",
							"proxy_user_name": "user2",
							"proxy_password":  "1234567",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"resolution_width":  "1080",
						"lock_resolution":   "on",
						"camera_redirect":   "on",
						"resolution_height": "1920",
						"clipboard":         "write",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"resolution_width":  "720",
					"resolution_height": "1280",
					"clipboard":         "readwrite",
					"net_redirect_policy": []map[string]interface{}{
						{
							"net_redirect":    "on",
							"custom_proxy":    "on",
							"proxy_type":      "socks5",
							"host_addr":       "123.0.0.1",
							"port":            "8084",
							"proxy_user_name": "user7",
							"proxy_password":  "1234567",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"resolution_width":  "720",
						"resolution_height": "1280",
						"clipboard":         "readwrite",
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

var AlicloudCloudPhonePolicyMap10062 = map[string]string{}

func AlicloudCloudPhonePolicyBasicDependence10062(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}


`, name)
}

// Test CloudPhone Policy. <<< Resource test cases, automatically generated.
