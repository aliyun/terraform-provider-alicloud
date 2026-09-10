// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"log"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAliCloudVpcTrafficMirrorFilter() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudVpcTrafficMirrorFilterCreate,
		Read:   resourceAliCloudVpcTrafficMirrorFilterRead,
		Update: resourceAliCloudVpcTrafficMirrorFilterUpdate,
		Delete: resourceAliCloudVpcTrafficMirrorFilterDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"dry_run": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"egress_rules": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"destination_port_range": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
						"action": {
							Type:         schema.TypeString,
							Required:     true,
							ForceNew:     true,
							ValidateFunc: StringInSlice([]string{"accept", "drop"}, false),
						},
						"source_port_range": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
						"priority": {
							Type:     schema.TypeInt,
							Optional: true,
							ForceNew: true,
						},
						"source_cidr_block": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
						"traffic_mirror_filter_rule_status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"destination_cidr_block": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
						"protocol": {
							Type:         schema.TypeString,
							Required:     true,
							ForceNew:     true,
							ValidateFunc: StringInSlice([]string{"ALL", "ICMP", "TCP", "UDP"}, false),
						},
					},
				},
			},
			"ingress_rules": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"destination_port_range": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
						"action": {
							Type:         schema.TypeString,
							Required:     true,
							ForceNew:     true,
							ValidateFunc: StringInSlice([]string{"accept", "drop"}, false),
						},
						"source_port_range": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
						"priority": {
							Type:     schema.TypeInt,
							Optional: true,
							ForceNew: true,
						},
						"source_cidr_block": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
						"traffic_mirror_filter_rule_status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"destination_cidr_block": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
						"protocol": {
							Type:         schema.TypeString,
							Required:     true,
							ForceNew:     true,
							ValidateFunc: StringInSlice([]string{"ALL", "ICMP", "TCP", "UDP"}, false),
						},
					},
				},
			},
			"resource_group_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"tags": tagsSchema(),
			"traffic_mirror_filter_description": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: StringLenBetween(1, 256),
			},
			"traffic_mirror_filter_name": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: StringLenBetween(1, 128),
			},
		},
	}
}

