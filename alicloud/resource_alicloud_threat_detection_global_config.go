package alicloud

import (
	"fmt"
	"log"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAlicloudThreatDetectionGlobalConfig() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudThreatDetectionGlobalConfigCreate,
		Read:   resourceAlicloudThreatDetectionGlobalConfigRead,
		Update: resourceAlicloudThreatDetectionGlobalConfigUpdate,
		Delete: resourceAlicloudThreatDetectionGlobalConfigDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"global_config_name": {
				Required: true,
				ForceNew: true,
				Type:     schema.TypeString,
			},
			"global_config_value": {
				Optional: true,
				Type:     schema.TypeString,
			},
			"lang": {
				Optional: true,
				Type:     schema.TypeString,
			},
			"role_for": {
				Optional: true,
				Type:     schema.TypeInt,
			},
		},
	}
}

// cloudSiemEndpoint resolves the endpoint for the cloud-siem (Cloud SIEM) product.
// The resource is center scope with different endpoints for china and intl sites.
func cloudSiemEndpoint(client *connectivity.AliyunClient) string {
	if client.IsInternationalAccount() {
		return "cloud-siem.ap-southeast-1.aliyuncs.com"
	}
	return "cloud-siem.aliyuncs.com"
}

func resourceAlicloudThreatDetectionGlobalConfigCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	request := make(map[string]interface{})

	if v, ok := d.GetOk("global_config_name"); ok {
		request["GlobalConfigName"] = v
	}
	if v, ok := d.GetOk("global_config_value"); ok {
		request["GlobalConfigValue"] = v
	}
	if v, ok := d.GetOk("lang"); ok {
		request["Lang"] = v
	}
	if v, ok := d.GetOk("role_for"); ok {
		request["RoleFor"] = v
	}
	request["RegionId"] = client.RegionId

	endpoint := cloudSiemEndpoint(client)
	action := "UpdateGlobalConfig"
	wait := incrementalWait(3*time.Second, 3*time.Second)
	var err error
	err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutCreate)), func() *resource.RetryError {
		resp, err := client.RpcPostWithEndpoint("cloud-siem", "2024-12-12", action, nil, request, false, endpoint)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(action, resp, request)
		return nil
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_threat_detection_global_config", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(d.Get("global_config_name")))
	return resourceAlicloudThreatDetectionGlobalConfigRead(d, meta)
}

func resourceAlicloudThreatDetectionGlobalConfigRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	threatDetectionServiceV2 := ThreatDetectionServiceV2{client}

	object, err := threatDetectionServiceV2.DescribeThreatDetectionGlobalConfig(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_threat_detection_global_config DescribeThreatDetectionGlobalConfig Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	d.Set("global_config_name", object["GlobalConfigName"])
	d.Set("global_config_value", object["GlobalConfigValue"])

	return nil
}

func resourceAlicloudThreatDetectionGlobalConfigUpdate(d *schema.ResourceData, meta interface{}) error {
	return resourceAlicloudThreatDetectionGlobalConfigCreate(d, meta)
}

func resourceAlicloudThreatDetectionGlobalConfigDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[WARN] Cannot destroy resource AliCloud Threat Detection Global Config. Terraform will remove this resource from the state file, however resources may remain.")
	return nil
}

// DescribeThreatDetectionGlobalConfig queries the Threat Detection Global Config by its name.
func (s *ThreatDetectionServiceV2) DescribeThreatDetectionGlobalConfig(id string) (object map[string]interface{}, err error) {
	client := s.client
	request := make(map[string]interface{})
	request["GlobalConfigName"] = id
	request["RegionId"] = client.RegionId

	endpoint := cloudSiemEndpoint(client)
	action := "GetGlobalConfig"
	var response map[string]interface{}
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		resp, err := client.RpcPostWithEndpoint("cloud-siem", "2024-12-12", action, nil, request, true, endpoint)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		response = resp
		addDebug(action, response, request)
		return nil
	})
	if err != nil {
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}
	v, err := jsonpath.Get("$.GlobalConfig", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.GlobalConfig", response)
	}
	if v == nil {
		return object, WrapErrorf(NotFoundErr("ThreatDetection:GlobalConfig", id), NotFoundWithResponse, response)
	}
	object, ok := v.(map[string]interface{})
	if !ok {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.GlobalConfig", response)
	}
	return object, nil
}
