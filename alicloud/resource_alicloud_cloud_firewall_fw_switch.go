// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"log"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAliCloudCloudFirewallFwSwitch() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudCloudFirewallFwSwitchCreate,
		Read:   resourceAliCloudCloudFirewallFwSwitchRead,
		Update: resourceAliCloudCloudFirewallFwSwitchUpdate,
		Delete: resourceAliCloudCloudFirewallFwSwitchDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"ali_uid": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"bind_instance_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"bind_instance_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"dry_run": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"internet_address": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"intranet_address": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"ip_version": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"lang": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"member_uid": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"note": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"protect_status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"region_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"region_status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"resource_instance_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"resource_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"sg_status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"sg_status_time": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"sync_status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"type": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceAliCloudCloudFirewallFwSwitchCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "PutEnableFwSwitch"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	if v, ok := d.GetOk("internet_address"); ok {
		localData, _ := jsonpath.Get("$", v)
		ipaddrListMapsArray := convertToInterfaceArray(localData)

		request["IpaddrList"] = ipaddrListMapsArray
	}

	request["ClientToken"] = buildClientToken(action)

	if v, ok := d.GetOk("dry_run"); ok {
		request["DryRun"] = v
	}
	if v, ok := d.GetOk("ip_version"); ok {
		request["IpVersion"] = v
	}
	if v, ok := d.GetOk("lang"); ok {
		request["Lang"] = v
	}
	if v, ok := d.GetOk("member_uid"); ok {
		request["MemberUid"] = v
	}

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.RpcPost("Cloudfw", "2017-12-07", action, query, request, true)
		if err != nil {
			if IsExpectedErrors(err, []string{"IdempotentToken.Processing"}) || NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_cloud_firewall_fw_switch", action, AlibabaCloudSdkGoERROR)
	}
	if abnormalList, ok := response["AbnormalResourceStatusList"].([]interface{}); ok && len(abnormalList) > 0 {
		return WrapError(fmt.Errorf("%s failed, response: %v", action, response))
	}

	id, _ := jsonpath.Get("IpaddrList[0]", request)
	d.SetId(fmt.Sprint(id))

	cloudFirewallServiceV2 := CloudFirewallServiceV2{client}
	stateConf := BuildStateConf([]string{}, []string{"open"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, cloudFirewallServiceV2.CloudFirewallFwSwitchStateRefreshFuncWithMemberUid(d.Id(), "ProtectStatus", []string{}, d.Get("member_uid").(int)))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return resourceAliCloudCloudFirewallFwSwitchRead(d, meta)
}

func resourceAliCloudCloudFirewallFwSwitchRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	cloudFirewallServiceV2 := CloudFirewallServiceV2{client}

	objectRaw, err := cloudFirewallServiceV2.DescribeCloudFirewallFwSwitchWithMemberUid(d.Id(), d.Get("member_uid").(int))
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_cloud_firewall_fw_switch DescribeCloudFirewallFwSwitch Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("ali_uid", objectRaw["AliUid"])
	d.Set("bind_instance_id", objectRaw["BindInstanceId"])
	d.Set("bind_instance_name", objectRaw["BindInstanceName"])
	d.Set("intranet_address", objectRaw["IntranetAddress"])
	if v := objectRaw["IpVersion"]; v != nil {
		d.Set("ip_version", fmt.Sprint(v))
	}
	d.Set("member_uid", objectRaw["MemberUid"])
	d.Set("name", objectRaw["Name"])
	d.Set("note", objectRaw["Note"])
	d.Set("protect_status", objectRaw["ProtectStatus"])
	d.Set("region_id", objectRaw["RegionID"])
	d.Set("region_status", objectRaw["RegionStatus"])
	d.Set("resource_instance_id", objectRaw["ResourceInstanceId"])
	d.Set("resource_type", objectRaw["ResourceType"])
	d.Set("sg_status", objectRaw["SgStatus"])
	d.Set("sg_status_time", objectRaw["SgStatusTime"])
	d.Set("sync_status", objectRaw["SyncStatus"])
	d.Set("type", objectRaw["Type"])

	d.Set("internet_address", d.Id())

	return nil
}

func resourceAliCloudCloudFirewallFwSwitchUpdate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[INFO] Cannot update resource Alicloud Resource Fw Switch.")
	return resourceAliCloudCloudFirewallFwSwitchRead(d, meta)
}

func resourceAliCloudCloudFirewallFwSwitchDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	action := "PutDisableFwSwitch"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	request["IpaddrList.1"] = d.Id()

	request["ClientToken"] = buildClientToken(action)

	if v, ok := d.GetOk("dry_run"); ok {
		request["DryRun"] = v
	}
	if v, ok := d.GetOk("ip_version"); ok {
		request["IpVersion"] = v
	}
	if v, ok := d.GetOk("lang"); ok {
		request["Lang"] = v
	}
	if v, ok := d.GetOk("member_uid"); ok {
		request["MemberUid"] = v
	}

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = client.RpcPost("Cloudfw", "2017-12-07", action, query, request, true)
		if err != nil {
			if IsExpectedErrors(err, []string{"IdempotentToken.Processing"}) || NeedRetry(err) {
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
	stateConf := BuildStateConf([]string{}, []string{"closed"}, d.Timeout(schema.TimeoutDelete), 5*time.Second, cloudFirewallServiceV2.CloudFirewallFwSwitchStateRefreshFuncWithMemberUid(d.Id(), "ProtectStatus", []string{}, d.Get("member_uid").(int)))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return nil
}
