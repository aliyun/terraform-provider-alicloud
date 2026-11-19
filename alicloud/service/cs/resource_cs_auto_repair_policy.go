package cs

import (
	"encoding/json"
	"fmt"
	"log"
	"reflect"
	"strings"
	"time"

	roacs "github.com/alibabacloud-go/cs-20151215/v8/client"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

const (
	ResourceAlicloudAutoRepairPolicy = "alicloud_cs_auto_repair_policy"

	autoRepairProcedureQuarantineNode   = "QuarantineNode"
	autoRepairProcedureEvict            = "Evict"
	autoRepairProcedureUnquarantineNode = "UnquarantineNode"
)

// ResourceAutoRepairPolicy returns the alicloud_cs_auto_repair_policy resource.
func ResourceAutoRepairPolicy() *schema.Resource {
	return &schema.Resource{
		Create: resourceAutoRepairPolicyCreate,
		Read:   resourceAutoRepairPolicyRead,
		Update: resourceAutoRepairPolicyUpdate,
		Delete: resourceAutoRepairPolicyDelete,
		Importer: &schema.ResourceImporter{
			State: resourceAutoRepairPolicyImport,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"cluster_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"resource_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Default:      "nodepool",
				ValidateFunc: validation.StringInSlice([]string{"nodepool"}, false),
			},
			"resource_sub_type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"ess", "lingjun"}, false),
			},
			"rules": {
				Type:     schema.TypeList,
				Required: true,
				MinItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"incidents": {
							Type:     schema.TypeList,
							Required: true,
							MinItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:     schema.TypeString,
										Required: true,
									},
								},
							},
						},
						"repair_procedure": {
							Type:     schema.TypeList,
							Required: true,
							MinItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:     schema.TypeString,
										Required: true,
									},
									"config": {
										Type:             schema.TypeString,
										Optional:         true,
										DiffSuppressFunc: autoRepairPolicyConfigDiffSuppress,
									},
									"intervention": {
										Type:     schema.TypeList,
										Optional: true,
										MaxItems: 1,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"enable": {
													Type:     schema.TypeBool,
													Optional: true,
													Default:  false,
												},
												"type": {
													Type:         schema.TypeString,
													Optional:     true,
													Default:      "label",
													ValidateFunc: validation.StringInSlice([]string{"label"}, false),
												},
												"inquiring_label": {
													Type:     schema.TypeList,
													Optional: true,
													MaxItems: 1,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"key": {
																Type:     schema.TypeString,
																Required: true,
															},
															"value": {
																Type:     schema.TypeString,
																Required: true,
															},
														},
													},
												},
												"approved_label": {
													Type:     schema.TypeList,
													Optional: true,
													MaxItems: 1,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"key": {
																Type:     schema.TypeString,
																Required: true,
															},
															"value": {
																Type:     schema.TypeString,
																Required: true,
															},
														},
													},
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
			"policy_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"resource_ids": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func resourceAutoRepairPolicyImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	parts, err := helper.ParseResourceId(d.Id(), 2)
	if err != nil {
		return nil, errs.WrapError(fmt.Errorf("the import ID must be '<cluster_id>:<policy_id>': %w", err))
	}
	if err = d.Set("cluster_id", parts[0]); err != nil {
		return nil, errs.WrapError(err)
	}
	d.SetId(fmt.Sprintf("%s%s%s", parts[0], COLON_SEPARATED, parts[1]))
	return []*schema.ResourceData{d}, nil
}

