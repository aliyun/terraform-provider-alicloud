package alicloud

import (
	"fmt"
	"log"
	"strings"
	"testing"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/yundun_bastionhost"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func init() {
	resource.AddTestSweepers("alicloud_yundun_bastionhost_instance", &resource.Sweeper{
		Name: "alicloud_yundun_bastionhost_instance",
		F:    testSweepBastionhostInstances,
	})
}

func testSweepBastionhostInstances(region string) error {
	if testSweepPreCheckWithRegions(region, true, connectivity.YundunBastionhostSupportedRegions) {
		log.Printf("[INFO] Skipping Bastionhost Instance unsupported region: %s", region)
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
	request := yundun_bastionhost.CreateDescribeInstanceBastionhostRequest()
	request.PageSize = requests.NewInteger(PageSizeSmall)
	request.CurrentPage = requests.NewInteger(1)
	var instances []yundun_bastionhost.Instance

	for {
		raw, err := client.WithBastionhostClient(func(bastionhostClient *yundun_bastionhost.Client) (interface{}, error) {
			return bastionhostClient.DescribeInstanceBastionhost(request)
		})
		if err != nil {
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_yundun_bastionhost", request.GetActionName(), AlibabaCloudSdkGoERROR)
		}

		addDebug(request.GetActionName(), raw)
		response, _ := raw.(*yundun_bastionhost.DescribeInstanceBastionhostResponse)
		if len(response.Instances) < 1 {
			break
		}

		instances = append(instances, response.Instances...)

		if len(response.Instances) < PageSizeSmall {
			break
		}

		currentPageNo := request.CurrentPage
		if err != nil {
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_yundun_bastionhost", request.GetActionName(), AlibabaCloudSdkGoERROR)
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
			log.Printf("[INFO] Skipping Bastionhost Instance: %s", name)
			continue
		}

		log.Printf("[INFO] Deleting Bastionhost Instance %s .", v.InstanceId)

		releaseReq := yundun_bastionhost.CreateRefundInstanceRequest()
		releaseReq.InstanceId = v.InstanceId
		_, err := client.WithBastionhostClient(func(bastionhostClient *yundun_bastionhost.Client) (interface{}, error) {
			return bastionhostClient.RefundInstance(releaseReq)
		})
		if err != nil {
			log.Printf("[ERROR] Deleting Instance %s got an error: %#v.", v.InstanceId, err)
		}
	}

	return nil
}

func TestAccAlicloudYundunBastionhostInstance_basic(t *testing.T) {
	var v yundun_bastionhost.Instance
	resourceId := "alicloud_yundun_bastionhost_instance.default"
	ra := resourceAttrInit(resourceId, bastionhostInstanceBasicMap)

	serviceFunc := func() interface{} {
		return &bastionhostService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000, 9999)
	name := fmt.Sprintf("tf_testAcc%d", rand)

	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceBastionhostInstanceDependence)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, connectivity.YundunBastionhostSupportedRegions)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"description":        "${var.name}",
					"license_code":       "bhah_ent_50_asset",
					"period":             "1",
					"vswitch_id":         "${alicloud_vswitch.default.id}",
					"security_group_ids": []string{"${alicloud_security_group.default.0.id}"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description":          name,
						"security_group_ids.#": "1",
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
					"license_code": "bhah_ent_100_asset",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"license_code": "bhah_ent_100_asset",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"security_group_ids": []string{"${alicloud_security_group.default.1.id}"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"security_group_ids.#": "1",
					}),
				),
			},
		},
	})
}

func TestAccAlicloudYundunBastionhostInstance_Multi(t *testing.T) {
	var v yundun_bastionhost.Instance

	resourceId := "alicloud_yundun_bastionhost_instance.default.1"
	ra := resourceAttrInit(resourceId, bastionhostInstanceBasicMap)

	serviceFunc := func() interface{} {
		return &bastionhostService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000, 9999)
	name := fmt.Sprintf("tf_testAcc%d", rand)

	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceBastionhostInstanceDependence)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, connectivity.YundunBastionhostSupportedRegions)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"count":              "2",
					"description":        "${var.name}-${count.index}",
					"period":             "1",
					"vswitch_id":         "${alicloud_vswitch.default.id}",
					"license_code":       "bhah_ent_50_asset",
					"security_group_ids": []string{"${alicloud_security_group.default.0.id}"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(nil),
				),
			},
		},
	})
}

func resourceBastionhostInstanceDependence(name string) string {
	return fmt.Sprintf(
		`data "alicloud_zones" "default" {
				  available_resource_creation = "VSwitch"
				}
				
				variable "name" {
				  default = "%s"
				}
				
				resource "alicloud_vpc" "default" {
				  name       = "${var.name}"
				  cidr_block = "172.16.0.0/12"
				}
				
				resource "alicloud_vswitch" "default" {
				  vpc_id            = "${alicloud_vpc.default.id}"
				  cidr_block        = "172.16.0.0/21"
				  availability_zone = "${data.alicloud_zones.default.zones.0.id}"
				  name              = "${var.name}"
				}
				
				resource "alicloud_security_group" "default" {
				  count  = 2
				  name   = "${var.name}"
				  vpc_id = "${alicloud_vpc.default.id}"
				}
				
				provider "alicloud" {
				  endpoints {
					bssopenapi = "business.aliyuncs.com"
				  }
				}`, name)
}

var bastionhostInstanceBasicMap = map[string]string{
	"description":          CHECKSET,
	"license_code":         "bhah_ent_50_asset",
	"period":               "1",
	"vswitch_id":           CHECKSET,
	"security_group_ids.#": "1",
}
