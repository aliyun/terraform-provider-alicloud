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

func resourceAliCloudCloudFirewallVpcCenTrFirewallPolicy() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudCloudFirewallVpcCenTrFirewallPolicyCreate,
		Read:   resourceAliCloudCloudFirewallVpcCenTrFirewallPolicyRead,
		Update: resourceAliCloudCloudFirewallVpcCenTrFirewallPolicyUpdate,
		Delete: resourceAliCloudCloudFirewallVpcCenTrFirewallPolicyDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(16 * time.Minute),
			Update: schema.DefaultTimeout(16 * time.Minute),
			Delete: schema.DefaultTimeout(16 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"dest_candidate_list": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"candidate_type": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"candidate_id": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"firewall_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"policy_description": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"policy_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"policy_type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: StringInSlice([]string{"fullmesh", "one_to_one", "end_to_end"}, false),
			},
			"should_recover": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: StringInSlice([]string{"true", "false"}, false),
			},
			"src_candidate_list": {
				Type:     schema.TypeSet,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"candidate_type": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"candidate_id": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"status": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: StringInSlice([]string{"creating", "deleting", "opening", "opened", "closing", "closed", "open", "close"}, false),
			},
			"tr_firewall_route_policy_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceAliCloudCloudFirewallVpcCenTrFirewallPolicyCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "CreateTrFirewallV2RoutePolicy"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	if v, ok := d.GetOk("firewall_id"); ok {
		request["FirewallId"] = v
	}

	if v, ok := d.GetOk("lang"); ok {
		request["Lang"] = v
	}
	if v, ok := d.GetOk("src_candidate_list"); ok {
		srcCandidateListMapsArray := make([]interface{}, 0)
		for _, dataLoop := range v.(*schema.Set).List() {
			dataLoopTmp := dataLoop.(map[string]interface{})
			dataLoopMap := make(map[string]interface{})
			dataLoopMap["CandidateType"] = dataLoopTmp["candidate_type"]
			dataLoopMap["CandidateId"] = dataLoopTmp["candidate_id"]
			srcCandidateListMapsArray = append(srcCandidateListMapsArray, dataLoopMap)
		}
		srcCandidateListMapsJson, err := json.Marshal(srcCandidateListMapsArray)
		if err != nil {
			return WrapError(err)
		}
		request["SrcCandidateList"] = string(srcCandidateListMapsJson)
	}

	if v, ok := d.GetOk("dest_candidate_list"); ok {
		destCandidateListMapsArray := make([]interface{}, 0)
		for _, dataLoop1 := range v.(*schema.Set).List() {
			dataLoop1Tmp := dataLoop1.(map[string]interface{})
			dataLoop1Map := make(map[string]interface{})
			dataLoop1Map["CandidateType"] = dataLoop1Tmp["candidate_type"]
			dataLoop1Map["CandidateId"] = dataLoop1Tmp["candidate_id"]
			destCandidateListMapsArray = append(destCandidateListMapsArray, dataLoop1Map)
		}
		destCandidateListMapsJson, err := json.Marshal(destCandidateListMapsArray)
		if err != nil {
			return WrapError(err)
		}
		request["DestCandidateList"] = string(destCandidateListMapsJson)
	}

	request["PolicyType"] = d.Get("policy_type")
	request["PolicyName"] = d.Get("policy_name")
	request["PolicyDescription"] = d.Get("policy_description")
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.RpcPost("Cloudfw", "2017-12-07", action, query, request, true)
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_cloud_firewall_vpc_cen_tr_firewall_policy", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprintf("%v:%v", request["FirewallId"], response["TrFirewallRoutePolicyId"]))

	cloudFirewallServiceV2 := CloudFirewallServiceV2{client}
	stateConf := BuildStateConf([]string{}, []string{"opened"}, d.Timeout(schema.TimeoutCreate), 30*time.Second, cloudFirewallServiceV2.CloudFirewallVpcCenTrFirewallPolicyStateRefreshFunc(d.Id(), "PolicyStatus", []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return resourceAliCloudCloudFirewallVpcCenTrFirewallPolicyUpdate(d, meta)
}

func resourceAliCloudCloudFirewallVpcCenTrFirewallPolicyRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	cloudFirewallServiceV2 := CloudFirewallServiceV2{client}

	objectRaw, err := cloudFirewallServiceV2.DescribeCloudFirewallVpcCenTrFirewallPolicy(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_cloud_firewall_vpc_cen_tr_firewall_policy DescribeCloudFirewallVpcCenTrFirewallPolicy Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("policy_description", objectRaw["PolicyDescription"])
	d.Set("policy_name", objectRaw["PolicyName"])
	d.Set("policy_type", objectRaw["PolicyType"])
	d.Set("status", objectRaw["PolicyStatus"])
	d.Set("tr_firewall_route_policy_id", objectRaw["TrFirewallRoutePolicyId"])

	destCandidateListRaw := objectRaw["DestCandidateList"]
	destCandidateListMaps := make([]map[string]interface{}, 0)
	if destCandidateListRaw != nil {
		for _, destCandidateListChildRaw := range destCandidateListRaw.([]interface{}) {
			destCandidateListMap := make(map[string]interface{})
			destCandidateListChildRaw := destCandidateListChildRaw.(map[string]interface{})
			destCandidateListMap["candidate_id"] = destCandidateListChildRaw["CandidateId"]
			destCandidateListMap["candidate_type"] = destCandidateListChildRaw["CandidateType"]

			destCandidateListMaps = append(destCandidateListMaps, destCandidateListMap)
		}
	}
	if err := d.Set("dest_candidate_list", destCandidateListMaps); err != nil {
		return err
	}
	srcCandidateListRaw := objectRaw["SrcCandidateList"]
	srcCandidateListMaps := make([]map[string]interface{}, 0)
	if srcCandidateListRaw != nil {
		for _, srcCandidateListChildRaw := range srcCandidateListRaw.([]interface{}) {
			srcCandidateListMap := make(map[string]interface{})
			srcCandidateListChildRaw := srcCandidateListChildRaw.(map[string]interface{})
			srcCandidateListMap["candidate_id"] = srcCandidateListChildRaw["CandidateId"]
			srcCandidateListMap["candidate_type"] = srcCandidateListChildRaw["CandidateType"]

			srcCandidateListMaps = append(srcCandidateListMaps, srcCandidateListMap)
		}
	}
	if err := d.Set("src_candidate_list", srcCandidateListMaps); err != nil {
		return err
	}

	parts := strings.Split(d.Id(), ":")
	d.Set("firewall_id", parts[0])

	return nil
}

func resourceAliCloudCloudFirewallVpcCenTrFirewallPolicyUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	update := false
	d.Partial(true)

	var err error
	parts := strings.Split(d.Id(), ":")
	action := "ModifyTrFirewallV2RoutePolicyScope"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["FirewallId"] = parts[0]
	request["TrFirewallRoutePolicyId"] = parts[1]

	if v, ok := d.GetOk("lang"); ok {
		request["Lang"] = v
	}
	if !d.IsNewResource() && d.HasChange("src_candidate_list") {
		update = true
	}
	if v, ok := d.GetOk("src_candidate_list"); ok || d.HasChange("src_candidate_list") {
		srcCandidateListMapsArray := make([]interface{}, 0)
		for _, dataLoop := range v.(*schema.Set).List() {
			dataLoopTmp := dataLoop.(map[string]interface{})
			dataLoopMap := make(map[string]interface{})
			dataLoopMap["CandidateType"] = dataLoopTmp["candidate_type"]
			dataLoopMap["CandidateId"] = dataLoopTmp["candidate_id"]
			srcCandidateListMapsArray = append(srcCandidateListMapsArray, dataLoopMap)
		}
		srcCandidateListMapsJson, err := json.Marshal(srcCandidateListMapsArray)
		if err != nil {
			return WrapError(err)
		}
		request["SrcCandidateList"] = string(srcCandidateListMapsJson)
	}

	if !d.IsNewResource() && d.HasChange("dest_candidate_list") {
		update = true
	}
	if v, ok := d.GetOk("dest_candidate_list"); ok || d.HasChange("dest_candidate_list") {
		destCandidateListMapsArray := make([]interface{}, 0)
		for _, dataLoop1 := range v.(*schema.Set).List() {
			dataLoop1Tmp := dataLoop1.(map[string]interface{})
			dataLoop1Map := make(map[string]interface{})
			dataLoop1Map["CandidateType"] = dataLoop1Tmp["candidate_type"]
			dataLoop1Map["CandidateId"] = dataLoop1Tmp["candidate_id"]
			destCandidateListMapsArray = append(destCandidateListMapsArray, dataLoop1Map)
		}
		destCandidateListMapsJson, err := json.Marshal(destCandidateListMapsArray)
		if err != nil {
			return WrapError(err)
		}
		request["DestCandidateList"] = string(destCandidateListMapsJson)
	}

	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("Cloudfw", "2017-12-07", action, query, request, true)
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
		cloudFirewallServiceV2 := CloudFirewallServiceV2{client}
		stateConf := BuildStateConf([]string{}, []string{"opened", "closed"}, d.Timeout(schema.TimeoutUpdate), 30*time.Second, cloudFirewallServiceV2.CloudFirewallVpcCenTrFirewallPolicyStateRefreshFunc(d.Id(), "PolicyStatus", []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
	}
	update = false
	parts = strings.Split(d.Id(), ":")
	action = "ModifyFirewallV2RoutePolicySwitch"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["FirewallId"] = parts[0]
	request["TrFirewallRoutePolicyId"] = parts[1]

	if v, ok := d.GetOk("lang"); ok {
		request["Lang"] = v
	}
	if d.HasChange("status") {
		update = true
	}
	if v, ok := d.GetOk("status"); ok || d.HasChange("status") {
		request["TrFirewallRoutePolicySwitchStatus"] = convertCloudFirewallVpcCenTrFirewallPolicyTrFirewallRoutePolicySwitchStatusRequest(v.(string))
	}
	if v, ok := d.GetOk("should_recover"); ok {
		request["ShouldRecover"] = v
	}
	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("Cloudfw", "2017-12-07", action, query, request, true)
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
		cloudFirewallServiceV2 := CloudFirewallServiceV2{client}
		stateConf := BuildStateConf([]string{}, []string{"opened", "closed"}, d.Timeout(schema.TimeoutUpdate), 30*time.Second, cloudFirewallServiceV2.CloudFirewallVpcCenTrFirewallPolicyStateRefreshFunc(d.Id(), "PolicyStatus", []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
	}

	d.Partial(false)
	return resourceAliCloudCloudFirewallVpcCenTrFirewallPolicyRead(d, meta)
}

func resourceAliCloudCloudFirewallVpcCenTrFirewallPolicyDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	parts := strings.Split(d.Id(), ":")
	action := "DeleteFirewallV2RoutePolicies"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	request["FirewallId"] = parts[0]
	request["TrFirewallRoutePolicyId"] = parts[1]

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = client.RpcPost("Cloudfw", "2017-12-07", action, query, request, true)

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

	cloudFirewallServiceV2 := CloudFirewallServiceV2{client}
	stateConf := BuildStateConf([]string{}, []string{}, d.Timeout(schema.TimeoutDelete), 30*time.Second, cloudFirewallServiceV2.CloudFirewallVpcCenTrFirewallPolicyStateRefreshFunc(d.Id(), "PolicyStatus", []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return nil
}

func convertCloudFirewallVpcCenTrFirewallPolicyTrFirewallRoutePolicySwitchStatusRequest(source interface{}) interface{} {
	source = fmt.Sprint(source)
	switch source {
	case "opened":
		return "open"
	case "closed":
		return "close"
	}
	return source
}
