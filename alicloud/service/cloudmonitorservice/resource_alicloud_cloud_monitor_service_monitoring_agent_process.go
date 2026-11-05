package cloudmonitorservice

import (
	"fmt"
	"log"
	"time"

	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAliCloudCloudMonitorServiceMonitoringAgentProcess() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudCloudMonitorServiceMonitoringAgentProcessCreate,
		Read:   resourceAliCloudCloudMonitorServiceMonitoringAgentProcessRead,
		Delete: resourceAliCloudCloudMonitorServiceMonitoringAgentProcessDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"process_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"process_user": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"process_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceAliCloudCloudMonitorServiceMonitoringAgentProcessCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "CreateMonitorAgentProcess"
	request := make(map[string]interface{})
	var err error

	request["InstanceId"] = d.Get("instance_id")
	request["ProcessName"] = d.Get("process_name")

	if v, ok := d.GetOk("process_user"); ok {
		request["ProcessUser"] = v
	}

	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutCreate)), func() *retry.RetryError {
		response, err = client.RpcPost("Cms", "2019-01-01", action, nil, request, false)
		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	logs.AddDebug(action, response, request)

	if err != nil {
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, "alicloud_cloud_monitor_service_monitoring_agent_process", action, errs.AlibabaCloudSdkGoERROR)
	}

	if fmt.Sprint(response["Success"]) == "false" {
		return errs.WrapError(fmt.Errorf("%s failed, response: %v", action, response))
	}

	d.SetId(fmt.Sprintf("%v:%v", request["InstanceId"], response["Id"]))

	return resourceAliCloudCloudMonitorServiceMonitoringAgentProcessRead(d, meta)
}

func resourceAliCloudCloudMonitorServiceMonitoringAgentProcessRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	cloudMonitorServiceServiceV2 := CloudMonitorServiceServiceV2{client}

	object, err := cloudMonitorServiceServiceV2.DescribeCloudMonitorServiceMonitoringAgentProcess(d.Id())
	if err != nil {
		if !d.IsNewResource() && errs.NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_cloud_monitor_service_monitoring_agent_process DescribeCloudMonitorServiceMonitoringAgentProcess Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return errs.WrapError(err)
	}

	d.Set("instance_id", object["InstanceId"])
	d.Set("process_id", object["ProcessId"])
	d.Set("process_name", object["ProcessName"])
	d.Set("process_user", object["ProcessUser"])

	return nil
}

func resourceAliCloudCloudMonitorServiceMonitoringAgentProcessDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	action := "DeleteMonitoringAgentProcess"
	var response map[string]interface{}

	var err error

	parts, err := helper.ParseResourceId(d.Id(), 2)
	if err != nil {
		return errs.WrapError(err)
	}

	request := map[string]interface{}{
		"InstanceId": parts[0],
		"ProcessId":  parts[1],
	}

	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutDelete)), func() *retry.RetryError {
		response, err = client.RpcPost("Cms", "2019-01-01", action, nil, request, false)
		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	logs.AddDebug(action, response, request)

	if err != nil {
		if errs.NotFoundError(err) {
			return nil
		}
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
	}

	if fmt.Sprint(response["Success"]) == "false" {
		return errs.WrapError(fmt.Errorf("%s failed, response: %v", action, response))
	}

	return nil
}
