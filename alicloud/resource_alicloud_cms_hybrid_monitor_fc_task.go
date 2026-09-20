package alicloud

import (
	"fmt"
	"log"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAlicloudCmsHybridMonitorFcTask() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudCmsHybridMonitorFcTaskCreate,
		Read:   resourceAlicloudCmsHybridMonitorFcTaskRead,
		Update: resourceAlicloudCmsHybridMonitorFcTaskUpdate,
		Delete: resourceAlicloudCmsHybridMonitorFcTaskDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(2 * time.Minute),
			Delete: schema.DefaultTimeout(2 * time.Minute),
			Update: schema.DefaultTimeout(2 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"hybrid_monitor_fc_task_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"namespace": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"target_user_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"yarm_config": {
				Type:     schema.TypeString,
				Required: true,
				StateFunc: func(v interface{}) string {
					yamlString, _ := normalizeYamlString(v)
					return yamlString
				},
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					equal, _ := compareCmsHybridMonitorFcTaskYamlConfigAreEquivalent(old, new)
					return equal
				},
				ValidateFunc: validateYamlString,
			},
		},
	}
}

func resourceAlicloudCmsHybridMonitorFcTaskCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "CreateHybridMonitorTask"
	request := make(map[string]interface{})
	var err error
	request["CollectTargetType"] = "aliyun_fc"
	request["Namespace"] = d.Get("namespace")
	if v, ok := d.GetOk("target_user_id"); ok {
		request["TargetUserId"] = v
	}
	request["TaskType"] = "aliyun_fc"
	request["YARMConfig"] = d.Get("yarm_config")
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.RpcPost("Cms", "2019-01-01", action, nil, request, false)
		if err != nil {
			if IsExpectedErrors(err, []string{"InternalError", "undefined"}) || NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_cms_hybrid_monitor_fc_task", action, AlibabaCloudSdkGoERROR)
	}
	if fmt.Sprint(response["Success"]) == "false" {
		return WrapError(fmt.Errorf("%s failed, response: %v", action, response))
	}

	d.SetId(fmt.Sprintf("%s:%s", response["TaskId"], request["Namespace"]))

	return resourceAlicloudCmsHybridMonitorFcTaskRead(d, meta)
}
func resourceAlicloudCmsHybridMonitorFcTaskRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	cmsService := CmsService{client}
	object, err := cmsService.DescribeCmsHybridMonitorFcTask(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_cms_hybrid_monitor_fc_task cmsService.DescribeCmsHybridMonitorFcTask Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}
	d.Set("hybrid_monitor_fc_task_id", parts[0])
	d.Set("namespace", parts[1])
	d.Set("target_user_id", object["TargetUserId"])
	d.Set("yarm_config", object["YARMConfig"])
	return nil
}
func resourceAlicloudCmsHybridMonitorFcTaskUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}

	if d.HasChange("yarm_config") {
		parts, err := ParseResourceId(d.Id(), 2)
		if err != nil {
			return WrapError(err)
		}
		// ModifyHybridMonitorTask does not accept YARMConfig, so the only way to
		// apply a yarm_config change is to destroy the old task and create a new
		// one with the updated config, then refresh the Terraform state ID.
		wait := incrementalWait(3*time.Second, 5*time.Second)

		// 1. Delete the old monitoring task.
		deleteAction := "DeleteHybridMonitorTask"
		deleteRequest := map[string]interface{}{
			"Namespace": parts[1],
		}
		if v, ok := d.GetOk("target_user_id"); ok {
			deleteRequest["TargetUserId"] = v
		}
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("Cms", "2019-01-01", deleteAction, nil, deleteRequest, false)
			if err != nil {
				if IsExpectedErrors(err, []string{"InternalError"}) || NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})
		addDebug(deleteAction, response, deleteRequest)
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), deleteAction, AlibabaCloudSdkGoERROR)
		}
		if fmt.Sprint(response["Success"]) == "false" {
			return WrapError(fmt.Errorf("%s failed, response: %v", deleteAction, response))
		}

		// 2. Create a new monitoring task with the updated yarm_config.
		createAction := "CreateHybridMonitorTask"
		createRequest := map[string]interface{}{
			"CollectTargetType": "aliyun_fc",
			"Namespace":         d.Get("namespace"),
			"TaskType":          "aliyun_fc",
			"YARMConfig":        d.Get("yarm_config"),
		}
		if v, ok := d.GetOk("target_user_id"); ok {
			createRequest["TargetUserId"] = v
		}
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("Cms", "2019-01-01", createAction, nil, createRequest, false)
			if err != nil {
				if IsExpectedErrors(err, []string{"InternalError", "undefined"}) || NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})
		addDebug(createAction, response, createRequest)
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), createAction, AlibabaCloudSdkGoERROR)
		}
		if fmt.Sprint(response["Success"]) == "false" {
			return WrapError(fmt.Errorf("%s failed, response: %v", createAction, response))
		}

		// 3. Refresh the state ID so it points to the newly created task.
		d.SetId(fmt.Sprintf("%s:%s", response["TaskId"], createRequest["Namespace"]))
	}

	return resourceAlicloudCmsHybridMonitorFcTaskRead(d, meta)
}
func resourceAlicloudCmsHybridMonitorFcTaskDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}
	action := "DeleteHybridMonitorTask"
	var response map[string]interface{}
	request := map[string]interface{}{
		"Namespace": parts[1],
	}
	if v, ok := d.GetOk("target_user_id"); ok {
		request["TargetUserId"] = v
	}
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = client.RpcPost("Cms", "2019-01-01", action, nil, request, false)
		if err != nil {
			if IsExpectedErrors(err, []string{"InternalError"}) || NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}
	if IsExpectedErrorCodes(fmt.Sprint(response["Code"]), []string{"ResourceNotFound"}) {
		return nil
	}
	if fmt.Sprint(response["Success"]) == "false" {
		return WrapError(fmt.Errorf("%s failed, response: %v", action, response))
	}
	return nil
}
