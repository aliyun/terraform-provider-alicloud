package alicloud

import (
	"fmt"
	"log"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAliCloudExpressConnectPhysicalConnection() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudExpressConnectPhysicalConnectionCreate,
		Read:   resourceAliCloudExpressConnectPhysicalConnectionRead,
		Update: resourceAliCloudExpressConnectPhysicalConnectionUpdate,
		Delete: resourceAliCloudExpressConnectPhysicalConnectionDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(2 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"access_point_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"line_operator": {
				Type:     schema.TypeString,
				Required: true,
			},
			"type": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"port_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: StringInSlice([]string{"1000Base-LX", "1000Base-T", "100Base-T", "10GBase-LR", "10GBase-T", "40GBase-LR", "100GBase-LR"}, false),
			},
			"bandwidth": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"circuit_code": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"peer_location": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"redundant_physical_connection_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"physical_connection_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"status": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: StringInSlice([]string{"Confirmed", "Enabled", "Canceled", "Terminated"}, false),
			},
			"period": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"pricing_cycle": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"order_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceAliCloudExpressConnectPhysicalConnectionCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	vpcService := VpcService{client}
	var response map[string]interface{}
	action := "CreatePhysicalConnection"
	request := make(map[string]interface{})
	var err error

	request["RegionId"] = client.RegionId
	request["ClientToken"] = buildClientToken("CreatePhysicalConnection")
	request["AccessPointId"] = d.Get("access_point_id")
	request["LineOperator"] = d.Get("line_operator")

	if v, ok := d.GetOk("type"); ok {
		request["Type"] = v
	}

	if v, ok := d.GetOk("port_type"); ok {
		request["PortType"] = v
	}

	if v, ok := d.GetOk("bandwidth"); ok {
		request["bandwidth"] = v
	}

	if v, ok := d.GetOk("circuit_code"); ok {
		request["CircuitCode"] = v
	}

	if v, ok := d.GetOk("peer_location"); ok {
		request["PeerLocation"] = v
	}

	if v, ok := d.GetOk("redundant_physical_connection_id"); ok {
		request["RedundantPhysicalConnectionId"] = v
	}

	if v, ok := d.GetOk("physical_connection_name"); ok {
		request["Name"] = v
	}

	if v, ok := d.GetOk("description"); ok {
		request["Description"] = v
	}

	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutCreate)), func() *resource.RetryError {
		response, err = client.RpcPost("Vpc", "2016-04-28", action, nil, request, true)
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_express_connect_physical_connection", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(response["PhysicalConnectionId"]))

	stateConf := BuildStateConf([]string{}, []string{"Allocated"}, d.Timeout(schema.TimeoutCreate), 1*time.Second, vpcService.ExpressConnectPhysicalConnectionStateRefreshFunc(d.Id(), []string{"Allocation Failed"}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return resourceAliCloudExpressConnectPhysicalConnectionUpdate(d, meta)
}

func resourceAliCloudExpressConnectPhysicalConnectionRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	vpcService := VpcService{client}

	object, err := vpcService.DescribeExpressConnectPhysicalConnection(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_express_connect_physical_connection vpcService.DescribeExpressConnectPhysicalConnection Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("access_point_id", object["AccessPointId"])
	d.Set("line_operator", object["LineOperator"])
	d.Set("type", object["Type"])
	d.Set("port_type", object["PortType"])
	d.Set("bandwidth", fmt.Sprint(formatInt(object["Bandwidth"])))
	d.Set("circuit_code", object["CircuitCode"])
	d.Set("peer_location", object["PeerLocation"])
	d.Set("redundant_physical_connection_id", object["RedundantPhysicalConnectionId"])
	d.Set("physical_connection_name", object["Name"])
	d.Set("description", object["Description"])
	d.Set("status", object["Status"])

	return nil
}

func resourceAliCloudExpressConnectPhysicalConnectionUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	vpcService := VpcService{client}
	var response map[string]interface{}
	var err error
	d.Partial(true)

	update := false
	request := map[string]interface{}{
		"RegionId":             client.RegionId,
		"ClientToken":          buildClientToken("ModifyPhysicalConnectionAttribute"),
		"PhysicalConnectionId": d.Id(),
	}

	if !d.IsNewResource() && d.HasChange("line_operator") {
		update = true
	}
	request["LineOperator"] = d.Get("line_operator")

	if !d.IsNewResource() && d.HasChange("bandwidth") {
		update = true

		if v, ok := d.GetOk("bandwidth"); ok {
			request["bandwidth"] = v
		}
	}

	if !d.IsNewResource() && d.HasChange("circuit_code") {
		update = true

		if v, ok := d.GetOk("circuit_code"); ok {
			request["CircuitCode"] = v
		}
	}

	if !d.IsNewResource() && d.HasChange("peer_location") {
		update = true

		if v, ok := d.GetOk("peer_location"); ok {
			request["PeerLocation"] = v
		}
	}

	if !d.IsNewResource() && d.HasChange("physical_connection_name") {
		update = true

		if v, ok := d.GetOk("physical_connection_name"); ok {
			request["Name"] = v
		}
	}

	if !d.IsNewResource() && d.HasChange("description") {
		update = true

		if v, ok := d.GetOk("description"); ok {
			request["Description"] = v
		}
	}

	if update {
		action := "ModifyPhysicalConnectionAttribute"
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *resource.RetryError {
			response, err = client.RpcPost("Vpc", "2016-04-28", action, nil, request, true)
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

		d.SetPartial("line_operator")
		d.SetPartial("port_type")
		d.SetPartial("bandwidth")
		d.SetPartial("circuit_code")
		d.SetPartial("peer_location")
		d.SetPartial("redundant_physical_connection_id")
		d.SetPartial("physical_connection_name")
		d.SetPartial("description")
	}

	if d.HasChange("status") {
		object, err := vpcService.DescribeExpressConnectPhysicalConnection(d.Id())
		if err != nil {
			return WrapError(err)
		}

		target := d.Get("status").(string)
		if object["Status"].(string) != target {

			if target == "Confirmed" {
				err := vpcService.ModifyExpressConnectPhysicalConnectionStatus(d, "Confirmed")
				if err != nil {
					return WrapError(err)
				}
			}

			if target == "Enabled" {
				if object["Status"].(string) != "Confirmed" {
					err := vpcService.ModifyExpressConnectPhysicalConnectionStatus(d, "Confirmed")
					if err != nil {
						return WrapError(err)
					}
				}

				err := vpcService.ModifyExpressConnectPhysicalConnectionStatus(d, "Enabled")
				if err != nil {
					return WrapError(err)
				}
			}

			if target == "Canceled" {
				err := vpcService.ModifyExpressConnectPhysicalConnectionStatus(d, "Canceled")
				if err != nil {
					return WrapError(err)
				}
			}

			if target == "Terminated" {
				err := vpcService.ModifyExpressConnectPhysicalConnectionStatus(d, "Terminated")
				if err != nil {
					return WrapError(err)
				}
			}

			d.SetPartial("status")
		}
	}

	d.Partial(false)

	return resourceAliCloudExpressConnectPhysicalConnectionRead(d, meta)
}

func resourceAliCloudExpressConnectPhysicalConnectionDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	action := "DeletePhysicalConnection"
	var response map[string]interface{}
	var err error

	request := map[string]interface{}{
		"RegionId":             client.RegionId,
		"ClientToken":          buildClientToken("DeletePhysicalConnection"),
		"PhysicalConnectionId": d.Id(),
	}

	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutDelete)), func() *resource.RetryError {
		response, err = client.RpcPost("Vpc", "2016-04-28", action, nil, request, true)
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

	return nil
}
