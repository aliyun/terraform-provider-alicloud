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
	ResourceAlicloudCSAutoRepairPolicy = "alicloud_cs_auto_repair_policy"

	autoRepairProcedureQuarantineNode   = "QuarantineNode"
	autoRepairProcedureEvict            = "Evict"
	autoRepairProcedureUnquarantineNode = "UnquarantineNode"
)

// ResourceCSAutoRepairPolicy returns the alicloud_cs_auto_repair_policy resource.
func ResourceCSAutoRepairPolicy() *schema.Resource {
	return &schema.Resource{
		Create: resourceCSAutoRepairPolicyCreate,
		Read:   resourceCSAutoRepairPolicyRead,
		Update: resourceCSAutoRepairPolicyUpdate,
		Delete: resourceCSAutoRepairPolicyDelete,
		Importer: &schema.ResourceImporter{
			State: resourceCSAutoRepairPolicyImport,
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

func resourceCSAutoRepairPolicyImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
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

func resourceCSAutoRepairPolicyCreate(d *schema.ResourceData, meta interface{}) error {
	client, clientErr := meta.(*connectivity.AliyunClient).NewRoaCsClient()
	if clientErr != nil {
		return errs.WrapError(clientErr)
	}

	clusterId := d.Get("cluster_id").(string)
	action := "CreateAutoRepairPolicy"

	if err := validateCSAutoRepairPolicyRules(d.Get("rules")); err != nil {
		return errs.WrapError(err)
	}

	request := &roacs.CreateAutoRepairPolicyRequest{
		Name:            tea.String(d.Get("name").(string)),
		ResourceType:    tea.String(d.Get("resource_type").(string)),
		ResourceSubType: tea.String(d.Get("resource_sub_type").(string)),
		Rules:           buildCSAutoRepairPolicyCreateRules(d.Get("rules").([]interface{})),
	}

	response, err := client.CreateAutoRepairPolicy(tea.String(clusterId), request)
	logs.AddDebug(action, response, request)
	if err != nil {
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, ResourceAlicloudCSAutoRepairPolicy, action, errs.AlibabaCloudSdkGoERROR)
	}
	if response == nil || response.Body == nil || tea.StringValue(response.Body.PolicyId) == "" {
		return errs.WrapErrorf(fmt.Errorf("no policy_id returned"), errs.DefaultErrorMsg, ResourceAlicloudCSAutoRepairPolicy, action, errs.AlibabaCloudSdkGoERROR)
	}

	policyId := tea.StringValue(response.Body.PolicyId)
	d.SetId(fmt.Sprintf("%s%s%s", clusterId, COLON_SEPARATED, policyId))

	return resourceCSAutoRepairPolicyRead(d, meta)
}

func resourceCSAutoRepairPolicyRead(d *schema.ResourceData, meta interface{}) error {
	client, clientErr := meta.(*connectivity.AliyunClient).NewRoaCsClient()
	if clientErr != nil {
		return errs.WrapError(clientErr)
	}

	clusterId, policyId, err := parseCSAutoRepairPolicyId(d.Id())
	if err != nil {
		return errs.WrapError(err)
	}

	action := "DescribeAutoRepairPolicy"
	response, err := client.DescribeAutoRepairPolicy(tea.String(clusterId), tea.String(policyId))
	logs.AddDebug(action, response)
	if err != nil {
		if errs.IsExpectedErrors(err, []string{"ResourceNotFound.AutoRepairPolicy"}) || errs.NotFoundError(err) {
			log.Printf("[DEBUG] Resource %s DescribeAutoRepairPolicy not found: %s", ResourceAlicloudCSAutoRepairPolicy, err)
			d.SetId("")
			return nil
		}
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, ResourceAlicloudCSAutoRepairPolicy, action, errs.AlibabaCloudSdkGoERROR)
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
	if err = d.Set("resource_ids", flattenCSAutoRepairPolicyResourceIds(body.ResourceIds)); err != nil {
		return errs.WrapError(err)
	}
	if err = d.Set("rules", flattenCSAutoRepairPolicyRules(body.Rules)); err != nil {
		return errs.WrapError(err)
	}

	return nil
}

func resourceCSAutoRepairPolicyUpdate(d *schema.ResourceData, meta interface{}) error {
	client, clientErr := meta.(*connectivity.AliyunClient).NewRoaCsClient()
	if clientErr != nil {
		return errs.WrapError(clientErr)
	}

	clusterId, policyId, err := parseCSAutoRepairPolicyId(d.Id())
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
	if err := validateCSAutoRepairPolicyRules(d.Get("rules")); err != nil {
		return errs.WrapError(err)
	}
	// ModifyAutoRepairPolicy replaces name and rules as a whole, so both are
	// always sent even when only one of them changed
	request.Name = tea.String(d.Get("name").(string))
	request.Rules = buildCSAutoRepairPolicyModifyRules(d.Get("rules").([]interface{}))

	action := "ModifyAutoRepairPolicy"
	response, err := client.ModifyAutoRepairPolicy(tea.String(clusterId), tea.String(policyId), request)
	logs.AddDebug(action, response, request)
	if err != nil {
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, ResourceAlicloudCSAutoRepairPolicy, action, errs.AlibabaCloudSdkGoERROR)
	}

	return resourceCSAutoRepairPolicyRead(d, meta)
}

func resourceCSAutoRepairPolicyDelete(d *schema.ResourceData, meta interface{}) error {
	client, clientErr := meta.(*connectivity.AliyunClient).NewRoaCsClient()
	if clientErr != nil {
		return errs.WrapError(clientErr)
	}

	clusterId, policyId, err := parseCSAutoRepairPolicyId(d.Id())
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
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, ResourceAlicloudCSAutoRepairPolicy, action, errs.AlibabaCloudSdkGoERROR)
	}

	return nil
}

