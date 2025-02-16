// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAliCloudGpdbJdbcDataSource() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudGpdbJdbcDataSourceCreate,
		Read:   resourceAliCloudGpdbJdbcDataSourceRead,
		Update: resourceAliCloudGpdbJdbcDataSourceUpdate,
		Delete: resourceAliCloudGpdbJdbcDataSourceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"data_source_description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"data_source_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"data_source_name": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"data_source_type": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"db_instance_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"jdbc_connection_string": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"jdbc_password": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"jdbc_user_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceAliCloudGpdbJdbcDataSourceCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "CreateJDBCDataSource"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	query["DBInstanceId"] = d.Get("db_instance_id")
	query["RegionId"] = client.RegionId

	if v, ok := d.GetOk("data_source_name"); ok {
		request["DataSourceName"] = v
	}
	if v, ok := d.GetOk("data_source_description"); ok {
		request["DataSourceDescription"] = v
	}
	if v, ok := d.GetOk("data_source_type"); ok {
		request["DataSourceType"] = v
	}
	if v, ok := d.GetOk("jdbc_connection_string"); ok {
		request["JDBCConnectionString"] = v
	}
	request["JDBCUserName"] = d.Get("jdbc_user_name")
	if v, ok := d.GetOk("jdbc_password"); ok {
		request["JDBCPassword"] = v
	}
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.RpcPost("gpdb", "2016-05-03", action, query, request, true)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(action, response, request)
		return nil
	})

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_gpdb_jdbc_data_source", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprintf("%v:%v", query["DBInstanceId"], response["DataSourceId"]))

	gpdbServiceV2 := GpdbServiceV2{client}
	stateConf := BuildStateConf([]string{}, []string{"running"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, gpdbServiceV2.GpdbJdbcDataSourceStateRefreshFunc(d.Id(), "DataSourceStatus", []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return resourceAliCloudGpdbJdbcDataSourceRead(d, meta)
}

func resourceAliCloudGpdbJdbcDataSourceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	gpdbServiceV2 := GpdbServiceV2{client}

	objectRaw, err := gpdbServiceV2.DescribeGpdbJdbcDataSource(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_gpdb_jdbc_data_source DescribeGpdbJdbcDataSource Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	if objectRaw["CreateTime"] != nil {
		d.Set("create_time", objectRaw["CreateTime"])
	}
	if objectRaw["DataSourceDescription"] != nil {
		d.Set("data_source_description", objectRaw["DataSourceDescription"])
	}
	if objectRaw["DataSourceName"] != nil {
		d.Set("data_source_name", objectRaw["DataSourceName"])
	}
	if objectRaw["DataSourceType"] != nil {
		d.Set("data_source_type", objectRaw["DataSourceType"])
	}
	if objectRaw["JDBCConnectionString"] != nil {
		d.Set("jdbc_connection_string", objectRaw["JDBCConnectionString"])
	}
	if objectRaw["JDBCPassword"] != nil {
		d.Set("jdbc_password", objectRaw["JDBCPassword"])
	}
	if objectRaw["JDBCUserName"] != nil {
		d.Set("jdbc_user_name", objectRaw["JDBCUserName"])
	}
	if objectRaw["DataSourceStatus"] != nil {
		d.Set("status", objectRaw["DataSourceStatus"])
	}
	if objectRaw["DataSourceId"] != nil {
		d.Set("data_source_id", objectRaw["DataSourceId"])
	}

	parts := strings.Split(d.Id(), ":")
	d.Set("db_instance_id", parts[0])
	d.Set("data_source_id", parts[1])

	return nil
}

func resourceAliCloudGpdbJdbcDataSourceUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	update := false
	parts := strings.Split(d.Id(), ":")
	action := "ModifyJDBCDataSource"
	var err error
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	query["DataSourceId"] = parts[1]
	query["DBInstanceId"] = parts[0]
	query["RegionId"] = client.RegionId
	if d.HasChange("data_source_description") {
		update = true
		request["DataSourceDescription"] = d.Get("data_source_description")
	}

	if d.HasChange("data_source_type") {
		update = true
		request["DataSourceType"] = d.Get("data_source_type")
	}

	if d.HasChange("jdbc_connection_string") {
		update = true
		request["JDBCConnectionString"] = d.Get("jdbc_connection_string")
	}

	if d.HasChange("jdbc_user_name") {
		update = true
	}
	request["JDBCUserName"] = d.Get("jdbc_user_name")
	if d.HasChange("jdbc_password") {
		update = true
	}
	request["JDBCPassword"] = d.Get("jdbc_password")

	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("gpdb", "2016-05-03", action, query, request, true)
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			addDebug(action, response, request)
			return nil
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
		gpdbServiceV2 := GpdbServiceV2{client}
		stateConf := BuildStateConf([]string{}, []string{"running"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, gpdbServiceV2.GpdbJdbcDataSourceStateRefreshFunc(d.Id(), "DataSourceStatus", []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
	}

	return resourceAliCloudGpdbJdbcDataSourceRead(d, meta)
}

func resourceAliCloudGpdbJdbcDataSourceDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	parts := strings.Split(d.Id(), ":")
	action := "DeleteJDBCDataSource"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	query["DataSourceId"] = parts[1]
	query["DBInstanceId"] = parts[0]
	query["RegionId"] = client.RegionId

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = client.RpcPost("gpdb", "2016-05-03", action, query, request, true)

		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(action, response, request)
		return nil
	})

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	return nil
}
