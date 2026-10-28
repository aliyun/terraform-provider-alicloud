package mscopensubscription_test

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/mscopensubscription"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccAliCloudMscSubContact_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_msc_sub_contact.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudMscSubContactMap0)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return mscopensubscription.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeContact")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	name := fmt.Sprintf("tftestacc")
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudMscSubContactBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
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

func TestAccAliCloudMscSubContact_basic1(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_msc_sub_contact.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudMscSubContactMap1)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return mscopensubscription.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeContact")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	name := fmt.Sprintf("tftestacc")
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudMscSubContactBasicDependence1)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"contact_name": "张三",
					"position":     "Other",
					"email":        "123@163.com",
					"mobile":       "12345257908",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"contact_name": "张三",
						"position":     "Other",
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

var AlicloudMscSubContactMap1 = map[string]string{}

func AlicloudMscSubContactBasicDependence1(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "张三"
}
`)
}
