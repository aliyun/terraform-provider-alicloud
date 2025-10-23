package alicloud

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/polardb"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAlicloudPolarDBZonalEndpoint() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudPolarDBZonalEndpointCreate,
		Read:   resourceAlicloudPolarDBZonalEndpointRead,
		Update: resourceAlicloudPolarDBZonalEndpointUpdate,
		Delete: resourceAlicloudPolarDBZonalEndpointDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"db_cluster_id": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},
			"endpoint_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: StringInSlice([]string{"Custom", "Primary", "Cluster"}, false),
				Computed:     true,
			},
			"nodes": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"nodes_key": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					oldKeyList, newKeyList := d.GetChange("nodes_key")
					oldIdsMap, newIdsMap := d.GetChange("db_cluster_nodes_ids")
					newkeySet := newKeyList.(*schema.Set)
					oldkeySet := oldKeyList.(*schema.Set)
					newSet := make(map[string]bool)

					for _, v := range newkeySet.List() {
						nodeId := (newIdsMap.(map[string]interface{}))[v.(string)]
						if nodeId == nil {
							nodeId = ""
						}
						newSet[nodeId.(string)] = true
					}

					for _, v := range oldkeySet.List() {
						nodeId := (oldIdsMap.(map[string]interface{}))[v.(string)]
						if nodeId == nil {
							nodeId = ""
						}
						if !newSet[nodeId.(string)] {
							return false
						}
					}
					return true
				},
			},
			"db_cluster_nodes_ids": {
				Type:     schema.TypeMap,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					oldIdsMap, newIdsMap := d.GetChange("db_cluster_nodes_ids")
					if len(newIdsMap.(map[string]interface{})) == 0 {
						return true
					}

					newSet := make(map[string]bool)
					for _, v := range newIdsMap.(map[string]interface{}) {
						newSet[v.(string)] = true
					}

					for _, v := range oldIdsMap.(map[string]interface{}) {
						if !newSet[v.(string)] {
							return false
						}
					}
					return true
				},
			},
			"read_write_mode": {
				Type:         schema.TypeString,
				ValidateFunc: StringInSlice([]string{"ReadWrite", "ReadOnly"}, false),
				Optional:     true,
			},
			"auto_add_new_nodes": {
				Type:         schema.TypeString,
				ValidateFunc: StringInSlice([]string{"Enable", "Disable"}, false),
				Computed:     true,
				Optional:     true,
			},
			"endpoint_config": {
				Type:     schema.TypeMap,
				Optional: true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					value := d.Get("endpoint_config").(map[string]interface{})
					if value == nil || len(value) == 0 {
						return true
					}
					return false
				},
			},
			"net_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: StringInSlice([]string{"Private"}, false),
				Computed:     true,
			},
			"db_endpoint_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"db_endpoint_description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"connection_prefix": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"port": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"vpc_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"vswitch_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
		},
	}
}

