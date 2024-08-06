// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
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

func resourceAliCloudGovernanceAccount() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudGovernanceAccountCreate,
		Read:   resourceAliCloudGovernanceAccountRead,
		Update: resourceAliCloudGovernanceAccountUpdate,
		Delete: resourceAliCloudGovernanceAccountDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"account_id": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"account_name_prefix": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"baseline_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"display_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"folder_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"payer_account_id": {
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

func resourceAliCloudGovernanceAccountCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "EnrollAccount"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	conn, err := client.NewGovernanceClient()
	if err != nil {
		return WrapError(err)
	}
	request = make(map[string]interface{})
	if v, ok := d.GetOk("account_id"); ok {
		query["AccountUid"] = v
	}
	query["RegionId"] = client.RegionId

	if v, ok := d.GetOk("account_name_prefix"); ok {
		request["AccountNamePrefix"] = v
	}
	request["BaselineId"] = d.Get("baseline_id")
	if v, ok := d.GetOk("display_name"); ok {
		request["DisplayName"] = v
	}
	if v, ok := d.GetOk("folder_id"); ok {
		request["FolderId"] = v
	}
	if v, ok := d.GetOk("payer_account_id"); ok {
		request["PayerAccountUid"] = v
	}
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2021-01-20"), StringPointer("AK"), query, request, &runtime)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(action, response, request)
		return nil
	})

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_governance_account", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(response["AccountUid"]))

	governanceServiceV2 := GovernanceServiceV2{client}
	stateConf := BuildStateConf([]string{}, []string{"Finished"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, governanceServiceV2.GovernanceAccountStateRefreshFunc(d.Id(), "Status", []string{"ScheduleFailed", "Failed"}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return resourceAliCloudGovernanceAccountRead(d, meta)
}

func resourceAliCloudGovernanceAccountRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	governanceServiceV2 := GovernanceServiceV2{client}

	objectRaw, err := governanceServiceV2.DescribeGovernanceAccount(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_governance_account DescribeGovernanceAccount Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	if objectRaw["BaselineId"] != nil {
		d.Set("baseline_id", objectRaw["BaselineId"])
	}
	if objectRaw["Status"] != nil {
		d.Set("status", objectRaw["Status"])
	}
	if objectRaw["AccountUid"] != nil {
		d.Set("account_id", objectRaw["AccountUid"])
	}

	d.Set("account_id", d.Id())

	return nil
}

func resourceAliCloudGovernanceAccountUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	update := false
	action := "EnrollAccount"
	conn, err := client.NewGovernanceClient()
	if err != nil {
		return WrapError(err)
	}
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	query["AccountUid"] = d.Id()
	query["RegionId"] = client.RegionId
	if d.HasChange("baseline_id") {
		update = true
	}
	request["BaselineId"] = d.Get("baseline_id")
	if update {
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2021-01-20"), StringPointer("AK"), query, request, &runtime)
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			addDebug(action, response, request)
			return nil
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
		governanceServiceV2 := GovernanceServiceV2{client}
		stateConf := BuildStateConf([]string{}, []string{"Finished"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, governanceServiceV2.GovernanceAccountStateRefreshFunc(d.Id(), "Status", []string{"Failed", "ScheduleFailed"}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
	}

	return resourceAliCloudGovernanceAccountRead(d, meta)
}

func resourceAliCloudGovernanceAccountDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[WARN] Cannot destroy resource AliCloud Resource Account. Terraform will remove this resource from the state file, however resources may remain.")
	return nil
}
