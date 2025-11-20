package alicloud

import (
	"fmt"
	"testing"

	sls "github.com/aliyun/aliyun-log-go-sdk"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAliCloudLogStoreIndex_basic(t *testing.T) {
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
					"project":               "${alicloud_log_project.default.name}",
					"logstore":              "${alicloud_log_store.default.name}",
					"log_reduce":            true,
					"log_reduce_black_list": []interface{}{"test"},
					"log_reduce_white_list": []interface{}{"name"},
					"max_text_len":          2048,
					"full_text": []map[string]interface{}{
						{
							"case_sensitive":  true,
							"token":           ` #$^*\r\n\t`,
							"include_chinese": true,
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"project":                 name,
						"logstore":                name,
						"log_reduce":              "true",
						"log_reduce_black_list.#": "1",
						"log_reduce_black_list.0": "test",
						"log_reduce_white_list.#": "1",
						"log_reduce_white_list.0": "name",
						"max_text_len":            "2048",
						"full_text.#":             "1",
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
					"log_reduce":            false,
					"log_reduce_black_list": []interface{}{},
					"log_reduce_white_list": []interface{}{},
					"max_text_len":          1024,
					"full_text":             REMOVEKEY,
					"field_search": []map[string]interface{}{
						{
							"name":             "${var.name}",
							"enable_analytics": true,
							"token":            ` #$^*\r\n\t`,
							"type":             "text",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"log_reduce":              "false",
						"log_reduce_black_list.#": "0",
						"log_reduce_black_list.0": REMOVEKEY,
						"log_reduce_white_list.#": "0",
						"log_reduce_white_list.0": REMOVEKEY,
						"max_text_len":            "1024",
						"full_text.#":             REMOVEKEY,
						"field_search.#":          "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"log_reduce":            REMOVEKEY,
					"log_reduce_black_list": REMOVEKEY,
					"log_reduce_white_list": REMOVEKEY,
					"max_text_len":          REMOVEKEY,
					"full_text": []map[string]interface{}{
						{
							"case_sensitive":  "true",
							"token":           `,`,
							"include_chinese": "false",
						},
					},
					"field_search": []map[string]interface{}{
						{
							"name":             "${var.name}-1",
							"enable_analytics": "true",
							"token":            ` #$^*\r\n\t`,
							"type":             "json",
							"json_keys": []map[string]interface{}{
								{"name": "key2222", "alias": "alisa22222", "doc_value": true, "type": "long"},
								{"name": "key1111", "alias": "alisa1111", "doc_value": false, "type": "text"},
							},
						},
						{
							"name":            "${var.name}-2",
							"token":           ` #$^*\r\n\t`,
							"type":            "json",
							"case_sensitive":  true,
							"include_chinese": true,
							"alias":           "json_alias",
							"json_keys": []map[string]interface{}{
								{"name": "key3333", "alias": "alisa3333"},
								{"name": "key4444", "alias": "alisa4444"},
							},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"log_reduce":              REMOVEKEY,
						"log_reduce_black_list.#": REMOVEKEY,
						"log_reduce_white_list.#": REMOVEKEY,
						"max_text_len":            REMOVEKEY,
						"full_text.#":             "1",
						"field_search.#":          "2",
					}),
				),
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
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
	    name = "${var.name}"
	    description = "tf unit test"
	}
	resource "alicloud_log_store" "default" {
	    project = "${alicloud_log_project.default.name}"
	    name = "${var.name}"
	    retention_period = "3000"
	    shard_count = 1
	}
	`, name)
}

var logStoreIndexMap = map[string]string{
	"project":  CHECKSET,
	"logstore": CHECKSET,
}
