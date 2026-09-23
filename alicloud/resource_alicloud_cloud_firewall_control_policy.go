package alicloud

import (
	"fmt"
	"log"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAliCloudCloudFirewallControlPolicy() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudCloudFirewallControlPolicyCreate,
		Read:   resourceAliCloudCloudFirewallControlPolicyRead,
		Update: resourceAliCloudCloudFirewallControlPolicyUpdate,
		Delete: resourceAliCloudCloudFirewallControlPolicyDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"direction": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: StringInSlice([]string{"in", "out"}, false),
			},
			"description": {
				Type:     schema.TypeString,
				Required: true,
			},
			"acl_action": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: StringInSlice([]string{"accept", "drop", "log"}, false),
			},
			"source": {
				Type:     schema.TypeString,
				Required: true,
			},
			"source_type": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: StringInSlice([]string{"net", "group", "location"}, false),
			},
			"destination": {
				Type:     schema.TypeString,
				Required: true,
			},
			"destination_type": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: StringInSlice([]string{"net", "group", "domain", "location"}, false),
			},
			"proto": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: StringInSlice([]string{"ANY", "TCP", "UDP", "ICMP"}, false),
			},
			"application_name": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: StringInSlice([]string{"ANY", "HTTP", "HTTPS", "MySQL", "SMTP", "SMTPS", "RDP", "VNC", "SSH", "Redis", "MQTT", "MongoDB", "Memcache", "SSL"}, false),
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
			"dest_port_type": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: StringInSlice([]string{"port", "group"}, false),
			},
			"ip_version": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"domain_resolve_type": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"repeat_type": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"start_time": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"end_time": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"repeat_start_time": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"repeat_end_time": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"repeat_days": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeInt},
			},
			"application_name_list": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"release": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"source_ip": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"lang": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: StringInSlice([]string{"zh", "en"}, false),
			},
			"acl_uuid": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceAliCloudCloudFirewallControlPolicyCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	var err error
	var endpoint string
	action := "AddControlPolicy"
	request := make(map[string]interface{})

	// order属性不透出
	request["NewOrder"] = "-1"
	request["Direction"] = d.Get("direction")
	request["Description"] = d.Get("description")
	request["AclAction"] = d.Get("acl_action")
	request["Source"] = d.Get("source")
	request["SourceType"] = d.Get("source_type")
	request["Destination"] = d.Get("destination")
	request["DestinationType"] = d.Get("destination_type")
	request["Proto"] = d.Get("proto")

	if v, ok := d.GetOk("application_name"); ok {
		request["ApplicationName"] = v
	}

	if v, ok := d.GetOk("dest_port"); ok {
		request["DestPort"] = v
	}

	if v, ok := d.GetOk("dest_port_group"); ok {
		request["DestPortGroup"] = v
	}

	if v, ok := d.GetOk("dest_port_type"); ok {
		request["DestPortType"] = v
	}

	if v, ok := d.GetOk("ip_version"); ok {
		request["IpVersion"] = v
	}

	if v, ok := d.GetOk("domain_resolve_type"); ok {
		request["DomainResolveType"] = v
	}

	if v, ok := d.GetOk("repeat_type"); ok {
		request["RepeatType"] = v
	}

	if v, ok := d.GetOkExists("start_time"); ok {
		request["StartTime"] = v
	}

	if v, ok := d.GetOkExists("end_time"); ok {
		request["EndTime"] = v
	}

	if v, ok := d.GetOk("repeat_start_time"); ok {
		request["RepeatStartTime"] = v
	}

	if v, ok := d.GetOk("repeat_end_time"); ok {
		request["RepeatEndTime"] = v
	}

	if v, ok := d.GetOk("repeat_days"); ok {
		request["RepeatDays"] = v
	}

	if v, ok := d.GetOk("application_name_list"); ok {
		request["ApplicationNameList"] = v
	}

	if v, ok := d.GetOk("source_ip"); ok {
		request["SourceIp"] = v
	}

	if v, ok := d.GetOk("lang"); ok {
		request["Lang"] = v
	}

	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutCreate)), func() *resource.RetryError {
		response, err = client.RpcPostWithEndpoint("Cloudfw", "2017-12-07", action, nil, request, false, endpoint)
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_cloud_firewall_control_policy", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprintf("%v:%v", response["AclUuid"], request["Direction"]))

	return resourceAliCloudCloudFirewallControlPolicyUpdate(d, meta)
}

