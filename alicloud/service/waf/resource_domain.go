package waf

import (
	"fmt"
	"log"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/convert"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func ResourceDomain() *schema.Resource {
	return &schema.Resource{
		Create: resourceDomainCreate,
		Read:   resourceDomainRead,
		Update: resourceDomainUpdate,
		Delete: resourceDomainDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		DeprecationMessage: "This resource has been deprecated since v1.211.0 and will be removed in the future. Please use 'alicloud_wafv3_domain' instead.",
		Schema: map[string]*schema.Schema{
			"cluster_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"PhysicalCluster", "VirtualCluster"}, false),
				Default:      "PhysicalCluster",
			},
			"cname": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"connection_time": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  5,
			},
			"domain_name": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ForceNew:      true,
				ConflictsWith: []string{"domain"},
			},
			"domain": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ForceNew:      true,
				Deprecated:    "Field 'domain' has been deprecated from version 1.94.0. Use 'domain_name' instead.",
				ConflictsWith: []string{"domain_name"},
			},
			"http2_port": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"http_port": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"http_to_user_ip": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"Off", "On"}, false),
				Default:      "Off",
			},
			"https_port": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"https_redirect": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"Off", "On"}, false),
				Default:      "Off",
			},
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"is_access_product": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"Off", "On"}, false),
			},
			"load_balancing": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"IpHash", "RoundRobin"}, false),
				Default:      "IpHash",
			},
			"log_headers": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"value": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"read_time": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  120,
			},
			"resource_group_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"source_ips": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"write_time": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  120,
			},
		},
	}
}

func resourceDomainCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	wafOpenapiService := NewService(client)
	var response map[string]interface{}
	action := "CreateDomain"
	request := make(map[string]interface{})
	var err error
	if v, ok := d.GetOk("cluster_type"); ok {
		request["ClusterType"] = convertClusterTypeRequest(v.(string))
	}

	if v, ok := d.GetOk("connection_time"); ok {
		request["ConnectionTime"] = v
	}

	if v, ok := d.GetOk("domain_name"); ok {
		request["Domain"] = v
	} else if v, ok := d.GetOk("domain"); ok {
		request["Domain"] = v
	} else {
		return errs.WrapError(errs.Error(`[ERROR] Argument "domain" or "domain_name" must be set one!`))
	}

	if v, ok := d.GetOk("http2_port"); ok {
		request["Http2Port"] = convert.ConvertListToJsonString(v.(*schema.Set).List())
	}

	if v, ok := d.GetOk("http_port"); ok {
		request["HttpPort"] = convert.ConvertListToJsonString(v.(*schema.Set).List())
	}

	if v, ok := d.GetOk("http_to_user_ip"); ok {
		request["HttpToUserIp"] = convertHttpToUserIpRequest(v.(string))
	}

	if v, ok := d.GetOk("https_port"); ok {
		request["HttpsPort"] = convert.ConvertListToJsonString(v.(*schema.Set).List())
	}

	if v, ok := d.GetOk("https_redirect"); ok {
		request["HttpsRedirect"] = convertHttpsRedirectRequest(v.(string))
	}

	request["InstanceId"] = d.Get("instance_id")

	request["IsAccessProduct"] = convertIsAccessProductRequest(d.Get("is_access_product").(string))
	if v, ok := d.GetOk("load_balancing"); ok {
		request["LoadBalancing"] = convertLoadBalancingRequest(v.(string))
	}

	if v, ok := d.GetOk("log_headers"); ok {
		logHeaders, err := wafOpenapiService.convertLogHeadersToString(v.(*schema.Set).List())
		if err != nil {
			return errs.WrapError(err)
		}
		request["LogHeaders"] = logHeaders
	}

	if v, ok := d.GetOk("read_time"); ok {
		request["ReadTime"] = v
	}

	if v, ok := d.GetOk("resource_group_id"); ok {
		request["ResourceGroupId"] = v
	}

	if v, ok := d.GetOk("source_ips"); ok {
		request["SourceIps"] = convert.ConvertListToJsonString(v.(*schema.Set).List())
	}

	if v, ok := d.GetOk("write_time"); ok {
		request["WriteTime"] = v
	}

	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(3*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("waf-openapi", "2019-09-10", action, nil, request, false)
		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		logs.AddDebug(action, response, request)
		return nil
	})
	if err != nil {
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, "alicloud_waf_domain", action, errs.AlibabaCloudSdkGoERROR)
	}
	logs.AddDebug(action, response, request)

	d.SetId(fmt.Sprint(request["InstanceId"], ":", request["Domain"]))

	return resourceDomainRead(d, meta)
}
func resourceDomainRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	wafOpenapiService := NewService(client)
	object, err := wafOpenapiService.DescribeWafDomain(d.Id())
	if err != nil {
		if errs.NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_waf_domain wafOpenapiService.DescribeWafDomain Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return errs.WrapError(err)
	}
	parts, err := helper.ParseResourceId(d.Id(), 2)
	if err != nil {
		return errs.WrapError(err)
	}
	d.Set("domain_name", parts[1])
	d.Set("domain", parts[1])
	d.Set("instance_id", parts[0])
	clusterType, err := convert.FormatInt(object["ClusterType"])
	if err != nil {
		return err
	}
	d.Set("cluster_type", convertClusterTypeResponse(clusterType))
	d.Set("cname", object["Cname"])
	connection_timeVal, err := convert.FormatInt(object["ConnectionTime"])
	if err != nil {
		return err
	}
	d.Set("connection_time", connection_timeVal)
	http2_portVal, err := convert.ConvertJsonStringToStringList(object["Http2Port"])
	if err != nil {
		return err
	}
	d.Set("http2_port", http2_portVal)
	http_portVal, err := convert.ConvertJsonStringToStringList(object["HttpPort"])
	if err != nil {
		return err
	}
	d.Set("http_port", http_portVal)
	httpToUserIp, err := convert.FormatInt(object["HttpToUserIp"])
	if err != nil {
		return err
	}
	d.Set("http_to_user_ip", convertHttpToUserIpResponse(httpToUserIp))
	https_portVal, err := convert.ConvertJsonStringToStringList(object["HttpsPort"])
	if err != nil {
		return err
	}
	d.Set("https_port", https_portVal)
	httpsRedirect, err := convert.FormatInt(object["HttpsRedirect"])
	if err != nil {
		return err
	}
	d.Set("https_redirect", convertHttpsRedirectResponse(httpsRedirect))
	isAccessProduct, err := convert.FormatInt(object["IsAccessProduct"])
	if err != nil {
		return err
	}
	d.Set("is_access_product", convertIsAccessProductResponse(isAccessProduct))
	loadBalancing, err := convert.FormatInt(object["LoadBalancing"])
	if err != nil {
		return err
	}
	d.Set("load_balancing", convertLoadBalancingResponse(loadBalancing))
	if v, ok := object["LogHeaders"].([]interface{}); ok {
		logHeaders := make([]map[string]interface{}, 0)
		for _, val := range v {
			item := val.(map[string]interface{})
			logHeaders = append(logHeaders, map[string]interface{}{
				"key":   item["k"].(string),
				"value": item["v"].(string),
			})
		}
		if err := d.Set("log_headers", logHeaders); err != nil {
			return errs.WrapError(err)
		}
	}
	read_timeVal, err := convert.FormatInt(object["ReadTime"])
	if err != nil {
		return err
	}
	d.Set("read_time", read_timeVal)
	d.Set("resource_group_id", object["ResourceGroupId"])
	d.Set("source_ips", object["SourceIps"])
	write_timeVal, err := convert.FormatInt(object["WriteTime"])
	if err != nil {
		return err
	}
	d.Set("write_time", write_timeVal)
	return nil
}
func resourceDomainUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	wafOpenapiService := NewService(client)
	var response map[string]interface{}
	parts, err := helper.ParseResourceId(d.Id(), 2)
	if err != nil {
		return errs.WrapError(err)
	}
	d.Partial(true)

	if d.HasChange("cluster_type") {
		request := map[string]interface{}{
			"Domain":     parts[1],
			"InstanceId": parts[0],
		}
		request["ClusterType"] = convertClusterTypeRequest(d.Get("cluster_type").(string))
		action := "ModifyDomainClusterType"
		wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
		err = retry.Retry(3*time.Minute, func() *retry.RetryError {
			response, err = client.RpcPost("waf-openapi", "2019-09-10", action, nil, request, false)
			if err != nil {
				if helper.NeedRetry(err) {
					wait()
					return retry.RetryableError(err)
				}
				return retry.NonRetryableError(err)
			}
			logs.AddDebug(action, response, request)
			return nil
		})
		if err != nil {
			return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
		}
	}
	update := false
	request := map[string]interface{}{
		"Domain":     parts[1],
		"InstanceId": parts[0],
	}
	if d.HasChange("is_access_product") {
		update = true
	}
	request["IsAccessProduct"] = convertIsAccessProductRequest(d.Get("is_access_product").(string))
	if d.HasChange("connection_time") {
		update = true
	}
	request["ConnectionTime"] = d.Get("connection_time")
	if d.HasChange("http2_port") {
		update = true
	}
	request["Http2Port"] = convert.ConvertListToJsonString(d.Get("http2_port").(*schema.Set).List())
	if d.HasChange("http_port") {
		update = true
	}
	request["HttpPort"] = convert.ConvertListToJsonString(d.Get("http_port").(*schema.Set).List())
	if d.HasChange("http_to_user_ip") {
		update = true
	}
	request["HttpToUserIp"] = convertHttpToUserIpRequest(d.Get("http_to_user_ip").(string))
	if d.HasChange("https_port") {
		update = true
	}
	request["HttpsPort"] = convert.ConvertListToJsonString(d.Get("https_port").(*schema.Set).List())
	if d.HasChange("https_redirect") {
		update = true
	}
	request["HttpsRedirect"] = convertHttpsRedirectRequest(d.Get("https_redirect").(string))
	if d.HasChange("load_balancing") {
		update = true
	}
	request["LoadBalancing"] = convertLoadBalancingRequest(d.Get("load_balancing").(string))
	if d.HasChange("log_headers") {
		update = true
	}
	logHeaders, err := wafOpenapiService.convertLogHeadersToString(d.Get("log_headers").(*schema.Set).List())
	if err != nil {
		return errs.WrapError(err)
	}
	request["LogHeaders"] = logHeaders
	if d.HasChange("read_time") {
		update = true
	}
	request["ReadTime"] = d.Get("read_time")
	if d.HasChange("source_ips") {
		update = true
	}
	request["SourceIps"] = convert.ConvertListToJsonString(d.Get("source_ips").(*schema.Set).List())
	if d.HasChange("write_time") {
		update = true
	}
	request["WriteTime"] = d.Get("write_time")
	if update {
		action := "ModifyDomain"
		wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
		err = retry.Retry(3*time.Minute, func() *retry.RetryError {
			response, err = client.RpcPost("waf-openapi", "2019-09-10", action, nil, request, false)
			if err != nil {
				if helper.NeedRetry(err) {
					wait()
					return retry.RetryableError(err)
				}
				return retry.NonRetryableError(err)
			}
			logs.AddDebug(action, response, request)
			return nil
		})
		if err != nil {
			return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
		}
	}
	d.Partial(false)
	return resourceDomainRead(d, meta)
}
func resourceDomainDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	parts, err := helper.ParseResourceId(d.Id(), 2)
	if err != nil {
		return errs.WrapError(err)
	}
	action := "DeleteDomain"
	var response map[string]interface{}
	request := map[string]interface{}{
		"Domain":     parts[1],
		"InstanceId": parts[0],
	}

	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(3*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("waf-openapi", "2019-09-10", action, nil, request, false)
		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		logs.AddDebug(action, response, request)
		return nil
	})
	if err != nil {
		if errs.IsExpectedErrors(err, []string{"DomainNotExist"}) {
			return nil
		}
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
	}
	return nil
}
func convertClusterTypeRequest(source string) int {
	switch source {
	case "PhysicalCluster":
		return 0
	case "VirtualCluster":
		return 1
	}
	return 0
}
func convertHttpToUserIpRequest(source string) int {
	switch source {
	case "Off":
		return 0
	case "On":
		return 1
	}
	return 0
}
func convertHttpsRedirectRequest(source string) int {
	switch source {
	case "Off":
		return 0
	case "On":
		return 1
	}
	return 0
}
func convertIsAccessProductRequest(source string) int {
	switch source {
	case "Off":
		return 0
	case "On":
		return 1
	}
	return 0
}
func convertLoadBalancingRequest(source string) int {
	switch source {
	case "IpHash":
		return 0
	case "RoundRobin":
		return 1
	}
	return 0
}
func convertClusterTypeResponse(source int) string {
	switch source {
	case 0:
		return "PhysicalCluster"
	case 1:
		return "VirtualCluster"
	}
	return ""
}
func convertHttpToUserIpResponse(source int) string {
	switch source {
	case 0:
		return "Off"
	case 1:
		return "On"
	}
	return ""
}
func convertHttpsRedirectResponse(source int) string {
	switch source {
	case 0:
		return "Off"
	case 1:
		return "On"
	}
	return ""
}
func convertIsAccessProductResponse(source int) string {
	switch source {
	case 0:
		return "Off"
	case 1:
		return "On"
	}
	return ""
}
func convertLoadBalancingResponse(source int) string {
	switch source {
	case 0:
		return "IpHash"
	case 1:
		return "RoundRobin"
	}
	return ""
}
