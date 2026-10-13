// Package alicloud. Hand-written for the DataWorks CheckerInstance resource
// (CloudSpec @terraform enable:true with only a CheckFileDeployment operation;
// see resource_alicloud_dataworks_checker_instance_test.go for details).
package alicloud

import (
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

// resourceAliCloudDataworksCheckerInstance models alicloud_dataworks_checker_instance.
//
// The backing OpenAPI CheckFileDeployment (DataWorks public, v2020-05-18) is a
// callback-response operation: when a file submitted in the DataStudio UI
// enters the "publish-check" state, DataWorks emits a deployment-check event
// to the customer (carrying a CheckerInstanceId); the customer invokes
// CheckFileDeployment with that CheckerInstanceId and a Status of OK / FAIL /
// WARN to record their decision back into DataWorks.
//
// The product exposes no Create / Read / Delete / List API for the
// CheckerInstance resource itself, so this Terraform resource is a
// command-style resource:
//   - Create  = invoke CheckFileDeployment to record the decision;
//   - Read    = no-op (the config is the source of truth; the product offers
//     no way to query a persisted decision afterwards);
//   - Update  = re-invoke CheckFileDeployment when Status or CheckDetailUrl
//     changes (the API is idempotent on CheckerInstanceId);
//   - Delete  = no-op (removing the resource from Terraform state does not
//     retract the decision; DataWorks retains the last decision received for
//     a CheckerInstanceId until the next CheckFileDeployment call replaces
//     it).
//
// RegionId is sent in the gateway Host header by the SDK (cspec parameter
// `in: "host"`), so it is not added to the request body.
func resourceAliCloudDataworksCheckerInstance() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudDataworksCheckerInstanceCreate,
		Read:   resourceAliCloudDataworksCheckerInstanceRead,
		Update: resourceAliCloudDataworksCheckerInstanceUpdate,
		Delete: resourceAliCloudDataworksCheckerInstanceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"checker_instance_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"status": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					"OK", "FAIL", "WARN",
				}, false),
			},
			"check_detail_url": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"region_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceAliCloudDataworksCheckerInstanceCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	action := "CheckFileDeployment"
	request := map[string]interface{}{
		"CheckerInstanceId": d.Get("checker_instance_id").(string),
		"Status":            d.Get("status").(string),
	}
	if v, ok := d.GetOk("check_detail_url"); ok && v.(string) != "" {
		request["CheckDetailUrl"] = v.(string)
	}
	var response map[string]interface{}
	var err error
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.RpcPost("dataworks-public", "2020-05-18", action, nil, request, true)
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
		return WrapErrorf(err, DefaultErrorMsg, d.Get("checker_instance_id").(string), action, AlibabaCloudSdkGoERROR)
	}
	d.SetId(d.Get("checker_instance_id").(string))
	return resourceAliCloudDataworksCheckerInstanceRead(d, meta)
}

func resourceAliCloudDataworksCheckerInstanceRead(d *schema.ResourceData, meta interface{}) error {
	// No server-side Read API: keep the config-derived values in state and
	// only refresh the Computed region_id. This matches the contract that
	// CheckFileDeployment is a one-shot decision-recording API with no
	// persisted queryable state afterwards.
	client := meta.(*connectivity.AliyunClient)
	if err := d.Set("checker_instance_id", d.Id()); err != nil {
		return WrapError(err)
	}
	if err := d.Set("region_id", client.RegionId); err != nil {
		return WrapError(err)
	}
	return nil
}

func resourceAliCloudDataworksCheckerInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
	// CheckFileDeployment is the single operation that records the customer
	// decision for a given CheckerInstanceId; re-invoking it on every change
	// is the contract (the API is idempotent on CheckerInstanceId).
	client := meta.(*connectivity.AliyunClient)
	action := "CheckFileDeployment"
	request := map[string]interface{}{
		"CheckerInstanceId": d.Get("checker_instance_id").(string),
		"Status":            d.Get("status").(string),
	}
	if v, ok := d.GetOk("check_detail_url"); ok && v.(string) != "" {
		request["CheckDetailUrl"] = v.(string)
	}
	var response map[string]interface{}
	var err error
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
		response, err = client.RpcPost("dataworks-public", "2020-05-18", action, nil, request, true)
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
	return resourceAliCloudDataworksCheckerInstanceRead(d, meta)
}

func resourceAliCloudDataworksCheckerInstanceDelete(d *schema.ResourceData, meta interface{}) error {
	// No server-side Delete API. The Terraform resource represented the
	// customer's decision for a CheckerInstanceId, not a piece of cloud
	// infrastructure. Removing it from state does not require a server call;
	// DataWorks retains the last decision it received for the CheckerInstanceId
	// until the next CheckFileDeployment call replaces it.
	return nil
}
