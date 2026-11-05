package dlfnext

import (
	"fmt"
	"log"
	"regexp"
	"sort"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func DataSourceCatalogs() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceCatalogsRead,
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
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"names": {
				Type:     schema.TypeList,
				Computed: true,
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
							Type:     schema.TypeString,
							Computed: true,
						},
						"created_by": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"updated_at": {
							Type:     schema.TypeString,
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

func dataSourceCatalogsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	service := NewService(client)

	var allCatalogs []interface{}
	catalogNamePattern, _ := d.Get("catalog_name_pattern").(string)

	nextToken := ""
	maxResults := 100
	for {
		objects, token, err := service.ListCatalogs(catalogNamePattern, nextToken, maxResults)
		if err != nil {
			return errs.WrapError(err)
		}
		allCatalogs = append(allCatalogs, objects...)
		if token == "" {
			break
		}
		nextToken = token
	}

	var catalogs []map[string]interface{}
	var matchedNames []string
	idSet := make(map[string]bool)
	if v, ok := d.GetOk("ids"); ok {
		for _, n := range v.([]interface{}) {
			idSet[n.(string)] = true
		}
	}
	hasIDFilter := len(idSet) > 0

	var nameRegex *regexp.Regexp
	if v, ok := d.GetOk("name_regex"); ok && v.(string) != "" {
		r, err := regexp.Compile(v.(string))
		if err != nil {
			return errs.WrapError(err)
		}
		nameRegex = r
	}

	for _, raw := range allCatalogs {
		object, ok := raw.(map[string]interface{})
		if !ok {
			continue
		}

		name := fmt.Sprint(object["name"])

		if hasIDFilter && !idSet[name] {
			continue
		}

		if nameRegex != nil && !nameRegex.MatchString(name) {
			continue
		}

		matchedNames = append(matchedNames, name)
		// The Terraform identity of a catalog is its name (the resource sets
		// its id to the name and Import takes the name), so expose the same
		// value as the list item id to keep ids/import usable.
		catalog := map[string]interface{}{
			"name":       name,
			"type":       object["type"],
			"options":    object["options"],
			"is_shared":  object["isShared"],
			"share_id":   object["shareId"],
			"id":         name,
			"region_id":  object["regionId"],
			"status":     object["status"],
			"owner":      object["owner"],
			"created_at": catalogTimestampString(object["createdAt"]),
			"created_by": object["createdBy"],
			"updated_at": catalogTimestampString(object["updatedAt"]),
			"updated_by": object["updatedBy"],
		}
		catalogs = append(catalogs, catalog)
	}

	// ListCatalogs may return the same set in a different order across calls
	// (pagination, backend sharding). Sort both the list and the matched-name
	// outputs so state stays stable across refreshes and catalogs.0.name
	// agrees with names.0; the data source id is derived from the same sorted
	// names.
	matchedNames = sortCatalogNames(matchedNames)
	catalogs = sortCatalogsByName(catalogs)
	d.SetId(catalogsID(matchedNames))
	if err := d.Set("ids", matchedNames); err != nil {
		return errs.WrapError(err)
	}
	if err := d.Set("names", matchedNames); err != nil {
		return errs.WrapError(err)
	}
	if err := d.Set("catalogs", catalogs); err != nil {
		return errs.WrapError(err)
	}

	if v, ok := d.GetOk("output_file"); ok && v.(string) != "" {
		if err := helper.WriteToFile(v.(string), catalogs); err != nil {
			return errs.WrapError(err)
		}
	}

	log.Printf("[DEBUG] DataSourceCatalogs: found %d catalogs", len(catalogs))

	return nil
}

// catalogsID returns a stable data source id from a set of matched catalog
// names. The names are sorted before hashing so the same set of catalogs
// yields the same id regardless of the order ListCatalogs returned them in.
func catalogsID(names []string) string {
	sorted := make([]string, len(names))
	copy(sorted, names)
	sort.Strings(sorted)
	return helper.DataResourceIdHash(sorted)
}

// sortCatalogNames returns a sorted copy of the matched catalog names so the
// order-sensitive ids/names lists are deterministic and agree with the sorted
// catalogs list.
func sortCatalogNames(names []string) []string {
	sorted := make([]string, len(names))
	copy(sorted, names)
	sort.Strings(sorted)
	return sorted
}

// sortCatalogsByName returns a copy of catalogs sorted by the catalog "name"
// field so the TypeList stays stable across refreshes.
func sortCatalogsByName(catalogs []map[string]interface{}) []map[string]interface{} {
	sorted := make([]map[string]interface{}, len(catalogs))
	copy(sorted, catalogs)
	sort.SliceStable(sorted, func(i, j int) bool {
		return fmt.Sprint(sorted[i]["name"]) < fmt.Sprint(sorted[j]["name"])
	})
	return sorted
}
