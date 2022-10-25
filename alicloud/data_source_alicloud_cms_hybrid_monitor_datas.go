package alicloud

import (
	"fmt"
	"strconv"
	"time"

	"github.com/PaesslerAG/jsonpath"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceAlicloudCmsHybridMonitorDatas() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudCmsHybridMonitorDatasRead,
		Schema: map[string]*schema.Schema{
			"namespace": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"prom_sql": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"start": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"end": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"period": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},

			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"datas": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"labels": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"key": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"value": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"metric_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"values": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"ts": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"value": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func dataSourceAlicloudCmsHybridMonitorDatasRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	action := "DescribeHybridMonitorDataList"
	request := make(map[string]interface{})
	request["End"] = d.Get("end")
	request["Namespace"] = d.Get("namespace")
	if v, ok := d.GetOk("period"); ok {
		request["Period"] = v
	}
	request["PromSQL"] = d.Get("prom_sql")
	request["Start"] = d.Get("start")
	var response map[string]interface{}
	conn, err := client.NewCmsClient()
	if err != nil {
		return WrapError(err)
	}
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(2*time.Minute, func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-01-01"), StringPointer("AK"), nil, request, &runtime)
		if err != nil {
			if IsExpectedErrors(err, []string{"InternalError", "BadRequest"}) || NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)
	if err != nil {
		return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_cms_hybrid_monitor_datas", action, AlibabaCloudSdkGoERROR)
	}
	if fmt.Sprint(response["Success"]) == "false" {
		return WrapError(fmt.Errorf("%s failed, response: %v", action, response))
	}
	resp, err := jsonpath.Get("$.TimeSeries", response)
	if err != nil {
		return WrapErrorf(err, FailedGetAttributeMsg, action, "$.TimeSeries", response)
	}
	s := make([]map[string]interface{}, 0)
	if timeSeriesList, ok := resp.([]interface{}); ok {
		for _, v := range timeSeriesList {
			if m1, ok := v.(map[string]interface{}); ok {
				temp1 := map[string]interface{}{
					"metric_name": m1["MetricName"],
				}
				if m1["Labels"] != nil {
					labelsMaps := make([]map[string]interface{}, 0)
					for _, labelsValue := range m1["Labels"].([]interface{}) {
						labels := labelsValue.(map[string]interface{})
						labelsMap := map[string]interface{}{
							"key":   labels["K"],
							"value": labels["V"],
						}
						labelsMaps = append(labelsMaps, labelsMap)
					}
					temp1["labels"] = labelsMaps
				}
				if m1["Values"] != nil {
					valuesMaps := make([]map[string]interface{}, 0)
					for _, valuesValue := range m1["Values"].([]interface{}) {
						values := valuesValue.(map[string]interface{})
						valuesMap := map[string]interface{}{
							"ts":    values["Ts"],
							"value": values["V"],
						}
						valuesMaps = append(valuesMaps, valuesMap)
					}
					temp1["values"] = valuesMaps
				}
				s = append(s, temp1)
			}
		}
	}
	d.SetId(strconv.FormatInt(time.Now().Unix(), 16))

	if err := d.Set("datas", s); err != nil {
		return WrapError(err)
	}
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}

	return nil
}
