package alicloud

import (
	"fmt"
	"os"
	"testing"

	sls "github.com/aliyun/aliyun-log-go-sdk"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAlicloudLogStore_basic(t *testing.T) {
	var v *sls.LogStore
	resourceId := "alicloud_log_store.default"
	ra := resourceAttrInit(resourceId, logStoreMap)
	serviceFunc := func() interface{} {
		return &LogService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testacc-log-store-%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceLogStoreConfigDependence)

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
					"name":                  name,
					"project":               "${alicloud_log_project.foo.name}",
					"shard_count":           "1",
					"auto_split":            "true",
					"max_split_shard_count": "1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":                  name,
						"project":               name,
						"shard_count":           "1",
						"auto_split":            "true",
						"max_split_shard_count": "1",
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"max_split_shard_count": "60",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"max_split_shard_count": "60",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"retention_period": "3000",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"retention_period": "3000",
					}),
				),
			},
			// TODO: because auto_split and max_split_shard_count affect each other, when auto_split = false, max_split_shard_count will be set to 0, and when updating auto_split = true, max_split_shard_count must be set to be greater than 0, so in the test, auto_split = true in step 0, omitting this step
			// {
			// 	Config: testAccConfig(map[string]interface{}{
			// 		"auto_split": "true",
			// 	}),
			// 	Check: resource.ComposeTestCheckFunc(
			// 		testAccCheck(map[string]string{
			// 			"auto_split": "true",
			// 		}),
			// 	),
			// },
			{
				Config: testAccConfig(map[string]interface{}{
					"append_meta": "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"append_meta": "false",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"enable_web_tracking": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"enable_web_tracking": "true",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"retention_period":      REMOVEKEY,
					"auto_split":            REMOVEKEY,
					"max_split_shard_count": REMOVEKEY,
					"append_meta":           REMOVEKEY,
					"enable_web_tracking":   REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"retention_period":      "30",
						"auto_split":            "false",
						"max_split_shard_count": "0",
						"append_meta":           "true",
						"enable_web_tracking":   "false",
					}),
				),
			},
		},
	})
}

func TestAccAlicloudLogStore_create_with_encrypt(t *testing.T) {
	var v *sls.LogStore
	resourceId := "alicloud_log_store.default"
	ra := resourceAttrInit(resourceId, logStoreMap)
	serviceFunc := func() interface{} {
		return &LogService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testacc-log-store-%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceLogStoreConfigDependenceWithEncrypt)

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
					"name":        name,
					"project":     "${alicloud_log_project.foo.name}",
					"shard_count": "1",
					"encrypt_conf": []map[string]interface{}{
						{
							"enable":       "true",
							"encrypt_type": "default",
							"user_cmk_info": []map[string]string{
								{
									"cmk_key_id": "${alicloud_kms_key.key.id}",
									"arn":        "acs:ram::${data.alicloud_account.default.id}:role/aliyunlogdefaultrole",
									"region_id":  os.Getenv("ALICLOUD_REGION"),
								},
							},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":           name,
						"project":        name,
						"shard_count":    "1",
						"encrypt_conf.#": "1",
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAlicloudLogStore_multi(t *testing.T) {
	var v *sls.LogStore
	resourceId := "alicloud_log_store.default.4"
	ra := resourceAttrInit(resourceId, logStoreMap)
	serviceFunc := func() interface{} {
		return &LogService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testacc-log-store-%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceLogStoreConfigDependence)

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
					"name":    name + "${count.index}",
					"project": "${alicloud_log_project.foo.name}",
					"count":   "5",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(nil),
				),
			},
		},
	})
}

func resourceLogStoreConfigDependence(name string) string {
	return fmt.Sprintf(`
	variable "name" {
	    default = "%s"
	}
	resource "alicloud_log_project" "foo" {
	    name = "${var.name}"
	    description = "tf unit test"
	}
	`, name)
}

func resourceLogStoreConfigDependenceWithEncrypt(name string) string {
	return fmt.Sprintf(`
	variable "name" {
	    default = "%s"
	}
	data "alicloud_account" "default"{
	}
	resource "alicloud_kms_key" "key" {
  		description             = "${var.name}"
  		pending_window_in_days  = "7"
  		status                  = "Enabled"
	}
	resource "alicloud_log_project" "foo" {
	    name = "${var.name}"
	    description = "tf unit test"
	}
	`, name)
}

var logStoreMap = map[string]string{
	"name":    CHECKSET,
	"project": CHECKSET,
}
