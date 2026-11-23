// collection-order: bypass schema.index_setting.routing_fields -- pre-existing attribute outside this change
// collection-order: bypass schema.index_sort.sorter -- pre-existing attribute outside this change
package ots_test

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"regexp"
	"sort"
	"strings"
	"sync"
	"testing"

	"github.com/aliyun/aliyun-tablestore-go-sdk/tablestore"
	"github.com/aliyun/aliyun-tablestore-go-sdk/tablestore/otsprotocol"
	"github.com/aliyun/credentials-go/credentials"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/ots"
	"github.com/golang/protobuf/proto"
	"github.com/hashicorp/go-cty/cty"
	sdkSchema "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	sdkterraform "github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccAliCloudOtsSearchIndex_basic(t *testing.T) {
	var v *tablestore.DescribeTableResponse

	resourceId := "alicloud_ots_search_index.default"
	ra := acctest.ResourceAttrInit(resourceId, searchIndexBasicMap)
	serviceFunc := func() interface{} {
		return ots.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}
	rc := acctest.ResourceCheckInit(resourceId, &v, serviceFunc)

	rac := acctest.ResourceAttrCheckInit(rc, ra)

	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("testAcc%d", rand)
	//testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, resourceOtsSearchIndexConfigDependence)

	var applied *sdkterraform.InstanceState

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
					captureSearchIndexState(resourceId, &applied),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateCheck:  searchIndexImportStateCheck(&applied, searchIndexTopLevelStoreNotReturned),
				ImportStateVerify: false,
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
	ra := acctest.ResourceAttrInit(resourceId, searchIndexBasicMap)
	serviceFunc := func() interface{} {
		return ots.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}
	rc := acctest.ResourceCheckInit(resourceId, &v, serviceFunc)

	rac := acctest.ResourceAttrCheckInit(rc, ra)

	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("testAcc%d", rand)

	var applied *sdkterraform.InstanceState

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
					captureSearchIndexState(resourceId, &applied),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateCheck:  searchIndexImportStateCheck(&applied, searchIndexTopLevelStoreNotReturned),
				ImportStateVerify: false,
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

func TestAccAliCloudOtsSearchIndex_complete(t *testing.T) {
	var v *tablestore.DescribeTableResponse

	resourceId := "alicloud_ots_search_index.default"
	ra := acctest.ResourceAttrInit(resourceId, searchIndexBasicMap)
	serviceFunc := func() interface{} {
		return ots.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}
	rc := acctest.ResourceCheckInit(resourceId, &v, serviceFunc)
	rac := acctest.ResourceAttrCheckInit(rc, ra)

	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("testAcc%d", rand)

	var applied *sdkterraform.InstanceState

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckWithRegions(t, false, connectivity.OtsCapacityNoSupportedRegions)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: resourceSearchIndexCompleteConfig(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(getSearchIndexCheckMap()),
					checkSearchIndexSchemaFields(),
					captureSearchIndexState(resourceId, &applied),
				),
			},
			{
				ResourceName:     resourceId,
				ImportState:      true,
				ImportStateCheck: searchIndexImportStateCheck(&applied, searchIndexCompleteStoreOverrides),
			},
		},
	})
}

// searchIndexTopLevelStoreNotReturned is the precise import exception for the
// documented service limitation behind OPEN-01: DescribeSearchIndex reports
// `store=false` (or omits the value) for a top-level field whose create request
// sent `store=true`. The run on head 66d591ea5 showed exactly this difference
// for the Text, Vector, FuzzyKeyword and virtual columns, while nested
// sub-fields keep their store value. Every other leaf is compared strictly.
var searchIndexTopLevelStoreNotReturned = map[string]bool{
	"col1": false,
}

// searchIndexCompleteStoreOverrides completes searchIndexTopLevelStoreNotReturned
// for the complete configuration: a nested sub-field whose argument is omitted
// follows the SDK default (true for every type except Text) and the service
// echoes it, while the applied state materializes the omitted boolean as false.
var searchIndexCompleteStoreOverrides = map[string]bool{
	"col1":                              false,
	"vector_col":                        false,
	"fuzzy_col":                         false,
	"virtual_keyword":                   false,
	"nested_col.nested_long":            true,
	"nested_col.nested_keyword_array":   true,
	"nested_col.nested_keyword_noindex": true,
	"nested_col.nested_date":            true,
	"nested_col.nested_virtual":         true,
}

