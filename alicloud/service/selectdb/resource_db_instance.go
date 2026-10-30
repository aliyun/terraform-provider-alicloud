package selectdb

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/convert"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/tags"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/validators"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceDbInstance() *schema.Resource {
	return &schema.Resource{
		Create: resourceDbInstanceCreate,
		Read:   resourceDbInstanceRead,
		Update: resourceDbInstanceUpdate,
		Delete: resourceDbInstanceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(60 * time.Minute),
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
				ValidateFunc: validators.StringInSlice([]string{"PayAsYouGo", "Subscription"}, false),
				Required:     true,
			},
			"period": {
				Type:             schema.TypeString,
				ValidateFunc:     validators.StringInSlice([]string{"Year", "Month"}, false),
				Optional:         true,
				DiffSuppressFunc: postPaidDiffSuppressFunc,
			},
			"period_time": {
				Type:             schema.TypeInt,
				ValidateFunc:     validators.IntInSlice([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 12, 24, 36}),
				Optional:         true,
				DiffSuppressFunc: postPaidDiffSuppressFunc,
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
			"engine_minor_version": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"upgraded_engine_minor_version"},
			},
			"upgraded_engine_minor_version": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"engine_minor_version"},
				Deprecated:    "Field `upgraded_engine_minor_version` has been deprecated from provider version 1.248.0. New field `engine_minor_version` instead.",
			},
			"tags": tags.TagsSchema(),
			// flag for public network and update
			"enable_public_network": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"admin_pass": {
				Type:      schema.TypeString,
				Sensitive: true,
				Optional:  true,
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

func resourceDbInstanceCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	selectDBService := NewService(client)
	request, err := buildCreateInstanceRequest(d, meta)
	if err != nil {
		return errs.WrapError(err)
	}
	action := "CreateDBInstance"
	response, err := selectDBService.RequestProcess(request, action, "POST")
	if err != nil {
		return errs.WrapError(err)
	}
	if resp, err := jsonpath.Get("$.Data", response); err != nil || resp == nil {
		return errs.WrapErrorf(err, errs.IdMsg, "alicloud_selectdb_db_instance")
	} else {
		instanceId := resp.(map[string]interface{})["DBInstanceId"].(string)
		d.SetId(fmt.Sprint(instanceId))
	}

	stateConfPreparing := helper.BuildStateConf([]string{"RESOURCE_PREPARING", "CREATING"}, []string{"ACTIVATION"}, d.Timeout(schema.TimeoutCreate), 1*time.Minute, selectDBService.DbInstanceStateRefreshFunc(d.Id(), []string{"DELETING"}))
	if _, err := stateConfPreparing.WaitForState(); err != nil {
		return errs.WrapErrorf(err, errs.IdMsg, d.Id())
	}

	enable_public_network, exist := d.GetOkExists("enable_public_network")
	if exist {
		if enable_public_network == true {
			if _, err := selectDBService.AllocateInstancePublicConnection(d.Id()); err != nil {
				return errs.WrapError(err)
			}
			stateConf := helper.BuildStateConf([]string{"NET_CREATING"}, []string{"ACTIVATION"}, d.Timeout(schema.TimeoutCreate), 1*time.Minute, selectDBService.DbInstanceStateRefreshFunc(d.Id(), []string{"DELETING"}))
			if _, err := stateConf.WaitForState(); err != nil {
				return errs.WrapErrorf(err, errs.IdMsg, d.Id())
			}
		}
	}

	return resourceDbInstanceUpdate(d, meta)
}

func resourceDbInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	selectDBService := NewService(client)
	d.Partial(true)

	if !d.IsNewResource() && d.HasChange("payment_type") {
		_, newPayment := d.GetChange("payment_type")

		request := make(map[string]string)
		payment := convert.ConvertPaymentTypeToChargeType(newPayment)

		if payment == "PostPaid" {
			request["payment_type"] = string("POSTPAY")
		} else if payment == "PrePaid" {
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
		if _, err := selectDBService.ModifyInstancePaymentType(d.Id(), request); err != nil {
			return errs.WrapError(err)
		}

	}

	if !d.IsNewResource() && d.HasChange("enable_public_network") {
		oldNetStatus, newNetStatus := d.GetChange("enable_public_network")
		if oldNetStatus.(bool) && !newNetStatus.(bool) {
			netResp, err := selectDBService.DescribeDbInstanceNetInfo(d.Id())
			if err != nil {
				return errs.WrapError(err)
			}
			connectionString := ""
			resultNet, _ := netResp["DBInstanceNetInfos"].([]interface{})
			for _, v := range resultNet {
				item := v.(map[string]interface{})
				if item["NetType"].(string) == "PUBLIC" && !strings.Contains(item["ConnectionString"].(string), "webui") {
					connectionString = item["ConnectionString"].(string)
				}
			}

			if _, err := selectDBService.ReleaseInstancePublicConnection(d.Id(), connectionString); err != nil {
				return errs.WrapError(err)
			}
			stateConf := helper.BuildStateConf([]string{"NET_DELETING"}, []string{"ACTIVATION"}, d.Timeout(schema.TimeoutUpdate), 1*time.Minute, selectDBService.DbInstanceStateRefreshFunc(d.Id(), []string{"DELETING"}))
			if _, err := stateConf.WaitForState(); err != nil {
				return errs.WrapErrorf(err, errs.IdMsg, d.Id())
			}
		} else if !oldNetStatus.(bool) && newNetStatus.(bool) {
			if _, err := selectDBService.AllocateInstancePublicConnection(d.Id()); err != nil {
				return errs.WrapError(err)
			}
			stateConf := helper.BuildStateConf([]string{"NET_CREATING"}, []string{"ACTIVATION"}, d.Timeout(schema.TimeoutUpdate), 1*time.Minute, selectDBService.DbInstanceStateRefreshFunc(d.Id(), []string{"DELETING"}))
			if _, err := stateConf.WaitForState(); err != nil {
				return errs.WrapErrorf(err, errs.IdMsg, d.Id())
			}
		}

	}

	if !d.IsNewResource() && (d.HasChange("engine_minor_version") || d.HasChange("upgraded_engine_minor_version")) {
		var newVersion interface{} = ""
		if d.HasChange("upgraded_engine_minor_version") {
			_, newVersion = d.GetChange("upgraded_engine_minor_version")
		} else if d.HasChange("engine_minor_version") {
			_, newVersion = d.GetChange("engine_minor_version")
		}

		instanceId := fmt.Sprint(d.Id())
		instanceResp, err := selectDBService.DescribeDbInstance(instanceId)
		if err != nil {
			return errs.WrapError(err)
		}
		upgradeTargetVersion := ""
		canUpgradeVersion := instanceResp["CanUpgradeVersions"].([]interface{})
		for _, version := range canUpgradeVersion {
			if strings.HasPrefix(version.(string), newVersion.(string)) {
				upgradeTargetVersion = version.(string)
				break
			}
		}

		if upgradeTargetVersion == "" {
			return errs.WrapErrorf(err, "Invalid upgrade version for %s, cannot upgrade to %s", d.Id(), newVersion.(string), errs.AlibabaCloudSdkGoERROR)
		} else {
			// todo maintaintime update
			_, err = selectDBService.UpgradeInstanceEngineVersion(d.Id(), upgradeTargetVersion, false)
			if err != nil {
				return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), "UpgradeInstanceEngineVersion", errs.AlibabaCloudSdkGoERROR)
			}
			stateConf := helper.BuildStateConf([]string{"MODULE_UPGRADING"}, []string{"ACTIVATION"}, d.Timeout(schema.TimeoutUpdate), 10*time.Second, selectDBService.DbInstanceStateRefreshFunc(d.Id(), []string{"DELETING"}))
			if _, err := stateConf.WaitForState(); err != nil {
				return errs.WrapErrorf(err, errs.IdMsg, d.Id())
			}
		}

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
		_, err := selectDBService.ModifyCluster(defaultBeId, newClass.(string), cache_size)
		if err != nil {
			return errs.WrapErrorf(err, errs.DefaultErrorMsg, defaultBeId, "ModifyDBCluster", errs.AlibabaCloudSdkGoERROR)
		}
		stateConf := helper.BuildStateConf([]string{"RESOURCE_PREPARING", "CLASS_CHANGING"}, []string{"ACTIVATION"}, d.Timeout(schema.TimeoutUpdate), 1*time.Minute, selectDBService.DbClusterStateRefreshFunc(defaultBeId, []string{"DELETING"}))
		if _, err := stateConf.WaitForState(); err != nil {
			return errs.WrapErrorf(err, errs.IdMsg, defaultBeId)
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
		_, err := selectDBService.ModifyCluster(defaultBeId, db_cluster_class, newCacheSize.(int))
		if err != nil {
			return errs.WrapErrorf(err, errs.DefaultErrorMsg, defaultBeId, "ModifyDBCluster", errs.AlibabaCloudSdkGoERROR)
		}
		stateConf := helper.BuildStateConf([]string{"RESOURCE_PREPARING", "CLASS_CHANGING"}, []string{"ACTIVATION"}, d.Timeout(schema.TimeoutUpdate), 1*time.Minute, selectDBService.DbClusterStateRefreshFunc(defaultBeId, []string{"DELETING"}))
		if _, err := stateConf.WaitForState(); err != nil {
			return errs.WrapErrorf(err, errs.IdMsg, defaultBeId)
		}
	}

	if !d.IsNewResource() && d.HasChange("db_instance_description") {
		_, newDesc := d.GetChange("db_instance_description")
		_, err := selectDBService.ModifyInstanceDescription(d.Id(), newDesc.(string))
		if err != nil {
			return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), "ModifyInstanceDescription", errs.AlibabaCloudSdkGoERROR)
		}
	}

	if d.HasChange("tags") {
		added, removed := tags.ParsingTags(d)
		if err := selectDBService.SetResourceTags(d.Id(), added, removed); err != nil {
			return errs.WrapError(err)
		}
	}

	if d.HasChange("admin_pass") {
		_, newPass := d.GetChange("admin_pass")
		if _, err := selectDBService.ModifyInstanceAdminPass(d.Id(), newPass.(string)); err != nil {
			return errs.WrapError(err)
		}
	}

	if d.HasChange("desired_security_ip_lists") {
		_, newDesc := d.GetChange("desired_security_ip_lists")
		for _, v := range newDesc.([]interface{}) {
			item := v.(map[string]interface{})
			_, err := selectDBService.ModifyDbInstanceSecurityIPList(d.Id(), item["group_name"].(string), item["security_ip_list"].(string))
			if err != nil {
				return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), "ModifySecurityIPList", errs.AlibabaCloudSdkGoERROR)
			}
		}
	}
	stateConf := helper.BuildStateConf([]string{"RESOURCE_PREPARING", "CREATING", "CLASS_CHANGING", "MODULE_UPGRADING", "NET_CREATING", "NET_DELETING"},
		[]string{"ACTIVATION"}, d.Timeout(schema.TimeoutUpdate), 10*time.Second, selectDBService.DbInstanceStateRefreshFunc(d.Id(), []string{"DELETING"}))

	if _, err := stateConf.WaitForState(); err != nil {
		return errs.WrapErrorf(err, errs.IdMsg, d.Id())
	}

	d.Partial(false)
	return resourceDbInstanceRead(d, meta)
}

func resourceDbInstanceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	selectDBService := NewService(client)

	instanceId := fmt.Sprint(d.Id())
	instanceResp, err := selectDBService.DescribeDbInstance(instanceId)
	if err != nil {
		if !d.IsNewResource() && errs.NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return errs.WrapError(err)
	}

	d.Set("engine", instanceResp["Engine"])
	d.Set("engine_minor_version", instanceResp["EngineMinorVersion"])
	d.Set("upgraded_engine_minor_version", instanceResp["EngineMinorVersion"])

	d.Set("region_id", instanceResp["RegionId"])
	d.Set("zone_id", instanceResp["ZoneId"])
	d.Set("vpc_id", instanceResp["VpcId"])

	d.Set("payment_type", convert.ConvertChargeTypeToPaymentType(instanceResp["ChargeType"]))

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
		tagsMap := make(map[string]interface{})
		for _, t := range resp.([]interface{}) {
			key := t.(map[string]interface{})["TagKey"].(string)
			value := t.(map[string]interface{})["TagValue"].(string)
			if !tags.IgnoredTags(key, value) {
				tagsMap[key] = value
			}
		}
		d.Set("tags", tagsMap)
	}

	netResp, err := selectDBService.DescribeDbInstanceNetInfo(d.Id())
	if err != nil {
		return errs.WrapError(err)
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

	securityIpArrayList, err := selectDBService.DescribeDbInstanceSecurityIPList(d.Id())
	if err != nil {
		return errs.WrapError(err)
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

func resourceDbInstanceDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	selectDBService := NewService(client)

	stateConf := helper.BuildStateConf([]string{"RESOURCE_PREPARING", "CREATING", "CLASS_CHANGING", "MODULE_UPGRADING", "NET_CREATING", "NET_DELETING"},
		[]string{"ACTIVATION"}, d.Timeout(schema.TimeoutDelete), 10*time.Second, selectDBService.DbInstanceStateRefreshFunc(d.Id(), []string{"DELETING"}))
	if _, err := stateConf.WaitForState(); err != nil {
		return errs.WrapErrorf(err, errs.IdMsg, d.Id())
	}

	payment := convert.ConvertPaymentTypeToChargeType(d.Get("payment_type"))

	if payment == "PrePaid" {
		request := make(map[string]string)
		request["payment_type"] = string("POSTPAY")
		if _, err := selectDBService.ModifyInstancePaymentType(d.Id(), request); err != nil {
			return errs.WrapError(err)
		}
	}

	_, err := selectDBService.DeleteInstance(d.Id())
	if err != nil {
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), "DeleteDBInstance", errs.AlibabaCloudSdkGoERROR)
	}
	return nil
}

