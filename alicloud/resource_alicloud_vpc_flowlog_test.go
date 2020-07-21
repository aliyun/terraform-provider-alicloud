package alicloud

import (
	"fmt"
	"log"
	"strings"
	"testing"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/vpc"

	"github.com/hashicorp/terraform/helper/acctest"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func init() {
	resource.AddTestSweepers("alicloud_vpc_flowlog", &resource.Sweeper{
		Name: "alicloud_vpc_flowlog",
		F:    testSweepVpcFlowlog,
	})
}

func testSweepVpcFlowlog(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting Alicloud client: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)

	prefixes := []string{
		fmt.Sprintf("tf-testAcc%s", region),
		fmt.Sprintf("tf_testAcc%s", region),
	}

	var insts []vpc.FlowLog
	request := vpc.CreateDescribeFlowLogsRequest()
	request.RegionId = client.RegionId
	request.PageSize = requests.NewInteger(PageSizeLarge)
	request.PageNumber = requests.NewInteger(1)

	for {
		var raw interface{}
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			raw, err = client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
				return vpcClient.DescribeFlowLogs(request)
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
			log.Printf("[ERROR] Failed to retrieving Vpc flowlog: %s", err)
			break
		}
		response, _ := raw.(*vpc.DescribeFlowLogsResponse)
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
		for _, prefix := range prefixes {
			if strings.HasPrefix(strings.ToLower(name), strings.ToLower(prefix)) {
				skip = false
				break
			}
		}
		if skip {
			log.Printf("[INFO] Skipping Vpc flowlog: %s (%s)", name, id)
			continue
		}
		sweeped = true
		log.Printf("[INFO] Deleting Vpc flowlog: %s (%s)", name, id)
		request := vpc.CreateDeleteFlowLogRequest()
		request.FlowLogId = id
		_, err := client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
			return vpcClient.DeleteFlowLog(request)
		})
		if err != nil {
			log.Printf("[ERROR] Failed to delete Vpc flowlog (%s (%s)): %s", name, id, err)
		}
	}
	if sweeped {
		// Waiting 5 seconds to eusure these flowlog have been deleted.
		time.Sleep(5 * time.Second)
	}
	return nil
}

func TestAccAlicloudVpcFlowlog_basic(t *testing.T) {
	var vpcFlowlog vpc.FlowLog

	resourceId := "alicloud_vpc_flowlog.default"
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testAcc%sFlowlog-%d", defaultRegionToTest, rand)
	nameUpdated := fmt.Sprintf("tf-testAcc%sFlowlog-%d-Update", defaultRegionToTest, rand)

	ra := resourceAttrInit(resourceId, map[string]string{
		"flow_log_name": name,
		"description":   name,
		"status":        "Inactive",
	})

	serviceFunc := func() interface{} {
		return &VpcService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &vpcFlowlog, serviceFunc)

	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceVpcFlowlogConfigDependence)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, connectivity.VpcFlowLogNoSkipRegions)
		},
		// module name
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"resource_id":    "${alicloud_vpc.default.id}",
					"resource_type":  "VPC",
					"traffic_type":   "All",
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

// One Vpc Instance only can have one flowlog.

func resourceVpcFlowlogConfigDependence(name string) string {
	return fmt.Sprintf(`
	variable "name" {
		default = "%s"
	}
	resource "alicloud_vpc" "default" {
  		cidr_block = "192.168.0.0/24"
  		name = "${var.name}"
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
