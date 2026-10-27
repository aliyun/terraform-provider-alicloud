package drds

import (
	"fmt"
	"log"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/validators"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourcePolardbxInstance() *schema.Resource {
	return &schema.Resource{
		Create: resourcePolardbxInstanceCreate,
		Read:   resourcePolardbxInstanceRead,
		Update: resourcePolardbxInstanceUpdate,
		Delete: resourcePolardbxInstanceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(61 * time.Minute),
			Update: schema.DefaultTimeout(61 * time.Minute),
			Delete: schema.DefaultTimeout(61 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"auto_renew": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
			},
			"cn_class": {
				Type:     schema.TypeString,
				Required: true,
			},
			"cn_node_count": {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: validators.IntAtLeast(2),
			},
			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"dn_class": {
				Type:     schema.TypeString,
				Required: true,
			},
			"dn_node_count": {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: validators.IntAtLeast(2),
			},
			"dn_storage_space": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"engine_version": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				ValidateFunc: validators.StringInSlice([]string{"5.7", "8.0"}, false),
			},
			"is_read_db_instance": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"pay_type": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				ValidateFunc: validators.StringInSlice([]string{"PREPAY", "POSTPAY"}, false),
			},
			"period": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validators.StringInSlice([]string{"Year", "Month"}, false),
			},
			"primary_db_instance_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"primary_zone": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
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
			"secondary_zone": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"specified_dn_scale": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"specified_dn_spec_map_json": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"storage_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Computed:     true,
				ValidateFunc: validators.StringInSlice([]string{"custom_local_ssd", "cloud_auto"}, false),
			},
			"switch_time": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"switch_time_mode": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"tertiary_zone": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"topology_type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validators.StringInSlice([]string{"1azone", "3azones"}, false),
			},
			"used_time": {
				Type:         schema.TypeInt,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validators.IntAtLeast(1),
			},
			"vswitch_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"vpc_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourcePolardbxInstanceCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "CreateDBInstance"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	request["RegionId"] = client.RegionId
	request["ClientToken"] = helper.BuildClientToken(action)

	if v, ok := d.GetOk("tertiary_zone"); ok {
		request["TertiaryZone"] = v
	}
	if v, ok := d.GetOk("resource_group_id"); ok {
		request["ResourceGroupId"] = v
	}
	if v, ok := d.GetOk("pay_type"); ok && v.(string) == "PREPAY" {
		request["PayType"] = "PREPAY"
		if v, ok := d.GetOk("period"); ok {
			request["Period"] = v
		}
		if v, ok := d.GetOk("used_time"); ok {
			request["UsedTime"] = v
		}
		if v, ok := d.GetOkExists("auto_renew"); ok {
			request["AutoRenew"] = v
		}
	} else {
		request["PayType"] = "POSTPAY"
	}
	if v, ok := d.GetOk("description"); ok {
		request["Description"] = v
	}
	request["VSwitchId"] = d.Get("vswitch_id")
	request["EngineVersion"] = d.Get("engine_version")
	if request["EngineVersion"] == "" {
		request["EngineVersion"] = "5.7"
	}
	request["NetworkType"] = "vpc"
	request["ZoneId"] = "null"
	request["PrimaryZone"] = d.Get("primary_zone")
	request["CnClass"] = d.Get("cn_class")
	request["TopologyType"] = d.Get("topology_type")
	request["DnClass"] = d.Get("dn_class")
	if v, ok := d.GetOk("secondary_zone"); ok {
		request["SecondaryZone"] = v
	}
	request["CNNodeCount"] = d.Get("cn_node_count")
	request["VPCId"] = d.Get("vpc_id")
	request["DNNodeCount"] = d.Get("dn_node_count")
	if v, ok := d.GetOkExists("is_read_db_instance"); ok {
		request["IsReadDBInstance"] = v
	}
	if v, ok := d.GetOk("primary_db_instance_name"); ok {
		request["PrimaryDBInstanceName"] = v
	}
	if v, ok := d.GetOk("storage_type"); ok {
		request["StorageType"] = v
	}
	if v, ok := d.GetOk("dn_storage_space"); ok {
		request["DnStorageSpace"] = v
	}
	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
		response, err = client.RpcPost("polardbx", "2020-02-02", action, query, request, true)
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
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, "alicloud_drds_polardbx_instance", action, errs.AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(response["DBInstanceName"]))

	drdsService := NewService(client)
	stateConf := helper.BuildStateConf([]string{}, []string{"Running"}, d.Timeout(schema.TimeoutCreate), 60*time.Second, drdsService.PolardbxInstanceStateRefreshFunc(d.Id(), "Status", []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return errs.WrapErrorf(err, errs.IdMsg, d.Id())
	}

	return resourcePolardbxInstanceUpdate(d, meta)
}

func resourcePolardbxInstanceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	drdsService := NewService(client)

	objectRaw, err := drdsService.DescribePolardbxInstance(d.Id())
	if err != nil {
		if !d.IsNewResource() && errs.NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_drds_polardbx_instance DescribePolardbxInstance Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return errs.WrapError(err)
	}

	d.Set("cn_class", objectRaw["CnNodeClassCode"])
	d.Set("cn_node_count", objectRaw["CnNodeCount"])
	d.Set("create_time", objectRaw["CreateTime"])
	d.Set("description", objectRaw["Description"])
	d.Set("dn_class", objectRaw["DnNodeClassCode"])
	d.Set("dn_node_count", objectRaw["DnNodeCount"])
	d.Set("dn_storage_space", objectRaw["DnStorageSpace"])
	d.Set("engine_version", objectRaw["EngineVersion"])
	d.Set("pay_type", convertPolardbxInstanceDBInstancePayTypeResponse(objectRaw["PayType"]))
	d.Set("primary_zone", objectRaw["PrimaryZone"])
	d.Set("region_id", objectRaw["RegionId"])
	d.Set("resource_group_id", objectRaw["ResourceGroupId"])
	d.Set("secondary_zone", objectRaw["SecondaryZone"])
	d.Set("status", objectRaw["Status"])
	d.Set("storage_type", objectRaw["StorageType"])
	d.Set("tertiary_zone", objectRaw["TertiaryZone"])
	d.Set("topology_type", objectRaw["TopologyType"])
	d.Set("vswitch_id", objectRaw["VSwitchId"])
	d.Set("vpc_id", objectRaw["VPCId"])

	return nil
}

func resourcePolardbxInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	update := false
	d.Partial(true)

	var err error
	action := "UpdatePolarDBXInstanceNode"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["DBInstanceName"] = d.Id()
	request["RegionId"] = client.RegionId
	request["ClientToken"] = helper.BuildClientToken(action)
	if !d.IsNewResource() && d.HasChange("cn_node_count") {
		update = true
	}
	request["CNNodeCount"] = d.Get("cn_node_count")
	if !d.IsNewResource() && d.HasChange("dn_node_count") {
		update = true
	}
	request["DNNodeCount"] = d.Get("dn_node_count")
	if update {
		wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
		err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
			response, err = client.RpcPost("polardbx", "2020-02-02", action, query, request, true)
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
		drdsService := NewService(client)
		stateConf := helper.BuildStateConf([]string{}, []string{"Running"}, d.Timeout(schema.TimeoutUpdate), 60*time.Second, drdsService.PolardbxInstanceStateRefreshFunc(d.Id(), "Status", []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return errs.WrapErrorf(err, errs.IdMsg, d.Id())
		}
	}
	update = false
	action = "ModifyDBInstanceClass"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["DBInstanceName"] = d.Id()
	request["RegionId"] = client.RegionId
	request["ClientToken"] = helper.BuildClientToken(action)
	if !d.IsNewResource() && d.HasChange("cn_class") {
		update = true
	}
	request["CnClass"] = d.Get("cn_class")
	if !d.IsNewResource() && d.HasChange("dn_class") {
		update = true
	}
	request["DnClass"] = d.Get("dn_class")
	if !d.IsNewResource() && d.HasChange("dn_storage_space") {
		update = true
		request["DnStorageSpace"] = d.Get("dn_storage_space")
	}
	if d.HasChange("specified_dn_scale") {
		if v, ok := d.GetOkExists("specified_dn_scale"); ok {
			request["SpecifiedDNScale"] = v
		}
	}
	if d.HasChange("specified_dn_spec_map_json") {
		if v, ok := d.GetOk("specified_dn_spec_map_json"); ok {
			request["SpecifiedDNSpecMapJson"] = v
		}
	}
	if d.HasChange("switch_time_mode") {
		if v, ok := d.GetOk("switch_time_mode"); ok {
			request["SwitchTimeMode"] = v
		}
	}
	if d.HasChange("switch_time") {
		if v, ok := d.GetOk("switch_time"); ok {
			request["SwitchTime"] = v
		}
	}
	if update {
		wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
		err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
			response, err = client.RpcPost("polardbx", "2020-02-02", action, query, request, true)
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
		drdsService := NewService(client)
		stateConf := helper.BuildStateConf([]string{}, []string{"Running"}, d.Timeout(schema.TimeoutUpdate), 60*time.Second, drdsService.PolardbxInstanceStateRefreshFunc(d.Id(), "Status", []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return errs.WrapErrorf(err, errs.IdMsg, d.Id())
		}
	}
	update = false
	action = "ChangeResourceGroup"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["ResourceId"] = d.Id()
	request["RegionId"] = client.RegionId
	if _, ok := d.GetOk("resource_group_id"); ok && !d.IsNewResource() && d.HasChange("resource_group_id") {
		update = true
	}
	request["NewResourceGroupId"] = d.Get("resource_group_id")
	request["ResourceType"] = "PolarDBXInstance"
	if update {
		wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
		err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
			response, err = client.RpcPost("polardbx", "2020-02-02", action, query, request, true)
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
	action = "ModifyDBInstanceDescription"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["DBInstanceName"] = d.Id()
	request["RegionId"] = client.RegionId
	if !d.IsNewResource() && d.HasChange("description") {
		update = true
	}
	request["DBInstanceDescription"] = d.Get("description")
	if update {
		wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
		err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
			response, err = client.RpcPost("polardbx", "2020-02-02", action, query, request, true)
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

	d.Partial(false)
	return resourcePolardbxInstanceRead(d, meta)
}

func resourcePolardbxInstanceDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	drdsService := NewService(client)

	// Determine the billing type to pick the right destroy path. Prepaid
	// (subscription) instances are unsubscribed and released via BSS OpenApi
	// RefundInstance; Postpaid (pay-as-you-go) instances are released via
	// DeleteDBInstance.
	object, err := drdsService.DescribePolardbxInstance(d.Id())
	if err != nil {
		if errs.NotFoundError(err) || errs.IsExpectedErrors(err, []string{"DBInstance.NotFound"}) {
			return nil
		}
		return errs.WrapError(err)
	}
	if fmt.Sprint(object["PayType"]) == "Prepaid" {
		// Subscription instances cannot be released by DeleteDBInstance before the
		// subscription period expires. Unsubscribe and release the cloud resource via
		// BSS OpenApi RefundInstance instead, so Terraform does not "forget" a
		// still-renewing, still-billing instance.
		// Verify the billing identity at runtime instead of assuming the RPC service
		// name equals the billing ProductCode. QueryAvailableInstanceByInstanceID
		// resolves the real ProductCode/ProductType for this instance solely from
		// its InstanceID via the domestic BSS QueryAvailableInstances API, with no
		// ProductCode filter and no international-endpoint fallback; the billing
		// ProductCode of a PolarDB-X 2.0 subscription instance is "drds", not
		// "polardbx". If it is not found under this account, refund is not supported
		// and the resource stays in state.
		billingInstance, err := drdsService.QueryAvailableInstanceByInstanceID(d.Id())
		if err != nil {
			return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), "QueryAvailableInstances", errs.AlibabaCloudSdkGoERROR)
		}
		productCode := fmt.Sprint(billingInstance["ProductCode"])
		productType := fmt.Sprint(billingInstance["ProductType"])

		action := "RefundInstance"
		var request map[string]interface{}
		var response map[string]interface{}
		query := make(map[string]interface{})
		request = make(map[string]interface{})
		request["InstanceId"] = d.Id()
		request["ClientToken"] = helper.BuildClientToken(action)
		request["ImmediatelyRelease"] = "1"
		request["ProductCode"] = productCode
		request["ProductType"] = productType
		var endpoint string
		if client.IsInternationalAccount() {
			endpoint = connectivity.BssOpenAPIEndpointInternational
		}
		wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
		err = retry.Retry(d.Timeout(schema.TimeoutDelete), func() *retry.RetryError {
			response, err = client.RpcPostWithEndpoint("BssOpenApi", "2017-12-14", action, query, request, true, endpoint)
			if err != nil {
				if helper.NeedRetry(err) || errs.IsExpectedErrors(err, []string{"ResourceStatusError"}) {
					wait()
					return retry.RetryableError(err)
				}
				return retry.NonRetryableError(err)
			}
			return nil
		})
		logs.AddDebug(action, response, request)
		if err != nil {
			// BSS refund-system NotFound (ResourceNotExists) does NOT prove the
			// PolarDB-X product instance is gone: the DescribePolardbxInstance call
			// above already confirmed a Prepaid instance exists. The refund record
			// may be missing while the product still runs and keeps billing. Treat
			// refund-system NotFound and ExistRefundingOrderError (a refund already
			// in progress — the instance may have expired, been refunded, or been
			// released) the same way: fall through to the product-side release
			// confirmation below. If the product is still present or the Describe
			// query fails, surface an executable error and keep the resource in
			// state; only a genuine product-side NotFound (confirmed by
			// DescribePolardbxInstance) lets destroy succeed. This avoids silently
			// forgetting a still-existing instance when the BSS refund record is
			// missing.
			if !errs.IsExpectedErrors(err, []string{"ResourceNotExists", "ExistRefundingOrderError"}) && !errs.NotFoundError(err) {
				// The commodity/account does not support API unsubscription (BSS
				// returns CommodityNotSupported, NotApplicable or NotAuthorized).
				// Surface an actionable error pointing to the console unsubscribe
				// page so the operator can finish by hand; the resource stays in
				// state (Delete returning an error keeps it in Terraform state).
				if errs.IsExpectedErrors(err, []string{"CommodityNotSupported", "NotApplicable", "NotAuthorized"}) {
					return fmt.Errorf("destroying PolarDB-X subscription instance %q failed: BssOpenApi RefundInstance rejected the refund (%w). "+
						"This commodity or account does not support API unsubscription. Please unsubscribe the instance in the BSS console (费用中心 > 退订管理, https://usercenter2.aliyun.com/) and then remove it from Terraform state with `terraform state rm alicloud_drds_polardbx_instance.<name>` (replacing <name> with the resource block label from the configuration). "+
						"The resource remains in Terraform state until it is released manually; it is not silently forgotten.", d.Id(), err)
				}
				// Refund failed for another reason: surface the executable error so
				// the resource stays in state instead of being silently forgotten.
				return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
			}
		} else {
			// Confirm the refund succeeded: RefundInstance returns a Success flag that
			// must be true; a nil error with Success != true means the refund was
			// rejected (e.g. ineligible SKU or missing bss:RefundBatchRemainRefund
			// permission), so surface an executable error and keep the resource in state.
			if fmt.Sprint(response["Success"]) == "null" || fmt.Sprint(response["Success"]) != "true" {
				return errs.WrapErrorf(fmt.Errorf("%s failed, response: %v", action, response), errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
			}
		}
		// Confirm the release before considering destroy complete: poll the instance
		// until DescribePolardbxInstance reports NotFound.
		err = retry.Retry(d.Timeout(schema.TimeoutDelete), func() *retry.RetryError {
			if _, e := drdsService.DescribePolardbxInstance(d.Id()); e != nil {
				if errs.NotFoundError(e) {
					return nil
				}
				return retry.NonRetryableError(e)
			}
			return retry.RetryableError(fmt.Errorf("waiting for PolarDB-X instance %q to be released after refund", d.Id()))
		})
		if err != nil {
			return errs.WrapErrorf(err, errs.IdMsg, d.Id(), "wait for release after refund")
		}
		return nil
	}

	action := "DeleteDBInstance"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	request = make(map[string]interface{})
	request["DBInstanceName"] = d.Id()
	request["RegionId"] = client.RegionId

	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutDelete), func() *retry.RetryError {
		response, err = client.RpcPost("polardbx", "2020-02-02", action, query, request, true)
		if err != nil {
			if errs.IsExpectedErrors(err, []string{"DBInstance.InOrder"}) || helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	logs.AddDebug(action, response, request)

	if err != nil {
		if errs.IsExpectedErrors(err, []string{"DBInstance.NotFound"}) || errs.NotFoundError(err) {
			return nil
		}
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
	}

	stateConf := helper.BuildStateConf([]string{}, []string{""}, d.Timeout(schema.TimeoutDelete), 60*time.Second, drdsService.DescribeAsyncPolardbxInstanceStateRefreshFunc(d, response, "$.Items", []string{}))
	if jobDetail, err := stateConf.WaitForState(); err != nil {
		return errs.WrapErrorf(err, errs.IdMsg, d.Id(), jobDetail)
	}

	return nil
}
