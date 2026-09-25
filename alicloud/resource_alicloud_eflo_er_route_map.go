// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
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

func resourceAliCloudEfloErRouteMap() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudEfloErRouteMapCreate,
		Read:   resourceAliCloudEfloErRouteMapRead,
		Update: resourceAliCloudEfloErRouteMapUpdate,
		Delete: resourceAliCloudEfloErRouteMapDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"action": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: StringInSlice([]string{"permit", "deny"}, false),
			},
			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"destination_cidr_block": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"er_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"er_route_map_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"er_route_map_num": {
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
			},
			"reception_instance_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"reception_instance_owner": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"reception_instance_type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: StringInSlice([]string{"VPD", "VCC"}, false),
			},
			"region_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"transmission_instance_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"transmission_instance_owner": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"transmission_instance_type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: StringInSlice([]string{"VPD", "VCC"}, false),
			},
		},
	}
}

func resourceAliCloudEfloErRouteMapCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "CreateErRouteMap"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	request["ErId"] = d.Get("er_id")
	request["RegionId"] = client.RegionId
	if v, ok := d.GetOk("transmission_instance_owner"); ok {
		request["TransmissionInstanceOwner"] = v
	}
	if v, ok := d.GetOk("description"); ok {
		request["Description"] = v
	}
	request["ReceptionInstanceId"] = d.Get("reception_instance_id")
	request["ReceptionInstanceType"] = d.Get("reception_instance_type")
	request["TransmissionInstanceType"] = d.Get("transmission_instance_type")
	request["RouteMapNum"] = d.Get("er_route_map_num")
	request["TransmissionInstanceId"] = d.Get("transmission_instance_id")
	request["RouteMapAction"] = d.Get("action")
	request["DestinationCidrBlock"] = d.Get("destination_cidr_block")
	if v, ok := d.GetOk("reception_instance_owner"); ok {
		request["ReceptionInstanceOwner"] = v
	}
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.RpcPost("eflo", "2022-05-30", action, query, request, true)
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_eflo_er_route_map", action, AlibabaCloudSdkGoERROR)
	}

	contentErRouteMapId, err := jsonpath.Get("$.Content.ErRouteMapId", response)
	if err != nil {
		return WrapErrorf(err, FailedGetAttributeMsg, d.Id(), "$.Content.ErRouteMapId", response)
	}
	d.SetId(fmt.Sprintf("%v:%v", request["ErId"], contentErRouteMapId))

	efloServiceV2 := EfloServiceV2{client}
	stateConf := BuildStateConf([]string{}, []string{"Available"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, efloServiceV2.EfloErRouteMapStateRefreshFunc(d.Id(), "Status", []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return resourceAliCloudEfloErRouteMapRead(d, meta)
}

func resourceAliCloudEfloErRouteMapRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	efloServiceV2 := EfloServiceV2{client}

	objectRaw, err := efloServiceV2.DescribeEfloErRouteMap(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_eflo_er_route_map DescribeEfloErRouteMap Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("action", objectRaw["Action"])
	d.Set("create_time", objectRaw["GmtCreate"])
	d.Set("description", objectRaw["Description"])
	d.Set("destination_cidr_block", objectRaw["DestinationCidrBlock"])
	d.Set("er_route_map_num", objectRaw["RouteMapNum"])
	d.Set("reception_instance_id", objectRaw["ReceptionInstanceId"])
	d.Set("reception_instance_owner", objectRaw["ReceptionInstanceOwner"])
	d.Set("reception_instance_type", objectRaw["ReceptionInstanceType"])
	d.Set("region_id", objectRaw["RegionId"])
	d.Set("status", objectRaw["Status"])
	d.Set("transmission_instance_id", objectRaw["TransmissionInstanceId"])
	d.Set("transmission_instance_owner", objectRaw["TransmissionInstanceOwner"])
	d.Set("transmission_instance_type", objectRaw["TransmissionInstanceType"])
	d.Set("er_id", objectRaw["ErId"])
	d.Set("er_route_map_id", objectRaw["ErRouteMapId"])

	return nil
}

func resourceAliCloudEfloErRouteMapUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	update := false

	var err error
	parts := strings.Split(d.Id(), ":")
	action := "UpdateErRouteMap"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	if len(parts) < 2 {
		return WrapError(fmt.Errorf("invalid Resource Id %s. Expected format <ErId>:<ErRouteMapId>", d.Id()))
	}
	request["ErRouteMapId"] = parts[1]
	request["ErId"] = parts[0]
	request["RegionId"] = client.RegionId
	if d.HasChange("description") {
		update = true
		request["Description"] = d.Get("description")
	}

	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("eflo", "2022-05-30", action, query, request, true)
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

	return resourceAliCloudEfloErRouteMapRead(d, meta)
}

func resourceAliCloudEfloErRouteMapDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	parts := strings.Split(d.Id(), ":")
	if len(parts) < 2 {
		return WrapError(fmt.Errorf("invalid Resource Id %s. Expected format <ErId>:<ErRouteMapId>", d.Id()))
	}
	action := "DeleteErRouteMap"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	request["ErRouteMapIds.1"] = parts[1]
	request["ErId"] = parts[0]
	request["RegionId"] = client.RegionId

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = client.RpcPost("eflo", "2022-05-30", action, query, request, true)
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
		if IsExpectedErrors(err, []string{"1003"}) || NotFoundError(err) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	return nil
}
