package alicloud

import (
	"encoding/json"
	"fmt"
	"regexp"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceAlicloudCmsContacts() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudCmsContactsRead,
		Schema: map[string]*schema.Schema{
			"ids": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"contact_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"phone": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"email": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"query_ungrouped_contacts": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"name_regex": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"contacts": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"contact_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"contact_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"email": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"phone": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"lang": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"workspace": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"im_user_ids": {
							Type:     schema.TypeMap,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},
		},
	}
}

func dataSourceAlicloudCmsContactsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	cmsServiceV2 := CmsServiceV2{client}
	objects, err := cmsServiceV2.ListCmsContacts(d, meta)
	if err != nil {
		return WrapError(err)
	}
	var nameRegex *regexp.Regexp
	if v, ok := d.GetOk("name_regex"); ok {
		if r, err := regexp.Compile(v.(string)); err == nil {
			nameRegex = r
		}
	}
	var ids []string
	var s []map[string]interface{}
	for _, object := range objects {
		contactName := fmt.Sprint(object["name"])
		if nameRegex != nil && !nameRegex.MatchString(contactName) {
			continue
		}
		mapping := map[string]interface{}{
			"contact_id":   object["contactId"],
			"contact_name": object["name"],
			"email":        object["email"],
			"phone":        object["phone"],
			"lang":         object["lang"],
			"workspace":    object["workspace"],
			"im_user_ids":  object["imUserIds"],
		}
		ids = append(ids, fmt.Sprint(object["contactId"]))
		s = append(s, mapping)
	}
	d.SetId(dataSourceAlicloudCmsContactsID(d, client, ids))
	d.Set("ids", ids)
	d.Set("contacts", s)
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}
	return nil
}

func dataSourceAlicloudCmsContactsID(d *schema.ResourceData, client *connectivity.AliyunClient, ids []string) string {
	parts := []string{client.RegionId}
	parts = append(parts, ids...)
	return dataResourceIdHash(parts)
}

// ListCmsContacts paginates the CMS ListContacts API.
func (s *CmsServiceV2) ListCmsContacts(d *schema.ResourceData, meta interface{}) (objects []map[string]interface{}, err error) {
	client := s.client
	action := "/contact"
	var response map[string]interface{}
	pageNumber := 1
	pageSize := 100
	for {
		query := make(map[string]*string)
		pageNumberStr := fmt.Sprintf("%d", pageNumber)
		pageSizeStr := fmt.Sprintf("%d", pageSize)
		query["pageNumber"] = &pageNumberStr
		query["pageSize"] = &pageSizeStr
		source := "OBS"
		query["source"] = &source
		if v, ok := d.GetOk("contact_name"); ok {
			nameStr := v.(string)
			query["name"] = &nameStr
		}
		if v, ok := d.GetOk("phone"); ok {
			phoneStr := v.(string)
			query["phone"] = &phoneStr
		}
		if v, ok := d.GetOk("email"); ok {
			emailStr := v.(string)
			query["email"] = &emailStr
		}
		if v, ok := d.GetOk("query_ungrouped_contacts"); ok {
			qucStr := fmt.Sprintf("%t", v.(bool))
			query["queryUngroupedContacts"] = &qucStr
		}
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(1*time.Minute, func() *resource.RetryError {
			response, err = client.RoaGet("Cms", "2024-03-30", action, query, nil, nil)
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})
		addDebug(action, response, query)
		if err != nil {
			return objects, WrapErrorf(err, DefaultErrorMsg, "data.alicloud_cms_contacts", action, AlibabaCloudSdkGoERROR)
		}
		var contactsArr []interface{}
		if dataRaw, ok := response["data"].(map[string]interface{}); ok {
			if arr, ok := dataRaw["data"].([]interface{}); ok {
				contactsArr = arr
			} else if arr, ok := dataRaw["contacts"].([]interface{}); ok {
				contactsArr = arr
			}
		} else if arr, ok := response["contacts"].([]interface{}); ok {
			contactsArr = arr
		}
		for _, raw := range contactsArr {
			if contact, ok := raw.(map[string]interface{}); ok {
				objects = append(objects, contact)
			}
		}
		var total int
		if dataRaw, ok := response["data"].(map[string]interface{}); ok {
			if t, ok := dataRaw["total"]; ok {
				total = cmsContactsTotalToInt(t)
			}
		} else if t, ok := response["total"]; ok {
			total = cmsContactsTotalToInt(t)
		}
		if len(objects) >= total || len(contactsArr) == 0 {
			break
		}
		pageNumber++
	}
	return objects, nil
}

// cmsContactsTotalToInt converts the ListContacts total field to int safely.
// The provider JSON decoder enables UseNumber, so total arrives as json.Number
// (not float64); a bare t.(float64) type assertion panics. Handle both forms
// for robustness, matching the provider convention (e.g. ess_scaling_rules).
func cmsContactsTotalToInt(t interface{}) int {
	if n, ok := t.(json.Number); ok {
		if v, err := n.Int64(); err == nil {
			return int(v)
		}
	} else if f, ok := t.(float64); ok {
		return int(f)
	}
	return 0
}
