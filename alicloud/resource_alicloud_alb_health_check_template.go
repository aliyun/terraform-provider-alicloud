package alicloud

import (
	"fmt"
	"log"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAliCloudAlbHealthCheckTemplate() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudAlbHealthCheckTemplateCreate,
		Read:   resourceAliCloudAlbHealthCheckTemplateRead,
		Update: resourceAliCloudAlbHealthCheckTemplateUpdate,
		Delete: resourceAliCloudAlbHealthCheckTemplateDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"dry_run": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"health_check_codes": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"health_check_connect_port": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: IntBetween(0, 65535),
			},
			"health_check_host": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"health_check_http_version": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: StringInSlice([]string{"HTTP1.0", "HTTP1.1"}, false),
			},
			"health_check_interval": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: IntBetween(0, 50),
			},
			"health_check_method": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: StringInSlice([]string{"HEAD", "GET"}, false),
			},
			"health_check_path": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"health_check_protocol": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: StringInSlice([]string{"HTTP", "TCP", "gRPC", "HTTPS"}, false),
			},
			"health_check_template_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"health_check_timeout": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: IntBetween(0, 300),
			},
			"healthy_threshold": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: IntBetween(2, 10),
			},
			"resource_group_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"tags": tagsSchema(),
			"unhealthy_threshold": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: IntBetween(2, 10),
			},
		},
	}
}

func resourceAliCloudAlbHealthCheckTemplateCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "CreateHealthCheckTemplate"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})

	request["ClientToken"] = buildClientToken(action)

	if v, ok := d.GetOk("tags"); ok {
		tagsMap := ConvertTags(v.(map[string]interface{}))
		request["Tags"] = tagsMap
	}

	if v, ok := d.GetOk("resource_group_id"); ok {
		request["ResourceGroupId"] = v
	}
	if v, ok := d.GetOkExists("health_check_connect_port"); ok {
		request["HealthCheckConnectPort"] = v
	}
	if v, ok := d.GetOk("health_check_host"); ok {
		request["HealthCheckHost"] = v
	}
	if v, ok := d.GetOkExists("health_check_timeout"); ok && v.(int) > 0 {
		request["HealthCheckTimeout"] = v
	}
	if v, ok := d.GetOk("health_check_path"); ok {
		request["HealthCheckPath"] = v
	}
	if v, ok := d.GetOk("health_check_method"); ok {
		request["HealthCheckMethod"] = v
	}
	if v, ok := d.GetOkExists("health_check_interval"); ok && v.(int) > 0 {
		request["HealthCheckInterval"] = v
	}
	if v, ok := d.GetOk("health_check_protocol"); ok {
		request["HealthCheckProtocol"] = v
	}
	request["HealthCheckTemplateName"] = d.Get("health_check_template_name")
	if v, ok := d.GetOk("health_check_codes"); ok {
		healthCheckCodesMapsArray := v.([]interface{})
		request["HealthCheckCodes"] = healthCheckCodesMapsArray
	}

	if v, ok := d.GetOkExists("unhealthy_threshold"); ok && v.(int) > 0 {
		request["UnhealthyThreshold"] = v
	}
	if v, ok := d.GetOkExists("healthy_threshold"); ok && v.(int) > 0 {
		request["HealthyThreshold"] = v
	}
	if v, ok := d.GetOkExists("dry_run"); ok {
		request["DryRun"] = v
	}
	if v, ok := d.GetOk("health_check_http_version"); ok {
		request["HealthCheckHttpVersion"] = v
	}
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.RpcPost("Alb", "2020-06-16", action, query, request, true)
		if err != nil {
			if IsExpectedErrors(err, []string{"QuotaExceeded.HealthCheckTemplatesNum", "SystemBusy", "IdempotenceProcessing"}) || NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_alb_health_check_template", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(response["HealthCheckTemplateId"]))

	return resourceAliCloudAlbHealthCheckTemplateRead(d, meta)
}

