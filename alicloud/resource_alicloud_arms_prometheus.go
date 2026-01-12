package alicloud

import (
	"fmt"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceAlicloudArmsPrometheus() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudArmsPrometheusCreate,
		Read:   resourceAliCloudArmsPrometheusRead,
		Update: resourceAliCloudArmsPrometheusUpdate,
		Delete: resourceAliCloudArmsPrometheusDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"cluster_type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"remote-write", "ecs", "global-view", "aliyun-cs"}, false),
			},
			"grafana_instance_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"vpc_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"vswitch_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"security_group_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"cluster_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"cluster_name": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"sub_clusters_json": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"resource_group_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"tags": tagsSchema(),
		},
	}
}

func resourceAliCloudArmsPrometheusCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "CreatePrometheusInstance"
	request := make(map[string]interface{})
	var err error

	request["RegionId"] = client.RegionId
	request["ClusterType"] = d.Get("cluster_type")
	request["GrafanaInstanceId"] = d.Get("grafana_instance_id")
	request["AllSubClustersSuccess"] = true

	if v, ok := d.GetOk("vpc_id"); ok {
		request["VpcId"] = v
	}

	if v, ok := d.GetOk("vswitch_id"); ok {
		request["VSwitchId"] = v
	}

	if v, ok := d.GetOk("security_group_id"); ok {
		request["SecurityGroupId"] = v
	}

	if v, ok := d.GetOk("cluster_id"); ok {
		request["ClusterId"] = v
	}

	if v, ok := d.GetOk("cluster_name"); ok {
		request["ClusterName"] = v
	}

	if v, ok := d.GetOk("sub_clusters_json"); ok {
		request["SubClustersJson"] = v
	}

	if v, ok := d.GetOk("resource_group_id"); ok {
		request["ResourceGroupId"] = v
	}

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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_arms_prometheus", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(response["Data"]))

	return resourceAliCloudArmsPrometheusUpdate(d, meta)
}

func resourceAliCloudArmsPrometheusRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	armsService := ArmsService{client}

	object, err := armsService.DescribeArmsPrometheus(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("cluster_type", object["ClusterType"])
	d.Set("grafana_instance_id", object["GrafanaInstanceId"])
	d.Set("vpc_id", object["VpcId"])
	d.Set("vswitch_id", object["VSwitchId"])
	d.Set("security_group_id", object["SecurityGroupId"])
	d.Set("cluster_id", object["ClusterId"])
	d.Set("cluster_name", object["ClusterName"])
	d.Set("sub_clusters_json", object["SubClustersJson"])
	d.Set("resource_group_id", object["ResourceGroupId"])

	listTagResourcesObject, err := armsService.ListTagResources(d.Id(), "PROMETHEUS")
	if err != nil {
		return WrapError(err)
	}

	d.Set("tags", tagsToMap(listTagResourcesObject))

	return nil
}

func resourceAliCloudArmsPrometheusUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	armsService := ArmsService{client}
	var response map[string]interface{}
	var err error
	d.Partial(true)

	if d.HasChange("tags") {
		if err := armsService.SetResourceTags(d, "PROMETHEUS"); err != nil {
			return WrapError(err)
		}
		d.SetPartial("tags")
	}

	update := false
	updatePrometheusGlobalViewReq := map[string]interface{}{
		"RegionId":              client.RegionId,
		"ClusterId":             d.Id(),
		"AllSubClustersSuccess": true,
	}

	if !d.IsNewResource() && d.HasChange("sub_clusters_json") {
		update = true
	}
	if v, ok := d.GetOk("sub_clusters_json"); ok {
		updatePrometheusGlobalViewReq["SubClustersJson"] = v
	}

	if update {
		action := "UpdatePrometheusGlobalView"
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *resource.RetryError {
			response, err = client.RpcPost("ARMS", "2019-08-08", action, nil, updatePrometheusGlobalViewReq, true)
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})
		addDebug(action, response, updatePrometheusGlobalViewReq)

		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}

		d.SetPartial("sub_clusters_json")
	}

	update = false
	bindPrometheusGrafanaInstanceReq := map[string]interface{}{
		"RegionId":  client.RegionId,
		"ClusterId": d.Id(),
	}

	if !d.IsNewResource() && d.HasChange("grafana_instance_id") {
		update = true
	}
	bindPrometheusGrafanaInstanceReq["GrafanaInstanceId"] = d.Get("grafana_instance_id")

	if update {
		action := "BindPrometheusGrafanaInstance"
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *resource.RetryError {
			response, err = client.RpcPost("ARMS", "2019-08-08", action, nil, bindPrometheusGrafanaInstanceReq, true)
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})
		addDebug(action, response, bindPrometheusGrafanaInstanceReq)

		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}

		d.SetPartial("grafana_instance_id")
	}

	update = false
	changeResourceGroupReq := map[string]interface{}{
		"RegionId":     client.RegionId,
		"ResourceId":   d.Id(),
		"ResourceType": "PROMETHEUS",
	}

	if !d.IsNewResource() && d.HasChange("resource_group_id") {
		update = true
	}
	if v, ok := d.GetOk("resource_group_id"); ok {
		changeResourceGroupReq["NewResourceGroupId"] = v
	}

	if update {
		action := "ChangeResourceGroup"
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *resource.RetryError {
			response, err = client.RpcPost("ARMS", "2019-08-08", action, nil, changeResourceGroupReq, true)
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})
		addDebug(action, response, changeResourceGroupReq)

		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}

		d.SetPartial("resource_group_id")
	}

	d.Partial(false)

	return resourceAliCloudArmsPrometheusRead(d, meta)
}

func resourceAliCloudArmsPrometheusDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	action := "UninstallPromCluster"
	var response map[string]interface{}

	var err error

	request := map[string]interface{}{
		"RegionId":  client.RegionId,
		"ClusterId": d.Id(),
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
