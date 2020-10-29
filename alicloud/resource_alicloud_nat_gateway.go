package alicloud

import (
	"strings"
	"time"

	util "github.com/alibabacloud-go/tea-utils/service"

	"github.com/denverdino/aliyungo/common"

	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"

	"strconv"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/vpc"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAliyunNatGateway() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliyunNatGatewayCreate,
		Read:   resourceAliyunNatGatewayRead,
		Update: resourceAliyunNatGatewayUpdate,
		Delete: resourceAliyunNatGatewayDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"vpc_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"spec": {
				Type:       schema.TypeString,
				Optional:   true,
				Deprecated: "Field 'spec' has been deprecated from provider version 1.7.1, and new field 'specification' can replace it.",
			},
			"specification": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"Small", "Middle", "Large"}, false),
				Default:      "Small",
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"nat_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"Normal", "Enhanced"}, false),
				Default:      "Normal",
			},
			"vswitch_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"bandwidth_package_ids": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"snat_table_ids": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"forward_table_ids": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"bandwidth_packages": {
				Type: schema.TypeList,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ip_count": {
							Type:     schema.TypeInt,
							Required: true,
						},
						"bandwidth": {
							Type:     schema.TypeInt,
							Required: true,
						},
						"zone": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"public_ip_addresses": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
				MaxItems: 4,
				Optional: true,
			},
			"instance_charge_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Computed:     true,
				ValidateFunc: validation.StringInSlice([]string{string(common.PrePaid), string(common.PostPaid)}, false),
			},

			"period": {
				Type:             schema.TypeInt,
				Optional:         true,
				ForceNew:         true,
				Default:          1,
				DiffSuppressFunc: PostPaidDiffSuppressFunc,
				ValidateFunc: validation.Any(
					validation.IntBetween(1, 9),
					validation.IntInSlice([]int{12, 24, 36})),
			},
		},
	}
}

func resourceAliyunNatGatewayCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	vpcService := VpcService{client}

	request := vpc.CreateCreateNatGatewayRequest()
	request.RegionId = string(client.Region)
	request.VpcId = string(d.Get("vpc_id").(string))
	request.Spec = string(d.Get("specification").(string))
	request.NatType = d.Get("nat_type").(string)
	request.InstanceChargeType = d.Get("instance_charge_type").(string)
	if request.InstanceChargeType == string(PrePaid) {
		period := d.Get("period").(int)
		request.Duration = strconv.Itoa(period)
		request.PricingCycle = string(Month)
		if period > 9 {
			request.Duration = strconv.Itoa(period / 12)
			request.PricingCycle = string(Year)
		}
		request.AutoPay = requests.NewBoolean(true)
	}
	request.ClientToken = buildClientToken(request.GetActionName())
	bandwidthPackages := []vpc.CreateNatGatewayBandwidthPackage{}
	for _, e := range d.Get("bandwidth_packages").([]interface{}) {
		pack := e.(map[string]interface{})
		bandwidthPackage := vpc.CreateNatGatewayBandwidthPackage{
			IpCount:   strconv.Itoa(pack["ip_count"].(int)),
			Bandwidth: strconv.Itoa(pack["bandwidth"].(int)),
		}
		if pack["zone"].(string) != "" {
			bandwidthPackage.Zone = pack["zone"].(string)
		}
		bandwidthPackages = append(bandwidthPackages, bandwidthPackage)
	}

	request.BandwidthPackage = &bandwidthPackages

	if v, ok := d.GetOk("name"); ok {
		request.Name = v.(string)
	}

	if v, ok := d.GetOk("description"); ok {
		request.Description = v.(string)
	}

	if v, ok := d.GetOk("vswitch_id"); ok {
		request.VSwitchId = v.(string)
	}

	if err := resource.Retry(3*time.Minute, func() *resource.RetryError {
		args := *request
		raw, err := client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
			return vpcClient.CreateNatGateway(&args)
		})
		if err != nil {
			if IsExpectedErrors(err, []string{"VswitchStatusError", "TaskConflict"}) {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(args.GetActionName(), raw, args.RpcRequest, args)
		response, _ := raw.(*vpc.CreateNatGatewayResponse)
		d.SetId(response.NatGatewayId)
		return nil
	}); err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_nat_gateway", request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	if err := vpcService.WaitForNatGateway(d.Id(), Available, DefaultTimeout*3); err != nil {
		return WrapError(err)
	}
	return resourceAliyunNatGatewayRead(d, meta)
}

