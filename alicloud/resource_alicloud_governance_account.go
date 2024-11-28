// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"encoding/json"
	"fmt"
	"log"
	"regexp"
	"time"

	"github.com/PaesslerAG/jsonpath"
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
			"account_tags": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
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
				ValidateFunc: StringMatch(regexp.MustCompile("^[A-Za-z0-9_.][A-Za-z0-9_.-]{0,49}[A-Za-z0-9_.].onaliyun.com$"), "The default_domain_name (with suffix) has a maximum length of 64 characters and must end with .onaliyun.com."),
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
			return WrapError(err)
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
		return nil
	})
	addDebug(action, response, request)

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
				return WrapError(err)
			}
			request["Tag"] = string(tagMapsJson)
		}
	}

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
			return nil
		})
		addDebug(action, response, request)
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
