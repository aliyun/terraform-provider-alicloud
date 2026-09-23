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

func resourceAliCloudArmsPrometheusMonitoring() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudArmsPrometheusMonitoringCreate,
		Read:   resourceAliCloudArmsPrometheusMonitoringRead,
		Update: resourceAliCloudArmsPrometheusMonitoringUpdate,
		Delete: resourceAliCloudArmsPrometheusMonitoringDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"cluster_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"config_yaml": {
				Type:     schema.TypeString,
				Required: true,
				StateFunc: func(v interface{}) string {
					yaml, _ := normalizeYamlString(v)
					return yaml
				},
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					equal, _ := compareYamlTemplateAreEquivalent(old, new)
					return equal
				},
				ValidateFunc: validateYamlString,
			},
			"monitoring_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: StringInSlice([]string{"stop", "run"}, false),
			},
			"type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: StringInSlice([]string{"serviceMonitor", "podMonitor", "customJob", "probe"}, false),
			},
		},
	}
}

func resourceAliCloudArmsPrometheusMonitoringCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	action := "CreatePrometheusMonitoring"
	var request map[string]interface{}
	var response map[string]interface{}
	var err error
	request = make(map[string]interface{})
	request["ClusterId"] = d.Get("cluster_id")
	request["Type"] = d.Get("type")
	request["RegionId"] = client.RegionId

	request["ConfigYaml"] = d.Get("config_yaml")
	if v, ok := d.GetOk("status"); ok {
		request["Status"] = v
	}
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.RpcPost("ARMS", "2019-08-08", action, nil, request, false)

		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(action, response, request)
		return nil
	})

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_arms_prometheus_monitoring", action, AlibabaCloudSdkGoERROR)
	}

	code, _ := jsonpath.Get("$.Code", response)
	if fmt.Sprint(code) != "200" {
		log.Printf("[DEBUG] Resource alicloud_arms_prometheus_monitoring CreatePrometheusMonitoring Failed!!! %s", response)
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_arms_prometheus_monitoring", action, AlibabaCloudSdkGoERROR, response)
	}
	d.SetId(fmt.Sprintf("%v:%v:%v", request["ClusterId"], response["Data"], request["Type"]))

	return resourceAliCloudArmsPrometheusMonitoringRead(d, meta)
}

func resourceAliCloudArmsPrometheusMonitoringRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	armsServiceV2 := ArmsServiceV2{client}

	objectRaw, err := armsServiceV2.DescribeArmsPrometheusMonitoring(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_arms_prometheus_monitoring DescribeArmsPrometheusMonitoring Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("config_yaml", objectRaw["ConfigYaml"])
	d.Set("status", objectRaw["Status"])
	d.Set("cluster_id", objectRaw["ClusterId"])
	d.Set("monitoring_name", objectRaw["MonitoringName"])
	d.Set("type", objectRaw["Type"])

	return nil
}

func resourceAliCloudArmsPrometheusMonitoringUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	update := false
	d.Partial(true)
	parts := strings.Split(d.Id(), ":")
	action := "UpdatePrometheusMonitoring"
	var err error
	request = make(map[string]interface{})
	request["Type"] = parts[2]
	request["ClusterId"] = parts[0]
	request["MonitoringName"] = parts[1]
	request["RegionId"] = client.RegionId
	if d.HasChange("config_yaml") {
		update = true
	}
	request["ConfigYaml"] = d.Get("config_yaml")
	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("ARMS", "2019-08-08", action, nil, request, false)

			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			addDebug(action, response, request)
			return nil
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
		d.SetPartial("config_yaml")
	}
	update = false
	parts = strings.Split(d.Id(), ":")
	action = "UpdatePrometheusMonitoringStatus"
	request = make(map[string]interface{})
	request["MonitoringName"] = parts[1]
	request["ClusterId"] = parts[0]
	request["Type"] = parts[2]
	request["RegionId"] = client.RegionId
	if d.HasChange("status") {
		update = true
		request["Status"] = d.Get("status")
	}

	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("ARMS", "2019-08-08", action, nil, request, false)

			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			addDebug(action, response, request)
			return nil
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
		d.SetPartial("status")
	}

	d.Partial(false)
	return resourceAliCloudArmsPrometheusMonitoringRead(d, meta)
}

func resourceAliCloudArmsPrometheusMonitoringDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	parts := strings.Split(d.Id(), ":")
	action := "DeletePrometheusMonitoring"
	var request map[string]interface{}
	var response map[string]interface{}
	var err error
	request = make(map[string]interface{})
	request["MonitoringName"] = parts[1]
	request["ClusterId"] = parts[0]
	request["Type"] = parts[2]
	request["RegionId"] = client.RegionId

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = client.RpcPost("ARMS", "2019-08-08", action, nil, request, false)

		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(action, response, request)
		return nil
	})

	if err != nil {
		if IsExpectedErrors(err, []string{"404"}) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	return nil
}
