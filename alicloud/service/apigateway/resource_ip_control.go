package apigateway

import (
	"fmt"
	"log"
	"regexp"
	"strings"
	"time"

	"github.com/PaesslerAG/jsonpath"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/validators"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// ResourceIpControl defines the alicloud_api_gateway_ip_control resource, which
// manages an IP access control (allowlist / blocklist) and its policy items on
// the API Gateway. The resource mirrors the lifecycle of CloudAPI/2016-07-14
// IpControl: the control itself is created/read/updated/deleted through the
// IpControl APIs, while the nested policy items are reconciled through the
// Add/Remove policy item operations (same pattern as the AccessControlList
// resource and its acl_entrys).
func ResourceIpControl() *schema.Resource {
	return &schema.Resource{
		Create: ResourceIpControlCreate,
		Read:   ResourceIpControlRead,
		Update: ResourceIpControlUpdate,
		Delete: ResourceIpControlDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"ip_control_name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validators.StringMatch(regexp.MustCompile(`^[A-Za-z0-9_-]+$`), "The name of the IP access control."),
			},
			"ip_control_type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validators.StringInSlice([]string{"ALLOW", "REFUSE"}, false),
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"ip_control_policy_items": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"policy_item_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"cidr_ip": {
							Type:     schema.TypeString,
							Required: true,
						},
						"app_id": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"region_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"modified_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func ResourceIpControlCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	action := "CreateIpControl"
	var request map[string]interface{}
	var response map[string]interface{}
	var err error
	request = make(map[string]interface{})

	request["IpControlName"] = d.Get("ip_control_name")
	request["IpControlType"] = d.Get("ip_control_type")
	if v, ok := d.GetOk("description"); ok {
		request["Description"] = v
	}

	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
		response, err = client.RpcPost("CloudAPI", "2016-07-14", action, nil, request, true)
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
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, "alicloud_api_gateway_ip_control", action, errs.AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(response["IpControlId"]))

	// The policy items are managed as an associated sub-collection, so they are
	// reconciled (added on first creation) through the Add operation.
	if err := resourceIpControlUpdatePolicyItems(d, meta); err != nil {
		return err
	}

	return ResourceIpControlRead(d, meta)
}

func ResourceIpControlRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	apiGatewayServiceV2 := ApiGatewayServiceV2{client}

	objectRaw, err := apiGatewayServiceV2.DescribeApiGatewayIpControl(d.Id())
	if err != nil {
		if !d.IsNewResource() && errs.NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_api_gateway_ip_control DescribeApiGatewayIpControl Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return errs.WrapError(err)
	}

	d.Set("ip_control_name", objectRaw["IpControlName"])
	d.Set("ip_control_type", objectRaw["IpControlType"])
	d.Set("description", objectRaw["Description"])
	d.Set("region_id", objectRaw["RegionId"])
	d.Set("create_time", objectRaw["CreateTime"])
	d.Set("modified_time", objectRaw["ModifiedTime"])

	policyItemsRaw, err := apiGatewayServiceV2.DescribeApiGatewayIpControlPolicyItems(d.Id())
	if err != nil {
		return errs.WrapError(err)
	}
	policyItemsMaps := make([]map[string]interface{}, 0)
	for _, item := range policyItemsRaw {
		itemMap := item.(map[string]interface{})
		m := make(map[string]interface{})
		m["policy_item_id"] = itemMap["PolicyItemId"]
		m["cidr_ip"] = itemMap["CidrIp"]
		// AppId is absent for REFUSE controls and for ALLOW items without an
		// app binding, so it is only set when the API actually returns it; an
		// absent key leaves the field at its zero value (""), matching a config
		// that omits app_id and keeping the set hash consistent.
		if v, ok := itemMap["AppId"]; ok && v != nil {
			m["app_id"] = v
		}
		policyItemsMaps = append(policyItemsMaps, m)
	}
	d.Set("ip_control_policy_items", policyItemsMaps)

	return nil
}

func ResourceIpControlUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	// IpControlType is ForceNew, so only the name and description are mutable on
	// the control itself; the policy items are reconciled separately below.
	if d.HasChange("ip_control_name") || d.HasChange("description") {
		action := "ModifyIpControl"
		var request map[string]interface{}
		var response map[string]interface{}
		var err error
		request = make(map[string]interface{})
		request["IpControlId"] = d.Id()
		request["IpControlName"] = d.Get("ip_control_name")
		request["Description"] = d.Get("description")

		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
		err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
			response, err = client.RpcPost("CloudAPI", "2016-07-14", action, nil, request, true)
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

	if err := resourceIpControlUpdatePolicyItems(d, meta); err != nil {
		return err
	}

	return ResourceIpControlRead(d, meta)
}

