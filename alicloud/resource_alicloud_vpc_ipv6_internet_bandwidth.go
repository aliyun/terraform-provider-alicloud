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

func resourceAliCloudVpcIpv6InternetBandwidth() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudVpcIpv6InternetBandwidthCreate,
		Read:   resourceAliCloudVpcIpv6InternetBandwidthRead,
		Update: resourceAliCloudVpcIpv6InternetBandwidthUpdate,
		Delete: resourceAliCloudVpcIpv6InternetBandwidthDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"bandwidth": {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: IntBetween(1, 5000),
			},
			"internet_charge_type": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				ValidateFunc: StringInSlice([]string{"PayByTraffic", "PayByBandwidth"}, false),
			},
			"ipv6_address_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"ipv6_gateway_id": {
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

func resourceAliCloudVpcIpv6InternetBandwidthCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	action := "AllocateIpv6InternetBandwidth"
	var request map[string]interface{}
	var response map[string]interface{}
	var err error
	request = make(map[string]interface{})
	request["RegionId"] = client.RegionId
	request["ClientToken"] = buildClientToken(action)

	request["Ipv6GatewayId"] = d.Get("ipv6_gateway_id")
	request["Ipv6AddressId"] = d.Get("ipv6_address_id")
	request["Bandwidth"] = d.Get("bandwidth")
	if v, ok := d.GetOk("internet_charge_type"); ok {
		request["InternetChargeType"] = v
	}
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.RpcPost("Vpc", "2016-04-28", action, nil, request, true)
		request["ClientToken"] = buildClientToken(action)

		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(action, response, request)
		return nil
	})

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_vpc_ipv6_internet_bandwidth", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(response["InternetBandwidthId"]))

	return resourceAliCloudVpcIpv6InternetBandwidthRead(d, meta)
}

func resourceAliCloudVpcIpv6InternetBandwidthRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	vpcServiceV2 := VpcServiceV2{client}

	objectRaw, err := vpcServiceV2.DescribeVpcIpv6InternetBandwidth(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_vpc_ipv6_internet_bandwidth DescribeVpcIpv6InternetBandwidth Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("ipv6_address_id", objectRaw["Ipv6AddressId"])
	d.Set("ipv6_gateway_id", objectRaw["Ipv6GatewayId"])
	ipv6InternetBandwidth1RawObj, _ := jsonpath.Get("$.Ipv6InternetBandwidth", objectRaw)
	ipv6InternetBandwidth1Raw := make(map[string]interface{})
	if ipv6InternetBandwidth1RawObj != nil {
		ipv6InternetBandwidth1Raw = ipv6InternetBandwidth1RawObj.(map[string]interface{})
	}
	d.Set("bandwidth", ipv6InternetBandwidth1Raw["Bandwidth"])
	d.Set("internet_charge_type", ipv6InternetBandwidth1Raw["InternetChargeType"])
	d.Set("status", ipv6InternetBandwidth1Raw["BusinessStatus"])

	return nil
}

func resourceAliCloudVpcIpv6InternetBandwidthUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	update := false
	action := "ModifyIpv6InternetBandwidth"
	var err error
	request = make(map[string]interface{})

	request["Ipv6InternetBandwidthId"] = d.Id()
	request["RegionId"] = client.RegionId
	request["ClientToken"] = buildClientToken(action)
	if d.HasChange("bandwidth") {
		update = true
	}
	request["Bandwidth"] = d.Get("bandwidth")
	if d.HasChange("ipv6_address_id") {
		update = true
	}
	request["Ipv6AddressId"] = d.Get("ipv6_address_id")
	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("Vpc", "2016-04-28", action, nil, request, true)
			request["ClientToken"] = buildClientToken(action)

			if err != nil {
				if NeedRetry(err) {
					wait()
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
	return resourceAliCloudVpcIpv6InternetBandwidthRead(d, meta)
}

func resourceAliCloudVpcIpv6InternetBandwidthDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "DeleteIpv6InternetBandwidth"
	var request map[string]interface{}
	var response map[string]interface{}
	var err error
	request = make(map[string]interface{})

	request["Ipv6InternetBandwidthId"] = d.Id()
	request["RegionId"] = client.RegionId

	request["Ipv6AddressId"] = d.Get("ipv6_address_id")
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = client.RpcPost("Vpc", "2016-04-28", action, nil, request, false)

		if err != nil {
			if NeedRetry(err) {
				wait()
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
