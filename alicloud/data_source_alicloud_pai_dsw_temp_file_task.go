// Package alicloud. This file is hand-written for PaiDsw TempFileTask data source.
package alicloud

import (
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceAliCloudPaiDswTempFileTask() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAliCloudPaiDswTempFileTaskRead,
		Schema: map[string]*schema.Schema{
			"temp_file_task_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"instance_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"owner_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"user_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"gmt_modified_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"gmt_expired_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"region_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func dataSourceAliCloudPaiDswTempFileTaskRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	paiDswService := PaiDswService{client}

	tempFileTaskId := d.Get("temp_file_task_id").(string)
	objectRaw, err := paiDswService.DescribePaiDswTempFileTask(tempFileTaskId)
	if err != nil {
		return WrapError(err)
	}

	d.SetId(tempFileTaskId)

	if objectRaw["InstanceId"] != nil {
		d.Set("instance_id", objectRaw["InstanceId"])
	}
	if objectRaw["OwnerId"] != nil {
		d.Set("owner_id", objectRaw["OwnerId"])
	}
	if objectRaw["UserId"] != nil {
		d.Set("user_id", objectRaw["UserId"])
	}
	if objectRaw["GmtCreateTime"] != nil {
		d.Set("create_time", objectRaw["GmtCreateTime"])
	}
	if objectRaw["GmtModifiedTime"] != nil {
		d.Set("gmt_modified_time", objectRaw["GmtModifiedTime"])
	}
	if objectRaw["GmtExpiredTime"] != nil {
		d.Set("gmt_expired_time", objectRaw["GmtExpiredTime"])
	}
	d.Set("region_id", client.RegionId)

	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), map[string]interface{}{
			"temp_file_task_id": tempFileTaskId,
			"instance_id":       objectRaw["InstanceId"],
			"owner_id":          objectRaw["OwnerId"],
			"user_id":           objectRaw["UserId"],
			"create_time":       objectRaw["GmtCreateTime"],
			"gmt_modified_time": objectRaw["GmtModifiedTime"],
			"gmt_expired_time":  objectRaw["GmtExpiredTime"],
			"region_id":         client.RegionId,
		})
	}

	return nil
}
