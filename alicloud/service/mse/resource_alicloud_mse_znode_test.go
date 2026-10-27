package mse_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/mse"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccAliCloudMSEZnode_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_mse_znode.default"
	acctest.CheckoutSupportedRegions(t, true, connectivity.MSESupportRegions)
	ra := acctest.ResourceAttrInit(resourceId, AlicloudMSEZnodeMap0)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return mse.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeMseZnode")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("/tf-testaccmseznode%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudMSEZnodeBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             testAccCheckMseZnodeDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"accept_language": "zh",
					"path":            "${var.name}",
					"data":            "${var.name}",
					"cluster_id":      "${alicloud_mse_cluster.default.cluster_id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"path":       name,
						"data":       name,
						"cluster_id": acctest.CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"data": "${var.name}_update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"data": name + "_update",
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

// A znode cannot outlive its parent cluster. After destroy the parent cluster
// reaches DESTROY_SUCCESS but remains queryable; QueryZnodeDetail then fails
// with InvalidParameter ("cluster not found") and DescribeMseZnode reports
// DefaultError instead of NotFound. testAccCheckMseZnodeDestroy does not
// blanket-swallow that error as NotFound: it bounded-retries DescribeMseZnode
// until the znode is NotFound, and only when the query fails it confirms the
// znode is gone by proving its parent cluster is destroyed (NotFound or
// DESTROY_SUCCESS); if the parent is still in a non-destroyed state it keeps
// waiting rather than swallowing the error.
func testAccCheckMseZnodeDestroy(s *terraform.State) error {
	client := acctest.Provider.Meta().(*connectivity.AliyunClient)
	mseService := mse.NewService(client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_mse_znode" {
			continue
		}
		if rs.Primary.ID == "" {
			continue
		}
		parts, err := helper.ParseResourceId(rs.Primary.ID, 2)
		if err != nil {
			return fmt.Errorf("invalid mse znode id %q: %s", rs.Primary.ID, err)
		}
		clusterId := parts[0]
		wait := helper.IncrementalWait(5*time.Second, 10*time.Second)
		if err := retry.Retry(10*time.Minute, func() *retry.RetryError {
			_, err := mseService.DescribeMseZnode(rs.Primary.ID)
			if err == nil {
				wait()
				return retry.RetryableError(fmt.Errorf("mse znode %s still exists after delete", rs.Primary.ID))
			}
			if errs.NotFoundError(err) {
				return nil
			}
			// QueryZnodeDetail failed (typically the parent cluster is
			// destroyed). Confirm the znode is gone by proving its parent
			// cluster is destroyed (NotFound or DESTROY_SUCCESS); otherwise
			// keep waiting rather than swallowing the error.
			cobj, cerr := mseService.DescribeMseCluster(clusterId)
			if cerr != nil {
				if errs.NotFoundError(cerr) {
					return nil
				}
				return retry.NonRetryableError(cerr)
			}
			if fmt.Sprint(cobj["InitStatus"]) == "DESTROY_SUCCESS" {
				return nil
			}
			wait()
			return retry.RetryableError(err)
		}); err != nil {
			return fmt.Errorf("mse znode %s was not destroyed: %s", rs.Primary.ID, err)
		}
	}
	return nil
}

var AlicloudMSEZnodeMap0 = map[string]string{
	"path":       acctest.CHECKSET,
	"cluster_id": acctest.CHECKSET,
}

func AlicloudMSEZnodeBasicDependence0(name string) string {
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
  cluster_type = "ZooKeeper"
  cluster_version = "ZooKeeper_3_8_0"
  instance_count = 1
  net_type = "privatenet"
  vswitch_id = alicloud_vswitch.default.id
  pub_network_flow = "1"
  cluster_alias_name= "tf-testaccmseznode"
  mse_version = "mse_dev"
}
`, name)
}
