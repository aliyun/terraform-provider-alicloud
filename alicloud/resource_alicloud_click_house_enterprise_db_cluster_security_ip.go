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

func resourceAliCloudClickHouseEnterpriseDbClusterSecurityIP() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudClickHouseEnterpriseDbClusterSecurityIPCreate,
		Read:   resourceAliCloudClickHouseEnterpriseDbClusterSecurityIPRead,
		Update: resourceAliCloudClickHouseEnterpriseDbClusterSecurityIPUpdate,
		Delete: resourceAliCloudClickHouseEnterpriseDbClusterSecurityIPDelete,
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
			"group_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"security_ip_list": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func resourceAliCloudClickHouseEnterpriseDbClusterSecurityIPCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "ModifySecurityIPList"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	if v, ok := d.GetOk("db_instance_id"); ok {
		request["DBInstanceId"] = v
	}
	if v, ok := d.GetOk("group_name"); ok {
		request["GroupName"] = v
	}
	request["RegionId"] = client.RegionId

	request["ModifyMode"] = "0"
	request["SecurityIPList"] = d.Get("security_ip_list")
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.RpcPost("clickhouse", "2023-05-22", action, query, request, true)
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_click_house_enterprise_db_cluster_security_ip", action, AlibabaCloudSdkGoERROR)
	}

	DataDBInstanceNameVar, _ := jsonpath.Get("$.Data.DBInstanceName", response)
	DataGroupNameVar, _ := jsonpath.Get("$.Data.GroupName", response)
	d.SetId(fmt.Sprintf("%v:%v", DataDBInstanceNameVar, DataGroupNameVar))

	return resourceAliCloudClickHouseEnterpriseDbClusterSecurityIPRead(d, meta)
}

func resourceAliCloudClickHouseEnterpriseDbClusterSecurityIPRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	clickHouseServiceV2 := ClickHouseServiceV2{client}

	objectRaw, err := clickHouseServiceV2.DescribeClickHouseEnterpriseDbClusterSecurityIP(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_click_house_enterprise_db_cluster_security_ip DescribeClickHouseEnterpriseDbClusterSecurityIP Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("security_ip_list", objectRaw["SecurityIPList"])
	d.Set("group_name", objectRaw["GroupName"])

	parts := strings.Split(d.Id(), ":")
	d.Set("db_instance_id", parts[0])

	return nil
}

func resourceAliCloudClickHouseEnterpriseDbClusterSecurityIPUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	update := false

	var err error
	parts := strings.Split(d.Id(), ":")
	action := "ModifySecurityIPList"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["DBInstanceId"] = parts[0]
	request["GroupName"] = parts[1]
	request["RegionId"] = client.RegionId
	request["ModifyMode"] = "0"
	if d.HasChange("security_ip_list") {
		update = true
	}
	request["SecurityIPList"] = d.Get("security_ip_list")
	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("clickhouse", "2023-05-22", action, query, request, true)
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
	}

	return resourceAliCloudClickHouseEnterpriseDbClusterSecurityIPRead(d, meta)
}

func resourceAliCloudClickHouseEnterpriseDbClusterSecurityIPDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	parts := strings.Split(d.Id(), ":")
	action := "ModifySecurityIPList"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	request["DBInstanceId"] = parts[0]
	request["GroupName"] = parts[1]
	request["RegionId"] = client.RegionId

	request["ModifyMode"] = "2"
	request["SecurityIPList"] = d.Get("security_ip_list")
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = client.RpcPost("clickhouse", "2023-05-22", action, query, request, true)

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

	return nil
}
