package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// Test DataWorks DataSourceSharedRule. >>> Resource test cases, automatically generated.
// Case DataSourceSharedRule-TF验收_成都 8955
func TestAccAliCloudDataWorksDataSourceSharedRule_basic8955(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_data_works_data_source_shared_rule.default"
	ra := resourceAttrInit(resourceId, AlicloudDataWorksDataSourceSharedRuleMap8955)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &DataWorksServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeDataWorksDataSourceSharedRule")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf_testacc_dwpt%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudDataWorksDataSourceSharedRuleBasicDependence8955)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-chengdu"})
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"target_project_id": "${alicloud_data_works_project.defaultasjsH5.id}",
					"data_source_id":    "${alicloud_data_works_data_source.defaultvzu0wG.data_source_id}",
					"env_type":          "Prod",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"target_project_id": CHECKSET,
						"data_source_id":    CHECKSET,
						"env_type":          "Prod",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{},
			},
		},
	})
}

var AlicloudDataWorksDataSourceSharedRuleMap8955 = map[string]string{
	"create_time":                CHECKSET,
	"data_source_shared_rule_id": CHECKSET,
}

func AlicloudDataWorksDataSourceSharedRuleBasicDependence8955(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

resource "alicloud_data_works_project" "defaultQeRfvU" {
  description  = "源项目"
  project_name = var.name
  display_name = "shared_source2"
  pai_task_enabled = true
}

resource "alicloud_data_works_project" "defaultasjsH5" {
  description  = "目标空间"
  project_name = format("%%s1", var.name)
  display_name = "shared_target2"
  pai_task_enabled = true
}

resource "alicloud_data_works_data_source" "defaultvzu0wG" {
  type                       = "hive"
  data_source_name           = format("%%s2", var.name)
  connection_properties      = jsonencode({ "address" : [{ "host" : "127.0.0.1", "port" : "1234" }], "database" : "hive_database", "metaType" : "HiveMetastore", "metastoreUris" : "thrift://123:123", "version" : "2.3.9", "loginMode" : "Anonymous", "securityProtocol" : "authTypeNone", "envType" : "Prod", "properties" : { "key1" : "value1" } })
  project_id                 = alicloud_data_works_project.defaultQeRfvU.id
  connection_properties_mode = "UrlMode"
}


`, name)
}

// Case DataSourceSharedRule-SharedUser_正式副本 8166
func TestAccAliCloudDataWorksDataSourceSharedRule_basic8166(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_data_works_data_source_shared_rule.default"
	ra := resourceAttrInit(resourceId, AlicloudDataWorksDataSourceSharedRuleMap8166)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &DataWorksServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeDataWorksDataSourceSharedRule")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf_testacc_dwpt%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudDataWorksDataSourceSharedRuleBasicDependence8166)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-chengdu"})
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"shared_user":       "300803888805783368",
					"target_project_id": "${alicloud_data_works_project.defaultGTFU1x.id}",
					"data_source_id":    "${alicloud_data_works_data_source.defaultsQybqp.data_source_id}",
					"env_type":          "Prod",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"shared_user":       "300803888805783368",
						"target_project_id": CHECKSET,
						"data_source_id":    CHECKSET,
						"env_type":          "Prod",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{},
			},
		},
	})
}

var AlicloudDataWorksDataSourceSharedRuleMap8166 = map[string]string{
	"create_time":                CHECKSET,
	"data_source_shared_rule_id": CHECKSET,
}

func AlicloudDataWorksDataSourceSharedRuleBasicDependence8166(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

resource "alicloud_data_works_project" "defaultFW8Oua" {
  description  = "源项目"
  project_name = var.name
  display_name = "shared_source1"
  pai_task_enabled = true
}

resource "alicloud_data_works_project" "defaultGTFU1x" {
  description  = "目标项目"
  project_name = format("%%s1", var.name)
  display_name = "target_source1"
  pai_task_enabled = true
}

resource "alicloud_data_works_data_source" "defaultsQybqp" {
  type                       = "hive"
  data_source_name           = format("%%s2", var.name)
  connection_properties      = jsonencode({ "address" : [{ "host" : "127.0.0.1", "port" : "1234" }], "database" : "hive_database", "metaType" : "HiveMetastore", "metastoreUris" : "thrift://123:123", "version" : "2.3.9", "loginMode" : "Anonymous", "securityProtocol" : "authTypeNone", "envType" : "Prod", "properties" : { "key1" : "value1" } })
  project_id                 = alicloud_data_works_project.defaultFW8Oua.id
  connection_properties_mode = "UrlMode"
}


`, name)
}

// Test DataWorks DataSourceSharedRule. <<< Resource test cases, automatically generated.
