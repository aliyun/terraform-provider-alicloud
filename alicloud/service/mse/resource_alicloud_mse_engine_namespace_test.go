package mse_test

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/mse"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccAliCloudMSEEngineNamespace_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_mse_engine_namespace.default"
	acctest.CheckoutSupportedRegions(t, true, connectivity.MSESupportRegions)
	ra := acctest.ResourceAttrInit(resourceId, AlicloudMSEEngineNamespaceMap0)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return mse.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeMseEngineNamespace")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc-mseenginenamespace%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudMSEEngineNamespaceBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_id":         "${alicloud_mse_cluster.default.id}",
					"cluster_id":          "${alicloud_mse_cluster.default.cluster_id}",
					"namespace_show_name": "${var.name}",
					"namespace_id":        "${var.name}",
					"accept_language":     "zh",
					"namespace_desc":      "test",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"cluster_id":          acctest.CHECKSET,
						"namespace_show_name": name,
						"namespace_id":        name,
						"namespace_desc":      "test",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"namespace_show_name": "${var.name}_update",
					"namespace_desc":      "test_update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"namespace_show_name": name + "_update",
						"namespace_desc":      "test_update",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"accept_language"},
			},
		},
	})
}

var AlicloudMSEEngineNamespaceMap0 = map[string]string{}

func AlicloudMSEEngineNamespaceBasicDependence0(name string) string {
	return fmt.Sprintf(` 
variable "name" {
  default = "%s"
}

data "alicloud_zones" "default" {
 available_resource_creation = "VSwitch"
}

resource "alicloud_vpc" "default" {
 vpc_name   = "default"
 cidr_block = "172.17.3.0/24"
}

resource "alicloud_vswitch" "default" {
 vswitch_name = "default"
 cidr_block   = "172.17.3.0/24"
 vpc_id       = alicloud_vpc.default.id
 zone_id      = data.alicloud_zones.default.zones.0.id
}

resource "alicloud_mse_cluster" "default" {
 cluster_specification = "MSE_SC_1_2_60_c"
 cluster_type          = "Nacos-Ans"
 cluster_version       = "NACOS_2_0_0"
 instance_count        = 3
 net_type              = "privatenet"
 vswitch_id            = alicloud_vswitch.default.id
 connection_type       = "slb"
 pub_network_flow      = "1"
 mse_version           = "mse_pro"
 vpc_id                = alicloud_vpc.default.id
}

data "alicloud_mse_clusters" "default" {
	name_regex = "tftest"
}
`, name)
}