func resourceAliCloudAlbHealthCheckTemplateRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	albServiceV2 := AlbServiceV2{client}

	objectRaw, err := albServiceV2.DescribeAlbHealthCheckTemplate(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_alb_health_check_template DescribeAlbHealthCheckTemplate Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("health_check_connect_port", objectRaw["HealthCheckConnectPort"])
	d.Set("health_check_host", objectRaw["HealthCheckHost"])
	d.Set("health_check_http_version", objectRaw["HealthCheckHttpVersion"])
	d.Set("health_check_interval", objectRaw["HealthCheckInterval"])
	d.Set("health_check_method", objectRaw["HealthCheckMethod"])
	d.Set("health_check_path", objectRaw["HealthCheckPath"])
	d.Set("health_check_protocol", objectRaw["HealthCheckProtocol"])
	d.Set("health_check_template_name", objectRaw["HealthCheckTemplateName"])
	d.Set("health_check_timeout", objectRaw["HealthCheckTimeout"])
	d.Set("healthy_threshold", objectRaw["HealthyThreshold"])
	d.Set("resource_group_id", objectRaw["ResourceGroupId"])
	d.Set("unhealthy_threshold", objectRaw["UnhealthyThreshold"])

	healthCheckCodesRaw := make([]interface{}, 0)
	if objectRaw["HealthCheckCodes"] != nil {
		healthCheckCodesRaw = objectRaw["HealthCheckCodes"].([]interface{})
	}

	d.Set("health_check_codes", healthCheckCodesRaw)
	tagsMaps := objectRaw["Tags"]
	d.Set("tags", tagsToMap(tagsMaps))

	return nil
}

func resourceAliCloudAlbHealthCheckTemplateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	update := false

	var err error
	action := "UpdateHealthCheckTemplateAttribute"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["HealthCheckTemplateId"] = d.Id()

	request["ClientToken"] = buildClientToken(action)
	if d.HasChange("health_check_connect_port") {
		update = true
		request["HealthCheckConnectPort"] = d.Get("health_check_connect_port")
	}

	if d.HasChange("health_check_host") {
		update = true
		request["HealthCheckHost"] = d.Get("health_check_host")
	}

	if d.HasChange("health_check_timeout") {
		update = true
		request["HealthCheckTimeout"] = d.Get("health_check_timeout")
	}

	if d.HasChange("health_check_path") {
		update = true
		request["HealthCheckPath"] = d.Get("health_check_path")
	}

	if d.HasChange("health_check_method") {
		update = true
		request["HealthCheckMethod"] = d.Get("health_check_method")
	}

	if d.HasChange("health_check_interval") {
		update = true
		request["HealthCheckInterval"] = d.Get("health_check_interval")
	}

	if d.HasChange("health_check_protocol") {
		update = true
		request["HealthCheckProtocol"] = d.Get("health_check_protocol")
	}

	if d.HasChange("health_check_template_name") {
		update = true
	}
	request["HealthCheckTemplateName"] = d.Get("health_check_template_name")
	if d.HasChange("health_check_codes") {
		update = true
		if v, ok := d.GetOk("health_check_codes"); ok || d.HasChange("health_check_codes") {
			healthCheckCodesMapsArray := v.([]interface{})
			request["HealthCheckCodes"] = healthCheckCodesMapsArray
		}
	}

	if d.HasChange("unhealthy_threshold") {
		update = true
		request["UnhealthyThreshold"] = d.Get("unhealthy_threshold")
	}

	if d.HasChange("healthy_threshold") {
		update = true
		request["HealthyThreshold"] = d.Get("healthy_threshold")
	}

	if d.HasChange("health_check_http_version") {
		update = true
		request["HealthCheckHttpVersion"] = d.Get("health_check_http_version")
	}

	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("Alb", "2020-06-16", action, query, request, true)
			if err != nil {
				if IsExpectedErrors(err, []string{"SystemBusy", "IncorrectStatus.HealthCheckTemplate", "IdempotenceProcessing"}) || NeedRetry(err) {
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

	if d.HasChange("tags") {
		albServiceV2 := AlbServiceV2{client}
		if err := albServiceV2.SetResourceTags(d, "healthchecktemplate"); err != nil {
			return WrapError(err)
		}
	}
	return resourceAliCloudAlbHealthCheckTemplateRead(d, meta)
}

func resourceAliCloudAlbHealthCheckTemplateDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	action := "DeleteHealthCheckTemplates"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	request["HealthCheckTemplateIds.1"] = d.Id()

	request["ClientToken"] = buildClientToken(action)

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = client.RpcPost("Alb", "2020-06-16", action, query, request, true)

		if err != nil {
			if IsExpectedErrors(err, []string{"SystemBusy", "IncorrectStatus.HealthCheckTemplate", "IdempotenceProcessing"}) || NeedRetry(err) {
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
