package alicloud

import (
	"fmt"
	"log"
	"time"

	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAlicloudVswitch() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudVswitchCreate,
		Read:   resourceAlicloudVswitchRead,
		Update: resourceAlicloudVswitchUpdate,
		Delete: resourceAlicloudVswitchDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"cidr_block": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"tags": tagsSchema(),
			"vswitch_name": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"name"},
			},
			"name": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"vswitch_name"},
				Deprecated:    "Field 'name' has been deprecated from provider version 1.119.0. New field 'vswitch_name' instead.",
			},
			"vpc_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"zone_id": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ForceNew:      true,
				ConflictsWith: []string{"availability_zone"},
			},
			"availability_zone": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ForceNew:      true,
				ConflictsWith: []string{"zone_id"},
				Deprecated:    "Field `availability_zone` has been deprecated from provider version 1.119.0. New field `zone_id` instead.",
			},
			"enable_ipv6": {
				Type:       schema.TypeBool,
				Optional:   true,
				Computed:   true,
				Deprecated: "Field `enable_ipv6` has been deprecated from provider version 1.203.0. If you want to enable ipv6, please set `ipv6_cidr_block_mask` to value, and if you want to disable ipv6, please set `ipv6_cidr_block_mask` to `-1`",
			},
			"ipv6_cidr_block_mask": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"ipv6_cidr_block": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceAlicloudVswitchCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	vpcService := VpcService{client}
	var response map[string]interface{}
	action := "CreateVSwitch"
	request := make(map[string]interface{})
	conn, err := client.NewVpcClient()
	if err != nil {
		return WrapError(err)
	}
	request["CidrBlock"] = d.Get("cidr_block")
	if v, ok := d.GetOk("description"); ok {
		request["Description"] = v
	}

	request["RegionId"] = client.RegionId
	if v, ok := d.GetOk("vswitch_name"); ok {
		request["VSwitchName"] = v
	} else if v, ok := d.GetOk("name"); ok {
		request["VSwitchName"] = v
	}

	request["VpcId"] = d.Get("vpc_id")
	if v, ok := d.GetOk("zone_id"); ok {
		request["ZoneId"] = v
	} else if v, ok := d.GetOk("availability_zone"); ok {
		request["ZoneId"] = v
	} else {
		return WrapError(Error(`[ERROR] Argument "availability_zone" or "zone_id" must be set one!`))
	}

	if v, ok := d.GetOkExists("ipv6_cidr_block_mask"); ok && fmt.Sprint(v) != "-1" {
		request["Ipv6CidrBlock"] = v
	}

	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutCreate)), func() *resource.RetryError {
		request["ClientToken"] = buildClientToken("CreateVSwitch")
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2016-04-28"), StringPointer("AK"), nil, request, &runtime)
		if err != nil {
			if IsExpectedErrors(err, []string{"InvalidStatus.RouteEntry", "OperationFailed.IdempotentTokenProcessing", "TaskConflict", "CreateVSwitch.IncorrectStatus.cbnStatus", "IncorrectStatus.cbnStatus", "IncorrectStatus.%s", "IncorrectVSwitchStatus", "OperationConflict", "OperationFailed.DistibuteLock"}) || NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_vswitch", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(response["VSwitchId"]))

	stateConf := BuildStateConf([]string{}, []string{"Available"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, vpcService.VswitchStateRefreshFunc(d.Id(), []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return resourceAlicloudVswitchUpdate(d, meta)
}

func resourceAlicloudVswitchRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	vpcService := VpcService{client}
	object, err := vpcService.DescribeVswitch(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_vswitch vpcService.DescribeVswitch Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("cidr_block", object["CidrBlock"])
	d.Set("description", object["Description"])
	d.Set("status", object["Status"])
	d.Set("vswitch_name", object["VSwitchName"])
	d.Set("name", object["VSwitchName"])
	d.Set("vpc_id", object["VpcId"])
	d.Set("zone_id", object["ZoneId"])
	d.Set("availability_zone", object["ZoneId"])

	listTagResourcesObject, err := vpcService.ListTagResources(d.Id(), "VSWITCH")
	if err != nil {
		return WrapError(err)
	}
	d.Set("tags", tagsToMap(listTagResourcesObject))

	if v, ok := object["Ipv6CidrBlock"]; ok && fmt.Sprint(v) != "" {
		ipv6CidrBlockMask, err := getIpv6CidrNum(v.(string))
		if err != nil {
			return WrapError(err)
		}

		d.Set("enable_ipv6", true)
		d.Set("ipv6_cidr_block", object["Ipv6CidrBlock"])
		d.Set("ipv6_cidr_block_mask", formatInt(ipv6CidrBlockMask))
	} else {
		d.Set("enable_ipv6", false)
		d.Set("ipv6_cidr_block", object["Ipv6CidrBlock"])
		d.Set("ipv6_cidr_block_mask", -1)
	}

	return nil
}

func resourceAlicloudVswitchUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	vpcService := VpcService{client}
	var response map[string]interface{}
	d.Partial(true)
	conn, err := client.NewVpcClient()
	if err != nil {
		return WrapError(err)
	}

	if d.HasChange("tags") {
		if err := vpcService.SetResourceTags(d, "VSWITCH"); err != nil {
			return WrapError(err)
		}
		d.SetPartial("tags")
	}

	update := false
	request := map[string]interface{}{
		"VSwitchId": d.Id(),
	}
	request["RegionId"] = client.RegionId

	if !d.IsNewResource() && d.HasChange("description") {
		update = true
	}
	if v, ok := d.GetOk("description"); ok {
		request["Description"] = v
	}

	if !d.IsNewResource() && d.HasChange("vswitch_name") {
		update = true
		if v, ok := d.GetOk("vswitch_name"); ok {
			request["VSwitchName"] = v
		}
	}

	if !d.IsNewResource() && d.HasChange("name") {
		update = true
		if v, ok := d.GetOk("name"); ok {
			request["VSwitchName"] = v
		}
	}

	if !d.IsNewResource() && d.HasChange("ipv6_cidr_block_mask") {
		err := vpcService.CancelIpv6(d)
		if err != nil {
			return WrapError(err)
		}

		if v, ok := d.GetOkExists("ipv6_cidr_block_mask"); ok && fmt.Sprint(v) != "-1" {
			update = true
			request["Ipv6CidrBlock"] = v
		}
	}

	if update {
		action := "ModifyVSwitchAttribute"
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2016-04-28"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
			if err != nil {
				if IsExpectedErrors(err, []string{"OperationConflict", "OperationFailed.LastTokenProcessing", "IncorrectStatus.VSwitch", "IncorrectStatus.VpcRouteEntry", "ServiceUnavailable"}) || NeedRetry(err) {
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

		stateConf := BuildStateConf([]string{}, []string{"Available"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, vpcService.VswitchStateRefreshFunc(d.Id(), []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}

		d.SetPartial("description")
		d.SetPartial("name")
		d.SetPartial("vswitch_name")
		d.SetPartial("ipv6_cidr_block_mask")
	}

	d.Partial(false)

	return resourceAlicloudVswitchRead(d, meta)
}

func resourceAlicloudVswitchDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	vpcService := VpcService{client}
	action := "DeleteVSwitch"
	var response map[string]interface{}
	conn, err := client.NewVpcClient()
	if err != nil {
		return WrapError(err)
	}
	request := map[string]interface{}{
		"VSwitchId": d.Id(),
	}

	request["RegionId"] = client.RegionId
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutDelete)), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2016-04-28"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
		if err != nil {
			if IsExpectedErrors(err, []string{"InvalidVSwitchId.NotFound", "InvalidVswitchID.NotFound"}) {
				return nil
			}
			if IsExpectedErrors(err, []string{"DependencyViolation.SnatEntry", "DependencyViolation.MulticastDomain", "DependencyViolation", "OperationConflict", "IncorrectRouteEntryStatus", "InternalError", "TaskConflict", "DependencyViolation.EnhancedNatgw", "DependencyViolation.RouteTable", "DependencyViolation.HaVip", "DeleteVSwitch.IncorrectStatus.cbnStatus", "SystemBusy", "IncorrectVSwitchStatus", "LastTokenProcessing", "OperationDenied.OtherSubnetProcessing", "DependencyViolation.SNAT", "DependencyViolation.NetworkAcl"}) || NeedRetry(err) {
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
	stateConf := BuildStateConf([]string{}, []string{}, d.Timeout(schema.TimeoutDelete), 5*time.Second, vpcService.VswitchStateRefreshFunc(d.Id(), []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}
	return nil
}
