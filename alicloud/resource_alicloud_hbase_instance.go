package alicloud

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/hbase"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

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
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"region_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"zone_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"engine": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "hbase",
			},
			"engine_version": {
				Type:     schema.TypeString,
				Required: true,
			},
			"master_instance_type": {
				Type:     schema.TypeString,
				Required: true,
			},
			"core_instance_type": {
				Type:     schema.TypeString,
				Required: true,
			},
			"core_instance_quantity": {
				Type:         schema.TypeInt,
				ValidateFunc: validation.IntBetween(2, 20),
				Optional:     true,
				Default:      2,
			},
			"core_disk_type": {
				Type:     schema.TypeString,
				Required: true,
			},
			"core_disk_size": {
				Type:         schema.TypeInt,
				ValidateFunc: validation.IntBetween(100, 2000),
				Optional:     true,
				Default:      100,
			},
			"pay_type": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  Postpaid,
			},
			"pricing_cycle": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "month",
			},
			"duration": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"auto_renew": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  "false",
			},
			"net_type": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"vswitch_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"is_cold_storage": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"security_ip_list": {
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Optional: true,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			// hbase 后续的参数，增加参数校验，
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
	request.Engine = "hbase"
	request.EngineVersion = Trim(d.Get("engine_version").(string))
	request.DbType = "hbase"
	request.MasterInstanceType = Trim(d.Get("master_instance_type").(string))
	request.CoreInstanceType = Trim(d.Get("core_instance_type").(string))

	request.CoreInstanceQuantity = strconv.Itoa(d.Get("core_instance_quantity").(int))
	request.CoreDiskType = Trim(d.Get("core_disk_type").(string))
	request.CoreDiskQuantity = "4"
	request.CoreDiskSize = strconv.Itoa(d.Get("core_disk_size").(int))
	request.PayType = Trim(d.Get("pay_type").(string))
	request.PricingCycle = Trim(d.Get("pricing_cycle").(string))
	request.Duration = Trim(d.Get("duration").(string))
	request.AutoRenew = strconv.FormatBool(d.Get("auto_renew").(bool))
	request.NetType = string("classic")
	vswitchId := Trim(d.Get("vswitch_id").(string))
	if vswitchId != "" {
		request.NetType = string("vpc")
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

	request.IsColdStorage = strconv.FormatBool(d.Get("is_cold_storage").(bool))

	request.SecurityIPList = LOCAL_HOST_IP
	if len(d.Get("security_ip_list").(*schema.Set).List()) > 0 {
		request.SecurityIPList = strings.Join(expandStringList(d.Get("security_ip_list").(*schema.Set).List())[:], COMMA_SEPARATED)
	}
	return request, nil
}

func resourceAlicloudHBaseInstanceCreate(d *schema.ResourceData, meta interface{}) error {
	log.Println("[DEBUG] hbase-log id=", d)
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

	return resourceAlicloudHBaseInstanceUpdate(d, meta)
}

func resourceAlicloudHBaseInstanceRead(d *schema.ResourceData, meta interface{}) error {
	log.Println("resourceAlicloudHBaseInstanceRead id=", d.Id())
	client := meta.(*connectivity.AliyunClient)
	hbaseService := HBaseService{client}

	instance, err := hbaseService.DescribeHBaseInstance(d.Id())
	fmt.Println(instance)
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	_ = d.Set("name", instance.InstanceName)
	_ = d.Set("region_id", instance.RegionId)
	_ = d.Set("zone_id", instance.ZoneId)
	_ = d.Set("engine", instance.Engine)
	_ = d.Set("engine_version", instance.MajorVersion)
	_ = d.Set("master_instance_type", instance.MasterInstanceType)
	_ = d.Set("core_instance_type", instance.CoreInstanceType)
	_ = d.Set("core_instance_quantity", instance.CoreNodeCount)
	_ = d.Set("core_disk_size", instance.CoreDiskSize)
	_ = d.Set("core_disk_type", instance.CoreDiskType)
	_ = d.Set("pay_type", instance.PayType)
	_ = d.Set("status", instance.Status)
	_ = d.Set("net_type", instance.NetworkType)
	_ = d.Set("vpc_id", instance.VpcId)
	_ = d.Set("vswitch_id", instance.VswitchId)
	log.Println("[DEBUG] hbase-log result d=", d)
	return nil
}

func resourceAlicloudHBaseInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
	log.Println("[DEBUG] hbase-log id=", d)
	client := meta.(*connectivity.AliyunClient)

	d.Partial(true)
	// hbase 更改更多的信息
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
		d.SetPartial("name")
	}

	d.Partial(false)
	return resourceAlicloudHBaseInstanceRead(d, meta)
}

func resourceAlicloudHBaseInstanceDelete(d *schema.ResourceData, meta interface{}) error {
	log.Println("[DEBUG] hbase-log id=", d)
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
				return resource.NonRetryableError(err)
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
	stateConf := BuildStateConf([]string{"Creating", "Deleting"}, []string{}, d.Timeout(schema.TimeoutDelete), 1*time.Minute, hbaseService.HBaseClusterStateRefreshFunc(d.Id(), []string{}))
	_, err = stateConf.WaitForState()
	return WrapError(err)
}
