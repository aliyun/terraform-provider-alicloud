package dataworks_test

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/dataworks"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// Test DataWorks DwResourceGroup. >>> Resource test cases, automatically generated.
// Case Dataworks资源组管理_TF验收_北京 8964
func TestAccAliCloudDataWorksDwResourceGroup_basic8964(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_data_works_dw_resource_group.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudDataWorksDwResourceGroupMap8964)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return dataworks.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeDataWorksDwResourceGroup")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf_testacc_dwpt%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudDataWorksDwResourceGroupBasicDependence8964)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-shenzhen"})
			acctest.PreCheck(t)
		},
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"payment_type":          "PostPaid",
					"default_vpc_id":        "${alicloud_vpc.defaulte4zhaL.id}",
					"remark":                "openapi_test",
					"resource_group_name":   name,
					"auto_renew":            "false",
					"default_vswitch_id":    "${alicloud_vswitch.default675v38.id}",
					"payment_duration_unit": "Month",
					"specification":         "500",
					"payment_duration":      "1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"payment_type":          "PostPaid",
						"default_vpc_id":        acctest.CHECKSET,
						"remark":                "openapi_test",
						"resource_group_name":   acctest.CHECKSET,
						"auto_renew":            "false",
						"default_vswitch_id":    acctest.CHECKSET,
						"payment_duration_unit": "Month",
						"specification":         "500",
						"payment_duration":      "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"remark":              "openapi_test_update",
					"resource_group_name": "openapi_pop2_test_resg_update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"remark":              "openapi_test_update",
						"resource_group_name": "openapi_pop2_test_resg_update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Group",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF",
						"tags.For":     "Group",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Group-Update",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF",
						"tags.For":     "Group-Update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.ids.0}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"resource_group_id": acctest.CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.ids.1}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"resource_group_id": acctest.CHECKSET,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"auto_renew", "payment_duration", "payment_duration_unit", "specification"},
			},
		},
	})
}

var AlicloudDataWorksDwResourceGroupMap8964 = map[string]string{
	"status":      acctest.CHECKSET,
	"create_time": acctest.CHECKSET,
}

func AlicloudDataWorksDwResourceGroupBasicDependence8964(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

resource "alicloud_data_works_project" "defaultZImuCO" {
  description  = "default_pj002"
  project_name = var.name
  display_name = "default_pj002"
  pai_task_enabled = false
}

resource "alicloud_vpc" "defaulte4zhaL" {
  description = "default_resgv2_vpc001"
  vpc_name    = format("%%s1", var.name)
  cidr_block  = "172.16.0.0/12"
}

resource "alicloud_vswitch" "default675v38" {
  description  = "default_resg_vsw001"
  vpc_id       = alicloud_vpc.defaulte4zhaL.id
  zone_id      = "cn-shenzhen-f"
  vswitch_name = format("%%s2", var.name)
  cidr_block   = "172.16.0.0/24"
}

data "alicloud_resource_manager_resource_groups" "default" {}
`, name)
}

// Test DataWorks DwResourceGroup. <<< Resource test cases, automatically generated.
