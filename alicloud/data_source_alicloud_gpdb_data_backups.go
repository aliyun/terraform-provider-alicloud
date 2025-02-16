package alicloud

import (
	"fmt"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceAliCloudGpdbDataBackups() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAliCloudGpdbDataBackupRead,
		Schema: map[string]*schema.Schema{
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"backup_mode": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"db_instance_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"data_backup_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"data_type": {
				Type:     schema.TypeString,
				Optional: true,
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
			"status": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"backups": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"backup_end_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"backup_end_time_local": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"backup_method": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"backup_mode": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"backup_set_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"backup_size": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"backup_start_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"backup_start_time_local": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"bakset_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"consistent_time": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"db_instance_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"data_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
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

func dataSourceAliCloudGpdbDataBackupRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	action := "DescribeDataBackups"
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

		if v, ok := d.GetOk("backup_mode"); ok {
			request["BackupMode"] = v
		}
		request["DBInstanceId"] = d.Get("db_instance_id")
		if v, ok := d.GetOk("data_backup_id"); ok {
			request["BackupId"] = v
		}
		if v, ok := d.GetOk("data_type"); ok {
			request["DataType"] = v
		}
		request["EndTime"] = d.Get("end_time")
		request["StartTime"] = d.Get("start_time")
		if v, ok := d.GetOk("status"); ok {
			request["BackupStatus"] = v
		}
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
			if len(idsMap) > 0 {
				if _, ok := idsMap[fmt.Sprint(item["DBInstanceId"])]; !ok {
					continue
				}
			}
			objects = append(objects, item)
		}

		if len(result) < request["PageSize"].(int) {
			break
		}
		request["PageNumber"] = request["PageNumber"].(int) + 1

	}

	ids := make([]string, 0)
	s := make([]map[string]interface{}, 0)
	for _, objectRaw := range objects {
		mapping := map[string]interface{}{}

		mapping["status"] = objectRaw["BackupStatus"]
		mapping["backup_end_time"] = objectRaw["BackupEndTime"]
		mapping["backup_end_time_local"] = objectRaw["BackupEndTimeLocal"]
		mapping["backup_method"] = objectRaw["BackupMethod"]
		mapping["backup_mode"] = objectRaw["BackupMode"]
		mapping["backup_set_id"] = objectRaw["BackupSetId"]
		mapping["backup_size"] = objectRaw["BackupSize"]
		mapping["backup_start_time"] = objectRaw["BackupStartTime"]
		mapping["backup_start_time_local"] = objectRaw["BackupStartTimeLocal"]
		mapping["bakset_name"] = objectRaw["BaksetName"]
		mapping["consistent_time"] = objectRaw["ConsistentTime"]
		mapping["db_instance_id"] = objectRaw["DBInstanceId"]
		mapping["data_type"] = objectRaw["DataType"]

		ids = append(ids, fmt.Sprint(mapping["db_instance_id"]))
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}
	if err := d.Set("backups", s); err != nil {
		return WrapError(err)
	}

	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}
	return nil
}
