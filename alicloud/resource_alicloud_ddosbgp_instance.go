package alicloud

import (
	"strconv"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/bssopenapi"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ddosbgp"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func resourceAlicloudDdosbgpInstance() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudDdosbgpInstanceCreate,
		Read:   resourceAlicloudDdosbgpInstanceRead,
		Update: resourceAlicloudDdosbgpInstanceUpdate,
		Delete: resourceAlicloudDdosbgpInstanceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"type": {
				Type:         schema.TypeString,
				Optional:     true,
				Required:     false,
				Default:      string(Enterprise),
				ForceNew:     true,
				ValidateFunc: validateAllowedStringValue([]string{string(Enterprise), string(Professional)}),
			},
			"name": {
				Type:         schema.TypeString,
				Optional:     true,
				Required:     false,
				ValidateFunc: validateDdosbgpInstanceName,
			},
			"base_bandwidth": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"bandwidth": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"ip_type": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validateAllowedStringValue([]string{string(IPv4), string(IPv6)}),
			},
			"ip_count": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"period": {
				Type:         schema.TypeInt,
				ValidateFunc: validateAllowedIntValue([]int{1, 2, 3}),
				Optional:     true,
				Default:      1,
				ForceNew:     true,
			},
		},
	}
}

func resourceAlicloudDdosbgpInstanceCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	request := buildDdosbgpCreateRequest(client.RegionId, d, meta)

	raw, err := client.WithBssopenapiClient(func(bssopenapiClient *bssopenapi.Client) (interface{}, error) {
		return bssopenapiClient.CreateInstance(request)
	})

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_ddosbgp_instance", request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)

	resp := raw.(*bssopenapi.CreateInstanceResponse)

	d.SetId(resp.Data.InstanceId)

	return resourceAlicloudDdosbgpInstanceUpdate(d, meta)
}

func resourceAlicloudDdosbgpInstanceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	ddosbgpService := DdosbgpService{client}
	insInfo, err := ddosbgpService.DescribeDdosbgpInstance(d.Id())
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}

		return WrapError(err)
	}

	specInfo, err := ddosbgpService.DescribeDdosbgpInstanceSpec(d.Id(), client.RegionId)
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}

		return WrapError(err)
	}

	ddosbgpInstanceType := string(Enterprise)
	if insInfo.InstanceType == "0" {
		ddosbgpInstanceType = string(Professional)
	}

	d.Set("name", insInfo.Remark)
	d.Set("region", specInfo.Region)
	d.Set("bandwidth", strconv.Itoa(specInfo.PackConfig.PackAdvThre))
	d.Set("base_bandwidth", strconv.Itoa(specInfo.PackConfig.PackBasicThre))
	d.Set("ip_type", insInfo.IpType)
	d.Set("ip_count", strconv.Itoa(specInfo.PackConfig.IpSpec))
	d.Set("type", ddosbgpInstanceType)

	return nil
}

func resourceAlicloudDdosbgpInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	if d.HasChange("name") {
		request := ddosbgp.CreateModifyRemarkRequest()
		request.InstanceId = d.Id()
		request.RegionId = client.RegionId
		request.ResourceRegionId = client.RegionId

		request.Remark = d.Get("name").(string)

		raw, err := client.WithDdosbgpClient(func(ddosbgpClient *ddosbgp.Client) (interface{}, error) {
			return ddosbgpClient.ModifyRemark(request)
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	}

	return resourceAlicloudDdosbgpInstanceRead(d, meta)
}

func resourceAlicloudDdosbgpInstanceDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	ddosbgpService := DdosbgpService{client}

	request := ddosbgp.CreateReleaseInstanceRequest()
	request.InstanceId = d.Id()
	request.RegionId = client.RegionId

	raw, err := client.WithDdosbgpClient(func(ddosbgpClient *ddosbgp.Client) (interface{}, error) {
		return ddosbgpClient.ReleaseInstance(request)
	})
	if err != nil {
		if IsExceptedError(err, DdosbgpInstanceNotFound) {
			return nil
		}

		return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)

	return WrapError(ddosbgpService.WaitForDdosbgpInstance(d.Id(), Deleted, DefaultTimeoutMedium))
}

func buildDdosbgpCreateRequest(region string, d *schema.ResourceData, meta interface{}) *bssopenapi.CreateInstanceRequest {
	request := bssopenapi.CreateCreateInstanceRequest()
	request.ProductCode = "ddos"
	request.ProductType = "ddosbgp"
	request.SubscriptionType = "Subscription"
	request.Period = requests.NewInteger(d.Get("period").(int))

	ddosbgpInstanceType := "1"
	if d.Get("type").(string) == string(Professional) {
		ddosbgpInstanceType = "0"
	}

	ddosbgpInstanceIpType := "v4"
	if d.Get("ip_type").(string) == string(IPv6) {
		ddosbgpInstanceIpType = "v6"
	}

	request.Parameter = &[]bssopenapi.CreateInstanceParameter{
		{
			Code:  "Type",
			Value: ddosbgpInstanceType,
		},
		{
			Code:  "Region",
			Value: region,
		},
		{
			Code:  "IpType",
			Value: ddosbgpInstanceIpType,
		},
		{
			Code:  "BaseBandwidth",
			Value: d.Get("base_bandwidth").(string),
		},
		{
			Code:  "Bandwidth",
			Value: d.Get("bandwidth").(string),
		},
		{
			Code:  "IpCount",
			Value: d.Get("ip_count").(string),
		},
	}

	return request
}
