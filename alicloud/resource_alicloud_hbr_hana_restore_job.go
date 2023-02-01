package alicloud

import (
	"fmt"
	"github.com/PaesslerAG/jsonpath"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"log"
	"time"
)

func resourceAlicloudHbrHanaRestoreJob() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudHbrHanaRestoreJobCreate,
		Read:   resourceAlicloudHbrHanaRestoreJobRead,
		Update: resourceAlicloudHbrHanaRestoreJobUpdate,
		Delete: resourceAlicloudHbrHanaRestoreJobDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"backup_id": {
				Optional: true,
				ForceNew: true,
				Type:     schema.TypeString,
			},
			"backup_prefix": {
				Required: true,
				ForceNew: true,
				Type:     schema.TypeString,
			},
			"check_access": {
				Optional: true,
				ForceNew: true,
				Type:     schema.TypeBool,
			},
			"clear_log": {
				Optional: true,
				ForceNew: true,
				Type:     schema.TypeBool,
			},
			"cluster_id": {
				Required: true,
				ForceNew: true,
				Type:     schema.TypeString,
			},
			"current_phase": {
				Computed: true,
				Type:     schema.TypeInt,
			},
			"current_progress": {
				Computed: true,
				Type:     schema.TypeString,
			},
			"database_name": {
				Required: true,
				ForceNew: true,
				Type:     schema.TypeString,
			},
			"database_restore_id": {
				Computed: true,
				Type:     schema.TypeString,
			},
			"end_time": {
				Computed: true,
				Type:     schema.TypeString,
			},
			"log_position": {
				Optional: true,
				ForceNew: true,
				Type:     schema.TypeString,
			},
			"max_phase": {
				Computed: true,
				Type:     schema.TypeInt,
			},
			"max_progress": {
				Computed: true,
				Type:     schema.TypeString,
			},
			"message": {
				Computed: true,
				ForceNew: true,
				Type:     schema.TypeString,
			},
			"mode": {
				Required: true,
				ForceNew: true,
				Type:     schema.TypeString,
			},
			"phase": {
				Computed: true,
				Type:     schema.TypeString,
			},
			"reached_time": {
				Computed: true,
				Type:     schema.TypeString,
			},
			"recovery_point_in_time": {
				Optional: true,
				ForceNew: true,
				Type:     schema.TypeString,
			},
			"restore_id": {
				Computed: true,
				Type:     schema.TypeString,
			},
			"sid_admin": {
				Optional: true,
				Type:     schema.TypeString,
			},
			"source": {
				Optional: true,
				ForceNew: true,
				Type:     schema.TypeString,
			},
			"source_cluster_id": {
				Optional: true,
				ForceNew: true,
				Type:     schema.TypeString,
			},
			"start_time": {
				Computed: true,
				Type:     schema.TypeString,
			},
			"state": {
				Computed: true,
				Type:     schema.TypeString,
			},
			"status": {
				Optional: true,
				Type:     schema.TypeString,
			},
			"system_copy": {
				Optional: true,
				ForceNew: true,
				Type:     schema.TypeBool,
			},
			"token": {
				Optional: true,
				Type:     schema.TypeString,
			},
			"use_catalog": {
				Optional: true,
				ForceNew: true,
				Type:     schema.TypeBool,
			},
			"use_delta": {
				Optional: true,
				ForceNew: true,
				Type:     schema.TypeBool,
			},
			"vault_id": {
				Optional: true,
				ForceNew: true,
				Type:     schema.TypeString,
			},
			"volume_id": {
				Optional: true,
				ForceNew: true,
				Type:     schema.TypeInt,
			},
		},
	}
}

func resourceAlicloudHbrHanaRestoreJobCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	request := make(map[string]interface{})
	conn, err := client.NewHbrClient()
	if err != nil {
		return WrapError(err)
	}

	if v, ok := d.GetOk("backup_id"); ok {
		request["BackupId"] = v
	}
	if v, ok := d.GetOk("backup_prefix"); ok {
		request["BackupPrefix"] = v
	}
	if v, ok := d.GetOk("check_access"); ok {
		request["CheckAccess"] = v
	}
	if v, ok := d.GetOk("clear_log"); ok {
		request["ClearLog"] = v
	}
	if v, ok := d.GetOk("cluster_id"); ok {
		request["ClusterId"] = v
	}
	if v, ok := d.GetOk("database_name"); ok {
		request["DatabaseName"] = v
	}
	if v, ok := d.GetOk("log_position"); ok {
		request["LogPosition"] = v
	}
	if v, ok := d.GetOk("mode"); ok {
		request["Mode"] = v
	}
	if v, ok := d.GetOk("recovery_point_in_time"); ok {
		request["RecoveryPointInTime"] = v
	}
	if v, ok := d.GetOk("sid_admin"); ok {
		request["SidAdmin"] = v
	}
	if v, ok := d.GetOk("source"); ok {
		request["Source"] = v
	}
	if v, ok := d.GetOk("source_cluster_id"); ok {
		request["SourceClusterId"] = v
	}
	if v, ok := d.GetOk("system_copy"); ok {
		request["SystemCopy"] = v
	}
	if v, ok := d.GetOk("use_catalog"); ok {
		request["UseCatalog"] = v
	}
	if v, ok := d.GetOk("use_delta"); ok {
		request["UseDelta"] = v
	}
	if v, ok := d.GetOk("volume_id"); ok {
		request["VolumeId"] = v
	}

	var response map[string]interface{}
	action := "CreateHanaRestore"
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutCreate)), func() *resource.RetryError {
		resp, err := conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2017-09-08"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		response = resp
		addDebug(action, response, request)
		return nil
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_hbr_hana_restore_job", action, AlibabaCloudSdkGoERROR)
	}

	if v, err := jsonpath.Get("$.RestoreId", response); err != nil || v == nil {
		return WrapErrorf(err, IdMsg, "alicloud_hbr_hana_restore_job")
	} else {
		d.SetId(fmt.Sprint(v))
	}

	return resourceAlicloudHbrHanaRestoreJobUpdate(d, meta)
}