func resourceAlicloudPolarDBZonalEndpointCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	polarDbServiceV2 := PolarDbServiceV2{client}
	clusterId := d.Get("db_cluster_id").(string)
	endpointType := "Custom"
	if d.Get("endpoint_type") != nil && d.Get("endpoint_type") != "" {
		endpointType = d.Get("endpoint_type").(string)
	}
	dbEndpointDescription := d.Get("db_endpoint_description").(string)
	request := CreateCreateDBClusterEndpointRequest()
	request.RegionId = client.RegionId
	request.DBClusterId = clusterId
	request.EndpointType = endpointType
	request.DBEndpointDescription = dbEndpointDescription
	if d.Get("vpc_id") != nil {
		request.VPCId = d.Get("vpc_id").(string)
	}
	if d.Get("vswitch_id") != nil {
		request.VPCId = d.Get("vpc_id").(string)
		request.VSwitchId = d.Get("vswitch_id").(string)
	}
	nodesKeyInter, nodesKeyOk := d.GetOk("nodes_key")
	nodesIdsMapInter, nodesIdsMapOk := d.GetOk("db_cluster_nodes_ids")
	if nodesKeyOk && nodesIdsMapOk {
		nodesIdsMap := nodesIdsMapInter.(map[string]interface{})
		var dbNodeParts []string
		nodesKeys := expandStringList(nodesKeyInter.(*schema.Set).List())
		for _, nodeKey := range nodesKeys {
			if v, exists := nodesIdsMap[nodeKey]; !exists {
				return WrapError(fmt.Errorf("node %s not found in db_cluster_nodes_ids", nodeKey))
			} else {
				dbNodeParts = append(dbNodeParts, v.(string))
			}
		}
		dbNodes := strings.Join(dbNodeParts, ",")
		request.Nodes = dbNodes
	}
	if readWriteMode, ok := d.GetOk("read_write_mode"); ok {
		request.ReadWriteMode = readWriteMode.(string)
	}
	if autoAddNewNodes, ok := d.GetOk("auto_add_new_nodes"); ok && autoAddNewNodes != "" {
		request.AutoAddNewNodes = autoAddNewNodes.(string)
	} else {
		request.AutoAddNewNodes = "Enable"
	}
	if endpointConfig, ok := d.GetOk("endpoint_config"); ok {
		endpointConfig, err := json.Marshal(endpointConfig)
		if err != nil {
			return WrapError(err)
		}
		request.EndpointConfig = string(endpointConfig)
	}

	endpoints, err := polarDbServiceV2.DescribeDBClusterEndpointsZonalList(client.RegionId, clusterId)
	if err != nil {
		return WrapError(err)
	}
	oldEndpoints := make([]interface{}, 0)
	for _, value := range *endpoints {
		oldEndpoints = append(oldEndpoints, value.DBEndpointId)
	}

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		err = polarDbServiceV2.CreateDBClusterEndpointZonal(request)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_polardb_endpoint", request.GetActionName(), AlibabaCloudSdkGoERROR)
	}

	oldEndpointIds := schema.NewSet(schema.HashString, oldEndpoints)
	dbEndpointId, err := polarDbServiceV2.WaitForPolarDBEndpoints(d, Active, oldEndpointIds, endpointType, DefaultTimeoutMedium)
	if err != nil {
		return WrapError(err)
	}
	d.SetId(fmt.Sprintf("%s%s%s", clusterId, COLON_SEPARATED, dbEndpointId))

	stateConf := BuildStateConf([]string{}, []string{"Running"}, d.Timeout(schema.TimeoutUpdate), 10*time.Second, polarDbServiceV2.PolarDbZonalClusterStateRefreshFunc(clusterId, []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return resourceAlicloudPolarDBZonalEndpointUpdate(d, meta)
}

func resourceAlicloudPolarDBZonalEndpointRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	polarDBServiceV2 := PolarDbServiceV2{client}
	parts, errParse := ParseResourceId(d.Id(), 2)
	if errParse != nil {
		return WrapError(errParse)
	}
	dbClusterId := parts[0]
	dbEndpointId := parts[1]
	object, err := polarDBServiceV2.DescribeDBClusterEndpointsZonal(d.Id())
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	d.Set("db_cluster_id", dbClusterId)
	d.Set("db_endpoint_id", dbEndpointId)
	d.Set("endpoint_type", object.EndpointType)
	d.Set("db_endpoint_description", object.DBEndpointDescription)
	nodes := strings.Split(object.Nodes, ",")
	nodesKeyList := make([]string, 0)
	if nodesIdsMap, nodesIdsMapOk := d.GetOk("db_cluster_nodes_ids"); nodesIdsMapOk {
		for _, node := range nodes {
			isMatched := false
			for nodeKey, nodeId := range nodesIdsMap.(map[string]interface{}) {
				if node == nodeId.(string) {
					nodesKeyList = append(nodesKeyList, nodeKey)
					isMatched = true
					continue
				}
			}
			if !isMatched {
				return WrapError(fmt.Errorf("node %s not found in db_cluster_nodes_ids", node))
			}
		}
		d.Set("nodes_key", nodesKeyList)
	} else {
		var dbClusterNodesIds map[string]interface{}
		dbClusterNodesIds = make(map[string]interface{}, 0)
		for i, node := range nodes {
			nodeKey := "db_node_" + strconv.Itoa(i+1)
			nodesKeyList = append(nodesKeyList, nodeKey)
			dbClusterNodesIds[nodeKey] = node
		}
		d.Set("nodes_key", nodesKeyList)
		d.Set("db_cluster_nodes_ids", dbClusterNodesIds)
	}
	d.Set("nodes", nodes)

	var documented map[string]interface{}
	configInter, ok := d.GetOk("endpoint_config")
	if !ok {
		documented = make(map[string]interface{}, 0)
		d.Set("endpoint_config", documented)
	}
	documented = configInter.(map[string]interface{})

	var endpointConfig = make(map[string]interface{})
	err = json.Unmarshal([]byte(object.EndpointConfig), &endpointConfig)
	if err != nil {
		return WrapError(err)
	}
	if endpointConfig["AutoAddNewNodes"] != nil {
		d.Set("auto_add_new_nodes", endpointConfig["AutoAddNewNodes"].(string))
	}
	if endpointConfig["AutoAddNewNodes"] != nil {
		d.Set("read_write_mode", endpointConfig["ReadWriteMode"].(string))
	}

	for k, v := range documented {
		if _, ok := endpointConfig[k]; ok {
			documented[k] = v
		}
	}
	if err := d.Set("endpoint_config", documented); err != nil {
		return WrapError(err)
	}

	for _, p := range object.AddressItems {
		if p.NetType == "Private" {
			d.Set("net_type", "Private")
			d.Set("port", p.Port)
			d.Set("connection_prefix", p.ConnectionString)
			d.Set("vpc_id", p.VPCId)
			d.Set("vswitch_id", p.VSwitchId)
		}
	}
	return nil
}

