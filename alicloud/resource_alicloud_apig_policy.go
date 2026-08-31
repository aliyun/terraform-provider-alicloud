package alicloud

import (
	"fmt"
	"log"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceAliCloudApigPolicy() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudApigPolicyCreate,
		Read:   resourceAliCloudApigPolicyRead,
		Update: resourceAliCloudApigPolicyUpdate,
		Delete: resourceAliCloudApigPolicyDelete,
		Importer: &schema.ResourceImporter{
			State: resourceAliCloudApigPolicyImport,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"policy_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"class_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"config": {
				Type:         schema.TypeString,
				Required:     true,
				Sensitive:    true,
				ValidateFunc: validation.StringIsJSON,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					equal, _ := compareJsonTemplateAreEquivalent(old, new)
					return equal
				},
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"attach_resource_ids": {
				Type:     schema.TypeSet,
				Required: true,
				MinItems: 1,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"attach_resource_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Default:      "GatewayRoute",
				ValidateFunc: StringInSlice([]string{"HttpApi", "GatewayRoute", "Operation", "GatewayService", "GatewayServicePort", "Gateway", "Domain"}, false),
			},
			"environment_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"gateway_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"policy_attachment_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceAliCloudApigPolicyImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	parts, err := apigParseCompositeID(d.Id(), 2)
	if err != nil {
		return nil, err
	}
	d.SetId(parts[0])
	if err := d.Set("policy_attachment_id", parts[1]); err != nil {
		return nil, err
	}
	return []*schema.ResourceData{d}, nil
}

func apigPolicyRequest(d *schema.ResourceData, includeClass bool) map[string]interface{} {
	body := map[string]interface{}{
		"name":               d.Get("policy_name"),
		"config":             d.Get("config"),
		"description":        d.Get("description"),
		"attachResourceIds":  apigSortedStringSet(d.Get("attach_resource_ids")),
		"attachResourceType": d.Get("attach_resource_type"),
		"environmentId":      d.Get("environment_id"),
		"gatewayId":          d.Get("gateway_id"),
	}
	if includeClass {
		body["className"] = d.Get("class_name")
	}
	return body
}

func resourceAliCloudApigPolicyCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	action := "/v1/policies"
	body := apigPolicyRequest(d, true)
	var response map[string]interface{}
	var err error
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.RoaPost("APIG", "2024-03-27", action, nil, nil, body, true)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_apig_policy", action, AlibabaCloudSdkGoERROR)
	}
	data, err := apigResponseData(response)
	if err != nil {
		return err
	}
	policyID, ok := data["policyId"].(string)
	if !ok || policyID == "" {
		return fmt.Errorf("APIG CreateAndAttachPolicy response does not contain policyId")
	}
	attachment, ok := data["attachment"].(map[string]interface{})
	if !ok {
		return fmt.Errorf("APIG CreateAndAttachPolicy response does not contain attachment metadata")
	}
	attachmentID, ok := attachment["policyAttachmentId"].(string)
	if !ok || attachmentID == "" {
		return fmt.Errorf("APIG CreateAndAttachPolicy response does not contain policyAttachmentId")
	}
	d.SetId(policyID)
	if err := d.Set("policy_attachment_id", attachmentID); err != nil {
		return err
	}
	if err := apigWaitForPolicyPresence(client, d, true, d.Timeout(schema.TimeoutCreate)); err != nil {
		return WrapErrorf(err, DefaultErrorMsg, policyID, action, AlibabaCloudSdkGoERROR)
	}
	return resourceAliCloudApigPolicyRead(d, meta)
}

