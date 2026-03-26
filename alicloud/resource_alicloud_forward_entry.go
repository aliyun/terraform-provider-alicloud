package alicloud

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAliCloudNatGatewayForwardEntry() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudNatGatewayForwardEntryCreate,
		Read:   resourceAliCloudNatGatewayForwardEntryRead,
		Update: resourceAliCloudNatGatewayForwardEntryUpdate,
		Delete: resourceAliCloudNatGatewayForwardEntryDelete,
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
			"forward_entry_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"forward_entry_name": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"name"},
			},
			"forward_table_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
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
				Required:     true,
				ValidateFunc: StringInSlice([]string{"any", "tcp", "udp"}, false),
			},
			"port_break": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				Deprecated:    "Field `name` has been deprecated from provider version 1.119.1. New field `forward_entry_name` instead.",
				ConflictsWith: []string{"forward_entry_name"},
			},
		},
	}
}

func resourceAliCloudNatGatewayForwardEntryCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "CreateForwardEntry"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	if v, ok := d.GetOk("forward_table_id"); ok {
		request["ForwardTableId"] = v
	}
	request["RegionId"] = client.RegionId
	request["ClientToken"] = buildClientToken(action)

	request["InternalIp"] = d.Get("internal_ip")
	request["IpProtocol"] = d.Get("ip_protocol")
	if v, ok := d.GetOkExists("port_break"); ok {
		request["PortBreak"] = v
	}
	request["ExternalIp"] = d.Get("external_ip")
	request["InternalPort"] = d.Get("internal_port")
	if v, ok := d.GetOk("forward_entry_name"); ok {
		request["ForwardEntryName"] = v
	} else if v, ok := d.GetOk("name"); ok {
		request["ForwardEntryName"] = v
	}
	request["ExternalPort"] = d.Get("external_port")
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.RpcPost("Vpc", "2016-04-28", action, query, request, true)
		if err != nil {
			if IsExpectedErrors(err, []string{"OperationConflict", "TaskConflict", "OperationUnsupported.EipInBinding", "IncorrectStatus", "InvalidIp.NotInNatgw"}) || NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_forward_entry", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprintf("%v:%v", request["ForwardTableId"], response["ForwardEntryId"]))

	nATGatewayServiceV2 := NATGatewayServiceV2{client}
	stateConf := BuildStateConf([]string{}, []string{"Available"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, nATGatewayServiceV2.NatGatewayForwardEntryStateRefreshFunc(d.Id(), "Status", []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return resourceAliCloudNatGatewayForwardEntryRead(d, meta)
}

func resourceAliCloudNatGatewayForwardEntryRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	nATGatewayServiceV2 := NATGatewayServiceV2{client}

	if !strings.Contains(d.Id(), ":") {
		d.SetId(fmt.Sprintf("%v:%v", d.Get("forward_table_id"), d.Id()))
	}

	objectRaw, err := nATGatewayServiceV2.DescribeNatGatewayForwardEntry(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_forward_entry DescribeNatGatewayForwardEntry Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("external_ip", objectRaw["ExternalIp"])
	d.Set("external_port", objectRaw["ExternalPort"])
	d.Set("forward_entry_name", objectRaw["ForwardEntryName"])
	d.Set("internal_ip", objectRaw["InternalIp"])
	d.Set("internal_port", objectRaw["InternalPort"])
	d.Set("ip_protocol", objectRaw["IpProtocol"])
	d.Set("status", objectRaw["Status"])
	d.Set("forward_entry_id", objectRaw["ForwardEntryId"])
	d.Set("forward_table_id", objectRaw["ForwardTableId"])
	d.Set("name", objectRaw["ForwardEntryName"])

	return nil
}

func resourceAliCloudNatGatewayForwardEntryUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	update := false

	if !strings.Contains(d.Id(), ":") {
		d.SetId(fmt.Sprintf("%v:%v", d.Get("forward_table_id"), d.Id()))
	}

	var err error
	parts := strings.Split(d.Id(), ":")
	action := "ModifyForwardEntry"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["ForwardTableId"] = parts[0]
	request["ForwardEntryId"] = parts[1]
	request["RegionId"] = client.RegionId
	request["ClientToken"] = buildClientToken(action)
	if d.HasChange("internal_ip") {
		update = true
	}
	request["InternalIp"] = d.Get("internal_ip")
	if d.HasChange("ip_protocol") {
		update = true
	}
	request["IpProtocol"] = d.Get("ip_protocol")
	if v, ok := d.GetOkExists("port_break"); ok {
		request["PortBreak"] = v
	}
	if d.HasChange("external_ip") {
		update = true
	}
	request["ExternalIp"] = d.Get("external_ip")
	if d.HasChange("internal_port") {
		update = true
	}
	request["InternalPort"] = d.Get("internal_port")
	if d.HasChange("forward_entry_name") {
		update = true

		if v, ok := d.GetOk("forward_entry_name"); ok {
			request["ForwardEntryName"] = v
		}
	}

	if d.HasChange("external_port") {
		update = true
	}
	request["ExternalPort"] = d.Get("external_port")

	if d.HasChange("name") {
		update = true

		if v, ok := d.GetOk("name"); ok {
			request["ForwardEntryName"] = v
		}
	}
	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("Vpc", "2016-04-28", action, query, request, true)
			if err != nil {
				if IsExpectedErrors(err, []string{"OperationUnsupported.EipInBinding", "IncorretForwardEntryStatus"}) || NeedRetry(err) {
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
		nATGatewayServiceV2 := NATGatewayServiceV2{client}
		stateConf := BuildStateConf([]string{}, []string{"Available"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, nATGatewayServiceV2.NatGatewayForwardEntryStateRefreshFunc(d.Id(), "Status", []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
	}

	return resourceAliCloudNatGatewayForwardEntryRead(d, meta)
}

func resourceAliCloudNatGatewayForwardEntryDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	if !strings.Contains(d.Id(), ":") {
		d.SetId(fmt.Sprintf("%v:%v", d.Get("forward_table_id"), d.Id()))
	}
	parts := strings.Split(d.Id(), ":")
	action := "DeleteForwardEntry"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	request["ForwardTableId"] = parts[0]
	request["ForwardEntryId"] = parts[1]
	request["RegionId"] = client.RegionId
	request["ClientToken"] = buildClientToken(action)

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = client.RpcPost("Vpc", "2016-04-28", action, query, request, true)
		if err != nil {
			if IsExpectedErrors(err, []string{"IncorrectStatus.NATGW", "OperationConflict", "UnknownError"}) || NeedRetry(err) {
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

	nATGatewayServiceV2 := NATGatewayServiceV2{client}
	stateConf := BuildStateConf([]string{}, []string{""}, d.Timeout(schema.TimeoutDelete), 5*time.Second, nATGatewayServiceV2.NatGatewayForwardEntryStateRefreshFunc(d.Id(), "Status", []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return nil
}
