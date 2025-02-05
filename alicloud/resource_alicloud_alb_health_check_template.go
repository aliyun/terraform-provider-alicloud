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
			Create: schema.DefaultTimeout(2 * time.Minute),
			Update: schema.DefaultTimeout(2 * time.Minute),
			Delete: schema.DefaultTimeout(2 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"health_check_template_name": {
				Type:     schema.TypeString,
				Required: true,
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
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					if v, ok := d.GetOk("health_check_protocol"); ok && v.(string) == "HTTP" {
						return false
					}
					return true
				},
			},
			"health_check_http_version": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: StringInSlice([]string{"HTTP1.0", "HTTP1.1"}, false),
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					if v, ok := d.GetOk("health_check_protocol"); ok && v.(string) == "HTTP" {
						return false
					}
					return true
				},
			},
			"health_check_interval": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: IntBetween(1, 50),
			},
			"health_check_method": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: StringInSlice([]string{"HEAD", "GET"}, false),
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					if v, ok := d.GetOk("health_check_protocol"); ok && v.(string) == "HTTP" {
						return false
					}
					return true
				},
			},
			"health_check_path": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					if v, ok := d.GetOk("health_check_protocol"); ok && v.(string) == "HTTP" {
						return false
					}
					return true
				},
			},
			"health_check_protocol": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: StringInSlice([]string{"HTTP", "TCP"}, false),
			},
			"health_check_timeout": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: IntBetween(1, 300),
			},
			"healthy_threshold": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: IntBetween(2, 10),
			},
			"unhealthy_threshold": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: IntBetween(2, 10),
			},
			"health_check_codes": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					if v, ok := d.GetOk("health_check_protocol"); ok && v.(string) == "HTTP" {
						return false
					}
					return true
				},
			},
			"dry_run": {
				Type:     schema.TypeBool,
				Optional: true,
			},
		},
	}
}

