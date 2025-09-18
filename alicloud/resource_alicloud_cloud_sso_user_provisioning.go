// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"github.com/PaesslerAG/jsonpath"
	"log"
	"strings"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAliCloudCloudSSOUserProvisioning() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudCloudSSOUserProvisioningCreate,
		Read:   resourceAliCloudCloudSSOUserProvisioningRead,
		Update: resourceAliCloudCloudSSOUserProvisioningUpdate,
		Delete: resourceAliCloudCloudSSOUserProvisioningDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"deletion_strategy": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"directory_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"duplication_strategy": {
				Type:     schema.TypeString,
				Required: true,
			},
			"principal_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"principal_type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"target_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"target_type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"user_provisioning_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"user_provisioning_statistics": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"failed_event_count": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"gmt_latest_sync": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func resourceAliCloudCloudSSOUserProvisioningCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "CreateUserProvisioning"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	if v, ok := d.GetOk("description"); ok {
		request["Description"] = v
	}
	request["DirectoryId"] = d.Get("directory_id")
	request["DuplicationStrategy"] = d.Get("duplication_strategy")
	request["TargetType"] = d.Get("target_type")
	request["TargetId"] = d.Get("target_id")
	request["PrincipalType"] = d.Get("principal_type")
	request["DeletionStrategy"] = d.Get("deletion_strategy")
	request["PrincipalId"] = d.Get("principal_id")
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_cloud_sso_user_provisioning", action, AlibabaCloudSdkGoERROR)
	}

	UserProvisioningDirectoryIdVar, _ := jsonpath.Get("$.UserProvisioning.DirectoryId", response)
	UserProvisioningUserProvisioningIdVar, _ := jsonpath.Get("$.UserProvisioning.UserProvisioningId", response)
	d.SetId(fmt.Sprintf("%v:%v", UserProvisioningDirectoryIdVar, UserProvisioningUserProvisioningIdVar))

	return resourceAliCloudCloudSSOUserProvisioningUpdate(d, meta)
}

func resourceAliCloudCloudSSOUserProvisioningRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	cloudSSOServiceV2 := CloudSSOServiceV2{client}

	objectRaw, err := cloudSSOServiceV2.DescribeCloudSSOUserProvisioning(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_cloud_sso_user_provisioning DescribeCloudSSOUserProvisioning Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("create_time", objectRaw["CreateTime"])
	d.Set("deletion_strategy", objectRaw["DeletionStrategy"])
	d.Set("description", objectRaw["Description"])
	d.Set("duplication_strategy", objectRaw["DuplicationStrategy"])
	d.Set("principal_id", objectRaw["PrincipalId"])
	d.Set("principal_type", objectRaw["PrincipalType"])
	d.Set("status", objectRaw["Status"])
	d.Set("target_id", objectRaw["TargetId"])
	d.Set("target_type", objectRaw["TargetType"])
	d.Set("directory_id", objectRaw["DirectoryId"])
	d.Set("user_provisioning_id", objectRaw["UserProvisioningId"])

	objectRaw, err = cloudSSOServiceV2.DescribeUserProvisioningGetUserProvisioningStatistics(d.Id())
	if err != nil && !NotFoundError(err) {
		return WrapError(err)
	}

	userProvisioningStatisticsMaps := make([]map[string]interface{}, 0)
	userProvisioningStatisticsMap := make(map[string]interface{})

	userProvisioningStatisticsMap["failed_event_count"] = objectRaw["FailedEventCount"]
	userProvisioningStatisticsMap["gmt_latest_sync"] = objectRaw["LatestAsyncTime"]

	userProvisioningStatisticsMaps = append(userProvisioningStatisticsMaps, userProvisioningStatisticsMap)
	if err := d.Set("user_provisioning_statistics", userProvisioningStatisticsMaps); err != nil {
		return err
	}

	return nil
}

func resourceAliCloudCloudSSOUserProvisioningUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	update := false

	var err error
	parts := strings.Split(d.Id(), ":")
	action := "UpdateUserProvisioning"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["UserProvisioningId"] = parts[1]
	request["DirectoryId"] = parts[0]

	if !d.IsNewResource() && d.HasChange("description") {
		update = true
		request["NewDescription"] = d.Get("description")
	}

	if !d.IsNewResource() && d.HasChange("duplication_strategy") {
		update = true
	}
	request["NewDuplicationStrategy"] = d.Get("duplication_strategy")
	if !d.IsNewResource() && d.HasChange("deletion_strategy") {
		update = true
	}
	request["NewDeletionStrategy"] = d.Get("deletion_strategy")
	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
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
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
	}

	return resourceAliCloudCloudSSOUserProvisioningRead(d, meta)
}

func resourceAliCloudCloudSSOUserProvisioningDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	parts := strings.Split(d.Id(), ":")
	action := "DeleteUserProvisioning"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	request["UserProvisioningId"] = parts[1]
	request["DirectoryId"] = parts[0]

	request["DeletionStrategy"] = d.Get("deletion_strategy")
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
		if NotFoundError(err) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	return nil
}
