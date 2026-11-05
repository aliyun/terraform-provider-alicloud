package dmsenterprise_test

import (
	"fmt"
	"testing"

	tfacctest "github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	dmsenterprise "github.com/aliyun/terraform-provider-alicloud/alicloud/service/dmsenterprise"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// Case 1
func SkipTestAccAlicloudDms_enterpriseLogicDatabase_basic1887(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_dms_enterprise_logic_database.default"
	ra := tfacctest.ResourceAttrInit(resourceId, AlicloudDms_enterpriseLogicDatabaseMap1887)
	rc := tfacctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return dmsenterprise.NewDmsEnterpriseService(tfacctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeDmsEnterpriseLogicDatabase")
	rac := tfacctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sDmsLogicDatabase%d", tfacctest.DefaultRegionToTest, rand)
	testAccConfig := tfacctest.ResourceTestAccConfig(resourceId, name, AlicloudDms_enterpriseLogicDatabaseBasicDependence1887)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tfacctest.PreCheck(t)
		},
		ProviderFactories: tfacctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"alias": "${var.name}",
					"database_ids": []string{
						"${data.alicloud_dms_enterprise_databases.test2.databases.0.id}",
						"${data.alicloud_dms_enterprise_databases.test3.databases.0.id}",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"alias":             tfacctest.CHECKSET,
						"database_ids.#":    "2",
						"logic_database_id": tfacctest.CHECKSET,
						"search_name":       tfacctest.CHECKSET,
						"schema_name":       tfacctest.CHECKSET,
						"db_type":           tfacctest.CHECKSET,
						"env_type":          tfacctest.CHECKSET,
						"owner_id_list.#":   "1",
						"owner_name_list.#": "1",
					}),
				),
			}, {
				Config: testAccConfig(map[string]interface{}{
					"alias": "${var.name}_update",
					"database_ids": []string{
						"${data.alicloud_dms_enterprise_databases.test4.databases.0.id}",
						"${data.alicloud_dms_enterprise_databases.test5.databases.0.id}",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"alias":          tfacctest.CHECKSET,
						"database_ids.#": "2",
					}),
				),
			}, {
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{},
			},
		},
	})
}

var AlicloudDms_enterpriseLogicDatabaseMap1887 = map[string]string{}

func AlicloudDms_enterpriseLogicDatabaseBasicDependence1887(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

data "alicloud_dms_enterprise_instances" "dms_enterprise_instances_ds" {
  instance_type = "mysql"
  search_key    = "tf-test-no-deleting"
}

data "alicloud_dms_enterprise_databases" "test2" {
  name_regex  = "test2"
  instance_id = data.alicloud_dms_enterprise_instances.dms_enterprise_instances_ds.instances.0.instance_id
}

data "alicloud_dms_enterprise_databases" "test3" {
  name_regex  = "test3"
  instance_id = data.alicloud_dms_enterprise_instances.dms_enterprise_instances_ds.instances.0.instance_id
}

data "alicloud_dms_enterprise_databases" "test4" {
  name_regex  = "test4"
  instance_id = data.alicloud_dms_enterprise_instances.dms_enterprise_instances_ds.instances.0.instance_id
}

data "alicloud_dms_enterprise_databases" "test5" {
  name_regex  = "test5"
  instance_id = data.alicloud_dms_enterprise_instances.dms_enterprise_instances_ds.instances.0.instance_id
}
`, name)
}
