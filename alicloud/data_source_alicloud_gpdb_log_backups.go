package alicloud

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceAliCloudGpdbLogbackups() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAliCloudGpdbLogbackupRead,
		Schema: map[string]*schema.Schema{
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"db_instance_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"end_time": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"page_number": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"page_size": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"start_time": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"logbackups": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"db_instance_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"log_backup_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"log_file_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"log_file_size": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"log_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"record_total": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"segment_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func dataSourceAliCloudGpdbLogbackupRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	action := "DescribeLogBackups"
	var err error

	var objects []map[string]interface{}
	for {
		request = make(map[string]interface{})
		query = make(map[string]interface{})
		if v, ok := d.GetOk("page_number"); ok && v.(int) > 0 {
			request["PageNumber"] = v.(int)
		} else {
			request["PageNumber"] = 1
		}
		if v, ok := d.GetOk("page_size"); ok && v.(int) > 0 {
			request["PageSize"] = v.(int)
		} else {
			request["PageSize"] = PageSizeLarge
		}
		request["DBInstanceId"] = d.Get("db_instance_id")
		request["EndTime"] = d.Get("end_time")
		request["StartTime"] = d.Get("start_time")
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("gpdb", "2016-05-03", action, query, request, true)

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

		resp, _ := jsonpath.Get("$.Items[*]", response)

		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})
			if nameRegex != nil && !nameRegex.MatchString(fmt.Sprint(item["Name"])) {
				continue
			}
			found := false
			for id, _ := range idsMap {
				if strings.Contains(fmt.Sprint(item["DBInstanceId"]), id) {
					found = true
					break
				}
			}
			if !found {
				continue
			}

			objects = append(objects, item)
		}

		if len(result) < request["PageSize"].(int) {
			break
		}
		request["PageNumber"] = request["PageNumber"].(int) + 1
	}

	ids := make([]string, 0)
	names := make([]interface{}, 0)
	s := make([]map[string]interface{}, 0)
	for _, objectRaw := range objects {
		mapping := map[string]interface{}{}

		mapping["db_instance_id"] = objectRaw["DBInstanceId"]
		mapping["log_backup_id"] = objectRaw["BackupId"]
		mapping["log_file_name"] = objectRaw["LogFileName"]
		mapping["log_file_size"] = objectRaw["LogFileSize"]
		mapping["log_time"] = objectRaw["LogTime"]
		mapping["segment_name"] = objectRaw["SegmentName"]

		ids = append(ids, fmt.Sprint(mapping["id"]))
		names = append(names, objectRaw["AlertName"])
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}

	if err := d.Set("logbackups", s); err != nil {
		return WrapError(err)
	}

	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}
	return nil
}
