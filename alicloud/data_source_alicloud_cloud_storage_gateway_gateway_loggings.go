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

func dataSourceAlicloudCloudStorageGatewayGatewayLoggings() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudCloudStorageGatewayGatewayLoggingsRead,
		Schema: map[string]*schema.Schema{
			"gateway_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"loggings": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"gateway_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"sls_logstore": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"sls_project": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceAlicloudCloudStorageGatewayGatewayLoggingsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	action := "DescribeGatewayLogging"
	request := map[string]interface{}{
		"GatewayId": d.Get("gateway_id"),
	}
	var response map[string]interface{}
	conn, err := client.NewHcsSgwClient()
	if err != nil {
		return WrapError(err)
	}
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2018-05-11"), StringPointer("AK"), nil, request, &runtime)
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
		return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_cloud_storage_gateway_gateway_loggings", action, AlibabaCloudSdkGoERROR)
	}
	if fmt.Sprint(response["Success"]) == "false" {
		return WrapError(fmt.Errorf("%s failed, response: %v", action, response))
	}
	resp, err := jsonpath.Get("$", response)
	if err != nil {
		return WrapErrorf(err, FailedGetAttributeMsg, action, "$", response)
	}
	result, _ := resp.(map[string]interface{})
	fmt.Println("result", result)
	s := make([]map[string]interface{}, 0)
	mapping := map[string]interface{}{
		"id":           fmt.Sprint(request["GatewayId"]),
		"gateway_id":   fmt.Sprint(request["GatewayId"]),
		"sls_logstore": result["SlsLogstore"],
		"sls_project":  result["SlsProject"],
		"status":       result["GatewayLoggingStatus"],
	}
	s = append(s, mapping)

	d.SetId(strconv.FormatInt(time.Now().Unix(), 16))
	if err := d.Set("loggings", s); err != nil {
		return WrapError(err)
	}
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}

	return nil
}