// resourceIpControlUpdatePolicyItems reconciles the nested policy items against
// the configured set: items no longer present are removed (batched by
// PolicyItemIds, semicolon separated) and new items are added one by one. The
// set is keyed by the user-specified cidr_ip/app_id, so a change to either
// manifests as a remove of the old item followed by an add of the new one,
// matching the AccessControlList acl_entrys behaviour.
func resourceIpControlUpdatePolicyItems(d *schema.ResourceData, meta interface{}) error {
	if !d.HasChange("ip_control_policy_items") {
		return nil
	}

	client := meta.(*connectivity.AliyunClient)
	oldEntry, newEntry := d.GetChange("ip_control_policy_items")
	oldEntrySet := oldEntry.(*schema.Set)
	newEntrySet := newEntry.(*schema.Set)
	removed := oldEntrySet.Difference(newEntrySet)
	added := newEntrySet.Difference(oldEntrySet)

	if removed.Len() > 0 {
		policyItemIds := make([]string, 0)
		for _, data := range removed.List() {
			m := data.(map[string]interface{})
			if pid := fmt.Sprint(m["policy_item_id"]); pid != "" && pid != "<nil>" {
				policyItemIds = append(policyItemIds, pid)
			}
		}
		if len(policyItemIds) > 0 {
			action := "RemoveIpControlPolicyItem"
			var request map[string]interface{}
			var response map[string]interface{}
			var err error
			request = make(map[string]interface{})
			request["IpControlId"] = d.Id()
			request["PolicyItemIds"] = strings.Join(policyItemIds, ";")

			runtime := util.RuntimeOptions{}
			runtime.SetAutoretry(true)
			wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
			err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
				response, err = client.RpcPost("CloudAPI", "2016-07-14", action, nil, request, true)
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
	}

	if added.Len() > 0 {
		action := "AddIpControlPolicyItem"
		for _, data := range added.List() {
			m := data.(map[string]interface{})
			var request map[string]interface{}
			var response map[string]interface{}
			var err error
			request = make(map[string]interface{})
			request["IpControlId"] = d.Id()
			request["CidrIp"] = m["cidr_ip"]
			if v := fmt.Sprint(m["app_id"]); v != "" && v != "<nil>" {
				request["AppId"] = v
			}

			runtime := util.RuntimeOptions{}
			runtime.SetAutoretry(true)
			wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
			err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
				response, err = client.RpcPost("CloudAPI", "2016-07-14", action, nil, request, true)
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
	}

	return nil
}

func ResourceIpControlDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	action := "DeleteIpControl"
	var request map[string]interface{}
	var response map[string]interface{}
	var err error
	request = make(map[string]interface{})
	request["IpControlId"] = d.Id()

	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutDelete), func() *retry.RetryError {
		response, err = client.RpcPost("CloudAPI", "2016-07-14", action, nil, request, true)

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
		if errs.IsExpectedErrors(err, []string{"NotFoundIpControl"}) {
			return nil
		}
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
	}

	return nil
}

// DescribeApiGatewayIpControl fetches a single IpControl by id. DescribeIpControls
// is a list API; passing IpControlId filters to the target control, so the
// matching entry is selected from the returned array.
func (s *ApiGatewayServiceV2) DescribeApiGatewayIpControl(id string) (object map[string]interface{}, err error) {
	client := s.client
	var request map[string]interface{}
	var response map[string]interface{}
	action := "DescribeIpControls"
	request = make(map[string]interface{})
	request["IpControlId"] = id
	request["PageNumber"] = 1
	request["PageSize"] = 10

	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(1*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("CloudAPI", "2016-07-14", action, nil, request, true)
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
		return object, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
	}

	v, err := jsonpath.Get("$.IpControlInfos.IpControlInfo[*]", response)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$.IpControlInfos.IpControlInfo[*]", response)
	}
	infos, ok := v.([]interface{})
	if !ok || len(infos) == 0 {
		return object, errs.WrapErrorf(errs.NotFoundErr("IpControl", id), errs.NotFoundMsg, response)
	}
	for _, raw := range infos {
		if fmt.Sprint(raw.(map[string]interface{})["IpControlId"]) == id {
			return raw.(map[string]interface{}), nil
		}
	}
	return object, errs.WrapErrorf(errs.NotFoundErr("IpControl", id), errs.NotFoundMsg, response)
}

// DescribeApiGatewayIpControlPolicyItems lists the policy items of an IpControl.
// An empty result (no items or absent array) is returned as an empty slice.
func (s *ApiGatewayServiceV2) DescribeApiGatewayIpControlPolicyItems(id string) (items []interface{}, err error) {
	client := s.client
	var request map[string]interface{}
	var response map[string]interface{}
	action := "DescribeIpControlPolicyItems"
	request = make(map[string]interface{})
	request["IpControlId"] = id
	request["PageNumber"] = 1
	request["PageSize"] = 100

	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(1*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("CloudAPI", "2016-07-14", action, nil, request, true)
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
		return items, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
	}

	v, e := jsonpath.Get("$.IpControlPolicyItems.IpControlPolicyItem[*]", response)
	if e != nil || v == nil {
		return items, nil
	}
	if list, ok := v.([]interface{}); ok {
		return list, nil
	}
	return items, nil
}
