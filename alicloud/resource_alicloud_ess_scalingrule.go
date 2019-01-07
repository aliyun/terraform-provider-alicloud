package alicloud

import (
	"fmt"
	"strings"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ess"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func resourceAlicloudEssScalingRule() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliyunEssScalingRuleCreate,
		Read:   resourceAliyunEssScalingRuleRead,
		Update: resourceAliyunEssScalingRuleUpdate,
		Delete: resourceAliyunEssScalingRuleDelete,

		Schema: map[string]*schema.Schema{
			"scaling_group_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"adjustment_type": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validateAllowedStringValue([]string{string(QuantityChangeInCapacity),
					string(PercentChangeInCapacity), string(TotalCapacity)}),
			},
			"adjustment_value": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"scaling_rule_name": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
			"ari": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"cooldown": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validateIntegerInRange(0, 86400),
			},
		},
	}
}

func resourceAliyunEssScalingRuleCreate(d *schema.ResourceData, meta interface{}) error {

	args, err := buildAlicloudEssScalingRuleArgs(d, meta)
	if err != nil {
		return err
	}

	client := meta.(*connectivity.AliyunClient)

	raw, err := client.WithEssClient(func(essClient *ess.Client) (interface{}, error) {
		return essClient.CreateScalingRule(args)
	})
	if err != nil {
		return err
	}
	rule, _ := raw.(*ess.CreateScalingRuleResponse)
	d.SetId(d.Get("scaling_group_id").(string) + COLON_SEPARATED + rule.ScalingRuleId)

	return resourceAliyunEssScalingRuleUpdate(d, meta)
}

func resourceAliyunEssScalingRuleRead(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	essService := EssService{client}
	ids := strings.Split(d.Id(), COLON_SEPARATED)

	rule, err := essService.DescribeScalingRuleById(ids[0], ids[1])
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error Describe ESS scaling rule Attribute: %#v", err)
	}

	d.Set("scaling_group_id", rule.ScalingGroupId)
	d.Set("ari", rule.ScalingRuleAri)
	d.Set("adjustment_type", rule.AdjustmentType)
	d.Set("adjustment_value", rule.AdjustmentValue)
	d.Set("scaling_rule_name", rule.ScalingRuleName)
	d.Set("cooldown", rule.Cooldown)

	return nil
}

func resourceAliyunEssScalingRuleDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	essService := EssService{client}
	ids := strings.Split(d.Id(), COLON_SEPARATED)

	return resource.Retry(2*time.Minute, func() *resource.RetryError {
		err := essService.DeleteScalingRuleById(ids[1])

		if err != nil {
			if IsExceptedErrors(err, []string{InvalidScalingRuleIdNotFound}) {
				return nil
			}
			return resource.RetryableError(fmt.Errorf("Delete scaling rule timeout and got an error:%#v.", err))
		}

		_, err = essService.DescribeScalingRuleById(ids[0], ids[1])
		if err != nil {
			if NotFoundError(err) {
				return nil
			}
			return resource.NonRetryableError(err)
		}

		return resource.RetryableError(fmt.Errorf("Delete scaling rule timeout and got an error:%#v.", err))
	})
}

func resourceAliyunEssScalingRuleUpdate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	ids := strings.Split(d.Id(), COLON_SEPARATED)

	args := ess.CreateModifyScalingRuleRequest()
	args.ScalingRuleId = ids[1]

	if d.HasChange("adjustment_type") {
		args.AdjustmentType = d.Get("adjustment_type").(string)
	}

	if d.HasChange("adjustment_value") {
		args.AdjustmentValue = requests.NewInteger(d.Get("adjustment_value").(int))
	}

	if d.HasChange("scaling_rule_name") {
		args.ScalingRuleName = d.Get("scaling_rule_name").(string)
	}

	if d.HasChange("cooldown") {
		args.Cooldown = requests.NewInteger(d.Get("cooldown").(int))
	}

	_, err := client.WithEssClient(func(essClient *ess.Client) (interface{}, error) {
		return essClient.ModifyScalingRule(args)
	})
	if err != nil {
		return err
	}

	return resourceAliyunEssScalingRuleRead(d, meta)
}

func buildAlicloudEssScalingRuleArgs(d *schema.ResourceData, meta interface{}) (*ess.CreateScalingRuleRequest, error) {
	args := ess.CreateCreateScalingRuleRequest()
	args.ScalingGroupId = d.Get("scaling_group_id").(string)
	args.AdjustmentType = d.Get("adjustment_type").(string)
	args.AdjustmentValue = requests.NewInteger(d.Get("adjustment_value").(int))

	if v := d.Get("scaling_rule_name").(string); v != "" {
		args.ScalingRuleName = v
	}

	if v := d.Get("cooldown").(int); v != 0 {
		args.Cooldown = requests.NewInteger(v)
	}

	return args, nil
}
