package alicloud

import (
	"fmt"
	"regexp"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func dataSourceAlicloudDataworksDataServiceAuthorizedApis() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudDataworksDataServiceAuthorizedApisRead,
		Schema: map[string]*schema.Schema{
			"api_name_keyword": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"project_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"tenant_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"page_number": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  1,
			},
			"page_size": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      PageSizeLarge,
				ValidateFunc: validation.IntAtLeast(1),
			},
			"name_regex": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsValidRegExp,
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
			"total_count": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"names": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"apis": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"api_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"api_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"api_path": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"project_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"tenant_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"group_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"region_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"grant_end_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"grant_created_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"grant_operator_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"create_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"modified_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceAlicloudDataworksDataServiceAuthorizedApisRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	action := "ListDataServiceAuthorizedApis"
	request := map[string]interface{}{}
	if v, ok := d.GetOk("api_name_keyword"); ok {
		request["ApiNameKeyword"] = v.(string)
	}
	request["ProjectId"] = d.Get("project_id").(string)
	if v, ok := d.GetOk("tenant_id"); ok {
		request["TenantId"] = v.(string)
	}
	if v, ok := d.GetOk("page_number"); ok && v.(int) > 0 {
		request["PageNumber"] = v.(int)
	} else {
		request["PageNumber"] = 1
	}
	if v, ok := d.GetOk("page_size"); ok && v.(int) > 0 {
		request["PageSize"] = v.(int)
	} else {
		request["PageSize"] = PageSizeLarge
	}

	var apiNameRegex *regexp.Regexp
	if v, ok := d.GetOk("name_regex"); ok {
		r, err := regexp.Compile(v.(string))
		if err != nil {
			return WrapError(err)
		}
		apiNameRegex = r
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
	pageSize := request["PageSize"].(int)
	for {
		response, err = client.RpcPost("dataworks-public", "2020-05-18", action, nil, request, true)
		if err != nil {
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_dataworks_data_service_authorized_apis", action, AlibabaCloudSdkGoERROR)
		}
		addDebug(action, response, request)

		resp, err := jsonpath.Get("$.Data.ApiAuthorizedList", response)
		if err != nil {
			return WrapErrorf(err, FailedGetAttributeMsg, action, "$.Data.ApiAuthorizedList", response)
		}
		result, _ := resp.([]interface{})
		if isPagingRequest(d) {
			for _, v := range result {
				if item, ok := v.(map[string]interface{}); ok {
					objects = append(objects, item)
				}
			}
			break
		}
		for _, v := range result {
			item, ok := v.(map[string]interface{})
			if !ok {
				continue
			}
			if apiNameRegex != nil {
				if !apiNameRegex.MatchString(fmt.Sprint(item["ApiName"])) {
					continue
				}
			}
			if len(idsMap) > 0 {
				if _, ok := idsMap[fmt.Sprint(item["ApiId"])]; !ok {
					continue
				}
			}
			objects = append(objects, item)
		}
		if len(result) < pageSize {
			break
		}
		request["PageNumber"] = request["PageNumber"].(int) + 1
	}

	totalCount := 0
	if tc, err := jsonpath.Get("$.Data.TotalCount", response); err == nil {
		totalCount = formatInt(tc)
	}

	ids := make([]string, 0)
	names := make([]string, 0)
	s := make([]map[string]interface{}, 0)
	for _, object := range objects {
		apiId := fmt.Sprint(object["ApiId"])
		mapping := map[string]interface{}{
			"api_id":             apiId,
			"api_name":           object["ApiName"],
			"api_path":           object["ApiPath"],
			"project_id":         object["ProjectId"],
			"tenant_id":          object["TenantId"],
			"group_id":           object["GroupId"],
			"region_id":          object["RegionId"],
			"status":             object["Status"],
			"grant_end_time":     object["GrantEndTime"],
			"grant_created_time": object["GrantCreatedTime"],
			"grant_operator_id":  object["GrantOperatorId"],
			"create_time":        object["CreateTime"],
			"modified_time":      object["ModifiedTime"],
			"id":                 apiId,
		}
		ids = append(ids, apiId)
		names = append(names, fmt.Sprint(object["ApiName"]))
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}
	if err := d.Set("names", names); err != nil {
		return WrapError(err)
	}
	if err := d.Set("apis", s); err != nil {
		return WrapError(err)
	}
	if err := d.Set("total_count", totalCount); err != nil {
		return WrapError(err)
	}
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}

	return nil
}
