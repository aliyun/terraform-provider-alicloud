package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAlicloudMscSubContact_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_msc_sub_contact.default"
	ra := resourceAttrInit(resourceId, AlicloudMscSubContactMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &MscOpenSubscriptionService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeMscSubContact")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	name := fmt.Sprintf("tftestacc")
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudMscSubContactBasicDependence0)
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
					"contact_name": "${var.name}",
					"position":     "CEO",
					"email":        "123@163.com",
					"mobile":       "12345257908",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"contact_name": name,
						"position":     "CEO",
						"email":        "123@163.com",
						"mobile":       "12345257908",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"contact_name": name + "New",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"contact_name": name + "New",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"email": "aba@163.com",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"email": "aba@163.com",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"mobile": "12345257911",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"mobile": "12345257911",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"contact_name": name,
					"email":        "123@163.com",
					"mobile":       "12345257908",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"contact_name": name,
						"email":        "123@163.com",
						"mobile":       "12345257908",
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

var AlicloudMscSubContactMap0 = map[string]string{}

func AlicloudMscSubContactBasicDependence0(name string) string {
	return fmt.Sprintf(` 
variable "name" {
  default = "%s"
}
`, name)
}
