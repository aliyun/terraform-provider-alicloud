package clickhouse_test

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/clickhouse"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// Test ClickHouse EnterpriseDbClusterPublicEndpoint. >>> Resource test cases, automatically generated.
// Case 企业版CK开公网1-线上 10561
func TestAccAliCloudClickHouseEnterpriseDbClusterPublicEndpoint_basic10561(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_click_house_enterprise_db_cluster_public_endpoint.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudClickHouseEnterpriseDbClusterPublicEndpointMap10561)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return clickhouse.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeClickHouseEnterpriseDbClusterPublicEndpoint")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccclickhouse%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudClickHouseEnterpriseDbClusterPublicEndpointBasicDependence10561)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-beijing"})
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"db_instance_id":           "${alicloud_click_house_enterprise_db_cluster.defaultaqnt22.id}",
					"net_type":                 "Public",
					"connection_string_prefix": "${alicloud_click_house_enterprise_db_cluster.defaultaqnt22.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"db_instance_id":           acctest.CHECKSET,
						"net_type":                 "Public",
						"connection_string_prefix": acctest.CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"connection_string_prefix": "${alicloud_click_house_enterprise_db_cluster.defaultaqnt22.id}8",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"connection_string_prefix": acctest.CHECKSET,
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

var AlicloudClickHouseEnterpriseDbClusterPublicEndpointMap10561 = map[string]string{}

func AlicloudClickHouseEnterpriseDbClusterPublicEndpointBasicDependence10561(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

variable "region_id" {
  default = "cn-beijing"
}

variable "vsw_ip_range_i" {
  default = "172.16.1.0/24"
}

variable "vpc_ip_range" {
  default = "172.16.0.0/12"
}

variable "zone_id_i" {
  default = "cn-beijing-i"
}

resource "alicloud_vpc" "defaultktKLuM" {
  cidr_block = var.vpc_ip_range
}

resource "alicloud_vswitch" "defaultTQWN3k" {
  vpc_id     = alicloud_vpc.defaultktKLuM.id
  zone_id    = var.zone_id_i
  cidr_block = var.vsw_ip_range_i
}

resource "alicloud_click_house_enterprise_db_cluster" "defaultaqnt22" {
  zone_id    = var.zone_id_i
  vpc_id     = alicloud_vpc.defaultktKLuM.id
  scale_min  = "8"
  scale_max  = "16"
  vswitch_id = alicloud_vswitch.defaultTQWN3k.id
}


`, name)
}

// Test ClickHouse EnterpriseDbClusterPublicEndpoint. <<< Resource test cases, automatically generated.
