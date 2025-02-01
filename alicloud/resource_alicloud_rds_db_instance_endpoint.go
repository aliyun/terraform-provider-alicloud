package alicloud

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/PaesslerAG/jsonpath"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAlicloudRdsDBInstanceEndpoint() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudRdsDBInstanceEndpointCreate,
		Read:   resourceAlicloudRdsDBInstanceEndpointRead,
		Update: resourceAlicloudRdsDBInstanceEndpointUpdate,
		Delete: resourceAlicloudRdsDBInstanceEndpointDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(20 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"db_instance_id": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},
			"vpc_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"vswitch_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"private_ip_address": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"db_instance_endpoint_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"connection_string_prefix": {
				Type:     schema.TypeString,
				Required: true,
			},
			"connection_string": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"port": {
				Type:     schema.TypeString,
				Required: true,
			},
			"db_instance_endpoint_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"db_instance_endpoint_description": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"ip_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"node_items": {
				Type:     schema.TypeSet,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"node_id": {
							Type:     schema.TypeString,
							Required: true,
						},
						"weight": {
							Type:     schema.TypeInt,
							Required: true,
						},
					},
				},
			},
		},
	}
}

func resourceAlicloudRdsDBInstanceEndpointCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	rdsService := RdsService{client}
	var response map[string]interface{}
	action := "CreateDBInstanceEndpoint"
	request := map[string]interface{}{
		"RegionId":               client.RegionId,
		"DBInstanceId":           Trim(d.Get("db_instance_id").(string)),
		"VPCId":                  d.Get("vpc_id"),
		"VSwitchId":              d.Get("vswitch_id"),
		"ConnectionStringPrefix": d.Get("connection_string_prefix"),
		"Port":                   d.Get("port"),
		"DBInstanceEndpointType": "Readonly",
		"SourceIp":               client.SourceIp,
	}

	var err error
	if v, ok := d.GetOk("db_instance_endpoint_description"); ok {
		request["DBInstanceEndpointDescription"] = v
	}
	if v, ok := d.GetOk("node_items"); ok {
		objects := make([]interface{}, 0)
		for _, dbNodes := range v.(*schema.Set).List() {
			dbNodesArg := dbNodes.(map[string]interface{})
			dbNodesMap := map[string]interface{}{}
			dbNodesMap["NodeId"] = dbNodesArg["node_id"]
			dbNodesMap["Weight"] = dbNodesArg["weight"]
			dbNodesMap["DBInstanceId"] = d.Get("db_instance_id").(string)

			objects = append(objects, dbNodesMap)
		}
		dbNode, err := json.Marshal(objects)
		if err != nil {
			return WrapError(err)
		}
		request["NodeItems"] = string(dbNode)
	}
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutCreate)), func() *resource.RetryError {
		response, err = client.RpcPost("Rds", "2014-08-15", action, nil, request, false)
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_rds_db_instance_endpoint", action, AlibabaCloudSdkGoERROR)
	}
	if v, err := jsonpath.Get("$.Data", response); err != nil || v == nil {
		return WrapErrorf(err, IdMsg, "alicloud_cloud_rds_db_instance_endpoint")
	} else {
		object := v.(map[string]interface{})
		d.SetId(fmt.Sprint(request["DBInstanceId"], ":", object["DBInstanceEndpointId"].(string)))
		d.Set("connection_string", object["ConnectionString"].(string))
	}
	stateConf := BuildStateConf([]string{}, []string{"Running"}, client.GetRetryTimeout(d.Timeout(schema.TimeoutCreate)), 1*time.Minute, rdsService.RdsDBInstanceStateRefreshFunc(request["DBInstanceId"].(string), []string{"Deleting"}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}
	return resourceAlicloudRdsDBInstanceEndpointUpdate(d, meta)
}

