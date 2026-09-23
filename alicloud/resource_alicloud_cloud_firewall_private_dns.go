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

func resourceAliCloudCloudFirewallPrivateDns() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudCloudFirewallPrivateDnsCreate,
		Read:   resourceAliCloudCloudFirewallPrivateDnsRead,
		Update: resourceAliCloudCloudFirewallPrivateDnsUpdate,
		Delete: resourceAliCloudCloudFirewallPrivateDnsDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"access_instance_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"access_instance_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"domain_name_list": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"firewall_type": {
				Type:     schema.TypeList,
				Required: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"ip_protocol": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"member_uid": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
			},
			"port": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
			},
			"primary_dns": {
				Type:     schema.TypeString,
				Optional: true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					if v, ok := d.GetOk("private_dns_type"); ok && v.(string) == "PrivateZone" {
						return new == ""
					}
					return false
				},
			},
			"primary_vswitch_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"primary_vswitch_ip": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"private_dns_type": {
				Type:     schema.TypeString,
				Required: true,
			},
			"region_no": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"standby_dns": {
				Type:     schema.TypeString,
				Optional: true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					if v, ok := d.GetOk("private_dns_type"); ok && v.(string) == "PrivateZone" {
						return new == ""
					}
					return false
				},
			},
			"standby_vswitch_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"standby_vswitch_ip": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"vpc_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceAliCloudCloudFirewallPrivateDnsCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "CreatePrivateDnsEndpoint"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	if v, ok := d.GetOk("region_no"); ok {
		request["RegionNo"] = v
	}

	privateDnsType := d.Get("private_dns_type").(string)

	// 当 private_dns_type 为 PrivateZone 时，不能传递 standby_dns 和 primary_dns
	if privateDnsType != "PrivateZone" {
		if v, ok := d.GetOk("standby_dns"); ok {
			request["StandbyDns"] = v
		}
	}
	if v, ok := d.GetOk("standby_vswitch_id"); ok {
		request["StandbyVSwitchId"] = v
	}
	if v, ok := d.GetOkExists("member_uid"); ok {
		request["MemberUid"] = v
	}
	request["AccessInstanceName"] = d.Get("access_instance_name")
	if v, ok := d.GetOk("primary_vswitch_ip"); ok {
		request["PrimaryVSwitchIp"] = v
	}
	if v, ok := d.GetOk("standby_vswitch_ip"); ok {
		request["StandbyVSwitchIp"] = v
	}
	request["PrivateDnsType"] = privateDnsType
	if v, ok := d.GetOk("ip_protocol"); ok {
		request["IpProtocol"] = v
	}
	request["VpcId"] = d.Get("vpc_id")
	if v, ok := d.GetOk("firewall_type"); ok {
		firewallTypeMapsArray := convertToInterfaceArray(v)

		request["FirewallType"] = firewallTypeMapsArray
	}

	// 当 private_dns_type 为 PrivateZone 时，不能传递 primary_dns
	if privateDnsType != "PrivateZone" {
		if v, ok := d.GetOk("primary_dns"); ok {
			request["PrimaryDns"] = v
		}
	}
	if v, ok := d.GetOkExists("port"); ok {
		request["Port"] = v
	}
	if v, ok := d.GetOk("primary_vswitch_id"); ok {
		request["PrimaryVSwitchId"] = v
	}
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_cloud_firewall_private_dns", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprintf("%v:%v", response["AccessInstanceId"], request["RegionNo"]))

	cloudFirewallServiceV2 := CloudFirewallServiceV2{client}
	stateConf := BuildStateConf([]string{}, []string{"normal"}, d.Timeout(schema.TimeoutCreate), 30*time.Second, cloudFirewallServiceV2.CloudFirewallPrivateDnsStateRefreshFunc(d.Id(), "Status", []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return resourceAliCloudCloudFirewallPrivateDnsUpdate(d, meta)
}

func resourceAliCloudCloudFirewallPrivateDnsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	cloudFirewallServiceV2 := CloudFirewallServiceV2{client}

	objectRaw, err := cloudFirewallServiceV2.DescribeCloudFirewallPrivateDns(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_cloud_firewall_private_dns DescribeCloudFirewallPrivateDns Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("access_instance_name", objectRaw["AccessInstanceName"])
	d.Set("ip_protocol", objectRaw["IpProtocol"])
	d.Set("member_uid", objectRaw["MemberUid"])
	d.Set("port", objectRaw["Port"])
	d.Set("primary_dns", objectRaw["PrimaryDns"])
	d.Set("primary_vswitch_id", objectRaw["PrimaryVSwitchId"])
	d.Set("primary_vswitch_ip", objectRaw["PrimaryVSwitchIp"])
	d.Set("private_dns_type", objectRaw["PrivateDnsType"])
	d.Set("standby_dns", objectRaw["StandbyDns"])
	d.Set("standby_vswitch_id", objectRaw["StandbyVSwitchId"])
	d.Set("standby_vswitch_ip", objectRaw["StandbyVSwitchIp"])
	d.Set("status", objectRaw["Status"])
	d.Set("vpc_id", objectRaw["VpcId"])
	d.Set("access_instance_id", objectRaw["AccessInstanceId"])
	d.Set("region_no", objectRaw["RegionNo"])

	firewallTypeRaw := make([]interface{}, 0)
	if objectRaw["FirewallType"] != nil {
		firewallTypeRaw = convertToInterfaceArray(objectRaw["FirewallType"])
	}

	d.Set("firewall_type", firewallTypeRaw)

	objectRaw, err = cloudFirewallServiceV2.DescribePrivateDnsDescribePrivateDnsDomainNameList(d.Id())
	if err != nil && !NotFoundError(err) {
		return WrapError(err)
	}

	domainNameListRaw, _ := jsonpath.Get("$.DomainNameList", objectRaw)

	d.Set("domain_name_list", domainNameListRaw)

	return nil
}

func resourceAliCloudCloudFirewallPrivateDnsUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	update := false
	d.Partial(true)

	var err error
	parts := strings.Split(d.Id(), ":")
	action := "ModifyPrivateDnsEndpoint"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["RegionNo"] = parts[1]
	request["AccessInstanceId"] = parts[0]

	privateDnsType := d.Get("private_dns_type").(string)

	if privateDnsType != "PrivateZone" {
		if !d.IsNewResource() && d.HasChange("standby_dns") {
			update = true
		}
		if v, ok := d.GetOk("standby_dns"); ok || d.HasChange("standby_dns") {
			request["StandbyDns"] = v
		}
		if !d.IsNewResource() && d.HasChange("primary_dns") {
			update = true
		}
		if v, ok := d.GetOk("primary_dns"); ok || d.HasChange("primary_dns") {
			request["PrimaryDns"] = v
		}
	}
	if !d.IsNewResource() && d.HasChange("access_instance_name") {
		update = true
	}
	request["AccessInstanceName"] = d.Get("access_instance_name")
	if !d.IsNewResource() && d.HasChange("private_dns_type") {
		update = true
	}
	request["PrivateDnsType"] = d.Get("private_dns_type")
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
		stateConf := BuildStateConf([]string{}, []string{"normal"}, d.Timeout(schema.TimeoutUpdate), 30*time.Second, cloudFirewallServiceV2.CloudFirewallPrivateDnsStateRefreshFunc(d.Id(), "Status", []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
	}
	update = false
	parts = strings.Split(d.Id(), ":")

	if d.HasChange("domain_name_list") {
		update = true
		oldValue, newValue := d.GetChange("domain_name_list")

		oldSet := oldValue.(*schema.Set)
		newSet := newValue.(*schema.Set)
		toRemove := oldSet.Difference(newSet)

		if toRemove.Len() > 0 {
			action = "DeletePrivateDnsDomainName"
			request = make(map[string]interface{})
			query = make(map[string]interface{})
			request["RegionNo"] = parts[1]
			request["AccessInstanceId"] = parts[0]
			request["DomainNameList"] = convertToInterfaceArray(toRemove)

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
		}

		toAdd := newSet.Difference(oldSet)

		if toAdd.Len() > 0 {
			action = "AddPrivateDnsDomainName"
			request = make(map[string]interface{})
			query = make(map[string]interface{})
			request["RegionNo"] = parts[1]
			request["AccessInstanceId"] = parts[0]
			request["DomainNameList"] = convertToInterfaceArray(toAdd)

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
		}
	}

	d.Partial(false)
	return resourceAliCloudCloudFirewallPrivateDnsRead(d, meta)
}

func resourceAliCloudCloudFirewallPrivateDnsDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	parts := strings.Split(d.Id(), ":")
	action := "DeletePrivateDnsAllDomainName"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	request["RegionNo"] = parts[1]
	request["AccessInstanceId"] = parts[0]

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

	action = "DeletePrivateDnsEndpoint"
	request = make(map[string]interface{})
	request["RegionNo"] = parts[1]
	request["AccessInstanceId"] = parts[0]

	wait = incrementalWait(3*time.Second, 5*time.Second)
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
	stateConf := BuildStateConf([]string{}, []string{""}, d.Timeout(schema.TimeoutDelete), 30*time.Second, cloudFirewallServiceV2.CloudFirewallPrivateDnsStateRefreshFunc(d.Id(), "AccessInstanceId", []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return nil
}
