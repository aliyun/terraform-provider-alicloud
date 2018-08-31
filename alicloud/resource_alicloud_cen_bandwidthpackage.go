package alicloud

import (
	"fmt"
	"strings"

	"time"

	"log"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/cbn"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
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

			"geographic_region_id": &schema.Schema{
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
				Default:  "POSTPAY",
				ValidateFunc: func(v interface{}, k string) (ws []string, errors []error) {
					value := v.(string)
					if value != "PREPAY" && value != "POSTPAY" {
						errors = append(errors, fmt.Errorf("%s must be one of: PREPAY or POSTPAY", k))
					}
					return
				},
			},

			"auto_pay": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},

			"period": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				Default:  1,
			},

			"pricing_cycle": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  "Month",
				ValidateFunc: func(v interface{}, k string) (ws []string, errors []error) {
					value := v.(string)
					if value != "Month" && value != "Year" {
						errors = append(errors, fmt.Errorf("%s must be one of: Month or Year", k))
					}
					return
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

	client := meta.(*AliyunClient)

	var cenbwp *cbn.CreateCenBandwidthPackageResponse
	err := resource.Retry(3*time.Minute, func() *resource.RetryError {
		args, err := buildAliCloudCenBandwidthPackageArgs(d, meta)
		if err != nil {
			return resource.NonRetryableError(fmt.Errorf("Building CreateCenBandwidthPackageRequest got an error: %#v", err))
		}
		resp, err := client.cenconn.CreateCenBandwidthPackage(args)
		if err != nil {
			if IsExceptedError(err, BwpForSameSpanQuotaExceeded) {
				return resource.NonRetryableError(fmt.Errorf("Each span only allows to create a bandwidth package."))
			}
			if IsExceptedError(err, UnknownError) {
				return resource.RetryableError(fmt.Errorf("Create vpc timeout and got an error: %#v.", err))
			}
			return resource.NonRetryableError(err)
		}
		cenbwp = resp
		return nil
	})
	if err != nil {
		return fmt.Errorf("Create cen bandwidth package got an error :%#v", err)
	}

	d.SetId(cenbwp.CenBandwidthPackageId)

	err = client.WaitForCenBandwidthPackage(d.Id(), Idle, 60)
	if err != nil {
		log.Printf("%s", err)
		return fmt.Errorf("Timeout when WaitForCenBandwidthPackageIdle")
	}

	return resourceAlicloudCenBandwidthPackageUpdate(d, meta)
}

func resourceAlicloudCenBandwidthPackageRead(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*AliyunClient)
	resp, err := client.DescribeCenBandwidthPackage(d.Id())
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

	d.Set("geographic_region_id", geographicRegionId)
	d.Set("bandwidth", resp.Bandwidth)
	d.Set("name", resp.Name)
	d.Set("description", resp.Description)
	d.Set("expired_time", resp.ExpiredTime)
	d.Set("status", resp.Status)
	d.Set("charge_type", resp.BandwidthPackageChargeType)

	return nil
}

func resourceAlicloudCenBandwidthPackageUpdate(d *schema.ResourceData, meta interface{}) error {

	d.Partial(true)

	client := meta.(*AliyunClient)
	attributeUpdate := false
	request1 := cbn.CreateModifyCenBandwidthPackageAttributeRequest()
	request1.CenBandwidthPackageId = d.Id()

	if d.HasChange("name") {
		d.SetPartial("name")
		request1.Name = d.Get("name").(string)

		attributeUpdate = true
	}

	if d.HasChange("description") {
		d.SetPartial("description")
		request1.Description = d.Get("description").(string)

		attributeUpdate = true
	}

	if attributeUpdate {
		if _, err := client.cenconn.ModifyCenBandwidthPackageAttribute(request1); err != nil {
			return err
		}
	}

	request2 := cbn.CreateModifyCenBandwidthPackageSpecRequest()
	request2.CenBandwidthPackageId = d.Id()

	if d.HasChange("bandwidth") {
		d.SetPartial("bandwidth")
		bandwidth := d.Get("bandwidth").(int)
		request2.Bandwidth = requests.NewInteger(bandwidth)

		if _, err := client.cenconn.ModifyCenBandwidthPackageSpec(request2); err != nil {
			return err
		}
		/* modify function may delay for a little while */
		if err := client.WaitForCenBandwidthPackageUpdate(d.Id(), bandwidth, 10); err != nil {
			return err
		}
	}

	d.Partial(false)

	return resourceAlicloudCenBandwidthPackageRead(d, meta)
}

func resourceAlicloudCenBandwidthPackageDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*AliyunClient)
	request := cbn.CreateDeleteCenBandwidthPackageRequest()
	request.CenBandwidthPackageId = d.Id()

	return resource.Retry(5*time.Minute, func() *resource.RetryError {
		_, err := client.cenconn.DeleteCenBandwidthPackage(request)

		if err != nil {
			if IsExceptedError(err, ForbiddenRelease) {
				return resource.NonRetryableError(err)
			}
			return resource.RetryableError(fmt.Errorf("Delete CEN timeout and got an error: %#v.", err))
		}

		if _, err := client.DescribeCenBandwidthPackage(d.Id()); err != nil {
			if NotFoundError(err) {
				return nil
			}
			return resource.NonRetryableError(err)
		}

		return nil
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

func buildAliCloudCenBandwidthPackageArgs(d *schema.ResourceData, meta interface{}) (*cbn.CreateCenBandwidthPackageRequest, error) {
	request := cbn.CreateCreateCenBandwidthPackageRequest()
	request.Bandwidth = requests.NewInteger(d.Get("bandwidth").(int))

	geographicRegionId := d.Get("geographic_region_id").(*schema.Set).List()
	if len(geographicRegionId) == 1 {
		request.GeographicRegionAId = geographicRegionId[0].(string)
		request.GeographicRegionBId = geographicRegionId[0].(string)
	} else if len(geographicRegionId) == 2 {
		request.GeographicRegionAId = geographicRegionId[0].(string)
		request.GeographicRegionBId = geographicRegionId[1].(string)
	} else {
		return request, fmt.Errorf("Gepgraphic Region Id error")
	}

	if v := d.Get("name").(string); v != "" {
		request.Name = v
	}

	if v := d.Get("description").(string); v != "" {
		request.Description = v
	}

	request.Period = requests.NewInteger(d.Get("period").(int))        //default: 1
	request.PricingCycle = d.Get("pricing_cycle").(string)             //default: "Month"
	request.BandwidthPackageChargeType = d.Get("charge_type").(string) //default: "POSTPAY"
	request.AutoPay = requests.NewBoolean(d.Get("auto_pay").(bool))    //default: false

	return request, nil
}
