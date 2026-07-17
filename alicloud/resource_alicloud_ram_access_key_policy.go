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
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

// ramAccessKeyPolicyResetDocument is the policy document used to clear the
// network access restriction policy on Delete. The SetAccessKeyPolicy API does
// NOT accept an empty document ("{}"), so an explicit disabled policy with no
// statements is sent to reset the access key to the unrestricted baseline.
const ramAccessKeyPolicyResetDocument = `{"Version":1,"Status":"Inactive","Statements":[]}`

// stripAccessKeyPolicyVersion removes the server-managed "Version" field so a
// user-supplied document (which usually omits it) compares equal to the value
// returned by GetAccessKeyPolicy (which always includes "Version":1).
func stripAccessKeyPolicyVersion(s string) string {
	var m map[string]interface{}
	if err := json.Unmarshal([]byte(s), &m); err != nil {
		return s
	}
	delete(m, "Version")
	out, err := json.Marshal(m)
	if err != nil {
		return s
	}
	return string(out)
}

// accessKeyPolicyEquivalent reports whether two access key policy documents are
// semantically equal, ignoring the server-managed "Version" field.
func accessKeyPolicyEquivalent(a, b string) bool {
	equal, _ := compareJsonTemplateAreEquivalent(stripAccessKeyPolicyVersion(a), stripAccessKeyPolicyVersion(b))
	return equal
}

// isEmptyAccessKeyPolicy reports whether a policy document represents the
// "no policy" baseline: an empty object "{}" (never configured), or a disabled
// policy carrying no statements (the state left behind after a reset/Delete).
func isEmptyAccessKeyPolicy(s string) bool {
	if strings.TrimSpace(s) == "" {
		return true
	}
	var m map[string]interface{}
	if err := json.Unmarshal([]byte(s), &m); err != nil {
		return false
	}
	if len(m) == 0 {
		return true
	}
	status, _ := m["Status"].(string)
	statements, hasStatements := m["Statements"].([]interface{})
	if strings.EqualFold(status, "Inactive") && (!hasStatements || len(statements) == 0) {
		return true
	}
	return false
}

// ramAccessKeyPolicyWaitConsistent polls GetAccessKeyPolicy until the returned
// document is semantically equivalent to the one just submitted via
// SetAccessKeyPolicy. The backend is eventually consistent: for a short window
// after a write, GetAccessKeyPolicy can keep returning the previous document,
// so reading immediately would persist a stale value into state and leave a
// non-empty plan after apply.
func ramAccessKeyPolicyWaitConsistent(client *connectivity.AliyunClient, id, expected string, timeout time.Duration) error {
	// A baseline/empty target can never converge to a "present" policy (it is
	// treated as "not exist"), so there is nothing to wait for.
	if isEmptyAccessKeyPolicy(expected) {
		return nil
	}
	ramServiceV2 := RamServiceV2{client}
	wait := incrementalWait(3*time.Second, 5*time.Second)
	return resource.Retry(timeout, func() *resource.RetryError {
		objectRaw, err := ramServiceV2.DescribeRamAccessKeyPolicy(id)
		if err != nil {
			if NotFoundError(err) {
				wait()
				return resource.RetryableError(fmt.Errorf("access key policy %s not converged yet (still empty)", id))
			}
			return resource.NonRetryableError(err)
		}
		if accessKeyPolicyEquivalent(fmt.Sprint(objectRaw["AccessKeyPolicy"]), expected) {
			return nil
		}
		wait()
		return resource.RetryableError(fmt.Errorf("access key policy %s not converged yet", id))
	})
}

func resourceAliCloudRamAccessKeyPolicy() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudRamAccessKeyPolicyCreate,
		Read:   resourceAliCloudRamAccessKeyPolicyRead,
		Update: resourceAliCloudRamAccessKeyPolicyUpdate,
		Delete: resourceAliCloudRamAccessKeyPolicyDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"user_access_key_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"user_principal_name": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"access_key_policy": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsJSON,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return accessKeyPolicyEquivalent(old, new)
				},
			},
		},
	}
}

