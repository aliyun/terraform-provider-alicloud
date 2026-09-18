package alicloud

import (
	"fmt"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAlicloudPolarDBPolarFsEndpoint() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudPolarDBPolarFsEndpointCreate,
		Read:   resourceAlicloudPolarDBPolarFsEndpointRead,
		Delete: resourceAlicloudPolarDBPolarFsEndpointDelete,
		Importer: &schema.ResourceImporter{
			State: resourceAlicloudPolarDBPolarFsEndpointImport,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"db_cluster_id":        {Type: schema.TypeString, Required: true, ForceNew: true},
			"polar_fs_instance_id": {Type: schema.TypeString, Required: true, ForceNew: true},
			"endpoint_type": {
				Type: schema.TypeString, Required: true, ForceNew: true,
				ValidateFunc: StringInSlice([]string{"Nas", "S3Gateway", "S3Pvtz"}, false),
			},
			"vpc_id":                  {Type: schema.TypeString, Optional: true, ForceNew: true},
			"vswitch_id":              {Type: schema.TypeString, Optional: true, ForceNew: true},
			"db_endpoint_description": {Type: schema.TypeString, Optional: true, Computed: true, ForceNew: true},
			"db_endpoint_id":          {Type: schema.TypeString, Computed: true},
			"address_items": {
				Type: schema.TypeList, Computed: true,
				Elem: &schema.Resource{Schema: map[string]*schema.Schema{
					"connection_string":              {Type: schema.TypeString, Computed: true},
					"private_zone_connection_string": {Type: schema.TypeString, Computed: true},
					"ip_address":                     {Type: schema.TypeString, Computed: true},
					"port":                           {Type: schema.TypeString, Computed: true},
					"vpc_id":                         {Type: schema.TypeString, Computed: true},
					"vswitch_id":                     {Type: schema.TypeString, Computed: true},
					"net_type":                       {Type: schema.TypeString, Computed: true},
				}},
			},
		},
	}
}

func resourceAlicloudPolarDBPolarFsEndpointCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	oldEndpoints, err := describePolarDBPolarFsEndpoints(client, d.Get("db_cluster_id").(string), d.Get("polar_fs_instance_id").(string), "")
	if err != nil {
		return WrapError(err)
	}
	oldIds := make(map[string]struct{}, len(oldEndpoints))
	for _, endpoint := range oldEndpoints {
		oldIds[fmt.Sprint(endpoint["DBEndpointId"])] = struct{}{}
	}

	action := "CreateDBClusterEndpoint"
	request := map[string]interface{}{
		"DBClusterId":       d.Get("db_cluster_id"),
		"PolarFsInstanceId": d.Get("polar_fs_instance_id"),
		"EndpointType":      d.Get("endpoint_type"),
		"ClientToken":       buildClientToken(action),
	}
	for field, apiField := range map[string]string{
		"vpc_id": "VPCId", "vswitch_id": "VSwitchId", "db_endpoint_description": "DBEndpointDescription",
	} {
		if value, ok := d.GetOk(field); ok {
			request[apiField] = value
		}
	}
	if _, err = polarDBRpcPost(client, action, request, d.Timeout(schema.TimeoutCreate)); err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_polardb_polar_fs_endpoint", action, AlibabaCloudSdkGoERROR)
	}

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"Pending"},
		Target:     []string{"Available"},
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      10 * time.Second,
		MinTimeout: 5 * time.Second,
		Refresh: func() (interface{}, string, error) {
			endpoints, refreshErr := describePolarDBPolarFsEndpoints(client, d.Get("db_cluster_id").(string), d.Get("polar_fs_instance_id").(string), "")
			if refreshErr != nil {
				return nil, "", refreshErr
			}
			for _, endpoint := range endpoints {
				id := fmt.Sprint(endpoint["DBEndpointId"])
				if _, exists := oldIds[id]; !exists && fmt.Sprint(endpoint["EndpointType"]) == d.Get("endpoint_type").(string) {
					return endpoint, "Available", nil
				}
			}
			return nil, "Pending", nil
		},
	}
	result, err := stateConf.WaitForState()
	if err != nil {
		return WrapErrorf(err, IdMsg, d.Get("polar_fs_instance_id").(string))
	}
	endpoint := result.(map[string]interface{})
	endpointId := fmt.Sprint(endpoint["DBEndpointId"])
	d.SetId(fmt.Sprintf("%s%s%s", d.Get("db_cluster_id").(string), COLON_SEPARATED, endpointId))
	return resourceAlicloudPolarDBPolarFsEndpointRead(d, meta)
}

