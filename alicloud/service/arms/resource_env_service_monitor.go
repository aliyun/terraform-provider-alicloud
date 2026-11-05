// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package arms

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/validators"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceEnvServiceMonitor() *schema.Resource {
	return &schema.Resource{
		Create: ResourceEnvServiceMonitorCreate,
		Read:   ResourceEnvServiceMonitorRead,
		Update: ResourceEnvServiceMonitorUpdate,
		Delete: ResourceEnvServiceMonitorDelete,
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
					yaml, _ := validators.NormalizeYamlString(v)
					return yaml
				},
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					equal, _ := helper.CompareYamlTemplateAreEquivalent(old, new)
					return equal
				},
				ValidateFunc: validators.ValidateYamlString,
			},
			"env_service_monitor_name": {
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

func ResourceEnvServiceMonitorCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "CreateEnvServiceMonitor"
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
	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
		response, err = client.RpcPost("ARMS", "2019-08-08", action, query, request, true)

		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		logs.AddDebug(action, response, request)
		return nil
	})

	if err != nil {
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, "alicloud_arms_env_service_monitor", action, errs.AlibabaCloudSdkGoERROR)
	}
	code, _ := jsonpath.Get("$.Code", response)
	if fmt.Sprint(code) != "200" {
		log.Printf("[DEBUG] Resource alicloud_arms_env_service_monitor CreateEnvServiceMonitor Failed!!! %s", response)
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, "alicloud_arms_env_service_monitor", action, errs.AlibabaCloudSdkGoERROR, response)
	}

	namespace, _ := jsonpath.Get("$.Data.Namespace", response)
	serviceMonitorName, _ := jsonpath.Get("$.Data.ServiceMonitorName", response)
	d.SetId(fmt.Sprintf("%v:%v:%v", query["EnvironmentId"], namespace, serviceMonitorName))

	return ResourceEnvServiceMonitorRead(d, meta)
}

func ResourceEnvServiceMonitorRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	armsServiceV2 := ServiceV2{client}

	objectRaw, err := armsServiceV2.DescribeArmsEnvServiceMonitor(d.Id())
	if err != nil {
		if !d.IsNewResource() && errs.NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_arms_env_service_monitor DescribeArmsEnvServiceMonitor Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return errs.WrapError(err)
	}

	d.Set("config_yaml", objectRaw["ConfigYaml"])
	d.Set("status", objectRaw["Status"])
	d.Set("env_service_monitor_name", objectRaw["ServiceMonitorName"])
	d.Set("environment_id", objectRaw["EnvironmentId"])
	d.Set("namespace", objectRaw["Namespace"])

	return nil
}

func ResourceEnvServiceMonitorUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	update := false
	parts := strings.Split(d.Id(), ":")
	action := "UpdateEnvServiceMonitor"
	var err error
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	query["Namespace"] = parts[1]
	query["EnvironmentId"] = parts[0]
	query["ServiceMonitorName"] = parts[2]
	request["RegionId"] = client.RegionId
	if d.HasChange("config_yaml") {
		update = true
	}
	request["ConfigYaml"] = d.Get("config_yaml")
	if v, ok := d.GetOk("aliyun_lang"); ok {
		request["AliyunLang"] = v
	}
	if update {
		wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
		err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
			response, err = client.RpcPost("ARMS", "2019-08-08", action, query, request, true)

			if err != nil {
				if helper.NeedRetry(err) {
					wait()
					return retry.RetryableError(err)
				}
				return retry.NonRetryableError(err)
			}
			logs.AddDebug(action, response, request)
			return nil
		})
		if err != nil {
			return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
		}
	}

	return ResourceEnvServiceMonitorRead(d, meta)
}

func ResourceEnvServiceMonitorDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	parts := strings.Split(d.Id(), ":")
	action := "DeleteEnvServiceMonitor"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	query["Namespace"] = parts[1]
	query["EnvironmentId"] = parts[0]
	query["ServiceMonitorName"] = parts[2]
	request["RegionId"] = client.RegionId

	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutDelete), func() *retry.RetryError {
		response, err = client.RpcPost("ARMS", "2019-08-08", action, query, request, true)

		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		logs.AddDebug(action, response, request)
		return nil
	})

	if err != nil {
		if errs.IsExpectedErrors(err, []string{"404"}) {
			return nil
		}
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
	}

	return nil
}
