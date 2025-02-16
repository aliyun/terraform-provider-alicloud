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

func resourceAliCloudGpdbRemoteADBDataSource() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudGpdbRemoteADBDataSourceCreate,
		Read:   resourceAliCloudGpdbRemoteADBDataSourceRead,
		Update: resourceAliCloudGpdbRemoteADBDataSourceUpdate,
		Delete: resourceAliCloudGpdbRemoteADBDataSourceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"data_source_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"local_database": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"local_db_instance_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"manager_user_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"manager_user_password": {
				Type:      schema.TypeString,
				Required:  true,
				Sensitive: true,
			},
			"remote_adb_data_source_id": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"remote_database": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"remote_db_instance_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"user_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"user_password": {
				Type:      schema.TypeString,
				Required:  true,
				Sensitive: true,
			},
		},
	}
}

func resourceAliCloudGpdbRemoteADBDataSourceCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "CreateRemoteADBDataSource"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	query["LocalDBInstanceId"] = d.Get("local_db_instance_id")

	request["LocalDatabase"] = d.Get("local_database")
	request["RemoteDatabase"] = d.Get("remote_database")
	request["UserName"] = d.Get("user_name")
	request["UserPassword"] = d.Get("user_password")
	request["ManagerUserName"] = d.Get("manager_user_name")
	request["ManagerUserPassword"] = d.Get("manager_user_password")
	if v, ok := d.GetOk("data_source_name"); ok {
		request["DataSourceName"] = v
	}
	request["RemoteDBInstanceId"] = d.Get("remote_db_instance_id")
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_gpdb_remote_adb_data_source", action, AlibabaCloudSdkGoERROR)
	}

	DataSourceItemLocalInstanceNameVar, _ := jsonpath.Get("$.DataSourceItem.LocalInstanceName", response)
	DataSourceItemIdVar, _ := jsonpath.Get("$.DataSourceItem.Id", response)
	d.SetId(fmt.Sprintf("%v:%v", DataSourceItemLocalInstanceNameVar, DataSourceItemIdVar))

	gpdbServiceV2 := GpdbServiceV2{client}
	stateConf := BuildStateConf([]string{}, []string{"ACTIVE"}, d.Timeout(schema.TimeoutCreate), 30*time.Second, gpdbServiceV2.GpdbRemoteADBDataSourceStateRefreshFunc(d.Id(), "Status", []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return resourceAliCloudGpdbRemoteADBDataSourceRead(d, meta)
}

func resourceAliCloudGpdbRemoteADBDataSourceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	gpdbServiceV2 := GpdbServiceV2{client}

	objectRaw, err := gpdbServiceV2.DescribeGpdbRemoteADBDataSource(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_gpdb_remote_adb_data_source DescribeGpdbRemoteADBDataSource Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	if objectRaw["DataSourceName"] != nil {
		d.Set("data_source_name", objectRaw["DataSourceName"])
	}
	if objectRaw["LocalDatabase"] != nil {
		d.Set("local_database", objectRaw["LocalDatabase"])
	}
	if objectRaw["ManagerUserName"] != nil {
		d.Set("manager_user_name", objectRaw["ManagerUserName"])
	}
	if objectRaw["RemoteDatabase"] != nil {
		d.Set("remote_database", objectRaw["RemoteDatabase"])
	}
	if objectRaw["RemoteInstanceName"] != nil {
		d.Set("remote_db_instance_id", objectRaw["RemoteInstanceName"])
	}
	if objectRaw["Status"] != nil {
		d.Set("status", objectRaw["Status"])
	}
	if objectRaw["UserName"] != nil {
		d.Set("user_name", objectRaw["UserName"])
	}
	if objectRaw["LocalInstanceName"] != nil {
		d.Set("local_db_instance_id", objectRaw["LocalInstanceName"])
	}
	if objectRaw["Id"] != nil {
		d.Set("remote_adb_data_source_id", objectRaw["Id"])
	}

	parts := strings.Split(d.Id(), ":")
	d.Set("local_db_instance_id", parts[0])
	d.Set("remote_adb_data_source_id", parts[1])

	return nil
}

func resourceAliCloudGpdbRemoteADBDataSourceUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	update := false
	parts := strings.Split(d.Id(), ":")
	action := "ModifyRemoteADBDataSource"
	var err error
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	query["LocalDBInstanceId"] = parts[0]
	query["DataSourceId"] = parts[1]

	if d.HasChange("user_name") {
		update = true
	}
	request["UserName"] = d.Get("user_name")
	if d.HasChange("user_password") {
		update = true
	}
	request["UserPassword"] = d.Get("user_password")
	if d.HasChange("data_source_name") {
		update = true
		request["DataSourceName"] = d.Get("data_source_name")
	}

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
		stateConf := BuildStateConf([]string{}, []string{"ACTIVE"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, gpdbServiceV2.GpdbRemoteADBDataSourceStateRefreshFunc(d.Id(), "Status", []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
	}

	return resourceAliCloudGpdbRemoteADBDataSourceRead(d, meta)
}

func resourceAliCloudGpdbRemoteADBDataSourceDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	parts := strings.Split(d.Id(), ":")
	action := "DeleteRemoteADBDataSource"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	query["LocalDBInstanceId"] = parts[0]
	query["DataSourceId"] = parts[1]

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

	gpdbServiceV2 := GpdbServiceV2{client}
	stateConf := BuildStateConf([]string{}, []string{}, d.Timeout(schema.TimeoutDelete), 30*time.Second, gpdbServiceV2.GpdbRemoteADBDataSourceStateRefreshFunc(d.Id(), "", []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return nil
}
