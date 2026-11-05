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

// Case Dataworks网络资源管理-TF验收_北京 8963
func TestAccAliCloudDataWorksNetwork_basic8963_modify(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_data_works_network.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudDataWorksNetworkMap8963)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return dataworks.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeDataWorksNetwork")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf_testacc_dwpt%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudDataWorksNetworkBasicDependence8963)
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
					"vpc_id":               "${data.alicloud_vpcs.default2.ids.0}",
					"vswitch_id":           "${data.alicloud_vswitches.default2.ids.0}",
					"dw_resource_group_id": "${alicloud_data_works_dw_resource_group.defaultVJvKvl.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"vpc_id":               acctest.CHECKSET,
						"vswitch_id":           acctest.CHECKSET,
						"dw_resource_group_id": acctest.CHECKSET,
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

var AlicloudDataWorksNetworkMap8963 = map[string]string{
	"status":      acctest.CHECKSET,
	"create_time": acctest.CHECKSET,
}

func AlicloudDataWorksNetworkBasicDependence8963(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

data "alicloud_vpcs" "default" {
	name_regex = "^default-NODELETING$"
}

data "alicloud_vswitches" "default" {
	vpc_id  = data.alicloud_vpcs.default.ids.0
}

data "alicloud_vpcs" "default2" {
	name_regex = "^default-NODELETING-2$"
}

data "alicloud_vswitches" "default2" {
	vpc_id  = data.alicloud_vpcs.default2.ids.0
}


resource "alicloud_data_works_dw_resource_group" "defaultVJvKvl" {
  payment_duration_unit = "Month"
  payment_type          = "PostPaid"
  specification         = "500"
  default_vswitch_id    = data.alicloud_vswitches.default.ids.0
  remark                = "OpenAPI测试用资源组"
  resource_group_name   = var.name
  default_vpc_id        = data.alicloud_vpcs.default.ids.0
  auto_renew = false
}

`, name)
}
