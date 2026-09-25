package alicloud

import (
	"fmt"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceAlicloudDataWorksTableLineage() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudDataWorksTableLineageRead,
		Schema: map[string]*schema.Schema{
			"direction": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: func(v interface{}, k string) (ws []string, es []error) {
					value := v.(string)
					if value != "up" && value != "down" {
						es = append(es, fmt.Errorf("%q must be one of [up, down], got %q", k, value))
					}
					return
				},
			},
			"table_guid": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"table_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"database_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"data_source_type": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: func(v interface{}, k string) (ws []string, es []error) {
					value := v.(string)
					if value != "odps" && value != "emr" {
						es = append(es, fmt.Errorf("%q must be one of [odps, emr], got %q", k, value))
					}
					return
				},
			},
			"cluster_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"page_size": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  100,
				ValidateFunc: func(v interface{}, k string) (ws []string, es []error) {
					value := v.(int)
					if value < 1 || value > 100 {
						es = append(es, fmt.Errorf("%q must be between 1 and 100, got %d", k, value))
					}
					return
				},
			},
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"lineages": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"table_guid": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"table_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"database_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"create_timestamp": {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceAlicloudDataWorksTableLineageRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	action := "GetMetaTableLineage"
	request := make(map[string]interface{})
	request["Direction"] = d.Get("direction").(string)
	if v, ok := d.GetOk("table_guid"); ok {
		request["TableGuid"] = v
	}
	if v, ok := d.GetOk("table_name"); ok {
		request["TableName"] = v
	}
	if v, ok := d.GetOk("database_name"); ok {
		request["DatabaseName"] = v
	}
	if v, ok := d.GetOk("data_source_type"); ok {
		request["DataSourceType"] = v
	}
	if v, ok := d.GetOk("cluster_id"); ok {
		request["ClusterId"] = v
	}
	pageSize := d.Get("page_size").(int)
	if pageSize < 1 {
		pageSize = 1
	}
	if pageSize > 100 {
		pageSize = 100
	}
	request["PageSize"] = pageSize

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
			response, err = client.RpcPost("dataworks-public", "2020-05-18", action, nil, request, true)
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})
		addDebug(action, response, request)
		if err != nil {
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_data_works_table_lineage", action, AlibabaCloudSdkGoERROR)
		}
		// jsonpath.Get requires a parsed JSON value; RpcPost returns response as
		// map[string]interface{} which is already the parsed JSON. Passing a raw
		// string would silently swallow the error and leave d.Set no-op.
		resp, perr := jsonpath.Get("$.Data.DataEntityList", response)
		if perr != nil {
			return WrapErrorf(perr, FailedGetAttributeMsg, action, "$.Data.DataEntityList", response)
		}
		result, _ := resp.([]interface{})
		for _, v := range result {
			item, ok := v.(map[string]interface{})
			if !ok {
				continue
			}
			tableGuid, _ := item["TableGuid"].(string)
			tableName, _ := item["TableName"].(string)
			databaseName, _ := item["DatabaseName"].(string)
			// CreateTimestamp may be a json.Number (RpcPost UseNumber decoding),
			// float64, or int64 depending on the SDK path; formatInt tolerates all.
			createTimestamp := formatInt(item["CreateTimestamp"])
			id := tableGuid
			if id == "" {
				// Fall back to a composite key when the API does not return a
				// TableGuid (e.g. EMR tables), so the element id remains unique
				// within the data source.
				id = fmt.Sprintf("%s:%s:%s", tableName, databaseName, d.Get("direction").(string))
			}
			if len(idsMap) > 0 {
				if _, ok := idsMap[id]; !ok {
					continue
				}
			}
			objects = append(objects, map[string]interface{}{
				"id":               id,
				"table_guid":       tableGuid,
				"table_name":       tableName,
				"database_name":    databaseName,
				"create_timestamp": createTimestamp,
			})
		}
		// GetMetaTableLineage paginates with HasNext + NextPrimaryKey (cursor
		// pagination), not PageNumber. Stop when HasNext is false/missing or
		// NextPrimaryKey is empty, or when the current page returned no items.
		hasNext := false
		if v, perr := jsonpath.Get("$.Data.HasNext", response); perr == nil {
			if hb, ok := v.(bool); ok {
				hasNext = hb
			}
		}
		nextPrimaryKey := ""
		if v, perr := jsonpath.Get("$.Data.NextPrimaryKey", response); perr == nil {
			if s, ok := v.(string); ok {
				nextPrimaryKey = s
			}
		}
		if !hasNext || nextPrimaryKey == "" || len(result) == 0 {
			break
		}
		request["NextPrimaryKey"] = nextPrimaryKey
	}

	ids := make([]string, 0)
	s := make([]map[string]interface{}, 0)
	for _, object := range objects {
		ids = append(ids, object["id"].(string))
		s = append(s, object)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}
	if err := d.Set("lineages", s); err != nil {
		return WrapError(err)
	}
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}
	return nil
}
