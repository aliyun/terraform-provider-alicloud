package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// Test DataWorks DataSource. >>> Resource test cases, automatically generated.
// Case DataSource资源测试_成都 8927
func TestAccAliCloudDataWorksDataSource_basic8927(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_data_works_data_source.default"
	ra := resourceAttrInit(resourceId, AlicloudDataWorksDataSourceMap8927)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &DataWorksServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeDataWorksDataSource")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf_testacc_dwpt%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudDataWorksDataSourceBasicDependence8927)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			//testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-chengdu"})
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"type":                       "hive",
					"data_source_name":           name,
					"connection_properties":      "{   \\\"address\\\": [     {       \\\"host\\\": \\\"127.0.0.1\\\",       \\\"port\\\": \\\"1234\\\"     }   ],   \\\"database\\\": \\\"hive_database\\\",   \\\"metaType\\\": \\\"HiveMetastore\\\",   \\\"metastoreUris\\\": \\\"thrift://123:123\\\",   \\\"version\\\": \\\"2.3.9\\\",   \\\"loginMode\\\": \\\"Anonymous\\\",   \\\"securityProtocol\\\": \\\"authTypeNone\\\",   \\\"envType\\\": \\\"Prod\\\",   \\\"properties\\\": {     \\\"key1\\\": \\\"value1\\\"   } }",
					"connection_properties_mode": "UrlMode",
					"project_id":                 "${alicloud_data_works_project.defaultkguw4R.id}",
					"description":                "描述信息-初始状态",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"type":                       "hive",
						"data_source_name":           name,
						"connection_properties":      "{   \"address\": [     {       \"host\": \"127.0.0.1\",       \"port\": \"1234\"     }   ],   \"database\": \"hive_database\",   \"metaType\": \"HiveMetastore\",   \"metastoreUris\": \"thrift://123:123\",   \"version\": \"2.3.9\",   \"loginMode\": \"Anonymous\",   \"securityProtocol\": \"authTypeNone\",   \"envType\": \"Prod\",   \"properties\": {     \"key1\": \"value1\"   } }",
						"connection_properties_mode": "UrlMode",
						"project_id":                 CHECKSET,
						"description":                "描述信息-初始状态",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"connection_properties":      "{     \\\"clusterIdentifier\\\": \\\"cdh_cluster\\\",     \\\"database\\\": \\\"hive_database\\\",     \\\"loginMode\\\": \\\"Anonymous\\\",     \\\"securityProtocol\\\": \\\"authTypeNone\\\",     \\\"envType\\\": \\\"Prod\\\" }",
					"connection_properties_mode": "CdhMode",
					"description":                "描述信息-状态1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"connection_properties":      "{     \"clusterIdentifier\": \"cdh_cluster\",     \"database\": \"hive_database\",     \"loginMode\": \"Anonymous\",     \"securityProtocol\": \"authTypeNone\",     \"envType\": \"Prod\" }",
						"connection_properties_mode": "CdhMode",
						"description":                "描述信息-状态1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"connection_properties":      "{   \\\"address\\\": [     {       \\\"host\\\": \\\"127.0.0.1\\\",       \\\"port\\\": \\\"1234\\\"     }   ],   \\\"database\\\": \\\"hive_database\\\",   \\\"metaType\\\": \\\"HiveMetastore\\\",   \\\"metastoreUris\\\": \\\"thrift://123:123\\\",   \\\"version\\\": \\\"2.3.9\\\",   \\\"loginMode\\\": \\\"Anonymous\\\",   \\\"securityProtocol\\\": \\\"authTypeNone\\\",   \\\"envType\\\": \\\"Prod\\\",   \\\"properties\\\": {     \\\"key1\\\": \\\"value1\\\"   } }",
					"connection_properties_mode": "UrlMode",
					"description":                "描述信息-最终状态",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"connection_properties":      "{   \"address\": [     {       \"host\": \"127.0.0.1\",       \"port\": \"1234\"     }   ],   \"database\": \"hive_database\",   \"metaType\": \"HiveMetastore\",   \"metastoreUris\": \"thrift://123:123\",   \"version\": \"2.3.9\",   \"loginMode\": \"Anonymous\",   \"securityProtocol\": \"authTypeNone\",   \"envType\": \"Prod\",   \"properties\": {     \"key1\": \"value1\"   } }",
						"connection_properties_mode": "UrlMode",
						"description":                "描述信息-最终状态",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"connection_properties"},
			},
		},
	})
}

var AlicloudDataWorksDataSourceMap8927 = map[string]string{
	"create_time":    CHECKSET,
	"data_source_id": CHECKSET,
}

func AlicloudDataWorksDataSourceBasicDependence8927(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

resource "alicloud_data_works_project" "defaultkguw4R" {
  description  = "匠承测试"
  project_name = var.name
  display_name = "jiangcheng_terraform_test2"
  pai_task_enabled = true
}


`, name)
}

// Test DataWorks DataSource. <<< Resource test cases, automatically generated.
