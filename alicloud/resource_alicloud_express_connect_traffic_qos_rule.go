// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAliCloudExpressConnectTrafficQosRule() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudExpressConnectTrafficQosRuleCreate,
		Read:   resourceAliCloudExpressConnectTrafficQosRuleRead,
		Update: resourceAliCloudExpressConnectTrafficQosRuleUpdate,
		Delete: resourceAliCloudExpressConnectTrafficQosRuleDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"dst_cidr": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"dst_ipv6_cidr": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"dst_port_range": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"match_dscp": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: IntBetween(-1, 63),
			},
			"priority": {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: IntBetween(0, 9000),
			},
			"protocol": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: StringInSlice([]string{"ALL", "ICMP(IPv4)", "ICMPv6(IPv6)", "TCP", "UDP", "GRE", "SSH", "Telnet", "HTTP", "HTTPS", "MS SQL", "Oracle", "MySql", "RDP", "Postgre SQL", "Redis"}, false),
			},
			"qos_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"queue_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"remarking_dscp": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: IntBetween(-1, 63),
			},
			"rule_description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"rule_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"rule_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"src_cidr": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"src_ipv6_cidr": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"src_port_range": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceAliCloudExpressConnectTrafficQosRuleCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "CreateExpressConnectTrafficQosRule"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	query["QosId"] = d.Get("qos_id")
	query["QueueId"] = d.Get("queue_id")
	request["RegionId"] = client.RegionId
	request["ClientToken"] = buildClientToken(action)

	if v, ok := d.GetOk("src_port_range"); ok {
		request["SrcPortRange"] = v
	}
	request["Protocol"] = d.Get("protocol")
	if v, ok := d.GetOk("match_dscp"); ok {
		request["MatchDscp"] = v
	}
	if v, ok := d.GetOk("rule_name"); ok {
		request["RuleName"] = v
	}
	if v, ok := d.GetOk("remarking_dscp"); ok {
		request["RemarkingDscp"] = v
	}
	if v, ok := d.GetOk("rule_description"); ok {
		request["RuleDescription"] = v
	}
	request["Priority"] = d.Get("priority")
	if v, ok := d.GetOk("dst_cidr"); ok {
		request["DstCidr"] = v
	}
	if v, ok := d.GetOk("src_cidr"); ok {
		request["SrcCidr"] = v
	}
	if v, ok := d.GetOk("dst_port_range"); ok {
		request["DstPortRange"] = v
	}
	if v, ok := d.GetOk("src_ipv6_cidr"); ok {
		request["SrcIPv6Cidr"] = v
	}
	if v, ok := d.GetOk("dst_ipv6_cidr"); ok {
		request["DstIPv6Cidr"] = v
	}
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.RpcPost("Vpc", "2016-04-28", action, query, request, true)
		if err != nil {
			if NeedRetry(err) || IsExpectedErrors(err, []string{"EcQoSConflict", "IncorrectStatus.Qos"}) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(action, response, request)
		return nil
	})

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_express_connect_traffic_qos_rule", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprintf("%v:%v:%v", response["QosId"], response["QueueId"], response["RuleId"]))

	expressConnectServiceV2 := ExpressConnectServiceV2{client}
	stateConf := BuildStateConf([]string{}, []string{"Normal"}, d.Timeout(schema.TimeoutCreate), 0, expressConnectServiceV2.DescribeAsyncExpressConnectTrafficQosRuleStateRefreshFunc(d, response, "$.QosList[0].Status", []string{}))
	if jobDetail, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id(), jobDetail)
	}

	return resourceAliCloudExpressConnectTrafficQosRuleRead(d, meta)
}

func resourceAliCloudExpressConnectTrafficQosRuleRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	expressConnectServiceV2 := ExpressConnectServiceV2{client}

	objectRaw, err := expressConnectServiceV2.DescribeExpressConnectTrafficQosRule(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_express_connect_traffic_qos_rule DescribeExpressConnectTrafficQosRule Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("dst_cidr", objectRaw["DstCidr"])
	d.Set("dst_ipv6_cidr", objectRaw["DstIPv6Cidr"])
	d.Set("dst_port_range", objectRaw["DstPortRange"])
	d.Set("match_dscp", objectRaw["MatchDscp"])
	d.Set("priority", objectRaw["Priority"])
	d.Set("protocol", objectRaw["Protocol"])
	d.Set("remarking_dscp", objectRaw["RemarkingDscp"])
	d.Set("rule_description", objectRaw["RuleDescription"])
	d.Set("rule_name", objectRaw["RuleName"])
	d.Set("src_cidr", objectRaw["SrcCidr"])
	d.Set("src_ipv6_cidr", objectRaw["SrcIPv6Cidr"])
	d.Set("src_port_range", objectRaw["SrcPortRange"])
	d.Set("status", objectRaw["Status"])
	d.Set("qos_id", objectRaw["QosId"])
	d.Set("queue_id", objectRaw["QueueId"])
	d.Set("rule_id", objectRaw["RuleId"])

	parts := strings.Split(d.Id(), ":")
	d.Set("qos_id", parts[0])
	d.Set("queue_id", parts[1])
	d.Set("rule_id", parts[2])

	return nil
}

