// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAliCloudCrScanRule() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudCrScanRuleCreate,
		Read:   resourceAliCloudCrScanRuleRead,
		Update: resourceAliCloudCrScanRuleUpdate,
		Delete: resourceAliCloudCrScanRuleDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"create_time": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"namespaces": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"repo_names": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"repo_tag_filter_pattern": {
				Type:     schema.TypeString,
				Required: true,
			},
			"rule_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"scan_rule_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"scan_scope": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: StringInSlice([]string{"NAMESPACE", "REPO", "INSTANCE"}, false),
			},
			"scan_type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"trigger_type": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: StringInSlice([]string{"MANUAL", "AUTO"}, false),
			},
		},
	}
}

func resourceAliCloudCrScanRuleCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "CreateScanRule"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	if v, ok := d.GetOk("instance_id"); ok {
		request["InstanceId"] = v
	}
	request["RegionId"] = client.RegionId

	request["RepoTagFilterPattern"] = d.Get("repo_tag_filter_pattern")
	if v, ok := d.GetOk("namespaces"); ok {
		namespacesMapsArray := convertToInterfaceArray(v)

		namespacesMapsJson, err := json.Marshal(namespacesMapsArray)
		if err != nil {
			return WrapError(err)
		}
		request["Namespaces"] = string(namespacesMapsJson)
	}

	request["TriggerType"] = d.Get("trigger_type")
	if v, ok := d.GetOk("repo_names"); ok {
		repoNamesMapsArray := convertToInterfaceArray(v)

		repoNamesMapsJson, err := json.Marshal(repoNamesMapsArray)
		if err != nil {
			return WrapError(err)
		}
		request["RepoNames"] = string(repoNamesMapsJson)
	}

	request["ScanType"] = d.Get("scan_type")
	request["RuleName"] = d.Get("rule_name")
	request["ScanScope"] = d.Get("scan_scope")
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.RpcPost("cr", "2018-12-01", action, query, request, true)
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_cr_scan_rule", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprintf("%v:%v", request["InstanceId"], response["ScanRuleId"]))

	return resourceAliCloudCrScanRuleRead(d, meta)
}

func resourceAliCloudCrScanRuleRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	crServiceV2 := CrServiceV2{client}

	objectRaw, err := crServiceV2.DescribeCrScanRule(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_cr_scan_rule DescribeCrScanRule Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("create_time", objectRaw["CreateTime"])
	d.Set("repo_tag_filter_pattern", objectRaw["RepoTagFilterPattern"])
	d.Set("rule_name", objectRaw["RuleName"])
	d.Set("scan_scope", objectRaw["ScanScope"])
	d.Set("scan_type", objectRaw["ScanType"])
	d.Set("trigger_type", objectRaw["TriggerType"])
	d.Set("instance_id", objectRaw["InstanceId"])
	d.Set("scan_rule_id", objectRaw["ScanRuleId"])

	namespacesRaw := make([]interface{}, 0)
	if objectRaw["Namespaces"] != nil {
		namespacesRaw = convertToInterfaceArray(objectRaw["Namespaces"])
	}

	d.Set("namespaces", namespacesRaw)
	repoNamesRaw := make([]interface{}, 0)
	if objectRaw["RepoNames"] != nil {
		repoNamesRaw = convertToInterfaceArray(objectRaw["RepoNames"])
	}

	d.Set("repo_names", repoNamesRaw)

	return nil
}

func resourceAliCloudCrScanRuleUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	update := false

	var err error
	parts := strings.Split(d.Id(), ":")
	action := "UpdateScanRule"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["InstanceId"] = parts[0]
	request["ScanRuleId"] = parts[1]
	request["RegionId"] = client.RegionId
	if d.HasChange("repo_tag_filter_pattern") {
		update = true
	}
	request["RepoTagFilterPattern"] = d.Get("repo_tag_filter_pattern")
	if d.HasChange("namespaces") {
		update = true
	}
	if v, ok := d.GetOk("namespaces"); ok || d.HasChange("namespaces") {
		namespacesMapsArray := convertToInterfaceArray(v)

		namespacesMapsJson, err := json.Marshal(namespacesMapsArray)
		if err != nil {
			return WrapError(err)
		}
		request["Namespaces"] = string(namespacesMapsJson)
	}

	if d.HasChange("trigger_type") {
		update = true
	}
	request["TriggerType"] = d.Get("trigger_type")
	if d.HasChange("repo_names") {
		update = true
	}
	if v, ok := d.GetOk("repo_names"); ok || d.HasChange("repo_names") {
		repoNamesMapsArray := convertToInterfaceArray(v)

		repoNamesMapsJson, err := json.Marshal(repoNamesMapsArray)
		if err != nil {
			return WrapError(err)
		}
		request["RepoNames"] = string(repoNamesMapsJson)
	}

	if d.HasChange("rule_name") {
		update = true
	}
	request["RuleName"] = d.Get("rule_name")
	if d.HasChange("scan_scope") {
		update = true
	}
	request["ScanScope"] = d.Get("scan_scope")
	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("cr", "2018-12-01", action, query, request, true)
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

	return resourceAliCloudCrScanRuleRead(d, meta)
}

func resourceAliCloudCrScanRuleDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	parts := strings.Split(d.Id(), ":")
	action := "DeleteScanRule"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	request["InstanceId"] = parts[0]
	request["ScanRuleId"] = parts[1]
	request["RegionId"] = client.RegionId

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = client.RpcPost("cr", "2018-12-01", action, query, request, true)
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
		if IsExpectedErrors(err, []string{"SCAN_RULE_NOT_EXIST"}) || NotFoundError(err) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	return nil
}
