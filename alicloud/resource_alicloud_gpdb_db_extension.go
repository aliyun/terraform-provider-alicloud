// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAliCloudGpdbDbExtension() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudGpdbDbExtensionCreate,
		Read:   resourceAliCloudGpdbDbExtensionRead,
		Update: resourceAliCloudGpdbDbExtensionUpdate,
		Delete: resourceAliCloudGpdbDbExtensionDelete,
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

func resourceAliCloudGpdbDbExtensionCreate(d *schema.ResourceData, meta interface{}) error {

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

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
		response, err = client.RpcPost("gpdb", "2016-05-03", action, query, request, true)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_gpdb_db_extension", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprintf("%v:%v:%v", request["DBInstanceId"], request["DBNames"], request["Extensions"]))

	gpdbServiceV2 := GpdbServiceV2{client}
	stateConf := BuildStateConf([]string{}, []string{"installed"}, d.Timeout(schema.TimeoutCreate), 20*time.Second, gpdbServiceV2.GpdbDbExtensionStateRefreshFunc(d.Id(), "Status", []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return resourceAliCloudGpdbDbExtensionUpdate(d, meta)
}

func resourceAliCloudGpdbDbExtensionRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	gpdbServiceV2 := GpdbServiceV2{client}

	objectRaw, err := gpdbServiceV2.DescribeGpdbDbExtension(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_gpdb_db_extension DescribeGpdbDbExtension Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
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

func resourceAliCloudGpdbDbExtensionUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}

	gpdbServiceV2 := GpdbServiceV2{client}
	objectRaw, err := gpdbServiceV2.DescribeGpdbDbExtension(d.Id())
	if err != nil {
		return WrapError(err)
	}

	initedIsLatestVersion := false
	if _, ok := d.GetOkExists("is_latest_version"); ok && d.IsNewResource() {
		initedIsLatestVersion = true
	}
	if initedIsLatestVersion || d.HasChange("is_latest_version") {
		target := d.Get("is_latest_version").(bool)

		currentStatus, err := jsonpath.Get("IsLatestVersion", objectRaw)
		if err != nil {
			return WrapErrorf(err, FailedGetAttributeMsg, d.Id(), "IsLatestVersion", objectRaw)
		}
		if formatBool(currentStatus) != target {
			if target == true {
				parts := strings.Split(d.Id(), ":")
				action := "UpgradeExtensions"
				request = make(map[string]interface{})
				query = make(map[string]interface{})
				request["Extensions"] = parts[2]
				request["DatabaseName"] = parts[1]
				request["DBInstanceId"] = parts[0]
				request["RegionId"] = client.RegionId
				wait := incrementalWait(3*time.Second, 5*time.Second)
				err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
					response, err = client.RpcPost("gpdb", "2016-05-03", action, query, request, true)
					if err != nil {
						if NeedRetry(err) {
							wait()
							return retry.RetryableError(err)
						}
						return retry.NonRetryableError(err)
					}
					return nil
				})
				addDebug(action, response, request)
				if err != nil {
					return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
				}

			}
		}
	}

	return resourceAliCloudGpdbDbExtensionRead(d, meta)
}

func resourceAliCloudGpdbDbExtensionDelete(d *schema.ResourceData, meta interface{}) error {

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

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutDelete), func() *retry.RetryError {
		response, err = client.RpcPost("gpdb", "2016-05-03", action, query, request, true)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)

	if err != nil {
		if IsExpectedErrors(err, []string{"Extension.NotInstalled"}) || NotFoundError(err) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	return nil
}
