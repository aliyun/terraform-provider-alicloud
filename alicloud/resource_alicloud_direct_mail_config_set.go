// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAliCloudDirectMailConfigSet() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudDirectMailConfigSetCreate,
		Read:   resourceAliCloudDirectMailConfigSetRead,
		Update: resourceAliCloudDirectMailConfigSetUpdate,
		Delete: resourceAliCloudDirectMailConfigSetDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"ip_pool_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"ip_pool_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"is_force": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"is_public_channel_backoff": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"validation_option": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"forbidden_sub_status_list": {
							Type:     schema.TypeList,
							Optional: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"enabled": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"forbidden_status_list": {
							Type:     schema.TypeList,
							Optional: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},
		},
	}
}

func resourceAliCloudDirectMailConfigSetCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "ConfigSetCreate"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	request["RegionId"] = client.RegionId

	if v, ok := d.GetOk("description"); ok {
		request["Description"] = v
	}
	if v, ok := d.GetOk("ip_pool_id"); ok {
		request["IpPoolId"] = v
	}
	// GetOkExists is unreliable for boolean zero values in plugin-sdk v2; use the
	// raw config to distinguish an explicit false from an omitted attribute.
	if rawConfig := d.GetRawConfig(); rawConfig != cty.NilVal && !rawConfig.IsNull() {
		if v, ok := rawConfig.AsValueMap()["is_public_channel_backoff"]; ok && !v.IsNull() && v.IsKnown() {
			request["IsPublicChannelBackoff"] = v.True()
		}
	}
	if v, ok := d.GetOk("name"); ok {
		request["Name"] = v
	}
	validationOption := make(map[string]interface{})

	if v := d.Get("validation_option"); !IsNil(v) {
		enabled1, _ := jsonpath.Get("$[0].enabled", v)
		if enabled1 != nil && enabled1 != "" {
			validationOption["Enabled"] = enabled1
		}
		forbiddenStatusList1, _ := jsonpath.Get("$[0].forbidden_status_list", v)
		if forbiddenStatusList1 != nil && forbiddenStatusList1 != "" {
			validationOption["ForbiddenStatusList"] = forbiddenStatusList1
		}
		forbiddenSubStatusList1, _ := jsonpath.Get("$[0].forbidden_sub_status_list", v)
		if forbiddenSubStatusList1 != nil && forbiddenSubStatusList1 != "" {
			validationOption["ForbiddenSubStatusList"] = forbiddenSubStatusList1
		}

		validationOptionJson, err := json.Marshal(validationOption)
		if err != nil {
			return WrapError(err)
		}
		request["ValidationOption"] = string(validationOptionJson)
	}

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
		response, err = client.RpcPost("Dm", "2015-11-23", action, query, request, true)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_direct_mail_config_set", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(response["Id"]))

	return resourceAliCloudDirectMailConfigSetRead(d, meta)
}

func resourceAliCloudDirectMailConfigSetRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	directMailServiceV2 := DirectMailServiceV2{client}

	objectRaw, err := directMailServiceV2.DescribeDirectMailConfigSet(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_direct_mail_config_set DescribeDirectMailConfigSet Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	detailRawObj, _ := jsonpath.Get("$.Detail", objectRaw)
	detailRaw := make(map[string]interface{})
	if detailRawObj != nil {
		detailRaw = detailRawObj.(map[string]interface{})
	}
	d.Set("description", detailRaw["Description"])
	d.Set("is_public_channel_backoff", detailRaw["IsPublicChannelBackoff"])
	d.Set("name", detailRaw["Name"])

	ipPoolRawObj, _ := jsonpath.Get("$.Detail.IpPool", objectRaw)
	ipPoolRaw := make(map[string]interface{})
	if ipPoolRawObj != nil {
		ipPoolRaw = ipPoolRawObj.(map[string]interface{})
	}
	d.Set("ip_pool_id", ipPoolRaw["IpPoolId"])
	d.Set("ip_pool_name", ipPoolRaw["IpPoolName"])

	validationOptionMaps := make([]map[string]interface{}, 0)
	validationOptionMap := make(map[string]interface{})
	validationOptionRaw := make(map[string]interface{})
	if detailRaw["ValidationOption"] != nil {
		if raw, ok := detailRaw["ValidationOption"].(map[string]interface{}); ok {
			validationOptionRaw = raw
		}
	}
	configuredOptions, _ := d.Get("validation_option").([]interface{})
	// Normalize the service's disabled empty option to an omitted configuration block.
	if len(validationOptionRaw) > 0 && (len(configuredOptions) > 0 || directMailValidationOptionKey(validationOptionRaw) != "false||") {
		validationOptionMap["enabled"] = validationOptionRaw["Enabled"]

		forbiddenStatusListRaw := make([]interface{}, 0)
		if validationOptionRaw["ForbiddenStatusList"] != nil {
			forbiddenStatusListRaw = convertToInterfaceArray(validationOptionRaw["ForbiddenStatusList"])
		}

		validationOptionMap["forbidden_status_list"] = forbiddenStatusListRaw
		forbiddenSubStatusListRaw := make([]interface{}, 0)
		if validationOptionRaw["ForbiddenSubStatusList"] != nil {
			forbiddenSubStatusListRaw = convertToInterfaceArray(validationOptionRaw["ForbiddenSubStatusList"])
		}

		validationOptionMap["forbidden_sub_status_list"] = forbiddenSubStatusListRaw
		validationOptionMaps = append(validationOptionMaps, validationOptionMap)
	}
	if err := d.Set("validation_option", validationOptionMaps); err != nil {
		return err
	}

	return nil
}

func resourceAliCloudDirectMailConfigSetUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	update := false

	var err error
	action := "ConfigSetUpdate"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["Id"] = d.Id()
	request["RegionId"] = client.RegionId
	// Resend current mutable values so unrelated updates preserve them.
	if d.HasChange("description") || d.HasChange("ip_pool_id") || d.HasChange("is_public_channel_backoff") || d.HasChange("name") || d.HasChange("validation_option") {
		update = true
	}
	if update {
		request["Name"] = d.Get("name")
		request["Description"] = d.Get("description")
		request["IsPublicChannelBackoff"] = d.Get("is_public_channel_backoff")
		if v, ok := d.GetOk("ip_pool_id"); ok {
			request["IpPoolId"] = v
		}
	}

	// An empty JSON object does not explicitly disable previously enabled validation.
	if update {
		validationOption := map[string]interface{}{
			"Enabled":                false,
			"ForbiddenStatusList":    []interface{}{},
			"ForbiddenSubStatusList": []interface{}{},
		}

		if v := d.Get("validation_option"); v != nil {
			enabled1, _ := jsonpath.Get("$[0].enabled", v)
			if enabled1 != nil && enabled1 != "" {
				validationOption["Enabled"] = enabled1
			}
			forbiddenStatusList1, _ := jsonpath.Get("$[0].forbidden_status_list", v)
			if forbiddenStatusList1 != nil && forbiddenStatusList1 != "" {
				validationOption["ForbiddenStatusList"] = forbiddenStatusList1
			}
			forbiddenSubStatusList1, _ := jsonpath.Get("$[0].forbidden_sub_status_list", v)
			if forbiddenSubStatusList1 != nil && forbiddenSubStatusList1 != "" {
				validationOption["ForbiddenSubStatusList"] = forbiddenSubStatusList1
			}
		}
		validationOptionJson, err := json.Marshal(validationOption)
		if err != nil {
			return WrapError(err)
		}
		request["ValidationOption"] = string(validationOptionJson)
	}

	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
			response, err = client.RpcPost("Dm", "2015-11-23", action, query, request, true)
			if err != nil {
				if NeedRetry(err) {
					wait()
					return retry.RetryableError(err)
				}
				return retry.NonRetryableError(err)
			}
			return nil
		})
		addDebug(action, response, request)
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
		if err := waitForDirectMailConfigSetUpdatePropagated(d, DirectMailServiceV2{client}); err != nil {
			return WrapError(err)
		}
	}

	return resourceAliCloudDirectMailConfigSetRead(d, meta)
}

// directMailDetailString returns the string form of a detail field, treating an
// absent or null field as the empty string: ConfigSetDetail omits unset fields
// (for example Description) instead of returning them with a zero value.
func directMailDetailString(detail map[string]interface{}, key string) string {
	if v, ok := detail[key]; ok && v != nil {
		return fmt.Sprint(v)
	}
	return ""
}

