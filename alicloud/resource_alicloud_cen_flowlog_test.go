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
