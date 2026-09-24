package alicloud

import (
	"fmt"
	"regexp"
	"strconv"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func dataSourceAliCloudEsaAigwInstances() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAliCloudEsaAigwInstancesRead,
		Schema: map[string]*schema.Schema{
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"name_regex": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsValidRegExp,
			},
			"names": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"fuzzy_search_key": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"page_number": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"page_size": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"instances": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"aigw_instance_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"aigw_instance_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"auth_key": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"comment": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"create_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"enable_auth": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"record_count": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"update_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceAliCloudEsaAigwInstancesRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	action := "ListAIGWInstances"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})

	if v, ok := d.GetOk("fuzzy_search_key"); ok {
		query["FuzzySearchKey"] = v
	}
	pageNumber, pageSize := 1, 20
	if v, ok := d.GetOk("page_number"); ok && v.(int) > 0 {
		pageNumber = v.(int)
	}
	if v, ok := d.GetOk("page_size"); ok && v.(int) > 0 {
		pageSize = v.(int)
	}
	query["PageNumber"] = strconv.Itoa(pageNumber)
	query["PageSize"] = strconv.Itoa(pageSize)

	ids := make([]string, 0)
	names := make([]string, 0)
	instances := make([]map[string]interface{}, 0)

	var nameRegex *regexp.Regexp
	if v, ok := d.GetOk("name_regex"); ok {
		nameRegex, err = regexp.Compile(v.(string))
		if err != nil {
			return WrapError(err)
		}
	}

	idFilter := map[string]bool{}
	if v, ok := d.GetOk("ids"); ok {
		for _, id := range v.([]interface{}) {
			idFilter[fmt.Sprint(id)] = true
		}
	}

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = client.RpcGet("ESA", "2024-09-10", action, query, request)
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
		return WrapErrorf(err, DefaultErrorMsg, "data.alicloud_esa_aigw_instances", action, AlibabaCloudSdkGoERROR)
	}

	content, err := jsonpath.Get("$.data.content", response)
	if err != nil {
		return WrapErrorf(err, FailedGetAttributeMsg, "data.alicloud_esa_aigw_instances", "$.data.content", response)
	}
	contentMap, ok := content.(map[string]interface{})
	if !ok {
		return WrapError(fmt.Errorf("failed to parse data.content from ListAIGWInstances response"))
	}

	itemsRaw, ok := contentMap["Items"].([]interface{})
	if !ok {
		itemsRaw = make([]interface{}, 0)
	}

	for _, item := range itemsRaw {
		itemMap, ok := item.(map[string]interface{})
		if !ok {
			continue
		}
		instanceId := fmt.Sprint(itemMap["InstanceId"])
		name := fmt.Sprint(itemMap["Name"])

		if len(idFilter) > 0 && !idFilter[instanceId] {
			continue
		}
		if nameRegex != nil && !nameRegex.MatchString(name) {
			continue
		}

		ids = append(ids, instanceId)
		names = append(names, name)
		instances = append(instances, map[string]interface{}{
			"aigw_instance_id":   instanceId,
			"aigw_instance_name": name,
			"auth_key":           itemMap["AuthKey"],
			"comment":            itemMap["Comment"],
			"create_time":        itemMap["CreateTime"],
			"enable_auth":        itemMap["AuthEnabled"],
			"record_count":       itemMap["RecordCount"],
			"status":             itemMap["Status"],
			"update_time":        itemMap["UpdateTime"],
		})
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}
	if err := d.Set("names", names); err != nil {
		return WrapError(err)
	}
	if err := d.Set("instances", instances); err != nil {
		return WrapError(err)
	}

	return nil
}
