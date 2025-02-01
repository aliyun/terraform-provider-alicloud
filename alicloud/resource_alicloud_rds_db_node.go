package alicloud

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAlicloudRdsDBNode() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudRdsDBNodeCreate,
		Read:   resourceAlicloudRdsDBNodeRead,
		Delete: resourceAlicloudRdsDBNodeDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(20 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"db_instance_id": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},
			"class_code": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"zone_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"node_role": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"node_region_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"node_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceAlicloudRdsDBNodeCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	rdsService := RdsService{client}
	var response map[string]interface{}
	action := "CreateDBNodes"
	request := map[string]interface{}{
		"RegionId":     client.RegionId,
		"DBInstanceId": Trim(d.Get("db_instance_id").(string)),
		"SourceIp":     client.SourceIp,
	}
	var err error
	objects := make([]interface{}, 0)
	dbNodesMap := make(map[string]interface{})
	dbNodesMap["classCode"] = d.Get("class_code").(string)
	dbNodesMap["zoneId"] = d.Get("zone_id").(string)
	objects = append(objects, dbNodesMap)
	dbNode, err := json.Marshal(objects)
	if err != nil {
		return WrapError(err)
	}
	request["DBNode"] = string(dbNode)
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutCreate)), func() *resource.RetryError {
		response, err = client.RpcPost("Rds", "2014-08-15", action, nil, request, false)
		if err != nil {
			if IsExpectedErrors(err, []string{"SYSTEM.CONCURRENT_OPERATE", "OperationDenied.DBInstanceStatus", "MissingParameter"}) || NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_rds_db_node", action, AlibabaCloudSdkGoERROR)
	}
	d.SetId(fmt.Sprint(request["DBInstanceId"], ":", response["NodeIds"].(string)))
	stateConf := BuildStateConf([]string{}, []string{"Running"}, client.GetRetryTimeout(d.Timeout(schema.TimeoutCreate)), 3*time.Minute, rdsService.RdsDBInstanceStateRefreshFunc(request["DBInstanceId"].(string), []string{"Deleting"}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}
	return resourceAlicloudRdsDBNodeRead(d, meta)
}

func resourceAlicloudRdsDBNodeRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	rdsService := RdsService{client}
	dbNode, nodeErr := rdsService.DescribeRdsNode(d.Id())
	if nodeErr != nil {
		if !d.IsNewResource() && NotFoundError(nodeErr) {
			d.SetId("")
			return nil
		}
		return WrapError(nodeErr)
	}
	d.Set("class_code", dbNode["ClassCode"])
	d.Set("zone_id", dbNode["NodeZoneId"])
	d.Set("node_role", dbNode["NodeRole"])
	d.Set("node_region_id", dbNode["NodeRegionId"])
	d.Set("node_id", dbNode["NodeId"])
	d.Set("db_instance_id", dbNode["DBInstanceId"])
	return nil
}

func resourceAlicloudRdsDBNodeDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	rdsService := RdsService{client}
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}
	_, err = rdsService.DescribeDBInstance(parts[0])
	if err != nil {
		if NotFoundError(err) {
			return nil
		}
		return WrapError(err)
	}
	stateConf := BuildStateConf([]string{}, []string{"Running"}, client.GetRetryTimeout(d.Timeout(schema.TimeoutDelete)), 3*time.Minute, rdsService.RdsDBInstanceStateRefreshFunc(parts[0], []string{"Deleting"}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}
	action := "DeleteDBNodes"
	request := map[string]interface{}{
		"RegionId":     client.RegionId,
		"DBInstanceId": parts[0],
		"SourceIp":     client.SourceIp,
	}
	nodeIds := make([]string, 0)
	nodeIds = append(nodeIds, parts[1])
	dbNodeIds, err := json.Marshal(nodeIds)
	if err != nil {
		return WrapError(err)
	}
	request["DBNodeId"] = string(dbNodeIds)
	var response map[string]interface{}
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutDelete)), func() *resource.RetryError {
		response, err = client.RpcPost("Rds", "2014-08-15", action, nil, request, false)
		if err != nil {
			if IsExpectedErrors(err, []string{"SYSTEM.CONCURRENT_OPERATE", "OperationDenied.DBInstanceStatus", "MissingParameter"}) || NeedRetry(err) {
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
	stateConf = BuildStateConf([]string{}, []string{"Running"}, client.GetRetryTimeout(d.Timeout(schema.TimeoutDelete)), 3*time.Minute, rdsService.RdsDBInstanceStateRefreshFunc(parts[0], []string{"Deleting"}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}
	return nil
}