func resourceAliCloudVpcTrafficMirrorFilterCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	action := "CreateTrafficMirrorFilter"
	var request map[string]interface{}
	var response map[string]interface{}
	var err error
	request = make(map[string]interface{})
	request["RegionId"] = client.RegionId
	request["ClientToken"] = buildClientToken(action)

	if v, ok := d.GetOk("traffic_mirror_filter_description"); ok {
		request["TrafficMirrorFilterDescription"] = v
	}
	if v, ok := d.GetOk("traffic_mirror_filter_name"); ok {
		request["TrafficMirrorFilterName"] = v
	}
	if v, ok := d.GetOk("ingress_rules"); ok {
		ingressRulesMaps := make([]map[string]interface{}, 0)
		for _, dataLoop := range v.([]interface{}) {
			dataLoopTmp := dataLoop.(map[string]interface{})
			dataLoopMap := make(map[string]interface{})
			dataLoopMap["Action"] = dataLoopTmp["action"]
			dataLoopMap["SourceCidrBlock"] = dataLoopTmp["source_cidr_block"]
			dataLoopMap["Protocol"] = dataLoopTmp["protocol"]
			dataLoopMap["DestinationPortRange"] = dataLoopTmp["destination_port_range"]
			dataLoopMap["Priority"] = dataLoopTmp["priority"]
			dataLoopMap["DestinationCidrBlock"] = dataLoopTmp["destination_cidr_block"]
			dataLoopMap["SourcePortRange"] = dataLoopTmp["source_port_range"]
			ingressRulesMaps = append(ingressRulesMaps, dataLoopMap)
		}
		request["IngressRules"] = ingressRulesMaps
	}

	if v, ok := d.GetOk("egress_rules"); ok {
		egressRulesMaps := make([]map[string]interface{}, 0)
		for _, dataLoop1 := range v.([]interface{}) {
			dataLoop1Tmp := dataLoop1.(map[string]interface{})
			dataLoop1Map := make(map[string]interface{})
			dataLoop1Map["Action"] = dataLoop1Tmp["action"]
			dataLoop1Map["SourceCidrBlock"] = dataLoop1Tmp["source_cidr_block"]
			dataLoop1Map["Protocol"] = dataLoop1Tmp["protocol"]
			dataLoop1Map["DestinationPortRange"] = dataLoop1Tmp["destination_port_range"]
			dataLoop1Map["Priority"] = dataLoop1Tmp["priority"]
			dataLoop1Map["DestinationCidrBlock"] = dataLoop1Tmp["destination_cidr_block"]
			dataLoop1Map["SourcePortRange"] = dataLoop1Tmp["source_port_range"]
			egressRulesMaps = append(egressRulesMaps, dataLoop1Map)
		}
		request["EgressRules"] = egressRulesMaps
	}

	if v, ok := d.GetOk("resource_group_id"); ok {
		request["ResourceGroupId"] = v
	}
	if v, ok := d.GetOkExists("dry_run"); ok {
		request["DryRun"] = v
	}
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.RpcPost("Vpc", "2016-04-28", action, nil, request, true)
		request["ClientToken"] = buildClientToken(action)

		if err != nil {
			if IsExpectedErrors(err, []string{"OperationConflict", "IncorrectStatus", "ServiceUnavailable", "SystemBusy"}) || NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(action, response, request)
		return nil
	})

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_vpc_traffic_mirror_filter", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(response["TrafficMirrorFilterId"]))

	vpcServiceV2 := VpcServiceV2{client}
	stateConf := BuildStateConf([]string{}, []string{"Created"}, d.Timeout(schema.TimeoutCreate), 0, vpcServiceV2.VpcTrafficMirrorFilterStateRefreshFunc(d.Id(), "TrafficMirrorFilterStatus", []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return resourceAliCloudVpcTrafficMirrorFilterUpdate(d, meta)
}

func resourceAliCloudVpcTrafficMirrorFilterRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	vpcServiceV2 := VpcServiceV2{client}

	objectRaw, err := vpcServiceV2.DescribeVpcTrafficMirrorFilter(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_vpc_traffic_mirror_filter DescribeVpcTrafficMirrorFilter Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("resource_group_id", objectRaw["ResourceGroupId"])
	d.Set("status", objectRaw["TrafficMirrorFilterStatus"])
	d.Set("traffic_mirror_filter_description", objectRaw["TrafficMirrorFilterDescription"])
	d.Set("traffic_mirror_filter_name", objectRaw["TrafficMirrorFilterName"])
	egressRules1Raw := objectRaw["EgressRules"]
	egressRulesMaps := make([]map[string]interface{}, 0)
	if egressRules1Raw != nil {
		for _, egressRulesChild1Raw := range egressRules1Raw.([]interface{}) {
			egressRulesMap := make(map[string]interface{})
			egressRulesChild1Raw := egressRulesChild1Raw.(map[string]interface{})
			egressRulesMap["action"] = egressRulesChild1Raw["Action"]
			egressRulesMap["destination_cidr_block"] = egressRulesChild1Raw["DestinationCidrBlock"]
			egressRulesMap["destination_port_range"] = egressRulesChild1Raw["DestinationPortRange"]
			egressRulesMap["priority"] = egressRulesChild1Raw["Priority"]
			egressRulesMap["protocol"] = egressRulesChild1Raw["Protocol"]
			egressRulesMap["source_cidr_block"] = egressRulesChild1Raw["SourceCidrBlock"]
			egressRulesMap["source_port_range"] = egressRulesChild1Raw["SourcePortRange"]
			egressRulesMap["traffic_mirror_filter_rule_status"] = egressRulesChild1Raw["TrafficMirrorFilterRuleStatus"]
			egressRulesMaps = append(egressRulesMaps, egressRulesMap)
		}
	}
	d.Set("egress_rules", egressRulesMaps)
	ingressRules1Raw := objectRaw["IngressRules"]
	ingressRulesMaps := make([]map[string]interface{}, 0)
	if ingressRules1Raw != nil {
		for _, ingressRulesChild1Raw := range ingressRules1Raw.([]interface{}) {
			ingressRulesMap := make(map[string]interface{})
			ingressRulesChild1Raw := ingressRulesChild1Raw.(map[string]interface{})
			ingressRulesMap["action"] = ingressRulesChild1Raw["Action"]
			ingressRulesMap["destination_cidr_block"] = ingressRulesChild1Raw["DestinationCidrBlock"]
			ingressRulesMap["destination_port_range"] = ingressRulesChild1Raw["DestinationPortRange"]
			ingressRulesMap["priority"] = ingressRulesChild1Raw["Priority"]
			ingressRulesMap["protocol"] = ingressRulesChild1Raw["Protocol"]
			ingressRulesMap["source_cidr_block"] = ingressRulesChild1Raw["SourceCidrBlock"]
			ingressRulesMap["source_port_range"] = ingressRulesChild1Raw["SourcePortRange"]
			ingressRulesMap["traffic_mirror_filter_rule_status"] = ingressRulesChild1Raw["TrafficMirrorFilterRuleStatus"]
			ingressRulesMaps = append(ingressRulesMaps, ingressRulesMap)
		}
	}
	d.Set("ingress_rules", ingressRulesMaps)
	tagsMaps := objectRaw["Tags"]
	d.Set("tags", tagsToMap(tagsMaps))

	return nil
}

func resourceAliCloudVpcTrafficMirrorFilterUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	update := false
	d.Partial(true)
	action := "UpdateTrafficMirrorFilterAttribute"
	var err error
	request = make(map[string]interface{})

	request["TrafficMirrorFilterId"] = d.Id()
	request["RegionId"] = client.RegionId
	request["ClientToken"] = buildClientToken(action)
	if !d.IsNewResource() && d.HasChange("traffic_mirror_filter_description") {
		update = true
		request["TrafficMirrorFilterDescription"] = d.Get("traffic_mirror_filter_description")
	}

	if !d.IsNewResource() && d.HasChange("traffic_mirror_filter_name") {
		update = true
		request["TrafficMirrorFilterName"] = d.Get("traffic_mirror_filter_name")
	}

	if v, ok := d.GetOkExists("dry_run"); ok {
		request["DryRun"] = v
	}
	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("Vpc", "2016-04-28", action, nil, request, false)
			request["ClientToken"] = buildClientToken(action)

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
		d.SetPartial("traffic_mirror_filter_description")
		d.SetPartial("traffic_mirror_filter_name")
	}
	update = false
	action = "MoveResourceGroup"
	request = make(map[string]interface{})

	request["ResourceId"] = d.Id()
	request["RegionId"] = client.RegionId
	if !d.IsNewResource() && d.HasChange("resource_group_id") {
		update = true
		request["NewResourceGroupId"] = d.Get("resource_group_id")
	}

	request["ResourceType"] = "TRAFFICMIRRORFILTER"
	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("Vpc", "2016-04-28", action, nil, request, false)

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
		d.SetPartial("resource_group_id")
	}

	update = false
	if d.HasChange("tags") {
		update = true
		vpcServiceV2 := VpcServiceV2{client}
		if err := vpcServiceV2.SetResourceTags(d, "TRAFFICMIRRORFILTER"); err != nil {
			return WrapError(err)
		}
		d.SetPartial("tags")
	}
	d.Partial(false)
	return resourceAliCloudVpcTrafficMirrorFilterRead(d, meta)
}

