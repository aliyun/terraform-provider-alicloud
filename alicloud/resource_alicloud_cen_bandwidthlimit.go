package alicloud

import (
	"fmt"

	"strings"

	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/cbn"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceAlicloudCenBandwidthLimit() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudCenBandwidthLimitCreate,
		Read:   resourceAlicloudCenBandwidthLimitRead,
		Update: resourceAlicloudCenBandwidthLimitUpdate,
		Delete: resourceAlicloudCenBandwidthLimitDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"cen_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"regions_id": &schema.Schema{
				Type:     schema.TypeSet,
				Required: true,
				ForceNew: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				MaxItems: 2,
				MinItems: 2,
			},
			"bandwidth_limit": &schema.Schema{
				Type:     schema.TypeInt,
				Required: true,
				ValidateFunc: func(v interface{}, k string) (ws []string, errors []error) {
					value := v.(int)
					if value < 0 {
						errors = append(errors, fmt.Errorf("%s cannot be smaller than 0Mbps", k))
					}

					return
				},
			},
		},
	}
}

func resourceAlicloudCenBandwidthLimitCreate(d *schema.ResourceData, meta interface{}) error {
	client := (meta).(*AliyunClient)
	cenId := d.Get("cen_id").(string)
	bandwidthLimit := d.Get("bandwidth_limit").(int)

	regionsId := d.Get("regions_id").(*schema.Set).List()
	localRegionId := regionsId[0].(string)
	oppositeRegionId := regionsId[1].(string)

	request := cbn.CreateSetCenInterRegionBandwidthLimitRequest()
	request.CenId = cenId
	request.LocalRegionId = localRegionId
	request.OppositeRegionId = oppositeRegionId
	request.BandwidthLimit = requests.NewInteger(bandwidthLimit)

	err := resource.Retry(5*time.Minute, func() *resource.RetryError {
		_, err := client.cenconn.SetCenInterRegionBandwidthLimit(request)
		if err != nil {
			if IsExceptedError(err, InvalidCenInstanceStatus) {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("CEN ID %s Set BandwidthLimit got an error: %#v.", cenId, err)
	}

	if err := client.WaitForCenInterRegionBandwidthLimitActive(cenId, localRegionId, oppositeRegionId, 30); err != nil {
		return nil
	}

	d.SetId(cenId + ":" + localRegionId + ":" + oppositeRegionId)
	return resourceAlicloudCenBandwidthLimitRead(d, meta)
}

func resourceAlicloudCenBandwidthLimitRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*AliyunClient)

	paras, err := getParaForCenBandwidthLimit(d.Id())
	if err != nil {
		return err
	}

	cenId := paras[0]
	localRegionId := paras[1]
	oppositeRegionId := paras[2]

	resp, err := client.DescribeCenBandwidthLimit(cenId, localRegionId, oppositeRegionId)
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return err
	}

	respRegionIds := make([]string, 0)
	respRegionIds = append(respRegionIds, resp.LocalRegionId)
	respRegionIds = append(respRegionIds, resp.OppositeRegionId)

	d.Set("regions_id", respRegionIds)
	d.Set("cen_id", resp.CenId)
	d.Set("bandwidth_limit", resp.BandwidthLimit)

	return nil
}

func resourceAlicloudCenBandwidthLimitUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*AliyunClient)
	cenId := d.Get("cen_id").(string)
	regionsId := d.Get("regions_id").(*schema.Set).List()
	localRegionId := regionsId[0].(string)
	oppositeRegionId := regionsId[1].(string)
	var bandwidthLimit int

	d.Partial(true)

	attributeUpdate := false

	if d.HasChange("bandwidth_limit") {
		d.SetPartial("bandwidth_limit")
		bandwidthLimit = d.Get("bandwidth_limit").(int)

		attributeUpdate = true
	}

	if attributeUpdate {
		if err := client.SetCenInterRegionBandwidthLimit(d, bandwidthLimit); err != nil {
			return err
		}

		if bandwidthLimit == 0 {
			if err := client.WaitForCenInterRegionBandwidthLimitDestroy(cenId, localRegionId, oppositeRegionId, 30); err != nil {
				return err
			}
		} else {
			if err := client.WaitForCenInterRegionBandwidthLimitActive(cenId, localRegionId, oppositeRegionId, 30); err != nil {
				return err
			}
		}
	}

	d.Partial(false)

	return resourceAlicloudCenBandwidthLimitRead(d, meta)
}

func resourceAlicloudCenBandwidthLimitDelete(d *schema.ResourceData, meta interface{}) error {
	client := (meta).(*AliyunClient)
	cenId := d.Get("cen_id").(string)
	regionsId := d.Get("regions_id").(*schema.Set).List()
	localRegionId := regionsId[0].(string)
	oppositeRegionId := regionsId[1].(string)

	request := cbn.CreateSetCenInterRegionBandwidthLimitRequest()
	request.CenId = cenId
	request.LocalRegionId = localRegionId
	request.OppositeRegionId = oppositeRegionId
	request.BandwidthLimit = requests.NewInteger(0)

	err := resource.Retry(5*time.Minute, func() *resource.RetryError {
		_, err := client.cenconn.SetCenInterRegionBandwidthLimit(request)
		if err != nil {
			if IsExceptedError(err, InvalidCenInstanceStatus) {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("CEN ID %s Delete BandwidthLimit got an error: %#v.", cenId, err)
	}

	if err := client.WaitForCenInterRegionBandwidthLimitDestroy(cenId, localRegionId, oppositeRegionId, 30); err != nil {
		return err
	}

	return nil
}

func getParaForCenBandwidthLimit(id string) (retString []string, err error) {
	parts := strings.Split(id, ":")

	if len(parts) != 3 {
		return retString, fmt.Errorf("invalid resource id")
	}

	return parts, nil
}
