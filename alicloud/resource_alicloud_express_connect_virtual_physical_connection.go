package alicloud

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAlicloudExpressConnectVirtualPhysicalConnection() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudExpressConnectVirtualPhysicalConnectionCreate,
		Read:   resourceAlicloudExpressConnectVirtualPhysicalConnectionRead,
		Update: resourceAlicloudExpressConnectVirtualPhysicalConnectionUpdate,
		Delete: resourceAlicloudExpressConnectVirtualPhysicalConnectionDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Hour),
			Update: schema.DefaultTimeout(5 * time.Hour),
		},
		Schema: map[string]*schema.Schema{
			"access_point_id": {
				Computed: true,
				Type:     schema.TypeString,
			},
			"ad_location": {
				Computed: true,
				Type:     schema.TypeString,
			},
			"bandwidth": {
				Computed: true,
				Type:     schema.TypeString,
			},
			"business_status": {
				Computed: true,
				Type:     schema.TypeString,
			},
			"circuit_code": {
				Computed: true,
				Type:     schema.TypeString,
			},
			"create_time": {
				Computed: true,
				Type:     schema.TypeString,
			},
			"description": {
				Optional: true,
				ForceNew: true,
				Type:     schema.TypeString,
			},
			"dry_run": {
				Optional: true,
				Type:     schema.TypeBool,
			},
			"enabled_time": {
				Computed: true,
				Type:     schema.TypeString,
			},
			"end_time": {
				Computed: true,
				Type:     schema.TypeString,
			},
			"expect_spec": {
				Optional: true,
				Type:     schema.TypeString,
			},
			"line_operator": {
				Computed: true,
				Type:     schema.TypeString,
			},
			"loa_status": {
				Computed: true,
				Type:     schema.TypeString,
			},
			"order_mode": {
				Required: true,
				ForceNew: true,
				Type:     schema.TypeString,
			},
			"parent_physical_connection_ali_uid": {
				Computed: true,
				Type:     schema.TypeString,
			},
			"parent_physical_connection_id": {
				Required: true,
				ForceNew: true,
				Type:     schema.TypeString,
			},
			"peer_location": {
				Computed: true,
				Type:     schema.TypeString,
			},
			"port_number": {
				Computed: true,
				Type:     schema.TypeString,
			},
			"port_type": {
				Computed: true,
				Type:     schema.TypeString,
			},
			"redundant_physical_connection_id": {
				Computed: true,
				Type:     schema.TypeString,
			},
			"resource_group_id": {
				Optional: true,
				Computed: true,
				ForceNew: true,
				Type:     schema.TypeString,
			},
			"spec": {
				Required:     true,
				ForceNew:     true,
				Type:         schema.TypeString,
				ValidateFunc: validation.StringInSlice([]string{"50M", "100M", "200M", "300M", "400M", "500M", "1G", "2G", "5G", "8G", "10G"}, false),
			},
			"status": {
				Computed: true,
				Type:     schema.TypeString,
			},
			"virtual_physical_connection_name": {
				Optional: true,
				ForceNew: true,
				Type:     schema.TypeString,
			},
			"virtual_physical_connection_status": {
				Computed: true,
				Type:     schema.TypeString,
			},
			"vlan_id": {
				Required:     true,
				Type:         schema.TypeInt,
				ValidateFunc: validation.IntBetween(0, 2999),
			},
			"vpconn_ali_uid": {
				Required: true,
				ForceNew: true,
				Type:     schema.TypeString,
			},
		},
	}
}

func resourceAlicloudExpressConnectVirtualPhysicalConnectionCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	request := map[string]interface{}{
		"RegionId": client.RegionId,
	}
	var err error

	if v, ok := d.GetOk("description"); ok {
		request["Description"] = v
	}
	if v, ok := d.GetOk("dry_run"); ok {
		request["DryRun"] = v
	}
	if v, ok := d.GetOk("order_mode"); ok {
		request["OrderMode"] = v
	}
	if v, ok := d.GetOk("resource_group_id"); ok {
		request["ResourceGroupId"] = v
	}
	if v, ok := d.GetOk("parent_physical_connection_id"); ok {
		request["PhysicalConnectionId"] = v
	}
	if v, ok := d.GetOk("spec"); ok {
		request["Spec"] = v
	}
	if v, ok := d.GetOk("virtual_physical_connection_name"); ok {
		request["Name"] = v
	}
	if v, ok := d.GetOk("vlan_id"); ok {
		request["VlanId"] = v
	}
	if v, ok := d.GetOk("vpconn_ali_uid"); ok {
		request["VpconnAliUid"] = v
	}

	var response map[string]interface{}
	action := "CreateVirtualPhysicalConnection"
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutCreate)), func() *resource.RetryError {
		resp, err := client.RpcPost("Vpc", "2016-04-28", action, nil, request, false)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		response = resp
		addDebug(action, response, request)
		return nil
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_express_connect_virtual_physical_connection", action, AlibabaCloudSdkGoERROR)
	}

	if v, err := jsonpath.Get("$.VirtualPhysicalConnection", response); err != nil || v == nil {
		return WrapErrorf(err, IdMsg, "alicloud_express_connect_virtual_physical_connection")
	} else {
		d.SetId(fmt.Sprint(v))
	}

	return resourceAlicloudExpressConnectVirtualPhysicalConnectionUpdate(d, meta)
}

