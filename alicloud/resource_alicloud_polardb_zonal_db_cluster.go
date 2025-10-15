// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"encoding/json"
	"fmt"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/polardb"
	"log"
	"strconv"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

type Pair struct {
	val string
	idx string
}

type DBNodeConfig struct {
	DBNodeClass    string `json:"db_node_class"`
	DBNodeRole     string `json:"db_node_role,omitempty"`
	HotReplicaMode string `json:"hot_replica_mode,omitempty"`
	ImciSwitch     string `json:"imci_switch,omitempty"`
}

type DBNodeAttribute struct {
	DBNodeId       string `json:"db_node_id"`
	DBNodeClass    string `json:"db_node_class"`
	DBNodeRole     string `json:"db_node_role"`
	HotReplicaMode string `json:"hot_replica_mode,omitempty"` //ENS不支持热备RO
	ImciSwitch     string `json:"imci_switch"`
}

func resourceAliCloudPolarDbZonalCluster() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudPolarDbZonalClusterCreate,
		Read:   resourceAliCloudPolarDbZonalClusterRead,
		Update: resourceAliCloudPolarDbZonalClusterUpdate,
		Delete: resourceAliCloudPolarDbZonalClusterDelete,
		Importer: &schema.ResourceImporter{
			State: func(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
				d.Set("creation_category", "___importing")
				return []*schema.ResourceData{d}, nil
			},
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(38 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},

		CustomizeDiff: func(d *schema.ResourceDiff, v interface{}) error {
			if d.HasChange("db_cluster_nodes_configs") {
				flagNodeAttributesUpdate := false
				flagNodeIdsUpdate := false
				nodeConfigsOld := make(map[string]DBNodeConfig)
				nodeConfigsNew := make(map[string]DBNodeConfig)
				configsOld, configsNew := d.GetChange("db_cluster_nodes_configs")
				for key, configJSON := range configsOld.(map[string]interface{}) {
					var dbNodeConfig DBNodeConfig
					err := json.Unmarshal([]byte(configJSON.(string)), &dbNodeConfig)
					if err != nil {
						return fmt.Errorf("failed to decode JSON for %s: %s", key, err)
					}
					if dbNodeConfig.ImciSwitch == "" {
						dbNodeConfig.ImciSwitch = "OFF"
					}
					if dbNodeConfig.DBNodeRole == "" {
						dbNodeConfig.DBNodeRole = "Reader"
					}
					nodeConfigsOld[key] = dbNodeConfig
				}
				for key, configJSON := range configsNew.(map[string]interface{}) {
					var dbNodeConfig DBNodeConfig
					err := json.Unmarshal([]byte(configJSON.(string)), &dbNodeConfig)
					if err != nil {
						return fmt.Errorf("failed to decode JSON for %s: %s", key, err)
					}
					if dbNodeConfig.ImciSwitch == "" {
						dbNodeConfig.ImciSwitch = "OFF"
					}
					if dbNodeConfig.DBNodeRole == "" {
						dbNodeConfig.DBNodeRole = "Reader"
					}
					nodeConfigsNew[key] = dbNodeConfig
				}
				for nodeKey, nodeConfig := range nodeConfigsNew {
					if v, exist := nodeConfigsOld[nodeKey]; !exist {
						flagNodeAttributesUpdate = true
						flagNodeIdsUpdate = true
						break
					} else {
						if v.DBNodeClass != nodeConfig.DBNodeClass || v.DBNodeRole != nodeConfig.DBNodeRole {
							flagNodeAttributesUpdate = true
						}
						if v.ImciSwitch != nodeConfig.ImciSwitch {
							flagNodeAttributesUpdate = true
							flagNodeIdsUpdate = true
						}
					}
				}
				for nodeKey, _ := range nodeConfigsOld {
					if _, exist := nodeConfigsNew[nodeKey]; !exist {
						flagNodeAttributesUpdate = true
						flagNodeIdsUpdate = true
						break
					}
				}
				if flagNodeAttributesUpdate {
					err := d.SetNewComputed("db_cluster_nodes_attributes")
					if err != nil {
						return err
					}
				}
				if flagNodeIdsUpdate {
					err := d.SetNewComputed("db_cluster_nodes_ids")
					if err != nil {
						return err
					}
				}
			}
			return nil
		},

		Schema: map[string]*schema.Schema{
			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"creation_category": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				ValidateFunc: StringInSlice([]string{"SENormal"}, false),
			},
			"description": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: StringLenBetween(2, 256),
			},
			"db_minor_version": {
				Type:         schema.TypeString,
				ValidateFunc: StringInSlice([]string{"8.0.1", "8.0.2"}, false),
				Default:      "8.0.2",
				ForceNew:     true,
				Optional:     true,
			},
			"db_node_class": {
				Type:     schema.TypeString,
				Required: true,
			},
			"db_type": {
				Type:         schema.TypeString,
				Default:      "MySQL",
				ValidateFunc: StringInSlice([]string{"MySQL"}, false),
				Optional:     true,
				ForceNew:     true,
			},
			"db_version": {
				Type:         schema.TypeString,
				ValidateFunc: StringInSlice([]string{"8.0"}, false),
				Optional:     true,
				Default:      "8.0",
				ForceNew:     true,
			},
			"ens_region_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"pay_type": {
				Type:         schema.TypeString,
				ValidateFunc: StringInSlice([]string{string(PostPaid), string(PrePaid)}, false),
				Optional:     true,
				ForceNew:     true,
				Default:      PrePaid,
			},
			"renewal_status": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  RenewNotRenewal,
				ValidateFunc: StringInSlice([]string{
					string(RenewAutoRenewal),
					string(RenewNormal),
					string(RenewNotRenewal)}, false),
				DiffSuppressFunc: polardbPostPaidDiffSuppressFunc,
			},
			"auto_renew_period": {
				Type:             schema.TypeInt,
				Optional:         true,
				Default:          1,
				ValidateFunc:     IntInSlice([]int{1, 2, 3, 6, 12, 24, 36}),
				DiffSuppressFunc: polardbPostPaidAndRenewDiffSuppressFunc,
			},
			"used_time": {
				Type:             schema.TypeInt,
				Default:          1,
				ValidateFunc:     IntInSlice([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 12, 24, 36}),
				Optional:         true,
				DiffSuppressFunc: polardbPostPaidDiffSuppressFunc,
			},
			"region_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"storage_pay_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"storage_space": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      20,
				ValidateFunc: IntBetween(20, 100000),
			},
			"storage_type": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "ESSDPL0",
				ValidateFunc: StringInSlice([]string{"ESSDPL0", "ESSDPL1"}, false),
				ForceNew:     true,
			},
			"target_minor_version": {
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
				Required: true,
				ForceNew: true,
			},
			"db_cluster_nodes_configs": {
				Type:     schema.TypeMap,
				Required: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"db_cluster_nodes_attributes": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"db_cluster_nodes_ids": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"cluster_version": {
				Type:             schema.TypeString,
				Optional:         true,
				Computed:         true,
				ForceNew:         true,
				DiffSuppressFunc: polardbDBClusterVersionDiffSuppressFunc,
			},
			"cluster_latest_version": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceAliCloudPolarDbZonalClusterCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "CreateDBCluster"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	request["RegionId"] = client.RegionId
	request["ClientToken"] = buildClientToken(action)

	if v, ok := d.GetOk("storage_type"); ok {
		request["StorageType"] = v
	} else {
		request["StorageType"] = "ESSDPL0"
	}
	if v, ok := d.GetOk("vswitch_id"); ok {
		request["VSwitchId"] = v
	}
	if v, ok := d.GetOk("db_minor_version"); ok {
		request["DBMinorVersion"] = v
	}
	request["DBNodeClass"] = d.Get("db_node_class")
	if v, ok := d.GetOk("description"); ok {
		request["DBClusterDescription"] = v
	}
	if v, ok := d.GetOk("creation_category"); ok {
		request["CreationCategory"] = v
	}
	request["DBType"] = d.Get("db_type")
	if v, ok := d.GetOk("target_minor_version"); ok {
		request["TargetMinorVersion"] = v
	}
	request["StoragePayType"] = "Prepaid"
	if v, ok := d.GetOk("storage_space"); ok {
		request["StorageSpace"] = v
	}
	if v, ok := d.GetOk("ens_region_id"); ok {
		request["EnsRegionId"] = v
	}
	if v, ok := d.GetOk("vpc_id"); ok {
		request["VPCId"] = v
	}
	request["DBVersion"] = d.Get("db_version")
	request["CloudProvider"] = "ENS"

	payType := Trim(d.Get("pay_type").(string))
	request["PayType"] = string(Postpaid)
	if payType == string(PrePaid) {
		request["PayType"] = string(Prepaid)
	}
	if PayType(request["PayType"].(string)) == Prepaid {
		period := d.Get("used_time").(int)
		request["UsedTime"] = strconv.Itoa(period)
		request["Period"] = string(Month)
		if period > 9 {
			request["UsedTime"] = strconv.Itoa(period / 12)
			request["Period"] = string(Year)
		}
		if d.Get("renewal_status").(string) != string(RenewNotRenewal) {
			request["AutoRenew"] = requests.Boolean(strconv.FormatBool(true))
		} else {
			request["AutoRenew"] = requests.Boolean(strconv.FormatBool(false))
		}
	}

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.RpcPost("polardb", "2017-08-01", action, query, request, true)
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_polardb_on_ens_db_cluster", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(response["DBClusterId"]))

	polarDbServiceV2 := PolarDbServiceV2{client}
	stateConf := BuildStateConf([]string{}, []string{"Running"}, d.Timeout(schema.TimeoutCreate), 10*time.Second, polarDbServiceV2.PolarDbZonalClusterStateRefreshFunc(d.Id(), []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return resourceAliCloudPolarDbZonalClusterUpdate(d, meta)
}

func resourceAliCloudPolarDbZonalClusterRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	polarDbServiceV2 := PolarDbServiceV2{client}

	objectObj, err := polarDbServiceV2.DescribePolarDbZonalCluster(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_polardb_on_ens_db_cluster DescribePolarDbZonalCluster Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	nodeConfigs, err := getNodeConfig(d)
	if err != nil {
		return err
	}
	nodeAttributes, err := getNodeAttributes(d)
	if err != nil {
		return err
	}

	isImporting := false
	if nodeIds, ok := d.GetOk("creation_category"); ok && nodeIds == "___importing" {
		isImporting = true
		d.Set("creation_category", nil)
	}
	if err = mergeNodeAttributeFromThreeWay(d, &nodeAttributes, &nodeConfigs, objectObj, isImporting); err != nil {
		return err
	}

	updatedNodeIds := make(map[string]interface{})
	for key, attribute := range nodeAttributes {
		updatedNodeIds[key] = attribute.DBNodeId
	}
	if err := d.Set("db_cluster_nodes_ids", updatedNodeIds); err != nil {
		return fmt.Errorf("failed to set nodes_ids: %s", err)
	}

	d.Set("create_time", objectObj.CreationTime)
	d.Set("creation_category", objectObj.Category)
	d.Set("description", objectObj.DBClusterDescription)
	d.Set("db_type", objectObj.DBType)
	d.Set("db_version", objectObj.DBVersion)
	d.Set("pay_type", getChargeType(objectObj.PayType))
	d.Set("region_id", objectObj.RegionId)
	d.Set("storage_pay_type", objectObj.StoragePayType)
	d.Set("storage_type", objectObj.StorageType)
	d.Set("vswitch_id", objectObj.VSwitchId)
	d.Set("vpc_id", objectObj.VPCId)
	d.Set("ens_region_id", objectObj.ZoneIds)

	dBNodesRaw := objectObj.DBNodes
	dbNodeMaps := make([]map[string]interface{}, 0)
	if dBNodesRaw != nil {
		for _, dBNodesChildRaw := range dBNodesRaw {
			dbNodeMap := make(map[string]interface{})
			dbNodeMap["db_node_id"] = dBNodesChildRaw.DBNodeId
			dbNodeMap["target_class"] = dBNodesChildRaw.DBNodeClass

			dbNodeMaps = append(dbNodeMaps, dbNodeMap)
		}
	}

	d.Set("storage_space", objectObj.StorageSpace/1024/1024)
	objectRaw, err := polarDbServiceV2.DescribeZonalClusterDescribeDBClusterVersionZonal(d.Id())
	if err != nil && !NotFoundError(err) {
		return WrapError(err)
	}
	d.Set("cluster_version", objectRaw["DBRevisionVersion"])
	d.Set("cluster_latest_version", objectRaw["DBLatestVersion"])
	d.Set("db_minor_version", objectRaw["DBMinorVersion"])

	if objectObj.PayType == string(Prepaid) {
		objectRaw, err = polarDbServiceV2.DescribeZonalClusterDescribeAutoRenewAttribute(d.Id())
		if err != nil && !NotFoundError(err) {
			return WrapError(err)
		}

		renewPeriod := 1
		if objectRaw["Duration"] != nil {
			duration, _ := strconv.Atoi(string(objectRaw["Duration"].(json.Number)))
			renewPeriod = duration
		}
		if objectRaw != nil && objectRaw["PeriodUnit"] == string(Year) {
			renewPeriod = renewPeriod * 12
		}
		d.Set("renewal_status", objectRaw["RenewalStatus"])
		d.Set("auto_renew_period", renewPeriod)
	}

	return nil
}

func resourceAliCloudPolarDbZonalClusterUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	d.Partial(true)
	polarDbServiceV2 := PolarDbServiceV2{client}

	var err error
	if (d.Get("pay_type").(string) == string(PrePaid)) &&
		(d.HasChange("renewal_status") || d.HasChange("auto_renew_period")) {
		status := d.Get("renewal_status").(string)
		var duration string
		var periodUnit string
		if status == string(RenewAutoRenewal) {
			period := d.Get("auto_renew_period").(int)
			duration = strconv.Itoa(period)
			periodUnit = string(Month)
			if period > 9 {
				duration = strconv.Itoa(period / 12)
				periodUnit = string(Year)
			}
		} else {
			duration = "0"
			periodUnit = string(Month)
		}
		//wait asynchronously cluster payType
		if err := polarDbServiceV2.WaitForPolarDBPayType(d.Id(), "Prepaid", DefaultLongTimeout); err != nil {
			return WrapError(err)
		}
		_, err = polarDbServiceV2.ModifyAutoRenewAttribute(client.RegionId, d.Id(), status, duration, periodUnit)
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), "ModifyAutoRenewAttribute", AlibabaCloudSdkGoERROR)
		}
		d.SetPartial("renewal_status")
		d.SetPartial("auto_renew_period")
	}

	if d.HasChange("db_cluster_nodes_configs") {
		nodeConfigs, err := getNodeConfig(d)
		if err != nil {
			return err
		}
		nodeAttributes, err := getNodeAttributes(d)
		if err != nil {
			return err
		}
		clusterAttribute, err := polarDbServiceV2.DescribePolarDbZonalCluster(d.Id())
		if err != nil {
			return WrapError(err)
		}

		nodeOperationPackage, err := reconcileNodeAttributeFromThreeWay(d, &nodeAttributes, &nodeConfigs, clusterAttribute)
		if err != nil {
			return err
		}

		log.Printf("reconcileNodeAttributeFromThreeWay  %v %s", nodeOperationPackage, d.Id())

		if len(*nodeOperationPackage.NodesToDelete) > 0 {
			if err = deleteDBNodes(d, nodeAttributes, meta, *nodeOperationPackage.NodesToDelete); err != nil {
				return err
			}
		}

		var addNormalNodeList []Pair
		var addImciNodeList []Pair
		for _, node := range *nodeOperationPackage.NodesToAdd {
			if node.IMCISwitch == "ON" {
				addImciNodeList = append(addImciNodeList, Pair{val: node.DBNodeClass, idx: node.Index})
				continue
			}
			addNormalNodeList = append(addNormalNodeList, Pair{val: node.DBNodeClass, idx: node.Index})
		}

		if _, err = createDBNodes(d, meta, addNormalNodeList, "OFF"); err != nil {
			return err
		}

		if _, err = createDBNodes(d, meta, addImciNodeList, "ON"); err != nil {
			return err
		}

		if err = modifyDBNodesClass(d, meta, *nodeOperationPackage.NodesToModifyClass); err != nil {
			return err
		}

		if err = modifyDBNodesRole(d, meta, nodeOperationPackage.NodesToSwitchOver); err != nil {
			return err
		}
	}

	if d.HasChange("storage_space") {
		action := "ModifyDBClusterStorageSpace"
		storageSpace := d.Get("storage_space").(int)
		request := map[string]interface{}{
			"DBClusterId":   d.Id(),
			"StorageSpace":  storageSpace,
			"CloudProvider": "ENS",
		}
		if v, ok := d.GetOk("planned_start_time"); ok && v.(string) != "" {
			request["PlannedStartTime"] = v
		}
		if v, ok := d.GetOk("planned_end_time"); ok && v.(string) != "" {
			request["PlannedEndTime"] = v
		}
		if v, ok := d.GetOk("sub_category"); ok && v.(string) != "" {
			request["SubCategory"] = convertPolarDBSubCategoryUpdateRequest(v.(string))
		}
		//request["ClientToken"] = buildClientToken(action)
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err := resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err := client.RpcPost("polardb", "2017-08-01", action, nil, request, false)
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				addDebug(action, response, request)
			}
			return nil
		})
		if err != nil {
			if IsExpectedErrors(err, []string{"InvalidDBCluster.NotFound"}) {
				return nil
			}
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, ProviderERROR)
		}
		// wait cluster status change from StorageExpanding to running
		stateConf := BuildStateConf([]string{}, []string{"Running"}, d.Timeout(schema.TimeoutUpdate), 10*time.Second, polarDbServiceV2.PolarDbZonalClusterStateRefreshFunc(d.Id(), []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
	}

	if d.HasChange("description") {
		description := d.Get("description").(string)
		if err := polarDbServiceV2.ModifyDBClusterDescriptionZonal(d.Id(), description); err != nil {
			return err
		}
	}

	if d.HasChange("cluster_version") {
		clusterVersion := d.Get("cluster_version").(string)
		objectRaw, err := polarDbServiceV2.DescribeZonalClusterDescribeDBClusterVersionZonal(d.Id())
		if err != nil && !NotFoundError(err) {
			return WrapError(err)
		}
		if clusterVersion != objectRaw["DBLatestVersion"] {
			return fmt.Errorf("must upgrade to the latest version %s", objectRaw["DBLatestVersion"])
		}

		if err := polarDbServiceV2.UpgradeDBClusterVersionZonal(d.Id()); err != nil {
			return err
		}
		// wait cluster status change from StorageExpanding to running
		stateConf := BuildStateConf([]string{}, []string{"Running"}, d.Timeout(schema.TimeoutUpdate), 10*time.Second, polarDbServiceV2.PolarDbZonalClusterStateRefreshFunc(d.Id(), []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
	}

	d.Partial(false)
	return resourceAliCloudPolarDbZonalClusterRead(d, meta)
}

func resourceAliCloudPolarDbZonalClusterDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	polarDbServiceV2 := PolarDbServiceV2{client}
	if d.Get("pay_type").(string) == string(PostPaid) {
		action := "DeleteDBCluster"
		var request map[string]interface{}
		var response map[string]interface{}
		query := make(map[string]interface{})
		var err error
		request = make(map[string]interface{})
		request["DBClusterId"] = d.Id()

		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
			response, err = client.RpcPost("polardb", "2017-08-01", action, query, request, true)

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
	} else {
		client := meta.(*connectivity.AliyunClient)
		action := "RefundInstance"
		var request map[string]interface{}
		var response map[string]interface{}
		var err error
		query := make(map[string]interface{})
		request = make(map[string]interface{})
		request["InstanceId"] = d.Id()

		request["ClientToken"] = buildClientToken(action)
		request["ProductCode"] = "polardb"
		request["ProductType"] = "polardb_edgesub_public_cn"
		if client.IsInternationalAccount() {
			request["ProductType"] = "polardb_edgesub_public_intl"
		}
		request["ImmediatelyRelease"] = "1"
		var endpoint string
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
			response, err = client.RpcPostWithEndpoint("BssOpenApi", "2017-12-14", action, query, request, true, endpoint)
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				if !client.IsInternationalAccount() && IsExpectedErrors(err, []string{"NotApplicable"}) {
					request["ProductType"] = "polardb_edgesub_public_cn"
					endpoint = connectivity.BssOpenAPIEndpointInternational
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})
		addDebug(action, response, request)

		if err != nil {
			if IsExpectedErrors(err, []string{"ResourceNotExists"}) || NotFoundError(err) {
				return nil
			}
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
	}
	if err := polarDbServiceV2.WaitForPolarDBDeleted(d.Id(), 600); err != nil {
		return err
	}
	return nil
}

func createDBNodes(d *schema.ResourceData, meta interface{}, addNodeList []Pair, ImciSwitch string) (map[string]string, error) {
	if len(addNodeList) == 0 {
		return nil, nil
	}
	nodeAttributesId := make(map[string]string)
	client := meta.(*connectivity.AliyunClient)
	polarDbServiceV2 := PolarDbServiceV2{client}

	dbNode := make([]polardb.CreateDBNodesDBNode, len(addNodeList))
	for i, nodeClass := range addNodeList {
		dbNode[i] = polardb.CreateDBNodesDBNode{
			TargetClass: nodeClass.val,
		}
	}
	var err error
	response, err := polarDbServiceV2.CreateDBNodes(client.RegionId, d.Id(), ImciSwitch, dbNode)
	if err != nil {
		return nil, WrapErrorf(
			err, DefaultErrorMsg, d.Id(), "CreateDBNodes", AlibabaCloudSdkGoERROR)
	}

	if err != nil {
		if IsExpectedErrors(err, []string{"InvalidDBNodeNum.Malformed"}) {
			return createDBNodesSerial(d, meta, addNodeList, ImciSwitch)
		}
		return nil, WrapErrorf(
			err, DefaultErrorMsg, d.Id(), "CreateDBNodes", AlibabaCloudSdkGoERROR)
	}
	dbNodeId, _ := response["DBNodeIds"].([]interface{})
	for i, nodeId := range dbNodeId {
		nodeAttributesId[addNodeList[i].idx] = nodeId.(string)
	}
	addDebug("CreateDBNodes", nodeAttributesId)
	stateConf := BuildStateConf([]string{}, []string{"Running"}, d.Timeout(schema.TimeoutUpdate), 10*time.Second, polarDbServiceV2.PolarDbZonalClusterStateRefreshFunc(d.Id(), []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return nil, WrapErrorf(err, IdMsg, d.Id())
	}
	return nodeAttributesId, nil
}

func createDBNodesSerial(d *schema.ResourceData, meta interface{}, addNodeList []Pair, ImciSwitch string) (map[string]string, error) {
	if len(addNodeList) == 0 {
		return nil, nil
	}
	nodeAttributesId := make(map[string]string)
	client := meta.(*connectivity.AliyunClient)
	polarDbServiceV2 := PolarDbServiceV2{client}
	for i, nodeClass := range addNodeList {
		dbNode := make([]polardb.CreateDBNodesDBNode, 1)
		dbNode[0] = polardb.CreateDBNodesDBNode{
			TargetClass: nodeClass.val,
		}

		var err error
		response, err := polarDbServiceV2.CreateDBNodes(client.RegionId, d.Id(), ImciSwitch, dbNode)
		if err != nil {
			return nil, WrapErrorf(
				err, DefaultErrorMsg, d.Id(), "CreateDBNodes", AlibabaCloudSdkGoERROR)
		}
		dbNodeId, _ := response["DBNodeIds"].([]interface{})
		nodeAttributesId[addNodeList[i].idx] = dbNodeId[0].(string)
		stateConf := BuildStateConf([]string{}, []string{"Running"}, d.Timeout(schema.TimeoutCreate), 10*time.Second, polarDbServiceV2.PolarDbZonalClusterStateRefreshFunc(d.Id(), []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return nil, WrapErrorf(err, IdMsg, d.Id())
		}
	}
	return nodeAttributesId, nil
}

func deleteDBNodes(d *schema.ResourceData, nodeAttributes map[string]DBNodeAttribute, meta interface{}, removeNodeList []Pair) error {
	if len(removeNodeList) == 0 {
		return nil
	}
	client := meta.(*connectivity.AliyunClient)
	polarDbServiceV2 := PolarDbServiceV2{client}
	removeNodes := make([]string, len(removeNodeList))
	for i, node := range removeNodeList {
		removeNodes[i] = node.val
	}
	var err error
	_, err = polarDbServiceV2.DeleteDBNodes(client.RegionId, d.Id(), removeNodes)

	if err != nil {
		if IsExpectedErrors(err, []string{"InvalidDBNodeNum.Malformed"}) {
			return deleteDBNodesSerial(d, meta, removeNodeList)
		}
		return WrapErrorf(
			err, DefaultErrorMsg, d.Id(), "DeleteDBNodes", AlibabaCloudSdkGoERROR)
	}
	for _, node := range removeNodeList {
		delete(nodeAttributes, node.idx)
	}
	stateConf := BuildStateConf([]string{}, []string{"Running"}, d.Timeout(schema.TimeoutDelete), 10*time.Second, polarDbServiceV2.PolarDbZonalClusterStateRefreshFunc(d.Id(), []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}
	return nil
}

func deleteDBNodesSerial(d *schema.ResourceData, meta interface{}, removeNodeList []Pair) error {
	if len(removeNodeList) == 0 {
		return nil
	}
	client := meta.(*connectivity.AliyunClient)
	polarDbServiceV2 := PolarDbServiceV2{client}
	for _, node := range removeNodeList {
		removeNodes := make([]string, 1)
		removeNodes[0] = node.val
		var err error
		_, err = polarDbServiceV2.DeleteDBNodes(client.RegionId, d.Id(), removeNodes)
		if err != nil {
			return WrapErrorf(
				err, DefaultErrorMsg, d.Id(), "DeleteDBNodes", AlibabaCloudSdkGoERROR)
		}
		stateConf := BuildStateConf([]string{}, []string{"Running"}, d.Timeout(schema.TimeoutDelete), 10*time.Second, polarDbServiceV2.PolarDbZonalClusterStateRefreshFunc(d.Id(), []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
	}
	return nil
}

func modifyDBNodesClass(d *schema.ResourceData, meta interface{}, modifyNodeClassList []polardb.ModifyDBNodesClassDBNode) error {
	if len(modifyNodeClassList) == 0 {
		return nil
	}
	client := meta.(*connectivity.AliyunClient)
	polarDbServiceV2 := PolarDbServiceV2{client}
	_, err := polarDbServiceV2.ModifyDBNodesClass(client.RegionId, d.Id(), "Upgrade", modifyNodeClassList)
	if err != nil {
		return WrapErrorf(
			err, DefaultErrorMsg, d.Id(), "ModifyDBNodesClass", AlibabaCloudSdkGoERROR)
	}
	stateConf := BuildStateConf([]string{}, []string{"Running"}, d.Timeout(schema.TimeoutUpdate), 10*time.Second, polarDbServiceV2.PolarDbZonalClusterStateRefreshFunc(d.Id(), []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}
	return nil
}

func modifyDBNodesRole(d *schema.ResourceData, meta interface{}, pair *Pair) error {
	if pair == nil {
		return nil
	}
	client := meta.(*connectivity.AliyunClient)
	polarDbServiceV2 := PolarDbServiceV2{client}
	err := polarDbServiceV2.FailoverDBClusterZonal(d.Id(), pair.val)
	if err != nil {
		return WrapErrorf(
			err, DefaultErrorMsg, d.Id(), "FailoverDBClusterZonal", AlibabaCloudSdkGoERROR)
	}
	stateConf := BuildStateConf([]string{}, []string{"Running"}, d.Timeout(schema.TimeoutUpdate), 10*time.Second, polarDbServiceV2.PolarDbZonalClusterStateRefreshFunc(d.Id(), []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}
	return nil
}

func flushNodeAttributes(d *schema.ResourceData, nodeAttributes map[string]DBNodeAttribute) error {
	updatedAttributes := make(map[string]interface{})
	for key, attribute := range nodeAttributes {
		attributesJSON, err := json.Marshal(attribute)
		if err != nil {
			return fmt.Errorf("failed to encode JSON for %s: %s", key, err)
		}
		updatedAttributes[key] = string(attributesJSON)
	}
	if err := d.Set("db_cluster_nodes_attributes", updatedAttributes); err != nil {
		return fmt.Errorf("failed to set node_attributes: %s", err)
	}
	return nil
}

func flushNodeConfigs(d *schema.ResourceData, nodeConfigs map[string]DBNodeConfig) error {
	updatedConfigs := make(map[string]interface{})
	for key, config := range nodeConfigs {
		configsJSON, err := json.Marshal(config)
		if err != nil {
			return fmt.Errorf("failed to encode JSON for %s: %s", key, err)
		}
		updatedConfigs[key] = string(configsJSON)
	}
	if err := d.Set("db_cluster_nodes_configs", updatedConfigs); err != nil {
		return fmt.Errorf("failed to set node_configs: %s", err)
	}
	return nil
}

func getNodeConfig(d *schema.ResourceData) (map[string]DBNodeConfig, error) {
	configsMap := d.Get("db_cluster_nodes_configs").(map[string]interface{})
	nodeConfigs := make(map[string]DBNodeConfig)
	for key, configJSON := range configsMap {
		var dbNodeConfig DBNodeConfig
		err := json.Unmarshal([]byte(configJSON.(string)), &dbNodeConfig)
		if err != nil {
			return nil, fmt.Errorf("failed to decode JSON for %s: %s", key, err)
		}
		if dbNodeConfig.ImciSwitch == "" {
			dbNodeConfig.ImciSwitch = "OFF"
		}
		if dbNodeConfig.DBNodeRole == "" {
			dbNodeConfig.DBNodeRole = "Reader"
		}
		nodeConfigs[key] = dbNodeConfig
	}
	return nodeConfigs, nil
}

func getNodeAttributes(d *schema.ResourceData) (map[string]DBNodeAttribute, error) {
	attributesOld, _ := d.GetChange("db_cluster_nodes_attributes")
	var attributesMap map[string]interface{}
	attributesMap = attributesOld.(map[string]interface{})
	nodeAttributes := make(map[string]DBNodeAttribute)
	for key, attributeJSON := range attributesMap {
		var dbNodeAttribute DBNodeAttribute
		err := json.Unmarshal([]byte(attributeJSON.(string)), &dbNodeAttribute)
		if err != nil {
			return nil, fmt.Errorf("failed to decode JSON for %s: %s", key, err)
		}
		nodeAttributes[key] = dbNodeAttribute
	}
	return nodeAttributes, nil
}

type AddRODBNode struct {
	Index       string
	DBNodeClass string
	IMCISwitch  string
}

type NodeOperationPackage struct {
	NodesToDelete      *[]Pair
	NodesToAdd         *[]AddRODBNode
	NodesToModifyClass *[]polardb.ModifyDBNodesClassDBNode
	NodesToSwitchOver  *Pair
}

func reconcileNodeAttributeFromThreeWay(d *schema.ResourceData, nodeAttributes *map[string]DBNodeAttribute, nodeConfigs *map[string]DBNodeConfig,
	serverResponse *DescribeDBClusterAttributeResponse) (*NodeOperationPackage, error) {
	var err error
	nodeAttributes, err = mergeFromServerResponse(nodeAttributes, nodeConfigs, serverResponse, true, false)
	if err != nil {
		return nil, err
	}

	nodeOperationPackage := NodeOperationPackage{}
	// merge from config
	removeNodeList := make([]Pair, 0)
	addNodeList := make([]AddRODBNode, 0)
	modifyNodeClassList := make([]polardb.ModifyDBNodesClassDBNode, 0)

	normalClassNodeNumToDelete := 0
	normalClassNodeConfigNum := 0
	normalClassNodeAttributeNum := 0
	for key, config := range *nodeConfigs {
		imciSwitch := config.ImciSwitch
		if imciSwitch == "" {
			imciSwitch = "OFF"
		}

		if config.DBNodeClass == d.Get("db_node_class") && imciSwitch == "OFF" {
			normalClassNodeAttributeNum++
		}
		if _, exists := (*nodeAttributes)[key]; exists {
			if config.DBNodeClass != (*nodeAttributes)[key].DBNodeClass {
				modifyNodeClassList = append(modifyNodeClassList, polardb.ModifyDBNodesClassDBNode{DBNodeId: (*nodeAttributes)[key].DBNodeId, TargetClass: config.DBNodeClass})
			}
			if config.DBNodeRole == "Writer" && config.DBNodeRole != (*nodeAttributes)[key].DBNodeRole {
				nodeOperationPackage.NodesToSwitchOver = &Pair{val: (*nodeAttributes)[key].DBNodeId, idx: config.DBNodeRole}
			}
			continue
		}
		dbNodeClass := config.DBNodeClass
		addNode := AddRODBNode{DBNodeClass: dbNodeClass, Index: key, IMCISwitch: imciSwitch}
		addNodeList = append(addNodeList, addNode)
	}

	for key, attribute := range *nodeAttributes {
		if attribute.DBNodeRole == "Writer" {
			continue
		}
		if attribute.DBNodeClass == d.Get("db_node_class") && attribute.ImciSwitch == "OFF" {
			normalClassNodeConfigNum++
		}
		if _, exists := (*nodeConfigs)[key]; !exists {
			if (*nodeAttributes)[key].DBNodeClass == d.Get("db_node_class") && (*nodeAttributes)[key].ImciSwitch == "OFF" {
				normalClassNodeNumToDelete++
			}
			removeNodeList = append(removeNodeList, Pair{val: (*nodeAttributes)[key].DBNodeId, idx: key})
		}
	}
	nodeOperationPackage.NodesToDelete = &removeNodeList
	nodeOperationPackage.NodesToAdd = &addNodeList
	nodeOperationPackage.NodesToModifyClass = &modifyNodeClassList

	return &nodeOperationPackage, nil
}

func mergeNodeAttributeFromThreeWay(d *schema.ResourceData, nodeAttributes *map[string]DBNodeAttribute, nodeConfigs *map[string]DBNodeConfig,
	serverResponse *DescribeDBClusterAttributeResponse, isImporting bool) error {
	var err error
	nodeAttributes, err = mergeFromServerResponse(nodeAttributes, nodeConfigs, serverResponse, false, isImporting)
	if err != nil {
		return err
	}

	// sync state from server response after delete node、modify class
	for key, nodeAttribute := range *nodeAttributes {
		foundNode := false
		for _, nodeStatus := range serverResponse.DBNodes {
			if nodeStatus.DBNodeId == nodeAttribute.DBNodeId {
				nodeAttribute.DBNodeClass = nodeStatus.DBNodeClass
				nodeAttribute.ImciSwitch = nodeStatus.ImciSwitch
				nodeAttribute.DBNodeRole = nodeStatus.DBNodeRole
				(*nodeAttributes)[key] = nodeAttribute
				if _, exist := (*nodeConfigs)[key]; exist {
					nodeConfig := (*nodeConfigs)[key]
					nodeConfig.DBNodeClass = nodeAttribute.DBNodeClass
					if !(nodeConfig.ImciSwitch == "" && nodeAttribute.ImciSwitch == "OFF") {
						nodeConfig.ImciSwitch = nodeAttribute.ImciSwitch
					}
					nodeConfig.DBNodeRole = nodeAttribute.DBNodeRole
					(*nodeConfigs)[key] = nodeConfig
				} else if isImporting {
					(*nodeConfigs)[key] = DBNodeConfig{
						DBNodeClass: nodeAttribute.DBNodeClass,
						DBNodeRole:  nodeAttribute.DBNodeRole,
					}
				}
				foundNode = true
				break
			}
		}
		if !foundNode {
			delete(*nodeAttributes, key)
		}
	}

	if err = flushNodeAttributes(d, *nodeAttributes); err != nil {
		return err
	}

	if err = flushNodeConfigs(d, *nodeConfigs); err != nil {
		return err
	}
	return nil
}

func mergeFromServerResponse(nodeAttributes *map[string]DBNodeAttribute, nodeConfigs *map[string]DBNodeConfig,
	serverResponse *DescribeDBClusterAttributeResponse, reportOrphan, isImporting bool) (*map[string]DBNodeAttribute, error) {
	for i, dbNode := range serverResponse.DBNodes {
		isCaught := false
		for _, nodeAttribute := range *nodeAttributes {
			if nodeAttribute.DBNodeId == dbNode.DBNodeId {
				isCaught = true
				break
			}
		}
		if isCaught {
			continue
		}

		imciSwitchServer := dbNode.ImciSwitch
		if imciSwitchServer == "" {
			imciSwitchServer = "OFF"
		}

		//catch nodes from server with new nodes claimed in config
		for nodeKey, nodeConfig := range *nodeConfigs {
			if _, exists := (*nodeAttributes)[nodeKey]; exists {
				continue
			}
			imciSwitch := nodeConfig.ImciSwitch
			if imciSwitch == "" {
				imciSwitch = "OFF"
			}
			if nodeConfig.DBNodeClass == dbNode.DBNodeClass &&
				nodeConfig.DBNodeRole == dbNode.DBNodeRole && imciSwitch == imciSwitchServer {
				isCaught = true
				nodeAttribute := (*nodeAttributes)[nodeKey]
				nodeAttribute.DBNodeId = dbNode.DBNodeId
				nodeAttribute.DBNodeClass = dbNode.DBNodeClass
				nodeAttribute.DBNodeRole = dbNode.DBNodeRole
				nodeAttribute.ImciSwitch = imciSwitchServer
				(*nodeAttributes)[nodeKey] = nodeAttribute
				log.Printf("catch node %s %v", nodeKey, nodeAttributes)
				break
			}
		}

		if !isCaught && isImporting {
			key := "db_node_" + strconv.Itoa(i+1)
			nodeAttribute := (*nodeAttributes)[key]
			nodeAttribute.DBNodeId = dbNode.DBNodeId
			nodeAttribute.DBNodeClass = dbNode.DBNodeClass
			nodeAttribute.DBNodeRole = dbNode.DBNodeRole
			nodeAttribute.ImciSwitch = imciSwitchServer
			(*nodeAttributes)[key] = nodeAttribute
			log.Printf("imported node %s", dbNode.DBNodeId)
		}

		//orphan node
		if !isCaught && reportOrphan {
			log.Printf("found orphan node %s %v", dbNode.DBNodeId, nodeConfigs)
			return nil, fmt.Errorf("orphan node from server %s", dbNode.DBNodeId)
		}
	}
	return nodeAttributes, nil
}
