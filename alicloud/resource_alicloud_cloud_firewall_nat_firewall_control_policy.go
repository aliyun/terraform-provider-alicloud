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

func resourceAliCloudCloudFirewallNatFirewallControlPolicy() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudCloudFirewallNatFirewallControlPolicyCreate,
		Read:   resourceAliCloudCloudFirewallNatFirewallControlPolicyRead,
		Update: resourceAliCloudCloudFirewallNatFirewallControlPolicyUpdate,
		Delete: resourceAliCloudCloudFirewallNatFirewallControlPolicyDelete,
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
			"application_name_list": {
				Type:     schema.TypeList,
				Required: true,
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
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					if d.Get("dest_port_type").(string) == "port" {
						return false
					}
					return true
				},
			},
			"dest_port_group": {
				Type:     schema.TypeString,
				Optional: true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					if d.Get("dest_port_type").(string) == "group" {
						return false
					}
					return true
				},
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
			"destination_type": {
				Type:     schema.TypeString,
				Required: true,
			},
			"direction": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"domain_resolve_type": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"end_time": {
				Type:     schema.TypeInt,
				Optional: true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					if InArray(d.Get("repeat_type").(string), []string{"None", "Daily", "Weekly", "Monthly"}) {
						return false
					}
					return true
				},
			},
			"ip_version": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"nat_gateway_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"new_order": {
				Type:     schema.TypeString,
				Required: true,
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
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					if InArray(d.Get("repeat_type").(string), []string{"Weekly", "Monthly"}) {
						return false
					}
					return true
				},
			},
			"repeat_end_time": {
				Type:     schema.TypeString,
				Optional: true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					if InArray(d.Get("repeat_type").(string), []string{"Daily", "Weekly", "Monthly"}) {
						return false
					}
					return true
				},
			},
			"repeat_start_time": {
				Type:     schema.TypeString,
				Optional: true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					if InArray(d.Get("repeat_type").(string), []string{"Daily", "Weekly", "Monthly"}) {
						return false
					}
					return true
				},
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
			"source_type": {
				Type:     schema.TypeString,
				Required: true,
			},
			"start_time": {
				Type:     schema.TypeInt,
				Optional: true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					if InArray(d.Get("repeat_type").(string), []string{"None", "Daily", "Weekly", "Monthly"}) {
						return false
					}
					return true
				},
			},
		},
	}
}

func resourceAliCloudCloudFirewallNatFirewallControlPolicyCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "CreateNatFirewallControlPolicy"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	var endpoint string
	request = make(map[string]interface{})
	query["NatGatewayId"] = d.Get("nat_gateway_id")
	query["Direction"] = d.Get("direction")

	if v, ok := d.GetOk("lang"); ok {
		request["Lang"] = v
	}
	request["AclAction"] = d.Get("acl_action")
	request["Description"] = d.Get("description")
	if v, ok := d.GetOk("dest_port"); ok {
		request["DestPort"] = v
	}
	request["Destination"] = d.Get("destination")
	request["DestinationType"] = d.Get("destination_type")
	request["Proto"] = d.Get("proto")
	request["Source"] = d.Get("source")
	request["SourceType"] = d.Get("source_type")
	request["NewOrder"] = d.Get("new_order")
	if v, ok := d.GetOk("dest_port_type"); ok {
		request["DestPortType"] = v
	}
	if v, ok := d.GetOk("dest_port_group"); ok {
		request["DestPortGroup"] = v
	}
	if v, ok := d.GetOk("release"); ok {
		request["Release"] = v
	}
	if v, ok := d.GetOk("domain_resolve_type"); ok {
		request["DomainResolveType"] = v
	}
	if v, ok := d.GetOk("ip_version"); ok {
		request["IpVersion"] = v
	}
	if v, ok := d.GetOk("repeat_type"); ok {
		request["RepeatType"] = v
	}
	if v, ok := d.GetOk("repeat_start_time"); ok {
		request["RepeatStartTime"] = v
	}
	if v, ok := d.GetOk("repeat_end_time"); ok {
		request["RepeatEndTime"] = v
	}
	if v, ok := d.GetOk("start_time"); ok {
		request["StartTime"] = v
	}
	if v, ok := d.GetOk("end_time"); ok {
		request["EndTime"] = v
	}
	if v, ok := d.GetOk("repeat_days"); ok {
		repeatDaysMaps := v.([]interface{})
		request = expandArrayToMap(request, repeatDaysMaps, "RepeatDays")
	}

	if v, ok := d.GetOk("application_name_list"); ok {
		applicationNameListMaps := v.([]interface{})
		request = expandArrayToMap(request, applicationNameListMaps, "ApplicationNameList")
	}

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.RpcPostWithEndpoint("Cloudfw", "2017-12-07", action, query, request, false, endpoint)
		if err != nil {
			if NeedRetry(err) || IsExpectedErrors(err, []string{"-200142"}) {
				wait()
				return resource.RetryableError(err)
			} else if IsExpectedErrors(err, []string{"not buy user"}) {
				endpoint = connectivity.CloudFirewallOpenAPIEndpointControlPolicy
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}

		addDebug(action, response, request)
		return nil
	})

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_cloud_firewall_nat_firewall_control_policy", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprintf("%v:%v:%v", response["AclUuid"], query["NatGatewayId"], query["Direction"]))

	return resourceAliCloudCloudFirewallNatFirewallControlPolicyRead(d, meta)
}

