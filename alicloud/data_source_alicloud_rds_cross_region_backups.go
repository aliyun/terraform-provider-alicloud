package alicloud

import (
	"fmt"
	"time"

	"github.com/PaesslerAG/jsonpath"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceAlicloudRdsCrossRegionBackups() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudRdsCrossRegionBackupsRead,
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
				ForceNew: true,
			},
			"cross_backup_region": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"end_time": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"start_time": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"cross_backup_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"backup_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"resource_group_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"backups": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"consistent_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"db_instance_storage_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"cross_backup_id": {
							Type:     schema.TypeString,
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
						"cross_backup_set_location": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"instance_id": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"cross_backup_download_link": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"engine_version": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"cross_backup_set_file": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"backup_set_scale": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"cross_backup_set_size": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"backup_set_status": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"cross_backup_region": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"category": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"engine": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"backup_method": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"restore_regions": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"recovery_begin_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"recovery_end_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceAlicloudRdsCrossRegionBackupsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	action := "DescribeCrossRegionBackups"
	request := make(map[string]interface{})
	request["SourceIp"] = client.SourceIp
	request["RegionId"] = client.RegionId
	request["DBInstanceId"] = d.Get("db_instance_id")
	if v, ok := d.GetOk("cross_backup_region"); ok {
		request["CrossBackupRegion"] = v
	}
	if v, ok := d.GetOk("cross_backup_id"); ok {
		request["CrossBackupId"] = v
	}
	if v, ok := d.GetOk("start_time"); ok {
		request["StartTime"] = v
	}
	if v, ok := d.GetOk("end_time"); ok {
		request["EndTime"] = v
	}
	if v, ok := d.GetOk("backup_id"); ok {
		request["BackupId"] = v
	}
	if v, ok := d.GetOk("resource__group_id"); ok {
		request["ResourceGroupId"] = v
	}
	request["PageSize"] = PageSizeLarge
	request["PageNumber"] = 1
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
	var response map[string]interface{}
	var err error
	for {
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			response, err = client.RpcPost("Rds", "2014-08-15", action, nil, request, true)
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
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_rds_cross_backups", action, AlibabaCloudSdkGoERROR)
		}
		resp, err := jsonpath.Get("$.Items.Item", response)
		if err != nil {
			return WrapErrorf(err, FailedGetAttributeMsg, action, "$.Items.Item", response)
		}
		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})
			if len(idsMap) > 0 {
				if _, ok := idsMap[fmt.Sprint(item["CrossBackupId"])]; !ok {
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
	for _, object := range objects {
		var response map[string]interface{}
		if object["Engine"].(string) != "mssql" {
			action := "DescribeAvailableRecoveryTime"
			request := map[string]interface{}{
				"RegionId":      client.RegionId,
				"SourceIp":      client.SourceIp,
				"CrossBackupId": formatInt(object["CrossBackupId"]),
			}
			runtime := util.RuntimeOptions{}
			runtime.SetAutoretry(true)
			wait := incrementalWait(3*time.Second, 3*time.Second)
			err = resource.Retry(5*time.Minute, func() *resource.RetryError {
				response, err = client.RpcPost("Rds", "2014-08-15", action, nil, request, true)
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
				return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_rds_cross_backups", action, AlibabaCloudSdkGoERROR)
			}
		}
		mapping := map[string]interface{}{
			"id":                         fmt.Sprint(object["CrossBackupId"]),
			"cross_backup_id":            fmt.Sprint(object["CrossBackupId"]),
			"db_instance_storage_type":   object["DBInstanceStorageType"],
			"cross_backup_set_location":  object["CrossBackupSetLocation"],
			"instance_id":                formatInt(object["InstanceId"]),
			"cross_backup_download_link": object["CrossBackupDownloadLink"],
			"engine_version":             fmt.Sprint(object["EngineVersion"]),
			"backup_start_time":          object["BackupStartTime"],
			"backup_end_time":            object["BackupEndTime"],
			"backup_type":                object["BackupType"],
			"consistent_time":            fmt.Sprint(object["ConsistentTime"]),
			"cross_backup_set_file":      object["CrossBackupSetFile"],
			"backup_set_scale":           object["BackupSetScale"],
			"backup_set_status":          formatInt(object["BackupSetStatus"]),
			"cross_backup_set_size":      formatInt(object["CrossBackupSetSize"]),
			"cross_backup_region":        object["CrossBackupRegion"],
			"category":                   object["Category"],
			"engine":                     object["Engine"],
			"backup_method":              object["BackupMethod"],
			"restore_regions":            object["RestoreRegions"].(map[string]interface{})["RestoreRegion"],
			"recovery_begin_time":        response["RecoveryBeginTime"],
			"recovery_end_time":          response["RecoveryEndTime"],
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
