// Package cbwp. This file is generated automatically. Please do not modify it manually, thank you!
package cbwp

import (
	"fmt"
	"log"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/tags"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/validators"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceCommonBandwidthPackage() *schema.Resource {
	return &schema.Resource{
		Create: resourceCommonBandwidthPackageCreate,
		Read:   resourceCommonBandwidthPackageRead,
		Update: resourceCommonBandwidthPackageUpdate,
		Delete: resourceCommonBandwidthPackageDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"bandwidth": {
				Type:     schema.TypeString,
				Required: true,
			},
			"bandwidth_package_name": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"name"},
				Computed:      true,
			},
			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"deletion_protection": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"description": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validators.StringLenBetween(2, 256),
			},
			"force": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"internet_charge_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Default:      "PayByTraffic",
				ValidateFunc: validators.StringInSlice([]string{"PayBy95", "PayByBandwidth", "PayByTraffic", "PayByDominantTraffic", "PayByOld95"}, false),
			},
			"isp": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				ValidateFunc: validators.StringInSlice([]string{"BGP", "BGP_PRO", "ChinaTelecom", "ChinaUnicom", "ChinaMobile", "ChinaTelecom_L2", "ChinaUnicom_L2", "ChinaMobile_L2", "BGP_FinanceCloud", "BGP_International"}, false),
			},
			"payment_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"ratio": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				ValidateFunc: validators.IntBetween(0, 100),
			},
			"region_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"resource_group_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"security_protection_types": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"tags": tags.TagsSchema(),
			"zone": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"name": {
				Type:       schema.TypeString,
				Optional:   true,
				Computed:   true,
				Deprecated: "Field 'name' has been deprecated since provider version 1.120.0. New field 'bandwidth_package_name' instead.",
			},
		},
	}
}

func resourceCommonBandwidthPackageCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "CreateCommonBandwidthPackage"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	request["RegionId"] = client.RegionId
	request["ClientToken"] = helper.BuildClientToken(action)

	if v, ok := d.GetOk("description"); ok && len(v.(string)) > 0 {
		request["Description"] = v
	}
	if v, ok := d.GetOk("resource_group_id"); ok {
		request["ResourceGroupId"] = v
	}
	request["Bandwidth"] = d.Get("bandwidth")
	if v, ok := d.GetOkExists("ratio"); ok {
		request["Ratio"] = v
	}
	if v, ok := d.GetOk("security_protection_types"); ok {
		securityProtectionTypesMapsArray := v.([]interface{})
		request["SecurityProtectionTypes"] = securityProtectionTypesMapsArray
	}

	if v, ok := d.GetOk("isp"); ok {
		request["ISP"] = v
	}
	if v, ok := d.GetOk("internet_charge_type"); ok {
		request["InternetChargeType"] = v
	}
	if v, ok := d.GetOk("zone"); ok {
		request["Zone"] = v
	}
	if v, ok := d.GetOk("name"); ok || d.HasChange("name") {
		request["Name"] = v
	}

	if v, ok := d.GetOk("bandwidth_package_name"); ok {
		request["Name"] = v
	}
	if v, ok := d.GetOk("tags"); ok {
		tagsMap := tags.ConvertTags(v.(map[string]interface{}))
		request = tags.ExpandTagsToMap(request, tagsMap)
	}

	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
		response, err = client.RpcPost("Vpc", "2016-04-28", action, query, request, true)
		if err != nil {
			if errs.IsExpectedErrors(err, []string{"BandwidthPackageOperation.conflict", "OperationConflict", "LastTokenProcessing", "IncorrectStatus", "SystemBusy", "ServiceUnavailable"}) || helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	logs.AddDebug(action, response, request)

	if err != nil {
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, "alicloud_common_bandwidth_package", action, errs.AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(response["BandwidthPackageId"]))

	cbwpService := NewService(client)
	stateConf := helper.BuildStateConf([]string{}, []string{"Available"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, cbwpService.CommonBandwidthPackageStateRefreshFunc(d.Id(), "Status", []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return errs.WrapErrorf(err, errs.IdMsg, d.Id())
	}

	return resourceCommonBandwidthPackageUpdate(d, meta)
}

func resourceCommonBandwidthPackageRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	cbwpService := NewService(client)

	objectRaw, err := cbwpService.DescribeCommonBandwidthPackage(d.Id())
	if err != nil {
		if !d.IsNewResource() && errs.NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_common_bandwidth_package DescribeCommonBandwidthPackage Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return errs.WrapError(err)
	}

	d.Set("bandwidth", objectRaw["Bandwidth"])
	d.Set("bandwidth_package_name", objectRaw["Name"])
	d.Set("create_time", objectRaw["CreationTime"])
	d.Set("deletion_protection", objectRaw["DeletionProtection"])
	d.Set("description", objectRaw["Description"])
	d.Set("internet_charge_type", objectRaw["InternetChargeType"])
	d.Set("isp", objectRaw["ISP"])
	d.Set("ratio", objectRaw["Ratio"])
	d.Set("region_id", objectRaw["RegionId"])
	d.Set("resource_group_id", objectRaw["ResourceGroupId"])
	d.Set("status", objectRaw["Status"])

	securityProtectionTypeRaw, _ := jsonpath.Get("$.SecurityProtectionTypes.SecurityProtectionType", objectRaw)
	d.Set("security_protection_types", securityProtectionTypeRaw)
	d.Set("payment_type", convertCommonBandwidthPackageInstanceChargeTypeResponse(objectRaw["InstanceChargeType"]))

	objectRaw, err = cbwpService.DescribeCommonBandwidthPackageListTagResources(d.Id())
	if err != nil && !errs.NotFoundError(err) {
		return errs.WrapError(err)
	}

	tagsMaps, _ := jsonpath.Get("$.TagResources.TagResource", objectRaw)
	d.Set("tags", tags.TagsToMap(tagsMaps))

	d.Set("name", d.Get("bandwidth_package_name"))
	return nil
}

func resourceCommonBandwidthPackageUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	update := false
	d.Partial(true)

	var err error
	action := "ModifyCommonBandwidthPackageAttribute"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["BandwidthPackageId"] = d.Id()
	request["RegionId"] = client.RegionId
	if !d.IsNewResource() && d.HasChange("description") {
		update = true
		request["Description"] = d.Get("description")
	}

	if !d.IsNewResource() && d.HasChange("name") {
		update = true
		request["Name"] = d.Get("name")
	}

	if !d.IsNewResource() && d.HasChange("bandwidth_package_name") {
		update = true
		request["Name"] = d.Get("bandwidth_package_name")
	}

	if update {
		wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
		err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
			response, err = client.RpcPost("Vpc", "2016-04-28", action, query, request, true)
			if err != nil {
				if helper.NeedRetry(err) {
					wait()
					return retry.RetryableError(err)
				}
				return retry.NonRetryableError(err)
			}
			return nil
		})
		logs.AddDebug(action, response, request)
		if err != nil {
			return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
		}
	}
	update = false
	action = "ModifyCommonBandwidthPackageSpec"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["BandwidthPackageId"] = d.Id()
	request["RegionId"] = client.RegionId
	if !d.IsNewResource() && d.HasChange("bandwidth") {
		update = true
	}
	request["Bandwidth"] = d.Get("bandwidth")
	if update {
		wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
		err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
			response, err = client.RpcPost("Vpc", "2016-04-28", action, query, request, true)
			if err != nil {
				if errs.IsExpectedErrors(err, []string{"BandwidthPackageOperation.conflict", "OperationConflict", "LastTokenProcessing", "IncorrectStatus", "SystemBusy", "ServiceUnavailable"}) || helper.NeedRetry(err) {
					wait()
					return retry.RetryableError(err)
				}
				return retry.NonRetryableError(err)
			}
			return nil
		})
		logs.AddDebug(action, response, request)
		if err != nil {
			return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
		}
		cbwpService := NewService(client)
		stateConf := helper.BuildStateConf([]string{}, []string{"Available"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, cbwpService.CommonBandwidthPackageStateRefreshFunc(d.Id(), "Status", []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return errs.WrapErrorf(err, errs.IdMsg, d.Id())
		}
	}
	update = false
	action = "MoveResourceGroup"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["ResourceId"] = d.Id()
	request["RegionId"] = client.RegionId
	if _, ok := d.GetOk("resource_group_id"); ok && !d.IsNewResource() && d.HasChange("resource_group_id") {
		update = true
	}
	request["NewResourceGroupId"] = d.Get("resource_group_id")
	request["ResourceType"] = "bandwidthpackage"
	if update {
		wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
		err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
			response, err = client.RpcPost("Vpc", "2016-04-28", action, query, request, true)
			if err != nil {
				if helper.NeedRetry(err) {
					wait()
					return retry.RetryableError(err)
				}
				return retry.NonRetryableError(err)
			}
			return nil
		})
		logs.AddDebug(action, response, request)
		if err != nil {
			return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
		}
	}
	update = false
	action = "DeletionProtection"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["InstanceId"] = d.Id()
	request["RegionId"] = client.RegionId
	request["ClientToken"] = helper.BuildClientToken(action)
	if d.HasChange("deletion_protection") {
		update = true
	}
	request["ProtectionEnable"] = d.Get("deletion_protection")
	request["Type"] = "CBWP"
	if update {
		wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
		err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
			response, err = client.RpcPost("Vpc", "2016-04-28", action, query, request, true)
			if err != nil {
				if helper.NeedRetry(err) {
					wait()
					return retry.RetryableError(err)
				}
				return retry.NonRetryableError(err)
			}
			return nil
		})
		logs.AddDebug(action, response, request)
		if err != nil {
			return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
		}
	}

	if d.HasChange("tags") {
		cbwpService := NewService(client)
		if err := cbwpService.SetResourceTags(d, "COMMONBANDWIDTHPACKAGE"); err != nil {
			return errs.WrapError(err)
		}
	}
	d.Partial(false)
	return resourceCommonBandwidthPackageRead(d, meta)
}

