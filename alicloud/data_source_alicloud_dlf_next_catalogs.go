// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"log"
	"regexp"
	"sort"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
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
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsValidRegExp,
			},
			"names": {
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
	var matchedNames []string
	nameSet := make(map[string]bool)
	if v, ok := d.GetOk("names"); ok {
		for _, n := range v.([]interface{}) {
			nameSet[n.(string)] = true
		}
	}
	hasNameFilter := len(nameSet) > 0

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

		if hasNameFilter && !nameSet[name] {
			continue
		}

		if nameRegex != nil && !nameRegex.MatchString(name) {
			continue
		}

		matchedNames = append(matchedNames, name)
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

	// Stable id derived from the matched catalog names (repo convention for
	// list data sources: dataResourceIdHash). The names are sorted before
	// hashing so identical inputs yield an identical id across runs even when
	// ListCatalogs returns the same set in a different order, instead of a
	// per-call ordering that would force a spurious diff on every apply.
	// ListCatalogs may return the same set of catalogs in a different order
	// across calls (pagination, backend sharding). Sort the catalogs written to
	// the TypeList by name so Terraform state stays stable across refreshes;
	// the data source id is derived separately from the matched names below.
	catalogs = sortDlfNextCatalogsByName(catalogs)
	d.SetId(dlfNextCatalogsID(matchedNames))
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

// dlfNextCatalogsID returns a stable data source id from a set of matched
// catalog names. The names are sorted before hashing so the same set of
// catalogs yields the same id regardless of the order ListCatalogs returned
// them in (avoiding a spurious diff on every apply). Pure function for unit
// testing.
func dlfNextCatalogsID(names []string) string {
	sorted := make([]string, len(names))
	copy(sorted, names)
	sort.Strings(sorted)
	return dataResourceIdHash(sorted)
}

// sortDlfNextCatalogsByName returns a copy of catalogs sorted by the catalog
// "name" field. ListCatalogs may return the same set of catalogs in a
// different order across calls (pagination, backend sharding), so the
// catalogs written to the TypeList are deterministically sorted by name to
// keep Terraform state stable across refreshes — otherwise the same set of
// catalogs returned in a different order would reorder the list in state and
// force a spurious diff. Pure function for unit testing; the data source id
// is derived separately via dlfNextCatalogsID from the matched names.
func sortDlfNextCatalogsByName(catalogs []map[string]interface{}) []map[string]interface{} {
	sorted := make([]map[string]interface{}, len(catalogs))
	copy(sorted, catalogs)
	sort.SliceStable(sorted, func(i, j int) bool {
		return fmt.Sprint(sorted[i]["name"]) < fmt.Sprint(sorted[j]["name"])
	})
	return sorted
}
