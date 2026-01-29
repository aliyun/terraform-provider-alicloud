// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAliCloudAlikafkaSaslAcl() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudAlikafkaSaslAclCreate,
		Read:   resourceAliCloudAlikafkaSaslAclRead,
		Update: resourceAliCloudAlikafkaSaslAclUpdate,
		Delete: resourceAliCloudAlikafkaSaslAclDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"acl_operation_type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: StringInSlice([]string{"Write", "Read", "Describe", "IdempotentWrite", "IDEMPOTENT_WRITE", "DESCRIBE_CONFIGS"}, false),
			},
			"acl_operation_types": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"acl_permission_type": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"acl_resource_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"acl_resource_pattern_type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"acl_resource_type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: StringInSlice([]string{"Topic", "Group", "Cluster", "TransactionalId"}, false),
			},
			"host": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"username": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceAliCloudAlikafkaSaslAclCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "CreateAcl"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	if v, ok := d.GetOk("acl_resource_name"); ok {
		request["AclResourceName"] = v
	}
	if v, ok := d.GetOk("username"); ok {
		request["Username"] = v
	}
	if v, ok := d.GetOk("instance_id"); ok {
		request["InstanceId"] = v
	}
	if v, ok := d.GetOk("acl_resource_type"); ok {
		request["AclResourceType"] = v
	}
	if v, ok := d.GetOk("acl_resource_pattern_type"); ok {
		request["AclResourcePatternType"] = v
	}
	if v, ok := d.GetOk("acl_operation_type"); ok {
		request["AclOperationType"] = v
	}
	request["RegionId"] = client.RegionId

	if v, ok := d.GetOk("host"); ok {
		request["Host"] = v
	}
	if v, ok := d.GetOk("acl_permission_type"); ok {
		request["AclPermissionType"] = v
	}
	if v, ok := d.GetOk("acl_operation_types"); ok {
		request["AclOperationTypes"] = v
	}
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.RpcPost("alikafka", "2019-09-16", action, query, request, true)
		if err != nil {
			if IsExpectedErrors(err, []string{"ONS_SYSTEM_FLOW_CONTROL"}) || NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_alikafka_sasl_acl", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprintf("%v:%v:%v:%v:%v:%v", request["InstanceId"], request["Username"], request["AclResourceType"], request["AclResourceName"], request["AclResourcePatternType"], request["AclOperationType"]))

	return resourceAliCloudAlikafkaSaslAclRead(d, meta)
}

func resourceAliCloudAlikafkaSaslAclRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	alikafkaServiceV2 := AlikafkaServiceV2{client}

	objectRaw, err := alikafkaServiceV2.DescribeAlikafkaSaslAcl(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_alikafka_sasl_acl DescribeAlikafkaSaslAcl Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("acl_permission_type", objectRaw["AclPermissionType"])
	d.Set("host", objectRaw["Host"])
	d.Set("acl_operation_type", convertAlikafkaSaslAclKafkaAclListKafkaAclVOAclOperationTypeResponse(objectRaw["AclOperationType"]))
	d.Set("acl_resource_name", objectRaw["AclResourceName"])
	d.Set("acl_resource_pattern_type", objectRaw["AclResourcePatternType"])
	d.Set("acl_resource_type", convertAlikafkaSaslAclKafkaAclListKafkaAclVOAclResourceTypeResponse(objectRaw["AclResourceType"]))
	d.Set("username", objectRaw["Username"])

	parts := strings.Split(d.Id(), ":")
	d.Set("instance_id", parts[0])

	return nil
}

func resourceAliCloudAlikafkaSaslAclUpdate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[INFO] Cannot update resource Alicloud Resource Sasl Acl.")
	return nil
}

func resourceAliCloudAlikafkaSaslAclDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	parts := strings.Split(d.Id(), ":")
	action := "DeleteAcl"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	request["AclResourceName"] = parts[3]
	request["Username"] = parts[1]
	request["InstanceId"] = parts[0]
	request["AclResourceType"] = parts[2]
	request["AclResourcePatternType"] = parts[4]
	request["AclOperationType"] = parts[5]
	request["RegionId"] = client.RegionId

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = client.RpcPost("alikafka", "2019-09-16", action, query, request, true)
		if err != nil {
			if IsExpectedErrors(err, []string{"ONS_SYSTEM_FLOW_CONTROL"}) || NeedRetry(err) {
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

func convertAlikafkaSaslAclKafkaAclListKafkaAclVOAclOperationTypeResponse(source interface{}) interface{} {
	source = fmt.Sprint(source)
	switch source {
	case "WRITE":
		return "Write"
	case "READ":
		return "Read"
	case "DESCRIBE":
		return "Describe"
	}
	return source
}

func convertAlikafkaSaslAclKafkaAclListKafkaAclVOAclResourceTypeResponse(source interface{}) interface{} {
	source = fmt.Sprint(source)
	switch source {
	case "TOPIC":
		return "Topic"
	case "GROUP":
		return "Group"
	case "CLUSTER":
		return "Cluster"
	case "TRANSACTIONAL_ID":
		return "TransactionalId"
	}
	return source
}
