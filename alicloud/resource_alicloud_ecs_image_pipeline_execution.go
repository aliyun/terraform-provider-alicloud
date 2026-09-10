// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"log"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAliCloudEcsImagePipelineExecution() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudEcsImagePipelineExecutionCreate,
		Read:   resourceAliCloudEcsImagePipelineExecutionRead,
		Update: resourceAliCloudEcsImagePipelineExecutionUpdate,
		Delete: resourceAliCloudEcsImagePipelineExecutionDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"image_pipeline_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"status": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		},
	}
}

func resourceAliCloudEcsImagePipelineExecutionCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "StartImagePipelineExecution"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	request["RegionId"] = client.RegionId
	request["ClientToken"] = buildClientToken(action)

	request["ImagePipelineId"] = d.Get("image_pipeline_id")
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.RpcPost("Ecs", "2014-05-26", action, query, request, true)
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_ecs_image_pipeline_execution", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(response["ExecutionId"]))

	ecsServiceV2 := EcsServiceV2{client}
	stateConf := BuildStateConf([]string{}, []string{"SUCCESS", "BUILDING", "PREPARING"}, d.Timeout(schema.TimeoutCreate), 20*time.Second, ecsServiceV2.EcsImagePipelineExecutionStateRefreshFunc(d.Id(), "Status", []string{"FAILED"}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return resourceAliCloudEcsImagePipelineExecutionUpdate(d, meta)
}

func resourceAliCloudEcsImagePipelineExecutionRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	ecsServiceV2 := EcsServiceV2{client}

	objectRaw, err := ecsServiceV2.DescribeEcsImagePipelineExecution(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_ecs_image_pipeline_execution DescribeEcsImagePipelineExecution Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	if objectRaw["CreationTime"] != nil {
		d.Set("create_time", objectRaw["CreationTime"])
	}
	if objectRaw["ImagePipelineId"] != nil {
		d.Set("image_pipeline_id", objectRaw["ImagePipelineId"])
	}
	if objectRaw["Status"] != nil {
		d.Set("status", objectRaw["Status"])
	}

	return nil
}

func resourceAliCloudEcsImagePipelineExecutionUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}

	if d.HasChange("status") {
		ecsServiceV2 := EcsServiceV2{client}
		object, err := ecsServiceV2.DescribeEcsImagePipelineExecution(d.Id())
		if err != nil {
			return WrapError(err)
		}

		target := d.Get("status").(string)
		if object["Status"].(string) != target {
			if target == "CANCELLED" {
				action := "CancelImagePipelineExecution"
				request = make(map[string]interface{})
				query = make(map[string]interface{})
				request["ExecutionId"] = d.Id()
				request["RegionId"] = client.RegionId
				wait := incrementalWait(3*time.Second, 5*time.Second)
				err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
					response, err = client.RpcPost("Ecs", "2014-05-26", action, query, request, true)
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
					return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
				}
				ecsServiceV2 := EcsServiceV2{client}
				stateConf := BuildStateConf([]string{}, []string{"CANCELLED"}, d.Timeout(schema.TimeoutUpdate), 10*time.Second, ecsServiceV2.EcsImagePipelineExecutionStateRefreshFunc(d.Id(), "Status", []string{}))
				if _, err := stateConf.WaitForState(); err != nil {
					return WrapErrorf(err, IdMsg, d.Id())
				}

			}
		}
	}

	return resourceAliCloudEcsImagePipelineExecutionRead(d, meta)
}

func resourceAliCloudEcsImagePipelineExecutionDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[WARN] Cannot destroy resource AliCloud Resource Image Pipeline Execution. Terraform will remove this resource from the state file, however resources may remain.")
	return nil
}
