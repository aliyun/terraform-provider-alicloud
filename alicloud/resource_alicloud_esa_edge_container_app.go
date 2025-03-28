// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
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

func resourceAliCloudEsaEdgeContainerApp() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudEsaEdgeContainerAppCreate,
		Read:   resourceAliCloudEsaEdgeContainerAppRead,
		Delete: resourceAliCloudEsaEdgeContainerAppDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(15 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"edge_container_app_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"health_check_fail_times": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
			},
			"health_check_host": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"health_check_http_code": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"health_check_interval": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"health_check_method": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"health_check_port": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"health_check_succ_times": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"health_check_timeout": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"health_check_type": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"health_check_uri": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"remarks": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"service_port": {
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"target_port": {
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceAliCloudEsaEdgeContainerAppCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "CreateEdgeContainerApp"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	request["RegionId"] = client.RegionId

	if v, ok := d.GetOkExists("health_check_fail_times"); ok {
		request["HealthCheckFailTimes"] = v
	}
	request["ServicePort"] = d.Get("service_port")
	if v, ok := d.GetOk("health_check_http_code"); ok {
		request["HealthCheckHttpCode"] = v
	}
	if v, ok := d.GetOk("health_check_host"); ok {
		request["HealthCheckHost"] = v
	}
	if v, ok := d.GetOkExists("health_check_timeout"); ok {
		request["HealthCheckTimeout"] = v
	}
	if v, ok := d.GetOk("health_check_method"); ok {
		request["HealthCheckMethod"] = v
	}
	request["Name"] = d.Get("edge_container_app_name")
	if v, ok := d.GetOk("health_check_type"); ok {
		request["HealthCheckType"] = v
	}
	if v, ok := d.GetOkExists("health_check_interval"); ok {
		request["HealthCheckInterval"] = v
	}
	if v, ok := d.GetOk("remarks"); ok {
		request["Remarks"] = v
	}
	if v, ok := d.GetOkExists("health_check_port"); ok {
		request["HealthCheckPort"] = v
	}
	if v, ok := d.GetOkExists("health_check_succ_times"); ok {
		request["HealthCheckSuccTimes"] = v
	}
	request["TargetPort"] = d.Get("target_port")
	if v, ok := d.GetOk("health_check_uri"); ok {
		request["HealthCheckURI"] = v
	}
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.RpcPost("ESA", "2024-09-10", action, query, request, true)
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_esa_edge_container_app", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(response["AppId"]))

	esaServiceV2 := EsaServiceV2{client}
	stateConf := BuildStateConf([]string{}, []string{"created"}, d.Timeout(schema.TimeoutCreate), 5*time.Minute, esaServiceV2.EsaEdgeContainerAppStateRefreshFunc(d.Id(), "$.App.Status", []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return resourceAliCloudEsaEdgeContainerAppRead(d, meta)
}

func resourceAliCloudEsaEdgeContainerAppRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	esaServiceV2 := EsaServiceV2{client}

	objectRaw, err := esaServiceV2.DescribeEsaEdgeContainerApp(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_esa_edge_container_app DescribeEsaEdgeContainerApp Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	appRawObj, _ := jsonpath.Get("$.App", objectRaw)
	appRaw := make(map[string]interface{})
	if appRawObj != nil {
		appRaw = appRawObj.(map[string]interface{})
	}
	d.Set("create_time", appRaw["CreateTime"])
	d.Set("edge_container_app_name", appRaw["Name"])
	d.Set("remarks", appRaw["Remarks"])
	d.Set("service_port", appRaw["ServicePort"])
	d.Set("status", appRaw["Status"])
	d.Set("target_port", appRaw["TargetPort"])

	healthCheckRawObj, _ := jsonpath.Get("$.App.HealthCheck", objectRaw)
	healthCheckRaw := make(map[string]interface{})
	if healthCheckRawObj != nil {
		healthCheckRaw = healthCheckRawObj.(map[string]interface{})
	}
	d.Set("health_check_fail_times", healthCheckRaw["FailTimes"])
	d.Set("health_check_host", healthCheckRaw["Host"])
	d.Set("health_check_http_code", healthCheckRaw["HttpCode"])
	d.Set("health_check_interval", healthCheckRaw["Interval"])
	d.Set("health_check_method", healthCheckRaw["Method"])
	d.Set("health_check_port", healthCheckRaw["Port"])
	d.Set("health_check_succ_times", healthCheckRaw["SuccTimes"])
	d.Set("health_check_timeout", healthCheckRaw["Timeout"])
	d.Set("health_check_type", healthCheckRaw["Type"])
	d.Set("health_check_uri", healthCheckRaw["Uri"])

	return nil
}

func resourceAliCloudEsaEdgeContainerAppDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	action := "DeleteEdgeContainerApp"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	request["AppId"] = d.Id()
	request["RegionId"] = client.RegionId

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = client.RpcPost("ESA", "2024-09-10", action, query, request, true)

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