func apigWaitForPolicyPresence(client *connectivity.AliyunClient, d *schema.ResourceData, present bool, timeout time.Duration) error {
	return resource.Retry(timeout, func() *resource.RetryError {
		policyAction := fmt.Sprintf("/v2/policies/%s", d.Id())
		_, err := client.RoaGet("APIG", "2024-03-27", policyAction, nil, nil, nil)
		if err != nil {
			if IsExpectedErrors(err, []string{"NotFound.PolicyNotFound"}) || NotFoundError(err) {
				if present {
					return resource.RetryableError(err)
				}
				return nil
			}
			if NeedRetry(err) {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		if !present {
			return resource.RetryableError(fmt.Errorf("APIG policy %s still exists", d.Id()))
		}

		attachmentID := d.Get("policy_attachment_id").(string)
		attachmentAction := fmt.Sprintf("/v1/policy-attachments/%s", attachmentID)
		_, err = client.RoaGet("APIG", "2024-03-27", attachmentAction, nil, nil, nil)
		if err != nil {
			if IsExpectedErrors(err, []string{"NotFound.PolicyAttachmentNotFound"}) || NotFoundError(err) || NeedRetry(err) {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
}

func resourceAliCloudApigPolicyRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	policyAction := fmt.Sprintf("/v2/policies/%s", d.Id())
	policyResponse, err := client.RoaGet("APIG", "2024-03-27", policyAction, nil, nil, nil)
	if err != nil {
		if !d.IsNewResource() && (IsExpectedErrors(err, []string{"NotFound.PolicyNotFound"}) || NotFoundError(err)) {
			log.Printf("[DEBUG] Resource alicloud_apig_policy %s was not found", d.Id())
			d.SetId("")
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), policyAction, AlibabaCloudSdkGoERROR)
	}
	policy, err := apigResponseData(policyResponse)
	if err != nil {
		return err
	}
	attachmentID := d.Get("policy_attachment_id").(string)
	if attachmentID == "" {
		return fmt.Errorf("APIG policy %s has no policy_attachment_id in Terraform state", d.Id())
	}
	attachmentAction := fmt.Sprintf("/v1/policy-attachments/%s", attachmentID)
	attachmentResponse, err := client.RoaGet("APIG", "2024-03-27", attachmentAction, nil, nil, nil)
	if err != nil {
		if !d.IsNewResource() && (IsExpectedErrors(err, []string{"NotFound.PolicyAttachmentNotFound"}) || NotFoundError(err)) {
			log.Printf("[DEBUG] APIG policy attachment %s was not found", attachmentID)
			d.SetId("")
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, attachmentID, attachmentAction, AlibabaCloudSdkGoERROR)
	}
	attachment, err := apigResponseData(attachmentResponse)
	if err != nil {
		return err
	}
	if attachment["policyId"] != d.Id() {
		return fmt.Errorf("APIG policy attachment %s belongs to a different policy", attachmentID)
	}
	attachment, err = apigFindPolicyAttachment(client, d.Id(), attachmentID, attachment)
	if err != nil {
		return err
	}
	if err := d.Set("policy_name", policy["name"]); err != nil {
		return err
	}
	if err := d.Set("class_name", policy["className"]); err != nil {
		return err
	}
	if err := d.Set("config", policy["config"]); err != nil {
		return err
	}
	if err := d.Set("description", policy["description"]); err != nil {
		return err
	}
	if err := d.Set("attach_resource_type", attachment["attachResourceType"]); err != nil {
		return err
	}
	if err := d.Set("environment_id", attachment["environmentId"]); err != nil {
		return err
	}
	if err := d.Set("gateway_id", attachment["gatewayId"]); err != nil {
		return err
	}
	resourceIDs := apigStringSlice(attachment["attachResourceIds"])
	if len(resourceIDs) == 0 {
		if resourceID, ok := attachment["attachResourceId"].(string); ok && resourceID != "" {
			resourceIDs = []string{resourceID}
		}
	}
	if len(resourceIDs) == 0 {
		return fmt.Errorf("APIG policy attachment %s contains no attached resource IDs", attachmentID)
	}
	if err := d.Set("attach_resource_ids", resourceIDs); err != nil {
		return err
	}
	return nil
}

func apigFindPolicyAttachment(client *connectivity.AliyunClient, policyID, attachmentID string, attachment map[string]interface{}) (map[string]interface{}, error) {
	query := map[string]*string{
		"withAttachments": StringPointer("true"),
	}
	for queryName, responseName := range map[string]string{
		"attachResourceId":   "attachResourceId",
		"attachResourceType": "attachResourceType",
		"environmentId":      "environmentId",
		"gatewayId":          "gatewayId",
	} {
		if value, ok := attachment[responseName].(string); ok && value != "" {
			query[queryName] = StringPointer(value)
		}
	}
	response, err := client.RoaGet("APIG", "2024-03-27", "/v1/policies", query, nil, nil)
	if err != nil {
		return nil, WrapErrorf(err, DefaultErrorMsg, policyID, "/v1/policies", AlibabaCloudSdkGoERROR)
	}
	data, err := apigResponseData(response)
	if err != nil {
		return nil, err
	}
	if candidateAttachment := apigSelectPolicyAttachment(apigObjectSlice(data["items"]), policyID, attachmentID); candidateAttachment != nil {
		return candidateAttachment, nil
	}
	return nil, fmt.Errorf("APIG policy attachment %s was not returned by ListPolicies with attachments", attachmentID)
}

func apigSelectPolicyAttachment(policies []map[string]interface{}, policyID, attachmentID string) map[string]interface{} {
	for _, candidatePolicy := range policies {
		if candidatePolicy["policyId"] != policyID {
			continue
		}
		for _, candidateAttachment := range apigObjectSlice(candidatePolicy["attachments"]) {
			if candidateAttachment["policyAttachmentId"] == attachmentID {
				return candidateAttachment
			}
		}
	}
	return nil
}

func resourceAliCloudApigPolicyUpdate(d *schema.ResourceData, meta interface{}) error {
	if !d.HasChanges("policy_name", "config", "description", "attach_resource_ids", "environment_id", "gateway_id") {
		return resourceAliCloudApigPolicyRead(d, meta)
	}
	client := meta.(*connectivity.AliyunClient)
	action := fmt.Sprintf("/v1/policies/%s", d.Id())
	body := apigPolicyRequest(d, false)
	var err error
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
		_, err = client.RoaPut("APIG", "2024-03-27", action, nil, nil, body, true)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}
	return resourceAliCloudApigPolicyRead(d, meta)
}

func resourceAliCloudApigPolicyDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	attachmentID := d.Get("policy_attachment_id").(string)
	if attachmentID != "" {
		action := fmt.Sprintf("/v1/policy-attachments/%s", attachmentID)
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err := resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
			_, err := client.RoaDelete("APIG", "2024-03-27", action, nil, nil, nil, true)
			if err != nil {
				if IsExpectedErrors(err, []string{"NotFound.PolicyAttachmentNotFound"}) || NotFoundError(err) {
					return nil
				}
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, attachmentID, action, AlibabaCloudSdkGoERROR)
		}
	}
	action := fmt.Sprintf("/v2/policies/%s", d.Id())
	var err error
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		_, err = client.RoaDelete("APIG", "2024-03-27", action, nil, nil, nil, true)
		if err != nil {
			if IsExpectedErrors(err, []string{"NotFound.PolicyNotFound"}) || NotFoundError(err) {
				return nil
			}
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}
	if err := apigWaitForPolicyPresence(client, d, false, d.Timeout(schema.TimeoutDelete)); err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}
	return nil
}
