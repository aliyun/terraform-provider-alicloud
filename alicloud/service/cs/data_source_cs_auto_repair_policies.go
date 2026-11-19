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

const DataSourceAlicloudAutoRepairPolicies = "alicloud_cs_auto_repair_policies"

// DataSourceAutoRepairPolicies returns the alicloud_cs_auto_repair_policies data source.
func DataSourceAutoRepairPolicies() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAutoRepairPoliciesRead,
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

func dataSourceAutoRepairPoliciesRead(d *schema.ResourceData, meta interface{}) error {
	client, clientErr := meta.(*connectivity.AliyunClient).NewRoaCsClient()
	if clientErr != nil {
		return errs.WrapError(clientErr)
	}

	clusterId := d.Get("cluster_id").(string)
	action := "ListAutoRepairPolicies"
	response, err := client.ListAutoRepairPolicies(tea.String(clusterId))
	logs.AddDebug(action, response)
	if err != nil {
		return errs.WrapErrorf(err, errs.DataDefaultErrorMsg, DataSourceAlicloudAutoRepairPolicies, action, errs.AlibabaCloudSdkGoERROR)
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

	type autoRepairPolicyEntry struct {
		id     string
		name   string
		policy map[string]interface{}
	}
	entries := make([]autoRepairPolicyEntry, 0)
	if response != nil && response.Body != nil {
		for _, item := range response.Body.Items {
			policyId := tea.StringValue(item.Id)
			name := tea.StringValue(item.Name)
			if policyId == "" {
				continue
			}
			id := fmt.Sprintf("%s%s%s", clusterId, COLON_SEPARATED, policyId)
			if len(idsFilter) > 0 && !helper.InArray(id, idsFilter) {
				continue
			}
			if nameRegex != nil && !nameRegex.MatchString(name) {
				continue
			}
			if hasTypeFilter && tea.StringValue(item.Type) != typeFilter.(string) {
				continue
			}
			resourcesIds := flattenAutoRepairPolicyResourceIds(item.ResourceIds)
			if hasSubTypeFilter && !matchAutoRepairResourceSubType(tea.StringValue(item.ResourceSubType), subTypeFilter.(string)) {
				continue
			}
			if hasNodepoolIdFilter && !helper.InArray(nodepoolIdFilter.(string), resourcesIds) {
				continue
			}

			rules, err := json.Marshal(item.Rules)
			if err != nil {
				return errs.WrapError(err)
			}
			entries = append(entries, autoRepairPolicyEntry{
				id:   id,
				name: name,
				policy: map[string]interface{}{
					"id":                id,
					"policy_id":         policyId,
					"name":              name,
					"type":              tea.StringValue(item.Type),
					"resource_type":     tea.StringValue(item.ResourceType),
					"resource_sub_type": tea.StringValue(item.ResourceSubType),
					"resource_ids":      resourcesIds,
					"rules":             string(rules),
				},
			})
		}
	}

	// ids, names and policies are parallel lists, so they must share one ordering
	sort.Slice(entries, func(i, j int) bool { return entries[i].id < entries[j].id })
	ids := make([]string, 0, len(entries))
	names := make([]string, 0, len(entries))
	policies := make([]map[string]interface{}, 0, len(entries))
	for _, entry := range entries {
		ids = append(ids, entry.id)
		names = append(names, entry.name)
		policies = append(policies, entry.policy)
	}

	d.Set("ids", ids)
	d.Set("names", names)
	if err = d.Set("policies", policies); err != nil {
		return errs.WrapError(err)
	}
	d.SetId(helper.DataResourceIdHash(ids))

	return nil
}

// matchAutoRepairResourceSubType treats a "|" separated value as a set, which
// is how builtin policies declare that they cover several node pool types.
func matchAutoRepairResourceSubType(raw, want string) bool {
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
