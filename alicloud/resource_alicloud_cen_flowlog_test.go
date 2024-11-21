package alicloud

import (
	"fmt"
	"log"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/cbn"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func init() {
	resource.AddTestSweepers("alicloud_cen_flowlog", &resource.Sweeper{
		Name: "alicloud_cen_flowlog",
		F:    testSweepCenFlowlog,
	})
}

func testSweepCenFlowlog(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting Alicloud client: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)

	prefixes := []string{
		"tf-testAcc",
		"tf_testAcc",
	}

	var insts []cbn.FlowLog
	request := cbn.CreateDescribeFlowlogsRequest()
	request.RegionId = client.RegionId
	request.PageSize = requests.NewInteger(PageSizeLarge)
	request.PageNumber = requests.NewInteger(1)

	for {
		var raw interface{}
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			raw, err = client.WithCbnClient(func(cbnClient *cbn.Client) (interface{}, error) {
				return cbnClient.DescribeFlowlogs(request)
			})
			if err != nil {
				if IsExpectedErrors(err, []string{ThrottlingUser}) {
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})
		if err != nil {
			log.Printf("[ERROR] Failed to retrieving CEN flowlog: %s", err)
			break
		}
		response, _ := raw.(*cbn.DescribeFlowlogsResponse)
		if len(response.FlowLogs.FlowLog) < 1 {
			break
		}
		insts = append(insts, response.FlowLogs.FlowLog...)

		if len(response.FlowLogs.FlowLog) < PageSizeLarge {
			break
		}

		if page, err := getNextpageNumber(request.PageNumber); err != nil {
			return err
		} else {
			request.PageNumber = page
		}
	}

	sweeped := false
	for _, v := range insts {
		name := v.FlowLogName
		id := v.FlowLogId
		skip := true
		if !sweepAll() {
			for _, prefix := range prefixes {
				if strings.HasPrefix(strings.ToLower(name), strings.ToLower(prefix)) {
					skip = false
					break
				}
			}
			if skip {
				log.Printf("[INFO] Skipping CEN flowlog: %s (%s)", name, id)
				continue
			}
		}
		sweeped = true
		log.Printf("[INFO] Deleting CEN flowlog: %s (%s)", name, id)
		request := cbn.CreateDeleteFlowlogRequest()
		request.FlowLogId = id
		_, err := client.WithCbnClient(func(cbnClient *cbn.Client) (interface{}, error) {
			return cbnClient.DeleteFlowlog(request)
		})
		if err != nil {
			log.Printf("[ERROR] Failed to delete CEN flowlog (%s (%s)): %s", name, id, err)
		}
	}
	if sweeped {
		// Waiting 5 seconds to eusure these flowlog have been deleted.
		time.Sleep(5 * time.Second)
	}
	return nil
}

func SkipTestAccAlicloudCenFlowlog_basic(t *testing.T) {
	// flow log has been offline
	t.Skip("From January 30, 2022, the cloud enterprise network will take the old console flow log function offline. If you need to continue to use the flow log function, you can enter the new version console to use the flow log function of the enterprise version forwarding router. The Enterprise Edition Forwarding Router Flow Log feature provides the same capabilities as the Legacy Console Flow Log feature")
	var cbnFlowlog cbn.FlowLog

	resourceId := "alicloud_cen_flowlog.default"
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testAcc%sFlowlog-%d", defaultRegionToTest, rand)
	nameUpdated := fmt.Sprintf("tf-testAcc%sFlowlog-%d-Update", defaultRegionToTest, rand)

	ra := resourceAttrInit(resourceId, map[string]string{
		"flow_log_name": name,
		"description":   name,
		"status":        "Inactive",
	})

	serviceFunc := func() interface{} {
		return &CbnService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &cbnFlowlog, serviceFunc)

	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceCenFlowlogConfigDependence)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithAccountSiteType(t, DomesticSite)
			testAccPreCheckWithRegions(t, true, connectivity.CenNoSkipRegions)
		},
		// module name
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"cen_id":         "${alicloud_cen_instance.default.id}",
					"project_name":   "${alicloud_log_project.default.name}",
					"log_store_name": "${alicloud_log_store.default.name}",
					"flow_log_name":  name,
					"description":    name,
					"status":         "Inactive",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(nil),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"status": "Active",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"status": "Active",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"flow_log_name": nameUpdated,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"flow_log_name": nameUpdated,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": nameUpdated,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": nameUpdated,
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{
					"flow_log_name": name,
					"description":   name,
					"status":        "Active",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"flow_log_name": name,
						"description":   name,
						"status":        "Active",
					}),
				),
			},
		},
	})
}

