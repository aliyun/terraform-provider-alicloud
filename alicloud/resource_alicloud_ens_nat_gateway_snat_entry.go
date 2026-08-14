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

func resourceAliCloudEnsNatGatewaySnatEntry() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudEnsNatGatewaySnatEntryCreate,
		Read:   resourceAliCloudEnsNatGatewaySnatEntryRead,
		Update: resourceAliCloudEnsNatGatewaySnatEntryUpdate,
		Delete: resourceAliCloudEnsNatGatewaySnatEntryDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"eip_affinity": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"idle_timeout": {
				Type:         schema.TypeInt,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: IntBetween(0, 86400),
			},
			"isp_affinity": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"nat_gateway_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"snat_entry_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"snat_ip": {
				Type:     schema.TypeString,
				Required: true,
			},
			"source_cidr": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceAliCloudEnsNatGatewaySnatEntryCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "CreateSnatEntry"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})

	request["NatGatewayId"] = d.Get("nat_gateway_id")
	request["SnatIp"] = d.Get("snat_ip")
	if v, ok := d.GetOk("snat_entry_name"); ok {
		request["SnatEntryName"] = v
	}
	if v, ok := d.GetOk("source_cidr"); ok {
		request["SourceCIDR"] = v
	}
	if v, ok := d.GetOkExists("idle_timeout"); ok {
		request["IdleTimeout"] = v
	}
	if v, ok := d.GetOkExists("isp_affinity"); ok {
		request["IspAffinity"] = v
	}
	if v, ok := d.GetOkExists("eip_affinity"); ok {
		request["EipAffinity"] = v
	}
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_ens_nat_gateway_snat_entry", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(response["SnatEntryId"]))

	ensServiceV2 := EnsServiceV2{client}
	stateConf := BuildStateConf([]string{}, []string{"Available"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, ensServiceV2.EnsNatGatewaySnatEntryStateRefreshFunc(d.Id(), "Status", []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return resourceAliCloudEnsNatGatewaySnatEntryRead(d, meta)
}

func resourceAliCloudEnsNatGatewaySnatEntryRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	ensServiceV2 := EnsServiceV2{client}

	objectRaw, err := ensServiceV2.DescribeEnsNatGatewaySnatEntry(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_ens_nat_gateway_snat_entry DescribeEnsNatGatewaySnatEntry Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("eip_affinity", objectRaw["EipAffinity"])
	d.Set("idle_timeout", formatInt(objectRaw["IdleTimeout"]))
	d.Set("isp_affinity", objectRaw["IspAffinity"])
	d.Set("nat_gateway_id", objectRaw["NatGatewayId"])
	d.Set("snat_entry_name", objectRaw["SnatEntryName"])
	d.Set("snat_ip", objectRaw["SnatIp"])
	d.Set("source_cidr", objectRaw["SourceCIDR"])
	d.Set("status", objectRaw["Status"])

	return nil
}

func resourceAliCloudEnsNatGatewaySnatEntryUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	update := false

	var err error
	action := "ModifySnatEntry"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["SnatEntryId"] = d.Id()

	if d.HasChange("snat_entry_name") {
		update = true
		request["SnatEntryName"] = d.Get("snat_entry_name")
	}

	if d.HasChange("isp_affinity") {
		update = true
		request["IspAffinity"] = d.Get("isp_affinity")
	}

	if d.HasChange("eip_affinity") {
		update = true
		request["EipAffinity"] = d.Get("eip_affinity")
	}

	if d.HasChange("snat_ip") {
		update = true
	}
	request["SnatIp"] = d.Get("snat_ip")
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
		stateConf := BuildStateConf([]string{}, []string{"Available"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, ensServiceV2.EnsNatGatewaySnatEntryStateRefreshFunc(d.Id(), "Status", []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
	}

	return resourceAliCloudEnsNatGatewaySnatEntryRead(d, meta)
}

func resourceAliCloudEnsNatGatewaySnatEntryDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	action := "DeleteSnatEntry"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	request["SnatEntryId"] = d.Id()

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
		if IsExpectedErrors(err, []string{"InvalidParameter.SnatNotFound"}) || NotFoundError(err) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	ensServiceV2 := EnsServiceV2{client}
	stateConf := BuildStateConf([]string{}, []string{""}, d.Timeout(schema.TimeoutDelete), 5*time.Second, ensServiceV2.EnsNatGatewaySnatEntryStateRefreshFunc(d.Id(), "$.SnatEntryId", []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return nil
}
