package alicloud

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAliCloudRdsAiInstance() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudRdsAiInstanceCreate,
		Read:   resourceAliCloudRdsAiInstanceRead,
		Update: resourceAliCloudRdsAiInstanceUpdate,
		Delete: resourceAliCloudRdsAiInstanceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(20 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"app_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"app_type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"auth_config_list": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"value": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"name": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"ca_type": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"db_instance_name": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"dashboard_password": {
				Type:      schema.TypeString,
				Optional:  true,
				Sensitive: true,
			},
			"database_password": {
				Type:      schema.TypeString,
				Optional:  true,
				Sensitive: true,
			},
			"initialize_with_existing_data": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"public_endpoint_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"public_network_access_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"ssl_enabled": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"server_cert": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"server_key": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"status": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"storage_config_list": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"value": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"name": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
		},
	}
}

func resourceAliCloudRdsAiInstanceCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "CreateAppInstance"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	request["RegionId"] = client.RegionId
	request["ClientToken"] = buildClientToken(action)

	if v, ok := d.GetOk("dashboard_password"); ok {
		request["DashboardPassword"] = v
	}
	request["DashboardUsername"] = "supabase"
	if v, ok := d.GetOkExists("initialize_with_existing_data"); ok {
		request["InitializeWithExistingData"] = v
	}
	request["InstanceClass"] = "rdsai.supabase.basic"
	request["AppName"] = d.Get("app_name")
	request["AppType"] = d.Get("app_type")
	if v, ok := d.GetOk("db_instance_name"); ok {
		request["DBInstanceName"] = v
	}
	if v, ok := d.GetOkExists("public_network_access_enabled"); ok {
		request["PublicNetworkAccessEnabled"] = v
	}
	if v, ok := d.GetOk("database_password"); ok {
		request["DatabasePassword"] = v
	}
	if v, ok := d.GetOkExists("public_endpoint_enabled"); ok {
		request["PublicEndpointEnabled"] = v
	}
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.RpcPost("RdsAi", "2025-05-07", action, query, request, true)
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_rds_ai_instance", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(response["InstanceName"]))

	rdsAiServiceV2 := RdsAiServiceV2{client}
	stateConf := BuildStateConf([]string{}, []string{"Running"}, d.Timeout(schema.TimeoutCreate), 10*time.Second, rdsAiServiceV2.DescribeAsyncRdsAiInstanceStateRefreshFunc(d, response, "$.Status", []string{}))
	if jobDetail, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id(), jobDetail)
	}

	return resourceAliCloudRdsAiInstanceUpdate(d, meta)
}

func resourceAliCloudRdsAiInstanceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	rdsAiServiceV2 := RdsAiServiceV2{client}

	objectRaw, err := rdsAiServiceV2.DescribeRdsAiInstance(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_rds_ai_instance DescribeRdsAiInstance Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("app_name", objectRaw["AppName"])
	d.Set("app_type", objectRaw["AppType"])
	d.Set("db_instance_name", objectRaw["DBInstanceName"])
	d.Set("status", objectRaw["Status"])

	objectRaw, err = rdsAiServiceV2.DescribeInstanceDescribeInstanceAuthInfo(d.Id())
	if err != nil && !NotFoundError(err) {
		return WrapError(err)
	}

	authConfigListNameArray := make([]string, 0)
	if v, ok := d.GetOk("auth_config_list"); ok {
		for _, dataLoop := range convertToInterfaceArray(v) {
			dataLoopTmp := dataLoop.(map[string]interface{})
			authConfigListNameArray = append(authConfigListNameArray, fmt.Sprint(dataLoopTmp["name"]))
		}
	}

	configListRaw := objectRaw["ConfigList"]
	authConfigListMaps := make([]map[string]interface{}, 0)
	if configListRaw != nil {
		for _, configListChildRaw := range convertToInterfaceArray(configListRaw) {
			authConfigListMap := make(map[string]interface{})
			configListChildRaw := configListChildRaw.(map[string]interface{})
			if !InArray(fmt.Sprint(configListChildRaw["Name"]), authConfigListNameArray) {
				continue
			}

			authConfigListMap["name"] = configListChildRaw["Name"]
			authConfigListMap["value"] = configListChildRaw["Value"]

			authConfigListMaps = append(authConfigListMaps, authConfigListMap)
		}
	}
	if err := d.Set("auth_config_list", authConfigListMaps); err != nil {
		return err
	}

	objectRaw, err = rdsAiServiceV2.DescribeInstanceDescribeInstanceSSL(d.Id())
	if err != nil && !NotFoundError(err) {
		return WrapError(err)
	}

	d.Set("ca_type", objectRaw["CAType"])
	d.Set("ssl_enabled", formatInt(objectRaw["SSLEnabled"]))
	d.Set("server_cert", objectRaw["ServerCert"])
	d.Set("server_key", objectRaw["ServerKey"])

	objectRaw, err = rdsAiServiceV2.DescribeInstanceDescribeInstanceStorageConfig(d.Id())
	if err != nil && !NotFoundError(err) {
		return WrapError(err)
	}

	storageConfigListNameArray := make([]string, 0)
	if v, ok := d.GetOk("storage_config_list"); ok {
		for _, dataLoop := range convertToInterfaceArray(v) {
			dataLoopTmp := dataLoop.(map[string]interface{})
			storageConfigListNameArray = append(storageConfigListNameArray, fmt.Sprint(dataLoopTmp["name"]))
		}

	}

	configListRaw = objectRaw["ConfigList"]
	storageConfigListMaps := make([]map[string]interface{}, 0)
	if configListRaw != nil {
		for _, configListChildRaw := range convertToInterfaceArray(configListRaw) {
			storageConfigListMap := make(map[string]interface{})
			configListChildRaw := configListChildRaw.(map[string]interface{})
			if !InArray(fmt.Sprint(configListChildRaw["Name"]), storageConfigListNameArray) {
				continue
			}

			storageConfigListMap["name"] = configListChildRaw["Name"]
			storageConfigListMap["value"] = configListChildRaw["Value"]

			storageConfigListMaps = append(storageConfigListMaps, storageConfigListMap)
		}
	}
	if err := d.Set("storage_config_list", storageConfigListMaps); err != nil {
		return err
	}

	return nil
}

func resourceAliCloudRdsAiInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	update := false
	d.Partial(true)

	rdsAiServiceV2 := RdsAiServiceV2{client}
	objectRaw, _ := rdsAiServiceV2.DescribeRdsAiInstance(d.Id())

	if d.HasChange("status") {
		var err error
		target := d.Get("status").(string)

		currentStatus, err := jsonpath.Get("Status", objectRaw)
		if err != nil {
			return WrapErrorf(err, FailedGetAttributeMsg, d.Id(), "Status", objectRaw)
		}
		if fmt.Sprint(currentStatus) != target {
			enableStopInstanceStopped := false
			checkValue00 := objectRaw["Status"]
			if checkValue00 == "Running" {
				enableStopInstanceStopped = true
			}
			if enableStopInstanceStopped && target == "Stopped" {
				action := "StopInstance"
				request = make(map[string]interface{})
				query = make(map[string]interface{})
				request["InstanceName"] = d.Id()
				request["RegionId"] = client.RegionId
				wait := incrementalWait(3*time.Second, 5*time.Second)
				err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
					response, err = client.RpcPost("RdsAi", "2025-05-07", action, query, request, true)
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
				rdsAiServiceV2 := RdsAiServiceV2{client}
				stateConf := BuildStateConf([]string{}, []string{"Stopped"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, rdsAiServiceV2.DescribeAsyncRdsAiInstanceStateRefreshFunc(d, response, "$.Status", []string{}))
				if jobDetail, err := stateConf.WaitForState(); err != nil {
					return WrapErrorf(err, IdMsg, d.Id(), jobDetail)
				}

			}
			enableStartInstanceRunning := false
			checkValue00 = objectRaw["Status"]
			if checkValue00 == "Stopped" {
				enableStartInstanceRunning = true
			}
			if enableStartInstanceRunning && target == "Running" {
				action := "StartInstance"
				request = make(map[string]interface{})
				query = make(map[string]interface{})
				request["InstanceName"] = d.Id()
				request["RegionId"] = client.RegionId
				wait := incrementalWait(3*time.Second, 5*time.Second)
				err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
					response, err = client.RpcPost("RdsAi", "2025-05-07", action, query, request, true)
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
				rdsAiServiceV2 := RdsAiServiceV2{client}
				stateConf := BuildStateConf([]string{}, []string{"Running"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, rdsAiServiceV2.DescribeAsyncRdsAiInstanceStateRefreshFunc(d, response, "$.Status", []string{}))
				if jobDetail, err := stateConf.WaitForState(); err != nil {
					return WrapErrorf(err, IdMsg, d.Id(), jobDetail)
				}

			}
		}
	}

	var err error
	action := "ResetInstancePassword"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["InstanceName"] = d.Id()
	request["RegionId"] = client.RegionId
	if !d.IsNewResource() && d.HasChange("dashboard_password") {
		update = true
		request["DashboardPassword"] = d.Get("dashboard_password")
	}

	if !d.IsNewResource() && d.HasChange("database_password") {
		update = true
		request["DatabasePassword"] = d.Get("database_password")
	}

	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("RdsAi", "2025-05-07", action, query, request, true)
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
		rdsAiServiceV2 := RdsAiServiceV2{client}
		stateConf := BuildStateConf([]string{}, []string{"Running"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, rdsAiServiceV2.DescribeAsyncRdsAiInstanceStateRefreshFunc(d, response, "$.Status", []string{}))
		if jobDetail, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id(), jobDetail)
		}
	}
	update = false
	action = "ModifyInstanceAuthConfig"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["InstanceName"] = d.Id()
	request["RegionId"] = client.RegionId
	if d.HasChange("auth_config_list") {
		update = true
		if v, ok := d.GetOk("auth_config_list"); ok || d.HasChange("auth_config_list") {
			configListMapsArray := make([]interface{}, 0)
			for _, dataLoop := range convertToInterfaceArray(v) {
				dataLoopTmp := dataLoop.(map[string]interface{})
				dataLoopMap := make(map[string]interface{})
				dataLoopMap["Value"] = dataLoopTmp["value"]
				dataLoopMap["Name"] = dataLoopTmp["name"]
				configListMapsArray = append(configListMapsArray, dataLoopMap)
			}
			configListMapsJson, err := json.Marshal(configListMapsArray)
			if err != nil {
				return WrapError(err)
			}
			request["ConfigList"] = string(configListMapsJson)
		}
	}

	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("RdsAi", "2025-05-07", action, query, request, true)
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
		rdsAiServiceV2 := RdsAiServiceV2{client}
		stateConf := BuildStateConf([]string{}, []string{"Running"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, rdsAiServiceV2.DescribeAsyncRdsAiInstanceStateRefreshFunc(d, response, "$.Status", []string{}))
		if jobDetail, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id(), jobDetail)
		}
	}
	update = false
	action = "ModifyInstanceStorageConfig"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["InstanceName"] = d.Id()
	request["RegionId"] = client.RegionId
	request["ClientToken"] = buildClientToken(action)
	if d.HasChange("storage_config_list") {
		update = true
		if v, ok := d.GetOk("storage_config_list"); ok || d.HasChange("storage_config_list") {
			configListMapsArray := make([]interface{}, 0)
			for _, dataLoop := range convertToInterfaceArray(v) {
				dataLoopTmp := dataLoop.(map[string]interface{})
				dataLoopMap := make(map[string]interface{})
				dataLoopMap["Value"] = dataLoopTmp["value"]
				dataLoopMap["Name"] = dataLoopTmp["name"]
				configListMapsArray = append(configListMapsArray, dataLoopMap)
			}
			configListMapsJson, err := json.Marshal(configListMapsArray)
			if err != nil {
				return WrapError(err)
			}
			request["ConfigList"] = string(configListMapsJson)
		}
	}

	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("RdsAi", "2025-05-07", action, query, request, true)
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
		rdsAiServiceV2 := RdsAiServiceV2{client}
		stateConf := BuildStateConf([]string{}, []string{"Running"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, rdsAiServiceV2.DescribeAsyncRdsAiInstanceStateRefreshFunc(d, response, "$.Status", []string{}))
		if jobDetail, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id(), jobDetail)
		}
	}
	update = false
	action = "ModifyInstanceSSL"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["InstanceName"] = d.Id()
	request["RegionId"] = client.RegionId
	if d.HasChange("server_cert") {
		update = true
		request["ServerCert"] = d.Get("server_cert")
	}

	if d.HasChange("server_key") {
		update = true
		request["ServerKey"] = d.Get("server_key")
	}

	if d.HasChange("ca_type") {
		update = true
		request["CAType"] = d.Get("ca_type")
	}

	if d.HasChange("ssl_enabled") {
		update = true
	}
	request["SSLEnabled"] = d.Get("ssl_enabled")
	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("RdsAi", "2025-05-07", action, query, request, true)
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
		rdsAiServiceV2 := RdsAiServiceV2{client}
		stateConf := BuildStateConf([]string{}, []string{"Running"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, rdsAiServiceV2.DescribeAsyncRdsAiInstanceStateRefreshFunc(d, response, "$.Status", []string{}))
		if jobDetail, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id(), jobDetail)
		}
	}

	d.Partial(false)
	return resourceAliCloudRdsAiInstanceRead(d, meta)
}

func resourceAliCloudRdsAiInstanceDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	action := "DeleteAppInstance"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	request["InstanceName"] = d.Id()
	request["RegionId"] = client.RegionId
	request["ClientToken"] = buildClientToken(action)

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = client.RpcPost("RdsAi", "2025-05-07", action, query, request, true)
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

	rdsAiServiceV2 := RdsAiServiceV2{client}
	stateConf := BuildStateConf([]string{}, []string{""}, d.Timeout(schema.TimeoutDelete), 5*time.Second, rdsAiServiceV2.DescribeAsyncRdsAiInstanceStateRefreshFunc(d, response, "$.Status", []string{}))
	if jobDetail, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id(), jobDetail)
	}

	return nil
}
