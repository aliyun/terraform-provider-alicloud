package alicloud

import (
	"fmt"
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
			// Basic instance information
			"business_endpoint": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"id": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateDdoscooInstanceName,
				ForceNew:     true,
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
			},
		},
	}
}

func resourceAlicloudDdoscooCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	request := buildDdoscooCreateRequest(d, meta)

	raw, err := client.WithBssopenapiClient(d.Get("business_endpoint").(string), func(bssopenapiClient *bssopenapi.Client) (interface{}, error) {
		return bssopenapiClient.CreateInstance(request)
	})

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_ddoscoo_instance", request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw)

	resp, _ := raw.(*bssopenapi.CreateInstanceResponse)
	d.SetId(resp.Data.InstanceId)

	ddoscooService := DdoscooService{client}
	err = ddoscooService.UpdateDdoscooInstanceName(d.Id(), d.Get("name").(string))
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_ddoscoo_instance", request.GetActionName(), AlibabaCloudSdkGoERROR)
	}

	return resourceAlicloudDdoscooRead(d, meta)
}

func resourceAlicloudDdoscooRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	ddoscooService := DdoscooService{client}

	resp, err := ddoscooService.DescribeDdoscooInstanceSpec(d.Id())
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}

		return WrapError(err)
	}

	d.Set("bandwidth", resp.InstanceSpecs[0].ElasticBandwidth)
	d.Set("base_bandwidth", resp.InstanceSpecs[0].BaseBandwidth)
	d.Set("domain_count", resp.InstanceSpecs[0].DomainLimit)
	d.Set("port_count", resp.InstanceSpecs[0].PortLimit)
	d.Set("service_bandwidth", resp.InstanceSpecs[0].BandwidthMbps)

	return nil
}

func resourceAlicloudDdoscooUpdate(d *schema.ResourceData, meta interface{}) error {
	d.Partial(true)
	if d.HasChange("base_bandwidth") || d.HasChange("bandwidth") {
		obaseBandwidth, nbaseBandwidth := d.GetChange("base_bandwidth")
		intOldbaseBandwidth, _ := strconv.Atoi(obaseBandwidth.(string))
		intNewbaseBandwidth, _ := strconv.Atoi(nbaseBandwidth.(string))

		oBandwidth, nBandwidth := d.GetChange("bandwidth")
		intOldBandwidth, _ := strconv.Atoi(oBandwidth.(string))
		intNewBandwidth, _ := strconv.Atoi(nBandwidth.(string))

		if intNewbaseBandwidth < intOldbaseBandwidth {
			return fmt.Errorf("The base bandwidth must greater than the current. The instance's current base bandwidth is %d.", intOldbaseBandwidth)
		}

		if intNewBandwidth < intOldBandwidth {
			return fmt.Errorf("The bandwidth must greater than the current. The instance's current bandwidth is %d.", intOldBandwidth)
		}

		if intNewbaseBandwidth > intOldbaseBandwidth || intNewBandwidth > intOldBandwidth {
			if err := UpdateBandwidth(d, meta); err != nil {
				return WrapError(err)
			}
		}

		d.SetPartial("base_bandwidth")
		d.SetPartial("bandwidth")
	}

	if d.HasChange("domain_count") {
		o, n := d.GetChange("domain_count")
		oi, _ := strconv.Atoi(o.(string))
		ni, _ := strconv.Atoi(n.(string))
		if ni < oi {
			if err := DowngradeDomainCount(d, meta); err != nil {
				return WrapError(err)
			}
		}

		if ni > oi {
			if err := UpgradeDomainCount(d, meta); err != nil {
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
			if err := DowngradePortCount(d, meta); err != nil {
				return WrapError(err)
			}
		}

		if ni > oi {
			if err := UpgradePortCount(d, meta); err != nil {
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
			if err := DowngradeServiceBandwidth(d, meta); err != nil {
				return WrapError(err)
			}
		}

		if ni > oi {
			if err := UpgradeServiceBandwidth(d, meta); err != nil {
				return WrapError(err)
			}
		}

		d.SetPartial("service_bandwidth")
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
			return WrapErrorf(err, DdoscooInstanceNotFound)
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
