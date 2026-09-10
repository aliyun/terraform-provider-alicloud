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

func resourceAliCloudEnsNatGatewayForwardEntry() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudEnsNatGatewayForwardEntryCreate,
		Read:   resourceAliCloudEnsNatGatewayForwardEntryRead,
		Update: resourceAliCloudEnsNatGatewayForwardEntryUpdate,
		Delete: resourceAliCloudEnsNatGatewayForwardEntryDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"external_ip": {
				Type:     schema.TypeString,
				Required: true,
			},
			"external_port": {
				Type:     schema.TypeString,
				Required: true,
			},
			"forward_entry_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"health_check_port": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"internal_ip": {
				Type:     schema.TypeString,
				Required: true,
			},
			"internal_port": {
				Type:     schema.TypeString,
				Required: true,
			},
			"ip_protocol": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: StringInSlice([]string{"TCP", "UDP", "Any"}, false),
			},
			"nat_gateway_id": {
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

func resourceAliCloudEnsNatGatewayForwardEntryCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "CreateForwardEntry"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})

	request["ExternalIp"] = d.Get("external_ip")
	request["ExternalPort"] = d.Get("external_port")
	request["InternalIp"] = d.Get("internal_ip")
	request["InternalPort"] = d.Get("internal_port")
	if v, ok := d.GetOk("forward_entry_name"); ok {
		request["ForwardEntryName"] = v
	}
	if v, ok := d.GetOk("ip_protocol"); ok {
		request["IpProtocol"] = v
	}
	if v, ok := d.GetOkExists("health_check_port"); ok {
		request["HealthCheckPort"] = v
	}
	request["NatGatewayId"] = d.Get("nat_gateway_id")
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.RpcPost("Ens", "2017-11-10", action, query, request, true)
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_ens_nat_gateway_forward_entry", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(response["ForwardEntryId"]))

	ensServiceV2 := EnsServiceV2{client}
	stateConf := BuildStateConf([]string{}, []string{"Available"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, ensServiceV2.EnsNatGatewayForwardEntryStateRefreshFunc(d.Id(), "Status", []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return resourceAliCloudEnsNatGatewayForwardEntryRead(d, meta)
}

func resourceAliCloudEnsNatGatewayForwardEntryRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	ensServiceV2 := EnsServiceV2{client}

	objectRaw, err := ensServiceV2.DescribeEnsNatGatewayForwardEntry(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_ens_nat_gateway_forward_entry DescribeEnsNatGatewayForwardEntry Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("external_ip", objectRaw["ExternalIp"])
	d.Set("external_port", objectRaw["ExternalPort"])
	d.Set("forward_entry_name", objectRaw["ForwardEntryName"])
	d.Set("health_check_port", formatInt(objectRaw["HealthCheckPort"]))
	d.Set("internal_ip", objectRaw["InternalIp"])
	d.Set("internal_port", objectRaw["InternalPort"])
	d.Set("ip_protocol", objectRaw["IpProtocol"])
	d.Set("nat_gateway_id", objectRaw["NatGatewayId"])
	d.Set("status", objectRaw["Status"])

	return nil
}

func resourceAliCloudEnsNatGatewayForwardEntryUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	update := false

	var err error
	action := "ModifyForwardEntry"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["ForwardEntryId"] = d.Id()

	if d.HasChange("forward_entry_name") {
		update = true
		request["ForwardEntryName"] = d.Get("forward_entry_name")
	}

	if d.HasChange("health_check_port") {
		update = true
		request["HealthCheckPort"] = d.Get("health_check_port")
	}

	if d.HasChange("external_ip") {
		update = true
	}
	request["ExternalIp"] = d.Get("external_ip")
	if d.HasChange("external_port") {
		update = true
	}
	request["ExternalPort"] = d.Get("external_port")
	if d.HasChange("internal_ip") {
		update = true
	}
	request["InternalIp"] = d.Get("internal_ip")
	if d.HasChange("internal_port") {
		update = true
	}
	request["InternalPort"] = d.Get("internal_port")
	if d.HasChange("ip_protocol") {
		update = true
		request["IpProtocol"] = d.Get("ip_protocol")
	}

	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("Ens", "2017-11-10", action, query, request, true)
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
		ensServiceV2 := EnsServiceV2{client}
		stateConf := BuildStateConf([]string{}, []string{"Available"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, ensServiceV2.EnsNatGatewayForwardEntryStateRefreshFunc(d.Id(), "Status", []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
	}

	return resourceAliCloudEnsNatGatewayForwardEntryRead(d, meta)
}

func resourceAliCloudEnsNatGatewayForwardEntryDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	action := "DeleteForwardEntry"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	request["ForwardEntryId"] = d.Id()

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = client.RpcPost("Ens", "2017-11-10", action, query, request, true)
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
		if IsExpectedErrors(err, []string{"InvalidParameter.DnatNotFound"}) || NotFoundError(err) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	ensServiceV2 := EnsServiceV2{client}
	stateConf := BuildStateConf([]string{}, []string{""}, d.Timeout(schema.TimeoutDelete), 5*time.Second, ensServiceV2.EnsNatGatewayForwardEntryStateRefreshFunc(d.Id(), "$.ForwardEntryId", []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return nil
}
