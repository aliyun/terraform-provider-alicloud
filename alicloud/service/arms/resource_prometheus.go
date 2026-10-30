package arms

import (
	"fmt"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/convert"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/tags"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func ResourcePrometheus() *schema.Resource {
	return &schema.Resource{
		Create: ResourcePrometheusCreate,
		Read:   ResourcePrometheusRead,
		Update: ResourcePrometheusUpdate,
		Delete: ResourcePrometheusDelete,
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
			"duration": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"archive_duration": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"payment_type": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringInSlice([]string{"POSTPAY", "POSTPAY_GB"}, false),
			},
			"tags": tags.TagsSchema(),
		},
	}
}

func ResourcePrometheusCreate(d *schema.ResourceData, meta interface{}) error {
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

	if v, ok := d.GetOkExists("duration"); ok {
		request["Duration"] = v
	}

	if v, ok := d.GetOkExists("archive_duration"); ok {
		request["ArchiveDuration"] = v
	}

	if v, ok := d.GetOk("payment_type"); ok {
		request["PaymentType"] = v
	}

	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutCreate)), func() *retry.RetryError {
		response, err = client.RpcPost("ARMS", "2019-08-08", action, nil, request, true)
		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	logs.AddDebug(action, response, request)

	if err != nil {
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, "alicloud_arms_prometheus", action, errs.AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(response["Data"]))

	return ResourcePrometheusUpdate(d, meta)
}

func ResourcePrometheusRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	armsService := Service{client}

	object, err := armsService.DescribeArmsPrometheus(d.Id())
	if err != nil {
		if !d.IsNewResource() && errs.NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return errs.WrapError(err)
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
	durationVal, err := convert.FormatInt(object["StorageDuration"])
	if err != nil {
		return err
	}
	d.Set("duration", durationVal)
	archive_durationVal, err := convert.FormatInt(object["ArchiveDuration"])
	if err != nil {
		return err
	}
	d.Set("archive_duration", archive_durationVal)
	d.Set("payment_type", object["PaymentType"])

	listTagResourcesObject, err := armsService.ListTagResources(d.Id(), "PROMETHEUS")
	if err != nil {
		return errs.WrapError(err)
	}

	d.Set("tags", tags.TagsToMap(listTagResourcesObject))

	return nil
}

func ResourcePrometheusUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	armsService := Service{client}
	var response map[string]interface{}
	var err error
	d.Partial(true)

	if d.HasChange("tags") {
		if err := armsService.SetResourceTags(d, "PROMETHEUS"); err != nil {
			return errs.WrapError(err)
		}
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
		wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
		err = retry.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *retry.RetryError {
			response, err = client.RpcPost("ARMS", "2019-08-08", action, nil, updatePrometheusGlobalViewReq, true)
			if err != nil {
				if helper.NeedRetry(err) {
					wait()
					return retry.RetryableError(err)
				}
				return retry.NonRetryableError(err)
			}
			return nil
		})
		logs.AddDebug(action, response, updatePrometheusGlobalViewReq)

		if err != nil {
			return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
		}

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
		wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
		err = retry.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *retry.RetryError {
			response, err = client.RpcPost("ARMS", "2019-08-08", action, nil, bindPrometheusGrafanaInstanceReq, true)
			if err != nil {
				if helper.NeedRetry(err) {
					wait()
					return retry.RetryableError(err)
				}
				return retry.NonRetryableError(err)
			}
			return nil
		})
		logs.AddDebug(action, response, bindPrometheusGrafanaInstanceReq)

		if err != nil {
			return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
		}

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
		wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
		err = retry.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *retry.RetryError {
			response, err = client.RpcPost("ARMS", "2019-08-08", action, nil, changeResourceGroupReq, true)
			if err != nil {
				if helper.NeedRetry(err) {
					wait()
					return retry.RetryableError(err)
				}
				return retry.NonRetryableError(err)
			}
			return nil
		})
		logs.AddDebug(action, response, changeResourceGroupReq)

		if err != nil {
			return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
		}

	}

	update = false
	updatePrometheusInstanceReq := map[string]interface{}{
		"RegionId":  client.RegionId,
		"ClusterId": d.Id(),
	}

	if !d.IsNewResource() && d.HasChange("duration") {
		update = true
	}
	if v, ok := d.GetOkExists("duration"); ok {
		updatePrometheusInstanceReq["StorageDuration"] = v
	}

	if !d.IsNewResource() && d.HasChange("archive_duration") {
		update = true
	}
	if v, ok := d.GetOkExists("archive_duration"); ok {
		updatePrometheusInstanceReq["ArchiveDuration"] = v
	}

	if !d.IsNewResource() && d.HasChange("payment_type") {
		update = true
	}
	if v, ok := d.GetOk("payment_type"); ok {
		updatePrometheusInstanceReq["PaymentType"] = v
	}

	if update {
		action := "UpdatePrometheusInstance"
		wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
		err = retry.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *retry.RetryError {
			response, err = client.RpcPost("ARMS", "2019-08-08", action, nil, updatePrometheusInstanceReq, true)
			if err != nil {
				if helper.NeedRetry(err) {
					wait()
					return retry.RetryableError(err)
				}
				return retry.NonRetryableError(err)
			}
			return nil
		})
		logs.AddDebug(action, response, updatePrometheusInstanceReq)

		if err != nil {
			return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
		}
	}

	d.Partial(false)

	return ResourcePrometheusRead(d, meta)
}

func ResourcePrometheusDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	action := "UninstallPromCluster"
	var response map[string]interface{}

	var err error

	request := map[string]interface{}{
		"RegionId":  client.RegionId,
		"ClusterId": d.Id(),
	}

	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutDelete)), func() *retry.RetryError {
		response, err = client.RpcPost("ARMS", "2019-08-08", action, nil, request, true)
		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	logs.AddDebug(action, response, request)

	if err != nil {
		if errs.NotFoundError(err) {
			return nil
		}
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
	}

	return nil
}
