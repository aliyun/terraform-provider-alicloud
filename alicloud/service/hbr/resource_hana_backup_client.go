package hbr

import (
	"fmt"
	"log"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func ResourceHanaBackupClient() *schema.Resource {
	return &schema.Resource{
		Create: ResourceHanaBackupClientCreate,
		Read:   ResourceHanaBackupClientRead,
		Update: ResourceHanaBackupClientUpdate,
		Delete: ResourceHanaBackupClientDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"vault_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"client_info": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"alert_setting": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Computed:     true,
				ValidateFunc: validation.StringInSlice([]string{"INHERITED"}, false),
			},
			"use_https": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
			},
			"client_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"instance_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"cluster_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func ResourceHanaBackupClientCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	hbrService := HbrService{client}
	var response map[string]interface{}
	action := "CreateClients"
	request := make(map[string]interface{})
	var err error

	request["VaultId"] = d.Get("vault_id")

	if v, ok := d.GetOk("client_info"); ok {
		request["ClientInfo"] = v
	}

	if v, ok := d.GetOk("alert_setting"); ok {
		request["AlertSetting"] = v
	}

	if v, ok := d.GetOk("use_https"); ok {
		request["UseHttps"] = v
	}

	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutCreate)), func() *retry.RetryError {
		response, err = client.RpcPost("hbr", "2017-09-08", action, nil, request, true)
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
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, "alicloud_hbr_hana_backup_client", action, errs.AlibabaCloudSdkGoERROR)
	}

	if fmt.Sprint(response["Success"]) == "false" {
		return errs.WrapError(fmt.Errorf("%s failed, response: %v", action, response))
	}

	taskId := fmt.Sprint(response["TaskId"])
	taskConf := helper.BuildStateConf([]string{}, []string{"ACTIVATED"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, hbrService.HbrTaskRefreshFunc(taskId, []string{"INSTALL_FAILED", "DEACTIVATED", "UNKNOWN"}))
	if _, err := taskConf.WaitForState(); err != nil {
		return errs.WrapErrorf(err, errs.IdMsg, d.Id())
	}

	taskResult, err := hbrService.DescribeHbrTask(taskId)
	if err != nil {
		return errs.WrapErrorf(err, errs.IdMsg, d.Id())
	}

	d.SetId(fmt.Sprintf("%v:%v", request["VaultId"], taskResult["ClientId"]))

	stateConf := helper.BuildStateConf([]string{}, []string{"ACTIVATED"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, hbrService.HbrHanaBackupClientStateRefreshFunc(d.Id(), []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return errs.WrapErrorf(err, errs.IdMsg, d.Id())
	}

	return ResourceHanaBackupClientRead(d, meta)
}

func ResourceHanaBackupClientRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	hbrService := HbrService{client}

	object, err := hbrService.DescribeHbrHanaBackupClient(d.Id())
	if err != nil {
		if errs.NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return errs.WrapError(err)
	}

	d.Set("vault_id", object["VaultId"])
	d.Set("client_id", object["ClientId"])
	d.Set("alert_setting", object["AlertSetting"])
	d.Set("use_https", object["UseHttps"])
	d.Set("instance_id", object["InstanceId"])
	d.Set("cluster_id", object["ClusterId"])
	d.Set("status", object["Status"])

	return nil
}

func ResourceHanaBackupClientUpdate(d *schema.ResourceData, meta interface{}) error {
	log.Println(fmt.Sprintf("[WARNING] The resouce has not update operation."))
	return ResourceHanaBackupClientRead(d, meta)
}

func ResourceHanaBackupClientDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	hbrService := HbrService{client}
	action := "DeleteClient"
	var response map[string]interface{}

	var err error

	parts, err := helper.ParseResourceId(d.Id(), 2)
	if err != nil {
		return errs.WrapError(err)
	}

	request := map[string]interface{}{
		"VaultId":  parts[0],
		"ClientId": parts[1],
	}

	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutDelete)), func() *retry.RetryError {
		response, err = client.RpcPost("hbr", "2017-09-08", action, nil, request, true)
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
		if errs.NotFoundError(err) {
			return nil
		}
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
	}

	stateConf := helper.BuildStateConf([]string{}, []string{}, d.Timeout(schema.TimeoutDelete), 5*time.Second, hbrService.HbrHanaBackupClientStateRefreshFunc(d.Id(), []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return errs.WrapErrorf(err, errs.IdMsg, d.Id())
	}

	return nil
}
