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
			"description": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateStringLengthInRange(1, 129),
			},
			"zone_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"specification": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"instance_charge_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateAllowedStringValue([]string{string(PostPaid), string(PrePaid)}),
				ForceNew:     true,
				Default:      PostPaid,
			},
			"vswitch_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"instance_series": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateAllowedStringValue([]string{string(drds4c8g), string(drds8c16g), string(drds16c32g), string(drds32c64g)}),
				ForceNew:     true,
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

	if req.VswitchId != "" {

		vpcService := VpcService{client}
		vsw, err := vpcService.DescribeVswitch(req.VswitchId)
		if err != nil {
			return fmt.Errorf("DescribeVSwitche got an error: %#v.", err)
		}

		req.VpcId = vsw.VpcId
	}
	req.ClientToken = buildClientToken("drds")
	if req.PayType == "PostPaid" {
		req.PayType = "drdsPost"
	}
	if req.PayType == "PrePaid" {
		req.PayType = "drdsPre"
	}
	response, err := drdsService.CreateDrdsInstance(req)
	idList := response.Data.DrdsInstanceIdList.DrdsInstanceId
	if err != nil || len(idList) != 1 {
		return fmt.Errorf("failed to create DRDS instance with error: %s", err)
	}
	d.SetId(idList[0])

	// wait instance status change from Creating to running
	//0 -> running for drds,1->creating,2->exception,3->expire,4->release,5->locked
	if err := drdsService.WaitForDrdsInstance(d.Id(), "0", DefaultLongTimeout); err != nil {
		return fmt.Errorf("WaitForInstance %s %s got error: %#v", Running, d.Id(), err)
	}

	return resourceAliCloudDRDSInstanceRead(d, meta)

}

func resourceAliCloudDRDSInstanceUpdate(d *schema.ResourceData, meta interface{}) error {

	if d.HasChange("description") {
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
	return resourceAliCloudDRDSInstanceRead(d, meta)
}

func resourceAliCloudDRDSInstanceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	drdsService := DrdsService{client}

	res, err := drdsService.DescribeDrdsInstance(d.Id())
	if err != nil {
		if NotFoundError(err) {
			return nil
		}
	}
	if res == nil {
		return nil
	}
	data := res.Data
	//other attribute not set,because these attribute from `data` can't  get
	d.Set("zone_id", data.ZoneId)
	d.Set("description", data.Description)

	return nil
}

func resourceAliCloudDRDSInstanceDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	drdsService := DrdsService{client}

	return resource.Retry(5*time.Minute, func() *resource.RetryError {
		_, err := drdsService.DescribeDrdsInstance(d.Id())
		if err != nil {
			if NotFoundError(err) {
				return nil
			} else {
				return resource.NonRetryableError(fmt.Errorf(" got an error: %#v", err))
			}
		}

		removeRes, removeErr := drdsService.RemoveDrdsInstance(d.Id())
		if removeErr != nil || (removeRes != nil && !removeRes.Success) {
			return resource.RetryableError(fmt.Errorf("failed to delete instance timeout "+
				"and got an error: %#v", err))
		}
		return nil
	})
}
