// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package gpdb

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceExternalDataService() *schema.Resource {
	return &schema.Resource{
		Create: resourceExternalDataServiceCreate,
		Read:   resourceExternalDataServiceRead,
		Update: resourceExternalDataServiceUpdate,
		Delete: resourceExternalDataServiceDelete,
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
			"service_description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"service_id": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"service_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"service_spec": {
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

func resourceExternalDataServiceCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "CreateExternalDataService"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	query["DBInstanceId"] = d.Get("db_instance_id")
	request["RegionId"] = client.RegionId

	request["ServiceName"] = d.Get("service_name")
	request["ServiceSpec"] = d.Get("service_spec")
	if v, ok := d.GetOk("service_description"); ok {
		request["ServiceDescription"] = v
	}
	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
		response, err = client.RpcPost("gpdb", "2016-05-03", action, query, request, true)
		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		logs.AddDebug(action, response, request)
		return nil
	})

	if err != nil {
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, "alicloud_gpdb_external_data_service", action, errs.AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprintf("%v:%v", query["DBInstanceId"], response["ServiceId"]))

	gpdbServiceV2 := NewServiceV2(client)
	stateConf := helper.BuildStateConf([]string{}, []string{"Running", "running"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, gpdbServiceV2.ExternalDataServiceStateRefreshFunc(d.Id(), "Status", []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return errs.WrapErrorf(err, errs.IdMsg, d.Id())
	}

	return resourceExternalDataServiceRead(d, meta)
}

func resourceExternalDataServiceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	gpdbServiceV2 := NewServiceV2(client)

	objectRaw, err := gpdbServiceV2.DescribeExternalDataService(d.Id())
	if err != nil {
		if !d.IsNewResource() && errs.NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_gpdb_external_data_service DescribeGpdbExternalDataService Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return errs.WrapError(err)
	}

	if objectRaw["CreateTime"] != nil {
		d.Set("create_time", objectRaw["CreateTime"])
	}
	if objectRaw["ServiceDescription"] != nil {
		d.Set("service_description", objectRaw["ServiceDescription"])
	}
	if objectRaw["ServiceName"] != nil {
		d.Set("service_name", objectRaw["ServiceName"])
	}
	if objectRaw["ServiceSpec"] != nil {
		d.Set("service_spec", objectRaw["ServiceSpec"])
	}
	if objectRaw["Status"] != nil {
		d.Set("status", objectRaw["Status"])
	}
	if objectRaw["ServiceId"] != nil {
		serviceId, err := strconv.ParseInt(fmt.Sprint(objectRaw["ServiceId"]), 10, 64)
		if err != nil {
			return errs.WrapError(err)
		}
		d.Set("service_id", serviceId)
	}

	parts := strings.Split(d.Id(), ":")
	d.Set("db_instance_id", parts[0])

	return nil
}

func resourceExternalDataServiceUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	update := false
	parts := strings.Split(d.Id(), ":")
	action := "ModifyExternalDataService"
	var err error
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	query["DBInstanceId"] = parts[0]
	query["ServiceId"] = parts[1]
	request["RegionId"] = client.RegionId
	if d.HasChange("service_spec") {
		update = true
	}
	request["ServiceSpec"] = d.Get("service_spec")
	if d.HasChange("service_description") {
		update = true
		request["ServiceDescription"] = d.Get("service_description")
	}

	if update {
		wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
		err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
			response, err = client.RpcPost("gpdb", "2016-05-03", action, query, request, true)
			if err != nil {
				if helper.NeedRetry(err) {
					wait()
					return retry.RetryableError(err)
				}
				return retry.NonRetryableError(err)
			}
			logs.AddDebug(action, response, request)
			return nil
		})
		if err != nil {
			return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
		}
		gpdbServiceV2 := NewServiceV2(client)
		stateConf := helper.BuildStateConf([]string{}, []string{"Running", "running"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, gpdbServiceV2.ExternalDataServiceStateRefreshFunc(d.Id(), "Status", []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return errs.WrapErrorf(err, errs.IdMsg, d.Id())
		}
	}

	return resourceExternalDataServiceRead(d, meta)
}

func resourceExternalDataServiceDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	parts := strings.Split(d.Id(), ":")
	action := "DeleteExternalDataService"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	query["ServiceId"] = parts[1]
	query["DBInstanceId"] = parts[0]
	request["RegionId"] = client.RegionId

	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutDelete), func() *retry.RetryError {
		response, err = client.RpcPost("gpdb", "2016-05-03", action, query, request, true)

		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		logs.AddDebug(action, response, request)
		return nil
	})

	if err != nil {
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
	}

	return nil
}
