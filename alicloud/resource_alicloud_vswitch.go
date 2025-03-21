package alicloud

import (
	"fmt"
	"log"
	"net"
	"strconv"
	"strings"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAliCloudVpcVswitch() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudVpcVswitchCreate,
		Read:   resourceAliCloudVpcVswitchRead,
		Update: resourceAliCloudVpcVswitchUpdate,
		Delete: resourceAliCloudVpcVswitchDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"cidr_block": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"enable_ipv6": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"ipv6_cidr_block": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"ipv6_cidr_block_mask": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"is_default": {
				Type:     schema.TypeBool,
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
			"vpc_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"zone_id": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ExactlyOneOf: []string{"zone_id", "availability_zone"},
				ForceNew:     true,
			},
			"name": {
				Type:       schema.TypeString,
				Optional:   true,
				Computed:   true,
				Deprecated: "Field 'name' has been deprecated from provider version 1.119.0. New field 'vswitch_name' instead.",
			},
			"availability_zone": {
				Type:       schema.TypeString,
				Optional:   true,
				Computed:   true,
				Deprecated: "Field 'availability_zone' has been deprecated from provider version 1.119.0. New field 'zone_id' instead.",
				ForceNew:   true,
			},
		},
	}
}

func resourceAliCloudVpcVswitchCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	isDefault := false
	if v, ok := d.GetOkExists("is_default"); ok {
		isDefault = v.(bool)
	}

	if isDefault {

		action := "CreateDefaultVSwitch"
		var request map[string]interface{}
		var response map[string]interface{}
		var err error
		request = make(map[string]interface{})
		request["RegionId"] = client.RegionId
		request["ClientToken"] = buildClientToken(action)

		if v, ok := d.GetOk("availability_zone"); ok {
			request["ZoneId"] = v
		}

		if v, ok := d.GetOk("zone_id"); ok {
			request["ZoneId"] = v
		}

		if v, ok := d.GetOk("ipv6_cidr_block_mask"); ok {
			request["Ipv6CidrBlock"] = v
		}

		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
			response, err = client.RpcPost("Vpc", "2016-04-28", action, nil, request, true)
			request["ClientToken"] = buildClientToken(action)

			if err != nil {
				if IsExpectedErrors(err, []string{"TaskConflict", "IncorrectStatus.cbnStatus", "InvalidStatus.RouteEntry", "OperationFailed.IdempotentTokenProcessing", "IncorrectStatus", "CreateVSwitch.IncorrectStatus.cbnStatus", "IncorrectVSwitchStatus", "OperationConflict", "OperationFailed.DistibuteLock", "OperationFailed.NotifyCenCreate"}) || NeedRetry(err) {
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

	} else {
		action := "CreateVSwitch"
		var request map[string]interface{}
		var response map[string]interface{}
		var err error
		request = make(map[string]interface{})
		request["RegionId"] = client.RegionId
		request["ClientToken"] = buildClientToken(action)

		request["CidrBlock"] = d.Get("cidr_block")
		request["VpcId"] = d.Get("vpc_id")
		if v, ok := d.GetOk("name"); ok {
			request["VSwitchName"] = v
		}
		if v, ok := d.GetOk("vswitch_name"); ok {
			request["VSwitchName"] = v
		}
		if v, ok := d.GetOk("description"); ok {
			request["Description"] = v
		}
		if v, ok := d.GetOk("availability_zone"); ok {
			request["ZoneId"] = v
		}
		if v, ok := d.GetOk("zone_id"); ok {
			request["ZoneId"] = v
		}
		if v, ok := d.GetOk("ipv6_cidr_block_mask"); ok {
			request["Ipv6CidrBlock"] = v
		}
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
			response, err = client.RpcPost("Vpc", "2016-04-28", action, nil, request, true)
			request["ClientToken"] = buildClientToken(action)

			if err != nil {
				if IsExpectedErrors(err, []string{"TaskConflict", "IncorrectStatus.cbnStatus", "InvalidStatus.RouteEntry", "OperationFailed.IdempotentTokenProcessing", "IncorrectStatus", "CreateVSwitch.IncorrectStatus.cbnStatus", "IncorrectVSwitchStatus", "OperationConflict", "OperationFailed.DistibuteLock", "OperationFailed.NotifyCenCreate"}) || NeedRetry(err) {
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

	}

	vpcServiceV2 := VpcServiceV2{client}
	stateConf := BuildStateConf([]string{}, []string{"Available"}, d.Timeout(schema.TimeoutCreate), 0, vpcServiceV2.VpcVswitchStateRefreshFunc(d.Id(), "Status", []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return resourceAliCloudVpcVswitchUpdate(d, meta)
}

func resourceAliCloudVpcVswitchRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	vpcServiceV2 := VpcServiceV2{client}

	objectRaw, err := vpcServiceV2.DescribeVpcVswitch(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_vswitch DescribeVpcVswitch Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("cidr_block", objectRaw["CidrBlock"])
	d.Set("create_time", objectRaw["CreationTime"])
	d.Set("description", objectRaw["Description"])
	d.Set("ipv6_cidr_block", objectRaw["Ipv6CidrBlock"])
	d.Set("status", objectRaw["Status"])
	d.Set("vswitch_name", objectRaw["VSwitchName"])
	d.Set("vpc_id", objectRaw["VpcId"])
	d.Set("zone_id", objectRaw["ZoneId"])

	tagsMaps, _ := jsonpath.Get("$.Tags.Tag", objectRaw)
	d.Set("tags", tagsToMap(tagsMaps))

	d.Set("name", d.Get("vswitch_name"))
	d.Set("availability_zone", d.Get("zone_id"))
	if v, ok := objectRaw["Ipv6CidrBlock"]; ok && fmt.Sprint(v) != "" {
		_, cidrBlock := GetIPv6SubnetAddr(v.(string))
		d.Set("ipv6_cidr_block_mask", cidrBlock)
	}

	if enableIpv6, ok := d.GetOkExists("enable_ipv6"); ok {
		d.Set("enable_ipv6", enableIpv6)
	}
	return nil
}

func resourceAliCloudVpcVswitchUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	update := false
	d.Partial(true)
	action := "ModifyVSwitchAttribute"
	var err error
	request = make(map[string]interface{})

	request["VSwitchId"] = d.Id()
	request["RegionId"] = client.RegionId
	if !d.IsNewResource() && d.HasChange("name") {
		update = true
		request["VSwitchName"] = d.Get("name")
	}
	if !d.IsNewResource() && d.HasChange("vswitch_name") {
		update = true
		request["VSwitchName"] = d.Get("vswitch_name")
	}

	if !d.IsNewResource() && d.HasChange("description") {
		update = true
		request["Description"] = d.Get("description")
	}

	if !d.IsNewResource() && d.HasChange("ipv6_cidr_block_mask") {
		err := CancelIpv6(d, meta)
		if err != nil {
			return WrapError(err)
		}
		if v, ok := d.GetOk("ipv6_cidr_block_mask"); ok {
			update = true
			request["EnableIPv6"] = true
			request["Ipv6CidrBlock"] = v
		}
	}

	if !d.IsNewResource() && d.HasChange("enable_ipv6") {
		update = true
		request["EnableIPv6"] = d.Get("enable_ipv6")
	}

	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("Vpc", "2016-04-28", action, nil, request, false)

			if err != nil {
				if IsExpectedErrors(err, []string{"OperationConflict", "OperationFailed.LastTokenProcessing", "IncorrectStatus.VSwitch", "IncorrectStatus.VpcRouteEntry", "ServiceUnavailable"}) || NeedRetry(err) {
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
		stateConf := BuildStateConf([]string{}, []string{"Available"}, d.Timeout(schema.TimeoutUpdate), 0, vpcServiceV2.VpcVswitchStateRefreshFunc(d.Id(), "Status", []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
		d.SetPartial("vswitch_name")
		d.SetPartial("description")
	}

	update = false
	if d.HasChange("tags") {
		update = true
		vpcServiceV2 := VpcServiceV2{client}
		if err := vpcServiceV2.SetResourceTags(d, "VSWITCH"); err != nil {
			return WrapError(err)
		}
		d.SetPartial("tags")
	}
	d.Partial(false)
	return resourceAliCloudVpcVswitchRead(d, meta)
}

func resourceAliCloudVpcVswitchDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "DeleteVSwitch"
	var request map[string]interface{}
	var response map[string]interface{}
	var err error
	request = make(map[string]interface{})

	request["VSwitchId"] = d.Id()
	request["RegionId"] = client.RegionId

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = client.RpcPost("Vpc", "2016-04-28", action, nil, request, false)

		if err != nil {
			if IsExpectedErrors(err, []string{"DependencyViolation", "DependencyViolation.SnatEntry", "DependencyViolation.MulticastDomain", "DependencyViolation", "OperationConflict", "IncorrectRouteEntryStatus", "InternalError", "TaskConflict", "DependencyViolation.EnhancedNatgw", "DependencyViolation.RouteTable", "DependencyViolation.HaVip", "DeleteVSwitch.IncorrectStatus.cbnStatus", "SystemBusy", "IncorrectVSwitchStatus", "LastTokenProcessing", "OperationDenied.OtherSubnetProcessing", "DependencyViolation.SNAT", "DependencyViolation.NetworkAcl"}) || NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(action, response, request)
		return nil
	})

	if err != nil {
		if IsExpectedErrors(err, []string{"InvalidVswitchID.NotFound", "InvalidVSwitchId.NotFound"}) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	vpcServiceV2 := VpcServiceV2{client}
	stateConf := BuildStateConf([]string{}, []string{""}, d.Timeout(schema.TimeoutDelete), 0, vpcServiceV2.VpcVswitchStateRefreshFunc(d.Id(), "Status", []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}
	return nil
}

func GetIPv6SubnetAddr(ipAddr string) (string, int) {
	// Split the IP address and subnet prefix length
	ip, prefix, err := net.ParseCIDR(ipAddr)
	if err != nil {
		return "", 0
	}
	mask, _ := strconv.Atoi(strings.Split(ipAddr, "/")[1])
	// Get the network address by masking the IP address with the subnet prefix
	netAddr := ip.Mask(prefix.Mask)
	// Convert the network address to a string
	netAddrStr := netAddr.String()
	// Get the last8 bits of the address
	last8Bits := netAddr[mask/8-1]
	// Convert the last8 bits to an integer
	return netAddrStr, int(last8Bits)
}

func CancelIpv6(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	vpcServiceV2 := VpcServiceV2{client}
	object, _ := vpcServiceV2.DescribeVpcVswitch(d.Id())

	if _, ok := d.GetOk("ipv6_cidr_block_mask"); !ok {
		return nil
	}

	if v, ok := object["Ipv6CidrBlock"]; ok && fmt.Sprint(v) != "" {
		var response map[string]interface{}
		var err error
		action := "ModifyVSwitchAttribute"
		request := map[string]interface{}{
			"RegionId":   client.RegionId,
			"VSwitchId":  d.Id(),
			"EnableIPv6": false,
		}
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			response, err = client.RpcPost("Vpc", "2016-04-28", action, nil, request, false)
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

		vpcServiceV2 := VpcServiceV2{client}
		stateConf := BuildStateConf([]string{}, []string{"Available"}, d.Timeout(schema.TimeoutUpdate), 0, vpcServiceV2.VpcVswitchStateRefreshFunc(d.Id(), "Status", []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
	}

	return nil
}