func parseCSAutoRepairPolicyId(id string) (string, string, error) {
	parts := strings.Split(id, COLON_SEPARATED)
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return "", "", fmt.Errorf("invalid resource id %q, expected '<cluster_id>:<policy_id>'", id)
	}
	return parts[0], parts[1], nil
}

func buildCSAutoRepairPolicyCreateRules(rules []interface{}) []*roacs.CreateAutoRepairPolicyRequestRules {
	result := make([]*roacs.CreateAutoRepairPolicyRequestRules, 0, len(rules))
	for _, rule := range rules {
		ruleMap := rule.(map[string]interface{})
		result = append(result, &roacs.CreateAutoRepairPolicyRequestRules{
			Incidents:       buildCSAutoRepairPolicyCreateIncidents(ruleMap["incidents"].([]interface{})),
			RepairProcedure: buildCSAutoRepairPolicyCreateProcedures(ruleMap["repair_procedure"].([]interface{})),
		})
	}
	return result
}

func buildCSAutoRepairPolicyCreateIncidents(incidents []interface{}) []*roacs.CreateAutoRepairPolicyRequestRulesIncidents {
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

func buildCSAutoRepairPolicyCreateProcedures(procedures []interface{}) []*roacs.CreateAutoRepairPolicyRequestRulesRepairProcedure {
	result := make([]*roacs.CreateAutoRepairPolicyRequestRulesRepairProcedure, 0, len(procedures))
	for _, procedure := range procedures {
		procedureMap := procedure.(map[string]interface{})
		item := &roacs.CreateAutoRepairPolicyRequestRulesRepairProcedure{
			Name: tea.String(procedureMap["name"].(string)),
		}
		if config := buildCSAutoRepairPolicyConfig(procedureMap["config"]); config != nil {
			item.Config = config
		}
		if intervention := buildCSAutoRepairPolicyCreateIntervention(procedureMap["intervention"]); intervention != nil {
			item.Intervention = intervention
		}
		result = append(result, item)
	}
	return result
}

func buildCSAutoRepairPolicyCreateIntervention(raw interface{}) *roacs.CreateAutoRepairPolicyRequestRulesRepairProcedureIntervention {
	interventionMap := firstMapOfList(raw)
	if interventionMap == nil {
		return nil
	}
	intervention := &roacs.CreateAutoRepairPolicyRequestRulesRepairProcedureIntervention{
		Enable: tea.Bool(toBool(interventionMap["enable"])),
		Type:   tea.String(toStr(interventionMap["type"])),
	}
	if label := buildCSAutoRepairPolicyCreateLabel(interventionMap["inquiring_label"]); label != nil {
		intervention.InquiringLabel = &roacs.CreateAutoRepairPolicyRequestRulesRepairProcedureInterventionInquiringLabel{
			Key:   label.Key,
			Value: label.Value,
		}
	}
	if label := buildCSAutoRepairPolicyCreateLabel(interventionMap["approved_label"]); label != nil {
		intervention.ApprovedLabel = &roacs.CreateAutoRepairPolicyRequestRulesRepairProcedureInterventionApprovedLabel{
			Key:   label.Key,
			Value: label.Value,
		}
	}
	return intervention
}

type csAutoRepairPolicyLabel struct {
	Key   *string
	Value *string
}

func buildCSAutoRepairPolicyCreateLabel(raw interface{}) *csAutoRepairPolicyLabel {
	labelMap := firstMapOfList(raw)
	if labelMap == nil {
		return nil
	}
	return &csAutoRepairPolicyLabel{
		Key:   tea.String(toStr(labelMap["key"])),
		Value: tea.String(toStr(labelMap["value"])),
	}
}

func buildCSAutoRepairPolicyModifyRules(rules []interface{}) []*roacs.ModifyAutoRepairPolicyRequestRules {
	result := make([]*roacs.ModifyAutoRepairPolicyRequestRules, 0, len(rules))
	for _, rule := range rules {
		ruleMap := rule.(map[string]interface{})
		result = append(result, &roacs.ModifyAutoRepairPolicyRequestRules{
			Incidents:       buildCSAutoRepairPolicyModifyIncidents(ruleMap["incidents"].([]interface{})),
			RepairProcedure: buildCSAutoRepairPolicyModifyProcedures(ruleMap["repair_procedure"].([]interface{})),
		})
	}
	return result
}

func buildCSAutoRepairPolicyModifyIncidents(incidents []interface{}) []*roacs.ModifyAutoRepairPolicyRequestRulesIncidents {
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

func buildCSAutoRepairPolicyModifyProcedures(procedures []interface{}) []*roacs.ModifyAutoRepairPolicyRequestRulesRepairProcedure {
	result := make([]*roacs.ModifyAutoRepairPolicyRequestRulesRepairProcedure, 0, len(procedures))
	for _, procedure := range procedures {
		procedureMap := procedure.(map[string]interface{})
		item := &roacs.ModifyAutoRepairPolicyRequestRulesRepairProcedure{
			Name: tea.String(procedureMap["name"].(string)),
		}
		if config := buildCSAutoRepairPolicyConfig(procedureMap["config"]); config != nil {
			item.Config = config
		}
		if intervention := buildCSAutoRepairPolicyModifyIntervention(procedureMap["intervention"]); intervention != nil {
			item.Intervention = intervention
		}
		result = append(result, item)
	}
	return result
}

func buildCSAutoRepairPolicyModifyIntervention(raw interface{}) *roacs.ModifyAutoRepairPolicyRequestRulesRepairProcedureIntervention {
	interventionMap := firstMapOfList(raw)
	if interventionMap == nil {
		return nil
	}
	intervention := &roacs.ModifyAutoRepairPolicyRequestRulesRepairProcedureIntervention{
		Enable: tea.Bool(toBool(interventionMap["enable"])),
		Type:   tea.String(toStr(interventionMap["type"])),
	}
	if label := buildCSAutoRepairPolicyCreateLabel(interventionMap["inquiring_label"]); label != nil {
		intervention.InquiringLabel = &roacs.ModifyAutoRepairPolicyRequestRulesRepairProcedureInterventionInquiringLabel{
			Key:   label.Key,
			Value: label.Value,
		}
	}
	if label := buildCSAutoRepairPolicyCreateLabel(interventionMap["approved_label"]); label != nil {
		intervention.ApprovedLabel = &roacs.ModifyAutoRepairPolicyRequestRulesRepairProcedureInterventionApprovedLabel{
			Key:   label.Key,
			Value: label.Value,
		}
	}
	return intervention
}

func validateCSAutoRepairPolicyRules(raw interface{}) error {
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

func buildCSAutoRepairPolicyConfig(raw interface{}) map[string]interface{} {
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

func flattenCSAutoRepairPolicyRules(rules []*roacs.DescribeAutoRepairPolicyResponseBodyRules) []map[string]interface{} {
	result := make([]map[string]interface{}, 0, len(rules))
	for _, rule := range rules {
		incidents := make([]map[string]interface{}, 0, len(rule.Incidents))
		for _, incident := range rule.Incidents {
			incidents = append(incidents, map[string]interface{}{
				"name": tea.StringValue(incident.Name),
			})
		}
		procedures := make([]map[string]interface{}, 0, len(rule.RepairProcedure))
		for _, procedure := range rule.RepairProcedure {
			item := map[string]interface{}{
				"name":   tea.StringValue(procedure.Name),
				"config": flattenCSAutoRepairPolicyConfig(procedure.Config),
			}
			// the API always returns an intervention object, even when manual
			// approval is not enabled; only a enabled one is written to state,
			// otherwise the block would diff against a configuration that omits it
			if procedure.Intervention != nil && tea.BoolValue(procedure.Intervention.Enable) {
				item["intervention"] = flattenCSAutoRepairPolicyIntervention(procedure.Intervention)
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

func flattenCSAutoRepairPolicyIntervention(intervention *roacs.DescribeAutoRepairPolicyResponseBodyRulesRepairProcedureIntervention) []map[string]interface{} {
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

func flattenCSAutoRepairPolicyConfig(config map[string]interface{}) string {
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

func flattenCSAutoRepairPolicyResourceIds(ids []*string) []string {
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
	oldConfig, err := normalizeCSAutoRepairPolicyConfig(oldValue)
	if err != nil {
		return false
	}
	newConfig, err := normalizeCSAutoRepairPolicyConfig(newValue)
	if err != nil {
		return false
	}
	return reflect.DeepEqual(oldConfig, newConfig)
}

func normalizeCSAutoRepairPolicyConfig(raw string) (interface{}, error) {
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
