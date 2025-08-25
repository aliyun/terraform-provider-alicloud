package alicloud

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceAlicloudClickHouseDbCluster() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudClickHouseDbClusterCreate,
		Read:   resourceAlicloudClickHouseDbClusterRead,
		Update: resourceAlicloudClickHouseDbClusterUpdate,
		Delete: resourceAlicloudClickHouseDbClusterDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(60 * time.Minute),
			Update: schema.DefaultTimeout(60 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"category": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"Basic", "HighAvailability"}, false),
				ForceNew:     true,
			},
			"db_cluster_access_white_list": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"db_cluster_ip_array_attribute": {
							Type:     schema.TypeString,
							Optional: true,
							Removed:  "Field 'db_cluster_ip_array_attribute' has been removed from provider",
						},
						"db_cluster_ip_array_name": {
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
			"db_cluster_class": {
				Type:     schema.TypeString,
				Required: true,
			},
			"db_cluster_network_type": {
				Type:         schema.TypeString,
				ValidateFunc: validation.StringInSlice([]string{"vpc"}, false),
				Required:     true,
				ForceNew:     true,
			},
			"db_cluster_version": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"19.15.2.2", "20.3.10.75", "20.8.7.15", "21.8.10.19", "22.8.5.29", "23.8", "25.3"}, false),
			},
			"db_node_storage": {
				Type:     schema.TypeString,
				Required: true,
			},
			"db_node_group_count": {
				Type:         schema.TypeInt,
				ValidateFunc: validation.IntBetween(1, 48),
				Required:     true,
			},

			"encryption_key": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"encryption_type": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"payment_type": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"PayAsYouGo", "Subscription"}, false),
				ForceNew:     true,
			},
			"period": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"Month", "Year"}, false),
			},
			"renewal_status": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: StringInSlice([]string{"AutoRenewal", "Normal"}, false),
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					if v, ok := d.GetOk("payment_type"); ok && v.(string) == "Subscription" {
						return false
					}
					return true
				},
			},
			"storage_type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"cloud_essd", "cloud_efficiency", "cloud_essd_pl2", "cloud_essd_pl3"}, false),
			},
			"used_time": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"vswitch_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},

			"db_cluster_description": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
			"status": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringInSlice([]string{"Creating", "Deleting", "Restarting", "Preparing", "Running"}, false),
			},
			"maintain_time": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
			"zone_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"multi_zone_vswitch_list": {
				Type:     schema.TypeSet,
				Optional: true,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"zone_id": {
							Optional: true,
							Computed: true,
							ForceNew: true,
							Type:     schema.TypeString,
						},
						"vswitch_id": {
							Required: true,
							ForceNew: true,
							Type:     schema.TypeString,
						},
					},
				},
			},
			"vpc_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"connection_string": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"port": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"resource_group_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"allocate_public_connection": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"public_connection_string": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"cold_storage": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringInSlice([]string{"ENABLE", "DISABLE"}, false),
			},
		},
	}
}

func resourceAlicloudClickHouseDbClusterCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	clickhouseService := ClickhouseService{client}
	var response map[string]interface{}
	action := "CreateDBInstance"
	request := make(map[string]interface{})
	var err error
	request["DBClusterCategory"] = d.Get("category")
	request["DBClusterClass"] = d.Get("db_cluster_class")
	if v, ok := d.GetOk("db_cluster_description"); ok {
		request["DBClusterDescription"] = v
	}
	request["DBClusterNetworkType"] = d.Get("db_cluster_network_type")
	request["DBClusterVersion"] = d.Get("db_cluster_version")
	request["DBNodeGroupCount"] = d.Get("db_node_group_count")
	request["DBNodeStorage"] = d.Get("db_node_storage")
	if v, ok := d.GetOk("renewal_status"); ok {
		switch v.(string) {
		case "Normal":
			request["AutoRenew"] = false
		case "AutoRenewal":
			request["AutoRenew"] = true
		default:
		}
	}
	if v, ok := d.GetOk("encryption_key"); ok {
		request["EncryptionKey"] = v
	}
	if v, ok := d.GetOk("encryption_type"); ok {
		request["EncryptionType"] = v
	}
	request["PayType"] = convertClickHouseDbClusterPaymentTypeRequest(d.Get("payment_type").(string))

	if v, ok := d.GetOk("period"); ok {
		request["Period"] = v
	}
	request["RegionId"] = client.RegionId
	request["DbNodeStorageType"] = d.Get("storage_type")
	if v, ok := d.GetOk("used_time"); ok {
		request["UsedTime"] = v
	}

	if v, ok := d.GetOk("vpc_id"); ok {
		request["VPCId"] = v
	}

	if v, ok := d.GetOk("zone_id"); ok {
		request["ZoneId"] = v
	}

	if v, ok := d.GetOk("vswitch_id"); ok {
		request["VSwitchId"] = v
	}

	if v, ok := d.GetOk("multi_zone_vswitch_list"); ok {
		vlist := v.(*schema.Set).List()
		if len(vlist) != 2 {
			return WrapError(fmt.Errorf("multi_zone_vswitch_list must have 2 different zones and vswitches, got: %d", len(vlist)))
		}
		vswitch1, vswitch2 := vlist[0].(map[string]interface{}), vlist[1].(map[string]interface{})
		if vswitch1["zone_id"] == vswitch2["zone_id"] || vswitch1["vswitch_id"] == vswitch2["vswitch_id"] {
			return WrapError(fmt.Errorf("multi_zone_vswitch_list must have 2 different zone ids and vswitch ids"))
		}
		request["ZoneIdBak"] = vswitch1["zone_id"]
		request["VSwitchBak"] = vswitch1["vswitch_id"]
		request["ZoneIdBak2"] = vswitch2["zone_id"]
		request["VSwitchBak2"] = vswitch2["vswitch_id"]
	}

	if v, ok := d.GetOk("resource_group_id"); ok {
		request["ResourceGroupId"] = v
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

	if request["ZoneIdBak"] == nil && request["VSwitchBak"] != nil {
		vpcService := VpcService{client}
		vsw, err := vpcService.DescribeVSwitchWithTeadsl(request["VSwitchBak"].(string))
		if err != nil {
			return WrapError(err)
		}
		request["ZoneIdBak"] = vsw["ZoneId"]
	}

	if request["ZondIdBak2"] == nil && request["VSwitchBak2"] != nil {
		vpcService := VpcService{client}
		vsw, err := vpcService.DescribeVSwitchWithTeadsl(request["VSwitchBak2"].(string))
		if err != nil {
			return WrapError(err)
		}
		request["ZondIdBak2"] = vsw["ZoneId"]
	}

	request["ClientToken"] = buildClientToken("CreateDBInstance")
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.RpcPost("clickhouse", "2019-11-11", action, nil, request, true)
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_click_house_db_cluster", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(response["DBClusterId"]))
	stateConf := BuildStateConf([]string{"Creating", "Deleting", "Restarting", "Preparing"}, []string{"Running"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, clickhouseService.ClickHouseDbClusterStateRefreshFunc(d.Id(), []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}
	if d.Get("payment_type").(string) == "Subscription" {
		stateConf = BuildStateConf([]string{""}, []string{"Normal", "AutoRenewal", "NotRenewal"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, clickhouseService.ClickHouseAutoRenewStatusRefreshFunc(d.Id(), []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
	}
	return resourceAlicloudClickHouseDbClusterUpdate(d, meta)
}
func resourceAlicloudClickHouseDbClusterRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	clickhouseService := ClickhouseService{client}
	object, err := clickhouseService.DescribeClickHouseDbCluster(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_click_house_db_cluster clickhouseService.DescribeClickHouseDbCluster Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	d.Set("db_cluster_version", object["EngineVersion"])
	d.Set("db_cluster_class", object["DBNodeClass"])
	d.Set("db_node_group_count", object["DBNodeCount"])
	d.Set("category", object["Category"])
	d.Set("db_cluster_description", object["DBClusterDescription"])
	d.Set("db_cluster_network_type", object["DBClusterNetworkType"])
	d.Set("db_node_storage", fmt.Sprint(formatInt(object["DBNodeStorage"])))
	d.Set("encryption_key", object["EncryptionKey"])
	d.Set("encryption_type", object["EncryptionType"])
	d.Set("maintain_time", object["MaintainTime"])
	d.Set("status", object["DBClusterStatus"])
	d.Set("payment_type", convertClickHouseDbClusterPaymentTypeResponse(object["PayType"].(string)))
	d.Set("storage_type", convertClickHouseDbClusterStorageTypeResponse(object["StorageType"].(string)))
	d.Set("vswitch_id", object["VSwitchId"])
	d.Set("zone_id", object["ZoneId"])
	d.Set("vpc_id", object["VpcId"])
	d.Set("connection_string", object["ConnectionString"])
	d.Set("port", object["Port"])
	d.Set("public_connection_string", object["PublicConnectionString"])
	if object["PublicConnectionString"].(string) == "" {
		d.Set("allocate_public_connection", false)
	} else {
		d.Set("allocate_public_connection", true)
	}
	d.Set("resource_group_id", object["ResourceGroupId"])

	if ZoneIdVswitchMap, ok := object["ZoneIdVswitchMap"]; ok {
		vMap := ZoneIdVswitchMap.(map[string]interface{})
		if _, ok := vMap[object["ZoneId"].(string)]; ok {
			delete(vMap, object["ZoneId"].(string))
		}
		vList := make([]map[string]interface{}, 0)
		for k, v := range vMap {
			vList = append(vList, map[string]interface{}{
				"zone_id":    k,
				"vswitch_id": v.(string),
			})
		}
		d.Set("multi_zone_vswitch_list", vList)
	}

	object, err = clickhouseService.DescribeClickHouseAutoRenewStatus(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_click_house_db_cluster clickhouseService.DescribeClickHouseAutoRenewStatus Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	if v, ok := object["RenewalStatus"]; ok {
		d.Set("renewal_status", v)
	}

	describeDBClusterAccessWhiteListObject, err := clickhouseService.DescribeDBClusterAccessWhiteList(d.Id())
	if err != nil {
		return WrapError(err)
	}
	if dBClusterAccessWhiteListMap, ok := describeDBClusterAccessWhiteListObject["DBClusterAccessWhiteList"].(map[string]interface{}); ok && dBClusterAccessWhiteListMap != nil {
		if iPArrayList, ok := dBClusterAccessWhiteListMap["IPArray"]; ok && iPArrayList != nil {
			dBClusterAccessWhiteListMaps := make([]map[string]interface{}, 0)
			for _, iPArrayListItem := range iPArrayList.([]interface{}) {
				if v, ok := iPArrayListItem.(map[string]interface{}); ok {
					if v["DBClusterIPArrayName"].(string) == "default" || v["DBClusterIPArrayName"].(string) == "dms" {
						continue
					}
					iPArrayListItemMap := make(map[string]interface{})
					iPArrayListItemMap["db_cluster_ip_array_attribute"] = v["DBClusterIPArrayAttribute"]
					iPArrayListItemMap["db_cluster_ip_array_name"] = v["DBClusterIPArrayName"]
					iPArrayListItemMap["security_ip_list"] = v["SecurityIPList"]
					dBClusterAccessWhiteListMaps = append(dBClusterAccessWhiteListMaps, iPArrayListItemMap)
				}
			}
			d.Set("db_cluster_access_white_list", dBClusterAccessWhiteListMaps)
		}
	}

	describeClickHouseOSSStorageObject, err := clickhouseService.DescribeClickHouseOSSStorage(d.Id())
	if err != nil {
		return WrapError(err)
	}
	d.Set("cold_storage", describeClickHouseOSSStorageObject["State"].(string))

	return nil
}
func resourceAlicloudClickHouseDbClusterUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	var err error
	d.Partial(true)

	update := false
	request := map[string]interface{}{
		"DBClusterId": d.Id(),
	}
	if !d.IsNewResource() && d.HasChange("db_cluster_description") {
		update = true
		if v, ok := d.GetOk("db_cluster_description"); ok {
			request["DBClusterDescription"] = v
		}
	}
	if update {
		action := "ModifyDBClusterDescription"
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("clickhouse", "2019-11-11", action, nil, request, false)
			if err != nil {
				if IsExpectedErrors(err, []string{"IncorrectDBInstanceState"}) || NeedRetry(err) {
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
		d.SetPartial("db_cluster_description")
	}
	update = false
	modifyDBClusterMaintainTimeReq := map[string]interface{}{
		"DBClusterId": d.Id(),
	}
	if !d.IsNewResource() && d.HasChange("maintain_time") {
		update = true
	}
	if v, ok := d.GetOk("maintain_time"); ok {
		modifyDBClusterMaintainTimeReq["MaintainTime"] = v
	}
	if update {
		action := "ModifyDBClusterMaintainTime"
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("clickhouse", "2019-11-11", action, nil, modifyDBClusterMaintainTimeReq, false)
			if err != nil {
				if IsExpectedErrors(err, []string{"IncorrectDBInstanceState"}) || NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})
		addDebug(action, response, modifyDBClusterMaintainTimeReq)
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
		d.SetPartial("maintain_time")
	}
	if d.HasChange("db_cluster_access_white_list") {

		oraw, nraw := d.GetChange("db_cluster_access_white_list")
		remove := oraw.(*schema.Set).Difference(nraw.(*schema.Set)).List()
		create := nraw.(*schema.Set).Difference(oraw.(*schema.Set)).List()
		if len(remove) > 0 {
			removeWhiteListReq := map[string]interface{}{
				"DBClusterId": d.Id(),
				"ModifyMode":  "Delete",
			}

			for _, whiteList := range remove {
				whiteListArg := whiteList.(map[string]interface{})
				removeWhiteListReq["DBClusterIPArrayAttribute"] = whiteListArg["db_cluster_ip_array_attribute"]
				removeWhiteListReq["DBClusterIPArrayName"] = whiteListArg["db_cluster_ip_array_name"]
				removeWhiteListReq["SecurityIps"] = whiteListArg["security_ip_list"]

				action := "ModifyDBClusterAccessWhiteList"
				wait := incrementalWait(3*time.Second, 5*time.Second)
				err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
					response, err = client.RpcPost("clickhouse", "2019-11-11", action, nil, removeWhiteListReq, false)
					if err != nil {
						if IsExpectedErrors(err, []string{"IncorrectDBInstanceState"}) || NeedRetry(err) {
							wait()
							return resource.RetryableError(err)
						}
						return resource.NonRetryableError(err)
					}
					return nil
				})
				addDebug(action, response, removeWhiteListReq)
				if err != nil {
					return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
				}
			}
		}
		if len(create) > 0 {
			createWhiteListReq := map[string]interface{}{
				"DBClusterId": d.Id(),
				"ModifyMode":  "Append",
			}

			for _, whiteList := range create {
				whiteListArg := whiteList.(map[string]interface{})
				createWhiteListReq["DBClusterIPArrayAttribute"] = whiteListArg["db_cluster_ip_array_attribute"]
				createWhiteListReq["DBClusterIPArrayName"] = whiteListArg["db_cluster_ip_array_name"]
				createWhiteListReq["SecurityIps"] = whiteListArg["security_ip_list"]

				action := "ModifyDBClusterAccessWhiteList"
				wait := incrementalWait(3*time.Second, 5*time.Second)
				err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
					response, err = client.RpcPost("clickhouse", "2019-11-11", action, nil, createWhiteListReq, false)
					if err != nil {
						if IsExpectedErrors(err, []string{"IncorrectDBInstanceState"}) || NeedRetry(err) {
							wait()
							return resource.RetryableError(err)
						}
						return resource.NonRetryableError(err)
					}
					return nil
				})
				addDebug(action, response, createWhiteListReq)
				if err != nil {
					return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
				}
			}
		}
		d.SetPartial("db_cluster_access_white_list")
	}
	if d.HasChange("status") {
		clickhouseService := ClickhouseService{client}
		object, err := clickhouseService.DescribeClickHouseDbCluster(d.Id())
		if err != nil {
			return WrapError(err)
		}
		target := d.Get("status").(string)
		if object["DBClusterStatus"].(string) != target {
			if target == "Running" {
				request := map[string]interface{}{
					"DBClusterId": d.Id(),
				}
				request["RegionId"] = client.RegionId
				action := "RestartInstance"
				wait := incrementalWait(3*time.Second, 3*time.Second)
				err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
					response, err = client.RpcPost("clickhouse", "2019-11-11", action, nil, request, false)
					if err != nil {
						if IsExpectedErrors(err, []string{"IncorrectDBInstanceState"}) || NeedRetry(err) {
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
			}
			d.SetPartial("status")
		}
		stateConf := BuildStateConf([]string{"RESTARTING"}, []string{"Running"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, clickhouseService.ClickHouseDbClusterStateRefreshFunc(d.Id(), []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
	}
	if !d.IsNewResource() && (d.HasChange("db_node_storage") || d.HasChange("db_node_group_count") || d.HasChange("db_cluster_class")) {
		clickhouseService := ClickhouseService{client}
		object, err := clickhouseService.DescribeClickHouseDbCluster(d.Id())
		if err != nil {
			return WrapError(err)
		}
		storageLocal, err := strconv.ParseInt(d.Get("db_node_storage").(string), 10, 64)
		if err != nil {
			return WrapError(err)
		}
		storageRemote, err := object["DBNodeStorage"].(json.Number).Int64()
		if err != nil {
			return WrapError(err)
		}
		if storageLocal < storageRemote {
			return WrapError(fmt.Errorf("downgrading storage is not supported"))
		}
		nodeCountLocal := d.Get("db_node_group_count").(int)
		nodeCountRemote, err := object["DBNodeCount"].(json.Number).Int64()
		if err != nil {
			return WrapError(err)
		}
		if int64(nodeCountLocal) < nodeCountRemote {
			return WrapError(fmt.Errorf("downgrading db_node_group_count is not supported"))
		}
		request := map[string]interface{}{
			"DBClusterId":       d.Id(),
			"DBNodeGroupCount":  fmt.Sprintf("%v", d.Get("db_node_group_count")),
			"DBNodeStorage":     fmt.Sprintf("%v", d.Get("db_node_storage")),
			"DBClusterClass":    fmt.Sprintf("%v", d.Get("db_cluster_class")),
			"DbNodeStorageType": convertClickHouseDbClusterStorageTypeRequest(d.Get("storage_type").(string)),
			"RegionId":          client.RegionId,
		}
		action := "ModifyDBCluster"
		var response map[string]interface{}
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("clickhouse", "2019-11-11", action, nil, request, false)
			if err != nil {
				if IsExpectedErrors(err, []string{"IncorrectDBInstanceState", "OperationDenied.OrderProcessing"}) || NeedRetry(err) {
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
		stateConf := BuildStateConf([]string{"ClassChanging", "SCALING_OUT"}, []string{"Running"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, clickhouseService.ClickHouseDbClusterStateRefreshFunc(d.Id(), []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
		d.SetPartial("db_node_storage")
		d.SetPartial("db_node_group_count")
		d.SetPartial("db_cluster_class")
	}
	if v, ok := d.GetOk("payment_type"); ok && v.(string) == "Subscription" && d.HasChange("renewal_status") && !d.IsNewResource() {
		action := "ModifyAutoRenewAttribute"
		if s, ok := d.GetOk("renewal_status"); ok {
			request := map[string]interface{}{
				"DBClusterIds":  d.Id(),
				"RenewalStatus": s,
				"RegionId":      client.RegionId,
			}
			wait := incrementalWait(3*time.Second, 3*time.Second)
			err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
				response, err = client.RpcPost("clickhouse", "2019-11-11", action, nil, request, false)
				if err != nil {
					if IsExpectedErrors(err, []string{"IncorrectDBInstanceState"}) || NeedRetry(err) {
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
			clickhouseService := ClickhouseService{client}
			stateConf := BuildStateConf([]string{""}, []string{"Normal", "AutoRenewal", "NotRenewal"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, clickhouseService.ClickHouseAutoRenewStatusRefreshFunc(d.Id(), []string{}))
			if _, err := stateConf.WaitForState(); err != nil {
				return WrapErrorf(err, IdMsg, d.Id())
			}
			d.SetPartial("renewal_status")
		}
	}
	if d.HasChange("resource_group_id") && !d.IsNewResource() {
		action := "ChangeResourceGroup"
		request := map[string]interface{}{
			"ResourceId":       d.Id(),
			"ResourceGroupId":  d.Get("resource_group_id"),
			"ResourceRegionId": client.RegionId,
		}
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("clickhouse", "2019-11-11", action, nil, request, false)
			if err != nil {
				if IsExpectedErrors(err, []string{"IncorrectDBInstanceState"}) || NeedRetry(err) {
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
		clickhouseService := ClickhouseService{client}
		stateConf := BuildStateConf([]string{}, []string{"Running"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, clickhouseService.ClickHouseDbClusterStateRefreshFunc(d.Id(), []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
		d.SetPartial("resource_group_id")
	}
	if d.HasChange("allocate_public_connection") {
		openPublicConnection := d.Get("allocate_public_connection").(bool)
		if openPublicConnection {
			action := "AllocateClusterPublicConnection"
			request := map[string]interface{}{
				"DBClusterId": d.Id(),
				"RegionId":    client.RegionId,
			}
			wait := incrementalWait(3*time.Second, 3*time.Second)
			err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
				response, err = client.RpcPost("clickhouse", "2019-11-11", action, nil, request, false)
				if err != nil {
					if IsExpectedErrors(err, []string{"IncorrectDBInstanceState"}) || NeedRetry(err) {
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
			clickhouseService := ClickhouseService{client}
			stateConf := BuildStateConf([]string{"NetAddressCreating"}, []string{"Running"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, clickhouseService.ClickHouseDbClusterStateRefreshFunc(d.Id(), []string{}))
			if _, err := stateConf.WaitForState(); err != nil {
				return WrapErrorf(err, IdMsg, d.Id())
			}
			d.SetPartial("allocate_public_connection")
		} else {
			action := "ReleaseClusterPublicConnection"
			//NetAddressDeleting
			request := map[string]interface{}{
				"DBClusterId": d.Id(),
				"RegionId":    client.RegionId,
			}
			wait := incrementalWait(3*time.Second, 3*time.Second)
			err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
				response, err = client.RpcPost("clickhouse", "2019-11-11", action, nil, request, false)
				if err != nil {
					if IsExpectedErrors(err, []string{"IncorrectDBInstanceState"}) || NeedRetry(err) {
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
			clickhouseService := ClickhouseService{client}
			stateConf := BuildStateConf([]string{"NetAddressDeleting"}, []string{"Running"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, clickhouseService.ClickHouseDbClusterStateRefreshFunc(d.Id(), []string{}))
			if _, err := stateConf.WaitForState(); err != nil {
				return WrapErrorf(err, IdMsg, d.Id())
			}
			d.SetPartial("allocate_public_connection")
		}
	}
	if v, ok := d.GetOk("cold_storage"); ok && v.(string) == "ENABLE" && d.HasChange("cold_storage") && !d.IsNewResource() {
		action := "CreateOSSStorage"
		wait := incrementalWait(3*time.Second, 3*time.Second)
		request := map[string]interface{}{
			"DBClusterId": d.Id(),
			"RegionId":    client.RegionId,
		}
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("clickhouse", "2019-11-11", action, nil, request, false)
			if err != nil {
				if IsExpectedErrors(err, []string{"IncorrectDBInstanceState"}) || NeedRetry(err) {
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
		clickhouseService := ClickhouseService{client}
		stateConf := BuildStateConf([]string{"CREATING"}, []string{"DISABLE", "ENABLE"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, clickhouseService.ClickHouseOSSStorageStateRefreshFunc(d.Id(), []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
		d.SetPartial("cold_storage")
	}

	d.Partial(false)
	return resourceAlicloudClickHouseDbClusterRead(d, meta)
}
func resourceAlicloudClickHouseDbClusterDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	action := "DeleteDBCluster"
	var response map[string]interface{}
	var err error
	request := map[string]interface{}{
		"DBClusterId": d.Id(),
	}
	if v, ok := d.GetOk("payment_type"); ok && v.(string) == "Subscription" {
		log.Printf("[WARN] Cannot destroy resourceClickHouseDbCluster. Because payment_type = 'Subscription'. Terraform will remove this resource from the state file, however resources may remain.")
		return nil
	}
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = client.RpcPost("clickhouse", "2019-11-11", action, nil, request, false)
		if err != nil {
			if IsExpectedErrors(err, []string{"IncorrectDBInstanceState"}) || NeedRetry(err) {
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
	return nil
}

func convertClickHouseDbClusterPaymentTypeRequest(source string) string {
	switch source {
	case "PayAsYouGo":
		return "Postpaid"
	case "Subscription":
		return "Prepaid"
	}
	return source
}

func convertClickHouseDbClusterPaymentTypeResponse(source string) string {
	switch source {
	case "Postpaid":
		return "PayAsYouGo"
	case "Prepaid":
		return "Subscription"
	}
	return source
}

func convertClickHouseDbClusterStorageTypeResponse(source string) string {
	switch source {
	case "CloudESSD":
		return "cloud_essd"
	case "CloudEfficiency":
		return "cloud_efficiency"
	case "CloudESSD_PL2":
		return "cloud_essd_pl2"
	case "CloudESSD_PL3":
		return "cloud_essd_pl3"

	}
	return source
}

func convertClickHouseDbClusterStorageTypeRequest(source string) string {
	switch source {
	case "cloud_essd":
		return "CloudESSD"
	case "cloud_efficiency":
		return "CloudEfficiency"
	case "cloud_essd_pl2":
		return "CloudESSD_PL2"
	case "cloud_essd_pl3":
		return "CloudESSD_PL3"

	}
	return source
}
