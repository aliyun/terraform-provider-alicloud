// Package alicloud
package alicloud

import (
	"fmt"
	"os"
	"regexp"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceAliCloudCmsContextStores() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAliCloudCmsContextStoresRead,
		Timeouts: &schema.ResourceTimeout{
			Read: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"workspace": {
				Type:     schema.TypeString,
				Required: true,
			},
			"context_store_name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"context_type": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"ids": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"names": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"context_stores": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"context_store_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"context_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"workspace": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"region_id": {
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
					},
				},
			},
			"name_regex": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func dataSourceAliCloudCmsContextStoresRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	workspace := d.Get("workspace").(string)
	action := fmt.Sprintf("/workspace/%s/contextstore", workspace)
	var response map[string]interface{}
	var objects []map[string]interface{}
	query := make(map[string]*string)
	var err error

	if v, ok := d.GetOk("context_store_name"); ok {
		query["contextStoreName"] = StringPointer(v.(string))
	}
	if v, ok := d.GetOk("context_type"); ok {
		query["contextType"] = StringPointer(v.(string))
	}
	query["maxResults"] = StringPointer(fmt.Sprintf("%d", PageSizeLarge))

	wait := incrementalWait(3*time.Second, 5*time.Second)
	for {
		err = resource.Retry(d.Timeout(schema.TimeoutRead), func() *resource.RetryError {
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
		addDebug(action, response, nil)
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
		if results, ok := response["results"].([]interface{}); ok {
			for _, v := range results {
				if item, ok := v.(map[string]interface{}); ok {
					objects = append(objects, item)
				}
			}
		}
		if nextToken, ok := response["nextToken"].(string); ok && nextToken != "" {
			query["nextToken"] = StringPointer(nextToken)
		} else {
			break
		}
	}

	ids := make([]string, 0)
	names := make([]string, 0)
	s := make([]map[string]interface{}, 0)
	var nameRegex *regexp.Regexp
	if v, ok := d.GetOk("name_regex"); ok && v.(string) != "" {
		if r, err := regexp.Compile(v.(string)); err == nil {
			nameRegex = r
		}
	}
	for _, objectRaw := range objects {
		contextStoreName := fmt.Sprint(objectRaw["contextStoreName"])
		if nameRegex != nil && !nameRegex.MatchString(contextStoreName) {
			continue
		}
		ids = append(ids, fmt.Sprintf("%s:%s", objectRaw["workspace"], contextStoreName))
		names = append(names, contextStoreName)
		s = append(s, map[string]interface{}{
			"context_store_name": contextStoreName,
			"context_type":       objectRaw["contextType"],
			"workspace":          objectRaw["workspace"],
			"description":        objectRaw["description"],
			"status":             objectRaw["status"],
			"region_id":          objectRaw["regionId"],
			"create_time":        objectRaw["createTime"],
			"update_time":        objectRaw["updateTime"],
		})
	}
	d.Set("ids", ids)
	d.Set("names", names)
	d.Set("context_stores", s)

	if v, ok := d.GetOk("output_file"); ok && v.(string) != "" {
		_ = os.WriteFile(v.(string), []byte(""), 0644)
	}

	return nil
}
