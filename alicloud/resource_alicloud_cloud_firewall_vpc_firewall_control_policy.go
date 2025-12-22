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

func resourceAliCloudCloudFirewallVpcFirewallControlPolicy() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudCloudFirewallVpcFirewallControlPolicyCreate,
		Read:   resourceAliCloudCloudFirewallVpcFirewallControlPolicyRead,
		Update: resourceAliCloudCloudFirewallVpcFirewallControlPolicyUpdate,
		Delete: resourceAliCloudCloudFirewallVpcFirewallControlPolicyDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"acl_action": {
				Type:     schema.TypeString,
				Required: true,
			},
			"acl_uuid": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"application_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"application_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"application_name_list": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"create_time": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Required: true,
			},
			"dest_port": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"dest_port_group": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"dest_port_group_ports": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"dest_port_type": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"destination": {
				Type:     schema.TypeString,
				Required: true,
			},
			"destination_group_cidrs": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"destination_group_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"destination_type": {
				Type:     schema.TypeString,
				Required: true,
			},
			"domain_resolve_type": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"end_time": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"hit_times": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"lang": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: StringInSlice([]string{"zh", "en"}, false),
			},
			"member_uid": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"order": {
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
			},
			"proto": {
				Type:     schema.TypeString,
				Required: true,
			},
			"release": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"repeat_days": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeInt},
			},
			"repeat_end_time": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"repeat_start_time": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"repeat_type": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"source": {
				Type:     schema.TypeString,
				Required: true,
			},
			"source_group_cidrs": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"source_group_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"source_type": {
				Type:     schema.TypeString,
				Required: true,
			},
			"start_time": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"vpc_firewall_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceAliCloudCloudFirewallVpcFirewallControlPolicyCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "CreateVpcFirewallControlPolicy"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	var endpoint string
	if v, ok := d.GetOk("vpc_firewall_id"); ok {
		request["VpcFirewallId"] = v
	}

	if v, ok := d.GetOk("application_name"); ok {
		request["ApplicationName"] = v
	}
	if v, ok := d.GetOk("lang"); ok {
		request["Lang"] = v
	}
	request["Proto"] = d.Get("proto")
	if v, ok := d.GetOkExists("start_time"); ok {
		request["StartTime"] = v
	}
	request["NewOrder"] = d.Get("order")
	if v, ok := d.GetOk("repeat_start_time"); ok {
		request["RepeatStartTime"] = v
	}
	if v, ok := d.GetOk("application_name_list"); ok {
		applicationNameListMapsArray := convertToInterfaceArray(v)

		request["ApplicationNameList"] = applicationNameListMapsArray
	}

	if v, ok := d.GetOk("release"); ok {
		request["Release"] = v
	}
	if v, ok := d.GetOk("domain_resolve_type"); ok {
		request["DomainResolveType"] = v
	}
	if v, ok := d.GetOk("repeat_days"); ok {
		repeatDaysMapsArray := convertToInterfaceArray(v)

		request["RepeatDays"] = repeatDaysMapsArray
	}

	if v, ok := d.GetOk("repeat_type"); ok {
		request["RepeatType"] = v
	}
	if v, ok := d.GetOkExists("end_time"); ok {
		request["EndTime"] = v
	}
	if v, ok := d.GetOk("member_uid"); ok {
		request["MemberUid"] = v
	}
	request["DestinationType"] = d.Get("destination_type")
	request["Description"] = d.Get("description")
	if v, ok := d.GetOk("dest_port"); ok {
		request["DestPort"] = v
	}
	if v, ok := d.GetOk("dest_port_group"); ok {
		request["DestPortGroup"] = v
	}
	request["Destination"] = d.Get("destination")
	if v, ok := d.GetOk("dest_port_type"); ok {
		request["DestPortType"] = v
	}
	if v, ok := d.GetOk("repeat_end_time"); ok {
		request["RepeatEndTime"] = v
	}
	request["AclAction"] = d.Get("acl_action")
	request["SourceType"] = d.Get("source_type")
	request["Source"] = d.Get("source")
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.RpcPostWithEndpoint("Cloudfw", "2017-12-07", action, query, request, true, endpoint)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			} else if IsExpectedErrors(err, []string{"not buy user"}) {
				endpoint = connectivity.CloudFirewallOpenAPIEndpointControlPolicy
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}

		return nil
	})
	addDebug(action, response, request)

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_cloud_firewall_vpc_firewall_control_policy", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprintf("%v:%v", request["VpcFirewallId"], response["AclUuid"]))

	return resourceAliCloudCloudFirewallVpcFirewallControlPolicyRead(d, meta)
}

