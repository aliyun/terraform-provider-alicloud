package alicloud

import (
	"fmt"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceAlicloudDataWorksDataServiceApplication() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudDataWorksDataServiceApplicationRead,
		Schema: map[string]*schema.Schema{
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"project_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"project_id_list": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"tenant_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"applications": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"application_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"application_key": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"application_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"application_secret": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"project_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"region_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"application_code": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceAlicloudDataWorksDataServiceApplicationRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	action := "ListDataServiceApplications"
	request := make(map[string]interface{})

	var projectIds []interface{}
	if v, ok := d.GetOk("project_id_list"); ok && len(v.([]interface{})) > 0 {
		projectIds = v.([]interface{})
	} else if v, ok := d.GetOk("project_id"); ok && v.(string) != "" {
		projectIds = []interface{}{v.(string)}
	} else {
		return WrapError(fmt.Errorf("One of `project_id` or `project_id_list` must be set to query Data Works Data Service Applications."))
	}
	if len(projectIds) == 0 {
		return WrapError(fmt.Errorf("At least one project id must be set to query Data Works Data Service Applications."))
	}
	request["ProjectIdList"] = convertListToCommaSeparate(projectIds)

	if v, ok := d.GetOk("tenant_id"); ok {
		request["TenantId"] = v
	}
	request["PageSize"] = PageSizeLarge
	request["PageNumber"] = 1

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
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_data_works_data_service_application", action, AlibabaCloudSdkGoERROR)
		}
		resp, err := jsonpath.Get("$.Data.Applications", response)
		if err != nil {
			return WrapErrorf(err, FailedGetAttributeMsg, action, "$.Data.Applications", response)
		}
		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})
			if len(idsMap) > 0 {
				if _, ok := idsMap[fmt.Sprint(item["ApplicationId"])]; !ok {
					continue
				}
			}
			objects = append(objects, item)
		}
		if len(result) < PageSizeLarge {
			break
		}
		request["PageNumber"] = request["PageNumber"].(int) + 1
	}

	ids := make([]string, 0)
	s := make([]map[string]interface{}, 0)
	for _, object := range objects {
		mapping := map[string]interface{}{
			"id":                 fmt.Sprint(object["ApplicationId"]),
			"application_id":     fmt.Sprint(object["ApplicationId"]),
			"application_key":    fmt.Sprint(object["ApplicationKey"]),
			"application_name":   fmt.Sprint(object["ApplicationName"]),
			"application_secret": fmt.Sprint(object["ApplicationSecret"]),
			"project_id":         fmt.Sprint(object["ProjectId"]),
			"region_id":          fmt.Sprint(object["RegionId"]),
			"application_code":   fmt.Sprint(object["ApplicationCode"]),
		}
		ids = append(ids, fmt.Sprint(mapping["id"]))
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}
	if err := d.Set("applications", s); err != nil {
		return WrapError(err)
	}
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}

	return nil
}
