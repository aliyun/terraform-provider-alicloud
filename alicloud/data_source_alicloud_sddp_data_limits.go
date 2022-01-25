package alicloud

import (
	"fmt"
	"time"

	"github.com/PaesslerAG/jsonpath"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func dataSourceAlicloudSddpDataLimits() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudSddpDataLimitsRead,
		Schema: map[string]*schema.Schema{
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"parent_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"resource_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"MaxCompute", "OSS", "RDS"}, false),
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"limits": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"audit_status": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"check_status": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"data_limit_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"engine_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"local_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"log_store_day": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"parent_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"port": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"resource_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"user_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceAlicloudSddpDataLimitsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	action := "DescribeDataLimits"
	request := make(map[string]interface{})
	if v, ok := d.GetOk("resource_type"); ok {
		request["ResourceType"] = convertSddpDataLimitResourceTypeRequest(v)
	}
	if v, ok := d.GetOk("parent_id"); ok {
		request["ParentId"] = v
	}
	request["PageSize"] = PageSizeLarge
	request["CurrentPage"] = 1
	var objects []map[string]interface{}

	idsMap := make(map[string]string)
	if v, ok := d.GetOk("ids"); ok {
		for _, vv := range v.([]interface{}) {
			if vv == nil {
				continue
			}
			idsMap[vv.(string)] = vv.(string)
		}
	}
	var response map[string]interface{}
	conn, err := client.NewSddpClient()
	if err != nil {
		return WrapError(err)
	}
	for {
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-01-03"), StringPointer("AK"), nil, request, &runtime)
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
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_sddp_data_limits", action, AlibabaCloudSdkGoERROR)
		}
		resp, err := jsonpath.Get("$.Items", response)
		if err != nil {
			return WrapErrorf(err, FailedGetAttributeMsg, action, "$.Items", response)
		}
		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})
			if len(idsMap) > 0 {
				if _, ok := idsMap[fmt.Sprint(item["Id"])]; !ok {
					continue
				}
			}
			objects = append(objects, item)
		}
		if len(result) < PageSizeLarge {
			break
		}
		request["CurrentPage"] = request["CurrentPage"].(int) + 1
	}
	ids := make([]string, 0)
	s := make([]map[string]interface{}, 0)
	for _, object := range objects {
		mapping := map[string]interface{}{
			"audit_status":  formatInt(object["AuditStatus"]),
			"check_status":  formatInt(object["CheckStatus"]),
			"id":            fmt.Sprint(object["Id"]),
			"data_limit_id": fmt.Sprint(object["Id"]),
			"engine_type":   object["EngineType"],
			"local_name":    object["LocalName"],
			"log_store_day": formatInt(object["LogStoreDay"]),
			"parent_id":     object["ParentId"],
			"port":          formatInt(object["Port"]),
			"resource_type": fmt.Sprint(object["ResourceTypeCode"]),
			"user_name":     object["UserName"],
		}
		ids = append(ids, fmt.Sprint(mapping["id"]))
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}

	if err := d.Set("limits", s); err != nil {
		return WrapError(err)
	}
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}

	return nil
}