func resourceSearchIndexCompleteConfig(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

resource "alicloud_ots_instance" "default" {
  name          = "tf-${var.name}"
  description   = "${var.name}"
  accessed_by   = "Any"
  instance_type = "%s"
  tags = {
    Created = "TF"
    For     = "acceptance test"
  }
}

resource "alicloud_ots_table" "default" {
  instance_name = alicloud_ots_instance.default.name
  table_name    = var.name
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
  defined_column {
    name = "vector_col"
    type = "String"
  }
  defined_column {
    name = "fuzzy_col"
    type = "String"
  }
  defined_column {
    name = "nested_col"
    type = "String"
  }
  defined_column {
    name = "date_col"
    type = "String"
  }
  time_to_live = -1
  max_version  = 1
}

resource "alicloud_ots_search_index" "default" {
  instance_name = alicloud_ots_instance.default.name
  table_name    = alicloud_ots_table.default.table_name
  index_name    = "%s_search_index"
  time_to_live  = -1

  schema {
    field_schema {
      field_name = "col1"
      field_type = "Text"
      analyzer   = "Split"
      store      = true
    }
    field_schema {
      field_name          = "col2"
      field_type          = "Long"
      enable_sort_and_agg = true
    }
    field_schema {
      field_name = "vector_col"
      field_type = "Vector"
      vector_options {
        data_type   = "float32"
        dimension   = 4
        metric_type = "cosine"
      }
    }
    field_schema {
      field_name = "fuzzy_col"
      field_type = "FuzzyKeyword"
    }
    field_schema {
      field_name = "nested_col"
      field_type = "Nested"
      index      = false
      sub_field_schemas {
        field_name          = "nested_long"
        field_type          = "Long"
        enable_sort_and_agg = true
      }
      sub_field_schemas {
        field_name = "nested_text"
        field_type = "Text"
        analyzer   = "SingleWord"
        store      = true
      }
      sub_field_schemas {
        field_name = "nested_keyword_array"
        field_type = "Keyword"
        is_array   = true
      }
      sub_field_schemas {
        field_name = "nested_keyword_noindex"
        field_type = "Keyword"
        index      = false
      }
      sub_field_schemas {
        field_name = "nested_keyword_nostore"
        field_type = "Keyword"
        store      = false
      }
      sub_field_schemas {
        field_name       = "nested_date"
        field_type       = "Date"
        date_format_list = ["epoch_millis"]
      }
      sub_field_schemas {
        field_name        = "nested_virtual"
        field_type        = "Keyword"
        is_virtual_field  = true
        source_field_name = "pk2"
      }
    }
    field_schema {
      field_name       = "date_col"
      field_type       = "Date"
      date_format_list = ["yyyy-MM-dd HH:mm:ss", "yyyy-MM-dd"]
    }
    field_schema {
      field_name        = "virtual_keyword"
      field_type        = "Keyword"
      is_virtual_field  = true
      source_field_name = "pk2"
    }
    field_schema {
      field_name = "pk2"
      field_type = "Text"
    }
  }
}
`, name, string(ots.OtsCapacity), name)
}

func checkSearchIndexSchemaFields() resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources["alicloud_ots_search_index.default"]
		if !ok {
			return fmt.Errorf("can not find alicloud_ots_search_index.default in state")
		}

		otsService := ots.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
		resp, err := otsService.DescribeOtsSearchIndex(rs.Primary.ID)
		if err != nil {
			return err
		}
		if resp == nil || resp.Schema == nil {
			return fmt.Errorf("search index schema is empty")
		}

		fields := make(map[string]*tablestore.FieldSchema)
		var walk func([]*tablestore.FieldSchema)
		walk = func(schemas []*tablestore.FieldSchema) {
			for _, fieldSchema := range schemas {
				if fieldSchema.FieldName != nil {
					fields[*fieldSchema.FieldName] = fieldSchema
				}
				walk(fieldSchema.FieldSchemas)
			}
		}
		walk(resp.Schema.FieldSchemas)

		expectedFieldTypes := map[string]tablestore.FieldType{
			"vector_col":             tablestore.FieldType_VECTOR,
			"fuzzy_col":              tablestore.FieldType_FUZZY_KEYWORD,
			"nested_col":             tablestore.FieldType_NESTED,
			"nested_long":            tablestore.FieldType_LONG,
			"nested_text":            tablestore.FieldType_TEXT,
			"nested_keyword_array":   tablestore.FieldType_KEYWORD,
			"nested_keyword_noindex": tablestore.FieldType_KEYWORD,
			"nested_keyword_nostore": tablestore.FieldType_KEYWORD,
			"col2":                   tablestore.FieldType_LONG,
			"nested_date":            tablestore.FieldType_DATE,
			"nested_virtual":         tablestore.FieldType_KEYWORD,
			"date_col":               tablestore.FieldType_DATE,
			"virtual_keyword":        tablestore.FieldType_KEYWORD,
		}
		for name, expectedFieldType := range expectedFieldTypes {
			fieldSchema, ok := fields[name]
			if !ok {
				return fmt.Errorf("field %s is missing from DescribeSearchIndex response", name)
			}
			if fieldSchema.FieldType != expectedFieldType {
				return fmt.Errorf("field %s type is %s, expected %s", name, fieldSchema.FieldType.String(), expectedFieldType.String())
			}
		}

		vectorField := fields["vector_col"]
		if vectorField.VectorOptions == nil {
			return fmt.Errorf("vector_col has no vector options in DescribeSearchIndex response")
		}
		if vectorField.VectorOptions.Dimension == nil || *vectorField.VectorOptions.Dimension != 4 {
			return fmt.Errorf("vector_col dimension is not 4")
		}
		if vectorField.VectorOptions.VectorDataType == nil || *vectorField.VectorOptions.VectorDataType != tablestore.VectorDataType_FLOAT_32 {
			return fmt.Errorf("vector_col data type is not float32")
		}
		if vectorField.VectorOptions.VectorMetricType == nil || *vectorField.VectorOptions.VectorMetricType != tablestore.VectorMetricType_COSINE {
			return fmt.Errorf("vector_col metric type is not cosine")
		}

		virtualField := fields["virtual_keyword"]
		if virtualField.IsVirtualField == nil || !*virtualField.IsVirtualField {
			return fmt.Errorf("virtual_keyword is not marked as a virtual field")
		}
		if len(virtualField.SourceFieldNames) != 1 || virtualField.SourceFieldNames[0] != "pk2" {
			return fmt.Errorf("virtual_keyword source field names are %v", virtualField.SourceFieldNames)
		}

		nestedVirtual := fields["nested_virtual"]
		if nestedVirtual.IsVirtualField == nil || !*nestedVirtual.IsVirtualField {
			return fmt.Errorf("nested_virtual is not marked as a virtual field")
		}
		if len(nestedVirtual.SourceFieldNames) != 1 || nestedVirtual.SourceFieldNames[0] != "pk2" {
			return fmt.Errorf("nested_virtual source field names are %v", nestedVirtual.SourceFieldNames)
		}

		expectedDateFormats := map[string]struct{}{
			"yyyy-MM-dd HH:mm:ss": {},
			"yyyy-MM-dd":          {},
		}
		dateFormats := fields["date_col"].DateFormats
		if len(dateFormats) != len(expectedDateFormats) {
			return fmt.Errorf("date_col date formats are %v", dateFormats)
		}
		for _, dateFormat := range dateFormats {
			if _, ok := expectedDateFormats[dateFormat]; !ok {
				return fmt.Errorf("date_col has unexpected date format %s", dateFormat)
			}
		}

		expectedNestedDateFormats := map[string]struct{}{
			"epoch_millis": {},
		}
		nestedDateFormats := fields["nested_date"].DateFormats
		if len(nestedDateFormats) != len(expectedNestedDateFormats) {
			return fmt.Errorf("nested_date date formats are %v", nestedDateFormats)
		}
		for _, dateFormat := range nestedDateFormats {
			if _, ok := expectedNestedDateFormats[dateFormat]; !ok {
				return fmt.Errorf("nested_date has unexpected date format %s", dateFormat)
			}
		}

		nestedContainer := fields["nested_col"]
		if nestedContainer.Index == nil || *nestedContainer.Index {
			return fmt.Errorf("nested_col index is not disabled")
		}

		nestedLong := fields["nested_long"]
		if nestedLong.EnableSortAndAgg == nil || !*nestedLong.EnableSortAndAgg {
			return fmt.Errorf("nested_long enable_sort_and_agg is not enabled")
		}

		// Text defaults to Store=false in the SDK, so store=true is a
		// discriminating read-back for the sub-field attribute.
		nestedText := fields["nested_text"]
		if nestedText.Analyzer == nil || *nestedText.Analyzer != tablestore.Analyzer_SingleWord {
			return fmt.Errorf("nested_text analyzer is not SingleWord")
		}
		if nestedText.Store == nil || !*nestedText.Store {
			return fmt.Errorf("nested_text store is not enabled")
		}

		nestedKeywordArray := fields["nested_keyword_array"]
		if nestedKeywordArray.IsArray == nil || !*nestedKeywordArray.IsArray {
			return fmt.Errorf("nested_keyword_array is_array is not enabled")
		}

		nestedKeywordNoIndex := fields["nested_keyword_noindex"]
		if nestedKeywordNoIndex.Index == nil || *nestedKeywordNoIndex.Index {
			return fmt.Errorf("nested_keyword_noindex index is not disabled")
		}

		// `store` is an optional boolean without a Schema default: an omitted
		// argument must reach the SDK as nil so a non-Text field falls back to
		// the SDK default (true), while an explicit false must stay false.
		if nestedKeywordArray.Store == nil || !*nestedKeywordArray.Store {
			return fmt.Errorf("nested_keyword_array must keep the SDK store default true when store is omitted")
		}
		nestedKeywordNoStore := fields["nested_keyword_nostore"]
		if nestedKeywordNoStore.Store == nil || *nestedKeywordNoStore.Store {
			return fmt.Errorf("nested_keyword_nostore must send the explicit store=false")
		}
		// The top-level field types that existed before this change keep the
		// legacy explicit store=false when the argument is omitted.
		if col2 := fields["col2"]; col2.Store == nil || *col2.Store {
			return fmt.Errorf("col2 must keep the legacy explicit store=false when store is omitted")
		}
		// col1 is a top-level Text field configured with store=true, but
		// DescribeSearchIndex returns store=false for it even though the same
		// argument on the Nested sub-field nested_text is returned as true. The
		// service therefore does not round-trip store for this top-level field,
		// so the wire shape is pinned by TestUnitOtsSearchIndexStoreWithoutRawConfig
		// instead of a read-back assertion.
		// The three states are distinguishable on Keyword fields: omitted keeps
		// the SDK default (true), an explicit false stays false and nested_text
		// above proves an explicit true. FuzzyKeyword, Vector and virtual fields
		// follow the same parser path but the service does not echo a
		// distinguishing store value for them, so their types are asserted above
		// instead of a store value that only reflects server-side defaults.

		return nil
	}
}

func TestAccAliCloudOtsSearchIndex_fieldOrder(t *testing.T) {
	var v *tablestore.DescribeTableResponse

	resourceId := "alicloud_ots_search_index.default"
	ra := acctest.ResourceAttrInit(resourceId, searchIndexBasicMap)
	serviceFunc := func() interface{} {
		return ots.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}
	rc := acctest.ResourceCheckInit(resourceId, &v, serviceFunc)
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	name := fmt.Sprintf("testAcc%d", tfacctest.RandIntRange(10000, 99999))
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, resourceOtsSearchIndexFieldOrderConfigDependence)

	var applied *sdkterraform.InstanceState

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckWithRegions(t, false, connectivity.OtsCapacityNoSupportedRegions)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_name": "${alicloud_ots_instance.default.name}",
					"table_name":    "${alicloud_ots_table.default.table_name}",
					"index_name":    "${var.name}_search_index",
					"time_to_live":  "-1",
					"schema": []map[string]interface{}{
						{
							"field_schema": []map[string]interface{}{
								{
									"field_name": "col1",
									"field_type": "Text",
									"analyzer":   "Split",
									"store":      "true",
								},
								{
									"field_name": "col2",
									"field_type": "Long",
								},
								{
									"field_name": "pk1",
									"field_type": "Long",
								},
							},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(getSearchIndexCheckMap()),
					captureSearchIndexState(resourceId, &applied),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateCheck:  searchIndexImportStateCheck(&applied, searchIndexTopLevelStoreNotReturned),
				ImportStateVerify: false,
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_name": "${alicloud_ots_instance.default.name}",
					"table_name":    "${alicloud_ots_table.default.table_name}",
					"index_name":    "${var.name}_search_index",
					"time_to_live":  "-1",
					"schema": []map[string]interface{}{
						{
							"field_schema": []map[string]interface{}{
								{
									"field_name": "col1",
									"field_type": "Text",
									"analyzer":   "Split",
									"store":      "true",
								},
								{
									"field_name": "col2",
									"field_type": "Long",
								},
								{
									"field_name": "pk1",
									"field_type": "Long",
								},
							},
						},
					},
				}),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"schema": []map[string]interface{}{
						{
							"field_schema": []map[string]interface{}{
								{
									"field_name": "pk1",
									"field_type": "Long",
								},
								{
									"field_name": "col1",
									"field_type": "Text",
									"analyzer":   "Split",
									"store":      "true",
								},
								{
									"field_name": "col2",
									"field_type": "Long",
								},
							},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(getSearchIndexCheckMap()),
				),
			},
		},
	})
}

func resourceOtsSearchIndexFieldOrderConfigDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

resource "alicloud_ots_instance" "default" {
  name          = "tf-${var.name}"
  description   = "${var.name}"
  accessed_by   = "Any"
  instance_type = "%s"
  tags = {
    Created = "TF"
    For     = "acceptance test"
  }
}

resource "alicloud_ots_table" "default" {
  instance_name = alicloud_ots_instance.default.name
  table_name    = var.name
  primary_key {
    name = "pk1"
    type = "Integer"
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
  max_version  = 1
}
`, name, string(ots.OtsCapacity))
}

