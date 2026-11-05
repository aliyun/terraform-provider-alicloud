package polardb_test

import (
	"fmt"
	"testing"

	sdkpolardb "github.com/aliyun/alibaba-cloud-sdk-go/services/polardb"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	polardb "github.com/aliyun/terraform-provider-alicloud/alicloud/service/polardb"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccAliCloudPolarDBEndpointAddressConfigUpdate(t *testing.T) {
	var v *sdkpolardb.Address
	rand := tfacctest.RandStringFromCharSet(10, tfacctest.CharSetAlphaNum)
	name := fmt.Sprintf("tf-testAccPolarDBconnection%s", rand)
	var basicMap = map[string]string{
		"db_cluster_id":  CHECKSET,
		"db_endpoint_id": CHECKSET,
		"net_type":       "Public",
	}
	resourceId := "alicloud_polardb_endpoint_address.default"
	ra := acctest.ResourceAttrInit(resourceId, basicMap)
	serviceFunc := func() interface{} {
		return polardb.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, serviceFunc, "DescribePolarDBConnection")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, resourcePolarDBEndpointAddressConfigDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},

		// module name

		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"db_cluster_id":  "${alicloud_polardb_cluster.default.id}",
					"db_endpoint_id": "${data.alicloud_polardb_endpoints.default1.endpoints[0].db_endpoint_id}",
					"net_type":       "Public",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(nil),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"connection_prefix": "tf-testacc",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"connection_prefix": "tf-testacc",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"port": "3307",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"port": "3307",
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

func resourcePolarDBEndpointAddressConfigDependence(name string) string {
	return fmt.Sprintf(`
	%s
	variable "name" {
		default = "%s"
	}

	resource "alicloud_polardb_cluster" "default" {
		db_type       = "MySQL"
		db_version    = "8.0"
		pay_type      = "PostPaid"
		db_node_class = "polar.mysql.x4.large"
		vswitch_id    = local.vswitch_id
		description   = "${var.name}"
		lifecycle {
			ignore_changes = [encrypt_new_tables]
		}
	}

	data "alicloud_polardb_endpoints" "default1" {
	  db_cluster_id    = "${alicloud_polardb_cluster.default.id}"
	}
	`, PolarDBCommonTestCase, name)
}
