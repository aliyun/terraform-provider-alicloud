// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceAliCloudMongodbBackups() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAliCloudMongodbBackupRead,
		Schema: map[string]*schema.Schema{
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"backup_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"db_instance_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"end_time": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"start_time": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"backups": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"backup_db_names": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"backup_download_url": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"backup_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"backup_intranet_download_url": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"backup_job_id": {
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
						"backup_size": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"backup_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"backup_start_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"backup_end_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"id": {
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
			"enable_details": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
		},
	}
}

func dataSourceAliCloudMongodbBackupRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	var objects []map[string]interface{}

	idsMap := make(map[string]string)
	if v, ok := d.GetOk("ids"); ok {
		for _, vv := range v.([]interface{}) {
			if vv == nil {
				continue
			}
			idsMap[vv.(string)] = vv.(string)
		}
	}

	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	action := "DescribeBackups"
	var err error
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["RegionId"] = client.RegionId
	if v, ok := d.GetOk("backup_id"); ok {
		request["BackupId"] = v
	}
	if v, ok := d.GetOk("db_instance_id"); ok {
		request["DBInstanceId"] = v
	}
	if v, ok := d.GetOk("end_time"); ok {
		request["EndTime"] = v
	}
	if v, ok := d.GetOk("start_time"); ok {
		request["StartTime"] = v
	}
	request["PageSize"] = PageSizeLarge
	request["PageNumber"] = 1
	for {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutRead), func() *resource.RetryError {
			response, err = client.RpcPost("Dds", "2015-12-01", action, query, request, true)

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

		resp, _ := jsonpath.Get("$.Backups.Backup[*]", response)

		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})
			if len(idsMap) > 0 {
				if _, ok := idsMap[fmt.Sprint(request["DBInstanceId"], ":", item["BackupId"])]; !ok {
					continue
				}
			}
			objects = append(objects, item)
		}

		if len(result) < PageSizeLarge {
			break
		}
		request["PageNumber"] = request["PageNumber"].(int) + 1
	}

	ids := make([]string, 0)
	s := make([]map[string]interface{}, 0)
	for _, objectRaw := range objects {
		mapping := map[string]interface{}{}

		mapping["id"] = fmt.Sprint(request["DBInstanceId"], ":", objectRaw["BackupId"])

		mapping["backup_db_names"] = objectRaw["BackupDBNames"]
		mapping["backup_download_url"] = objectRaw["BackupDownloadURL"]
		mapping["backup_intranet_download_url"] = objectRaw["BackupIntranetDownloadURL"]
		mapping["backup_method"] = objectRaw["BackupMethod"]
		mapping["backup_mode"] = objectRaw["BackupMode"]
		mapping["backup_size"] = objectRaw["BackupSize"]
		mapping["backup_type"] = objectRaw["BackupType"]
		mapping["backup_start_time"] = objectRaw["BackupStartTime"]
		mapping["backup_end_time"] = objectRaw["BackupEndTime"]
		mapping["status"] = objectRaw["BackupStatus"]
		mapping["backup_id"] = objectRaw["BackupId"]

		if detailedEnabled := d.Get("enable_details"); !detailedEnabled.(bool) {
			ids = append(ids, fmt.Sprint(mapping["id"]))
			s = append(s, mapping)
			continue
		}

		id := fmt.Sprint(request["DBInstanceId"], ":", objectRaw["BackupId"])
		mapping, err = dataSourceAliCloudMongodbBackupReadDescription(d, id, mapping, meta)
		if err != nil {
			return WrapError(err)
		}

		ids = append(ids, fmt.Sprint(mapping["id"]))
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

func dataSourceAliCloudMongodbBackupReadDescription(d *schema.ResourceData, id string, object map[string]interface{}, meta interface{}) (map[string]interface{}, error) {
	client := meta.(*connectivity.AliyunClient)

	mongodbServiceV2 := MongodbServiceV2{client}
	getResp, err := mongodbServiceV2.DescribeMongodbBackup(id)
	if err != nil {
		return nil, WrapError(err)
	}

	// Merge additional fields from Get API response to mapping
	// Reuse the response mapping template from Resource's read function
	mapping := object
	objectRaw := getResp

	mapping["backup_db_names"] = objectRaw["BackupDBNames"]
	mapping["backup_download_url"] = objectRaw["BackupDownloadURL"]
	mapping["backup_intranet_download_url"] = objectRaw["BackupIntranetDownloadURL"]
	mapping["backup_job_id"] = objectRaw["BackupJobId"]
	mapping["backup_method"] = objectRaw["BackupMethod"]
	mapping["backup_mode"] = objectRaw["BackupMode"]
	mapping["backup_size"] = objectRaw["BackupSize"]
	mapping["backup_type"] = objectRaw["BackupType"]
	mapping["backup_start_time"] = objectRaw["BackupStartTime"]
	mapping["backup_end_time"] = objectRaw["BackupEndTime"]
	mapping["status"] = objectRaw["BackupStatus"]
	mapping["backup_id"] = objectRaw["BackupId"]

	return mapping, nil
}
