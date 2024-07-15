package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAliCloudRamLoginProfile_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ram_login_profile.default"
	ra := resourceAttrInit(resourceId, AliCloudRamLoginProfileMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &RamService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeRamLoginProfile")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sRamLoginProfileConfig%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudRamLoginProfileBasicDependence0)
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
					"user_name": "${alicloud_ram_user.default.name}",
					"password":  "YourPassword123!",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"user_name": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"password_reset_required": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"password_reset_required": "true",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"mfa_bind_required": "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"mfa_bind_required": "false",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"password_reset_required": "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"password_reset_required": "false",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"mfa_bind_required": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"mfa_bind_required": "true",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"password"},
			},
		},
	})
}

func TestAccAliCloudRamLoginProfile_basic0_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ram_login_profile.default"
	ra := resourceAttrInit(resourceId, AliCloudRamLoginProfileMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &RamService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeRamLoginProfile")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sRamLoginProfileConfig%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudRamLoginProfileBasicDependence0)
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
					"user_name":               "${alicloud_ram_user.default.name}",
					"password":                "YourPassword123!",
					"password_reset_required": "true",
					"mfa_bind_required":       "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"user_name":               CHECKSET,
						"password_reset_required": "true",
						"mfa_bind_required":       "true",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"password"},
			},
		},
	})
}

var AliCloudRamLoginProfileMap0 = map[string]string{}

func AliCloudRamLoginProfileBasicDependence0(name string) string {
	return fmt.Sprintf(` 
	variable "name" {
  		default = "%s"
	}

	resource "alicloud_ram_user" "default" {
  		name         = var.name
  		display_name = "displayname"
  		mobile       = "86-18888888888"
  		email        = "hello.uuu@aaa.com"
  		comments     = "yoyoyo"
	}
`, name)
}
