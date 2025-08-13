package alicloud

import (
	"encoding/json"
	"fmt"
	"sort"

	"github.com/alibabacloud-go/tea/tea"
	"github.com/hashicorp/terraform-plugin-sdk/helper/hashcode"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceAliCloudRamPolicyDocument() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAliCloudRamPolicyDocumentRead,
		Schema: map[string]*schema.Schema{
			"version": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "1",
				ValidateFunc: StringInSlice([]string{"1"}, false),
			},
			"statement": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"effect": {
							Type:         schema.TypeString,
							Optional:     true,
							Default:      "Allow",
							ValidateFunc: StringInSlice([]string{"Allow", "Deny"}, false),
						},
						"action": {
							Type:     schema.TypeList,
							Required: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"resource": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"principal": {
							Type:     schema.TypeSet,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"entity": {
										Type:         schema.TypeString,
										Required:     true,
										ValidateFunc: StringInSlice([]string{"RAM", "Service", "Federated"}, false),
									},
									"identifiers": {
										Type:     schema.TypeList,
										Required: true,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
								},
							},
						},
						"condition": {
							Type:     schema.TypeSet,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"operator": {
										Type:     schema.TypeString,
										Required: true,
									},
									"variable": {
										Type:     schema.TypeString,
										Required: true,
									},
									"values": {
										Type:     schema.TypeList,
										Required: true,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
								},
							},
						},
					},
				},
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"document": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceAliCloudRamPolicyDocumentRead(d *schema.ResourceData, meta interface{}) error {

	if v, ok := d.GetOk("statement"); ok {
		doc, err := AssembleDataSourcePolicyDocument(v.([]interface{}), d.Get("version").(string))
		if err != nil {
			return WrapError(err)
		}

		if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
			writeToFile(output.(string), doc)
		}

		d.Set("document", doc)

		d.SetId(tea.ToString(hashcode.String(doc)))
	}

	return nil
}

func AssembleDataSourcePolicyDocument(statements []interface{}, version string) (string, error) {
	var statementsList []PolicyDocumentStatement
	for _, v := range statements {
		statementMap := v.(map[string]interface{})

		actions, err := getOneStringOrAllStringSlice(statementMap["action"].([]interface{}), "action")
		if err != nil {
			return "", WrapError(err)
		}

		statement := PolicyDocumentStatement{
			Effect: Effect(statementMap["effect"].(string)),
			Action: actions,
		}

		if resources := statementMap["resource"].([]interface{}); len(resources) > 0 {
			resource, err := getOneStringOrAllStringSlice(resources, "resource")
			if err != nil {
				return "", WrapError(err)
			}

			statement.Resource = resource
		}

		principalSlice := make(PolicyDocumentStatementPrincipalSet, 0)
		if principals := statementMap["principal"].(*schema.Set).List(); len(principals) > 0 {
			for _, principal := range principals {
				principalArg := principal.(map[string]interface{})
				principalObject := PolicyDocumentStatementPrincipal{}
				principalObject.Entity = principalArg["entity"].(string)
				identifiersSlice := make([]string, 0)
				for _, v := range principalArg["identifiers"].([]interface{}) {
					identifiersSlice = append(identifiersSlice, v.(string))
				}
				identifiers := identifiersSlice
				principalObject.Identifiers = identifiers
				principalSlice = append(principalSlice, principalObject)
			}
			statement.Principal = principalSlice
		}

		conditionSlice := make(PolicyDocumentStatementConditionSet, 0)
		if conditions := statementMap["condition"].(*schema.Set).List(); len(conditions) > 0 {
			for _, condition := range conditions {
				conditionArg := condition.(map[string]interface{})
				conditionObject := PolicyDocumentStatementCondition{}

				conditionObject.Operator = conditionArg["operator"].(string)
				conditionObject.Variable = conditionArg["variable"].(string)
				values, err := getOneStringOrAllStringSlice(conditionArg["values"].([]interface{}), "values")
				if err != nil {
					return "", WrapError(err)
				}

				conditionObject.Values = values
				conditionSlice = append(conditionSlice, conditionObject)
			}
			statement.Condition = conditionSlice
		}
		statementsList = append(statementsList, statement)
	}

	policyDocument := PolicyDocument{
		Version:   version,
		Statement: statementsList,
	}

	data, err := json.Marshal(policyDocument)
	if err != nil {
		return "", WrapError(err)
	}

	return string(data), nil
}

