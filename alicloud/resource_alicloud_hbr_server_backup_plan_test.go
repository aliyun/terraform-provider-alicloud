package alicloud

import (
	"fmt"
	"log"
	"strings"
	"testing"
	"time"

	"github.com/PaesslerAG/jsonpath"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func init() {
	resource.AddTestSweepers(
		"alicloud_hbr_server_backup_plan",
		&resource.Sweeper{
			Name: "alicloud_hbr_server_backup_plan",
			F:    testSweepHbrServerBackupPlan,
		})
}

func testSweepHbrServerBackupPlan(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting Alicloud client: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)
	prefixes := []string{
		"tf-testAcc",
		"tf_testAcc",
	}
	action := "DescribeBackupPlans"
	request := map[string]interface{}{
		"SourceType": "UDM_ECS",
		"PageNumber": 1,
		"PageSize":   PageSizeLarge,
	}

	var hbrPlans []interface{}

	conn, err := client.NewHbrClient()

	if err != nil {
		log.Printf("[ERROR] %s get an error: %#v", action, err)
		return nil
	}

	for {
		var response map[string]interface{}

		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2017-09-08"), StringPointer("AK"), nil, request, &runtime)
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})
		addDebug(action, response, request)
		if err != nil {
			log.Printf("[ERROR] %s get an error: %#v", action, err)
			return nil
		}
		resp, err := jsonpath.Get("$.BackupPlans.BackupPlan", response)
		if err != nil {
			log.Printf("[ERROR] Getting resource %s attribute by path %s failed!!! Body: %v.", "$.BackupPlans.BackupPlan", action, err)
			return nil
		}
		result, _ := resp.([]interface{})

		if len(result) < 1 {
			break
		}
		hbrPlans = append(hbrPlans, result...)
		if len(result) < PageSizeLarge {
			break
		}

		if page, err := getNextpageNumber(request["PageNumber"].(requests.Integer)); err != nil {
			log.Printf("[ERROR] %s get an error: %#v", "getNextpageNumber", err)
			break
		} else {
			request["PageNumber"] = page
		}
	}

	sweeped := false
	for _, v := range hbrPlans {
		item := v.(map[string]interface{})
		skip := true
		for _, prefix := range prefixes {
			if strings.HasPrefix(strings.ToLower(item["PlanName"].(string)), strings.ToLower(prefix)) {
				skip = false
			}
		}
		if skip {
			log.Printf("[INFO] Skipping HBR Ecs Server Backup Plan : %s", item["PlanName"].(string))
			continue
		}

		sweeped = true
		action := "DeleteBackupPlan"
		request := map[string]interface{}{
			"PlanId": item["PlanId"],
		}
		_, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2017-09-08"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
		if err != nil {
			log.Printf("[ERROR] Failed to delete HBR Ecs Server Backup Plan (%s): %s", item["PlanName"].(string), err)
		}

		if sweeped {
			// Waiting 5 seconds to ensure Direct Mail Domain have been deleted.
			time.Sleep(5 * time.Second)
		}
		log.Printf("[INFO] Delete HBR Ecs Server Backup Plan success: %s ", item["PlanName"].(string))
	}
	return nil
}

