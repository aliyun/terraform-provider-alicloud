package alicloud

import (
	"fmt"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/bssopenapi"
	waf_openapi "github.com/aliyun/alibaba-cloud-sdk-go/services/waf-openapi"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func resourceAlicloudWafInstance() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudWafInstanceCreate,
		Read:   resourceAlicloudWafInstanceRead,
		Update: resourceAlicloudWafInstanceUpdate,
		Delete: resourceAlicloudWafInstanceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"big_screen": {
				Type:     schema.TypeString,
				Required: true,
			},
			"exclusive_ip_package": {
				Type:     schema.TypeString,
				Required: true,
			},
			"ext_bandwidth": {
				Type:     schema.TypeString,
				Required: true,
			},
			"ext_domain_package": {
				Type:     schema.TypeString,
				Required: true,
			},
			"log_storage": {
				Type:     schema.TypeString,
				Required: true,
			},
			"log_time": {
				Type:     schema.TypeString,
				Required: true,
			},
			"modify_type": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"package_code": {
				Type:     schema.TypeString,
				Required: true,
			},
			"period": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
			},
			"prefessional_service": {
				Type:     schema.TypeString,
				Required: true,
			},
			"renew_period": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
			},
			"renewal_status": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"resource_group_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"status": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"subscription_type": {
				Type:     schema.TypeString,
				Required: true,
			},
			"waf_log": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func resourceAlicloudWafInstanceCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	request := bssopenapi.CreateCreateInstanceRequest()
	if v, ok := d.GetOk("period"); ok {
		request.Period = requests.NewInteger(v.(int))
	}
	request.ProductCode = "waf"
	request.ProductType = "waf"
	if v, ok := d.GetOk("renew_period"); ok {
		request.RenewPeriod = requests.NewInteger(v.(int))
	}
	if v, ok := d.GetOk("renewal_status"); ok {
		request.RenewalStatus = v.(string)
	}
	request.SubscriptionType = d.Get("subscription_type").(string)
	request.Parameter = &[]bssopenapi.CreateInstanceParameter{
		{
			Code:  "BigScreen",
			Value: d.Get("big_screen").(string),
		},
		{
			Code:  "ExclusiveIpPackage",
			Value: d.Get("exclusive_ip_package").(string),
		},
		{
			Code:  "ExtBandwidth",
			Value: d.Get("ext_bandwidth").(string),
		},
		{
			Code:  "ExtDomainPackage",
			Value: d.Get("ext_domain_package").(string),
		},
		{
			Code:  "LogStorage",
			Value: d.Get("log_storage").(string),
		},
		{
			Code:  "LogTime",
			Value: d.Get("log_time").(string),
		},
		{
			Code:  "PackageCode",
			Value: d.Get("package_code").(string),
		},
		{
			Code:  "PrefessionalService",
			Value: d.Get("prefessional_service").(string),
		},
		{
			Code:  "Region",
			Value: client.RegionId,
		},
		{
			Code:  "WafLog",
			Value: d.Get("waf_log").(string),
		},
	}
	raw, err := client.WithBssopenapiClient(func(bssopenapiClient *bssopenapi.Client) (interface{}, error) {
		return bssopenapiClient.CreateInstance(request)
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_waf_instance", request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw)
	response, _ := raw.(*bssopenapi.CreateInstanceResponse)
	if !response.Success {
		return WrapErrorf(fmt.Errorf("%v", response), DefaultErrorMsg, "alicloud_waf_instance", request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	d.SetId(fmt.Sprintf("%v", response.Data.InstanceId))

	return resourceAlicloudWafInstanceUpdate(d, meta)
}
func resourceAlicloudWafInstanceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	waf_openapiService := Waf_openapiService{client}
	object, err := waf_openapiService.DescribeWafInstance(d.Id())
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("status", object.InstanceInfo.Status)
	d.Set("subscription_type", object.InstanceInfo.SubscriptionType)
	return nil
}
func resourceAlicloudWafInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	update := false
	request := bssopenapi.CreateModifyInstanceRequest()
	request.InstanceId = d.Id()
	if d.HasChange("modify_type") {
		update = true
	}
	if !d.IsNewResource() && d.HasChange("subscription_type") {
		update = true
	}
	request.ProductType = "waf"
	request.SubscriptionType = d.Get("subscription_type").(string)
	request.ProductCode = "waf"
	request.ModifyType = d.Get("modify_type").(string)
	request.Parameter = &[]bssopenapi.ModifyInstanceParameter{
		{
			Code:  "BigScreen",
			Value: d.Get("big_screen").(string),
		},
		{
			Code:  "ExclusiveIpPackage",
			Value: d.Get("exclusive_ip_package").(string),
		},
		{
			Code:  "ExtBandwidth",
			Value: d.Get("ext_bandwidth").(string),
		},
		{
			Code:  "ExtDomainPackage",
			Value: d.Get("ext_domain_package").(string),
		},
		{
			Code:  "LogStorage",
			Value: d.Get("log_storage").(string),
		},
		{
			Code:  "LogTime",
			Value: d.Get("log_time").(string),
		},
		{
			Code:  "PackageCode",
			Value: d.Get("package_code").(string),
		},
		{
			Code:  "PrefessionalService",
			Value: d.Get("prefessional_service").(string),
		},
		{
			Code:  "WafLog",
			Value: d.Get("waf_log").(string),
		},
	}
	if update {
		raw, err := client.WithBssopenapiClient(func(bssopenapiClient *bssopenapi.Client) (interface{}, error) {
			return bssopenapiClient.ModifyInstance(request)
		})
		addDebug(request.GetActionName(), raw)
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		response, _ := raw.(*bssopenapi.ModifyInstanceResponse)
		if !response.Success {
			return WrapErrorf(fmt.Errorf("%v", response), DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
		}

	}
	return resourceAlicloudWafInstanceRead(d, meta)
}
func resourceAlicloudWafInstanceDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	request := waf_openapi.CreateDeleteInstanceRequest()
	request.InstanceId = d.Id()
	if v, ok := d.GetOk("resource_group_id"); ok {
		request.ResourceGroupId = v.(string)
	}
	raw, err := client.WithWafOpenapiClient(func(waf_openapiClient *waf_openapi.Client) (interface{}, error) {
		return waf_openapiClient.DeleteInstance(request)
	})
	addDebug(request.GetActionName(), raw)
	if err != nil {
		if IsExpectedErrors(err, []string{"ComboError"}) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	return nil
}
