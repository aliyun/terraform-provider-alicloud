package alicloud

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAlicloudEcsInvocation() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudEcsInvocationCreate,
		Read:   resourceAlicloudEcsInvocationRead,
		Delete: resourceAlicloudEcsInvocationDelete,
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(1 * time.Minute),
			Delete: schema.DefaultTimeout(1 * time.Minute),
		},
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"command_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"instance_id": {
				Type:     schema.TypeList,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				ForceNew: true,
			},
			"repeat_mode": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Computed:     true,
				ValidateFunc: validation.StringInSlice([]string{"Once", "Period", "NextRebootOnly", "EveryReboot"}, false),
			},
			"timed": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"frequency": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return d.Get("repeat_mode").(string) != "Period"
				},
			},
			"parameters": {
				Type:     schema.TypeMap,
				Optional: true,
				ForceNew: true,
			},
			"username": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"windows_password_name": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceAlicloudEcsInvocationCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "InvokeCommand"
	request := make(map[string]interface{})
	conn, err := client.NewEcsClient()
	if err != nil {
		return WrapError(err)
	}
	request["CommandId"] = d.Get("command_id")
	request["InstanceId"] = d.Get("instance_id")
	request["RegionId"] = client.RegionId
	if v, ok := d.GetOk("repeat_mode"); ok {
		request["RepeatMode"] = v
	}
	if v, ok := d.GetOkExists("timed"); ok {
		request["Timed"] = v
	}
	if v, ok := d.GetOk("frequency"); ok {
		request["Frequency"] = v
	}
	if v, ok := d.GetOk("parameters"); ok {
		parameters, err := convertMaptoJsonString(v.(map[string]interface{}))
		if err != nil {
			return WrapError(err)
		}
		request["Parameters"] = parameters
	}
	if v, ok := d.GetOk("username"); ok {
		request["Username"] = v
	}
	if v, ok := d.GetOk("windows_password_name"); ok {
		request["WindowsPasswordName"] = v
	}
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2014-05-26"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_ecs_invocation", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(response["InvokeId"]))

	ecsService := EcsService{client}
	stateConf := BuildStateConf([]string{}, []string{"Scheduled", "Success"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, ecsService.EcsInvocationStateRefreshFunc(d.Id(), []string{"Failed", "PartialFailed", "Stopped"}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return resourceAlicloudEcsInvocationRead(d, meta)
}
func resourceAlicloudEcsInvocationRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	ecsService := EcsService{client}
	object, err := ecsService.DescribeEcsInvocation(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_ecs_invocation ecsService.DescribeEcsInvocation Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	d.Set("command_id", object["CommandId"])
	d.Set("frequency", object["Frequency"])
	d.Set("repeat_mode", object["RepeatMode"])

	parametersMap := make(map[string]string)
	err = json.Unmarshal([]byte(object["Parameters"].(string)), &parametersMap)
	if err != nil {
		return WrapError(err)
	}

	d.Set("parameters", parametersMap)
	d.Set("timed", object["Timed"])
	d.Set("username", object["Username"])
	d.Set("status", object["InvocationStatus"])
	instanceIdItems := make([]string, 0)
	if invokeInstances, ok := object["InvokeInstances"]; ok && invokeInstances != nil {
		if invokeInstance, ok := invokeInstances.(map[string]interface{})["InvokeInstance"]; ok && invokeInstance != nil {
			for _, invokeInstanceItem := range invokeInstance.([]interface{}) {
				if instanceId, ok := invokeInstanceItem.(map[string]interface{})["InstanceId"]; ok && instanceId != nil {
					instanceIdItems = append(instanceIdItems, fmt.Sprint(instanceId))
				}
			}
		}

	}
	d.Set("instance_id", instanceIdItems)
	return nil
}
func resourceAlicloudEcsInvocationDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "StopInvocation"
	request := make(map[string]interface{})
	conn, err := client.NewEcsClient()
	if err != nil {
		return WrapError(err)
	}
	request["InvokeId"] = d.Id()
	request["RegionId"] = client.RegionId
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2014-05-26"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_ecs_invocation", action, AlibabaCloudSdkGoERROR)
	}

	log.Printf("[WARN] Cannot destroy resourceAlicloudEcsInvocation. Terraform will remove this resource from the state file, however resources may remain.")
	return nil
}
