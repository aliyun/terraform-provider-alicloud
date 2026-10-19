package alicloud

import (
	"fmt"
	"os"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// TestAccAliCloud DataWorksConnection. >>> Resource test cases.
// Case 0: covers create -> update (env_type, content, description, status) -> downgrade -> import.
func TestAccAliCloudDataWorksConnection_basic0(t *testing.T) {
	if v := os.Getenv("ALICLOUD_DATAWORKS_PROJECT_ID"); v == "" {
		t.Skip("Skipping: ALICLOUD_DATAWORKS_PROJECT_ID env var is required to run this test.")
	}
	var v map[string]interface{}
	resourceId := "alicloud_data_works_connection.default"
	ra := resourceAttrInit(resourceId, AlicloudDataWorksConnectionMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &DataworksPublicService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeDataWorksConnection")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sdataworksconnection%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudDataWorksConnectionBasicDependence0)
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
					"project_id":      "${var.project_id}",
					"connection_name": name,
					"connection_type": "mysql",
					"sub_type":        "mysql",
					"env_type":        0,
					"content":         "{\\\"database\\\":\\\"tf_example\\\",\\\"host\\\":\\\"127.0.0.1\\\",\\\"password\\\":\\\"tf_example_pw\\\",\\\"port\\\":\\\"3306\\\",\\\"username\\\":\\\"tf_example_user\\\"}",
					"description":     "tf-example-connection-desc-1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"connection_name": name,
						"connection_type": "mysql",
						"sub_type":        "mysql",
						"env_type":        "0",
						"description":     "tf-example-connection-desc-1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"project_id":      "${var.project_id}",
					"connection_name": name,
					"connection_type": "mysql",
					"sub_type":        "mysql",
					"env_type":        1,
					"content":         "{\\\"database\\\":\\\"tf_example_2\\\",\\\"host\\\":\\\"127.0.0.1\\\",\\\"password\\\":\\\"tf_example_pw_2\\\",\\\"port\\\":\\\"3306\\\",\\\"username\\\":\\\"tf_example_user_2\\\"}",
					"description":     "tf-example-connection-desc-2",
					"status":          "2",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"env_type":    "1",
						"description": "tf-example-connection-desc-2",
						"status":      "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"project_id":      "${var.project_id}",
					"connection_name": name,
					"connection_type": "mysql",
					"sub_type":        "mysql",
					"env_type":        0,
					"content":         "{\\\"database\\\":\\\"tf_example\\\",\\\"host\\\":\\\"127.0.0.1\\\",\\\"password\\\":\\\"tf_example_pw\\\",\\\"port\\\":\\\"3306\\\",\\\"username\\\":\\\"tf_example_user\\\"}",
					"description":     "tf-example-connection-desc-1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"env_type":    "0",
						"description": "tf-example-connection-desc-1",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"content", "description", "status"},
			},
		},
	})
}

var AlicloudDataWorksConnectionMap0 = map[string]string{
	"connection_id":          CHECKSET,
	"project_id":             CHECKSET,
	"connection_name":        CHECKSET,
	"connection_type":        CHECKSET,
	"sub_type":               CHECKSET,
	"env_type":               CHECKSET,
	"content":                NOSET,
	"description":            NOSET,
	"status":                 NOSET,
	"operator":               CHECKSET,
	"connect_status":         CHECKSET,
	"binding_calc_engine_id": CHECKSET,
	"gmt_modified":           CHECKSET,
	"sequence":               CHECKSET,
	"shared":                 CHECKSET,
	"default_engine":         CHECKSET,
	"create_time":            CHECKSET,
	"tenant_id":              CHECKSET,
	"region_id":              CHECKSET,
}

func AlicloudDataWorksConnectionBasicDependence0(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

variable "project_id" {
  default = "%s"
}
`, name, os.Getenv("ALICLOUD_DATAWORKS_PROJECT_ID"))
}
