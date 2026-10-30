package hbr

import (
	"fmt"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceOtsSnapshots() *schema.Resource {
	return &schema.Resource{
		Read: DataSourceOtsSnapshotsRead,
		Schema: map[string]*schema.Schema{
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"start_time": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"end_time": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"snapshots": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"snapshot_hash": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"vault_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"backup_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"create_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"actual_bytes": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"source_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"bytes_total": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"complete_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"retention": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"created_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"parent_snapshot_hash": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"start_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"updated_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"snapshot_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"job_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"instance_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"table_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"range_start": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"range_end": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func DataSourceOtsSnapshotsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	action := "DescribeOtsTableSnapshots"
	request := make(map[string]interface{})

	if v, ok := d.GetOk("start_time"); ok {
		request["StartTime"] = v
	}
	if v, ok := d.GetOk("end_time"); ok {
		request["EndTime"] = v
	}

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
		wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
		err = retry.Retry(5*time.Minute, func() *retry.RetryError {
			response, err = client.RpcGet("hbr", "2017-09-08", action, request, nil)
			if err != nil {
				if helper.NeedRetry(err) {
					wait()
					return retry.RetryableError(err)
				}
				return retry.NonRetryableError(err)
			}
			return nil
		})
		logs.AddDebug(action, response, request)
		if err != nil {
			return errs.WrapErrorf(err, errs.DataDefaultErrorMsg, "alicloud_hbr_ots_backup_snapshots", action, errs.AlibabaCloudSdkGoERROR)
		}
		if fmt.Sprint(response["Success"]) == "false" {
			return errs.WrapError(fmt.Errorf("%s failed, response: %v", action, response))
		}
		resp, err := jsonpath.Get("$.Snapshots", response)
		if err != nil {
			return errs.WrapErrorf(err, errs.FailedGetAttributeMsg, action, "$.Snapshots", response)
		}
		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})
			if len(idsMap) > 0 {
				if _, ok := idsMap[fmt.Sprint(item["SnapshotId"])]; !ok {
					continue
				}
			}
			objects = append(objects, item)
		}
		if nextToken, ok := response["NextToken"].(string); ok && nextToken != "" {
			request["NextToken"] = nextToken
		} else {
			break
		}
	}
	ids := make([]string, 0)
	s := make([]map[string]interface{}, 0)
	for _, object := range objects {
		mapping := map[string]interface{}{
			"status":               object["Status"],
			"snapshot_hash":        object["SnapshotHash"],
			"vault_id":             object["VaultId"],
			"backup_type":          object["BackupType"],
			"create_time":          fmt.Sprint(object["CreateTime"]),
			"actual_bytes":         fmt.Sprint(object["ActualBytes"]),
			"source_type":          object["SourceType"],
			"bytes_total":          fmt.Sprint(object["BytesTotal"]),
			"complete_time":        fmt.Sprint(object["CompleteTime"]),
			"retention":            fmt.Sprint(object["Retention"]),
			"created_time":         fmt.Sprint(object["CreatedTime"]),
			"parent_snapshot_hash": object["ParentSnapshotHash"],
			"start_time":           fmt.Sprint(object["StartTime"]),
			"updated_time":         fmt.Sprint(object["UpdatedTime"]),
			"id":                   fmt.Sprint(object["SnapshotId"]),
			"snapshot_id":          fmt.Sprint(object["SnapshotId"]),
			"job_id":               object["JobId"],
			"instance_name":        object["InstanceName"],
			"table_name":           fmt.Sprint(object["TableName"]),
			"range_start":          fmt.Sprint(object["RangeStart"]),
			"range_end":            fmt.Sprint(object["RangeEnd"]),
		}
		ids = append(ids, fmt.Sprint(mapping["id"]))
		s = append(s, mapping)
	}

	d.SetId(helper.DataResourceIdHash(ids))
	if err := d.Set("ids", ids); err != nil {
		return errs.WrapError(err)
	}

	if err := d.Set("snapshots", s); err != nil {
		return errs.WrapError(err)
	}
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		helper.WriteToFile(output.(string), s)
	}
	return nil
}