func resourceAlicloudRdsDBInstanceEndpointUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var err error
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}
	rdsService := RdsService{client}
	if !d.IsNewResource() && d.HasChanges("db_instance_endpoint_description", "node_items") {
		action := "ModifyDBInstanceEndpoint"
		request := map[string]interface{}{
			"RegionId":             client.RegionId,
			"DBInstanceEndpointId": parts[1],
			"DBInstanceId":         parts[0],
			"SourceIp":             client.SourceIp,
		}
		if v, ok := d.GetOk("db_instance_endpoint_description"); ok {
			request["DBInstanceEndpointDescription"] = v
		}
		if v, ok := d.GetOk("node_items"); ok {
			objects := make([]interface{}, 0)
			for _, dbNodes := range v.(*schema.Set).List() {
				dbNodesArg := dbNodes.(map[string]interface{})
				dbNodesMap := map[string]interface{}{}
				dbNodesMap["NodeId"] = dbNodesArg["node_id"]
				dbNodesMap["Weight"] = dbNodesArg["weight"]
				dbNodesMap["DBInstanceId"] = parts[0]
				objects = append(objects, dbNodesMap)
			}
			dbNode, err := json.Marshal(objects)
			if err != nil {
				return WrapError(err)
			}
			request["NodeItems"] = string(dbNode)
		}
		var response map[string]interface{}
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *resource.RetryError {
			response, err = client.RpcPost("Rds", "2014-08-15", action, nil, request, false)
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
			return WrapErrorf(err, DefaultErrorMsg, "alicloud_rds_db_instance_endpoint", action, AlibabaCloudSdkGoERROR)
		}
		stateConf := BuildStateConf([]string{}, []string{"Running"}, client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), 3*time.Minute, rdsService.RdsDBInstanceStateRefreshFunc(parts[0], []string{"Deleting"}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
	}
	if !d.IsNewResource() && d.HasChanges("vpc_id", "vswitch_id") {
		action := "ModifyDBInstanceEndpointAddress"
		request := map[string]interface{}{
			"RegionId":             client.RegionId,
			"DBInstanceEndpointId": parts[1],
			"DBInstanceId":         parts[0],
			"ConnectionString":     d.Get("connection_string"),
			"SourceIp":             client.SourceIp,
		}
		if v, ok := d.GetOk("vpc_id"); ok {
			request["VpcId"] = v
		}
		if v, ok := d.GetOk("vswitch_id"); ok {
			request["VSwitchId"] = v
		}
		var response map[string]interface{}
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *resource.RetryError {
			response, err = client.RpcPost("Rds", "2014-08-15", action, nil, request, false)
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
			return WrapErrorf(err, DefaultErrorMsg, "alicloud_rds_db_instance_endpoint", action, AlibabaCloudSdkGoERROR)
		}
		stateConf := BuildStateConf([]string{}, []string{"Running"}, client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), 3*time.Minute, rdsService.RdsDBInstanceStateRefreshFunc(parts[0], []string{"Deleting"}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
	}
	if !d.IsNewResource() && d.HasChanges("connection_string_prefix", "port") {
		action := "ModifyDBInstanceEndpointAddress"
		request := map[string]interface{}{
			"RegionId":             client.RegionId,
			"DBInstanceEndpointId": parts[1],
			"DBInstanceId":         parts[0],
			"ConnectionString":     d.Get("connection_string"),
			"SourceIp":             client.SourceIp,
		}
		if v, ok := d.GetOk("connection_string_prefix"); ok {
			request["ConnectionStringPrefix"] = v
		}
		if v, ok := d.GetOk("port"); ok {
			request["Port"] = v
		}
		var response map[string]interface{}
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *resource.RetryError {
			response, err = client.RpcPost("Rds", "2014-08-15", action, nil, request, false)
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
			return WrapErrorf(err, DefaultErrorMsg, "alicloud_rds_db_instance_endpoint", action, AlibabaCloudSdkGoERROR)
		}
		stateConf := BuildStateConf([]string{}, []string{"Running"}, client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), 3*time.Minute, rdsService.RdsDBInstanceStateRefreshFunc(parts[0], []string{"Deleting"}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
	}
	return resourceAlicloudRdsDBInstanceEndpointRead(d, meta)
}

func resourceAlicloudRdsDBInstanceEndpointRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	rdsService := RdsService{client}
	object, endpointErr := rdsService.DescribeDBInstanceEndpoints(d.Id())
	if endpointErr != nil {
		if !d.IsNewResource() && NotFoundError(endpointErr) {
			log.Printf("[DEBUG] Resource alicloud_rds_db_instance_endpoint rdsService.DescribeDBInstanceEndpoints Failed!!! %s", endpointErr)
			d.SetId("")
			return nil
		}
		return WrapError(endpointErr)
	}
	d.Set("db_instance_endpoint_description", object["EndpointDescription"])
	d.Set("db_instance_endpoint_type", object["EndpointType"])
	d.Set("node_items", object["NodeItems"])
	d.Set("vpc_id", object["VpcId"])
	d.Set("ip_type", object["IpType"])
	d.Set("private_ip_address", object["IpAddress"])
	d.Set("vswitch_id", object["VSwitchId"])
	d.Set("port", object["Port"])
	d.Set("connection_string", object["ConnectionString"])
	d.Set("connection_string_prefix", object["ConnectionStringPrefix"])
	d.Set("db_instance_id", object["DBInstanceId"])
	d.Set("db_instance_endpoint_id", object["DBInstanceEndpointId"])
	return nil
}

func resourceAlicloudRdsDBInstanceEndpointDelete(d *schema.ResourceData, meta interface{}) error {
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
	action := "DeleteDBInstanceEndpoint"
	request := map[string]interface{}{
		"RegionId":             client.RegionId,
		"DBInstanceId":         parts[0],
		"DBInstanceEndpointId": parts[1],
		"SourceIp":             client.SourceIp,
	}
	var response map[string]interface{}
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutDelete)), func() *resource.RetryError {
		response, err = client.RpcPost("Rds", "2014-08-15", action, nil, request, false)
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
	stateConf := BuildStateConf([]string{}, []string{"Running"}, client.GetRetryTimeout(d.Timeout(schema.TimeoutDelete)), 1*time.Minute, rdsService.RdsDBInstanceStateRefreshFunc(parts[0], []string{"Deleting"}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}
	return nil
}
