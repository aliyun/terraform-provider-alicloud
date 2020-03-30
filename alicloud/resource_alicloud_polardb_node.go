package alicloud

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/polardb"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func resourceAlicloudPolarDBNode() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudPolarDBNodeCreate,
		Read:   resourceAlicloudPolarDBNodeRead,
		Update: resourceAlicloudPolarDBNodeUpdate,
		Delete: resourceAlicloudPolarDBNodeDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Update: schema.DefaultTimeout(20 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"db_cluster_id": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},
			"db_node_class": {
				Type:     schema.TypeString,
				Required: true,
			},
			"modify_type": {
				Type:         schema.TypeString,
				ValidateFunc: validation.StringInSlice([]string{"Upgrade", "Downgrade"}, false),
				Optional:     true,
			},
		},
	}
}

func resourceAlicloudPolarDBNodeCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	polarDBService := PolarDBService{client}

	request := polardb.CreateCreateDBNodesRequest()
	request.RegionId = client.RegionId
	request.DBClusterId = d.Get("db_cluster_id").(string)
	request.DBNode = &[]polardb.CreateDBNodesDBNode{
		{
			TargetClass: d.Get("db_node_class").(string),
		},
	}
	raw, err := client.WithPolarDBClient(func(polarDBClient *polardb.Client) (interface{}, error) {
		return polarDBClient.CreateDBNodes(request)
	})

	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	if err != nil {
		return WrapErrorf(
			err, DefaultErrorMsg, "alicloud_polardb_node", request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	response, _ := raw.(*polardb.CreateDBNodesResponse)
	d.SetId(fmt.Sprintf("%s%s%s", request.DBClusterId, COLON_SEPARATED, response.DBNodeIds.DBNodeId[0]))

	// wait cluster status change from DBNodeCreating to running
	stateConf := BuildStateConf(
		[]string{"DBNodeCreating"},
		[]string{"Running"},
		d.Timeout(schema.TimeoutCreate),
		1*time.Minute,
		polarDBService.PolarDBClusterStateRefreshFunc(response.DBClusterId, []string{"Deleting"}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, response.DBClusterId)
	}

	return resourceAlicloudPolarDBNodeRead(d, meta)
}

func resourceAlicloudPolarDBNodeRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	polarDBService := PolarDBService{client}
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}
	dbClusterId := parts[0]
	dbNodeId := parts[1]
	object, err := polarDBService.DescribePolarDBCluster(dbClusterId)
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	dbNodes := object.DBNodes.DBNode

	for _, value := range dbNodes {
		if value.DBNodeId == dbNodeId {
			d.Set("db_cluster_id", dbClusterId)
			d.Set("db_node_class", value.DBNodeClass)
			break
		}
	}

	return nil
}

func resourceAlicloudPolarDBNodeUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	polarDBService := PolarDBService{client}

	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}
	dbClusterId := parts[0]
	// make sure DB cluster status is running before updating node
	if err := polarDBService.WaitForPolarDBInstance(dbClusterId, Running, 1800); err != nil {
		return WrapError(err)
	}
	d.Partial(true)

	if d.HasChange("db_node_class") && d.HasChange("modify_type") {
		request := polardb.CreateModifyDBNodeClassRequest()
		request.RegionId = client.RegionId
		request.DBClusterId = dbClusterId
		request.DBNodeTargetClass = d.Get("db_node_class").(string)
		request.ModifyType = d.Get("modify_type").(string)

		raw, err := client.WithPolarDBClient(func(polarDBClient *polardb.Client) (interface{}, error) {
			return polarDBClient.ModifyDBNodeClass(request)
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, dbClusterId, request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		stateConf := BuildStateConf(
			[]string{"ClassChanging"},
			[]string{"Running"},
			d.Timeout(schema.TimeoutCreate),
			1*time.Minute,
			polarDBService.PolarDBClusterStateRefreshFunc(dbClusterId, []string{"Deleting"}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, dbClusterId)
		}
		d.SetPartial("modify_type")
		d.SetPartial("db_node_class")
	}

	d.Partial(false)
	return resourceAlicloudPolarDBNodeRead(d, meta)
}

func resourceAlicloudPolarDBNodeDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	polarDBService := PolarDBService{client}

	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}
	dbClusterId := parts[0]
	DBNodeId := parts[1]
	// make sure DB cluster status is running before deleting node
	if err := polarDBService.WaitForPolarDBInstance(dbClusterId, Running, 1800); err != nil {
		return WrapError(err)
	}

	request := polardb.CreateDeleteDBNodesRequest()
	request.RegionId = client.RegionId
	request.DBClusterId = dbClusterId
	request.DBNodeId = &[]string{
		DBNodeId,
	}

	raw, err := client.WithPolarDBClient(func(polarDBClient *polardb.Client) (interface{}, error) {
		return polarDBClient.DeleteDBNodes(request)
	})

	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	stateConf := BuildStateConf(
		[]string{"DBNodeDeleting"},
		[]string{"Running"},
		d.Timeout(schema.TimeoutDelete),
		1*time.Minute,
		polarDBService.PolarDBClusterStateRefreshFunc(dbClusterId, []string{"Deleting"}))
	if _, err = stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, dbClusterId)
	}

	return nil
}
