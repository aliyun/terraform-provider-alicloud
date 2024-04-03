package alicloud

import (
	"github.com/PaesslerAG/jsonpath"
	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"time"

	util "github.com/alibabacloud-go/tea-utils/v2/service"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"

	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

const maxWaitTime = 60

func dataSourceAlicloudLogService() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudLogServiceRead,

		Schema: map[string]*schema.Schema{
			"enable": {
				Type:         schema.TypeString,
				ValidateFunc: validation.StringInSlice([]string{"On", "Off"}, false),
				Optional:     true,
				Default:      "Off",
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}
func dataSourceAlicloudLogServiceRead(d *schema.ResourceData, meta interface{}) error {
	if v, ok := d.GetOk("enable"); !ok || v.(string) != "On" {
		d.SetId("LogServiceHasNotBeenOpened")
		d.Set("status", "")
		return nil
	}
	client := meta.(*connectivity.AliyunClient)
	conn, err := client.NewSlsClient()

	if err != nil {
		return WrapError(err)
	}
	isNotOpened, err := waitServiceReady(conn, false)
	if err == nil {
		d.SetId("LogServiceHasBeenOpened")
		d.Set("status", "Opened")
		return nil
	}
	if isNotOpened {
		action := fmt.Sprintf("/slsservice")
		response, err := conn.Execute(genRoaParam("OpenSlsService", "POST", "2020-12-30", action), &openapi.OpenApiRequest{Query: nil, Body: nil, HostMap: nil}, &util.RuntimeOptions{})
		addDebug("OpenSlsService", response, nil)
		if err != nil {
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_log_service", "OpenLogService", AlibabaCloudSdkGoERROR)
		}
		code, _ := jsonpath.Get("$.statusCode", response)
		if code == nil || fmt.Sprint(code) != "200" {
			return WrapErrorf(fmt.Errorf("%s", response), DataDefaultErrorMsg, "alicloud_log_service", "OpenLogService", AlibabaCloudSdkGoERROR)
		}
		_, err = waitServiceReady(conn, true)
		if err == nil {
			d.SetId("LogServiceHasBeenOpened")
			d.Set("status", "Opened")
			return nil
		}
		return WrapError(err)
	}
	return WrapError(err)
}

func waitServiceReady(conn *openapi.Client, hasOpened bool) (bool, error) {
	beginTime := time.Now().Unix()
	for {
		action := fmt.Sprintf("/slsservice")
		response, err := conn.Execute(genRoaParam("GetSlsService", "GET", "2020-12-30", action), &openapi.OpenApiRequest{Query: nil, Body: nil, HostMap: nil}, &util.RuntimeOptions{})
		addDebug("GetSlsService", response, nil)
		if err != nil {
			return false, WrapErrorf(err, DataDefaultErrorMsg, "alicloud_log_service", "GetLogService", AlibabaCloudSdkGoERROR)
		}
		if response["success"] != nil && !response["success"].(bool) {
			return false, WrapErrorf(fmt.Errorf("%s", response), DataDefaultErrorMsg, "alicloud_log_service", "GetLogService", AlibabaCloudSdkGoERROR)
		}
		res := response["body"].(map[string]interface{})
		status := res["status"].(string)
		if "Opened" == status {
			return false, nil
		}
		if hasOpened || "Opening" == status {
			if time.Now().Unix()-beginTime >= maxWaitTime {
				return false, fmt.Errorf("wait until the maxWaitTime(60s) is still in the %s state", status)
			}
			time.Sleep(time.Second)
			continue
		}
		return true, fmt.Errorf("incorrect status: %s", status)
	}

}