func resourceCommonBandwidthPackageDelete(d *schema.ResourceData, meta interface{}) error {

	enableDelete := true
	if v, ok := d.GetOk("internet_charge_type"); ok {
		if helper.InArray(fmt.Sprint(v), []string{"PayBy95"}) {
			enableDelete = false
			log.Printf("[WARN] Cannot destroy resource alicloud_common_bandwidth_package which internet_charge_type valued PayBy95. Terraform will remove this resource from the state file, however resources may remain.")
		}
	}
	if enableDelete {
		client := meta.(*connectivity.AliyunClient)
		action := "DeleteCommonBandwidthPackage"
		var request map[string]interface{}
		var response map[string]interface{}
		query := make(map[string]interface{})
		var err error
		request = make(map[string]interface{})
		request["BandwidthPackageId"] = d.Id()
		request["RegionId"] = client.RegionId

		if v, ok := d.GetOk("force"); ok {
			request["Force"] = v
		}
		wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
		err = retry.Retry(d.Timeout(schema.TimeoutDelete), func() *retry.RetryError {
			response, err = client.RpcPost("Vpc", "2016-04-28", action, query, request, true)

			if err != nil {
				if errs.IsExpectedErrors(err, []string{"BandwidthPackageOperation.conflict", "OperationConflict", "LastTokenProcessing", "IncorrectStatus", "SystemBusy", "ServiceUnavailable"}) || helper.NeedRetry(err) {
					wait()
					return retry.RetryableError(err)
				}
				return retry.NonRetryableError(err)
			}
			return nil
		})
		logs.AddDebug(action, response, request)

		if err != nil {
			if errs.NotFoundError(err) {
				return nil
			}
			return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
		}

		cbwpService := NewService(client)
		stateConf := helper.BuildStateConf([]string{}, []string{""}, d.Timeout(schema.TimeoutDelete), 15*time.Second, cbwpService.CommonBandwidthPackageStateRefreshFunc(d.Id(), "Status", []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return errs.WrapErrorf(err, errs.IdMsg, d.Id())
		}

	}
	return nil
}

func convertCommonBandwidthPackageInstanceChargeTypeResponse(source interface{}) interface{} {
	source = fmt.Sprint(source)
	switch source {
	case "PrePaid":
		return "Subscription"
	case "PostPaid":
		return "PayAsYouGo"
	}
	return source
}
