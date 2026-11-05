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

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceDbExtension() *schema.Resource {
	return &schema.Resource{
		Create: resourceDbExtensionCreate,
		Read:   resourceDbExtensionRead,
		Update: resourceDbExtensionUpdate,
		Delete: resourceDbExtensionDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"current_version": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"db_instance_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"database_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"extension_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"extension_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"is_install_need_restart": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"is_latest_version": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"latest_version": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceDbExtensionCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "CreateExtensions"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	if v, ok := d.GetOk("extension_name"); ok {
		request["Extensions"] = v
	}
	if v, ok := d.GetOk("database_name"); ok {
		request["DBNames"] = v
	}
	if v, ok := d.GetOk("db_instance_id"); ok {
		request["DBInstanceId"] = v
	}
	request["RegionId"] = client.RegionId

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
		return nil
	})
	logs.AddDebug(action, response, request)

	if err != nil {
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, "alicloud_gpdb_db_extension", action, errs.AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprintf("%v:%v:%v", request["DBInstanceId"], request["DBNames"], request["Extensions"]))

	gpdbServiceV2 := NewServiceV2(client)
	stateConf := helper.BuildStateConf([]string{}, []string{"installed"}, d.Timeout(schema.TimeoutCreate), 20*time.Second, gpdbServiceV2.DbExtensionStateRefreshFunc(d.Id(), "Status", []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return errs.WrapErrorf(err, errs.IdMsg, d.Id())
	}

	return resourceDbExtensionUpdate(d, meta)
}

func resourceDbExtensionRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	gpdbServiceV2 := NewServiceV2(client)

	objectRaw, err := gpdbServiceV2.DescribeDbExtension(d.Id())
	if err != nil {
		if !d.IsNewResource() && errs.NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_gpdb_db_extension DescribeGpdbDbExtension Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return errs.WrapError(err)
	}

	d.Set("current_version", objectRaw["CurrentVersion"])
	d.Set("description", objectRaw["Description"])
	d.Set("extension_id", objectRaw["ExtensionId"])
	d.Set("is_install_need_restart", objectRaw["IsInstallNeedRestart"])
	d.Set("is_latest_version", objectRaw["IsLatestVersion"])
	d.Set("latest_version", objectRaw["LatestVersion"])
	d.Set("status", objectRaw["Status"])
	d.Set("extension_name", objectRaw["ExtensionName"])

	parts := strings.Split(d.Id(), ":")
	d.Set("db_instance_id", parts[0])
	d.Set("database_name", parts[1])

	return nil
}

func resourceDbExtensionUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}

	gpdbServiceV2 := NewServiceV2(client)
	objectRaw, err := gpdbServiceV2.DescribeDbExtension(d.Id())
	if err != nil {
		return errs.WrapError(err)
	}

	initedIsLatestVersion := false
	if _, ok := d.GetOkExists("is_latest_version"); ok && d.IsNewResource() {
		initedIsLatestVersion = true
	}
	if initedIsLatestVersion || d.HasChange("is_latest_version") {
		target := d.Get("is_latest_version").(bool)

		currentStatus, err := jsonpath.Get("IsLatestVersion", objectRaw)
		if err != nil {
			return errs.WrapErrorf(err, errs.FailedGetAttributeMsg, d.Id(), "IsLatestVersion", objectRaw)
		}
		currentStatusValue, formatErr := convert.FormatBool(currentStatus)
		if formatErr != nil {
			return formatErr
		}
		if currentStatusValue != target {
			if target == true {
				parts := strings.Split(d.Id(), ":")
				action := "UpgradeExtensions"
				request = make(map[string]interface{})
				query = make(map[string]interface{})
				request["Extensions"] = parts[2]
				request["DatabaseName"] = parts[1]
				request["DBInstanceId"] = parts[0]
				request["RegionId"] = client.RegionId
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
					return nil
				})
				logs.AddDebug(action, response, request)
				if err != nil {
					return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
				}

			}
		}
	}

	return resourceDbExtensionRead(d, meta)
}

func resourceDbExtensionDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	parts := strings.Split(d.Id(), ":")
	action := "DeleteExtension"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	request["Extension"] = parts[2]
	request["DBNames"] = parts[1]
	request["DBInstanceId"] = parts[0]
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
		return nil
	})
	logs.AddDebug(action, response, request)

	if err != nil {
		if errs.IsExpectedErrors(err, []string{"Extension.NotInstalled"}) || errs.NotFoundError(err) {
			return nil
		}
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
	}

	return nil
}
