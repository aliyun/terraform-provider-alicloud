package alicloud

import (
	"fmt"
	"log"
	"strings"
	"testing"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/gpdb"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

// Acceptance Test
// Create by yewei.oyyw@alibaba-inc.com on 2019-05-31

func init() {
	resource.AddTestSweepers("alicloud_gpdb_instance", &resource.Sweeper{
		Name: "alicloud_gpdb_instance",
		F:    testSweepGpdbInstances,
	})
}

func testSweepGpdbInstances(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return WrapError(err)
	}
	client := rawClient.(*connectivity.AliyunClient)

	prefixes := []string{
		"tf-testAcc",
		"tf_testAcc",
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
			return WrapErrorf(err, DefaultErrorMsg, "testSweepGpdbInstances", request.GetActionName(), AlibabaCloudSdkGoERROR)
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

		sweeper = true

		// Delete Instance
		request := gpdb.CreateDeleteDBInstanceRequest()
		request.DBInstanceId = id
		raw, err := client.WithGpdbClient(func(gpdbClient *gpdb.Client) (interface{}, error) {
			return gpdbClient.DeleteDBInstance(request)
		})
		if err != nil {
			log.Printf("[error] Failed to delete GPDB instance, ID:%v(%v)\n", id, request.GetActionName())
		}
		addDebug(request.GetActionName(), raw)
	}
	if sweeper {
		// Waiting 30 seconds to ensure these DB instances have been deleted.
		time.Sleep(30 * time.Second)
	}
	return nil
}

func TestAccAlicloudGpdbInstance_classic(t *testing.T) {
	var v gpdb.DBInstanceAttribute
	resourceId := "alicloud_gpdb_instance.default"
	serverFunc := func() interface{} {
		return &GpdbService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, serverFunc, "DescribeGpdbInstance")
	ra := resourceAttrInit(resourceId, nil)
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	testAccConfig := resourceTestAccConfigFunc(resourceId, "", resourceGpdbClassicConfigDependence)

	resource.Test(t, resource.TestCase{
		PreCheck:      func() { testAccPreCheckWithRegions(t, false, connectivity.GpdbClassicNoSupportedRegions) },
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"availability_zone":    "${data.alicloud_zones.default.zones.0.id}",
					"engine":               "gpdb",
					"engine_version":       "4.3",
					"instance_class":       "gpdb.group.segsdx2",
					"instance_group_count": "2",
					"description":          "tf-testAccGpdbInstance_new",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": fmt.Sprintf("tf-testAccGpdbInstance_new"),
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
			// change description
			{
				Config: testAccConfig(map[string]interface{}{
					"description": "tf-testAccGpdbInstance_test",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": "tf-testAccGpdbInstance_test",
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
						"security_ip_list.#":          "1",
						"security_ip_list.4095458986": "10.168.1.12",
					}),
				),
			},
		}})
}

func TestAccAlicloudGpdbInstance_vpc(t *testing.T) {
	var v gpdb.DBInstanceAttribute
	resourceId := "alicloud_gpdb_instance.default"
	serverFunc := func() interface{} {
		return &GpdbService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, serverFunc, "DescribeGpdbInstance")
	ra := resourceAttrInit(resourceId, nil)
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	testAccConfig := resourceTestAccConfigFunc(resourceId, "", resourceGpdbVpcConfigDependence)

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
					"availability_zone":    "${data.alicloud_zones.default.zones.0.id}",
					"vswitch_id":           "${alicloud_vswitch.default.id}",
					"engine":               "gpdb",
					"engine_version":       "4.3",
					"instance_class":       "gpdb.group.segsdx2",
					"instance_group_count": "2",
					"description":          "tf-testAccGpdbInstance_new",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": fmt.Sprintf("tf-testAccGpdbInstance_new"),
					}),
				),
			},
			// change description
			{
				Config: testAccConfig(map[string]interface{}{
					"description": "tf-testAccGpdbInstance_test",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": "tf-testAccGpdbInstance_test",
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
						"security_ip_list.#":          "1",
						"security_ip_list.4095458986": "10.168.1.12",
					}),
				),
			},
		}})
}

func resourceGpdbClassicConfigDependence(s string) string {
	return fmt.Sprintf(`
        data "alicloud_zones" "default" {
            available_resource_creation = "Gpdb"
        }`)
}

func resourceGpdbVpcConfigDependence(s string) string {
	return fmt.Sprintf(`
        data "alicloud_zones" "default" {
            available_resource_creation = "Gpdb"
        }
        variable "name" {
            default                = "tf-testAccGpdbInstance_vpc"
        }
        resource "alicloud_vpc" "default" {
            name                   = "${var.name}"
            cidr_block             = "172.16.0.0/16"
        }
        resource "alicloud_vswitch" "default" {
            availability_zone      = "${data.alicloud_zones.default.zones.0.id}"
            vpc_id                 = "${alicloud_vpc.default.id}"
            cidr_block             = "172.16.0.0/24"
            name                   = "${var.name}"
        }`)
}