// lintignore: R001
func resourceAlicloudPolarDBPolarFsEndpointRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}
	endpoints, err := describePolarDBPolarFsEndpoints(client, parts[0], d.Get("polar_fs_instance_id").(string), parts[1])
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	if len(endpoints) == 0 {
		d.SetId("")
		return nil
	}
	object := endpoints[0]
	for field, apiField := range map[string]string{
		"db_endpoint_id": "DBEndpointId", "endpoint_type": "EndpointType",
		"db_endpoint_description": "DBEndpointDescription",
	} {
		if value, ok := object[apiField]; ok && value != nil {
			if err = d.Set(field, value); err != nil {
				return WrapError(err)
			}
		}
	}
	if err = d.Set("address_items", flattenPolarDBPolarFsEndpointAddresses(object["AddressItems"])); err != nil {
		return WrapError(err)
	}
	return nil
}

func resourceAlicloudPolarDBPolarFsEndpointDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}
	request := map[string]interface{}{
		"DBClusterId":       parts[0],
		"DBEndpointId":      parts[1],
		"PolarFsInstanceId": d.Get("polar_fs_instance_id"),
	}
	if _, err = polarDBRpcPost(client, "DeleteDBClusterEndpoint", request, d.Timeout(schema.TimeoutDelete)); err != nil && !NotFoundError(err) {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), "DeleteDBClusterEndpoint", AlibabaCloudSdkGoERROR)
	}
	stateConf := &resource.StateChangeConf{
		Pending: []string{"Exists"}, Target: []string{"Deleted"}, Timeout: d.Timeout(schema.TimeoutDelete), MinTimeout: 5 * time.Second,
		Refresh: func() (interface{}, string, error) {
			endpoints, refreshErr := describePolarDBPolarFsEndpoints(client, parts[0], d.Get("polar_fs_instance_id").(string), parts[1])
			if refreshErr != nil {
				if NotFoundError(refreshErr) {
					return nil, "Deleted", nil
				}
				return nil, "", refreshErr
			}
			if len(endpoints) == 0 {
				return nil, "Deleted", nil
			}
			return endpoints[0], "Exists", nil
		},
	}
	if _, err = stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, DeleteTimeoutMsg, d.Id(), "DeleteDBClusterEndpoint", ProviderERROR)
	}
	d.SetId("")
	return nil
}

func resourceAlicloudPolarDBPolarFsEndpointImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	parts, err := ParseResourceId(d.Id(), 3)
	if err != nil {
		return nil, WrapError(fmt.Errorf("the import ID must be '<db_cluster_id>:<polar_fs_instance_id>:<db_endpoint_id>': %w", err))
	}
	if err = d.Set("db_cluster_id", parts[0]); err != nil {
		return nil, WrapError(err)
	}
	if err = d.Set("polar_fs_instance_id", parts[1]); err != nil {
		return nil, WrapError(err)
	}
	d.SetId(fmt.Sprintf("%s%s%s", parts[0], COLON_SEPARATED, parts[2]))
	return []*schema.ResourceData{d}, nil
}

func describePolarDBPolarFsEndpoints(client *connectivity.AliyunClient, dbClusterId, polarFsInstanceId, endpointId string) ([]map[string]interface{}, error) {
	request := map[string]interface{}{"DBClusterId": dbClusterId, "PolarFsInstanceId": polarFsInstanceId}
	if endpointId != "" {
		request["DBEndpointId"] = endpointId
	}
	response, err := polarDBRpcPost(client, "DescribeDBClusterEndpoints", request, 5*time.Minute)
	if err != nil {
		return nil, err
	}
	raw, ok := response["Items"].([]interface{})
	if !ok {
		return nil, nil
	}
	result := make([]map[string]interface{}, 0, len(raw))
	for _, value := range raw {
		if item, ok := value.(map[string]interface{}); ok {
			if endpointId != "" && fmt.Sprint(item["DBEndpointId"]) != endpointId {
				continue
			}
			result = append(result, item)
		}
	}
	return result, nil
}

func flattenPolarDBPolarFsEndpointAddresses(raw interface{}) []map[string]interface{} {
	values, ok := raw.([]interface{})
	if !ok {
		return nil
	}
	result := make([]map[string]interface{}, 0, len(values))
	for _, value := range values {
		if item, ok := value.(map[string]interface{}); ok {
			result = append(result, map[string]interface{}{
				"connection_string": item["ConnectionString"], "private_zone_connection_string": item["PrivateZoneConnectionString"],
				"ip_address": item["IPAddress"], "port": item["Port"], "vpc_id": item["VPCId"],
				"vswitch_id": item["VSwitchId"], "net_type": item["NetType"],
			})
		}
	}
	return result
}
