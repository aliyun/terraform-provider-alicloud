package polardb_test

import (
	"fmt"
	"log"
	"strings"
	"testing"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	polardb "github.com/aliyun/terraform-provider-alicloud/alicloud/service/polardb"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/sweep"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func init() {
	resource.AddTestSweepers("alicloud_polardb_global_database_network", &resource.Sweeper{
		Name: "alicloud_polardb_global_database_network",
		F:    testSweepPolarDBGlobalDatabaseNetwork,
	})
}

func testSweepPolarDBGlobalDatabaseNetwork(region string) error {
	rawClient, err := sweep.SharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting Alicloud client: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)
	prefixes := []string{
		"tf-testAcc",
		"tf_testAcc",
	}
	action := "DescribeGlobalDatabaseNetworks"
	request := make(map[string]interface{})
	request["PageSize"] = helper.PageSizeLarge
	request["PageNumber"] = 1
	var response map[string]interface{}
	PolarDBGlobalDatabaseNetworkIds := make([]string, 0)
	for {
		wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
		err = retry.Retry(5*time.Minute, func() *retry.RetryError {
			response, err = client.RpcPost("polardb", "2017-08-01", action, nil, request, true)
			if err != nil {
				if helper.NeedRetry(err) {
					wait()
					return retry.RetryableError(err)
				}
				return retry.NonRetryableError(err)
			}
			return nil
		})
		logs.AddDebug(action, response, request)
		if err != nil {
			return errs.WrapErrorf(err, errs.DataDefaultErrorMsg, "alicloud_polardb_global_database_network", action, errs.AlibabaCloudSdkGoERROR)
		}
		resp, err := jsonpath.Get("$.Items", response)
		if err != nil {
			return errs.WrapErrorf(err, errs.FailedGetAttributeMsg, action, "$.Items", response)
		}
		result, _ := resp.([]interface{})
		for _, v := range result {
			skip := true
			item := v.(map[string]interface{})
			for _, prefix := range prefixes {
				if strings.HasPrefix(strings.ToLower(fmt.Sprint(item["GDNDescription"])), strings.ToLower(prefix)) {
					skip = false
					break
				}
			}
			if skip {
				log.Printf("[INFO] Skipping PolarDBGlobalDatabaseNetwork Instance: %v", item["GDNDescription"])
				continue
			}
			PolarDBGlobalDatabaseNetworkIds = append(PolarDBGlobalDatabaseNetworkIds, fmt.Sprint(item["GDNId"]))
		}
		if len(result) < helper.PageSizeLarge {
			break
		}
		request["PageNumber"] = request["PageNumber"].(int) + 1
	}

	for _, id := range PolarDBGlobalDatabaseNetworkIds {
		log.Printf("[INFO] Deleting PolarDBGlobalDatabaseNetwork Instance: %s", id)
		deleteAction := "DeleteGlobalDatabaseNetwork"
		if err != nil {
			return errs.WrapError(err)
		}
		request = map[string]interface{}{
			"GDNId": id,
		}
		wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
		err = retry.Retry(3*time.Minute, func() *retry.RetryError {
			_, err = client.RpcPost("polardb", "2017-08-01", deleteAction, nil, request, true)
			if err != nil {
				if helper.NeedRetry(err) {
					wait()
					return retry.RetryableError(err)
				}
				return retry.NonRetryableError(err)
			}
			return nil
		})
		if err != nil {
			log.Printf("[ERROR] Failed to delete PolarDBGlobalDatabaseNetwork Instance (%s): %s", PolarDBGlobalDatabaseNetworkIds, err)
		}
	}
	return nil
}

func TestAccAliCloudPolarDBGlobalDatabaseNetwork_basic00(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_polardb_global_database_network.default"
	ra := acctest.ResourceAttrInit(resourceId, ResourceGlobalDatabaseNetworkMap)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return polardb.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribePolarDBGlobalDatabaseNetwork")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAcc%sPolarDBGlobalDatabaseNetwork-name%d", defaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, ResourceGlobalDatabaseNetworkBasicDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"db_cluster_id": "${alicloud_polardb_cluster.default.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"db_cluster_id": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": "tf-testAcc",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": "tf-testAcc",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": "update-tf-testAcc",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": "update-tf-testAcc",
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

var ResourceGlobalDatabaseNetworkMap = map[string]string{
	"status": CHECKSET,
}

func ResourceGlobalDatabaseNetworkBasicDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
	default = "%s"
}

data "alicloud_vpcs" "default" {
	name_regex = "^default-NODELETING$"
}

data "alicloud_polardb_zones" "default" {}

resource "alicloud_vpc" "default" {
	vpc_name = var.name
}

resource "alicloud_vswitch" "default" {
	zone_id = data.alicloud_polardb_zones.default.ids[length(data.alicloud_polardb_zones.default.ids)-1]
	vpc_id = alicloud_vpc.default.id
	cidr_block = cidrsubnet(alicloud_vpc.default.cidr_block, 8, 4)
}

resource "alicloud_polardb_cluster" "default" {
	db_type       = "MySQL"
	db_version    = "8.0"
	pay_type      = "PostPaid"
	db_node_class = "polar.mysql.x4.medium"
	vswitch_id    = alicloud_vswitch.default.id
	description   = "${var.name}"
}

`, name)
}
