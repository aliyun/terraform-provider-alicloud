package alicloud

import (
	"fmt"
	"log"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAliCloudCloudStorageGatewayGateway() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudCloudStorageGatewayGatewayCreate,
		Read:   resourceAliCloudCloudStorageGatewayGatewayRead,
		Update: resourceAliCloudCloudStorageGatewayGatewayUpdate,
		Delete: resourceAliCloudCloudStorageGatewayGatewayDelete,
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
				ValidateFunc: StringInSlice([]string{"File", "Iscsi"}, false),
			},
			"location": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: StringInSlice([]string{"Cloud", "On_Premise"}, false),
			},
			"gateway_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"gateway_class": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: StringInSlice([]string{"Basic", "Standard", "Enhanced", "Advanced"}, false),
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
				ValidateFunc: IntBetween(5, 200),
			},
			"payment_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Default:      "PayAsYouGo",
				ValidateFunc: StringInSlice([]string{"PayAsYouGo", "Subscription"}, false),
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

func resourceAliCloudCloudStorageGatewayGatewayCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	sgwService := SgwService{client}
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

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutCreate)), func() *resource.RetryError {
		response, err = client.RpcPost("sgw", "2018-05-11", action, nil, request, true)
		if err != nil {
			if IsExpectedErrors(err, []string{"BadRequest"}) || NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_cloud_storage_gateway_gateway", action, AlibabaCloudSdkGoERROR)
	}

	if fmt.Sprint(response["Success"]) == "false" {
		return WrapError(fmt.Errorf("%s failed, response: %v", action, response))
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

		err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
			response, err = client.RpcPost("sgw", "2018-05-11", action, nil, request, true)
			if err != nil {
				if IsExpectedErrors(err, []string{"BadRequest"}) || NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})
		addDebug(action, response, request)

		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, "alicloud_cloud_storage_gateway_gateway", action, AlibabaCloudSdkGoERROR)
		}

		if fmt.Sprint(response["Success"]) == "false" {
			return WrapError(fmt.Errorf("%s failed, response: %v", action, response))
		}

		stateConf := BuildStateConf([]string{}, []string{"task.state.completed"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, sgwService.CloudStorageGatewayTaskStateRefreshFunc(d.Id(), fmt.Sprint(response["TaskId"]), []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
	}

	return resourceAliCloudCloudStorageGatewayGatewayRead(d, meta)
}

func resourceAliCloudCloudStorageGatewayGatewayRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	sgwService := SgwService{client}

	object, err := sgwService.DescribeCloudStorageGatewayGateway(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_cloud_storage_gateway_gateway sgwService.DescribeCloudStorageGatewayGateway Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("storage_bundle_id", object["StorageBundleId"])
	d.Set("type", object["Type"])
	d.Set("location", object["Location"])
	d.Set("gateway_name", object["Name"])
	d.Set("gateway_class", object["GatewayClass"])
	d.Set("vswitch_id", object["VSwitchId"])
	d.Set("payment_type", convertCsgGatewayPaymentTypeResp(formatBool(object["IsPostPaid"])))
	d.Set("description", object["Description"])
	d.Set("release_after_expiration", object["IsReleaseAfterExpiration"])
	d.Set("status", object["Status"])

	if publicNetworkBandwidth, ok := object["PublicNetworkBandwidth"]; ok && fmt.Sprint(publicNetworkBandwidth) != "0" {
		d.Set("public_network_bandwidth", formatInt(publicNetworkBandwidth))
	}

	return nil
}

func resourceAliCloudCloudStorageGatewayGatewayUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	sgwService := SgwService{client}
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
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *resource.RetryError {
			response, err = client.RpcPost("sgw", "2018-05-11", action, nil, request, true)
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

		if fmt.Sprint(response["Success"]) == "false" {
			return WrapError(fmt.Errorf("%s failed, response: %v", action, response))
		}

		d.SetPartial("gateway_name")
		d.SetPartial("description")
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
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *resource.RetryError {
			response, err = client.RpcPost("sgw", "2018-05-11", action, nil, modifyGatewayClassReq, true)
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})
		addDebug(action, response, modifyGatewayClassReq)

		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}

		if fmt.Sprint(response["Success"]) == "false" {
			return WrapError(fmt.Errorf("%s failed, response: %v", action, response))
		}

		stateConf := BuildStateConf([]string{}, []string{"task.state.completed"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, sgwService.CloudStorageGatewayTaskStateRefreshFunc(d.Id(), fmt.Sprint(response["TaskId"]), []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}

		d.SetPartial("gateway_class")
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
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *resource.RetryError {
			response, err = client.RpcPost("sgw", "2018-05-11", action, nil, expandGatewayNetworkBandwidthReq, false)
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})
		addDebug(action, response, expandGatewayNetworkBandwidthReq)

		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}

		if fmt.Sprint(response["Success"]) == "false" {
			return WrapError(fmt.Errorf("%s failed, response: %v", action, response))
		}

		stateConf := BuildStateConf([]string{}, []string{"task.state.completed"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, sgwService.CloudStorageGatewayTaskStateRefreshFunc(d.Id(), fmt.Sprint(response["TaskId"]), []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}

		d.SetPartial("public_network_bandwidth")
	}

	d.Partial(false)

	return resourceAliCloudCloudStorageGatewayGatewayRead(d, meta)
}

func resourceAliCloudCloudStorageGatewayGatewayDelete(d *schema.ResourceData, meta interface{}) error {
	if v, ok := d.GetOk("payment_type"); ok && fmt.Sprint(v) == "Subscription" {
		log.Printf("[WARN] Cannot destroy resource alicloud_cloud_storage_gateway_gateway which payment_type valued Subscription. Terraform will remove this resource from the state file, however resources may remain.")
		return nil
	}

	client := meta.(*connectivity.AliyunClient)
	sgwService := SgwService{client}
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
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutDelete)), func() *resource.RetryError {
		response, err = client.RpcPost("sgw", "2018-05-11", action, nil, request, true)
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
		if IsExpectedErrors(err, []string{"GatewayNotExist"}) || NotFoundError(err) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	if fmt.Sprint(response["Success"]) == "false" {
		return WrapError(fmt.Errorf("%s failed, response: %v", action, response))
	}

	if isCloud {
		stateConf := BuildStateConf([]string{}, []string{"task.state.completed"}, d.Timeout(schema.TimeoutDelete), 5*time.Second, sgwService.CloudStorageGatewayTaskStateRefreshFunc(d.Id(), fmt.Sprint(response["TaskId"]), []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
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
