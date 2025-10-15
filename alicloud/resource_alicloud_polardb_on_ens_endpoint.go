package alicloud

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/polardb"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAlicloudPolarDBOnENSEndpoint() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudPolarDBOnENSEndpointCreate,
		Read:   resourceAlicloudPolarDBOnENSEndpointRead,
		Update: resourceAlicloudPolarDBOnENSEndpointUpdate,
		Delete: resourceAlicloudPolarDBOnENSEndpointDelete,
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
				Default:      "Custom",
			},
			"nodes": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"nodes_key": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"db_cluster_nodes_ids": {
				Type:     schema.TypeMap,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"read_write_mode": {
				Type:         schema.TypeString,
				ValidateFunc: StringInSlice([]string{"ReadWrite", "ReadOnly"}, false),
				Optional:     true,
				Computed:     true,
			},
			"auto_add_new_nodes": {
				Type:         schema.TypeString,
				ValidateFunc: StringInSlice([]string{"Enable", "Disable"}, false),
				Default:      "Enable",
				Optional:     true,
			},
			"endpoint_config": {
				Type:     schema.TypeMap,
				Optional: true,
				Computed: true,
			},
			"net_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: StringInSlice([]string{"Private"}, false),
				Default:      "Private",
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
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: StringMatch(regexp.MustCompile(`^[a-z][a-z0-9\\-]{4,28}[a-z0-9]$`), "The prefix must be 6 to 30 characters in length, and can contain lowercase letters, digits, and hyphens (-),  must start with a letter and end with a digit or letter."),
			},
			"port": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"vpc_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"vswitch_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		},
	}
}

func resourceAlicloudPolarDBOnENSEndpointCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	polarDbServiceV2 := PolarDbServiceV2{client}
	clusterId := d.Get("db_cluster_id").(string)
	endpointType := d.Get("endpoint_type").(string)
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
	if autoAddNewNodes, ok := d.GetOk("auto_add_new_nodes"); ok {
		request.AutoAddNewNodes = autoAddNewNodes.(string)
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
	dbEndpointId, err := polarDbServiceV2.WaitForPolarDBEndpoints(d, Active, oldEndpointIds, DefaultTimeoutMedium)
	if err != nil {
		return WrapError(err)
	}
	d.SetId(fmt.Sprintf("%s%s%s", clusterId, COLON_SEPARATED, dbEndpointId))

	stateConf := BuildStateConf([]string{}, []string{"Running"}, d.Timeout(schema.TimeoutUpdate), 10*time.Second, polarDbServiceV2.PolarDbZonalClusterStateRefreshFunc(clusterId, []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return resourceAlicloudPolarDBOnENSEndpointUpdate(d, meta)
}

func resourceAlicloudPolarDBOnENSEndpointRead(d *schema.ResourceData, meta interface{}) error {
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
	d.Set("auto_add_new_nodes", endpointConfig["AutoAddNewNodes"].(string))
	d.Set("read_write_mode", endpointConfig["ReadWriteMode"].(string))

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
			d.Set("port", p.Port)
			prefix := strings.Split(p.ConnectionString, ".")
			d.Set("connection_prefix", prefix[0])
			d.Set("vpc_id", p.VPCId)
			d.Set("vswitch_id", p.VSwitchId)
		}
	}
	return nil
}

func resourceAlicloudPolarDBOnENSEndpointUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	polarDBServiceV2 := PolarDbServiceV2{client}

	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}
	dbClusterId := parts[0]
	dbEndpointId := parts[1]

	if !d.IsNewResource() && d.HasChanges("nodes", "read_write_mode", "auto_add_new_nodes", "endpoint_config", "db_endpoint_description", "nodes_key", "db_cluster_nodes_ids") {
		modifyEndpointRequest := polardb.CreateModifyDBClusterEndpointRequest()
		modifyEndpointRequest.RegionId = client.RegionId
		modifyEndpointRequest.DBClusterId = dbClusterId
		modifyEndpointRequest.DBEndpointId = dbEndpointId

		configItem := make(map[string]string)
		nodes := expandStringList(d.Get("nodes").(*schema.Set).List())
		dbNodes := strings.Join(nodes, ",")
		modifyEndpointRequest.Nodes = dbNodes
		configItem["Nodes"] = dbNodes
		nodesKey, nodesKeyOk := d.GetOk("nodes_key")
		nodesIdsMap, nodesIdsMapOk := d.GetOk("db_cluster_nodes_ids")
		if nodesKeyOk && nodesIdsMapOk && (d.HasChange("nodes_key") || d.HasChange("db_cluster_nodes_ids")) {
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
			addDebug("ModifyDBClusterEndpointZonal", modifyEndpointRequest.RpcRequest, modifyEndpointRequest)
			return nil
		}); err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), modifyEndpointRequest.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		if _, err := polarDBServiceV2.WaitForPolarDBEndpoints(d, Active, nil, DefaultTimeoutMedium); err != nil {
			return err
		}

		stateConf := BuildStateConf([]string{}, []string{"Running"}, d.Timeout(schema.TimeoutUpdate), 10*time.Second, polarDBServiceV2.PolarDbZonalClusterStateRefreshFunc(dbClusterId, []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
	}

	return resourceAlicloudPolarDBOnENSEndpointRead(d, meta)
}

func resourceAlicloudPolarDBOnENSEndpointDelete(d *schema.ResourceData, meta interface{}) error {
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
	if _, err := polarDBServiceV2.WaitForPolarDBEndpoints(d, Deleted, nil, DefaultTimeoutMedium); err != nil {
		return WrapErrorf(err, DeleteTimeoutMsg, d.Id(), "DeleteDBClusterEndpointZonal", ProviderERROR)
	}

	stateConf := BuildStateConf([]string{}, []string{"Running"}, d.Timeout(schema.TimeoutUpdate), 10*time.Second, polarDBServiceV2.PolarDbZonalClusterStateRefreshFunc(DBClusterId, []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}
	return nil
}