func resourceAliCloudCloudFirewallControlPolicyRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	cloudFirewallServiceV2 := CloudFirewallServiceV2{client}

	objectRaw, err := cloudFirewallServiceV2.DescribeCloudFirewallControlPolicy(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_cloud_firewall_control_policy DescribeCloudFirewallControlPolicy Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	if objectRaw["AclAction"] != nil {
		d.Set("acl_action", objectRaw["AclAction"])
	}
	if objectRaw["ApplicationName"] != nil {
		d.Set("application_name", objectRaw["ApplicationName"])
	}
	if objectRaw["CreateTime"] != nil {
		d.Set("create_time", objectRaw["CreateTime"])
	}
	if objectRaw["Description"] != nil {
		d.Set("description", objectRaw["Description"])
	}
	if objectRaw["DestPort"] != nil {
		d.Set("dest_port", objectRaw["DestPort"])
	}
	if objectRaw["DestPortGroup"] != nil {
		d.Set("dest_port_group", objectRaw["DestPortGroup"])
	}
	if objectRaw["DestPortType"] != nil {
		d.Set("dest_port_type", objectRaw["DestPortType"])
	}
	if objectRaw["Destination"] != nil {
		d.Set("destination", objectRaw["Destination"])
	}
	if objectRaw["DestinationType"] != nil {
		d.Set("destination_type", objectRaw["DestinationType"])
	}
	if objectRaw["DomainResolveType"] != nil {
		d.Set("domain_resolve_type", objectRaw["DomainResolveType"])
	}
	if objectRaw["EndTime"] != nil {
		d.Set("end_time", objectRaw["EndTime"])
	}
	if objectRaw["IpVersion"] != nil {
		d.Set("ip_version", objectRaw["IpVersion"])
	}
	if objectRaw["Proto"] != nil {
		d.Set("proto", objectRaw["Proto"])
	}
	if objectRaw["Release"] != nil {
		d.Set("release", objectRaw["Release"])
	}
	if objectRaw["RepeatEndTime"] != nil {
		d.Set("repeat_end_time", objectRaw["RepeatEndTime"])
	}
	if objectRaw["RepeatStartTime"] != nil {
		d.Set("repeat_start_time", objectRaw["RepeatStartTime"])
	}
	if objectRaw["RepeatType"] != nil {
		d.Set("repeat_type", objectRaw["RepeatType"])
	}
	if objectRaw["Source"] != nil {
		d.Set("source", objectRaw["Source"])
	}
	if objectRaw["SourceType"] != nil {
		d.Set("source_type", objectRaw["SourceType"])
	}
	if objectRaw["StartTime"] != nil {
		d.Set("start_time", objectRaw["StartTime"])
	}
	if objectRaw["AclUuid"] != nil {
		d.Set("acl_uuid", objectRaw["AclUuid"])
	}
	if objectRaw["Direction"] != nil {
		d.Set("direction", objectRaw["Direction"])
	}

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

	return nil
}

