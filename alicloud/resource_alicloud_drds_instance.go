package alicloud

import (
	"github.com/aliyun/alibaba-cloud-sdk-go/services/drds"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
	"time"
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

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
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

	request := drds.CreateCreateDrdsInstanceRequest()
	request.Description = d.Get("description").(string)
	request.Type = "1"
	request.ZoneId = d.Get("zone_id").(string)
	request.Specification = d.Get("specification").(string)
	request.PayType = d.Get("instance_charge_type").(string)
	request.VswitchId = d.Get("vswitch_id").(string)
	request.InstanceSeries = d.Get("instance_series").(string)
	request.Quantity = "1"

	if request.VswitchId != "" {

		vpcService := VpcService{client}
		vsw, err := vpcService.DescribeVSwitch(request.VswitchId)
		if err != nil {
			return WrapError(err)
		}

		request.VpcId = vsw.VpcId
	}
	request.ClientToken = buildClientToken(request.GetActionName())
	if request.PayType == string(PostPaid) {
		request.PayType = "drdsPost"
	}
	if request.PayType == string(PrePaid) {
		request.PayType = "drdsPre"
	}
	raw, err := client.WithDrdsClient(func(drdsClient *drds.Client) (interface{}, error) {
		return drdsClient.CreateDrdsInstance(request)
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_drds_instance", request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw)
	response, _ := raw.(*drds.CreateDrdsInstanceResponse)
	idList := response.Data.DrdsInstanceIdList.DrdsInstanceId
	if len(idList) != 1 {
		return WrapError(Error("failed to create DRDS instance"))
	}
	d.SetId(idList[0])

	// wait instance status change from Creating to running
	//0 -> running for drds,1->creating,2->exception,3->expire,4->release,5->locked
	stateConf := BuildStateConf([]string{"1"}, []string{"0"}, d.Timeout(schema.TimeoutCreate), 1*time.Minute, drdsService.DrdsInstanceStateRefreshFunc(d.Id(), []string{"2"}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapError(err)
	}

	return resourceAliCloudDRDSInstanceUpdate(d, meta)

}

func resourceAliCloudDRDSInstanceUpdate(d *schema.ResourceData, meta interface{}) error {

	if d.HasChange("description") {
		request := drds.CreateModifyDrdsInstanceDescriptionRequest()
		request.DrdsInstanceId = d.Id()
		request.Description = d.Get("description").(string)
		client := meta.(*connectivity.AliyunClient)
		raw, err := client.WithDrdsClient(func(drdsClient *drds.Client) (interface{}, error) {
			return drdsClient.ModifyDrdsInstanceDescription(request)
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		addDebug(request.GetActionName(), raw)
	}
	return resourceAliCloudDRDSInstanceRead(d, meta)
}

func resourceAliCloudDRDSInstanceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	drdsService := DrdsService{client}

	object, err := drdsService.DescribeDrdsInstance(d.Id())
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	data := object.Data
	//other attribute not set,because these attribute from `data` can't  get
	d.Set("zone_id", data.ZoneId)
	d.Set("description", data.Description)

	return nil
}

func resourceAliCloudDRDSInstanceDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	drdsService := DrdsService{client}
	request := drds.CreateRemoveDrdsInstanceRequest()
	request.DrdsInstanceId = d.Id()
	raw, err := client.WithDrdsClient(func(drdsClient *drds.Client) (interface{}, error) {
		return drdsClient.RemoveDrdsInstance(request)
	})
	if err != nil {
		if IsExceptedError(err, InvalidDRDSInstanceIdNotFound) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw)
	response, _ := raw.(*drds.RemoveDrdsInstanceResponse)
	if !response.Success {
		return WrapError(Error("failed to delete instance timeout "+"and got an error: %#v", err))
	}

	stateConf := BuildStateConf([]string{"0", "1", "2", "3", "4", "5", "6"}, []string{}, d.Timeout(schema.TimeoutDelete), 3*time.Second, drdsService.DrdsInstanceStateRefreshFunc(d.Id(), []string{}))
	_, err = stateConf.WaitForState()
	return WrapError(err)
}
