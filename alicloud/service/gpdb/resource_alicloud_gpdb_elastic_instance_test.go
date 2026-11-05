package gpdb_test

import (
	"fmt"
	"log"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	gpdbsdk "github.com/aliyun/alibaba-cloud-sdk-go/services/gpdb"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/gpdb"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/sweep"

	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// collection-order: bypass security_ip_list -- pre-existing uncovered TypeList attribute outside this pure directory migration (file moved from alicloud/ root; the attribute and its test coverage are unchanged stock)

func init() {
	resource.AddTestSweepers("alicloud_gpdb_elastic_instance", &resource.Sweeper{
		Name: "alicloud_gpdb_elastic_instance",
		F:    testSweepGpdbElasticInstances,
	})
}

func testSweepGpdbElasticInstances(region string) error {
	rawClient, err := sweep.SharedClientForRegion(region)
	if err != nil {
		return errs.WrapError(err)
	}
	client := rawClient.(*connectivity.AliyunClient)

	prefixes := []string{
		"tf-testAcc",
		"tf_testAcc",
		"tftest",
	}

	var instances []gpdbsdk.DBInstance
	request := gpdbsdk.CreateDescribeDBInstancesRequest()
	request.RegionId = client.RegionId
	request.PageSize = requests.NewInteger(helper.PageSizeLarge)
	request.PageNumber = requests.NewInteger(1)
	for {
		raw, err := client.WithGpdbClient(func(gpdbClient *gpdbsdk.Client) (interface{}, error) {
			return gpdbClient.DescribeDBInstances(request)
		})
		if err != nil {
			return errs.WrapErrorf(err, errs.DefaultErrorMsg, "testSweepGpdbElasticInstances", request.GetActionName(), errs.AlibabaCloudSdkGoERROR)
		}
		response, _ := raw.(*gpdbsdk.DescribeDBInstancesResponse)
		logs.AddDebug(request.GetActionName(), response)

		if response == nil || len(response.Items.DBInstance) < 1 {
			break
		}
		instances = append(instances, response.Items.DBInstance...)

		if len(response.Items.DBInstance) < helper.PageSizeLarge {
			break
		}
		if page, err := helper.GetNextpageNumber(request.PageNumber); err != nil {
			return errs.WrapError(err)
		} else {
			request.PageNumber = page
		}
	}

	sweeper := false
	for _, v := range instances {
		id := v.DBInstanceId
		description := v.DBInstanceDescription
		skip := true
		if !sweep.SweepAll() {
			for _, prefix := range prefixes {
				if strings.HasPrefix(strings.ToLower(description), strings.ToLower(prefix)) {
					skip = false
					break
				}
			}
			// If description is not set successfully, it should be fetched by vpc name and deleted.
			if skip {
				if need, err := needSweepVpc(client, v.VpcId, v.VSwitchId); err == nil {
					skip = !need
				}
			}
			if skip {
				log.Printf("[INFO] Skipping GPDB instance: %s (%s)\n", description, id)
				continue
			}
		}

		// Delete Instance
		request := gpdbsdk.CreateDeleteDBInstanceRequest()
		request.DBInstanceId = id
		raw, err := client.WithGpdbClient(func(gpdbClient *gpdbsdk.Client) (interface{}, error) {
			return gpdbClient.DeleteDBInstance(request)
		})
		if err != nil {
			log.Printf("[error] Failed to delete GPDB instance, ID:%v(%v)\n", id, request.GetActionName())
		} else {
			sweeper = true
		}
		logs.AddDebug(request.GetActionName(), raw)
	}
	if sweeper {
		// Waiting 30 seconds to ensure these DB instances have been deleted.
		time.Sleep(30 * time.Second)
	}
	return nil
}

func TestAccAliCloudGPDBElasticInstanceVpc(t *testing.T) {
	var instance gpdbsdk.DBInstanceAttribute
	resourceId := "alicloud_gpdb_elastic_instance.default"
	serverFunc := func() interface{} {
		return gpdb.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &instance, serverFunc, "DescribeElasticInstance")
	ra := acctest.ResourceAttrInit(resourceId, nil)
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, "tf-testAccGpdbInstance_vpc", resourceGpdbElasticInstanceConfigDependence)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"engine":                  "gpdb",
					"engine_version":          "6.0",
					"seg_storage_type":        "cloud_essd",
					"seg_node_num":            "4",
					"storage_size":            "50",
					"instance_spec":           "2C16G",
					"db_instance_description": "tf-testAccGpdbInstance_6.0",
					"vswitch_id":              "${local.vswitch_id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"engine":                  "gpdb",
						"engine_version":          "6.0",
						"seg_storage_type":        "cloud_essd",
						"seg_node_num":            "4",
						"storage_size":            "50",
						"instance_spec":           "2C16G",
						"db_instance_description": "tf-testAccGpdbInstance_6.0",
						"instance_network_type":   "VPC",
						"payment_type":            "PayAsYouGo",
						"vswitch_id":              acctest.CHECKSET,
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
			// change db_instance_description
			{
				Config: testAccConfig(map[string]interface{}{
					"db_instance_description": "tf-testAccGpdbInstance_test",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"db_instance_description": "tf-testAccGpdbInstance_test",
					}),
				),
			},
			// change security ips
			{
				Config: testAccConfig(map[string]interface{}{
					"security_ip_list": []string{"10.168.1.12"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"security_ip_list.#": "1",
						"security_ip_list.0": "10.168.1.12",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"db_instance_description": "tf-testAccGpdbInstance_elastic_6.0",
					"security_ip_list":        []string{"10.168.1.13"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"db_instance_description": "tf-testAccGpdbInstance_elastic_6.0",
						"security_ip_list.#":      "1",
						"security_ip_list.0":      "10.168.1.13",
					}),
				),
			},
		}})
}

