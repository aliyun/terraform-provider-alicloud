package alicloud

import (
	"strings"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/rds"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func resourceAlicloudDBReadonlyInstance() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudDBReadonlyInstanceCreate,
		Read:   resourceAlicloudDBReadonlyInstanceRead,
		Update: resourceAlicloudDBReadonlyInstanceUpdate,
		Delete: resourceAlicloudDBReadonlyInstanceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"engine_version": &schema.Schema{
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},

			"master_db_instance_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			"instance_name": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateDBInstanceName,
				Computed:     true,
			},

			"instance_type": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			"instance_storage": &schema.Schema{
				Type:     schema.TypeInt,
				Required: true,
			},

			"zone_id": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},

			"vswitch_id": &schema.Schema{
				Type:     schema.TypeString,
				ForceNew: true,
				Optional: true,
			},

			"parameters": &schema.Schema{
				Type: schema.TypeSet,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": &schema.Schema{
							Type:     schema.TypeString,
							Required: true,
						},
						"value": &schema.Schema{
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
				Set:      parameterToHash,
				Optional: true,
				Computed: true,
			},

			"engine": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"connection_string": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},

			"port": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceAlicloudDBReadonlyInstanceCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	rdsService := RdsService{client}

	request, err := buildDBReadonlyCreateRequest(d, meta)
	if err != nil {
		return WrapError(err)
	}

	raw, err := client.WithRdsClient(func(rdsClient *rds.Client) (interface{}, error) {
		return rdsClient.CreateReadOnlyDBInstance(request)
	})

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw)
	resp, _ := raw.(*rds.CreateReadOnlyDBInstanceResponse)
	d.SetId(resp.DBInstanceId)

	// wait instance status change from Creating to running
	if err := rdsService.WaitForDBInstance(d.Id(), Running, DefaultLongTimeout); err != nil {
		return WrapError(err)
	}

	return resourceAlicloudDBReadonlyInstanceUpdate(d, meta)
}

func resourceAlicloudDBReadonlyInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	rdsService := RdsService{client}
	d.Partial(true)

	if d.HasChange("parameters") {
		if err := rdsService.ModifyParameters(d, "parameters"); err != nil {
			return WrapError(err)
		}
	}

	if d.IsNewResource() {
		d.Partial(false)
		return resourceAlicloudDBInstanceRead(d, meta)
	}

	if d.HasChange("instance_name") {
		request := rds.CreateModifyDBInstanceDescriptionRequest()
		request.DBInstanceId = d.Id()
		request.DBInstanceDescription = d.Get("instance_name").(string)

		err := resource.Retry(5*time.Minute, func() *resource.RetryError {
			raw, err := client.WithRdsClient(func(rdsClient *rds.Client) (interface{}, error) {
				return rdsClient.ModifyDBInstanceDescription(request)
			})
			if err != nil {
				if IsExceptedErrors(err, []string{"OperationDenied.DBInstanceStatus", "OperationDenied.MasterDBInstanceState"}) {
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}

			addDebug(request.GetActionName(), raw)

			d.SetPartial("instance_name")
			return nil
		})

		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
		}

	}

	update := false
	request := rds.CreateModifyDBInstanceSpecRequest()
	request.DBInstanceId = d.Id()
	request.PayType = string(Postpaid)

	if d.HasChange("instance_type") {
		request.DBInstanceClass = d.Get("instance_type").(string)
		update = true
	}

	if d.HasChange("instance_storage") {
		request.DBInstanceStorage = requests.NewInteger(d.Get("instance_storage").(int))
		update = true
	}

	if update {
		// wait instance status is running before modifying
		if err := rdsService.WaitForDBInstance(d.Id(), Running, DefaultTimeoutMedium); err != nil {
			return WrapError(err)
		}

		err := resource.Retry(5*time.Minute, func() *resource.RetryError {
			raw, err := client.WithRdsClient(func(rdsClient *rds.Client) (interface{}, error) {
				return rdsClient.ModifyDBInstanceSpec(request)
			})
			if err != nil {
				if IsExceptedErrors(err, []string{"InvalidOrderTask.NotSupport", "OperationDenied.DBInstanceStatus", "OperationDenied.MasterDBInstanceState"}) {
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			addDebug(request.GetActionName(), raw)
			d.SetPartial("instance_type")
			d.SetPartial("instance_storage")
			return nil
		})

		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
		}

		// wait instance status is running after modifying
		if err := rdsService.WaitForDBInstance(d.Id(), Running, 2*DefaultLongTimeout); err != nil {
			return WrapError(err)
		}
	}

	d.Partial(false)
	return resourceAlicloudDBReadonlyInstanceRead(d, meta)
}

func resourceAlicloudDBReadonlyInstanceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	rdsService := RdsService{client}

	instance, err := rdsService.DescribeDBInstance(d.Id())
	if err != nil {
		if rdsService.NotFoundDBInstance(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("engine", instance.Engine)
	d.Set("master_db_instance_id", instance.MasterInstanceId)
	d.Set("engine_version", instance.EngineVersion)
	d.Set("instance_type", instance.DBInstanceClass)
	d.Set("port", instance.Port)
	d.Set("instance_storage", instance.DBInstanceStorage)
	d.Set("zone_id", instance.ZoneId)
	d.Set("vswitch_id", instance.VSwitchId)
	d.Set("connection_string", instance.ConnectionString)
	d.Set("instance_name", instance.DBInstanceDescription)

	if err = rdsService.RefreshParameters(d, "parameters"); err != nil {
		return err
	}

	return nil
}

func resourceAlicloudDBReadonlyInstanceDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	rdsService := RdsService{client}

	instance, err := rdsService.DescribeDBInstance(d.Id())
	if err != nil {
		if rdsService.NotFoundDBInstance(err) {
			return nil
		}
		return WrapError(err)
	}
	if PayType(instance.PayType) == Prepaid {
		return WrapError(Error("At present, 'Prepaid' instance cannot be deleted and must wait it to be expired and release it automatically."))
	}

	request := rds.CreateDeleteDBInstanceRequest()
	request.DBInstanceId = d.Id()

	err = resource.Retry(5*time.Minute, func() *resource.RetryError {

		raw, err := client.WithRdsClient(func(rdsClient *rds.Client) (interface{}, error) {
			return rdsClient.DeleteDBInstance(request)
		})

		if err != nil {
			if IsExceptedErrors(err, []string{"RwSplitNetType.Exist", "OperationDenied.DBInstanceStatus", "OperationDenied.MasterDBInstanceState"}) {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}

		addDebug(request.GetActionName(), raw)

		return nil
	})

	if err != nil {
		if rdsService.NotFoundDBInstance(err) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
	}

	return WrapError(rdsService.WaitForDBInstance(d.Id(), Deleted, DefaultTimeoutMedium))

}

func buildDBReadonlyCreateRequest(d *schema.ResourceData, meta interface{}) (*rds.CreateReadOnlyDBInstanceRequest, error) {
	client := meta.(*connectivity.AliyunClient)
	vpcService := VpcService{client}
	request := rds.CreateCreateReadOnlyDBInstanceRequest()
	request.RegionId = string(client.Region)
	request.DBInstanceId = Trim(d.Get("master_db_instance_id").(string))
	request.EngineVersion = Trim(d.Get("engine_version").(string))
	request.DBInstanceStorage = requests.NewInteger(d.Get("instance_storage").(int))
	request.DBInstanceClass = Trim(d.Get("instance_type").(string))
	request.DBInstanceDescription = d.Get("instance_name").(string)

	if zone, ok := d.GetOk("zone_id"); ok && Trim(zone.(string)) != "" {
		request.ZoneId = Trim(zone.(string))
	}

	vswitchId := Trim(d.Get("vswitch_id").(string))

	request.InstanceNetworkType = string(Classic)

	if vswitchId != "" {
		request.VSwitchId = vswitchId
		request.InstanceNetworkType = strings.ToUpper(string(Vpc))

		// check vswitchId in zone
		vsw, err := vpcService.DescribeVSwitch(vswitchId)
		if err != nil {
			return nil, WrapError(err)
		}

		if request.ZoneId == "" {
			request.ZoneId = vsw.ZoneId
		} else if strings.Contains(request.ZoneId, MULTI_IZ_SYMBOL) {
			zonestr := strings.Split(strings.SplitAfter(request.ZoneId, "(")[1], ")")[0]
			if !strings.Contains(zonestr, string([]byte(vsw.ZoneId)[len(vsw.ZoneId)-1])) {
				return nil, WrapError(Error("The specified vswitch %s isn't in the multi zone %s.", vsw.VSwitchId, request.ZoneId))
			}
		} else if request.ZoneId != vsw.ZoneId {
			return nil, WrapError(Error("The specified vswitch %s isn't in the multi zone %s.", vsw.VSwitchId, request.ZoneId))
		}

		request.VPCId = vsw.VpcId
	}

	request.PayType = string(Postpaid)
	request.ClientToken = buildClientToken(request.GetActionName())

	return request, nil
}
