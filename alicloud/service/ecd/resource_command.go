package ecd

import (
	"encoding/base64"
	"fmt"
	"log"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func ResourceCommand() *schema.Resource {
	return &schema.Resource{
		Create: resourceCommandCreate,
		Read:   resourceCommandRead,
		Update: resourceCommandUpdate,
		Delete: resourceCommandDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"command_content": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"command_type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"RunBatScript", "RunPowerShellScript"}, false),
			},
			"content_encoding": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"Base64", "PlainText"}, false),
			},
			"desktop_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"timeout": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceCommandCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	ecdService := NewService(client)
	var response map[string]interface{}
	action := "RunCommand"
	request := make(map[string]interface{})
	var err error
	request["CommandContent"] = d.Get("command_content")
	request["Type"] = d.Get("command_type")
	if v, ok := d.GetOk("content_encoding"); ok {
		request["ContentEncoding"] = v
	}
	request["DesktopId.1"] = d.Get("desktop_id")
	request["RegionId"] = client.RegionId
	if v, ok := d.GetOk("timeout"); ok {
		request["Timeout"] = v
	}
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
		response, err = client.RpcPost("ecd", "2020-09-30", action, nil, request, false)
		if err != nil {
			if helper.NeedRetry(err) || errs.IsExpectedErrors(err, []string{"CloudAssistant.NotReady"}) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	logs.AddDebug(action, response, request)
	if err != nil {
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, "alicloud_ecd_command", action, errs.AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(response["InvokeId"]))
	stateConf := helper.BuildStateConf([]string{}, []string{"Success"}, d.Timeout(schema.TimeoutCreate), 10*time.Second, ecdService.EcdCommandStateRefreshFunc(d.Id(), []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return errs.WrapErrorf(err, errs.IdMsg, d.Id())
	}
	return resourceCommandRead(d, meta)
}
func resourceCommandRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	ecdService := NewService(client)
	object, err := ecdService.DescribeEcdCommand(d.Id())
	if err != nil {
		if errs.NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_ecd_command ecdService.DescribeEcdCommand Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return errs.WrapError(err)
	}
	commandContent, _ := base64.StdEncoding.DecodeString(object["CommandContent"].(string))
	d.Set("command_content", string(commandContent))
	d.Set("command_type", object["CommandType"])
	d.Set("status", object["InvocationStatus"])
	temp1 := map[string]interface{}{}
	if invokeDesktopsList, ok := object["InvokeDesktops"].([]interface{}); ok {
		for _, v := range invokeDesktopsList {
			if m1, ok := v.(map[string]interface{}); ok {
				temp1 = map[string]interface{}{
					"desktop_id": m1["DesktopId"],
				}
			}
		}
	}
	d.Set("desktop_id", temp1["desktop_id"])
	return nil
}
func resourceCommandUpdate(d *schema.ResourceData, meta interface{}) error {
	log.Println(fmt.Sprintf("[WARNING] The resouce has not update operation."))
	return nil
}
func resourceCommandDelete(d *schema.ResourceData, meta interface{}) error {
	log.Println(fmt.Sprintf("[WARNING] The resouce has not delete operation."))
	return nil

}
