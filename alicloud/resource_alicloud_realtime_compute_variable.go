package alicloud

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAliCloudRealtimeComputeVariable() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudRealtimeComputeVariableCreate,
		Read:   resourceAliCloudRealtimeComputeVariableRead,
		Update: resourceAliCloudRealtimeComputeVariableUpdate,
		Delete: resourceAliCloudRealtimeComputeVariableDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"workspace": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"namespace": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"kind": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: StringInSlice([]string{"Plain"}, false),
			},
			"value": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"region_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceAliCloudRealtimeComputeVariableCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	namespace := d.Get("namespace")
	name := d.Get("name")
	action := fmt.Sprintf("/api/v2/namespaces/%s/variables", namespace)
	var request map[string]interface{}
	var response map[string]interface{}
	header := make(map[string]*string)
	query := make(map[string]*string)
	body := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	header["workspace"] = StringPointer(d.Get("workspace").(string))
	request["name"] = name
	request["kind"] = d.Get("kind")
	request["value"] = d.Get("value")
	if v, ok := d.GetOk("description"); ok {
		request["description"] = v
	}
	body = request
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.RoaPost("ververica", "2022-07-18", action, query, header, body, true)
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_realtime_compute_variable", action, AlibabaCloudSdkGoERROR)
	}

	if fmt.Sprint(response["success"]) == "false" {
		return WrapError(fmt.Errorf("%s failed, response: %v", action, response))
	}

	d.SetId(fmt.Sprintf("%v:%v:%v", d.Get("workspace").(string), namespace, name))

	return resourceAliCloudRealtimeComputeVariableRead(d, meta)
}

func resourceAliCloudRealtimeComputeVariableRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	realtimeComputeServiceV2 := RealtimeComputeServiceV2{client}

	objectRaw, err := realtimeComputeServiceV2.DescribeRealtimeComputeVariable(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_realtime_compute_variable DescribeRealtimeComputeVariable Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("kind", objectRaw["kind"])
	d.Set("value", objectRaw["value"])
	d.Set("description", objectRaw["description"])

	parts := strings.Split(d.Id(), ":")
	d.Set("workspace", parts[0])
	d.Set("namespace", parts[1])
	d.Set("name", parts[2])
	d.Set("region_id", client.RegionId)

	return nil
}

func resourceAliCloudRealtimeComputeVariableUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]*string
	var header map[string]*string
	var body map[string]interface{}
	update := false

	var err error
	parts := strings.Split(d.Id(), ":")
	namespace := parts[1]
	name := parts[2]
	action := fmt.Sprintf("/api/v2/namespaces/%s/variables/%s", namespace, name)
	request = make(map[string]interface{})
	query = make(map[string]*string)
	header = make(map[string]*string)
	body = make(map[string]interface{})
	header["workspace"] = StringPointer(parts[0])
	request["name"] = name

	if d.HasChange("value") || d.HasChange("description") {
		update = true
	}
	request["kind"] = d.Get("kind")
	request["value"] = d.Get("value")
	if v, ok := d.GetOk("description"); ok {
		request["description"] = v
	}
	body = request
	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RoaPatch("ververica", "2022-07-18", action, query, header, body, true)
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
		if fmt.Sprint(response["success"]) == "false" {
			return WrapError(fmt.Errorf("%s failed, response: %v", action, response))
		}
	}

	return resourceAliCloudRealtimeComputeVariableRead(d, meta)
}

func resourceAliCloudRealtimeComputeVariableDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	parts := strings.Split(d.Id(), ":")
	namespace := parts[1]
	name := parts[2]
	action := fmt.Sprintf("/api/v2/namespaces/%s/variables/%s", namespace, name)
	var request map[string]interface{}
	var response map[string]interface{}
	header := make(map[string]*string)
	query := make(map[string]*string)
	var err error
	request = make(map[string]interface{})
	header["workspace"] = StringPointer(parts[0])

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = client.RoaDelete("ververica", "2022-07-18", action, query, header, nil, true)
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
		if IsExpectedErrors(err, []string{"990301"}) || NotFoundError(err) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	return nil
}
