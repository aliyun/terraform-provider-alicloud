package alicloud

import (
	"fmt"
	"regexp"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceAlicloudRealtimeComputeVariables() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudRealtimeComputeVariablesRead,
		Schema: map[string]*schema.Schema{
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"workspace": {
				Type:     schema.TypeString,
				Required: true,
			},
			"namespace": {
				Type:     schema.TypeString,
				Required: true,
			},
			"name_regex": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"variables": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"workspace": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"namespace": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"kind": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"value": {
							Type:      schema.TypeString,
							Computed:  true,
							Sensitive: true,
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceAlicloudRealtimeComputeVariablesRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	workspace := d.Get("workspace").(string)
	namespace := d.Get("namespace").(string)
	realtimeComputeServiceV2 := RealtimeComputeServiceV2{client}

	// ListVariables is paginated inside the service function (pageSize/pageIndex with
	// totalSize), so this data source never misses variables beyond the first page.
	objects, err := realtimeComputeServiceV2.DescribeRealtimeComputeVariables(fmt.Sprintf("%v:%v", workspace, namespace))
	if err != nil {
		if NotFoundError(err) {
			d.SetId(fmt.Sprintf("%v:%v", workspace, namespace))
			return nil
		}
		return WrapError(err)
	}

	var nameRegex *regexp.Regexp
	if v, ok := d.GetOk("name_regex"); ok && v.(string) != "" {
		r, e := regexp.Compile(v.(string))
		if e != nil {
			return WrapError(e)
		}
		nameRegex = r
	}

	ids := make([]string, 0)
	filteredObjects := make([]map[string]interface{}, 0)
	for _, item := range objects {
		name := fmt.Sprint(item["name"])
		if nameRegex != nil && !nameRegex.MatchString(name) {
			continue
		}
		variableMap := map[string]interface{}{
			"id":          fmt.Sprintf("%v:%v:%v", workspace, namespace, name),
			"workspace":   workspace,
			"namespace":   namespace,
			"name":        name,
			"kind":        item["kind"],
			"value":       item["value"],
			"description": item["description"],
		}
		filteredObjects = append(filteredObjects, variableMap)
		ids = append(ids, variableMap["id"].(string))
	}

	d.SetId(fmt.Sprintf("%v:%v", workspace, namespace))
	d.Set("variables", filteredObjects)
	d.Set("ids", ids)

	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		if err := writeToFile(output.(string), filteredObjects); err != nil {
			return WrapError(err)
		}
	}

	return nil
}