// ConfigSetDetail is eventually consistent after ConfigSetUpdate; wait until the
// detail reflects the updated configuration before refreshing state. Propagation
// has been observed to take up to ~5 minutes, so allow a generous window.
func waitForDirectMailConfigSetUpdatePropagated(d *schema.ResourceData, service DirectMailServiceV2) error {
	return retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
		object, err := service.DescribeDirectMailConfigSet(d.Id())
		if err != nil {
			return retry.NonRetryableError(err)
		}
		detail, _ := object["Detail"].(map[string]interface{})
		if detail == nil {
			return retry.RetryableError(fmt.Errorf("waiting for ConfigSet %s update to propagate", d.Id()))
		}
		if directMailDetailString(detail, "Name") != fmt.Sprint(d.Get("name")) {
			return retry.RetryableError(fmt.Errorf("waiting for ConfigSet %s update to propagate", d.Id()))
		}
		if directMailDetailString(detail, "Description") != fmt.Sprint(d.Get("description")) {
			return retry.RetryableError(fmt.Errorf("waiting for ConfigSet %s update to propagate", d.Id()))
		}
		// ConfigSetDetail omits the flag while the update is still propagating, so an
		// absent value must compare as false instead of skipping the comparison.
		backoff := false
		if v, ok := detail["IsPublicChannelBackoff"]; ok && v != nil {
			backoff = fmt.Sprint(v) == "true"
		}
		if backoff != d.Get("is_public_channel_backoff").(bool) {
			return retry.RetryableError(fmt.Errorf("waiting for ConfigSet %s update to propagate", d.Id()))
		}
		ipPool, _ := detail["IpPool"].(map[string]interface{})
		if directMailDetailString(ipPool, "IpPoolId") != fmt.Sprint(d.Get("ip_pool_id")) {
			return retry.RetryableError(fmt.Errorf("waiting for ConfigSet %s update to propagate", d.Id()))
		}
		validationOption, _ := detail["ValidationOption"].(map[string]interface{})
		want := directMailConfiguredValidationOptionKey(d)
		got := directMailValidationOptionKey(validationOption)
		if got != want {
			return retry.RetryableError(fmt.Errorf("waiting for ConfigSet %s validation options: expected %s, got %s", d.Id(), want, got))
		}
		return nil
	})
}

// List order must match Schema and Read when checking update propagation.
func directMailValidationOptionKey(option map[string]interface{}) string {
	enabled := "false"
	if v, ok := option["Enabled"]; ok && v != nil {
		enabled = strings.ToLower(fmt.Sprint(v))
	}
	return directMailValidationOptionKeyOf(enabled, option["ForbiddenStatusList"], option["ForbiddenSubStatusList"])
}

func directMailConfiguredValidationOptionKey(d *schema.ResourceData) string {
	options, _ := d.Get("validation_option").([]interface{})
	// Removed blocks can leave stale child attributes in the SDK diff, and an
	// empty block can surface as a nil or non-map element; both mean disabled.
	if len(options) == 0 || options[0] == nil {
		return "false||"
	}
	option, ok := options[0].(map[string]interface{})
	if !ok {
		return "false||"
	}
	return directMailValidationOptionKeyOf(fmt.Sprint(option["enabled"]), option["forbidden_status_list"], option["forbidden_sub_status_list"])
}

func directMailValidationOptionKeyOf(enabled string, statusList, subStatusList interface{}) string {
	normalize := func(raw interface{}) string {
		items := make([]string, 0)
		if raw != nil {
			for _, item := range convertToInterfaceArray(raw) {
				items = append(items, fmt.Sprint(item))
			}
		}
		return strings.Join(items, ",")
	}
	return enabled + "|" + normalize(statusList) + "|" + normalize(subStatusList)
}

func resourceAliCloudDirectMailConfigSetDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	action := "ConfigSetDelete"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	request["Ids"] = d.Id()
	request["RegionId"] = client.RegionId

	// IsForce uses the merged state value (which retains the last configured
	// value) rather than the raw config: the raw config is unavailable when a
	// removed block is destroyed. Explicit false is equivalent to the default.
	request["IsForce"] = d.Get("is_force")
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutDelete), func() *retry.RetryError {
		response, err = client.RpcPost("Dm", "2015-11-23", action, query, request, true)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
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

	directMailServiceV2 := DirectMailServiceV2{client}
	return retry.Retry(d.Timeout(schema.TimeoutDelete), func() *retry.RetryError {
		_, err := directMailServiceV2.DescribeDirectMailConfigSet(d.Id())
		if err != nil {
			if NotFoundError(err) {
				return nil
			}
			return retry.NonRetryableError(err)
		}
		return retry.RetryableError(fmt.Errorf("waiting for ConfigSet %s to be deleted", d.Id()))
	})
}
