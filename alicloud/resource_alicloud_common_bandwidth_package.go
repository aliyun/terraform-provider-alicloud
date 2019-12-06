package alicloud

import (
	"strconv"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/vpc"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func resourceAliyunCommonBandwidthPackage() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliyunCommonBandwidthPackageCreate,
		Read:   resourceAliyunCommonBandwidthPackageRead,
		Update: resourceAliyunCommonBandwidthPackageUpdate,
		Delete: resourceAliyunCommonBandwidthPackageDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"description": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringLenBetween(2, 256),
			},
			"name": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringLenBetween(2, 128),
			},
			"bandwidth": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"internet_charge_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Default:      PayByTraffic,
				ValidateFunc: validation.StringInSlice([]string{"PayByBandwidth", "PayBy95", "PayByTraffic"}, false),
			},
			"ratio": {
				Type:         schema.TypeInt,
				Optional:     true,
				ForceNew:     true,
				Default:      100,
				ValidateFunc: validation.IntBetween(10, 100),
			},
			"resource_group_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
		},
	}
}

func resourceAliyunCommonBandwidthPackageCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	vpcService := VpcService{client}

	request := vpc.CreateCreateCommonBandwidthPackageRequest()
	request.RegionId = client.RegionId

	request.Bandwidth = requests.NewInteger(d.Get("bandwidth").(int))
	request.Name = d.Get("name").(string)
	request.Description = d.Get("description").(string)
	request.ResourceGroupId = d.Get("resource_group_id").(string)
	request.InternetChargeType = d.Get("internet_charge_type").(string)
	request.Ratio = requests.NewInteger(d.Get("ratio").(int))
	request.ClientToken = buildClientToken(request.GetActionName())
	raw, err := client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
		return vpcClient.CreateCommonBandwidthPackage(request)
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_common_bandwidth_package", request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	response, _ := raw.(*vpc.CreateCommonBandwidthPackageResponse)
	d.SetId(response.BandwidthPackageId)
	if err = vpcService.WaitForCommonBandwidthPackage(response.BandwidthPackageId, Available, DefaultTimeout); err != nil {
		return WrapError(err)
	}

	return resourceAliyunCommonBandwidthPackageRead(d, meta)
}

func resourceAliyunCommonBandwidthPackageRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	vpcService := VpcService{client}
	object, err := vpcService.DescribeCommonBandwidthPackage(d.Id())
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	bandwidth, err := strconv.Atoi(object.Bandwidth)
	if err != nil {
		return WrapError(err)
	}
	d.Set("bandwidth", bandwidth)
	d.Set("name", object.Name)
	d.Set("description", object.Description)
	d.Set("internet_charge_type", object.InternetChargeType)
	d.Set("ratio", object.Ratio)
	d.Set("resource_group_id", object.ResourceGroupId)
	return nil
}

func resourceAliyunCommonBandwidthPackageUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	d.Partial(true)
	update := false
	request := vpc.CreateModifyCommonBandwidthPackageAttributeRequest()
	request.RegionId = client.RegionId
	request.BandwidthPackageId = d.Id()
	if d.HasChange("description") {
		request.Description = d.Get("description").(string)
		update = true
	}

	if d.HasChange("name") {
		request.Name = d.Get("name").(string)
		update = true
	}

	if update {
		raw, err := client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
			return vpcClient.ModifyCommonBandwidthPackageAttribute(request)
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		d.SetPartial("description")
		d.SetPartial("name")
	}

	if d.HasChange("bandwidth") {
		request := vpc.CreateModifyCommonBandwidthPackageSpecRequest()
		request.RegionId = string(client.Region)
		request.BandwidthPackageId = d.Id()
		request.Bandwidth = strconv.Itoa(d.Get("bandwidth").(int))
		raw, err := client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
			return vpcClient.ModifyCommonBandwidthPackageSpec(request)
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		d.SetPartial("bandwidth")
	}

	d.Partial(false)
	return resourceAliyunCommonBandwidthPackageRead(d, meta)
}

func resourceAliyunCommonBandwidthPackageDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	vpcService := VpcService{client}

	request := vpc.CreateDeleteCommonBandwidthPackageRequest()
	request.RegionId = client.RegionId
	request.BandwidthPackageId = d.Id()
	raw, err := client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
		return vpcClient.DeleteCommonBandwidthPackage(request)
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)

	return WrapError(vpcService.WaitForCommonBandwidthPackage(d.Id(), Deleted, DefaultTimeoutMedium))
}
