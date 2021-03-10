package alicloud

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/PaesslerAG/jsonpath"
	util "github.com/alibabacloud-go/tea-utils/service"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func init() {
	resource.AddTestSweepers("alicloud_dms_enterprise_instance", &resource.Sweeper{
		Name: "alicloud_dms_enterprise_instance",
		F:    testSweepDMSEnterpriseInstances,
	})
}

func testSweepDMSEnterpriseInstances(region string) error {

	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return WrapErrorf(err, "Error getting Alicloud client.")
	}
	client := rawClient.(*connectivity.AliyunClient)

	prefixes := []string{
		"tf-testAcc",
		"tf_testAcc",
	}
	request := make(map[string]interface{})
	var response map[string]interface{}
	action := "ListInstances"
	conn, err := client.NewDmsenterpriseClient()
	if err != nil {
		return WrapError(err)
	}
	request["PageSize"] = PageSizeLarge
	request["PageNumber"] = 1
	for {
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2018-11-01"), StringPointer("AK"), nil, request, &runtime)
		if err != nil {
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_dms_enterprise_instances", action, AlibabaCloudSdkGoERROR)
		}
		addDebug(action, response, request)

		resp, err := jsonpath.Get("$.InstanceList.Instance", response)
		if err != nil {
			return WrapErrorf(err, FailedGetAttributeMsg, action, "$.InstanceList.Instance", response)
		}

		sweeped := false
		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})
			id := item["Host"].(string) + ":" + item["Port"].(json.Number).String()

			skip := true

			for _, prefix := range prefixes {
				if item["InstanceAlias"] != nil {
					if strings.HasPrefix(strings.ToLower(item["InstanceAlias"].(string)), strings.ToLower(prefix)) {
						skip = false
						break
					}
				}
			}
			if skip || item["InstanceAlias"] == nil {
				log.Printf("[INFO] Skipping DMS Enterprise Instances: %s", id)
				continue
			}
			sweeped = true
			action = "DeleteInstance"
			request := map[string]interface{}{
				"Host": item["Host"],
				"Port": item["Port"],
			}
			_, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2018-11-01"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
			if err != nil {
				log.Printf("[ERROR] Failed to delete DMS Enterprise Instance (%s (%s)): %s", item["InstanceAlias"].(string), id, err)
				continue
			}
			if sweeped {
				// Waiting 30 seconds to ensure these DMS Enterprise Instance have been deleted.
				time.Sleep(30 * time.Second)
			}
			log.Printf("[INFO] Delete DMS Enterprise Instance Success: %s ", item["InstanceAlias"].(string))
		}
		if len(result) < PageSizeLarge {
			break
		}
		request["PageNumber"] = request["PageNumber"].(int) + 1
	}
	return nil
}

func TestAccAlicloudDmsEnterprise(t *testing.T) {
	resourceId := "alicloud_dms_enterprise_instance.default"
	var v map[string]interface{}
	ra := resourceAttrInit(resourceId, testAccCheckKeyValueInMapsForDMS)

	serviceFunc := func() interface{} {
		return &Dms_enterpriseService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)

	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000, 9999)
	name := fmt.Sprintf("tf-testAccDmsEnterpriseInstance%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceDmsConfigDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"dba_uid":           "${tonumber(data.alicloud_account.current.id)}",
					"host":              "${alicloud_db_instance.instance.connection_string}",
					"port":              "3306",
					"network_type":      "VPC",
					"safe_rule":         "自由操作",
					"tid":               "13429",
					"instance_type":     "mysql",
					"instance_source":   "RDS",
					"env_type":          "test",
					"database_user":     "${alicloud_db_account.account.name}",
					"database_password": "${alicloud_db_account.account.password}",
					"instance_alias":    name,
					"query_timeout":     "70",
					"export_timeout":    "2000",
					"ecs_region":        os.Getenv("ALICLOUD_REGION"),
					"ddl_online":        "0",
					"use_dsql":          "0",
					"data_link_name":    "",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"dba_uid":         CHECKSET,
						"host":            CHECKSET,
						"port":            "3306",
						"network_type":    "VPC",
						"safe_rule":       "自由操作",
						"tid":             "13429",
						"instance_type":   "mysql",
						"instance_source": "RDS",
						"env_type":        "test",
						"database_user":   CHECKSET,
						"instance_alias":  name,
						"query_timeout":   "70",
						"export_timeout":  "2000",
						"ecs_region":      os.Getenv("ALICLOUD_REGION"),
						"ddl_online":      "0",
						"use_dsql":        "0",
						"data_link_name":  "",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"database_password", "dba_uid", "network_type", "port", "safe_rule", "tid"},
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"env_type": "dev",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"env_type": "dev",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_alias": "other_name",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_alias": "other_name",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"query_timeout": "77",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"query_timeout": "77",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{
					"dba_uid":           "${tonumber(data.alicloud_account.current.id)}",
					"host":              "${alicloud_db_instance.instance.connection_string}",
					"port":              "3306",
					"network_type":      "VPC",
					"safe_rule":         "自由操作",
					"tid":               "13429",
					"instance_type":     "mysql",
					"instance_source":   "RDS",
					"env_type":          "test",
					"database_user":     "${alicloud_db_account.account.name}",
					"database_password": "${alicloud_db_account.account.password}",
					"instance_alias":    name,
					"query_timeout":     "70",
					"export_timeout":    "2000",
					"ecs_region":        os.Getenv("ALICLOUD_REGION"),
					"ddl_online":        "0",
					"use_dsql":          "0",
					"data_link_name":    "",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"dba_uid":         CHECKSET,
						"host":            CHECKSET,
						"port":            "3306",
						"network_type":    "VPC",
						"safe_rule":       "自由操作",
						"tid":             "13429",
						"instance_type":   "mysql",
						"instance_source": "RDS",
						"env_type":        "test",
						"database_user":   CHECKSET,
						"instance_alias":  name,
						"query_timeout":   "70",
						"export_timeout":  "2000",
						"ecs_region":      os.Getenv("ALICLOUD_REGION"),
						"ddl_online":      "0",
						"use_dsql":        "0",
						"data_link_name":  "",
					}),
				),
			},
		},
	})
}

var testAccCheckKeyValueInMapsForDMS = map[string]string{}

func resourceDmsConfigDependence(name string) string {
	return fmt.Sprintf(`
	data "alicloud_account" "current" {
	}
	
	data "alicloud_vpcs" "default" {
	is_default = true
	}
	data "alicloud_vswitches" "default" {
	ids = [
	  data.alicloud_vpcs.default.vpcs.0.vswitch_ids.0]
	}
	
	resource "alicloud_security_group" "default" {
	name = "%[1]s"
	vpc_id = "${data.alicloud_vpcs.default.ids.0}"
	}
	
	resource "alicloud_db_instance" "instance" {
	engine           = "MySQL"
	engine_version   = "5.7"
	instance_type    = "rds.mysql.t1.small"
	instance_storage = "10"
	vswitch_id       = "${data.alicloud_vswitches.default.ids.0}"
	instance_name    = "%[1]s"
	security_ips     = ["100.104.5.0/24","192.168.0.6"]
	}
	
	resource "alicloud_db_account" "account" {
	instance_id = "${alicloud_db_instance.instance.id}"
	name        = "tftestnormal"
	password    = "Test12345"
	type        = "Normal"
	}`, name)
}