func resourceAlicloudHbrHanaRestoreJobRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	hbrService := HbrService{client}

	object, err := hbrService.DescribeHbrHanaRestoreJob(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_hbr_hana_restore_job hbrService.DescribeHbrHanaRestoreJob Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	d.Set("backup_prefix", object["HanaRestore.HanaRestores[*].BackupPrefix"])
	d.Set("check_access", object["HanaRestore.HanaRestores[*].CheckAccess"])
	d.Set("clear_log", object["HanaRestore.HanaRestores[*].ClearLog"])
	d.Set("cluster_id", object["HanaRestore.HanaRestores[*].ClusterId"])
	d.Set("current_phase", object["HanaRestore.HanaRestores[*].CurrentPhase"])
	d.Set("current_progress", object["HanaRestore.HanaRestores[*].CurrentProgress"])
	d.Set("database_name", object["HanaRestore.HanaRestores[*].DatabaseName"])
	d.Set("database_restore_id", object["HanaRestore.HanaRestores[*].DatabaseRestoreId"])
	d.Set("end_time", object["HanaRestore.HanaRestores[*].EndTime"])
	d.Set("log_position", object["HanaRestore.HanaRestores[*].LogPosition"])
	d.Set("max_phase", object["HanaRestore.HanaRestores[*].MaxPhase"])
	d.Set("max_progress", object["HanaRestore.HanaRestores[*].MaxProgress"])
	d.Set("message", object["Message"])
	d.Set("mode", object["HanaRestore.HanaRestores[*].Mode"])
	d.Set("phase", object["HanaRestore.HanaRestores[*].Phase"])
	d.Set("reached_time", object["HanaRestore.HanaRestores[*].ReachedTime"])
	d.Set("recovery_point_in_time", object["HanaRestore.HanaRestores[*].RecoveryPointInTime"])
	d.Set("source", object["HanaRestore.HanaRestores[*].Source"])
	d.Set("source_cluster_id", object["HanaRestore.HanaRestores[*].SourceClusterId"])
	d.Set("start_time", object["HanaRestore.HanaRestores[*].StartTime"])
	d.Set("state", object["HanaRestore.HanaRestores[*].State"])
	d.Set("status", object["HanaRestore.HanaRestores[*].Status"])
	d.Set("system_copy", object["HanaRestore.HanaRestores[*].SystemCopy"])
	d.Set("use_catalog", object["HanaRestore.HanaRestores[*].UseCatalog"])
	d.Set("use_delta", object["HanaRestore.HanaRestores[*].UseDelta"])
	d.Set("vault_id", object["HanaRestore.HanaRestores[*].VaultId"])
	d.Set("volume_id", object["HanaRestore.HanaRestores[*].VolumeId"])

	return nil
}

func resourceAlicloudHbrHanaRestoreJobUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	hbrService := HbrService{client}
	conn, err := client.NewHbrClient()
	if err != nil {
		return WrapError(err)
	}
	update := false
	request := map[string]interface{}{}

	if !d.IsNewResource() && d.HasChange("cluster_id") {
		update = true
	}
	request["ClusterId"] = d.Get("cluster_id")
	if !d.IsNewResource() && d.HasChange("database_name") {
		update = true
	}
	request["DatabaseName"] = d.Get("database_name")
	if v, ok := d.GetOk("token"); ok {
		request["Token"] = v
	}
	if d.HasChange("vault_id") {
		update = true
		if v, ok := d.GetOk("vault_id"); ok {
			request["VaultId"] = v
		}
	}

	if update {
		action := "CancelHanaRestore"
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *resource.RetryError {
			resp, err := conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2017-09-08"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			addDebug(action, resp, request)
			return nil
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
	}

	return resourceAlicloudHbrHanaRestoreJobRead(d, meta)
}

func resourceAlicloudHbrHanaRestoreJobDelete(d *schema.ResourceData, meta interface{}) error {
	return nil
}
