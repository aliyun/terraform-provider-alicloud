package sls_test

import (
	"fmt"
	"testing"

	sls "github.com/aliyun/aliyun-log-go-sdk"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	slssvc "github.com/aliyun/terraform-provider-alicloud/alicloud/service/sls"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccAliCloudLogStoreIndex_basic(t *testing.T) {
	var v *sls.Index
	resourceId := "alicloud_log_store_index.default"
	ra := acctest.ResourceAttrInit(resourceId, logStoreIndexMap)
	serviceFunc := func() interface{} {
		return slssvc.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}
	rc := acctest.ResourceCheckInit(resourceId, &v, serviceFunc)
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testacclogstoreindex-%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, resourceLogStoreIndexConfigDependence)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
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
					"full_text":             acctest.REMOVEKEY,
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
						"log_reduce_black_list.0": acctest.REMOVEKEY,
						"log_reduce_white_list.#": "0",
						"log_reduce_white_list.0": acctest.REMOVEKEY,
						"max_text_len":            "1024",
						"full_text.#":             acctest.REMOVEKEY,
						"field_search.#":          "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"log_reduce":            acctest.REMOVEKEY,
					"log_reduce_black_list": acctest.REMOVEKEY,
					"log_reduce_white_list": acctest.REMOVEKEY,
					"max_text_len":          acctest.REMOVEKEY,
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
						"log_reduce":              acctest.REMOVEKEY,
						"log_reduce_black_list.#": acctest.REMOVEKEY,
						"log_reduce_white_list.#": acctest.REMOVEKEY,
						"max_text_len":            acctest.REMOVEKEY,
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

// TestAccAliCloudLogStoreIndex_logReduceListsReorder covers TypeListOrderCoverage
// for the top-level TypeList paths log_reduce_black_list and log_reduce_white_list.
// step0 applies the base config with both lists set to >=2 distinct members in
// order A (plus the Required attrs). Each following step is a delta that reorders
// exactly one list (a genuine permutation of order A); the builder inherits every
// other attribute unchanged, so each consecutive pair changes only one list and
// the post-apply plan converges on the reordered members.
func TestAccAliCloudLogStoreIndex_logReduceListsReorder(t *testing.T) {
	var v *sls.Index
	resourceId := "alicloud_log_store_index.default"
	ra := acctest.ResourceAttrInit(resourceId, logStoreIndexMap)
	serviceFunc := func() interface{} {
		return slssvc.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}
	rc := acctest.ResourceCheckInit(resourceId, &v, serviceFunc)
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testacclogstoreindex-%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, resourceLogStoreIndexConfigDependence)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"project":               "${alicloud_log_project.default.name}",
					"logstore":              "${alicloud_log_store.default.name}",
					"log_reduce":            true,
					"log_reduce_black_list": []interface{}{"b1", "b2", "b3"},
					"log_reduce_white_list": []interface{}{"w1", "w2", "w3"},
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
						"log_reduce":              "true",
						"log_reduce_black_list.#": "3",
						"log_reduce_white_list.#": "3",
						"max_text_len":            "2048",
						"full_text.#":             "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"log_reduce_black_list": []interface{}{"b3", "b2", "b1"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"log_reduce_black_list.#": "3",
						"log_reduce_white_list.#": "3",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"log_reduce_white_list": []interface{}{"w3", "w2", "w1"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"log_reduce_black_list.#": "3",
						"log_reduce_white_list.#": "3",
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
  name        = "${var.name}"
  description = "tf unit test"
}
resource "alicloud_log_store" "default" {
  project          = "${alicloud_log_project.default.name}"
  name             = "${var.name}"
  retention_period = "3000"
  shard_count      = 1
}
	`, name)
}

var logStoreIndexMap = map[string]string{
	"project":  acctest.CHECKSET,
	"logstore": acctest.CHECKSET,
}