func (s PolicyDocumentStatementPrincipalSet) MarshalJSON() ([]byte, error) {
	raw := map[string]interface{}{}

	for _, p := range s {
		switch i := p.Identifiers.(type) {
		case []string:
			switch v := raw[p.Entity].(type) {
			case nil:
				raw[p.Entity] = make([]string, 0, len(i))
			case string:
				raw[p.Entity] = make([]string, 0, len(i)+1)
				raw[p.Entity] = append(raw[p.Entity].([]string), v)
			}
			sort.Sort(sort.Reverse(sort.StringSlice(i)))
			raw[p.Entity] = append(raw[p.Entity].([]string), i...)
		case string:
			switch v := raw[p.Entity].(type) {
			case nil:
				raw[p.Entity] = i
			case string:
				// Convert to []string to stop drop of principals
				raw[p.Entity] = make([]string, 0, 2)
				raw[p.Entity] = append(raw[p.Entity].([]string), v)
				raw[p.Entity] = append(raw[p.Entity].([]string), i)
			case []string:
				raw[p.Entity] = append(raw[p.Entity].([]string), i)
			}
		default:
			return []byte{}, fmt.Errorf("Unsupported data type %T for PolicyDocumentStatementPrincipalSet", i)
		}
	}

	return json.Marshal(&raw)
}

func (s *PolicyDocumentStatementPrincipalSet) UnmarshalJSON(b []byte) error {
	var out PolicyDocumentStatementPrincipalSet

	var data interface{}
	if err := json.Unmarshal(b, &data); err != nil {
		return err
	}

	switch t := data.(type) {
	case string:
		out = append(out, PolicyDocumentStatementPrincipal{Entity: "*", Identifiers: []string{"*"}})
	case map[string]interface{}:
		for key, value := range data.(map[string]interface{}) {
			switch vt := value.(type) {
			case string:
				out = append(out, PolicyDocumentStatementPrincipal{Entity: key, Identifiers: value.(string)})
			case []interface{}:
				var values []string
				for _, v := range value.([]interface{}) {
					values = append(values, v.(string))
				}
				out = append(out, PolicyDocumentStatementPrincipal{Entity: key, Identifiers: values})
			default:
				return fmt.Errorf("Unsupported data type %T for PolicyDocumentStatementPrincipalSet.Identifiers", vt)
			}
		}
	default:
		return fmt.Errorf("Unsupported data type %T for PolicyDocumentStatementPrincipalSet", t)
	}

	*s = out

	return nil
}

func (s PolicyDocumentStatementConditionSet) MarshalJSON() ([]byte, error) {
	raw := map[string]map[string]interface{}{}

	for _, c := range s {
		if _, ok := raw[c.Operator]; !ok {
			raw[c.Operator] = map[string]interface{}{}
		}
		switch i := c.Values.(type) {
		case []string:
			if _, ok := raw[c.Operator][c.Variable]; !ok {
				raw[c.Operator][c.Variable] = make([]string, 0, len(i))
			}
			// order matters with values so not sorting here
			raw[c.Operator][c.Variable] = append(raw[c.Operator][c.Variable].([]string), i...)
		case string:
			raw[c.Operator][c.Variable] = i
		default:
			return nil, fmt.Errorf("Unsupported data type for PolicyStatementConditionSet: %s", i)
		}
	}

	return json.Marshal(&raw)
}

func (s *PolicyDocumentStatementConditionSet) UnmarshalJSON(b []byte) error {
	var out PolicyDocumentStatementConditionSet

	var data map[string]map[string]interface{}
	if err := json.Unmarshal(b, &data); err != nil {
		return err
	}

	for operator_key, operator_value := range data {
		for var_key, var_values := range operator_value {
			switch var_values := var_values.(type) {
			case string:
				out = append(out, PolicyDocumentStatementCondition{Operator: operator_key, Variable: var_key, Values: []string{var_values}})
			case []interface{}:
				values := []string{}
				for _, v := range var_values {
					values = append(values, v.(string))
				}
				out = append(out, PolicyDocumentStatementCondition{Operator: operator_key, Variable: var_key, Values: values})
			}
		}
	}

	*s = out

	return nil
}
