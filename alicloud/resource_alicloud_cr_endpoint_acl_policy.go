package alicloud

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
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
			// Description is a read-only legacy field in the corrected CloudSpec
			// model (@readonly): the per-entry comment now lives in entries[*].
			// comment, and this top-level field is populated from the first
			// AclEntry's Comment for backward-compatible display only.
			"description": {
				Type:     schema.TypeString,
				Computed: true,
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
			"entries": {
				Type:     schema.TypeSet,
				Required: true,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"entry": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
						"comment": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
					},
				},
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

// buildCrEndpointAclEntries expands the schema entries set into the
// AccessControlEntry-shaped [{Entry, Comment}] list that the CR API expects.
// includeComment controls whether Comment is attached: create sends it so the
// server records each entry's description; delete matches entries by Entry and
// omits Comment.
func buildCrEndpointAclEntries(d *schema.ResourceData, includeComment bool) []map[string]interface{} {
	entriesSet := d.Get("entries").(*schema.Set)
	entries := make([]map[string]interface{}, 0, entriesSet.Len())
	for _, e := range entriesSet.List() {
		em := e.(map[string]interface{})
		item := map[string]interface{}{"Entry": em["entry"]}
		if includeComment {
			if c, ok := em["comment"].(string); ok && c != "" {
				item["Comment"] = c
			}
		}
		entries = append(entries, item)
	}
	return entries
}

// expectedCrEndpointAclEntries returns the user-configured Entry values so the
// create path can wait until every entry has propagated into AclEntries.
func expectedCrEndpointAclEntries(d *schema.ResourceData) []string {
	entriesSet := d.Get("entries").(*schema.Set)
	expected := make([]string, 0, entriesSet.Len())
	for _, e := range entriesSet.List() {
		em := e.(map[string]interface{})
		expected = append(expected, fmt.Sprint(em["entry"]))
	}
	return expected
}

func resourceAlicloudCrEndpointAclPolicyCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	action := "CreateInstanceEndpointAclPolicy"
	// Entries is a jsonArray (style:json) parameter; marshal the expanded
	// AccessControlEntry list to a JSON string so it is sent verbatim on the
	// query regardless of the SDK's value coercion.
	entriesJSON, err := json.Marshal(buildCrEndpointAclEntries(d, true))
	if err != nil {
		return WrapError(err)
	}
	request := map[string]interface{}{
		"EndpointType": convertCrEndpointType(d.Get("endpoint_type").(string)),
		"InstanceId":   d.Get("instance_id"),
		"Entries":      string(entriesJSON),
	}
	if v, ok := d.GetOk("module_name"); ok {
		request["ModuleName"] = v
	}
	request["RegionId"] = client.RegionId

	// The internet endpoint must reach RUNNING before it accepts ACL policy
	// creation. When the endpoint is enabled in the same apply (via the
	// alicloud_cr_endpoint_acl_service data source or the
	// alicloud_cr_internet_endpoint resource), CreateInstanceEndpointAclPolicy
	// issued while the endpoint is still CREATING is silently dropped by the
	// server and the ACL entries never propagate into GetInstanceEndpoint's
	// AclEntries. Poll GetInstanceEndpoint until Status=RUNNING before issuing
	// the ACL creation so the entries reliably propagate.
	crService := CrService{client}
	if err := crService.WaitCrEndpointRunning(d.Get("instance_id").(string), d.Get("endpoint_type").(string), d.Timeout(schema.TimeoutCreate)); err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_cr_endpoint_acl_policy", "GetInstanceEndpoint", AlibabaCloudSdkGoERROR)
	}

	var response map[string]interface{}
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
		response, err = client.RpcPost("cr", "2018-12-01", action, nil, request, false)
		if err != nil {
			if IsExpectedErrors(err, []string{"SLB_SERVICE_ERROR"}) || NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_cr_endpoint_acl_policy", action, AlibabaCloudSdkGoERROR)
	}
	if v, ok := response["IsSuccess"]; !ok || fmt.Sprint(v) == "false" {
		return WrapError(fmt.Errorf("%s failed, response: %v", action, response))
	}

	// The resource identity is the (instance, endpoint) pair: one policy per
	// endpoint carrying the full Entries list, matching the corrected
	// CloudSpec identifyDefinition (InstanceId + EndpointType).
	d.SetId(fmt.Sprint(request["InstanceId"], ":", d.Get("endpoint_type").(string)))

	// CreateInstanceEndpointAclPolicy propagates ACL entries asynchronously;
	// now that the endpoint is RUNNING (gated above), poll GetInstanceEndpoint
	// until every created entry (matched by value) appears in AclEntries. This
	// prevents the subsequent Read from racing the propagation and misreading
	// the just-created resource as absent.
	if err := crService.WaitCrEndpointAclEntriesPropagate(d.Get("instance_id").(string), d.Get("endpoint_type").(string), expectedCrEndpointAclEntries(d), d.Timeout(schema.TimeoutCreate)); err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), "GetInstanceEndpoint", AlibabaCloudSdkGoERROR)
	}

	return resourceAlicloudCrEndpointAclPolicyRead(d, meta)
}

