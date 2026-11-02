package rkvstore

import (
	"fmt"
	"log"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/convert"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceAlicloudKvstoreAuditLogConfig() *schema.Resource {
	return &schema.Resource{
		Create: ResourceAlicloudKvstoreAuditLogConfigCreate,
		Read:   ResourceAlicloudKvstoreAuditLogConfigRead,
		Update: ResourceAlicloudKvstoreAuditLogConfigUpdate,
		Delete: ResourceAlicloudKvstoreAuditLogConfigDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(1 * time.Minute),
			Update: schema.DefaultTimeout(1 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"db_audit": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"retention": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func ResourceAlicloudKvstoreAuditLogConfigCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	var err error
	action := "ModifyAuditLogConfig"
	request := make(map[string]interface{})
	if v, ok := d.GetOkExists("db_audit"); ok {
		request["DbAudit"] = v
	}
	request["InstanceId"] = d.Get("instance_id")
	if v, ok := d.GetOk("retention"); ok {
		request["Retention"] = v
	} else if v, ok := d.GetOkExists("db_audit"); ok && v.(bool) {
		return errs.WrapError(fmt.Errorf("attribute '%s' is required when '%s' is %v ", "retention", "db_audit", d.Get("db_audit")))
	}
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
		response, err = client.RpcPost("R-kvstore", "2015-01-01", action, nil, request, false)
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
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, "alicloud_kvstore_audit_log_config", action, errs.AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(request["InstanceId"]))
	rKvstoreService := RKvstoreService{client}
	stateConf := helper.BuildStateConf([]string{}, []string{"Normal"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, rKvstoreService.KvstoreAuditLogConfigStateRefreshFunc(d.Id(), []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return errs.WrapErrorf(err, errs.IdMsg, d.Id())
	}

	return ResourceAlicloudKvstoreAuditLogConfigRead(d, meta)
}
func ResourceAlicloudKvstoreAuditLogConfigRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	rKvstoreService := RKvstoreService{client}
	object, err := rKvstoreService.DescribeKvstoreAuditLogConfig(d.Id())
	if err != nil {
		if errs.NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_kvstore_audit_log_config rKvstoreService.DescribeKvstoreAuditLogConfig Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return errs.WrapError(err)
	}

	d.Set("instance_id", d.Id())
	d.Set("db_audit", object["DbAudit"])
	if v, ok := object["Retention"]; ok && v.(string) != "" {
		retentionVal, err := convert.FormatInt(v)
		if err != nil {
			return err
		}
		d.Set("retention", retentionVal)
	}
	describeInstanceAttributeObject, err := rKvstoreService.DescribeInstanceAttribute(d.Id())
	if err != nil {
		return errs.WrapError(err)
	}
	d.Set("create_time", describeInstanceAttributeObject["CreateTime"])
	d.Set("status", describeInstanceAttributeObject["InstanceStatus"])
	return nil
}
func ResourceAlicloudKvstoreAuditLogConfigUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	rKvstoreService := RKvstoreService{client}
	var response map[string]interface{}
	var err error
	update := false
	request := map[string]interface{}{
		"InstanceId": d.Id(),
	}
	if d.HasChange("db_audit") || d.IsNewResource() {
		update = true
		if v, ok := d.GetOkExists("db_audit"); ok {
			request["DbAudit"] = v
		}
	}
	if d.HasChange("retention") {
		update = true
		if v, ok := d.GetOk("retention"); ok {
			request["Retention"] = v
		} else if v, ok := d.GetOkExists("db_audit"); ok && v.(bool) {
			return errs.WrapError(fmt.Errorf("attribute '%s' is required when '%s' is %v ", "retention", "db_audit", d.Get("db_audit")))
		}
	}
	if update {
		action := "ModifyAuditLogConfig"
		wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
		err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
			response, err = client.RpcPost("R-kvstore", "2015-01-01", action, nil, request, false)
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
			return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
		}
		stateConf := helper.BuildStateConf([]string{}, []string{"Normal"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, rKvstoreService.KvstoreAuditLogConfigStateRefreshFunc(d.Id(), []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return errs.WrapErrorf(err, errs.IdMsg, d.Id())
		}
	}
	return ResourceAlicloudKvstoreAuditLogConfigRead(d, meta)
}
func ResourceAlicloudKvstoreAuditLogConfigDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[WARN] Cannot destroy ResourceAlicloudKvstoreAuditLogConfig. Terraform will remove this resource from the state file, however resources may remain.")
	return nil
}
