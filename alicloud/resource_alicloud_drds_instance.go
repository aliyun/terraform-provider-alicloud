package alicloud

import (
	"fmt"
	"log"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/drds"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func resourceAliCloudDRDSInstance() *schema.Resource {
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
			"type": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateAllowedStringValue([]string{string(PrivateType), string(PrivateType_)}),
			},
			"zone_id": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"specification": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"pay_type": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateAllowedStringValue([]string{string(DRDSInstancePostPayType)}),
			},
			"vswitch_id": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"instance_series": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"instance_id": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceAliCloudDRDSInstanceCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	drdsService := DrdsService{client}

	req := drds.CreateCreateDrdsInstanceRequest()
	req.Description = d.Get("description").(string)
	req.Type = d.Get("type").(string)
	req.ZoneId = d.Get("zone_id").(string)
	req.Specification = d.Get("specification").(string)
	req.PayType = d.Get("pay_type").(string)
	req.VswitchId = d.Get("vswitch_id").(string)
	req.InstanceSeries = d.Get("instance_series").(string)
	req.Quantity = "1"
	response, err := drdsService.CreateDrdsInstance(req)
	idList := response.Data.DrdsInstanceIdList.DrdsInstanceId
	if err != nil || len(idList) != 1 {
		return fmt.Errorf("failed to create DRDS instance with error: %s", err)
	}
	d.SetId(idList[0])
	return resourceAliCloudDRDSInstanceUpdate(d, meta)

}

func resourceAliCloudDRDSInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	drdsService := DrdsService{client}
	update := false
	req := drds.CreateModifyDrdsInstanceDescriptionRequest()
	req.DrdsInstanceId = d.Id()
	if d.HasChange("description") && !d.IsNewResource() {
		update = true
		req.Description = d.Get("description").(string)
	}
	if update {
		_, err := drdsService.ModifyDrdsInstanceDescription(req)
		if err != nil {
			return fmt.Errorf("failed to update Drds instance with error: %s", err)
		}
	}
	return resourceAliCloudDRDSInstanceRead(d, meta)
}

func resourceAliCloudDRDSInstanceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	drdsService := DrdsService{client}

	res, err := drdsService.DescribeDrdsInstance(d.Id())
	data := res.Data
	if err != nil || res == nil || data.DrdsInstanceId == "" {
		log.Printf("[WARN] Failed to describe DRDS instance with error: %s", err)
		d.SetId("")
		return nil
	}
	// `description` isn't returned somehow, reported a bug https://connect.aliyun.com/suggestion/39734.
	//d.Set("description", data.Description)
	// As `describe` only return `type` 0 or 1, convert `type`. https://help.aliyun.com/document_detail/51126.html
	d.Set("type", convertTypeValue(data.Type, d.Get("type").(string)))
	d.Set("zone_id", data.ZoneId)
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
