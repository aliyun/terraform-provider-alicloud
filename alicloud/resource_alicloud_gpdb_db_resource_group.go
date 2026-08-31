// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/tidwall/sjson"
)

func resourceAliCloudGpdbDbResourceGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudGpdbDbResourceGroupCreate,
		Read:   resourceAliCloudGpdbDbResourceGroupRead,
		Update: resourceAliCloudGpdbDbResourceGroupUpdate,
		Delete: resourceAliCloudGpdbDbResourceGroupDelete,
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
					equal, _ := compareJsonTemplateAreEquivalent(old, new)
					return equal
				},
				StateFunc: func(v interface{}) string {
					jsonString, _ := normalizeJsonString(v)
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

func resourceAliCloudGpdbDbResourceGroupCreate(d *schema.ResourceData, meta interface{}) error {

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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_gpdb_db_resource_group", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprintf("%v:%v", query["DBInstanceId"], query["ResourceGroupName"]))

	return resourceAliCloudGpdbDbResourceGroupUpdate(d, meta)
}

func resourceAliCloudGpdbDbResourceGroupRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	gpdbServiceV2 := GpdbServiceV2{client}

	objectRaw, err := gpdbServiceV2.DescribeGpdbDbResourceGroup(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_gpdb_db_resource_group DescribeGpdbDbResourceGroup Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
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

func resourceAliCloudGpdbDbResourceGroupUpdate(d *schema.ResourceData, meta interface{}) error {
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
	query["ResourceGroupItems"], _ = convertArrayObjectToJsonString(request["ResourceGroupItems"])
	if err != nil {
		return WrapError(err)
	}
	if _, ok := d.GetOk("resource_group_config"); ok && d.HasChange("resource_group_config") {
		update = true
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
		resourceGroupConfig, _ := normalizeJsonString(fmt.Sprint(d.Get("resource_group_config")))
		stateConf := BuildStateConf([]string{}, []string{resourceGroupConfig}, d.Timeout(schema.TimeoutUpdate), 30*time.Second, gpdbServiceV2.GpdbDbResourceGroupStateRefreshFunc(d.Id(), "ResourceGroupConfig", []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
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
			query["RoleList"] = convertListToCommaSeparate(removed.List())
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

		}

		if added.Len() > 0 {
			parts := strings.Split(d.Id(), ":")
			action := "BindDBResourceGroupWithRole"
			request = make(map[string]interface{})
			query = make(map[string]interface{})
			query["DBInstanceId"] = parts[0]
			query["ResourceGroupName"] = parts[1]
			query["RoleList"] = convertListToCommaSeparate(added.List())

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

		}

	}
	return resourceAliCloudGpdbDbResourceGroupRead(d, meta)
}

func resourceAliCloudGpdbDbResourceGroupDelete(d *schema.ResourceData, meta interface{}) error {

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
