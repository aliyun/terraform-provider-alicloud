// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
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
			},
			"ip_access_rule": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"ipv6": {
				Type:     schema.TypeString,
				Optional: true,
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
				Type:     schema.TypeInt,
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
		for _, dataLoop := range v.([]interface{}) {
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
			return WrapError(err)
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
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.RpcPost("ESA", "2024-09-10", action, query, request, true)
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_esa_transport_layer_application", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprintf("%v:%v", request["SiteId"], response["ApplicationId"]))

	esaServiceV2 := EsaServiceV2{client}
	stateConf := BuildStateConf([]string{}, []string{"active"}, d.Timeout(schema.TimeoutCreate), 15*time.Minute, esaServiceV2.EsaTransportLayerApplicationStateRefreshFunc(d.Id(), "Status", []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return resourceAliCloudEsaTransportLayerApplicationUpdate(d, meta)
}

func resourceAliCloudEsaTransportLayerApplicationRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	esaServiceV2 := EsaServiceV2{client}

	objectRaw, err := esaServiceV2.DescribeEsaTransportLayerApplication(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_esa_transport_layer_application DescribeEsaTransportLayerApplication Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("cross_border_optimization", objectRaw["CrossBorderOptimization"])
	d.Set("ip_access_rule", objectRaw["IpAccessRule"])
	d.Set("ipv6", objectRaw["Ipv6"])
	d.Set("record_name", objectRaw["RecordName"])
	d.Set("status", objectRaw["Status"])
	d.Set("application_id", objectRaw["ApplicationId"])
	d.Set("site_id", objectRaw["SiteId"])

	rulesRaw := objectRaw["Rules"]
	rulesMaps := make([]map[string]interface{}, 0)
	if rulesRaw != nil {
		for _, rulesChildRaw := range rulesRaw.([]interface{}) {
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

	if !d.IsNewResource() && d.HasChange("rules") {
		update = true
	}
	if v, ok := d.GetOk("rules"); ok || d.HasChange("rules") {
		rulesMapsArray := make([]interface{}, 0)
		for _, dataLoop := range v.([]interface{}) {
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
			return WrapError(err)
		}
		request["Rules"] = string(rulesMapsJson)
	}

	if !d.IsNewResource() && d.HasChange("ip_access_rule") {
		update = true
		request["IpAccessRule"] = d.Get("ip_access_rule")
	}

	if !d.IsNewResource() && d.HasChange("ipv6") {
		update = true
		request["Ipv6"] = d.Get("ipv6")
	}

	if !d.IsNewResource() && d.HasChange("cross_border_optimization") {
		update = true
		request["CrossBorderOptimization"] = d.Get("cross_border_optimization")
	}

	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("ESA", "2024-09-10", action, query, request, true)
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
		esaServiceV2 := EsaServiceV2{client}
		stateConf := BuildStateConf([]string{}, []string{"active"}, d.Timeout(schema.TimeoutUpdate), 10*time.Minute, esaServiceV2.EsaTransportLayerApplicationStateRefreshFunc(d.Id(), "Status", []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
	}

	return resourceAliCloudEsaTransportLayerApplicationRead(d, meta)
}

func resourceAliCloudEsaTransportLayerApplicationDelete(d *schema.ResourceData, meta interface{}) error {

	enableDelete := false
	if v, ok := d.GetOk("status"); ok {
		if InArray(fmt.Sprint(v), []string{"active"}) {
			enableDelete = true
		}
	}
	if enableDelete {
		client := meta.(*connectivity.AliyunClient)
		parts := strings.Split(d.Id(), ":")
		action := "DeleteTransportLayerApplication"
		var request map[string]interface{}
		var response map[string]interface{}
		query := make(map[string]interface{})
		var err error
		request = make(map[string]interface{})
		request["SiteId"] = parts[0]
		request["ApplicationId"] = parts[1]

		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
			response, err = client.RpcPost("ESA", "2024-09-10", action, query, request, true)

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
			if NotFoundError(err) {
				return nil
			}
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}

	}
	return nil
}