func resourceAliCloudCloudFirewallNatFirewallControlPolicyRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	cloudFirewallServiceV2 := CloudFirewallServiceV2{client}

	objectRaw, err := cloudFirewallServiceV2.DescribeCloudFirewallNatFirewallControlPolicy(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_cloud_firewall_nat_firewall_control_policy DescribeCloudFirewallNatFirewallControlPolicy Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("acl_action", objectRaw["AclAction"])
	d.Set("create_time", objectRaw["CreateTime"])
	d.Set("description", objectRaw["Description"])
	d.Set("dest_port", objectRaw["DestPort"])
	d.Set("dest_port_group", objectRaw["DestPortGroup"])
	d.Set("dest_port_type", objectRaw["DestPortType"])
	d.Set("destination", objectRaw["Destination"])
	d.Set("destination_type", objectRaw["DestinationType"])
	d.Set("domain_resolve_type", objectRaw["DomainResolveType"])
	d.Set("end_time", objectRaw["EndTime"])
	d.Set("new_order", objectRaw["Order"])
	d.Set("proto", objectRaw["Proto"])
	d.Set("release", objectRaw["Release"])
	d.Set("repeat_end_time", objectRaw["RepeatEndTime"])
	d.Set("repeat_start_time", objectRaw["RepeatStartTime"])
	d.Set("repeat_type", objectRaw["RepeatType"])
	d.Set("source", objectRaw["Source"])
	d.Set("source_type", objectRaw["SourceType"])
	d.Set("start_time", objectRaw["StartTime"])
	d.Set("acl_uuid", objectRaw["AclUuid"])
	d.Set("nat_gateway_id", objectRaw["NatGatewayId"])

	applicationNameList1Raw := make([]interface{}, 0)
	if objectRaw["ApplicationNameList"] != nil {
		applicationNameList1Raw = objectRaw["ApplicationNameList"].([]interface{})
	}

	d.Set("application_name_list", applicationNameList1Raw)

	repeatDays1Raw := make([]interface{}, 0)
	if objectRaw["RepeatDays"] != nil {
		repeatDays1Raw = objectRaw["RepeatDays"].([]interface{})
	}

	d.Set("repeat_days", repeatDays1Raw)

	parts := strings.Split(d.Id(), ":")
	d.Set("acl_uuid", parts[0])
	d.Set("nat_gateway_id", parts[1])
	d.Set("direction", parts[2])

	return nil
}

