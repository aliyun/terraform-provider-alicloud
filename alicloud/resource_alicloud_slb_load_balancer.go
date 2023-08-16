// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"log"
	"time"

	"github.com/PaesslerAG/jsonpath"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAliCloudSlbLoadBalancer() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudSlbLoadBalancerCreate,
		Read:   resourceAliCloudSlbLoadBalancerRead,
		Update: resourceAliCloudSlbLoadBalancerUpdate,
		Delete: resourceAliCloudSlbLoadBalancerDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"address": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"address_ip_version": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"address_type": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"auto_pay": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"bandwidth": {
				Type:     schema.TypeInt,
				Computed: true,
				Optional: true,
			},
			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"delete_protection": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: StringInSlice([]string{"on", "off"}, false),
			},
			"duration": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"instance_charge_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: StringInSlice([]string{"PayBySpec", "PayByCLCU"}, false),
			},
			"internet_charge_type": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: StringInSlice([]string{"paybybandwidth", "paybytraffic"}, false),
			},
			"load_balancer_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"load_balancer_spec": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"master_zone_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"modification_protection_reason": {
				Type:     schema.TypeString,
				Optional: true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					if v, ok := d.GetOk("modification_protection_status"); ok && v.(string) == "NonProtection" {
						return true
					}
					return false
				},
			},
			"modification_protection_status": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
			"payment_type": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: StringInSlice([]string{"Subscription", "PayAsYouGo"}, false),
			},
			"pricing_cycle": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"resource_group_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"slave_zone_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"status": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: StringInSlice([]string{"inactive", "active", "locked"}, false),
			},
			"tags": tagsSchema(),
			"vswitch_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"vpc_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
		},
	}
}

func resourceAliCloudSlbLoadBalancerCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	action := "CreateLoadBalancer"
	var request map[string]interface{}
	var response map[string]interface{}
	conn, err := client.NewSlbClient()
	if err != nil {
		return WrapError(err)
	}
	request = make(map[string]interface{})
	request["RegionId"] = client.RegionId
	request["ClientToken"] = buildClientToken(action)

	if v, ok := d.GetOk("address_type"); ok {
		request["AddressType"] = v
	}
	if v, ok := d.GetOk("internet_charge_type"); ok {
		request["InternetChargeType"] = convertSlbInternetChargeTypeRequest(v.(string))
	}
	if v, ok := d.GetOk("bandwidth"); ok {
		request["Bandwidth"] = v
	}
	if v, ok := d.GetOk("load_balancer_name"); ok {
		request["LoadBalancerName"] = v
	}
	if v, ok := d.GetOk("vpc_id"); ok {
		request["VpcId"] = v
	}
	if v, ok := d.GetOk("vswitch_id"); ok {
		request["VSwitchId"] = v
	}
	if v, ok := d.GetOk("master_zone_id"); ok {
		request["MasterZoneId"] = v
	}
	if v, ok := d.GetOk("slave_zone_id"); ok {
		request["SlaveZoneId"] = v
	}
	if v, ok := d.GetOk("resource_group_id"); ok {
		request["ResourceGroupId"] = v
	}
	if v, ok := d.GetOk("address"); ok {
		request["Address"] = v
	}
	if v, ok := d.GetOk("delete_protection"); ok {
		request["DeleteProtection"] = v
	}
	if v, ok := d.GetOk("payment_type"); ok {
		request["PayType"] = convertSlbPayTypeRequest(v.(string))
	}
	if v, ok := d.GetOk("modification_protection_status"); ok {
		request["ModificationProtectionStatus"] = v
	}
	if v, ok := d.GetOk("modification_protection_reason"); ok {
		request["ModificationProtectionReason"] = v
	}
	if v, ok := d.GetOk("address_ip_version"); ok {
		request["AddressIPVersion"] = v
	}
	if v, ok := d.GetOk("load_balancer_spec"); ok {
		request["LoadBalancerSpec"] = v
	}
	if v, ok := d.GetOk("pricing_cycle"); ok {
		request["PricingCycle"] = v
	}
	if v, ok := d.GetOk("duration"); ok {
		request["Duration"] = v
	}
	if v, ok := d.GetOkExists("auto_pay"); ok {
		request["AutoPay"] = v
	}
	if v, ok := d.GetOk("instance_charge_type"); ok {
		request["InstanceChargeType"] = v
	}
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2014-05-15"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
		request["ClientToken"] = buildClientToken(action)

		if err != nil {
			if IsExpectedErrors(err, []string{"OperationFailed.TokenIsProcessing"}) || NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(action, response, request)
		return nil
	})

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_slb_load_balancer", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(response["LoadBalancerId"]))

	slbServiceV2 := SlbServiceV2{client}
	stateConf := BuildStateConf([]string{}, []string{"active"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, slbServiceV2.SlbLoadBalancerStateRefreshFunc(d.Id(), "LoadBalancerStatus", []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return resourceAliCloudSlbLoadBalancerUpdate(d, meta)
}

func resourceAliCloudSlbLoadBalancerRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	slbServiceV2 := SlbServiceV2{client}

	objectRaw, err := slbServiceV2.DescribeSlbLoadBalancer(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_slb_load_balancer DescribeSlbLoadBalancer Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("address", objectRaw["Address"])
	d.Set("address_ip_version", objectRaw["AddressIPVersion"])
	d.Set("address_type", objectRaw["AddressType"])
	d.Set("bandwidth", objectRaw["Bandwidth"])
	d.Set("create_time", objectRaw["CreateTime"])
	d.Set("delete_protection", objectRaw["DeleteProtection"])
	d.Set("instance_charge_type", objectRaw["InstanceChargeType"])
	d.Set("internet_charge_type", convertSlbInternetChargeTypeResponse(objectRaw["InternetChargeType"]))
	d.Set("load_balancer_name", objectRaw["LoadBalancerName"])
	d.Set("load_balancer_spec", objectRaw["LoadBalancerSpec"])
	d.Set("master_zone_id", objectRaw["MasterZoneId"])
	d.Set("modification_protection_reason", objectRaw["ModificationProtectionReason"])
	d.Set("modification_protection_status", objectRaw["ModificationProtectionStatus"])
	d.Set("payment_type", convertSlbPayTypeResponse(objectRaw["PayType"]))
	d.Set("resource_group_id", objectRaw["ResourceGroupId"])
	d.Set("slave_zone_id", objectRaw["SlaveZoneId"])
	d.Set("status", objectRaw["LoadBalancerStatus"])
	d.Set("vswitch_id", objectRaw["VSwitchId"])
	d.Set("vpc_id", objectRaw["VpcId"])

	objectRaw, err = slbServiceV2.DescribeListTagResources(d.Id())
	if err != nil {
		return WrapError(err)
	}

	tagsMaps, _ := jsonpath.Get("$.TagResources.TagResource", objectRaw)
	d.Set("tags", tagsToMap(tagsMaps))

	return nil
}

func resourceAliCloudSlbLoadBalancerUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	update := false
	d.Partial(true)
	action := "ModifyLoadBalancerInstanceSpec"
	conn, err := client.NewSlbClient()
	if err != nil {
		return WrapError(err)
	}
	request = make(map[string]interface{})
	request["LoadBalancerId"] = d.Id()
	request["RegionId"] = client.RegionId
	if !d.IsNewResource() && d.HasChange("load_balancer_spec") {
		update = true
		request["LoadBalancerSpec"] = d.Get("load_balancer_spec")
	}

	if v, ok := d.GetOkExists("auto_pay"); ok {
		request["AutoPay"] = v
	}
	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2014-05-15"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})

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
		slbServiceV2 := SlbServiceV2{client}
		stateConf := BuildStateConf([]string{}, []string{"active"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, slbServiceV2.SlbLoadBalancerStateRefreshFunc(d.Id(), "LoadBalancerStatus", []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
		d.SetPartial("load_balancer_spec")
	}
	update = false
	action = "SetLoadBalancerName"
	conn, err = client.NewSlbClient()
	if err != nil {
		return WrapError(err)
	}
	request = make(map[string]interface{})
	request["LoadBalancerId"] = d.Id()
	request["RegionId"] = client.RegionId
	if !d.IsNewResource() && d.HasChange("load_balancer_name") {
		update = true
		request["LoadBalancerName"] = d.Get("load_balancer_name")
	}

	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2014-05-15"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})

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
		d.SetPartial("load_balancer_name")
	}
	update = false
	action = "SetLoadBalancerDeleteProtection"
	conn, err = client.NewSlbClient()
	if err != nil {
		return WrapError(err)
	}
	request = make(map[string]interface{})
	request["LoadBalancerId"] = d.Id()
	request["RegionId"] = client.RegionId
	if !d.IsNewResource() && d.HasChange("delete_protection") {
		update = true
		request["DeleteProtection"] = d.Get("delete_protection")
	}

	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2014-05-15"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})

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
		d.SetPartial("delete_protection")
	}
	update = false
	action = "ModifyLoadBalancerInternetSpec"
	conn, err = client.NewSlbClient()
	if err != nil {
		return WrapError(err)
	}
	request = make(map[string]interface{})
	request["LoadBalancerId"] = d.Id()
	request["RegionId"] = client.RegionId
	if !d.IsNewResource() && d.HasChange("internet_charge_type") {
		update = true
		request["InternetChargeType"] = d.Get("internet_charge_type")
	}

	if !d.IsNewResource() && d.HasChange("bandwidth") {
		update = true
		request["Bandwidth"] = d.Get("bandwidth")
	}

	if v, ok := d.GetOkExists("auto_pay"); ok {
		request["AutoPay"] = v
	}
	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2014-05-15"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})

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
		slbServiceV2 := SlbServiceV2{client}
		stateConf := BuildStateConf([]string{}, []string{"active"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, slbServiceV2.SlbLoadBalancerStateRefreshFunc(d.Id(), "LoadBalancerStatus", []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
		d.SetPartial("internet_charge_type")
		d.SetPartial("bandwidth")
	}
	update = false
	action = "SetLoadBalancerModificationProtection"
	conn, err = client.NewSlbClient()
	if err != nil {
		return WrapError(err)
	}
	request = make(map[string]interface{})
	request["LoadBalancerId"] = d.Id()
	request["RegionId"] = client.RegionId
	if !d.IsNewResource() && d.HasChange("modification_protection_status") {
		update = true
	}
	request["ModificationProtectionStatus"] = d.Get("modification_protection_status")

	if !d.IsNewResource() && d.HasChange("modification_protection_reason") {
		update = true
		request["ModificationProtectionReason"] = d.Get("modification_protection_reason")
	}

	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2014-05-15"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})

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
		d.SetPartial("modification_protection_status")
		d.SetPartial("modification_protection_reason")
	}
	update = false
	action = "ModifyLoadBalancerPayType"
	conn, err = client.NewSlbClient()
	if err != nil {
		return WrapError(err)
	}
	request = make(map[string]interface{})
	request["LoadBalancerId"] = d.Id()
	request["RegionId"] = client.RegionId
	if !d.IsNewResource() && d.HasChange("payment_type") {
		update = true
		request["PayType"] = convertSlbPayTypeRequest(d.Get("payment_type").(string))
	}

	if v, ok := d.GetOk("pricing_cycle"); ok {
		request["PricingCycle"] = v
	}
	if v, ok := d.GetOk("duration"); ok {
		request["Duration"] = v
	}
	if v, ok := d.GetOkExists("auto_pay"); ok {
		request["AutoPay"] = v
	}
	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2014-05-15"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})

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
		d.SetPartial("payment_type")
	}
	update = false
	action = "SetLoadBalancerStatus"
	conn, err = client.NewSlbClient()
	if err != nil {
		return WrapError(err)
	}
	request = make(map[string]interface{})
	request["LoadBalancerId"] = d.Id()
	request["RegionId"] = client.RegionId
	if d.HasChange("status") {
		update = true
		request["LoadBalancerStatus"] = d.Get("status")
	}

	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2014-05-15"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})

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
		slbServiceV2 := SlbServiceV2{client}
		stateConf := BuildStateConf([]string{}, []string{"active", "inactive"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, slbServiceV2.SlbLoadBalancerStateRefreshFunc(d.Id(), "LoadBalancerStatus", []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
		d.SetPartial("status")
	}
	update = false
	action = "ModifyLoadBalancerInstanceChargeType"
	conn, err = client.NewSlbClient()
	if err != nil {
		return WrapError(err)
	}
	request = make(map[string]interface{})
	request["LoadBalancerId"] = d.Id()
	request["RegionId"] = client.RegionId
	if !d.IsNewResource() && d.HasChange("load_balancer_name") {
		update = true
		request["OwnerId"] = d.Get("load_balancer_name")
	}

	if !d.IsNewResource() && d.HasChange("internet_charge_type") {
		update = true
		request["InternetChargeType"] = d.Get("internet_charge_type")
	}

	if !d.IsNewResource() && d.HasChange("instance_charge_type") {
		update = true
	}
	request["InstanceChargeType"] = d.Get("instance_charge_type")

	if !d.IsNewResource() && d.HasChange("load_balancer_spec") {
		update = true
		request["LoadBalancerSpec"] = d.Get("load_balancer_spec")
	}

	if !d.IsNewResource() && d.HasChange("bandwidth") {
		update = true
		request["Bandwidth"] = d.Get("bandwidth")
	}

	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2014-05-15"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})

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
		d.SetPartial("load_balancer_name")
		d.SetPartial("internet_charge_type")
		d.SetPartial("instance_charge_type")
		d.SetPartial("load_balancer_spec")
		d.SetPartial("bandwidth")
	}
	update = false
	action = "MoveResourceGroup"
	conn, err = client.NewSlbClient()
	if err != nil {
		return WrapError(err)
	}
	request = make(map[string]interface{})
	request["ResourceId"] = d.Id()
	request["RegionId"] = client.RegionId
	if !d.IsNewResource() && d.HasChange("resource_group_id") {
		update = true
		request["NewResourceGroupId"] = d.Get("resource_group_id")
	}

	request["ResourceType"] = "instance"
	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2014-05-15"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})

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
		d.SetPartial("resource_group_id")
	}

	if d.HasChange("tags") {
		slbServiceV2 := SlbServiceV2{client}
		if err := slbServiceV2.SetResourceTags(d, "instance"); err != nil {
			return WrapError(err)
		}
		d.SetPartial("tags")
	}
	d.Partial(false)
	return resourceAliCloudSlbLoadBalancerRead(d, meta)
}

func resourceAliCloudSlbLoadBalancerDelete(d *schema.ResourceData, meta interface{}) error {

	if v, ok := d.GetOk("payment_type"); ok {
		if v == "Subscription" {
			log.Printf("[WARN] Cannot destroy resource alicloud_slb_load_balancer which payment_type valued Subscription. Terraform will remove this resource from the state file, however resources may remain.")
			return nil
		}
	}
	client := meta.(*connectivity.AliyunClient)
	action := "DeleteLoadBalancer"
	var request map[string]interface{}
	var response map[string]interface{}
	conn, err := client.NewSlbClient()
	if err != nil {
		return WrapError(err)
	}
	request = make(map[string]interface{})
	request["LoadBalancerId"] = d.Id()
	request["RegionId"] = client.RegionId

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2014-05-15"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})

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
		if IsExpectedErrors(err, []string{"InvalidLoadBalancerId.NotFound"}) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	slbServiceV2 := SlbServiceV2{client}
	stateConf := BuildStateConf([]string{}, []string{}, d.Timeout(schema.TimeoutDelete), 5*time.Second, slbServiceV2.SlbLoadBalancerStateRefreshFunc(d.Id(), "LoadBalancerStatus", []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}
	return nil
}

func convertSlbInternetChargeTypeResponse(source interface{}) interface{} {
	switch source {
	}
	return source
}
func convertSlbPayTypeResponse(source interface{}) interface{} {
	switch source {
	case "PayOnDemand":
		return "PayAsYouGo"
	case "PrePay":
		return "Subscription"
	}
	return source
}
func convertSlbPayTypeRequest(source interface{}) interface{} {
	switch source {
	case "PayAsYouGo":
		return "PayOnDemand"
	case "Subscription":
		return "PrePay"
	}
	return source
}
func convertSlbInternetChargeTypeRequest(source interface{}) interface{} {
	switch source {
	case "PayByBandwidth":
		return "paybybandwidth"
	case "PayByTraffic":
		return "paybytraffic"
	}
	return source
}

func convertSlbLoadBalancerInternetChargeTypeRequest(source interface{}) interface{} {
	switch source {
	case "PayByBandwidth":
		return "paybybandwidth"
	case "PayByTraffic":
		return "paybytraffic"
	}
	return source
}
func convertSlbLoadBalancerPaymentTypeRequest(source interface{}) interface{} {
	switch source {
	case "PayAsYouGo":
		return "PayOnDemand"
	case "Subscription":
		return "PrePay"
	}
	return source
}
func convertSlbLoadBalancerInstanceChargeTypeRequest(source interface{}) interface{} {
	switch source {
	case "PostPaid":
		return "PayOnDemand"
	case "PrePaid":
		return "PrePay"
	}
	return source
}
func convertSlbLoadBalancerInternetChargeTypeResponse(source interface{}) interface{} {
	switch source {
	case "paybybandwidth":
		return "PayByBandwidth"
	case "paybytraffic":
		return "PayByTraffic"
	}
	return source
}
func convertSlbLoadBalancerPaymentTypeResponse(source interface{}) interface{} {
	switch source {
	case "PayOnDemand":
		return "PayAsYouGo"
	case "PrePay":
		return "Subscription"
	}
	return source
}
func convertSlbLoadBalancerInstanceChargeTypeResponse(source interface{}) interface{} {
	switch source {
	case "PayOnDemand":
		return "PostPaid"
	case "PrePay":
		return "PrePaid"
	}
	return source
}
