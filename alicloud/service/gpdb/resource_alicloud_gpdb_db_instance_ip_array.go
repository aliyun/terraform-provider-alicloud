// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package gpdb

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/convert"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"

	"github.com/blues/jsonata-go"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceDBInstanceIPArray() *schema.Resource {
	return &schema.Resource{
		Create: resourceDBInstanceIPArrayCreate,
		Read:   resourceDBInstanceIPArrayRead,
		Update: resourceDBInstanceIPArrayUpdate,
		Delete: resourceDBInstanceIPArrayDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"db_instance_ip_array_attribute": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"db_instance_ip_array_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"db_instance_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"modify_mode": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"security_ip_list": {
				Type:     schema.TypeSet,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func resourceDBInstanceIPArrayCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "CreateDBInstanceIPArray"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	query["DBInstanceId"] = d.Get("db_instance_id")
	query["IPArrayName"] = d.Get("db_instance_ip_array_name")

	if v, ok := d.GetOk("db_instance_ip_array_attribute"); ok {
		request["IPArrayAttribute"] = v
	}
	jsonPathResult1, err := jsonpath.Get("$", d.Get("security_ip_list"))
	if err == nil {
		request["SecurityIPList"] = convert.ConvertListToCommaSeparate(jsonPathResult1.(*schema.Set).List())
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
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, "alicloud_gpdb_db_instance_ip_array", action, errs.AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprintf("%v:%v", query["DBInstanceId"], query["IPArrayName"]))

	return resourceDBInstanceIPArrayRead(d, meta)
}

func resourceDBInstanceIPArrayRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	gpdbServiceV2 := NewServiceV2(client)

	objectRaw, err := gpdbServiceV2.DescribeDBInstanceIPArray(d.Id())
	if err != nil {
		if !d.IsNewResource() && errs.NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_gpdb_db_instance_ip_array DescribeGpdbDBInstanceIPArray Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return errs.WrapError(err)
	}

	if objectRaw["DBInstanceIPArrayAttribute"] != nil {
		d.Set("db_instance_ip_array_attribute", objectRaw["DBInstanceIPArrayAttribute"])
	}
	if objectRaw["DBInstanceIPArrayName"] != nil {
		d.Set("db_instance_ip_array_name", objectRaw["DBInstanceIPArrayName"])
	}

	e := jsonata.MustCompile("$split($.SecurityIPList, \",\")")
	evaluation, _ := e.Eval(objectRaw)
	d.Set("security_ip_list", evaluation)

	parts := strings.Split(d.Id(), ":")
	d.Set("db_instance_id", parts[0])
	d.Set("db_instance_ip_array_name", parts[1])

	return nil
}

func resourceDBInstanceIPArrayUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	update := false
	parts := strings.Split(d.Id(), ":")
	action := "ModifySecurityIps"
	var err error
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	query["DBInstanceIPArrayName"] = parts[1]
	query["DBInstanceId"] = parts[0]

	if v, ok := d.GetOk("modify_mode"); ok {
		request["ModifyMode"] = v
	}
	if d.HasChange("security_ip_list") {
		update = true
	}
	jsonPathResult1, err := jsonpath.Get("$", d.Get("security_ip_list"))
	if err == nil {
		request["SecurityIPList"] = convert.ConvertListToCommaSeparate(jsonPathResult1.(*schema.Set).List())
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
	}

	return resourceDBInstanceIPArrayRead(d, meta)
}

func resourceDBInstanceIPArrayDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	parts := strings.Split(d.Id(), ":")
	action := "DeleteDBInstanceIPArray"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	query["DBInstanceId"] = parts[0]
	query["IPArrayName"] = parts[1]

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
