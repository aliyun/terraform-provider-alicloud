package alicloud

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/vpc"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
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
				ValidateFunc: validateInstanceDescription,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: func(v interface{}, k string) (ws []string, errors []error) {
					value := v.(string)
					if len(value) < 2 || len(value) > 128 {
						errors = append(errors, fmt.Errorf("%s cannot be longer than 128 characters", k))
					}

					if strings.HasPrefix(value, "http://") || strings.HasPrefix(value, "https://") {
						errors = append(errors, fmt.Errorf("%s cannot starts with http:// or https://", k))
					}

					return
				},
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
				ValidateFunc: validateCommonBandwidthPackageChargeType,
			},
			"ratio": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      100,
				ValidateFunc: validateRatio,
			},
		},
	}
}

func resourceAliyunCommonBandwidthPackageCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	commonBandwidthPackageService := CommonBandwidthPackageService{client}

	request := vpc.CreateCreateCommonBandwidthPackageRequest()
	request.RegionId = client.RegionId

	request.Bandwidth = requests.NewInteger(d.Get("bandwidth").(int))
	request.Name = d.Get("name").(string)
	request.Description = d.Get("description").(string)
	request.InternetChargeType = d.Get("internet_charge_type").(string)
	request.Ratio = requests.NewInteger(d.Get("ratio").(int))
	request.ClientToken = buildClientToken("TF-AllocateCommonBandwidthPackage")

	raw, err := client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
		return vpcClient.CreateCommonBandwidthPackage(request)
	})
	if err != nil {
		return err
	}
	commonBandwidthPackage, _ := raw.(*vpc.CreateCommonBandwidthPackageResponse)
	d.SetId(commonBandwidthPackage.BandwidthPackageId)

	if err := commonBandwidthPackageService.WaitForCommonBandwidthPackage(commonBandwidthPackage.BandwidthPackageId, DefaultTimeout); err != nil {
		return fmt.Errorf("Wait for common bandwidth package got error: %#v", err)
	}

	return resourceAliyunCommonBandwidthPackageRead(d, meta)
}

func resourceAliyunCommonBandwidthPackageRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	commonBandwidthPackageService := CommonBandwidthPackageService{client}

	resp, err := commonBandwidthPackageService.DescribeCommonBandwidthPackage(d.Id())
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error Describe Common Bandwidth Package Attribute: %#v", err)
	}

	if bandwidth, err := strconv.Atoi(resp.Bandwidth); err != nil {
		return fmt.Errorf("Convertting bandwidth from string to int got an error: %#v.", err)
	} else {
		d.Set("bandwidth", bandwidth)
	}
	d.Set("name", resp.Name)
	d.Set("description", resp.Description)
	d.Set("internet_charge_type", resp.InternetChargeType)
	d.Set("ratio", resp.Ratio)
	return nil
}

func resourceAliyunCommonBandwidthPackageUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	d.Partial(true)
	update := false
	request := vpc.CreateModifyCommonBandwidthPackageAttributeRequest()
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
		_, err := client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
			return vpcClient.ModifyCommonBandwidthPackageAttribute(request)
		})
		if err != nil {
			return err
		}
		d.SetPartial("description")
		d.SetPartial("name")
	}

	if d.HasChange("bandwidth") {
		request := vpc.CreateModifyCommonBandwidthPackageSpecRequest()
		request.RegionId = string(client.Region)
		request.BandwidthPackageId = d.Id()
		request.Bandwidth = strconv.Itoa(d.Get("bandwidth").(int))
		_, err := client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
			return vpcClient.ModifyCommonBandwidthPackageSpec(request)
		})
		if err != nil {
			return fmt.Errorf("ModifyCommonBandwidthPackageSpec got an error: %#v", err)
		}
		d.SetPartial("bandwidth")
	}

	d.Partial(false)
	return resourceAliyunCommonBandwidthPackageRead(d, meta)
}

func resourceAliyunCommonBandwidthPackageDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	commonBandwidthPackageService := CommonBandwidthPackageService{client}

	request := vpc.CreateDeleteCommonBandwidthPackageRequest()
	request.BandwidthPackageId = d.Id()

	return resource.Retry(5*time.Minute, func() *resource.RetryError {
		_, err := client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
			return vpcClient.DeleteCommonBandwidthPackage(request)
		})
		if err != nil {
			return resource.NonRetryableError(err)
		}
		if _, err := commonBandwidthPackageService.DescribeCommonBandwidthPackage(d.Id()); err != nil {
			if NotFoundError(err) {
				return nil
			}
			return resource.RetryableError(fmt.Errorf("Error describing common bandwidth package failed when deleting common bandwidth package: %#v", err))
		}
		return resource.RetryableError(fmt.Errorf("Delete Common Bandwidth Package timeout."))
	})
}