func resourceAliyunNatGatewayRead(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	vpcService := VpcService{client}

	object, err := vpcService.DescribeNatGateway(d.Id())
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("name", object.Name)
	d.Set("specification", object.Spec)
	d.Set("bandwidth_package_ids", strings.Join(object.BandwidthPackageIds.BandwidthPackageId, ","))
	d.Set("snat_table_ids", strings.Join(object.SnatTableIds.SnatTableId, ","))
	d.Set("forward_table_ids", strings.Join(object.ForwardTableIds.ForwardTableId, ","))
	d.Set("description", object.Description)
	d.Set("vpc_id", object.VpcId)
	d.Set("nat_type", object.NatType)
	d.Set("instance_charge_type", object.InstanceChargeType)
	d.Set("vswitch_id", object.NatGatewayPrivateInfo.VswitchId)
	if object.InstanceChargeType == "PrePaid" {
		period, err := computePeriodByUnit(object.CreationTime, object.ExpiredTime, d.Get("period").(int), "Month")
		if err != nil {
			return WrapError(err)
		}
		d.Set("period", period)
	}

	bindWidthPackages, err := flattenBandWidthPackages(d.Id(), meta)
	if err != nil {
		return WrapError(err)
	} else {
		d.Set("bandwidth_packages", bindWidthPackages)
	}

	return nil
}

func resourceAliyunNatGatewayUpdate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	vpcService := VpcService{client}

	natGateway, err := vpcService.DescribeNatGateway(d.Id())
	if err != nil {
		return WrapError(err)
	}

	d.Partial(true)
	attributeUpdate := false
	modifyNatGatewayAttributeRequest := vpc.CreateModifyNatGatewayAttributeRequest()
	modifyNatGatewayAttributeRequest.RegionId = natGateway.RegionId
	modifyNatGatewayAttributeRequest.NatGatewayId = natGateway.NatGatewayId

	if d.HasChange("name") {
		d.SetPartial("name")
		var name string
		if v, ok := d.GetOk("name"); ok {
			name = v.(string)
		} else {
			return WrapError(Error("cann't change name to empty string"))
		}
		modifyNatGatewayAttributeRequest.Name = name

		attributeUpdate = true
	}

	if d.HasChange("description") {
		d.SetPartial("description")
		var description string
		if v, ok := d.GetOk("description"); ok {
			description = v.(string)
		} else {
			return WrapError(Error("can to change description to empty string"))
		}

		modifyNatGatewayAttributeRequest.Description = description

		attributeUpdate = true
	}

	if attributeUpdate {
		raw, err := client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
			return vpcClient.ModifyNatGatewayAttribute(modifyNatGatewayAttributeRequest)
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), modifyNatGatewayAttributeRequest.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		addDebug(modifyNatGatewayAttributeRequest.GetActionName(), raw, modifyNatGatewayAttributeRequest.RpcRequest, modifyNatGatewayAttributeRequest)
	}

	if d.HasChange("specification") {
		d.SetPartial("specification")
		modifyNatGatewaySpecRequest := vpc.CreateModifyNatGatewaySpecRequest()
		modifyNatGatewaySpecRequest.RegionId = natGateway.RegionId
		modifyNatGatewaySpecRequest.NatGatewayId = natGateway.NatGatewayId
		modifyNatGatewaySpecRequest.Spec = d.Get("specification").(string)

		raw, err := client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
			return vpcClient.ModifyNatGatewaySpec(modifyNatGatewaySpecRequest)
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), modifyNatGatewaySpecRequest.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		addDebug(modifyNatGatewaySpecRequest.GetActionName(), raw, modifyNatGatewaySpecRequest.RpcRequest, modifyNatGatewaySpecRequest)
	}
	d.Partial(false)
	if err := vpcService.WaitForNatGateway(d.Id(), Available, DefaultTimeout); err != nil {
		return WrapError(err)
	}
	return resourceAliyunNatGatewayRead(d, meta)
}

func resourceAliyunNatGatewayDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	vpcService := VpcService{client}
	err := deleteBandwidthPackages(d, meta)
	if err != nil {
		return WrapError(err)
	}
	request := vpc.CreateDeleteNatGatewayRequest()
	request.RegionId = string(client.Region)
	request.NatGatewayId = d.Id()
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		request := vpc.CreateDeleteNatGatewayRequest()
		request.RegionId = string(client.Region)
		request.NatGatewayId = d.Id()
		raw, err := client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
			return vpcClient.DeleteNatGateway(request)
		})
		if err != nil {
			if IsExpectedErrors(err, []string{"DependencyViolation.BandwidthPackages"}) {
				return resource.RetryableError(err)
			}
			if IsExpectedErrors(err, []string{"InvalidNatGatewayId.NotFound"}) {
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
	return WrapError(vpcService.WaitForNatGateway(d.Id(), Deleted, DefaultTimeoutMedium))
}

func deleteBandwidthPackages(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	packages, err := DescribeNatBandwidthPackages(d.Id(), meta)
	if err != nil {
		return WrapError(err)
	}

	conn, err := meta.(*connectivity.AliyunClient).NewVpcClient()
	if err != nil {
		return WrapError(err)
	}

	for _, val := range packages {
		pg := val.(map[string]interface{})
		var response map[string]interface{}
		action := "DeleteBandwidthPackage"
		request := map[string]interface{}{
			"RegionId":           client.RegionId,
			"BandwidthPackageId": pg["BandwidthPackageId"].(string),
		}

		// If the API supports
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		err = resource.Retry(3*time.Minute, func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2016-04-28"), StringPointer("AK"), nil, request, &runtime)
			if err != nil {
				if IsExpectedErrors(err, []string{"Invalid.RegionId"}) {
					time.Sleep(5 * time.Second)
					return resource.RetryableError(err)
				} else if IsExpectedErrors(err, []string{"INSTANCE_NOT_EXISTS"}) {
					return nil
				}
				return resource.NonRetryableError(err)
			}
			addDebug(action, response, request)
			return nil
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), pg["BandwidthPackageId"].(string), AlibabaCloudSdkGoERROR)
		}
	}
	return nil
}

func flattenBandWidthPackages(natGatewayId string, meta interface{}) ([]map[string]interface{}, error) {
	packages, err := DescribeNatBandwidthPackages(natGatewayId, meta)
	if err != nil {
		return nil, WrapError(err)
	}

	var result []map[string]interface{}
	for _, val := range packages {
		pg := val.(map[string]interface{})
		var ipAddress []string
		publicIp := pg["PublicIpAddresses"].(map[string]interface{})["PublicIpAddresse"]
		publicIpAddresses := publicIp.([]interface{})
		for _, val := range publicIpAddresses {
			ipAddress = append(ipAddress, val.(map[string]interface{})["IpAddress"].(string))
		}

		ipCount, err1 := strconv.Atoi(pg["IpCount"].(string))
		if err1 != nil {
			return nil, WrapError(err1)
		}

		bandwidth, err2 := strconv.Atoi(pg["Bandwidth"].(string))
		if err2 != nil {
			return nil, WrapError(err2)
		}

		l := map[string]interface{}{
			"ip_count":            ipCount,
			"bandwidth":           bandwidth,
			"zone":                pg["ZoneId"].(string),
			"public_ip_addresses": strings.Join(ipAddress, ","),
		}
		result = append(result, l)
	}
	return result, nil
}

func DescribeNatBandwidthPackages(natGatewayId string, meta interface{}) ([]interface{}, error) {
	client := meta.(*connectivity.AliyunClient)
	addDebug("DescribeBandwidthPackages", natGatewayId, natGatewayId, "")
	var response map[string]interface{}
	action := "DescribeBandwidthPackages"
	request := map[string]interface{}{
		"RegionId":     client.RegionId,
		"NatGatewayId": natGatewayId,
	}

	conn, err := meta.(*connectivity.AliyunClient).NewVpcClient()
	if err != nil {
		return nil, WrapError(err)
	}
	// If the API supports
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	err = resource.Retry(3*time.Minute, func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2016-04-28"), StringPointer("AK"), nil, request, &runtime)
		if err != nil {
			if IsExpectedErrors(err, []string{"TaskConflict", "UnknownError", Throttling}) {
				time.Sleep(5 * time.Second)
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(action, response, request)
		return nil
	})
	if err != nil {
		return nil, WrapError(err)
	}

	packages := response["BandwidthPackages"].(map[string]interface{})
	return packages["BandwidthPackage"].([]interface{}), nil
}
