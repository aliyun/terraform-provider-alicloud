package alicloud

import (
	"encoding/json"
	"time"

	util "github.com/alibabacloud-go/tea-utils/service"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAlicloudRdsDBNodes() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudRdsDBNodesCreate,
		Read:   resourceAlicloudRdsDBNodesRead,
		Delete: resourceAlicloudRdsDBNodesDelete,
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
			"db_node": {
				Type:     schema.TypeSet,
				Required: true,
				ForceNew: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"class_code": {
							Type:     schema.TypeString,
							Required: true,
						},
						"zone_id": {
							Type:     schema.TypeString,
							Required: true,
						},
						"node_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"node_role": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"node_region_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func resourceAlicloudRdsDBNodesCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	rdsService := RdsService{client}
	var response map[string]interface{}
	action := "CreateDBNodes"
	request := map[string]interface{}{
		"RegionId":     client.RegionId,
		"DBInstanceId": Trim(d.Get("db_instance_id").(string)),
		"SourceIp":     client.SourceIp,
	}
	conn, err := client.NewRdsClient()
	if err != nil {
		return WrapError(err)
	}
	if v, ok := d.GetOk("db_node"); ok {
		objects := make([]interface{}, 0)
		for _, dbNodes := range v.(*schema.Set).List() {
			dbNodesArg := dbNodes.(map[string]interface{})
			dbNodesMap := map[string]interface{}{}
			dbNodesMap["classCode"] = dbNodesArg["class_code"]
			dbNodesMap["zoneId"] = dbNodesArg["zone_id"]

			objects = append(objects, dbNodesMap)
		}
		dbNode, err := json.Marshal(objects)
		if err != nil {
			return WrapError(err)
		}
		request["DBNode"] = string(dbNode)
	}
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutCreate)), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2014-08-15"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_rds_db_nodes", action, AlibabaCloudSdkGoERROR)
	}
	d.SetId(request["DBInstanceId"].(string))
	d.SetPartial("db_instance_id")
	stateConf := BuildStateConf([]string{}, []string{"Running"}, client.GetRetryTimeout(d.Timeout(schema.TimeoutCreate)), 3*time.Minute, rdsService.RdsDBInstanceStateRefreshFunc(d.Id(), []string{"Deleting"}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}
	return resourceAlicloudRdsDBNodesRead(d, meta)
}

func resourceAlicloudRdsDBNodesRead(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	rdsService := RdsService{client}
	instance, proxyErr := rdsService.DescribeDBInstance(d.Id())
	if proxyErr != nil {
		if NotFoundError(proxyErr) {
			d.SetId("")
			return nil
		}
		return WrapError(proxyErr)
	}
	d.Set("db_instance_id", d.Id())
	if dbNodesList, ok := instance["DBClusterNodes"]; ok && dbNodesList != nil {
		dbNodesMaps := make([]map[string]interface{}, 0)
		if v, ok := dbNodesList.(map[string]interface{})["DBClusterNode"].([]interface{}); ok && len(v) > 2 {
			dbNodeLen := len(v)
			for i := 2; i < dbNodeLen; i++ {
				dbNodesListItemMap := map[string]interface{}{}
				dbNodesListItemMap["class_code"] = v[i].(map[string]interface{})["ClassCode"]
				dbNodesListItemMap["zone_id"] = v[i].(map[string]interface{})["NodeZoneId"]
				dbNodesListItemMap["node_id"] = v[i].(map[string]interface{})["NodeId"]
				dbNodesListItemMap["node_role"] = v[i].(map[string]interface{})["NodeRole"]
				dbNodesListItemMap["node_region_id"] = v[i].(map[string]interface{})["NodeRegionId"]
				dbNodesMaps = append(dbNodesMaps, dbNodesListItemMap)
			}
		}
		d.Set("db_node", dbNodesMaps)
	}
	return nil
}

func resourceAlicloudRdsDBNodesDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	rdsService := RdsService{client}
	_, err := rdsService.DescribeDBInstance(d.Id())
	if err != nil {
		if NotFoundError(err) {
			return nil
		}
		return WrapError(err)
	}
	action := "DeleteDBNodes"
	request := map[string]interface{}{
		"RegionId":     client.RegionId,
		"DBInstanceId": d.Id(),
		"SourceIp":     client.SourceIp,
	}
	if v, ok := d.GetOk("db_node"); ok && len(v.(*schema.Set).List()) > 0 {
		nodeIds := make([]string, 0, len(v.(*schema.Set).List()))
		for _, node := range v.(*schema.Set).List() {
			if node == nil {
				continue
			}
			nodeMap := node.(map[string]interface{})
			nodeIds = append(nodeIds, nodeMap["node_id"].(string))
		}
		dbNodeIds, err := json.Marshal(nodeIds)
		if err != nil {
			return WrapError(err)
		}
		request["DBNodeId"] = string(dbNodeIds)
	}
	var response map[string]interface{}
	conn, err := client.NewRdsClient()
	if err != nil {
		return WrapError(err)
	}
	runtime := util.RuntimeOptions{}
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutDelete)), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2014-08-15"), StringPointer("AK"), nil, request, &runtime)
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
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}
	stateConf := BuildStateConf([]string{}, []string{"Running"}, client.GetRetryTimeout(d.Timeout(schema.TimeoutDelete)), 3*time.Minute, rdsService.RdsDBInstanceStateRefreshFunc(d.Id(), []string{"Deleting"}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}
	return nil
}
