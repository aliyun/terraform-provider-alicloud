package alicloud

import (
	"fmt"
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
			"ip_address": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"instance": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
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
	request.ClientToken = buildClientToken("TF-AllocateEip")

	raw, err := client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
		return vpcClient.AllocateEipAddress(request)
	})
	if err != nil {
		if IsExceptedError(err, COMMODITYINVALID_COMPONENT) && request.InternetChargeType == string(PayByBandwidth) {
			return fmt.Errorf("Your account is international and it can only create '%s' elastic IP. Please change it and try again.", PayByTraffic)
		}
		return err
	}
	eip, _ := raw.(*vpc.AllocateEipAddressResponse)
	err = vpcService.WaitForEip(eip.AllocationId, Available, 60)
	if err != nil {
		return fmt.Errorf("Error Waitting for EIP available: %#v", err)
	}

	d.SetId(eip.AllocationId)

	return resourceAliyunEipUpdate(d, meta)
}

func resourceAliyunEipRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	vpcService := VpcService{client}

	eip, err := vpcService.DescribeEipAddress(d.Id())
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error Describe Eip Attribute: %#v", err)
	}

	// Output parameter 'instance' would be deprecated in the next version.
	if eip.InstanceId != "" {
		d.Set("instance", eip.InstanceId)
	} else {
		d.Set("instance", "")
	}

	d.Set("name", eip.Name)
	d.Set("description", eip.Descritpion)
	bandwidth, _ := strconv.Atoi(eip.Bandwidth)
	d.Set("bandwidth", bandwidth)
	d.Set("internet_charge_type", eip.InternetChargeType)
	d.Set("instance_charge_type", eip.ChargeType)
	d.Set("ip_address", eip.IpAddress)
	d.Set("status", eip.Status)

	return nil
}

func resourceAliyunEipUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	d.Partial(true)

	update := false
	request := vpc.CreateModifyEipAddressAttributeRequest()
	request.AllocationId = d.Id()

	if d.HasChange("bandwidth") && !d.IsNewResource() {
		update = true
		request.Bandwidth = strconv.Itoa(d.Get("bandwidth").(int))
		d.SetPartial("bandwidth")
	}
	if d.HasChange("name") {
		update = true
		request.Name = d.Get("name").(string)
		d.SetPartial("name")
	}
	if d.HasChange("description") {
		update = true
		request.Description = d.Get("description").(string)
		d.SetPartial("description")
	}
	if update {
		_, err := client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
			return vpcClient.ModifyEipAddressAttribute(request)
		})
		if err != nil {
			return err
		}
	}

	d.Partial(false)

	return resourceAliyunEipRead(d, meta)
}

func resourceAliyunEipDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	vpcService := VpcService{client}

	request := vpc.CreateReleaseEipAddressRequest()
	request.AllocationId = d.Id()

	return resource.Retry(5*time.Minute, func() *resource.RetryError {
		_, err := client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
			return vpcClient.ReleaseEipAddress(request)
		})
		if err != nil {
			if IsExceptedError(err, EipIncorrectStatus) {
				return resource.RetryableError(fmt.Errorf("Delete EIP timeout and got an error:%#v.", err))
			}
			return resource.NonRetryableError(err)

		}

		eip, descErr := vpcService.DescribeEipAddress(d.Id())

		if descErr != nil {
			if NotFoundError(descErr) {
				return nil
			}
			return resource.NonRetryableError(descErr)
		} else if eip.AllocationId == d.Id() {
			return resource.RetryableError(fmt.Errorf("Delete EIP timeout and it still exists."))
		}
		return nil
	})
}
