// collection-order: bypass schema.field_schema -- pre-existing attribute outside this change
// collection-order: bypass schema.index_setting.routing_fields -- pre-existing attribute outside this change
// collection-order: bypass schema.index_sort.sorter -- pre-existing attribute outside this change
package ots_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/aliyun/aliyun-tablestore-go-sdk/tablestore"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/ots"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
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

func TestAccAliCloudOtsSearchIndex_complete(t *testing.T) {
	var v *tablestore.DescribeSearchIndexResponse

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
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: resourceOtsSearchIndexCompleteConfig(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(getSearchIndexCheckMap()),
					resource.TestCheckTypeSetElemNestedAttrs(resourceId, "schema.*", map[string]string{
						"field_schema.#":                              "8",
						"field_schema.2.field_type":                   "Vector",
						"field_schema.2.vector_options.#":             "1",
						"field_schema.2.vector_options.0.metric_type": "cosine",
						"field_schema.3.field_type":                   "FuzzyKeyword",
						"field_schema.4.field_type":                   "Nested",
						"field_schema.4.sub_field_schemas.#":          "6",
						"field_schema.5.date_format_list.#":           "2",
						"field_schema.6.is_virtual_field":             "true",
						"field_schema.6.source_field_name":            "pk2",
						"field_schema.7.field_type":                   "Text",
					}),
					checkSearchIndexSchemaFields(),
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

func TestAccAliCloudOtsSearchIndex_invalidSchema(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckWithRegions(t, false, connectivity.OtsCapacityNoSupportedRegions)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: resourceOtsSearchIndexInvalidConfig(`
    field_schema {
      field_name = "col1"
      field_type = "Vector"
    }`),
				ExpectError: regexp.MustCompile("requires vector_options"),
			},
			{
				Config: resourceOtsSearchIndexInvalidConfig(`
    field_schema {
      field_name       = "col1"
      field_type       = "Keyword"
      is_virtual_field = true
    }`),
				ExpectError: regexp.MustCompile("requires source_field_name"),
			},
			{
				Config: resourceOtsSearchIndexInvalidConfig(`
    field_schema {
      field_name = "col1"
      field_type = "Nested"
    }`),
				ExpectError: regexp.MustCompile("requires sub_field_schemas"),
			},
			{
				Config: resourceOtsSearchIndexInvalidConfig(`
    field_schema {
      field_name = "col1"
      field_type = "Date"
    }`),
				ExpectError: regexp.MustCompile("requires date_format_list"),
			},
			{
				Config: resourceOtsSearchIndexInvalidConfig(`
    field_schema {
      field_name        = "col1"
      field_type        = "Keyword"
      source_field_name = "col2"
    }`),
				ExpectError: regexp.MustCompile("is_virtual_field is not true"),
			},
			{
				Config: resourceOtsSearchIndexInvalidConfig(`
    field_schema {
      field_name = "col1"
      field_type = "Unknown"
    }`),
				ExpectError: regexp.MustCompile("expected .*field_type to be one of"),
			},
			{
				Config: resourceOtsSearchIndexInvalidConfig(`
    field_schema {
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
				Config: resourceOtsSearchIndexInvalidConfig(`
    field_schema {
      field_name          = "col1"
      field_type          = "Text"
      enable_sort_and_agg = true
    }`),
				ExpectError: regexp.MustCompile("does not support enable_sort_and_agg"),
			},
			{
				Config: resourceOtsSearchIndexInvalidConfig(`
    field_schema {
      field_name          = "col1"
      field_type          = "FuzzyKeyword"
      enable_sort_and_agg = true
    }`),
				ExpectError: regexp.MustCompile("does not support enable_sort_and_agg"),
			},
			{
				Config: resourceOtsSearchIndexInvalidConfig(`
    field_schema {
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
				Config: resourceOtsSearchIndexInvalidConfig(`
    field_schema {
      field_name       = "col1"
      field_type       = "Nested"
      is_virtual_field = true
      sub_field_schemas {
        field_name = "nested_long"
        field_type = "Long"
      }
    }`),
				ExpectError: regexp.MustCompile("cannot be a virtual field"),
			},
			{
				Config: resourceOtsSearchIndexInvalidConfig(`
    field_schema {
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
				Config: resourceOtsSearchIndexInvalidConfig(`
    field_schema {
      field_name = "col1"
      field_type = "Nested"
      index      = false
      sub_field_schemas {
        field_name = "nested_long"
        field_type = "Long"
      }
    }`),
				ExpectError: regexp.MustCompile("does not support disabling index"),
			},
			{
				Config: resourceOtsSearchIndexInvalidConfig(`
    field_schema {
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
				Config: resourceOtsSearchIndexInvalidConfig(`
    field_schema {
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
				Config: resourceOtsSearchIndexInvalidConfig(`
    field_schema {
      field_name = "col1"
      field_type = "Nested"
      sub_field_schemas {
        field_name = "nested_vector"
        field_type = "Vector"
      }
    }`),
				ExpectError: regexp.MustCompile("field_type to be one of"),
			},
			{
				Config: resourceOtsSearchIndexInvalidConfig(`
    field_schema {
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
				Config: resourceOtsSearchIndexInvalidConfig(`
    field_schema {
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
				Config: resourceOtsSearchIndexInvalidConfig(`
    field_schema {
      field_name = "col1"
      field_type = "Nested"
      sub_field_schemas {
        field_name        = "nested_virtual"
        field_type        = "Keyword"
        source_field_name = "pk2"
      }
    }`),
				ExpectError: regexp.MustCompile(`field "nested_virtual" sets source_field_name but is_virtual_field is not true`),
			},
			{
				Config: resourceOtsSearchIndexInvalidConfig(`
    field_schema {
      field_name = "col1"
      field_type = "Nested"
      sub_field_schemas {
        field_name          = "nested_text"
        field_type          = "Text"
        enable_sort_and_agg = true
      }
    }`),
				ExpectError: regexp.MustCompile(`field "nested_text" of type TEXT does not support enable_sort_and_agg`),
			},
			{
				Config: resourceOtsSearchIndexInvalidConfig(`
    field_schema {
      field_name = "col1"
      field_type = "Nested"
      sub_field_schemas {
        field_name          = "nested_fuzzy"
        field_type          = "FuzzyKeyword"
        enable_sort_and_agg = true
      }
    }`),
				ExpectError: regexp.MustCompile(`field "nested_fuzzy" of type FUZZY_KEYWORD does not support enable_sort_and_agg`),
			},
		},
	})
}

func resourceOtsSearchIndexCompleteConfig(name string) string {
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
        field_name       = "nested_date"
        field_type       = "Date"
        date_format_list = ["datetime", "epoch_millis"]
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

    index_setting {
      routing_fields = ["pk1", "pk2"]
    }
  }
}
`, name, string(ots.OtsCapacity), name)
}

func resourceOtsSearchIndexInvalidConfig(fieldSchema string) string {
	return fmt.Sprintf(`
resource "alicloud_ots_search_index" "default" {
  instance_name = "tf-test-invalid"
  table_name    = "tf_test_invalid"
  index_name    = "tf_test_invalid_index"
  time_to_live  = -1

  schema {
%s
  }
}
`, fieldSchema)
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

		nestedLong := fields["nested_long"]
		if nestedLong.EnableSortAndAgg == nil || !*nestedLong.EnableSortAndAgg {
			return fmt.Errorf("nested_long enable_sort_and_agg is not enabled")
		}

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

		nestedDateFormats := fields["nested_date"].DateFormats
		expectedNestedDateFormats := map[string]struct{}{
			"datetime":     {},
			"epoch_millis": {},
		}
		if len(nestedDateFormats) != len(expectedNestedDateFormats) {
			return fmt.Errorf("nested_date date formats are %v", nestedDateFormats)
		}
		for _, dateFormat := range nestedDateFormats {
			if _, ok := expectedNestedDateFormats[dateFormat]; !ok {
				return fmt.Errorf("nested_date has unexpected date format %s", dateFormat)
			}
		}

		nestedVirtual := fields["nested_virtual"]
		if nestedVirtual.IsVirtualField == nil || !*nestedVirtual.IsVirtualField {
			return fmt.Errorf("nested_virtual is not marked as a virtual field")
		}
		if len(nestedVirtual.SourceFieldNames) != 1 || nestedVirtual.SourceFieldNames[0] != "pk2" {
			return fmt.Errorf("nested_virtual source field names are %v", nestedVirtual.SourceFieldNames)
		}

		return nil
	}
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
