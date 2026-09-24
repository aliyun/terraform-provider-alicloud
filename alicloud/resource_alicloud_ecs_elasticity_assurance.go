// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAliCloudEcsElasticityAssurance() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudEcsElasticityAssuranceCreate,
		Read:   resourceAliCloudEcsElasticityAssuranceRead,
		Update: resourceAliCloudEcsElasticityAssuranceUpdate,
		Delete: resourceAliCloudEcsElasticityAssuranceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		CustomizeDiff: func(ctx context.Context, d *schema.ResourceDiff, v interface{}) error {
			if !d.NewValueKnown("period_unit") {
				return nil
			}
			periodUnit, _ := d.Get("period_unit").(string)

			if !d.NewValueKnown("period") {
				return nil
			}
			if period, ok := d.GetOk("period"); ok {
				if err := validateEcsElasticityAssurancePeriod(periodUnit, period.(int)); err != nil {
					return err
				}
			}

			if !ecsElasticityAssuranceRecurrenceRulesWhollyKnown(d) {
				return nil
			}
			if err := validateEcsElasticityAssuranceRecurrenceRuleAttributes(d.GetRawConfig()); err != nil {
				return err
			}
			return validateEcsElasticityAssuranceRecurrenceRules(periodUnit, ecsElasticityAssuranceRecurrenceRules(d.Get("recurrence_rules")))
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"assurance_times": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"auto_renew": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"auto_renew_period": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"auto_renew_period_unit": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"elasticity_assurance_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"end_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"instance_amount": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"instance_charge_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"instance_type": {
				Type:     schema.TypeList,
				Required: true,
				ForceNew: true,
				MaxItems: 1,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"period": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: IntBetween(1, 365),
			},
			"period_unit": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: StringInSlice([]string{"Month", "Year", "Day"}, false),
			},
			"private_pool_options_match_criteria": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"private_pool_options_name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"recurrence_rules": {
				Type:     schema.TypeSet,
				Optional: true,
				ForceNew: true,
				MaxItems: 20,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"start_hour": {
							Type:         schema.TypeInt,
							Optional:     true,
							ForceNew:     true,
							ValidateFunc: IntBetween(0, 24),
						},
						"recurrence_type": {
							Type:         schema.TypeString,
							Optional:     true,
							ForceNew:     true,
							ValidateFunc: StringInSlice([]string{"Daily", "Weekly", "Monthly"}, false),
						},
						"end_hour": {
							Type:         schema.TypeInt,
							Optional:     true,
							ForceNew:     true,
							ValidateFunc: IntBetween(0, 24),
						},
						"recurrence_value": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
					},
				},
			},
			"region_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"resource_group_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"start_time": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return old != "" && new != "" && strings.HasPrefix(new, strings.Trim(old, "00Z"))
				},
			},
			"start_time_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"tags": tagsSchema(),
			"used_assurance_times": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"zone_ids": {
				Type:     schema.TypeList,
				Required: true,
				ForceNew: true,
				MaxItems: 1,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func resourceAliCloudEcsElasticityAssuranceCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "CreateElasticityAssurance"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	request["RegionId"] = client.RegionId
	request["ClientToken"] = buildClientToken(action)

	if v, ok := d.GetOk("instance_type"); ok {
		instanceTypeMapsArray := v.([]interface{})
		request["InstanceType"] = instanceTypeMapsArray
	}

	if v, ok := d.GetOkExists("auto_renew"); ok {
		request["AutoRenew"] = v
	}
	if v, ok := d.GetOk("period_unit"); ok {
		request["PeriodUnit"] = v
	}
	if v, ok := d.GetOk("resource_group_id"); ok {
		request["ResourceGroupId"] = v
	}
	if v, ok := d.GetOk("tags"); ok {
		tagsMap := ConvertTags(v.(map[string]interface{}))
		request = expandTagsToMap(request, tagsMap)
	}

	if v, ok := d.GetOk("private_pool_options_name"); ok {
		request["PrivatePoolOptions.Name"] = v
	}
	if v, ok := d.GetOk("assurance_times"); ok {
		request["AssuranceTimes"] = v
	}
	if v, ok := d.GetOk("description"); ok {
		request["Description"] = v
	}
	if v, ok := d.GetOk("start_time"); ok {
		request["StartTime"] = v
	}
	if v, ok := d.GetOk("private_pool_options_match_criteria"); ok {
		request["PrivatePoolOptions.MatchCriteria"] = v
	}
	if v, ok := d.GetOkExists("auto_renew_period"); ok {
		request["AutoRenewPeriod"] = v
	}
	if v, ok := d.GetOkExists("instance_amount"); ok {
		request["InstanceAmount"] = v
	}
	if v, ok := d.GetOkExists("period"); ok {
		request["Period"] = v
	}
	periodUnit := fmt.Sprint(d.Get("period_unit"))
	if period, ok := d.GetOk("period"); ok {
		if err := validateEcsElasticityAssurancePeriod(periodUnit, period.(int)); err != nil {
			return WrapError(err)
		}
	}
	recurrenceRules := ecsElasticityAssuranceRecurrenceRules(d.Get("recurrence_rules"))
	if err := validateEcsElasticityAssuranceRecurrenceRuleAttributes(d.GetRawConfig()); err != nil {
		return WrapError(err)
	}
	if err := validateEcsElasticityAssuranceRecurrenceRules(periodUnit, recurrenceRules); err != nil {
		return WrapError(err)
	}
	if v, ok := d.GetOk("recurrence_rules"); ok {
		recurrenceRulesMapsArray, expandErr := expandEcsElasticityAssuranceRecurrenceRules(ecsElasticityAssuranceRecurrenceRules(v))
		if expandErr != nil {
			return WrapError(expandErr)
		}
		request["RecurrenceRules"] = recurrenceRulesMapsArray
	}
	if v, ok := d.GetOk("zone_ids"); ok {
		zoneIdMapsArray := v.([]interface{})
		request["ZoneId"] = zoneIdMapsArray
	}

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
		response, err = client.RpcPost("Ecs", "2014-05-26", action, query, request, true)
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_ecs_elasticity_assurance", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(response["PrivatePoolOptionsId"]))

	ecsServiceV2 := EcsServiceV2{client}
	stateConf := BuildStateConf([]string{}, []string{"Active", "Prepared"}, d.Timeout(schema.TimeoutCreate), 20*time.Second, ecsServiceV2.EcsElasticityAssuranceStateRefreshFunc(d.Id(), "Status", []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return resourceAliCloudEcsElasticityAssuranceUpdate(d, meta)
}

func expandEcsElasticityAssuranceRecurrenceRules(recurrenceRules []interface{}) ([]interface{}, error) {
	recurrenceRulesMapsArray := make([]interface{}, 0, len(recurrenceRules))
	for i, recurrenceRuleRaw := range recurrenceRules {
		recurrenceRule, ok := recurrenceRuleRaw.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("invalid recurrence_rules.%d: expected an object, got %T", i, recurrenceRuleRaw)
		}
		recurrenceRulesMapsArray = append(recurrenceRulesMapsArray, map[string]interface{}{
			"StartHour":       recurrenceRule["start_hour"],
			"EndHour":         recurrenceRule["end_hour"],
			"RecurrenceType":  recurrenceRule["recurrence_type"],
			"RecurrenceValue": recurrenceRule["recurrence_value"],
		})
	}
	return recurrenceRulesMapsArray, nil
}

func ecsElasticityAssuranceRecurrenceRules(v interface{}) []interface{} {
	switch rules := v.(type) {
	case *schema.Set:
		return rules.List()
	case []interface{}:
		return rules
	}
	return nil
}

// ecsElasticityAssuranceRecurrenceRulesWhollyKnown reports whether every recurrence rule attribute
// already has its final value. The block is a set, so a single unknown attribute makes the whole set
// unknown at plan time; validation is then deferred to Create, which re-runs the same checks.
func ecsElasticityAssuranceRecurrenceRulesWhollyKnown(d *schema.ResourceDiff) bool {
	if !d.NewValueKnown("recurrence_rules") {
		return false
	}
	rawConfig := d.GetRawConfig()
	if rawConfig.IsNull() || !rawConfig.IsKnown() || !rawConfig.Type().IsObjectType() || !rawConfig.Type().HasAttribute("recurrence_rules") {
		return false
	}
	return rawConfig.GetAttr("recurrence_rules").IsWhollyKnown()
}

// The attributes of a recurrence rule are declared Optional so that adding the block is not treated as a
// breaking change by scripts/compatibility (new required attributes are rejected). They remain mandatory for
// every rule, which is checked against the raw config because a missing attribute is null there, while the
// zero value of the decoded rule is indistinguishable from an omitted one.
func validateEcsElasticityAssuranceRecurrenceRuleAttributes(rawConfig cty.Value) error {
	if rawConfig.IsNull() || !rawConfig.IsKnown() || !rawConfig.Type().IsObjectType() || !rawConfig.Type().HasAttribute("recurrence_rules") {
		return nil
	}
	recurrenceRules := rawConfig.GetAttr("recurrence_rules")
	if recurrenceRules.IsNull() || !recurrenceRules.IsKnown() || !recurrenceRules.CanIterateElements() {
		return nil
	}
	index := 0
	for iterator := recurrenceRules.ElementIterator(); iterator.Next(); {
		_, recurrenceRule := iterator.Element()
		if recurrenceRule.IsKnown() && !recurrenceRule.IsNull() && recurrenceRule.Type().IsObjectType() {
			for _, attribute := range []string{"start_hour", "end_hour", "recurrence_type", "recurrence_value"} {
				if recurrenceRule.GetAttr(attribute).IsNull() {
					return fmt.Errorf("invalid recurrence_rules.%d: %s must be specified", index, attribute)
				}
			}
		}
		index++
	}
	return nil
}

func validateEcsElasticityAssurancePeriod(periodUnit string, period int) error {
	switch strings.ToUpper(periodUnit) {
	case "DAY":
		if period < 1 || period > 365 {
			return fmt.Errorf("invalid period %d: when period_unit is Day, the value must be between 1 and 365", period)
		}
	case "MONTH":
		if period < 1 || period > 9 {
			return fmt.Errorf("invalid period %d: when period_unit is Month, the value must be between 1 and 9", period)
		}
	default:
		if period < 1 || period > 5 {
			return fmt.Errorf("invalid period %d: when period_unit is Year, the value must be between 1 and 5", period)
		}
	}
	return nil
}

func validateEcsElasticityAssuranceRecurrenceValue(recurrenceType, recurrenceValue string) error {
	switch recurrenceType {
	case "Daily":
		value, err := strconv.Atoi(recurrenceValue)
		if err != nil || value < 1 || value > 31 {
			return fmt.Errorf("recurrence_value %q must be a single integer between 1 and 31 when recurrence_type is Daily", recurrenceValue)
		}
	case "Weekly":
		for _, part := range strings.Split(recurrenceValue, ",") {
			value, err := strconv.Atoi(strings.TrimSpace(part))
			if err != nil || value < 0 || value > 6 {
				return fmt.Errorf("recurrence_value %q must be comma-separated integers between 0 and 6 when recurrence_type is Weekly", recurrenceValue)
			}
		}
	case "Monthly":
		parts := strings.Split(recurrenceValue, "-")
		if len(parts) != 2 {
			return fmt.Errorf("recurrence_value %q must be in the form A-B when recurrence_type is Monthly", recurrenceValue)
		}
		start, startErr := strconv.Atoi(strings.TrimSpace(parts[0]))
		end, endErr := strconv.Atoi(strings.TrimSpace(parts[1]))
		if startErr != nil || endErr != nil || start < 1 || end > 31 || start > end {
			return fmt.Errorf("recurrence_value %q must be in the form A-B with 1 <= A <= B <= 31 when recurrence_type is Monthly", recurrenceValue)
		}
	default:
		return fmt.Errorf("recurrence_type must be one of Daily, Weekly and Monthly")
	}
	return nil
}

func validateEcsElasticityAssuranceRecurrenceRules(periodUnit string, recurrenceRules []interface{}) error {
	if len(recurrenceRules) == 0 {
		if periodUnit == "Day" {
			return fmt.Errorf("recurrence_rules must be specified when period_unit is Day")
		}
		return nil
	}
	if periodUnit != "Day" {
		return fmt.Errorf("recurrence_rules can only be specified when period_unit is Day")
	}
	for i, recurrenceRuleRaw := range recurrenceRules {
		recurrenceRule, ok := recurrenceRuleRaw.(map[string]interface{})
		if !ok {
			return fmt.Errorf("invalid recurrence_rules.%d: expected an object, got %T", i, recurrenceRuleRaw)
		}
		startHour, startHourOk := recurrenceRule["start_hour"].(int)
		endHour, endHourOk := recurrenceRule["end_hour"].(int)
		if !startHourOk || !endHourOk {
			return fmt.Errorf("invalid recurrence_rules.%d: start_hour and end_hour must be specified", i)
		}
		if endHour-startHour < 4 {
			return fmt.Errorf("invalid recurrence_rules.%d: the difference between end_hour and start_hour must be at least 4 hours", i)
		}
		recurrenceType, _ := recurrenceRule["recurrence_type"].(string)
		recurrenceValue, _ := recurrenceRule["recurrence_value"].(string)
		if err := validateEcsElasticityAssuranceRecurrenceValue(recurrenceType, recurrenceValue); err != nil {
			return fmt.Errorf("invalid recurrence_rules.%d: %s", i, err)
		}
	}
	return nil
}

func resourceAliCloudEcsElasticityAssuranceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	ecsServiceV2 := EcsServiceV2{client}

	objectRaw, err := ecsServiceV2.DescribeEcsElasticityAssurance(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_ecs_elasticity_assurance DescribeEcsElasticityAssurance Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("assurance_times", objectRaw["TotalAssuranceTimes"])
	d.Set("description", objectRaw["Description"])
	d.Set("end_time", objectRaw["EndTime"])
	d.Set("instance_charge_type", objectRaw["InstanceChargeType"])
	d.Set("private_pool_options_match_criteria", objectRaw["PrivatePoolOptionsMatchCriteria"])
	d.Set("private_pool_options_name", objectRaw["PrivatePoolOptionsName"])
	d.Set("region_id", objectRaw["RegionId"])
	d.Set("resource_group_id", objectRaw["ResourceGroupId"])
	d.Set("start_time_type", objectRaw["StartTimeType"])
	d.Set("status", objectRaw["Status"])
	d.Set("used_assurance_times", objectRaw["UsedAssuranceTimes"])
	d.Set("elasticity_assurance_id", objectRaw["PrivatePoolOptionsId"])
	d.Set("start_time", objectRaw["StartTime"])

	if v, ok := objectRaw["Tags"].(map[string]interface{}); ok {
		d.Set("tags", tagsToMap(v["Tag"]))
	}

	if v, ok := objectRaw["AllocatedResources"]; ok {
		allocatedResources := v.(map[string]interface{})
		if v, ok := allocatedResources["AllocatedResource"]; ok && len(v.([]interface{})) > 0 {
			allocatedResourceMap := v.([]interface{})[0].(map[string]interface{})
			d.Set("instance_type", []string{fmt.Sprint(allocatedResourceMap["InstanceType"])})
			d.Set("instance_amount", allocatedResourceMap["TotalAmount"])
			d.Set("zone_ids", []string{fmt.Sprint(allocatedResourceMap["zoneId"])})
		}
	}

	recurrenceRulesMaps := make([]map[string]interface{}, 0)
	if v, ok := objectRaw["RecurrenceRules"].(map[string]interface{}); ok {
		for _, recurrenceRuleChildRaw := range convertToInterfaceArray(v["RecurrenceRule"]) {
			recurrenceRule, ok := recurrenceRuleChildRaw.(map[string]interface{})
			if !ok {
				continue
			}
			recurrenceRulesMap := make(map[string]interface{})
			recurrenceRulesMap["start_hour"] = recurrenceRule["StartHour"]
			recurrenceRulesMap["end_hour"] = recurrenceRule["EndHour"]
			recurrenceRulesMap["recurrence_type"] = convertEcsElasticityAssuranceRecurrenceTypeResponse(recurrenceRule["RecurrenceType"])
			recurrenceRulesMap["recurrence_value"] = recurrenceRule["RecurrenceValue"]
			recurrenceRulesMaps = append(recurrenceRulesMaps, recurrenceRulesMap)
		}
	}
	if err := d.Set("recurrence_rules", recurrenceRulesMaps); err != nil {
		return WrapError(err)
	}

	objectRaw, err = ecsServiceV2.DescribeElasticityAssuranceDescribeElasticityAssuranceAutoRenewAttribute(d.Id())
	if err != nil && !NotFoundError(err) {
		return WrapError(err)
	}

	auto_renewVal, err := formatBool(convertEcsElasticityAssuranceElasticityAssuranceRenewAttributesElasticityAssuranceRenewAttributeRenewalStatusResponse(objectRaw["RenewalStatus"]))
	if err != nil {
		return err
	}
	d.Set("auto_renew", auto_renewVal)
	d.Set("auto_renew_period", objectRaw["Period"])
	d.Set("auto_renew_period_unit", objectRaw["PeriodUnit"])
	d.Set("elasticity_assurance_id", objectRaw["PrivatePoolOptionsId"])

	return nil
}

func resourceAliCloudEcsElasticityAssuranceUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	update := false
	d.Partial(true)

	var err error
	action := "ModifyElasticityAssurance"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["PrivatePoolOptions.Id"] = d.Id()
	request["RegionId"] = client.RegionId
	request["ClientToken"] = buildClientToken(action)
	if !d.IsNewResource() && d.HasChange("private_pool_options_name") {
		update = true
		request["PrivatePoolOptions.Name"] = d.Get("private_pool_options_name")
	}

	if !d.IsNewResource() && d.HasChange("description") {
		update = true
		request["Description"] = d.Get("description")
	}

	if !d.IsNewResource() && d.HasChange("instance_amount") {
		update = true
		request["InstanceAmount"] = d.Get("instance_amount")
	}

	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
			response, err = client.RpcPost("Ecs", "2014-05-26", action, query, request, true)
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
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
		ecsServiceV2 := EcsServiceV2{client}

		if d.HasChange("instance_amount") {
			stateConf := BuildStateConf([]string{}, []string{fmt.Sprint(d.Get("instance_amount"))}, d.Timeout(schema.TimeoutUpdate), 10*time.Second, ecsServiceV2.EcsElasticityAssuranceStateRefreshFunc(d.Id(), "InstanceAmount", []string{}))
			if _, err := stateConf.WaitForState(); err != nil {
				return WrapErrorf(err, IdMsg, d.Id())
			}
		}
	}
	update = false
	action = "ModifyElasticityAssuranceAutoRenewAttribute"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["PrivatePoolOptions.Id.1"] = d.Id()
	request["RegionId"] = client.RegionId

	if !d.IsNewResource() && d.HasChange("auto_renew") {
		update = true
	}
	isAutoRenew := d.Get("auto_renew").(bool)
	request["RenewalStatus"] = convertEcsElasticityAssuranceRenewalStatusRequest(isAutoRenew)

	if isAutoRenew && d.HasChange("auto_renew_period_unit") {
		update = true
	}
	request["PeriodUnit"] = d.Get("auto_renew_period_unit")

	if !d.IsNewResource() && d.HasChange("auto_renew_period") {
		update = true
	}
	request["Period"] = d.Get("auto_renew_period")

	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
			response, err = client.RpcPost("Ecs", "2014-05-26", action, query, request, true)
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
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
	}

	if d.HasChange("tags") {
		ecsServiceV2 := EcsServiceV2{client}
		if err := ecsServiceV2.SetResourceTags(d, "ELASTICITYASSURANCE"); err != nil {
			return WrapError(err)
		}
	}
	d.Partial(false)
	return resourceAliCloudEcsElasticityAssuranceRead(d, meta)
}