func resourceAliCloudCloudFirewallNatFirewallControlPolicyUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	update := false
	d.Partial(true)
	parts := strings.Split(d.Id(), ":")
	action := "ModifyNatFirewallControlPolicy"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	query["NatGatewayId"] = parts[1]
	query["AclUuid"] = parts[0]
	query["Direction"] = parts[2]
	if d.HasChange("acl_action") {
		update = true
	}
	request["AclAction"] = d.Get("acl_action")
	if d.HasChange("description") {
		update = true
	}
	request["Description"] = d.Get("description")
	if d.HasChange("dest_port") {
		update = true
	}
	if v, ok := d.GetOk("dest_port"); ok && d.Get("dest_port_type").(string) == "port" {
		request["DestPort"] = v
	}
	if d.HasChange("destination") {
		update = true
	}
	request["Destination"] = d.Get("destination")
	if d.HasChange("destination_type") {
		update = true
	}
	request["DestinationType"] = d.Get("destination_type")
	if d.HasChange("proto") {
		update = true
	}
	request["Proto"] = d.Get("proto")
	if d.HasChange("source") {
		update = true
	}
	request["Source"] = d.Get("source")
	if d.HasChange("source_type") {
		update = true
	}
	request["SourceType"] = d.Get("source_type")
	if d.HasChange("dest_port_type") {
		update = true
	}
	if v, ok := d.GetOk("dest_port_type"); ok {
		request["DestPortType"] = v
	}
	if d.HasChange("dest_port_group") {
		update = true
	}
	if v, ok := d.GetOk("dest_port_group"); ok {
		if d.Get("dest_port_type").(string) == "group" {
			request["DestPortGroup"] = v
		}
	}
	if d.HasChange("release") {
		update = true
	}
	if v, ok := d.GetOk("release"); ok {
		request["Release"] = v
	}
	if d.HasChange("domain_resolve_type") {
		update = true
	}
	if v, ok := d.GetOk("domain_resolve_type"); ok {
		request["DomainResolveType"] = v
	}
	if d.HasChange("repeat_type") {
		update = true
	}
	if v, ok := d.GetOk("repeat_type"); ok {
		request["RepeatType"] = v
	}
	if d.HasChange("repeat_start_time") {
		update = true
	}
	if v, ok := d.GetOk("repeat_start_time"); ok {
		if InArray(d.Get("repeat_type").(string), []string{"Daily", "Weekly", "Monthly"}) {
			request["RepeatStartTime"] = v
		}
	}
	if d.HasChange("repeat_end_time") {
		update = true
	}
	if v, ok := d.GetOk("repeat_end_time"); ok {
		if InArray(d.Get("repeat_type").(string), []string{"Daily", "Weekly", "Monthly"}) {
			request["RepeatEndTime"] = v
		}
	}
	if d.HasChange("start_time") {
		update = true
	}
	if v, ok := d.GetOk("start_time"); ok {
		if InArray(d.Get("repeat_type").(string), []string{"None", "Daily", "Weekly", "Monthly"}) {
			request["StartTime"] = v
		}
	}
	if d.HasChange("end_time") {
		update = true
	}
	if v, ok := d.GetOk("end_time"); ok {
		if InArray(d.Get("repeat_type").(string), []string{"None", "Daily", "Weekly", "Monthly"}) {
			request["EndTime"] = v
		}
	}
	if d.HasChange("repeat_days") {
		update = true
	}
	if v, ok := d.GetOk("repeat_days"); ok {
		if InArray(d.Get("repeat_type").(string), []string{"Weekly", "Monthly"}) {
			repeatDaysMaps := v.([]interface{})
			request = expandArrayToMap(request, repeatDaysMaps, "RepeatDays")
		}
	}

	if d.HasChange("application_name_list") {
		update = true
	}
	if v, ok := d.GetOk("application_name_list"); ok {
		applicationNameListMaps := v.([]interface{})
		request = expandArrayToMap(request, applicationNameListMaps, "ApplicationNameList")
	}

	if v, ok := d.GetOk("lang"); ok {
		request["Lang"] = v
	}
	if update {
		var err error
		var endpoint string
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPostWithEndpoint("Cloudfw", "2017-12-07", action, query, request, false, endpoint)
			if err != nil {
				if NeedRetry(err) || IsExpectedErrors(err, []string{"-200511"}) {
					wait()
					return resource.RetryableError(err)
				} else if IsExpectedErrors(err, []string{"not buy user"}) {
					endpoint = connectivity.CloudFirewallOpenAPIEndpointControlPolicy
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
	update = false
	parts = strings.Split(d.Id(), ":")
	action = "ModifyNatFirewallControlPolicyPosition"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	query["NatGatewayId"] = parts[1]
	query["AclUuid"] = parts[0]
	query["Direction"] = parts[2]
	if v, ok := d.GetOk("lang"); ok {
		request["Lang"] = v
	}
	if d.HasChange("new_order") {
		update = true
	}
	request["NewOrder"] = d.Get("new_order")
	if update {
		var err error
		var endpoint string
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPostWithEndpoint("Cloudfw", "2017-12-07", action, query, request, false, endpoint)
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

			addDebug(action, response, request)
			return nil
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
	}

	d.Partial(false)
	return resourceAliCloudCloudFirewallNatFirewallControlPolicyRead(d, meta)
}

func resourceAliCloudCloudFirewallNatFirewallControlPolicyDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	parts := strings.Split(d.Id(), ":")
	action := "DeleteNatFirewallControlPolicy"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	var endpoint string
	request = make(map[string]interface{})
	query["AclUuid"] = parts[0]
	query["NatGatewayId"] = parts[1]
	query["Direction"] = parts[2]

	if v, ok := d.GetOk("lang"); ok {
		request["Lang"] = v
	}
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = client.RpcPostWithEndpoint("Cloudfw", "2017-12-07", action, query, request, false, endpoint)

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

		addDebug(action, response, request)
		return nil
	})

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	return nil
}
