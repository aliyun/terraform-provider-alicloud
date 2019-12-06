package alicloud

import (
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/bssopenapi"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ddoscoo"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
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
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(1, 64),
			},
			"base_bandwidth": {
				Type:     schema.TypeString,
				Required: true,
			},
			"bandwidth": {
				Type:     schema.TypeString,
				Required: true,
			},
			"service_bandwidth": {
				Type:     schema.TypeString,
				Required: true,
			},
			"port_count": {
				Type:     schema.TypeString,
				Required: true,
			},
			"domain_count": {
				Type:     schema.TypeString,
				Required: true,
			},
			"period": {
				Type:         schema.TypeInt,
				ValidateFunc: validation.IntInSlice([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 12, 24, 36}),
				Optional:     true,
				Default:      1,
				ForceNew:     true,
			},
		},
	}
}

func resourceAlicloudDdoscooInstanceCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	request := buildDdoscooCreateRequest(d, meta)
	request.RegionId = client.RegionId

	raw, err := client.WithBssopenapiClient(func(bssopenapiClient *bssopenapi.Client) (interface{}, error) {
		return bssopenapiClient.CreateInstance(request)
	})

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_ddoscoo_instance", request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)

	response := raw.(*bssopenapi.CreateInstanceResponse)
	// execute errors including in the bssopenapi response
	if !response.Success {
		return WrapError(Error(response.Message))
	}

	d.SetId(response.Data.InstanceId)

	return resourceAlicloudDdoscooInstanceUpdate(d, meta)
}

func resourceAlicloudDdoscooInstanceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	ddoscooService := DdoscooService{client}

	insInfo, err := ddoscooService.DescribeDdoscooInstance(d.Id())
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}

		return WrapError(err)
	}

	specInfo, err := ddoscooService.DescribeDdoscooInstanceSpec(d.Id())
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}

		return WrapError(err)
	}

	d.Set("name", insInfo.Remark)
	d.Set("bandwidth", strconv.Itoa(specInfo.ElasticBandwidth))
	d.Set("base_bandwidth", strconv.Itoa(specInfo.BaseBandwidth))
	d.Set("domain_count", strconv.Itoa(specInfo.DomainLimit))
	d.Set("port_count", strconv.Itoa(specInfo.PortLimit))
	d.Set("service_bandwidth", strconv.Itoa(specInfo.BandwidthMbps))

	return nil
}

func resourceAlicloudDdoscooInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	ddoscooService := DdoscooService{client}

	d.Partial(true)

	if d.HasChange("name") {
		if err := ddoscooService.UpdateDdoscooInstanceName(d.Id(), d.Get("name").(string)); err != nil {
			return WrapError(err)
		}
		d.SetPartial("name")
	}

	if d.IsNewResource() {
		d.Partial(false)
		return resourceAlicloudDdoscooInstanceRead(d, meta)
	}

	if d.HasChange("bandwidth") {
		if err := ddoscooService.UpdateInstanceSpec("bandwidth", "Bandwidth", d, meta); err != nil {
			return WrapError(err)
		}

		d.SetPartial("bandwidth")
	}

	if d.HasChange("base_bandwidth") {
		if err := ddoscooService.UpdateInstanceSpec("base_bandwidth", "BaseBandwidth", d, meta); err != nil {
			return WrapError(err)
		}

		d.SetPartial("base_bandwidth")
	}

	if d.HasChange("domain_count") {
		if err := ddoscooService.UpdateInstanceSpec("domain_count", "DomainCount", d, meta); err != nil {
			return WrapError(err)
		}

		d.SetPartial("domain_count")
	}

	if d.HasChange("port_count") {
		if err := ddoscooService.UpdateInstanceSpec("port_count", "PortCount", d, meta); err != nil {
			return WrapError(err)
		}

		d.SetPartial("port_count")
	}

	if d.HasChange("service_bandwidth") {
		if err := ddoscooService.UpdateInstanceSpec("service_bandwidth", "ServiceBandwidth", d, meta); err != nil {
			return WrapError(err)
		}

		d.SetPartial("service_bandwidth")
	}

	d.Partial(false)
	return resourceAlicloudDdoscooInstanceRead(d, meta)
}

func resourceAlicloudDdoscooInstanceDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	request := ddoscoo.CreateReleaseInstanceRequest()
	request.RegionId = client.RegionId
	request.InstanceId = d.Id()

	raw, err := client.WithDdoscooClient(func(ddoscooClient *ddoscoo.Client) (interface{}, error) {
		return ddoscooClient.ReleaseInstance(request)
	})
	if err != nil {
		if IsExceptedError(err, DdoscooInstanceNotFound) {
			return nil
		}

		return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	return nil
}

func buildDdoscooCreateRequest(d *schema.ResourceData, meta interface{}) *bssopenapi.CreateInstanceRequest {
	request := bssopenapi.CreateCreateInstanceRequest()
	request.ProductCode = "ddos"
	request.ProductType = "ddoscoo"
	request.SubscriptionType = "Subscription"
	request.Period = requests.NewInteger(d.Get("period").(int))

	request.Parameter = &[]bssopenapi.CreateInstanceParameter{
		{
			Code:  "ServicePartner",
			Value: "coop-line-001",
		},
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
			Code:  "PortCount",
			Value: d.Get("port_count").(string),
		},
		{
			Code:  "ServiceBandwidth",
			Value: d.Get("service_bandwidth").(string),
		},
	}

	return request
}
