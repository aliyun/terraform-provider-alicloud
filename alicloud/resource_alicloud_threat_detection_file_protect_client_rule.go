// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"log"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAliCloudThreatDetectionFileProtectClientRule() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudThreatDetectionFileProtectClientRuleCreate,
		Read:   resourceAliCloudThreatDetectionFileProtectClientRuleRead,
		Update: resourceAliCloudThreatDetectionFileProtectClientRuleUpdate,
		Delete: resourceAliCloudThreatDetectionFileProtectClientRuleDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"alert_level": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: IntInSlice([]int{0, 1, 2, 3}),
			},
			"exclude_users": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"file_ops": {
				Type:     schema.TypeList,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"file_paths": {
				Type:     schema.TypeList,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"file_types": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"platform": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: StringInSlice([]string{"windows", "linux"}, false),
			},
			"proc_paths": {
				Type:     schema.TypeList,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"rule_action": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: StringInSlice([]string{"pass", "monitor", "block"}, false),
			},
			"rule_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"status": {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: IntInSlice([]int{0, 1}),
			},
			"switch_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceAliCloudThreatDetectionFileProtectClientRuleCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "CreateFileProtectClientRule"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})

	request["ClientToken"] = buildClientToken(action)

	if v, ok := d.GetOk("proc_paths"); ok {
		procPathsMapsArray := convertToInterfaceArray(v)

		request["ProcPaths"] = procPathsMapsArray
	}

	if v, ok := d.GetOk("file_types"); ok {
		fileTypesMapsArray := convertToInterfaceArray(v)

		request["FileTypes"] = fileTypesMapsArray
	}

	if v, ok := d.GetOk("file_paths"); ok {
		filePathsMapsArray := convertToInterfaceArray(v)

		request["FilePaths"] = filePathsMapsArray
	}

	request["RuleAction"] = d.Get("rule_action")
	if v, ok := d.GetOkExists("alert_level"); ok {
		request["AlertLevel"] = v
	}
	request["Status"] = d.Get("status")
	request["RuleName"] = d.Get("rule_name")
	if v, ok := d.GetOk("platform"); ok {
		request["Platform"] = v
	}
	if v, ok := d.GetOk("exclude_users"); ok {
		excludeUsersMapsArray := convertToInterfaceArray(v)

		request["ExcludeUsers"] = excludeUsersMapsArray
	}
	if v, ok := d.GetOk("file_ops"); ok {
		fileOpsMapsArray := convertToInterfaceArray(v)

		request["FileOps"] = fileOpsMapsArray
	}

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.RpcPost("Sas", "2018-12-03", action, query, request, true)
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_threat_detection_file_protect_client_rule", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(response["Id"]))

	return resourceAliCloudThreatDetectionFileProtectClientRuleRead(d, meta)
}

func resourceAliCloudThreatDetectionFileProtectClientRuleRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	threatDetectionServiceV2 := ThreatDetectionServiceV2{client}

	objectRaw, err := threatDetectionServiceV2.DescribeThreatDetectionFileProtectClientRule(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_threat_detection_file_protect_client_rule DescribeThreatDetectionFileProtectClientRule Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("alert_level", objectRaw["AlertLevel"])
	d.Set("platform", objectRaw["Platform"])
	d.Set("rule_action", objectRaw["RuleAction"])
	d.Set("rule_name", objectRaw["RuleName"])
	d.Set("status", objectRaw["Status"])
	d.Set("switch_id", objectRaw["SwitchId"])

	excludeUsersRaw := make([]interface{}, 0)
	if objectRaw["ExcludeUsers"] != nil {
		excludeUsersRaw = convertToInterfaceArray(objectRaw["ExcludeUsers"])
	}

	d.Set("exclude_users", excludeUsersRaw)
	fileOpsRaw := make([]interface{}, 0)
	if objectRaw["FileOps"] != nil {
		fileOpsRaw = convertToInterfaceArray(objectRaw["FileOps"])
	}

	d.Set("file_ops", fileOpsRaw)
	filePathsRaw := make([]interface{}, 0)
	if objectRaw["FilePaths"] != nil {
		filePathsRaw = convertToInterfaceArray(objectRaw["FilePaths"])
	}

	d.Set("file_paths", filePathsRaw)
	fileTypesRaw := make([]interface{}, 0)
	if objectRaw["FileTypes"] != nil {
		fileTypesRaw = convertToInterfaceArray(objectRaw["FileTypes"])
	}

	d.Set("file_types", fileTypesRaw)
	procPathsRaw := make([]interface{}, 0)
	if objectRaw["ProcPaths"] != nil {
		procPathsRaw = convertToInterfaceArray(objectRaw["ProcPaths"])
	}

	d.Set("proc_paths", procPathsRaw)

	return nil
}

