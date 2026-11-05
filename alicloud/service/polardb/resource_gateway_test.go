package polardb_test

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	polardb "github.com/aliyun/terraform-provider-alicloud/alicloud/service/polardb"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccAliCloudPolarDBGateway_basic(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_polardb_gateway.default"
	ra := acctest.ResourceAttrInit(resourceId, map[string]string{
		"status": CHECKSET,
	})
	serviceFunc := func() interface{} {
		return polardb.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, serviceFunc, "DescribePolarDBGatewayAttribute")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAccPolarDBGateway-%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, resourcePolarDBGatewayConfigDependence)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{connectivity.Beijing})
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"zone_id":           "${data.alicloud_vswitches.default.vswitches.0.zone_id}",
					"db_cluster_class":  "polar.app.g2.small",
					"pay_type":          "Prepaid",
					"auto_renew":        "false",
					"period":            "Month",
					"used_time":         "1",
					"vpc_id":            "${data.alicloud_vpcs.default.ids.0}",
					"vswitch_id":        "${data.alicloud_vswitches.default.ids.0}",
					"security_group_id": "${alicloud_security_group.default.id}",
					"db_type":           "PostgreSQL",
				}),
				ExpectNonEmptyPlan: true,
				PlanOnly:           true,
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"zone_id":           "${data.alicloud_vswitches.default.vswitches.0.zone_id}",
					"db_cluster_class":  "polar.app.g2.small",
					"pay_type":          "Postpaid",
					"auto_renew":        "false",
					"vpc_id":            "${data.alicloud_vpcs.default.ids.0}",
					"vswitch_id":        "${data.alicloud_vswitches.default.ids.0}",
					"security_group_id": "${alicloud_security_group.default.id}",
					"db_type":           "PostgreSQL",
				}),
				Check: resource.ComposeTestCheckFunc(
					rac.ResourceAttrMapUpdateSet()(map[string]string{
						"region_id": CHECKSET, "zone_id": CHECKSET, "db_cluster_class": "polar.app.g2.small",
						"pay_type": "Postpaid", "auto_renew": "false", "vpc_id": CHECKSET,
						"vswitch_id": CHECKSET, "security_group_id": CHECKSET, "db_type": "PostgreSQL",
						"status": CHECKSET, "create_time": CHECKSET, "expired": CHECKSET,
						"endpoints.#": "2", "endpoints.0.address": CHECKSET, "endpoints.0.port": CHECKSET,
						"security_ip_arrays.#": "0",
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"zone_id", "db_cluster_class", "pay_type", "auto_renew", "period", "used_time",
					"vpc_id", "vswitch_id", "security_group_id", "db_type", "modify_time",
				},
			},
		},
	})
}

func resourcePolarDBGatewayConfigDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

data "alicloud_vpcs" "default" {
  is_default = true
}

data "alicloud_vswitches" "default" {
  vpc_id  = data.alicloud_vpcs.default.ids.0
  zone_id = "cn-beijing-l"
}

resource "alicloud_security_group" "default" {
  security_group_name = var.name
  vpc_id              = data.alicloud_vpcs.default.ids.0
}
`, name)
}
