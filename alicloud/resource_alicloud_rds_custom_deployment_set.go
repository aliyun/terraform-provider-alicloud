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

func resourceAliCloudRdsCustomDeploymentSet() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudRdsCustomDeploymentSetCreate,
		Read:   resourceAliCloudRdsCustomDeploymentSetRead,
		Update: resourceAliCloudRdsCustomDeploymentSetUpdate,
		Delete: resourceAliCloudRdsCustomDeploymentSetDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"custom_deployment_set_name": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"group_count": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: IntInSlice([]int{0, 3}),
			},
			"on_unable_to_redeploy_failed_instance": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: StringInSlice([]string{"CancelMembershipAndStart"}, false),
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"strategy": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				ValidateFunc: StringInSlice([]string{"Availability"}, false),
			},
		},
	}
}

func resourceAliCloudRdsCustomDeploymentSetCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "CreateRCDeploymentSet"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	request["RegionId"] = client.RegionId
	request["ClientToken"] = buildClientToken(action)

	if v, ok := d.GetOk("custom_deployment_set_name"); ok {
		request["DeploymentSetName"] = v
	}
	if v, ok := d.GetOk("on_unable_to_redeploy_failed_instance"); ok {
		request["OnUnableToRedeployFailedInstance"] = v
	}
	if v, ok := d.GetOk("description"); ok {
		request["Description"] = v
	}
	if v, ok := d.GetOk("strategy"); ok {
		request["Strategy"] = v
	}
	if v, ok := d.GetOkExists("group_count"); ok {
		request["GroupCount"] = v
	}
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.RpcPost("Rds", "2014-08-15", action, query, request, true)
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_rds_custom_deployment_set", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(response["DeploymentSetId"]))

	return resourceAliCloudRdsCustomDeploymentSetRead(d, meta)
}

func resourceAliCloudRdsCustomDeploymentSetRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	rdsServiceV2 := RdsServiceV2{client}

	objectRaw, err := rdsServiceV2.DescribeRdsCustomDeploymentSet(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_rds_custom_deployment_set DescribeRdsCustomDeploymentSet Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	if objectRaw["DeploymentSetName"] != nil {
		d.Set("custom_deployment_set_name", objectRaw["DeploymentSetName"])
	}
	if objectRaw["DeploymentSetDescription"] != nil {
		d.Set("description", objectRaw["DeploymentSetDescription"])
	}
	if objectRaw["Domain"] != nil {
		d.Set("status", objectRaw["Domain"])
	}
	if objectRaw["Strategy"] != nil {
		d.Set("strategy", objectRaw["Strategy"])
	}

	return nil
}

func resourceAliCloudRdsCustomDeploymentSetUpdate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[INFO] Cannot update resource Alicloud Resource Custom Deployment Set.")
	return nil
}

func resourceAliCloudRdsCustomDeploymentSetDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	action := "DeleteRCDeploymentSet"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	request["DeploymentSetId"] = d.Id()
	request["RegionId"] = client.RegionId

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = client.RpcPost("Rds", "2014-08-15", action, query, request, false)

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
