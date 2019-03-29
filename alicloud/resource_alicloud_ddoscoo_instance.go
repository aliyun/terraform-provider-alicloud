package alicloud

import (
	"strconv"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/bssopenapi"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ddoscoo"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func resourceAlicloudDdoscoo() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudDdoscooCreate,
		Read:   resourceAlicloudDdoscooRead,
		Update: resourceAlicloudDdoscooUpdate,
		Delete: resourceAlicloudDdoscooDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"id": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateDdoscooInstanceName,
			},
			"base_bandwidth": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"bandwidth": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"service_bandwidth": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"port_count": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"domain_count": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"period": {
				Type:         schema.TypeInt,
				ValidateFunc: validateAllowedIntValue([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 12, 24, 36}),
				Optional:     true,
				Default:      1,
				ForceNew:     true,
			},
		},
	}
}

func resourceAlicloudDdoscooCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	request := buildDdoscooCreateRequest(d, meta)

	raw, err := client.WithBssopenapiClient(func(bssopenapiClient *bssopenapi.Client) (interface{}, error) {
		return bssopenapiClient.CreateInstance(request)
	})

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_ddoscoo_instance", request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw)

	resp := raw.(*bssopenapi.CreateInstanceResponse)
	// execute errors including in the bssopenapi response
	if !resp.Success {
		return Error(resp.Message)
	}

	d.SetId(resp.Data.InstanceId)

	return resourceAlicloudDdoscooUpdate(d, meta)
}

func resourceAlicloudDdoscooRead(d *schema.ResourceData, meta interface{}) error {
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

	d.Set("id", insInfo.Instances[0].InstanceId)
	d.Set("name", insInfo.Instances[0].Remark)
	d.Set("bandwidth", specInfo.InstanceSpecs[0].ElasticBandwidth)
	d.Set("base_bandwidth", specInfo.InstanceSpecs[0].BaseBandwidth)
	d.Set("domain_count", specInfo.InstanceSpecs[0].DomainLimit)
	d.Set("port_count", specInfo.InstanceSpecs[0].PortLimit)
	d.Set("service_bandwidth", specInfo.InstanceSpecs[0].BandwidthMbps)
	d.Set("period", d.Get("period"))

	return nil
}

func resourceAlicloudDdoscooUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	ddoscooService := DdoscooService{client}

	d.Partial(true)
	if d.HasChange("base_bandwidth") || d.HasChange("bandwidth") {
		if err := ddoscooService.UpdateBandwidth(d, meta); err != nil {
			return WrapError(err)
		}

		d.SetPartial("base_bandwidth")
		d.SetPartial("bandwidth")
	}

	if d.HasChange("domain_count") {
		o, n := d.GetChange("domain_count")
		oi, _ := strconv.Atoi(o.(string))
		ni, _ := strconv.Atoi(n.(string))
		if ni < oi {
			if err := ddoscooService.DowngradeDomainCount(d, meta); err != nil {
				return WrapError(err)
			}
		}

		if ni > oi {
			if err := ddoscooService.UpgradeDomainCount(d, meta); err != nil {
				return WrapError(err)
			}
		}

		d.SetPartial("domain_count")
	}

	if d.HasChange("port_count") {
		o, n := d.GetChange("port_count")
		oi, _ := strconv.Atoi(o.(string))
		ni, _ := strconv.Atoi(n.(string))
		if ni < oi {
			if err := ddoscooService.DowngradePortCount(d, meta); err != nil {
				return WrapError(err)
			}
		}

		if ni > oi {
			if err := ddoscooService.UpgradePortCount(d, meta); err != nil {
				return WrapError(err)
			}
		}

		d.SetPartial("port_count")
	}

	if d.HasChange("service_bandwidth") {
		o, n := d.GetChange("service_bandwidth")
		oi, _ := strconv.Atoi(o.(string))
		ni, _ := strconv.Atoi(n.(string))
		if ni < oi {
			if err := ddoscooService.DowngradeServiceBandwidth(d, meta); err != nil {
				return WrapError(err)
			}
		}

		if ni > oi {
			if err := ddoscooService.UpgradeServiceBandwidth(d, meta); err != nil {
				return WrapError(err)
			}
		}

		d.SetPartial("service_bandwidth")
	}

	if d.HasChange("name") {
		if err := ddoscooService.UpdateDdoscooInstanceName(d.Id(), d.Get("name").(string)); err != nil {
			return WrapError(err)
		}
		d.SetPartial("name")
	}

	d.Partial(false)
	return resourceAlicloudDdoscooRead(d, meta)
}

func resourceAlicloudDdoscooDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	request := ddoscoo.CreateReleaseInstanceRequest()
	request.InstanceId = d.Id()

	_, err := client.WithDdoscooClient(func(ddoscooClient *ddoscoo.Client) (interface{}, error) {
		return ddoscooClient.ReleaseInstance(request)
	})
	if err != nil {
		if IsExceptedError(err, DdoscooInstanceNotFound) {
			return nil
		}

		return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
	}

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
