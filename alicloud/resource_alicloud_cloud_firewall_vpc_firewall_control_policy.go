package alicloud

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	"github.com/PaesslerAG/jsonpath"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAlicloudCloudFirewallVpcFirewallControlPolicy() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudCloudFirewallVpcFirewallControlPolicyCreate,
		Read:   resourceAlicloudCloudFirewallVpcFirewallControlPolicyRead,
		Update: resourceAlicloudCloudFirewallVpcFirewallControlPolicyUpdate,
		Delete: resourceAlicloudCloudFirewallVpcFirewallControlPolicyDelete,
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
				Required:     true,
				Type:         schema.TypeString,
				ValidateFunc: validation.StringInSlice([]string{"accept", "drop", "log"}, false),
			},
			"acl_uuid": {
				Computed: true,
				Type:     schema.TypeString,
			},
			"application_id": {
				Computed: true,
				Type:     schema.TypeString,
			},
			"application_name": {
				Required:     true,
				Type:         schema.TypeString,
				ValidateFunc: validation.StringInSlice([]string{"FTP", "HTTP", "HTTPS", "MySQL", "SMTP", "SMTPS", "RDP", "VNC", "SSH", "Redis", "MQTT", "MongoDB", "Memcache", "SSL", "ANY"}, false),
			},
			"description": {
				Required: true,
				Type:     schema.TypeString,
			},
			"dest_port": {
				Optional: true,
				Type:     schema.TypeString,
			},
			"dest_port_group": {
				Optional: true,
				Type:     schema.TypeString,
			},
			"dest_port_group_ports": {
				Computed: true,
				Type:     schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"dest_port_type": {
				Optional:     true,
				Type:         schema.TypeString,
				ValidateFunc: validation.StringInSlice([]string{"port", "group"}, false),
			},
			"destination": {
				Required: true,
				Type:     schema.TypeString,
			},
			"destination_group_cidrs": {
				Computed: true,
				Type:     schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"destination_group_type": {
				Computed: true,
				Type:     schema.TypeString,
			},
			"destination_type": {
				Required:     true,
				Type:         schema.TypeString,
				ValidateFunc: validation.StringInSlice([]string{"net", "group", "domain"}, false),
			},
			"hit_times": {
				Computed: true,
				Type:     schema.TypeInt,
			},
			"lang": {
				Optional:     true,
				Type:         schema.TypeString,
				ValidateFunc: validation.StringInSlice([]string{"zh", "en"}, false),
			},
			"member_uid": {
				Optional: true,
				ForceNew: true,
				Computed: true,
				Type:     schema.TypeString,
			},
			"order": {
				Required: true,
				ForceNew: true,
				Type:     schema.TypeInt,
			},
			"proto": {
				Required:     true,
				Type:         schema.TypeString,
				ValidateFunc: validation.StringInSlice([]string{"ANY", "TCP", "UDP", "ICMP"}, false),
			},
			"release": {
				Optional: true,
				Computed: true,
				Type:     schema.TypeBool,
			},
			"source": {
				Required: true,
				Type:     schema.TypeString,
			},
			"source_group_cidrs": {
				Computed: true,
				Type:     schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"source_group_type": {
				Computed: true,
				Type:     schema.TypeString,
			},
			"source_type": {
				Required:     true,
				Type:         schema.TypeString,
				ValidateFunc: validation.StringInSlice([]string{"net", "group"}, false),
			},
			"vpc_firewall_id": {
				Required: true,
				ForceNew: true,
				Type:     schema.TypeString,
			},
		},
	}
}

func resourceAlicloudCloudFirewallVpcFirewallControlPolicyCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	request := make(map[string]interface{})
	conn, err := client.NewCloudfirewallClient()
	if err != nil {
		return WrapError(err)
	}

	if v, ok := d.GetOk("acl_action"); ok {
		request["AclAction"] = v
	}
	if v, ok := d.GetOk("application_name"); ok {
		request["ApplicationName"] = v
	}
	if v, ok := d.GetOk("description"); ok {
		request["Description"] = v
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
	if v, ok := d.GetOk("destination"); ok {
		request["Destination"] = v
	}
	if v, ok := d.GetOk("destination_type"); ok {
		request["DestinationType"] = v
	}
	if v, ok := d.GetOk("lang"); ok {
		request["Lang"] = v
	}
	if v, ok := d.GetOk("member_uid"); ok {
		request["MemberUid"] = v
	}
	if v, ok := d.GetOk("order"); ok {
		request["NewOrder"] = v
	}
	if v, ok := d.GetOk("proto"); ok {
		request["Proto"] = v
	}
	if v, ok := d.GetOkExists("release"); ok {
		request["Release"] = v
	}
	if v, ok := d.GetOk("source"); ok {
		request["Source"] = v
	}
	if v, ok := d.GetOk("source_type"); ok {
		request["SourceType"] = v
	}
	request["VpcFirewallId"] = d.Get("vpc_firewall_id")

	var response map[string]interface{}
	action := "CreateVpcFirewallControlPolicy"
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutCreate)), func() *resource.RetryError {
		resp, err := conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2017-12-07"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		response = resp
		addDebug(action, response, request)
		return nil
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_cloud_firewall_vpc_firewall_control_policy", action, AlibabaCloudSdkGoERROR)
	}
	aclUuidValue, err := jsonpath.Get("$.AclUuid", response)
	if err != nil || aclUuidValue == nil {
		return WrapErrorf(err, IdMsg, "alicloud_cloud_firewall_vpc_firewall_control_policy")
	}

	d.SetId(fmt.Sprint(request["VpcFirewallId"], ":", aclUuidValue))

	return resourceAlicloudCloudFirewallVpcFirewallControlPolicyRead(d, meta)
}

