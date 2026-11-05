// Package mongodb. This file is generated automatically. Please do not modify it manually, thank you!
package mongodb

import (
	"fmt"
	"log"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAliCloudMongodbPrivateSrvNetworkAddress() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudMongodbPrivateSrvNetworkAddressCreate,
		Read:   resourceAliCloudMongodbPrivateSrvNetworkAddressRead,
		Delete: resourceAliCloudMongodbPrivateSrvNetworkAddressDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(6 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"db_instance_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"private_srv_connection_string_uri": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceAliCloudMongodbPrivateSrvNetworkAddressCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "AllocateDBInstanceSrvNetworkAddress"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	if v, ok := d.GetOk("db_instance_id"); ok {
		request["DBInstanceId"] = v
	}
	request["RegionId"] = client.RegionId

	request["SrvConnectionType"] = "vpc"
	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
		response, err = client.RpcPost("Dds", "2015-12-01", action, query, request, true)
		if err != nil {
			if errs.IsExpectedErrors(err, []string{"InvalidStatus.NotFound"}) || helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	logs.AddDebug(action, response, request)

	if err != nil {
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, "alicloud_mongodb_private_srv_network_address", action, errs.AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(request["DBInstanceId"]))

	mongodbServiceV2 := NewMongodbServiceV2(client)
	stateConf := helper.BuildStateConf([]string{}, []string{"#CHECKSET"}, d.Timeout(schema.TimeoutCreate), 30*time.Second, mongodbServiceV2.MongodbPrivateSrvNetworkAddressStateRefreshFunc(d.Id(), "#PrivateSrvConnectionStringUri", []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return errs.WrapErrorf(err, errs.IdMsg, d.Id())
	}

	return resourceAliCloudMongodbPrivateSrvNetworkAddressRead(d, meta)
}

func resourceAliCloudMongodbPrivateSrvNetworkAddressRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	mongodbServiceV2 := NewMongodbServiceV2(client)

	objectRaw, err := mongodbServiceV2.DescribeMongodbPrivateSrvNetworkAddress(d.Id())
	if err != nil {
		if !d.IsNewResource() && errs.NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_mongodb_private_srv_network_address DescribeMongodbPrivateSrvNetworkAddress Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return errs.WrapError(err)
	}

	d.Set("private_srv_connection_string_uri", objectRaw["PrivateSrvConnectionStringUri"])

	d.Set("db_instance_id", d.Id())

	return nil
}

func resourceAliCloudMongodbPrivateSrvNetworkAddressDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	action := "ReleaseNodePrivateNetworkAddress"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	request["DBInstanceId"] = d.Id()
	request["RegionId"] = client.RegionId

	request["NetworkType"] = "VPC"
	request["ConnectionType"] = "SRV"
	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutDelete), func() *retry.RetryError {
		response, err = client.RpcPost("Dds", "2015-12-01", action, query, request, true)
		if err != nil {
			if errs.IsExpectedErrors(err, []string{"OperationDenied.DBInstanceStatus"}) || helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	logs.AddDebug(action, response, request)

	if err != nil {
		if errs.NotFoundError(err) {
			return nil
		}
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
	}

	return nil
}
