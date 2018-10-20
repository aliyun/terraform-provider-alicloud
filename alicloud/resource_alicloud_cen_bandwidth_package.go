package alicloud

import (
	"fmt"
	"strings"

	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/cbn"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
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
			"bandwidth": &schema.Schema{
				Type:     schema.TypeInt,
				Required: true,
				ValidateFunc: func(v interface{}, k string) (ws []string, errors []error) {
					value := v.(int)
					if value < 1 {
						errors = append(errors, fmt.Errorf("%s cannot be smaller than 1Mbps", k))
					}

					return
				},
			},

			"geographic_region_ids": &schema.Schema{
				Type:     schema.TypeSet,
				Required: true,
				ForceNew: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				MaxItems: 2,
				MinItems: 1,
			},

			"name": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: func(v interface{}, k string) (ws []string, errors []error) {
					value := v.(string)
					if len(value) < 2 || len(value) > 128 {
						errors = append(errors, fmt.Errorf("%s cannot be shorter than 2 characters or longer than 128 characters", k))
					}

					if strings.HasPrefix(value, "http://") || strings.HasPrefix(value, "https://") {
						errors = append(errors, fmt.Errorf("%s cannot starts with http:// or https://", k))
					}

					return
				},
			},

			"description": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: func(v interface{}, k string) (ws []string, errors []error) {
					value := v.(string)
					if len(value) < 2 || len(value) > 256 {
						errors = append(errors, fmt.Errorf("%s cannot be shorter than 2 characters or longer than 256 characters", k))
					}

					if strings.HasPrefix(value, "http://") || strings.HasPrefix(value, "https://") {
						errors = append(errors, fmt.Errorf("%s cannot starts with http:// or https://", k))
					}

					return
				},
			},

			"charge_type": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  PostPaid,
				ValidateFunc: func(v interface{}, k string) (ws []string, errors []error) {
					value := PayType(v.(string))
					if value != PrePaid && value != PostPaid {
						errors = append(errors, fmt.Errorf("%s must be one of: %s or %s", k, string(PrePaid), string(PostPaid)))
					}
					return
				},
			},

			"period": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				Default:  1,
				ValidateFunc: func(v interface{}, k string) (ws []string, errors []error) {
					value := v.(int)
					if value != 1 && value != 2 && value != 3 && value != 6 && value != 12 {
						errors = append(errors, fmt.Errorf("%s must be one of: 1, 2, 3, 6, 12", k))
					}
					return
				},
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return PayType(d.Get("charge_type").(string)) == PostPaid
				},
			},

			"expired_time": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},

			"status": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceAlicloudCenBandwidthPackageCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	cenService := CenService{client}
	var cenbwp *cbn.CreateCenBandwidthPackageResponse
	err := resource.Retry(3*time.Minute, func() *resource.RetryError {
		args := buildAliCloudCenBandwidthPackageArgs(d, meta)
		raw, err := client.WithCenClient(func(cbnClient *cbn.Client) (interface{}, error) {
			return cbnClient.CreateCenBandwidthPackage(args)
		})
		if err != nil {
			if IsExceptedError(err, OperationBlocking) {
				return resource.RetryableError(fmt.Errorf("Create CEN bandwidth package timeout and got an error: %#v.", err))
			}
			return resource.NonRetryableError(err)
		}
		resp, _ := raw.(*cbn.CreateCenBandwidthPackageResponse)
		cenbwp = resp
		return nil
	})
	if err != nil {
		return fmt.Errorf("Create CEN bandwidth package got an error :%#v", err)
	}

	d.SetId(cenbwp.CenBandwidthPackageId)

	err = cenService.WaitForCenBandwidthPackage(d.Id(), Idle, DefaultCenTimeout)
	if err != nil {
		return fmt.Errorf("Timeout when WaitForCenBandwidthPackageIdle, bandwidth package ID %s", err)
	}

	return resourceAlicloudCenBandwidthPackageRead(d, meta)
}

func resourceAlicloudCenBandwidthPackageRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	cenService := CenService{client}
	resp, err := cenService.DescribeCenBandwidthPackage(d.Id())
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return err
	}

	geographicRegionId := make([]string, 0)
	geographicRegionId = append(geographicRegionId, convertGeographicRegionId(resp.GeographicRegionAId))
	geographicRegionId = append(geographicRegionId, convertGeographicRegionId(resp.GeographicRegionBId))

	d.Set("geographic_region_ids", geographicRegionId)
	d.Set("bandwidth", resp.Bandwidth)
	d.Set("name", resp.Name)
	d.Set("description", resp.Description)
	d.Set("expired_time", resp.ExpiredTime)
	d.Set("status", resp.Status)

	if resp.BandwidthPackageChargeType == "POSTPAY" {
		d.Set("charge_type", PostPaid)
	} else {
		d.Set("charge_type", PrePaid)
	}

	return nil
}

func resourceAlicloudCenBandwidthPackageUpdate(d *schema.ResourceData, meta interface{}) error {
	d.Partial(true)

	client := meta.(*connectivity.AliyunClient)
	cenService := CenService{client}
	attributeUpdate := false
	request1 := cbn.CreateModifyCenBandwidthPackageAttributeRequest()
	request1.CenBandwidthPackageId = d.Id()

	if d.HasChange("name") {
		request1.Name = d.Get("name").(string)
		attributeUpdate = true
	}

	if d.HasChange("description") {
		request1.Description = d.Get("description").(string)
		attributeUpdate = true
	}

	if attributeUpdate {
		_, err := client.WithCenClient(func(cbnClient *cbn.Client) (interface{}, error) {
			return cbnClient.ModifyCenBandwidthPackageAttribute(request1)
		})
		if err != nil {
			return err
		}
		d.SetPartial("name")
		d.SetPartial("description")
	}

	if d.HasChange("bandwidth") {
		bandwidth := d.Get("bandwidth").(int)
		request2 := cbn.CreateModifyCenBandwidthPackageSpecRequest()
		request2.CenBandwidthPackageId = d.Id()
		request2.Bandwidth = requests.NewInteger(bandwidth)

		_, err := client.WithCenClient(func(cbnClient *cbn.Client) (interface{}, error) {
			return cbnClient.ModifyCenBandwidthPackageSpec(request2)
		})
		if err != nil {
			return err
		}
		/* modify function may delay for a while */
		if err := cenService.WaitForCenBandwidthPackageUpdate(d.Id(), bandwidth, DefaultCenTimeout); err != nil {
			return err
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

	return resource.Retry(5*time.Minute, func() *resource.RetryError {
		_, err := client.WithCenClient(func(cbnClient *cbn.Client) (interface{}, error) {
			return cbnClient.DeleteCenBandwidthPackage(request)
		})
		if err != nil {
			if IsExceptedErrors(err, []string{ForbiddenRelease, InvalidCenBandwidthLimitsNotZero}) {
				return resource.NonRetryableError(err)
			}
			if IsExceptedError(err, ParameterBwpInstanceId) {
				return nil
			}
			return resource.RetryableError(fmt.Errorf("Delete CEN bandwidth package %s timeout and got an error: %#v.", cenBwpId, err))
		}

		if _, err := cenService.DescribeCenBandwidthPackage(cenBwpId); err != nil {
			if NotFoundError(err) {
				return nil
			}
			return resource.NonRetryableError(err)
		}

		return resource.RetryableError(fmt.Errorf("Delete CEN bandwidth package %s, but not completed", cenBwpId))
	})
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
		request.GeographicRegionAId = geographicRegionId[0].(string)
		request.GeographicRegionBId = geographicRegionId[1].(string)
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
	request.ClientToken = buildClientToken("TF-CreateCenBandwidthPackage")

	return request
}
