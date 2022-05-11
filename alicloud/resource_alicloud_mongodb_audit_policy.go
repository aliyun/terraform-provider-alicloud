package alicloud

import (
	"fmt"
	"log"
	"time"

	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceAlicloudMongodbAuditPolicy() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudMongodbAuditPolicyCreate,
		Read:   resourceAlicloudMongodbAuditPolicyRead,
		Update: resourceAlicloudMongodbAuditPolicyUpdate,
		Delete: resourceAlicloudMongodbAuditPolicyDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"db_instance_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"audit_status": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"disabled", "enable"}, false),
			},
			"storage_period": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntBetween(1, 30),
			},
		},
	}
}

func resourceAlicloudMongodbAuditPolicyCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "ModifyAuditPolicy"
	request := make(map[string]interface{})
	conn, err := client.NewDdsClient()
	if err != nil {
		return WrapError(err)
	}
	request["DBInstanceId"] = d.Get("db_instance_id")
	request["AuditStatus"] = d.Get("audit_status")
	request["ServiceType"] = "Standard"
	if v, ok := d.GetOk("storage_period"); ok {
		request["StoragePeriod"] = v
	}
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2015-12-01"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
		if err != nil {
			if NeedRetry(err) || IsExpectedErrors(err, []string{"OperationDenied.DBInstanceStatus"}) {
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
	ddsService := MongoDBService{client}
	stateConf := BuildStateConf([]string{"AuditLogOpening"}, []string{"Running"}, d.Timeout(schema.TimeoutCreate), 1*time.Minute, ddsService.RdsMongodbDBInstanceStateRefreshFunc(d.Id(), []string{"AuditLogClosing"}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapError(err)
	}

	return resourceAlicloudMongodbAuditPolicyRead(d, meta)
}
func resourceAlicloudMongodbAuditPolicyRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	MongoDBService := MongoDBService{client}
	object, err := MongoDBService.DescribeMongodbAuditPolicy(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_mongodb_audit_policy MongoDBService.DescribeMongodbAuditPolicy Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("db_instance_id", d.Id())
	d.Set("audit_status", convertMongodbAuditPolicyResponse(object["LogAuditStatus"].(string)))
	return nil
}
func resourceAlicloudMongodbAuditPolicyUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	conn, err := client.NewDdsClient()
	if err != nil {
		return WrapError(err)
	}
	var response map[string]interface{}
	update := false
	request := map[string]interface{}{
		"DBInstanceId": d.Id(),
	}
	request["AuditStatus"] = d.Get("audit_status")
	if d.HasChange("audit_status") {
		update = true
	}
	if v, ok := d.GetOk("storage_period"); ok {
		request["StoragePeriod"] = v
		update = true
	}
	if update {
		action := "ModifyAuditPolicy"
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2015-12-01"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
			if err != nil {
				if NeedRetry(err) || IsExpectedErrors(err, []string{"OperationDenied.DBInstanceStatus"}) {
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
	d.SetId(fmt.Sprint(request["DBInstanceId"]))
	ddsService := MongoDBService{client}
	stateConf := BuildStateConf([]string{"AuditLogOpening", "AuditLogClosing"}, []string{"Running"}, d.Timeout(schema.TimeoutUpdate), 1*time.Minute, ddsService.RdsMongodbDBInstanceStateRefreshFunc(d.Id(), []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapError(err)
	}

	return resourceAlicloudMongodbAuditPolicyRead(d, meta)
}
func resourceAlicloudMongodbAuditPolicyDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[WARN] Cannot destroy resource Alicloud Mongodb AuditPolicy. Terraform will remove this resource from the state file, however resources may remain.")
	return nil
}

func convertMongodbAuditPolicyResponse(source string) string {
	switch source {
	case "Enable":
		return "enable"
	case "Disabled":
		return "disabled"
	}

	return source
}
