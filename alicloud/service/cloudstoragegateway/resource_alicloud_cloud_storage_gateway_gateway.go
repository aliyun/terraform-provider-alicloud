package cloudstoragegateway

import (
	"fmt"
	"log"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/convert"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/validators"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceAliCloudCloudStorageGatewayGateway() *schema.Resource {
	return &schema.Resource{
		Create: ResourceAliCloudCloudStorageGatewayGatewayCreate,
		Read:   ResourceAliCloudCloudStorageGatewayGatewayRead,
		Update: ResourceAliCloudCloudStorageGatewayGatewayUpdate,
		Delete: ResourceAliCloudCloudStorageGatewayGatewayDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"storage_bundle_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validators.StringInSlice([]string{"File", "Iscsi"}, false),
			},
			"location": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validators.StringInSlice([]string{"Cloud", "On_Premise"}, false),
			},
			"gateway_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"gateway_class": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validators.StringInSlice([]string{"Basic", "Standard", "Enhanced", "Advanced"}, false),
			},
			"vswitch_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"public_network_bandwidth": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validators.IntBetween(5, 200),
			},
			"payment_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Default:      "PayAsYouGo",
				ValidateFunc: validators.StringInSlice([]string{"PayAsYouGo", "Subscription"}, false),
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"release_after_expiration": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
			},
			"reason_type": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"reason_detail": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func ResourceAliCloudCloudStorageGatewayGatewayCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	sgwService := NewService(client)
	var response map[string]interface{}
	action := "CreateGateway"
	request := make(map[string]interface{})
	var err error

	request["StorageBundleId"] = d.Get("storage_bundle_id")
	request["Type"] = d.Get("type")
	request["Location"] = d.Get("location")
	request["Name"] = d.Get("gateway_name")

	if v, ok := d.GetOk("gateway_class"); ok {
		request["GatewayClass"] = v
	}

	if v, ok := d.GetOk("vswitch_id"); ok {
		request["VSwitchId"] = v
	}

	if v, ok := d.GetOk("public_network_bandwidth"); ok {
		request["PublicNetworkBandwidth"] = v
	}

	if v, ok := d.GetOk("payment_type"); ok {
		request["PostPaid"] = convertCsgGatewayPaymentTypeReq(v.(string))
	}

	if v, ok := d.GetOk("description"); ok {
		request["Description"] = v
	}

	if v, ok := d.GetOkExists("release_after_expiration"); ok {
		request["ReleaseAfterExpiration"] = v
	}

	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutCreate)), func() *retry.RetryError {
		response, err = client.RpcPost("sgw", "2018-05-11", action, nil, request, true)
		if err != nil {
			if errs.IsExpectedErrors(err, []string{"BadRequest"}) || helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	logs.AddDebug(action, response, request)

	if err != nil {
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, "alicloud_cloud_storage_gateway_gateway", action, errs.AlibabaCloudSdkGoERROR)
	}

	if fmt.Sprint(response["Success"]) == "false" {
		return errs.WrapError(fmt.Errorf("%s failed, response: %v", action, response))
	}

	d.SetId(fmt.Sprint(response["GatewayId"]))

	if d.Get("location").(string) == "Cloud" && d.Get("payment_type").(string) == "PayAsYouGo" {
		action = "DeployGateway"
		request = map[string]interface{}{
			"GatewayId": d.Id(),
		}

		if id, ok := response["GatewayId"]; ok {
			request["GatewayId"] = fmt.Sprint(id)
		}

		if v, ok := d.GetOk("gateway_class"); ok {
			request["GatewayClass"] = v
		}

		err = retry.Retry(d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
			response, err = client.RpcPost("sgw", "2018-05-11", action, nil, request, true)
			if err != nil {
				if errs.IsExpectedErrors(err, []string{"BadRequest"}) || helper.NeedRetry(err) {
					wait()
					return retry.RetryableError(err)
				}
				return retry.NonRetryableError(err)
			}
			return nil
		})
		logs.AddDebug(action, response, request)

		if err != nil {
			return errs.WrapErrorf(err, errs.DefaultErrorMsg, "alicloud_cloud_storage_gateway_gateway", action, errs.AlibabaCloudSdkGoERROR)
		}

		if fmt.Sprint(response["Success"]) == "false" {
			return errs.WrapError(fmt.Errorf("%s failed, response: %v", action, response))
		}

		stateConf := helper.BuildStateConf([]string{}, []string{"task.state.completed"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, sgwService.CloudStorageGatewayTaskStateRefreshFunc(d.Id(), fmt.Sprint(response["TaskId"]), []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return errs.WrapErrorf(err, errs.IdMsg, d.Id())
		}
	}

	return ResourceAliCloudCloudStorageGatewayGatewayRead(d, meta)
}

func ResourceAliCloudCloudStorageGatewayGatewayRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	sgwService := NewService(client)

	object, err := sgwService.DescribeCloudStorageGatewayGateway(d.Id())
	if err != nil {
		if !d.IsNewResource() && errs.NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_cloud_storage_gateway_gateway sgwService.DescribeCloudStorageGatewayGateway Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return errs.WrapError(err)
	}

	d.Set("storage_bundle_id", object["StorageBundleId"])
	d.Set("type", object["Type"])
	d.Set("location", object["Location"])
	d.Set("gateway_name", object["Name"])
	d.Set("gateway_class", object["GatewayClass"])
	d.Set("vswitch_id", object["VSwitchId"])
	isPostPaid, err := convert.FormatBool(object["IsPostPaid"])
	if err != nil {
		return err
	}
	d.Set("payment_type", convertCsgGatewayPaymentTypeResp(isPostPaid))
	d.Set("description", object["Description"])
	d.Set("release_after_expiration", object["IsReleaseAfterExpiration"])
	d.Set("status", object["Status"])

	if publicNetworkBandwidth, ok := object["PublicNetworkBandwidth"]; ok && fmt.Sprint(publicNetworkBandwidth) != "0" {
		public_network_bandwidthVal, err := convert.FormatInt(publicNetworkBandwidth)
		if err != nil {
			return err
		}
		d.Set("public_network_bandwidth", public_network_bandwidthVal)
	}

	return nil
}

