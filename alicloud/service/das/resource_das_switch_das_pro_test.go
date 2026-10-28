package das_test

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/das"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// SkipTestAccAliCloudDasSwitchDasPro_basic0 remains disabled (Skip prefix
// preserved from the flat layout). The resource is DEPRECATED since v1.249.0
// and its ACC was already skipped pre-migration; the migration does not
// re-enable it. Kept compiling so the package builds and the known coverage
// gap stays explicit.
func SkipTestAccAliCloudDasSwitchDasPro_basic0(t *testing.T) {
	var v map[string]interface{}
	acctest.CheckoutSupportedRegions(t, true, connectivity.TestSalveRegions)
	resourceId := "alicloud_das_switch_das_pro.default"
	ra := acctest.ResourceAttrInit(resourceId, AliCloudDasSwitchDasProMap0)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return das.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeDasSwitchDasPro")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAccDasSwitchDasPro-name%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AliCloudDasSwitchDasProBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_id":   "${alicloud_polardb_cluster.default.id}",
					"sql_retention": "30",
					"user_id":       "${data.alicloud_account.default.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_id":   acctest.CHECKSET,
						"sql_retention": "30",
						"user_id":       acctest.CHECKSET,
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

var AliCloudDasSwitchDasProMap0 = map[string]string{
	"status": acctest.CHECKSET,
}

func AliCloudDasSwitchDasProBasicDependence0(name string) string {
	return fmt.Sprintf(`
	variable "name" {
  		default = "%s"
	}

	data "alicloud_account" "default" {
	}

	data "alicloud_polardb_zones" "default" {
	}

	data "alicloud_polardb_node_classes" "default" {
  		db_type    = "MySQL"
  		db_version = "8.0"
  		pay_type   = "PostPaid"
  		zone_id    = data.alicloud_polardb_zones.default.ids[length(data.alicloud_polardb_zones.default.ids) - 1]
  		category   = "Normal"
	}

	resource "alicloud_vpc" "default" {
  		vpc_name   = var.name
  		cidr_block = "172.16.0.0/16"
	}

	resource "alicloud_vswitch" "default" {
  		zone_id      = data.alicloud_polardb_zones.default.ids[length(data.alicloud_polardb_zones.default.ids) - 1]
  		vpc_id       = alicloud_vpc.default.id
  		vswitch_name = var.name
  		cidr_block   = "172.16.0.0/24"
	}

	resource "alicloud_polardb_cluster" "default" {
  		db_type       = "MySQL"
  		db_version    = "8.0"
  		pay_type      = "PostPaid"
  		db_node_class = data.alicloud_polardb_node_classes.default.classes.0.supported_engines.0.available_resources.0.db_node_class
  		vswitch_id    = alicloud_vswitch.default.id
  		description   = var.name
  		encrypt_new_tables = "OFF"
	}
`, name)
}
