// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package resourcemanager

import (
	"fmt"
	"log"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var resourceManagerHandshakeRpcPost = func(client *connectivity.AliyunClient, apiProductCode string, apiVersion string, apiName string, query map[string]interface{}, body map[string]interface{}, autoRetry bool) (map[string]interface{}, error) {
	return client.RpcPost(apiProductCode, apiVersion, apiName, query, body, autoRetry)
}

func ResourceHandshake() *schema.Resource {
	return &schema.Resource{
		Create: ResourceHandshakeCreate,
		Read:   ResourceHandshakeRead,
		Delete: ResourceHandshakeDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"expire_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"master_account_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"master_account_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"modify_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"note": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"resource_directory_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"target_entity": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"target_type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func ResourceHandshakeCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "InviteAccountToResourceDirectory"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})

	request["TargetType"] = d.Get("target_type")
	if v, ok := d.GetOk("note"); ok {
		request["Note"] = v
	}
	request["TargetEntity"] = d.Get("target_entity")
	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
		response, err = resourceManagerHandshakeRpcPost(client, "ResourceManager", "2020-03-31", action, query, request, true)
		if err != nil {
			if errs.IsExpectedErrors(err, []string{"ConcurrentCallNotSupported", "LimitExceeded.InvitationRate", "LimitExceeded.SameTargetInvitationRate"}) || helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	logs.AddDebug(action, response, request)

	if err != nil {
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, "alicloud_resource_manager_handshake", action, errs.AlibabaCloudSdkGoERROR)
	}

	id, _ := jsonpath.Get("$.Handshake.HandshakeId", response)
	d.SetId(fmt.Sprint(id))

	return ResourceHandshakeRead(d, meta)
}

func ResourceHandshakeRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	resourceManagerServiceV2 := ResourceManagerServiceV2{client}

	objectRaw, err := resourceManagerServiceV2.DescribeResourceManagerHandshake(d.Id())
	if err != nil {
		if !d.IsNewResource() && errs.NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_resource_manager_handshake DescribeResourceManagerHandshake Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return errs.WrapError(err)
	}

	d.Set("create_time", objectRaw["CreateTime"])
	d.Set("expire_time", objectRaw["ExpireTime"])
	d.Set("master_account_id", objectRaw["MasterAccountId"])
	d.Set("master_account_name", objectRaw["MasterAccountName"])
	d.Set("modify_time", objectRaw["ModifyTime"])
	d.Set("note", objectRaw["Note"])
	d.Set("resource_directory_id", objectRaw["ResourceDirectoryId"])
	d.Set("status", objectRaw["Status"])
	d.Set("target_entity", objectRaw["TargetEntity"])
	d.Set("target_type", objectRaw["TargetType"])

	return nil
}

func ResourceHandshakeDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	action := "CancelHandshake"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	request["HandshakeId"] = d.Id()

	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutDelete), func() *retry.RetryError {
		response, err = resourceManagerHandshakeRpcPost(client, "ResourceManager", "2020-03-31", action, query, request, true)

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
		if errs.IsExpectedErrors(err, []string{"HandshakeStatusMismatch"}) {
			return ResourceHandshakeRemoveAcceptedAccount(d, meta)
		}
		if errs.IsExpectedErrors(err, []string{"EntityNotExists.Handshake"}) || errs.NotFoundError(err) {
			return nil
		}
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
	}

	return nil
}

func ResourceHandshakeRemoveAcceptedAccount(d *schema.ResourceData, meta interface{}) error {
	targetType := d.Get("target_type").(string)
	accountId := d.Get("target_entity").(string)
	if targetType != "Account" || accountId == "" {
		log.Printf("[WARN] Cannot remove accepted Resource Manager handshake %s because target_type is %q and target_entity is %q. Terraform will remove this resource from the state file, however resources may remain.", d.Id(), targetType, accountId)
		return nil
	}

	client := meta.(*connectivity.AliyunClient)
	action := "RemoveCloudAccount"
	request := map[string]interface{}{
		"AccountId": accountId,
	}
	query := make(map[string]interface{})
	var response map[string]interface{}
	var err error

	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutDelete), func() *retry.RetryError {
		response, err = resourceManagerHandshakeRpcPost(client, "ResourceManager", "2020-03-31", action, query, request, true)

		if err != nil {
			if errs.IsExpectedErrors(err, []string{"ConcurrentCallNotSupported"}) || helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			if errs.IsExpectedErrors(err, []string{"EntityNotExists.Account", "EntityNotExists.ResourceDirectory"}) || errs.NotFoundError(err) {
				return nil
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	logs.AddDebug(action, response, request)

	if err != nil {
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, accountId, action, errs.AlibabaCloudSdkGoERROR)
	}

	action = "GetAccount"
	request = map[string]interface{}{
		"AccountId": accountId,
	}
	response = nil
	wait = helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutDelete), func() *retry.RetryError {
		response, err = resourceManagerHandshakeRpcPost(client, "ResourceManager", "2020-03-31", action, query, request, true)
		if err != nil {
			if errs.IsExpectedErrors(err, []string{"EntityNotExists.Account", "EntityNotExists.ResourceDirectory"}) || errs.NotFoundError(err) {
				return nil
			}
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		wait()
		return retry.RetryableError(fmt.Errorf("Resource Manager member account %s is still attached", accountId))
	})
	logs.AddDebug(action, response, request)
	if err != nil {
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, accountId, "GetAccount", errs.AlibabaCloudSdkGoERROR)
	}

	return nil
}
