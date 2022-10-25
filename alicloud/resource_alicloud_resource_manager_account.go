package alicloud

import (
	"fmt"
	"log"
	"time"

	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAlicloudResourceManagerAccount() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudResourceManagerAccountCreate,
		Read:   resourceAlicloudResourceManagerAccountRead,
		Update: resourceAlicloudResourceManagerAccountUpdate,
		Delete: resourceAlicloudResourceManagerAccountDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(3 * time.Minute),
			Update: schema.DefaultTimeout(3 * time.Minute),
			Delete: schema.DefaultTimeout(3 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"account_name_prefix": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
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
			"payer_account_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"abandon_able_check_id": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"tags": tagsSchema(),
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
			"resource_directory_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"type": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceAlicloudResourceManagerAccountCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "CreateResourceAccount"
	request := make(map[string]interface{})
	conn, err := client.NewResourcemanagerClient()
	if err != nil {
		return WrapError(err)
	}
	if v, ok := d.GetOk("account_name_prefix"); ok {
		request["AccountNamePrefix"] = v
	}
	request["DisplayName"] = d.Get("display_name")
	if v, ok := d.GetOk("folder_id"); ok {
		request["ParentFolderId"] = v
	}
	if v, ok := d.GetOk("payer_account_id"); ok {
		request["PayerAccountId"] = v
	}

	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutCreate)), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-03-31"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_resource_manager_account", action, AlibabaCloudSdkGoERROR)
	}
	responseAccount := response["Account"].(map[string]interface{})
	d.SetId(fmt.Sprint(responseAccount["AccountId"]))

	return resourceAlicloudResourceManagerAccountUpdate(d, meta)
}

func resourceAlicloudResourceManagerAccountRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	resourcemanagerService := ResourcemanagerService{client}
	object, err := resourcemanagerService.DescribeResourceManagerAccount(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_resource_manager_account resourcemanagerService.DescribeResourceManagerAccount Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("display_name", object["DisplayName"])
	d.Set("folder_id", object["FolderId"])
	d.Set("join_method", object["JoinMethod"])
	d.Set("join_time", object["JoinTime"])
	d.Set("modify_time", object["ModifyTime"])
	d.Set("resource_directory_id", object["ResourceDirectoryId"])
	d.Set("status", object["Status"])
	d.Set("type", object["Type"])

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

func resourceAlicloudResourceManagerAccountUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	resourcemanagerService := ResourcemanagerService{client}
	conn, err := client.NewResourcemanagerClient()
	if err != nil {
		return WrapError(err)
	}
	var response map[string]interface{}
	d.Partial(true)

	if d.HasChange("tags") {
		if err := resourcemanagerService.SetResourceTags(d, "Account"); err != nil {
			return WrapError(err)
		}
		d.SetPartial("tags")
	}
	if d.HasChange("folder_id") && !d.IsNewResource() {
		request := map[string]interface{}{
			"AccountId": d.Id(),
		}
		request["DestinationFolderId"] = d.Get("folder_id")
		action := "MoveAccount"
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-03-31"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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
		d.SetPartial("folder_id")
	}
	update := false
	request := map[string]interface{}{
		"AccountId": d.Id(),
	}
	if d.HasChange("display_name") && !d.IsNewResource() {
		update = true
	}
	request["NewDisplayName"] = d.Get("display_name")
	if update {
		action := "UpdateAccount"
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-03-31"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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
		d.SetPartial("display_name")
	}
	d.Partial(false)
	return resourceAlicloudResourceManagerAccountRead(d, meta)
}

func resourceAlicloudResourceManagerAccountDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	action := "DeleteAccount"
	var response map[string]interface{}
	conn, err := client.NewResourcemanagerClient()
	if err != nil {
		return WrapError(err)
	}

	request := map[string]interface{}{
		"AccountId": d.Id(),
	}

	if v, ok := d.GetOk("abandon_able_check_id"); ok {
		request["AbandonableCheckId"] = v
	}

	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutDelete)), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-03-31"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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

	return nil
}
