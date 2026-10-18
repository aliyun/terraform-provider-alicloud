package alicloud

import (
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAlicloudPolarDBPolarFs() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudPolarDBPolarFsCreate,
		Read:   resourceAlicloudPolarDBPolarFsRead,
		Delete: resourceAlicloudPolarDBPolarFsDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(60 * time.Minute),
			Delete: schema.DefaultTimeout(20 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"region_id": {Type: schema.TypeString, Computed: true},
			"storage_type": {
				Type: schema.TypeString, Optional: true, Computed: true, ForceNew: true,
				ValidateFunc: StringInSlice([]string{"local_redundancy", "city_redundancy", "essdpl1", "essdpl0"}, false),
			},
			"db_type": {
				Type: schema.TypeString, Optional: true, Computed: true, ForceNew: true,
				ValidateFunc: StringInSlice([]string{"MySQL", "PostgreSQL"}, false),
			},
			"vpc_id":        {Type: schema.TypeString, Optional: true, Computed: true, ForceNew: true},
			"vswitch_id":    {Type: schema.TypeString, Optional: true, Computed: true, ForceNew: true},
			"zone_id":       {Type: schema.TypeString, Optional: true, Computed: true, ForceNew: true},
			"storage_space": {Type: schema.TypeInt, Optional: true, Computed: true, ForceNew: true, ValidateFunc: IntBetween(10, 100000)},
			"db_cluster_id": {Type: schema.TypeString, Optional: true, Computed: true, ForceNew: true},
			"pay_type": {
				Type: schema.TypeString, Optional: true, Computed: true, ForceNew: true,
				ValidateFunc: StringInSlice([]string{"Postpaid"}, false),
			},
			"accelerate_switch": {Type: schema.TypeString, Optional: true, ForceNew: true, ValidateFunc: StringInSlice([]string{"ONLY", "ON"}, false)},
			"accelerate_storage_size": {
				Type: schema.TypeInt, Optional: true, ForceNew: true,
			},
			"creation_category":               {Type: schema.TypeString, Optional: true, Computed: true, ForceNew: true, ValidateFunc: StringInSlice([]string{"basic", "cold", "high_performance"}, false)},
			"accelerate_type":                 {Type: schema.TypeString, Optional: true, Computed: true, ForceNew: true, ValidateFunc: StringInSlice([]string{"juice", "alluxio"}, false)},
			"polar_fs_path":                   {Type: schema.TypeString, Computed: true},
			"status":                          {Type: schema.TypeString, Computed: true},
			"polar_fs_version":                {Type: schema.TypeString, Computed: true},
			"description":                     {Type: schema.TypeString, Computed: true},
			"security_group_id":               {Type: schema.TypeString, Computed: true},
			"create_time":                     {Type: schema.TypeString, Computed: true},
			"expire_time":                     {Type: schema.TypeString, Computed: true},
			"expired":                         {Type: schema.TypeString, Computed: true},
			"polar_fs_type":                   {Type: schema.TypeString, Computed: true},
			"storage_used":                    {Type: schema.TypeFloat, Computed: true},
			"bandwidth":                       {Type: schema.TypeFloat, Computed: true},
			"bandwidth_base_line":             {Type: schema.TypeFloat, Computed: true},
			"lock_mode":                       {Type: schema.TypeString, Computed: true},
			"accelerating_enable":             {Type: schema.TypeString, Computed: true},
			"accelerated_storage_space":       {Type: schema.TypeFloat, Computed: true},
			"minor_version":                   {Type: schema.TypeString, Computed: true},
			"client_download_path":            {Type: schema.TypeString, Computed: true},
			"relative_pfs_cluster_id":         {Type: schema.TypeString, Computed: true},
			"bucket_id":                       {Type: schema.TypeString, Computed: true},
			"file_system_id":                  {Type: schema.TypeString, Computed: true},
			"meta_instance_name":              {Type: schema.TypeString, Computed: true},
			"db_endpoint_id":                  {Type: schema.TypeString, Computed: true},
			"maxscale_endpoint_id":            {Type: schema.TypeString, Computed: true},
			"meta_connection_string":          {Type: schema.TypeString, Computed: true},
			"meta_maxscale_connection_string": {Type: schema.TypeString, Computed: true},
			"authorized_user_arn_ids":         {Type: schema.TypeString, Computed: true},
		},
	}
}

func resourceAlicloudPolarDBPolarFsCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	action := "CreatePolarFs"
	request := buildPolarDBPolarFsCreateRequest(d, client.RegionId)
	response, err := polarDBRpcPost(client, action, request, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_polardb_polar_fs", action, AlibabaCloudSdkGoERROR)
	}
	instanceId := fmt.Sprint(response["PolarFsInstanceId"])
	if instanceId == "" || instanceId == "<nil>" {
		return WrapError(fmt.Errorf("%s returned an empty PolarFsInstanceId", action))
	}
	if err = d.Set("polar_fs_path", response["PolarFsPath"]); err != nil {
		return WrapError(err)
	}
	d.SetId(instanceId)

	stateConf := BuildStateConf([]string{"Pending", "Creating"}, []string{"Running", "Activation"}, d.Timeout(schema.TimeoutCreate), 15*time.Second, polarDBPolarFsStateRefreshFunc(client, d.Id(), d.Get("db_cluster_id").(string)))
	if _, err = stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}
	return resourceAlicloudPolarDBPolarFsRead(d, meta)
}

func buildPolarDBPolarFsCreateRequest(d *schema.ResourceData, regionId string) map[string]interface{} {
	request := map[string]interface{}{"RegionId": regionId}
	fields := map[string]string{
		"storage_type": "StorageType", "vpc_id": "VPCId", "vswitch_id": "VSwitchId",
		"zone_id": "ZoneId", "storage_space": "StorageSpace", "db_cluster_id": "DBClusterId", "pay_type": "PayType",
		"accelerate_switch":       "AccelerateSwitch",
		"accelerate_storage_size": "AccelerateStorageSize", "creation_category": "CreationCategory",
		"accelerate_type": "AccelerateType",
	}
	for field, apiField := range fields {
		if value, ok := d.GetOk(field); ok {
			request[apiField] = value
		}
	}
	if value, ok := d.GetOk("db_type"); ok {
		request["DBType"] = expandPolarDBPolarFsDBType(value.(string))
	}
	return request
}

// lintignore: R001
func resourceAlicloudPolarDBPolarFsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	object, err := describePolarDBPolarFsAttribute(client, d.Id(), d.Get("db_cluster_id").(string))
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	return flattenPolarDBPolarFsAttribute(d, object)
}

// lintignore: R001
func flattenPolarDBPolarFsAttribute(d *schema.ResourceData, object map[string]interface{}) error {
	values := map[string]string{
		"region_id": "RegionId", "storage_type": "StorageType", "db_type": "DBType", "vpc_id": "VPCId",
		"vswitch_id": "VSwitchId", "zone_id": "ZoneId", "storage_space": "StorageSpace",
		"db_cluster_id": "RelativeDbClusterId", "pay_type": "PayType", "status": "PolarFsStatus",
		"polar_fs_version": "PolarFsVersion", "description": "PolarFsInstanceDescription",
		"security_group_id": "SecurityGroupId", "create_time": "CreateTime", "expire_time": "ExpireTime",
		"expired": "Expired", "polar_fs_type": "PolarFsType", "storage_used": "StorageUsed",
		"bandwidth": "Bandwidth", "bandwidth_base_line": "BandwidthBaseLine", "creation_category": "Category",
		"lock_mode": "LockMode", "accelerating_enable": "AcceleratingEnable",
		"accelerated_storage_space": "AcceleratedStorageSpace", "minor_version": "MinorVersion",
		"client_download_path": "ClientDownloadPath", "relative_pfs_cluster_id": "RelativePfsClusterId",
		"bucket_id": "BucketId", "file_system_id": "FileSystemId",
		"accelerate_type": "AccelerateType", "meta_instance_name": "MetaInstanceName",
		"db_endpoint_id": "DBEndpointId", "maxscale_endpoint_id": "MaxscaleEndpointId",
		"meta_connection_string": "MetaConnString", "meta_maxscale_connection_string": "MetaMxsConnString",
		"authorized_user_arn_ids": "AuthorizedUserArnIds",
	}
	for field, apiField := range values {
		if value, ok := object[apiField]; ok && value != nil {
			if field == "storage_space" {
				integer, err := normalizePolarDBPolarFsInteger(value)
				if err != nil {
					return WrapErrorf(err, "failed to decode %s from %v", apiField, object[apiField])
				}
				value = integer
			}
			if field == "db_type" {
				value = flattenPolarDBPolarFsDBType(fmt.Sprint(value))
			}
			if field == "zone_id" {
				value = strings.TrimSuffix(fmt.Sprint(value), "-aliyun")
			}
			if err := d.Set(field, value); err != nil {
				return WrapError(err)
			}
		}
	}
	return nil
}

func expandPolarDBPolarFsDBType(value string) string {
	switch value {
	case "MySQL":
		return "polardb_mysql"
	case "PostgreSQL":
		return "polardb_pg"
	default:
		return value
	}
}

func flattenPolarDBPolarFsDBType(value string) string {
	switch strings.ToLower(value) {
	case "polardb_mysql", "mysql":
		return "MySQL"
	case "polardb_pg", "postgresql":
		return "PostgreSQL"
	default:
		return value
	}
}

