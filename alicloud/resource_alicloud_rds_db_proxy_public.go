package alicloud

import (
	"strings"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAliCloudRdsDBProxyPublic() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudRdsDBProxyPublicCreate,
		Read:   resourceAliCloudRdsDBProxyPublicRead,
		Update: resourceAliCloudRdsDBProxyPublicUpdate,
		Delete: resourceAliCloudRdsDBProxyPublicDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(60 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(20 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"db_instance_id": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},
			"db_proxy_endpoint_id": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},
			"connection_string_prefix": {
				Type:     schema.TypeString,
				Required: true,
			},
			"db_proxy_connection_string_net_type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: StringInSlice([]string{"Public"}, false),
			},
			"db_proxy_new_connect_string_port": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		},
	}
}

func resourceAliCloudRdsDBProxyPublicCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	rdsService := RdsService{client}
	action := "CreateDBProxyEndpointAddress"
	request := map[string]interface{}{
		"RegionId":                    client.RegionId,
		"DBInstanceId":                Trim(d.Get("db_instance_id").(string)),
		"DBProxyEndpointId":           Trim(d.Get("db_proxy_endpoint_id").(string)),
		"ConnectionStringPrefix":      Trim(d.Get("connection_string_prefix").(string)),
		"DBProxyConnectStringNetType": Trim(d.Get("db_proxy_connection_string_net_type").(string)),
	}
	dBProxyNewConnectStringPort, ok := d.GetOk("db_proxy_new_connect_string_port")
	if ok && dBProxyNewConnectStringPort.(string) != "" {
		request["DBProxyNewConnectStringPort"] = dBProxyNewConnectStringPort
	}
	if err := resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutCreate)), func() *resource.RetryError {
		response, err := client.RpcPost("Rds", "2014-08-15", action, nil, request, false)
		if err != nil {
			if NeedRetry(err) {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(action, response, request)
		return nil
	}); err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}
	d.SetId(request["DBInstanceId"].(string))
	stateConf := BuildStateConf([]string{}, []string{"Running"}, d.Timeout(schema.TimeoutCreate), 3*time.Minute, rdsService.RdsDBProxyStateRefreshFunc(request["DBInstanceId"].(string), []string{"Deleting"}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}
	return resourceAliCloudRdsDBProxyPublicRead(d, meta)
}

func resourceAliCloudRdsDBProxyPublicRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	rdsService := RdsService{client}
	proxy, proxyErr := rdsService.DescribeDBProxy(d.Id())
	if proxyErr != nil {
		if NotFoundError(proxyErr) {
			d.SetId("")
			return nil
		}
		return WrapError(proxyErr)
	}
	d.Set("db_instance_id", d.Id())
	if endpointId, ok := proxy["DBProxyInstanceName"].(string); ok {
		d.Set("db_proxy_endpoint_id", endpointId)
	}
	d.Set("db_proxy_connection_string_net_type", "Public")
	if connectString, ok := proxy["DBProxyConnectString"].(string); ok {
		parts := strings.Split(connectString, ".")
		if len(parts) > 0 {
			d.Set("connection_string_prefix", parts[0])
		}
	}
	if port, ok := proxy["DBProxyConnectStringPort"].(string); ok {
		d.Set("db_proxy_new_connect_string_port", port)
	}
	return nil
}

func resourceAliCloudRdsDBProxyPublicUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	rdsService := RdsService{client}
	if d.IsNewResource() {
		d.Partial(false)
		return resourceAliCloudRdsDBProxyPublicRead(d, meta)
	}

	if d.HasChanges("connection_string_prefix", "db_proxy_new_connect_string_port") {
		action := "ModifyDBProxyEndpointAddress"
		request := map[string]interface{}{
			"RegionId":                    client.RegionId,
			"DBInstanceId":                d.Id(),
			"DBProxyEndpointId":           d.Get("db_proxy_endpoint_id"),
			"DBProxyConnectStringNetType": "Public",
		}
		portAddressUpdate := false
		if v, ok := d.GetOk("connection_string_prefix"); ok {
			request["DBProxyNewConnectString"] = v
			portAddressUpdate = true
		}
		if v, ok := d.GetOk("db_proxy_new_connect_string_port"); ok {
			request["DBProxyNewConnectStringPort"] = v
			portAddressUpdate = true
		}
		if portAddressUpdate {
			if err := resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *resource.RetryError {
				response, err := client.RpcPost("Rds", "2014-08-15", action, nil, request, false)
				if err != nil {
					if NeedRetry(err) {
						return resource.RetryableError(err)
					}
					return resource.NonRetryableError(err)
				}
				addDebug(action, response, request)
				return nil
			}); err != nil {
				return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
			}
			stateConf := BuildStateConf([]string{}, []string{"Running"}, d.Timeout(schema.TimeoutUpdate), 3*time.Minute, rdsService.RdsDBProxyStateRefreshFunc(d.Id(), []string{"Deleting"}))
			if _, err := stateConf.WaitForState(); err != nil {
				return WrapErrorf(err, IdMsg, d.Id())
			}
		}
	}
	return resourceAliCloudRdsDBProxyPublicRead(d, meta)
}
func resourceAliCloudRdsDBProxyPublicDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	rdsService := RdsService{client}
	if d.Id() == "" {
		return nil
	}
	action := "DeleteDBProxyEndpointAddress"
	request := map[string]interface{}{
		"RegionId":                    client.RegionId,
		"DBInstanceId":                d.Id(),
		"DBProxyEndpointId":           d.Get("db_proxy_endpoint_id"),
		"DBProxyConnectStringNetType": "Public",
	}
	if err := resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutDelete)), func() *resource.RetryError {
		response, err := client.RpcPost("Rds", "2014-08-15", action, nil, request, false)
		if err != nil {
			if IsExpectedErrors(err, OperationDeniedDBStatus) || NeedRetry(err) {
				return resource.RetryableError(err)
			}
			if NotFoundError(err) {
				return nil
			}
			return resource.NonRetryableError(err)
		}
		addDebug(action, response, request)
		return nil
	}); err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}
	stateConf := BuildStateConf([]string{}, []string{"Running"}, d.Timeout(schema.TimeoutDelete), 3*time.Minute, rdsService.RdsDBProxyStateRefreshFunc(d.Id(), []string{"Deleted"}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}
	return nil
}