func resourceAliCloudEcsElasticityAssuranceDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	action := "RefundInstance"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	request["InstanceId"] = d.Id()

	request["ClientToken"] = buildClientToken(action)

	request["ProductCode"] = "ecsrep"
	request["ImmediatelyRelease"] = "1"
	var endpoint string
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutDelete), func() *retry.RetryError {
		response, err = client.RpcPostWithEndpoint("BssOpenApi", "2017-12-14", action, query, request, true, endpoint)
		request["ClientToken"] = buildClientToken(action)

		if err != nil {
			if IsExpectedErrors(err, []string{"ServiceUnavailable"}) || NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			if !client.IsInternationalAccount() && IsExpectedErrors(err, []string{""}) {
				endpoint = connectivity.BssOpenAPIEndpointInternational
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
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

	ecsServiceV2 := EcsServiceV2{client}
	stateConf := BuildStateConf([]string{}, []string{}, d.Timeout(schema.TimeoutDelete), 5*time.Second, ecsServiceV2.EcsElasticityAssuranceStateRefreshFunc(d.Id(), "Status", []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return nil
}

func convertEcsElasticityAssuranceElasticityAssuranceRenewAttributesElasticityAssuranceRenewAttributeRenewalStatusResponse(source interface{}) interface{} {
	source = fmt.Sprint(source)
	switch source {
	case "AutoRenewal":
		return "true"
	case "Normal":
		return "false"
	case "NotRenewal":
		return "false"
	}
	return source
}
func convertEcsElasticityAssuranceRenewalStatusRequest(source interface{}) interface{} {
	source = fmt.Sprint(source)
	switch source {
	case "true":
		return "AutoRenewal"
	case "false":
		return "Normal"
	}
	return source
}

func convertEcsElasticityAssuranceRecurrenceTypeResponse(source interface{}) interface{} {
	recurrenceType := fmt.Sprint(source)
	switch strings.ToUpper(recurrenceType) {
	case "DAILY":
		return "Daily"
	case "WEEKLY":
		return "Weekly"
	case "MONTHLY":
		return "Monthly"
	}
	return recurrenceType
}
