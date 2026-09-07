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

func resourceAliCloudPolarDbExtension() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudPolarDbExtensionCreate,
		Read:   resourceAliCloudPolarDbExtensionRead,
		Update: resourceAliCloudPolarDbExtensionUpdate,
		Delete: resourceAliCloudPolarDbExtensionDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"account_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"db_cluster_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"db_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"default_version": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"extension_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"installed_version": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		},
	}
}

func resourceAliCloudPolarDbExtensionCreate(d *schema.ResourceData, meta interface{}) error {

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
	if v, ok := d.GetOk("account_name"); ok {
		request["AccountName"] = v
	}
	if v, ok := d.GetOk("db_cluster_id"); ok {
		request["DBClusterId"] = v
	}
	if v, ok := d.GetOk("db_name"); ok {
		request["DBNames"] = v
	}
	request["RegionId"] = client.RegionId
	request["ClientToken"] = buildClientToken(action)

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.RpcPost("polardb", "2017-08-01", action, query, request, true)
		if err != nil {
			if IsExpectedErrors(err, []string{"OperationDenied.DBClusterStatus"}) || NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_polar_db_extension", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprintf("%v:%v:%v:%v", request["DBClusterId"], request["AccountName"], request["DBNames"], request["Extensions"]))

	polarDbServiceV2 := PolarDbServiceV2{client}
	stateConf := BuildStateConf([]string{}, []string{"#CHECKSET"}, d.Timeout(schema.TimeoutCreate), 30*time.Second, polarDbServiceV2.PolarDbExtensionStateRefreshFunc(d.Id(), "#$.Name", []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return resourceAliCloudPolarDbExtensionUpdate(d, meta)
}

func resourceAliCloudPolarDbExtensionRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	polarDbServiceV2 := PolarDbServiceV2{client}

	objectRaw, err := polarDbServiceV2.DescribePolarDbExtension(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_polar_db_extension DescribePolarDbExtension Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("default_version", objectRaw["DefaultVersion"])
	d.Set("installed_version", objectRaw["InstalledVersion"])
	d.Set("account_name", objectRaw["Owner"])
	d.Set("extension_name", objectRaw["Name"])

	parts := strings.Split(d.Id(), ":")
	d.Set("db_cluster_id", parts[0])
	d.Set("db_name", parts[2])

	return nil
}

func resourceAliCloudPolarDbExtensionUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}

	polarDbServiceV2 := PolarDbServiceV2{client}
	objectRaw, _ := polarDbServiceV2.DescribePolarDbExtension(d.Id())

	if d.HasChange("installed_version") {
		var err error
		target := d.Get("installed_version").(string)
		currentInstalledVersion, err := jsonpath.Get("$.InstalledVersion", objectRaw)
		latestInstalledVersion, err := jsonpath.Get("$.DefaultVersion", objectRaw)
		if currentInstalledVersion != target {
			if target != latestInstalledVersion {
				return WrapErrorf(err, InvalidAttributeValue, "installed_version", latestInstalledVersion)
			} else {
				parts := strings.Split(d.Id(), ":")
				action := "UpdateExtensions"
				request = make(map[string]interface{})
				query = make(map[string]interface{})
				request["Extensions"] = parts[3]
				request["DBClusterId"] = parts[0]
				request["DBNames"] = parts[2]
				request["RegionId"] = client.RegionId
				request["ClientToken"] = buildClientToken(action)
				wait := incrementalWait(3*time.Second, 5*time.Second)
				err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
					response, err = client.RpcPost("polardb", "2017-08-01", action, query, request, true)
					if err != nil {
						if IsExpectedErrors(err, []string{"OperationDenied.DBClusterStatus"}) || NeedRetry(err) {
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
				polarDbServiceV2 := PolarDbServiceV2{client}
				stateConf := BuildStateConf([]string{}, []string{fmt.Sprint(d.Get("installed_version"))}, d.Timeout(schema.TimeoutUpdate), 10*time.Second, polarDbServiceV2.PolarDbExtensionStateRefreshFunc(d.Id(), "$.InstalledVersion", []string{}))
				if _, err := stateConf.WaitForState(); err != nil {
					return WrapErrorf(err, IdMsg, d.Id())
				}

			}
		}
	}

	return resourceAliCloudPolarDbExtensionRead(d, meta)
}

func resourceAliCloudPolarDbExtensionDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	parts := strings.Split(d.Id(), ":")
	action := "DeleteExtensions"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	request["Extensions"] = parts[3]
	request["DBClusterId"] = parts[0]
	request["DBNames"] = parts[2]
	request["RegionId"] = client.RegionId
	request["ClientToken"] = buildClientToken(action)

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = client.RpcPost("polardb", "2017-08-01", action, query, request, true)
		if err != nil {
			if IsExpectedErrors(err, []string{"OperationDenied.DBClusterStatus"}) || NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)

	if err != nil {
		if IsExpectedErrors(err, []string{"Extension.NotInstall"}) || NotFoundError(err) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	polarDbServiceV2 := PolarDbServiceV2{client}
	stateConf := BuildStateConf([]string{}, []string{""}, d.Timeout(schema.TimeoutDelete), 10*time.Second, polarDbServiceV2.PolarDbExtensionStateRefreshFunc(d.Id(), "$.Name", []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return nil
}
