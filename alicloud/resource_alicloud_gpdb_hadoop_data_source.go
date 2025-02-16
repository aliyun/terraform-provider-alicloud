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

func resourceAliCloudGpdbHadoopDataSource() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudGpdbHadoopDataSourceCreate,
		Read:   resourceAliCloudGpdbHadoopDataSourceRead,
		Update: resourceAliCloudGpdbHadoopDataSourceUpdate,
		Delete: resourceAliCloudGpdbHadoopDataSourceDelete,
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
			"db_instance_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"data_source_description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"data_source_id": {
				Type:     schema.TypeInt,
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
			"emr_instance_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"hadoop_core_conf": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"hadoop_create_type": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"hadoop_hosts_address": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"hdfs_conf": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"hive_conf": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"map_reduce_conf": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"yarn_conf": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceAliCloudGpdbHadoopDataSourceCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "CreateHadoopDataSource"
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
	if v, ok := d.GetOk("hadoop_create_type"); ok {
		request["HadoopCreateType"] = v
	}
	if v, ok := d.GetOk("hadoop_hosts_address"); ok {
		request["HadoopHostsAddress"] = v
	}
	if v, ok := d.GetOk("hadoop_core_conf"); ok {
		request["HadoopCoreConf"] = v
	}
	if v, ok := d.GetOk("hive_conf"); ok {
		request["HiveConf"] = v
	}
	if v, ok := d.GetOk("yarn_conf"); ok {
		request["YarnConf"] = v
	}
	if v, ok := d.GetOk("map_reduce_conf"); ok {
		request["MapReduceConf"] = v
	}
	if v, ok := d.GetOk("hdfs_conf"); ok {
		request["HDFSConf"] = v
	}
	if v, ok := d.GetOk("emr_instance_id"); ok {
		request["EmrInstanceId"] = v
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_gpdb_hadoop_data_source", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprintf("%v:%v", query["DBInstanceId"], response["DataSourceId"]))

	gpdbServiceV2 := GpdbServiceV2{client}
	stateConf := BuildStateConf([]string{}, []string{"running"}, d.Timeout(schema.TimeoutCreate), 30*time.Second, gpdbServiceV2.GpdbHadoopDataSourceStateRefreshFunc(d.Id(), "DataSourceStatus", []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return resourceAliCloudGpdbHadoopDataSourceRead(d, meta)
}

func resourceAliCloudGpdbHadoopDataSourceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	gpdbServiceV2 := GpdbServiceV2{client}

	objectRaw, err := gpdbServiceV2.DescribeGpdbHadoopDataSource(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_gpdb_hadoop_data_source DescribeGpdbHadoopDataSource Failed!!! %s", err)
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
	if objectRaw["EmrInstanceId"] != nil {
		d.Set("emr_instance_id", objectRaw["EmrInstanceId"])
	}
	if objectRaw["HadoopCoreConf"] != nil {
		d.Set("hadoop_core_conf", objectRaw["HadoopCoreConf"])
	}
	if objectRaw["HadoopCreateType"] != nil {
		d.Set("hadoop_create_type", objectRaw["HadoopCreateType"])
	}
	if objectRaw["HadoopHostsAddress"] != nil {
		d.Set("hadoop_hosts_address", objectRaw["HadoopHostsAddress"])
	}
	if objectRaw["HDFSConf"] != nil {
		d.Set("hdfs_conf", objectRaw["HDFSConf"])
	}
	if objectRaw["HiveConf"] != nil {
		d.Set("hive_conf", objectRaw["HiveConf"])
	}
	if objectRaw["MapReduceConf"] != nil {
		d.Set("map_reduce_conf", objectRaw["MapReduceConf"])
	}
	if objectRaw["DataSourceStatus"] != nil {
		d.Set("status", objectRaw["DataSourceStatus"])
	}
	if objectRaw["YarnConf"] != nil {
		d.Set("yarn_conf", objectRaw["YarnConf"])
	}
	if objectRaw["DataSourceId"] != nil {
		d.Set("data_source_id", formatInt(objectRaw["DataSourceId"]))
	}

	parts := strings.Split(d.Id(), ":")
	d.Set("db_instance_id", parts[0])
	d.Set("data_source_id", formatInt(parts[1]))

	return nil
}

func resourceAliCloudGpdbHadoopDataSourceUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	update := false
	parts := strings.Split(d.Id(), ":")
	action := "ModifyHadoopDataSource"
	var err error
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	query["DBInstanceId"] = parts[0]
	query["DataSourceId"] = parts[1]
	query["RegionId"] = client.RegionId
	if d.HasChange("data_source_description") {
		update = true
		request["DataSourceDescription"] = d.Get("data_source_description")
	}

	if d.HasChange("data_source_type") {
		update = true
		request["DataSourceType"] = d.Get("data_source_type")
	}

	if d.HasChange("hadoop_create_type") {
		update = true
		request["HadoopCreateType"] = d.Get("hadoop_create_type")
	}

	if d.HasChange("hadoop_hosts_address") {
		update = true
		request["HadoopHostsAddress"] = d.Get("hadoop_hosts_address")
	}

	if d.HasChange("hadoop_core_conf") {
		update = true
		request["HadoopCoreConf"] = d.Get("hadoop_core_conf")
	}

	if d.HasChange("hive_conf") {
		update = true
		request["HiveConf"] = d.Get("hive_conf")
	}

	if d.HasChange("yarn_conf") {
		update = true
		request["YarnConf"] = d.Get("yarn_conf")
	}

	if d.HasChange("map_reduce_conf") {
		update = true
		request["MapReduceConf"] = d.Get("map_reduce_conf")
	}

	if d.HasChange("hdfs_conf") {
		update = true
		request["HDFSConf"] = d.Get("hdfs_conf")
	}

	if d.HasChange("emr_instance_id") {
		update = true
		request["EmrInstanceId"] = d.Get("emr_instance_id")
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
		stateConf := BuildStateConf([]string{}, []string{"running"}, d.Timeout(schema.TimeoutUpdate), 30*time.Second, gpdbServiceV2.GpdbHadoopDataSourceStateRefreshFunc(d.Id(), "DataSourceStatus", []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
	}

	return resourceAliCloudGpdbHadoopDataSourceRead(d, meta)
}

func resourceAliCloudGpdbHadoopDataSourceDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	parts := strings.Split(d.Id(), ":")
	action := "DeleteHadoopDataSource"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	query["DBInstanceId"] = parts[0]
	query["DataSourceId"] = parts[1]
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