func resourceAutoRepairPolicyCreate(d *schema.ResourceData, meta interface{}) error {
	client, clientErr := meta.(*connectivity.AliyunClient).NewRoaCsClient()
	if clientErr != nil {
		return errs.WrapError(clientErr)
	}

	clusterId := d.Get("cluster_id").(string)
	action := "CreateAutoRepairPolicy"

	if err := validateAutoRepairPolicyRules(d.Get("rules")); err != nil {
		return errs.WrapError(err)
	}

	request := &roacs.CreateAutoRepairPolicyRequest{
		Name:            tea.String(d.Get("name").(string)),
		ResourceType:    tea.String(d.Get("resource_type").(string)),
		ResourceSubType: tea.String(d.Get("resource_sub_type").(string)),
		Rules:           buildAutoRepairPolicyCreateRules(d.Get("rules").([]interface{})),
	}

	response, err := client.CreateAutoRepairPolicy(tea.String(clusterId), request)
	logs.AddDebug(action, response, request)
	if err != nil {
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, ResourceAlicloudAutoRepairPolicy, action, errs.AlibabaCloudSdkGoERROR)
	}
	if response == nil || response.Body == nil || tea.StringValue(response.Body.PolicyId) == "" {
		return errs.WrapErrorf(fmt.Errorf("no policy_id returned"), errs.DefaultErrorMsg, ResourceAlicloudAutoRepairPolicy, action, errs.AlibabaCloudSdkGoERROR)
	}

	policyId := tea.StringValue(response.Body.PolicyId)
	d.SetId(fmt.Sprintf("%s%s%s", clusterId, COLON_SEPARATED, policyId))

	return resourceAutoRepairPolicyRead(d, meta)
}

func resourceAutoRepairPolicyRead(d *schema.ResourceData, meta interface{}) error {
	client, clientErr := meta.(*connectivity.AliyunClient).NewRoaCsClient()
	if clientErr != nil {
		return errs.WrapError(clientErr)
	}

	clusterId, policyId, err := parseAutoRepairPolicyId(d.Id())
	if err != nil {
		return errs.WrapError(err)
	}

	action := "DescribeAutoRepairPolicy"
	response, err := client.DescribeAutoRepairPolicy(tea.String(clusterId), tea.String(policyId))
	logs.AddDebug(action, response)
	if err != nil {
		if errs.IsExpectedErrors(err, []string{"ResourceNotFound.AutoRepairPolicy"}) || errs.NotFoundError(err) {
			log.Printf("[DEBUG] Resource %s DescribeAutoRepairPolicy not found: %s", ResourceAlicloudAutoRepairPolicy, err)
			d.SetId("")
			return nil
		}
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, ResourceAlicloudAutoRepairPolicy, action, errs.AlibabaCloudSdkGoERROR)
	}
	if response == nil || response.Body == nil {
		d.SetId("")
		return nil
	}

	body := response.Body
	d.Set("cluster_id", clusterId)
	d.Set("policy_id", policyId)
	d.Set("name", tea.StringValue(body.Name))
	d.Set("resource_type", tea.StringValue(body.ResourceType))
	d.Set("resource_sub_type", tea.StringValue(body.ResourceSubType))
	d.Set("type", tea.StringValue(body.Type))
	if err = d.Set("resource_ids", flattenAutoRepairPolicyResourceIds(body.ResourceIds)); err != nil {
		return errs.WrapError(err)
	}
	if err = d.Set("rules", flattenAutoRepairPolicyRules(body.Rules, d.Get("rules"))); err != nil {
		return errs.WrapError(err)
	}

	return nil
}

func resourceAutoRepairPolicyUpdate(d *schema.ResourceData, meta interface{}) error {
	client, clientErr := meta.(*connectivity.AliyunClient).NewRoaCsClient()
	if clientErr != nil {
		return errs.WrapError(clientErr)
	}

	clusterId, policyId, err := parseAutoRepairPolicyId(d.Id())
	if err != nil {
		return errs.WrapError(err)
	}

	request := &roacs.ModifyAutoRepairPolicyRequest{}
	update := false
	if d.HasChange("name") {
		update = true
	}
	if d.HasChange("rules") {
		update = true
	}
	if !update {
		return nil
	}
	if err := validateAutoRepairPolicyRules(d.Get("rules")); err != nil {
		return errs.WrapError(err)
	}
	// ModifyAutoRepairPolicy replaces name and rules as a whole, so both are
	// always sent even when only one of them changed
	request.Name = tea.String(d.Get("name").(string))
	request.Rules = buildAutoRepairPolicyModifyRules(d.Get("rules").([]interface{}))

	action := "ModifyAutoRepairPolicy"
	response, err := client.ModifyAutoRepairPolicy(tea.String(clusterId), tea.String(policyId), request)
	logs.AddDebug(action, response, request)
	if err != nil {
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, ResourceAlicloudAutoRepairPolicy, action, errs.AlibabaCloudSdkGoERROR)
	}

	return resourceAutoRepairPolicyRead(d, meta)
}

