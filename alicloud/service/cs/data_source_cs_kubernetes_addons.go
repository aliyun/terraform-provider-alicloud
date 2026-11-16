package cs

import (
	"fmt"
	"regexp"
	"sort"

	"github.com/alibabacloud-go/tea/tea"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/convert"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func DataSourceCSKubernetesAddons() *schema.Resource {
	return &schema.Resource{
		Read: dataCSKubernetesAddonsRead,

		Schema: map[string]*schema.Schema{
			"cluster_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Computed: true,
			},
			"name_regex": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsValidRegExp,
			},
			"names": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"addons": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"current_version": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"next_version": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"required": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"current_config": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"supported_actions": {
							Type:     schema.TypeSet,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},
		},
	}
}

func dataCSKubernetesAddonsRead(d *schema.ResourceData, meta interface{}) error {
	clusterId := d.Get("cluster_id").(string)
	addons, err := describeAvailableAddons(d, meta)
	if err != nil {
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, ResourceAlicloudCSKubernetesAddon, "describeAddonsMeta", err)
	}

	var ids, names []string
	for _, addon := range addons {
		ids = append(ids, fmt.Sprintf("%s%s%s", clusterId, COLON_SEPARATED, addon.ComponentName))
		names = append(names, addon.ComponentName)
	}

	result := fetchAddonsMetadata(addons)

	d.Set("cluster_id", clusterId)
	d.Set("ids", ids)
	d.Set("names", names)
	d.Set("addons", result)

	d.SetId(tea.ToString(helper.Hashcode(clusterId)))
	return nil
}

// addonsNameMatcher turns the ids and name_regex arguments into a filter that is applied before
// the per-addon requests, instead of dropping their results afterwards.
func addonsNameMatcher(d *schema.ResourceData, clusterId string) (func(name string) bool, error) {
	var nameRegex *regexp.Regexp
	if nameRegexArg, ok := d.GetOk("name_regex"); ok {
		regex, err := regexp.Compile(nameRegexArg.(string))
		if err != nil {
			return nil, err
		}
		nameRegex = regex
	}

	var wantedIds map[string]bool
	if idsArg, ok := d.GetOk("ids"); ok {
		wantedIds = make(map[string]bool)
		for _, id := range convert.ExpandStringList(idsArg.([]interface{})) {
			wantedIds[id] = true
		}
	}

	if nameRegex == nil && wantedIds == nil {
		return nil, nil
	}

	return func(name string) bool {
		if nameRegex != nil && !nameRegex.MatchString(name) {
			return false
		}
		if wantedIds != nil && !wantedIds[fmt.Sprintf("%s%s%s", clusterId, COLON_SEPARATED, name)] {
			return false
		}
		return true
	}, nil
}

func describeAvailableAddons(d *schema.ResourceData, meta interface{}) (map[string]*Component, error) {
	clusterId := d.Get("cluster_id").(string)

	match, err := addonsNameMatcher(d, clusterId)
	if err != nil {
		return nil, err
	}

	client, err := meta.(*connectivity.AliyunClient).NewRoaCsClient()
	if err != nil {
		return nil, err
	}
	csClient := CsClient{client}

	availableAddons, err := csClient.DescribeCsKubernetesAllAvailableAddons(clusterId, match)
	if err != nil {
		return nil, errs.WrapErrorf(err, errs.DefaultErrorMsg, ResourceAlicloudCSKubernetesAddon, "DescribeCsKubernetesAllAvailableAddons", err)
	}

	return availableAddons, nil
}

func fetchAddonsMetadata(addonsMap map[string]*Component) []map[string]interface{} {
	result := []map[string]interface{}{}
	for name, addon := range addonsMap {
		state := map[string]interface{}{}
		state["name"] = name
		state["current_version"] = addon.Version
		state["next_version"] = addon.NextVersion
		state["required"] = addon.Required
		state["current_config"] = addon.Config
		state["supported_actions"] = addon.SupportedActions
		result = append(result, state)
	}
	sort.Slice(result, func(i, j int) bool {
		return result[i]["name"].(string) > result[j]["name"].(string)
	})
	return result
}
