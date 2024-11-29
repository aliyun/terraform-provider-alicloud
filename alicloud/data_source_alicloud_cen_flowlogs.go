// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"regexp"
	"time"

	"github.com/PaesslerAG/jsonpath"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func dataSourceAliCloudCenFlowLogs() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAliCloudCenFlowLogRead,
		Schema: map[string]*schema.Schema{
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"name_regex": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.ValidateRegexp,
				ForceNew:     true,
			},
			"names": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"cen_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"flow_log_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"flow_log_name": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"flow_log_version": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"interval": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
			},
			"log_store_name": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"page_number": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
			},
			"page_size": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
			},
			"project_name": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"region_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"status": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"transit_router_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"flowlogs": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"cen_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"create_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"flow_log_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"flow_log_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"flow_log_version": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"interval": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"log_format_string": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"log_store_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"project_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"record_total": {
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
						"tags": {
							Type:     schema.TypeMap,
							Computed: true,
						},
						"transit_router_attachment_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"transit_router_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
		},
	}
}

func dataSourceAliCloudCenFlowLogRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	action := "DescribeFlowlogs"
	conn, err := client.NewCenClient()
	if err != nil {
		return WrapError(err)
	}
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["RegionId"] = client.RegionId
	request["ClientToken"] = buildClientToken(action)
	request["CenId"] = d.Get("cen_id")
	if v, ok := d.GetOk("description"); ok {
		request["Description"] = v
	}
	if v, ok := d.GetOk("flow_log_id"); ok {
		request["FlowLogId"] = v
	}
	if v, ok := d.GetOk("flow_log_name"); ok {
		request["FlowLogName"] = v
	}
	if v, ok := d.GetOk("flow_log_version"); ok {
		request["FlowLogVersion"] = v
	}
	if v, ok := d.GetOkExists("interval"); ok {
		request["Interval"] = v
	}
	request["LogStoreName"] = d.Get("log_store_name")
	request["ProjectName"] = d.Get("project_name")
	if v, ok := d.GetOk("status"); ok {
		request["Status"] = v
	}
	if v, ok := d.GetOk("transit_router_id"); ok {
		request["TransitRouterId"] = v
	}
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2017-09-12"), StringPointer("AK"), query, request, &runtime)
		request["ClientToken"] = buildClientToken(action)

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

	var objects []map[string]interface{}
	var nameRegex *regexp.Regexp
	if v, ok := d.GetOk("name_regex"); ok {
		r, err := regexp.Compile(v.(string))
		if err != nil {
			return WrapError(err)
		}
		nameRegex = r
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

	resp, _ := jsonpath.Get("$.FlowLogs.FlowLog[*]", response)

	result, _ := resp.([]interface{})
	for _, v := range result {
		item := v.(map[string]interface{})
		if nameRegex != nil && !nameRegex.MatchString(fmt.Sprint(item["FlowLogName"])) {
			continue
		}
		if len(idsMap) > 0 {
			if _, ok := idsMap[fmt.Sprint(item["FlowLogId"])]; !ok {
				continue
			}
		}
		objects = append(objects, item)
	}

	ids := make([]string, 0)
	names := make([]interface{}, 0)
	s := make([]map[string]interface{}, 0)
	for _, objectRaw := range objects {
		mapping := map[string]interface{}{}

		mapping["id"] = objectRaw["FlowLogId"]
		mapping["cen_id"] = objectRaw["CenId"]
		mapping["create_time"] = objectRaw["CreationTime"]
		mapping["description"] = objectRaw["Description"]
		mapping["flow_log_id"] = objectRaw["FlowLogId"]
		mapping["flow_log_name"] = objectRaw["FlowLogName"]
		mapping["interval"] = objectRaw["Interval"]
		mapping["log_format_string"] = objectRaw["LogFormatString"]
		mapping["log_store_name"] = objectRaw["LogStoreName"]
		mapping["project_name"] = objectRaw["ProjectName"]
		mapping["region_id"] = objectRaw["RegionId"]
		mapping["status"] = objectRaw["Status"]
		mapping["transit_router_attachment_id"] = objectRaw["TransitRouterAttachmentId"]
		mapping["transit_router_id"] = objectRaw["TransitRouterId"]

		tagsMaps, _ := jsonpath.Get("$.Tags.Tag", objectRaw)
		mapping["tags"] = tagsToMap(tagsMaps)

		ids = append(ids, fmt.Sprint(mapping["id"]))
		names = append(names, objectRaw["FlowLogName"])
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}

	if err := d.Set("names", names); err != nil {
		return WrapError(err)
	}
	if err := d.Set("flowlogs", s); err != nil {
		return WrapError(err)
	}

	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}
	return nil
}
