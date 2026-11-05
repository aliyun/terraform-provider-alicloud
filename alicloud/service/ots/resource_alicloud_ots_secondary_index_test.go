// collection-order: bypass defined_columns -- pre-existing attribute outside this change
// collection-order: bypass primary_keys -- pre-existing attribute outside this change
package ots_test

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/aliyun/aliyun-tablestore-go-sdk/tablestore"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/ots"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccAliCloudOtsSecondaryIndex_basic(t *testing.T) {
	var v *tablestore.DescribeTableResponse

	resourceId := "alicloud_ots_secondary_index.default"
	ra := acctest.ResourceAttrInit(resourceId, otsSecondaryIndexBasicMap)
	serviceFunc := func() interface{} {
		return ots.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}
	rc := acctest.ResourceCheckInit(resourceId, &v, serviceFunc)

	rac := acctest.ResourceAttrCheckInit(rc, ra)

	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("testAcc%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, resourceOtsSecondaryIndexConfigDependence)
	pks := []string{"pk1", "pk2"}
	cols := []string{"col2", "col1"}

	randIndexType := string(random(ots.Local, ots.Global).(ots.SecondaryIndexTypeString))
	randIncludeBaseData := random(true, false).(bool)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckWithRegions(t, false, connectivity.OtsCapacityNoSupportedRegions)
		},
		// module name
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_name":     "${alicloud_ots_instance.default.name}",
					"table_name":        "${alicloud_ots_table.default.table_name}",
					"index_name":        "${var.name}" + "_index",
					"index_type":        randIndexType,
					"include_base_data": strconv.FormatBool(randIncludeBaseData),
					//"${var.pks}"
					"primary_keys": "${var.pks}",
					//"${var.cols}"
					"defined_columns": "${var.cols}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(getCheckMap(name, randIndexType, pks, cols)),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"include_base_data"},
			},
		},
	})
}

func resourceOtsSecondaryIndexConfigDependence(name string) string {
	return fmt.Sprintf(`
	variable "name" {
	  default = "%s"
	}
	variable "pks" {
	  default = ["pk1", "pk2"]
	}
	variable "cols" {
	  default = ["col2", "col1"]
	}
	resource "alicloud_ots_instance" "default" {
	  name = "tf-${var.name}"
	  description = "${var.name}"
	  accessed_by = "Any"
	  instance_type = "%s"
	  tags = {
	    Created = "TF"
	    For = "acceptance test"
	  }
	}
	resource "alicloud_ots_table" "default" {
	  instance_name = "${alicloud_ots_instance.default.name}"
	  table_name = "${var.name}"
	  primary_key {
          name = "pk1"
	      type = "Integer"
	  }
	  primary_key {
          name = "pk2"
          type = "String"
      }

	  defined_column {
          name = "col1"
          type = "String"
      }
	  defined_column {
          name = "col2"
          type = "Integer"
      }

	  time_to_live = -1
	  max_version = 1
	}
	`, name, string(ots.OtsCapacity))
}

var otsSecondaryIndexBasicMap = map[string]string{}

func random(values ...interface{}) interface{} {
	r := tfacctest.RandInt() % len(values)
	return values[r]
}

func getCheckMap(name string, indexType string, pks []string, cols []string) map[string]string {
	accMap := make(map[string]string)
	accMap["instance_name"] = "tf-" + name
	accMap["table_name"] = name
	accMap["index_name"] = name + "_index"
	accMap["index_type"] = indexType
	accMap["primary_keys.#"] = "2"
	accMap["defined_columns.#"] = "2"
	for i, pk := range pks {
		accMap["primary_keys."+strconv.Itoa(i)] = pk
	}
	for i, col := range cols {
		accMap["defined_columns."+strconv.Itoa(i)] = col
	}

	return accMap
}
