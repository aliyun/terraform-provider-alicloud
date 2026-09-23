// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"log"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAliCloudCloudMonitorServiceAgentConfig() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudCloudMonitorServiceAgentConfigCreate,
		Read:   resourceAliCloudCloudMonitorServiceAgentConfigRead,
		Update: resourceAliCloudCloudMonitorServiceAgentConfigUpdate,
		Delete: resourceAliCloudCloudMonitorServiceAgentConfigDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"enable_install_agent_new_ecs": {
				Type:     schema.TypeBool,
				Optional: true,
			},
		},
	}
}

func resourceAliCloudCloudMonitorServiceAgentConfigCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "PutMonitoringConfig"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})

	if v, ok := d.GetOkExists("enable_install_agent_new_ecs"); ok {
		request["EnableInstallAgentNewECS"] = v
	}
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.RpcPost("Cms", "2019-01-01", action, query, request, true)
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_cloud_monitor_service_agent_config", action, AlibabaCloudSdkGoERROR)
	}

	accountId, err := client.AccountId()
	d.SetId(accountId)

	return resourceAliCloudCloudMonitorServiceAgentConfigRead(d, meta)
}

func resourceAliCloudCloudMonitorServiceAgentConfigRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	cloudMonitorServiceServiceV2 := CloudMonitorServiceServiceV2{client}

	objectRaw, err := cloudMonitorServiceServiceV2.DescribeCloudMonitorServiceAgentConfig(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_cloud_monitor_service_agent_config DescribeCloudMonitorServiceAgentConfig Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("enable_install_agent_new_ecs", objectRaw["EnableInstallAgentNewECS"])

	return nil
}

func resourceAliCloudCloudMonitorServiceAgentConfigUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	update := false

	var err error
	action := "PutMonitoringConfig"
	request = make(map[string]interface{})
	query = make(map[string]interface{})

	if d.HasChange("enable_install_agent_new_ecs") {
		update = true
		request["EnableInstallAgentNewECS"] = d.Get("enable_install_agent_new_ecs")
	}

	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("Cms", "2019-01-01", action, query, request, true)
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

	return resourceAliCloudCloudMonitorServiceAgentConfigRead(d, meta)
}

func resourceAliCloudCloudMonitorServiceAgentConfigDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[WARN] Cannot destroy resource AliCloud Resource Agent Config. Terraform will remove this resource from the state file, however resources may remain.")
	return nil
}
