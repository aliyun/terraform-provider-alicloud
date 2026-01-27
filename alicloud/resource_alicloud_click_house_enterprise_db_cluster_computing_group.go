// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAliCloudClickHouseEnterpriseDbClusterComputingGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudClickHouseEnterpriseDbClusterComputingGroupCreate,
		Read:   resourceAliCloudClickHouseEnterpriseDbClusterComputingGroupRead,
		Update: resourceAliCloudClickHouseEnterpriseDbClusterComputingGroupUpdate,
		Delete: resourceAliCloudClickHouseEnterpriseDbClusterComputingGroupDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(60 * time.Minute),
			Update: schema.DefaultTimeout(60 * time.Minute),
			Delete: schema.DefaultTimeout(60 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"computing_group_description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"computing_group_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"db_instance_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"is_readonly": {
				Type:     schema.TypeBool,
				Required: true,
			},
			"node_count": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"node_scale_max": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"node_scale_min": {
				Type:     schema.TypeInt,
				Required: true,
			},
		},
	}
}

func resourceAliCloudClickHouseEnterpriseDbClusterComputingGroupCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "CreateComputingGroup"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	if v, ok := d.GetOk("db_instance_id"); ok {
		request["DBInstanceId"] = v
	}
	request["RegionId"] = client.RegionId

	request["NodeScaleMax"] = d.Get("node_scale_max")
	request["NodeScaleMin"] = d.Get("node_scale_min")
	request["IsReadonly"] = d.Get("is_readonly")
	request["NodeCount"] = d.Get("node_count")
	if v, ok := d.GetOk("computing_group_description"); ok {
		request["ComputingGroupDescription"] = v
	}
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_click_house_enterprise_db_cluster_computing_group", action, AlibabaCloudSdkGoERROR)
	}

	DataDBInstanceNameVar, _ := jsonpath.Get("$.Data.DBInstanceName", response)
	DataComputingGroupIdVar, _ := jsonpath.Get("$.Data.ComputingGroupId", response)
	d.SetId(fmt.Sprintf("%v:%v", DataDBInstanceNameVar, DataComputingGroupIdVar))

	clickHouseServiceV2 := ClickHouseServiceV2{client}
	stateConf := BuildStateConf([]string{}, []string{"activation"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, clickHouseServiceV2.ClickHouseEnterpriseDbClusterComputingGroupStateRefreshFunc(d.Id(), "$.ComputingGroupStatus", []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return resourceAliCloudClickHouseEnterpriseDbClusterComputingGroupUpdate(d, meta)
}

func resourceAliCloudClickHouseEnterpriseDbClusterComputingGroupRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	clickHouseServiceV2 := ClickHouseServiceV2{client}

	objectRaw, err := clickHouseServiceV2.DescribeClickHouseEnterpriseDbClusterComputingGroup(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_click_house_enterprise_db_cluster_computing_group DescribeClickHouseEnterpriseDbClusterComputingGroup Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("computing_group_description", objectRaw["ComputingGroupDescription"])
	d.Set("is_readonly", objectRaw["IsReadonly"])
	d.Set("node_count", objectRaw["NodeCount"])
	d.Set("node_scale_max", objectRaw["NodeScaleMax"])
	d.Set("node_scale_min", objectRaw["NodeScaleMin"])
	d.Set("computing_group_id", objectRaw["ComputingGroupId"])

	parts := strings.Split(d.Id(), ":")
	d.Set("db_instance_id", parts[0])

	return nil
}

func resourceAliCloudClickHouseEnterpriseDbClusterComputingGroupUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	update := false
	d.Partial(true)

	var err error
	parts := strings.Split(d.Id(), ":")
	action := "ModifyDBInstanceClass"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["ComputingGroupId"] = parts[1]
	request["DBInstanceId"] = parts[0]
	request["RegionId"] = client.RegionId
	if !d.IsNewResource() && d.HasChange("node_scale_max") {
		update = true
	}
	request["NodeScaleMax"] = d.Get("node_scale_max")
	if !d.IsNewResource() && d.HasChange("node_scale_min") {
		update = true
	}
	request["NodeScaleMin"] = d.Get("node_scale_min")
	if !d.IsNewResource() && d.HasChange("node_count") {
		update = true
	}
	request["NodeCount"] = d.Get("node_count")
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
		stateConf := BuildStateConf([]string{}, []string{"activation"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, clickHouseServiceV2.ClickHouseEnterpriseDbClusterComputingGroupStateRefreshFunc(d.Id(), "$.ComputingGroupStatus", []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
	}
	update = false
	parts = strings.Split(d.Id(), ":")
	action = "ModifyComputingGroupDescription"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["ComputingGroupId"] = parts[1]
	request["DBInstanceId"] = parts[0]
	request["RegionId"] = client.RegionId
	if !d.IsNewResource() && d.HasChange("computing_group_description") {
		update = true
		request["ComputingGroupDescription"] = d.Get("computing_group_description")
	}

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
	}
	update = false
	parts = strings.Split(d.Id(), ":")
	action = "ModifyComputingGroupRWConfig"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["ComputingGroupId"] = parts[1]
	request["DBInstanceId"] = parts[0]
	request["RegionId"] = client.RegionId
	if !d.IsNewResource() && d.HasChange("is_readonly") {
		update = true
	}
	request["IsReadonly"] = d.Get("is_readonly")
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
		stateConf := BuildStateConf([]string{}, []string{"activation"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, clickHouseServiceV2.ClickHouseEnterpriseDbClusterComputingGroupStateRefreshFunc(d.Id(), "$.ComputingGroupStatus", []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
	}

	d.Partial(false)
	return resourceAliCloudClickHouseEnterpriseDbClusterComputingGroupRead(d, meta)
}

func resourceAliCloudClickHouseEnterpriseDbClusterComputingGroupDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	parts := strings.Split(d.Id(), ":")
	action := "DeleteComputingGroup"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	request["ComputingGroupId"] = parts[1]
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
		if IsExpectedErrors(err, []string{"InvalidComputingGroupId.NotFound"}) || NotFoundError(err) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	clickHouseServiceV2 := ClickHouseServiceV2{client}
	stateConf := BuildStateConf([]string{}, []string{"ACTIVATION"}, d.Timeout(schema.TimeoutDelete), 5*time.Second, clickHouseServiceV2.DescribeAsyncClickHouseEnterpriseDbClusterComputingGroupStateRefreshFunc(d, response, "$.Data.Status", []string{}))
	if jobDetail, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id(), jobDetail)
	}

	return nil
}
