package alicloud

import (
	"fmt"
	"strings"
	"time"

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
			"instance_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"region_ids": &schema.Schema{
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
					if value < 1 {
						errors = append(errors, fmt.Errorf("%s should be at least than 1 Mbps", k))
					}

					return
				},
			},
		},
	}
}

func resourceAlicloudCenBandwidthLimitCreate(d *schema.ResourceData, meta interface{}) error {
	cenId := d.Get("instance_id").(string)
	regionIds := d.Get("region_ids").(*schema.Set).List()
	localRegionId := regionIds[0].(string)
	oppositeRegionId := regionIds[1].(string)

	if strings.Compare(localRegionId, oppositeRegionId) <= 0 {
		d.SetId(cenId + COLON_SEPARATED + localRegionId + COLON_SEPARATED + oppositeRegionId)
	} else {
		d.SetId(cenId + COLON_SEPARATED + oppositeRegionId + COLON_SEPARATED + localRegionId)
	}

	return resourceAlicloudCenBandwidthLimitUpdate(d, meta)
}

func resourceAlicloudCenBandwidthLimitRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*AliyunClient)
	paras, err := getCenAndRegionIds(d.Id())
	if err != nil {
		return err
	}

	cenId := paras[0]
	localRegionId := paras[1]
	oppositeRegionId := paras[2]
	if strings.Compare(localRegionId, oppositeRegionId) > 0 {
		d.SetId(cenId + COLON_SEPARATED + oppositeRegionId + COLON_SEPARATED + localRegionId)
	}

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

	d.Set("region_ids", respRegionIds)
	d.Set("instance_id", resp.CenId)
	d.Set("bandwidth_limit", resp.BandwidthLimit)

	return nil
}

func resourceAlicloudCenBandwidthLimitUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*AliyunClient)
	cenId := d.Get("instance_id").(string)
	regionIds := d.Get("region_ids").(*schema.Set).List()
	localRegionId := regionIds[0].(string)
	oppositeRegionId := regionIds[1].(string)
	var bandwidthLimit int

	attributeUpdate := false
	if d.HasChange("bandwidth_limit") {
		attributeUpdate = true
		d.SetPartial("bandwidth_limit")
		bandwidthLimit = d.Get("bandwidth_limit").(int)
		if bandwidthLimit == 0 {
			return fmt.Errorf("the bandwidth limit should be at least than 1 Mbps")
		}
	}

	if attributeUpdate {
		err := resource.Retry(5*time.Minute, func() *resource.RetryError {
			err := client.SetCenInterRegionBandwidthLimit(cenId, localRegionId, oppositeRegionId, bandwidthLimit)
			if err != nil {
				if IsExceptedError(err, InvalidCenInstanceStatus) {
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})

		if err != nil {
			return fmt.Errorf("Create/Update bandwidth Limit CEN ID %s localRegionId %s oppositeRegionId %s got an error: %#v.",
				cenId, localRegionId, oppositeRegionId, err)
		}

		if err = client.WaitForCenInterRegionBandwidthLimitActive(cenId, localRegionId, oppositeRegionId, DefaultCenTimeout); err != nil {
			return err
		}
	}

	return resourceAlicloudCenBandwidthLimitRead(d, meta)
}

func resourceAlicloudCenBandwidthLimitDelete(d *schema.ResourceData, meta interface{}) error {
	client := (meta).(*AliyunClient)
	cenId := d.Get("instance_id").(string)
	regionIds := d.Get("region_ids").(*schema.Set).List()
	localRegionId := regionIds[0].(string)
	oppositeRegionId := regionIds[1].(string)

	err := resource.Retry(5*time.Minute, func() *resource.RetryError {
		err := client.SetCenInterRegionBandwidthLimit(cenId, localRegionId, oppositeRegionId, 0)
		if err != nil {
			if IsExceptedError(err, InvalidCenInstanceStatus) {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("delete bandwidth Limit CEN ID %s localRegionId %s oppositeRegionId %s got an error: %#v.",
			cenId, localRegionId, oppositeRegionId, err)
	}

	if err := client.WaitForCenInterRegionBandwidthLimitDestroy(cenId, localRegionId, oppositeRegionId, DefaultCenTimeout); err != nil {
		return err
	}

	return nil
}
