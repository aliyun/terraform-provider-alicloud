package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAliCloudMseNacosConfig_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_mse_nacos_config.default"
	checkoutSupportedRegions(t, true, connectivity.MSESupportRegions)

	ra := resourceAttrInit(resourceId, AlicloudMseNacosConfigMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &MseService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeMseNacosConfig")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc-msenacosconfig%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudMseNacosConfigBasicDependence0)
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
					"instance_id":     "${alicloud_mse_cluster.example.id}",
					"data_id":         "${var.name}:dataId",
					"group":           "${var.name}:group",
					"namespace_id":    "${alicloud_mse_engine_namespace.example.namespace_id}",
					"content":         "test",
					"app_name":        "test",
					"desc":            "test",
					"accept_language": "zh",
					"type":            "text",
					"tags":            "test",
					"beta_ips":        "",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"data_id":  name + ":dataId",
						"group":    name + ":group",
						"content":  "test",
						"app_name": "test",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"content": "test_update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"content": "test_update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"app_name": "test_update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"app_name": "test_update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"desc": "test_update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"desc": "test_update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": "test_update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags": "test_update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"type": "xml",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"type": "xml",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"beta_ips": "test_update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"beta_ips": "test_update",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"accept_language", "beta_ips"},
			},
		},
	})
}

var AlicloudMseNacosConfigMap0 = map[string]string{}

func AlicloudMseNacosConfigBasicDependence0(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

data "alicloud_zones" "example" {
  available_resource_creation = "VSwitch"
}

resource "alicloud_vpc" "example" {
  vpc_name   = "terraform-example"
  cidr_block = "172.17.3.0/24"
}

resource "alicloud_vswitch" "example" {
  vswitch_name = "terraform-example"
  cidr_block   = "172.17.3.0/24"
  vpc_id       = alicloud_vpc.example.id
  zone_id      = data.alicloud_zones.example.zones.0.id
}

resource "alicloud_mse_cluster" "example" {
  connection_type       = "slb"
  net_type              = "privatenet"
  vswitch_id            = alicloud_vswitch.example.id
  cluster_specification = "MSE_SC_1_2_60_c"
  cluster_version       = "NACOS_2_0_0"
  instance_count        = "3"
  pub_network_flow      = "1"
  cluster_alias_name    = var.name
  mse_version           = "mse_pro"
  cluster_type          = "Nacos-Ans"
}

resource "alicloud_mse_engine_namespace" "example" {
  instance_id          = alicloud_mse_cluster.example.id
  namespace_show_name  = var.name
  namespace_id         = var.name
}

`, name)
}
