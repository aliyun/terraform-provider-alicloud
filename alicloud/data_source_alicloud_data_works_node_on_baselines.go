package alicloud

import (
	"fmt"
	"regexp"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceAlicloudDataWorksNodeOnBaselines() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudDataWorksNodeOnBaselinesRead,
		Schema: map[string]*schema.Schema{
			"baseline_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"node_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"owner": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"project_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"name_regex": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"nodes": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"baseline_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"node_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"node_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"owner": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"project_id": {
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

func dataSourceAlicloudDataWorksNodeOnBaselinesRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	action := "ListNodesByBaseline"
	request := map[string]interface{}{
		"BaselineId": d.Get("baseline_id").(string),
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

	var (
		nodeIDFilter    = d.Get("node_id").(string)
		ownerFilter     = d.Get("owner").(string)
		projectIDFilter = d.Get("project_id").(string)
		nameRegexStr    = d.Get("name_regex").(string)
		nameRegex       *regexp.Regexp
	)
	if nameRegexStr != "" {
		if r, e := regexp.Compile(nameRegexStr); e == nil {
			nameRegex = r
		}
	}

	var response map[string]interface{}
	var err error
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = client.RpcPost("dataworks-public", "2020-05-18", action, nil, request, true)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)
	if err != nil {
		return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_data_works_node_on_baselines", action, AlibabaCloudSdkGoERROR)
	}

	resp, err := jsonpath.Get("$.Data", response)
	if err != nil {
		return WrapErrorf(err, FailedGetAttributeMsg, action, "$.Data", response)
	}
	result, _ := resp.([]interface{})

	var objects []map[string]interface{}
	ids := make([]string, 0)
	for _, v := range result {
		item, ok := v.(map[string]interface{})
		if !ok {
			continue
		}
		nodeID := fmt.Sprint(item["NodeId"])
		if nodeIDFilter != "" && nodeID != nodeIDFilter {
			continue
		}
		if ownerFilter != "" && fmt.Sprint(item["Owner"]) != ownerFilter {
			continue
		}
		if projectIDFilter != "" && fmt.Sprint(item["ProjectId"]) != projectIDFilter {
			continue
		}
		if nameRegex != nil {
			nodeName := fmt.Sprint(item["NodeName"])
			if !nameRegex.MatchString(nodeName) {
				continue
			}
		}
		if len(idsMap) > 0 {
			if _, ok := idsMap[nodeID]; !ok {
				continue
			}
		}
		mapping := map[string]interface{}{
			"id":          fmt.Sprintf("%s:%s", d.Get("baseline_id").(string), nodeID),
			"baseline_id": d.Get("baseline_id").(string),
			"node_id":     nodeID,
			"node_name":   item["NodeName"],
			"owner":       item["Owner"],
			"project_id":  item["ProjectId"],
			"region_id":   client.RegionId,
		}
		ids = append(ids, mapping["id"].(string))
		objects = append(objects, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}
	if err := d.Set("nodes", objects); err != nil {
		return WrapError(err)
	}

	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), objects)
	}

	return nil
}
