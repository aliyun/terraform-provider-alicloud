// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/PaesslerAG/jsonpath"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAliCloudArmsEnvCustomJob() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudArmsEnvCustomJobCreate,
		Read:   resourceAliCloudArmsEnvCustomJobRead,
		Update: resourceAliCloudArmsEnvCustomJobUpdate,
		Delete: resourceAliCloudArmsEnvCustomJobDelete,
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
			"config_yaml": {
				Type:     schema.TypeString,
				Required: true,
				StateFunc: func(v interface{}) string {
					yaml, _ := normalizeYamlString(v)
					return yaml
				},
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					equal, _ := compareYamlTemplateAreEquivalent(old, new)
					return equal
				},
				ValidateFunc: validateYamlString,
			},
			"env_custom_job_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"environment_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"status": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: StringInSlice([]string{"run", "stop"}, false),
			},
		},
	}
}

func resourceAliCloudArmsEnvCustomJobCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "CreateEnvCustomJob"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	query["EnvironmentId"] = d.Get("environment_id")
	query["CustomJobName"] = d.Get("env_custom_job_name")
	request["RegionId"] = client.RegionId

	request["ConfigYaml"] = d.Get("config_yaml")
	if v, ok := d.GetOk("aliyun_lang"); ok {
		request["AliyunLang"] = v
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_arms_env_custom_job", action, AlibabaCloudSdkGoERROR)
	}
	code, _ := jsonpath.Get("$.Code", response)
	if fmt.Sprint(code) != "200" {
		log.Printf("[DEBUG] Resource alicloud_arms_env_custom_job CreateEnvCustomJob Failed!!! %s", response)
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_arms_env_custom_job", action, AlibabaCloudSdkGoERROR, response)
	}

	d.SetId(fmt.Sprintf("%v:%v", query["EnvironmentId"], response["Data"]))

	return resourceAliCloudArmsEnvCustomJobUpdate(d, meta)
}

func resourceAliCloudArmsEnvCustomJobRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	armsServiceV2 := ArmsServiceV2{client}

	objectRaw, err := armsServiceV2.DescribeArmsEnvCustomJob(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_arms_env_custom_job DescribeArmsEnvCustomJob Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("config_yaml", objectRaw["ConfigYaml"])
	d.Set("status", objectRaw["Status"])
	d.Set("env_custom_job_name", objectRaw["CustomJobName"])
	d.Set("environment_id", objectRaw["EnvironmentId"])

	return nil
}

func resourceAliCloudArmsEnvCustomJobUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	update := false
	parts := strings.Split(d.Id(), ":")
	action := "UpdateEnvCustomJob"
	var err error
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	query["EnvironmentId"] = parts[0]
	query["CustomJobName"] = parts[1]
	request["RegionId"] = client.RegionId
	if !d.IsNewResource() && d.HasChange("config_yaml") {
		update = true
	}
	request["ConfigYaml"] = d.Get("config_yaml")
	if d.HasChange("status") {
		update = true
		request["Status"] = d.Get("status")
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
	}

	return resourceAliCloudArmsEnvCustomJobRead(d, meta)
}

func resourceAliCloudArmsEnvCustomJobDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	parts := strings.Split(d.Id(), ":")
	action := "DeleteEnvCustomJob"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	query["EnvironmentId"] = parts[0]
	query["CustomJobName"] = parts[1]
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
