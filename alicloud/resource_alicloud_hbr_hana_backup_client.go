package alicloud

import (
	"fmt"
	"log"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceAlicloudHbrHanaBackupClient() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudHbrHanaBackupClientCreate,
		Read:   resourceAlicloudHbrHanaBackupClientRead,
		Update: resourceAlicloudHbrHanaBackupClientUpdate,
		Delete: resourceAlicloudHbrHanaBackupClientDelete,
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

func resourceAlicloudHbrHanaBackupClientCreate(d *schema.ResourceData, meta interface{}) error {
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

	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutCreate)), func() *resource.RetryError {
		response, err = client.RpcPost("hbr", "2017-09-08", action, nil, request, true)
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_hbr_hana_backup_client", action, AlibabaCloudSdkGoERROR)
	}

	if fmt.Sprint(response["Success"]) == "false" {
		return WrapError(fmt.Errorf("%s failed, response: %v", action, response))
	}

	taskId := fmt.Sprint(response["TaskId"])
	taskConf := BuildStateConf([]string{}, []string{"ACTIVATED"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, hbrService.HbrTaskRefreshFunc(taskId, []string{"INSTALL_FAILED", "DEACTIVATED", "UNKNOWN"}))
	if _, err := taskConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	taskResult, err := hbrService.DescribeHbrTask(taskId)
	if err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	d.SetId(fmt.Sprintf("%v:%v", request["VaultId"], taskResult["ClientId"]))

	stateConf := BuildStateConf([]string{}, []string{"ACTIVATED"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, hbrService.HbrHanaBackupClientStateRefreshFunc(d.Id(), []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return resourceAlicloudHbrHanaBackupClientRead(d, meta)
}

func resourceAlicloudHbrHanaBackupClientRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	hbrService := HbrService{client}

	object, err := hbrService.DescribeHbrHanaBackupClient(d.Id())
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
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

func resourceAlicloudHbrHanaBackupClientUpdate(d *schema.ResourceData, meta interface{}) error {
	log.Println(fmt.Sprintf("[WARNING] The resouce has not update operation."))
	return resourceAlicloudHbrHanaBackupClientRead(d, meta)
}

func resourceAlicloudHbrHanaBackupClientDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	hbrService := HbrService{client}
	action := "DeleteClient"
	var response map[string]interface{}

	var err error

	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}

	request := map[string]interface{}{
		"VaultId":  parts[0],
		"ClientId": parts[1],
	}

	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutDelete)), func() *resource.RetryError {
		response, err = client.RpcPost("hbr", "2017-09-08", action, nil, request, true)
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
		if NotFoundError(err) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	stateConf := BuildStateConf([]string{}, []string{}, d.Timeout(schema.TimeoutDelete), 5*time.Second, hbrService.HbrHanaBackupClientStateRefreshFunc(d.Id(), []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return nil
}
