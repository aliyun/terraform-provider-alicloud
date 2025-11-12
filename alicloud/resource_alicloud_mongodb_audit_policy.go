// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"log"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAliCloudMongodbAuditPolicy() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudMongodbAuditPolicyCreate,
		Read:   resourceAliCloudMongodbAuditPolicyRead,
		Update: resourceAliCloudMongodbAuditPolicyUpdate,
		Delete: resourceAliCloudMongodbAuditPolicyDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"audit_status": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: StringInSlice([]string{"enable", "disabled"}, false),
			},
			"db_instance_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"storage_period": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
		},
	}
}

func resourceAliCloudMongodbAuditPolicyCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "ModifyAuditPolicy"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	if v, ok := d.GetOk("db_instance_id"); ok {
		request["DBInstanceId"] = v
	}
	request["RegionId"] = client.RegionId

	if v, ok := d.GetOkExists("storage_period"); ok {
		request["StoragePeriod"] = v
	}
	request["ServiceType"] = "Standard"
	request["AuditStatus"] = convertMongodbAuditPolicyAuditStatusRequest(d.Get("audit_status").(string))
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.RpcPost("Dds", "2015-12-01", action, query, request, true)
		if err != nil {
			if IsExpectedErrors(err, []string{"OperationDenied.DBInstanceStatus"}) || NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_mongodb_audit_policy", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(request["DBInstanceId"]))

	mongodbServiceV2 := MongodbServiceV2{client}
	stateConf := BuildStateConf([]string{}, []string{"[Running]"}, d.Timeout(schema.TimeoutCreate), 5*time.Minute, mongodbServiceV2.DescribeAsyncMongodbAuditPolicyStateRefreshFunc(d, response, "$.DBInstances.DBInstance[*].DBInstanceStatus", []string{}))
	if jobDetail, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id(), jobDetail)
	}

	return resourceAliCloudMongodbAuditPolicyRead(d, meta)
}

func resourceAliCloudMongodbAuditPolicyRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	mongodbServiceV2 := MongodbServiceV2{client}

	objectRaw, err := mongodbServiceV2.DescribeMongodbAuditPolicy(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_mongodb_audit_policy DescribeMongodbAuditPolicy Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("storage_period", objectRaw["TtlForStandard"])

	objectRaw, err = mongodbServiceV2.DescribeAuditPolicyDescribeAuditPolicy(d.Id())
	if err != nil && !NotFoundError(err) {
		return WrapError(err)
	}

	d.Set("audit_status", convertMongodbAuditPolicyLogAuditStatusResponse(objectRaw["LogAuditStatus"]))

	d.Set("db_instance_id", d.Id())

	return nil
}

func resourceAliCloudMongodbAuditPolicyUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	update := false

	var err error
	action := "ModifyAuditPolicy"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["DBInstanceId"] = d.Id()
	request["RegionId"] = client.RegionId
	if d.HasChange("storage_period") {
		update = true
	}
	if v, ok := d.GetOk("storage_period"); ok {
		request["StoragePeriod"] = v
	}
	if d.HasChange("audit_status") {
		update = true
	}
	request["AuditStatus"] = convertMongodbAuditPolicyAuditStatusRequest(d.Get("audit_status").(string))
	request["ServiceType"] = "Standard"
	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("Dds", "2015-12-01", action, query, request, true)
			if err != nil {
				if IsExpectedErrors(err, []string{"OperationDenied.DBInstanceStatus"}) || NeedRetry(err) {
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
		mongodbServiceV2 := MongodbServiceV2{client}
		stateConf := BuildStateConf([]string{}, []string{"[Running]"}, d.Timeout(schema.TimeoutUpdate), 5*time.Minute, mongodbServiceV2.DescribeAsyncMongodbAuditPolicyStateRefreshFunc(d, response, "$.DBInstances.DBInstance[*].DBInstanceStatus", []string{}))
		if jobDetail, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id(), jobDetail)
		}
	}

	return resourceAliCloudMongodbAuditPolicyRead(d, meta)
}

func resourceAliCloudMongodbAuditPolicyDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[WARN] Cannot destroy resource AliCloud Resource Audit Policy. Terraform will remove this resource from the state file, however resources may remain.")
	return nil
}

func convertMongodbAuditPolicyLogAuditStatusResponse(source interface{}) interface{} {
	source = fmt.Sprint(source)
	switch source {
	case "Enable":
		return "enable"
	case "Disabled":
		return "disabled"
	}
	return source
}
func convertMongodbAuditPolicyAuditStatusRequest(source interface{}) interface{} {
	source = fmt.Sprint(source)
	switch source {
	case "enable":
		return "Enable"
	case " disabled":
		return "Disabled"
	}
	return source
}

func convertMongodbAuditPolicyResponse(source string) string {
	source = fmt.Sprint(source)
	switch source {
	case "Enable":
		return "enable"
	case "Disabled":
		return "disabled"
	}
	return source
}
