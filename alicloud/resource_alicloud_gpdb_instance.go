package alicloud

import (
	"strings"
	"time"

	"github.com/hashicorp/terraform/helper/resource"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/gpdb"
	"github.com/denverdino/aliyungo/common"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func resourceAlicloudGpdbInstance() *schema.Resource {
	return &schema.Resource{
		Read:   resourceAlicloudGpdbInstanceRead,
		Create: resourceAlicloudGpdbInstanceCreate,
		Update: resourceAlicloudGpdbInstanceUpdate,
		Delete: resourceAlicloudGpdbInstanceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"availability_zone": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"instance_class": {
				Type:     schema.TypeString,
				Required: true,
			},
			"instance_group_count": {
				Type:     schema.TypeString,
				Required: true,
			},
			"instance_network_type": {
				Type:         schema.TypeString,
				ValidateFunc: validateAllowedStringValue([]string{string(common.Classic), string(common.VPC)}),
				Optional:     true,
				ForceNew:     true,
				Computed:     true,
			},
			"instance_charge_type": {
				Type:         schema.TypeString,
				ValidateFunc: validateAllowedStringValue([]string{string(PostPaid)}),
				Optional:     true,
				ForceNew:     true,
				Computed:     true,
			},
			"description": {
				Type:         schema.TypeString,
				ValidateFunc: validateDBInstanceName,
				Optional:     true,
			},
			"vpc_id": {
				Type:     schema.TypeString,
				ForceNew: true,
				Optional: true,
			},
			"vswitch_id": {
				Type:     schema.TypeString,
				ForceNew: true,
				Optional: true,
				Computed: true,
			},
			"security_ip_list": {
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
				Optional: true,
			},
			"engine": {
				Type:         schema.TypeString,
				ValidateFunc: validateAllowedStringValue([]string{string(GPDB)}),
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
			},
			"engine_version": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
		},
	}
}

func resourceAlicloudGpdbInstanceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	gpdbService := GpdbService{client}

	instance, err := gpdbService.DescribeGpdbInstance(d.Id())
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("instance_id", instance.DBInstanceId)
	d.Set("region_id", instance.RegionId)
	d.Set("availability_zone", instance.ZoneId)
	d.Set("engine", instance.Engine)
	d.Set("engine_version", instance.EngineVersion)
	d.Set("status", instance.DBInstanceStatus)
	d.Set("description", instance.DBInstanceDescription)
	d.Set("instance_class", instance.DBInstanceClass)
	d.Set("instance_group_count", instance.DBInstanceGroupCount)
	d.Set("instance_network_type", instance.InstanceNetworkType)
	security_ips, err := gpdbService.GetSecurityIps(d.Id())
	if err != nil {
		return WrapError(err)
	}
	d.Set("security_ip_list", security_ips)
	d.Set("create_time", instance.CreationTime)
	d.Set("instance_charge_type", instance.PayType)

	return nil
}

func resourceAlicloudGpdbInstanceCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	gpdbService := GpdbService{client}

	request, err := buildGpdbCreateRequest(d, meta)
	if err != nil {
		return WrapError(err)
	}
	raw, err := client.WithGpdbClient(func(client *gpdb.Client) (interface{}, error) {
		return client.CreateDBInstance(request)
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_gpdb_instance", request.GetActionName(), AlibabaCloudSdkGoERROR)
	}

	response, _ := raw.(*gpdb.CreateDBInstanceResponse)
	addDebug(request.GetActionName(), response)
	d.SetId(response.DBInstanceId)
	if err := gpdbService.WaitForGpdbInstance(d.Id(), Running, DefaultLongTimeout); err != nil {
		return WrapError(err)
	}

	return resourceAlicloudGpdbInstanceUpdate(d, meta)
}

func resourceAlicloudGpdbInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	gpdbService := GpdbService{client}

	// Begin Update
	d.Partial(true)

	// Update Instance Description
	if d.HasChange("description") {
		request := gpdb.CreateModifyDBInstanceDescriptionRequest()
		request.DBInstanceId = d.Id()
		request.DBInstanceDescription = d.Get("description").(string)
		raw, err := client.WithGpdbClient(func(gpdbClient *gpdb.Client) (interface{}, error) {
			return gpdbClient.ModifyDBInstanceDescription(request)
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		addDebug(request.GetActionName(), raw)
		d.SetPartial("description")
	}

	// Update Security Ips
	if d.HasChange("security_ip_list") {
		ipList := expandStringList(d.Get("security_ip_list").(*schema.Set).List())
		ipStr := strings.Join(ipList[:], COMMA_SEPARATED)
		// default disable connect from outside
		if ipStr == "" {
			ipStr = LOCAL_HOST_IP
		}
		if err := gpdbService.ModifyGpdbSecurityIps(d.Id(), ipStr); err != nil {
			return WrapError(err)
		}
		d.SetPartial("security_ip_list")
	}

	// Finish Update
	d.Partial(false)

	return resourceAlicloudGpdbInstanceRead(d, meta)
}

func resourceAlicloudGpdbInstanceDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	gpdbService := GpdbService{client}

	request := gpdb.CreateDeleteDBInstanceRequest()
	request.DBInstanceId = d.Id()

	err := resource.Retry(10*5*time.Minute, func() *resource.RetryError {
		raw, err := client.WithGpdbClient(func(client *gpdb.Client) (interface{}, error) {
			return client.DeleteDBInstance(request)
		})

		if err != nil {
			if IsExceptedErrors(err, []string{InvalidGpdbInstanceStatus}) {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(request.GetActionName(), raw)
		return nil
	})
	if err != nil {
		if gpdbService.NotFoundGpdbInstance(err) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	return WrapError(gpdbService.WaitForGpdbInstance(d.Id(), Deleted, DefaultLongTimeout))
}

func buildGpdbCreateRequest(d *schema.ResourceData, meta interface{}) (*gpdb.CreateDBInstanceRequest, error) {
	client := meta.(*connectivity.AliyunClient)
	request := gpdb.CreateCreateDBInstanceRequest()
	request.RegionId = string(client.Region)
	request.ZoneId = Trim(d.Get("availability_zone").(string))
	request.PayType = d.Get("instance_charge_type").(string)
	request.VSwitchId = Trim(d.Get("vswitch_id").(string))
	request.DBInstanceDescription = d.Get("description").(string)
	request.DBInstanceClass = Trim(d.Get("instance_class").(string))
	request.DBInstanceGroupCount = Trim(d.Get("instance_group_count").(string))
	request.Engine = Trim(d.Get("engine").(string))
	request.EngineVersion = Trim(d.Get("engine_version").(string))

	// Instance NetWorkType
	request.InstanceNetworkType = string(Classic)
	if request.VSwitchId != "" {
		// check vswitchId in zone
		vpcService := VpcService{client}
		object, err := vpcService.DescribeVSwitch(request.VSwitchId)
		if err != nil {
			return nil, WrapError(err)
		}

		if request.ZoneId == "" {
			request.ZoneId = object.ZoneId
		} else if strings.Contains(request.ZoneId, MULTI_IZ_SYMBOL) {
			zoneStr := strings.Split(strings.SplitAfter(request.ZoneId, "(")[1], ")")[0]
			if !strings.Contains(zoneStr, string([]byte(object.ZoneId)[len(object.ZoneId)-1])) {
				return nil, WrapError(Error("The specified vswitch %s isn't in the multi zone %s.", object.VSwitchId, request.ZoneId))
			}
		} else if request.ZoneId != object.ZoneId {
			return nil, WrapError(Error("The specified vswitch %s isn't in the zone %s.", object.VSwitchId, request.ZoneId))
		}

		request.VPCId = object.VpcId
		request.InstanceNetworkType = strings.ToUpper(string(Vpc))
	}

	// Security Ips
	request.SecurityIPList = LOCAL_HOST_IP
	if len(d.Get("security_ip_list").(*schema.Set).List()) > 0 {
		request.SecurityIPList = strings.Join(expandStringList(d.Get("security_ip_list").(*schema.Set).List())[:], COMMA_SEPARATED)
	}

	// ClientToken
	request.ClientToken = buildClientToken(request.GetActionName())

	return request, nil
}
