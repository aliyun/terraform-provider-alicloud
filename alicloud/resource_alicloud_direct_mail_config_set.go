// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
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
				Optional: true,
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
	if v, ok := d.GetOkExists("is_public_channel_backoff"); ok {
		request["IsPublicChannelBackoff"] = v
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
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.RpcPost("Dm", "2015-11-23", action, query, request, true)
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
		validationOptionRaw = detailRaw["ValidationOption"].(map[string]interface{})
	}
	if len(validationOptionRaw) > 0 {
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
	// ConfigSetUpdate follows full-replace semantics (the API marks Name as required):
	// fields omitted from the request are cleared server-side, so the flat mutable
	// fields must be sent on every update.
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

	if update && d.HasChange("validation_option") {
		validationOption := make(map[string]interface{})

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
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("Dm", "2015-11-23", action, query, request, true)
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
		if err := waitForDirectMailConfigSetUpdatePropagated(d, DirectMailServiceV2{client}); err != nil {
			return WrapError(err)
		}
	}

	return resourceAliCloudDirectMailConfigSetRead(d, meta)
}

// ConfigSetDetail is eventually consistent after ConfigSetUpdate; wait until the
// detail reflects the updated configuration before refreshing state. Propagation
// has been observed to take up to ~5 minutes, so allow a generous window.
func waitForDirectMailConfigSetUpdatePropagated(d *schema.ResourceData, service DirectMailServiceV2) error {
	return resource.Retry(10*time.Minute, func() *resource.RetryError {
		object, err := service.DescribeDirectMailConfigSet(d.Id())
		if err != nil {
			return resource.NonRetryableError(err)
		}
		detail, _ := object["Detail"].(map[string]interface{})
		if detail == nil {
			return resource.RetryableError(fmt.Errorf("waiting for ConfigSet %s update to propagate", d.Id()))
		}
		if fmt.Sprint(detail["Name"]) != d.Get("name") {
			return resource.RetryableError(fmt.Errorf("waiting for ConfigSet %s update to propagate", d.Id()))
		}
		if fmt.Sprint(detail["Description"]) != d.Get("description") {
			return resource.RetryableError(fmt.Errorf("waiting for ConfigSet %s update to propagate", d.Id()))
		}
		if fmt.Sprint(detail["IsPublicChannelBackoff"]) != fmt.Sprint(d.Get("is_public_channel_backoff")) {
			return resource.RetryableError(fmt.Errorf("waiting for ConfigSet %s update to propagate", d.Id()))
		}
		if ipPool, ok := detail["IpPool"].(map[string]interface{}); ok && ipPool != nil {
			if fmt.Sprint(ipPool["IpPoolId"]) != d.Get("ip_pool_id") {
				return resource.RetryableError(fmt.Errorf("waiting for ConfigSet %s update to propagate", d.Id()))
			}
		}
		return nil
	})
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

	if v, ok := d.GetOkExists("is_force"); ok {
		request["IsForce"] = v
	}
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = client.RpcPost("Dm", "2015-11-23", action, query, request, true)
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

	directMailServiceV2 := DirectMailServiceV2{client}
	return resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		_, err := directMailServiceV2.DescribeDirectMailConfigSet(d.Id())
		if err != nil {
			if NotFoundError(err) {
				return nil
			}
			return resource.NonRetryableError(err)
		}
		return resource.RetryableError(fmt.Errorf("waiting for ConfigSet %s to be deleted", d.Id()))
	})
}
