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

func resourceAlicloudEcsImagePipelineExecution() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudEcsImagePipelineExecutionCreate,
		Read:   resourceAlicloudEcsImagePipelineExecutionRead,
		Delete: resourceAlicloudEcsImagePipelineExecutionDelete,
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
		},
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"image_pipeline_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceAlicloudEcsImagePipelineExecutionCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "StartImagePipelineExecution"
	request := make(map[string]interface{})
	conn, err := client.NewEcsClient()
	if err != nil {
		return WrapError(err)
	}
	request["ImagePipelineId"] = d.Get("image_pipeline_id")
	request["RegionId"] = client.RegionId

	request["ClientToken"] = buildClientToken("StartImagePipelineExecution")
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2014-05-26"), StringPointer("AK"), nil, request, &runtime)
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
	ecsService := EcsService{client}
	stateConf := BuildStateConf([]string{}, []string{"SUCCESS", "FAILED"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, ecsService.EcsImagePipelineExecutionStateRefreshFunc(d.Id(), []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return resourceAlicloudEcsImagePipelineExecutionRead(d, meta)
}
func resourceAlicloudEcsImagePipelineExecutionRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	ecsService := EcsService{client}
	object, err := ecsService.DescribeEcsImagePipelineExecution(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_ecs_image_pipeline_execution ecsService.DescribeEcsImagePipelineExecution Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	d.Set("image_pipeline_id", object["ImagePipelineId"])
	d.Set("status", object["Status"])
	return nil
}
func resourceAlicloudEcsImagePipelineExecutionDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[WARN] Cannot destroy resourceAlicloudEcsImagePipelineExecution. Terraform will remove this resource from the state file, however resources may remain.")
	return nil
}
