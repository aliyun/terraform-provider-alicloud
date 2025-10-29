// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAliCloudConfigAggregator() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudConfigAggregatorCreate,
		Read:   resourceAliCloudConfigAggregatorRead,
		Update: resourceAliCloudConfigAggregatorUpdate,
		Delete: resourceAliCloudConfigAggregatorDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"aggregator_accounts": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"account_id": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"account_type": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: StringInSlice([]string{"ResourceDirectory"}, false),
						},
						"account_name": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"aggregator_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"aggregator_type": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"create_time": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Required: true,
			},
			"folder_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceAliCloudConfigAggregatorCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "CreateAggregator"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})

	request["ClientToken"] = buildClientToken(action)

	if v, ok := d.GetOk("aggregator_accounts"); ok {
		aggregatorAccountsMapsArray := make([]interface{}, 0)
		for _, dataLoop := range convertToInterfaceArray(v) {
			dataLoopTmp := dataLoop.(map[string]interface{})
			dataLoopMap := make(map[string]interface{})
			dataLoopMap["AccountName"] = dataLoopTmp["account_name"]
			dataLoopMap["AccountId"] = dataLoopTmp["account_id"]
			dataLoopMap["AccountType"] = dataLoopTmp["account_type"]
			aggregatorAccountsMapsArray = append(aggregatorAccountsMapsArray, dataLoopMap)
		}
		aggregatorAccountsMapsJson, err := json.Marshal(aggregatorAccountsMapsArray)
		if err != nil {
			return WrapError(err)
		}
		request["AggregatorAccounts"] = string(aggregatorAccountsMapsJson)
	}

	if v, ok := d.GetOk("folder_id"); ok {
		request["FolderId"] = v
	}
	if v, ok := d.GetOk("aggregator_type"); ok {
		request["AggregatorType"] = v
	}
	request["Description"] = d.Get("description")
	request["AggregatorName"] = d.Get("aggregator_name")
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.RpcPost("Config", "2020-09-07", action, query, request, true)
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_config_aggregator", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(response["AggregatorId"]))

	configServiceV2 := ConfigServiceV2{client}
	stateConf := BuildStateConf([]string{}, []string{"1"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, configServiceV2.DescribeAsyncConfigAggregatorStateRefreshFunc(d, response, "$.Aggregator.AggregatorStatus", []string{}))
	if jobDetail, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id(), jobDetail)
	}

	return resourceAliCloudConfigAggregatorRead(d, meta)
}

func resourceAliCloudConfigAggregatorRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	configServiceV2 := ConfigServiceV2{client}

	objectRaw, err := configServiceV2.DescribeConfigAggregator(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_config_aggregator DescribeConfigAggregator Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("aggregator_name", objectRaw["AggregatorName"])
	d.Set("aggregator_type", objectRaw["AggregatorType"])
	d.Set("create_time", formatInt(objectRaw["AggregatorCreateTimestamp"]))
	d.Set("description", objectRaw["Description"])
	d.Set("folder_id", objectRaw["FolderId"])
	d.Set("status", convertConfigAggregatorAggregatorAggregatorStatusResponse(objectRaw["AggregatorStatus"]))

	aggregatorAccountsRaw := objectRaw["AggregatorAccounts"]
	aggregatorAccountsMaps := make([]map[string]interface{}, 0)
	if aggregatorAccountsRaw != nil {
		for _, aggregatorAccountsChildRaw := range convertToInterfaceArray(aggregatorAccountsRaw) {
			aggregatorAccountsMap := make(map[string]interface{})
			aggregatorAccountsChildRaw := aggregatorAccountsChildRaw.(map[string]interface{})
			aggregatorAccountsMap["account_id"] = aggregatorAccountsChildRaw["AccountId"]
			aggregatorAccountsMap["account_name"] = aggregatorAccountsChildRaw["AccountName"]
			aggregatorAccountsMap["account_type"] = aggregatorAccountsChildRaw["AccountType"]

			aggregatorAccountsMaps = append(aggregatorAccountsMaps, aggregatorAccountsMap)
		}
	}
	if err := d.Set("aggregator_accounts", aggregatorAccountsMaps); err != nil {
		return err
	}

	return nil
}

func resourceAliCloudConfigAggregatorUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	update := false

	var err error
	action := "UpdateAggregator"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["AggregatorId"] = d.Id()

	request["ClientToken"] = buildClientToken(action)
	if d.HasChange("aggregator_accounts") {
		update = true
	}
	if v, ok := d.GetOk("aggregator_accounts"); ok && d.Get("aggregator_type") == "CUSTOM" {
		aggregatorAccountsMapsArray := make([]interface{}, 0)
		for _, dataLoop := range convertToInterfaceArray(v) {
			dataLoopTmp := dataLoop.(map[string]interface{})
			dataLoopMap := make(map[string]interface{})
			dataLoopMap["AccountName"] = dataLoopTmp["account_name"]
			dataLoopMap["AccountId"] = dataLoopTmp["account_id"]
			dataLoopMap["AccountType"] = dataLoopTmp["account_type"]
			aggregatorAccountsMapsArray = append(aggregatorAccountsMapsArray, dataLoopMap)
		}
		aggregatorAccountsMapsJson, err := json.Marshal(aggregatorAccountsMapsArray)
		if err != nil {
			return WrapError(err)
		}
		request["AggregatorAccounts"] = string(aggregatorAccountsMapsJson)
	}

	if d.HasChange("folder_id") {
		update = true
	}
	if v, ok := d.GetOk("folder_id"); ok {
		request["FolderId"] = v
	}

	if d.HasChange("description") {
		update = true
	}
	request["Description"] = d.Get("description")
	if d.HasChange("aggregator_name") {
		update = true
	}
	request["AggregatorName"] = d.Get("aggregator_name")
	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("Config", "2020-09-07", action, query, request, true)
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
	}

	return resourceAliCloudConfigAggregatorRead(d, meta)
}

func resourceAliCloudConfigAggregatorDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	action := "DeleteAggregators"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	request["AggregatorIds"] = d.Id()

	request["ClientToken"] = buildClientToken(action)

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = client.RpcPost("Config", "2020-09-07", action, query, request, true)
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
		if IsExpectedErrors(err, []string{"Invalid.AggregatorIds.Empty", "AccountNotExisted"}) || NotFoundError(err) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	return nil
}

func convertConfigAggregatorAggregatorAggregatorStatusResponse(source interface{}) interface{} {
	source = fmt.Sprint(source)
	switch source {
	case "0":
		return "Creating"
	case "1":
		return "Normal"
	case "2":
		return "Deleting"
	}
	return source
}

func convertConfigAggregatorStatusResponse(source interface{}) interface{} {
	switch source {
	case 0:
		return "Creating"
	case 2:
		return "Deleting"
	case 1:
		return "Normal"
	}
	return ""
}
