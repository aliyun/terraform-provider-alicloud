package alicloud

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAliCloudSelectDBDbInstance() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudSelectDBDbInstanceCreate,
		Read:   resourceAliCloudSelectDBDbInstanceRead,
		Update: resourceAliCloudSelectDBDbInstanceUpdate,
		Delete: resourceAliCloudSelectDBDbInstanceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"db_instance_description": {
				Type:     schema.TypeString,
				Required: true,
			},
			"db_instance_class": {
				Type:     schema.TypeString,
				Required: true,
			},
			"payment_type": {
				Type:         schema.TypeString,
				ValidateFunc: StringInSlice([]string{"PayAsYouGo", "Subscription"}, false),
				Required:     true,
			},
			"period": {
				Type:             schema.TypeString,
				ValidateFunc:     StringInSlice([]string{string(Year), string(Month)}, false),
				Optional:         true,
				DiffSuppressFunc: selectdbPostPaidDiffSuppressFunc,
			},
			"period_time": {
				Type:             schema.TypeInt,
				ValidateFunc:     IntInSlice([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 12, 24, 36}),
				Optional:         true,
				DiffSuppressFunc: selectdbPostPaidDiffSuppressFunc,
			},
			"zone_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"vpc_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"vswitch_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"cache_size": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"tags": tagsSchema(),
			// flag for public network and update
			"enable_public_network": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"upgraded_engine_minor_version": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"desired_security_ip_lists": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"group_name": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"security_ip_list": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},

			// Computed values
			"region_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"engine": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"engine_minor_version": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"cpu_prepaid": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"memory_prepaid": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"cache_size_prepaid": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"cluster_count_prepaid": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"cpu_postpaid": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"memory_postpaid": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"cache_size_postpaid": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"cluster_count_postpaid": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"sub_domain": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"gmt_created": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"gmt_modified": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"gmt_expired": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"lock_mode": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"lock_reason": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"instance_net_infos": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"db_ip": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"vpc_instance_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"connection_string": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"net_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"vswitch_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"port_list": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"port": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"protocol": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
					},
				},
			},
			"security_ip_lists": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"group_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"group_tag": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"security_ip_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"security_ip_list": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"list_net_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func resourceAliCloudSelectDBDbInstanceCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	selectDBService := SelectDBService{client}
	request, err := buildSelectDBCreateInstanceRequest(d, meta)
	if err != nil {
		return WrapError(err)
	}
	action := "CreateDBInstance"
	response, err := selectDBService.RequestProcessForSelectDB(request, action, "POST")
	if err != nil {
		return WrapError(err)
	}
	if resp, err := jsonpath.Get("$.Data", response); err != nil || resp == nil {
		return WrapErrorf(err, IdMsg, "alicloud_selectdb_db_instance")
	} else {
		instanceId := resp.(map[string]interface{})["DBInstanceId"].(string)
		d.SetId(fmt.Sprint(instanceId))
	}

	stateConfPreparing := BuildStateConf([]string{"RESOURCE_PREPARING", "CREATING"}, []string{"ACTIVE"}, d.Timeout(schema.TimeoutCreate), 1*time.Minute, selectDBService.SelectDBDbInstanceStateRefreshFunc(d.Id(), []string{"DELETING"}))
	if _, err := stateConfPreparing.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	enable_public_network, exist := d.GetOkExists("enable_public_network")
	if exist {
		if enable_public_network == true {
			if _, err := selectDBService.AllocateSelectDBInstancePublicConnection(d.Id()); err != nil {
				return WrapError(err)
			}
			stateConf := BuildStateConf([]string{"NET_CREATING"}, []string{"ACTIVE"}, d.Timeout(schema.TimeoutCreate), 1*time.Minute, selectDBService.SelectDBDbInstanceStateRefreshFunc(d.Id(), []string{"DELETING"}))
			if _, err := stateConf.WaitForState(); err != nil {
				return WrapErrorf(err, IdMsg, d.Id())
			}
		}
	}

	return resourceAliCloudSelectDBDbInstanceUpdate(d, meta)
}

func resourceAliCloudSelectDBDbInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	selectDBService := SelectDBService{client}
	d.Partial(true)

	if !d.IsNewResource() && d.HasChange("payment_type") {
		_, newPayment := d.GetChange("payment_type")

		request := make(map[string]string)
		payment := convertPaymentTypeToChargeType(newPayment)

		if payment == string(PostPaid) {
			request["payment_type"] = string("POSTPAY")
		} else if payment == string(PrePaid) {
			request["payment_type"] = string("PREPAY")
			if v, ok := d.GetOk("period"); ok {
				request["period"] = v.(string)
			} else {
				request["period"] = "Month"
			}
			if v, ok := d.GetOk("period_time"); ok && v.(int) > 0 {
				request["period_time"] = fmt.Sprint(v.(int))
			} else {
				request["period_time"] = "1"
			}
		}
		if _, err := selectDBService.ModifySelectDBInstancePaymentType(d.Id(), request); err != nil {
			return WrapError(err)
		}
		d.SetPartial("payment_type")

	}

	if !d.IsNewResource() && d.HasChange("enable_public_network") {
		oldNetStatus, newNetStatus := d.GetChange("enable_public_network")
		if oldNetStatus.(bool) && !newNetStatus.(bool) {
			netResp, err := selectDBService.DescribeSelectDBDbInstanceNetInfo(d.Id())
			if err != nil {
				return WrapError(err)
			}
			connectionString := ""
			resultNet, _ := netResp["DBInstanceNetInfos"].([]interface{})
			for _, v := range resultNet {
				item := v.(map[string]interface{})
				if item["NetType"].(string) == "PUBLIC" && !strings.Contains(item["ConnectionString"].(string), "webui") {
					connectionString = item["ConnectionString"].(string)
				}
			}

			if _, err := selectDBService.ReleaseSelectDBInstancePublicConnection(d.Id(), connectionString); err != nil {
				return WrapError(err)
			}
			stateConf := BuildStateConf([]string{"NET_DELETING"}, []string{"ACTIVE"}, d.Timeout(schema.TimeoutUpdate), 1*time.Minute, selectDBService.SelectDBDbInstanceStateRefreshFunc(d.Id(), []string{"DELETING"}))
			if _, err := stateConf.WaitForState(); err != nil {
				return WrapErrorf(err, IdMsg, d.Id())
			}
			d.SetPartial("enable_public_network")
			d.SetPartial("instance_net_infos")
		} else if !oldNetStatus.(bool) && newNetStatus.(bool) {
			if _, err := selectDBService.AllocateSelectDBInstancePublicConnection(d.Id()); err != nil {
				return WrapError(err)
			}
			stateConf := BuildStateConf([]string{"NET_CREATING"}, []string{"ACTIVE"}, d.Timeout(schema.TimeoutUpdate), 1*time.Minute, selectDBService.SelectDBDbInstanceStateRefreshFunc(d.Id(), []string{"DELETING"}))
			if _, err := stateConf.WaitForState(); err != nil {
				return WrapErrorf(err, IdMsg, d.Id())
			}
			d.SetPartial("enable_public_network")
			d.SetPartial("instance_net_infos")
		}

	}

	if !d.IsNewResource() && d.HasChange("upgraded_engine_minor_version") && d.Get("upgraded_engine_minor_version") != "" {
		_, newVersion := d.GetChange("upgraded_engine_minor_version")
		instanceId := fmt.Sprint(d.Id())
		instanceResp, err := selectDBService.DescribeSelectDBDbInstance(instanceId)
		if err != nil {
			return WrapError(err)
		}
		upgradeTargetVersion := ""
		canUpgradeVersion := instanceResp["CanUpgradeVersions"].([]interface{})
		for _, version := range canUpgradeVersion {
			if newVersion.(string) == version.(string) {
				upgradeTargetVersion = newVersion.(string)
				break
			}
		}
		if upgradeTargetVersion == "" {
			return WrapErrorf(err, "Invalid upgrade version for %s, cannot upgrade to %s", d.Id(), newVersion.(string), AlibabaCloudSdkGoERROR)
		}

		// todo maintaintime update
		_, err = selectDBService.UpgradeSelectDBInstanceEngineVersion(d.Id(), upgradeTargetVersion, false)
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), "UpgradeSelectDBInstanceEngineVersion", AlibabaCloudSdkGoERROR)
		}
		stateConf := BuildStateConf([]string{"MODULE_UPGRADING"}, []string{"ACTIVE"}, d.Timeout(schema.TimeoutUpdate), 10*time.Second, selectDBService.SelectDBDbInstanceStateRefreshFunc(d.Id(), []string{"DELETING"}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
		d.SetPartial("upgraded_engine_minor_version")

	}

	cacheSizeModified := false
	if !d.IsNewResource() && (d.HasChange("db_instance_class")) {
		defaultBeId := fmt.Sprint(d.Id() + ":" + d.Id() + "-be")
		_, newClass := d.GetChange("db_instance_class")
		cache_size := 0
		if d.HasChange("cache_size") {
			_, newCacheSize := d.GetChange("cache_size")
			cache_size = newCacheSize.(int)
			cacheSizeModified = true
		}
		_, err := selectDBService.ModifySelectDBCluster(defaultBeId, newClass.(string), cache_size)
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, defaultBeId, "ModifyDBCluster", AlibabaCloudSdkGoERROR)
		}
		stateConf := BuildStateConf([]string{"RESOURCE_PREPARING", "CLASS_CHANGING"}, []string{"ACTIVATION"}, d.Timeout(schema.TimeoutUpdate), 1*time.Minute, selectDBService.SelectDBDbClusterStateRefreshFunc(defaultBeId, []string{"DELETING"}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, defaultBeId)
		}
		d.SetPartial("db_instance_class")
		if d.HasChange("cache_size") {
			d.SetPartial("cache_size")
		}
	}

	if !d.IsNewResource() && d.HasChange("cache_size") && !cacheSizeModified {
		defaultBeId := fmt.Sprint(d.Id() + ":" + d.Id() + "-be")
		_, newCacheSize := d.GetChange("cache_size")
		db_cluster_class := d.Get("db_instance_class").(string)
		if d.HasChange("db_instance_class") {
			_, newClass := d.GetChange("db_instance_class")
			db_cluster_class = newClass.(string)
		}
		_, err := selectDBService.ModifySelectDBCluster(defaultBeId, db_cluster_class, newCacheSize.(int))
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, defaultBeId, "ModifyDBCluster", AlibabaCloudSdkGoERROR)
		}
		stateConf := BuildStateConf([]string{"RESOURCE_PREPARING", "CLASS_CHANGING"}, []string{"ACTIVATION"}, d.Timeout(schema.TimeoutUpdate), 1*time.Minute, selectDBService.SelectDBDbClusterStateRefreshFunc(defaultBeId, []string{"DELETING"}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, defaultBeId)
		}
		d.SetPartial("cache_size")
		if d.HasChange("db_instance_class") {
			d.SetPartial("db_instance_class")
		}
	}

	if !d.IsNewResource() && d.HasChange("db_instance_description") {
		_, newDesc := d.GetChange("db_instance_description")
		_, err := selectDBService.ModifySelectDBInstanceDescription(d.Id(), newDesc.(string))
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), "ModifySelectDBInstanceDescription", AlibabaCloudSdkGoERROR)
		}
		d.SetPartial("db_instance_description")
	}

	if d.HasChange("tags") {
		added, removed := parsingTags(d)
		if err := selectDBService.SetResourceTags(d.Id(), added, removed); err != nil {
			return WrapError(err)
		}
		d.SetPartial("tags")
	}

	if d.HasChange("desired_security_ip_lists") {
		_, newDesc := d.GetChange("desired_security_ip_lists")
		for _, v := range newDesc.([]interface{}) {
			item := v.(map[string]interface{})
			_, err := selectDBService.ModifySelectDBDbInstanceSecurityIPList(d.Id(), item["group_name"].(string), item["security_ip_list"].(string))
			if err != nil {
				return WrapErrorf(err, DefaultErrorMsg, d.Id(), "ModifySecurityIPList", AlibabaCloudSdkGoERROR)
			}
		}
		d.SetPartial("desired_security_ip_lists")
	}
	stateConf := BuildStateConf([]string{"RESOURCE_PREPARING", "CREATING", "CLASS_CHANGING", "MODULE_UPGRADING", "NET_CREATING", "NET_DELETING"},
		[]string{"ACTIVE"}, d.Timeout(schema.TimeoutUpdate), 10*time.Second, selectDBService.SelectDBDbInstanceStateRefreshFunc(d.Id(), []string{"DELETING"}))

	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	d.Partial(false)
	return resourceAliCloudSelectDBDbInstanceRead(d, meta)
}

func resourceAliCloudSelectDBDbInstanceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	selectDBService := SelectDBService{client}

	instanceId := fmt.Sprint(d.Id())
	instanceResp, err := selectDBService.DescribeSelectDBDbInstance(instanceId)
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("engine", instanceResp["Engine"])
	d.Set("engine_minor_version", instanceResp["EngineMinorVersion"])

	d.Set("region_id", instanceResp["RegionId"])
	d.Set("zone_id", instanceResp["ZoneId"])
	d.Set("vpc_id", instanceResp["VpcId"])

	d.Set("payment_type", convertChargeTypeToPaymentType(instanceResp["ChargeType"]))

	d.Set("db_instance_description", instanceResp["Description"])
	d.Set("status", instanceResp["Status"])
	d.Set("sub_domain", instanceResp["SubDomain"])
	d.Set("gmt_created", instanceResp["CreateTime"])
	d.Set("gmt_modified", instanceResp["GmtModified"])
	d.Set("gmt_expired", instanceResp["ExpiredTime"])
	d.Set("lock_mode", instanceResp["LockMode"])
	d.Set("lock_reason", instanceResp["LockReason"])

	default_cache_size := 0

	clusterResp := instanceResp["DBClusterList"]
	result, _ := clusterResp.([]interface{})
	defaultBeClusterId := d.Id() + "-be"
	for _, v := range result {
		item := v.(map[string]interface{})
		if item["DbClusterId"].(string) == defaultBeClusterId {
			d.Set("db_instance_class", item["DbClusterClass"])
			cache, _ := item["CacheStorageSizeGB"].(json.Number).Int64()
			default_cache_size = int(cache)
		}
	}
	d.Set("cache_size", default_cache_size)

	cpuPrepaid := 0
	cpuPostpaid := 0
	memPrepaid := 0
	memPostpaid := 0
	cachePrepaid := 0
	cachePostpaid := 0

	clusterPrepaidCount := 0
	clusterPostpaidCount := 0

	for _, v := range clusterResp.([]interface{}) {
		item := v.(map[string]interface{})
		if item["ChargeType"].(string) == "Postpaid" {
			cpuP, _ := item["CpuCores"].(json.Number).Int64()
			cpuPostpaid += int(cpuP)
			memP, _ := item["Memory"].(json.Number).Int64()
			memPostpaid += int(memP)
			cacheP, _ := item["CacheStorageSizeGB"].(json.Number).Int64()
			cachePostpaid += int(cacheP)
			clusterPostpaidCount += 1
		}
		if item["ChargeType"].(string) == "Prepaid" {
			cpuP, _ := item["CpuCores"].(json.Number).Int64()
			cpuPrepaid += int(cpuP)
			memP, _ := item["Memory"].(json.Number).Int64()
			memPrepaid += int(memP)
			cacheP, _ := item["CacheStorageSizeGB"].(json.Number).Int64()
			cachePrepaid += int(cacheP)
			clusterPrepaidCount += 1
		}
	}
	d.Set("cpu_prepaid", cpuPrepaid)
	d.Set("memory_prepaid", memPrepaid)
	d.Set("cache_size_prepaid", cachePrepaid)
	d.Set("cpu_postpaid", cpuPostpaid)
	d.Set("memory_postpaid", memPostpaid)
	d.Set("cache_size_postpaid", cachePostpaid)

	d.Set("cluster_count_prepaid", clusterPrepaidCount)
	d.Set("cluster_count_postpaid", clusterPostpaidCount)

	if resp, err := jsonpath.Get("$.Tags", instanceResp); err == nil || resp != nil {
		tags := make(map[string]interface{})
		for _, t := range resp.([]interface{}) {
			key := t.(map[string]interface{})["TagKey"].(string)
			value := t.(map[string]interface{})["TagValue"].(string)
			if !ignoredTags(key, value) {
				tags[key] = value
			}
		}
		d.Set("tags", tags)
	}

	netResp, err := selectDBService.DescribeSelectDBDbInstanceNetInfo(d.Id())
	if err != nil {
		return WrapError(err)
	}
	instanceNetArray := make([]map[string]interface{}, 0)
	resultInstanceNet, _ := netResp["DBInstanceNetInfos"].([]interface{})
	new_vswitch_id := ""
	for _, v := range resultInstanceNet {
		item := v.(map[string]interface{})
		port_list := make([]map[string]interface{}, 0)
		for _, vv := range item["PortList"].([]interface{}) {
			port_map := map[string]interface{}{
				"port":     vv.(map[string]interface{})["Port"],
				"protocol": vv.(map[string]interface{})["Protocol"],
			}
			port_list = append(port_list, port_map)
		}
		mapping := map[string]interface{}{
			"db_ip":             item["Ip"],
			"vpc_instance_id":   item["VpcInstanceId"],
			"connection_string": item["ConnectionString"],
			"net_type":          item["NetType"],
			"vswitch_id":        item["VswitchId"],
			"port_list":         port_list,
		}
		if item["VswitchId"].(string) != "" {
			new_vswitch_id = item["VswitchId"].(string)
		}
		instanceNetArray = append(instanceNetArray, mapping)
	}
	if new_vswitch_id != "" {
		d.Set("vswitch_id", new_vswitch_id)
	}
	d.Set("instance_net_infos", instanceNetArray)

	securityIpArrayList, err := selectDBService.DescribeSelectDBDbInstanceSecurityIPList(d.Id())
	if err != nil {
		return WrapError(err)
	}
	securityIpArray := make([]map[string]interface{}, 0)
	for _, v := range securityIpArrayList {
		if m1, ok := v.(map[string]interface{}); ok {
			temp1 := map[string]interface{}{
				"group_name":       m1["GroupName"],
				"group_tag":        m1["GroupTag"],
				"security_ip_type": m1["AecurityIPType"],
				"security_ip_list": m1["SecurityIPList"],
				"list_net_type":    m1["WhitelistNetType"],
			}
			securityIpArray = append(securityIpArray, temp1)
		}
	}
	d.Set("security_ip_lists", securityIpArray)

	return nil
}

func resourceAliCloudSelectDBDbInstanceDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	selectDBService := SelectDBService{client}

	stateConf := BuildStateConf([]string{"RESOURCE_PREPARING", "CREATING", "CLASS_CHANGING", "MODULE_UPGRADING", "NET_CREATING", "NET_DELETING"},
		[]string{"ACTIVE"}, d.Timeout(schema.TimeoutDelete), 10*time.Second, selectDBService.SelectDBDbInstanceStateRefreshFunc(d.Id(), []string{"DELETING"}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	payment := convertPaymentTypeToChargeType(d.Get("payment_type"))

	if payment == string(PrePaid) {
		request := make(map[string]string)
		request["payment_type"] = string("POSTPAY")
		if _, err := selectDBService.ModifySelectDBInstancePaymentType(d.Id(), request); err != nil {
			return WrapError(err)
		}
	}

	_, err := selectDBService.DeleteSelectDBInstance(d.Id())
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), "DeleteDBInstance", AlibabaCloudSdkGoERROR)
	}
	return nil
}

func buildSelectDBCreateInstanceRequest(d *schema.ResourceData, meta interface{}) (map[string]interface{}, error) {
	client := meta.(*connectivity.AliyunClient)

	request := map[string]interface{}{
		"Engine":                "SelectDB",
		"EngineVersion":         "3.0",
		"DBInstanceClass":       d.Get("db_instance_class").(string),
		"RegionId":              client.RegionId,
		"ZoneId":                d.Get("zone_id").(string),
		"VpcId":                 d.Get("vpc_id").(string),
		"VSwitchId":             d.Get("vswitch_id").(string),
		"DBInstanceDescription": d.Get("db_instance_description").(string),
	}
	cache_size, exist := d.GetOkExists("cache_size")
	if exist {
		request["CacheSize"] = cache_size.(int)
	}

	payType := convertPaymentTypeToChargeType(d.Get("payment_type"))

	if payType == string(PostPaid) {
		request["ChargeType"] = string("POSTPAY")
	} else if payType == string(PrePaid) {
		request["ChargeType"] = string("PREPAY")
		request["Period"] = d.Get("period").(string)
		period_time, _ := d.GetOkExists("period_time")
		request["UsedTime"] = strconv.Itoa(period_time.(int))
	}

	return request, nil
}
