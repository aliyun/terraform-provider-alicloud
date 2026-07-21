// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAliCloudVpcTrafficMirrorFilterIngressRule() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudVpcTrafficMirrorFilterIngressRuleCreate,
		Read:   resourceAliCloudVpcTrafficMirrorFilterIngressRuleRead,
		Update: resourceAliCloudVpcTrafficMirrorFilterIngressRuleUpdate,
		Delete: resourceAliCloudVpcTrafficMirrorFilterIngressRuleDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"action": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ExactlyOneOf: []string{"action", "rule_action"},
				ValidateFunc: StringInSlice([]string{"accept", "drop"}, false),
			},
			"destination_cidr_block": {
				Type:     schema.TypeString,
				Required: true,
			},
			"destination_port_range": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					if v, ok := d.GetOk("protocol"); ok && v.(string) == "ICMP" {
						return true
					}
					return false
				},
			},
			"dry_run": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"priority": {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: IntBetween(1, 10),
			},
			"protocol": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: StringInSlice([]string{"ALL", "ICMP", "TCP", "UDP"}, false),
			},
			"source_cidr_block": {
				Type:     schema.TypeString,
				Required: true,
			},
			"source_port_range": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					if v, ok := d.GetOk("protocol"); ok && v.(string) == "ICMP" {
						return true
					}
					return false
				},
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"traffic_mirror_filter_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"traffic_mirror_filter_ingress_rule_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"rule_action": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				Deprecated:   "Field 'rule_action' has been deprecated since provider version 1.211.0. New field 'action' instead.",
				ValidateFunc: StringInSlice([]string{"accept", "drop"}, false),
			},
		},
	}
}

