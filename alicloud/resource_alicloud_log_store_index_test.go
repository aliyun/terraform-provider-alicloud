package alicloud

import (
	"fmt"
	"testing"

	sls "github.com/aliyun/aliyun-log-go-sdk"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func TestAccAlicloudLogStoreIndex_basic(t *testing.T) {
	var v *sls.Index
	resourceId := "alicloud_log_store_index.default"
	ra := resourceAttrInit(resourceId, logStoreIndexMap)
	serviceFunc := func() interface{} {
		return &LogService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testacclogstoreindex-%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceLogStoreIndexConfigDependence)

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
					"project":  alicloud_log_project.default.name,
					"logstore": alicloud_log_store.default.name,
					"full_text": []map[string]interface{}{
						{
							"case_sensitive": "true",
							"token":          ` #$^*\r\n\t`,
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"project":     name,
						"logstore":    name,
						"full_text.#": "1",
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
					"full_text": REMOVEKEY,
					"field_search": []map[string]interface{}{
						{
							"name":             var.name,
							"enable_analytics": "true",
							"token":            ` #$^*\r\n\t`,
							"type":             "text",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"full_text.#":    REMOVEKEY,
						"field_search.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{

					"full_text": []map[string]interface{}{
						{
							"case_sensitive": "true",
							"token":          ` #$^*\r\n\t`,
						},
					},
					"field_search": []map[string]interface{}{
						{
							"name":             "${var.name}-1",
							"enable_analytics": "true",
							"token":            ` #$^*\r\n\t`,
							"type":             "json",
							"json_keys": []map[string]interface{}{
								{"name": "key2222", "alias": "alisa22222"},
								{"name": "key1111", "alias": "alisa1111"},
							},
						},
						{
							"name":  "${var.name}-2",
							"token": ` #$^*\r\n\t`,
							"type":  "json",
							"json_keys": []map[string]interface{}{
								{"name": "key3333", "alias": "alisa3333"},
								{"name": "key4444", "alias": "alisa4444"},
							},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"full_text.#":    "1",
						"field_search.#": "2",
					}),
				),
			},
		},
	})
}

func resourceLogStoreIndexConfigDependence(name string) string {
	return fmt.Sprintf(`
	variable "name" {
	    default = "%s"
	}
	resource "alicloud_log_project" "default" {
	    name = var.name
	    description = "tf unit test"
	}
	resource "alicloud_log_store" "default" {
	    project = alicloud_log_project.default.name
	    name = var.name
	    retention_period = "3000"
	    shard_count = 1
	}
	`, name)
}

var logStoreIndexMap = map[string]string{
	"project":  CHECKSET,
	"logstore": CHECKSET,
}