func TestAccAlicloudHBRServerBackupPlan_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_hbr_server_backup_plan.default"
	checkoutSupportedRegions(t, true, connectivity.HbrSupportRegions)
	ra := resourceAttrInit(resourceId, AlicloudHBRServerBackupPlanMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &HbrService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeHbrServerBackupPlan")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAcc%shbrecsbackupplan%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudHBRServerBackupPlanBasicDependence0)
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
					"instance_id":                 "${data.alicloud_instances.default.instances.0.id}",
					"schedule":                    "I|1602673264|PT2H",
					"ecs_server_backup_plan_name": "tf-testAcc-hbr-backup-plan",
					"retention":                   "1",
					"detail": []map[string]interface{}{
						{
							"app_consistent": "false",
							"snapshot_group": "true",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"schedule":                    "I|1602673264|PT2H",
						"ecs_server_backup_plan_name": "tf-testAcc-hbr-backup-plan",
						"retention":                   "1",
						"detail.#":                    "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"detail": []map[string]interface{}{
						{
							"app_consistent":     "true",
							"snapshot_group":     "true",
							"enable_fs_freeze":   "true",
							"pre_script_path":    "",
							"post_script_path":   "",
							"timeout_in_seconds": "60",
							"disk_id_list":       []string{},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"detail.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"detail": []map[string]interface{}{
						{
							"app_consistent":     "true",
							"snapshot_group":     "false",
							"enable_fs_freeze":   "true",
							"pre_script_path":    "",
							"post_script_path":   "",
							"timeout_in_seconds": "60",
							"disk_id_list":       []string{},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"detail.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"detail": []map[string]interface{}{
						{
							"app_consistent":     "true",
							"snapshot_group":     "false",
							"enable_fs_freeze":   "false",
							"pre_script_path":    "",
							"post_script_path":   "",
							"timeout_in_seconds": "60",
							"disk_id_list":       []string{},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"detail.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"detail": []map[string]interface{}{
						{
							"app_consistent":     "true",
							"snapshot_group":     "false",
							"enable_fs_freeze":   "false",
							"pre_script_path":    "./example.js",
							"post_script_path":   "",
							"timeout_in_seconds": "60",
							"disk_id_list":       []string{},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"detail.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"detail": []map[string]interface{}{
						{
							"app_consistent":     "true",
							"snapshot_group":     "false",
							"enable_fs_freeze":   "false",
							"pre_script_path":    "./example.js",
							"post_script_path":   "./example.js",
							"timeout_in_seconds": "60",
							"disk_id_list":       []string{},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"detail.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"detail": []map[string]interface{}{
						{
							"app_consistent":     "true",
							"snapshot_group":     "false",
							"enable_fs_freeze":   "false",
							"pre_script_path":    "./example.js",
							"post_script_path":   "./example.js",
							"timeout_in_seconds": "180",
							"disk_id_list":       []string{},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"detail.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"detail": []map[string]interface{}{
						{
							"app_consistent":     "true",
							"snapshot_group":     "false",
							"enable_fs_freeze":   "false",
							"pre_script_path":    "./example.js",
							"post_script_path":   "./example.js",
							"timeout_in_seconds": "180",
							"disk_id_list":       []string{"/home"},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"detail.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"detail": []map[string]interface{}{
						{
							"app_consistent":     "true",
							"snapshot_group":     "false",
							"enable_fs_freeze":   "false",
							"pre_script_path":    "./example.js",
							"post_script_path":   "./example.js",
							"timeout_in_seconds": "180",
							"disk_id_list":       []string{"/home", "/var"},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"detail.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"detail": []map[string]interface{}{
						{
							"app_consistent":     "true",
							"snapshot_group":     "false",
							"enable_fs_freeze":   "false",
							"pre_script_path":    "./example.js",
							"post_script_path":   "./example.js",
							"timeout_in_seconds": "180",
							"disk_id_list":       []string{},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"detail.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"detail": []map[string]interface{}{
						{
							"app_consistent":        "true",
							"snapshot_group":        "false",
							"enable_fs_freeze":      "false",
							"pre_script_path":       "./example.js",
							"post_script_path":      "./example.js",
							"timeout_in_seconds":    "180",
							"disk_id_list":          []string{},
							"do_copy":               "false",
							"destination_region_id": "cn-hangzhou",
							"destination_retention": "1",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"detail.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"detail": []map[string]interface{}{
						{
							"app_consistent":        "true",
							"snapshot_group":        "false",
							"enable_fs_freeze":      "false",
							"pre_script_path":       "./example.js",
							"post_script_path":      "./example.js",
							"timeout_in_seconds":    "180",
							"disk_id_list":          []string{},
							"do_copy":               "true",
							"destination_region_id": "cn-shanghai",
							"destination_retention": "2",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"detail.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"detail": []map[string]interface{}{
						{
							"app_consistent":        "true",
							"snapshot_group":        "false",
							"enable_fs_freeze":      "false",
							"pre_script_path":       "./example.js",
							"post_script_path":      "./example.js",
							"timeout_in_seconds":    "180",
							"disk_id_list":          []string{},
							"do_copy":               "false",
							"destination_region_id": "cn-beijing",
							"destination_retention": "1",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"detail.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"detail": []map[string]interface{}{
						{
							"app_consistent":        "true",
							"snapshot_group":        "false",
							"enable_fs_freeze":      "false",
							"pre_script_path":       "./example.js",
							"post_script_path":      "./example.js",
							"timeout_in_seconds":    "180",
							"disk_id_list":          []string{},
							"do_copy":               "false",
							"destination_retention": "1",
							"destination_region_id": "cn-hangzhou",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"detail.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"retention": "2",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"retention": "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"ecs_server_backup_plan_name": "tf-testAcc-hbr-backup-plan2",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"ecs_server_backup_plan_name": "tf-testAcc-hbr-backup-plan2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"disabled": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"disabled": "true",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"disabled": "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"disabled": "false",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"schedule": "I|1602673264|P1D",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"schedule": "I|1602673264|P1D",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"disabled": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"disabled": "true",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"ecs_server_backup_plan_name": "tf-testAcc-hbr-backup-plan3",
					"schedule":                    "I|1602673264|PT2H",
					"retention":                   "4",
					"disabled":                    "true",
					"detail": []map[string]interface{}{
						{
							"app_consistent":        "true",
							"snapshot_group":        "false",
							"enable_fs_freeze":      "false",
							"pre_script_path":       "./example.js",
							"post_script_path":      "./example.js",
							"timeout_in_seconds":    "180",
							"disk_id_list":          []string{},
							"do_copy":               "false",
							"destination_retention": "1",
							"destination_region_id": "cn-hangzhou",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"ecs_server_backup_plan_name": "tf-testAcc-hbr-backup-plan3",
						"schedule":                    "I|1602673264|PT2H",
						"retention":                   "4",
						"disabled":                    "true",
						"detail.#":                    "1",
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

var AlicloudHBRServerBackupPlanMap0 = map[string]string{
	"instance_id":                 CHECKSET,
	"schedule":                    CHECKSET,
	"ecs_server_backup_plan_name": CHECKSET,
}

func AlicloudHBRServerBackupPlanBasicDependence0(name string) string {
	return fmt.Sprintf(` 
variable "name" {
  default = "%s"
}

data "alicloud_instances" "default" {
  name_regex = "no-deleteing-hbr-ecs-server-backup-plan"
  status = "Running"
}
`, name)
}
