// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package gpdb

import (
	"encoding/json"
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
	"github.com/aliyun/terraform-provider-alicloud/alicloud/validators"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/tidwall/sjson"
)

func ResourceDbResourceGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceDbResourceGroupCreate,
		Read:   resourceDbResourceGroupRead,
		Update: resourceDbResourceGroupUpdate,
		Delete: resourceDbResourceGroupDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"db_instance_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"resource_group_config": {
				Type:     schema.TypeString,
				Required: true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					equal, _ := helper.CompareJsonTemplateAreEquivalent(old, new)
					return equal
				},
				StateFunc: func(v interface{}) string {
					jsonString, _ := validators.NormalizeJsonString(v)
					return jsonString
				},
			},
			"resource_group_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"role_list": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func resourceDbResourceGroupCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "CreateDBResourceGroup"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	query["DBInstanceId"] = d.Get("db_instance_id")
	query["ResourceGroupName"] = d.Get("resource_group_name")

	request["ResourceGroupConfig"] = d.Get("resource_group_config")
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
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, "alicloud_gpdb_db_resource_group", action, errs.AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprintf("%v:%v", query["DBInstanceId"], query["ResourceGroupName"]))

	return resourceDbResourceGroupUpdate(d, meta)
}

func resourceDbResourceGroupRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	gpdbServiceV2 := NewServiceV2(client)

	objectRaw, err := gpdbServiceV2.DescribeDbResourceGroup(d.Id())
	if err != nil {
		if !d.IsNewResource() && errs.NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_gpdb_db_resource_group DescribeGpdbDbResourceGroup Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return errs.WrapError(err)
	}

	if objectRaw["ResourceGroupConfig"] != nil {
		d.Set("resource_group_config", objectRaw["ResourceGroupConfig"])
	}
	if objectRaw["ResourceGroupName"] != nil {
		d.Set("resource_group_name", objectRaw["ResourceGroupName"])
	}

	role1Raw, _ := jsonpath.Get("$.RoleList.Role", objectRaw)
	d.Set("role_list", role1Raw)

	parts := strings.Split(d.Id(), ":")
	d.Set("db_instance_id", parts[0])
	d.Set("resource_group_name", parts[1])

	return nil
}

func resourceDbResourceGroupUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	update := false
	parts := strings.Split(d.Id(), ":")
	action := "ModifyDBResourceGroup"
	var err error
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	query["DBInstanceId"] = parts[0]
	jsonString := "{}"
	jsonString, _ = sjson.Set(jsonString, "ResourceGroupItems.0.ResourceGroupName", parts[1])
	jsonString, _ = sjson.Set(jsonString, "ResourceGroupItems.0.ResourceGroupConfig", d.Get("resource_group_config"))
	err = json.Unmarshal([]byte(jsonString), &request)
	query["ResourceGroupItems"], _ = convert.ConvertArrayObjectToJsonString(request["ResourceGroupItems"])
	if err != nil {
		return errs.WrapError(err)
	}
	if _, ok := d.GetOk("resource_group_config"); ok && d.HasChange("resource_group_config") {
		update = true
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
		resourceGroupConfig, _ := validators.NormalizeJsonString(fmt.Sprint(d.Get("resource_group_config")))
		stateConf := helper.BuildStateConf([]string{}, []string{resourceGroupConfig}, d.Timeout(schema.TimeoutUpdate), 30*time.Second, gpdbServiceV2.DbResourceGroupStateRefreshFunc(d.Id(), "ResourceGroupConfig", []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return errs.WrapErrorf(err, errs.IdMsg, d.Id())
		}
	}

	if d.HasChange("role_list") {
		oldEntry, newEntry := d.GetChange("role_list")
		oldEntrySet := oldEntry.(*schema.Set)
		newEntrySet := newEntry.(*schema.Set)
		removed := oldEntrySet.Difference(newEntrySet)
		added := newEntrySet.Difference(oldEntrySet)

		if removed.Len() > 0 {
			parts := strings.Split(d.Id(), ":")
			action := "UnbindDBResourceGroupWithRole"
			request = make(map[string]interface{})
			query = make(map[string]interface{})
			query["DBInstanceId"] = parts[0]
			query["ResourceGroupName"] = parts[1]
			query["RoleList"] = convert.ConvertListToCommaSeparate(removed.List())
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

		if added.Len() > 0 {
			parts := strings.Split(d.Id(), ":")
			action := "BindDBResourceGroupWithRole"
			request = make(map[string]interface{})
			query = make(map[string]interface{})
			query["DBInstanceId"] = parts[0]
			query["ResourceGroupName"] = parts[1]
			query["RoleList"] = convert.ConvertListToCommaSeparate(added.List())

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

	}
	return resourceDbResourceGroupRead(d, meta)
}

func resourceDbResourceGroupDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	parts := strings.Split(d.Id(), ":")
	action := "DeleteDBResourceGroup"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	query["DBInstanceId"] = parts[0]
	query["ResourceGroupName"] = parts[1]

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