func ResourceAliCloudCloudStorageGatewayGatewayUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	sgwService := NewService(client)
	var response map[string]interface{}
	var err error
	d.Partial(true)

	update := false
	request := map[string]interface{}{
		"GatewayId": d.Id(),
	}

	if d.HasChange("gateway_name") {
		update = true
	}
	request["Name"] = d.Get("gateway_name")

	if d.HasChange("description") {
		update = true
	}
	if v, ok := d.GetOk("description"); ok {
		request["Description"] = v
	}

	if update {
		action := "ModifyGateway"
		wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
		err = retry.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *retry.RetryError {
			response, err = client.RpcPost("sgw", "2018-05-11", action, nil, request, true)
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

		if fmt.Sprint(response["Success"]) == "false" {
			return errs.WrapError(fmt.Errorf("%s failed, response: %v", action, response))
		}

	}

	update = false
	modifyGatewayClassReq := map[string]interface{}{
		"GatewayId": d.Id(),
	}

	if d.HasChange("gateway_class") {
		update = true
	}
	if v, ok := d.GetOk("gateway_class"); ok {
		modifyGatewayClassReq["GatewayClass"] = v
	}

	if update {
		action := "ModifyGatewayClass"
		wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
		err = retry.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *retry.RetryError {
			response, err = client.RpcPost("sgw", "2018-05-11", action, nil, modifyGatewayClassReq, true)
			if err != nil {
				if helper.NeedRetry(err) {
					wait()
					return retry.RetryableError(err)
				}
				return retry.NonRetryableError(err)
			}
			return nil
		})
		logs.AddDebug(action, response, modifyGatewayClassReq)

		if err != nil {
			return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
		}

		if fmt.Sprint(response["Success"]) == "false" {
			return errs.WrapError(fmt.Errorf("%s failed, response: %v", action, response))
		}

		stateConf := helper.BuildStateConf([]string{}, []string{"task.state.completed"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, sgwService.CloudStorageGatewayTaskStateRefreshFunc(d.Id(), fmt.Sprint(response["TaskId"]), []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return errs.WrapErrorf(err, errs.IdMsg, d.Id())
		}

	}

	update = false
	expandGatewayNetworkBandwidthReq := map[string]interface{}{
		"GatewayId": d.Id(),
	}

	if d.HasChange("public_network_bandwidth") {
		update = true
	}
	if v, ok := d.GetOk("public_network_bandwidth"); ok {
		expandGatewayNetworkBandwidthReq["NewNetworkBandwidth"] = v
	}

	if update {
		action := "ExpandGatewayNetworkBandwidth"
		wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
		err = retry.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *retry.RetryError {
			response, err = client.RpcPost("sgw", "2018-05-11", action, nil, expandGatewayNetworkBandwidthReq, false)
			if err != nil {
				if helper.NeedRetry(err) {
					wait()
					return retry.RetryableError(err)
				}
				return retry.NonRetryableError(err)
			}
			return nil
		})
		logs.AddDebug(action, response, expandGatewayNetworkBandwidthReq)

		if err != nil {
			return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
		}

		if fmt.Sprint(response["Success"]) == "false" {
			return errs.WrapError(fmt.Errorf("%s failed, response: %v", action, response))
		}

		stateConf := helper.BuildStateConf([]string{}, []string{"task.state.completed"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, sgwService.CloudStorageGatewayTaskStateRefreshFunc(d.Id(), fmt.Sprint(response["TaskId"]), []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return errs.WrapErrorf(err, errs.IdMsg, d.Id())
		}

	}

	d.Partial(false)

	return ResourceAliCloudCloudStorageGatewayGatewayRead(d, meta)
}

func ResourceAliCloudCloudStorageGatewayGatewayDelete(d *schema.ResourceData, meta interface{}) error {
	if v, ok := d.GetOk("payment_type"); ok && fmt.Sprint(v) == "Subscription" {
		return errs.WrapError(errs.Error("'Subscription' instance cannot be destroyed: please unsubscribe it in the console or wait for it to expire. To remove it from the state only, run 'terraform state rm'."))
	}

	client := meta.(*connectivity.AliyunClient)
	sgwService := NewService(client)
	action := "DeleteGateway"
	var response map[string]interface{}
	var err error

	request := map[string]interface{}{
		"GatewayId": d.Id(),
	}

	if v, ok := d.GetOk("reason_type"); ok {
		request["ReasonType"] = v
	}

	if v, ok := d.GetOk("reason_detail"); ok {
		request["ReasonDetail"] = v
	}

	isCloud := d.Get("location") == "Cloud"
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutDelete)), func() *retry.RetryError {
		response, err = client.RpcPost("sgw", "2018-05-11", action, nil, request, true)
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
		if errs.IsExpectedErrors(err, []string{"GatewayNotExist"}) || errs.NotFoundError(err) {
			return nil
		}
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
	}

	if fmt.Sprint(response["Success"]) == "false" {
		return errs.WrapError(fmt.Errorf("%s failed, response: %v", action, response))
	}

	if isCloud {
		stateConf := helper.BuildStateConf([]string{}, []string{"task.state.completed"}, d.Timeout(schema.TimeoutDelete), 5*time.Second, sgwService.CloudStorageGatewayTaskStateRefreshFunc(d.Id(), fmt.Sprint(response["TaskId"]), []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return errs.WrapErrorf(err, errs.IdMsg, d.Id())
		}
	}

	return nil
}

func convertCsgGatewayPaymentTypeReq(source interface{}) interface{} {
	switch source {
	case "PayAsYouGo":
		return true
	case "Subscription":
		return false
	}

	return source
}

func convertCsgGatewayPaymentTypeResp(source interface{}) interface{} {
	switch source {
	case true:
		return "PayAsYouGo"
	case false:
		return "Subscription"
	}

	return source
}
