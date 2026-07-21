package alicloud

import (
	"fmt"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceAlicloudDmsEnterpriseLogicDatabases() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudDmsEnterpriseLogicDatabasesRead,
		Schema: map[string]*schema.Schema{
			"ids": {
				Optional: true,
				Computed: true,
				Type:     schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"output_file": {
				Optional: true,
				Type:     schema.TypeString,
			},
			"databases": {
				Computed: true,
				Type:     schema.TypeList,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"alias": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"database_ids": {
							Computed: true,
							Type:     schema.TypeList,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"db_type": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"env_type": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"logic": {
							Computed: true,
							Type:     schema.TypeBool,
						},
						"logic_database_id": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"owner_id_list": {
							Computed: true,
							Type:     schema.TypeList,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"owner_name_list": {
							Computed: true,
							Type:     schema.TypeList,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"schema_name": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"search_name": {
							Computed: true,
							Type:     schema.TypeString,
						},
					},
				},
			},
		},
	}
}

func dataSourceAlicloudDmsEnterpriseLogicDatabasesRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	request := map[string]interface{}{
		"RegionId": client.RegionId,
	}

	if v, ok := d.GetOk("page_number"); ok && v.(int) > 0 {
		request["PageNumber"] = v.(int)
	} else {
		request["PageNumber"] = 1
	}
	if v, ok := d.GetOk("page_size"); ok && v.(int) > 0 {
		request["PageSize"] = v.(int)
	} else {
		request["PageSize"] = 10
	}

	idsMap := make(map[string]string)
	if v, ok := d.GetOk("ids"); ok {
		for _, vv := range v.([]interface{}) {
			if vv == nil {
				continue
			}
			idsMap[vv.(string)] = vv.(string)
		}
	}

	var err error
	var objects []interface{}
	var response map[string]interface{}

	for {
		action := "ListLogicDatabases"
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			resp, err := client.RpcPost("dms-enterprise", "2018-11-01", action, nil, request, true)
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			response = resp
			addDebug(action, response, request)
			return nil
		})
		if err != nil {
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_dms_enterprise_logic_databases", action, AlibabaCloudSdkGoERROR)
		}
		resp, err := jsonpath.Get("$.LogicDatabaseList.LogicDatabase", response)
		if err != nil {
			return WrapErrorf(err, FailedGetAttributeMsg, action, "$.LogicDatabaseList.LogicDatabase", response)
		}
		result, _ := resp.([]interface{})
		if isPagingRequest(d) {
			objects = result
			break
		}
		for _, v := range result {
			item := v.(map[string]interface{})
			if len(idsMap) > 0 {
				if _, ok := idsMap[fmt.Sprint(item["DatabaseId"])]; !ok {
					continue
				}
			}
			objects = append(objects, item)
		}
		if len(result) < request["PageSize"].(int) {
			break
		}
		request["PageNumber"] = request["PageNumber"].(int) + 1
	}

	ids := make([]string, 0)
	s := make([]map[string]interface{}, 0)
	for _, v := range objects {
		object := v.(map[string]interface{})
		mapping := map[string]interface{}{
			"id":                fmt.Sprint(object["DatabaseId"]),
			"alias":             object["Alias"],
			"db_type":           object["DbType"],
			"env_type":          object["EnvType"],
			"logic":             object["Logic"],
			"logic_database_id": object["DatabaseId"],
			"schema_name":       object["SchemaName"],
			"search_name":       object["SearchName"],
		}

		databaseIds, _ := jsonpath.Get("$.DatabaseIds.DatabaseIds", object)
		mapping["database_ids"] = databaseIds
		ownerIdList, _ := jsonpath.Get("$.OwnerIdList.OwnerIds", object)
		mapping["owner_id_list"] = ownerIdList
		ownerNameList, _ := jsonpath.Get("$.OwnerNameList.OwnerNames", object)
		mapping["owner_name_list"] = ownerNameList

		ids = append(ids, fmt.Sprint(object["DatabaseId"]))
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}

	if err := d.Set("databases", s); err != nil {
		return WrapError(err)
	}
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}
	return nil
}
