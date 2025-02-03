package alicloud

import (
	"fmt"
	"regexp"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceAlicloudDmsEnterpriseDatabases() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudDmsEnterpriseDatabasesRead,
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Required: true,
				ForceNew: true,
				Type:     schema.TypeString,
			},
			"name_regex": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.ValidateRegexp,
				ForceNew:     true,
			},
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
						"catalog_name": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"database_id": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"db_type": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"dba_id": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"dba_name": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"encoding": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"env_type": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"host": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"instance_id": {
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
						"port": {
							Computed: true,
							Type:     schema.TypeInt,
						},
						"schema_name": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"search_name": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"sid": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"state": {
							Computed: true,
							Type:     schema.TypeString,
						},
					},
				},
			},
		},
	}
}

func dataSourceAlicloudDmsEnterpriseDatabasesRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	request := map[string]interface{}{
		"RegionId": client.RegionId,
	}

	if v, ok := d.GetOk("instance_id"); ok {
		request["InstanceId"] = v
	}
	request["PageNumber"] = 1
	request["PageSize"] = PageSizeSmall

	var instanceNameRegex *regexp.Regexp
	if v, ok := d.GetOk("name_regex"); ok {
		r, err := regexp.Compile(v.(string))
		if err != nil {
			return WrapError(err)
		}
		instanceNameRegex = r
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
		action := "ListDatabases"
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
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_dms_enterprise_databases", action, AlibabaCloudSdkGoERROR)
		}
		resp, err := jsonpath.Get("$.DatabaseList.Database", response)
		if err != nil {
			return WrapErrorf(err, FailedGetAttributeMsg, action, "$.DatabaseList.Database", response)
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
		if instanceNameRegex != nil {
			if !instanceNameRegex.MatchString(fmt.Sprint(object["SchemaName"])) {
				continue
			}
		}
		mapping := map[string]interface{}{
			"id":           fmt.Sprint(object["DatabaseId"]),
			"catalog_name": object["CatalogName"],
			"database_id":  object["DatabaseId"],
			"db_type":      object["DbType"],
			"dba_id":       object["DbaId"],
			"dba_name":     object["DbaName"],
			"encoding":     object["Encoding"],
			"env_type":     object["EnvType"],
			"host":         object["Host"],
			"instance_id":  object["InstanceId"],
			"port":         object["Port"],
			"schema_name":  object["SchemaName"],
			"search_name":  object["SearchName"],
			"sid":          object["Sid"],
			"state":        object["State"],
		}

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
