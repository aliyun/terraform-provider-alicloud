package alicloud

import (
	"encoding/json"

	"github.com/alibabacloud-go/tea/tea"
	"github.com/hashicorp/terraform-plugin-sdk/helper/hashcode"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceAliCloudMaxComputeRolePolicyDocument() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAliCloudMaxComputeRolePolicyDocumentRead,
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
							Required: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
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

func dataSourceAliCloudMaxComputeRolePolicyDocumentRead(d *schema.ResourceData, meta interface{}) error {

	if v, ok := d.GetOk("statement"); ok {
		doc, err := AssembleDataSourceMaxComputeRolePolicyDocument(v.([]interface{}), d.Get("version").(string))
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

func AssembleDataSourceMaxComputeRolePolicyDocument(statements []interface{}, version string) (string, error) {
	var statementsList []PolicyDocumentStatement
	for _, v := range statements {
		statementMap := v.(map[string]interface{})

		actions, err := getOneStringOrAllStringSlice(statementMap["action"].([]interface{}), "action", true)
		if err != nil {
			return "", WrapError(err)
		}

		statement := PolicyDocumentStatement{
			Effect: Effect(statementMap["effect"].(string)),
			Action: actions,
		}

		if resources := statementMap["resource"].([]interface{}); len(resources) > 0 {
			resource, err := getOneStringOrAllStringSlice(resources, "resource", true)
			if err != nil {
				return "", WrapError(err)
			}

			statement.Resource = resource
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