func resourceAliCloudExpressConnectTrafficQosRuleUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	update := false
	parts := strings.Split(d.Id(), ":")
	action := "ModifyExpressConnectTrafficQosRule"
	var err error
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	query["QosId"] = parts[0]
	query["RuleId"] = parts[2]
	query["QueueId"] = parts[1]
	request["RegionId"] = client.RegionId
	request["ClientToken"] = buildClientToken(action)
	if d.HasChange("src_port_range") {
		update = true
		request["SrcPortRange"] = d.Get("src_port_range")
	}

	if d.HasChange("protocol") {
		update = true
	}
	request["Protocol"] = d.Get("protocol")
	if d.HasChange("match_dscp") {
		update = true
		request["MatchDscp"] = d.Get("match_dscp")
	}

	if d.HasChange("rule_name") {
		update = true
		request["RuleName"] = d.Get("rule_name")
	}

	if d.HasChange("remarking_dscp") {
		update = true
		request["RemarkingDscp"] = d.Get("remarking_dscp")
	}

	if d.HasChange("rule_description") {
		update = true
		request["RuleDescription"] = d.Get("rule_description")
	}

	if d.HasChange("priority") {
		update = true
	}
	request["Priority"] = d.Get("priority")
	if d.HasChange("dst_cidr") {
		update = true
		request["DstCidr"] = d.Get("dst_cidr")
	}

	if d.HasChange("src_cidr") {
		update = true
		request["SrcCidr"] = d.Get("src_cidr")
	}

	if d.HasChange("dst_port_range") {
		update = true
		request["DstPortRange"] = d.Get("dst_port_range")
	}

	if d.HasChange("src_ipv6_cidr") {
		update = true
		request["SrcIPv6Cidr"] = d.Get("src_ipv6_cidr")
	}

	if d.HasChange("dst_ipv6_cidr") {
		update = true
		request["DstIPv6Cidr"] = d.Get("dst_ipv6_cidr")
	}

	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("Vpc", "2016-04-28", action, query, request, true)
			if err != nil {
				if NeedRetry(err) || IsExpectedErrors(err, []string{"EcQoSConflict", "IncorrectStatus.Qos"}) {
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
		expressConnectServiceV2 := ExpressConnectServiceV2{client}
		stateConf := BuildStateConf([]string{}, []string{"Normal"}, d.Timeout(schema.TimeoutCreate), 0, expressConnectServiceV2.DescribeAsyncExpressConnectTrafficQosRuleStateRefreshFunc(d, response, "$.QosList[0].Status", []string{}))
		if jobDetail, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id(), jobDetail)
		}
	}

	return resourceAliCloudExpressConnectTrafficQosRuleRead(d, meta)
}

func resourceAliCloudExpressConnectTrafficQosRuleDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	parts := strings.Split(d.Id(), ":")
	action := "DeleteExpressConnectTrafficQosRule"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	query["RuleId"] = parts[2]
	query["QueueId"] = parts[1]
	query["QosId"] = parts[0]
	request["RegionId"] = client.RegionId

	request["ClientToken"] = buildClientToken(action)

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = client.RpcPost("Vpc", "2016-04-28", action, query, request, true)
		request["ClientToken"] = buildClientToken(action)

		if err != nil {
			if NeedRetry(err) || IsExpectedErrors(err, []string{"EcQoSConflict", "IncorrectStatus.Qos"}) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(action, response, request)
		return nil
	})

	if err != nil {
		if IsExpectedErrors(err, []string{"IllegalParam.%s"}) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	expressConnectServiceV2 := ExpressConnectServiceV2{client}
	stateConf := BuildStateConf([]string{}, []string{"Normal"}, d.Timeout(schema.TimeoutCreate), 0, expressConnectServiceV2.DescribeAsyncExpressConnectTrafficQosRuleStateRefreshFunc(d, response, "$.QosList[0].Status", []string{}))
	if jobDetail, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id(), jobDetail)
	}
	return nil
}