func buildCreateInstanceRequest(d *schema.ResourceData, meta interface{}) (map[string]interface{}, error) {
	client := meta.(*connectivity.AliyunClient)

	request := map[string]interface{}{
		"Engine":                "SelectDB",
		"DBInstanceClass":       d.Get("db_instance_class").(string),
		"RegionId":              client.RegionId,
		"ZoneId":                d.Get("zone_id").(string),
		"VpcId":                 d.Get("vpc_id").(string),
		"VSwitchId":             d.Get("vswitch_id").(string),
		"DBInstanceDescription": d.Get("db_instance_description").(string),
	}

	if v, ok := d.GetOk("engine_minor_version"); ok {
		request["EngineVersion"] = v.(string)
	} else if v, ok := d.GetOk("upgraded_engine_minor_version"); ok {
		request["EngineVersion"] = v.(string)
	}

	cache_size, exist := d.GetOkExists("cache_size")
	if exist {
		request["CacheSize"] = cache_size.(int)
	}

	payType := convert.ConvertPaymentTypeToChargeType(d.Get("payment_type"))

	if payType == "PostPaid" {
		request["ChargeType"] = string("POSTPAY")
	} else if payType == "PrePaid" {
		request["ChargeType"] = string("PREPAY")
		request["Period"] = d.Get("period").(string)
		period_time, _ := d.GetOkExists("period_time")
		request["UsedTime"] = strconv.Itoa(period_time.(int))
	}

	return request, nil
}
