// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"log"
	"sort"
	"strings"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAliCloudThreatDetectionMonitorAccount() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudThreatDetectionMonitorAccountCreate,
		Read:   resourceAliCloudThreatDetectionMonitorAccountRead,
		Update: resourceAliCloudThreatDetectionMonitorAccountUpdate,
		Delete: resourceAliCloudThreatDetectionMonitorAccountDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"account_ids": {
				Type:             schema.TypeString,
				Optional:         true,
				DiffSuppressFunc: monitorAccountIdsDiffSuppressFunc,
			},
		},
	}
}

// monitorAccountIdsDiffSuppressFunc compares two comma-separated account id lists as sets,
// because the backend returns the member account ids in its own order.
func monitorAccountIdsDiffSuppressFunc(k, old, new string, d *schema.ResourceData) bool {
	oldSet := make(map[string]struct{})
	for _, v := range strings.Split(old, ",") {
		v = strings.TrimSpace(v)
		if v != "" {
			oldSet[v] = struct{}{}
		}
	}
	newSet := make(map[string]struct{})
	for _, v := range strings.Split(new, ",") {
		v = strings.TrimSpace(v)
		if v != "" {
			newSet[v] = struct{}{}
		}
	}
	if len(oldSet) != len(newSet) {
		return false
	}
	for v := range oldSet {
		if _, ok := newSet[v]; !ok {
			return false
		}
	}
	return true
}

func resourceAliCloudThreatDetectionMonitorAccountCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	action := "CreateMonitorAccount"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})

	if v, ok := d.GetOk("account_ids"); ok {
		request["AccountIds"] = v
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_threat_detection_monitor_account", action, AlibabaCloudSdkGoERROR)
	}

	accountId, err := client.AccountId()
	if err != nil {
		return WrapError(err)
	}
	d.SetId(accountId)

	return resourceAliCloudThreatDetectionMonitorAccountRead(d, meta)
}

func resourceAliCloudThreatDetectionMonitorAccountRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	threatDetectionServiceV2 := ThreatDetectionServiceV2{client}

	objectRaw, err := threatDetectionServiceV2.DescribeThreatDetectionMonitorAccount(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_threat_detection_monitor_account DescribeThreatDetectionMonitorAccount Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	accountIds := make([]string, 0)
	if v, ok := objectRaw["AccountIds"]; ok && v != nil {
		for _, item := range v.([]interface{}) {
			accountIds = append(accountIds, fmt.Sprint(item))
		}
	}
	sort.Strings(accountIds)
	d.Set("account_ids", strings.Join(accountIds, ","))

	return nil
}

func resourceAliCloudThreatDetectionMonitorAccountUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	update := false

	var err error
	action := "CreateMonitorAccount"
	request = make(map[string]interface{})
	query = make(map[string]interface{})

	if d.HasChange("account_ids") {
		update = true
		request["AccountIds"] = d.Get("account_ids")
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

	return resourceAliCloudThreatDetectionMonitorAccountRead(d, meta)
}

func resourceAliCloudThreatDetectionMonitorAccountDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[WARN] Cannot destroy resource AliCloud Resource Monitor Account. Terraform will remove this resource from the state file, however resources may remain.")
	return nil
}
