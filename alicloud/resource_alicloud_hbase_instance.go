package alicloud

import (
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/hbase"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func payTypePostPaidDiffSuppressFunc(k, old, new string, d *schema.ResourceData) bool {
	return strings.ToLower(d.Get("pay_type").(string)) == "postpaid"
}

func resourceAlicloudHBaseInstance() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudHBaseInstanceCreate,
		Read:   resourceAlicloudHBaseInstanceRead,
		Update: resourceAlicloudHBaseInstanceUpdate,
		Delete: resourceAlicloudHBaseInstanceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(2, 128),
			},
			"zone_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"engine": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"hbase"}, false),
				Default:      "hbase",
			},
			"engine_version": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"master_instance_type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"core_instance_type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"core_instance_quantity": {
				Type:         schema.TypeInt,
				ValidateFunc: validation.IntBetween(2, 20),
				Optional:     true,
				ForceNew:     true,
				Default:      2,
			},
			"core_disk_type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"cloud_ssd", "cloud_efficiency", "local_hdd_pro", "local_ssd_pro"}, false),
			},
			"core_disk_size": {
				Type:         schema.TypeInt,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.IntBetween(100, 2000),
				Default:      100,
			},
			"pay_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{string(PostPaid), string(PrePaid)}, false),
				Default:      PostPaid,
			},
			"duration": {
				Type:             schema.TypeInt,
				Optional:         true,
				ForceNew:         true,
				Computed:         true,
				ValidateFunc:     validation.IntInSlice([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 12, 24, 36, 60}),
				DiffSuppressFunc: payTypePostPaidDiffSuppressFunc,
			},
			"auto_renew": {
				Type:             schema.TypeBool,
				Optional:         true,
				ForceNew:         true,
				Computed:         true,
				DiffSuppressFunc: payTypePostPaidDiffSuppressFunc,
			},
			"vswitch_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"cold_storage_size": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
		},
	}
}

func buildHBaseCreateRequest(d *schema.ResourceData, meta interface{}) (*hbase.CreateInstanceRequest, error) {
	client := meta.(*connectivity.AliyunClient)
	vpcService := VpcService{client}
	request := hbase.CreateCreateInstanceRequest()
	request.ClusterName = Trim(d.Get("name").(string))
	request.RegionId = string(client.Region)
	request.ZoneId = Trim(d.Get("zone_id").(string))
	request.Engine = Trim(d.Get("engine").(string))
	request.DbType = request.Engine
	request.EngineVersion = Trim(d.Get("engine_version").(string))
	request.DbType = request.Engine
	request.MasterInstanceType = Trim(d.Get("master_instance_type").(string))
	request.CoreInstanceType = Trim(d.Get("core_instance_type").(string))
	request.CoreInstanceQuantity = strconv.Itoa(d.Get("core_instance_quantity").(int))
	request.CoreDiskType = Trim(d.Get("core_disk_type").(string))
	request.CoreDiskQuantity = "4"
	request.CoreDiskSize = strconv.Itoa(d.Get("core_disk_size").(int))
	request.PayType = Trim(d.Get("pay_type").(string))
	request.PricingCycle = "month"
	request.Duration = strconv.Itoa(d.Get("duration").(int))
	request.AutoRenew = strconv.FormatBool(d.Get("auto_renew").(bool))
	request.NetType = "classic"
	vswitchId := Trim(d.Get("vswitch_id").(string))
	if vswitchId != "" {
		request.NetType = "vpc"
		request.VSwitchId = vswitchId
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
			return nil, WrapError(Error("The specified vswitch %s isn't in the zone %s.", vsw.VSwitchId, request.ZoneId))
		}

		request.VpcId = vsw.VpcId
	}

	if d.Get("cold_storage_size").(int) > 0 {
		request.IsColdStorage = "true"
		request.ColdStorageSize = strconv.Itoa(d.Get("cold_storage_size").(int))
	} else {
		request.IsColdStorage = "false"
	}

	request.SecurityIPList = LOCAL_HOST_IP
	return request, nil
}

func resourceAlicloudHBaseInstanceCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	hBaseService := HBaseService{client}

	request, err := buildHBaseCreateRequest(d, meta)
	if err != nil {
		return WrapError(err)
	}

	raw, err := client.WithHbaseClient(func(client *hbase.Client) (interface{}, error) {
		return client.CreateInstance(request)
	})

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_hbase_instance", request.GetActionName(), AlibabaCloudSdkGoERROR)
	}

	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	response, _ := raw.(*hbase.CreateInstanceResponse)

	d.SetId(response.ClusterId)

	stateConf := BuildStateConf([]string{Hb_LAUNCHING, Hb_CREATING}, []string{Hb_ACTIVATION}, d.Timeout(schema.TimeoutCreate),
		10*time.Minute, hBaseService.HBaseClusterStateRefreshFunc(d.Id(), []string{Hb_CREATE_FAILED}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapError(err)
	}

	return resourceAlicloudHBaseInstanceRead(d, meta)
}

func resourceAlicloudHBaseInstanceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	hbaseService := HBaseService{client}

	instance, err := hbaseService.DescribeHBaseInstance(d.Id())
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("name", instance.InstanceName)
	d.Set("zone_id", instance.ZoneId)
	d.Set("engine", instance.Engine)
	d.Set("engine_version", instance.MajorVersion)
	d.Set("master_instance_type", instance.MasterInstanceType)
	d.Set("core_instance_type", instance.CoreInstanceType)
	d.Set("core_instance_quantity", instance.CoreNodeCount)
	d.Set("core_disk_size", instance.CoreDiskSize)
	d.Set("core_disk_type", instance.CoreDiskType)
	// Postpaid -> PostPaid
	if instance.PayType == string(Postpaid) {
		d.Set("pay_type", string(PostPaid))
	} else if instance.PayType == string(Prepaid) {
		d.Set("pay_type", string(PrePaid))
		period, err := computePeriodByUnit(instance.CreatedTimeUTC, instance.ExpireTimeUTC, d.Get("duration").(int), "Month")
		if err != nil {
			return WrapError(err)
		}
		d.Set("duration", period)
	}
	// now sdk can not get right value, "auto_renew", "is_cold_storage".
	d.Set("auto_renew", d.Get("auto_renew"))
	d.Set("cold_storage_size", d.Get("cold_storage_size"))
	d.Set("vpc_id", instance.VpcId)
	d.Set("vswitch_id", instance.VswitchId)
	return nil
}

func resourceAlicloudHBaseInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	if d.HasChange("name") {
		request := hbase.CreateModifyInstanceNameRequest()
		request.ClusterId = d.Id()
		request.ClusterName = d.Get("name").(string)

		raw, err := client.WithHbaseClient(func(client *hbase.Client) (interface{}, error) {
			return client.ModifyInstanceName(request)
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	}

	return resourceAlicloudHBaseInstanceRead(d, meta)
}

func resourceAlicloudHBaseInstanceDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	hbaseService := HBaseService{client}

	request := hbase.CreateDeleteInstanceRequest()
	request.ClusterId = d.Id()
	err := resource.Retry(10*5*time.Minute, func() *resource.RetryError {
		raw, err := client.WithHbaseClient(func(hbaseClient *hbase.Client) (interface{}, error) {
			return hbaseClient.DeleteInstance(request)
		})

		if err != nil {
			if hbaseService.NotFoundHBaseInstance(err) {
				return nil
			}
			return resource.RetryableError(err)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		return nil
	})

	if err != nil {
		if IsExceptedErrors(err, []string{"InvalidDBInstanceId.NotFound"}) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	stateConf := BuildStateConf([]string{Hb_DELETING}, []string{}, d.Timeout(schema.TimeoutDelete), 1*time.Minute, hbaseService.HBaseClusterStateRefreshFunc(d.Id(), []string{}))
	_, err = stateConf.WaitForState()
	return WrapError(err)
}