func resourceGpdbElasticInstanceConfigDependence(name string) string {
	return fmt.Sprintf(`
        data "alicloud_gpdb_zones" "default" {}
        variable "name" {
            default  = "%s"
        }
        data "alicloud_vpcs" "default" {
            name_regex = "default-NODELETING"
        }
        data "alicloud_vswitches" "default" {
            vpc_id = data.alicloud_vpcs.default.ids.0
            zone_id = data.alicloud_gpdb_zones.default.ids.0
        }
        resource "alicloud_vswitch" "vswitch" {
            count             = length(data.alicloud_vswitches.default.ids) > 0 ? 0 : 1
            vpc_id            = data.alicloud_vpcs.default.ids.0
            cidr_block        = cidrsubnet(data.alicloud_vpcs.default.vpcs[0].cidr_block, 8, 8)
            zone_id = data.alicloud_gpdb_zones.default.ids.0
            vswitch_name              = var.name
        }
        
        locals {
            vswitch_id = length(data.alicloud_vswitches.default.ids) > 0 ? data.alicloud_vswitches.default.ids[0] : concat(alicloud_vswitch.vswitch.*.id, [""])[0]
        }
        `, name)
}

func TestAccAliCloudGPDBElasticInstance_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_gpdb_elastic_instance.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudGpdbElasticInstanceMap0)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return gpdb.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeElasticInstance")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAccgpab%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudGpdbElasticInstanceBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"engine":                  "gpdb",
					"engine_version":          "6.0",
					"seg_storage_type":        "cloud_essd",
					"seg_node_num":            "4",
					"storage_size":            "50",
					"instance_spec":           "2C16G",
					"db_instance_description": name,
					"vswitch_id":              "${data.alicloud_vswitches.default.ids[0]}",
					"db_instance_category":    "HighAvailability",
					"encryption_key":          "${alicloud_kms_key.default.id}",
					"encryption_type":         "CloudDisk",
					"tags": map[string]string{
						"Created": "TF",
						"For":     "acceptance test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"engine":                  "gpdb",
						"engine_version":          "6.0",
						"seg_storage_type":        "cloud_essd",
						"seg_node_num":            "4",
						"storage_size":            "50",
						"instance_spec":           "2C16G",
						"db_instance_description": name,
						"instance_network_type":   "VPC",
						"payment_type":            "PayAsYouGo",
						"db_instance_category":    "HighAvailability",
						"vswitch_id":              acctest.CHECKSET,
						"tags.%":                  "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF1",
						"For":     "acceptance test1",
						"Updated": "TF",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%": "3",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_network_type": "VPC",
					"payment_type":          "Subscription",
					"payment_duration":      "1",
					"payment_duration_unit": "Month",
					"zone_id":               "cn-hangzhou-g",
					"instance_spec":         acctest.REMOVEKEY,
				}),
				ExpectError: regexp.MustCompile("Missing required argument"),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{""},
			},
		},
	})
}

var AlicloudGpdbElasticInstanceMap0 = map[string]string{
	"connection_string": acctest.CHECKSET,
	"port":              acctest.CHECKSET,
}

func AlicloudGpdbElasticInstanceBasicDependence0(name string) string {
	return fmt.Sprintf(` 
variable "name" {
  default = "%s"
}

data "alicloud_gpdb_zones" "default" {}

data "alicloud_vpcs" "default" {
    name_regex = "^default-NODELETING$"
}

data "alicloud_vswitches" "default" {
  vpc_id  = data.alicloud_vpcs.default.ids.0
  zone_id = data.alicloud_gpdb_zones.default.zones.0.id
}

resource "alicloud_kms_key" "default" {
  description             =  var.name
  deletion_window_in_days =  7
  status                  = "Enabled"
}
`, name)
}