func resourceAlicloudCrEndpointAclPolicyRead(d *schema.ResourceData, meta interface{}) error {
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}
	if !isCrEndpointType(parts[1]) {
		return WrapError(fmt.Errorf("invalid endpoint_type %q in resource ID %q: the ID format is <instance_id>:<endpoint_type> (e.g. cri-abc123:internet)", parts[1], d.Id()))
	}
	client := meta.(*connectivity.AliyunClient)
	crService := CrService{client}
	object, err := crService.DescribeCrEndpointAclPolicy(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_cr_endpoint_acl_policy crService.DescribeCrEndpointAclPolicy Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	d.Set("endpoint_type", strings.ToLower(parts[1]))
	d.Set("instance_id", parts[0])

	// Map the server's AclEntries array back into the entries set. The
	// resource owns the user-managed entry list for the (instance, endpoint)
	// pair. GetInstanceEndpoint always returns an auto-added system loopback
	// ACL (entry 127.0.0.1/32, comment "default") once the endpoint is
	// enabled; it is not creatable/deletable via the ACL policy APIs, so
	// including it in state would create a perpetual plan diff against the
	// user's config — the sibling alicloud_cr_internet_endpoint resource
	// filters it for the same reason.
	aclEntries, _ := object["AclEntries"].([]interface{})
	entries := make([]map[string]interface{}, 0, len(aclEntries))
	for _, e := range aclEntries {
		em, ok := e.(map[string]interface{})
		if !ok {
			continue
		}
		if fmt.Sprint(em["Entry"]) == "127.0.0.1/32" && fmt.Sprint(em["Comment"]) == "default" {
			continue
		}
		entries = append(entries, map[string]interface{}{
			"entry":   em["Entry"],
			"comment": em["Comment"],
		})
	}
	d.Set("entries", entries)
	// Description is a read-only legacy field populated from the first
	// user-managed entry's Comment for backward-compatible display.
	if len(entries) > 0 {
		d.Set("description", entries[0]["comment"])
	}
	return nil
}

func resourceAlicloudCrEndpointAclPolicyDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}
	action := "DeleteInstanceEndpointAclPolicy"
	entriesJSON, err := json.Marshal(buildCrEndpointAclEntries(d, false))
	if err != nil {
		return WrapError(err)
	}
	request := map[string]interface{}{
		"EndpointType": convertCrEndpointType(parts[1]),
		"InstanceId":   parts[0],
		"Entries":      string(entriesJSON),
	}
	if v, ok := d.GetOk("module_name"); ok {
		request["ModuleName"] = v
	}

	// The CloudSpec delete mapping marks SLB_SERVICE_ERROR as a
	// resource-not-exist code: the entry cannot be removed (e.g. the system
	// loopback default 127.0.0.1/32) or is already gone, so it is treated as
	// deleted instead of being retried until the delete timeout.
	slbNotExist := false
	var response map[string]interface{}
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutDelete), func() *retry.RetryError {
		response, err = client.RpcPost("cr", "2018-12-01", action, nil, request, false)
		if err != nil {
			if IsExpectedErrors(err, []string{"SLB_SERVICE_ERROR"}) {
				slbNotExist = true
				return nil
			}
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
	if slbNotExist {
		return nil
	}
	if v, ok := response["IsSuccess"]; !ok || fmt.Sprint(v) == "false" {
		return WrapError(fmt.Errorf("%s failed, response: %v", action, response))
	}
	return nil
}
