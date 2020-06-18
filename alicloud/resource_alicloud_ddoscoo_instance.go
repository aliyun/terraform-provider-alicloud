package alicloud

import (
	"fmt"
	"strconv"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/bssopenapi"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ddoscoo"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func resourceAlicloudDdoscooInstance() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudDdoscooInstanceCreate,
		Read:   resourceAlicloudDdoscooInstanceRead,
		Update: resourceAlicloudDdoscooInstanceUpdate,
		Delete: resourceAlicloudDdoscooInstanceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"bandwidth": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"base_bandwidth": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"domain_count": {
				Type:     schema.TypeString,
				Required: true,
			},
			"edition": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"coop"}, false),
			},
			"function_version": {
				Type:     schema.TypeString,
				Required: true,
			},
			"instance_spec": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"bandwidth_mbps": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"instance_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"port_limit": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"qps_limit": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"site_limit": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"base_bandwidth": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"defense_count": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"domain_limit": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"elastic_bandwidth": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"function_version": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"modify_type": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"normal_qps": {
				Type:     schema.TypeString,
				Required: true,
			},
			"period": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
			},
			"port_count": {
				Type:     schema.TypeString,
				Required: true,
			},
			"remark": {
				Type:     schema.TypeString,
				Optional: true,
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
			"service_bandwidth": {
				Type:     schema.TypeString,
				Required: true,
			},
			"service_partner": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"status": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"tags": tagsSchema(),
		},
	}
}

func resourceAlicloudDdoscooInstanceCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	request := bssopenapi.CreateCreateInstanceRequest()
	if v, ok := d.GetOk("period"); ok {
		request.Period = requests.NewInteger(v.(int))
	}
	request.ProductCode = "ddos"
	request.ProductType = "ddoscoo"
	if v, ok := d.GetOk("renew_period"); ok {
		request.RenewPeriod = requests.NewInteger(v.(int))
	}
	if v, ok := d.GetOk("renewal_status"); ok {
		request.RenewalStatus = v.(string)
	}
	request.SubscriptionType = "Subscription"
	request.Parameter = &[]bssopenapi.CreateInstanceParameter{
		{
			Code:  "Bandwidth",
			Value: d.Get("bandwidth").(string),
		},
		{
			Code:  "BaseBandwidth",
			Value: d.Get("base_bandwidth").(string),
		},
		{
			Code:  "DomainCount",
			Value: d.Get("domain_count").(string),
		},
		{
			Code:  "Edition",
			Value: d.Get("edition").(string),
		},
		{
			Code:  "FunctionVersion",
			Value: d.Get("function_version").(string),
		},
		{
			Code:  "NormalQps",
			Value: d.Get("normal_qps").(string),
		},
		{
			Code:  "PortCount",
			Value: d.Get("port_count").(string),
		},
		{
			Code:  "ServiceBandwidth",
			Value: d.Get("service_bandwidth").(string),
		},
		{
			Code:  "ServicePartner",
			Value: d.Get("service_partner").(string),
		},
	}
	raw, err := client.WithBssopenapiClient(func(bssopenapiClient *bssopenapi.Client) (interface{}, error) {
		return bssopenapiClient.CreateInstance(request)
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_ddoscoo_instance", request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw)
	response, _ := raw.(*bssopenapi.CreateInstanceResponse)
	if !response.Success {
		return WrapErrorf(fmt.Errorf("%v", response), DefaultErrorMsg, "alicloud_ddoscoo_instance", request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	d.SetId(fmt.Sprintf("%v", response.Data.InstanceId))

	return resourceAlicloudDdoscooInstanceUpdate(d, meta)
}
func resourceAlicloudDdoscooInstanceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	ddoscooService := DdoscooService{client}
	object, err := ddoscooService.DescribeDdoscooInstance(d.Id())
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	base_bandwidth := strconv.Itoa(object.BaseBandwidth)
	d.Set("base_bandwidth", base_bandwidth)
	d.Set("function_version", convertFunctionVersionResponse(object.FunctionVersion))

	describeInstancesObject, err := ddoscooService.DescribeInstances(d.Id())
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	d.Set("edition", convertEditionResponse(describeInstancesObject.Edition))
	d.Set("remark", describeInstancesObject.Remark)
	d.Set("status", describeInstancesObject.Status)

	describeTagResourcesObject, err := ddoscooService.DescribeTagResources(d.Id())
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	tags := make(map[string]string)
	for _, t := range describeTagResourcesObject.TagResources.TagResource {
		tags[t.TagKey] = t.TagValue
	}
	d.Set("tags", tags)
	return nil
}
func resourceAlicloudDdoscooInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	ddoscooService := DdoscooService{client}
	d.Partial(true)

	if d.HasChange("remark") {
		request := ddoscoo.CreateModifyInstanceRemarkRequest()
		request.InstanceId = d.Id()
		request.Remark = d.Get("remark").(string)
		raw, err := client.WithDdoscooClient(func(ddoscooClient *ddoscoo.Client) (interface{}, error) {
			return ddoscooClient.ModifyInstanceRemark(request)
		})
		addDebug(request.GetActionName(), raw)
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		d.SetPartial("remark")
	}
	if d.HasChange("tags") {
		if err := ddoscooService.SetResourceTags(d, "INSTANCE"); err != nil {
			return WrapError(err)
		}
		d.SetPartial("tags")
	}
	update := false
	request := bssopenapi.CreateModifyInstanceRequest()
	request.InstanceId = d.Id()
	if d.HasChange("modify_type") {
		update = true
	}
	request.ProductType = "ddoscoo"
	request.SubscriptionType = "Subscription"
	request.ProductCode = "ddos"
	request.ModifyType = d.Get("modify_type").(string)
	request.Parameter = &[]bssopenapi.ModifyInstanceParameter{
		{
			Code:  "DomainCount",
			Value: d.Get("domain_count").(string),
		},
		{
			Code:  "Edition",
			Value: d.Get("edition").(string),
		},
		{
			Code:  "FunctionVersion",
			Value: d.Get("function_version").(string),
		},
		{
			Code:  "NormalQps",
			Value: d.Get("normal_qps").(string),
		},
		{
			Code:  "PortCount",
			Value: d.Get("port_count").(string),
		},
		{
			Code:  "ServiceBandwidth",
			Value: d.Get("service_bandwidth").(string),
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

		d.SetPartial("modify_type")
	}
	d.Partial(false)
	return resourceAlicloudDdoscooInstanceRead(d, meta)
}
func resourceAlicloudDdoscooInstanceDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	request := ddoscoo.CreateReleaseInstanceRequest()
	request.InstanceId = d.Id()
	raw, err := client.WithDdoscooClient(func(ddoscooClient *ddoscoo.Client) (interface{}, error) {
		return ddoscooClient.ReleaseInstance(request)
	})
	addDebug(request.GetActionName(), raw)
	if err != nil {
		if IsExpectedErrors(err, []string{"InstanceNotExpire"}) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	return nil
}
func convertFunctionVersionResponse(source string) string {
	switch source {
	case "default":
		return "0"
	case "enhance":
		return "1"
	}
	return ""
}
func convertEditionResponse(source int) string {
	switch source {
	case 9:
		return "coop"
	}
	return ""
}
