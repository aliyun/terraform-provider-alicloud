package alicloud

import (
	"fmt"
	"log"
	"time"

	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceAlicloudVpcFlowLog() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudVpcFlowLogCreate,
		Read:   resourceAlicloudVpcFlowLogRead,
		Update: resourceAlicloudVpcFlowLogUpdate,
		Delete: resourceAlicloudVpcFlowLogDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},
		Schema: map[string]*schema.Schema{

			"flow_log_name": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"resource_type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"NetworkInterface", "VSwitch", "VPC"}, false),
			},

			"resource_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"traffic_type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"Allow", "All", "Drop"}, false),
			},

			"project_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"log_store_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"status": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringInSlice([]string{"Active", "Inactive"}, false),
			},
		},
	}
}

func resourceAlicloudVpcFlowLogCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	VpcFlowLogService := VpcService{client}

	var response map[string]interface{}
	action := "CreateFlowLog"
	request := map[string]interface{}{
		"RegionId": client.RegionId,
	}

	if v, ok := d.GetOk("flow_log_name"); ok {
		request["FlowLogName"] = v.(string)
	}
	if v, ok := d.GetOk("description"); ok {
		request["Description"] = v.(string)
	}

	request["ResourceType"] = d.Get("resource_type").(string)
	request["ResourceId"] = d.Get("resource_id").(string)
	request["TrafficType"] = d.Get("traffic_type").(string)
	request["ProjectName"] = d.Get("project_name").(string)
	request["LogStoreName"] = d.Get("log_store_name").(string)
	conn, err := meta.(*connectivity.AliyunClient).NewVpcClient()
	if err != nil {
		return WrapError(err)
	}
	// If the API supports
	runtime := util.RuntimeOptions{}
	wait := incrementalWait(3*time.Second, 1*time.Second)
	runtime.SetAutoretry(true)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2016-04-28"), StringPointer("AK"), nil, request, &runtime)
		if err != nil {
			if IsExpectedErrors(err, []string{"TaskConflict", "UnknownError", Throttling}) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(action, response, request)
		return nil
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_vpc_flow_log", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprintf("%v", response["FlowLogId"]))

	if err := VpcFlowLogService.WaitForVpcFlowLog(d.Id(), Active, 2*DefaultTimeout); err != nil {
		return WrapError(err)
	}

	return resourceAlicloudVpcFlowLogUpdate(d, meta)
}

func resourceAlicloudVpcFlowLogRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	VpcFlowLogService := VpcService{client}

	object, err := VpcFlowLogService.DescribeFlowLogs(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_vpc_flow_log VpcService.DescribeFlowLogs Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	d.Set("flow_log_name", object["FlowLogName"].(string))
	d.Set("description", object["Description"].(string))
	d.Set("resource_type", object["ResourceType"].(string))
	d.Set("resource_id", object["ResourceId"].(string))
	d.Set("traffic_type", object["TrafficType"].(string))
	d.Set("project_name", object["ProjectName"].(string))
	d.Set("log_store_name", object["LogStoreName"].(string))
	d.Set("status", object["Status"].(string))

	return nil
}

func resourceAlicloudVpcFlowLogUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	VpcFlowLogService := VpcService{client}
	var response map[string]interface{}
	d.Partial(true)

	request := map[string]interface{}{
		"RegionId": client.RegionId,
	}
	if !d.IsNewResource() && (d.HasChange("flow_log_name") || d.HasChange("description")) {
		request["FlowLogId"] = d.Id()
		request["FlowLogName"] = d.Get("flow_log_name").(string)
		request["Description"] = d.Get("description").(string)
		action := "ModifyFlowLogAttribute"
		conn, err := client.NewVpcClient()
		if err != nil {
			return WrapError(err)
		}
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2016-04-28"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			addDebug(action, response, request)
			return nil
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
		d.SetPartial("flow_log_name")
		d.SetPartial("description")
	}

	if d.HasChange("status") {
		var action string
		var status Status
		if v := d.Get("status").(string); v == "Active" {
			action = "ActiveFlowLog"
			status = Active
		} else {
			action = "DeactiveFlowLog"
			status = Inactive
		}

		var response map[string]interface{}
		request := map[string]interface{}{
			"RegionId": client.RegionId,
		}
		request["FlowLogId"] = d.Id()
		conn, err := meta.(*connectivity.AliyunClient).NewVpcClient()
		if err != nil {
			return WrapError(err)
		}
		// If the API supports
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		err = resource.Retry(3*time.Minute, func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2016-04-28"), StringPointer("AK"), nil, request, &runtime)
			if err != nil {
				if IsExpectedErrors(err, []string{"TaskConflict", "UnknownError", Throttling}) {
					time.Sleep(5 * time.Second)
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			addDebug(action, response, request)
			return nil
		})
		d.SetPartial("status")
		if err := VpcFlowLogService.WaitForVpcFlowLog(d.Id(), status, 2*DefaultTimeout); err != nil {
			return WrapError(err)
		}
	}
	d.Partial(false)
	return resourceAlicloudVpcFlowLogRead(d, meta)
}

func resourceAlicloudVpcFlowLogDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	VpcFlowLogService := VpcService{client}
	action := "DeleteFlowLog"
	var response map[string]interface{}
	conn, err := client.NewVpcClient()
	if err != nil {
		return WrapError(err)
	}
	request := map[string]interface{}{
		"RegionId": client.RegionId,
	}

	request["FlowLogId"] = d.Id()
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2016-04-28"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(action, response, request)
		return nil
	})
	if err != nil {
		if IsExpectedErrors(err, []string{"Instance.IsNotAvailable", "Instance.IsNotPostPay"}) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}
	return WrapError(VpcFlowLogService.WaitForVpcFlowLog(d.Id(), Deleted, DefaultTimeoutMedium))
}