func resourceAliCloudCloudFirewallVpcFirewallControlPolicyRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	cloudFirewallServiceV2 := CloudFirewallServiceV2{client}

	objectRaw, err := cloudFirewallServiceV2.DescribeCloudFirewallVpcFirewallControlPolicy(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_cloud_firewall_vpc_firewall_control_policy DescribeCloudFirewallVpcFirewallControlPolicy Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("acl_action", objectRaw["AclAction"])
	d.Set("application_id", objectRaw["ApplicationId"])
	d.Set("application_name", objectRaw["ApplicationName"])
	d.Set("create_time", objectRaw["CreateTime"])
	d.Set("description", objectRaw["Description"])
	d.Set("dest_port", objectRaw["DestPort"])
	d.Set("dest_port_group", objectRaw["DestPortGroup"])
	d.Set("dest_port_type", objectRaw["DestPortType"])
	d.Set("destination", objectRaw["Destination"])
	d.Set("destination_group_type", objectRaw["DestinationGroupType"])
	d.Set("destination_type", objectRaw["DestinationType"])
	d.Set("domain_resolve_type", objectRaw["DomainResolveType"])
	d.Set("end_time", objectRaw["EndTime"])
	d.Set("hit_times", objectRaw["HitTimes"])
	d.Set("member_uid", objectRaw["MemberUid"])
	d.Set("order", objectRaw["Order"])
	d.Set("proto", objectRaw["Proto"])
	d.Set("release", objectRaw["Release"])
	d.Set("repeat_end_time", objectRaw["RepeatEndTime"])
	d.Set("repeat_start_time", objectRaw["RepeatStartTime"])
	d.Set("repeat_type", objectRaw["RepeatType"])
	d.Set("source", objectRaw["Source"])
	d.Set("source_group_type", objectRaw["SourceGroupType"])
	d.Set("source_type", objectRaw["SourceType"])
	d.Set("start_time", objectRaw["StartTime"])
	d.Set("acl_uuid", objectRaw["AclUuid"])

	applicationNameListRaw := make([]interface{}, 0)
	if objectRaw["ApplicationNameList"] != nil {
		applicationNameListRaw = convertToInterfaceArray(objectRaw["ApplicationNameList"])
	}

	d.Set("application_name_list", applicationNameListRaw)
	destPortGroupPortsRaw := make([]interface{}, 0)
	if objectRaw["DestPortGroupPorts"] != nil {
		destPortGroupPortsRaw = convertToInterfaceArray(objectRaw["DestPortGroupPorts"])
	}

	d.Set("dest_port_group_ports", destPortGroupPortsRaw)
	destinationGroupCidrsRaw := make([]interface{}, 0)
	if objectRaw["DestinationGroupCidrs"] != nil {
		destinationGroupCidrsRaw = convertToInterfaceArray(objectRaw["DestinationGroupCidrs"])
	}

	d.Set("destination_group_cidrs", destinationGroupCidrsRaw)
	repeatDaysRaw := make([]interface{}, 0)
	if objectRaw["RepeatDays"] != nil {
		repeatDaysRaw = convertToInterfaceArray(objectRaw["RepeatDays"])
	}

	d.Set("repeat_days", repeatDaysRaw)
	sourceGroupCidrsRaw := make([]interface{}, 0)
	if objectRaw["SourceGroupCidrs"] != nil {
		sourceGroupCidrsRaw = convertToInterfaceArray(objectRaw["SourceGroupCidrs"])
	}

	d.Set("source_group_cidrs", sourceGroupCidrsRaw)

	parts := strings.Split(d.Id(), ":")
	d.Set("vpc_firewall_id", parts[0])

	return nil
}