func resourceAutoRepairPolicyDelete(d *schema.ResourceData, meta interface{}) error {
	client, clientErr := meta.(*connectivity.AliyunClient).NewRoaCsClient()
	if clientErr != nil {
		return errs.WrapError(clientErr)
	}

	clusterId, policyId, err := parseAutoRepairPolicyId(d.Id())
	if err != nil {
		return errs.WrapError(err)
	}

	action := "DeleteAutoRepairPolicy"
	response, err := client.DeleteAutoRepairPolicy(tea.String(clusterId), tea.String(policyId))
	logs.AddDebug(action, response)
	if err != nil {
		if errs.IsExpectedErrors(err, []string{"ResourceNotFound.AutoRepairPolicy"}) || errs.NotFoundError(err) {
			return nil
		}
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, ResourceAlicloudAutoRepairPolicy, action, errs.AlibabaCloudSdkGoERROR)
	}

	return nil
}

func parseAutoRepairPolicyId(id string) (string, string, error) {
	parts := strings.Split(id, COLON_SEPARATED)
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return "", "", fmt.Errorf("invalid resource id %q, expected '<cluster_id>:<policy_id>'", id)
	}
	return parts[0], parts[1], nil
}

func buildAutoRepairPolicyCreateRules(rules []interface{}) []*roacs.CreateAutoRepairPolicyRequestRules {
	result := make([]*roacs.CreateAutoRepairPolicyRequestRules, 0, len(rules))
	for _, rule := range rules {
		ruleMap := rule.(map[string]interface{})
		result = append(result, &roacs.CreateAutoRepairPolicyRequestRules{
			Incidents:       buildAutoRepairPolicyCreateIncidents(ruleMap["incidents"].([]interface{})),
			RepairProcedure: buildAutoRepairPolicyCreateProcedures(ruleMap["repair_procedure"].([]interface{})),
		})
	}
	return result
}

func buildAutoRepairPolicyCreateIncidents(incidents []interface{}) []*roacs.CreateAutoRepairPolicyRequestRulesIncidents {
	result := make([]*roacs.CreateAutoRepairPolicyRequestRulesIncidents, 0, len(incidents))
	for _, incident := range incidents {
		incidentMap := incident.(map[string]interface{})
		result = append(result, &roacs.CreateAutoRepairPolicyRequestRulesIncidents{
			Name: tea.String(incidentMap["name"].(string)),
			// the public API only accepts system incidents, custom incidents are
			// internal and not exposed by this resource
			Type: tea.String("system"),
		})
	}
	return result
}

func buildAutoRepairPolicyCreateProcedures(procedures []interface{}) []*roacs.CreateAutoRepairPolicyRequestRulesRepairProcedure {
	result := make([]*roacs.CreateAutoRepairPolicyRequestRulesRepairProcedure, 0, len(procedures))
	for _, procedure := range procedures {
		procedureMap := procedure.(map[string]interface{})
		item := &roacs.CreateAutoRepairPolicyRequestRulesRepairProcedure{
			Name: tea.String(procedureMap["name"].(string)),
		}
		if config := buildAutoRepairPolicyConfig(procedureMap["config"]); config != nil {
			item.Config = config
		}
		if intervention := buildAutoRepairPolicyCreateIntervention(procedureMap["intervention"]); intervention != nil {
			item.Intervention = intervention
		}
		result = append(result, item)
	}
	return result
}

