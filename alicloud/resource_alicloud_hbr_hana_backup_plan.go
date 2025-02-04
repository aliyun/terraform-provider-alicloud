package alicloud

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAlicloudHbrHanaBackupPlan() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudHbrHanaBackupPlanCreate,
		Read:   resourceAlicloudHbrHanaBackupPlanRead,
		Update: resourceAlicloudHbrHanaBackupPlanUpdate,
		Delete: resourceAlicloudHbrHanaBackupPlanDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(1 * time.Minute),
			Delete: schema.DefaultTimeout(1 * time.Minute),
			Update: schema.DefaultTimeout(1 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"backup_prefix": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"backup_type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"COMPLETE", "INCREMENTAL", "DIFFERENTIAL"}, false),
			},
			"cluster_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"database_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"plan_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"plan_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"resource_group_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"schedule": {
				Type:     schema.TypeString,
				Required: true,
			},
			"vault_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"status": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringInSlice([]string{"Disabled", "Enabled"}, false),
			},
		},
	}
}

func resourceAlicloudHbrHanaBackupPlanCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "CreateHanaBackupPlan"
	request := make(map[string]interface{})
	var err error
	if v, ok := d.GetOk("backup_prefix"); ok {
		request["BackupPrefix"] = v
	}
	request["BackupType"] = d.Get("backup_type")
	request["ClusterId"] = d.Get("cluster_id")
	request["DatabaseName"] = d.Get("database_name")
	request["PlanName"] = d.Get("plan_name")
	if v, ok := d.GetOk("resource_group_id"); ok {
		request["ResourceGroupId"] = v
	}
	request["Schedule"] = d.Get("schedule")
	request["VaultId"] = d.Get("vault_id")
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.RpcPost("hbr", "2017-09-08", action, nil, request, false)
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_hbr_hana_backup_plan", action, AlibabaCloudSdkGoERROR)
	}
	if fmt.Sprint(response["Success"]) == "false" {
		return WrapError(fmt.Errorf("%s failed, response: %v", action, response))
	}

	d.SetId(fmt.Sprint(response["PlanId"], ":", request["VaultId"], ":", request["ClusterId"]))

	return resourceAlicloudHbrHanaBackupPlanUpdate(d, meta)
}
func resourceAlicloudHbrHanaBackupPlanRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	hbrService := HbrService{client}
	object, err := hbrService.DescribeHbrHanaBackupPlan(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_hbr_hana_backup_plan hbrService.DescribeHbrHanaBackupPlan Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	parts, err := ParseResourceId(d.Id(), 3)
	if err != nil {
		return WrapError(err)
	}
	d.Set("cluster_id", parts[2])
	d.Set("plan_id", parts[0])
	d.Set("vault_id", parts[1])
	d.Set("backup_prefix", object["BackupPrefix"])
	d.Set("backup_type", object["BackupType"])
	d.Set("database_name", object["DatabaseName"])
	d.Set("plan_name", object["PlanName"])
	d.Set("schedule", object["Schedule"])
	d.Set("status", convertStatusByDisabled(object["Disabled"].(bool)))
	return nil
}
func resourceAlicloudHbrHanaBackupPlanUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	hbrService := HbrService{client}
	var err error
	var response map[string]interface{}
	parts, err := ParseResourceId(d.Id(), 3)
	if err != nil {
		return WrapError(err)
	}
	d.Partial(true)

	update := false
	request := map[string]interface{}{
		"ClusterId": parts[2],
		"PlanId":    parts[0],
		"VaultId":   parts[1],
	}
	if !d.IsNewResource() && d.HasChange("backup_prefix") {
		update = true
		if v, ok := d.GetOk("backup_prefix"); ok {
			request["BackupPrefix"] = v
		}
	}
	if !d.IsNewResource() && d.HasChange("plan_name") {
		update = true
		request["PlanName"] = d.Get("plan_name")
	}
	if !d.IsNewResource() && d.HasChange("schedule") {
		update = true
		request["Schedule"] = d.Get("schedule")
	}
	if update {
		action := "UpdateHanaBackupPlan"
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("hbr", "2017-09-08", action, nil, request, false)
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
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
		if fmt.Sprint(response["Success"]) == "false" {
			return WrapError(fmt.Errorf("%s failed, response: %v", action, response))
		}
		d.SetPartial("backup_prefix")
		d.SetPartial("plan_name")
		d.SetPartial("schedule")
	}
	if d.HasChange("status") {
		object, err := hbrService.DescribeHbrHanaBackupPlan(d.Id())
		if err != nil {
			return WrapError(err)
		}
		status := convertStatusByDisabled(object["Disabled"].(bool))
		target := d.Get("status").(string)
		if status != target {
			if target == "Enabled" {
				request := map[string]interface{}{
					"ClusterId": parts[2],
					"PlanId":    parts[0],
					"VaultId":   parts[1],
				}
				action := "EnableHanaBackupPlan"
				wait := incrementalWait(3*time.Second, 3*time.Second)
				err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
					response, err = client.RpcPost("hbr", "2017-09-08", action, nil, request, false)
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
					return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
				}
				if fmt.Sprint(response["Success"]) == "false" {
					return WrapError(fmt.Errorf("%s failed, response: %v", action, response))
				}
			}
			if target == "Disabled" {
				request := map[string]interface{}{
					"ClusterId": parts[2],
					"PlanId":    parts[0],
					"VaultId":   parts[1],
				}
				action := "DisableHanaBackupPlan"
				wait := incrementalWait(3*time.Second, 3*time.Second)
				err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
					response, err = client.RpcPost("hbr", "2017-09-08", action, nil, request, false)
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
					return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
				}
				if fmt.Sprint(response["Success"]) == "false" {
					return WrapError(fmt.Errorf("%s failed, response: %v", action, response))
				}
			}
			d.SetPartial("status")
		}
	}
	d.Partial(false)
	return resourceAlicloudHbrHanaBackupPlanRead(d, meta)
}
func resourceAlicloudHbrHanaBackupPlanDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	parts, err := ParseResourceId(d.Id(), 3)
	if err != nil {
		return WrapError(err)
	}
	action := "DeleteHanaBackupPlan"
	var response map[string]interface{}
	request := map[string]interface{}{
		"ClusterId": parts[2],
		"PlanId":    parts[0],
		"VaultId":   parts[1],
	}

	if v, ok := d.GetOk("resource_group_id"); ok {
		request["ResourceGroupId"] = v
	}
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = client.RpcPost("hbr", "2017-09-08", action, nil, request, false)
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
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}
	if fmt.Sprint(response["Success"]) == "false" {
		return WrapError(fmt.Errorf("%s failed, response: %v", action, response))
	}
	return nil
}

func convertStatusByDisabled(disabled bool) (status string) {
	if disabled {
		status = "Disabled"
	} else {
		status = "Enabled"
	}
	return
}
