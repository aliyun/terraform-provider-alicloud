// Package alicloud. This file is hand-written from the Cms 2024-03-30 CloudSpec Umodel definition.
package alicloud

import (
	"fmt"
	"log"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceAlicloudCmsUmodel() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudCmsUmodelRead,
		Schema: map[string]*schema.Schema{
			"workspace": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"common_schema_ref": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"group": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"items": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"version": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"region_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func dataSourceAlicloudCmsUmodelRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	cmsServiceV2 := CmsServiceV2{client}

	workspace := d.Get("workspace").(string)
	objectRaw, err := cmsServiceV2.DescribeCmsUmodel(workspace)
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Data source alicloud_cms_umodel DescribeCmsUmodel Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.SetId(workspace)
	d.Set("description", objectRaw["description"])
	d.Set("region_id", objectRaw["regionId"])
	d.Set("workspace", objectRaw["workspace"])

	if v, ok := objectRaw["commonSchemaRef"]; ok && v != nil {
		if refs, ok := v.([]interface{}); ok {
			refList := make([]map[string]interface{}, 0, len(refs))
			for _, ref := range refs {
				refMap, ok := ref.(map[string]interface{})
				if !ok {
					continue
				}
				m := make(map[string]interface{})
				if g, ok := refMap["group"]; ok {
					m["group"] = g
				}
				if rawItems, ok := refMap["items"].([]interface{}); ok {
					items := make([]string, 0, len(rawItems))
					for _, it := range rawItems {
						items = append(items, fmt.Sprint(it))
					}
					m["items"] = items
				}
				if ver, ok := refMap["version"]; ok {
					m["version"] = ver
				}
				refList = append(refList, m)
			}
			d.Set("common_schema_ref", refList)
		}
	}

	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		if err := writeToFile(output.(string), flattenCmsUmodelMap(objectRaw)); err != nil {
			return WrapError(err)
		}
	}

	return nil
}

func flattenCmsUmodelMap(object map[string]interface{}) map[string]interface{} {
	result := map[string]interface{}{
		"workspace":   object["workspace"],
		"description": object["description"],
		"region_id":   object["regionId"],
	}
	if v, ok := object["commonSchemaRef"]; ok && v != nil {
		if refs, ok := v.([]interface{}); ok {
			refList := make([]map[string]interface{}, 0, len(refs))
			for _, ref := range refs {
				if refMap, ok := ref.(map[string]interface{}); ok {
					m := map[string]interface{}{
						"group":   refMap["group"],
						"version": refMap["version"],
					}
					if rawItems, ok := refMap["items"].([]interface{}); ok {
						items := make([]string, 0, len(rawItems))
						for _, it := range rawItems {
							items = append(items, fmt.Sprint(it))
						}
						m["items"] = items
					}
					refList = append(refList, m)
				}
			}
			result["common_schema_ref"] = refList
		}
	}
	return result
}
