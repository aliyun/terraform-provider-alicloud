// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
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

func resourceAliCloudArmsEnvironment() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudArmsEnvironmentCreate,
		Read:   resourceAliCloudArmsEnvironmentRead,
		Update: resourceAliCloudArmsEnvironmentUpdate,
		Delete: resourceAliCloudArmsEnvironmentDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"aliyun_lang": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"bind_resource_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"drop_metrics": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"environment_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"environment_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"environment_sub_type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: StringInSlice([]string{"ECS", "ACK", "Cloud", "ManagedKubernetes", "Kubernetes", "ExternalKubernetes", "One"}, true),
			},
			"environment_type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: StringInSlice([]string{"ECS", "CS", "Cloud"}, true),
			},
			"managed_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Computed:     true,
				ValidateFunc: StringInSlice([]string{"none", "agent", "agent-exporter"}, true),
			},
			"resource_group_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"tags": tagsSchema(),
		},
	}
}

func resourceAliCloudArmsEnvironmentCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "CreateEnvironment"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	request["RegionId"] = client.RegionId

	if v, ok := d.GetOk("resource_group_id"); ok {
		request["ResourceGroupId"] = v
	}
	if v, ok := d.GetOk("environment_name"); ok {
		request["EnvironmentName"] = v
	}
	if v, ok := d.GetOk("tags"); ok {
		tagsMap := ConvertTags(v.(map[string]interface{}))
		request["Tags"] = tagsMap
	}

	request["EnvironmentType"] = d.Get("environment_type")
	if v, ok := d.GetOk("bind_resource_id"); ok {
		request["BindResourceId"] = v
	}
	request["EnvironmentSubType"] = d.Get("environment_sub_type")
	if v, ok := d.GetOk("aliyun_lang"); ok {
		request["AliyunLang"] = v
	}
	if v, ok := d.GetOk("managed_type"); ok {
		request["ManagedType"] = v
	}
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.RpcPost("ARMS", "2019-08-08", action, query, request, true)

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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_arms_environment", action, AlibabaCloudSdkGoERROR)
	}
	code, _ := jsonpath.Get("$.Code", response)
	if fmt.Sprint(code) != "200" {
		log.Printf("[DEBUG] Resource alicloud_arms_environment CreateEnvironment Failed!!! %s", response)
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_arms_environment", action, AlibabaCloudSdkGoERROR, response)
	}

	d.SetId(fmt.Sprint(response["Data"]))

	action = "InitEnvironment"
	request = make(map[string]interface{})
	query["EnvironmentId"] = fmt.Sprint(response["Data"])
	request["RegionId"] = client.RegionId

	if v, ok := d.GetOk("aliyun_lang"); ok {
		request["AliyunLang"] = v
	}
	if v, ok := d.GetOk("managed_type"); ok {
		request["ManagedType"] = v
	}
	wait = incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.RpcPost("ARMS", "2019-08-08", action, query, request, true)

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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_arms_environment", action, AlibabaCloudSdkGoERROR)
	}
	code, _ = jsonpath.Get("$.Code", response)
	if fmt.Sprint(code) != "200" {
		log.Printf("[DEBUG] Resource alicloud_arms_environment InitEnvironment Failed!!! %s", response)
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_arms_environment", action, AlibabaCloudSdkGoERROR, response)
	}

	d.SetId(fmt.Sprint(query["EnvironmentId"]))

	return resourceAliCloudArmsEnvironmentRead(d, meta)
}

func resourceAliCloudArmsEnvironmentRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	armsServiceV2 := ArmsServiceV2{client}

	objectRaw, err := armsServiceV2.DescribeArmsEnvironment(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_arms_environment DescribeArmsEnvironment Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("bind_resource_id", objectRaw["BindResourceId"])
	d.Set("environment_name", objectRaw["EnvironmentName"])
	d.Set("environment_sub_type", objectRaw["EnvironmentSubType"])
	d.Set("environment_type", objectRaw["EnvironmentType"])
	d.Set("managed_type", objectRaw["ManagedType"])
	d.Set("resource_group_id", objectRaw["ResourceGroupId"])
	d.Set("environment_id", objectRaw["EnvironmentId"])

	tagsMaps := objectRaw["Tags"]
	d.Set("tags", tagsToMap(tagsMaps))

	d.Set("environment_id", d.Id())

	return nil
}

func resourceAliCloudArmsEnvironmentUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	update := false
	d.Partial(true)
	action := "UpdateEnvironment"
	var err error
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	query["EnvironmentId"] = d.Id()
	request["RegionId"] = client.RegionId
	if d.HasChange("environment_name") {
		update = true
		request["EnvironmentName"] = d.Get("environment_name")
	}

	if v, ok := d.GetOk("aliyun_lang"); ok {
		request["AliyunLang"] = v
	}
	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("ARMS", "2019-08-08", action, query, request, true)

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
		d.SetPartial("environment_name")
	}
	update = false
	action = "ChangeResourceGroup"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	query["ResourceId"] = d.Id()
	request["RegionId"] = client.RegionId
	if _, ok := d.GetOk("resource_group_id"); ok && d.HasChange("resource_group_id") {
		update = true
		request["NewResourceGroupId"] = d.Get("resource_group_id")
	}

	request["ResourceType"] = "environment"
	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("ARMS", "2019-08-08", action, query, request, true)

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
		d.SetPartial("resource_group_id")
	}

	if d.HasChange("tags") {
		armsServiceV2 := ArmsServiceV2{client}
		if err := armsServiceV2.SetResourceTags(d, "environment"); err != nil {
			return WrapError(err)
		}
		d.SetPartial("tags")
	}
	d.Partial(false)
	return resourceAliCloudArmsEnvironmentRead(d, meta)
}

func resourceAliCloudArmsEnvironmentDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	action := "DeleteEnvironment"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	query["EnvironmentId"] = d.Id()
	query["DeletePromInstance"] = true
	request["RegionId"] = client.RegionId

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = client.RpcPost("ARMS", "2019-08-08", action, query, request, true)

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
		if IsExpectedErrors(err, []string{"404"}) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	return nil
}
