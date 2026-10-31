package cloudmonitorservice

import (
	"encoding/json"
	"fmt"
	"log"
	"reflect"
	"sort"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/validators"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"gopkg.in/yaml.v2"
)

func resourceAlicloudCmsHybridMonitorFcTask() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudCmsHybridMonitorFcTaskCreate,
		Read:   resourceAlicloudCmsHybridMonitorFcTaskRead,
		Update: resourceAlicloudCmsHybridMonitorFcTaskUpdate,
		Delete: resourceAlicloudCmsHybridMonitorFcTaskDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(2 * time.Minute),
			Delete: schema.DefaultTimeout(2 * time.Minute),
			Update: schema.DefaultTimeout(2 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"hybrid_monitor_fc_task_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"namespace": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"target_user_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"yarm_config": {
				Type:     schema.TypeString,
				Required: true,
				StateFunc: func(v interface{}) string {
					yamlString, _ := validators.NormalizeYamlString(v)
					return yamlString
				},
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					equal, _ := compareCmsHybridMonitorFcTaskYamlConfigAreEquivalent(old, new)
					return equal
				},
				ValidateFunc: validators.ValidateYamlString,
			},
		},
	}
}

func resourceAlicloudCmsHybridMonitorFcTaskCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "CreateHybridMonitorTask"
	request := make(map[string]interface{})
	var err error
	request["CollectTargetType"] = "aliyun_fc"
	request["Namespace"] = d.Get("namespace")
	if v, ok := d.GetOk("target_user_id"); ok {
		request["TargetUserId"] = v
	}
	request["TaskType"] = "aliyun_fc"
	request["YARMConfig"] = d.Get("yarm_config")
	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
		response, err = client.RpcPost("Cms", "2019-01-01", action, nil, request, false)
		if err != nil {
			if errs.IsExpectedErrors(err, []string{"InternalError", "undefined"}) || helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	logs.AddDebug(action, response, request)
	if err != nil {
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, "alicloud_cms_hybrid_monitor_fc_task", action, errs.AlibabaCloudSdkGoERROR)
	}
	if fmt.Sprint(response["Success"]) == "false" {
		return errs.WrapError(fmt.Errorf("%s failed, response: %v", action, response))
	}

	d.SetId(fmt.Sprintf("%s:%s", response["TaskId"], request["Namespace"]))

	return resourceAlicloudCmsHybridMonitorFcTaskRead(d, meta)
}
func resourceAlicloudCmsHybridMonitorFcTaskRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	cmsService := CmsService{client}
	object, err := cmsService.DescribeCmsHybridMonitorFcTask(d.Id())
	if err != nil {
		if errs.NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_cms_hybrid_monitor_fc_task cmsService.DescribeCmsHybridMonitorFcTask Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return errs.WrapError(err)
	}
	parts, err := helper.ParseResourceId(d.Id(), 2)
	if err != nil {
		return errs.WrapError(err)
	}
	d.Set("hybrid_monitor_fc_task_id", parts[0])
	d.Set("namespace", parts[1])
	d.Set("target_user_id", object["TargetUserId"])
	d.Set("yarm_config", object["YARMConfig"])
	return nil
}
func resourceAlicloudCmsHybridMonitorFcTaskUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var err error
	var response map[string]interface{}
	parts, err := helper.ParseResourceId(d.Id(), 2)
	if err != nil {
		return errs.WrapError(err)
	}
	request := map[string]interface{}{
		"TargetUserId": parts[0],
		"Namespace":    parts[1],
	}
	request["CollectTargetType"] = "aliyun_fc"
	request["TaskType"] = "aliyun_fc"
	if d.HasChange("yarm_config") {
		request["YARMConfig"] = d.Get("yarm_config")
		action := "CreateHybridMonitorTask"

		wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
		err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
			response, err = client.RpcPost("Cms", "2019-01-01", action, nil, request, false)
			if err != nil {
				if errs.IsExpectedErrors(err, []string{"InternalError", "undefined"}) || helper.NeedRetry(err) {
					wait()
					return retry.RetryableError(err)
				}
				return retry.NonRetryableError(err)
			}
			return nil
		})
		logs.AddDebug(action, response, request)
		if err != nil {
			return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
		}
		if fmt.Sprint(response["Success"]) == "false" {
			return errs.WrapError(fmt.Errorf("%s failed, response: %v", action, response))
		}
	}

	return resourceAlicloudCmsHybridMonitorFcTaskRead(d, meta)
}
func resourceAlicloudCmsHybridMonitorFcTaskDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	parts, err := helper.ParseResourceId(d.Id(), 2)
	if err != nil {
		return errs.WrapError(err)
	}
	action := "DeleteHybridMonitorTask"
	var response map[string]interface{}
	request := map[string]interface{}{
		"TargetUserId": parts[0],
		"Namespace":    parts[1],
	}
	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutDelete), func() *retry.RetryError {
		response, err = client.RpcPost("Cms", "2019-01-01", action, nil, request, false)
		if err != nil {
			if errs.IsExpectedErrors(err, []string{"InternalError"}) || helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	logs.AddDebug(action, response, request)
	if err != nil {
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
	}
	if errs.IsExpectedErrorCodes(fmt.Sprint(response["Code"]), []string{"ResourceNotFound"}) {
		return nil
	}
	if fmt.Sprint(response["Success"]) == "false" {
		return errs.WrapError(fmt.Errorf("%s failed, response: %v", action, response))
	}
	return nil
}

// compareCmsHybridMonitorFcTaskYamlConfigAreEquivalent compares two YAML
// configs of the Cms hybrid-monitor FC task for semantic equivalence
// (independent of product/metric ordering), returning false on parse errors.
func compareCmsHybridMonitorFcTaskYamlConfigAreEquivalent(tem1, tem2 string) (bool, error) {
	type MetricList struct {
		MetricList []string `yaml:"metric_list"`
	}
	type Product struct {
		MetricInfo []MetricList `yaml:"metric_info"`
		Namespace  string       `yaml:"namespace"`
	}
	type Products struct {
		Products []Product
	}

	var P1 Products
	err := yaml.Unmarshal([]byte(tem1), &P1)
	if err != nil {
		return false, err
	}

	y1 := make([]string, 0)
	for _, product := range P1.Products {
		s1, _ := json.Marshal(product)
		y1 = append(y1, string(s1))
	}

	var P2 Products
	err = yaml.Unmarshal([]byte(tem2), &P2)
	if err != nil {
		return false, err
	}

	y2 := make([]string, 0)
	for _, product := range P2.Products {
		s2, _ := json.Marshal(product)
		y2 = append(y2, string(s2))
	}

	sort.Strings(y1)
	sort.Strings(y2)
	return reflect.DeepEqual(y1, y2), nil
}