func TestAccAliCloudOtsSearchIndex_invalidSchema(t *testing.T) {
	var v *tablestore.DescribeTableResponse

	resourceId := "alicloud_ots_search_index.default"
	serviceFunc := func() interface{} {
		return ots.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}
	rc := acctest.ResourceCheckInit(resourceId, &v, serviceFunc)

	name := fmt.Sprintf("testAcc%d", tfacctest.RandIntRange(10000, 99999))

	var applied *sdkterraform.InstanceState

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckWithRegions(t, false, connectivity.OtsCapacityNoSupportedRegions)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: resourceSearchIndexInvalidConfig(name, `field_schema {
      field_name = "col1"
      field_type = "Nested"
    }`),
				ExpectError: regexp.MustCompile("requires sub_field_schemas"),
			},
			{
				Config: resourceSearchIndexInvalidConfig(name, `field_schema {
      field_name = "col1"
      field_type = "Nested"
      is_array   = true
      sub_field_schemas {
        field_name = "nested_long"
        field_type = "Long"
      }
    }`),
				ExpectError: regexp.MustCompile("does not support is_array"),
			},
			{
				Config: resourceSearchIndexInvalidConfig(name, `field_schema {
      field_name = "col1"
      field_type = "Nested"
      store      = true
      sub_field_schemas {
        field_name = "nested_long"
        field_type = "Long"
      }
    }`),
				ExpectError: regexp.MustCompile("does not support store"),
			},
			{
				Config: resourceSearchIndexInvalidConfig(name, `field_schema {
      field_name = "col1"
      field_type = "Nested"
      analyzer   = "SingleWord"
      sub_field_schemas {
        field_name = "nested_long"
        field_type = "Long"
      }
    }`),
				ExpectError: regexp.MustCompile("does not support analyzer"),
			},
			{
				Config: resourceSearchIndexInvalidConfig(name, `field_schema {
      field_name          = "col1"
      field_type          = "Nested"
      enable_sort_and_agg = true
      sub_field_schemas {
        field_name = "nested_long"
        field_type = "Long"
      }
    }`),
				ExpectError: regexp.MustCompile("does not support enable_sort_and_agg"),
			},
			{
				Config: resourceSearchIndexInvalidConfig(name, `field_schema {
      field_name        = "col1"
      field_type        = "Nested"
      is_virtual_field  = true
      source_field_name = "col2"
      sub_field_schemas {
        field_name = "nested_long"
        field_type = "Long"
      }
    }`),
				ExpectError: regexp.MustCompile("cannot be a virtual field"),
			},
			{
				Config: resourceSearchIndexInvalidConfig(name, `field_schema {
      field_name = "col1"
      field_type = "Keyword"
      sub_field_schemas {
        field_name = "nested_long"
        field_type = "Long"
      }
    }`),
				ExpectError: regexp.MustCompile("does not support sub_field_schemas"),
			},
			{
				Config: resourceSearchIndexInvalidConfig(name, `field_schema {
      field_name = "col1"
      field_type = "Vector"
    }`),
				ExpectError: regexp.MustCompile("requires vector_options"),
			},
			{
				Config: resourceSearchIndexInvalidConfig(name, `field_schema {
      field_name = "col1"
      field_type = "Vector"
      vector_options {
        dimension   = 4
        metric_type = "cosine"
      }
    }`),
				ExpectError: regexp.MustCompile("vector_options requires data_type"),
			},
			{
				Config: resourceSearchIndexInvalidConfig(name, `field_schema {
      field_name = "col1"
      field_type = "Vector"
      vector_options {
        data_type   = "float32"
        metric_type = "cosine"
      }
    }`),
				ExpectError: regexp.MustCompile("vector_options requires dimension"),
			},
			{
				Config: resourceSearchIndexInvalidConfig(name, `field_schema {
      field_name = "col1"
      field_type = "Vector"
      vector_options {
        data_type = "float32"
        dimension = 4
      }
    }`),
				ExpectError: regexp.MustCompile("vector_options requires metric_type"),
			},
			{
				Config: resourceSearchIndexInvalidConfig(name, `field_schema {
      field_name = "col1"
      field_type = "Keyword"
      vector_options {
        data_type   = "float32"
        dimension   = 4
        metric_type = "cosine"
      }
    }`),
				ExpectError: regexp.MustCompile("does not support vector_options"),
			},
			{
				Config: resourceSearchIndexInvalidConfig(name, `field_schema {
      field_name        = "col1"
      field_type        = "Keyword"
      source_field_name = "col2"
    }`),
				ExpectError: regexp.MustCompile("is_virtual_field is not true"),
			},
			{
				Config: resourceSearchIndexInvalidConfig(name, `field_schema {
      field_name       = "col1"
      field_type       = "Keyword"
      is_virtual_field = true
    }`),
				ExpectError: regexp.MustCompile("requires source_field_name"),
			},
			{
				Config: resourceSearchIndexInvalidConfig(name, `field_schema {
      field_name = "col1"
      field_type = "Vector"
      is_array   = true
      vector_options {
        data_type   = "float32"
        dimension   = 4
        metric_type = "cosine"
      }
    }`),
				ExpectError: regexp.MustCompile(`field "col1" of type Vector does not support is_array`),
			},
			{
				Config: resourceSearchIndexInvalidConfig(name, `field_schema {
      field_name = "col1"
      field_type = "FuzzyKeyword"
      is_array   = true
    }`),
				ExpectError: regexp.MustCompile(`field "col1" of type FuzzyKeyword does not support is_array`),
			},
			{
				Config: resourceSearchIndexInvalidConfig(name, `field_schema {
      field_name = "col1"
      field_type = "Vector"
      analyzer   = "SingleWord"
      vector_options {
        data_type   = "float32"
        dimension   = 4
        metric_type = "cosine"
      }
    }`),
				ExpectError: regexp.MustCompile(`field "col1" of type Vector does not support analyzer`),
			},
			{
				Config: resourceSearchIndexInvalidConfig(name, `field_schema {
      field_name = "col1"
      field_type = "FuzzyKeyword"
      analyzer   = "SingleWord"
    }`),
				ExpectError: regexp.MustCompile(`field "col1" of type FuzzyKeyword does not support analyzer`),
			},
			{
				Config: resourceSearchIndexInvalidConfig(name, `field_schema {
      field_name = "col1"
      field_type = "Nested"
      sub_field_schemas {
        field_name = "nested_long"
        field_type = "Long"
        analyzer   = "SingleWord"
      }
    }`),
				ExpectError: regexp.MustCompile(`field "nested_long" of type Long does not support analyzer`),
			},
			{
				Config: resourceSearchIndexInvalidConfig(name, `field_schema {
      field_name        = "col1"
      field_type        = "Keyword"
      date_format_list  = ["epoch_millis"]
    }`),
				ExpectError: regexp.MustCompile(`field "col1" of type Keyword does not support date_format_list`),
			},
			{
				Config: resourceSearchIndexInvalidConfig(name, `field_schema {
      field_name        = "col1"
      field_type        = "Keyword"
      index             = false
      is_virtual_field  = true
      source_field_name = "col2"
    }`),
				ExpectError: regexp.MustCompile(`virtual field "col1" requires index = true`),
			},
			{
				Config: resourceSearchIndexInvalidConfig(name, `field_schema {
      field_name = "col1"
      field_type = "Unknown"
    }`),
				ExpectError: regexp.MustCompile("expected .*field_type to be one of"),
			},
			{
				Config: resourceSearchIndexInvalidConfig(name, `field_schema {
      field_name          = "col1"
      field_type          = "Text"
      enable_sort_and_agg = true
    }`),
				ExpectError: regexp.MustCompile("does not support enable_sort_and_agg"),
			},
			{
				Config: resourceSearchIndexInvalidConfig(name, `field_schema {
      field_name          = "col1"
      field_type          = "FuzzyKeyword"
      enable_sort_and_agg = true
    }`),
				ExpectError: regexp.MustCompile("does not support enable_sort_and_agg"),
			},
			{
				Config: resourceSearchIndexInvalidConfig(name, `field_schema {
      field_name = "col1"
      field_type = "Nested"
      sub_field_schemas {
        field_name = "nested_parent"
        field_type = "Nested"
      }
    }`),
				ExpectError: regexp.MustCompile("expected .*field_type to be one of"),
			},
			{
				Config: resourceSearchIndexInvalidConfig(name, `field_schema {
      field_name = "col1"
      field_type = "Nested"
      sub_field_schemas {
        field_name = "nested_vector"
        field_type = "Vector"
      }
    }`),
				ExpectError: regexp.MustCompile("expected .*field_type to be one of"),
			},
			{
				Config: resourceSearchIndexInvalidConfig(name, `field_schema {
      field_name = "legacy_date"
      field_type = "Date"
    }`),
				ExpectError: regexp.MustCompile(`(?s)OTSParameterInvalid.*date_format must be set when field\[legacy_date\] field_type is 'date'`),
			},
			{
				Config: resourceSearchIndexInvalidConfig(name, `field_schema {
      field_name = "col1"
      field_type = "Nested"
      sub_field_schemas {
        field_name = "nested_date"
        field_type = "Date"
      }
    }`),
				ExpectError: regexp.MustCompile(`field "nested_date" of type Date requires date_format_list`),
			},
			{
				Config: resourceSearchIndexInvalidConfig(name, `field_schema {
      field_name = "col1"
      field_type = "Nested"
      sub_field_schemas {
        field_name       = "nested_virtual"
        field_type       = "Keyword"
        is_virtual_field = true
      }
    }`),
				ExpectError: regexp.MustCompile(`virtual field "nested_virtual" requires source_field_name`),
			},
			{
				Config: resourceSearchIndexInvalidConfig(name, `field_schema {
      field_name = "col1"
      field_type = "Nested"
      sub_field_schemas {
        field_name        = "nested_keyword"
        field_type        = "Keyword"
        source_field_name = "pk2"
      }
    }`),
				ExpectError: regexp.MustCompile(`field "nested_keyword" sets source_field_name but is_virtual_field is not true`),
			},
			{
				Config: resourceSearchIndexInvalidConfig(name, `field_schema {
      field_name = "col1"
      field_type = "Nested"
      sub_field_schemas {
        field_name          = "nested_text"
        field_type          = "Text"
        enable_sort_and_agg = true
      }
    }`),
				ExpectError: regexp.MustCompile(`field "nested_text" of type Text does not support enable_sort_and_agg`),
			},
			{
				Config: resourceSearchIndexInvalidConfig(name, `field_schema {
      field_name = "col1"
      field_type = "Nested"
      sub_field_schemas {
        field_name          = "nested_fuzzy"
        field_type          = "FuzzyKeyword"
        enable_sort_and_agg = true
      }
    }`),
				ExpectError: regexp.MustCompile(`field "nested_fuzzy" of type FuzzyKeyword does not support enable_sort_and_agg`),
			},
			{
				// A Nested field cannot be combined with an explicit index_sort;
				// the guard runs after the table existence check on the service
				// side and before CreateSearchIndex.
				Config: resourceSearchIndexInvalidConfig(name, `field_schema {
      field_name = "nested_col"
      field_type = "Nested"
      sub_field_schemas {
        field_name = "nested_long"
        field_type = "Long"
      }
    }
    index_sort {
      sorter {
        sorter_type = "PrimaryKeySort"
        order       = "Asc"
      }
    }`),
				ExpectError: regexp.MustCompile("does not support index_sort"),
			},
			{
				// Close with a valid configuration so the post-test destroy has
				// a planable configuration; the rejected steps above apply
				// before their error and would otherwise leave that check
				// pointing at an invalid schema.
				Config: resourceOtsSearchIndexConfig(name),
				Check:  captureSearchIndexState(resourceId, &applied),
			},
			{
				ResourceName:     resourceId,
				ImportState:      true,
				ImportStateCheck: searchIndexImportStateCheck(&applied, searchIndexTopLevelStoreNotReturned),
			},
		},
	})
}

