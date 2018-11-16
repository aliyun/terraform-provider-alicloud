package alicloud

import (
	"fmt"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/drds"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func resourceAlicloudDRDSInstance() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudDRDSInstanceCreate,
		Read:   resourceAliCloudDRDSInstanceRead,
		Update: resourceAliCloudDRDSInstanceUpdate,
		Delete: resourceAliCloudDRDSInstanceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"description": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateStringLengthInRange(1, 129),
			},
			"zone_id": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"specification": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"instance_charge_type": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateAllowedStringValue([]string{string(Postpaid), string(Prepaid)}),
			},
			"vswitch_id": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"instance_series": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateAllowedStringValue([]string{string(drds4c8g), string(drds8c16g), string(drds16c32g), string(drds32c64g)}),
			},
		},
	}
}

func resourceAliCloudDRDSInstanceCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	drdsService := DrdsService{client}

	req := drds.CreateCreateDrdsInstanceRequest()
	req.Description = d.Get("description").(string)
	req.Type = "1"
	req.ZoneId = d.Get("zone_id").(string)
	req.Specification = d.Get("specification").(string)
	req.PayType = d.Get("instance_charge_type").(string)
	req.VswitchId = d.Get("vswitch_id").(string)
	req.InstanceSeries = d.Get("instance_series").(string)
	req.Quantity = "1"
	response, err := drdsService.CreateDrdsInstance(req)
	idList := response.Data.DrdsInstanceIdList.DrdsInstanceId
	if err != nil || len(idList) != 1 {
		return fmt.Errorf("failed to create DRDS instance with error: %s", err)
	}
	d.SetId(idList[0])

	// wait instance status change from Creating to running
	//0 -> running for drds
	//https://help.aliyun.com/document_detail/51126.html?spm=a2c4g.11174283.6.757.31eb73543ixaAc
	if err := drdsService.WaitForDrdsInstance(d.Id(), "0", DefaultLongTimeout); err != nil {
		return fmt.Errorf("WaitForInstance %s got error: %#v", Running, err)
	}

	return resourceAliCloudDRDSInstanceRead(d, meta)

}

func resourceAliCloudDRDSInstanceUpdate(d *schema.ResourceData, meta interface{}) error {

	update := false

	if d.HasChange("description") {
		update = true
	}
	if update {
		req := drds.CreateModifyDrdsInstanceDescriptionRequest()
		req.DrdsInstanceId = d.Id()
		req.Description = d.Get("description").(string)
		client := meta.(*connectivity.AliyunClient)
		drdsService := DrdsService{client}
		_, err := drdsService.ModifyDrdsInstanceDescription(req)
		if err != nil {
			return fmt.Errorf("failed to update Drds instance with error: %s", err)
		}
	}
	return nil
}

func resourceAliCloudDRDSInstanceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	drdsService := DrdsService{client}

	res, err := drdsService.DescribeDrdsInstance(d.Id())
	data := res.Data
	if err != nil {
		if NotFoundError(err) {
			return nil
		}
	}
	d.Set("zone_id", data.ZoneId)
	d.Set("drdsInstanceId", data.DrdsInstanceId)
	d.Set("status", data.Status)
	d.Set("netWorkType", data.NetworkType)
	d.Set("region_id", data.RegionId)
	d.Set("specification", data.Specification)
	d.Set("instance_charge_type", data.Type)
	d.Set("description", data.Description)
	return nil
}

func resourceAliCloudDRDSInstanceDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	drdsService := DrdsService{client}

	return resource.Retry(5*time.Minute, func() *resource.RetryError {
		req := drds.CreateDescribeDrdsInstanceRequest()
		req.DrdsInstanceId = d.Id()
		res, err := drdsService.DescribeDrdsInstance(d.Id())
		if err != nil {
			if NotFoundError(err) {
				return nil
			}
		}
		if res == nil || res.Data.DrdsInstanceId == "" {
			return nil
		}
		removeReq := drds.CreateRemoveDrdsInstanceRequest()
		removeReq.DrdsInstanceId = d.Id()
		removeRes, removeErr := drdsService.RemoveDrdsInstance(d.Id())
		if removeErr != nil || (removeRes != nil && !removeRes.Success) {
			return resource.RetryableError(fmt.Errorf("failed to delete instance timeout "+
				"and got an error: %#v", err))
		}
		return nil
	})
}