// One Cen Instance only can have one flowlog.

func resourceCenFlowlogConfigDependence(name string) string {
	return fmt.Sprintf(`
	variable "name" {
		default = "%s"
	}
	resource "alicloud_cen_instance" "default" {
		name = "${var.name}"
		description = "tf-testAccCenConfigDescription"
	}
	resource "alicloud_log_project" "default"{
		name = "${lower(var.name)}"
		description = "create by terraform"
	}
	resource "alicloud_log_store" "default"{
		project = "${alicloud_log_project.default.name}"
		name = "${lower(var.name)}"
		retention_period = 3650
		shard_count = 3
		auto_split = true
		max_split_shard_count = 60
		append_meta = true
	}
`, name)
}

// Test Cen FlowLog. >>> Resource test cases, automatically generated.
// Case attachment flowlog资源测试_副本1730702969870 8645
func TestAccAliCloudCenFlowLog_basic8645(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cen_flowlog.default"
	ra := resourceAttrInit(resourceId, AlicloudCenFlowLogMap8645)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CenServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeCenFlowLog")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%scenflowlog%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudCenFlowLogBasicDependence8645)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"status":                       "Active",
					"flow_log_name":                name,
					"description":                  "tr flowlog",
					"log_store_name":               "${alicloud_log_store.default.logstore_name}",
					"project_name":                 "${alicloud_log_store.default.project_name}",
					"log_format_string":            "$${srcaddr}$${dstaddr}$${bytes}",
					"cen_id":                       "${alicloud_cen_instance.defaultc5kxyC.id}",
					"interval":                     "60",
					"transit_router_id":            "${alicloud_cen_transit_router.defaultVw2U9u.transit_router_id}",
					"transit_router_attachment_id": "${alicloud_cen_transit_router_vpc_attachment.defaultW6LSKa.transit_router_attachment_id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"status":                       "Active",
						"flow_log_name":                name,
						"description":                  "tr flowlog",
						"log_store_name":               CHECKSET,
						"project_name":                 CHECKSET,
						"log_format_string":            "${srcaddr}${dstaddr}${bytes}",
						"cen_id":                       CHECKSET,
						"interval":                     "60",
						"transit_router_id":            CHECKSET,
						"transit_router_attachment_id": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"flow_log_name": name + "_update",
					"description":   "flowlog-resource-test-1",
					"interval":      "600",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"flow_log_name": name + "_update",
						"description":   "flowlog-resource-test-1",
						"interval":      "600",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF",
						"tags.For":     "Test",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF-update",
						"For":     "Test-update",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF-update",
						"tags.For":     "Test-update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "0",
						"tags.Created": REMOVEKEY,
						"tags.For":     REMOVEKEY,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{},
			},
		},
	})
}

var AlicloudCenFlowLogMap8645 = map[string]string{
	"create_time": CHECKSET,
	"region_id":   CHECKSET,
}

