package alicloud

import (
	"fmt"
	"log"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAlicloudDataWorksComponent() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudDataWorksComponentCreate,
		Read:   resourceAlicloudDataWorksComponentRead,
		Update: resourceAlicloudDataWorksComponentUpdate,
		Delete: resourceAlicloudDataWorksComponentDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"component_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"project_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"spec": {
				Type:     schema.TypeString,
				Required: true,
			},
			"component_type": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"source": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
		},
	}
}

func resourceAlicloudDataWorksComponentCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "CreateComponent"
	request := make(map[string]interface{})
	var err error
	request["ProjectId"] = d.Get("project_id")
	request["Spec"] = d.Get("spec")
	if v, ok := d.GetOk("component_type"); ok {
		request["ComponentType"] = v
	}
	if v, ok := d.GetOk("source"); ok {
		request["Source"] = v
	}
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.RpcPost("dataworks-public", "2020-05-18", action, nil, request, false)
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_data_works_component", action, AlibabaCloudSdkGoERROR)
	}
	d.SetId(fmt.Sprint(response["ComponentId"], ":", request["ProjectId"]))
	return resourceAlicloudDataWorksComponentRead(d, meta)
}

func resourceAlicloudDataWorksComponentRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	dataworksPublicService := DataworksPublicService{client}
	object, err := dataworksPublicService.DescribeDataWorksComponent(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_data_works_component dataworksPublicService.DescribeDataWorksComponent Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}
	d.Set("component_id", parts[0])
	d.Set("project_id", parts[1])
	d.Set("spec", object["Spec"])
	d.Set("component_type", object["ComponentType"])
	d.Set("source", object["Source"])
	return nil
}

func resourceAlicloudDataWorksComponentUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}
	action := "UpdateComponent"
	request := map[string]interface{}{
		"ComponentId": parts[0],
		"ProjectId":   parts[1],
		"Spec":        d.Get("spec"),
	}
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
		response, err = client.RpcPost("dataworks-public", "2020-05-18", action, nil, request, false)
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
	return resourceAlicloudDataWorksComponentRead(d, meta)
}

func resourceAlicloudDataWorksComponentDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}
	action := "DeleteComponent"
	var response map[string]interface{}
	request := map[string]interface{}{
		"ComponentId": parts[0],
		"ProjectId":   parts[1],
	}
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = client.RpcPost("dataworks-public", "2020-05-18", action, nil, request, false)
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
	return nil
}
