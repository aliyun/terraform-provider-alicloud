package alicloud

import (
	"fmt"
	"strconv"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceAliCloudCmsEntityGroups() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAliCloudCmsEntityGroupsRead,
		Schema: map[string]*schema.Schema{
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"entity_group_type": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"workspace": {
				Type:     schema.TypeString,
				Required: true,
			},
			"groups": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"entity_group_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"entity_group_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"workspace": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"resource_group_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"entity_rules": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"resource_group_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"tags": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"op":         {Type: schema.TypeString, Computed: true},
												"tag_key":    {Type: schema.TypeString, Computed: true},
												"tag_values": {Type: schema.TypeList, Computed: true, Elem: &schema.Schema{Type: schema.TypeString}},
											},
										},
									},
									"labels": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"op":         {Type: schema.TypeString, Computed: true},
												"tag_key":    {Type: schema.TypeString, Computed: true},
												"tag_values": {Type: schema.TypeList, Computed: true, Elem: &schema.Schema{Type: schema.TypeString}},
											},
										},
									},
									"ip_match_rule": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"ip_field_key": {Type: schema.TypeString, Computed: true},
												"ip_cidr":      {Type: schema.TypeString, Computed: true},
											},
										},
									},
									"instance_ids": {
										Type:     schema.TypeList,
										Computed: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
									"field_rules": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"field_key":    {Type: schema.TypeString, Computed: true},
												"op":           {Type: schema.TypeString, Computed: true},
												"field_values": {Type: schema.TypeList, Computed: true, Elem: &schema.Schema{Type: schema.TypeString}},
											},
										},
									},
									"entity_types": {
										Type:     schema.TypeList,
										Computed: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
									"annotations": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"op":         {Type: schema.TypeString, Computed: true},
												"tag_key":    {Type: schema.TypeString, Computed: true},
												"tag_values": {Type: schema.TypeList, Computed: true, Elem: &schema.Schema{Type: schema.TypeString}},
											},
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
		},
	}
}

func dataSourceAliCloudCmsEntityGroupsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	idsMap := make(map[string]string)
	if v, ok := d.GetOk("ids"); ok {
		for _, vv := range v.([]interface{}) {
			if vv == nil {
				continue
			}
			idsMap[vv.(string)] = vv.(string)
		}
	}

	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]*string
	action := "/entity-groups"
	var err error
	request = make(map[string]interface{})
	query = make(map[string]*string)
	query["regionId"] = StringPointer(client.RegionId)
	query["workspace"] = StringPointer(d.Get("workspace").(string))
	if v, ok := d.GetOk("entity_group_type"); ok && v.(string) != "" {
		query["entityGroupType"] = StringPointer(v.(string))
	}
	query["MaxResults"] = StringPointer(strconv.Itoa(PageSizeLarge))

	var objects []map[string]interface{}
	for {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutRead), func() *resource.RetryError {
			response, err = client.RoaGet("Cms", "2024-03-30", action, query, nil, nil)
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			addDebug(action, response, request)
			return nil
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}

		resp, _ := jsonpath.Get("$.entityGroups[*]", response)
		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})
			if len(idsMap) > 0 {
				if _, ok := idsMap[fmt.Sprint(item["entityGroupId"])]; !ok {
					continue
				}
			}
			objects = append(objects, item)
		}

		if nextToken, ok := response["nextToken"].(string); ok && nextToken != "" {
			query["nextToken"] = StringPointer(nextToken)
		} else {
			break
		}
	}

	var nameFilter string
	if v, ok := d.GetOk("name"); ok {
		nameFilter = v.(string)
	}

	ids := make([]string, 0)
	s := make([]map[string]interface{}, 0)
	for _, objectRaw := range objects {
		entityGroupName := fmt.Sprint(objectRaw["entityGroupName"])
		if nameFilter != "" && nameFilter != entityGroupName {
			continue
		}

		mapping := map[string]interface{}{}
		entityGroupId := fmt.Sprint(objectRaw["entityGroupId"])
		mapping["entity_group_id"] = entityGroupId
		mapping["entity_group_name"] = objectRaw["entityGroupName"]
		mapping["description"] = objectRaw["description"]
		mapping["workspace"] = objectRaw["workspace"]

		resourceGroupId := ""
		if entityRules, ok := objectRaw["entityRules"].(map[string]interface{}); ok {
			if rgId, ok := entityRules["resourceGroupId"].(string); ok {
				resourceGroupId = rgId
			}
		}
		mapping["resource_group_id"] = resourceGroupId
		mapping["entity_rules"] = flattenCmsEntityRules(objectRaw["entityRules"])

		ids = append(ids, entityGroupId)
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}
	if err := d.Set("groups", s); err != nil {
		return WrapError(err)
	}

	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		if err := writeToFile(output.(string), s); err != nil {
			return WrapError(err)
		}
	}
	return nil
}
