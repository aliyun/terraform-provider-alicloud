package alicloud

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/PaesslerAG/jsonpath"
	util "github.com/alibabacloud-go/tea-utils/service"
	log "github.com/sirupsen/logrus"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func init() {
	resource.AddTestSweepers("alicloud_hbr_ots_backup_plan", &resource.Sweeper{
		Name: "alicloud_hbr_ots_backup_plan",
		F:    testSweepHBROtsBackupPlan,
	})
}

func testSweepHBROtsBackupPlan(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting Alicloud client: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)
	prefixes := []string{
		"testAcc",
		"tf-testAcc",
		"tf_testAcc",
	}
	action := "DescribeBackupPlans"
	request := make(map[string]interface{})
	request["SourceType"] = "OTS"

	request["PageSize"] = PageSizeLarge
	request["PageNumber"] = 1

	var response map[string]interface{}
	conn, err := client.NewHbrClient()
	if err != nil {
		log.Printf("[ERROR] %s get an error: %#v", "NewHbrClient", err)
		return nil
	}
	for {
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(1*time.Minute, func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2017-09-08"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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
		for _, v := range result {
			item := v.(map[string]interface{})

			if _, ok := item["PlanName"]; !ok {
				continue
			}
			skip := true
			for _, prefix := range prefixes {
				if strings.HasPrefix(strings.ToLower(item["PlanName"].(string)), strings.ToLower(prefix)) {
					skip = false
				}
			}
			if skip {
				log.Printf("[INFO] Skipping Hbr Ots Backup Plan: %s", item["PlanName"].(string))
				continue
			}
			action := "DeleteBackupPlan"
			request := map[string]interface{}{
				"PlanId":     item["PlanId"],
				"SourceType": "OTS",
				"VaultId":    item["VaultId"],
			}
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2017-09-08"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
			if err != nil {
				log.Printf("[ERROR] Failed to delete HBR Ots BackUp Plan (%s): %s", item["ProductId"], err)
			}
			log.Printf("[INFO] Delete HBR Ots BackUp Plan success: %s ", item["ProductId"])
		}
		if len(result) < PageSizeLarge {
			break
		}
		request["PageNumber"] = request["PageNumber"].(int) + 1
	}
	return nil
}

func TestAccAlicloudHBROtsBackupPlan_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_hbr_ots_backup_plan.default"
	ra := resourceAttrInit(resourceId, AlicloudHBROtsBackupPlanMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &HbrService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeHbrOtsBackupPlan")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1, 9999)
	name := fmt.Sprintf("testAcc%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudHBROtsBackupPlanBasicDependence0)
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
					"vault_id":             "${alicloud_hbr_vault.default.id}",
					"ots_backup_plan_name": name,
					"backup_type":          "COMPLETE",
					"schedule":             "I|1602673264|PT2H",
					"retention":            "1",
					"instance_name":        "${alicloud_ots_instance.foo.name}",
					"ots_detail": []map[string]interface{}{
						{
							"table_names": []string{
								"${alicloud_ots_table.basic.table_name}",
							},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"backup_type":          "COMPLETE",
						"schedule":             "I|1602673264|PT2H",
						"ots_backup_plan_name": name,
						"retention":            "1",
						"ots_detail.#":         "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"ots_backup_plan_name": name + "update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"ots_backup_plan_name": name + "update",
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
					"retention": "2",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"retention": "2",
					}),
				),
			},

			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"update_paths"},
			},
		},
	})
}

func TestAccAlicloudHBROtsBackupPlan_basic1(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_hbr_ots_backup_plan.default"
	ra := resourceAttrInit(resourceId, AlicloudHBROtsBackupPlanMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &HbrService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeHbrOtsBackupPlan")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1, 9999)
	name := fmt.Sprintf("testAcc%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudHBROtsBackupPlanBasicDependence0)
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
					"vault_id":             "${alicloud_hbr_vault.default.id}",
					"ots_backup_plan_name": name,
					"backup_type":          "COMPLETE",
					"retention":            "1",
					"instance_name":        "${alicloud_ots_instance.foo.name}",
					"ots_detail": []map[string]interface{}{
						{
							"table_names": []string{
								"${alicloud_ots_table.basic.table_name}",
							},
						},
					},
					"rules": []map[string]interface{}{
						{
							"schedule":    fmt.Sprintf("I|%v|P1D", time.Now().Unix()),
							"retention":   "1",
							"disabled":    "false",
							"rule_name":   name,
							"backup_type": "COMPLETE",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"backup_type":          "COMPLETE",
						"rules.#":              "1",
						"ots_backup_plan_name": name,
						"retention":            "1",
						"ots_detail.#":         "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"rules": []map[string]interface{}{
						{
							"schedule":    fmt.Sprintf("I|%v|P1D", time.Now().Unix()),
							"retention":   "1",
							"disabled":    "false",
							"rule_name":   name,
							"backup_type": "COMPLETE",
						},
						{
							"schedule":    fmt.Sprintf("I|%v|P1D", time.Now().Unix()),
							"retention":   "1",
							"disabled":    "false",
							"rule_name":   name + "update",
							"backup_type": "COMPLETE",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"backup_type":          "COMPLETE",
						"rules.#":              "2",
						"ots_backup_plan_name": name,
						"retention":            "1",
						"ots_detail.#":         "1",
					}),
				),
			},
		},
	})
}

var AlicloudHBROtsBackupPlanMap0 = map[string]string{}

func AlicloudHBROtsBackupPlanBasicDependence0(name string) string {
	return fmt.Sprintf(` 
variable "name" {
  default = "%s"
}
resource "alicloud_hbr_vault" "default" {
  vault_name = var.name
  vault_type = "OTS_BACKUP"
}

resource "alicloud_ots_instance" "foo" {
  name        = var.name
  description = var.name
  accessed_by = "Any"
  tags = {
    Created = "TF"
    For     = "acceptance test"
  }
}

resource "alicloud_ots_table" "basic" {
  instance_name = alicloud_ots_instance.foo.name
  table_name    = var.name
  primary_key {
    name = "pk1"
    type = "Integer"
  }
  time_to_live                  = -1
  max_version                   = 1
  deviation_cell_version_in_sec = 1
}


`, name)
}