func resourceAliCloudCloudFirewallControlPolicyUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	update := false

	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}

	request := map[string]interface{}{
		"AclUuid":   parts[0],
		"Direction": parts[1],
	}

	if !d.IsNewResource() && d.HasChange("description") {
		update = true
	}
	request["Description"] = d.Get("description")

	if !d.IsNewResource() && d.HasChange("acl_action") {
		update = true
	}
	request["AclAction"] = d.Get("acl_action")

	if !d.IsNewResource() && d.HasChange("source") {
		update = true
	}
	request["Source"] = d.Get("source")

	if !d.IsNewResource() && d.HasChange("source_type") {
		update = true
	}
	request["SourceType"] = d.Get("source_type")

	if !d.IsNewResource() && d.HasChange("destination") {
		update = true
	}
	request["Destination"] = d.Get("destination")

	if !d.IsNewResource() && d.HasChange("destination_type") {
		update = true
	}
	request["DestinationType"] = d.Get("destination_type")

	if !d.IsNewResource() && d.HasChange("proto") {
		update = true
	}
	request["Proto"] = d.Get("proto")

	if !d.IsNewResource() && d.HasChange("application_name") {
		update = true
	}
	if v, ok := d.GetOk("application_name"); ok {
		request["ApplicationName"] = v
	}

	if !d.IsNewResource() && d.HasChange("dest_port") {
		update = true
	}
	if v, ok := d.GetOk("dest_port"); ok {
		request["DestPort"] = v
	}

	if !d.IsNewResource() && d.HasChange("dest_port_group") {
		update = true
	}
	if v, ok := d.GetOk("dest_port_group"); ok {
		request["DestPortGroup"] = v
	}

	if !d.IsNewResource() && d.HasChange("dest_port_type") {
		update = true
	}
	if v, ok := d.GetOk("dest_port_type"); ok {
		request["DestPortType"] = v
	}

	if !d.IsNewResource() && d.HasChange("domain_resolve_type") {
		update = true
	}
	if v, ok := d.GetOk("domain_resolve_type"); ok {
		request["DomainResolveType"] = v
	}

	if !d.IsNewResource() && d.HasChange("repeat_type") {
		update = true
	}
	if v, ok := d.GetOk("repeat_type"); ok {
		request["RepeatType"] = v
	}

	if !d.IsNewResource() && d.HasChange("start_time") {
		update = true
	}
	if v, ok := d.GetOkExists("start_time"); ok {
		request["StartTime"] = v
	}

	if !d.IsNewResource() && d.HasChange("end_time") {
		update = true
	}
	if v, ok := d.GetOkExists("end_time"); ok {
		request["EndTime"] = v
	}

	if !d.IsNewResource() && d.HasChange("repeat_start_time") {
		update = true
	}
	if v, ok := d.GetOk("repeat_start_time"); ok {
		request["RepeatStartTime"] = v
	}

	if !d.IsNewResource() && d.HasChange("repeat_end_time") {
		update = true
	}
	if v, ok := d.GetOk("repeat_end_time"); ok {
		request["RepeatEndTime"] = v
	}

	if !d.IsNewResource() && d.HasChange("repeat_days") {
		update = true
	}
	if v, ok := d.GetOk("repeat_days"); ok {
		request["RepeatDays"] = v
	}

	if !d.IsNewResource() && d.HasChange("application_name_list") {
		update = true
	}
	if v, ok := d.GetOk("application_name_list"); ok {
		request["ApplicationNameList"] = v
	}

	if d.HasChange("release") {
		update = true
	}
	if v, ok := d.GetOk("release"); ok {
		request["Release"] = v
	}

	if v, ok := d.GetOk("lang"); ok {
		request["Lang"] = v
	}

	if update {
		action := "ModifyControlPolicy"
		var endpoint string
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *resource.RetryError {
			response, err = client.RpcPostWithEndpoint("Cloudfw", "2017-12-07", action, nil, request, false, endpoint)
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

	return resourceAliCloudCloudFirewallControlPolicyRead(d, meta)
}

func resourceAliCloudCloudFirewallControlPolicyDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	action := "DeleteControlPolicy"
	var response map[string]interface{}
	var endpoint string
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}

	request := map[string]interface{}{
		"AclUuid":   parts[0],
		"Direction": parts[1],
	}

	if v, ok := d.GetOk("source_ip"); ok {
		request["SourceIp"] = v
	}

	if v, ok := d.GetOk("lang"); ok {
		request["Lang"] = v
	}

	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutDelete)), func() *resource.RetryError {
		response, err = client.RpcPostWithEndpoint("Cloudfw", "2017-12-07", action, nil, request, false, endpoint)
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
