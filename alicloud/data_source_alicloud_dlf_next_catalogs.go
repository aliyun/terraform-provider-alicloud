// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"log"
	"regexp"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceAliCloudDlfNextCatalogs() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAliCloudDlfNextCatalogsRead,
		Schema: map[string]*schema.Schema{
			"catalog_name_pattern": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"name_regex": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"catalogs": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"options": {
							Type:     schema.TypeMap,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"is_shared": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"share_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"region_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"owner": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"created_at": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"created_by": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"updated_at": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"updated_by": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceAliCloudDlfNextCatalogsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	dlfNextServiceV2 := DlfNextServiceV2{client}

	var allCatalogs []interface{}
	catalogNamePattern, _ := d.Get("catalog_name_pattern").(string)

	nextToken := ""
	maxResults := 100
	for {
		objects, token, err := dlfNextServiceV2.ListDlfNextCatalogs(catalogNamePattern, nextToken, maxResults)
		if err != nil {
			return WrapError(err)
		}
		allCatalogs = append(allCatalogs, objects...)
		if token == "" {
			break
		}
		nextToken = token
	}

	var catalogs []map[string]interface{}
	var ids []string
	idSet := make(map[string]bool)
	if v, ok := d.GetOk("ids"); ok {
		for _, id := range v.([]interface{}) {
			ids = append(ids, id.(string))
			idSet[id.(string)] = true
		}
	}

	var nameRegex *regexp.Regexp
	if v, ok := d.GetOk("name_regex"); ok && v.(string) != "" {
		r, err := regexp.Compile(v.(string))
		if err != nil {
			return WrapError(err)
		}
		nameRegex = r
	}

	for _, raw := range allCatalogs {
		object, ok := raw.(map[string]interface{})
		if !ok {
			continue
		}

		name := fmt.Sprint(object["name"])

		if len(ids) > 0 && !idSet[name] {
			continue
		}

		if nameRegex != nil && !nameRegex.MatchString(name) {
			continue
		}

		catalog := map[string]interface{}{
			"name":       name,
			"type":       object["type"],
			"options":    object["options"],
			"is_shared":  object["isShared"],
			"share_id":   object["shareId"],
			"id":         object["id"],
			"region_id":  object["regionId"],
			"status":     object["status"],
			"owner":      object["owner"],
			"created_at": object["createdAt"],
			"created_by": object["createdBy"],
			"updated_at": object["updatedAt"],
			"updated_by": object["updatedBy"],
		}
		catalogs = append(catalogs, catalog)
	}

	d.SetId(fmt.Sprintf("dlf_next_catalogs_%d", time.Now().Unix()))
	if err := d.Set("catalogs", catalogs); err != nil {
		return WrapError(err)
	}

	if v, ok := d.GetOk("output_file"); ok && v.(string) != "" {
		if err := writeToFile(v.(string), catalogs); err != nil {
			return WrapError(err)
		}
	}

	log.Printf("[DEBUG] dataSourceAliCloudDlfNextCatalogs: found %d catalogs", len(catalogs))

	return nil
}
