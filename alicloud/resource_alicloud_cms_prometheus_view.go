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

func resourceAliCloudCmsPrometheusView() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudCmsPrometheusViewCreate,
		Read:   resourceAliCloudCmsPrometheusViewRead,
		Update: resourceAliCloudCmsPrometheusViewUpdate,
		Delete: resourceAliCloudCmsPrometheusViewDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"auth_free_read_policy": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"enable_auth_free_read": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"prometheus_instances": {
				Type:     schema.TypeSet,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"user_id": {
							Type:     schema.TypeString,
							Required: true,
						},
						"prometheus_instance_id": {
							Type:     schema.TypeString,
							Required: true,
						},
						"region_id": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
			"prometheus_view_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"region_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"version": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"workspace": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceAliCloudCmsPrometheusViewCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := fmt.Sprintf("/prometheus-views")
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]*string)
	body := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})

	if v, ok := d.GetOk("prometheus_instances"); ok {
		prometheusInstancesMapsArray := make([]interface{}, 0)
		for _, dataLoop := range convertToInterfaceArray(v) {
			dataLoopTmp := dataLoop.(map[string]interface{})
			dataLoopMap := make(map[string]interface{})
			dataLoopMap["userId"] = dataLoopTmp["user_id"]
			dataLoopMap["prometheusInstanceId"] = dataLoopTmp["prometheus_instance_id"]
			dataLoopMap["regionId"] = dataLoopTmp["region_id"]
			prometheusInstancesMapsArray = append(prometheusInstancesMapsArray, dataLoopMap)
		}
		request["prometheusInstances"] = prometheusInstancesMapsArray
	}

	if v, ok := d.GetOkExists("enable_auth_free_read"); ok {
		request["enableAuthFreeRead"] = v
	}
	if v, ok := d.GetOk("auth_free_read_policy"); ok {
		request["authFreeReadPolicy"] = v
	}
	request["prometheusViewName"] = d.Get("prometheus_view_name")
	request["version"] = d.Get("version")
	request["workspace"] = d.Get("workspace")
	body = request
	wait := incrementalWait(10*time.Second, 10*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.RoaPost("Cms", "2024-03-30", action, query, nil, body, true)
		if err != nil {
			if IsExpectedErrors(err, []string{"500"}) || NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_cms_prometheus_view", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(response["prometheusViewId"]))

	return resourceAliCloudCmsPrometheusViewRead(d, meta)
}

func resourceAliCloudCmsPrometheusViewRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	cmsServiceV2 := CmsServiceV2{client}

	objectRaw, err := cmsServiceV2.DescribeCmsPrometheusView(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_cms_prometheus_view DescribeCmsPrometheusView Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("auth_free_read_policy", objectRaw["authFreeReadPolicy"])
	d.Set("create_time", objectRaw["createTime"])
	d.Set("enable_auth_free_read", objectRaw["enableAuthFreeRead"])
	d.Set("prometheus_view_name", objectRaw["prometheusViewName"])
	d.Set("region_id", objectRaw["regionId"])
	d.Set("version", objectRaw["version"])
	d.Set("workspace", objectRaw["workspace"])

	prometheusInstancesRaw := objectRaw["prometheusInstances"]
	prometheusInstancesMaps := make([]map[string]interface{}, 0)
	if prometheusInstancesRaw != nil {
		for _, prometheusInstancesChildRaw := range convertToInterfaceArray(prometheusInstancesRaw) {
			prometheusInstancesMap := make(map[string]interface{})
			prometheusInstancesChildRaw := prometheusInstancesChildRaw.(map[string]interface{})
			prometheusInstancesMap["prometheus_instance_id"] = prometheusInstancesChildRaw["prometheusInstanceId"]
			prometheusInstancesMap["region_id"] = prometheusInstancesChildRaw["regionId"]
			prometheusInstancesMap["user_id"] = prometheusInstancesChildRaw["userId"]

			prometheusInstancesMaps = append(prometheusInstancesMaps, prometheusInstancesMap)
		}
	}
	if err := d.Set("prometheus_instances", prometheusInstancesMaps); err != nil {
		return err
	}

	return nil
}

func resourceAliCloudCmsPrometheusViewUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]*string
	var body map[string]interface{}
	update := false

	var err error
	prometheusViewId := d.Id()
	action := fmt.Sprintf("/prometheus-views/%s", prometheusViewId)
	request = make(map[string]interface{})
	query = make(map[string]*string)
	body = make(map[string]interface{})

	if d.HasChange("prometheus_instances") {
		update = true
	}
	if v, ok := d.GetOk("prometheus_instances"); ok || d.HasChange("prometheus_instances") {
		prometheusInstancesMapsArray := make([]interface{}, 0)
		for _, dataLoop := range convertToInterfaceArray(v) {
			dataLoopTmp := dataLoop.(map[string]interface{})
			dataLoopMap := make(map[string]interface{})
			dataLoopMap["userId"] = dataLoopTmp["user_id"]
			dataLoopMap["prometheusInstanceId"] = dataLoopTmp["prometheus_instance_id"]
			dataLoopMap["regionId"] = dataLoopTmp["region_id"]
			prometheusInstancesMapsArray = append(prometheusInstancesMapsArray, dataLoopMap)
		}
		request["prometheusInstances"] = prometheusInstancesMapsArray
	}

	if d.HasChange("enable_auth_free_read") {
		update = true
	}
	if v, ok := d.GetOkExists("enable_auth_free_read"); ok || d.HasChange("enable_auth_free_read") {
		request["enableAuthFreeRead"] = v
	}
	if d.HasChange("auth_free_read_policy") {
		update = true
	}
	if v, ok := d.GetOk("auth_free_read_policy"); ok || d.HasChange("auth_free_read_policy") {
		request["authFreeReadPolicy"] = v
	}
	if d.HasChange("prometheus_view_name") {
		update = true
	}
	request["prometheusViewName"] = d.Get("prometheus_view_name")
	body = request
	if update {
		wait := incrementalWait(10*time.Second, 10*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RoaPut("Cms", "2024-03-30", action, query, nil, body, true)
			if err != nil {
				if IsExpectedErrors(err, []string{"500"}) || NeedRetry(err) {
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

	return resourceAliCloudCmsPrometheusViewRead(d, meta)
}

func resourceAliCloudCmsPrometheusViewDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	prometheusViewId := d.Id()
	action := fmt.Sprintf("/prometheus-views/%s", prometheusViewId)
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]*string)
	var err error
	request = make(map[string]interface{})

	wait := incrementalWait(10*time.Second, 10*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = client.RoaDelete("Cms", "2024-03-30", action, query, nil, nil, true)
		if err != nil {
			if IsExpectedErrors(err, []string{"500"}) || NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)

	if err != nil {
		if IsExpectedErrors(err, []string{"404"}) || NotFoundError(err) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	return nil
}