func resourceAliCloudCloudFirewallVpcFirewallControlPolicyUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	update := false

	var endpoint string
	var err error
	parts := strings.Split(d.Id(), ":")
	action := "ModifyVpcFirewallControlPolicy"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["VpcFirewallId"] = parts[0]
	request["AclUuid"] = parts[1]

	if d.HasChange("application_name") {
		update = true

		request["ApplicationName"] = d.Get("application_name")
	}

	if d.HasChange("proto") {
		update = true
	}
	request["Proto"] = d.Get("proto")
	if v, ok := d.GetOk("lang"); ok {
		request["Lang"] = v
	}
	if d.HasChange("destination_type") {
		update = true
	}
	request["DestinationType"] = d.Get("destination_type")
	if d.HasChange("description") {
		update = true
	}
	request["Description"] = d.Get("description")
	if d.HasChange("start_time") {
		update = true
		request["StartTime"] = d.Get("start_time")
	}

	if d.HasChange("dest_port") {
		update = true
	}
	if v, ok := d.GetOk("dest_port"); ok {
		request["DestPort"] = v
	}

	if d.HasChange("dest_port_group") {
		update = true
	}
	if v, ok := d.GetOk("dest_port_group"); ok {
		request["DestPortGroup"] = v
	}

	if d.HasChange("repeat_start_time") {
		update = true
		request["RepeatStartTime"] = d.Get("repeat_start_time")
	}

	if d.HasChange("destination") {
		update = true
	}
	request["Destination"] = d.Get("destination")
	if d.HasChange("application_name_list") {
		update = true
		if v, ok := d.GetOk("application_name_list"); ok || d.HasChange("application_name_list") {
			applicationNameListMapsArray := convertToInterfaceArray(v)

			request["ApplicationNameList"] = applicationNameListMapsArray
		}
	}

	if d.HasChange("dest_port_type") {
		update = true
	}
	if v, ok := d.GetOk("dest_port_type"); ok {
		request["DestPortType"] = v
	}

	if d.HasChange("repeat_end_time") {
		update = true
		request["RepeatEndTime"] = d.Get("repeat_end_time")
	}

	if d.HasChange("release") {
		update = true
	}
	if v, ok := d.GetOk("release"); ok {
		request["Release"] = v
	}

	if d.HasChange("acl_action") {
		update = true
	}
	request["AclAction"] = d.Get("acl_action")
	if d.HasChange("source_type") {
		update = true
	}
	request["SourceType"] = d.Get("source_type")
	if d.HasChange("domain_resolve_type") {
		update = true
		request["DomainResolveType"] = d.Get("domain_resolve_type")
	}

	if d.HasChange("source") {
		update = true
	}
	request["Source"] = d.Get("source")
	if d.HasChange("repeat_days") {
		update = true
		if v, ok := d.GetOk("repeat_days"); ok || d.HasChange("repeat_days") {
			repeatDaysMapsArray := convertToInterfaceArray(v)

			request["RepeatDays"] = repeatDaysMapsArray
		}
	}

	if d.HasChange("repeat_type") {
		update = true
		request["RepeatType"] = d.Get("repeat_type")
	}

	if d.HasChange("end_time") {
		update = true
		request["EndTime"] = d.Get("end_time")
	}

	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPostWithEndpoint("Cloudfw", "2017-12-07", action, query, request, true, endpoint)
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				} else if IsExpectedErrors(err, []string{"not buy user"}) {
					endpoint = connectivity.CloudFirewallOpenAPIEndpointControlPolicy
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

	return resourceAliCloudCloudFirewallVpcFirewallControlPolicyRead(d, meta)
}

func resourceAliCloudCloudFirewallVpcFirewallControlPolicyDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	parts := strings.Split(d.Id(), ":")
	action := "DeleteVpcFirewallControlPolicy"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	var endpoint string
	request = make(map[string]interface{})
	request["VpcFirewallId"] = parts[0]
	request["AclUuid"] = parts[1]

	if v, ok := d.GetOk("lang"); ok {
		request["Lang"] = v
	}
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = client.RpcPostWithEndpoint("Cloudfw", "2017-12-07", action, query, request, true, endpoint)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			} else if IsExpectedErrors(err, []string{"not buy user"}) {
				endpoint = connectivity.CloudFirewallOpenAPIEndpointControlPolicy
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
