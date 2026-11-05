// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package expressconnect

import (
	"fmt"
	"log"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/convert"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAliCloudExpressConnectBgpPeer() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudExpressConnectBgpPeerCreate,
		Read:   resourceAliCloudExpressConnectBgpPeerRead,
		Update: resourceAliCloudExpressConnectBgpPeerUpdate,
		Delete: resourceAliCloudExpressConnectBgpPeerDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"bfd_multi_hop": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: IntBetween(1, 255),
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					if v, ok := d.GetOk("enable_bfd"); ok && fmt.Sprint(v) == "true" {
						return false
					}
					return true
				},
			},
			"bgp_group_id": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},
			"bgp_peer_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"enable_bfd": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"ip_version": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"peer_ip_address": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceAliCloudExpressConnectBgpPeerCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "CreateBgpPeer"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	request["RegionId"] = client.RegionId
	request["ClientToken"] = helper.BuildClientToken(action)

	request["BgpGroupId"] = d.Get("bgp_group_id")
	if v, ok := d.GetOk("peer_ip_address"); ok {
		request["PeerIpAddress"] = v
	}
	if v, ok := d.GetOkExists("enable_bfd"); ok {
		request["EnableBfd"] = v
	}
	if v, ok := d.GetOk("ip_version"); ok {
		request["IpVersion"] = v
	}
	if v, ok := d.GetOk("bfd_multi_hop"); ok {
		request["BfdMultiHop"] = v
	}
	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
		response, err = client.RpcPost("Vpc", "2016-04-28", action, query, request, true)
		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		logs.AddDebug(action, response, request)
		return nil
	})

	if err != nil {
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, "alicloud_vpc_bgp_peer", action, errs.AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(response["BgpPeerId"]))

	expressConnectServiceV2 := ExpressConnectServiceV2{client}
	stateConf := helper.BuildStateConf([]string{}, []string{"Available"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, expressConnectServiceV2.ExpressConnectBgpPeerStateRefreshFunc(d.Id(), "Status", []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return errs.WrapErrorf(err, errs.IdMsg, d.Id())
	}

	return resourceAliCloudExpressConnectBgpPeerRead(d, meta)
}

func resourceAliCloudExpressConnectBgpPeerRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	expressConnectServiceV2 := ExpressConnectServiceV2{client}

	objectRaw, err := expressConnectServiceV2.DescribeExpressConnectBgpPeer(d.Id())
	if err != nil {
		if !d.IsNewResource() && errs.NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_vpc_bgp_peer DescribeExpressConnectBgpPeer Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return errs.WrapError(err)
	}

	if v, ok := objectRaw["BfdMultiHop"]; ok && fmt.Sprint(v) != "0" {
		bfd_multi_hopVal, err := convert.FormatInt(v)
		if err != nil {
			return err
		}
		d.Set("bfd_multi_hop", bfd_multi_hopVal)
	}
	d.Set("bgp_group_id", objectRaw["BgpGroupId"])
	d.Set("bgp_peer_name", objectRaw["Name"])
	d.Set("enable_bfd", objectRaw["EnableBfd"])
	d.Set("ip_version", objectRaw["IpVersion"])
	d.Set("peer_ip_address", objectRaw["PeerIpAddress"])
	d.Set("status", objectRaw["Status"])

	return nil
}

func resourceAliCloudExpressConnectBgpPeerUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	update := false
	action := "ModifyBgpPeerAttribute"
	var err error
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	query["BgpPeerId"] = d.Id()
	request["RegionId"] = client.RegionId
	request["ClientToken"] = helper.BuildClientToken(action)
	if d.HasChange("bgp_group_id") {
		update = true
	}
	request["BgpGroupId"] = d.Get("bgp_group_id")
	if d.HasChange("peer_ip_address") {
		update = true
		request["PeerIpAddress"] = d.Get("peer_ip_address")
	}

	if d.HasChange("enable_bfd") {
		update = true
		request["EnableBfd"] = d.Get("enable_bfd")
	}

	if d.HasChange("bfd_multi_hop") {
		update = true
		request["BfdMultiHop"] = d.Get("bfd_multi_hop")
	}

	if update {
		wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
		err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
			response, err = client.RpcPost("Vpc", "2016-04-28", action, query, request, true)
			if err != nil {
				if helper.NeedRetry(err) {
					wait()
					return retry.RetryableError(err)
				}
				return retry.NonRetryableError(err)
			}
			logs.AddDebug(action, response, request)
			return nil
		})
		if err != nil {
			return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
		}
		expressConnectServiceV2 := ExpressConnectServiceV2{client}
		stateConf := helper.BuildStateConf([]string{}, []string{"Available"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, expressConnectServiceV2.ExpressConnectBgpPeerStateRefreshFunc(d.Id(), "Status", []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return errs.WrapErrorf(err, errs.IdMsg, d.Id())
		}
	}

	return resourceAliCloudExpressConnectBgpPeerRead(d, meta)
}

func resourceAliCloudExpressConnectBgpPeerDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	action := "DeleteBgpPeer"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	query["BgpPeerId"] = d.Id()
	request["RegionId"] = client.RegionId

	request["ClientToken"] = helper.BuildClientToken(action)

	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutDelete), func() *retry.RetryError {
		response, err = client.RpcPost("Vpc", "2016-04-28", action, query, request, true)
		request["ClientToken"] = helper.BuildClientToken(action)

		if err != nil || errs.IsExpectedErrors(err, []string{"DependencyViolation.BgpPeer"}) {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		logs.AddDebug(action, response, request)
		return nil
	})

	if err != nil {
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
	}

	expressConnectServiceV2 := ExpressConnectServiceV2{client}
	stateConf := helper.BuildStateConf([]string{}, []string{}, d.Timeout(schema.TimeoutDelete), 5*time.Second, expressConnectServiceV2.ExpressConnectBgpPeerStateRefreshFunc(d.Id(), "Status", []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return errs.WrapErrorf(err, errs.IdMsg, d.Id())
	}
	return nil
}
