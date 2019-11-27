package alicloud

import (
	"fmt"
	"log"
	"strings"
	"testing"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/yundun_dbaudit"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func init() {
	resource.AddTestSweepers("alicloud_yundun_dbaudit_instance", &resource.Sweeper{
		Name: "alicloud_yundun_dbaudit_instance",
		F:    testSweepDbauditInstances,
	})
}

func testSweepDbauditInstances(region string) error {
	if testSweepPreCheckWithRegions(region, true, connectivity.YundunDbauditSupportedRegions) {
		log.Printf("[INFO] Skipping Dbaudit Instance unsupported region: %s", region)
		return nil
	}
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting Alicloud client: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)

	prefixes := []string{
		fmt.Sprintf("tf-testAcc%s", region),
		fmt.Sprintf("tf_testAcc%s", region),
	}
	request := yundun_dbaudit.CreateDescribeInstancesRequest()
	request.PageSize = requests.NewInteger(PageSizeSmall)
	request.CurrentPage = requests.NewInteger(1)
	var instances []yundun_dbaudit.Instance

	for {
		raw, err := client.WithDbauditClient(func(dbauditClient *yundun_dbaudit.Client) (interface{}, error) {
			return dbauditClient.DescribeInstances(request)
		})
		if err != nil {
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_yundun_dbaudit", request.GetActionName(), AlibabaCloudSdkGoERROR)
		}

		addDebug(request.GetActionName(), raw)
		response, _ := raw.(*yundun_dbaudit.DescribeInstancesResponse)
		if len(response.Instances) < 1 {
			break
		}

		instances = append(instances, response.Instances...)

		if len(response.Instances) < PageSizeSmall {
			break
		}

		currentPageNo := request.CurrentPage
		if err != nil {
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_yundun_dbaudit", request.GetActionName(), AlibabaCloudSdkGoERROR)
		}

		if page, err := getNextpageNumber(currentPageNo); err != nil {
			return WrapError(err)
		} else {
			request.CurrentPage = page
		}
	}

	for _, v := range instances {
		name := v.Description
		skip := true
		for _, prefix := range prefixes {
			if name != "" && strings.HasPrefix(strings.ToLower(name), strings.ToLower(prefix)) {
				skip = false
				break
			}
		}
		if skip {
			log.Printf("[INFO] Skipping Dbaudit Instance: %s", name)
			continue
		}

		log.Printf("[INFO] Deleting Dbaudit Instance %s .", v.InstanceId)

		releaseReq := yundun_dbaudit.CreateRefundInstanceRequest()
		releaseReq.InstanceId = v.InstanceId
		//releaseReq.InstanceId = v.InstanceId
		//
		_, err := client.WithDbauditClient(func(dbauditClient *yundun_dbaudit.Client) (interface{}, error) {
			return dbauditClient.RefundInstance(releaseReq)
		})
		if err != nil {
			log.Printf("[ERROR] Deleting Instance %s got an error: %#v.", v.InstanceId, err)
		}
	}

	return nil
}

func TestAccAlicloudYundunDbauditInstance_basic(t *testing.T) {
	var v yundun_dbaudit.Instance
	resourceId := "alicloud_yundun_dbaudit_instance.default"
	ra := resourceAttrInit(resourceId, dbauditInstanceBasicMap)

	serviceFunc := func() interface{} {
		return &DbauditService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000, 9999)
	name := fmt.Sprintf("tf_testAcc%d", rand)

	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceDbauditInstanceDependence)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, connectivity.YundunDbauditSupportedRegions)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"description": "${var.name}",
					"plan_code":   "alpha.professional",
					"period":      "1",
					"vswitch_id":  "${alicloud_vswitch.default.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": name,
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: false,
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": "${var.name}_update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": name + "_update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": "${var.name}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": name,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"plan_code": "alpha.basic",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"plan_code": "alpha.basic",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"plan_code": "alpha.premium",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"plan_code": "alpha.premium",
					}),
				),
			},
		},
	})
}

func TestAccAlicloudYundunDbauditInstance_Multi(t *testing.T) {
	var v yundun_dbaudit.Instance

	resourceId := "alicloud_yundun_dbaudit_instance.default.1"
	ra := resourceAttrInit(resourceId, dbauditInstanceBasicMap)

	serviceFunc := func() interface{} {
		return &DbauditService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000, 9999)
	name := fmt.Sprintf("tf_testAcc%d", rand)

	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceDbauditInstanceDependence)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, connectivity.YundunDbauditSupportedRegions)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"count":       "5",
					"description": "${var.name}-${count.index}",
					"plan_code":   "alpha.professional",
					"period":      "1",
					"vswitch_id":  "${alicloud_vswitch.default.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(nil),
				),
			},
		},
	})
}

func resourceDbauditInstanceDependence(name string) string {
	return fmt.Sprintf(
		`  data "alicloud_zones" "default" {
    				available_resource_creation = "VSwitch"
			  }

			  variable "name" {
				default = "%s"
			  }

			  resource "alicloud_vpc" "default" {
				name = "${var.name}"
				cidr_block = "172.16.0.0/12"
			  }

			  resource "alicloud_vswitch" "default" {
				vpc_id = "${alicloud_vpc.default.id}"
				cidr_block = "172.16.0.0/21"
				availability_zone = "${data.alicloud_zones.default.zones.0.id}"
				name = "${var.name}"
			  }
			
			  provider "alicloud" {
				endpoints {
					bssopenapi = "business.aliyuncs.com"
					}
			  }`, name)
}

var dbauditInstanceBasicMap = map[string]string{
	"description": CHECKSET,
	"plan_code":   "alpha.professional",
	"period":      "1",
	"vswitch_id":  CHECKSET,
}
