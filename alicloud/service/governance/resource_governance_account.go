// Package governance. This file is generated automatically. Please do not modify it manually, thank you!
package governance

import (
	"encoding/json"
	"fmt"
	"log"
	"regexp"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/convert"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/validators"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceGovernanceAccount() *schema.Resource {
	return &schema.Resource{
		Create: resourceGovernanceAccountCreate,
		Read:   resourceGovernanceAccountRead,
		Update: resourceGovernanceAccountUpdate,
		Delete: resourceGovernanceAccountDelete,
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
			"account_tags": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"tag_key": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"tag_value": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
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
			"default_domain_name": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validators.StringMatch(regexp.MustCompile("^[A-Za-z0-9_.][A-Za-z0-9_.-]{0,49}[A-Za-z0-9_.].onaliyun.com$"), "The default_domain_name (with suffix) has a maximum length of 64 characters and must end with .onaliyun.com."),
			},
		},
	}
}

func resourceGovernanceAccountCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "EnrollAccount"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	if v, ok := d.GetOk("account_id"); ok {
		request["AccountUid"] = v
	}
	request["RegionId"] = client.RegionId

	if v, ok := d.GetOk("account_name_prefix"); ok {
		request["AccountNamePrefix"] = v
	}
	if v, ok := d.GetOk("folder_id"); ok {
		request["FolderId"] = v
	}
	if v, ok := d.GetOk("account_tags"); ok {
		tagMapsArray := make([]interface{}, 0)
		for _, dataLoop := range v.([]interface{}) {
			dataLoopTmp := dataLoop.(map[string]interface{})
			dataLoopMap := make(map[string]interface{})
			dataLoopMap["Value"] = dataLoopTmp["tag_value"]
			dataLoopMap["Key"] = dataLoopTmp["tag_key"]
			tagMapsArray = append(tagMapsArray, dataLoopMap)
		}
		tagMapsJson, err := json.Marshal(tagMapsArray)
		if err != nil {
			return errs.WrapError(err)
		}
		request["Tag"] = string(tagMapsJson)
	}

	request["BaselineId"] = d.Get("baseline_id")
	if v, ok := d.GetOkExists("payer_account_id"); ok {
		request["PayerAccountUid"] = v
	}
	if v, ok := d.GetOk("display_name"); ok {
		request["DisplayName"] = v
	}
	baselineItemsMaps := make([]interface{}, 0)
	if v, ok := d.GetOk("default_domain_name"); ok {
		baselineItemConfig := make(map[string]interface{})
		baselineItemConfig["DefaultDomainName"] = v
		config, _ := json.Marshal(baselineItemConfig)
		baselineItem := make(map[string]interface{})
		baselineItem["Config"] = string(config)
		baselineItem["Name"] = "ACS-BP_ACCOUNT_FACTORY_RAM_DEFAULT_DOMAIN"
		baselineItemsMaps = append(baselineItemsMaps, baselineItem)
	}
	request["BaselineItems"] = baselineItemsMaps
	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
		response, err = client.RpcPost("governance", "2021-01-20", action, query, request, false)
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
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, "alicloud_governance_account", action, errs.AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(response["AccountUid"]))

	service := NewService(client)
	stateConf := helper.BuildStateConf([]string{}, []string{"Finished"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, service.GovernanceAccountStateRefreshFunc(d.Id(), "Status", []string{"ScheduleFailed", "Failed"}))
	if _, err := stateConf.WaitForState(); err != nil {
		return errs.WrapErrorf(err, errs.IdMsg, d.Id())
	}

	return resourceGovernanceAccountRead(d, meta)
}

func resourceGovernanceAccountRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	service := NewService(client)

	objectRaw, err := service.DescribeGovernanceAccount(d.Id())
	if err != nil {
		if !d.IsNewResource() && errs.NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_governance_account DescribeGovernanceAccount Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return errs.WrapError(err)
	}

	if objectRaw["BaselineId"] != nil {
		d.Set("baseline_id", objectRaw["BaselineId"])
	}
	if objectRaw["Status"] != nil {
		d.Set("status", objectRaw["Status"])
	}
	if objectRaw["AccountUid"] != nil {
		account_idVal, err := convert.FormatInt(objectRaw["AccountUid"])
		if err != nil {
			return err
		}
		d.Set("account_id", account_idVal)
	}

	tag1Raw, _ := jsonpath.Get("$.Inputs.Tag", objectRaw)
	accountTagsMaps := make([]map[string]interface{}, 0)
	if tag1Raw != nil {
		for _, tagChild1Raw := range tag1Raw.([]interface{}) {
			accountTagsMap := make(map[string]interface{})
			tagChild1Raw := tagChild1Raw.(map[string]interface{})
			accountTagsMap["tag_key"] = tagChild1Raw["Key"]
			accountTagsMap["tag_value"] = tagChild1Raw["Value"]

			accountTagsMaps = append(accountTagsMaps, accountTagsMap)
		}
	}
	if tag1Raw != nil {
		if err := d.Set("account_tags", accountTagsMaps); err != nil {
			return err
		}
	}

	accountIdValue, err := convert.FormatInt(d.Id())
	if err != nil {
		return err
	}
	d.Set("account_id", accountIdValue)

	return nil
}

func resourceGovernanceAccountUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	update := false

	action := "EnrollAccount"
	var err error
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["AccountUid"] = d.Id()
	request["RegionId"] = client.RegionId
	if d.HasChange("account_tags") {
		update = true
		if v, ok := d.GetOk("account_tags"); ok || d.HasChange("account_tags") {
			tagMapsArray := make([]interface{}, 0)
			for _, dataLoop := range v.([]interface{}) {
				dataLoopTmp := dataLoop.(map[string]interface{})
				dataLoopMap := make(map[string]interface{})
				dataLoopMap["Value"] = dataLoopTmp["tag_value"]
				dataLoopMap["Key"] = dataLoopTmp["tag_key"]
				tagMapsArray = append(tagMapsArray, dataLoopMap)
			}
			tagMapsJson, err := json.Marshal(tagMapsArray)
			if err != nil {
				return errs.WrapError(err)
			}
			request["Tag"] = string(tagMapsJson)
		}
	}

	if d.HasChange("baseline_id") {
		update = true
	}
	request["BaselineId"] = d.Get("baseline_id")

	if update {
		wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
		err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
			response, err = client.RpcPost("governance", "2021-01-20", action, query, request, false)
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
		service := NewService(client)
		stateConf := helper.BuildStateConf([]string{}, []string{"Finished"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, service.GovernanceAccountStateRefreshFunc(d.Id(), "Status", []string{"Failed", "ScheduleFailed"}))
		if _, err := stateConf.WaitForState(); err != nil {
			return errs.WrapErrorf(err, errs.IdMsg, d.Id())
		}
	}

	return resourceGovernanceAccountRead(d, meta)
}

func resourceGovernanceAccountDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[WARN] Cannot destroy resource AliCloud Resource Account. Terraform will remove this resource from the state file, however resources may remain.")
	return nil
}