func normalizePolarDBPolarFsInteger(value interface{}) (int, error) {
	text := strings.TrimSpace(fmt.Sprint(value))
	if text == "" || text == "<nil>" {
		return 0, nil
	}
	number, err := strconv.ParseFloat(text, 64)
	if err != nil {
		return 0, fmt.Errorf("%q is not a number: %w", text, err)
	}
	if math.IsNaN(number) || math.IsInf(number, 0) || math.Trunc(number) != number {
		return 0, fmt.Errorf("%q is not an integer", text)
	}
	integer := int(number)
	if float64(integer) != number {
		return 0, fmt.Errorf("%q is outside the supported integer range", text)
	}
	return integer, nil
}

func resourceAlicloudPolarDBPolarFsDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	request := map[string]interface{}{"PolarFsInstanceId": d.Id()}
	if value, ok := d.GetOk("db_cluster_id"); ok {
		request["DBClusterId"] = value
	}
	if _, err := polarDBRpcPost(client, "DeletePolarFs", request, d.Timeout(schema.TimeoutDelete)); err != nil && !NotFoundError(err) {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), "DeletePolarFs", AlibabaCloudSdkGoERROR)
	}
	stateConf := BuildStateConf([]string{"Pending", "Running", "Creating", "Mounting", "Unmounting", "Mounted", "Unmounted", "Activation", "Deleting", "NetCreating", "NetDeleting"}, []string{}, d.Timeout(schema.TimeoutDelete), 10*time.Second, polarDBPolarFsStateRefreshFunc(client, d.Id(), d.Get("db_cluster_id").(string)))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, DeleteTimeoutMsg, d.Id(), "DeletePolarFs", ProviderERROR)
	}
	d.SetId("")
	return nil
}

func polarDBRpcPost(client *connectivity.AliyunClient, action string, request map[string]interface{}, timeout time.Duration) (map[string]interface{}, error) {
	var response map[string]interface{}
	var err error
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(timeout, func() *retry.RetryError {
		response, err = client.RpcPost("polardb", "2017-08-01", action, nil, request, false)
		if err != nil {
			if polarDBRpcPostRetryable(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, redactPolarDBSensitiveValue(response), redactPolarDBSensitiveValue(request))
	return response, err
}

func polarDBRpcPostRetryable(err error) bool {
	return IsExpectedErrors(err, []string{"LockTimeout"}) || NeedRetry(err)
}

func redactPolarDBSensitiveValue(value interface{}) interface{} {
	sensitiveKeys := map[string]struct{}{
		"CustomOssAk": {}, "CustomOssSk": {}, "Token": {}, "UserDefaultAccSk": {},
	}
	switch typed := value.(type) {
	case map[string]interface{}:
		redacted := make(map[string]interface{}, len(typed))
		for key, child := range typed {
			if _, sensitive := sensitiveKeys[key]; sensitive {
				redacted[key] = "***"
			} else {
				redacted[key] = redactPolarDBSensitiveValue(child)
			}
		}
		return redacted
	case []interface{}:
		redacted := make([]interface{}, len(typed))
		for index, child := range typed {
			redacted[index] = redactPolarDBSensitiveValue(child)
		}
		return redacted
	case []map[string]interface{}:
		redacted := make([]map[string]interface{}, len(typed))
		for index, child := range typed {
			redacted[index] = redactPolarDBSensitiveValue(child).(map[string]interface{})
		}
		return redacted
	default:
		return value
	}
}

func describePolarDBPolarFsAttribute(client *connectivity.AliyunClient, instanceId, dbClusterId string) (map[string]interface{}, error) {
	request := map[string]interface{}{"PolarFsInstanceId": instanceId}
	if dbClusterId != "" {
		request["DBClusterId"] = dbClusterId
	}
	response, err := polarDBRpcPost(client, "DescribePolarFsAttribute", request, 5*time.Minute)
	if err != nil {
		if IsExpectedErrors(err, []string{"InvalidPolarFsInstanceId.NotFound"}) {
			return nil, WrapErrorf(GetNotFoundErrorFromString(err.Error()), NotFoundMsg, AlibabaCloudSdkGoERROR)
		}
		return nil, err
	}
	if fmt.Sprint(response["PolarFsInstanceId"]) == "" || fmt.Sprint(response["PolarFsInstanceId"]) == "<nil>" {
		return nil, WrapErrorf(GetNotFoundErrorFromString(fmt.Sprintf("PolarFs %s was not found", instanceId)), NotFoundMsg, ProviderERROR)
	}
	return response, nil
}

func polarDBPolarFsStateRefreshFunc(client *connectivity.AliyunClient, instanceId, dbClusterId string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := describePolarDBPolarFsAttribute(client, instanceId, dbClusterId)
		if err != nil {
			if NotFoundError(err) {
				return nil, "Deleted", nil
			}
			return nil, "", err
		}
		return object, fmt.Sprint(object["PolarFsStatus"]), nil
	}
}