func resourceAlicloudCloudFirewallVpcFirewallControlPolicyRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	cloudfwService := CloudfwService{client}

	object, err := cloudfwService.DescribeCloudFirewallVpcFirewallControlPolicy(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_cloud_firewall_vpc_firewall_control_policy cloudfwService.DescribeCloudFirewallVpcFirewallControlPolicy Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}
	d.Set("vpc_firewall_id", parts[0])
	d.Set("acl_uuid", object["AclUuid"])
	d.Set("acl_action", object["AclAction"])
	d.Set("application_id", object["ApplicationId"])
	d.Set("application_name", object["ApplicationName"])
	d.Set("description", object["Description"])
	d.Set("dest_port", object["DestPort"])
	d.Set("dest_port_group", object["DestPortGroup"])
	destPortGroupPorts, _ := jsonpath.Get("$.DestPortGroupPorts", object)
	d.Set("dest_port_group_ports", destPortGroupPorts)
	d.Set("dest_port_type", object["DestPortType"])
	d.Set("destination", object["Destination"])
	destinationGroupCidrs, _ := jsonpath.Get("$.DestinationGroupCidrs", object)
	d.Set("destination_group_cidrs", destinationGroupCidrs)
	d.Set("destination_group_type", object["DestinationGroupType"])
	d.Set("destination_type", object["DestinationType"])
	d.Set("hit_times", object["HitTimes"])
	d.Set("member_uid", object["MemberUid"])
	d.Set("order", object["Order"])
	d.Set("proto", object["Proto"])
	d.Set("release", Interface2Bool(object["Release"]))
	d.Set("source", object["Source"])
	sourceGroupCidrs, _ := jsonpath.Get("$.SourceGroupCidrs", object)
	d.Set("source_group_cidrs", sourceGroupCidrs)
	d.Set("source_group_type", object["SourceGroupType"])
	d.Set("source_type", object["SourceType"])

	return nil
}

func resourceAlicloudCloudFirewallVpcFirewallControlPolicyUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	conn, err := client.NewCloudfirewallClient()
	if err != nil {
		return WrapError(err)
	}
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}
	update := false
	request := map[string]interface{}{
		"VpcFirewallId": parts[0],
		"AclUuid":       parts[1],
	}

	if d.HasChange("acl_action") {
		update = true
	}
	request["AclAction"] = d.Get("acl_action")
	if d.HasChange("application_name") {
		update = true
	}
	request["ApplicationName"] = d.Get("application_name")
	if d.HasChange("description") {
		update = true
	}
	request["Description"] = d.Get("description")
	if v, ok := d.GetOk("dest_port"); ok {
		request["DestPort"] = v
	}
	if d.HasChange("dest_port") {
		update = true
	}
	if v, ok := d.GetOk("dest_port_group"); ok {
		request["DestPortGroup"] = v
	}
	if d.HasChange("dest_port_group") {
		update = true
	}
	if d.HasChange("dest_port_type") {
		update = true
	}
	if v, ok := d.GetOk("dest_port_type"); ok {
		request["DestPortType"] = v
	}
	if d.HasChange("destination") {
		update = true
	}
	request["Destination"] = d.Get("destination")
	if d.HasChange("destination_type") {
		update = true
	}
	request["DestinationType"] = d.Get("destination_type")
	if v, ok := d.GetOk("lang"); ok {
		request["Lang"] = v
	}
	if d.HasChange("proto") {
		update = true
	}
	request["Proto"] = d.Get("proto")
	if d.HasChange("release") {
		update = true
	}
	if v, ok := d.GetOkExists("release"); ok {
		request["Release"] = v
	}
	if d.HasChange("source") {
		update = true
	}
	request["Source"] = d.Get("source")
	if d.HasChange("source_type") {
		update = true
	}
	request["SourceType"] = d.Get("source_type")

	if update {
		action := "ModifyVpcFirewallControlPolicy"
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *resource.RetryError {
			resp, err := conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2017-12-07"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			addDebug(action, resp, request)
			return nil
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
	}

	return resourceAlicloudCloudFirewallVpcFirewallControlPolicyRead(d, meta)
}

func resourceAlicloudCloudFirewallVpcFirewallControlPolicyDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	conn, err := client.NewCloudfirewallClient()
	if err != nil {
		return WrapError(err)
	}
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}

	request := map[string]interface{}{
		"VpcFirewallId": parts[0],
		"AclUuid":       parts[1],
	}

	if v, ok := d.GetOk("lang"); ok {
		request["Lang"] = v
	}

	action := "DeleteVpcFirewallControlPolicy"
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutDelete)), func() *resource.RetryError {
		resp, err := conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2017-12-07"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(action, resp, request)
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