func buildAutoRepairPolicyCreateIntervention(raw interface{}) *roacs.CreateAutoRepairPolicyRequestRulesRepairProcedureIntervention {
	interventionMap := firstMapOfList(raw)
	if interventionMap == nil {
		return nil
	}
	intervention := &roacs.CreateAutoRepairPolicyRequestRulesRepairProcedureIntervention{
		Enable: tea.Bool(toBool(interventionMap["enable"])),
		Type:   tea.String(toStr(interventionMap["type"])),
	}
	if label := buildAutoRepairPolicyCreateLabel(interventionMap["inquiring_label"]); label != nil {
		intervention.InquiringLabel = &roacs.CreateAutoRepairPolicyRequestRulesRepairProcedureInterventionInquiringLabel{
			Key:   label.Key,
			Value: label.Value,
		}
	}
	if label := buildAutoRepairPolicyCreateLabel(interventionMap["approved_label"]); label != nil {
		intervention.ApprovedLabel = &roacs.CreateAutoRepairPolicyRequestRulesRepairProcedureInterventionApprovedLabel{
			Key:   label.Key,
			Value: label.Value,
		}
	}
	return intervention
}

type autoRepairPolicyLabel struct {
	Key   *string
	Value *string
}

func buildAutoRepairPolicyCreateLabel(raw interface{}) *autoRepairPolicyLabel {
	labelMap := firstMapOfList(raw)
	if labelMap == nil {
		return nil
	}
	return &autoRepairPolicyLabel{
		Key:   tea.String(toStr(labelMap["key"])),
		Value: tea.String(toStr(labelMap["value"])),
	}
}

func buildAutoRepairPolicyModifyRules(rules []interface{}) []*roacs.ModifyAutoRepairPolicyRequestRules {
	result := make([]*roacs.ModifyAutoRepairPolicyRequestRules, 0, len(rules))
	for _, rule := range rules {
		ruleMap := rule.(map[string]interface{})
		result = append(result, &roacs.ModifyAutoRepairPolicyRequestRules{
			Incidents:       buildAutoRepairPolicyModifyIncidents(ruleMap["incidents"].([]interface{})),
			RepairProcedure: buildAutoRepairPolicyModifyProcedures(ruleMap["repair_procedure"].([]interface{})),
		})
	}
	return result
}

func buildAutoRepairPolicyModifyIncidents(incidents []interface{}) []*roacs.ModifyAutoRepairPolicyRequestRulesIncidents {
	result := make([]*roacs.ModifyAutoRepairPolicyRequestRulesIncidents, 0, len(incidents))
	for _, incident := range incidents {
		incidentMap := incident.(map[string]interface{})
		result = append(result, &roacs.ModifyAutoRepairPolicyRequestRulesIncidents{
			Name: tea.String(incidentMap["name"].(string)),
			Type: tea.String("system"),
		})
	}
	return result
}

func buildAutoRepairPolicyModifyProcedures(procedures []interface{}) []*roacs.ModifyAutoRepairPolicyRequestRulesRepairProcedure {
	result := make([]*roacs.ModifyAutoRepairPolicyRequestRulesRepairProcedure, 0, len(procedures))
	for _, procedure := range procedures {
		procedureMap := procedure.(map[string]interface{})
		item := &roacs.ModifyAutoRepairPolicyRequestRulesRepairProcedure{
			Name: tea.String(procedureMap["name"].(string)),
		}
		if config := buildAutoRepairPolicyConfig(procedureMap["config"]); config != nil {
			item.Config = config
		}
		if intervention := buildAutoRepairPolicyModifyIntervention(procedureMap["intervention"]); intervention != nil {
			item.Intervention = intervention
		}
		result = append(result, item)
	}
	return result
}

