package alicloud

import (
	"fmt"
	"testing"

	sls "github.com/aliyun/aliyun-log-go-sdk"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAlicloudLogSavedSearch_basic(t *testing.T) {
	var v *sls.SavedSearch
	resourceId := "alicloud_log_saved_search.default"
	ra := resourceAttrInit(resourceId, logSavedSearchMap)
	serviceFunc := func() interface{} {
		return &LogService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("test-log-saved-search-%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceLogSavedSearchConfigDependence)

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
					"project_name":  name,
					"logstore_name": name,
					"search_name":   "test_saved_search",
					"display_name":  "test-log",
					"search_query":  "* | select count(*) as c,__time__ as t group by t order by t DESC limit 10",
					"topic":         "sls-test",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"project_name":  name,
						"logstore_name": name,
						"display_name":  "test-log",
						"search_name":   "test_saved_search",
						"search_query":  "* | select count(*) as c,__time__ as t group by t order by t DESC limit 10",
						"topic":         "sls-test",
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
					"search_query": "* | select count(*) as c,__time__ as t group by t order by t DESC limit 20",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"search_query": "* | select count(*) as c,__time__ as t group by t order by t DESC limit 20",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"display_name": "test-saved-search",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"display_name": "test-saved-search",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"topic": "test-log-saved-search",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"topic": "test-log-saved-search",
					}),
				),
			},
		},
	})
}

func resourceLogSavedSearchConfigDependence(name string) string {
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

var logSavedSearchMap = map[string]string{
	"project_name":  CHECKSET,
	"logstore_name": CHECKSET,
	"search_name":   CHECKSET,
}