func resourceAliCloudThreatDetectionFileProtectClientRuleUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	update := false

	var err error
	action := "UpdateFileProtectClientRule"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["Id"] = d.Id()

	if d.HasChange("proc_paths") {
		update = true
	}
	if v, ok := d.GetOk("proc_paths"); ok || d.HasChange("proc_paths") {
		procPathsMapsArray := convertToInterfaceArray(v)

		request["ProcPaths"] = procPathsMapsArray
	}

	if d.HasChange("file_types") {
		update = true
		if v, ok := d.GetOk("file_types"); ok || d.HasChange("file_types") {
			fileTypesMapsArray := convertToInterfaceArray(v)

			request["FileTypes"] = fileTypesMapsArray
		}
	}

	if d.HasChange("file_paths") {
		update = true
	}
	if v, ok := d.GetOk("file_paths"); ok || d.HasChange("file_paths") {
		filePathsMapsArray := convertToInterfaceArray(v)

		request["FilePaths"] = filePathsMapsArray
	}

	if d.HasChange("rule_action") {
		update = true
	}
	request["RuleAction"] = d.Get("rule_action")
	if d.HasChange("alert_level") {
		update = true
		request["AlertLevel"] = d.Get("alert_level")
	}

	if d.HasChange("rule_name") {
		update = true
	}
	request["RuleName"] = d.Get("rule_name")
	request["Status"] = d.Get("status")
	if d.HasChange("exclude_users") {
		update = true
		if v, ok := d.GetOk("exclude_users"); ok || d.HasChange("exclude_users") {
			excludeUsersMapsArray := convertToInterfaceArray(v)

			request["ExcludeUsers"] = excludeUsersMapsArray
		}
	}

	if d.HasChange("file_ops") {
		update = true
	}
	if v, ok := d.GetOk("file_ops"); ok || d.HasChange("file_ops") {
		fileOpsMapsArray := convertToInterfaceArray(v)

		request["FileOps"] = fileOpsMapsArray
	}

	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("Sas", "2018-12-03", action, query, request, true)
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
	update = false
	action = "UpdateFileProtectClientRuleStatus"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["IdList.1"] = d.Id()
	request["SelectAll"] = false

	if d.HasChange("rule_action") {
		update = true
	}
	request["RuleAction"] = d.Get("rule_action")
	if d.HasChange("alert_level") {
		update = true
		request["AlertLevel"] = d.Get("alert_level")
	}

	if d.HasChange("status") {
		update = true
	}
	request["Status"] = d.Get("status")
	if d.HasChange("rule_name") {
		update = true
	}
	request["RuleName"] = d.Get("rule_name")
	if d.HasChange("platform") {
		update = true
		request["Platform"] = d.Get("platform")
	}

	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("Sas", "2018-12-03", action, query, request, true)
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

	return resourceAliCloudThreatDetectionFileProtectClientRuleRead(d, meta)
}

func resourceAliCloudThreatDetectionFileProtectClientRuleDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	action := "DeleteFileProtectClientRule"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	request["IdList.1"] = d.Id()
	request["SelectAll"] = false

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = client.RpcPost("Sas", "2018-12-03", action, query, request, true)
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