func buildAutoRepairPolicyModifyIntervention(raw interface{}) *roacs.ModifyAutoRepairPolicyRequestRulesRepairProcedureIntervention {
	interventionMap := firstMapOfList(raw)
	if interventionMap == nil {
		return nil
	}
	intervention := &roacs.ModifyAutoRepairPolicyRequestRulesRepairProcedureIntervention{
		Enable: tea.Bool(toBool(interventionMap["enable"])),
		Type:   tea.String(toStr(interventionMap["type"])),
	}
	if label := buildAutoRepairPolicyCreateLabel(interventionMap["inquiring_label"]); label != nil {
		intervention.InquiringLabel = &roacs.ModifyAutoRepairPolicyRequestRulesRepairProcedureInterventionInquiringLabel{
			Key:   label.Key,
			Value: label.Value,
		}
	}
	if label := buildAutoRepairPolicyCreateLabel(interventionMap["approved_label"]); label != nil {
		intervention.ApprovedLabel = &roacs.ModifyAutoRepairPolicyRequestRulesRepairProcedureInterventionApprovedLabel{
			Key:   label.Key,
			Value: label.Value,
		}
	}
	return intervention
}

func validateAutoRepairPolicyRules(raw interface{}) error {
	rules, ok := raw.([]interface{})
	if !ok {
		return nil
	}
	for i, rule := range rules {
		ruleMap, ok := rule.(map[string]interface{})
		if !ok {
			continue
		}
		procedures, _ := ruleMap["repair_procedure"].([]interface{})
		definedAt := make(map[string]int, len(procedures))
		for j, procedure := range procedures {
			procedureMap, ok := procedure.(map[string]interface{})
			if !ok {
				continue
			}
			name := strings.TrimSpace(toStr(procedureMap["name"]))
			if name != "" {
				if first, duplicated := definedAt[name]; duplicated {
					return fmt.Errorf("rules.%d.repair_procedure.%d.name: repair procedure %s is duplicated, it is already defined at rules.%d.repair_procedure.%d", i, j, name, i, first)
				}
				definedAt[name] = j
			}
			if name == autoRepairProcedureEvict || name == autoRepairProcedureUnquarantineNode {
				if _, ok := definedAt[autoRepairProcedureQuarantineNode]; !ok {
					return fmt.Errorf("rules.%d.repair_procedure.%d.name: repair procedure %s requires %s to be defined before it", i, j, name, autoRepairProcedureQuarantineNode)
				}
			}
			config := strings.TrimSpace(toStr(procedureMap["config"]))
			if config == "" {
				continue
			}
			var value map[string]interface{}
			if err := json.Unmarshal([]byte(config), &value); err != nil {
				return fmt.Errorf("rules.%d.repair_procedure.%d.config must be a valid JSON object: %w", i, j, err)
			}
		}
	}
	return nil
}

func buildAutoRepairPolicyConfig(raw interface{}) map[string]interface{} {
	config := toStr(raw)
	if strings.TrimSpace(config) == "" {
		return nil
	}
	var result map[string]interface{}
	if err := json.Unmarshal([]byte(config), &result); err != nil {
		// the value is validated before the request is built
		return nil
	}
	if len(result) == 0 {
		return map[string]interface{}{}
	}
	return result
}

func flattenAutoRepairPolicyRules(rules []*roacs.DescribeAutoRepairPolicyResponseBodyRules, current interface{}) []map[string]interface{} {
	currentRules, _ := current.([]interface{})
	result := make([]map[string]interface{}, 0, len(rules))
	for i, rule := range rules {
		incidents := make([]map[string]interface{}, 0, len(rule.Incidents))
		for _, incident := range rule.Incidents {
			incidents = append(incidents, map[string]interface{}{
				"name": tea.StringValue(incident.Name),
			})
		}
		procedures := make([]map[string]interface{}, 0, len(rule.RepairProcedure))
		for _, procedure := range rule.RepairProcedure {
			name := tea.StringValue(procedure.Name)
			item := map[string]interface{}{
				"name":   name,
				"config": flattenAutoRepairPolicyConfig(procedure.Config),
			}
			// the API always returns an intervention object, even when manual
			// approval is not enabled; an enabled one is always written to state,
			// while a disabled one is written only when it is already there, so
			// that an explicitly configured "enable = false" survives the refresh
			if procedure.Intervention != nil && tea.BoolValue(procedure.Intervention.Enable) {
				item["intervention"] = flattenAutoRepairPolicyIntervention(procedure.Intervention)
			} else if intervention := interventionFromState(currentRules, i, name); len(intervention) > 0 {
				item["intervention"] = intervention
			}
			procedures = append(procedures, item)
		}
		result = append(result, map[string]interface{}{
			"incidents":        incidents,
			"repair_procedure": procedures,
		})
	}
	return result
}

