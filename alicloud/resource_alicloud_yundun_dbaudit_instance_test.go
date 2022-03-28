package alicloud

import (
	"fmt"
	"log"
	"strings"
	"testing"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/yundun_dbaudit"
	sls "github.com/aliyun/aliyun-log-go-sdk"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
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
		"tf-testAcc",
		"tf_testAcc",
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
		// 释放产生的 sls project
		_, err = client.WithLogClient(func(slsClient *sls.Client) (interface{}, error) {
			return nil, slsClient.DeleteProject(v.InstanceId)
		})
		if err != nil {
			log.Printf("[ERROR] Failed to delete Log Project (%s): %s", v.InstanceId, err)
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
		CheckDestroy:  nil,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"description":       "${var.name}",
					"plan_code":         "alpha.professional",
					"period":            "1",
					"vswitch_id":        "${data.alicloud_vswitches.default.ids.0}",
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.ids.0}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": name,
						"period":      "1",
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
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.ids.1}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"resource_group_id": CHECKSET,
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
					"tags": map[string]string{
						"Created": "TF",
						"For":     "acceptance test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF",
						"tags.For":     "acceptance test",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF",
						"For":     "acceptance test",
						"Updated": "TF",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "3",
						"tags.Created": "TF",
						"tags.For":     "acceptance test",
						"tags.Updated": "TF",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       REMOVEKEY,
						"tags.Created": REMOVEKEY,
						"tags.For":     REMOVEKEY,
						"tags.Updated": REMOVEKEY,
					}),
				),
			},
		},
	})
}

func resourceDbauditInstanceDependence(name string) string {
	return fmt.Sprintf(`
data "alicloud_zones" "default" {
		available_resource_creation = "VSwitch"
  }
	
data "alicloud_resource_manager_resource_groups" "default"{
	status="OK"
}

  variable "name" {
	default = "%s"
  }

  data "alicloud_vpcs" "default" {
	  name_regex = "default-NODELETING"
  }
  data "alicloud_vswitches" "default" {
	  vpc_id = data.alicloud_vpcs.default.ids.0
	  zone_id      = data.alicloud_zones.default.zones.0.id
  }

 `, name)
}

var dbauditInstanceBasicMap = map[string]string{
	"description": CHECKSET,
	"plan_code":   "alpha.professional",
	"period":      "1",
	"vswitch_id":  CHECKSET,
}
