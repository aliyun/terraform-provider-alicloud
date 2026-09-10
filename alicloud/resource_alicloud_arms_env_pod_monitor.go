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

func resourceAliCloudArmsEnvPodMonitor() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudArmsEnvPodMonitorCreate,
		Read:   resourceAliCloudArmsEnvPodMonitorRead,
		Update: resourceAliCloudArmsEnvPodMonitorUpdate,
		Delete: resourceAliCloudArmsEnvPodMonitorDelete,
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
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: StringInSlice([]string{"zh", "en"}, false),
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
			"env_pod_monitor_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"environment_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"namespace": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceAliCloudArmsEnvPodMonitorCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "CreateEnvPodMonitor"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	query["EnvironmentId"] = d.Get("environment_id")
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_arms_env_pod_monitor", action, AlibabaCloudSdkGoERROR)
	}
	code, _ := jsonpath.Get("$.Code", response)
	if fmt.Sprint(code) != "200" {
		log.Printf("[DEBUG] Resource alicloud_arms_env_pod_monitor CreateEnvPodMonitor Failed!!! %s", response)
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_arms_env_pod_monitor", action, AlibabaCloudSdkGoERROR, response)
	}

	namespace, _ := jsonpath.Get("$.Data.Namespace", response)
	podMonitorName, _ := jsonpath.Get("$.Data.PodMonitorName", response)
	d.SetId(fmt.Sprintf("%v:%v:%v", query["EnvironmentId"], namespace, podMonitorName))

	return resourceAliCloudArmsEnvPodMonitorRead(d, meta)
}

func resourceAliCloudArmsEnvPodMonitorRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	armsServiceV2 := ArmsServiceV2{client}

	objectRaw, err := armsServiceV2.DescribeArmsEnvPodMonitor(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_arms_env_pod_monitor DescribeArmsEnvPodMonitor Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("config_yaml", objectRaw["ConfigYaml"])
	d.Set("status", objectRaw["Status"])
	d.Set("env_pod_monitor_name", objectRaw["PodMonitorName"])
	d.Set("environment_id", objectRaw["EnvironmentId"])
	d.Set("namespace", objectRaw["Namespace"])

	return nil
}

func resourceAliCloudArmsEnvPodMonitorUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	update := false
	parts := strings.Split(d.Id(), ":")
	action := "UpdateEnvPodMonitor"
	var err error
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	query["Namespace"] = parts[1]
	query["EnvironmentId"] = parts[0]
	query["PodMonitorName"] = parts[2]
	request["RegionId"] = client.RegionId
	if d.HasChange("config_yaml") {
		update = true
	}
	request["ConfigYaml"] = d.Get("config_yaml")
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

	return resourceAliCloudArmsEnvPodMonitorRead(d, meta)
}

func resourceAliCloudArmsEnvPodMonitorDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	parts := strings.Split(d.Id(), ":")
	action := "DeleteEnvPodMonitor"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	query["Namespace"] = parts[1]
	query["EnvironmentId"] = parts[0]
	query["PodMonitorName"] = parts[2]
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