func resourceAliCloudVpcTrafficMirrorFilterDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "DeleteTrafficMirrorFilter"
	var request map[string]interface{}
	var response map[string]interface{}
	var err error
	request = make(map[string]interface{})

	request["TrafficMirrorFilterId"] = d.Id()
	request["RegionId"] = client.RegionId

	request["ClientToken"] = buildClientToken(action)

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = client.RpcPost("Vpc", "2016-04-28", action, nil, request, true)
		request["ClientToken"] = buildClientToken(action)

		if err != nil {
			if IsExpectedErrors(err, []string{"IncorrectStatus.TrafficMirrorFilter", "IncorrectStatus.TrafficMirrorRule", "OperationConflict", "IncorrectStatus", "ServiceUnavailable", "SystemBusy"}) || NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(action, response, request)
		return nil
	})

	if err != nil {
		if IsExpectedErrors(err, []string{"ResourceNotFound.TrafficMirrorFilter"}) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	vpcServiceV2 := VpcServiceV2{client}
	stateConf := BuildStateConf([]string{}, []string{}, d.Timeout(schema.TimeoutDelete), 0, vpcServiceV2.VpcTrafficMirrorFilterStateRefreshFunc(d.Id(), "TrafficMirrorFilterStatus", []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}
	return nil
}