func resourceSearchIndexInvalidConfig(name, fieldSchema string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

resource "alicloud_ots_instance" "default" {
  name          = "tf-${var.name}"
  description   = "${var.name}"
  accessed_by   = "Any"
  instance_type = "%s"
  tags = {
    Created = "TF"
    For     = "acceptance test"
  }
}

resource "alicloud_ots_table" "default" {
  instance_name = alicloud_ots_instance.default.name
  table_name    = var.name
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
  defined_column {
    name = "legacy_date"
    type = "String"
  }
  time_to_live = -1
  max_version  = 1
}

resource "alicloud_ots_search_index" "default" {
  instance_name = alicloud_ots_instance.default.name
  table_name    = alicloud_ots_table.default.table_name
  index_name    = "%s_search_index"
  time_to_live  = -1

  schema {
    field_schema {
      field_name = "pk1"
      field_type = "Long"
    }
    %s
  }
}
`, name, string(ots.OtsCapacity), name, fieldSchema)
}

var searchIndexBasicMap = map[string]string{}

func getSearchIndexCheckMap() map[string]string {
	accMap := make(map[string]string)
	accMap["index_id"] = acctest.CHECKSET
	accMap["create_time"] = acctest.CHECKSET
	accMap["sync_phase"] = acctest.CHECKSET
	accMap["current_sync_timestamp"] = acctest.CHECKSET
	accMap["instance_name"] = acctest.CHECKSET
	accMap["table_name"] = acctest.CHECKSET
	accMap["index_name"] = acctest.CHECKSET
	accMap["time_to_live"] = acctest.CHECKSET

	return accMap
}

// captureSearchIndexState stores the state produced by the configuration so the
// import step can compare the imported state with it attribute by attribute.
func captureSearchIndexState(resourceId string, target **sdkterraform.InstanceState) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceId]
		if !ok || rs.Primary == nil {
			return fmt.Errorf("%s is missing from state", resourceId)
		}
		attributes := make(map[string]string, len(rs.Primary.Attributes))
		for key, value := range rs.Primary.Attributes {
			attributes[key] = value
		}
		*target = &sdkterraform.InstanceState{ID: rs.Primary.ID, Attributes: attributes}
		return nil
	}
}

// searchIndexImportStateCheck compares the imported state with the state the
// configuration produced. ImportStateVerifyIgnore matches flatmap keys by
// prefix, which cannot address a single leaf inside a set element that carries a
// dynamic hash, so this check decodes both states through the resource schema
// and compares every attribute instead of ignoring the whole schema.
func searchIndexImportStateCheck(applied **sdkterraform.InstanceState, storeOverrides map[string]bool) resource.ImportStateCheckFunc {
	return func(states []*terraform.InstanceState) error {
		if applied == nil || *applied == nil {
			return fmt.Errorf("the applied state was not captured before the import step")
		}
		if len(states) != 1 {
			return fmt.Errorf("expected one imported state, got %d", len(states))
		}
		imported := &sdkterraform.InstanceState{ID: states[0].ID, Attributes: states[0].Attributes}
		return compareSearchIndexStates(*applied, imported, storeOverrides)
	}
}

func compareSearchIndexStates(applied, imported *sdkterraform.InstanceState, storeOverrides map[string]bool) error {
	// The sync statistics move forward between the apply and the import step;
	// the schema is decoded and compared separately below.
	ignored := map[string]bool{
		"sync_phase":             true,
		"current_sync_timestamp": true,
	}

	keys := make(map[string]struct{}, len(applied.Attributes)+len(imported.Attributes))
	for key := range applied.Attributes {
		keys[key] = struct{}{}
	}
	for key := range imported.Attributes {
		keys[key] = struct{}{}
	}

	var diffs []string
	for key := range keys {
		if ignored[key] || key == "schema.#" || strings.HasPrefix(key, "schema.") {
			continue
		}
		if applied.Attributes[key] != imported.Attributes[key] {
			diffs = append(diffs, fmt.Sprintf("%s: applied %q, imported %q", key, applied.Attributes[key], imported.Attributes[key]))
		}
	}

	appliedSchema, err := decodeSearchIndexSchema(applied)
	if err != nil {
		return fmt.Errorf("decode applied state: %w", err)
	}
	importedSchema, err := decodeSearchIndexSchema(imported)
	if err != nil {
		return fmt.Errorf("decode imported state: %w", err)
	}
	diffs = append(diffs, compareSearchIndexSchemas(appliedSchema, importedSchema, storeOverrides)...)

	sort.Strings(diffs)
	if len(diffs) > 0 {
		return fmt.Errorf("imported state does not match the applied state:\n%s", strings.Join(diffs, "\n"))
	}
	return nil
}

type searchIndexFieldModel struct {
	name        string
	fieldType   string
	isArray     bool
	index       bool
	analyzer    string
	sortAndAgg  bool
	store       bool
	dateFormats []string
	virtual     bool
	sourceField string
	vector      map[string]interface{}
	subFields   map[string]*searchIndexFieldModel
}

type searchIndexSorterModel struct {
	sorterType string
	order      string
	fieldName  string
}

type searchIndexSchemaModel struct {
	fields        map[string]*searchIndexFieldModel
	routingFields []string
	sorters       []searchIndexSorterModel
}

func decodeSearchIndexSchema(state *sdkterraform.InstanceState) (*searchIndexSchemaModel, error) {
	data := ots.ResourceOtsSearchIndex().Data(&sdkterraform.InstanceState{ID: state.ID, Attributes: state.Attributes})
	return decodeSearchIndexResourceData(data)
}

func decodeSearchIndexResourceData(data *sdkSchema.ResourceData) (*searchIndexSchemaModel, error) {
	blocks, _ := data.Get("schema").(*sdkSchema.Set)
	if blocks == nil || blocks.Len() != 1 {
		count := 0
		if blocks != nil {
			count = blocks.Len()
		}
		return nil, fmt.Errorf("schema has %d blocks, want exactly one", count)
	}

	block := blocks.List()[0].(map[string]interface{})
	model := &searchIndexSchemaModel{fields: make(map[string]*searchIndexFieldModel)}

	fieldList, ok := block["field_schema"].([]interface{})
	if !ok {
		return nil, fmt.Errorf("field_schema is missing from the schema block")
	}
	for _, item := range fieldList {
		field, err := decodeSearchIndexField(item)
		if err != nil {
			return nil, err
		}
		model.fields[field.name] = field
	}

	if settings, ok := block["index_setting"].(*sdkSchema.Set); ok && settings != nil && settings.Len() > 0 {
		setting := settings.List()[0].(map[string]interface{})
		if routing, ok := setting["routing_fields"].([]interface{}); ok {
			for _, value := range routing {
				model.routingFields = append(model.routingFields, value.(string))
			}
		}
	}

	if sorts, ok := block["index_sort"].(*sdkSchema.Set); ok && sorts != nil && sorts.Len() > 0 {
		sortBlock := sorts.List()[0].(map[string]interface{})
		if sorters, ok := sortBlock["sorter"].([]interface{}); ok {
			for _, item := range sorters {
				sorter := item.(map[string]interface{})
				model.sorters = append(model.sorters, searchIndexSorterModel{
					sorterType: sorter["sorter_type"].(string),
					order:      sorter["order"].(string),
					fieldName:  sorter["field_name"].(string),
				})
			}
		}
	}

	return model, nil
}

func decodeSearchIndexField(value interface{}) (*searchIndexFieldModel, error) {
	member, ok := value.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("unexpected field_schema value of type %T", value)
	}
	field := &searchIndexFieldModel{
		name:        member["field_name"].(string),
		fieldType:   member["field_type"].(string),
		isArray:     member["is_array"].(bool),
		index:       member["index"].(bool),
		analyzer:    member["analyzer"].(string),
		sortAndAgg:  member["enable_sort_and_agg"].(bool),
		store:       member["store"].(bool),
		virtual:     member["is_virtual_field"].(bool),
		sourceField: member["source_field_name"].(string),
	}
	if formats, ok := member["date_format_list"].(*sdkSchema.Set); ok && formats != nil {
		for _, value := range formats.List() {
			field.dateFormats = append(field.dateFormats, value.(string))
		}
		sort.Strings(field.dateFormats)
	}
	if options, ok := member["vector_options"].([]interface{}); ok && len(options) > 0 {
		field.vector = options[0].(map[string]interface{})
	}
	if subFields, ok := member["sub_field_schemas"].(*sdkSchema.Set); ok && subFields != nil && subFields.Len() > 0 {
		field.subFields = make(map[string]*searchIndexFieldModel, subFields.Len())
		for _, item := range subFields.List() {
			subField, err := decodeSearchIndexField(item)
			if err != nil {
				return nil, err
			}
			field.subFields[subField.name] = subField
		}
	}
	return field, nil
}

func compareSearchIndexSchemas(applied, imported *searchIndexSchemaModel, storeOverrides map[string]bool) []string {
	var diffs []string
	if len(applied.fields) != len(imported.fields) {
		diffs = append(diffs, fmt.Sprintf("schema.field_schema: applied %d fields, imported %d", len(applied.fields), len(imported.fields)))
	}
	for name, appliedField := range applied.fields {
		importedField, ok := imported.fields[name]
		if !ok {
			diffs = append(diffs, fmt.Sprintf("schema.field_schema.%s: missing after import", name))
			continue
		}
		diffs = append(diffs, compareSearchIndexFields(name, appliedField, importedField, storeOverrides)...)
	}
	for name := range imported.fields {
		if _, ok := applied.fields[name]; !ok {
			diffs = append(diffs, fmt.Sprintf("schema.field_schema.%s: unexpected after import", name))
		}
	}

	appliedRouting := append([]string(nil), applied.routingFields...)
	importedRouting := append([]string(nil), imported.routingFields...)
	sort.Strings(appliedRouting)
	sort.Strings(importedRouting)
	if strings.Join(appliedRouting, ",") != strings.Join(importedRouting, ",") {
		diffs = append(diffs, fmt.Sprintf("schema.index_setting.routing_fields: applied %v, imported %v", applied.routingFields, imported.routingFields))
	}

	appliedSorters := normalizeSearchIndexSorters(applied.sorters)
	importedSorters := normalizeSearchIndexSorters(imported.sorters)
	if !searchIndexSortersEqual(appliedSorters, importedSorters) {
		diffs = append(diffs, fmt.Sprintf("schema.index_sort.sorter: applied %v, imported %v", appliedSorters, importedSorters))
	}
	return diffs
}

func compareSearchIndexFields(path string, applied, imported *searchIndexFieldModel, storeOverrides map[string]bool) []string {
	var diffs []string
	report := func(attribute string, appliedValue, importedValue interface{}) {
		diffs = append(diffs, fmt.Sprintf("schema.field_schema.%s.%s: applied %v, imported %v", path, attribute, appliedValue, importedValue))
	}

	if applied.fieldType != imported.fieldType {
		report("field_type", applied.fieldType, imported.fieldType)
	}
	if applied.isArray != imported.isArray {
		report("is_array", applied.isArray, imported.isArray)
	}
	if applied.index != imported.index {
		report("index", applied.index, imported.index)
	}
	if applied.analyzer != imported.analyzer {
		report("analyzer", applied.analyzer, imported.analyzer)
	}
	if applied.sortAndAgg != imported.sortAndAgg {
		report("enable_sort_and_agg", applied.sortAndAgg, imported.sortAndAgg)
	}

	// The service does not return `store` for every field: an omitted argument
	// is materialized as false in the applied state, while the described index
	// reports the wire default or drops the value. The request shape is pinned
	// by the public-create fixture tests; the caller lists the requested values
	// that differ from the applied state.
	expectedStore := applied.store
	if override, ok := storeOverrides[path]; ok {
		expectedStore = override
	}
	if expectedStore != imported.store {
		report("store", expectedStore, imported.store)
	}

	if applied.virtual != imported.virtual {
		report("is_virtual_field", applied.virtual, imported.virtual)
	}
	if applied.sourceField != imported.sourceField {
		report("source_field_name", applied.sourceField, imported.sourceField)
	}
	if strings.Join(applied.dateFormats, ",") != strings.Join(imported.dateFormats, ",") {
		report("date_format_list", applied.dateFormats, imported.dateFormats)
	}

	if (applied.vector == nil) != (imported.vector == nil) {
		report("vector_options", applied.vector != nil, imported.vector != nil)
	} else if applied.vector != nil {
		for _, attribute := range []string{"data_type", "dimension", "metric_type"} {
			if fmt.Sprint(applied.vector[attribute]) != fmt.Sprint(imported.vector[attribute]) {
				report("vector_options."+attribute, applied.vector[attribute], imported.vector[attribute])
			}
		}
	}

	if len(applied.subFields) != len(imported.subFields) {
		report("sub_field_schemas", len(applied.subFields), len(imported.subFields))
	}
	for name, appliedSubField := range applied.subFields {
		importedSubField, ok := imported.subFields[name]
		if !ok {
			diffs = append(diffs, fmt.Sprintf("schema.field_schema.%s.sub_field_schemas.%s: missing after import", path, name))
			continue
		}
		diffs = append(diffs, compareSearchIndexFields(path+"."+name, appliedSubField, importedSubField, storeOverrides)...)
	}
	for name := range imported.subFields {
		if _, ok := applied.subFields[name]; !ok {
			diffs = append(diffs, fmt.Sprintf("schema.field_schema.%s.sub_field_schemas.%s: unexpected after import", path, name))
		}
	}

	return diffs
}

// normalizeSearchIndexSorters treats the implicit primary-key sorter the
// service materializes for an index without a custom sort as the same empty
// configuration. The FieldSort mode is not read back by DescribeSearchIndex and
// is already dropped when the sorter is decoded.
func normalizeSearchIndexSorters(sorters []searchIndexSorterModel) []searchIndexSorterModel {
	if len(sorters) == 1 && sorters[0].sorterType == string(ots.OtsSearchPrimaryKeySort) && sorters[0].order == string(ots.OtsSearchSortOrderAsc) {
		return nil
	}
	return sorters
}

func searchIndexSortersEqual(applied, imported []searchIndexSorterModel) bool {
	if len(applied) != len(imported) {
		return false
	}
	for index := range applied {
		if applied[index] != imported[index] {
			return false
		}
	}
	return true
}

type searchIndexRawOverride struct {
	unknownFieldNames map[string]bool
	unknownStores     map[string]bool
}

type searchIndexFixtureOptions struct {
	rawOverride    *searchIndexRawOverride
	describeSchema *otsprotocol.IndexSchema
	skipRawConfig  bool
}

// searchIndexFixture serves the endpoints the public resource create and read
// paths call. CreateSearchIndex returns the request the SDK actually sent and
// DescribeSearchIndex echoes the created schema unless a test overrides it.
type searchIndexFixture struct {
	t              *testing.T
	server         *httptest.Server
	client         *connectivity.AliyunClient
	mutex          sync.Mutex
	instanceName   string
	tableName      string
	indexName      string
	createRequest  *otsprotocol.CreateSearchIndexRequest
	describeSchema *otsprotocol.IndexSchema
	handlerErr     error
}

func newSearchIndexFixture(t *testing.T, instanceName, tableName, indexName string, describeSchema *otsprotocol.IndexSchema) *searchIndexFixture {
	t.Helper()
	t.Setenv("NO_PROXY", "127.0.0.1,localhost")

	fixture := &searchIndexFixture{
		t:              t,
		instanceName:   instanceName,
		tableName:      tableName,
		indexName:      indexName,
		describeSchema: describeSchema,
	}
	fixture.server = httptest.NewServer(http.HandlerFunc(fixture.handle))
	t.Cleanup(fixture.server.Close)

	credential, err := credentials.NewCredential(new(credentials.Config).
		SetType("access_key").
		SetAccessKeyId("test-access-key").
		SetAccessKeySecret("test-secret-key"))
	if err != nil {
		t.Fatalf("failed to build the test credential: %v", err)
	}
	endpoints := &sync.Map{}
	endpoints.Store("tablestore", strings.TrimPrefix(fixture.server.URL, "http://"))
	client, err := connectivity.NewClient(&connectivity.Config{
		AccessKey:            "test-access-key",
		SecretKey:            "test-secret-key",
		Region:               connectivity.Region("cn-beijing"),
		RegionId:             "cn-beijing",
		AccountType:          "Domestic",
		Protocol:             "HTTP",
		ClientReadTimeout:    5000,
		ClientConnectTimeout: 5000,
		SkipRegionValidation: true,
		OtsEndpoint:          fixture.server.URL,
		Endpoints:            endpoints,
		Credential:           credential,
	})
	if err != nil {
		t.Fatalf("failed to build the test client: %v", err)
	}
	fixture.client = client
	return fixture
}

func (f *searchIndexFixture) handle(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/v2/openapi/getinstance":
		w.Header().Set("Content-Type", "application/json")
		_, _ = io.WriteString(w, fmt.Sprintf(`{"InstanceName":%q,"InstanceStatus":"Normal"}`, f.instanceName))
	case "/CreateSearchIndex":
		body, err := io.ReadAll(r.Body)
		if err != nil {
			f.fail(w, err)
			return
		}
		request := new(otsprotocol.CreateSearchIndexRequest)
		if err := request.XXX_Unmarshal(body); err != nil {
			f.fail(w, err)
			return
		}
		f.mutex.Lock()
		f.createRequest = request
		f.mutex.Unlock()
		f.writeProto(w, &otsprotocol.CreateSearchIndexResponse{})
	case "/DescribeSearchIndex":
		f.mutex.Lock()
		request := f.createRequest
		describeSchema := f.describeSchema
		f.mutex.Unlock()
		if describeSchema == nil {
			if request == nil {
				f.fail(w, fmt.Errorf("CreateSearchIndex was not called before DescribeSearchIndex"))
				return
			}
			describeSchema = request.Schema
		}
		phase := otsprotocol.SyncPhase_INCR
		var timeToLive *int32
		if request != nil {
			timeToLive = request.TimeToLive
		}
		f.writeProto(w, &otsprotocol.DescribeSearchIndexResponse{
			Schema:     describeSchema,
			SyncStat:   &otsprotocol.SyncStat{SyncPhase: &phase, CurrentSyncTimestamp: proto.Int64(0)},
			CreateTime: proto.Int64(1600000000),
			TimeToLive: timeToLive,
		})
	case "/DescribeTable":
		f.writeProto(w, describeTableFixtureResponse(map[string]interface{}{
			"table_name":   f.tableName,
			"time_to_live": -1,
			"max_version":  1,
			"primary_key": []interface{}{
				map[string]interface{}{"name": "pk1", "type": "Integer"},
			},
		}))
	default:
		http.NotFound(w, r)
	}
}

func (f *searchIndexFixture) fail(w http.ResponseWriter, err error) {
	f.mutex.Lock()
	f.handlerErr = err
	f.mutex.Unlock()
	http.Error(w, err.Error(), http.StatusInternalServerError)
}

func (f *searchIndexFixture) writeProto(w http.ResponseWriter, message proto.Message) {
	payload, err := proto.Marshal(message)
	if err != nil {
		f.fail(w, err)
		return
	}
	w.Header().Set("Content-Type", "application/x-protobuf")
	_, _ = w.Write(payload)
}

// create runs the public Resource.Diff and Resource.Apply path, including the
// raw configuration a real plan provides, and returns the decoded
// CreateSearchIndex request the SDK sent over HTTP.
func (f *searchIndexFixture) create(config map[string]interface{}, options searchIndexFixtureOptions) (*otsprotocol.CreateSearchIndexRequest, *sdkterraform.InstanceState, error) {
	resourceOtsSearchIndex := ots.ResourceOtsSearchIndex()
	diff, err := resourceOtsSearchIndex.Diff(context.Background(), nil, sdkterraform.NewResourceConfigRaw(config), f.client)
	if err != nil {
		return nil, nil, err
	}
	if diff == nil {
		return nil, nil, fmt.Errorf("no diff produced for the test configuration")
	}
	if !options.skipRawConfig {
		diff.RawConfig = searchIndexRawConfigValue(f.t, resourceOtsSearchIndex.CoreConfigSchema().ImpliedType(), config, options.rawOverride)
	}
	state, diags := resourceOtsSearchIndex.Apply(context.Background(), nil, diff, f.client)
	if diags.HasError() {
		return nil, nil, fmt.Errorf("apply failed: %v", diags)
	}
	f.mutex.Lock()
	request := f.createRequest
	f.mutex.Unlock()
	if request == nil {
		if f.handlerErr != nil {
			return nil, nil, f.handlerErr
		}
		return nil, nil, fmt.Errorf("no CreateSearchIndex request reached the endpoint")
	}
	return request, state, nil
}

// read runs the public read path on an empty state, which is the import path
// that flattens the schema returned by DescribeSearchIndex.
func (f *searchIndexFixture) read() (*sdkSchema.ResourceData, error) {
	id := fmt.Sprintf("%s:%s:%s:Search", f.instanceName, f.tableName, f.indexName)
	data := ots.ResourceOtsSearchIndex().Data(&sdkterraform.InstanceState{ID: id})
	if err := ots.ResourceOtsSearchIndex().Read(data, f.client); err != nil {
		return nil, err
	}
	return data, nil
}

func searchIndexFixtureBaseConfig(fields []interface{}) map[string]interface{} {
	return map[string]interface{}{
		"instance_name": "test-instance",
		"table_name":    "test_table",
		"index_name":    "test_index",
		"time_to_live":  -1,
		"schema": []interface{}{map[string]interface{}{
			"field_schema": fields,
		}},
	}
}

func searchIndexFixtureVectorField(store interface{}) map[string]interface{} {
	field := map[string]interface{}{
		"field_name": "vector_col",
		"field_type": "Vector",
		"vector_options": []interface{}{map[string]interface{}{
			"data_type":   "float32",
			"dimension":   4,
			"metric_type": "cosine",
		}},
	}
	if store != nil {
		field["store"] = store
	}
	return field
}

func searchIndexFixtureVirtualField(store interface{}) map[string]interface{} {
	field := map[string]interface{}{
		"field_name":        "virtual_col",
		"field_type":        "Keyword",
		"is_virtual_field":  true,
		"source_field_name": "pk2",
	}
	if store != nil {
		field["store"] = store
	}
	return field
}

func searchIndexFixtureNestedField(subFields ...map[string]interface{}) map[string]interface{} {
	values := make([]interface{}, 0, len(subFields))
	for _, subField := range subFields {
		values = append(values, subField)
	}
	return map[string]interface{}{
		"field_name":        "nested_col",
		"field_type":        "Nested",
		"sub_field_schemas": values,
	}
}

func searchIndexFixtureSubField(store interface{}) map[string]interface{} {
	subField := map[string]interface{}{
		"field_name": "nested_sub",
		"field_type": "Keyword",
	}
	if store != nil {
		subField["store"] = store
	}
	return subField
}

func searchIndexFixturePbField(schema *otsprotocol.IndexSchema, name string) *otsprotocol.FieldSchema {
	var walk func([]*otsprotocol.FieldSchema) *otsprotocol.FieldSchema
	walk = func(fields []*otsprotocol.FieldSchema) *otsprotocol.FieldSchema {
		for _, field := range fields {
			if field.GetFieldName() == name {
				return field
			}
			if found := walk(field.GetFieldSchemas()); found != nil {
				return found
			}
		}
		return nil
	}
	return walk(schema.GetFieldSchemas())
}

// TestUnitOtsSearchIndexStoreRequestShape pins the wire value of the optional
// `store` argument through the public create path and the actual
// CreateSearchIndex HTTP request.
func TestUnitOtsSearchIndexStoreRequestShape(t *testing.T) {
	cases := []struct {
		name      string
		field     map[string]interface{}
		fieldName string
		wantStore bool
	}{
		{"vector omitted", searchIndexFixtureVectorField(nil), "vector_col", true},
		{"vector explicit false", searchIndexFixtureVectorField(false), "vector_col", false},
		{"vector explicit true", searchIndexFixtureVectorField(true), "vector_col", true},
		{"fuzzy keyword omitted", map[string]interface{}{"field_name": "fuzzy_col", "field_type": "FuzzyKeyword"}, "fuzzy_col", true},
		{"fuzzy keyword explicit false", map[string]interface{}{"field_name": "fuzzy_col", "field_type": "FuzzyKeyword", "store": false}, "fuzzy_col", false},
		{"fuzzy keyword explicit true", map[string]interface{}{"field_name": "fuzzy_col", "field_type": "FuzzyKeyword", "store": true}, "fuzzy_col", true},
		{"virtual field omitted", searchIndexFixtureVirtualField(nil), "virtual_col", true},
		{"virtual field explicit false", searchIndexFixtureVirtualField(false), "virtual_col", false},
		{"virtual field explicit true", searchIndexFixtureVirtualField(true), "virtual_col", true},
		{"nested sub-field omitted", searchIndexFixtureNestedField(searchIndexFixtureSubField(nil)), "nested_sub", true},
		{"nested sub-field explicit false", searchIndexFixtureNestedField(searchIndexFixtureSubField(false)), "nested_sub", false},
		{"nested sub-field explicit true", searchIndexFixtureNestedField(searchIndexFixtureSubField(true)), "nested_sub", true},
		{"legacy long omitted", map[string]interface{}{"field_name": "legacy_long", "field_type": "Long"}, "legacy_long", false},
		{"legacy long explicit false", map[string]interface{}{"field_name": "legacy_long", "field_type": "Long", "store": false}, "legacy_long", false},
		{"legacy long explicit true", map[string]interface{}{"field_name": "legacy_long", "field_type": "Long", "store": true}, "legacy_long", true},
		{"top-level text omitted", map[string]interface{}{"field_name": "text_col", "field_type": "Text"}, "text_col", false},
		{"top-level text explicit false", map[string]interface{}{"field_name": "text_col", "field_type": "Text", "analyzer": "Split", "store": false}, "text_col", false},
		{"top-level text explicit true", map[string]interface{}{"field_name": "text_col", "field_type": "Text", "analyzer": "Split", "store": true}, "text_col", true},
	}

	for _, testCase := range cases {
		t.Run(testCase.name, func(t *testing.T) {
			fixture := newSearchIndexFixture(t, "test-instance", "test_table", "test_index", nil)
			request, _, err := fixture.create(searchIndexFixtureBaseConfig([]interface{}{testCase.field}), searchIndexFixtureOptions{})
			if err != nil {
				t.Fatalf("create through the public resource failed: %v", err)
			}
			field := searchIndexFixturePbField(request.GetSchema(), testCase.fieldName)
			if field == nil {
				t.Fatalf("%s was not present in the CreateSearchIndex request", testCase.fieldName)
			}
			if field.Store == nil {
				t.Fatalf("%s store was not serialized, want %v", testCase.fieldName, testCase.wantStore)
			}
			if field.GetStore() != testCase.wantStore {
				t.Fatalf("%s wire store is %v, want %v", testCase.fieldName, field.GetStore(), testCase.wantStore)
			}
		})
	}
}

func TestUnitOtsSearchIndexFieldListStoreCorrelation(t *testing.T) {
	t.Run("known names", func(t *testing.T) {
		fixture := newSearchIndexFixture(t, "test-instance", "test_table", "test_index", nil)
		request, _, err := fixture.create(searchIndexFixtureBaseConfig([]interface{}{
			map[string]interface{}{"field_name": "text_col", "field_type": "Text", "analyzer": "Split", "store": true},
			map[string]interface{}{"field_name": "long_col", "field_type": "Long", "store": false},
			searchIndexFixtureVectorField(false),
		}), searchIndexFixtureOptions{})
		if err != nil {
			t.Fatalf("create through the public resource failed: %v", err)
		}
		for name, wantStore := range map[string]bool{"text_col": true, "long_col": false, "vector_col": false} {
			field := searchIndexFixturePbField(request.GetSchema(), name)
			if field == nil {
				t.Fatalf("%s was not present in the CreateSearchIndex request", name)
			}
			if field.Store == nil || field.GetStore() != wantStore {
				t.Fatalf("%s wire store is %v, want %v", name, field.Store, wantStore)
			}
		}
	})

	t.Run("unique unknown name", func(t *testing.T) {
		fixture := newSearchIndexFixture(t, "test-instance", "test_table", "test_index", nil)
		request, _, err := fixture.create(searchIndexFixtureBaseConfig([]interface{}{
			map[string]interface{}{"field_name": "keyword_col", "field_type": "Keyword", "store": false},
		}), searchIndexFixtureOptions{rawOverride: &searchIndexRawOverride{
			unknownFieldNames: map[string]bool{"keyword_col": true},
		}})
		if err != nil {
			t.Fatalf("unique unknown field name failed: %v", err)
		}
		field := searchIndexFixturePbField(request.GetSchema(), "keyword_col")
		if field == nil || field.Store == nil || field.GetStore() {
			t.Fatalf("keyword_col explicit false store was not preserved: %#v", field)
		}
	})

	for _, testCase := range []struct {
		name        string
		firstStore  interface{}
		secondStore interface{}
		wantFirst   bool
		wantSecond  bool
	}{
		{name: "consistent unknown names", firstStore: false, secondStore: false},
		{name: "false and true", firstStore: false, secondStore: true, wantSecond: true},
		{name: "true and false", firstStore: true, secondStore: false, wantFirst: true},
		{name: "omitted and false", secondStore: false},
		{name: "omitted and true", secondStore: true, wantSecond: true},
	} {
		t.Run(testCase.name, func(t *testing.T) {
			first := map[string]interface{}{"field_name": "first_col", "field_type": "Keyword"}
			second := map[string]interface{}{"field_name": "second_col", "field_type": "Keyword"}
			if testCase.firstStore != nil {
				first["store"] = testCase.firstStore
			}
			if testCase.secondStore != nil {
				second["store"] = testCase.secondStore
			}
			fixture := newSearchIndexFixture(t, "test-instance", "test_table", "test_index", nil)
			request, _, err := fixture.create(searchIndexFixtureBaseConfig([]interface{}{first, second}), searchIndexFixtureOptions{
				rawOverride: &searchIndexRawOverride{unknownFieldNames: map[string]bool{"first_col": true, "second_col": true}},
			})
			if err != nil {
				t.Fatalf("unknown-name store configuration failed: %v", err)
			}
			for name, wantStore := range map[string]bool{"first_col": testCase.wantFirst, "second_col": testCase.wantSecond} {
				field := searchIndexFixturePbField(request.GetSchema(), name)
				if field == nil || field.Store == nil || field.GetStore() != wantStore {
					t.Fatalf("%s wire store does not match the configured value, want %v: %#v", name, wantStore, field)
				}
			}
		})
	}
}

// TestUnitOtsSearchIndexComputedStoreRequestShape covers a raw configuration
// whose store value is unknown at plan time: the materialized value must reach
// the request.
func TestUnitOtsSearchIndexComputedStoreRequestShape(t *testing.T) {
	cases := []struct {
		name      string
		field     map[string]interface{}
		fieldName string
		wantStore bool
	}{
		{"legacy true", map[string]interface{}{"field_name": "legacy_long", "field_type": "Long", "store": true}, "legacy_long", true},
		{"legacy false", map[string]interface{}{"field_name": "legacy_long", "field_type": "Long", "store": false}, "legacy_long", false},
		{"vector true", searchIndexFixtureVectorField(true), "vector_col", true},
		{"vector false", searchIndexFixtureVectorField(false), "vector_col", false},
	}

	for _, testCase := range cases {
		t.Run(testCase.name, func(t *testing.T) {
			fixture := newSearchIndexFixture(t, "test-instance", "test_table", "test_index", nil)
			request, _, err := fixture.create(searchIndexFixtureBaseConfig([]interface{}{testCase.field}), searchIndexFixtureOptions{
				rawOverride: &searchIndexRawOverride{unknownStores: map[string]bool{testCase.fieldName: true}},
			})
			if err != nil {
				t.Fatalf("create through the public resource failed: %v", err)
			}
			field := searchIndexFixturePbField(request.GetSchema(), testCase.fieldName)
			if field == nil || field.Store == nil {
				t.Fatalf("%s store was not serialized", testCase.fieldName)
			}
			if field.GetStore() != testCase.wantStore {
				t.Fatalf("%s wire store is %v, want %v", testCase.fieldName, field.GetStore(), testCase.wantStore)
			}
		})
	}
}

// TestUnitOtsSearchIndexStoreWithoutRawConfig covers the apply path when the
// raw configuration is unavailable.
func TestUnitOtsSearchIndexStoreWithoutRawConfig(t *testing.T) {
	cases := []struct {
		name      string
		field     map[string]interface{}
		fieldName string
		wantStore bool
	}{
		{"explicit true", map[string]interface{}{"field_name": "col1", "field_type": "Text", "analyzer": "Split", "store": true}, "col1", true},
		{"omitted", map[string]interface{}{"field_name": "col2", "field_type": "Long"}, "col2", false},
	}

	for _, testCase := range cases {
		t.Run(testCase.name, func(t *testing.T) {
			fixture := newSearchIndexFixture(t, "test-instance", "test_table", "test_index", nil)
			request, _, err := fixture.create(searchIndexFixtureBaseConfig([]interface{}{testCase.field}), searchIndexFixtureOptions{skipRawConfig: true})
			if err != nil {
				t.Fatalf("create through the public resource failed: %v", err)
			}
			field := searchIndexFixturePbField(request.GetSchema(), testCase.fieldName)
			if field == nil || field.Store == nil || field.GetStore() != testCase.wantStore {
				t.Fatalf("%s wire store is %v, want %v", testCase.fieldName, field.Store, testCase.wantStore)
			}
		})
	}
}

func TestUnitOtsSearchIndexUnknownSubFieldStore(t *testing.T) {
	cases := []struct {
		name         string
		store        interface{}
		unknownStore bool
		wantStore    bool
	}{
		{name: "omitted", wantStore: true},
		{name: "explicit true", store: true, wantStore: true},
		{name: "explicit false", store: false, wantStore: false},
		{name: "computed true", store: true, unknownStore: true, wantStore: true},
		{name: "computed false", store: false, unknownStore: true, wantStore: false},
	}
	for _, testCase := range cases {
		t.Run(testCase.name, func(t *testing.T) {
			override := &searchIndexRawOverride{unknownFieldNames: map[string]bool{"nested_sub": true}}
			if testCase.unknownStore {
				override.unknownStores = map[string]bool{"nested_sub": true}
			}
			fixture := newSearchIndexFixture(t, "test-instance", "test_table", "test_index", nil)
			request, _, err := fixture.create(searchIndexFixtureBaseConfig([]interface{}{
				searchIndexFixtureNestedField(searchIndexFixtureSubField(testCase.store)),
			}), searchIndexFixtureOptions{rawOverride: override})
			if err != nil {
				t.Fatalf("create through the public resource failed: %v", err)
			}
			field := searchIndexFixturePbField(request.GetSchema(), "nested_sub")
			if field == nil || field.Store == nil || field.GetStore() != testCase.wantStore {
				t.Fatalf("nested_sub wire store is %v, want %v", field.Store, testCase.wantStore)
			}
		})
	}
}

func TestUnitOtsSearchIndexUnknownSubFieldStoreCorrelation(t *testing.T) {
	cases := []struct {
		name          string
		firstStore    interface{}
		secondStore   interface{}
		unknownStores bool
		wantStore     bool
		wantError     bool
	}{
		{name: "all true", firstStore: true, secondStore: true, wantStore: true},
		{name: "all false", firstStore: false, secondStore: false, wantStore: false},
		{name: "all computed false", firstStore: false, secondStore: false, unknownStores: true, wantStore: false},
		{name: "mixed true and false", firstStore: false, secondStore: true, wantError: true},
		{name: "mixed omitted and false", secondStore: false, wantError: true},
	}
	for _, testCase := range cases {
		t.Run(testCase.name, func(t *testing.T) {
			first := map[string]interface{}{"field_name": "nested_sub", "field_type": "Keyword"}
			second := map[string]interface{}{"field_name": "other_sub", "field_type": "Long"}
			if testCase.firstStore != nil {
				first["store"] = testCase.firstStore
			}
			if testCase.secondStore != nil {
				second["store"] = testCase.secondStore
			}
			override := &searchIndexRawOverride{unknownFieldNames: map[string]bool{"nested_sub": true, "other_sub": true}}
			if testCase.unknownStores {
				override.unknownStores = map[string]bool{"nested_sub": true, "other_sub": true}
			}

			fixture := newSearchIndexFixture(t, "test-instance", "test_table", "test_index", nil)
			request, _, err := fixture.create(searchIndexFixtureBaseConfig([]interface{}{
				searchIndexFixtureNestedField(first, second),
			}), searchIndexFixtureOptions{rawOverride: override})
			if testCase.wantError {
				if err == nil || !strings.Contains(err.Error(), "cannot reliably correlate sub_field_schemas") {
					t.Fatalf("ambiguous store configuration returned %v", err)
				}
				return
			}
			if err != nil {
				t.Fatalf("consistent store configuration failed: %v", err)
			}
			for _, name := range []string{"nested_sub", "other_sub"} {
				field := searchIndexFixturePbField(request.GetSchema(), name)
				if field == nil || field.Store == nil || field.GetStore() != testCase.wantStore {
					t.Fatalf("%s wire store is %v, want %v", name, field.Store, testCase.wantStore)
				}
			}
		})
	}
}

func TestUnitOtsSearchIndexBooleanRequestShape(t *testing.T) {
	cases := []struct {
		name        string
		field       map[string]interface{}
		wantIsArray bool
		wantSortAgg bool
	}{
		{"long omitted", map[string]interface{}{"field_name": "col", "field_type": "Long"}, false, false},
		{"long true", map[string]interface{}{"field_name": "col", "field_type": "Long", "is_array": true, "enable_sort_and_agg": true}, true, true},
	}

	for _, testCase := range cases {
		t.Run(testCase.name, func(t *testing.T) {
			fixture := newSearchIndexFixture(t, "test-instance", "test_table", "test_index", nil)
			request, _, err := fixture.create(searchIndexFixtureBaseConfig([]interface{}{testCase.field}), searchIndexFixtureOptions{})
			if err != nil {
				t.Fatalf("create through the public resource failed: %v", err)
			}
			field := searchIndexFixturePbField(request.GetSchema(), "col")
			if field == nil {
				t.Fatal("col was not present in the CreateSearchIndex request")
			}
			if field.IsArray == nil || field.GetIsArray() != testCase.wantIsArray {
				t.Fatalf("is_array is %v, want %v", field.IsArray, testCase.wantIsArray)
			}
			if field.SortAndAgg == nil || field.GetSortAndAgg() != testCase.wantSortAgg {
				t.Fatalf("enable_sort_and_agg is %v, want %v", field.SortAndAgg, testCase.wantSortAgg)
			}
		})
	}
}

func TestUnitOtsSearchIndexDateFormatCompatibility(t *testing.T) {
	t.Run("top-level date may omit formats", func(t *testing.T) {
		fixture := newSearchIndexFixture(t, "test-instance", "test_table", "test_index", nil)
		request, _, err := fixture.create(searchIndexFixtureBaseConfig([]interface{}{
			map[string]interface{}{"field_name": "date_col", "field_type": "Date"},
		}), searchIndexFixtureOptions{})
		if err != nil {
			t.Fatalf("create through the public resource failed: %v", err)
		}
		field := searchIndexFixturePbField(request.GetSchema(), "date_col")
		if field == nil {
			t.Fatal("date_col was not present in the CreateSearchIndex request")
		}
		if len(field.DateFormats) != 0 {
			t.Fatalf("date_col formats = %v, want none", field.DateFormats)
		}
	})

	t.Run("date sub-field requires formats", func(t *testing.T) {
		fixture := newSearchIndexFixture(t, "test-instance", "test_table", "test_index", nil)
		_, _, err := fixture.create(searchIndexFixtureBaseConfig([]interface{}{
			searchIndexFixtureNestedField(map[string]interface{}{
				"field_name": "nested_date",
				"field_type": "Date",
			}),
		}), searchIndexFixtureOptions{})
		if err == nil || !strings.Contains(err.Error(), `field "nested_date" of type Date requires date_format_list`) {
			t.Fatalf("date sub-field without formats returned %v", err)
		}
	})
}

func TestUnitOtsSearchIndexRejectsExtraSchemaBlocks(t *testing.T) {
	config := searchIndexFixtureBaseConfig([]interface{}{
		map[string]interface{}{"field_name": "pk1", "field_type": "Long"},
	})
	config["schema"] = []interface{}{
		map[string]interface{}{"field_schema": []interface{}{
			map[string]interface{}{"field_name": "pk1", "field_type": "Long"},
		}},
		map[string]interface{}{"field_schema": []interface{}{
			map[string]interface{}{"field_name": "pk2", "field_type": "Text"},
		}},
	}

	fixture := newSearchIndexFixture(t, "test-instance", "test_table", "test_index", nil)
	_, _, err := fixture.create(config, searchIndexFixtureOptions{})
	if err == nil || !strings.Contains(err.Error(), "only one schema block is supported") {
		t.Fatalf("multiple schema blocks must be rejected, got %v", err)
	}
}

// TestUnitOtsSearchIndexReadFlattensSchema covers the import-time flattening of
// every supported field kind through the public Read path.
func TestUnitOtsSearchIndexReadFlattensSchema(t *testing.T) {
	dimension := int32(128)
	describeSchema := &otsprotocol.IndexSchema{
		IndexSetting: &otsprotocol.IndexSetting{RoutingFields: []string{"pk1"}},
		IndexSort: &otsprotocol.Sort{Sorter: []*otsprotocol.Sorter{
			{PkSort: &otsprotocol.PrimaryKeySort{Order: otsprotocol.SortOrder_SORT_ORDER_ASC.Enum()}},
			{FieldSort: &otsprotocol.FieldSort{FieldName: proto.String("long_col"), Order: otsprotocol.SortOrder_SORT_ORDER_DESC.Enum()}},
		}},
		FieldSchemas: []*otsprotocol.FieldSchema{
			{
				FieldName: proto.String("text_col"), FieldType: otsprotocol.FieldType_TEXT.Enum(),
				Index: proto.Bool(true), Store: proto.Bool(true), Analyzer: proto.String("split"),
			},
			{
				FieldName: proto.String("nested_col"), FieldType: otsprotocol.FieldType_NESTED.Enum(),
				FieldSchemas: []*otsprotocol.FieldSchema{{
					FieldName: proto.String("date_col"), FieldType: otsprotocol.FieldType_DATE.Enum(),
					DateFormats: []string{"epoch_second"},
				}},
			},
			{
				FieldName: proto.String("vector_col"), FieldType: otsprotocol.FieldType_VECTOR.Enum(),
				VectorOptions: &otsprotocol.VectorOptions{
					DataType:   otsprotocol.VectorDataType_VD_FLOAT_32.Enum(),
					Dimension:  &dimension,
					MetricType: otsprotocol.VectorMetricType_VM_COSINE.Enum(),
				},
			},
			{
				FieldName: proto.String("virtual_col"), FieldType: otsprotocol.FieldType_KEYWORD.Enum(),
				IsVirtualField: proto.Bool(true), SourceFieldNames: []string{"pk2"},
			},
		},
	}

	fixture := newSearchIndexFixture(t, "test-instance", "test_table", "test_index", describeSchema)
	data, err := fixture.read()
	if err != nil {
		t.Fatalf("read through the public resource failed: %v", err)
	}
	model, err := decodeSearchIndexResourceData(data)
	if err != nil {
		t.Fatalf("decode flattened schema: %v", err)
	}

	if len(model.fields) != 4 {
		t.Fatalf("flattened %d fields, want 4", len(model.fields))
	}
	if text := model.fields["text_col"]; text == nil || text.analyzer != "Split" || !text.store {
		t.Fatalf("text_col was not flattened as expected: %#v", text)
	}
	nested := model.fields["nested_col"]
	if nested == nil || nested.subFields["date_col"] == nil || strings.Join(nested.subFields["date_col"].dateFormats, ",") != "epoch_second" {
		t.Fatalf("nested date sub-field was not flattened: %#v", nested)
	}
	vector := model.fields["vector_col"]
	if vector == nil || vector.vector == nil || fmt.Sprint(vector.vector["dimension"]) != "128" {
		t.Fatalf("vector options were not flattened: %#v", vector)
	}
	virtual := model.fields["virtual_col"]
	if virtual == nil || !virtual.virtual || virtual.sourceField != "pk2" {
		t.Fatalf("virtual field was not flattened: %#v", virtual)
	}
	if strings.Join(model.routingFields, ",") != "pk1" {
		t.Fatalf("routing fields are %v, want [pk1]", model.routingFields)
	}
	if len(model.sorters) != 2 || model.sorters[0].sorterType != string(ots.OtsSearchPrimaryKeySort) || model.sorters[1].fieldName != "long_col" {
		t.Fatalf("sorters were not flattened: %#v", model.sorters)
	}
}

func TestUnitOtsSearchIndexFieldOrderRefreshDiff(t *testing.T) {
	fields := func() []interface{} {
		return []interface{}{
			map[string]interface{}{"field_name": "pk1", "field_type": "Long", "store": true},
			map[string]interface{}{"field_name": "col1", "field_type": "Keyword", "store": false},
			map[string]interface{}{"field_name": "col2", "field_type": "Text", "analyzer": "Split"},
		}
	}
	for _, testCase := range []struct {
		name         string
		remoteOrder  []int
		configOrder  []int
		remoteChange string
		configChange string
		changed      bool
	}{
		{name: "same response", remoteOrder: []int{0, 1, 2}, configOrder: []int{0, 1, 2}},
		{name: "reverse response", remoteOrder: []int{2, 1, 0}, configOrder: []int{0, 1, 2}},
		{name: "rotated response and config", remoteOrder: []int{2, 0, 1}, configOrder: []int{1, 2, 0}},
		{name: "remote member added", remoteOrder: []int{0, 1, 2}, configOrder: []int{0, 1, 2}, remoteChange: "add", changed: true},
		{name: "remote member removed", remoteOrder: []int{0, 1, 2}, configOrder: []int{0, 1, 2}, remoteChange: "remove", changed: true},
		{name: "remote value changed", remoteOrder: []int{0, 1, 2}, configOrder: []int{0, 1, 2}, remoteChange: "value", changed: true},
		{name: "configured reorder", remoteOrder: []int{0, 1, 2}, configOrder: []int{2, 1, 0}},
		{name: "configured member added", remoteOrder: []int{0, 1, 2}, configOrder: []int{0, 1, 2}, configChange: "add", changed: true},
		{name: "configured member removed", remoteOrder: []int{0, 1, 2}, configOrder: []int{0, 1, 2}, configChange: "remove", changed: true},
		{name: "configured value changed", remoteOrder: []int{0, 1, 2}, configOrder: []int{0, 1, 2}, configChange: "value", changed: true},
	} {
		t.Run(testCase.name, func(t *testing.T) {
			fixture := newSearchIndexFixture(t, "test-instance", "test_table", "test_index", nil)
			request, _, err := fixture.create(searchIndexFixtureBaseConfig(fields()), searchIndexFixtureOptions{})
			if err != nil {
				t.Fatalf("create through the public resource failed: %v", err)
			}

			remoteSchema := proto.Clone(request.GetSchema()).(*otsprotocol.IndexSchema)
			remoteFields := make([]*otsprotocol.FieldSchema, 0, len(testCase.remoteOrder))
			for _, index := range testCase.remoteOrder {
				remoteFields = append(remoteFields, remoteSchema.FieldSchemas[index])
			}
			remoteSchema.FieldSchemas = remoteFields
			switch testCase.remoteChange {
			case "add":
				added := proto.Clone(remoteSchema.FieldSchemas[0]).(*otsprotocol.FieldSchema)
				added.FieldName = proto.String("remote_added")
				remoteSchema.FieldSchemas = append(remoteSchema.FieldSchemas, added)
			case "remove":
				remoteSchema.FieldSchemas = remoteSchema.FieldSchemas[:len(remoteSchema.FieldSchemas)-1]
			case "value":
				remoteSchema.FieldSchemas[0].FieldType = otsprotocol.FieldType_DOUBLE.Enum()
			}
			fixture.mutex.Lock()
			fixture.describeSchema = remoteSchema
			fixture.mutex.Unlock()

			imported, err := fixture.read()
			if err != nil {
				t.Fatalf("read through the public resource failed: %v", err)
			}
			plannedFields := fields()
			orderedFields := make([]interface{}, 0, len(testCase.configOrder))
			for _, index := range testCase.configOrder {
				orderedFields = append(orderedFields, plannedFields[index])
			}
			plannedFields = orderedFields
			switch testCase.configChange {
			case "add":
				plannedFields = append(plannedFields, map[string]interface{}{"field_name": "config_added", "field_type": "Long"})
			case "remove":
				plannedFields = plannedFields[:len(plannedFields)-1]
			case "value":
				plannedFields[1].(map[string]interface{})["field_type"] = "Double"
			}
			config := searchIndexFixtureBaseConfig(plannedFields)
			diff, err := ots.ResourceOtsSearchIndex().Diff(context.Background(), imported.State(), sdkterraform.NewResourceConfigRaw(config), fixture.client)
			if err != nil {
				t.Fatalf("diff after import failed: %v", err)
			}
			changed := diff != nil && !diff.Empty()
			if changed != testCase.changed {
				t.Fatalf("changed=%v want=%v diff=%#v", changed, testCase.changed, diff)
			}
			if changed && !diff.RequiresNew() {
				t.Fatalf("field membership or value change must require replacement: %#v", diff)
			}
		})
	}
}

func TestUnitOtsSearchIndexFieldOrderReadsLegacyHashState(t *testing.T) {
	fields := []interface{}{
		map[string]interface{}{"field_name": "pk1", "field_type": "Long", "store": true},
		map[string]interface{}{"field_name": "col1", "field_type": "Keyword", "store": false},
		map[string]interface{}{"field_name": "col2", "field_type": "Text", "analyzer": "Split"},
	}
	config := searchIndexFixtureBaseConfig(fields)
	fixture := newSearchIndexFixture(t, "test-instance", "test_table", "test_index", nil)
	_, applied, err := fixture.create(config, searchIndexFixtureOptions{})
	if err != nil {
		t.Fatalf("create through the public resource failed: %v", err)
	}
	legacyState := searchIndexLegacyFieldListState(t, applied, fields)
	model, err := decodeSearchIndexSchema(legacyState)
	if err != nil {
		t.Fatalf("decode legacy list state through the set schema: %v", err)
	}
	if len(model.fields) != len(fields) {
		t.Fatalf("decoded %d legacy fields, want %d", len(model.fields), len(fields))
	}
	resourceOtsSearchIndex := ots.ResourceOtsSearchIndex()
	refreshed, diags := resourceOtsSearchIndex.RefreshWithoutUpgrade(context.Background(), legacyState, fixture.client)
	if diags.HasError() {
		t.Fatalf("refresh legacy list state failed: %v", diags)
	}
	diff, err := resourceOtsSearchIndex.Diff(context.Background(), refreshed, sdkterraform.NewResourceConfigRaw(config), fixture.client)
	if err != nil {
		t.Fatalf("diff legacy list state failed: %v", err)
	}
	if diff != nil && !diff.Empty() {
		t.Fatalf("legacy list state must converge without replacement: %#v", diff)
	}
}

func searchIndexLegacyFieldListState(t *testing.T, state *sdkterraform.InstanceState, fields []interface{}) *sdkterraform.InstanceState {
	t.Helper()
	resourceSchema := ots.ResourceOtsSearchIndex().Schema["schema"]
	model, err := decodeSearchIndexSchema(state)
	if err != nil {
		t.Fatalf("decode applied state before constructing legacy state: %v", err)
	}
	data := ots.ResourceOtsSearchIndex().Data(state)
	blocks := data.Get("schema").(*sdkSchema.Set)
	block := blocks.List()[0].(map[string]interface{})
	legacyHash := sdkSchema.HashResource(resourceSchema.Elem.(*sdkSchema.Resource))(block)
	currentHash := resourceSchema.Set(block)
	if legacyHash == currentHash {
		t.Fatalf("legacy and current schema hashes are both %d; the test does not exercise migration", legacyHash)
	}
	if len(model.fields) != len(fields) {
		t.Fatalf("applied state has %d fields, want %d", len(model.fields), len(fields))
	}

	currentPrefix := fmt.Sprintf("schema.%d.", currentHash)
	legacyPrefix := fmt.Sprintf("schema.%d.", legacyHash)
	attributes := make(map[string]string, len(state.Attributes))
	rewritten := 0
	for key, value := range state.Attributes {
		if strings.HasPrefix(key, currentPrefix) {
			attributes[legacyPrefix+strings.TrimPrefix(key, currentPrefix)] = value
			rewritten++
			continue
		}
		attributes[key] = value
	}
	if rewritten == 0 {
		t.Fatalf("no state attributes matched current schema prefix %q", currentPrefix)
	}
	for key := range attributes {
		if strings.HasPrefix(key, currentPrefix) {
			t.Fatalf("current schema key %q remains in legacy state", key)
		}
	}
	for index, value := range fields {
		field := value.(map[string]interface{})
		key := fmt.Sprintf("%sfield_schema.%d.field_name", legacyPrefix, index)
		if attributes[key] != field["field_name"] {
			t.Fatalf("legacy numeric field key %q=%q, want %q", key, attributes[key], field["field_name"])
		}
	}
	return &sdkterraform.InstanceState{ID: state.ID, Attributes: attributes, Meta: state.Meta}
}

// searchIndexRawConfigValue builds the raw configuration a real plan produces,
// including the synthetic unknowns of optional arguments.
func searchIndexRawConfigValue(t *testing.T, impliedType cty.Type, config map[string]interface{}, override *searchIndexRawOverride) cty.Value {
	t.Helper()
	return searchIndexRawTypedValue(t, impliedType, config, override)
}

func searchIndexRawTypedValue(t *testing.T, typ cty.Type, value interface{}, override *searchIndexRawOverride) cty.Value {
	t.Helper()

	switch {
	case typ.IsObjectType():
		configured, ok := value.(map[string]interface{})
		if !ok {
			t.Fatalf("expected an object value for %s, got %T", typ.FriendlyName(), value)
		}
		fieldName, _ := configured["field_name"].(string)
		attributes := make(map[string]cty.Value, len(typ.AttributeTypes()))
		for name, attributeType := range typ.AttributeTypes() {
			attributeValue, present := configured[name]
			if !present {
				attributes[name] = cty.NullVal(attributeType)
				continue
			}
			if override != nil && name == "field_name" && override.unknownFieldNames[fieldName] {
				attributes[name] = cty.UnknownVal(attributeType)
				continue
			}
			if override != nil && name == "store" && override.unknownStores[fieldName] {
				attributes[name] = cty.UnknownVal(attributeType)
				continue
			}
			attributes[name] = searchIndexRawTypedValue(t, attributeType, attributeValue, override)
		}
		return cty.ObjectVal(attributes)
	case typ.IsListType() || typ.IsSetType():
		items, ok := value.([]interface{})
		if !ok {
			t.Fatalf("expected a collection value for %s, got %T", typ.FriendlyName(), value)
		}
		elementType := typ.ElementType()
		elements := make([]cty.Value, 0, len(items))
		for _, item := range items {
			elements = append(elements, searchIndexRawTypedValue(t, elementType, item, override))
		}
		switch {
		case typ.IsSetType() && len(elements) == 0:
			return cty.SetValEmpty(elementType)
		case typ.IsSetType():
			return cty.SetVal(elements)
		case len(elements) == 0:
			return cty.ListValEmpty(elementType)
		default:
			return cty.ListVal(elements)
		}
	case typ.Equals(cty.String):
		stringValue, ok := value.(string)
		if !ok {
			t.Fatalf("expected a string value for %s, got %T", typ.FriendlyName(), value)
		}
		return cty.StringVal(stringValue)
	case typ.Equals(cty.Bool):
		boolValue, ok := value.(bool)
		if !ok {
			t.Fatalf("expected a bool value for %s, got %T", typ.FriendlyName(), value)
		}
		return cty.BoolVal(boolValue)
	case typ.Equals(cty.Number):
		intValue, ok := value.(int)
		if !ok {
			t.Fatalf("expected an int value for %s, got %T", typ.FriendlyName(), value)
		}
		return cty.NumberIntVal(int64(intValue))
	default:
		t.Fatalf("unsupported raw config type %s", typ.FriendlyName())
		return cty.NilVal
	}
}

// TestUnitOtsSearchIndexReadRejectsUnsupportedSorter covers the public read
// path when DescribeSearchIndex returns a sorter the provider does not support.
func TestUnitOtsSearchIndexReadRejectsUnsupportedSorter(t *testing.T) {
	describeSchema := &otsprotocol.IndexSchema{
		IndexSetting: &otsprotocol.IndexSetting{},
		IndexSort: &otsprotocol.Sort{Sorter: []*otsprotocol.Sorter{
			{ScoreSort: &otsprotocol.ScoreSort{}},
		}},
		FieldSchemas: []*otsprotocol.FieldSchema{
			{FieldName: proto.String("pk1"), FieldType: otsprotocol.FieldType_LONG.Enum()},
		},
	}
	fixture := newSearchIndexFixture(t, "test-instance", "test_table", "test_index", describeSchema)
	if _, err := fixture.read(); err == nil || !strings.Contains(err.Error(), "unknown index sort type") {
		t.Fatalf("unsupported sorter must be rejected, got %v", err)
	}
}
