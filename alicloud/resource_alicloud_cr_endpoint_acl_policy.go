package alicloud

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceAlicloudCrEndpointAclPolicy() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudCrEndpointAclPolicyCreate,
		Read:   resourceAlicloudCrEndpointAclPolicyRead,
		Delete: resourceAlicloudCrEndpointAclPolicyDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"endpoint_type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				// The CR endpoint APIs document EndpointType as "Internet";
				// accept that spelling too and normalize to the lowercase
				// form used in state and resource IDs.
				StateFunc:    func(v interface{}) string { return strings.ToLower(v.(string)) },
				ValidateFunc: validation.StringInSlice([]string{"internet", "Internet"}, false),
			},
			"entry": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"module_name": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"Registry"}, false),
			},
		},
	}
}

func resourceAlicloudCrEndpointAclPolicyCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "CreateInstanceEndpointAclPolicy"
	request := make(map[string]interface{})
	var err error
	if v, ok := d.GetOk("description"); ok {
		request["Comment"] = v
	}
	request["EndpointType"] = convertCrEndpointType(d.Get("endpoint_type").(string))
	request["Entry"] = d.Get("entry")
	request["InstanceId"] = d.Get("instance_id")
	if v, ok := d.GetOk("module_name"); ok {
		request["ModuleName"] = v
	}
	request["RegionId"] = client.RegionId

	// The internet endpoint must reach RUNNING before it accepts ACL policy
	// creation. When the endpoint is enabled in the same apply (via the
	// alicloud_cr_endpoint_acl_service data source or the
	// alicloud_cr_internet_endpoint resource), CreateInstanceEndpointAclPolicy
	// issued while the endpoint is still CREATING is silently dropped by the
	// server and the ACL entry never propagates into GetInstanceEndpoint's
	// AclEntries. Poll GetInstanceEndpoint until Status=RUNNING before issuing
	// the ACL creation so the entry reliably propagates.
	crService := CrService{client}
	if err = crService.WaitCrEndpointRunning(d.Get("instance_id").(string), d.Get("endpoint_type").(string), d.Timeout(schema.TimeoutCreate)); err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_cr_instance_endpoint_acl_policy", "GetInstanceEndpoint", AlibabaCloudSdkGoERROR)
	}

	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.RpcPost("cr", "2018-12-01", action, nil, request, false)
		if err != nil {
			if IsExpectedErrors(err, []string{"SLB_SERVICE_ERROR"}) || NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_cr_instance_endpoint_acl_policy", action, AlibabaCloudSdkGoERROR)
	}

	if v, ok := response["IsSuccess"]; !ok || fmt.Sprint(v) == "false" {
		return WrapError(fmt.Errorf("%s failed, response: %v", action, response))
	}

	d.SetId(fmt.Sprint(request["InstanceId"], ":", d.Get("endpoint_type").(string), ":", request["Entry"]))

	// CreateInstanceEndpointAclPolicy propagates the ACL entry asynchronously;
	// now that the endpoint is RUNNING (gated above), poll GetInstanceEndpoint
	// until the new entry (matched by its value) appears in AclEntries. This
	// prevents the subsequent Read from racing the propagation and misreading
	// the just-created resource as absent.
	parts, err := ParseResourceId(d.Id(), 3)
	if err != nil {
		return WrapError(err)
	}
	if err = crService.WaitCrEndpointAclEntryPropagate(parts[0], parts[1], parts[2], d.Timeout(schema.TimeoutCreate)); err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), "GetInstanceEndpoint", AlibabaCloudSdkGoERROR)
	}

	return resourceAlicloudCrEndpointAclPolicyRead(d, meta)
}
func resourceAlicloudCrEndpointAclPolicyRead(d *schema.ResourceData, meta interface{}) error {
	parts, err := ParseResourceId(d.Id(), 3)
	if err != nil {
		return WrapError(err)
	}
	if !isCrEndpointType(parts[1]) {
		return WrapError(fmt.Errorf("invalid endpoint_type %q in resource ID %q: the ID format is <instance_id>:<endpoint_type>:<entry> (e.g. cri-abc123:internet:10.0.0.0/8)", parts[1], d.Id()))
	}
	client := meta.(*connectivity.AliyunClient)
	crService := CrService{client}
	object, err := crService.DescribeCrEndpointAclPolicy(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_cr_instance_endpoint_acl_policy crService.DescribeCrEndpointAclPolicy Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	d.Set("endpoint_type", strings.ToLower(parts[1]))
	d.Set("entry", parts[2])
	d.Set("instance_id", parts[0])
	d.Set("description", object["Comment"])
	return nil
}
func resourceAlicloudCrEndpointAclPolicyDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	parts, err := ParseResourceId(d.Id(), 3)
	if err != nil {
		return WrapError(err)
	}
	action := "DeleteInstanceEndpointAclPolicy"
	var response map[string]interface{}
	request := map[string]interface{}{
		"EndpointType": convertCrEndpointType(parts[1]),
		"Entry":        parts[2],
		"InstanceId":   parts[0],
	}

	if v, ok := d.GetOk("module_name"); ok {
		request["ModuleName"] = v
	}
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = client.RpcPost("cr", "2018-12-01", action, nil, request, false)
		if err != nil {
			if IsExpectedErrors(err, []string{"SLB_SERVICE_ERROR"}) || NeedRetry(err) {
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
	if v, ok := response["IsSuccess"]; !ok || fmt.Sprint(v) == "false" {
		return WrapError(fmt.Errorf("%s failed, response: %v", action, response))
	}
	return nil
}
