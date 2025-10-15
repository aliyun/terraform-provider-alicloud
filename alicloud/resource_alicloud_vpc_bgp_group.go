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

func resourceAliCloudExpressConnectBgpGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudExpressConnectBgpGroupCreate,
		Read:   resourceAliCloudExpressConnectBgpGroupRead,
		Update: resourceAliCloudExpressConnectBgpGroupUpdate,
		Delete: resourceAliCloudExpressConnectBgpGroupDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"auth_key": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"bgp_group_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"clear_auth_key": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"ip_version": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"is_fake_asn": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"local_asn": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"peer_asn": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"region_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"route_limit": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"router_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceAliCloudExpressConnectBgpGroupCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "CreateBgpGroup"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	request["RegionId"] = client.RegionId
	request["ClientToken"] = buildClientToken(action)

	if v, ok := d.GetOkExists("is_fake_asn"); ok {
		request["IsFakeAsn"] = v
	}
	if v, ok := d.GetOk("bgp_group_name"); ok {
		request["Name"] = v
	}
	request["RouterId"] = d.Get("router_id")
	request["PeerAsn"] = d.Get("peer_asn")
	if v, ok := d.GetOk("description"); ok {
		request["Description"] = v
	}
	if v, ok := d.GetOkExists("route_limit"); ok {
		request["RouteQuota"] = v
	}
	if v, ok := d.GetOk("ip_version"); ok {
		request["IpVersion"] = v
	}
	if v, ok := d.GetOk("auth_key"); ok {
		request["AuthKey"] = v
	}
	if v, ok := d.GetOk("local_asn"); ok {
		request["LocalAsn"] = v
	}
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.RpcPost("Vpc", "2016-04-28", action, query, request, true)
		if err != nil {
			if NeedRetry(err) || IsExpectedErrors(err, []string{"OperationConflict"}) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_vpc_bgp_group", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(response["BgpGroupId"]))

	expressConnectServiceV2 := ExpressConnectServiceV2{client}
	stateConf := BuildStateConf([]string{}, []string{"Available"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, expressConnectServiceV2.ExpressConnectBgpGroupStateRefreshFunc(d.Id(), "Status", []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return resourceAliCloudExpressConnectBgpGroupUpdate(d, meta)
}

func resourceAliCloudExpressConnectBgpGroupRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	expressConnectServiceV2 := ExpressConnectServiceV2{client}

	objectRaw, err := expressConnectServiceV2.DescribeExpressConnectBgpGroup(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_vpc_bgp_group DescribeExpressConnectBgpGroup Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("auth_key", objectRaw["AuthKey"])
	d.Set("bgp_group_name", objectRaw["Name"])
	d.Set("description", objectRaw["Description"])
	d.Set("ip_version", objectRaw["IpVersion"])
	d.Set("is_fake_asn", formatBool(objectRaw["IsFake"]))
	d.Set("local_asn", objectRaw["LocalAsn"])
	d.Set("peer_asn", objectRaw["PeerAsn"])
	d.Set("region_id", objectRaw["RegionId"])
	d.Set("route_limit", formatInt(objectRaw["RouteLimit"]))
	d.Set("router_id", objectRaw["RouterId"])
	d.Set("status", objectRaw["Status"])

	return nil
}

func resourceAliCloudExpressConnectBgpGroupUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	update := false

	var err error
	action := "ModifyBgpGroupAttribute"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["BgpGroupId"] = d.Id()
	request["RegionId"] = client.RegionId
	request["ClientToken"] = buildClientToken(action)
	if !d.IsNewResource() && d.HasChange("is_fake_asn") {
		update = true
		request["IsFakeAsn"] = d.Get("is_fake_asn")
	}

	if !d.IsNewResource() && d.HasChange("bgp_group_name") {
		update = true
		request["Name"] = d.Get("bgp_group_name")
	}

	if !d.IsNewResource() && d.HasChange("peer_asn") {
		update = true
	}
	request["PeerAsn"] = d.Get("peer_asn")
	if !d.IsNewResource() && d.HasChange("description") {
		update = true
		request["Description"] = d.Get("description")
	}

	if v, ok := d.GetOk("clear_auth_key"); ok {
		request["ClearAuthKey"] = v
	}
	if !d.IsNewResource() && d.HasChange("route_limit") {
		update = true
		request["RouteQuota"] = d.Get("route_limit")
	}

	if !d.IsNewResource() && d.HasChange("auth_key") {
		update = true
		request["AuthKey"] = d.Get("auth_key")
	}

	if !d.IsNewResource() && d.HasChange("local_asn") {
		update = true
		request["LocalAsn"] = d.Get("local_asn")
	}

	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("Vpc", "2016-04-28", action, query, request, true)
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
		expressConnectServiceV2 := ExpressConnectServiceV2{client}
		stateConf := BuildStateConf([]string{}, []string{"Available"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, expressConnectServiceV2.ExpressConnectBgpGroupStateRefreshFunc(d.Id(), "Status", []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
	}

	return resourceAliCloudExpressConnectBgpGroupRead(d, meta)
}

func resourceAliCloudExpressConnectBgpGroupDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	action := "DeleteBgpGroup"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	request["BgpGroupId"] = d.Id()
	request["RegionId"] = client.RegionId
	request["ClientToken"] = buildClientToken(action)

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = client.RpcPost("Vpc", "2016-04-28", action, query, request, true)
		request["ClientToken"] = buildClientToken(action)

		if err != nil {
			if NeedRetry(err) || IsExpectedErrors(err, []string{"DependencyViolation.BgpPeer"}) {
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

	expressConnectServiceV2 := ExpressConnectServiceV2{client}
	stateConf := BuildStateConf([]string{}, []string{""}, d.Timeout(schema.TimeoutDelete), 5*time.Second, expressConnectServiceV2.ExpressConnectBgpGroupStateRefreshFunc(d.Id(), "Status", []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return nil
}
