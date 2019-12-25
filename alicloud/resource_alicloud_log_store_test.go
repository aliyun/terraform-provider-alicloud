package alicloud

import (
	"fmt"
	"testing"

	sls "github.com/aliyun/aliyun-log-go-sdk"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
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
					"name":        name,
					"project":     alicloud_log_project.foo.name,
					"shard_count": "1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":        name,
						"project":     name,
						"shard_count": "1",
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
					"retention_period": "3000",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"retention_period": "3000",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"auto_split": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"auto_split": "true",
					}),
				),
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
					"project": alicloud_log_project.foo.name,
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
	    name = var.name
	    description = "tf unit test"
	}
	`, name)
}

var logStoreMap = map[string]string{
	"name":                  CHECKSET,
	"project":               CHECKSET,
	"retention_period":      "30",
	"shard_count":           CHECKSET,
	"shards.#":              CHECKSET,
	"auto_split":            "false",
	"max_split_shard_count": "0",
	"append_meta":           "true",
	"enable_web_tracking":   "false",
}
