package alicloud

import (
	"fmt"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAlicloudArmsIntegrationExporter() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudArmsIntegrationExporterCreate,
		Read:   resourceAlicloudArmsIntegrationExporterRead,
		Update: resourceAlicloudArmsIntegrationExporterUpdate,
		Delete: resourceAlicloudArmsIntegrationExporterDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(3 * time.Minute),
			Update: schema.DefaultTimeout(3 * time.Minute),
			Delete: schema.DefaultTimeout(3 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"cluster_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"integration_type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"param": {
				Type:     schema.TypeString,
				Required: true,
			},
			"instance_id": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func resourceAlicloudArmsIntegrationExporterCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "AddPrometheusIntegration"
	request := make(map[string]interface{})
	var err error

	request["RegionId"] = client.RegionId
	request["ClusterId"] = d.Get("cluster_id")
	request["IntegrationType"] = d.Get("integration_type")
	request["Param"] = d.Get("param")

	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutCreate)), func() *resource.RetryError {
		response, err = client.RpcPost("ARMS", "2019-08-08", action, nil, request, true)
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_arms_integration_exporter", action, AlibabaCloudSdkGoERROR)
	}

	if fmt.Sprint(response["Code"]) != "200" {
		return WrapErrorf(Error(GetCreateFailedMessage("Arms:IntegrationExporter")), NotFoundWithResponse, response)
	}

	if resp, err := jsonpath.Get("$.Data", response); err != nil || resp == nil {
		return WrapErrorf(err, IdMsg, "alicloud_arms_integration_exporter")
	} else {
		instanceId := resp.(map[string]interface{})["InstanceId"]
		d.SetId(fmt.Sprintf("%v:%v:%v", request["ClusterId"], request["IntegrationType"], instanceId))
	}

	return resourceAlicloudArmsIntegrationExporterRead(d, meta)
}

func resourceAlicloudArmsIntegrationExporterRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	armsService := ArmsService{client}

	object, err := armsService.DescribeArmsIntegrationExporter(d.Id())
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("cluster_id", object["ClusterId"])
	d.Set("integration_type", object["IntegrationType"])
	d.Set("param", object["Param"])
	d.Set("instance_id", formatInt(object["InstanceId"]))

	return nil
}

func resourceAlicloudArmsIntegrationExporterUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	update := false

	parts, err := ParseResourceId(d.Id(), 3)
	if err != nil {
		return WrapError(err)
	}

	request := map[string]interface{}{
		"RegionId":        client.RegionId,
		"ClusterId":       parts[0],
		"IntegrationType": parts[1],
		"InstanceId":      parts[2],
	}

	if d.HasChange("param") {
		update = true
	}
	request["Param"] = d.Get("param")

	if update {
		action := "UpdatePrometheusIntegration"
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *resource.RetryError {
			response, err = client.RpcPost("ARMS", "2019-08-08", action, nil, request, true)
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

	return resourceAlicloudArmsIntegrationExporterRead(d, meta)
}

func resourceAlicloudArmsIntegrationExporterDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	action := "DeletePrometheusIntegration"
	var response map[string]interface{}

	var err error

	parts, err := ParseResourceId(d.Id(), 3)
	if err != nil {
		return WrapError(err)
	}

	request := map[string]interface{}{
		"RegionId":        client.RegionId,
		"ClusterId":       parts[0],
		"IntegrationType": parts[1],
		"InstanceId":      parts[2],
	}

	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutDelete)), func() *resource.RetryError {
		response, err = client.RpcPost("ARMS", "2019-08-08", action, nil, request, true)
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
