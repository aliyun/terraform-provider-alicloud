package alicloud

import (
	"fmt"
	"log"
	"regexp"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAliCloudEcsDeploymentSet() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudEcsDeploymentSetCreate,
		Read:   resourceAliCloudEcsDeploymentSetRead,
		Update: resourceAliCloudEcsDeploymentSetUpdate,
		Delete: resourceAliCloudEcsDeploymentSetDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"strategy": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Computed:     true,
				ValidateFunc: StringInSlice([]string{"Availability", "AvailabilityGroup", "LowLatency"}, false),
			},
			"deployment_set_name": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: StringMatch(regexp.MustCompile(`^([\w\\:\-]){2,128}$`), "\t\nThe name of the deployment set.\n\nThe name must be 2 to 128 characters in length and can contain letters, digits, colons (:), underscores (_), and hyphens (-)."),
			},
			"description": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: StringLenBetween(2, 256),
			},
			"on_unable_to_redeploy_failed_instance": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: StringInSlice([]string{"CancelMembershipAndStart", "KeepStopped"}, false),
			},
			"domain": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Computed:     true,
				ValidateFunc: StringInSlice([]string{"Default"}, false),
				Deprecated:   "Field `domain` has been deprecated from provider version 1.243.0.",
			},
			"granularity": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Computed:     true,
				ValidateFunc: StringInSlice([]string{"Host"}, false),
				Deprecated:   "Field `granularity` has been deprecated from provider version 1.243.0.",
			},
		},
	}
}

func resourceAliCloudEcsDeploymentSetCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "CreateDeploymentSet"
	request := make(map[string]interface{})
	var err error

	request["RegionId"] = client.RegionId
	request["ClientToken"] = buildClientToken("CreateDeploymentSet")

	if v, ok := d.GetOk("strategy"); ok {
		request["Strategy"] = v
	}

	if v, ok := d.GetOk("deployment_set_name"); ok {
		request["DeploymentSetName"] = v
	}

	if v, ok := d.GetOk("description"); ok {
		request["Description"] = v
	}

	if v, ok := d.GetOk("on_unable_to_redeploy_failed_instance"); ok {
		request["OnUnableToRedeployFailedInstance"] = v
	}

	if v, ok := d.GetOk("domain"); ok {
		request["Domain"] = v
	}

	if v, ok := d.GetOk("granularity"); ok {
		request["Granularity"] = v
	}

	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutCreate)), func() *resource.RetryError {
		response, err = client.RpcPost("Ecs", "2014-05-26", action, nil, request, true)
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_ecs_deployment_set", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(response["DeploymentSetId"]))

	return resourceAliCloudEcsDeploymentSetRead(d, meta)
}

func resourceAliCloudEcsDeploymentSetRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	ecsService := EcsService{client}

	object, err := ecsService.DescribeEcsDeploymentSet(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_ecs_deployment_set ecsService.DescribeEcsDeploymentSet Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("strategy", object["DeploymentStrategy"])
	d.Set("deployment_set_name", object["DeploymentSetName"])
	d.Set("description", object["DeploymentSetDescription"])
	d.Set("domain", convertEcsDeploymentSetDomainResponse(object["Domain"]))
	d.Set("granularity", convertEcsDeploymentSetGranularityResponse(object["Granularity"]))

	return nil
}

func resourceAliCloudEcsDeploymentSetUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	var err error
	update := false

	request := map[string]interface{}{
		"RegionId":        client.RegionId,
		"DeploymentSetId": d.Id(),
	}

	if d.HasChange("deployment_set_name") {
		update = true
	}
	if v, ok := d.GetOk("deployment_set_name"); ok {
		request["DeploymentSetName"] = v
	}

	if d.HasChange("description") {
		update = true
	}
	if v, ok := d.GetOk("description"); ok {
		request["Description"] = v
	}

	if update {
		action := "ModifyDeploymentSetAttribute"
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *resource.RetryError {
			response, err = client.RpcPost("Ecs", "2014-05-26", action, nil, request, false)
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

	return resourceAliCloudEcsDeploymentSetRead(d, meta)
}

func resourceAliCloudEcsDeploymentSetDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	action := "DeleteDeploymentSet"
	var response map[string]interface{}
	var err error

	request := map[string]interface{}{
		"RegionId":        client.RegionId,
		"DeploymentSetId": d.Id(),
	}

	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutDelete)), func() *resource.RetryError {
		response, err = client.RpcPost("Ecs", "2014-05-26", action, nil, request, false)
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

	return nil
}

func convertEcsDeploymentSetDomainResponse(source interface{}) interface{} {
	switch source {
	case "default":
		return "Default"
	}

	return source
}

func convertEcsDeploymentSetGranularityResponse(source interface{}) interface{} {
	switch source {
	case "host":
		return "Host"
	}

	return source
}
