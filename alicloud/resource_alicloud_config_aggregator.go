package alicloud

import (
	"fmt"
	"log"
	"time"

	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceAlicloudConfigAggregator() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudConfigAggregatorCreate,
		Read:   resourceAlicloudConfigAggregatorRead,
		Update: resourceAlicloudConfigAggregatorUpdate,
		Delete: resourceAlicloudConfigAggregatorDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(1 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"aggregator_accounts": {
				Type:     schema.TypeSet,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"account_id": {
							Type:     schema.TypeInt,
							Required: true,
						},
						"account_name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"account_type": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
			"aggregator_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"aggregator_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"CUSTOM", "RD"}, false),
				Default:      "CUSTOM",
			},
			"description": {
				Type:     schema.TypeString,
				Required: true,
			},
			"status": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func resourceAlicloudConfigAggregatorCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	configService := ConfigService{client}
	var response map[string]interface{}
	action := "CreateAggregator"
	request := make(map[string]interface{})
	conn, err := client.NewConfigClient()
	if err != nil {
		return WrapError(err)
	}
	aggregatorAccounts, err := configService.convertAggregatorAccountsToString(d.Get("aggregator_accounts").(*schema.Set).List())
	if err != nil {
		return WrapError(err)
	}
	request["AggregatorAccounts"] = aggregatorAccounts
	request["AggregatorName"] = d.Get("aggregator_name")
	if v, ok := d.GetOk("aggregator_type"); ok {
		request["AggregatorType"] = v
	}

	request["Description"] = d.Get("description")
	request["ClientToken"] = buildClientToken("CreateAggregator")
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-09-07"), StringPointer("AK"), nil, request, &runtime)
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_config_aggregator", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(response["AggregatorId"]))
	stateConf := BuildStateConf([]string{}, []string{"1"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, configService.ConfigAggregatorStateRefreshFunc(d.Id(), []string{"0"}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return resourceAlicloudConfigAggregatorRead(d, meta)
}
func resourceAlicloudConfigAggregatorRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	configService := ConfigService{client}
	object, err := configService.DescribeConfigAggregator(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_config_aggregator configService.DescribeConfigAggregator Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	if v, ok := object["AggregatorAccounts"].([]interface{}); ok {
		aggregatorAccounts := make([]map[string]interface{}, 0)
		for _, val := range v {
			item := val.(map[string]interface{})
			aggregatorAccounts = append(aggregatorAccounts, map[string]interface{}{
				"account_id":   item["AccountId"],
				"account_name": item["AccountName"],
				"account_type": item["AccountType"],
			})
		}
		if err := d.Set("aggregator_accounts", aggregatorAccounts); err != nil {
			return WrapError(err)
		}
	}
	d.Set("aggregator_name", object["AggregatorName"])
	d.Set("aggregator_type", object["AggregatorType"])
	d.Set("description", object["Description"])
	d.Set("status", formatInt(object["AggregatorStatus"]))
	return nil
}
func resourceAlicloudConfigAggregatorUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	configService := ConfigService{client}
	var response map[string]interface{}
	update := false
	request := map[string]interface{}{
		"AggregatorId": d.Id(),
	}
	if d.HasChange("aggregator_accounts") {
		update = true
	}
	aggregatorAccounts, err := configService.convertAggregatorAccountsToString(d.Get("aggregator_accounts").(*schema.Set).List())
	if err != nil {
		return WrapError(err)
	}
	request["AggregatorAccounts"] = aggregatorAccounts
	if d.HasChange("aggregator_name") {
		update = true
	}
	request["AggregatorName"] = d.Get("aggregator_name")
	if d.HasChange("description") {
		update = true
	}
	request["Description"] = d.Get("description")
	if update {
		action := "UpdateAggregator"
		conn, err := client.NewConfigClient()
		if err != nil {
			return WrapError(err)
		}
		request["ClientToken"] = buildClientToken("UpdateAggregator")
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-09-07"), StringPointer("AK"), nil, request, &runtime)
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
	}
	return resourceAlicloudConfigAggregatorRead(d, meta)
}
func resourceAlicloudConfigAggregatorDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	action := "DeleteAggregators"
	var response map[string]interface{}
	conn, err := client.NewConfigClient()
	if err != nil {
		return WrapError(err)
	}
	request := map[string]interface{}{
		"AggregatorIds": d.Id(),
	}

	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-09-07"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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
		if IsExpectedErrors(err, []string{"AccountNotExisted", "Invalid.AggregatorIds.Empty"}) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}
	return nil
}
