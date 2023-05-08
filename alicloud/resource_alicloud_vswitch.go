// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"log"
	"net"
	"strconv"
	"strings"
	"time"

	"github.com/PaesslerAG/jsonpath"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAliCloudVpcVswitch() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudVpcVswitchCreate,
		Read:   resourceAlicloudVpcVswitchRead,
		Update: resourceAlicloudVpcVswitchUpdate,
		Delete: resourceAlicloudVpcVswitchDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"all": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"available_ip_address_count": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"cidr_block": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
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
			},
			"is_default": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"network_acl_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"resource_group_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"route_table_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"tags": tagsSchema(),
			"vswitch_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"vswitch_name": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"name"},
			},
			"vpc_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
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
			},
		},
	}
}

func resourceAlicloudVpcVswitchCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	action := "CreateVSwitch"
	var request map[string]interface{}
	var response map[string]interface{}
	conn, err := client.NewVpcClient()
	if err != nil {
		return WrapError(err)
	}
	request = make(map[string]interface{})
	request["RegionId"] = client.RegionId
	request["ClientToken"] = buildClientToken(action)

	if v, ok := d.GetOk("cidr_block"); ok {
		request["CidrBlock"] = v
	}

	if v, ok := d.GetOk("vpc_id"); ok {
		request["VpcId"] = v
	}

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
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2016-04-28"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
		request["ClientToken"] = buildClientToken(action)

		if err != nil {
			if IsExpectedErrors(err, []string{"TaskConflict", "IncorrectStatus.cbnStatus", "InvalidStatus.RouteEntry", "OperationFailed.IdempotentTokenProcessing", "IncorrectStatus.%s", "CreateVSwitch.IncorrectStatus.cbnStatus", "IncorrectVSwitchStatus", "OperationConflict", "OperationFailed.DistibuteLock", "OperationFailed.NotifyCenCreate"}) || NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(action, response, request)
		return nil
	})

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_vswitch", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(response["VSwitchId"]))

	vpcServiceV2 := VpcServiceV2{client}
	stateConf := BuildStateConf([]string{}, []string{"Available"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, vpcServiceV2.VpcVswitchStateRefreshFunc(d.Id(), "Status", []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return resourceAlicloudVpcVswitchUpdate(d, meta)
}

func resourceAlicloudVpcVswitchRead(d *schema.ResourceData, meta interface{}) error {
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

	d.Set("available_ip_address_count", objectRaw["AvailableIpAddressCount"])
	d.Set("cidr_block", objectRaw["CidrBlock"])
	d.Set("create_time", objectRaw["CreationTime"])
	d.Set("description", objectRaw["Description"])
	d.Set("ipv6_cidr_block", objectRaw["Ipv6CidrBlock"])
	d.Set("is_default", objectRaw["IsDefault"])
	d.Set("network_acl_id", objectRaw["NetworkAclId"])
	d.Set("resource_group_id", objectRaw["ResourceGroupId"])
	d.Set("status", objectRaw["Status"])
	d.Set("vswitch_name", objectRaw["VSwitchName"])
	d.Set("vpc_id", objectRaw["VpcId"])
	d.Set("zone_id", objectRaw["ZoneId"])
	d.Set("vswitch_id", objectRaw["VSwitchId"])
	routeTable2RawObj, _ := jsonpath.Get("$.RouteTable", objectRaw)
	routeTable2Raw := make(map[string]interface{})
	if routeTable2RawObj != nil {
		routeTable2Raw = routeTable2RawObj.(map[string]interface{})
	}
	d.Set("route_table_id", routeTable2Raw["RouteTableId"])
	tag2Raw, _ := jsonpath.Get("$.Tags.Tag", objectRaw)
	tagsMaps := make([]map[string]interface{}, 0)
	if tag2Raw != nil {
		for _, tagChild2Raw := range tag2Raw.([]interface{}) {
			tagsMap := make(map[string]interface{})
			tagChild2Raw := tagChild2Raw.(map[string]interface{})
			tagsMap["tag_key"] = tagChild2Raw["Key"]
			tagsMap["tag_value"] = tagChild2Raw["Value"]
			tagsMaps = append(tagsMaps, tagsMap)
		}
	}
	d.Set("tags", tagsToMap(tagsMaps))

	d.Set("name", d.Get("vswitch_name"))
	d.Set("availability_zone", d.Get("zone_id"))
	if v, ok := objectRaw["Ipv6CidrBlock"]; ok && fmt.Sprint(v) != "" {
		_, cidrBlock := GetIPv6SubnetAddr(v.(string))
		d.Set("ipv6_cidr_block_mask", cidrBlock)
	}
	return nil
}

func resourceAlicloudVpcVswitchUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	update := false
	d.Partial(true)
	update = false
	action := "ModifyVSwitchAttribute"
	conn, err := client.NewVpcClient()
	if err != nil {
		return WrapError(err)
	}
	request = make(map[string]interface{})

	request["VSwitchId"] = d.Id()
	request["RegionId"] = client.RegionId

	if !d.IsNewResource() && d.HasChange("name") {
		update = true
		if v, ok := d.GetOk("name"); ok {
			request["VSwitchName"] = v
		}
	}
	if !d.IsNewResource() && d.HasChange("vswitch_name") {
		update = true
		if v, ok := d.GetOk("vswitch_name"); ok {
			request["VSwitchName"] = v
		}
	}

	if !d.IsNewResource() && d.HasChange("description") {
		update = true
		if v, ok := d.GetOk("description"); ok {
			request["Description"] = v
		}
	}

	if d.HasChange("enable_ipv6") {
		request["Ipv6CidrBlock"] = d.Get("ipv6_cidr_block_mask")
		if v, ok := d.GetOkExists("enable_ipv6"); ok {
			request["EnableIPv6"] = v
		}
	}

	if !d.IsNewResource() && d.HasChange("ipv6_cidr_block_mask") {
		update = true
		if v, ok := d.GetOk("ipv6_cidr_block_mask"); ok {
			request["Ipv6CidrBlock"] = v
		}
	}

	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2016-04-28"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})

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
		{
			vpcServiceV2 := VpcServiceV2{client}
			stateConf := BuildStateConf([]string{}, []string{"Available"}, d.Timeout(schema.TimeoutDelete), 5*time.Second, vpcServiceV2.VpcVswitchStateRefreshFunc(d.Id(), "Status", []string{}))
			if _, err := stateConf.WaitForState(); err != nil {
				return WrapErrorf(err, IdMsg, d.Id())
			}
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
	return resourceAlicloudVpcVswitchRead(d, meta)
}

func resourceAlicloudVpcVswitchDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "DeleteVSwitch"
	var request map[string]interface{}
	var response map[string]interface{}
	conn, err := client.NewVpcClient()
	if err != nil {
		return WrapError(err)
	}
	request = make(map[string]interface{})

	request["VSwitchId"] = d.Id()
	request["RegionId"] = client.RegionId

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2016-04-28"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})

		if err != nil {
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
		if IsExpectedErrors(err, []string{"InvalidVswitchID.NotFound", "InvalidVSwitchId.NotFound"}) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	vpcServiceV2 := VpcServiceV2{client}
	stateConf := BuildStateConf([]string{}, []string{""}, d.Timeout(schema.TimeoutDelete), 5*time.Second, vpcServiceV2.VpcVswitchStateRefreshFunc(d.Id(), "Status", []string{}))
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
