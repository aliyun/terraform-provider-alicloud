package alicloud

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAliCloudLindormInstance() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudLindormInstanceCreate,
		Read:   resourceAliCloudLindormInstanceRead,
		Update: resourceAliCloudLindormInstanceUpdate,
		Delete: resourceAliCloudLindormInstanceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(90 * time.Minute),
			Update: schema.DefaultTimeout(180 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"vswitch_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"disk_category": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: StringInSlice([]string{"cloud_efficiency", "cloud_ssd", "cloud_essd", "cloud_essd_pl0", "capacity_cloud_storage", "local_ssd_pro", "local_hdd_pro"}, false),
			},
			"payment_type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: StringInSlice([]string{"PayAsYouGo", "Subscription"}, false),
			},
			"cold_storage": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.Any(IntInSlice([]int{0}), IntBetween(800, 1000000)),
			},
			"core_num": {
				Type:     schema.TypeInt,
				Optional: true,
				Removed:  "Field `core_num` has been deprecated from provider version 1.188.0, and it has been removed from provider version 1.207.0.",
			},
			"core_spec": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"deletion_proection": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"duration": {
				Type:     schema.TypeString,
				Optional: true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return d.Get("payment_type").(string) != "Subscription"
				},
			},
			"file_engine_node_count": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: IntAtLeast(2),
			},
			"file_engine_specification": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: StringInSlice([]string{"lindorm.c.xlarge"}, false),
			},
			"group_name": {
				Type:     schema.TypeString,
				Optional: true,
				Removed:  "Field `group_name` has been removed from provider version 1.211.0.",
			},
			"instance_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"instance_storage": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"ip_white_list": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"lts_node_count": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"lts_node_specification": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: StringInSlice([]string{"lindorm.g.xlarge", "lindorm.g.2xlarge"}, false),
			},
			"phoenix_node_count": {
				Type:     schema.TypeInt,
				Optional: true,
				Removed:  "Field `phoenix_node_count` has been removed from provider version 1.211.0.",
			},
			"phoenix_node_specification": {
				Type:     schema.TypeString,
				Optional: true,
				Removed:  "Field `phoenix_node_specification` has been removed from provider version 1.211.0.",
			},
			"pricing_cycle": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: StringInSlice([]string{"Month", "Year"}, false),
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return d.Get("payment_type").(string) != "Subscription"
				},
			},
			"search_engine_node_count": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: IntAtLeast(2),
			},
			"search_engine_specification": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: StringInSlice([]string{"lindorm.g.xlarge", "lindorm.g.2xlarge", "lindorm.g.4xlarge", "lindorm.g.8xlarge"}, false),
			},
			"table_engine_node_count": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: IntAtLeast(2),
			},
			"table_engine_specification": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"time_series_engine_node_count": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: IntAtLeast(2),
			},
			"time_serires_engine_specification": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"time_series_engine_specification"},
				Deprecated:    "Field `time_serires_engine_specification` has been deprecated from provider version 1.182.0. New field `time_series_engine_specification` instead.",
			},
			"time_series_engine_specification": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"time_serires_engine_specification"},
			},
			"stream_engine_node_count": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"stream_engine_specification": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"upgrade_type": {
				Type:     schema.TypeString,
				Optional: true,
				Removed:  "Field `upgrade_type` has been deprecated from provider version 1.163.0, and it has been removed from provider version 1.207.0.",
			},
			"vpc_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"zone_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"resource_group_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"log_num": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: IntBetween(4, 400),
			},
			"log_single_storage": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: IntBetween(400, 64000),
			},
			"arbiter_zone_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"multi_zone_combination": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: StringInSlice([]string{"ap-southeast-5abc-aliyun", "cn-hangzhou-ehi-aliyun", "cn-beijing-acd-aliyun", "ap-southeast-1-abc-aliyun", "cn-zhangjiakou-abc-aliyun", "cn-shanghai-efg-aliyun", "cn-shanghai-abd-aliyun", "cn-hangzhou-bef-aliyun", "cn-hangzhou-bce-aliyun", "cn-beijing-fgh-aliyun", "cn-shenzhen-abc-aliyun"}, false),
			},
			"arbiter_vswitch_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"standby_zone_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"log_spec": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: StringInSlice([]string{"lindorm.sn1.large", "lindorm.sn1.2xlarge"}, false),
			},
			"log_disk_category": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: StringInSlice([]string{"cloud_efficiency", "cloud_ssd"}, false),
			},
			"core_single_storage": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: IntBetween(400, 64000),
			},
			"standby_vswitch_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"arch_version": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Computed:     true,
				ValidateFunc: StringInSlice([]string{"1.0", "2.0"}, false),
			},
			"primary_vswitch_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"primary_zone_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"tags": tagsSchema(),
			"service_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"enabled_file_engine": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"enabled_time_serires_engine": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"enabled_table_engine": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"enabled_search_engine": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"enabled_lts_engine": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"enabled_stream_engine": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceAliCloudLindormInstanceCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "CreateLindormInstance"
	request := make(map[string]interface{})
	var err error
	if v, ok := d.GetOk("cold_storage"); ok {
		request["ColdStorage"] = v
	}
	request["DiskCategory"] = d.Get("disk_category")
	if v, ok := d.GetOk("duration"); ok {
		request["Duration"] = v
	}
	if v, ok := d.GetOk("file_engine_node_count"); ok {
		request["FilestoreNum"] = v
	}
	if v, ok := d.GetOk("file_engine_specification"); ok {
		request["FilestoreSpec"] = v
	}
	if v, ok := d.GetOk("instance_name"); ok {
		request["InstanceAlias"] = v
	}
	if v, ok := d.GetOk("instance_storage"); ok {
		request["InstanceStorage"] = v
	}
	request["PayType"] = convertLindormInstancePaymentTypeRequest(d.Get("payment_type").(string))
	if v, ok := d.GetOk("pricing_cycle"); ok {
		request["PricingCycle"] = v
	}
	request["RegionId"] = client.RegionId
	if v, ok := d.GetOk("search_engine_node_count"); ok {
		request["SolrNum"] = v
	}
	if v, ok := d.GetOk("search_engine_specification"); ok {
		request["SolrSpec"] = v
	}
	if v, ok := d.GetOk("table_engine_node_count"); ok {
		request["LindormNum"] = v
	}
	if v, ok := d.GetOk("table_engine_specification"); ok {
		request["LindormSpec"] = v
	}
	if v, ok := d.GetOk("time_series_engine_node_count"); ok {
		request["TsdbNum"] = v
	}
	if v, ok := d.GetOk("time_serires_engine_specification"); ok {
		request["TsdbSpec"] = v
	} else if v, ok := d.GetOk("time_series_engine_specification"); ok {
		request["TsdbSpec"] = v
	}

	if v, ok := d.GetOkExists("stream_engine_node_count"); ok {
		request["StreamNum"] = v
	}

	if v, ok := d.GetOk("stream_engine_specification"); ok {
		request["StreamSpec"] = v
	}

	if v, ok := d.GetOk("zone_id"); ok {
		request["ZoneId"] = v
	}

	if v, ok := d.GetOk("vpc_id"); ok {
		request["VPCId"] = v
	}

	if v, ok := d.GetOk("vswitch_id"); ok {
		request["VSwitchId"] = v
	}

	if v, ok := d.GetOk("resource_group_id"); ok {
		request["ResourceGroupId"] = v
	}

	if v, ok := d.GetOk("core_spec"); ok {
		request["CoreSpec"] = v
	}

	if v, ok := d.GetOk("log_num"); ok {
		request["LogNum"] = v
	}

	if v, ok := d.GetOk("log_single_storage"); ok {
		request["LogSingleStorage"] = v
	}

	if v, ok := d.GetOk("arbiter_zone_id"); ok {
		request["ArbiterZoneId"] = v
	}

	if v, ok := d.GetOk("multi_zone_combination"); ok {
		request["MultiZoneCombination"] = v
	}

	if v, ok := d.GetOk("arbiter_vswitch_id"); ok {
		request["ArbiterVSwitchId"] = v
	}

	if v, ok := d.GetOk("standby_zone_id"); ok {
		request["StandbyZoneId"] = v
	}

	if v, ok := d.GetOk("log_spec"); ok {
		request["LogSpec"] = v
	}

	if v, ok := d.GetOk("log_disk_category"); ok {
		request["LogDiskCategory"] = v
	}

	if v, ok := d.GetOk("core_single_storage"); ok {
		request["CoreSingleStorage"] = v
	}

	if v, ok := d.GetOk("standby_vswitch_id"); ok {
		request["StandbyVSwitchId"] = v
	}

	if v, ok := d.GetOk("primary_zone_id"); ok {
		request["PrimaryZoneId"] = v
	}

	if v, ok := d.GetOk("primary_vswitch_id"); ok {
		request["PrimaryVSwitchId"] = v
	}

	if v, ok := d.GetOk("arch_version"); ok {
		request["ArchVersion"] = v
	}

	if (request["ZoneId"] == nil || request["VpcId"] == nil) && request["VSwitchId"] != nil {
		vpcService := VpcService{client}
		vsw, err := vpcService.DescribeVSwitchWithTeadsl(request["VSwitchId"].(string))
		if err != nil {
			return WrapError(err)
		}
		if v, ok := request["VPCId"].(string); !ok || v == "" {
			request["VPCId"] = vsw["VpcId"]
		}
		if v, ok := request["ZoneId"].(string); !ok || v == "" {
			request["ZoneId"] = vsw["ZoneId"]
		}
	}

	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutCreate)), func() *resource.RetryError {
		response, err = client.RpcPost("hitsdb", "2020-06-15", action, nil, request, false)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_lindorm_instance", action, AlibabaCloudSdkGoERROR)
	}
	d.SetId(fmt.Sprint(response["InstanceId"]))
	hitsdbService := HitsdbService{client}
	stateConf := BuildStateConf([]string{}, []string{"ACTIVATION"}, d.Timeout(schema.TimeoutCreate), 1*time.Minute, hitsdbService.LindormInstanceStateRefreshFunc(d.Id(), []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return resourceAliCloudLindormInstanceUpdate(d, meta)
}

func resourceAliCloudLindormInstanceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	hitsdbService := HitsdbService{client}
	object, err := hitsdbService.DescribeLindormInstance(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_lindorm_instance hitsdbService.DescribeLindormInstance Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	if v, ok := object["ColdStorage"]; ok {
		d.Set("cold_storage", formatInt(v))
	}

	d.Set("deletion_proection", object["DeletionProtection"])
	d.Set("disk_category", object["DiskCategory"])
	d.Set("instance_name", object["InstanceAlias"])
	d.Set("payment_type", convertLindormInstancePaymentTypeResponse(object["PayType"]))
	d.Set("status", object["InstanceStatus"])
	d.Set("vswitch_id", object["VswitchId"])
	d.Set("zone_id", object["ZoneId"])
	d.Set("primary_zone_id", object["PrimaryZoneId"])
	d.Set("primary_vswitch_id", object["PrimaryVSwitchId"])
	d.Set("resource_group_id", object["ResourceGroupId"])
	d.Set("vpc_id", object["VpcId"])

	//if object["ServiceType"] == "lindorm_multizone" {
	d.Set("log_num", object["LogNum"])
	d.Set("log_single_storage", object["LogSingleStorage"])
	d.Set("arbiter_zone_id", object["ArbiterZoneId"])
	d.Set("multi_zone_combination", object["MultiZoneCombination"])
	d.Set("arbiter_vswitch_id", object["ArbiterVSwitchId"])
	d.Set("standby_zone_id", object["StandbyZoneId"])
	d.Set("arbiter_vswitch_id", object["ArbiterVSwitchId"])
	d.Set("log_spec", object["LogSpec"])
	d.Set("log_disk_category", object["LogDiskCategory"])
	d.Set("standby_vswitch_id", object["StandbyVSwitchId"])
	//}
	d.Set("arch_version", object["ArchVersion"])
	d.Set("core_spec", object["CoreSpec"])
	//}

	//if object["DiskCategory"] != "local_ssd_pro" && object["DiskCategory"] != "local_hdd_pro" {
	d.Set("core_single_storage", object["CoreSingleStorage"])
	d.Set("instance_storage", object["InstanceStorage"])
	//}

	engineType := formatInt(object["EngineType"])
	d.Set("enabled_file_engine", engineType&0x08 == 8)
	d.Set("enabled_time_serires_engine", engineType&0x02 == 2)
	d.Set("enabled_table_engine", engineType&0x04 == 4)
	d.Set("enabled_search_engine", engineType&0x01 == 1)
	d.Set("enabled_lts_engine", object["EnableBDS"])
	d.Set("enabled_stream_engine", object["EnableStream"])

	ipWhite, err := hitsdbService.GetInstanceIpWhiteList(d.Id())
	if err != nil {
		return WrapError(err)
	}
	d.Set("ip_white_list", ipWhite)
	getLindormInstanceEngineInfoObject, err := hitsdbService.GetLindormInstanceEngineInfo(d.Id())
	if err != nil {
		return WrapError(err)
	}
	if v, ok := getLindormInstanceEngineInfoObject["FileEngineNodeCount"]; ok {
		d.Set("file_engine_node_count", formatInt(v))
	}
	d.Set("file_engine_specification", getLindormInstanceEngineInfoObject["FileEngineSpecification"])
	if v, ok := getLindormInstanceEngineInfoObject["LtsNodeCount"]; ok {
		d.Set("lts_node_count", formatInt(v))
	}
	d.Set("lts_node_specification", getLindormInstanceEngineInfoObject["LtsNodeSpecification"])

	if v, ok := getLindormInstanceEngineInfoObject["SearchEngineNodeCount"]; ok {
		d.Set("search_engine_node_count", formatInt(v))
	}
	d.Set("search_engine_specification", getLindormInstanceEngineInfoObject["SearchEngineSpecification"])
	if v, ok := getLindormInstanceEngineInfoObject["TableEngineNodeCount"]; ok {
		d.Set("table_engine_node_count", formatInt(v))
	}
	d.Set("table_engine_specification", getLindormInstanceEngineInfoObject["TableEngineSpecification"])
	if v, ok := getLindormInstanceEngineInfoObject["TimeSeriesNodeCount"]; ok {
		d.Set("time_series_engine_node_count", formatInt(v))
	}
	d.Set("time_serires_engine_specification", getLindormInstanceEngineInfoObject["TimeSeriesSpecification"])
	d.Set("time_series_engine_specification", getLindormInstanceEngineInfoObject["TimeSeriesSpecification"])

	if v, ok := getLindormInstanceEngineInfoObject["StreamNodeCount"]; ok {
		d.Set("stream_engine_node_count", formatInt(v))
	}
	d.Set("stream_engine_specification", getLindormInstanceEngineInfoObject["StreamSpecification"])

	listTagResourcesObject, err := hitsdbService.ListTagResources(d.Id(), "INSTANCE")
	if err != nil {
		return WrapError(err)
	}
	d.Set("tags", tagsToMap(listTagResourcesObject))
	d.Set("service_type", object["ServiceType"])

	return nil
}

func resourceAliCloudLindormInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	hitsdbService := HitsdbService{client}
	var response map[string]interface{}
	var err error
	d.Partial(true)

	if d.HasChange("tags") {
		if err := hitsdbService.SetResourceTags(d, "INSTANCE"); err != nil {
			return WrapError(err)
		}
		d.SetPartial("tags")
	}
	update := false
	request := map[string]interface{}{
		"InstanceId": d.Id(),
	}
	if d.HasChange("ip_white_list") {
		update = true
		if v, ok := d.GetOk("ip_white_list"); ok && v != nil {
			request["SecurityIpList"] = convertListToCommaSeparate(v.(*schema.Set).List())
		}
	}
	if update {
		action := "UpdateInstanceIpWhiteList"
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *resource.RetryError {
			response, err = client.RpcPost("hitsdb", "2020-06-15", action, nil, request, false)
			if err != nil {
				if IsExpectedErrors(err, []string{"Instance.IsNotValid"}) || NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})
		addDebug(action, response, request)
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
		d.SetPartial("ip_white_list")
	}
	update = false
	updateLindormInstanceAttributeReq := map[string]interface{}{
		"InstanceId": d.Id(),
	}
	if d.HasChange("instance_name") && !d.IsNewResource() {
		update = true
		if v, ok := d.GetOk("instance_name"); ok {
			updateLindormInstanceAttributeReq["InstanceAlias"] = v
		}
	}
	if d.HasChange("deletion_proection") {
		update = true
		if v, ok := d.GetOkExists("deletion_proection"); ok {
			updateLindormInstanceAttributeReq["DeletionProtection"] = v
		}
	}
	if update {
		action := "UpdateLindormInstanceAttribute"

		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *resource.RetryError {
			response, err = client.RpcPost("hitsdb", "2020-06-15", action, nil, updateLindormInstanceAttributeReq, false)
			if err != nil {
				if IsExpectedErrors(err, []string{"Instance.IsNotValid"}) || NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})
		addDebug(action, response, updateLindormInstanceAttributeReq)
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
		d.SetPartial("instance_name")
		d.SetPartial("deletion_proection")
	}

	update = false
	upgradeLindormLogReq := map[string]interface{}{}

	if !d.IsNewResource() && d.HasChange("log_single_storage") {
		update = true
		upgradeLindormLogReq["UpgradeType"] = "upgrade-disk-size"
		upgradeLindormLogReq["LogSingleStorage"] = d.Get("log_single_storage")
	}

	if !d.IsNewResource() && d.HasChange("log_spec") {
		update = true
		upgradeLindormLogReq["UpgradeType"] = "upgrade-lindorm-engine"
		upgradeLindormLogReq["LogSpec"] = d.Get("log_spec")
	}
	if !d.IsNewResource() && d.HasChange("core_single_storage") {
		update = true
		upgradeLindormLogReq["UpgradeType"] = "upgrade-disk-size"
		upgradeLindormLogReq["CoreSingleStorage"] = d.Get("core_single_storage")
	}
	if update {
		err := UpgradeLindormInstance(d, meta, upgradeLindormLogReq)
		if err != nil {
			return WrapError(err)
		}
		d.SetPartial("log_single_storage")
		d.SetPartial("core_single_storage")
		d.SetPartial("log_spec")
	}

	update = false
	upgradeLindormInstanceColdStorageReq := map[string]interface{}{
		"UpgradeType": "upgrade-cold-storage",
	}
	if d.HasChange("cold_storage") && !d.IsNewResource() {
		update = true
		if v, ok := d.GetOk("cold_storage"); ok {
			upgradeLindormInstanceColdStorageReq["ColdStorage"] = v
		}
		d.SetPartial("cold_storage")
	}
	if update {
		err := UpgradeLindormInstance(d, meta, upgradeLindormInstanceColdStorageReq)
		if err != nil {
			return WrapError(err)
		}
	}

	update = false
	upgradeLindormInstanceFilestoreNumReq := map[string]interface{}{
		"UpgradeType": "upgrade-file-core-num",
	}
	if d.HasChange("file_engine_node_count") && !d.IsNewResource() {
		update = true
		if v, ok := d.GetOk("file_engine_node_count"); ok {
			upgradeLindormInstanceFilestoreNumReq["FilestoreNum"] = v
		}
	}
	if update {
		if v, ok := d.GetOk("instance_storage"); ok {
			upgradeLindormInstanceFilestoreNumReq["ClusterStorage"] = v
		}
		err := UpgradeLindormInstance(d, meta, upgradeLindormInstanceFilestoreNumReq)
		if err != nil {
			return WrapError(err)
		}
		d.SetPartial("file_engine_node_count")
	}

	update = false
	upgradeLindormInstanceFilestoreSpecReq := map[string]interface{}{
		"UpgradeType": "upgrade-file-engine",
	}
	if d.HasChange("file_engine_specification") && !d.IsNewResource() {
		update = true
		if v, ok := d.GetOk("file_engine_specification"); ok {
			upgradeLindormInstanceFilestoreSpecReq["FilestoreSpec"] = v
		}
	}
	if update {
		if v, ok := d.GetOk("instance_storage"); ok {
			upgradeLindormInstanceFilestoreNumReq["ClusterStorage"] = v
		}
		err := UpgradeLindormInstance(d, meta, upgradeLindormInstanceFilestoreSpecReq)
		if err != nil {
			return WrapError(err)
		}
		d.SetPartial("file_engine_specification")
	}

	if (d.HasChange("search_engine_node_count") || d.HasChange("search_engine_specification")) && !d.IsNewResource() {
		newSolrSpec := d.Get("search_engine_specification")
		newSolrNum := d.Get("search_engine_node_count")
		enabled := d.Get("enabled_search_engine").(bool)
		currentInstanceStorage := formatInt(d.Get("instance_storage"))
		if !enabled {
			upgradeLindormInstanceSearchReq := map[string]interface{}{}
			upgradeLindormInstanceSearchReq["UpgradeType"] = "open-search-engine"
			upgradeLindormInstanceSearchReq["SolrSpec"] = newSolrSpec
			upgradeLindormInstanceSearchReq["SolrNum"] = newSolrNum
			upgradeLindormInstanceSearchReq["ClusterStorage"] = currentInstanceStorage
			err := UpgradeLindormInstance(d, meta, upgradeLindormInstanceSearchReq)
			if err != nil {
				return WrapError(err)
			}
		}

		if enabled && d.HasChange("search_engine_specification") {
			upgradeLindormInstanceSearchReq := map[string]interface{}{}
			upgradeLindormInstanceSearchReq["UpgradeType"] = "upgrade-search-engine"
			upgradeLindormInstanceSearchReq["SolrSpec"] = newSolrSpec
			upgradeLindormInstanceSearchReq["ClusterStorage"] = currentInstanceStorage
			err := UpgradeLindormInstance(d, meta, upgradeLindormInstanceSearchReq)
			if err != nil {
				return WrapError(err)
			}
		}

		if enabled && d.HasChange("search_engine_node_count") {
			upgradeLindormInstanceSearchNumReq := map[string]interface{}{}
			upgradeLindormInstanceSearchNumReq["UpgradeType"] = "upgrade-search-core-num"
			upgradeLindormInstanceSearchNumReq["SolrNum"] = newSolrNum
			upgradeLindormInstanceSearchNumReq["ClusterStorage"] = currentInstanceStorage
			err := UpgradeLindormInstance(d, meta, upgradeLindormInstanceSearchNumReq)
			if err != nil {
				return WrapError(err)
			}
		}

		d.SetPartial("search_engine_specification")
		d.SetPartial("search_engine_node_count")
	}

	if (d.HasChange("table_engine_node_count") || d.HasChange("table_engine_specification")) && !d.IsNewResource() {
		newLindormSpec := d.Get("table_engine_specification")
		newLindormNum := d.Get("table_engine_node_count")
		enabled := d.Get("enabled_table_engine").(bool)
		currentInstanceStorage := formatInt(d.Get("instance_storage"))
		if !enabled {
			upgradeLindormInstanceTableReq := map[string]interface{}{}
			upgradeLindormInstanceTableReq["UpgradeType"] = "open-lindorm-engine"
			upgradeLindormInstanceTableReq["LindormSpec"] = newLindormSpec
			upgradeLindormInstanceTableReq["LindormNum"] = newLindormNum
			upgradeLindormInstanceTableReq["ClusterStorage"] = currentInstanceStorage
			err := UpgradeLindormInstance(d, meta, upgradeLindormInstanceTableReq)
			if err != nil {
				return WrapError(err)
			}
		}

		if enabled && d.HasChange("table_engine_specification") {
			upgradeLindormInstanceTableReq := map[string]interface{}{}
			upgradeLindormInstanceTableReq["UpgradeType"] = "upgrade-lindorm-engine"
			upgradeLindormInstanceTableReq["LindormSpec"] = newLindormSpec
			upgradeLindormInstanceTableReq["ClusterStorage"] = currentInstanceStorage
			err := UpgradeLindormInstance(d, meta, upgradeLindormInstanceTableReq)
			if err != nil {
				return WrapError(err)
			}
		}

		if enabled && d.HasChange("table_engine_node_count") {
			if !d.IsNewResource() && d.HasChange("log_num") && d.HasChange("table_engine_node_count") {
				upgradeLindormLogNumReq := map[string]interface{}{}
				upgradeLindormLogNumReq["UpgradeType"] = "upgrade-lindorm-core-num"
				upgradeLindormLogNumReq["LogNum"] = d.Get("log_num")
				upgradeLindormLogNumReq["LindormNum"] = d.Get("table_engine_node_count")
				err := UpgradeLindormInstance(d, meta, upgradeLindormLogNumReq)
				if err != nil {
					return WrapError(err)
				}
			} else {
				upgradeLindormInstanceTableNumReq := map[string]interface{}{}
				upgradeLindormInstanceTableNumReq["UpgradeType"] = "upgrade-lindorm-core-num"
				upgradeLindormInstanceTableNumReq["LindormNum"] = newLindormNum
				upgradeLindormInstanceTableNumReq["ClusterStorage"] = currentInstanceStorage
				err := UpgradeLindormInstance(d, meta, upgradeLindormInstanceTableNumReq)
				if err != nil {
					return WrapError(err)
				}
			}
		}

		d.SetPartial("table_engine_specification")
		d.SetPartial("table_engine_node_count")
	}

	if (d.HasChange("time_series_engine_node_count") || d.HasChange("time_serires_engine_specification") || d.HasChange("time_series_engine_specification")) && !d.IsNewResource() {
		var newTsdbSpec interface{}

		if d.HasChange("time_serires_engine_specification") {
			newTsdbSpec = d.Get("time_serires_engine_specification")
		}

		if d.HasChange("time_series_engine_specification") {
			newTsdbSpec = d.Get("time_series_engine_specification")
		}

		newTsdbNum := d.Get("time_series_engine_node_count")
		enabled := d.Get("enabled_time_serires_engine").(bool)
		currentInstanceStorage := formatInt(d.Get("instance_storage"))
		if !enabled {
			upgradeLindormInstanceSearchReq := map[string]interface{}{}
			upgradeLindormInstanceSearchReq["UpgradeType"] = "open-tsdb-engine"
			upgradeLindormInstanceSearchReq["TsdbSpec"] = newTsdbSpec
			upgradeLindormInstanceSearchReq["TsdbNum"] = newTsdbNum
			upgradeLindormInstanceSearchReq["ClusterStorage"] = currentInstanceStorage
			err := UpgradeLindormInstance(d, meta, upgradeLindormInstanceSearchReq)
			if err != nil {
				return WrapError(err)
			}
		}

		if enabled && (d.HasChange("time_serires_engine_specification") || d.HasChange("time_series_engine_specification")) {
			upgradeLindormInstanceSearchReq := map[string]interface{}{}
			upgradeLindormInstanceSearchReq["UpgradeType"] = "upgrade-tsdb-engine"
			upgradeLindormInstanceSearchReq["TsdbSpec"] = newTsdbSpec
			upgradeLindormInstanceSearchReq["ClusterStorage"] = currentInstanceStorage
			err := UpgradeLindormInstance(d, meta, upgradeLindormInstanceSearchReq)
			if err != nil {
				return WrapError(err)
			}
		}

		if enabled && d.HasChange("time_series_engine_node_count") {
			upgradeLindormInstanceSearchNumReq := map[string]interface{}{}
			upgradeLindormInstanceSearchNumReq["UpgradeType"] = "upgrade-tsdb-core-num"
			upgradeLindormInstanceSearchNumReq["TsdbNum"] = newTsdbNum
			upgradeLindormInstanceSearchNumReq["ClusterStorage"] = currentInstanceStorage
			err := UpgradeLindormInstance(d, meta, upgradeLindormInstanceSearchNumReq)
			if err != nil {
				return WrapError(err)
			}
		}

		d.SetPartial("time_serires_engine_specification")
		d.SetPartial("time_series_engine_specification")
		d.SetPartial("time_series_engine_node_count")
	}

	if d.HasChange("lts_node_count") || d.HasChange("lts_node_specification") {
		newLtsCoreSpec := d.Get("lts_node_specification")
		newLtsCoreNum := d.Get("lts_node_count")
		enabled := d.Get("enabled_lts_engine").(bool)
		currentInstanceStorage := formatInt(d.Get("instance_storage"))
		if !enabled {
			upgradeLindormInstanceLtsReq := map[string]interface{}{}
			upgradeLindormInstanceLtsReq["UpgradeType"] = "open-bds-transfer-only"
			upgradeLindormInstanceLtsReq["LtsCoreSpec"] = newLtsCoreSpec
			upgradeLindormInstanceLtsReq["LtsCoreNum"] = newLtsCoreNum
			upgradeLindormInstanceLtsReq["ClusterStorage"] = currentInstanceStorage
			err := UpgradeLindormInstance(d, meta, upgradeLindormInstanceLtsReq)
			if err != nil {
				return WrapError(err)
			}
		}

		if enabled && d.HasChange("lts_node_specification") {
			upgradeLindormInstanceLtsReq := map[string]interface{}{}
			upgradeLindormInstanceLtsReq["UpgradeType"] = "upgrade-bds-transfer"
			upgradeLindormInstanceLtsReq["LtsCoreSpec"] = newLtsCoreSpec
			upgradeLindormInstanceLtsReq["ClusterStorage"] = currentInstanceStorage
			err := UpgradeLindormInstance(d, meta, upgradeLindormInstanceLtsReq)
			if err != nil {
				return WrapError(err)
			}
		}

		if enabled && d.HasChange("lts_node_count") {
			upgradeLindormInstanceLtsNumReq := map[string]interface{}{}
			upgradeLindormInstanceLtsNumReq["UpgradeType"] = "upgrade-bds-core-num"
			upgradeLindormInstanceLtsNumReq["LtsCoreNum"] = newLtsCoreNum
			upgradeLindormInstanceLtsNumReq["ClusterStorage"] = currentInstanceStorage
			err := UpgradeLindormInstance(d, meta, upgradeLindormInstanceLtsNumReq)
			if err != nil {
				return WrapError(err)
			}
		}

		d.SetPartial("lts_node_specification")
		d.SetPartial("lts_node_count")
	}

	if (d.HasChange("stream_engine_node_count") || d.HasChange("stream_engine_specification")) && !d.IsNewResource() {
		newStreamCoreNum := d.Get("stream_engine_node_count")
		newStreamCoreSpec := d.Get("stream_engine_specification")
		enabled := d.Get("enabled_stream_engine").(bool)
		currentInstanceStorage := formatInt(d.Get("instance_storage"))

		if !enabled {
			openStreamEngineReq := map[string]interface{}{}
			openStreamEngineReq["UpgradeType"] = "open-stream-engine"
			openStreamEngineReq["StreamNum"] = newStreamCoreNum
			openStreamEngineReq["StreamSpec"] = newStreamCoreSpec
			openStreamEngineReq["ClusterStorage"] = currentInstanceStorage

			err := UpgradeLindormInstance(d, meta, openStreamEngineReq)
			if err != nil {
				return WrapError(err)
			}
		}

		if enabled && d.HasChange("stream_engine_node_count") {
			upgradeStreamCoreNumReq := map[string]interface{}{}
			upgradeStreamCoreNumReq["UpgradeType"] = "upgrade-stream-core-num"
			upgradeStreamCoreNumReq["StreamNum"] = newStreamCoreNum
			upgradeStreamCoreNumReq["ClusterStorage"] = currentInstanceStorage

			err := UpgradeLindormInstance(d, meta, upgradeStreamCoreNumReq)
			if err != nil {
				return WrapError(err)
			}
		}

		if enabled && d.HasChange("stream_engine_specification") {
			upgradeStreamEngineReq := map[string]interface{}{}
			upgradeStreamEngineReq["UpgradeType"] = "upgrade-stream-engine"
			upgradeStreamEngineReq["StreamSpec"] = newStreamCoreSpec

			err := UpgradeLindormInstance(d, meta, upgradeStreamEngineReq)
			if err != nil {
				return WrapError(err)
			}
		}

		d.SetPartial("stream_engine_node_count")
		d.SetPartial("stream_engine_specification")
	}

	update = false
	upgradeLindormInstanceClusterStorageReq := map[string]interface{}{
		"UpgradeType": "upgrade-disk-size",
	}
	if d.HasChange("instance_storage") && !d.IsNewResource() {
		object, err := hitsdbService.DescribeLindormInstance(d.Id())
		if err != nil {
			return WrapError(err)
		}

		currentInstanceStorage := fmt.Sprint(object["InstanceStorage"])
		chanageInstanceStorage := fmt.Sprint(d.Get("instance_storage"))

		if currentInstanceStorage != chanageInstanceStorage {
			update = true
			upgradeLindormInstanceClusterStorageReq["ClusterStorage"] = chanageInstanceStorage
		}
	}
	if update {
		err := UpgradeLindormInstance(d, meta, upgradeLindormInstanceClusterStorageReq)
		if err != nil {
			return WrapError(err)
		}
		d.SetPartial("instance_storage")
	}

	d.Partial(false)

	return resourceAliCloudLindormInstanceRead(d, meta)
}

func resourceAliCloudLindormInstanceDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	hitsdbService := HitsdbService{client}
	action := "ReleaseLindormInstance"
	var response map[string]interface{}
	var err error
	object, err := hitsdbService.DescribeLindormInstance(d.Id())
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	if fmt.Sprint(object["PayType"]) == "PREPAY" {
		log.Printf("[WARN] Cannot destroy resource alicloud_lindorm_instance. Terraform will remove this resource from the state file, however resources may remain.")
		return nil
	}

	request := map[string]interface{}{
		"InstanceId": d.Id(),
	}

	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutDelete)), func() *resource.RetryError {
		response, err = client.RpcPost("hitsdb", "2020-06-15", action, nil, request, false)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)

	if err != nil {
		if IsExpectedErrors(err, []string{"Lindorm.Errorcode.InstanceNotFound", "Instance.IsDeleted"}) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	return nil
}

func convertLindormInstancePaymentTypeRequest(source interface{}) interface{} {
	switch source {
	case "PayAsYouGo":
		return "POSTPAY"
	case "Subscription":
		return "PREPAY"
	}
	return source
}

func convertLindormInstancePaymentTypeResponse(source interface{}) interface{} {
	switch source {
	case "POSTPAY":
		return "PayAsYouGo"
	case "PREPAY":
		return "Subscription"
	}
	return source
}

func UpgradeLindormInstance(d *schema.ResourceData, meta interface{}, request map[string]interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	var err error
	request["InstanceId"] = d.Id()
	request["RegionId"] = client.RegionId
	request["ZoneId"] = d.Get("zone_id")

	action := "UpgradeLindormInstance"
	wait := incrementalWait(3*time.Second, 30*time.Second)
	err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *resource.RetryError {
		response, err = client.RpcPost("hitsdb", "2020-06-15", action, nil, request, false)
		if err != nil {
			if IsExpectedErrors(err, []string{"Instance.NotActive", "OperationDenied.OrderProcessing"}) || NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	hitsdbService := HitsdbService{client}
	stateConf := BuildStateConf([]string{}, []string{"ACTIVATION"}, d.Timeout(schema.TimeoutUpdate), 1*time.Minute, hitsdbService.LindormInstanceStateRefreshFunc(d.Id(), []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return nil
}
