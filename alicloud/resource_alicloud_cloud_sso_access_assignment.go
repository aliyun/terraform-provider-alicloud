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

func resourceAliCloudCloudSSOAccessAssignment() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudCloudSSOAccessAssignmentCreate,
		Read:   resourceAliCloudCloudSSOAccessAssignmentRead,
		Update: resourceAliCloudCloudSSOAccessAssignmentUpdate,
		Delete: resourceAliCloudCloudSSOAccessAssignmentDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"access_configuration_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"directory_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"principal_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"principal_type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: StringInSlice([]string{"User", "Group"}, false),
			},
			"target_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"target_type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: StringInSlice([]string{"RD-Account"}, false),
			},
			"deprovision_strategy": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: StringInSlice([]string{"DeprovisionForLastAccessAssignmentOnAccount", "None"}, false),
			},
		},
	}
}

func resourceAliCloudCloudSSOAccessAssignmentCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "CreateAccessAssignment"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})

	request["AccessConfigurationId"] = d.Get("access_configuration_id")
	request["DirectoryId"] = d.Get("directory_id")
	request["PrincipalId"] = d.Get("principal_id")
	request["PrincipalType"] = d.Get("principal_type")
	request["TargetId"] = d.Get("target_id")
	request["TargetType"] = d.Get("target_type")

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.RpcPost("cloudsso", "2021-05-15", action, query, request, true)
		if err != nil {
			if IsExpectedErrors(err, []string{"OperationConflict.Task"}) || NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_cloud_sso_access_assignment", action, AlibabaCloudSdkGoERROR)
	}

	resp, err := jsonpath.Get("$.Task", response)
	if err != nil || resp == nil {
		return WrapErrorf(err, IdMsg, "alicloud_cloud_sso_access_assignment")
	}

	accessConfigurationId := resp.(map[string]interface{})["AccessConfigurationId"]
	targetType := resp.(map[string]interface{})["TargetType"]
	targetId := resp.(map[string]interface{})["TargetId"]
	principalType := resp.(map[string]interface{})["PrincipalType"]
	principalId := resp.(map[string]interface{})["PrincipalId"]

	d.SetId(fmt.Sprintf("%v:%v:%v:%v:%v:%v", request["DirectoryId"], accessConfigurationId, targetType, targetId, principalType, principalId))

	cloudSSOServiceV2 := CloudSSOServiceV2{client}
	stateConf := BuildStateConf([]string{}, []string{"Success"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, cloudSSOServiceV2.DescribeAsyncCloudSSOAccessAssignmentStateRefreshFunc(d, response, "TaskStatus.Status", []string{}))
	if jobDetail, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id(), jobDetail)
	}

	return resourceAliCloudCloudSSOAccessAssignmentRead(d, meta)
}

func resourceAliCloudCloudSSOAccessAssignmentRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	cloudSSOServiceV2 := CloudSSOServiceV2{client}

	objectRaw, err := cloudSSOServiceV2.DescribeCloudSSOAccessAssignment(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_cloud_sso_access_assignment DescribeCloudSSOAccessAssignment Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("create_time", objectRaw["CreateTime"])
	d.Set("access_configuration_id", objectRaw["AccessConfigurationId"])
	d.Set("principal_id", objectRaw["PrincipalId"])
	d.Set("principal_type", objectRaw["PrincipalType"])
	d.Set("target_id", objectRaw["TargetId"])
	d.Set("target_type", objectRaw["TargetType"])

	parts := strings.Split(d.Id(), ":")
	d.Set("directory_id", parts[0])

	return nil
}

func resourceAliCloudCloudSSOAccessAssignmentUpdate(d *schema.ResourceData, meta interface{}) error {
	log.Println(fmt.Sprintf("[WARNING] The resouce has not update operation."))
	return resourceAliCloudCloudSSOAccessAssignmentRead(d, meta)
}

func resourceAliCloudCloudSSOAccessAssignmentDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	parts := strings.Split(d.Id(), ":")
	action := "DeleteAccessAssignment"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	request["DirectoryId"] = parts[0]
	request["AccessConfigurationId"] = parts[1]
	request["TargetType"] = parts[2]
	request["TargetId"] = parts[3]
	request["PrincipalType"] = parts[4]
	request["PrincipalId"] = parts[5]

	if v, ok := d.GetOk("deprovision_strategy"); ok {
		request["DeprovisionStrategy"] = v
	} else {
		request["DeprovisionStrategy"] = "DeprovisionForLastAccessAssignmentOnAccount"
	}

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = client.RpcPost("cloudsso", "2021-05-15", action, query, request, true)

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
		if IsExpectedErrors(err, []string{"EntityNotExists.AccessAssignment"}) || NotFoundError(err) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	cloudSSOServiceV2 := CloudSSOServiceV2{client}
	stateConf := BuildStateConf([]string{}, []string{"Success"}, d.Timeout(schema.TimeoutDelete), 5*time.Second, cloudSSOServiceV2.DescribeAsyncCloudSSOAccessAssignmentStateRefreshFunc(d, response, "TaskStatus.Status", []string{}))
	if jobDetail, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id(), jobDetail)
	}

	return nil
}
