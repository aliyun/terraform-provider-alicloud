// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"log"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAliCloudThreatDetectionServiceLinkedRole() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudThreatDetectionServiceLinkedRoleCreate,
		Read:   resourceAliCloudThreatDetectionServiceLinkedRoleRead,
		Delete: resourceAliCloudThreatDetectionServiceLinkedRoleDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"role_status": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"service_linked_role": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				ValidateFunc: StringInSlice([]string{"AliyunServiceRoleForSas", "AliyunServiceRoleForSasCspm", "AliyunServiceRoleForSasAgentless", "AliyunServiceRoleForOssMfd", "AliyunServiceRoleForAntiRansomwareMssp", "AliyunServiceRoleForSasSecllm", "AliyunServiceRoleForSasAI", "AliyunLogETLRole"}, false),
			},
		},
	}
}

func resourceAliCloudThreatDetectionServiceLinkedRoleCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "CreateServiceLinkedRole"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	if v, ok := d.GetOk("service_linked_role"); ok {
		request["ServiceLinkedRole"] = v
	}

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.RpcPost("Sas", "2018-12-03", action, query, request, true)
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_threat_detection_service_linked_role", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(request["ServiceLinkedRole"]))

	return resourceAliCloudThreatDetectionServiceLinkedRoleRead(d, meta)
}

func resourceAliCloudThreatDetectionServiceLinkedRoleRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	threatDetectionServiceV2 := ThreatDetectionServiceV2{client}

	objectRaw, err := threatDetectionServiceV2.DescribeThreatDetectionServiceLinkedRole(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_threat_detection_service_linked_role DescribeThreatDetectionServiceLinkedRole Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("role_status", objectRaw["Status"])

	d.Set("service_linked_role", d.Id())

	return nil
}

func resourceAliCloudThreatDetectionServiceLinkedRoleDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	action := "DeleteServiceLinkedRole"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	request["RoleName"] = d.Id()

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = client.RpcPost("ResourceManager", "2020-03-31", action, query, request, true)
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

	threatDetectionServiceV2 := ThreatDetectionServiceV2{client}
	stateConf := BuildStateConf([]string{}, []string{"false"}, d.Timeout(schema.TimeoutDelete), 10*time.Second, threatDetectionServiceV2.DescribeAsyncThreatDetectionServiceLinkedRoleStateRefreshFunc(d, response, "$.RoleStatus.Status", []string{}))
	if jobDetail, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id(), jobDetail)
	}

	return nil
}
