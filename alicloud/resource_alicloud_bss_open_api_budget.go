// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"log"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAliCloudBssOpenApiBudget() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudBssOpenApiBudgetCreate,
		Read:   resourceAliCloudBssOpenApiBudgetRead,
		Update: resourceAliCloudBssOpenApiBudgetUpdate,
		Delete: resourceAliCloudBssOpenApiBudgetDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"budget_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"budget_type": {
				Type:     schema.TypeString,
				Required: true,
			},
			"comment": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"cycle_end_period": {
				Type:     schema.TypeString,
				Required: true,
			},
			"cycle_quota": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"cycle_period": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"quota": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"cycle_start_period": {
				Type:     schema.TypeString,
				Required: true,
			},
			"cycle_type": {
				Type:     schema.TypeString,
				Required: true,
			},
			"metric": {
				Type:     schema.TypeString,
				Required: true,
			},
			"nbid": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"query_filter": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"select_type": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"values": {
							Type:     schema.TypeList,
							Optional: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"code": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"quota": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"quota_type": {
				Type:     schema.TypeString,
				Required: true,
			},
			"warn_confs": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"msc_channels": {
							Type:     schema.TypeList,
							Optional: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"comment": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"threshold_value": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"msc_contacts": {
							Type:     schema.TypeList,
							Optional: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"event_bridge": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"warn_target": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"name": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"threshold_type": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
		},
	}
}

func resourceAliCloudBssOpenApiBudgetCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "CreateBudget"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	if v, ok := d.GetOk("budget_name"); ok {
		request["BudgetName"] = v
	}

	if v, ok := d.GetOk("comment"); ok {
		request["Comment"] = v
	}
	request["WarnConfs"] = buildBudgetWarnConfs(d.Get("warn_confs"))
	request["Quota"] = d.Get("quota")
	request["CycleQuota"] = buildBudgetCycleQuota(d.Get("cycle_quota"))
	request["QueryFilter"] = buildBudgetQueryFilter(d.Get("query_filter"))
	request["BudgetType"] = d.Get("budget_type")
	request["QuotaType"] = d.Get("quota_type")
	request["CycleEndPeriod"] = d.Get("cycle_end_period")
	if v, ok := d.GetOk("nbid"); ok {
		request["Nbid"] = v
	}
	request["CycleType"] = d.Get("cycle_type")
	request["Metric"] = d.Get("metric")
	request["CycleStartPeriod"] = d.Get("cycle_start_period")
	var endpoint string
	request["ProductCode"] = ""
	request["ProductType"] = ""
	if client.IsInternationalAccount() {
		request["ProductType"] = ""
	}
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.RpcPostWithEndpoint("BssOpenApi", "2023-09-30", action, query, request, true, endpoint)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			if !client.IsInternationalAccount() && IsExpectedErrors(err, []string{""}) {
				request["ProductCode"] = ""
				request["ProductType"] = ""
				endpoint = connectivity.BssOpenAPIEndpointInternational
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_bss_open_api_budget", action, AlibabaCloudSdkGoERROR)
	}

	id, _ := jsonpath.Get("$.data.budgetName", response)
	d.SetId(fmt.Sprint(id))

	return resourceAliCloudBssOpenApiBudgetRead(d, meta)
}

func resourceAliCloudBssOpenApiBudgetRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	bssOpenApiServiceV2 := BssOpenApiServiceV2{client}

	objectRaw, err := bssOpenApiServiceV2.DescribeBssOpenApiBudget(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_bss_open_api_budget DescribeBssOpenApiBudget Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("budget_type", objectRaw["budgetType"])
	d.Set("comment", objectRaw["comment"])
	d.Set("cycle_end_period", objectRaw["cycleEndPeriod"])
	d.Set("cycle_start_period", objectRaw["cycleStartPeriod"])
	d.Set("cycle_type", objectRaw["cycleType"])
	d.Set("metric", objectRaw["metric"])
	d.Set("quota", objectRaw["quota"])
	d.Set("quota_type", objectRaw["quotaType"])
	d.Set("budget_name", objectRaw["budgetName"])

	cycleQuotaRaw := objectRaw["cycleQuota"]
	cycleQuotaMaps := make([]map[string]interface{}, 0)
	if cycleQuotaRaw != nil {
		for _, cycleQuotaChildRaw := range convertToInterfaceArray(cycleQuotaRaw) {
			cycleQuotaMap := make(map[string]interface{})
			cycleQuotaChildRaw := cycleQuotaChildRaw.(map[string]interface{})
			cycleQuotaMap["cycle_period"] = cycleQuotaChildRaw["cyclePeriod"]
			cycleQuotaMap["quota"] = cycleQuotaChildRaw["quota"]

			cycleQuotaMaps = append(cycleQuotaMaps, cycleQuotaMap)
		}
	}
	if err := d.Set("cycle_quota", cycleQuotaMaps); err != nil {
		return err
	}
	queryFilterRaw := objectRaw["queryFilter"]
	queryFilterMaps := make([]map[string]interface{}, 0)
	if queryFilterRaw != nil {
		for _, queryFilterChildRaw := range convertToInterfaceArray(queryFilterRaw) {
			queryFilterMap := make(map[string]interface{})
			queryFilterChildRaw := queryFilterChildRaw.(map[string]interface{})
			queryFilterMap["code"] = queryFilterChildRaw["code"]
			queryFilterMap["select_type"] = queryFilterChildRaw["selectType"]

			valuesRaw := make([]interface{}, 0)
			if queryFilterChildRaw["values"] != nil {
				valuesRaw = convertToInterfaceArray(queryFilterChildRaw["values"])
			}

			queryFilterMap["values"] = valuesRaw
			queryFilterMaps = append(queryFilterMaps, queryFilterMap)
		}
	}
	if err := d.Set("query_filter", queryFilterMaps); err != nil {
		return err
	}
	warnConfRaw := objectRaw["warnConf"]
	warnConfsMaps := make([]map[string]interface{}, 0)
	if warnConfRaw != nil {
		for _, warnConfChildRaw := range convertToInterfaceArray(warnConfRaw) {
			warnConfsMap := make(map[string]interface{})
			warnConfChildRaw := warnConfChildRaw.(map[string]interface{})
			warnConfsMap["comment"] = warnConfChildRaw["comment"]
			warnConfsMap["event_bridge"] = warnConfChildRaw["eventBridge"]
			warnConfsMap["name"] = warnConfChildRaw["name"]
			warnConfsMap["threshold_type"] = warnConfChildRaw["thresholdType"]
			warnConfsMap["threshold_value"] = warnConfChildRaw["thresholdValue"]
			warnConfsMap["warn_target"] = warnConfChildRaw["warnTarget"]

			mscChannelsRaw := make([]interface{}, 0)
			if warnConfChildRaw["mscChannels"] != nil {
				mscChannelsRaw = convertToInterfaceArray(warnConfChildRaw["mscChannels"])
			}

			warnConfsMap["msc_channels"] = mscChannelsRaw
			mscContactsRaw := make([]interface{}, 0)
			if warnConfChildRaw["mscContacts"] != nil {
				mscContactsRaw = convertToInterfaceArray(warnConfChildRaw["mscContacts"])
			}

			warnConfsMap["msc_contacts"] = mscContactsRaw
			warnConfsMaps = append(warnConfsMaps, warnConfsMap)
		}
	}
	if err := d.Set("warn_confs", warnConfsMaps); err != nil {
		return err
	}

	d.Set("budget_name", d.Id())

	return nil
}

func resourceAliCloudBssOpenApiBudgetUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	update := false

	var err error
	action := "UpdateBudget"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	// UpdateBudget maps BudgetName to both originalBudgetName (old) and budgetName (new).
	request["OriginalBudgetName"] = d.Id()
	if d.HasChange("budget_name") {
		request["BudgetName"] = d.Get("budget_name")
	} else {
		request["BudgetName"] = d.Id()
	}

	if d.HasChange("comment") {
		update = true
		request["Comment"] = d.Get("comment")
	}

	if d.HasChange("warn_confs") {
		update = true
		request["WarnConfs"] = buildBudgetWarnConfs(d.Get("warn_confs"))
	}

	if d.HasChange("quota") {
		update = true
		request["Quota"] = d.Get("quota")
	}

	if d.HasChange("cycle_quota") {
		update = true
		request["CycleQuota"] = buildBudgetCycleQuota(d.Get("cycle_quota"))
	}

	if d.HasChange("query_filter") {
		update = true
		request["QueryFilter"] = buildBudgetQueryFilter(d.Get("query_filter"))
	}

	if d.HasChange("budget_type") {
		update = true
	}
	request["BudgetType"] = d.Get("budget_type")
	if d.HasChange("quota_type") {
		update = true
	}
	request["QuotaType"] = d.Get("quota_type")
	if d.HasChange("cycle_end_period") {
		update = true
	}
	request["CycleEndPeriod"] = d.Get("cycle_end_period")
	if v, ok := d.GetOk("nbid"); ok {
		request["Nbid"] = v
	}
	if d.HasChange("cycle_type") {
		update = true
	}
	request["CycleType"] = d.Get("cycle_type")
	if d.HasChange("metric") {
		update = true
	}
	request["Metric"] = d.Get("metric")
	if d.HasChange("cycle_start_period") {
		update = true
	}
	request["CycleStartPeriod"] = d.Get("cycle_start_period")
	var endpoint string
	request["ProductCode"] = ""
	request["ProductType"] = ""
	if client.IsInternationalAccount() {
		request["ProductType"] = ""
	}
	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPostWithEndpoint("BssOpenApi", "2023-09-30", action, query, request, true, endpoint)
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				if !client.IsInternationalAccount() && IsExpectedErrors(err, []string{""}) {
					request["ProductCode"] = ""
					request["ProductType"] = ""
					endpoint = connectivity.BssOpenAPIEndpointInternational
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

	return resourceAliCloudBssOpenApiBudgetRead(d, meta)
}

func resourceAliCloudBssOpenApiBudgetDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	action := "DeleteBudget"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	request["BudgetName"] = d.Id()

	if v, ok := d.GetOk("nbid"); ok {
		request["Nbid"] = v
	}
	var endpoint string
	request["ProductCode"] = ""
	request["ProductType"] = ""
	if client.IsInternationalAccount() {
		request["ProductType"] = ""
	}
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = client.RpcPostWithEndpoint("BssOpenApi", "2023-09-30", action, query, request, true, endpoint)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			if !client.IsInternationalAccount() && IsExpectedErrors(err, []string{""}) {
				request["ProductCode"] = ""
				request["ProductType"] = ""
				endpoint = connectivity.BssOpenAPIEndpointInternational
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)

	if err != nil {
		if NotFoundError(err) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	return nil
}

// buildBudgetWarnConfs converts the warn_confs schema list into the API request array.
func buildBudgetWarnConfs(raw interface{}) []interface{} {
	result := make([]interface{}, 0)
	if raw == nil {
		return result
	}
	list, ok := raw.([]interface{})
	if !ok {
		return result
	}
	for _, itemRaw := range list {
		item, ok := itemRaw.(map[string]interface{})
		if !ok {
			continue
		}
		m := make(map[string]interface{})
		if v, ok := item["name"].(string); ok && v != "" {
			m["name"] = v
		}
		if v, ok := item["warn_target"].(string); ok && v != "" {
			m["warnTarget"] = v
		}
		if v, ok := item["threshold_type"].(string); ok && v != "" {
			m["thresholdType"] = v
		}
		if v, ok := item["threshold_value"].(string); ok && v != "" {
			m["thresholdValue"] = v
		}
		if v := item["msc_contacts"]; v != nil {
			if contacts, ok := v.([]interface{}); ok && len(contacts) > 0 {
				m["mscContacts"] = contacts
			}
		}
		if v := item["msc_channels"]; v != nil {
			if channels, ok := v.([]interface{}); ok && len(channels) > 0 {
				m["mscChannels"] = channels
			}
		}
		if v, ok := item["event_bridge"]; ok && v != nil {
			m["eventBridge"] = v
		}
		if v, ok := item["comment"].(string); ok && v != "" {
			m["comment"] = v
		}
		result = append(result, m)
	}
	return result
}

// buildBudgetCycleQuota converts the cycle_quota schema list into the API request array.
func buildBudgetCycleQuota(raw interface{}) []interface{} {
	result := make([]interface{}, 0)
	if raw == nil {
		return result
	}
	list, ok := raw.([]interface{})
	if !ok {
		return result
	}
	for _, itemRaw := range list {
		item, ok := itemRaw.(map[string]interface{})
		if !ok {
			continue
		}
		m := make(map[string]interface{})
		if v, ok := item["cycle_period"].(string); ok && v != "" {
			m["cyclePeriod"] = v
		}
		if v, ok := item["quota"].(string); ok && v != "" {
			m["quota"] = v
		}
		result = append(result, m)
	}
	return result
}

// buildBudgetQueryFilter converts the query_filter schema list into the API request array.
func buildBudgetQueryFilter(raw interface{}) []interface{} {
	result := make([]interface{}, 0)
	if raw == nil {
		return result
	}
	list, ok := raw.([]interface{})
	if !ok {
		return result
	}
	for _, itemRaw := range list {
		item, ok := itemRaw.(map[string]interface{})
		if !ok {
			continue
		}
		m := make(map[string]interface{})
		if v, ok := item["code"].(string); ok && v != "" {
			m["code"] = v
		}
		if v, ok := item["select_type"].(string); ok && v != "" {
			m["selectType"] = v
		}
		if v := item["values"]; v != nil {
			if values, ok := v.([]interface{}); ok && len(values) > 0 {
				m["values"] = values
			}
		}
		result = append(result, m)
	}
	return result
}
