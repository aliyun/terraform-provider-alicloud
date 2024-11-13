package alicloud

import (
	"fmt"
	"log"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAliCloudSddpRule() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudSddpRuleCreate,
		Read:   resourceAliCloudSddpRuleRead,
		Update: resourceAliCloudSddpRuleUpdate,
		Delete: resourceAliCloudSddpRuleDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"rule_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"category": {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: IntInSlice([]int{0, 2}),
			},
			"content": {
				Type:     schema.TypeString,
				Required: true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					equal, _ := compareArrayJsonTemplateAreEquivalent(old, new)
					return equal
				},
			},
			"content_category": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"risk_level_id": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: StringInSlice([]string{"2", "3", "4", "5"}, false),
			},
			"rule_type": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: IntInSlice([]int{1, 2, 3}),
			},
			"product_code": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: StringInSlice([]string{"ODPS", "OSS", "RDS"}, false),
			},
			"product_id": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: StringInSlice([]string{"1", "2", "5"}, false),
			},
			"warn_level": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: IntInSlice([]int{1, 2, 3}),
			},
			"stat_express": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"target": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"status": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: IntInSlice([]int{0, 1}),
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"lang": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "zh",
				ValidateFunc: StringInSlice([]string{"zh", "en"}, false),
			},
			"custom_type": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func resourceAliCloudSddpRuleCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "CreateRule"
	request := make(map[string]interface{})
	var err error

	request["Name"] = d.Get("rule_name")
	request["Category"] = d.Get("category").(int)
	request["Content"] = d.Get("content")

	if v, ok := d.GetOk("content_category"); ok {
		request["ContentCategory"] = v
	}

	if v, ok := d.GetOk("risk_level_id"); ok {
		request["RiskLevelId"] = v
	}

	if v, ok := d.GetOkExists("rule_type"); ok {
		request["RuleType"] = v
	}

	if v, ok := d.GetOk("product_code"); ok {
		request["ProductCode"] = v
	}

	if v, ok := d.GetOk("product_id"); ok {
		request["ProductId"] = v
	}

	if v, ok := d.GetOkExists("warn_level"); ok {
		request["WarnLevel"] = v
	}

	if v, ok := d.GetOk("stat_express"); ok {
		request["StatExpress"] = v
	}

	if v, ok := d.GetOk("target"); ok {
		request["Target"] = v
	}

	if v, ok := d.GetOkExists("status"); ok {
		request["Status"] = v
	}

	if v, ok := d.GetOk("description"); ok {
		request["Description"] = v
	}

	if v, ok := d.GetOk("lang"); ok {
		request["Lang"] = v
	}

	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutCreate)), func() *resource.RetryError {
		response, err = client.RpcPost("Sddp", "2019-01-03", action, nil, request, true)
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_sddp_rule", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(response["Id"]))

	return resourceAliCloudSddpRuleRead(d, meta)
}

func resourceAliCloudSddpRuleRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	sddpService := SddpService{client}

	object, err := sddpService.DescribeSddpRule(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_sddp_rule sddpService.DescribeSddpRule Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("rule_name", object["Name"])
	d.Set("category", convertSddpRuleCategoryResponse(formatInt(object["Category"])))
	d.Set("content", object["Content"])
	d.Set("content_category", object["ContentCategory"])
	d.Set("product_code", object["ProductCode"])
	d.Set("warn_level", formatInt(object["WarnLevel"]))
	d.Set("stat_express", object["StatExpress"])
	d.Set("target", object["Target"])
	d.Set("status", formatInt(object["Status"]))
	d.Set("description", object["Description"])
	d.Set("custom_type", formatInt(object["CustomType"]))

	if riskLevelId, ok := object["RiskLevelId"]; ok {
		d.Set("risk_level_id", fmt.Sprint(riskLevelId))
	}

	if productId, ok := object["ProductId"]; ok {
		d.Set("product_id", fmt.Sprint(productId))
	}

	return nil
}

func resourceAliCloudSddpRuleUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	d.Partial(true)

	update := false
	modifyRuleReq := map[string]interface{}{
		"Id": d.Id(),
	}

	if d.HasChange("rule_name") {
		update = true
	}
	modifyRuleReq["Name"] = d.Get("rule_name")

	if d.HasChange("category") {
		update = true
	}
	modifyRuleReq["Category"] = d.Get("category")

	if d.HasChange("content") {
		update = true
	}
	modifyRuleReq["Content"] = d.Get("content")

	if d.HasChange("risk_level_id") {
		update = true
	}
	if v, ok := d.GetOk("risk_level_id"); ok {
		modifyRuleReq["RiskLevelId"] = v
	}

	if d.HasChange("rule_type") {
		update = true

		if v, ok := d.GetOkExists("rule_type"); ok {
			modifyRuleReq["RuleType"] = v
		}
	}

	if d.HasChange("product_code") {
		update = true
	}
	if v, ok := d.GetOk("product_code"); ok {
		modifyRuleReq["ProductCode"] = v
	}

	if d.HasChange("product_id") {
		update = true
	}
	if v, ok := d.GetOk("product_id"); ok {
		modifyRuleReq["ProductId"] = v
	}

	if d.HasChange("warn_level") {
		update = true

		if v, ok := d.GetOkExists("warn_level"); ok {
			modifyRuleReq["WarnLevel"] = v
		}
	}

	if v, ok := d.GetOk("lang"); ok {
		modifyRuleReq["Lang"] = convertSddpRuleModifyRuleLangRequest(v)
	}

	if update {
		action := "ModifyRule"
		var err error
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *resource.RetryError {
			response, err = client.RpcPost("Sddp", "2019-01-03", action, nil, modifyRuleReq, false)
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})
		addDebug(action, response, modifyRuleReq)

		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}

		d.SetPartial("rule_name")
		d.SetPartial("category")
		d.SetPartial("content")
		d.SetPartial("risk_level_id")
		d.SetPartial("rule_type")
		d.SetPartial("product_code")
		d.SetPartial("product_id")
		d.SetPartial("warn_level")
	}

	update = false
	modifyRuleStatusReq := map[string]interface{}{
		"Id": d.Id(),
	}

	if d.HasChange("status") {
		update = true

		if v, ok := d.GetOkExists("status"); ok {
			modifyRuleStatusReq["Status"] = v
		}
	}

	if v, ok := d.GetOk("lang"); ok {
		modifyRuleStatusReq["Lang"] = v
	}

	if update {
		action := "ModifyRuleStatus"
		var err error
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *resource.RetryError {
			response, err = client.RpcPost("Sddp", "2019-01-03", action, nil, modifyRuleStatusReq, false)
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})
		addDebug(action, response, modifyRuleStatusReq)

		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}

		d.SetPartial("status")
	}

	d.Partial(false)

	return resourceAliCloudSddpRuleRead(d, meta)
}

func resourceAliCloudSddpRuleDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	action := "DeleteRule"
	var response map[string]interface{}
	var err error
	request := map[string]interface{}{
		"Id": d.Id(),
	}

	if v, ok := d.GetOk("lang"); ok {
		request["Lang"] = v
	}

	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutDelete)), func() *resource.RetryError {
		response, err = client.RpcPost("Sddp", "2019-01-03", action, nil, request, true)
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

	return nil
}

func convertSddpRuleCategoryResponse(source interface{}) interface{} {
	switch source {
	case 5:
		return 0
	}

	return source
}

func convertSddpRuleModifyRuleLangRequest(source interface{}) interface{} {
	switch source {
	case "zh":
		return "zh_cn"
	case "en":
		return "en_us"
	}

	return source
}
