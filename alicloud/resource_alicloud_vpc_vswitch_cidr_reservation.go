// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
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

func resourceAliCloudVpcVswitchCidrReservation() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudVpcVswitchCidrReservationCreate,
		Read:   resourceAlicloudVpcVswitchCidrReservationRead,
		Update: resourceAlicloudVpcVswitchCidrReservationUpdate,
		Delete: resourceAlicloudVpcVswitchCidrReservationDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"cidr_reservation_cidr": {
				Type:         schema.TypeString,
				ExactlyOneOf: []string{"cidr_reservation_cidr", "cidr_reservation_mask"},
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
			},
			"cidr_reservation_description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"cidr_reservation_mask": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"cidr_reservation_type": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				ValidateFunc: StringInSlice([]string{"Prefix"}, false),
			},
			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"ip_version": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				ValidateFunc: StringInSlice([]string{"IPv4", "IPv6"}, false),
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"vswitch_cidr_reservation_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"vswitch_cidr_reservation_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"vswitch_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"vpc_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceAlicloudVpcVswitchCidrReservationCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	action := "CreateVSwitchCidrReservation"
	var request map[string]interface{}
	var response map[string]interface{}
	var err error
	request = make(map[string]interface{})
	request["VSwitchId"] = d.Get("vswitch_id")
	request["RegionId"] = client.RegionId
	request["ClientToken"] = buildClientToken(action)

	if v, ok := d.GetOk("ip_version"); ok {
		request["IpVersion"] = v
	}

	if v, ok := d.GetOk("cidr_reservation_description"); ok {
		request["VSwitchCidrReservationDescription"] = v
	}

	if v, ok := d.GetOk("cidr_reservation_cidr"); ok {
		request["VSwitchCidrReservationCidr"] = v
	}

	if v, ok := d.GetOk("vswitch_cidr_reservation_name"); ok {
		request["VSwitchCidrReservationName"] = v
	}

	if v, ok := d.GetOk("cidr_reservation_mask"); ok {
		request["VSwitchCidrReservationMask"] = v
	}

	if v, ok := d.GetOk("cidr_reservation_type"); ok {
		request["VSwitchCidrReservationType"] = v
	}

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.RpcPost("Vpc", "2016-04-28", action, nil, request, true)
		request["ClientToken"] = buildClientToken(action)

		if err != nil {
			if IsExpectedErrors(err, []string{"IncorrectStatus.Vpc", "OperationConflict", "IncorrectStatus", "ServiceUnavailable", "SystemBusy", "LastTokenProcessing"}) || NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(action, response, request)
		return nil
	})

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_vpc_vswitch_cidr_reservation", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprintf("%v:%v", request["VSwitchId"], response["VSwitchCidrReservationId"]))

	vpcServiceV2 := VpcServiceV2{client}
	stateConf := BuildStateConf([]string{}, []string{"Assigned"}, d.Timeout(schema.TimeoutCreate), 0, vpcServiceV2.VpcVswitchCidrReservationStateRefreshFunc(d.Id(), "Status", []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return resourceAlicloudVpcVswitchCidrReservationRead(d, meta)
}

func resourceAlicloudVpcVswitchCidrReservationRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	vpcServiceV2 := VpcServiceV2{client}

	objectRaw, err := vpcServiceV2.DescribeVpcVswitchCidrReservation(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_vpc_vswitch_cidr_reservation DescribeVpcVswitchCidrReservation Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("cidr_reservation_cidr", objectRaw["VSwitchCidrReservationCidr"])
	d.Set("cidr_reservation_description", objectRaw["VSwitchCidrReservationDescription"])
	d.Set("cidr_reservation_type", objectRaw["Type"])
	d.Set("create_time", objectRaw["CreationTime"])
	d.Set("ip_version", objectRaw["IpVersion"])
	d.Set("status", objectRaw["Status"])
	d.Set("vswitch_cidr_reservation_name", objectRaw["VSwitchCidrReservationName"])
	d.Set("vpc_id", objectRaw["VpcId"])
	d.Set("vswitch_cidr_reservation_id", objectRaw["VSwitchCidrReservationId"])
	d.Set("vswitch_id", objectRaw["VSwitchId"])

	return nil
}

func resourceAlicloudVpcVswitchCidrReservationUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	update := false
	action := "ModifyVSwitchCidrReservationAttribute"
	var err error
	request = make(map[string]interface{})

	parts := strings.Split(d.Id(), ":")

	request["VSwitchCidrReservationId"] = parts[1]
	request["RegionId"] = client.RegionId

	if d.HasChange("cidr_reservation_description") {
		update = true
		if v, ok := d.GetOk("cidr_reservation_description"); ok {
			request["VSwitchCidrReservationDescription"] = v
		}
	}

	if d.HasChange("vswitch_cidr_reservation_name") {
		update = true
		if v, ok := d.GetOk("vswitch_cidr_reservation_name"); ok {
			request["VSwitchCidrReservationName"] = v
		}
	}

	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
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
		d.SetPartial("cidr_reservation_description")
		d.SetPartial("vswitch_cidr_reservation_name")
	}

	return resourceAlicloudVpcVswitchCidrReservationRead(d, meta)
}

func resourceAlicloudVpcVswitchCidrReservationDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "DeleteVSwitchCidrReservation"
	var request map[string]interface{}
	var response map[string]interface{}
	var err error
	request = make(map[string]interface{})

	parts := strings.Split(d.Id(), ":")

	request["VSwitchCidrReservationId"] = parts[1]
	request["RegionId"] = client.RegionId

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = client.RpcPost("Vpc", "2016-04-28", action, nil, request, false)

		if err != nil {
			if IsExpectedErrors(err, []string{"IncorrectStatus.Vpc", "OperationConflict", "IncorrectStatus", "ServiceUnavailable", "SystemBusy", "LastTokenProcessing"}) || NeedRetry(err) {
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

	vpcServiceV2 := VpcServiceV2{client}
	stateConf := BuildStateConf([]string{}, []string{}, d.Timeout(schema.TimeoutDelete), 0, vpcServiceV2.VpcVswitchCidrReservationStateRefreshFunc(d.Id(), "Status", []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}
	return nil
}
