package alicloud

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceAliCloudThreatDetectionCustomCheckStandardPolicy() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudThreatDetectionCustomCheckStandardPolicyCreate,
		Read:   resourceAliCloudThreatDetectionCustomCheckStandardPolicyRead,
		Update: resourceAliCloudThreatDetectionCustomCheckStandardPolicyUpdate,
		Delete: resourceAliCloudThreatDetectionCustomCheckStandardPolicyDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"check_type": {
				Computed: true,
				Type:     schema.TypeString,
			},
			"dependent_policy_id": {
				ForceNew: true,
				Optional: true,
				Type:     schema.TypeString,
			},
			"policy_id": {
				Computed: true,
				Type:     schema.TypeString,
			},
			"policy_show_name": {
				Required: true,
				Type:     schema.TypeString,
			},
			"policy_type": {
				ForceNew:     true,
				Required:     true,
				Type:         schema.TypeString,
				ValidateFunc: validation.StringInSlice([]string{"STANDARD", "REQUIREMENT", "SECTION"}, false),
			},
			"type": {
				Optional:     true,
				Type:         schema.TypeString,
				ValidateFunc: validation.StringInSlice([]string{"AISPM", "IDENTITY_PERMISSION", "RISK", "COMPLIANCE"}, false),
			},
		},
	}
}

func resourceAliCloudThreatDetectionCustomCheckStandardPolicyCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	request := make(map[string]interface{})
	var err error

	if v, ok := d.GetOk("policy_show_name"); ok {
		request["PolicyShowName"] = v
	}
	if v, ok := d.GetOk("policy_type"); ok {
		request["PolicyType"] = v
	}
	if v, ok := d.GetOk("dependent_policy_id"); ok && v.(string) != "" {
		request["DependentPolicyId"] = v
	}
	if v, ok := d.GetOk("type"); ok && v.(string) != "" {
		request["Type"] = v
	}

	var response map[string]interface{}
	action := "CreateCheckPolicy"
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutCreate)), func() *resource.RetryError {
		resp, err := client.RpcPost("Sas", "2018-12-03", action, nil, request, false)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		response = resp
		addDebug(action, response, request)
		return nil
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_threat_detection_custom_check_standard_policy", action, AlibabaCloudSdkGoERROR)
	}

	v, err := jsonpath.Get("$.CheckCustomPolicy.PolicyId", response)
	if err != nil || v == nil {
		return WrapErrorf(err, IdMsg, "alicloud_threat_detection_custom_check_standard_policy")
	}
	policyId := fmt.Sprint(v)
	d.SetId(fmt.Sprintf("%s:%s", policyId, d.Get("policy_type").(string)))

	return resourceAliCloudThreatDetectionCustomCheckStandardPolicyRead(d, meta)
}

func resourceAliCloudThreatDetectionCustomCheckStandardPolicyRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	sasService := SasService{client}
	object, err := sasService.DescribeThreatDetectionCustomCheckStandardPolicy(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_threat_detection_custom_check_standard_policy sasService.DescribeThreatDetectionCustomCheckStandardPolicy Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("policy_id", fmt.Sprint(object["PolicyId"]))
	d.Set("policy_show_name", object["PolicyShowName"])
	d.Set("policy_type", object["PolicyType"])
	d.Set("check_type", object["CheckType"])
	if object["DependentPolicyId"] != nil && fmt.Sprint(object["DependentPolicyId"]) != "0" {
		d.Set("dependent_policy_id", fmt.Sprint(object["DependentPolicyId"]))
	}
	if v, ok := object["Type"]; ok && v != nil {
		d.Set("type", fmt.Sprint(v))
	}

	return nil
}

func resourceAliCloudThreatDetectionCustomCheckStandardPolicyUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	parts := strings.Split(d.Id(), ":")
	if len(parts) != 2 {
		return WrapError(fmt.Errorf("invalid resource id %s", d.Id()))
	}
	policyId := parts[0]
	policyType := parts[1]

	request := map[string]interface{}{
		"PolicyId":   policyId,
		"PolicyType": policyType,
	}
	update := false

	if !d.IsNewResource() && d.HasChange("policy_show_name") {
		update = true
	}
	if v, ok := d.GetOk("policy_show_name"); ok {
		request["PolicyShowName"] = v
	}
	if !d.IsNewResource() && d.HasChange("type") {
		update = true
	}
	if v, ok := d.GetOk("type"); ok && v.(string) != "" {
		request["Type"] = v
	}

	if update {
		action := "UpdateCheckPolicy"
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err := resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *resource.RetryError {
			resp, err := client.RpcPost("Sas", "2018-12-03", action, nil, request, false)
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			addDebug(action, resp, request)
			return nil
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
	}

	return resourceAliCloudThreatDetectionCustomCheckStandardPolicyRead(d, meta)
}

func resourceAliCloudThreatDetectionCustomCheckStandardPolicyDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	parts := strings.Split(d.Id(), ":")
	if len(parts) != 2 {
		return WrapError(fmt.Errorf("invalid resource id %s", d.Id()))
	}
	policyId := parts[0]
	policyType := parts[1]

	request := map[string]interface{}{
		"PolicyIds.1": policyId,
		"PolicyType":  policyType,
	}

	action := "DeleteCheckPolicy"
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err := resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutDelete)), func() *resource.RetryError {
		resp, err := client.RpcPost("Sas", "2018-12-03", action, nil, request, false)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(action, resp, request)
		return nil
	})
	if err != nil {
		if NotFoundError(err) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}
	return nil
}
