package alicloud

import (
	"fmt"
	"log"
	"time"

	util "github.com/alibabacloud-go/tea-utils/service"

	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAliyunVpnCustomerGateway() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliyunVpnCustomerGatewayCreate,
		Read:   resourceAliyunVpnCustomerGatewayRead,
		Update: resourceAliyunVpnCustomerGatewayUpdate,
		Delete: resourceAliyunVpnCustomerGatewayDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(1 * time.Minute),
			Delete: schema.DefaultTimeout(1 * time.Minute),
			Update: schema.DefaultTimeout(1 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"ip_address": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.SingleIP(),
			},
			"name": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringLenBetween(2, 128),
			},
			"description": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringLenBetween(2, 256),
			},
			"asn": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
		},
	}
}

func resourceAliyunVpnCustomerGatewayCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "CreateCustomerGateway"
	request := make(map[string]interface{})
	conn, err := client.NewVpcClient()
	if err != nil {
		return WrapError(err)
	}

	request["RegionId"] = client.RegionId
	request["IpAddress"] = d.Get("ip_address").(string)

	if v, ok := d.GetOk("name"); ok {
		request["Name"] = v
	}

	if v, ok := d.GetOk("description"); ok {
		request["Description"] = v
	}

	if v, ok := d.GetOk("asn"); ok {
		request["Asn"] = v
	}

	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		request["ClientToken"] = buildClientToken("CreateCustomerGateway")
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2016-04-28"), StringPointer("AK"), nil, request, &runtime)
		if err != nil {
			if IsExpectedErrors(err, []string{"OperationConflict"}) || NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(action, response, request)
		return nil
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_vpn_customer_gateway", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(response["CustomerGatewayId"]))
	return resourceAliyunVpnCustomerGatewayRead(d, meta)
}

func resourceAliyunVpnCustomerGatewayRead(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	vpcService := VpcService{client}
	object, err := vpcService.DescribeVpnCustomerGateway(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_vpn_customer_gateway vpcService.DescribeVpnCustomerGateway Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("ip_address", object["IpAddress"])
	d.Set("name", object["Name"])
	d.Set("description", object["Description"])
	d.Set("asn", object["Asn"])

	return nil
}

func resourceAliyunVpnCustomerGatewayUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	action := "ModifyCustomerGatewayAttribute"
	var response map[string]interface{}
	d.Partial(false)
	conn, err := client.NewVpcClient()
	if err != nil {
		return WrapError(err)
	}
	request := map[string]interface{}{
		"CustomerGatewayId": d.Id(),
	}

	request["RegionId"] = client.RegionId

	update := false
	if d.HasChange("name") {
		update = true
		if v, ok := d.GetOk("name"); ok {
			request["Name"] = v
		}
	}

	if d.HasChange("description") {
		update = true
		if v, ok := d.GetOk("description"); ok {
			request["Description"] = v
		}
	}

	if update {
		request["ClientToken"] = buildClientToken(action)
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2016-04-28"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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
	}

	return resourceAliyunVpnCustomerGatewayRead(d, meta)
}

func resourceAliyunVpnCustomerGatewayDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	action := "DeleteCustomerGateway"
	var response map[string]interface{}
	conn, err := client.NewVpcClient()
	if err != nil {
		return WrapError(err)
	}
	request := map[string]interface{}{
		"CustomerGatewayId": d.Id(),
	}

	request["RegionId"] = client.RegionId
	request["ClientToken"] = buildClientToken(action)
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2016-04-28"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
		if err != nil {
			if IsExpectedErrors(err, []string{"VpnGateway.Configuring"}) || NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(action, response, request)
		return nil
	})

	if err != nil {
		if IsExpectedErrors(err, []string{"InvalidCustomerGatewayInstanceId.NotFound"}) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	return nil
}
