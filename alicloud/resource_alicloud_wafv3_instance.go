package alicloud

import (
	"fmt"
	"log"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAlicloudWafv3Instance() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudWafv3InstanceCreate,
		Read:   resourceAlicloudWafv3InstanceRead,
		Delete: resourceAlicloudWafv3InstanceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(1 * time.Minute),
			Delete: schema.DefaultTimeout(1 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"create_time": {
				Computed: true,
				Type:     schema.TypeString,
			},
			"instance_id": {
				Computed: true,
				Type:     schema.TypeString,
			},
			"status": {
				Computed: true,
				Type:     schema.TypeString,
			},
		},
	}
}

func resourceAlicloudWafv3InstanceCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	request := make(map[string]interface{})
	wafOpenapiService := WafOpenapiService{client}
	var err error

	var response map[string]interface{}
	action := "CreatePostpaidInstance"
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutCreate)), func() *resource.RetryError {
		response, err = client.RpcPost("waf-openapi", "2021-10-01", action, nil, request, false)
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_wafv3_instance", action, AlibabaCloudSdkGoERROR)
	}

	if v, err := jsonpath.Get("$.InstanceId", response); err != nil || v == nil {
		return WrapErrorf(err, IdMsg, "alicloud_wafv3_instance")
	} else {
		d.SetId(fmt.Sprint(v))
	}
	stateConf := BuildStateConf([]string{}, []string{"1"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, wafOpenapiService.Wafv3InstanceStateRefreshFunc(d.Id(), []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return resourceAlicloudWafv3InstanceRead(d, meta)
}

func resourceAlicloudWafv3InstanceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	wafOpenapiService := WafOpenapiService{client}

	object, err := wafOpenapiService.DescribeWafv3Instance(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_wafv3_instance wafOpenapiService.DescribeWafv3Instance Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	startTime91 := object["StartTime"]
	d.Set("create_time", startTime91)
	d.Set("instance_id", object["InstanceId"])

	status82 := object["Status"]
	d.Set("status", status82)

	return nil
}

func resourceAlicloudWafv3InstanceDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	wafOpenapiService := WafOpenapiService{client}
	var response map[string]interface{}
	var err error

	request := map[string]interface{}{

		"InstanceId": d.Id(),
		"RegionId":   client.RegionId,
	}

	action := "ReleaseInstance"
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutDelete)), func() *resource.RetryError {
		response, err = client.RpcPost("waf-openapi", "2021-10-01", action, nil, request, false)
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
		if NotFoundError(err) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}
	stateConf := BuildStateConf([]string{}, []string{}, d.Timeout(schema.TimeoutDelete), 5*time.Second, wafOpenapiService.Wafv3InstanceStateRefreshFunc(d.Id(), []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}
	return nil
}
