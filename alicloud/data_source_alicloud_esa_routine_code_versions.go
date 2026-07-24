package alicloud

import (
	"fmt"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceAliCloudEsaRoutineCodeVersions() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAliCloudEsaRoutineCodeVersionsRead,
		Schema: map[string]*schema.Schema{
			"routine_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"search_key_word": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"versions": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"code_version": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"code_description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"create_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"deploy_env": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"build_id": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"has_env_vars": {
							Type:     schema.TypeBool,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceAliCloudEsaRoutineCodeVersionsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	action := "ListRoutineCodeVersions"
	request := make(map[string]interface{})
	request["Name"] = d.Get("routine_name")
	request["PageSize"] = 20
	request["PageNumber"] = 1
	if v, ok := d.GetOk("search_key_word"); ok {
		request["SearchKeyWord"] = v
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

	var objects []map[string]interface{}
	var response map[string]interface{}
	var err error

	// ListRoutineCodeVersions caps at PageNumber 2 and PageSize 20 (up to 40
	// versions). The loop respects that server-side ceiling.
	for {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			response, err = client.RpcPost("ESA", "2024-09-10", action, nil, request, true)
			if err != nil {
				if IsExpectedErrors(err, []string{"Site.ServiceBusy", "TooManyRequests"}) || NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})
		addDebug(action, response, request)
		if err != nil {
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_esa_routine_code_versions", action, AlibabaCloudSdkGoERROR)
		}

		resp, err := jsonpath.Get("$.CodeVersions", response)
		if err != nil {
			return WrapErrorf(err, FailedGetAttributeMsg, action, "$.CodeVersions", response)
		}
		result, _ := resp.([]interface{})
		for _, v := range result {
			item, ok := v.(map[string]interface{})
			if !ok {
				continue
			}
			if len(idsMap) > 0 {
				if _, ok := idsMap[fmt.Sprint(item["CodeVersion"])]; !ok {
					continue
				}
			}
			objects = append(objects, item)
		}

		pageNumber := request["PageNumber"].(int)
		if len(result) < 20 || pageNumber >= 2 {
			break
		}
		request["PageNumber"] = pageNumber + 1
	}

	ids := make([]string, 0)
	s := make([]map[string]interface{}, 0)
	for _, object := range objects {
		codeVersion := fmt.Sprint(object["CodeVersion"])
		mapping := map[string]interface{}{
			"id":               codeVersion,
			"code_version":     codeVersion,
			"code_description": object["CodeDescription"],
			"create_time":      object["CreateTime"],
			"status":           object["Status"],
			"deploy_env":       object["DeployEnv"],
		}
		if v, ok := object["BuildId"]; ok && v != nil {
			mapping["build_id"] = formatInt(v)
		}
		if v, ok := object["HasEnvVars"]; ok && v != nil {
			mapping["has_env_vars"] = v
		}

		ids = append(ids, codeVersion)
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))

	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}
	if err := d.Set("versions", s); err != nil {
		return WrapError(err)
	}

	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}

	return nil
}
