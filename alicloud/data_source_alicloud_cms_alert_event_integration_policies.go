package alicloud

import (
	"fmt"
	"regexp"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceAlicloudCmsAlertEventIntegrationPolicies() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudCmsAlertEventIntegrationPoliciesRead,
		Schema: map[string]*schema.Schema{
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"name_regex": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"workspace": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"enable": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"names": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"policies": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"alert_event_integration_policy_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"workspace": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"integration_setting": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"enable": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"region_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"create_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"update_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"user_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"filter_setting": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"conditions": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"field": {Type: schema.TypeString, Computed: true},
												"value": {Type: schema.TypeString, Computed: true},
												"op":    {Type: schema.TypeString, Computed: true},
											},
										},
									},
									"expression": {Type: schema.TypeString, Computed: true},
									"relation":   {Type: schema.TypeString, Computed: true},
								},
							},
						},
						"transformer_setting": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"filter_setting": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"conditions": {
													Type:     schema.TypeList,
													Computed: true,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"field": {Type: schema.TypeString, Computed: true},
															"value": {Type: schema.TypeString, Computed: true},
															"op":    {Type: schema.TypeString, Computed: true},
														},
													},
												},
												"expression": {Type: schema.TypeString, Computed: true},
												"relation":   {Type: schema.TypeString, Computed: true},
											},
										},
									},
									"label_key": {Type: schema.TypeString, Computed: true},
									"mapping": {
										Type:     schema.TypeMap,
										Elem:     &schema.Schema{Type: schema.TypeString},
										Computed: true,
									},
									"reg_exp":  {Type: schema.TypeString, Computed: true},
									"source":   {Type: schema.TypeString, Computed: true},
									"target":   {Type: schema.TypeString, Computed: true},
									"type":     {Type: schema.TypeString, Computed: true},
									"value":    {Type: schema.TypeString, Computed: true},
									"variable": {Type: schema.TypeString, Computed: true},
								},
							},
						},
					},
				},
			},
			"next_token": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceAlicloudCmsAlertEventIntegrationPoliciesRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	cmsServiceV2 := CmsServiceV2{client}

	workspace := d.Get("workspace").(string)
	args := &struct {
		AlertEventIntegrationPolicyName *string
		AlertEventIntegrationPolicyId   *string
		Enable                          *bool
		NextToken                       *string
		MaxResults                      *int
	}{}

	if v, ok := d.GetOk("enable"); ok {
		enableVal := v.(bool)
		args.Enable = &enableVal
	}

	response, err := cmsServiceV2.ListCmsAlertEventIntegrationPolicies(workspace, args)
	if err != nil {
		return WrapError(err)
	}

	listRaw, _ := jsonpath.Get("$.dataList", response)
	policies := make([]map[string]interface{}, 0)
	if listRaw != nil {
		for _, item := range convertToInterfaceArray(listRaw) {
			p, ok := item.(map[string]interface{})
			if !ok {
				continue
			}
			policies = append(policies, p)
		}
	}

	ids := make([]string, 0)
	names := make([]string, 0)
	out := make([]map[string]interface{}, 0)

	var nameRegex *regexp.Regexp
	if v, ok := d.GetOk("name_regex"); ok {
		if r, err := regexp.Compile(v.(string)); err == nil {
			nameRegex = r
		}
	}

	filterIds := make(map[string]bool)
	if v, ok := d.GetOk("ids"); ok {
		for _, id := range v.([]interface{}) {
			filterIds[fmt.Sprint(id)] = true
		}
	}

	for _, p := range policies {
		policyId := fmt.Sprint(p["alertEventIntegrationPolicyId"])
		policyName := fmt.Sprint(p["alertEventIntegrationPolicyName"])

		if len(filterIds) > 0 && !filterIds[policyId] {
			continue
		}
		if nameRegex != nil && !nameRegex.MatchString(policyName) {
			continue
		}

		ids = append(ids, policyId)
		names = append(names, policyName)
		entry := map[string]interface{}{
			"id":                                  policyId,
			"alert_event_integration_policy_name": policyName,
			"description":                         p["description"],
			"type":                                p["type"],
			"workspace":                           p["workspace"],
			"integration_setting":                 p["integrationSetting"],
			"enable":                              p["enabled"],
			"region_id":                           p["regionId"],
			"create_time":                         p["createTime"],
			"update_time":                         p["updateTime"],
			"user_id":                             p["userId"],
		}
		if fsList := flattenCmsAlertEventIntegrationPolicyFilterSetting(p["filterSetting"]); len(fsList) > 0 {
			entry["filter_setting"] = fsList
		} else {
			entry["filter_setting"] = []map[string]interface{}{}
		}
		if tsList := flattenCmsAlertEventIntegrationPolicyTransformerSetting(p["transformerSetting"]); len(tsList) > 0 {
			entry["transformer_setting"] = tsList
		} else {
			entry["transformer_setting"] = []map[string]interface{}{}
		}
		out = append(out, entry)
	}

	d.SetId(dataSourceAlicloudCmsAlertEventIntegrationPoliciesId(workspace))
	d.Set("policies", out)
	d.Set("names", names)

	if nextToken, _ := jsonpath.Get("$.nextToken", response); nextToken != nil {
		d.Set("next_token", fmt.Sprint(nextToken))
	}

	if v, ok := d.GetOk("output_file"); ok && v.(string) != "" {
		if err := writeToFile(v.(string), out); err != nil {
			return WrapError(err)
		}
	}

	return nil
}

func dataSourceAlicloudCmsAlertEventIntegrationPoliciesId(workspace string) string {
	if workspace != "" {
		return fmt.Sprintf("cms_alert_event_integration_policies:%s", workspace)
	}
	return "cms_alert_event_integration_policies"
}
