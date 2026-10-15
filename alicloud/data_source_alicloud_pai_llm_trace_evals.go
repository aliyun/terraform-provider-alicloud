package alicloud

import (
	"fmt"
	"regexp"
	"sort"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func dataSourceAlicloudPaiLlmTraceEvals() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudPaiLlmTraceEvalsRead,
		Schema: map[string]*schema.Schema{
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"name_regex": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsValidRegExp,
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"names": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"evals": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"app_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"data_source": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"eval_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"eval_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"gmt_create_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"metadata": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"record_count": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"region_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceAlicloudPaiLlmTraceEvalsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	paillmtraceServiceV2 := PaillmtraceServiceV2{client}

	ids := make([]string, 0)
	names := make([]string, 0)
	objects := make([]map[string]interface{}, 0)

	pageNumber := 1
	pageSize := 50
	for {
		rawObjects, totalCount, err := paillmtraceServiceV2.ListPaillmtraceEvals(pageNumber, pageSize)
		if err != nil {
			return WrapError(err)
		}

		for _, raw := range rawObjects {
			object, ok := raw.(map[string]interface{})
			if !ok {
				continue
			}
			evalId := fmt.Sprint(object["EvalId"])
			evalName := fmt.Sprint(object["EvaluationName"])

			if v, ok := d.GetOk("ids"); ok {
				filterIds := v.([]interface{})
				matched := false
				for _, fid := range filterIds {
					if evalId == fid.(string) {
						matched = true
						break
					}
				}
				if !matched {
					continue
				}
			}

			if v, ok := d.GetOk("name_regex"); ok {
				r := regexp.MustCompile(v.(string))
				if !r.MatchString(evalName) {
					continue
				}
			}

			ids = append(ids, evalId)
			names = append(names, evalName)
			objects = append(objects, map[string]interface{}{
				"app_name":        object["AppName"],
				"data_source":     object["DataSource"],
				"description":     object["Description"],
				"eval_id":         evalId,
				"eval_name":       evalName,
				"gmt_create_time": object["GmtCreateTime"],
				"id":              evalId,
				"metadata":        object["Metadata"],
				"record_count":    object["RecordCount"],
				"region_id":       client.RegionId,
			})
		}

		if pageNumber*pageSize >= totalCount {
			break
		}
		pageNumber++
	}

	sort.Strings(names)

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}
	if err := d.Set("names", names); err != nil {
		return WrapError(err)
	}

	if v, ok := d.GetOk("output_file"); ok && v.(string) != "" {
		if err := writeToFile(v.(string), objects); err != nil {
			return WrapError(err)
		}
	}

	return d.Set("evals", objects)
}
