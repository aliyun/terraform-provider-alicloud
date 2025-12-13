package alicloud

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAliCloudResourceManagerAccount() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudResourceManagerAccountCreate,
		Read:   resourceAliCloudResourceManagerAccountRead,
		Update: resourceAliCloudResourceManagerAccountUpdate,
		Delete: resourceAliCloudResourceManagerAccountDelete,
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
				ValidateFunc: StringInSlice([]string{"resell", "non_resell"}, false),
			},
			"resource_directory_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"tags": tagsSchema(),
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

func resourceAliCloudResourceManagerAccountCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "CreateResourceAccount"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})

	if v, ok := d.GetOk("tags"); ok {
		tagsMap := ConvertTags(v.(map[string]interface{}))
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
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.RpcPost("ResourceManager", "2020-03-31", action, query, request, true)
		if err != nil {
			if IsExpectedErrors(err, []string{"ConcurrentCallNotSupported"}) || NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_resource_manager_account", action, AlibabaCloudSdkGoERROR)
	}

	id, _ := jsonpath.Get("$.Account.AccountId", response)
	d.SetId(fmt.Sprint(id))

	return resourceAliCloudResourceManagerAccountUpdate(d, meta)
}

func resourceAliCloudResourceManagerAccountRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	resourceManagerServiceV2 := ResourceManagerServiceV2{client}

	objectRaw, err := resourceManagerServiceV2.DescribeResourceManagerAccount(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_resource_manager_account DescribeResourceManagerAccount Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
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
		return WrapError(err)
	}
	d.Set("payer_account_id", getPayerForAccountObject["PayerAccountId"])

	listTagResourcesObject, err := resourcemanagerService.ListTagResources(d.Id(), "Account")
	if err != nil {
		return WrapError(err)
	}
	d.Set("tags", tagsToMap(listTagResourcesObject))

	return nil
}

func resourceAliCloudResourceManagerAccountUpdate(d *schema.ResourceData, meta interface{}) error {
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
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("ResourceManager", "2020-03-31", action, query, request, true)
			if err != nil {
				if IsExpectedErrors(err, []string{"ConcurrentCallNotSupported"}) || NeedRetry(err) {
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
	}
	request["NewDisplayName"] = d.Get("display_name")
	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("ResourceManager", "2020-03-31", action, query, request, true)
			if err != nil {
				if IsExpectedErrors(err, []string{"ConcurrentCallNotSupported"}) || NeedRetry(err) {
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
	}
	update = false
	action = "UpdatePayerForAccount"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["AccountId"] = d.Id()

	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("ResourceDirectoryMaster", "2022-04-19", action, query, request, true)
			if err != nil {
				if IsExpectedErrors(err, []string{"ConcurrentCallNotSupported"}) || NeedRetry(err) {
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
	}

	if d.HasChange("tags") {
		resourceManagerServiceV2 := ResourceManagerServiceV2{client}
		if err := resourceManagerServiceV2.SetResourceTags(d, "Account"); err != nil {
			return WrapError(err)
		}
	}
	d.Partial(false)
	return resourceAliCloudResourceManagerAccountRead(d, meta)
}

func resourceAliCloudResourceManagerAccountDelete(d *schema.ResourceData, meta interface{}) error {

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
			return WrapError(err)
		}
		request["AbandonableCheckId"] = string(abandonableCheckIdMapsJson)
	}

	if v, ok := d.GetOk("abandon_able_check_id"); ok {
		request["AbandonableCheckId"] = convertListToJsonString(v.([]interface{}))
	}

	if v, ok := d.GetOkExists("force_delete"); ok && v.(bool) {
		abandonableCheckIds, err := preCheckResourceManagerAccountDelete(d, meta)
		request["AbandonableCheckId"] = convertListToJsonString(convertListStringToListInterface(abandonableCheckIds))
		if err != nil {
			return WrapError(err)
		}
	}

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = client.RpcPost("ResourceManager", "2020-03-31", action, query, request, true)

		if err != nil {
			if IsExpectedErrors(err, []string{"ConcurrentCallNotSupported", "NotSupportedOperation.PreCheckingAccount"}) || NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)

	if err != nil {
		if IsExpectedErrors(err, []string{"EntityNotExists.ResourceDirectory"}) || NotFoundError(err) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	resourceManagerServiceV2 := ResourceManagerServiceV2{client}
	stateConf := BuildStateConf([]string{}, []string{"Success", "Deleting"}, d.Timeout(schema.TimeoutDelete), 5*time.Second, resourceManagerServiceV2.DescribeAsyncResourceManagerAccountStateRefreshFunc(d, response, "$.RdAccountDeletionStatus.Status", []string{"CheckFailed", "DeleteFailed"}))
	if jobDetail, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id(), jobDetail)
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
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutCreate)), func() *resource.RetryError {
		response, err = client.RpcPost("ResourceManager", "2020-03-31", action, nil, request, false)
		if err != nil {
			if IsExpectedErrors(err, []string{"ConcurrentCallNotSupported"}) || NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)
	if err != nil {
		return nil, WrapErrorf(err, DefaultErrorMsg, "alicloud_resource_manager_account", action, AlibabaCloudSdkGoERROR)
	}

	stateConf := BuildStateConf([]string{}, []string{"PreCheckComplete"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, resourceManagerService.ResourceManagerAccountDeletionCheckTaskStateRefreshFunc(d.Id(), []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return nil, WrapErrorf(err, IdMsg, d.Id())
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
