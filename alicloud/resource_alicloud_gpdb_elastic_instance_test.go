package alicloud

import (
	"fmt"
	"log"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/gpdb"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func init() {
	resource.AddTestSweepers("alicloud_gpdb_elastic_instance", &resource.Sweeper{
		Name: "alicloud_gpdb_elastic_instance",
		F:    testSweepGpdbElasticInstances,
	})
}

func testSweepGpdbElasticInstances(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return WrapError(err)
	}
	client := rawClient.(*connectivity.AliyunClient)

	prefixes := []string{
		"tf-testAcc",
		"tf_testAcc",
		"tftest",
	}

	var instances []gpdb.DBInstance
	request := gpdb.CreateDescribeDBInstancesRequest()
	request.RegionId = client.RegionId
	request.PageSize = requests.NewInteger(PageSizeLarge)
	request.PageNumber = requests.NewInteger(1)
	for {
		raw, err := client.WithGpdbClient(func(gpdbClient *gpdb.Client) (interface{}, error) {
			return gpdbClient.DescribeDBInstances(request)
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, "testSweepGpdbElasticInstances", request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		response, _ := raw.(*gpdb.DescribeDBInstancesResponse)
		addDebug(request.GetActionName(), response)

		if response == nil || len(response.Items.DBInstance) < 1 {
			break
		}
		instances = append(instances, response.Items.DBInstance...)

		if len(response.Items.DBInstance) < PageSizeLarge {
			break
		}
		if page, err := getNextpageNumber(request.PageNumber); err != nil {
			return WrapError(err)
		} else {
			request.PageNumber = page
		}
	}

	sweeper := false
	service := VpcService{client}
	for _, v := range instances {
		id := v.DBInstanceId
		description := v.DBInstanceDescription
		skip := true
		for _, prefix := range prefixes {
			if strings.HasPrefix(strings.ToLower(description), strings.ToLower(prefix)) {
				skip = false
				break
			}
		}
		// If description is not set successfully, it should be fetched by vpc name and deleted.
		if skip {
			if need, err := service.needSweepVpc(v.VpcId, v.VSwitchId); err == nil {
				skip = !need
			}
		}
		if skip {
			log.Printf("[INFO] Skipping GPDB instance: %s (%s)\n", description, id)
			continue
		}

		// Delete Instance
		request := gpdb.CreateDeleteDBInstanceRequest()
		request.DBInstanceId = id
		raw, err := client.WithGpdbClient(func(gpdbClient *gpdb.Client) (interface{}, error) {
			return gpdbClient.DeleteDBInstance(request)
		})
		if err != nil {
			log.Printf("[error] Failed to delete GPDB instance, ID:%v(%v)\n", id, request.GetActionName())
		} else {
			sweeper = true
		}
		addDebug(request.GetActionName(), raw)
	}
	if sweeper {
		// Waiting 30 seconds to ensure these DB instances have been deleted.
		time.Sleep(30 * time.Second)
	}
	return nil
}

func TestAccAlicloudGPDBElasticInstanceVpc(t *testing.T) {
	var instance gpdb.DBInstanceAttribute
	resourceId := "alicloud_gpdb_elastic_instance.default"
	serverFunc := func() interface{} {
		return &GpdbService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInitWithDescribeMethod(resourceId, &instance, serverFunc, "DescribeGpdbElasticInstance")
	ra := resourceAttrInit(resourceId, nil)
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	testAccConfig := resourceTestAccConfigFunc(resourceId, "tf-testAccGpdbInstance_vpc", resourceGpdbElasticInstanceConfigDependence)

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
						"vswitch_id":              CHECKSET,
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

func TestAccAlicloudGPDBElasticInstance_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_gpdb_elastic_instance.default"
	ra := resourceAttrInit(resourceId, AlicloudGpdbElasticInstanceMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &GpdbService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeGpdbElasticInstance")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAccgpab%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudGpdbElasticInstanceBasicDependence0)
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
						"vswitch_id":              CHECKSET,
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
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{""},
			},
		},
	})
}

var AlicloudGpdbElasticInstanceMap0 = map[string]string{}

func AlicloudGpdbElasticInstanceBasicDependence0(name string) string {
	return fmt.Sprintf(` 
variable "name" {
  default = "%s"
}

data "alicloud_gpdb_zones" "default" {}

data "alicloud_vpcs" "default" {
  name_regex = "default-NODELETING"
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
