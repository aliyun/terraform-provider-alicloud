package cs

import (
	"encoding/json"
	"fmt"
	"regexp"
	"sort"
	"strings"

	"github.com/alibabacloud-go/tea/tea"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/convert"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

const DataSourceAlicloudCSAutoRepairPolicies = "alicloud_cs_auto_repair_policies"

// DataSourceCSAutoRepairPolicies returns the alicloud_cs_auto_repair_policies data source.
func DataSourceCSAutoRepairPolicies() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceCSAutoRepairPoliciesRead,
		Schema: map[string]*schema.Schema{
			"cluster_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"name_regex": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsValidRegExp,
			},
			"type": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"custom", "builtin"}, false),
			},
			"resource_sub_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"ess", "lingjun"}, false),
			},
			"nodepool_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"names": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
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
						"policy_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"resource_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"resource_sub_type": {
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
						"rules": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceCSAutoRepairPoliciesRead(d *schema.ResourceData, meta interface{}) error {
	client, clientErr := meta.(*connectivity.AliyunClient).NewRoaCsClient()
	if clientErr != nil {
		return errs.WrapError(clientErr)
	}

	clusterId := d.Get("cluster_id").(string)
	action := "ListAutoRepairPolicies"
	response, err := client.ListAutoRepairPolicies(tea.String(clusterId))
	logs.AddDebug(action, response)
	if err != nil {
		return errs.WrapErrorf(err, errs.DataDefaultErrorMsg, DataSourceAlicloudCSAutoRepairPolicies, action, errs.AlibabaCloudSdkGoERROR)
	}

	var nameRegex *regexp.Regexp
	if v, ok := d.GetOk("name_regex"); ok {
		nameRegex, err = regexp.Compile(v.(string))
		if err != nil {
			return errs.WrapError(err)
		}
	}
	idsFilter := convert.ExpandStringList(d.Get("ids").([]interface{}))
	typeFilter, hasTypeFilter := d.GetOk("type")
	subTypeFilter, hasSubTypeFilter := d.GetOk("resource_sub_type")
	nodepoolIdFilter, hasNodepoolIdFilter := d.GetOk("nodepool_id")

	var ids, names []string
	policies := make([]map[string]interface{}, 0)
	if response != nil && response.Body != nil {
		for _, item := range response.Body.Items {
			policyId := tea.StringValue(item.Id)
			name := tea.StringValue(item.Name)
			if policyId == "" {
				continue
			}
			if len(idsFilter) > 0 && !containsCSAutoRepairPolicyString(idsFilter, fmt.Sprintf("%s%s%s", clusterId, COLON_SEPARATED, policyId)) {
				continue
			}
			if nameRegex != nil && !nameRegex.MatchString(name) {
				continue
			}
			if hasTypeFilter && tea.StringValue(item.Type) != typeFilter.(string) {
				continue
			}
			resourcesIds := flattenCSAutoRepairPolicyResourceIds(item.ResourceIds)
			if hasSubTypeFilter && !matchCSAutoRepairResourceSubType(tea.StringValue(item.ResourceSubType), subTypeFilter.(string)) {
				continue
			}
			if hasNodepoolIdFilter && !containsCSAutoRepairPolicyString(resourcesIds, nodepoolIdFilter.(string)) {
				continue
			}

			rules, err := json.Marshal(item.Rules)
			if err != nil {
				return errs.WrapError(err)
			}
			ids = append(ids, fmt.Sprintf("%s%s%s", clusterId, COLON_SEPARATED, policyId))
			names = append(names, name)
			policies = append(policies, map[string]interface{}{
				"id":                fmt.Sprintf("%s%s%s", clusterId, COLON_SEPARATED, policyId),
				"policy_id":         policyId,
				"name":              name,
				"type":              tea.StringValue(item.Type),
				"resource_type":     tea.StringValue(item.ResourceType),
				"resource_sub_type": tea.StringValue(item.ResourceSubType),
				"resource_ids":      resourcesIds,
				"rules":             string(rules),
			})
		}
	}

	sort.Strings(ids)
	d.Set("ids", ids)
	d.Set("names", names)
	if err = d.Set("policies", policies); err != nil {
		return errs.WrapError(err)
	}
	d.SetId(helper.DataResourceIdHash(ids))

	return nil
}

// matchCSAutoRepairResourceSubType treats a "|" separated value as a set, which
// is how builtin policies declare that they cover several node pool types.
func matchCSAutoRepairResourceSubType(raw, want string) bool {
	if raw == want {
		return true
	}
	for _, part := range strings.Split(raw, "|") {
		if strings.TrimSpace(part) == want {
			return true
		}
	}
	return false
}

func containsCSAutoRepairPolicyString(values []string, target string) bool {
	for _, value := range values {
		if value == target {
			return true
		}
	}
	return false
}
