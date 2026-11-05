// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package esa

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	convert "github.com/aliyun/terraform-provider-alicloud/alicloud/convert"
	errs "github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	helper "github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	logs "github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAliCloudEsaTransportLayerApplication() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudEsaTransportLayerApplicationCreate,
		Read:   resourceAliCloudEsaTransportLayerApplicationRead,
		Update: resourceAliCloudEsaTransportLayerApplicationUpdate,
		Delete: resourceAliCloudEsaTransportLayerApplicationDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(25 * time.Minute),
			Update: schema.DefaultTimeout(17 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"application_id": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"cross_border_optimization": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"ip_access_rule": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"ipv6": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"record_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"rules": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"comment": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"edge_port": {
							Type:     schema.TypeString,
							Required: true,
						},
						"source_type": {
							Type:     schema.TypeString,
							Required: true,
						},
						"rule_id": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"protocol": {
							Type:     schema.TypeString,
							Required: true,
						},
						"source_port": {
							Type:     schema.TypeString,
							Required: true,
						},
						"client_ip_pass_through_mode": {
							Type:     schema.TypeString,
							Required: true,
						},
						"source": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
			"site_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceAliCloudEsaTransportLayerApplicationCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "CreateTransportLayerApplication"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	if v, ok := d.GetOk("site_id"); ok {
		request["SiteId"] = v
	}

	if v, ok := d.GetOk("rules"); ok {
		rulesMapsArray := make([]interface{}, 0)
		for _, dataLoop := range convert.ConvertToInterfaceArray(v) {
			dataLoopTmp := dataLoop.(map[string]interface{})
			dataLoopMap := make(map[string]interface{})
			dataLoopMap["EdgePort"] = dataLoopTmp["edge_port"]
			dataLoopMap["SourceType"] = dataLoopTmp["source_type"]
			dataLoopMap["Source"] = dataLoopTmp["source"]
			dataLoopMap["SourcePort"] = dataLoopTmp["source_port"]
			dataLoopMap["Comment"] = dataLoopTmp["comment"]
			dataLoopMap["Protocol"] = dataLoopTmp["protocol"]
			dataLoopMap["ClientIPPassThroughMode"] = dataLoopTmp["client_ip_pass_through_mode"]
			rulesMapsArray = append(rulesMapsArray, dataLoopMap)
		}
		rulesMapsJson, err := json.Marshal(rulesMapsArray)
		if err != nil {
			return errs.WrapError(err)
		}
		request["Rules"] = string(rulesMapsJson)
	}

	if v, ok := d.GetOk("ip_access_rule"); ok {
		request["IpAccessRule"] = v
	}
	if v, ok := d.GetOk("ipv6"); ok {
		request["Ipv6"] = v
	}
	request["RecordName"] = d.Get("record_name")
	if v, ok := d.GetOk("cross_border_optimization"); ok {
		request["CrossBorderOptimization"] = v
	}
	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
		response, err = client.RpcPost("ESA", "2024-09-10", action, query, request, true)
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
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, "alicloud_esa_transport_layer_application", action, errs.AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprintf("%v:%v", request["SiteId"], response["ApplicationId"]))

	esaServiceV2 := EsaServiceV2{client}
	stateConf := helper.BuildStateConf([]string{}, []string{"active"}, d.Timeout(schema.TimeoutCreate), 3*time.Minute, esaServiceV2.EsaTransportLayerApplicationStateRefreshFunc(d.Id(), "Status", []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return errs.WrapErrorf(err, errs.IdMsg, d.Id())
	}

	return resourceAliCloudEsaTransportLayerApplicationRead(d, meta)
}

func resourceAliCloudEsaTransportLayerApplicationRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	esaServiceV2 := EsaServiceV2{client}

	objectRaw, err := esaServiceV2.DescribeEsaTransportLayerApplication(d.Id())
	if err != nil {
		if !d.IsNewResource() && errs.NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_esa_transport_layer_application DescribeEsaTransportLayerApplication Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return errs.WrapError(err)
	}

	d.Set("cross_border_optimization", objectRaw["CrossBorderOptimization"])
	d.Set("ip_access_rule", objectRaw["IpAccessRule"])
	d.Set("ipv6", objectRaw["Ipv6"])
	d.Set("record_name", objectRaw["RecordName"])
	d.Set("status", objectRaw["Status"])
	applicationId, err := strconv.ParseInt(objectRaw["ApplicationId"].(json.Number).String(), 10, 64)
	if err != nil {
		return errs.WrapError(err)
	}
	d.Set("application_id", applicationId)
	if v, ok := objectRaw["SiteId"]; ok {
		d.Set("site_id", v)
	}

	rulesRaw := objectRaw["Rules"]
	rulesMaps := make([]map[string]interface{}, 0)
	if rulesRaw != nil {
		for _, rulesChildRaw := range convert.ConvertToInterfaceArray(rulesRaw) {
			rulesMap := make(map[string]interface{})
			rulesChildRaw := rulesChildRaw.(map[string]interface{})
			rulesMap["client_ip_pass_through_mode"] = rulesChildRaw["ClientIPPassThroughMode"]
			rulesMap["comment"] = rulesChildRaw["Comment"]
			rulesMap["edge_port"] = rulesChildRaw["EdgePort"]
			rulesMap["protocol"] = rulesChildRaw["Protocol"]
			rulesMap["rule_id"] = rulesChildRaw["RuleId"]
			rulesMap["source"] = rulesChildRaw["Source"]
			rulesMap["source_port"] = rulesChildRaw["SourcePort"]
			rulesMap["source_type"] = rulesChildRaw["SourceType"]

			rulesMaps = append(rulesMaps, rulesMap)
		}
	}
	if err := d.Set("rules", rulesMaps); err != nil {
		return err
	}

	return nil
}

