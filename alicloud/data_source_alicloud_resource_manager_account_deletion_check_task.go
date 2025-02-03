package alicloud

import (
	"fmt"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceAlicloudResourceManagerAccountDeletionCheckTask() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudResourceManagerAccountDeletionCheckTaskRead,
		Schema: map[string]*schema.Schema{
			"account_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"allow_delete": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"not_allow_reason": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"check_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"check_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"abandon_able_checks": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"check_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"check_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceAlicloudResourceManagerAccountDeletionCheckTaskRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	resourceManagerService := ResourcemanagerService{client}
	var response map[string]interface{}
	action := "CheckAccountDelete"
	request := make(map[string]interface{})
	var err error

	request["AccountId"] = d.Get("account_id")

	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutCreate)), func() *resource.RetryError {
		response, err = client.RpcPost("ResourceManager", "2020-03-31", action, nil, request, false)
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_resource_manager_account_deletion_check_task", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(request["AccountId"]))

	stateConf := BuildStateConf([]string{}, []string{"PreCheckComplete"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, resourceManagerService.ResourceManagerAccountDeletionCheckTaskStateRefreshFunc(d.Id(), []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	object, err := resourceManagerService.DescribeResourceManagerAccountDeletionCheckTask(d.Id())
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("account_id", d.Id())

	d.Set("allow_delete", object["AllowDelete"])

	d.Set("status", object["Status"])

	if notAllowReasonList, ok := object["NotAllowReason"].([]interface{}); ok {
		notAllowReasonMaps := make([]map[string]interface{}, 0)
		for _, notAllowReason := range notAllowReasonList {
			notAllowReasonMap := map[string]interface{}{}
			notAllowReasonArg := notAllowReason.(map[string]interface{})

			if notAllowReasonCheckId, ok := notAllowReasonArg["CheckId"]; ok {
				notAllowReasonMap["check_id"] = notAllowReasonCheckId
			}

			if notAllowReasonCheckName, ok := notAllowReasonArg["CheckName"]; ok {
				notAllowReasonMap["check_name"] = notAllowReasonCheckName
			}

			if notAllowReasonDescription, ok := notAllowReasonArg["Description"]; ok {
				notAllowReasonMap["description"] = notAllowReasonDescription
			}

			notAllowReasonMaps = append(notAllowReasonMaps, notAllowReasonMap)
		}
		d.Set("not_allow_reason", notAllowReasonMaps)
	}

	if abandonAbleChecksList, ok := object["AbandonableChecks"].([]interface{}); ok {
		abandonAbleChecksMaps := make([]map[string]interface{}, 0)
		for _, abandonAbleChecks := range abandonAbleChecksList {
			abandonAbleChecksMap := map[string]interface{}{}
			abandonAbleChecksArg := abandonAbleChecks.(map[string]interface{})

			if abandonAbleChecksCheckId, ok := abandonAbleChecksArg["CheckId"]; ok {
				abandonAbleChecksMap["check_id"] = abandonAbleChecksCheckId
			}

			if abandonAbleChecksCheckName, ok := abandonAbleChecksArg["CheckName"]; ok {
				abandonAbleChecksMap["check_name"] = abandonAbleChecksCheckName
			}

			if abandonAbleChecksDescription, ok := abandonAbleChecksArg["Description"]; ok {
				abandonAbleChecksMap["description"] = abandonAbleChecksDescription
			}

			abandonAbleChecksMaps = append(abandonAbleChecksMaps, abandonAbleChecksMap)
		}
		d.Set("abandon_able_checks", abandonAbleChecksMaps)
	}

	return nil
}
