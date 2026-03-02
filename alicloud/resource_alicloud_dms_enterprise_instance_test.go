package alicloud

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"testing"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
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
	rpcClient := rawClient.(*connectivity.AliyunClient)

	prefixes := []string{
		"tf-testAcc",
		"tf_testAcc",
		"testacc",
	}
	request := map[string]interface{}{
		"InstanceState": "NORMAL",
		"PageSize":      PageSizeXLarge,
		"PageNumber":    1,
	}
	var response map[string]interface{}
	action := "ListInstances"
	for {
		response, err = rpcClient.RpcPost("dms-enterprise", "2018-11-01", action, nil, request, true)
		if err != nil {
			log.Printf("[ERROR] Failed to list DMS Enterprise Instance : %s", err)
			return nil
		}
		addDebug(action, response, request)

		resp, err := jsonpath.Get("$.InstanceList.Instance", response)
		if err != nil {
			log.Printf("[ERROR] Failed to parse DMS Enterprise Instance response: %s", err)
			return nil
		}

		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})
			id := item["Host"].(string) + ":" + item["Port"].(json.Number).String()
			skip := true

			for _, prefix := range prefixes {
				if item["InstanceAlias"] != nil &&
					strings.HasPrefix(strings.ToLower(item["InstanceAlias"].(string)), strings.ToLower(prefix)) {
					skip = false
					break
				}
			}
			if skip || item["InstanceAlias"] == nil {
				log.Printf("[INFO] Skipping DMS Enterprise Instances: %s", id)
				continue
			}
			action = "DeleteInstance"
			request = map[string]interface{}{
				"Host": item["Host"],
				"Port": item["Port"],
			}
			_, err = rpcClient.RpcPost("dms-enterprise", "2018-11-01", action, nil, request, false)
			if err != nil {
				log.Printf("[ERROR] Failed to delete DMS Enterprise Instance (%s (%s)): %s", item["InstanceAlias"].(string), id, err)
				continue
			}
			log.Printf("[INFO] Delete DMS Enterprise Instance Success: %s ", item["InstanceAlias"].(string))
		}
		if len(result) < PageSizeXLarge {
			break
		}
		request["PageNumber"] = request["PageNumber"].(int) + 1
	}
	return nil
}

func TestAccAlicloudDMSEnterprise(t *testing.T) {
	resourceId := "alicloud_dms_enterprise_instance.default"
	rac := resourceAttrCheckInit(
		resourceCheckInit(resourceId, &map[string]interface{}{},
			func() interface{} {
				return &DmsEnterpriseService{testAccProvider.Meta().(*connectivity.AliyunClient)}
			}),
		resourceAttrInit(resourceId, map[string]string{}))

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
		CheckDestroy:  resourceDMSEnterpriseCheckDestroy(rac.resourceCheck),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"dba_uid":           "${tonumber(data.alicloud_account.current.id)}",
					"host":              "${alicloud_db_instance.instance.connection_string}",
					"port":              "3306",
					"network_type":      "VPC",
					"safe_rule":         "自由操作",
					"tid":               "${data.alicloud_dms_user_tenants.default.ids.0}",
					"instance_type":     "mysql",
					"instance_source":   "RDS",
					"env_type":          "test",
					"database_user":     "${alicloud_db_account.account.name}",
					"database_password": "${alicloud_db_account.account.password}",
					"instance_name":     name,
					"query_timeout":     "70",
					"export_timeout":    "2000",
					"ecs_region":        os.Getenv("ALICLOUD_REGION"),
					"ddl_online":        "0",
					"sell_trust":        "false",
				}),
				ExpectNonEmptyPlan: true, // database_user will be encrypted if sell_trust=false
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"dba_uid":         CHECKSET,
						"host":            CHECKSET,
						"port":            "3306",
						"network_type":    "VPC",
						"safe_rule":       "自由操作",
						"tid":             CHECKSET,
						"instance_type":   "mysql",
						"instance_source": "RDS",
						"env_type":        "test",
						"database_user":   CHECKSET,
						"instance_name":   name,
						"query_timeout":   "70",
						"export_timeout":  "2000",
						"ecs_region":      os.Getenv("ALICLOUD_REGION"),
						"ddl_online":      "0",
						"use_dsql":        CHECKSET,
						"data_link_name":  CHECKSET,
						"sell_trust":      "false",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"database_password", "dba_uid", "network_type", "port", "safe_rule", "tid", "database_user"},
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"env_type":   "dev",
					"sell_trust": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"env_type":   "dev",
						"sell_trust": "true",
						// database_user will be plain text if sell_trust=true
						"database_user": CHECKSET,
					}),
				),
			},
			// instance_name cannot update
			//{
			//	Config: testAccConfig(map[string]interface{}{
			//		"instance_name": name + "update",
			//	}),
			//	Check: resource.ComposeTestCheckFunc(
			//		testAccCheck(map[string]string{
			//			"instance_name": name + "update",
			//		}),
			//	),
			//},
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
					"tid":               "${data.alicloud_dms_user_tenants.default.ids.0}",
					"instance_type":     "mysql",
					"instance_source":   "RDS",
					"env_type":          "test",
					"database_user":     "${alicloud_db_account.account.name}",
					"database_password": "${alicloud_db_account.account.password}",
					"instance_name":     name,
					"query_timeout":     "70",
					"export_timeout":    "2000",
					"ecs_region":        os.Getenv("ALICLOUD_REGION"),
					"ddl_online":        "0",
					"use_dsql":          "0",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"dba_uid":         CHECKSET,
						"host":            CHECKSET,
						"port":            "3306",
						"network_type":    "VPC",
						"safe_rule":       "自由操作",
						"tid":             CHECKSET,
						"instance_type":   "mysql",
						"instance_source": "RDS",
						"env_type":        "test",
						"database_user":   CHECKSET,
						"instance_name":   name,
						"query_timeout":   "70",
						"export_timeout":  "2000",
						"ecs_region":      os.Getenv("ALICLOUD_REGION"),
						"ddl_online":      "0",
						"use_dsql":        "0",
					}),
				),
			},
		},
	})
}

func resourceDmsConfigDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}
data "alicloud_account" "current" {
}
data "alicloud_dms_user_tenants" "default" {
  status = "ACTIVE"
}

data "alicloud_db_zones" "default" {
  engine                   = "MySQL"
  engine_version           = "8.0"
  instance_charge_type     = "PostPaid"
  category                 = "HighAvailability"
  db_instance_storage_type = "cloud_essd"
}

data "alicloud_db_instance_classes" "default" {
  zone_id                  = data.alicloud_db_zones.default.zones.0.id
  engine                   = "MySQL"
  engine_version           = "8.0"
  category                 = "HighAvailability"
  db_instance_storage_type = "cloud_essd"
  instance_charge_type     = "PostPaid"
}

data "alicloud_vpcs" "default" {
  name_regex = "^default-NODELETING"
}
data "alicloud_vswitches" "default" {
  vpc_id  = data.alicloud_vpcs.default.ids.0
  zone_id = data.alicloud_db_zones.default.zones.0.id
}

resource "alicloud_security_group" "default" {
  name   = var.name
  vpc_id = data.alicloud_vpcs.default.ids.0
}

resource "alicloud_db_instance" "instance" {
  engine                   = "MySQL"
  engine_version           = "8.0"
  db_instance_storage_type = "cloud_essd"
  instance_type            = data.alicloud_db_instance_classes.default.instance_classes.0.instance_class
  instance_storage         = data.alicloud_db_instance_classes.default.instance_classes.0.storage_range.min
  vswitch_id               = data.alicloud_vswitches.default.ids.0
  instance_name            = var.name
  security_ips             = ["100.104.5.0/24", "192.168.0.6"]
  tags = {
    "key1" = "value1"
    "key2" = "value2"
  }
}

resource "alicloud_db_account" "account" {
  instance_id = "${alicloud_db_instance.instance.id}"
  name        = "tftestnormal"
  password    = "Test12345"
  type        = "Normal"
}`, name)
}

func resourceDMSEnterpriseCheckDestroy(rc *resourceCheck) func(*terraform.State) error {
	return func(s *terraform.State) error {
		for _, rs := range s.RootModule().Resources {
			if rs.Type != "alicloud_dms_enterprise_instance" {
				continue
			}

			aliClient := testAccProvider.Meta().(*connectivity.AliyunClient)
			dmsEnterpriseService := DmsEnterpriseService{aliClient}
			resp, err := dmsEnterpriseService.DescribeDmsEnterpriseInstance(rs.Primary.ID)
			if err != nil {
				if NotFoundError(err) {
					continue
				}
				return WrapError(err)
			}

			// DMS enterprise resource deletion check depend on the response status, not the HTTP code 404
			if v, ok := resp["State"]; ok && v == "DELETED" {
				log.Printf(" Resource (alicloud_dms_enterprise_instance | %s) has been deleted.", rs.Primary.ID)
				continue
			}
			return WrapError(Error("the resource %s %s was not destroyed ! ", rc.resourceId, rs.Primary.ID))
		}
		return nil
	}
}
