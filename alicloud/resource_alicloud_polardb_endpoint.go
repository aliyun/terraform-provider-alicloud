package alicloud

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/polardb"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceAlicloudPolarDBEndpoint() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudPolarDBEndpointCreate,
		Read:   resourceAlicloudPolarDBEndpointRead,
		Update: resourceAlicloudPolarDBEndpointUpdate,
		Delete: resourceAlicloudPolarDBEndpointDelete,
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
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"Custom"}, false),
				ForceNew:     true,
			},
			"nodes": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"read_write_mode": {
				Type:         schema.TypeString,
				ValidateFunc: validation.StringInSlice([]string{"ReadWrite", "ReadOnly"}, false),
				Optional:     true,
				Default:      "ReadOnly",
			},
			"auto_add_new_nodes": {
				Type:         schema.TypeString,
				ValidateFunc: validation.StringInSlice([]string{"Enable", "Disable"}, false),
				Optional:     true,
				Default:      "Disable",
			},
			"endpoint_config": {
				Type:     schema.TypeMap,
				Optional: true,
				Computed: true,
			},
		},
	}
}

func resourceAlicloudPolarDBEndpointCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	polarDBService := PolarDBService{client}
	clusterId := d.Get("db_cluster_id").(string)
	endpointType := d.Get("endpoint_type").(string)
	request := polardb.CreateCreateDBClusterEndpointRequest()
	request.RegionId = client.RegionId
	request.DBClusterId = clusterId
	request.EndpointType = endpointType
	if nodes, ok := d.GetOk("nodes"); ok {
		nodes := expandStringList(nodes.(*schema.Set).List())
		dbNodes := strings.Join(nodes, ",")
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

	enpoints, err := polarDBService.DescribePolarDBInstanceNetInfo(clusterId)
	if err != nil {
		return WrapError(err)
	}
	oldEndpoints := make([]interface{}, 0)
	for _, value := range enpoints {
		oldEndpoints = append(oldEndpoints, value.DBEndpointId)
	}
	oldEndpointIds := schema.NewSet(schema.HashString, oldEndpoints)

	var raw interface{}
	err = resource.Retry(8*time.Minute, func() *resource.RetryError {
		raw, err = client.WithPolarDBClient(func(polarDBClient *polardb.Client) (interface{}, error) {
			return polarDBClient.CreateDBClusterEndpoint(request)
		})
		if err != nil {
			OperationDeniedDBStatus = append(OperationDeniedDBStatus, "ClusterEndpoint.StatusNotValid")
			if IsExpectedErrors(err, OperationDeniedDBStatus) {
				return resource.RetryableError(err)
			}

			return resource.NonRetryableError(err)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		return nil
	})

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_polardb_endpoint", request.GetActionName(), AlibabaCloudSdkGoERROR)
	}

	dbEndpointId, err := polarDBService.WaitForPolarDBEndpoints(d, Active, oldEndpointIds, DefaultTimeoutMedium)
	if err != nil {
		return WrapError(err)
	}
	d.SetId(fmt.Sprintf("%s%s%s", clusterId, COLON_SEPARATED, dbEndpointId))

	return resourceAlicloudPolarDBEndpointRead(d, meta)
}

func resourceAlicloudPolarDBEndpointRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	polarDBService := PolarDBService{client}
	parts, errParse := ParseResourceId(d.Id(), 2)
	if errParse != nil {
		return WrapError(errParse)
	}
	dbClusterId := parts[0]
	object, err := polarDBService.DescribePolarDBClusterEndpoint(d.Id())

	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	d.Set("db_cluster_id", dbClusterId)
	d.Set("endpoint_type", object.EndpointType)
	nodes := strings.Split(object.Nodes, ",")
	d.Set("nodes", nodes)
	d.Set("auto_add_new_nodes", object.AutoAddNewNodes)
	d.Set("read_write_mode", object.ReadWriteMode)

	if err = polarDBService.RefreshEndpointConfig(d); err != nil {
		return WrapError(err)
	}

	return nil
}

func resourceAlicloudPolarDBEndpointUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	polarDBService := PolarDBService{client}

	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}
	request := polardb.CreateModifyDBClusterEndpointRequest()
	request.RegionId = client.RegionId
	request.DBClusterId = parts[0]
	request.DBEndpointId = parts[1]

	configItem := make(map[string]string)
	if d.HasChange("nodes") {
		nodes := expandStringList(d.Get("nodes").(*schema.Set).List())
		dbNodes := strings.Join(nodes, ",")
		request.Nodes = dbNodes
		configItem["Nodes"] = dbNodes
	}
	if d.HasChange("read_write_mode") {
		request.ReadWriteMode = d.Get("read_write_mode").(string)
		configItem["ReadWriteMode"] = d.Get("read_write_mode").(string)
	}
	if d.HasChange("auto_add_new_nodes") {
		request.AutoAddNewNodes = d.Get("auto_add_new_nodes").(string)
		configItem["AutoAddNewNodes"] = d.Get("auto_add_new_nodes").(string)
	}
	if d.HasChange("endpoint_config") {
		endpointConfig, err := json.Marshal(d.Get("endpoint_config"))
		if err != nil {
			return WrapError(err)
		}
		request.EndpointConfig = string(endpointConfig)
		configItem["EndpointConfig"] = string(endpointConfig)
	}

	if err := resource.Retry(8*time.Minute, func() *resource.RetryError {
		raw, err := client.WithPolarDBClient(func(polarDBClient *polardb.Client) (interface{}, error) {
			return polarDBClient.ModifyDBClusterEndpoint(request)
		})
		if err != nil {
			if IsExpectedErrors(err, []string{"EndpointStatus.NotSupport", "OperationDenied.DBClusterStatus"}) {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		return nil
	}); err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
	}

	// wait cluster endpoint config modified
	if err := polarDBService.WaitPolardbEndpointConfigEffect(
		d.Id(), configItem, DefaultTimeoutMedium); err != nil {
		return WrapError(err)
	}

	return resourceAlicloudPolarDBEndpointRead(d, meta)
}

func resourceAlicloudPolarDBEndpointDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	polarDBService := PolarDBService{client}

	parts, errParse := ParseResourceId(d.Id(), 2)
	if errParse != nil {
		return WrapError(errParse)
	}
	request := polardb.CreateDeleteDBClusterEndpointRequest()
	request.RegionId = client.RegionId
	request.DBClusterId = parts[0]
	request.DBEndpointId = parts[1]

	err := resource.Retry(5*time.Minute, func() *resource.RetryError {
		raw, err := client.WithPolarDBClient(func(polarDBClient *polardb.Client) (interface{}, error) {
			return polarDBClient.DeleteDBClusterEndpoint(request)
		})
		if err != nil {
			if IsExpectedErrors(err, []string{"OperationDenied.DBClusterStatus", "EndpointStatus.NotSupport", "ClusterEndpoint.StatusNotValid"}) {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		return nil
	})

	if err != nil {
		if IsExpectedErrors(err, []string{"InvalidDBClusterId.NotFound", "InvalidCurrentConnectionString.NotFound", "AtLeastOneNetTypeExists"}) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	endpointIds := schema.NewSet(schema.HashString, make([]interface{}, 0))
	dbEndpoint, err := polarDBService.WaitForPolarDBEndpoints(d, Deleted, endpointIds, DefaultTimeoutMedium)
	if dbEndpoint != "" || err != nil {
		return WrapErrorf(err, DeleteTimeoutMsg, d.Id(), request.GetActionName(), ProviderERROR)
	}
	return nil
}
