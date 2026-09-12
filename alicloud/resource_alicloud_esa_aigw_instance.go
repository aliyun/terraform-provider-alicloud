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

func resourceAliCloudEsaAigwInstance() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudEsaAigwInstanceCreate,
		Read:   resourceAliCloudEsaAigwInstanceRead,
		Update: resourceAliCloudEsaAigwInstanceUpdate,
		Delete: resourceAliCloudEsaAigwInstanceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"aigw_instance_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"aigw_instance_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"auth_key": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"comment": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"enable_auth": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"record_count": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"update_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceAliCloudEsaAigwInstanceCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	action := "CreateAIGWInstance"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	request["Name"] = d.Get("aigw_instance_name")
	if v, ok := d.GetOk("comment"); ok && v.(string) != "" {
		request["Comment"] = v
	}
	request["ClientToken"] = buildClientToken(action)

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.RpcPost("ESA", "2024-09-10", action, query, request, true)
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_esa_aigw_instance", action, AlibabaCloudSdkGoERROR)
	}

	instanceId, err := jsonpath.Get("$.data.content.InstanceId", response)
	if err != nil {
		return WrapErrorf(err, FailedGetAttributeMsg, "alicloud_esa_aigw_instance", "$.data.content.InstanceId", response)
	}
	d.SetId(fmt.Sprint(instanceId))

	return resourceAliCloudEsaAigwInstanceRead(d, meta)
}

func resourceAliCloudEsaAigwInstanceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	esaServiceV2 := EsaServiceV2{client}

	objectRaw, err := esaServiceV2.DescribeEsaAigwInstance(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_esa_aigw_instance DescribeEsaAigwInstance Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("aigw_instance_id", objectRaw["InstanceId"])
	d.Set("aigw_instance_name", objectRaw["Name"])
	d.Set("auth_key", objectRaw["AuthKey"])
	d.Set("comment", objectRaw["Comment"])
	d.Set("create_time", objectRaw["CreateTime"])
	d.Set("enable_auth", objectRaw["AuthEnabled"])
	d.Set("record_count", objectRaw["RecordCount"])
	d.Set("status", objectRaw["Status"])
	d.Set("update_time", objectRaw["UpdateTime"])

	return nil
}

func resourceAliCloudEsaAigwInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	action := "UpdateAIGWInstance"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	update := false
	var err error
	request = make(map[string]interface{})
	request["InstanceId"] = d.Id()

	if d.HasChange("comment") {
		update = true
		request["Comment"] = d.Get("comment")
	}
	if d.HasChange("enable_auth") {
		update = true
		request["EnableAuth"] = d.Get("enable_auth")
	}

	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("ESA", "2024-09-10", action, query, request, true)
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
	}

	return resourceAliCloudEsaAigwInstanceRead(d, meta)
}

func resourceAliCloudEsaAigwInstanceDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	action := "DeleteAIGWInstance"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	request["InstanceId"] = d.Id()

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = client.RpcPost("ESA", "2024-09-10", action, query, request, true)
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

	return nil
}
