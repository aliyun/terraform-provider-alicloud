package alicloud

import (
	"fmt"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"log"
	"strings"
	"time"
)

func resourceAlicloudGpdbInstance() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudGpdbDbInstanceCreate,
		Read:   resourceAlicloudGpdbDbInstanceRead,
		Update: resourceAlicloudGpdbDbInstanceUpdate,
		Delete: resourceAlicloudGpdbDbInstanceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(60 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"db_instance_category": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"db_instance_class": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"db_instance_mode": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"StorageElastic", "Serverless", "Classic"}, false),
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"engine": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"engine_version": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"ip_whitelist": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ip_group_attribute": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"ip_group_name": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"security_ip_list": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
				ConflictsWith: []string{"security_ip_list"},
			},
			"security_ip_list": {
				Type:          schema.TypeSet,
				Elem:          &schema.Schema{Type: schema.TypeString},
				Optional:      true,
				Deprecated:    "Field 'security_ip_list' has been deprecated from version 1.187.0. Use 'ip_whitelist' instead.",
				ConflictsWith: []string{"ip_whitelist"},
			},
			"instance_network_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Computed:     true,
				ValidateFunc: validation.StringInSlice([]string{"VPC"}, false),
			},
			"instance_spec": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"2C16G", "4C32G", "16C128G", "2C8G", "4C16G", "8C32G", "16C64G"}, false),
			},
			"maintain_end_time": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"maintain_start_time": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"master_node_num": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntInSlice([]int{1, 2}),
				Default:      1,
			},
			"instance_group_count": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
			},
			"payment_type": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"PayAsYouGo", "Subscription"}, false),
			},
			"instance_charge_type": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ForceNew:      true,
				Deprecated:    "Field 'instance_charge_type' has been deprecated from version 1.187.0. Use 'payment_type' instead.",
				ValidateFunc:  validation.StringInSlice([]string{"Prepaid", "PostPaid"}, false),
				ConflictsWith: []string{"payment_type"},
			},
			"period": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"Month", "Week"}, false),
			},
			"private_ip_address": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"resource_group_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"seg_node_num": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.IntBetween(2, 512),
			},
			"seg_storage_type": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"storage_size": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"tags": tagsSchema(),
			"used_time": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"vswitch_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"vpc_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"zone_id": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ForceNew:      true,
				ConflictsWith: []string{"availability_zone"},
			},
			"availability_zone": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ForceNew:      true,
				Deprecated:    "Field 'availability_zone' has been deprecated from version 1.187.0. Use 'zone_id' instead.",
				ConflictsWith: []string{"zone_id"},
			},
			"create_sample_data": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
		},
	}
}

func resourceAlicloudGpdbDbInstanceCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "CreateDBInstance"
	request := make(map[string]interface{})
	conn, err := client.NewGpdbClient()
	if err != nil {
		return WrapError(err)
	}
	if v, ok := d.GetOk("db_instance_category"); ok {
		request["DBInstanceCategory"] = v
	}
	if v, ok := d.GetOk("db_instance_class"); ok {
		request["DBInstanceClass"] = v
	}
	if v, ok := d.GetOk("db_instance_mode"); ok {
		request["DBInstanceMode"] = v
	}
	if v, ok := d.GetOk("description"); ok {
		request["DBInstanceDescription"] = v
	}
	request["Engine"] = d.Get("engine")
	request["EngineVersion"] = d.Get("engine_version")

	if v, ok := d.GetOk("zone_id"); ok {
		request["ZoneId"] = v
	} else if v, ok := d.GetOk("availability_zone"); ok {
		request["ZoneId"] = v
	}

	if v, ok := d.GetOk("instance_network_type"); ok {
		request["InstanceNetworkType"] = v
	}
	if v, ok := d.GetOk("instance_spec"); ok {
		request["InstanceSpec"] = v
	}
	if v, ok := d.GetOk("master_node_num"); ok {
		request["MasterNodeNum"] = v
	}
	if v, ok := d.GetOk("instance_group_count"); ok {
		request["DBInstanceGroupCount"] = v
	}

	if v, ok := d.GetOk("payment_type"); ok {
		request["PayType"] = convertGpdbDbInstancePaymentTypeRequest(v.(string))
	} else if v, ok := d.GetOk("instance_charge_type"); ok {
		request["PayType"] = v
	}

	if v, ok := d.GetOk("period"); ok {
		request["Period"] = v
	}
	if v, ok := d.GetOk("private_ip_address"); ok {
		request["PrivateIpAddress"] = v
	}
	request["RegionId"] = client.RegionId
	if v, ok := d.GetOk("resource_group_id"); ok {
		request["ResourceGroupId"] = v
	}
	if v, ok := d.GetOk("seg_node_num"); ok {
		request["SegNodeNum"] = v
	}
	if v, ok := d.GetOk("seg_storage_type"); ok {
		request["SegStorageType"] = v
	}
	if v, ok := d.GetOk("storage_size"); ok {
		request["StorageSize"] = v
	}
	if v, ok := d.GetOk("used_time"); ok {
		request["UsedTime"] = v
	}

	if v, ok := d.GetOk("vpc_id"); ok {
		request["VPCId"] = v
	}

	if v, ok := d.GetOk("vswitch_id"); ok {
		request["VSwitchId"] = v
	}

	if v, ok := d.GetOk("create_sample_data"); ok {
		request["CreateSampleData"] = v
	}

	if request["VpcId"] == nil && request["VSwitchId"] != nil {
		vpcService := VpcService{client}
		vsw, err := vpcService.DescribeVSwitchWithTeadsl(request["VSwitchId"].(string))
		if err != nil {
			return WrapError(err)
		}
		if v, ok := request["VPCId"].(string); !ok || v == "" {
			request["VPCId"] = vsw["VpcId"]
		}
	}

	request["SecurityIPList"] = LOCAL_HOST_IP

	if v, ok := d.GetOk("ip_whitelist"); ok {
		for _, iPWhitelist := range v.(*schema.Set).List() {
			iPWhitelistArg := iPWhitelist.(map[string]interface{})
			request["IPGroupAttribute"] = iPWhitelistArg["ip_group_attribute"]
			request["IPGroupName"] = iPWhitelistArg["ip_group_name"]
			request["SecurityIPList"] = iPWhitelistArg["security_ip_list"]
		}
	} else {
		if len(d.Get("security_ip_list").(*schema.Set).List()) > 0 {
			request["SecurityIpList"] = strings.Join(expandStringList(d.Get("security_ip_list").(*schema.Set).List())[:], COMMA_SEPARATED)
		}
	}

	request["ClientToken"] = buildClientToken("CreateDBInstance")
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutCreate)), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2016-05-03"), StringPointer("AK"), nil, request, &runtime)
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_gpdb_instance", action, AlibabaCloudSdkGoERROR)
	}
	d.SetId(fmt.Sprint(response["DBInstanceId"]))
	gpdbService := GpdbService{client}
	stateConf := BuildStateConf([]string{}, []string{"Running"}, client.GetRetryTimeout(d.Timeout(schema.TimeoutCreate)), 5*time.Second, gpdbService.GpdbDbInstanceStateRefreshFunc(d.Id(), []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return resourceAlicloudGpdbDbInstanceUpdate(d, meta)
}
func resourceAlicloudGpdbDbInstanceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	gpdbService := GpdbService{client}
	object, err := gpdbService.DescribeGpdbDbInstance(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_gpdb_db_instance gpdbService.DescribeGpdbDbInstance Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	d.Set("db_instance_category", object["DBInstanceCategory"])
	d.Set("db_instance_mode", object["DBInstanceMode"])
	d.Set("description", object["DBInstanceDescription"])
	d.Set("engine", object["Engine"])
	d.Set("engine_version", object["EngineVersion"])
	d.Set("instance_network_type", object["InstanceNetworkType"])
	d.Set("maintain_end_time", object["MaintainEndTime"])
	d.Set("maintain_start_time", object["MaintainStartTime"])
	d.Set("master_node_num", formatInt(object["MasterNodeNum"]))
	if v, ok := object["SegmentCounts"]; ok && fmt.Sprint(v) != "0" {
		d.Set("node_num", formatInt(v))
	}
	d.Set("payment_type", convertGpdbDbInstancePaymentTypeResponse(object["PayType"]))
	d.Set("instance_charge_type", convertGpdbDbInstancePaymentTypeResponse(object["PayType"]))
	d.Set("seg_node_num", formatInt(object["SegNodeNum"]))
	d.Set("status", object["DBInstanceStatus"])
	if v, ok := object["StorageSize"]; ok && fmt.Sprint(v) != "0" {
		d.Set("storage_size", formatInt(v))
	}
	if v, ok := object["Tags"].(map[string]interface{}); ok {
		d.Set("tags", tagsToMap(v["Tag"]))
	}
	d.Set("vswitch_id", object["VSwitchId"])
	d.Set("vpc_id", object["VpcId"])
	d.Set("zone_id", object["ZoneId"])
	d.Set("availability_zone", object["ZoneId"])
	describeDBInstanceIPArrayListObject, err := gpdbService.DescribeDBInstanceIPArrayList(d.Id())
	if err != nil {
		return WrapError(err)
	}
	if iPWhitelistMap, ok := describeDBInstanceIPArrayListObject["Items"].(map[string]interface{}); ok && iPWhitelistMap != nil {
		if dBInstanceIPArrayList, ok := iPWhitelistMap["DBInstanceIPArray"]; ok && dBInstanceIPArrayList != nil {
			iPWhitelistMaps := make([]map[string]interface{}, 0)
			for _, dBInstanceIPArrayListItem := range dBInstanceIPArrayList.([]interface{}) {
				if dBInstanceIPArrayListItemMap, ok := dBInstanceIPArrayListItem.(map[string]interface{}); ok {
					if dBInstanceIPArrayListItem.(map[string]interface{})["DBInstanceIPArrayAttribute"] == "hidden" {
						continue
					}
					dBInstanceIPArrayListMap := map[string]interface{}{}
					dBInstanceIPArrayListMap["ip_group_attribute"] = dBInstanceIPArrayListItemMap["DBInstanceIPArrayAttribute"]
					dBInstanceIPArrayListMap["ip_group_name"] = dBInstanceIPArrayListItemMap["DBInstanceIPArrayName"]
					dBInstanceIPArrayListMap["security_ip_list"] = dBInstanceIPArrayListItemMap["SecurityIPList"]
					iPWhitelistMaps = append(iPWhitelistMaps, dBInstanceIPArrayListMap)
				}
			}
			d.Set("ip_whitelist", iPWhitelistMaps)
		}
	}
	return nil
}
func resourceAlicloudGpdbDbInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	gpdbService := GpdbService{client}
	conn, err := client.NewGpdbClient()
	if err != nil {
		return WrapError(err)
	}
	var response map[string]interface{}
	request := make(map[string]interface{})
	d.Partial(true)

	if d.HasChange("tags") {
		if err := gpdbService.SetResourceTags(d, "ALIYUN::GPDB::INSTANCE"); err != nil {
			return WrapError(err)
		}
		d.SetPartial("tags")
	}
	update := false
	request = map[string]interface{}{
		"DBInstanceId": d.Id(),
	}
	if !d.IsNewResource() && d.HasChange("description") {
		update = true
	}
	if v, ok := d.GetOk("description"); ok {
		request["DBInstanceDescription"] = v
	}
	if update {
		action := "ModifyDBInstanceDescription"
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2016-05-03"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
		d.SetPartial("description")
	}
	update = false
	request = map[string]interface{}{
		"DBInstanceId": d.Id(),
	}
	if !d.IsNewResource() && d.HasChange("resource_group_id") {
		update = true
	}
	if v, ok := d.GetOk("resource_group_id"); ok {
		request["NewResourceGroupId"] = v
	}
	if update {
		action := "ModifyDBInstanceResourceGroup"
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2016-05-03"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
		d.SetPartial("resource_group_id")
	}
	update = false
	request = map[string]interface{}{
		"DBInstanceId": d.Id(),
	}
	if !d.IsNewResource() && d.HasChange("maintain_end_time") {
		update = true
	}
	if v, ok := d.GetOk("maintain_end_time"); ok {
		request["EndTime"] = v
	}
	if !d.IsNewResource() && d.HasChange("maintain_start_time") {
		update = true
	}
	if v, ok := d.GetOk("maintain_start_time"); ok {
		request["StartTime"] = v
	}
	if update {
		action := "ModifyDBInstanceMaintainTime"
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2016-05-03"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
		d.SetPartial("maintain_end_time")
		d.SetPartial("maintain_start_time")
	}
	update = false
	request = map[string]interface{}{
		"DBInstanceId": d.Id(),
	}
	if !d.IsNewResource() && d.HasChange("ip_whitelist") {
		update = true
		if v, ok := d.GetOk("ip_whitelist"); ok {
			for _, iPWhitelist := range v.(*schema.Set).List() {
				iPWhitelistArg := iPWhitelist.(map[string]interface{})
				request["IPGroupAttribute"] = iPWhitelistArg["ip_group_attribute"]
				request["IPGroupName"] = iPWhitelistArg["ip_group_name"]
				request["SecurityIPList"] = iPWhitelistArg["security_ip_list"]
			}
		} else {
			if len(d.Get("security_ip_list").(*schema.Set).List()) > 0 {
				request["SecurityIpList"] = strings.Join(expandStringList(d.Get("security_ip_list").(*schema.Set).List())[:], COMMA_SEPARATED)
			}
		}
	}
	if update {
		action := "ModifySecurityIps"
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2016-05-03"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
		d.SetPartial("ip_whitelist")
	}
	update = false
	request = map[string]interface{}{
		"DBInstanceId": d.Id(),
	}

	request["RegionId"] = client.RegionId
	if !d.IsNewResource() && d.HasChange("seg_node_num") {
		update = true
		if v, ok := d.GetOk("seg_node_num"); ok {
			request["UpgradeType"] = 0
			request["SegNodeNum"] = v
		}
	}
	if update {
		action := "UpgradeDBInstance"
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2016-05-03"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
		stateConf := BuildStateConf([]string{}, []string{"Running"}, client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), 5*time.Second, gpdbService.GpdbDbInstanceStateRefreshFunc(d.Id(), []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
		d.SetPartial("seg_node_num")
		d.SetPartial("storage_size")
	}
	update = false
	request = map[string]interface{}{
		"DBInstanceId": d.Id(),
	}

	request["RegionId"] = client.RegionId
	if !d.IsNewResource() && d.HasChange("master_node_num") {
		update = true
		if v, ok := d.GetOk("master_node_num"); ok {
			request["UpgradeType"] = 2
			request["MasterNodeNum"] = v
		}
	}
	if update {
		action := "UpgradeDBInstance"
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2016-05-03"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
		stateConf := BuildStateConf([]string{}, []string{"Running"}, client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), 5*time.Second, gpdbService.GpdbDbInstanceStateRefreshFunc(d.Id(), []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
		d.SetPartial("master_node_num")
	}
	update = false
	request = map[string]interface{}{
		"DBInstanceId": d.Id(),
	}

	request["RegionId"] = client.RegionId
	if !d.IsNewResource() && d.HasChange("instance_spec") {
		update = true
		if v, ok := d.GetOk("instance_spec"); ok {
			request["UpgradeType"] = 1
			request["InstanceSpec"] = v
		}
	}
	if update {
		action := "UpgradeDBInstance"
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2016-05-03"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
		stateConf := BuildStateConf([]string{}, []string{"Running"}, client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), 5*time.Second, gpdbService.GpdbDbInstanceStateRefreshFunc(d.Id(), []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
		d.SetPartial("instance_spec")
	}
	update = false
	request = map[string]interface{}{
		"DBInstanceId": d.Id(),
	}

	request["RegionId"] = client.RegionId
	if !d.IsNewResource() && d.HasChange("storage_size") {
		update = true
		if v, ok := d.GetOk("storage_size"); ok {
			request["UpgradeType"] = 1
			request["StorageSize"] = v
		}
	}
	if update {
		action := "UpgradeDBInstance"
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2016-05-03"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
		stateConf := BuildStateConf([]string{}, []string{"Running"}, client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), 5*time.Second, gpdbService.GpdbDbInstanceStateRefreshFunc(d.Id(), []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
		d.SetPartial("storage_size")
	}
	d.Partial(false)
	return resourceAlicloudGpdbDbInstanceRead(d, meta)
}
func resourceAlicloudGpdbDbInstanceDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	action := "DeleteDBInstance"
	var response map[string]interface{}
	conn, err := client.NewGpdbClient()
	if err != nil {
		return WrapError(err)
	}
	request := map[string]interface{}{
		"DBInstanceId": d.Id(),
	}
	if v, ok := d.GetOk("payment_type"); ok && v.(string) == "Subscription" {
		log.Printf("[WARN] Cannot destroy resourceGpdbDbInstance. Because payment_type = 'Subscription'. Terraform will remove this resource from the state file, however resources may remain.")
		return nil
	}
	request["ClientToken"] = buildClientToken("DeleteDBInstance")
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutDelete)), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2016-05-03"), StringPointer("AK"), nil, request, &runtime)
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
		if IsExpectedErrors(err, []string{"InvalidDBInstanceId.NotFound", "OperationDenied.DBInstancePayType"}) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}
	return nil
}
func convertGpdbDbInstancePaymentTypeRequest(source interface{}) interface{} {
	switch source {
	case "Subscription":
		return "Prepaid"
	case "PayAsYouGo":
		return "Postpaid"
	}
	return source
}
func convertGpdbDbInstancePaymentTypeResponse(source interface{}) interface{} {
	switch source {
	case "Prepaid":
		return "Subscription"
	case "Postpaid":
		return "PayAsYouGo"
	}

	return source
}
