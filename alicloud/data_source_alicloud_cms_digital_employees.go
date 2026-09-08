package alicloud

import (
	"fmt"
	"regexp"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceAliCloudCmsDigitalEmployees() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAliCloudCmsDigitalEmployeesRead,
		Schema: map[string]*schema.Schema{
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"name_regex": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"digital_employee_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"display_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"employee_type": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"resource_group_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"digital_employees": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"digital_employee_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"role_arn": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"default_rule": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"display_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"resource_group_id": {
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
						"employee_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"region_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"tags": {
							Type:     schema.TypeSet,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"key": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"value": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func dataSourceAliCloudCmsDigitalEmployeesRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	cmsServiceV2 := CmsServiceV2{client}

	var nameRegex *regexp.Regexp
	if v, ok := d.GetOk("name_regex"); ok {
		r, err := regexp.Compile(v.(string))
		if err != nil {
			return WrapError(err)
		}
		nameRegex = r
	}

	idsMap := make(map[string]string)
	if v, ok := d.GetOk("ids"); ok {
		for _, vv := range v.([]interface{}) {
			if vv == nil {
				continue
			}
			idsMap[vv.(string)] = vv.(string)
		}
	}

	digitalEmployeeNameFilter, _ := d.Get("digital_employee_name").(string)
	displayNameFilter, _ := d.Get("display_name").(string)
	employeeTypeFilter, _ := d.Get("employee_type").(string)
	resourceGroupIdFilter, _ := d.Get("resource_group_id").(string)

	query := make(map[string]*string)
	query["maxResults"] = StringPointer(fmt.Sprintf("%d", PageSizeLarge))
	if digitalEmployeeNameFilter != "" {
		query["name"] = StringPointer(digitalEmployeeNameFilter)
	}
	if displayNameFilter != "" {
		query["displayName"] = StringPointer(displayNameFilter)
	}
	if employeeTypeFilter != "" {
		query["employeeType"] = StringPointer(employeeTypeFilter)
	}
	if resourceGroupIdFilter != "" {
		query["resourceGroupId"] = StringPointer(resourceGroupIdFilter)
	}

	objects, _, err := cmsServiceV2.ListCmsDigitalEmployees(query)
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_cms_digital_employees", "ListDigitalEmployees", AlibabaCloudSdkGoERROR)
	}

	filtered := make([]map[string]interface{}, 0, len(objects))
	for _, item := range objects {
		name := fmt.Sprintf("%v", item["name"])
		if nameRegex != nil && !nameRegex.MatchString(name) {
			continue
		}
		if len(idsMap) > 0 {
			if _, ok := idsMap[name]; !ok {
				continue
			}
		}
		filtered = append(filtered, item)
	}

	ids := make([]string, 0, len(filtered))
	mapping := make([]map[string]interface{}, 0, len(filtered))
	for _, objectRaw := range filtered {
		m := map[string]interface{}{
			"id":                    fmt.Sprintf("%v", objectRaw["name"]),
			"digital_employee_name": objectRaw["name"],
			"role_arn":              objectRaw["roleArn"],
			"default_rule":          objectRaw["defaultRule"],
			"description":           objectRaw["description"],
			"display_name":          objectRaw["displayName"],
			"resource_group_id":     objectRaw["resourceGroupId"],
			"create_time":           objectRaw["createTime"],
			"update_time":           objectRaw["updateTime"],
			"employee_type":         objectRaw["employeeType"],
			"region_id":             objectRaw["regionId"],
		}
		if tags, ok := objectRaw["tags"].([]interface{}); ok && len(tags) > 0 {
			m["tags"] = flattenCmsDigitalEmployeeTags(tags)
		}
		ids = append(ids, m["id"].(string))
		mapping = append(mapping, m)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}
	if err := d.Set("digital_employees", mapping); err != nil {
		return WrapError(err)
	}

	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		if err := writeToFile(output.(string), mapping); err != nil {
			return WrapError(err)
		}
	}

	// ensure non-empty output when no results (avoid nil-set warning)
	return nil
}