// interventionFromState returns the intervention block held in state for the
// given rule and procedure name, which keeps an explicitly disabled block.
func interventionFromState(rules []interface{}, ruleIndex int, procedureName string) []interface{} {
	if ruleIndex < 0 || ruleIndex >= len(rules) {
		return nil
	}
	rule, ok := rules[ruleIndex].(map[string]interface{})
	if !ok {
		return nil
	}
	procedures, _ := rule["repair_procedure"].([]interface{})
	for _, procedure := range procedures {
		procedureMap, ok := procedure.(map[string]interface{})
		if !ok || toStr(procedureMap["name"]) != procedureName {
			continue
		}
		intervention, _ := procedureMap["intervention"].([]interface{})
		return intervention
	}
	return nil
}

func flattenAutoRepairPolicyIntervention(intervention *roacs.DescribeAutoRepairPolicyResponseBodyRulesRepairProcedureIntervention) []map[string]interface{} {
	item := map[string]interface{}{
		"enable": tea.BoolValue(intervention.Enable),
		"type":   tea.StringValue(intervention.Type),
	}
	if intervention.InquiringLabel != nil {
		item["inquiring_label"] = []map[string]interface{}{
			{
				"key":   tea.StringValue(intervention.InquiringLabel.Key),
				"value": tea.StringValue(intervention.InquiringLabel.Value),
			},
		}
	}
	if intervention.ApprovedLabel != nil {
		item["approved_label"] = []map[string]interface{}{
			{
				"key":   tea.StringValue(intervention.ApprovedLabel.Key),
				"value": tea.StringValue(intervention.ApprovedLabel.Value),
			},
		}
	}
	return []map[string]interface{}{item}
}

func flattenAutoRepairPolicyConfig(config map[string]interface{}) string {
	if len(config) == 0 {
		return "{}"
	}
	encoded, err := json.Marshal(config)
	if err != nil {
		log.Printf("[WARN] failed to marshal auto repair policy config: %s", err)
		return "{}"
	}
	return string(encoded)
}

func flattenAutoRepairPolicyResourceIds(ids []*string) []string {
	result := make([]string, 0, len(ids))
	for _, id := range ids {
		result = append(result, tea.StringValue(id))
	}
	return result
}

func autoRepairPolicyConfigDiffSuppress(_, oldValue, newValue string, _ *schema.ResourceData) bool {
	if oldValue == newValue {
		return true
	}
	oldConfig, err := normalizeAutoRepairPolicyConfig(oldValue)
	if err != nil {
		return false
	}
	newConfig, err := normalizeAutoRepairPolicyConfig(newValue)
	if err != nil {
		return false
	}
	return reflect.DeepEqual(oldConfig, newConfig)
}

func normalizeAutoRepairPolicyConfig(raw string) (interface{}, error) {
	raw = strings.TrimSpace(raw)
	if raw == "" || raw == "null" {
		return map[string]interface{}{}, nil
	}
	var value interface{}
	if err := json.Unmarshal([]byte(raw), &value); err != nil {
		return nil, err
	}
	if value == nil {
		return map[string]interface{}{}, nil
	}
	return value, nil
}

func firstMapOfList(raw interface{}) map[string]interface{} {
	list, ok := raw.([]interface{})
	if !ok || len(list) == 0 {
		return nil
	}
	result, _ := list[0].(map[string]interface{})
	return result
}

func toStr(raw interface{}) string {
	if raw == nil {
		return ""
	}
	if value, ok := raw.(string); ok {
		return value
	}
	return fmt.Sprintf("%v", raw)
}

func toBool(raw interface{}) bool {
	value, _ := raw.(bool)
	return value
}
