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
	if v, ok := d.GetOk("db_instance_id"); ok {
		request["DBInstanceId"] = v
	}
	request["RegionId"] = client.RegionId

	if v, ok := d.GetOk("data_source_description"); ok {
		request["DataSourceDescription"] = v
	}
	if v, ok := d.GetOk("data_source_type"); ok {
		request["DataSourceType"] = v
	}
	if v, ok := d.GetOk("jdbc_password"); ok {
		request["JDBCPassword"] = v
	}
	if v, ok := d.GetOk("jdbc_connection_string"); ok {
		request["JDBCConnectionString"] = v
	}
	if v, ok := d.GetOk("data_source_name"); ok {
		request["DataSourceName"] = v
	}
	request["JDBCUserName"] = d.Get("jdbc_user_name")
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
		return nil
	})
	addDebug(action, response, request)

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_gpdb_jdbc_data_source", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprintf("%v:%v", request["DBInstanceId"], response["DataSourceId"]))

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

	d.Set("create_time", objectRaw["CreateTime"])
	d.Set("data_source_description", objectRaw["DataSourceDescription"])
	d.Set("data_source_name", objectRaw["DataSourceName"])
	d.Set("data_source_type", objectRaw["DataSourceType"])
	d.Set("jdbc_connection_string", objectRaw["JDBCConnectionString"])
	d.Set("jdbc_password", objectRaw["JDBCPassword"])
	d.Set("jdbc_user_name", objectRaw["JDBCUserName"])
	d.Set("status", objectRaw["DataSourceStatus"])
	d.Set("data_source_id", objectRaw["DataSourceId"])

	parts := strings.Split(d.Id(), ":")
	d.Set("db_instance_id", parts[0])

	return nil
}

func resourceAliCloudGpdbJdbcDataSourceUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	update := false

	var err error
	parts := strings.Split(d.Id(), ":")
	action := "ModifyJDBCDataSource"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["DataSourceId"] = parts[1]
	request["DBInstanceId"] = parts[0]
	request["RegionId"] = client.RegionId
	if d.HasChange("data_source_description") {
		update = true
		request["DataSourceDescription"] = d.Get("data_source_description")
	}

	if d.HasChange("data_source_type") {
		update = true
		request["DataSourceType"] = d.Get("data_source_type")
	}

	if d.HasChange("jdbc_password") {
		update = true
		request["JDBCPassword"] = d.Get("jdbc_password")
	}

	if d.HasChange("jdbc_connection_string") {
		update = true
		request["JDBCConnectionString"] = d.Get("jdbc_connection_string")
	}

	if d.HasChange("jdbc_user_name") {
		update = true
	}
	request["JDBCUserName"] = d.Get("jdbc_user_name")
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
			return nil
		})
		addDebug(action, response, request)
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
	request["DataSourceId"] = parts[1]
	request["DBInstanceId"] = parts[0]
	request["RegionId"] = client.RegionId

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
		return nil
	})
	addDebug(action, response, request)

	if err != nil {
		if NotFoundError(err) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	return nil
}
