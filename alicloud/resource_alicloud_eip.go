package alicloud

import (
	"strconv"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/vpc"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func resourceAliyunEip() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliyunEipCreate,
		Read:   resourceAliyunEipRead,
		Update: resourceAliyunEipUpdate,
		Delete: resourceAliyunEipDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateInstanceName,
			},
			"description": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateInstanceDescription,
			},
			"bandwidth": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  5,
			},
			"internet_charge_type": {
				Type:         schema.TypeString,
				Default:      "PayByTraffic",
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validateInternetChargeType,
			},
			"instance_charge_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateInstanceChargeType,
				Default:      PostPaid,
				ForceNew:     true,
			},
			"period": {
				Type:             schema.TypeInt,
				Optional:         true,
				Default:          1,
				ForceNew:         true,
				ValidateFunc:     validateEipChargeTypePeriod,
				DiffSuppressFunc: ecsPostPaidDiffSuppressFunc,
			},
			"tags": tagsSchema(),
			"ip_address": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"isp": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"resource_group_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
		},
	}
}

func resourceAliyunEipCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	vpcService := VpcService{client}

	request := vpc.CreateAllocateEipAddressRequest()
	request.RegionId = string(client.Region)
	request.Bandwidth = strconv.Itoa(d.Get("bandwidth").(int))
	request.InternetChargeType = d.Get("internet_charge_type").(string)
	request.InstanceChargeType = d.Get("instance_charge_type").(string)
	request.ResourceGroupId = d.Get("resource_group_id").(string)
	request.ISP = d.Get("isp").(string)
	if request.InstanceChargeType == string(PrePaid) {
		period := d.Get("period").(int)
		request.Period = requests.NewInteger(period)
		request.PricingCycle = string(Month)
		if period > 9 {
			request.Period = requests.NewInteger(period / 12)
			request.PricingCycle = string(Year)
		}
		request.AutoPay = requests.NewBoolean(true)
	}
	request.ClientToken = buildClientToken(request.GetActionName())

	raw, err := client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
		return vpcClient.AllocateEipAddress(request)
	})
	if err != nil {
		if IsExceptedError(err, COMMODITYINVALID_COMPONENT) && request.InternetChargeType == string(PayByBandwidth) {
			return WrapErrorf(err, "Your account is international and it can only create '%s' elastic IP. Please change it and try again. %s", PayByTraffic, AlibabaCloudSdkGoERROR)
		}
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_eip", request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	response, _ := raw.(*vpc.AllocateEipAddressResponse)
	d.SetId(response.AllocationId)
	err = vpcService.WaitForEip(d.Id(), Available, DefaultTimeoutMedium)
	if err != nil {
		return WrapError(err)
	}
	return resourceAliyunEipUpdate(d, meta)
}

func resourceAliyunEipRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	vpcService := VpcService{client}

	object, err := vpcService.DescribeEip(d.Id())
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("name", object.Name)
	d.Set("description", object.Descritpion)
	bandwidth, _ := strconv.Atoi(object.Bandwidth)
	d.Set("bandwidth", bandwidth)
	d.Set("internet_charge_type", object.InternetChargeType)
	d.Set("instance_charge_type", object.ChargeType)
	d.Set("isp", object.ISP)
	d.Set("ip_address", object.IpAddress)
	d.Set("status", object.Status)
	d.Set("resource_group_id", object.ResourceGroupId)
	tags, err := vpcService.DescribeTags(d.Id(), nil, TagResourceEip)
	if err != nil {
		return WrapError(err)
	}
	d.Set("tags", vpcService.tagsToMap(tags))

	return nil
}

func resourceAliyunEipUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	vpcService := VpcService{client}
	if err := vpcService.setInstanceTags(d, TagResourceEip); err != nil {
		return WrapError(err)
	}

	update := false
	request := vpc.CreateModifyEipAddressAttributeRequest()
	request.RegionId = client.RegionId
	request.AllocationId = d.Id()

	if d.HasChange("bandwidth") && !d.IsNewResource() {
		update = true
		request.Bandwidth = strconv.Itoa(d.Get("bandwidth").(int))
	}
	if d.HasChange("name") {
		update = true
		request.Name = d.Get("name").(string)
	}
	if d.HasChange("description") {
		update = true
		request.Description = d.Get("description").(string)
	}
	if update {
		raw, err := client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
			return vpcClient.ModifyEipAddressAttribute(request)
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	}
	return resourceAliyunEipRead(d, meta)
}

func resourceAliyunEipDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	vpcService := VpcService{client}

	request := vpc.CreateReleaseEipAddressRequest()
	request.AllocationId = d.Id()
	request.RegionId = client.RegionId
	err := resource.Retry(5*time.Minute, func() *resource.RetryError {
		raw, err := client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
			return vpcClient.ReleaseEipAddress(request)
		})
		if err != nil {
			if IsExceptedError(err, EipIncorrectStatus) {
				return resource.RetryableError(err)
			} else if IsExceptedError(err, AllocationIdNotFound) {
				return nil
			}
			return resource.NonRetryableError(err)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		return nil
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	return WrapError(vpcService.WaitForEip(d.Id(), Deleted, DefaultTimeout))
}