func AlicloudCenFlowLogBasicDependence8645(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

resource "alicloud_cen_instance" "defaultc5kxyC" {
  cen_instance_name = var.name
}

resource "alicloud_cen_transit_router" "defaultVw2U9u" {
  cen_id = alicloud_cen_instance.defaultc5kxyC.id
}

resource "alicloud_vpc" "defaultpybqKI" {
  cidr_block = "172.16.0.0/12"
  vpc_name   = format("%%s2", var.name)
}

resource "alicloud_vswitch" "defaultlOw4Ns" {
  vpc_id     = alicloud_vpc.defaultpybqKI.id
  cidr_block = "172.16.1.0/24"
  zone_id    = "cn-hangzhou-h"
}

resource "alicloud_log_project" "default" {
  project_name = var.name
  description  = "terraform-example"
}

resource "alicloud_log_store" "default" {
  project_name          = alicloud_log_project.default.project_name
  logstore_name         = var.name
  shard_count           = 3
  auto_split            = true
  max_split_shard_count = 60
  append_meta           = true
}

resource "alicloud_cen_transit_router_vpc_attachment" "defaultW6LSKa" {
  vpc_id = alicloud_vpc.defaultpybqKI.id
  cen_id = alicloud_cen_transit_router.defaultVw2U9u.cen_id
  zone_mappings {
    vswitch_id = alicloud_vswitch.defaultlOw4Ns.id
    zone_id    = alicloud_vswitch.defaultlOw4Ns.zone_id
  }
}


`, name)
}

// Case tr flowlog资源测试_副本1730797018057 8667
func TestAccAliCloudCenFlowLog_basic8667(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cen_flowlog.default"
	ra := resourceAttrInit(resourceId, AlicloudCenFlowLogMap8667)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CenServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeCenFlowLog")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%scenflowlog%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudCenFlowLogBasicDependence8667)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"status":            "Active",
					"flow_log_name":     name,
					"description":       "tr flowlog",
					"log_store_name":    "${alicloud_log_store.default.logstore_name}",
					"project_name":      "${alicloud_log_store.default.project_name}",
					"log_format_string": "$${srcaddr}$${dstaddr}$${bytes}",
					"cen_id":            "${alicloud_cen_instance.defaultc5kxyC.id}",
					"interval":          "60",
					"transit_router_id": "${alicloud_cen_transit_router.defaultVw2U9u.transit_router_id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"status":            "Active",
						"flow_log_name":     name,
						"description":       "tr flowlog",
						"log_store_name":    CHECKSET,
						"project_name":      CHECKSET,
						"log_format_string": "${srcaddr}${dstaddr}${bytes}",
						"cen_id":            CHECKSET,
						"interval":          "60",
						"transit_router_id": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"flow_log_name": name + "_update",
					"description":   "flowlog-resource-test-1",
					"interval":      "600",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"flow_log_name": name + "_update",
						"description":   "flowlog-resource-test-1",
						"interval":      "600",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"status": "Inactive",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"status": "Inactive",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"status": "Active",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"status": "Active",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF",
						"tags.For":     "Test",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF-update",
						"For":     "Test-update",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF-update",
						"tags.For":     "Test-update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "0",
						"tags.Created": REMOVEKEY,
						"tags.For":     REMOVEKEY,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{},
			},
		},
	})
}

var AlicloudCenFlowLogMap8667 = map[string]string{
	"create_time": CHECKSET,
	"region_id":   CHECKSET,
}

func AlicloudCenFlowLogBasicDependence8667(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

resource "alicloud_cen_instance" "defaultc5kxyC" {
  cen_instance_name = var.name
}

resource "alicloud_cen_transit_router" "defaultVw2U9u" {
  cen_id = alicloud_cen_instance.defaultc5kxyC.id
}

resource "alicloud_log_project" "default" {
  project_name = var.name
  description  = "terraform-example"
}

resource "alicloud_log_store" "default" {
  project_name          = alicloud_log_project.default.project_name
  logstore_name         = var.name
  shard_count           = 3
  auto_split            = true
  max_split_shard_count = 60
  append_meta           = true
}

`, name)
}

// Test Cen FlowLog. <<< Resource test cases, automatically generated.
