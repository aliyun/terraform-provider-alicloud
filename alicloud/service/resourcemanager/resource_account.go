package resourcemanager

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/convert"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/tags"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/validators"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceAccount() *schema.Resource {
	return &schema.Resource{
		Create: ResourceAccountCreate,
		Read:   ResourceAccountRead,
		Update: ResourceAccountUpdate,
		Delete: ResourceAccountDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"abandonable_check_id": {
				Type:          schema.TypeList,
				Optional:      true,
				ConflictsWith: []string{"abandon_able_check_id"},
				Elem:          &schema.Schema{Type: schema.TypeString},
			},
			"account_name_prefix": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"display_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"folder_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"join_method": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"join_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"modify_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"payer_account_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"resell_account_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validators.StringInSlice([]string{"resell", "non_resell"}, false),
			},
			"resource_directory_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"tags": tags.TagsSchema(),
			"type": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"force_delete": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"abandon_able_check_id": {
				Type:       schema.TypeList,
				Optional:   true,
				Deprecated: "Field 'abandon_able_check_id' has been deprecated since provider version 1.248.0. New field 'abandonable_check_id' instead.",
				Elem:       &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func ResourceAccountCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "CreateResourceAccount"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})

	if v, ok := d.GetOk("tags"); ok {
		tagsMap := tags.ConvertTags(v.(map[string]interface{}))
		request["Tags"] = tagsMap
	}

	if v, ok := d.GetOk("payer_account_id"); ok {
		request["PayerAccountId"] = v
	}
	request["DisplayName"] = d.Get("display_name")
	if v, ok := d.GetOk("resell_account_type"); ok {
		request["ResellAccountType"] = v
	}
	if v, ok := d.GetOk("account_name_prefix"); ok {
		request["AccountNamePrefix"] = v
	}
	if v, ok := d.GetOk("folder_id"); ok {
		request["ParentFolderId"] = v
	}
	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
		response, err = client.RpcPost("ResourceManager", "2020-03-31", action, query, request, true)
		if err != nil {
			if errs.IsExpectedErrors(err, []string{"ConcurrentCallNotSupported"}) || helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	logs.AddDebug(action, response, request)

	if err != nil {
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, "alicloud_resource_manager_account", action, errs.AlibabaCloudSdkGoERROR)
	}

	id, _ := jsonpath.Get("$.Account.AccountId", response)
	d.SetId(fmt.Sprint(id))

	return ResourceAccountUpdate(d, meta)
}

func ResourceAccountRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	resourceManagerServiceV2 := ResourceManagerServiceV2{client}

	objectRaw, err := resourceManagerServiceV2.DescribeResourceManagerAccount(d.Id())
	if err != nil {
		if !d.IsNewResource() && errs.NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_resource_manager_account DescribeResourceManagerAccount Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return errs.WrapError(err)
	}

	d.Set("display_name", objectRaw["DisplayName"])
	d.Set("folder_id", objectRaw["FolderId"])
	d.Set("join_method", objectRaw["JoinMethod"])
	d.Set("join_time", objectRaw["JoinTime"])
	d.Set("modify_time", objectRaw["ModifyTime"])
	d.Set("resource_directory_id", objectRaw["ResourceDirectoryId"])
	d.Set("status", objectRaw["Status"])
	d.Set("type", objectRaw["Type"])

	resourcemanagerService := ResourcemanagerService{client}
	getPayerForAccountObject, err := resourcemanagerService.GetPayerForAccount(d.Id())
	if err != nil {
		return errs.WrapError(err)
	}
	d.Set("payer_account_id", getPayerForAccountObject["PayerAccountId"])

	listTagResourcesObject, err := resourcemanagerService.ListTagResources(d.Id(), "Account")
	if err != nil {
		return errs.WrapError(err)
	}
	d.Set("tags", tags.TagsToMap(listTagResourcesObject))

	return nil
}

func ResourceAccountUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	update := false
	d.Partial(true)

	var err error
	action := "MoveAccount"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["AccountId"] = d.Id()

	if !d.IsNewResource() && d.HasChange("folder_id") {
		update = true
	}
	request["DestinationFolderId"] = d.Get("folder_id")
	if update {
		wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
		err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
			response, err = client.RpcPost("ResourceManager", "2020-03-31", action, query, request, true)
			if err != nil {
				if errs.IsExpectedErrors(err, []string{"ConcurrentCallNotSupported"}) || helper.NeedRetry(err) {
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
	}
	update = false
	action = "UpdateAccount"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["AccountId"] = d.Id()

	if d.HasChange("type") {
		update = true
		request["NewAccountType"] = d.Get("type")
	}

	if !d.IsNewResource() && d.HasChange("display_name") {
		update = true
		request["NewDisplayName"] = d.Get("display_name")
	}
	if update {
		wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
		err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
			response, err = client.RpcPost("ResourceManager", "2020-03-31", action, query, request, true)
			if err != nil {
				if errs.IsExpectedErrors(err, []string{"ConcurrentCallNotSupported"}) || helper.NeedRetry(err) {
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
	}
	update = false
	action = "UpdatePayerForAccount"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["AccountId"] = d.Id()

	if !d.IsNewResource() && d.HasChange("payer_account_id") {
		update = true
	}
	request["PayerAccountId"] = d.Get("payer_account_id")
	if update {
		wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
		err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
			response, err = client.RpcPost("ResourceDirectoryMaster", "2022-04-19", action, query, request, true)
			if err != nil {
				if errs.IsExpectedErrors(err, []string{"ConcurrentCallNotSupported"}) || helper.NeedRetry(err) {
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
	}

	if d.HasChange("tags") {
		resourceManagerServiceV2 := ResourceManagerServiceV2{client}
		if err := resourceManagerServiceV2.SetResourceTags(d, "Account"); err != nil {
			return errs.WrapError(err)
		}
	}
	d.Partial(false)
	return ResourceAccountRead(d, meta)
}

func ResourceAccountDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	action := "DeleteAccount"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	request["AccountId"] = d.Id()

	if v, ok := d.GetOk("abandonable_check_id"); ok {
		abandonableCheckIdMapsArray := v.([]interface{})
		abandonableCheckIdMapsJson, err := json.Marshal(abandonableCheckIdMapsArray)
		if err != nil {
			return errs.WrapError(err)
		}
		request["AbandonableCheckId"] = string(abandonableCheckIdMapsJson)
	}

	if v, ok := d.GetOk("abandon_able_check_id"); ok {
		request["AbandonableCheckId"] = convert.ConvertListToJsonString(v.([]interface{}))
	}

	if v, ok := d.GetOkExists("force_delete"); ok && v.(bool) {
		abandonableCheckIds, err := preCheckResourceManagerAccountDelete(d, meta)
		request["AbandonableCheckId"] = convert.ConvertListToJsonString(convert.ConvertListStringToListInterface(abandonableCheckIds))
		if err != nil {
			return errs.WrapError(err)
		}
	}

	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutDelete), func() *retry.RetryError {
		response, err = client.RpcPost("ResourceManager", "2020-03-31", action, query, request, true)

		if err != nil {
			if errs.IsExpectedErrors(err, []string{"ConcurrentCallNotSupported", "NotSupportedOperation.PreCheckingAccount"}) || helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	logs.AddDebug(action, response, request)

	if err != nil {
		if errs.IsExpectedErrors(err, []string{"EntityNotExists.ResourceDirectory"}) || errs.NotFoundError(err) {
			return nil
		}
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
	}

	resourceManagerServiceV2 := ResourceManagerServiceV2{client}
	stateConf := helper.BuildStateConf([]string{}, []string{"Success", "Deleting"}, d.Timeout(schema.TimeoutDelete), 5*time.Second, resourceManagerServiceV2.DescribeAsyncResourceManagerAccountStateRefreshFunc(d, response, "$.RdAccountDeletionStatus.Status", []string{"CheckFailed", "DeleteFailed"}))
	if jobDetail, err := stateConf.WaitForState(); err != nil {
		return errs.WrapErrorf(err, errs.IdMsg, d.Id(), jobDetail)
	}

	return nil
}

func preCheckResourceManagerAccountDelete(d *schema.ResourceData, meta interface{}) ([]string, error) {
	client := meta.(*connectivity.AliyunClient)
	resourceManagerService := ResourcemanagerService{client}
	var response map[string]interface{}
	action := "CheckAccountDelete"
	request := make(map[string]interface{})
	var err error
	request["AccountId"] = d.Id()
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutCreate)), func() *retry.RetryError {
		response, err = client.RpcPost("ResourceManager", "2020-03-31", action, nil, request, false)
		if err != nil {
			if errs.IsExpectedErrors(err, []string{"ConcurrentCallNotSupported"}) || helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	logs.AddDebug(action, response, request)
	if err != nil {
		return nil, errs.WrapErrorf(err, errs.DefaultErrorMsg, "alicloud_resource_manager_account", action, errs.AlibabaCloudSdkGoERROR)
	}

	stateConf := helper.BuildStateConf([]string{}, []string{"PreCheckComplete"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, resourceManagerService.ResourceManagerAccountDeletionCheckTaskStateRefreshFunc(d.Id(), []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return nil, errs.WrapErrorf(err, errs.IdMsg, d.Id())
	}

	object, err := resourceManagerService.DescribeResourceManagerAccountDeletionCheckTask(d.Id())
	abandonAbleCheckIds := make([]string, 0)
	if abandonAbleChecksList, ok := object["AbandonableChecks"].([]interface{}); ok {
		for _, abandonAbleChecks := range abandonAbleChecksList {
			abandonAbleChecksArg := abandonAbleChecks.(map[string]interface{})
			if abandonAbleChecksCheckId, ok := abandonAbleChecksArg["CheckId"]; ok {
				abandonAbleCheckIds = append(abandonAbleCheckIds, fmt.Sprint(abandonAbleChecksCheckId))
			}
		}
	}
	return abandonAbleCheckIds, nil
}