func resourceAliCloudEsaTransportLayerApplicationUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	update := false

	var err error
	parts := strings.Split(d.Id(), ":")
	action := "UpdateTransportLayerApplication"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["SiteId"] = parts[0]
	request["ApplicationId"] = parts[1]

	if d.HasChange("rules") {
		update = true
	}
	if v, ok := d.GetOk("rules"); ok || d.HasChange("rules") {
		rulesMapsArray := make([]interface{}, 0)
		for _, dataLoop := range convert.ConvertToInterfaceArray(v) {
			dataLoopTmp := dataLoop.(map[string]interface{})
			dataLoopMap := make(map[string]interface{})
			dataLoopMap["EdgePort"] = dataLoopTmp["edge_port"]
			dataLoopMap["SourceType"] = dataLoopTmp["source_type"]
			dataLoopMap["Source"] = dataLoopTmp["source"]
			dataLoopMap["SourcePort"] = dataLoopTmp["source_port"]
			dataLoopMap["Comment"] = dataLoopTmp["comment"]
			dataLoopMap["Protocol"] = dataLoopTmp["protocol"]
			dataLoopMap["ClientIPPassThroughMode"] = dataLoopTmp["client_ip_pass_through_mode"]
			rulesMapsArray = append(rulesMapsArray, dataLoopMap)
		}
		rulesMapsJson, err := json.Marshal(rulesMapsArray)
		if err != nil {
			return errs.WrapError(err)
		}
		request["Rules"] = string(rulesMapsJson)
	}

	if d.HasChange("ip_access_rule") {
		update = true
		request["IpAccessRule"] = d.Get("ip_access_rule")
	}

	if d.HasChange("ipv6") {
		update = true
		request["Ipv6"] = d.Get("ipv6")
	}

	if d.HasChange("cross_border_optimization") {
		update = true
		request["CrossBorderOptimization"] = d.Get("cross_border_optimization")
	}

	if update {
		wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
		err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
			response, err = client.RpcPost("ESA", "2024-09-10", action, query, request, true)
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
		esaServiceV2 := EsaServiceV2{client}
		stateConf := helper.BuildStateConf([]string{}, []string{"active"}, d.Timeout(schema.TimeoutUpdate), 3*time.Minute, esaServiceV2.EsaTransportLayerApplicationStateRefreshFunc(d.Id(), "Status", []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return errs.WrapErrorf(err, errs.IdMsg, d.Id())
		}
	}

	return resourceAliCloudEsaTransportLayerApplicationRead(d, meta)
}

func resourceAliCloudEsaTransportLayerApplicationDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	parts := strings.Split(d.Id(), ":")
	enableDelete := false
	if v, ok := d.GetOkExists("status"); ok {
		if helper.InArray(fmt.Sprint(v), []string{"active"}) {
			enableDelete = true
		}
	}
	if enableDelete {
		action := "DeleteTransportLayerApplication"
		var request map[string]interface{}
		var response map[string]interface{}
		query := make(map[string]interface{})
		var err error
		request = make(map[string]interface{})
		request["SiteId"] = parts[0]
		request["ApplicationId"] = parts[1]

		wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
		err = retry.Retry(d.Timeout(schema.TimeoutDelete), func() *retry.RetryError {
			response, err = client.RpcPost("ESA", "2024-09-10", action, query, request, true)
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
			if errs.NotFoundError(err) {
				return nil
			}
			return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
		}
		esaServiceV2 := EsaServiceV2{client}
		stateConf := helper.BuildStateConf([]string{}, []string{""}, d.Timeout(schema.TimeoutDelete), 3*time.Minute, esaServiceV2.EsaTransportLayerApplicationStateRefreshFunc(d.Id(), "Status", []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return errs.WrapErrorf(err, errs.IdMsg, d.Id())
		}
	}
	return nil
}
