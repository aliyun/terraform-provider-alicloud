package alicloud

import (
	"fmt"
	"log"
	"time"

	util "github.com/alibabacloud-go/tea-utils/service"
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
		Schema: map[string]*schema.Schema{
			"direction": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: StringInSlice([]string{"in", "out"}, false),
			},
			"application_name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: StringInSlice([]string{"ANY", "HTTP", "HTTPS", "MQTT", "Memcache", "MongoDB", "MySQL", "RDP", "Redis", "SMTP", "SMTPS", "SSH", "SSL", "VNC"}, false),
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
		},
	}
}

func resourceAliCloudCloudFirewallControlPolicyCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "AddControlPolicy"
	request := make(map[string]interface{})
	conn, err := client.NewCloudfwClient()
	if err != nil {
		return WrapError(err)
	}

	// order属性不透出
	request["NewOrder"] = "-1"
	request["Direction"] = d.Get("direction")
	request["ApplicationName"] = d.Get("application_name")
	request["Description"] = d.Get("description")
	request["AclAction"] = d.Get("acl_action")
	request["Source"] = d.Get("source")
	request["SourceType"] = d.Get("source_type")
	request["Destination"] = d.Get("destination")
	request["DestinationType"] = d.Get("destination_type")
	request["Proto"] = d.Get("proto")

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

	if v, ok := d.GetOk("source_ip"); ok {
		request["SourceIp"] = v
	}

	if v, ok := d.GetOk("lang"); ok {
		request["Lang"] = v
	}

	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutCreate)), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2017-12-07"), StringPointer("AK"), nil, request, &runtime)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}

		if fmt.Sprint(response["Message"]) == "not buy user" {
			conn.Endpoint = String(connectivity.CloudFirewallOpenAPIEndpointControlPolicy)
			return resource.RetryableError(fmt.Errorf("%s", response))
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
	cloudfwService := CloudfwService{client}

	object, err := cloudfwService.DescribeCloudFirewallControlPolicy(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_cloud_firewall_control_policy cloudfwService.DescribeCloudFirewallControlPolicy Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("direction", object["Direction"])
	d.Set("application_name", object["ApplicationName"])
	d.Set("description", object["Description"])
	d.Set("acl_action", object["AclAction"])
	d.Set("source", object["Source"])
	d.Set("source_type", object["SourceType"])
	d.Set("destination", object["Destination"])
	d.Set("destination_type", object["DestinationType"])
	d.Set("proto", object["Proto"])
	d.Set("dest_port", object["DestPort"])
	d.Set("dest_port_group", object["DestPortGroup"])
	d.Set("dest_port_type", object["DestPortType"])
	d.Set("ip_version", object["IpVersion"])
	d.Set("release", object["Release"])
	d.Set("acl_uuid", object["AclUuid"])

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

	if !d.IsNewResource() && d.HasChange("application_name") {
		update = true
	}
	request["ApplicationName"] = d.Get("application_name")

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
		conn, err := client.NewCloudfwClient()
		if err != nil {
			return WrapError(err)
		}

		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2017-12-07"), StringPointer("AK"), nil, request, &runtime)
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}

			if fmt.Sprint(response["Message"]) == "not buy user" {
				conn.Endpoint = String(connectivity.CloudFirewallOpenAPIEndpointControlPolicy)
				return resource.RetryableError(fmt.Errorf("%s", response))
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
	conn, err := client.NewCloudfwClient()
	if err != nil {
		return WrapError(err)
	}

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

	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutDelete)), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2017-12-07"), StringPointer("AK"), nil, request, &runtime)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}

		if fmt.Sprint(response["Message"]) == "not buy user" {
			conn.Endpoint = String(connectivity.CloudFirewallOpenAPIEndpointControlPolicy)
			return resource.RetryableError(fmt.Errorf("%s", response))
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
