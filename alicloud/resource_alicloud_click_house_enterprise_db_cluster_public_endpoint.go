// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/blues/jsonata-go"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAliCloudClickHouseEnterpriseDbClusterPublicEndpoint() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudClickHouseEnterpriseDbClusterPublicEndpointCreate,
		Read:   resourceAliCloudClickHouseEnterpriseDbClusterPublicEndpointRead,
		Update: resourceAliCloudClickHouseEnterpriseDbClusterPublicEndpointUpdate,
		Delete: resourceAliCloudClickHouseEnterpriseDbClusterPublicEndpointDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"connection_string_prefix": {
				Type:     schema.TypeString,
				Required: true,
			},
			"db_instance_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"net_type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: StringInSlice([]string{"Public"}, false),
			},
		},
	}
}

func resourceAliCloudClickHouseEnterpriseDbClusterPublicEndpointCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "CreateEndpoint"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	if v, ok := d.GetOk("net_type"); ok {
		request["DBInstanceNetType"] = v
	}
	if v, ok := d.GetOk("db_instance_id"); ok {
		request["DBInstanceId"] = v
	}
	request["RegionId"] = client.RegionId

	request["ConnectionPrefix"] = d.Get("connection_string_prefix")
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.RpcPost("clickhouse", "2023-05-22", action, query, request, true)
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_click_house_enterprise_db_cluster_public_endpoint", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprintf("%v:%v", request["DBInstanceId"], request["DBInstanceNetType"]))

	clickHouseServiceV2 := ClickHouseServiceV2{client}
	stateConf := BuildStateConf([]string{}, []string{"ACTIVATION", "ACTIVE"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, clickHouseServiceV2.DescribeAsyncClickHouseEnterpriseDbClusterPublicEndpointStateRefreshFunc(d, response, "$.Data.Status", []string{}))
	if jobDetail, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id(), jobDetail)
	}

	return resourceAliCloudClickHouseEnterpriseDbClusterPublicEndpointUpdate(d, meta)
}

func resourceAliCloudClickHouseEnterpriseDbClusterPublicEndpointRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	clickHouseServiceV2 := ClickHouseServiceV2{client}

	objectRaw, err := clickHouseServiceV2.DescribeClickHouseEnterpriseDbClusterPublicEndpoint(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_click_house_enterprise_db_cluster_public_endpoint DescribeClickHouseEnterpriseDbClusterPublicEndpoint Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("net_type", convertClickHouseEnterpriseDbClusterPublicEndpointDataEndpointsNetTypeResponse(objectRaw["NetType"]))

	e := jsonata.MustCompile("$substringBefore($.ConnectionString, '.')")
	evaluation, _ := e.Eval(objectRaw)
	d.Set("connection_string_prefix", evaluation)

	parts := strings.Split(d.Id(), ":")
	d.Set("db_instance_id", parts[0])

	return nil
}

func resourceAliCloudClickHouseEnterpriseDbClusterPublicEndpointUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	update := false

	var err error
	parts := strings.Split(d.Id(), ":")
	action := "ModifyDBInstanceConnectionString"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["DBInstanceNetType"] = parts[1]
	request["DBInstanceId"] = parts[0]
	request["RegionId"] = client.RegionId
	if !d.IsNewResource() && d.HasChange("connection_string_prefix") {
		update = true
	}
	request["ConnectionStringPrefix"] = d.Get("connection_string_prefix")
	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("clickhouse", "2023-05-22", action, query, request, true)
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
		clickHouseServiceV2 := ClickHouseServiceV2{client}
		stateConf := BuildStateConf([]string{}, []string{"ACTIVATION", "ACTIVE"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, clickHouseServiceV2.DescribeAsyncClickHouseEnterpriseDbClusterPublicEndpointStateRefreshFunc(d, response, "$.Data.Status", []string{}))
		if jobDetail, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id(), jobDetail)
		}
	}

	return resourceAliCloudClickHouseEnterpriseDbClusterPublicEndpointRead(d, meta)
}

func resourceAliCloudClickHouseEnterpriseDbClusterPublicEndpointDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	parts := strings.Split(d.Id(), ":")
	action := "DeleteEndpoint"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	request["DBInstanceNetType"] = parts[1]
	request["DBInstanceId"] = parts[0]
	request["RegionId"] = client.RegionId

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = client.RpcPost("clickhouse", "2023-05-22", action, query, request, true)

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
		if NotFoundError(err) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	clickHouseServiceV2 := ClickHouseServiceV2{client}
	stateConf := BuildStateConf([]string{}, []string{"ACTIVATION", "ACTIVE"}, d.Timeout(schema.TimeoutDelete), 5*time.Second, clickHouseServiceV2.DescribeAsyncClickHouseEnterpriseDbClusterPublicEndpointStateRefreshFunc(d, response, "$.Data.Status", []string{}))
	if jobDetail, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id(), jobDetail)
	}

	return nil
}

func convertClickHouseEnterpriseDbClusterPublicEndpointDataEndpointsNetTypeResponse(source interface{}) interface{} {
	source = fmt.Sprint(source)
	switch source {
	case "PUBLIC":
		return "Public"
	}
	return source
}
