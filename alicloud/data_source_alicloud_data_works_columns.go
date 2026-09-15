package alicloud

import (
	"fmt"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

// dataSourceType values supported by GetMetaTableColumn DataSourceType parameter.
var dataWorksColumnDataSourceTypes = []string{"maxcompute", "dlf", "hms", "holo", "mysql"}

func dataSourceAlicloudDataWorksColumns() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudDataWorksColumnsRead,
		Schema: map[string]*schema.Schema{
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"table_guid": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"data_source_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice(dataWorksColumnDataSourceTypes, false),
			},
			"cluster_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"database_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"table_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"page_size": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  PageSizeLarge,
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"columns": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"column_guid": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"column_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"comment": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"column_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"is_primary_key": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"region_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceAlicloudDataWorksColumnsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	action := "GetMetaTableColumn"
	query := make(map[string]interface{})
	if v, ok := d.GetOk("table_guid"); ok {
		query["TableGuid"] = v
	}
	if v, ok := d.GetOk("data_source_type"); ok {
		query["DataSourceType"] = v
	}
	if v, ok := d.GetOk("cluster_id"); ok {
		query["ClusterId"] = v
	}
	if v, ok := d.GetOk("database_name"); ok {
		query["DatabaseName"] = v
	}
	if v, ok := d.GetOk("table_name"); ok {
		query["TableName"] = v
	}

	pageSize := PageSizeLarge
	if v, ok := d.GetOk("page_size"); ok {
		pageSize = v.(int)
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = PageSizeLarge
	}
	query["PageSize"] = pageSize
	query["PageNum"] = 1

	idsMap := make(map[string]string)
	if v, ok := d.GetOk("ids"); ok {
		for _, vv := range v.([]interface{}) {
			if vv == nil {
				continue
			}
			idsMap[vv.(string)] = vv.(string)
		}
	}

	var objects []map[string]interface{}
	var response map[string]interface{}
	var err error
	for {
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			response, err = client.RpcGet("dataworks-public", "2020-05-18", action, query, nil)
			if err != nil {
				if NeedRetry(err) || IsExpectedErrors(err, []string{"Throttling.Api", "Throttling.Api.Risk", "Throttling.System", "Throttling.User", "InternalError.System", "InternalError.Meta.TenantTimeOut", "InternalError.Meta.Unknown"}) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})
		addDebug(action, response, query)
		if err != nil {
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_data_works_columns", action, AlibabaCloudSdkGoERROR)
		}
		resp, err := jsonpath.Get("$.Data.ColumnList", response)
		if err != nil {
			return WrapErrorf(err, FailedGetAttributeMsg, action, "$.Data.ColumnList", response)
		}
		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})
			if len(idsMap) > 0 {
				if _, ok := idsMap[fmt.Sprint(item["ColumnGuid"])]; !ok {
					continue
				}
			}
			objects = append(objects, item)
		}
		if len(result) < pageSize {
			break
		}
		query["PageNum"] = query["PageNum"].(int) + 1
	}

	regionId := client.RegionId
	ids := make([]string, 0)
	s := make([]map[string]interface{}, 0)
	for _, object := range objects {
		mapping := map[string]interface{}{
			"column_guid":    fmt.Sprint(object["ColumnGuid"]),
			"column_name":    fmt.Sprint(object["ColumnName"]),
			"comment":        fmt.Sprint(object["Comment"]),
			"column_type":    fmt.Sprint(object["ColumnType"]),
			"is_primary_key": object["IsPrimaryKey"],
			"region_id":      regionId,
		}
		ids = append(ids, fmt.Sprint(mapping["column_guid"]))
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}
	if err := d.Set("columns", s); err != nil {
		return WrapError(err)
	}
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}

	return nil
}
