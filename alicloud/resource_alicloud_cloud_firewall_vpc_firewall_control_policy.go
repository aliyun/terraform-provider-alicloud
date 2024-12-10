package alicloud

import (
	"fmt"
	"log"
	"time"

	"github.com/PaesslerAG/jsonpath"
	util "github.com/alibabacloud-go/tea-utils/service"
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
			"vpc_firewall_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"application_name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: StringInSlice([]string{"FTP", "HTTP", "HTTPS", "MySQL", "SMTP", "SMTPS", "RDP", "VNC", "SSH", "Redis", "MQTT", "MongoDB", "Memcache", "SSL", "ANY"}, false),
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
				ValidateFunc: StringInSlice([]string{"net", "group"}, false),
			},
			"destination": {
				Type:     schema.TypeString,
				Required: true,
			},
			"destination_type": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: StringInSlice([]string{"net", "group", "domain"}, false),
			},
			"proto": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: StringInSlice([]string{"ANY", "TCP", "UDP", "ICMP"}, false),
			},
			"order": {
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
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
			"release": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"member_uid": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
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
			"application_id": {
				Type:     schema.TypeString,
				Computed: true,
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
			"destination_group_cidrs": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"destination_group_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"dest_port_group_ports": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"hit_times": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func resourceAliCloudCloudFirewallVpcFirewallControlPolicyCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "CreateVpcFirewallControlPolicy"
	request := make(map[string]interface{})
	var err error
	var endpoint string

	request["VpcFirewallId"] = d.Get("vpc_firewall_id")
	request["ApplicationName"] = d.Get("application_name")
	request["Description"] = d.Get("description")
	request["AclAction"] = d.Get("acl_action")
	request["Source"] = d.Get("source")
	request["SourceType"] = d.Get("source_type")
	request["Destination"] = d.Get("destination")
	request["DestinationType"] = d.Get("destination_type")
	request["Proto"] = d.Get("proto")
	request["NewOrder"] = d.Get("order")

	if v, ok := d.GetOk("dest_port"); ok {
		request["DestPort"] = v
	}

	if v, ok := d.GetOk("dest_port_group"); ok {
		request["DestPortGroup"] = v
	}

	if v, ok := d.GetOk("dest_port_type"); ok {
		request["DestPortType"] = v
	}

	if v, ok := d.GetOkExists("release"); ok {
		request["Release"] = v
	}

	if v, ok := d.GetOk("member_uid"); ok {
		request["MemberUid"] = v
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_cloud_firewall_vpc_firewall_control_policy", action, AlibabaCloudSdkGoERROR)
	}

	aclUuidValue, err := jsonpath.Get("$.AclUuid", response)
	if err != nil || aclUuidValue == nil {
		return WrapErrorf(err, IdMsg, "alicloud_cloud_firewall_vpc_firewall_control_policy")
	}

	d.SetId(fmt.Sprint(request["VpcFirewallId"], ":", aclUuidValue))

	return resourceAliCloudCloudFirewallVpcFirewallControlPolicyRead(d, meta)
}

func resourceAliCloudCloudFirewallVpcFirewallControlPolicyRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	cloudfwService := CloudfwService{client}

	object, err := cloudfwService.DescribeCloudFirewallVpcFirewallControlPolicy(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
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
	d.Set("application_name", object["ApplicationName"])
	d.Set("description", object["Description"])
	d.Set("acl_action", object["AclAction"])
	d.Set("source", object["Source"])
	d.Set("source_type", object["SourceType"])
	d.Set("destination", object["Destination"])
	d.Set("destination_type", object["DestinationType"])
	d.Set("proto", object["Proto"])
	d.Set("order", object["Order"])
	d.Set("dest_port", object["DestPort"])
	d.Set("dest_port_group", object["DestPortGroup"])
	d.Set("dest_port_type", object["DestPortType"])
	d.Set("release", Interface2Bool(object["Release"]))
	d.Set("member_uid", object["MemberUid"])
	d.Set("acl_uuid", object["AclUuid"])
	d.Set("application_id", object["ApplicationId"])
	d.Set("source_group_type", object["SourceGroupType"])
	d.Set("destination_group_type", object["DestinationGroupType"])
	d.Set("hit_times", object["HitTimes"])

	sourceGroupCidrs, _ := jsonpath.Get("$.SourceGroupCidrs", object)
	d.Set("source_group_cidrs", sourceGroupCidrs)

	destinationGroupCidrs, _ := jsonpath.Get("$.DestinationGroupCidrs", object)
	d.Set("destination_group_cidrs", destinationGroupCidrs)

	destPortGroupPorts, _ := jsonpath.Get("$.DestPortGroupPorts", object)
	d.Set("dest_port_group_ports", destPortGroupPorts)

	return nil
}

func resourceAliCloudCloudFirewallVpcFirewallControlPolicyUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	update := false

	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}

	request := map[string]interface{}{
		"VpcFirewallId": parts[0],
		"AclUuid":       parts[1],
	}

	if d.HasChange("application_name") {
		update = true
	}
	request["ApplicationName"] = d.Get("application_name")

	if d.HasChange("description") {
		update = true
	}
	request["Description"] = d.Get("description")

	if d.HasChange("acl_action") {
		update = true
	}
	request["AclAction"] = d.Get("acl_action")

	if d.HasChange("source") {
		update = true
	}
	request["Source"] = d.Get("source")

	if d.HasChange("source_type") {
		update = true
	}
	request["SourceType"] = d.Get("source_type")

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

	if d.HasChange("dest_port_type") {
		update = true
	}
	if v, ok := d.GetOk("dest_port_type"); ok {
		request["DestPortType"] = v
	}

	if d.HasChange("release") {
		update = true
	}
	if v, ok := d.GetOkExists("release"); ok {
		request["Release"] = v
	}

	if update {
		if v, ok := d.GetOk("lang"); ok {
			request["Lang"] = v
		}

		action := "ModifyVpcFirewallControlPolicy"
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

	return resourceAliCloudCloudFirewallVpcFirewallControlPolicyRead(d, meta)
}

func resourceAliCloudCloudFirewallVpcFirewallControlPolicyDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	action := "DeleteVpcFirewallControlPolicy"
	var response map[string]interface{}
	var err error
	var endpoint string

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

	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
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
