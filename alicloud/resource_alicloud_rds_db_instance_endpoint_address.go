package alicloud

import (
	"fmt"
	"log"
	"time"

	"github.com/PaesslerAG/jsonpath"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAlicloudRdsDBInstanceEndpointAddress() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudRdsDBInstanceEndpointAddressCreate,
		Read:   resourceAlicloudRdsDBInstanceEndpointAddressRead,
		Update: resourceAlicloudRdsDBInstanceEndpointAddressUpdate,
		Delete: resourceAlicloudRdsDBInstanceEndpointAddressDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(20 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"db_instance_id": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},
			"db_instance_endpoint_id": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},
			"ip_address": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"connection_string_prefix": {
				Type:     schema.TypeString,
				Required: true,
			},
			"connection_string": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"port": {
				Type:     schema.TypeString,
				Required: true,
			},
			"ip_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceAlicloudRdsDBInstanceEndpointAddressCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	rdsService := RdsService{client}
	var response map[string]interface{}
	action := "CreateDBInstanceEndpointAddress"
	request := map[string]interface{}{
		"RegionId":               client.RegionId,
		"DBInstanceId":           Trim(d.Get("db_instance_id").(string)),
		"DBInstanceEndpointId":   d.Get("db_instance_endpoint_id"),
		"ConnectionStringPrefix": d.Get("connection_string_prefix"),
		"Port":                   d.Get("port"),
		"IpType":                 "Public",
		"SourceIp":               client.SourceIp,
	}

	var err error
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutCreate)), func() *resource.RetryError {
		response, err = client.RpcPost("Rds", "2014-08-15", action, nil, request, false)
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_rds_db_instance_endpoint_address", action, AlibabaCloudSdkGoERROR)
	}
	if v, err := jsonpath.Get("$.Data", response); err != nil || v == nil {
		return WrapErrorf(err, IdMsg, "alicloud_cloud_rds_db_instance_endpoint_address")
	} else {
		object := v.(map[string]interface{})
		d.SetId(fmt.Sprint(request["DBInstanceId"], ":", object["DBInstanceEndpointId"].(string)))
		d.Set("connection_string", object["ConnectionString"].(string))
	}
	stateConf := BuildStateConf([]string{}, []string{"Running"}, client.GetRetryTimeout(d.Timeout(schema.TimeoutCreate)), 1*time.Minute, rdsService.RdsDBInstanceStateRefreshFunc(request["DBInstanceId"].(string), []string{"Deleting"}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}
	return resourceAlicloudRdsDBInstanceEndpointAddressUpdate(d, meta)
}

func resourceAlicloudRdsDBInstanceEndpointAddressUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var err error
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}
	rdsService := RdsService{client}
	if !d.IsNewResource() && d.HasChanges("connection_string_prefix", "port") {
		action := "ModifyDBInstanceEndpointAddress"
		request := map[string]interface{}{
			"RegionId":             client.RegionId,
			"DBInstanceEndpointId": parts[1],
			"DBInstanceId":         parts[0],
			"ConnectionString":     d.Get("connection_string"),
			"SourceIp":             client.SourceIp,
		}
		if v, ok := d.GetOk("connection_string_prefix"); ok {
			request["ConnectionStringPrefix"] = v
		}
		if v, ok := d.GetOk("port"); ok {
			request["Port"] = v
		}
		var response map[string]interface{}
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *resource.RetryError {
			response, err = client.RpcPost("Rds", "2014-08-15", action, nil, request, false)
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
			return WrapErrorf(err, DefaultErrorMsg, "alicloud_rds_db_instance_endpoint_address", action, AlibabaCloudSdkGoERROR)
		}
		stateConf := BuildStateConf([]string{}, []string{"Running"}, client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), 3*time.Minute, rdsService.RdsDBInstanceStateRefreshFunc(parts[0], []string{"Deleting"}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
	}
	return resourceAlicloudRdsDBInstanceEndpointAddressRead(d, meta)
}

func resourceAlicloudRdsDBInstanceEndpointAddressRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	rdsService := RdsService{client}
	object, err := rdsService.DescribeDBInstanceEndpointPublicAddress(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_rds_db_instance_endpoint_address rdsService.DescribeDBInstanceEndpointPublicAddress Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	d.Set("ip_type", object["IpType"])
	d.Set("ip_address", object["IpAddress"])
	d.Set("port", object["Port"])
	d.Set("connection_string", object["ConnectionString"])
	d.Set("connection_string_prefix", object["ConnectionStringPrefix"])
	d.Set("db_instance_id", object["DBInstanceId"])
	d.Set("db_instance_endpoint_id", object["DBInstanceEndpointId"])
	return nil
}

func resourceAlicloudRdsDBInstanceEndpointAddressDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	rdsService := RdsService{client}
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}
	_, err = rdsService.DescribeDBInstance(parts[0])
	if err != nil {
		if NotFoundError(err) {
			return nil
		}
		return WrapError(err)
	}
	action := "DeleteDBInstanceEndpointAddress"
	request := map[string]interface{}{
		"RegionId":             client.RegionId,
		"DBInstanceId":         parts[0],
		"DBInstanceEndpointId": parts[1],
		"ConnectionString":     d.Get("connection_string"),
		"SourceIp":             client.SourceIp,
	}
	var response map[string]interface{}
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutDelete)), func() *resource.RetryError {
		response, err = client.RpcPost("Rds", "2014-08-15", action, nil, request, false)
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
	stateConf := BuildStateConf([]string{}, []string{"Running"}, client.GetRetryTimeout(d.Timeout(schema.TimeoutDelete)), 1*time.Minute, rdsService.RdsDBInstanceStateRefreshFunc(parts[0], []string{"Deleting"}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}
	return nil
}
