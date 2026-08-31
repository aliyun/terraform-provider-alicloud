package alicloud

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAliCloudSelectDBDbCluster() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudSelectDBDbClusterCreate,
		Read:   resourceAliCloudSelectDBDbClusterRead,
		Update: resourceAliCloudSelectDBDbClusterUpdate,
		Delete: resourceAliCloudSelectDBDbClusterDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"db_instance_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"payment_type": {
				Type:         schema.TypeString,
				ValidateFunc: StringInSlice([]string{"PayAsYouGo", "Subscription"}, false),
				Required:     true,
				ForceNew:     true,
			},
			"db_cluster_class": {
				Type:     schema.TypeString,
				Required: true,
			},
			"cache_size": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"db_cluster_description": {
				Type:     schema.TypeString,
				Required: true,
			},
			"desired_params": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"value": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"desired_status": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: StringInSlice([]string{"STOPPING", "STARTING", "RESTART"}, false),
			},

			// computed
			"db_cluster_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"engine": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"engine_version": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"cpu": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"memory": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"region_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"zone_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"vpc_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"param_change_logs": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"old_value": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"new_value": {
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
						"config_id": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"is_applied": {
							Type:     schema.TypeBool,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func resourceAliCloudSelectDBDbClusterCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	selectDBService := SelectDBService{client}

	request, err := buildSelectDBCreateClusterRequest(d, meta)
	if err != nil {
		return WrapError(err)
	}
	action := "CreateDBCluster"
	response, err := selectDBService.RequestProcessForSelectDB(request, action, "POST")
	if err != nil {
		return WrapError(err)
	}
	if resp, err := jsonpath.Get("$.Data", response); err != nil || resp == nil {
		return WrapErrorf(err, IdMsg, "alicloud_selectdb_db_clusters")
	} else {
		clusterId := resp.(map[string]interface{})["ClusterId"].(string)
		d.SetId(fmt.Sprint(d.Get("db_instance_id").(string) + ":" + clusterId))
	}

	stateConf := BuildStateConf([]string{"RESOURCE_PREPARING", "CREATING"}, []string{"ACTIVATION"}, d.Timeout(schema.TimeoutCreate), 20*time.Second, selectDBService.SelectDBDbClusterStateRefreshFunc(d.Id(), []string{"DELETING"}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}
	return resourceAliCloudSelectDBDbClusterUpdate(d, meta)
}

func resourceAliCloudSelectDBDbClusterUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	selectDBService := SelectDBService{client}
	d.Partial(true)

	cacheSizeModified := false
	if !d.IsNewResource() && (d.HasChange("db_cluster_class")) {
		_, newClass := d.GetChange("db_cluster_class")
		cache_size := 0
		if d.HasChange("cache_size") {
			_, newCacheSize := d.GetChange("cache_size")
			cache_size = newCacheSize.(int)
			cacheSizeModified = true
		}
		_, err := selectDBService.ModifySelectDBCluster(d.Id(), newClass.(string), cache_size)
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), "ModifyDBCluster", AlibabaCloudSdkGoERROR)
		}
		stateConf := BuildStateConf([]string{"RESOURCE_PREPARING", "CLASS_CHANGING"}, []string{"ACTIVATION"}, d.Timeout(schema.TimeoutUpdate), 1*time.Minute, selectDBService.SelectDBDbClusterStateRefreshFunc(d.Id(), []string{"DELETING"}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
		d.SetPartial("db_cluster_class")
		if d.HasChange("cache_size") {
			d.SetPartial("cache_size")
		}
	}

	if !d.IsNewResource() && d.HasChange("cache_size") && !cacheSizeModified {
		_, newCacheSize := d.GetChange("cache_size")
		db_cluster_class := d.Get("db_cluster_class").(string)
		if d.HasChange("db_cluster_class") {
			_, newClass := d.GetChange("db_cluster_class")
			db_cluster_class = newClass.(string)
		}
		_, err := selectDBService.ModifySelectDBCluster(d.Id(), db_cluster_class, newCacheSize.(int))
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), "ModifyDBCluster", AlibabaCloudSdkGoERROR)
		}
		stateConf := BuildStateConf([]string{"RESOURCE_PREPARING", "CLASS_CHANGING"}, []string{"ACTIVATION"}, d.Timeout(schema.TimeoutUpdate), 1*time.Minute, selectDBService.SelectDBDbClusterStateRefreshFunc(d.Id(), []string{"DELETING"}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
		d.SetPartial("cache_size")
		if d.HasChange("db_cluster_class") {
			d.SetPartial("db_cluster_class")
		}
	}

	if !d.IsNewResource() && d.HasChange("db_cluster_description") {
		_, newDesc := d.GetChange("db_cluster_description")
		_, err := selectDBService.ModifySelectDBClusterDescription(d.Id(), newDesc.(string))
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), "ModifyBEClusterAttribute", AlibabaCloudSdkGoERROR)
		}
		d.SetPartial("db_cluster_description")
	}

	if !d.IsNewResource() && d.HasChange("desired_status") {
		_, newStatus := d.GetChange("desired_status")
		oldStatus := d.Get("status")
		if oldStatus.(string) != "" && newStatus.(string) != "" {
			_, err := selectDBService.UpdateSelectDBClusterStatus(d.Id(), newStatus.(string))
			if err != nil {
				return WrapErrorf(err, DefaultErrorMsg, d.Id(), "UpdateSelectDBClusterStatus", AlibabaCloudSdkGoERROR)
			}
			newStatusFinal := convertSelectDBClusterStatusActionFinal(newStatus.(string))
			if newStatusFinal == "" {
				return WrapErrorf(err, DefaultErrorMsg, d.Id(), "UpdateSelectDBClusterStatus", AlibabaCloudSdkGoERROR)
			}
			stateConf := BuildStateConf([]string{newStatus.(string)}, []string{newStatusFinal}, d.Timeout(schema.TimeoutUpdate), 1*time.Minute, selectDBService.SelectDBDbClusterStateRefreshFunc(d.Id(), []string{}))
			if _, err := stateConf.WaitForState(); err != nil {
				return WrapErrorf(err, IdMsg, d.Id())
			}
			d.SetPartial("desired_status")
		}
	}

	if d.HasChange("desired_params") {
		oldConfig, newConfig := d.GetChange("desired_params")
		oldConfigMap := oldConfig.([]interface{})
		newConfigMap := newConfig.([]interface{})
		oldConfigMapIndex := make(map[string]string)
		for _, v := range oldConfigMap {
			item := v.(map[string]interface{})
			oldConfigMapIndex[item["name"].(string)] = item["value"].(string)
		}
		newConfigMapIndex := make(map[string]string)
		for _, v := range newConfigMap {
			item := v.(map[string]interface{})
			newConfigMapIndex[item["name"].(string)] = item["value"].(string)
		}

		diffConfig := make(map[string]string)
		for k, v := range newConfigMapIndex {
			if oldConfigMapIndex[k] != v {
				diffConfig[k] = v
			}
		}

		if _, err := selectDBService.UpdateSelectDBDbClusterConfig(d.Id(), diffConfig); err != nil {
			return WrapError(err)
		}
		d.SetPartial("desired_params")

		stateConf := BuildStateConf([]string{"RESTARTING", "MODIFY_PARAM"}, []string{"ACTIVATION"}, d.Timeout(schema.TimeoutUpdate), 10*time.Second, selectDBService.SelectDBDbClusterStateRefreshFunc(d.Id(), []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}

	}

	d.Partial(false)
	return resourceAliCloudSelectDBDbClusterRead(d, meta)
}

func resourceAliCloudSelectDBDbClusterRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	selectDBService := SelectDBService{client}

	clusterResp, err := selectDBService.DescribeSelectDBDbCluster(d.Id())
	if err != nil {
		return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_selectdb_db_cluster", AlibabaCloudSdkGoERROR)
	}
	cpu, _ := clusterResp["CpuCores"].(json.Number).Int64()
	memory, _ := clusterResp["Memory"].(json.Number).Int64()
	cache, _ := clusterResp["CacheStorageSizeGB"].(json.Number).Int64()

	d.Set("status", clusterResp["Status"])
	d.Set("create_time", clusterResp["CreatedTime"])
	d.Set("db_cluster_description", clusterResp["DbClusterName"])
	d.Set("payment_type", convertChargeTypeToPaymentType(clusterResp["ChargeType"]))
	d.Set("db_instance_id", clusterResp["DbInstanceName"])
	d.Set("db_cluster_class", clusterResp["DbClusterClass"])
	d.Set("cpu", cpu)
	d.Set("memory", memory)
	d.Set("cache_size", cache)

	d.Set("engine", fmt.Sprint(clusterResp["Engine"]))
	d.Set("engine_version", fmt.Sprint(clusterResp["EngineVersion"]))
	d.Set("vpc_id", fmt.Sprint(clusterResp["VpcId"]))
	d.Set("zone_id", fmt.Sprint(clusterResp["ZoneId"]))
	d.Set("region_id", fmt.Sprint(clusterResp["RegionId"]))

	configChangeArrayList, err := selectDBService.DescribeSelectDBDbClusterConfigChangeLog(d.Id())
	if err != nil {
		return WrapError(err)
	}
	configChangeArray := make([]map[string]interface{}, 0)
	for _, v := range configChangeArrayList {
		m1 := v.(map[string]interface{})
		ConfigId, _ := m1["Id"].(json.Number).Int64()

		temp1 := map[string]interface{}{
			"name":         m1["Name"].(string),
			"old_value":    m1["OldValue"].(string),
			"new_value":    m1["NewValue"].(string),
			"is_applied":   m1["IsApplied"].(bool),
			"gmt_created":  m1["GmtCreated"].(string),
			"gmt_modified": m1["GmtModified"].(string),
			"config_id":    ConfigId,
		}
		configChangeArray = append(configChangeArray, temp1)
	}
	d.Set("param_change_logs", configChangeArray)
	return nil
}

func resourceAliCloudSelectDBDbClusterDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	selectDBService := SelectDBService{client}

	stateConf := BuildStateConf([]string{"RESOURCE_PREPARING", "CLASS_CHANGING", "CREATING", "STOPPING", "STARTING", "RESTARTING", "RESTART", "MODIFY_PARAM"},
		[]string{"ACTIVATION"}, d.Timeout(schema.TimeoutUpdate), 10*time.Second, selectDBService.SelectDBDbClusterStateRefreshFunc(d.Id(), []string{}))

	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	_, err := selectDBService.DescribeSelectDBDbCluster(d.Id())
	if err != nil {
		if NotFoundError(err) {
			return nil
		}
		return WrapError(err)
	}

	_, err = selectDBService.DeleteSelectDBCluster(d.Id())
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), "DeleteDBCluster", AlibabaCloudSdkGoERROR)
	}

	instance_id := d.Get("db_instance_id").(string)
	// cluster deleting cannot be checked, use instance from class changing to active instead.
	// cluster deleting = related instance update
	stateConf = BuildStateConf([]string{"CLASS_CHANGING"}, []string{"ACTIVATION"}, d.Timeout(schema.TimeoutDelete), 10*time.Second, selectDBService.SelectDBDbInstanceStateRefreshFunc(instance_id, []string{"DELETING"}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}
	return nil
}

func buildSelectDBCreateClusterRequest(d *schema.ResourceData, meta interface{}) (map[string]interface{}, error) {
	client := meta.(*connectivity.AliyunClient)
	selectDBService := SelectDBService{client}

	instanceResp, err := selectDBService.DescribeSelectDBDbInstance(d.Get("db_instance_id").(string))
	if err != nil {
		return nil, WrapErrorf(err, DefaultErrorMsg, d.Id())
	}

	vswitchId := ""
	netResp, err := selectDBService.DescribeSelectDBDbInstanceNetInfo(d.Get("db_instance_id").(string))
	if err != nil {
		return nil, WrapErrorf(err, DefaultErrorMsg, d.Get("db_instance_id").(string))
	}
	resultClusterNet, _ := netResp["DBInstanceNetInfos"].([]interface{})
	for _, v := range resultClusterNet {
		item := v.(map[string]interface{})["VswitchId"].(string)
		if item != "" {
			vswitchId = item
			break
		}
	}

	cache_size, exist := d.GetOkExists("cache_size")
	if !exist {
		return nil, WrapErrorf(err, DefaultErrorMsg, d.Id())
	}

	request := map[string]interface{}{
		"DBInstanceId":         d.Get("db_instance_id").(string),
		"Engine":               "SelectDB",
		"EngineVersion":        instanceResp["EngineVersion"],
		"DBClusterClass":       d.Get("db_cluster_class").(string),
		"RegionId":             client.RegionId,
		"ZoneId":               instanceResp["ZoneId"],
		"VpcId":                instanceResp["VpcId"],
		"VSwitchId":            vswitchId,
		"CacheSize":            cache_size.(int),
		"DBClusterDescription": Trim(d.Get("db_cluster_description").(string)),
	}

	payType := convertPaymentTypeToChargeType(d.Get("payment_type"))

	if payType == string(PostPaid) {
		request["ChargeType"] = string("Postpaid")
	} else if payType == string(PrePaid) {
		period_time, _ := d.GetOkExists("period_time")
		request["ChargeType"] = string("Prepaid")
		request["Period"] = d.Get("period").(string)
		request["UsedTime"] = strconv.Itoa(period_time.(int))
	}

	return request, nil
}

func convertSelectDBClusterStatusActionFinal(source string) string {
	action := ""
	switch source {
	case "STOPPING", "STOPPED":
		action = "STOPPED"
	case "STARTING":
		action = "ACTIVATION"
	case "RESTART", "RESTARTING":
		action = "ACTIVATION"
	}
	return action
}