func resourceAliCloudVpcTrafficMirrorFilterIngressRuleCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "CreateTrafficMirrorFilterRules"
	var request map[string]interface{}
	var response map[string]interface{}
	var err error
	request = make(map[string]interface{})
	request["TrafficMirrorFilterId"] = d.Get("traffic_mirror_filter_id")
	request["RegionId"] = client.RegionId
	request["ClientToken"] = buildClientToken(action)

	requestIngressRules := make(map[string]interface{})
	requestIngressRulesMap := make([]interface{}, 0)
	if v, ok := d.GetOk("rule_action"); ok {
		requestIngressRules["Action"] = v
	}
	if v, ok := d.GetOk("action"); ok {
		requestIngressRules["Action"] = v
	}
	requestIngressRules["DestinationCidrBlock"] = d.Get("destination_cidr_block")
	requestIngressRules["Priority"] = d.Get("priority")
	requestIngressRules["Protocol"] = d.Get("protocol")
	requestIngressRules["SourceCidrBlock"] = d.Get("source_cidr_block")
	if fmt.Sprint(d.Get("protocol")) != "ICMP" {
		if v, ok := d.GetOk("source_port_range"); ok {
			requestIngressRules["SourcePortRange"] = v
		}
		if v, ok := d.GetOk("destination_port_range"); ok {
			requestIngressRules["DestinationPortRange"] = v
		}
	}
	requestIngressRulesMap = append(requestIngressRulesMap, requestIngressRules)
	request["IngressRules"] = requestIngressRulesMap

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.RpcPost("Vpc", "2016-04-28", action, nil, request, true)
		request["ClientToken"] = buildClientToken(action)

		if err != nil {
			if IsExpectedErrors(err, []string{"IncorrectStatus.TrafficMirrorSession", "OperationConflict", "IncorrectStatus", "SystemBusy", "LastTokenProcessing", "OperationFailed.LastTokenProcessing", "ServiceUnavailable", "IncorrectStatus.TrafficMirrorFilter"}) || NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(action, response, request)
		return nil
	})

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_vpc_traffic_mirror_filter_ingress_rule", action, AlibabaCloudSdkGoERROR)
	}

	v, err := jsonpath.Get("$.IngressRules", response)
	if err != nil || len(v.([]interface{})) < 1 {
		return WrapErrorf(err, IdMsg, d.Id())
	}
	response = v.([]interface{})[0].(map[string]interface{})
	d.SetId(fmt.Sprint(request["TrafficMirrorFilterId"], ":", response["InstanceId"]))

	vpcServiceV2 := VpcServiceV2{client}
	stateConf := BuildStateConf([]string{}, []string{"Created"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, vpcServiceV2.VpcTrafficMirrorFilterIngressRuleStateRefreshFunc(d.Id(), "TrafficMirrorFilterRuleStatus", []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return resourceAliCloudVpcTrafficMirrorFilterIngressRuleUpdate(d, meta)
}

func resourceAliCloudVpcTrafficMirrorFilterIngressRuleRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	vpcServiceV2 := VpcServiceV2{client}

	objectRaw, err := vpcServiceV2.DescribeVpcTrafficMirrorFilterIngressRule(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_vpc_traffic_mirror_filter_ingress_rule DescribeVpcTrafficMirrorFilterIngressRule Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("action", objectRaw["Action"])
	d.Set("destination_cidr_block", objectRaw["DestinationCidrBlock"])
	d.Set("destination_port_range", objectRaw["DestinationPortRange"])
	d.Set("priority", objectRaw["Priority"])
	d.Set("protocol", objectRaw["Protocol"])
	d.Set("source_cidr_block", objectRaw["SourceCidrBlock"])
	d.Set("source_port_range", objectRaw["SourcePortRange"])
	d.Set("status", objectRaw["TrafficMirrorFilterRuleStatus"])
	d.Set("traffic_mirror_filter_id", objectRaw["TrafficMirrorFilterId"])
	d.Set("traffic_mirror_filter_ingress_rule_id", objectRaw["TrafficMirrorFilterRuleId"])

	d.Set("rule_action", d.Get("action"))
	return nil
}

func resourceAliCloudVpcTrafficMirrorFilterIngressRuleUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	update := false
	parts := strings.Split(d.Id(), ":")
	action := "UpdateTrafficMirrorFilterRuleAttribute"
	var err error
	request = make(map[string]interface{})
	request["TrafficMirrorFilterRuleId"] = parts[1]
	request["RegionId"] = client.RegionId
	request["ClientToken"] = buildClientToken(action)
	if !d.IsNewResource() && d.HasChange("priority") {
		update = true
	}
	request["Priority"] = d.Get("priority")
	if !d.IsNewResource() && d.HasChange("protocol") {
		update = true
	}
	request["Protocol"] = d.Get("protocol")
	if !d.IsNewResource() && d.HasChange("destination_cidr_block") {
		update = true
	}
	request["DestinationCidrBlock"] = d.Get("destination_cidr_block")
	if !d.IsNewResource() && d.HasChange("source_cidr_block") {
		update = true
	}
	request["SourceCidrBlock"] = d.Get("source_cidr_block")
	if !d.IsNewResource() && d.HasChange("destination_port_range") {
		update = true
		request["DestinationPortRange"] = d.Get("destination_port_range")
	}

	if !d.IsNewResource() && d.HasChange("source_port_range") {
		update = true
		request["SourcePortRange"] = d.Get("source_port_range")
	}

	if !d.IsNewResource() && d.HasChange("rule_action") {
		update = true
		request["RuleAction"] = d.Get("rule_action")
	}

	if !d.IsNewResource() && d.HasChange("action") {
		update = true
		request["RuleAction"] = d.Get("action")
	}

	if v, ok := d.GetOkExists("dry_run"); ok {
		request["DryRun"] = v
	}

	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("Vpc", "2016-04-28", action, nil, request, true)
			request["ClientToken"] = buildClientToken(action)

			if err != nil {
				if IsExpectedErrors(err, []string{"IncorrectStatus.TrafficMirrorSession", "OperationConflict", "IncorrectStatus", "SystemBusy", "LastTokenProcessing", "OperationFailed.LastTokenProcessing", "ServiceUnavailable"}) || NeedRetry(err) {
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
		vpcServiceV2 := VpcServiceV2{client}
		stateConf := BuildStateConf([]string{}, []string{"Created"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, vpcServiceV2.VpcTrafficMirrorFilterIngressRuleStateRefreshFunc(d.Id(), "TrafficMirrorFilterRuleStatus", []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
	}

	return resourceAliCloudVpcTrafficMirrorFilterIngressRuleRead(d, meta)
}

func resourceAliCloudVpcTrafficMirrorFilterIngressRuleDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	parts := strings.Split(d.Id(), ":")
	action := "DeleteTrafficMirrorFilterRules"
	var request map[string]interface{}
	var response map[string]interface{}
	var err error
	request = make(map[string]interface{})
	request["TrafficMirrorFilterRuleIds.1"] = parts[1]
	request["TrafficMirrorFilterId"] = parts[0]
	request["RegionId"] = client.RegionId

	request["ClientToken"] = buildClientToken(action)
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = client.RpcPost("Vpc", "2016-04-28", action, nil, request, true)
		request["ClientToken"] = buildClientToken(action)

		if err != nil {
			if IsExpectedErrors(err, []string{"IncorrectStatus.TrafficMirrorSession", "OperationConflict", "IncorrectStatus", "SystemBusy", "LastTokenProcessing", "OperationFailed.LastTokenProcessing", "ServiceUnavailable", "IncorrectStatus.TrafficMirrorRule", "IncorrectStatus.TrafficMirrorFilter"}) || NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(action, response, request)
		return nil
	})

	if err != nil {
		if NotFoundError(err) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	return nil
}