func resourceAliCloudAlbHealthCheckTemplateCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "CreateHealthCheckTemplate"
	request := make(map[string]interface{})
	var err error

	request["ClientToken"] = buildClientToken("CreateHealthCheckTemplate")
	request["HealthCheckTemplateName"] = d.Get("health_check_template_name")

	if v, ok := d.GetOkExists("health_check_connect_port"); ok {
		request["HealthCheckConnectPort"] = v
	}

	if v, ok := d.GetOk("health_check_host"); ok {
		request["HealthCheckHost"] = v
	}

	if v, ok := d.GetOk("health_check_http_version"); ok {
		request["HealthCheckHttpVersion"] = v
	}

	if v, ok := d.GetOkExists("health_check_interval"); ok {
		request["HealthCheckInterval"] = v
	}

	if v, ok := d.GetOk("health_check_method"); ok {
		request["HealthCheckMethod"] = v
	}

	if v, ok := d.GetOk("health_check_path"); ok {
		request["HealthCheckPath"] = v
	}

	if v, ok := d.GetOk("health_check_protocol"); ok {
		request["HealthCheckProtocol"] = v
	}

	if v, ok := d.GetOkExists("health_check_timeout"); ok {
		request["HealthCheckTimeout"] = v
	}

	if v, ok := d.GetOkExists("healthy_threshold"); ok {
		request["HealthyThreshold"] = v
	}

	if v, ok := d.GetOkExists("unhealthy_threshold"); ok {
		request["UnhealthyThreshold"] = v
	}

	if v, ok := d.GetOk("health_check_codes"); ok {
		request["HealthCheckCodes"] = v
	}

	if v, ok := d.GetOkExists("dry_run"); ok {
		request["DryRun"] = v
	}

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutCreate)), func() *resource.RetryError {
		response, err = client.RpcPost("Alb", "2020-06-16", action, nil, request, true)
		if err != nil {
			if IsExpectedErrors(err, []string{"IdempotenceProcessing", "QuotaExceeded.HealthCheckTemplatesNum", "SystemBusy"}) || NeedRetry(err) {
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
	albService := AlbService{client}

	object, err := albService.DescribeAlbHealthCheckTemplate(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_alb_health_check_template albService.DescribeAlbHealthCheckTemplate Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("health_check_template_name", object["HealthCheckTemplateName"])
	d.Set("health_check_host", object["HealthCheckHost"])
	d.Set("health_check_http_version", object["HealthCheckHttpVersion"])
	d.Set("health_check_method", object["HealthCheckMethod"])
	d.Set("health_check_path", object["HealthCheckPath"])
	d.Set("health_check_protocol", object["HealthCheckProtocol"])
	d.Set("health_check_codes", object["HealthCheckCodes"])

	if healthCheckConnectPort, ok := object["HealthCheckConnectPort"]; ok {
		d.Set("health_check_connect_port", formatInt(healthCheckConnectPort))
	}

	if healthCheckInterval, ok := object["HealthCheckInterval"]; ok && fmt.Sprint(healthCheckInterval) != "0" {
		d.Set("health_check_interval", formatInt(healthCheckInterval))
	}

	if healthCheckTimeout, ok := object["HealthCheckTimeout"]; ok && fmt.Sprint(healthCheckTimeout) != "0" {
		d.Set("health_check_timeout", formatInt(healthCheckTimeout))
	}

	if healthyThreshold, ok := object["HealthyThreshold"]; ok && fmt.Sprint(healthyThreshold) != "0" {
		d.Set("healthy_threshold", formatInt(healthyThreshold))
	}

	if unhealthyThreshold, ok := object["UnhealthyThreshold"]; ok && fmt.Sprint(unhealthyThreshold) != "0" {
		d.Set("unhealthy_threshold", formatInt(unhealthyThreshold))
	}

	return nil
}

func resourceAliCloudAlbHealthCheckTemplateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	var err error
	update := false

	request := map[string]interface{}{
		"ClientToken":           buildClientToken("UpdateHealthCheckTemplateAttribute"),
		"HealthCheckTemplateId": d.Id(),
	}

	if d.HasChange("health_check_template_name") {
		update = true
	}
	request["HealthCheckTemplateName"] = d.Get("health_check_template_name")

	if d.HasChange("health_check_connect_port") {
		update = true

		if v, ok := d.GetOkExists("health_check_connect_port"); ok {
			request["HealthCheckConnectPort"] = v
		}
	}

	if d.HasChange("health_check_host") {
		update = true
	}
	if v, ok := d.GetOk("health_check_host"); ok {
		request["HealthCheckHost"] = v
	}

	if d.HasChange("health_check_http_version") {
		update = true
	}
	if v, ok := d.GetOk("health_check_http_version"); ok {
		request["HealthCheckHttpVersion"] = v
	}

	if d.HasChange("health_check_interval") {
		update = true

		if v, ok := d.GetOkExists("health_check_interval"); ok {
			request["HealthCheckInterval"] = v
		}
	}

	if d.HasChange("health_check_method") {
		update = true
	}
	if v, ok := d.GetOk("health_check_method"); ok {
		request["HealthCheckMethod"] = v
	}

	if d.HasChange("health_check_path") {
		update = true
	}
	if v, ok := d.GetOk("health_check_path"); ok {
		request["HealthCheckPath"] = v
	}

	if d.HasChange("health_check_protocol") {
		update = true
	}
	if v, ok := d.GetOk("health_check_protocol"); ok {
		request["HealthCheckProtocol"] = v
	}

	if d.HasChange("health_check_timeout") {
		update = true

		if v, ok := d.GetOkExists("health_check_timeout"); ok {
			request["HealthCheckTimeout"] = v
		}
	}

	if d.HasChange("healthy_threshold") {
		update = true

		if v, ok := d.GetOkExists("healthy_threshold"); ok {
			request["HealthyThreshold"] = v
		}
	}

	if d.HasChange("unhealthy_threshold") {
		update = true

		if v, ok := d.GetOkExists("unhealthy_threshold"); ok {
			request["UnhealthyThreshold"] = v
		}
	}

	if d.HasChange("health_check_codes") {
		update = true
	}
	if v, ok := d.GetOk("health_check_codes"); ok {
		request["HealthCheckCodes"] = v
	}

	if v, ok := d.GetOkExists("dry_run"); ok {
		request["DryRun"] = v
	}

	if update {
		action := "UpdateHealthCheckTemplateAttribute"
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *resource.RetryError {
			response, err = client.RpcPost("Alb", "2020-06-16", action, nil, request, true)
			if err != nil {
				if IsExpectedErrors(err, []string{"IdempotenceProcessing", "IncorrectStatus.HealthCheckTemplate", "SystemBusy"}) || NeedRetry(err) {
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

	return resourceAliCloudAlbHealthCheckTemplateRead(d, meta)
}

func resourceAliCloudAlbHealthCheckTemplateDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	action := "DeleteHealthCheckTemplates"
	var response map[string]interface{}
	var err error

	request := map[string]interface{}{
		"ClientToken":              buildClientToken("DeleteHealthCheckTemplates"),
		"HealthCheckTemplateIds.1": d.Id(),
	}

	if v, ok := d.GetOkExists("dry_run"); ok {
		request["DryRun"] = v
	}

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutDelete)), func() *resource.RetryError {
		response, err = client.RpcPost("Alb", "2020-06-16", action, nil, request, true)
		if err != nil {
			if IsExpectedErrors(err, []string{"IdempotenceProcessing", "IncorrectStatus.HealthCheckTemplate", "SystemBusy"}) || NeedRetry(err) {
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
