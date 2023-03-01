package alicloud

import (
	"fmt"
	"log"
	"regexp"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAlicloudVpcIpv4Gateway() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudVpcIpv4GatewayCreate,
		Read:   resourceAlicloudVpcIpv4GatewayRead,
		Update: resourceAlicloudVpcIpv4GatewayUpdate,
		Delete: resourceAlicloudVpcIpv4GatewayDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(1 * time.Minute),
			Delete: schema.DefaultTimeout(1 * time.Minute),
			Update: schema.DefaultTimeout(1 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"dry_run": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"ipv4_gateway_description": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.All(validation.StringLenBetween(2, 256), validation.StringDoesNotMatch(regexp.MustCompile(`(^http://.*)|(^https://.*)`), "It cannot begin with \"http://\", \"https://\".")),
			},
			"ipv4_gateway_name": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringMatch(regexp.MustCompile(`^[a-zA-Z][a-zA-Z0-9-._]{1,127}$`), "The name must be 2 to 128 characters in length, and can contain letters, digits, periods (.), underscores (_), and hyphens (-). The name must start with a letter."),
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"vpc_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"enabled": {
				Type:     schema.TypeBool,
				Optional: true,
			},
		},
	}
}

func resourceAlicloudVpcIpv4GatewayCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	vpcService := VpcService{client}
	action := "CreateIpv4Gateway"
	request := make(map[string]interface{})
	conn, err := client.NewVpcClient()
	if err != nil {
		return WrapError(err)
	}
	if v, ok := d.GetOkExists("dry_run"); ok {
		request["DryRun"] = v
	}
	if v, ok := d.GetOk("ipv4_gateway_description"); ok {
		request["Ipv4GatewayDescription"] = v
	}
	if v, ok := d.GetOk("ipv4_gateway_name"); ok {
		request["Ipv4GatewayName"] = v
	}
	request["RegionId"] = client.RegionId
	request["VpcId"] = d.Get("vpc_id")
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutCreate)), func() *resource.RetryError {
		request["ClientToken"] = buildClientToken("CreateIpv4Gateway")
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2016-04-28"), StringPointer("AK"), nil, request, &runtime)
		if err != nil {
			if IsExpectedErrors(err, []string{"OperationConflict", "IncorrectStatus.Vpc", "IncorrectStatus.%s", "ServiceUnavailable", "LastTokenProcessing", "SystemBusy"}) || NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_vpc_ipv4_gateway", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(response["Ipv4GatewayId"]))

	stateConf := BuildStateConf([]string{"Creating"}, []string{"Created"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, vpcService.VpcIpv4GatewayStateRefreshFunc(d.Id(), []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return resourceAlicloudVpcIpv4GatewayUpdate(d, meta)
}

func resourceAlicloudVpcIpv4GatewayRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	vpcService := VpcService{client}
	object, err := vpcService.DescribeVpcIpv4Gateway(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_vpc_ipv4_gateway vpcService.DescribeVpcIpv4Gateway Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	d.Set("ipv4_gateway_description", object["Ipv4GatewayDescription"])
	d.Set("ipv4_gateway_name", object["Ipv4GatewayName"])
	d.Set("status", object["Status"])
	d.Set("vpc_id", object["VpcId"])
	d.Set("enabled", object["Enabled"])
	return nil
}

func resourceAlicloudVpcIpv4GatewayUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	vpcService := VpcService{client}
	var response map[string]interface{}
	conn, err := client.NewVpcClient()
	if err != nil {
		return WrapError(err)
	}
	update := false
	request := map[string]interface{}{
		"Ipv4GatewayId": d.Id(),
	}
	request["RegionId"] = client.RegionId
	if !d.IsNewResource() && d.HasChange("ipv4_gateway_description") {
		update = true
		if v, ok := d.GetOk("ipv4_gateway_description"); ok {
			request["Ipv4GatewayDescription"] = v
		}
	}
	if !d.IsNewResource() && d.HasChange("ipv4_gateway_name") {
		update = true
		if v, ok := d.GetOk("ipv4_gateway_name"); ok {
			request["Ipv4GatewayName"] = v
		}
	}
	if update {
		if v, ok := d.GetOkExists("dry_run"); ok {
			request["DryRun"] = v
		}
		action := "UpdateIpv4GatewayAttribute"
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *resource.RetryError {
			request["ClientToken"] = buildClientToken("UpdateIpv4GatewayAttribute")
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2016-04-28"), StringPointer("AK"), nil, request, &runtime)
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
	}

	if d.HasChange("enabled") {
		object, err := vpcService.DescribeVpcIpv4Gateway(d.Id())
		if err != nil {
			return WrapError(err)
		}
		target := d.Get("enabled")
		if object["Enabled"].(bool) != true {
			if target.(bool) == true {
				request = map[string]interface{}{
					"Ipv4GatewayId": d.Id(),
				}
				request["RegionId"] = client.RegionId
				action := "EnableVpcIpv4Gateway"
				request["ClientToken"] = buildClientToken("EnableVpcIpv4Gateway")
				runtime := util.RuntimeOptions{}
				runtime.SetAutoretry(true)
				wait := incrementalWait(3*time.Second, 3*time.Second)
				err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *resource.RetryError {
					response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2016-04-28"), StringPointer("AK"), nil, request, &runtime)
					if err != nil {
						if IsExpectedErrors(err, []string{"OperationConflict", "OperationFailed.LastTokenProcessing", "IncorrectStatus.%s", "ServiceUnavailable", "LastTokenProcessing", "SystemBusy"}) || NeedRetry(err) {
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
				stateConf := BuildStateConf([]string{}, []string{"Created"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, vpcService.VpcIpv4GatewayStateRefreshFunc(d.Id(), []string{}))
				if _, err := stateConf.WaitForState(); err != nil {
					return WrapErrorf(err, IdMsg, d.Id())
				}
			}
		}
	}

	return resourceAlicloudVpcIpv4GatewayRead(d, meta)
}

func resourceAlicloudVpcIpv4GatewayDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	action := "DeleteIpv4Gateway"
	vpcService := VpcService{client}
	var response map[string]interface{}
	conn, err := client.NewVpcClient()
	if err != nil {
		return WrapError(err)
	}
	request := map[string]interface{}{
		"Ipv4GatewayId": d.Id(),
	}

	if v, ok := d.GetOkExists("dry_run"); ok {
		request["DryRun"] = v
	}
	request["RegionId"] = client.RegionId
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutDelete)), func() *resource.RetryError {
		request["ClientToken"] = buildClientToken("DeleteIpv4Gateway")
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2016-04-28"), StringPointer("AK"), nil, request, &runtime)
		if err != nil {
			if IsExpectedErrors(err, []string{"OperationConflict", "IncorrectStatus.Ipv4Gateway", "IncorrectStatus.%s", "ServiceUnavailable", "LastTokenProcessing", "SystemBusy"}) || NeedRetry(err) {
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
	stateConf := BuildStateConf([]string{}, []string{}, d.Timeout(schema.TimeoutDelete), 5*time.Second, vpcService.VpcIpv4GatewayStateRefreshFunc(d.Id(), []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}
	return nil
}
