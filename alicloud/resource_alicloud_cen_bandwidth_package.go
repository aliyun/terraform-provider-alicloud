package alicloud

import (
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/cbn"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func resourceAlicloudCenBandwidthPackage() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudCenBandwidthPackageCreate,
		Read:   resourceAlicloudCenBandwidthPackageRead,
		Update: resourceAlicloudCenBandwidthPackageUpdate,
		Delete: resourceAlicloudCenBandwidthPackageDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"bandwidth": {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: validation.IntAtLeast(1),
			},

			"geographic_region_ids": {
				Type:     schema.TypeSet,
				Required: true,
				ForceNew: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				MaxItems: 2,
				MinItems: 1,
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

			"charge_type": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      PostPaid,
				ValidateFunc: validation.StringInSlice([]string{string(PrePaid), string(PostPaid)}, false),
			},

			"period": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      1,
				ValidateFunc: validation.IntInSlice([]int{1, 2, 3, 6, 12}),
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return PayType(d.Get("charge_type").(string)) == PostPaid
				},
			},

			"expired_time": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceAlicloudCenBandwidthPackageCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	cenService := CenService{client}
	var response *cbn.CreateCenBandwidthPackageResponse
	request := buildAliCloudCenBandwidthPackageArgs(d, meta)
	request.RegionId = client.RegionId
	bandwidth, _ := strconv.Atoi(string(request.Bandwidth))

	req := *request
	err := resource.Retry(3*time.Minute, func() *resource.RetryError {
		raw, err := client.WithCenClient(func(cbnClient *cbn.Client) (interface{}, error) {
			return cbnClient.CreateCenBandwidthPackage(&req)
		})
		if err != nil {
			if IsExceptedError(err, OperationBlocking) {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		response = raw.(*cbn.CreateCenBandwidthPackageResponse)
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		return nil
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
	}

	d.SetId(response.CenBandwidthPackageId)
	err = cenService.WaitForCenBandwidthPackage(d.Id(), Idle, bandwidth, DefaultCenTimeout)
	if err != nil {
		return WrapError(err)
	}

	return resourceAlicloudCenBandwidthPackageRead(d, meta)
}

func resourceAlicloudCenBandwidthPackageRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	cenService := CenService{client}
	object, err := cenService.DescribeCenBandwidthPackage(d.Id())
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	geographicRegionId := make([]string, 0)
	geographicRegionId = append(geographicRegionId, convertGeographicRegionId(object.GeographicRegionAId))
	geographicRegionId = append(geographicRegionId, convertGeographicRegionId(object.GeographicRegionBId))

	d.Set("geographic_region_ids", geographicRegionId)
	d.Set("bandwidth", object.Bandwidth)
	d.Set("name", object.Name)
	d.Set("description", object.Description)
	d.Set("expired_time", object.ExpiredTime)
	d.Set("status", object.Status)

	if object.BandwidthPackageChargeType == "POSTPAY" {
		d.Set("charge_type", PostPaid)
	} else {
		d.Set("charge_type", PrePaid)
		period, err := computePeriodByUnit(object.CreationTime, object.ExpiredTime, "Month")
		if err != nil {
			return WrapError(err)
		}
		d.Set("period", period)
	}

	return nil
}

func resourceAlicloudCenBandwidthPackageUpdate(d *schema.ResourceData, meta interface{}) error {
	d.Partial(true)

	client := meta.(*connectivity.AliyunClient)
	cenService := CenService{client}
	attributeUpdate := false
	modifyCenBandwidthPackageAttributeRequest := cbn.CreateModifyCenBandwidthPackageAttributeRequest()
	modifyCenBandwidthPackageAttributeRequest.CenBandwidthPackageId = d.Id()

	if d.HasChange("name") {
		modifyCenBandwidthPackageAttributeRequest.Name = d.Get("name").(string)
		attributeUpdate = true
	}

	if d.HasChange("description") {
		modifyCenBandwidthPackageAttributeRequest.Description = d.Get("description").(string)
		attributeUpdate = true
	}

	if attributeUpdate {
		raw, err := client.WithCenClient(func(cbnClient *cbn.Client) (interface{}, error) {
			return cbnClient.ModifyCenBandwidthPackageAttribute(modifyCenBandwidthPackageAttributeRequest)
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), modifyCenBandwidthPackageAttributeRequest.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		addDebug(modifyCenBandwidthPackageAttributeRequest.GetActionName(), raw)
		d.SetPartial("name")
		d.SetPartial("description")
	}

	if d.HasChange("bandwidth") {
		bandwidth := d.Get("bandwidth").(int)
		modifyCenBandwidthPackageSpecRequest := cbn.CreateModifyCenBandwidthPackageSpecRequest()
		modifyCenBandwidthPackageSpecRequest.CenBandwidthPackageId = d.Id()
		modifyCenBandwidthPackageSpecRequest.Bandwidth = requests.NewInteger(bandwidth)

		raw, err := client.WithCenClient(func(cbnClient *cbn.Client) (interface{}, error) {
			return cbnClient.ModifyCenBandwidthPackageSpec(modifyCenBandwidthPackageSpecRequest)
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), modifyCenBandwidthPackageSpecRequest.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		addDebug(modifyCenBandwidthPackageSpecRequest.GetActionName(), raw)
		// modify function may delay for a while
		if err := cenService.WaitForCenBandwidthPackage(d.Id(), Idle, bandwidth, DefaultCenTimeout); err != nil {
			return WrapError(err)
		}
		d.SetPartial("bandwidth")
	}

	d.Partial(false)

	return resourceAlicloudCenBandwidthPackageRead(d, meta)
}

func resourceAlicloudCenBandwidthPackageDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	cenService := CenService{client}
	cenBwpId := d.Id()
	request := cbn.CreateDeleteCenBandwidthPackageRequest()
	request.CenBandwidthPackageId = cenBwpId

	err := resource.Retry(5*time.Minute, func() *resource.RetryError {
		raw, err := client.WithCenClient(func(cbnClient *cbn.Client) (interface{}, error) {
			return cbnClient.DeleteCenBandwidthPackage(request)
		})
		if err != nil {
			if IsExceptedErrors(err, []string{ForbiddenRelease, InvalidCenBandwidthLimitsNotZero, ParameterBwpInstanceId}) {
				return resource.NonRetryableError(err)
			}
			return resource.RetryableError(err)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		return nil
	})
	if err != nil {
		if IsExceptedError(err, ParameterBwpInstanceId) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
	}

	// set bandwidth "-1" here to use WaitForCenBandwidthPackage, actually determined by status.
	return WrapError(cenService.WaitForCenBandwidthPackage(d.Id(), Deleted, -1, DefaultCenTimeout))
}

func convertGeographicRegionId(regionId string) (retStr string) {
	switch regionId {
	case "china":
		retStr = "China"
	case "north-america":
		retStr = "North-America"
	case "asia-pacific":
		retStr = "Asia-Pacific"
	case "europe":
		retStr = "Europe"
	case "middle-east":
		retStr = "Middle-East"
	}

	return
}

func buildAliCloudCenBandwidthPackageArgs(d *schema.ResourceData, meta interface{}) *cbn.CreateCenBandwidthPackageRequest {
	request := cbn.CreateCreateCenBandwidthPackageRequest()
	request.Bandwidth = requests.NewInteger(d.Get("bandwidth").(int))

	geographicRegionId := d.Get("geographic_region_ids").(*schema.Set).List()
	if len(geographicRegionId) == 1 {
		request.GeographicRegionAId = geographicRegionId[0].(string)
		request.GeographicRegionBId = geographicRegionId[0].(string)
	} else if len(geographicRegionId) == 2 {
		if geographicRegionId[1].(string) == "China" {
			request.GeographicRegionAId = geographicRegionId[1].(string)
			request.GeographicRegionBId = geographicRegionId[0].(string)
		} else {
			request.GeographicRegionAId = geographicRegionId[0].(string)
			request.GeographicRegionBId = geographicRegionId[1].(string)
		}
	}

	if v := d.Get("name").(string); v != "" {
		request.Name = v
	}

	if v := d.Get("description").(string); v != "" {
		request.Description = v
	}

	changeType := d.Get("charge_type").(string)
	if changeType == string(PostPaid) {
		changeType = "POSTPAY"
	} else {
		changeType = "PREPAY"
		request.Period = requests.NewInteger(d.Get("period").(int))
		request.PricingCycle = "Month"
	}

	request.BandwidthPackageChargeType = changeType
	request.AutoPay = requests.NewBoolean(true)
	request.ClientToken = buildClientToken(request.GetActionName())

	return request
}