// parseRamAccessKeyPolicyId splits the resource ID into the optional user
// principal name and the user access key id. The ID is composed as
// "<user_principal_name>:<user_access_key_id>" when a principal name is set,
// otherwise it is just "<user_access_key_id>". Both the RAM login name and the
// access key id never contain a colon, so ":" is a safe separator.
func parseRamAccessKeyPolicyId(id string) (userPrincipalName string, userAccessKeyId string) {
	parts := strings.SplitN(id, ":", 2)
	if len(parts) == 2 {
		return parts[0], parts[1]
	}
	return "", parts[0]
}

func resourceAliCloudRamAccessKeyPolicyCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "SetAccessKeyPolicy"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})

	userAccessKeyId := d.Get("user_access_key_id").(string)
	request["UserAccessKeyId"] = userAccessKeyId
	request["AccessKeyPolicy"] = d.Get("access_key_policy").(string)
	userPrincipalName := ""
	if v, ok := d.GetOk("user_principal_name"); ok {
		userPrincipalName = v.(string)
		request["UserPrincipalName"] = userPrincipalName
	}

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.RpcPost("Ims", "2019-08-15", action, query, request, true)
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_ram_access_key_policy", action, AlibabaCloudSdkGoERROR)
	}

	if userPrincipalName != "" {
		d.SetId(fmt.Sprintf("%s:%s", userPrincipalName, userAccessKeyId))
	} else {
		d.SetId(userAccessKeyId)
	}

	// Wait for the write to be readable before Read persists it to state,
	// otherwise the eventually-consistent GetAccessKeyPolicy may return a stale
	// document and leave a non-empty plan.
	if err := ramAccessKeyPolicyWaitConsistent(client, d.Id(), d.Get("access_key_policy").(string), d.Timeout(schema.TimeoutCreate)); err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	return resourceAliCloudRamAccessKeyPolicyRead(d, meta)
}

func resourceAliCloudRamAccessKeyPolicyRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	ramServiceV2 := RamServiceV2{client}

	userPrincipalName, userAccessKeyId := parseRamAccessKeyPolicyId(d.Id())

	objectRaw, err := ramServiceV2.DescribeRamAccessKeyPolicy(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_ram_access_key_policy DescribeRamAccessKeyPolicy Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("access_key_policy", objectRaw["AccessKeyPolicy"])
	d.Set("user_access_key_id", userAccessKeyId)
	if userPrincipalName != "" {
		d.Set("user_principal_name", userPrincipalName)
	}

	return nil
}

func resourceAliCloudRamAccessKeyPolicyUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	update := false

	var err error
	action := "SetAccessKeyPolicy"
	userPrincipalName, userAccessKeyId := parseRamAccessKeyPolicyId(d.Id())
	request = make(map[string]interface{})
	request["UserAccessKeyId"] = userAccessKeyId
	if userPrincipalName != "" {
		request["UserPrincipalName"] = userPrincipalName
	}

	if d.HasChange("access_key_policy") {
		update = true
	}
	request["AccessKeyPolicy"] = d.Get("access_key_policy").(string)

	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("Ims", "2019-08-15", action, query, request, true)
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

		// Wait for the write to be readable before Read persists it to state
		// (GetAccessKeyPolicy is eventually consistent after a write).
		if err := ramAccessKeyPolicyWaitConsistent(client, d.Id(), d.Get("access_key_policy").(string), d.Timeout(schema.TimeoutUpdate)); err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
	}

	return resourceAliCloudRamAccessKeyPolicyRead(d, meta)
}

func resourceAliCloudRamAccessKeyPolicyDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	action := "SetAccessKeyPolicy"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	userPrincipalName, userAccessKeyId := parseRamAccessKeyPolicyId(d.Id())
	request = make(map[string]interface{})
	request["UserAccessKeyId"] = userAccessKeyId
	// There is no dedicated DeleteAccessKeyPolicy API. The network access
	// restriction policy is reset by setting a disabled policy with no
	// statements via SetAccessKeyPolicy, which clears all whitelist rules.
	// Note: the API rejects an empty document ("{}"), so an explicit disabled
	// policy document must be sent instead.
	request["AccessKeyPolicy"] = ramAccessKeyPolicyResetDocument
	if userPrincipalName != "" {
		request["UserPrincipalName"] = userPrincipalName
	}

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = client.RpcPost("Ims", "2019-08-15", action, query, request, true)

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
		if IsExpectedErrors(err, []string{"EntityNotExist.AccessKeyPolicy"}) || NotFoundError(err) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	return nil
}
