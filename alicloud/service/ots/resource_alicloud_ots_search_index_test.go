// collection-order: bypass schema.field_schema -- pre-existing attribute outside this change
// collection-order: bypass schema.index_setting.routing_fields -- pre-existing attribute outside this change
// collection-order: bypass schema.index_sort.sorter -- pre-existing attribute outside this change
package ots_test

import (
	"fmt"
	"testing"

	"github.com/aliyun/aliyun-tablestore-go-sdk/tablestore"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/ots"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccAliCloudOtsSearchIndex_basic(t *testing.T) {
	var v *tablestore.DescribeTableResponse

	resourceId := "alicloud_ots_search_index.default"
	ra := acctest.ResourceAttrInit(resourceId, otsSearchIndexBasicMap)
	serviceFunc := func() interface{} {
		return ots.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}
	rc := acctest.ResourceCheckInit(resourceId, &v, serviceFunc)

	rac := acctest.ResourceAttrCheckInit(rc, ra)

	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("testAcc%d", rand)
	//testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, resourceOtsSearchIndexConfigDependence)

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
				Config: resourceOtsSearchIndexConfig(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(getSearchIndexCheckMap()),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"schema", "time_to_live", "index_name", "table_name", "instance_name"},
			},
		},
	})
}

func resourceOtsSearchIndexConfig(name string) string {
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

  resource "alicloud_ots_search_index" "default" {
	  instance_name = "${alicloud_ots_instance.default.name}"
	  table_name = "${alicloud_ots_table.default.table_name}"

      index_name = "%s_search_index"
      time_to_live = -1
      schema {
          field_schema {
            field_name = "col1"
            field_type = "Text"
            is_array = false
            index = true
            analyzer = "Split"
            store = true
          }
          field_schema {
            field_name =  "col2"
             field_type = "Long"
             enable_sort_and_agg = true
          }
       

          field_schema {
            field_name =  "pk1"
            field_type = "Long"
            
          }
          field_schema {
            field_name =  "pk2"
            field_type = "Text"
            
          }


          index_setting {
            routing_fields = [ "pk1", "pk2"]
          }

          index_sort {
            sorter {
              sorter_type = "PrimaryKeySort"
              order = "Asc"
            }
            sorter {
              sorter_type = "FieldSort"
              order = "Desc"
              field_name =  "col2"
              mode = "Max"
            }
          }
    }
  }
`, name, string(ots.OtsCapacity), name)
}

func TestAccAliCloudOtsSearchIndex_ttl(t *testing.T) {
	var v *tablestore.DescribeTableResponse

	resourceId := "alicloud_ots_search_index.default"
	ra := acctest.ResourceAttrInit(resourceId, otsSearchIndexBasicMap)
	serviceFunc := func() interface{} {
		return ots.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}
	rc := acctest.ResourceCheckInit(resourceId, &v, serviceFunc)

	rac := acctest.ResourceAttrCheckInit(rc, ra)

	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("testAcc%d", rand)

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
				Config: resourceOtsSearchIndexTtlConfig(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(getSearchIndexCheckMap()),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"schema", "time_to_live", "index_name", "table_name", "instance_name"},
			},
		},
	})
}

func resourceOtsSearchIndexTtlConfig(name string) string {
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
      allow_update = false
  }

  resource "alicloud_ots_search_index" "default" {
	  instance_name = "${alicloud_ots_instance.default.name}"
	  table_name = "${alicloud_ots_table.default.table_name}"

      index_name = "%s_search_index"
      time_to_live = 86400
      schema {
          field_schema {
            field_name = "col1"
            field_type = "Text"
            is_array = false
            index = true
            analyzer = "Split"
            store = true
          }
          field_schema {
            field_name =  "col2"
             field_type = "Long"
             enable_sort_and_agg = true
          }
       

          field_schema {
            field_name =  "pk1"
            field_type = "Long"
            
          }
          field_schema {
            field_name =  "pk2"
            field_type = "Text"
            
          }
    }
  }
`, name, string(ots.OtsCapacity), name)
}

var otsSearchIndexBasicMap = map[string]string{}

func getSearchIndexCheckMap() map[string]string {
	accMap := make(map[string]string)
	accMap["index_id"] = acctest.CHECKSET
	accMap["create_time"] = acctest.CHECKSET
	accMap["sync_phase"] = acctest.CHECKSET
	accMap["current_sync_timestamp"] = acctest.CHECKSET

	return accMap
}
