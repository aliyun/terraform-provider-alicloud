package alicloud

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceAlicloudEcsInvocation() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudEcsInvocationCreate,
		Read:   resourceAlicloudEcsInvocationRead,
		Delete: resourceAlicloudEcsInvocationDelete,
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(1 * time.Minute),
			Delete: schema.DefaultTimeout(1 * time.Minute),
		},
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"command_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"instance_id": {
				Type:          schema.TypeList,
				Optional:      true,
				Computed:      true,
				Elem:          &schema.Schema{Type: schema.TypeString},
				ForceNew:      true,
				ConflictsWith: []string{"resource_tag"},
				AtLeastOneOf:  []string{"instance_id", "resource_tag"},
			},
			"repeat_mode": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Computed:     true,
				ValidateFunc: validation.StringInSlice([]string{"Once", "Period", "NextRebootOnly", "EveryReboot"}, false),
			},
			"timed": {
				Type:       schema.TypeBool,
				Optional:   true,
				ForceNew:   true,
				Computed:   true,
				Deprecated: "Field `timed` has been deprecated because it is deprecated in the ECS API and takes no effect. To run a command periodically, set `repeat_mode` to `Period` and specify `frequency`.",
			},
			"frequency": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return d.Get("repeat_mode").(string) != "Period"
				},
			},
			// lintignore: S006
			"parameters": {
				Type:     schema.TypeMap,
				Optional: true,
				ForceNew: true,
			},
			"username": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"windows_password_name": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"resource_tag": {
				Type:          schema.TypeSet,
				Optional:      true,
				ForceNew:      true,
				ConflictsWith: []string{"instance_id"},
				AtLeastOneOf:  []string{"instance_id", "resource_tag"},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
						"value": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
					},
				},
			},
			// resource_group_id is a create-only field. InvokeCommand accepts it
			// and DescribeInvocations uses it as a query filter, but the Invocation
			// response object does not return it, so it cannot be read back. It is
			// not Computed to avoid Set(nil) drift that would force replacement on
			// every refresh; the configured value is preserved across reads and the
			// field is not recoverable on import.
			"resource_group_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"oss_output_delivery": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"working_dir": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceAlicloudEcsInvocationCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "InvokeCommand"
	request := make(map[string]interface{})
	var err error
	request["CommandId"] = d.Get("command_id")
	if v, ok := d.GetOk("instance_id"); ok {
		request["InstanceId"] = v
	}
	if v, ok := d.GetOk("resource_tag"); ok {
		resourceTags := make([]map[string]interface{}, 0)
		for _, item := range v.(*schema.Set).List() {
			m := item.(map[string]interface{})
			resourceTags = append(resourceTags, map[string]interface{}{
				"Key":   m["key"],
				"Value": m["value"],
			})
		}
		request["ResourceTag"] = resourceTags
	}
	if v, ok := d.GetOk("resource_group_id"); ok {
		request["ResourceGroupId"] = v
	}
	if v, ok := d.GetOk("oss_output_delivery"); ok {
		request["OssOutputDelivery"] = v
	}
	if v, ok := d.GetOk("working_dir"); ok {
		request["WorkingDir"] = v
	}
	request["RegionId"] = client.RegionId
	if v, ok := d.GetOk("repeat_mode"); ok {
		request["RepeatMode"] = v
	}
	if v, ok := d.GetOkExists("timed"); ok {
		request["Timed"] = v
	}
	if v, ok := d.GetOk("frequency"); ok {
		request["Frequency"] = v
	}
	if v, ok := d.GetOk("parameters"); ok {
		parameters, err := convertMaptoJsonString(v.(map[string]interface{}))
		if err != nil {
			return WrapError(err)
		}
		request["Parameters"] = parameters
	}
	if v, ok := d.GetOk("username"); ok {
		request["Username"] = v
	}
	if v, ok := d.GetOk("windows_password_name"); ok {
		request["WindowsPasswordName"] = v
	}
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
		response, err = client.RpcPost("Ecs", "2014-05-26", action, nil, request, false)
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_ecs_invocation", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(response["InvokeId"]))

	ecsService := EcsService{client}
	stateConf := BuildStateConf([]string{}, []string{"Scheduled", "Success"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, ecsService.EcsInvocationStateRefreshFunc(d.Id(), []string{"Failed", "PartialFailed", "Stopped"}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return resourceAlicloudEcsInvocationRead(d, meta)
}
func resourceAlicloudEcsInvocationRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	ecsService := EcsService{client}
	object, err := ecsService.DescribeEcsInvocation(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_ecs_invocation ecsService.DescribeEcsInvocation Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	d.Set("command_id", object["CommandId"])
	d.Set("frequency", object["Frequency"])
	d.Set("repeat_mode", object["RepeatMode"])

	parametersMap := make(map[string]string)
	err = json.Unmarshal([]byte(object["Parameters"].(string)), &parametersMap)
	if err != nil {
		return WrapError(err)
	}

	d.Set("parameters", parametersMap)
	d.Set("timed", object["Timed"])
	d.Set("username", object["Username"])
	d.Set("status", object["InvocationStatus"])
	instanceIdItems := make([]string, 0)
	if invokeInstances, ok := object["InvokeInstances"]; ok && invokeInstances != nil {
		if invokeInstance, ok := invokeInstances.(map[string]interface{})["InvokeInstance"]; ok && invokeInstance != nil {
			for _, invokeInstanceItem := range invokeInstance.([]interface{}) {
				if instanceId, ok := invokeInstanceItem.(map[string]interface{})["InstanceId"]; ok && instanceId != nil {
					instanceIdItems = append(instanceIdItems, fmt.Sprint(instanceId))
				}
			}
		}

	}
	// Absorb API response ordering drift for instance_id. The DescribeInvocations
	// response returns InvokeInstances in non-deterministic order, but instance_id
	// is a TypeList (order-sensitive) ForceNew field, so an unchanged membership
	// returned in a different order would otherwise force replacement on every
	// refresh. Reorder the response to follow the ordering already present in the
	// prior state; a refresh against the same membership then becomes a no-op.
	// When there is no prior ordering (first read / ID-only import) the response
	// order is kept and the field is backfilled normally — no GetOk guard that
	// would suppress backfill. Real member additions/deletions are still detected
	// because the reordered result differs from the prior state, and an explicit
	// change to the configured List order is preserved because the reference is
	// the prior state, not the new config.
	referenceItems, _ := d.Get("instance_id").([]interface{})
	reference := make([]string, 0, len(referenceItems))
	for _, item := range referenceItems {
		reference = append(reference, fmt.Sprint(item))
	}
	if err := d.Set("instance_id", reorderInstancesByReference(instanceIdItems, reference)); err != nil {
		return WrapError(err)
	}
	// resource_group_id is create-only and absent from the DescribeInvocations
	// response; preserve the configured value instead of clearing it. Setting it
	// to nil on a ForceNew field would otherwise trigger replacement on refresh.
	// working_dir and oss_output_delivery are response properties; only set them
	// when the response actually carries the key and handle Set errors.
	if v, ok := object["ResourceGroupId"]; ok && v != nil {
		if err := d.Set("resource_group_id", v); err != nil {
			return WrapError(err)
		}
	}
	if v, ok := object["WorkingDir"]; ok && v != nil {
		if err := d.Set("working_dir", v); err != nil {
			return WrapError(err)
		}
	}
	if v, ok := object["OssOutputDelivery"]; ok && v != nil {
		if err := d.Set("oss_output_delivery", v); err != nil {
			return WrapError(err)
		}
	}
	return nil
}

// reorderInstancesByReference reorders the instance IDs returned by
// DescribeInvocations to follow the ordering already recorded in the prior
// state (reference). It absorbs pure ordering drift so a refresh of an
// unchanged membership is a no-op, while still surfacing real changes:
//
//   - A reference ID still present in the response is emitted in reference order.
//   - A reference ID missing from the response (real deletion) is dropped.
//   - A response ID not in the reference (real addition) is appended in
//     response order.
//
// When the reference is empty (first read / ID-only import, where there is no
// prior user ordering to recover) the response is returned unchanged so the
// field is backfilled from the API instead of being left empty.
func reorderInstancesByReference(response, reference []string) []string {
	if len(reference) == 0 {
		return response
	}
	present := make(map[string]bool, len(response))
	for _, id := range response {
		present[id] = true
	}
	ordered := make([]string, 0, len(response))
	for _, id := range reference {
		if present[id] {
			ordered = append(ordered, id)
			delete(present, id)
		}
	}
	for _, id := range response {
		if present[id] {
			ordered = append(ordered, id)
		}
	}
	return ordered
}
func resourceAlicloudEcsInvocationDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "StopInvocation"
	request := make(map[string]interface{})
	var err error
	request["InvokeId"] = d.Id()
	request["RegionId"] = client.RegionId
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutDelete), func() *retry.RetryError {
		response, err = client.RpcPost("Ecs", "2014-05-26", action, nil, request, false)
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_ecs_invocation", action, AlibabaCloudSdkGoERROR)
	}

	log.Printf("[WARN] Cannot destroy resourceAlicloudEcsInvocation. Terraform will remove this resource from the state file, however resources may remain.")
	return nil
}
