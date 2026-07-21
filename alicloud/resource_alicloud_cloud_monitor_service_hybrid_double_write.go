package alicloud

import (
	"fmt"
	"log"
	"time"

	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAliCloudCloudMonitorServiceHybridDoubleWrite() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudCloudMonitorServiceHybridDoubleWriteCreate,
		Read:   resourceAliCloudCloudMonitorServiceHybridDoubleWriteRead,
		Delete: resourceAliCloudCloudMonitorServiceHybridDoubleWriteDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"source_namespace": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"source_user_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"namespace": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"user_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceAliCloudCloudMonitorServiceHybridDoubleWriteCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "CreateHybridDoubleWrite"
	request := make(map[string]interface{})
	var err error

	request["SourceNamespace"] = d.Get("source_namespace")
	request["SourceUserId"] = d.Get("source_user_id")
	request["Namespace"] = d.Get("namespace")
	request["UserId"] = d.Get("user_id")

	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutCreate)), func() *resource.RetryError {
		response, err = client.RpcPost("Cms", "2018-03-08", action, nil, request, false)
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_cloud_monitor_service_hybrid_double_write", action, AlibabaCloudSdkGoERROR)
	}
	d.SetId(fmt.Sprintf("%v:%v", request["SourceNamespace"], request["SourceUserId"]))

	return resourceAliCloudCloudMonitorServiceHybridDoubleWriteRead(d, meta)
}

func resourceAliCloudCloudMonitorServiceHybridDoubleWriteRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	cloudMonitorServiceServiceV2 := CloudMonitorServiceServiceV2{client}

	object, err := cloudMonitorServiceServiceV2.DescribeCloudMonitorServiceHybridDoubleWrite(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_cloud_monitor_service_hybrid_double_write DescribeCloudMonitorServiceHybridDoubleWrite Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("source_namespace", object["SourceNamespace"])
	d.Set("source_user_id", object["SourceUserId"])
	d.Set("namespace", object["Namespace"])
	d.Set("user_id", object["UserId"])

	return nil
}

func resourceAliCloudCloudMonitorServiceHybridDoubleWriteDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	action := "DeleteHybridDoubleWrite"
	var response map[string]interface{}

	var err error

	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}

	request := map[string]interface{}{
		"SourceNamespace": parts[0],
		"SourceUserId":    parts[1],
	}

	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutDelete)), func() *resource.RetryError {
		response, err = client.RpcPost("Cms", "2018-03-08", action, nil, request, false)
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
		if NotFoundError(err) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	if fmt.Sprint(response["Success"]) == "false" {
		return WrapError(fmt.Errorf("%s failed, response: %v", action, response))
	}

	return nil
}