func resourceAlicloudExpressConnectVirtualPhysicalConnectionRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	expressConnectService := ExpressConnectService{client}

	object, err := expressConnectService.DescribeExpressConnectVirtualPhysicalConnection(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_express_connect_virtual_physical_connection vpcService.DescribeExpressConnectVirtualPhysicalConnection Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	d.Set("access_point_id", object["AccessPointId"])
	d.Set("ad_location", object["AdLocation"])
	d.Set("bandwidth", object["Bandwidth"])
	d.Set("business_status", object["BusinessStatus"])
	d.Set("circuit_code", object["CircuitCode"])
	d.Set("create_time", object["CreationTime"])
	d.Set("description", object["Description"])
	d.Set("enabled_time", object["EnabledTime"])
	d.Set("end_time", object["EndTime"])
	d.Set("expect_spec", object["ExpectSpec"])
	d.Set("line_operator", object["LineOperator"])
	d.Set("loa_status", object["LoaStatus"])
	d.Set("order_mode", object["OrderMode"])
	d.Set("parent_physical_connection_ali_uid", object["ParentPhysicalConnectionAliUid"])
	d.Set("parent_physical_connection_id", object["ParentPhysicalConnectionId"])
	d.Set("peer_location", object["PeerLocation"])
	d.Set("port_number", object["PortNumber"])
	d.Set("port_type", object["PortType"])
	d.Set("redundant_physical_connection_id", object["RedundantPhysicalConnectionId"])
	d.Set("resource_group_id", object["ResourceGroupId"])
	d.Set("spec", object["Spec"])
	d.Set("status", object["Status"])
	d.Set("virtual_physical_connection_name", object["Name"])
	d.Set("virtual_physical_connection_status", object["VirtualPhysicalConnectionStatus"])
	d.Set("vlan_id", formatInt(object["VlanId"]))
	d.Set("vpconn_ali_uid", object["AliUid"])

	return nil
}

func resourceAlicloudExpressConnectVirtualPhysicalConnectionUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	var err error
	update := false
	request := map[string]interface{}{
		"InstanceId": d.Id(),
		"RegionId":   client.RegionId,
	}

	if v, ok := d.GetOk("dry_run"); ok {
		request["DryRun"] = v
	}
	if d.HasChange("expect_spec") {
		update = true
		if v, ok := d.GetOk("expect_spec"); ok {
			request["ExpectSpec"] = v
		}
	}
	if !d.IsNewResource() && d.HasChange("vlan_id") {
		update = true
	}
	request["VlanId"] = d.Get("vlan_id")

	if update {
		action := "UpdateVirtualPhysicalConnection"
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *resource.RetryError {
			resp, err := client.RpcPost("Vpc", "2016-04-28", action, nil, request, false)
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			addDebug(action, resp, request)
			return nil
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
	}

	return resourceAlicloudExpressConnectVirtualPhysicalConnectionRead(d, meta)
}

func resourceAlicloudExpressConnectVirtualPhysicalConnectionDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var err error

	request := map[string]interface{}{

		"PhysicalConnectionId": d.Id(),
		"RegionId":             client.RegionId,
	}

	action := "DeletePhysicalConnection"
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutDelete)), func() *resource.RetryError {
		request["ClientToken"] = buildClientToken("DeletePhysicalConnection")
		resp, err := client.RpcPost("Vpc", "2016-04-28", action, nil, request, true)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(action, resp, request)
		return nil
	})
	if err != nil {
		if NotFoundError(err) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}
	return nil
}