func resourceAlicloudPolarDBZonalEndpointUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	polarDBServiceV2 := PolarDbServiceV2{client}

	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}
	dbClusterId := parts[0]
	dbEndpointId := parts[1]

	if !d.IsNewResource() && d.HasChanges("read_write_mode", "auto_add_new_nodes", "endpoint_config", "db_endpoint_description", "nodes_key", "db_cluster_nodes_ids") {
		modifyEndpointRequest := polardb.CreateModifyDBClusterEndpointRequest()
		modifyEndpointRequest.RegionId = client.RegionId
		modifyEndpointRequest.DBClusterId = dbClusterId
		modifyEndpointRequest.DBEndpointId = dbEndpointId

		configItem := make(map[string]string)
		nodesKey, nodesKeyOk := d.GetOk("nodes_key")
		nodesIdsMap, nodesIdsMapOk := d.GetOk("db_cluster_nodes_ids")
		if nodesKeyOk && nodesIdsMapOk {
			var dbNodeParts []string
			nodesKey := expandStringList(nodesKey.(*schema.Set).List())
			nodeIdsMap := nodesIdsMap.(map[string]interface{})
			for _, nodeKey := range nodesKey {
				if v, exists := nodeIdsMap[nodeKey]; !exists {
					return WrapError(fmt.Errorf("node %s not found in db_cluster_nodes_ids", nodeKey))
				} else {
					dbNodeParts = append(dbNodeParts, v.(string))
				}
			}
			dbNodes := strings.Join(dbNodeParts, ",")
			modifyEndpointRequest.Nodes = dbNodes
			configItem["Nodes"] = dbNodes
		}

		modifyEndpointRequest.ReadWriteMode = d.Get("read_write_mode").(string)
		configItem["ReadWriteMode"] = d.Get("read_write_mode").(string)

		modifyEndpointRequest.AutoAddNewNodes = d.Get("auto_add_new_nodes").(string)
		configItem["AutoAddNewNodes"] = d.Get("auto_add_new_nodes").(string)

		endpointConfig, err := json.Marshal(d.Get("endpoint_config"))
		if err != nil {
			return WrapError(err)
		}
		modifyEndpointRequest.EndpointConfig = string(endpointConfig)
		configItem["EndpointConfig"] = string(endpointConfig)

		modifyEndpointRequest.DBEndpointDescription = d.Get("db_endpoint_description").(string)
		configItem["DBEndpointDescription"] = d.Get("db_endpoint_description").(string)

		if err := resource.Retry(8*time.Minute, func() *resource.RetryError {
			err = polarDBServiceV2.ModifyDBClusterEndpointZonal(modifyEndpointRequest)
			if err != nil {
				if NeedRetry(err) {
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		}); err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), modifyEndpointRequest.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		if _, err := polarDBServiceV2.WaitForPolarDBEndpoints(d, Active, nil, "", DefaultTimeoutMedium); err != nil {
			return err
		}

		stateConf := BuildStateConf([]string{}, []string{"Running"}, d.Timeout(schema.TimeoutUpdate), 10*time.Second, polarDBServiceV2.PolarDbZonalClusterStateRefreshFunc(dbClusterId, []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
	}

	return resourceAlicloudPolarDBZonalEndpointRead(d, meta)
}

func resourceAlicloudPolarDBZonalEndpointDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	polarDBServiceV2 := PolarDbServiceV2{client}

	parts, errParse := ParseResourceId(d.Id(), 2)
	DBClusterId := parts[0]
	DBEndpointId := parts[1]
	if errParse != nil {
		return WrapError(errParse)
	}
	object, err := polarDBServiceV2.DescribeDBClusterEndpointsZonal(d.Id())
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	if object.EndpointType != "Custom" {
		return WrapErrorf(Error("%s type endpoint can not be deleted.", object.EndpointType), DefaultErrorMsg, d.Id(), "DeleteDBClusterEndpoint", ProviderERROR)
	}

	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		err = polarDBServiceV2.DeleteDBClusterEndpointZonal(client.RegionId, DBClusterId, DBEndpointId)
		if err != nil {
			if NeedRetry(err) {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})

	if err != nil {
		if IsExpectedErrors(err, []string{"InvalidDBClusterId.NotFound", "InvalidCurrentConnectionString.NotFound", "AtLeastOneNetTypeExists"}) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), "DeleteDBClusterEndpointZonal", AlibabaCloudSdkGoERROR)
	}
	if _, err := polarDBServiceV2.WaitForPolarDBEndpoints(d, Deleted, nil, object.EndpointType, DefaultTimeoutMedium); err != nil {
		return WrapErrorf(err, DeleteTimeoutMsg, d.Id(), "DeleteDBClusterEndpointZonal", ProviderERROR)
	}

	stateConf := BuildStateConf([]string{}, []string{"Running"}, d.Timeout(schema.TimeoutUpdate), 10*time.Second, polarDBServiceV2.PolarDbZonalClusterStateRefreshFunc(DBClusterId, []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}
	return nil
}
