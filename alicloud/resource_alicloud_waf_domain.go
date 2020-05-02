package alicloud

import (
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	waf_openapi "github.com/aliyun/alibaba-cloud-sdk-go/services/waf-openapi"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func resourceAlicloudWafDomain() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudWafDomainCreate,
		Read:   resourceAlicloudWafDomainRead,
		Update: resourceAlicloudWafDomainUpdate,
		Delete: resourceAlicloudWafDomainDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
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
			"domain": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
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
				ForceNew: true,
			},
			"source_ips": {
				Type:     schema.TypeSet,
				Required: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"status": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"write_time": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  120,
			},
		},
	}
}

func resourceAlicloudWafDomainCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	waf_openapiService := Waf_openapiService{client}

	request := waf_openapi.CreateCreateDomainRequest()
	if v, ok := d.GetOk("cluster_type"); ok {
		request.ClusterType = requests.NewInteger(convertClusterTypeRequest(v.(string)))
	}
	if v, ok := d.GetOk("connection_time"); ok {
		request.ConnectionTime = requests.NewInteger(v.(int))
	}
	request.Domain = d.Get("domain").(string)
	if v, ok := d.GetOk("http2_port"); ok {
		request.Http2Port = convertListToJsonString(v.(*schema.Set).List())
	}
	if v, ok := d.GetOk("http_port"); ok {
		request.HttpPort = convertListToJsonString(v.(*schema.Set).List())
	}
	if v, ok := d.GetOk("http_to_user_ip"); ok {
		request.HttpToUserIp = requests.NewInteger(convertHttpToUserIpRequest(v.(string)))
	}
	if v, ok := d.GetOk("https_port"); ok {
		request.HttpsPort = convertListToJsonString(v.(*schema.Set).List())
	}
	if v, ok := d.GetOk("https_redirect"); ok {
		request.HttpsRedirect = requests.NewInteger(convertHttpsRedirectRequest(v.(string)))
	}
	request.InstanceId = d.Get("instance_id").(string)
	request.IsAccessProduct = requests.NewInteger(convertIsAccessProductRequest(d.Get("is_access_product").(string)))
	if v, ok := d.GetOk("load_balancing"); ok {
		request.LoadBalancing = requests.NewInteger(convertLoadBalancingRequest(v.(string)))
	}
	if v, ok := d.GetOk("log_headers"); ok {
		logHeaders, err := waf_openapiService.convertLogHeadersToString(v.(*schema.Set).List())
		if err != nil {
			return WrapError(err)
		}
		request.LogHeaders = logHeaders
	}
	if v, ok := d.GetOk("read_time"); ok {
		request.ReadTime = requests.NewInteger(v.(int))
	}
	if v, ok := d.GetOk("resource_group_id"); ok {
		request.ResourceGroupId = v.(string)
	}
	request.SourceIps = convertListToJsonString(d.Get("source_ips").(*schema.Set).List())
	if v, ok := d.GetOk("write_time"); ok {
		request.WriteTime = requests.NewInteger(v.(int))
	}
	raw, err := client.WithWafOpenapiClient(func(waf_openapiClient *waf_openapi.Client) (interface{}, error) {
		return waf_openapiClient.CreateDomain(request)
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_waf_domain", request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw)
	d.SetId(d.Get("instance_id").(string) + ":" + d.Get("domain").(string))

	return resourceAlicloudWafDomainRead(d, meta)
}
func resourceAlicloudWafDomainRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	waf_openapiService := Waf_openapiService{client}
	object, err := waf_openapiService.DescribeWafDomain(d.Id())
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}
	d.Set("domain", parts[1])
	d.Set("instance_id", parts[0])
	d.Set("cluster_type", convertClusterTypeResponse(object.ClusterType))
	d.Set("cname", object.Cname)
	d.Set("connection_time", object.ConnectionTime)
	d.Set("http2_port", object.Http2Port)
	d.Set("http_port", object.HttpPort)
	d.Set("http_to_user_ip", convertHttpToUserIpResponse(object.HttpToUserIp))
	d.Set("https_port", object.HttpsPort)
	d.Set("https_redirect", convertHttpsRedirectResponse(object.HttpsRedirect))
	d.Set("is_access_product", convertIsAccessProductResponse(object.IsAccessProduct))
	d.Set("load_balancing", convertLoadBalancingResponse(object.LoadBalancing))
	logHeaders := make([]map[string]interface{}, len(object.LogHeaders))
	for i, v := range object.LogHeaders {
		logHeaders[i] = map[string]interface{}{
			"key":   v.K,
			"value": v.V,
		}
	}
	if err := d.Set("log_headers", logHeaders); err != nil {
		return WrapError(err)
	}
	d.Set("read_time", object.ReadTime)
	d.Set("resource_group_id", object.ResourceGroupId)
	d.Set("source_ips", object.SourceIps)
	d.Set("write_time", object.WriteTime)
	return nil
}
func resourceAlicloudWafDomainUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	waf_openapiService := Waf_openapiService{client}
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}
	d.Partial(true)

	if d.HasChange("cluster_type") {
		request := waf_openapi.CreateModifyDomainClusterTypeRequest()
		request.Domain = parts[1]
		request.InstanceId = parts[0]
		request.ClusterType = requests.NewInteger(convertClusterTypeRequest(d.Get("cluster_type").(string)))
		raw, err := client.WithWafOpenapiClient(func(waf_openapiClient *waf_openapi.Client) (interface{}, error) {
			return waf_openapiClient.ModifyDomainClusterType(request)
		})
		addDebug(request.GetActionName(), raw)
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		d.SetPartial("cluster_type")
	}
	update := false
	request := waf_openapi.CreateModifyDomainRequest()
	request.Domain = parts[1]
	request.InstanceId = parts[0]
	if d.HasChange("is_access_product") {
		update = true
	}
	request.IsAccessProduct = requests.NewInteger(convertIsAccessProductRequest(d.Get("is_access_product").(string)))
	if d.HasChange("source_ips") {
		update = true
	}
	request.SourceIps = convertListToJsonString(d.Get("source_ips").(*schema.Set).List())
	request.ClusterType = requests.NewInteger(convertClusterTypeRequest(d.Get("cluster_type").(string)))
	if d.HasChange("connection_time") {
		update = true
	}
	request.ConnectionTime = requests.NewInteger(d.Get("connection_time").(int))
	if d.HasChange("http2_port") {
		update = true
	}
	request.Http2Port = convertListToJsonString(d.Get("http2_port").(*schema.Set).List())
	if d.HasChange("http_port") {
		update = true
	}
	request.HttpPort = convertListToJsonString(d.Get("http_port").(*schema.Set).List())
	if d.HasChange("http_to_user_ip") {
		update = true
	}
	request.HttpToUserIp = requests.NewInteger(convertHttpToUserIpRequest(d.Get("http_to_user_ip").(string)))
	if d.HasChange("https_port") {
		update = true
	}
	request.HttpsPort = convertListToJsonString(d.Get("https_port").(*schema.Set).List())
	if d.HasChange("https_redirect") {
		update = true
	}
	request.HttpsRedirect = requests.NewInteger(convertHttpsRedirectRequest(d.Get("https_redirect").(string)))
	if d.HasChange("load_balancing") {
		update = true
	}
	request.LoadBalancing = requests.NewInteger(convertLoadBalancingRequest(d.Get("load_balancing").(string)))
	if d.HasChange("log_headers") {
		update = true
	}
	logHeaders, err := waf_openapiService.convertLogHeadersToString(d.Get("log_headers").(*schema.Set).List())
	if err != nil {
		return WrapError(err)
	}
	request.LogHeaders = logHeaders
	if d.HasChange("read_time") {
		update = true
	}
	request.ReadTime = requests.NewInteger(d.Get("read_time").(int))
	if d.HasChange("write_time") {
		update = true
	}
	request.WriteTime = requests.NewInteger(d.Get("write_time").(int))
	if update {
		raw, err := client.WithWafOpenapiClient(func(waf_openapiClient *waf_openapi.Client) (interface{}, error) {
			return waf_openapiClient.ModifyDomain(request)
		})
		addDebug(request.GetActionName(), raw)
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		d.SetPartial("is_access_product")
		d.SetPartial("source_ips")
		d.SetPartial("connection_time")
		d.SetPartial("http2_port")
		d.SetPartial("http_port")
		d.SetPartial("http_to_user_ip")
		d.SetPartial("https_port")
		d.SetPartial("https_redirect")
		d.SetPartial("load_balancing")
		d.SetPartial("log_headers")
		d.SetPartial("read_time")
		d.SetPartial("write_time")
	}
	d.Partial(false)
	return resourceAlicloudWafDomainRead(d, meta)
}
func resourceAlicloudWafDomainDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}
	request := waf_openapi.CreateDeleteDomainRequest()
	request.Domain = parts[1]
	request.InstanceId = parts[0]
	raw, err := client.WithWafOpenapiClient(func(waf_openapiClient *waf_openapi.Client) (interface{}, error) {
		return waf_openapiClient.DeleteDomain(request)
	})
	addDebug(request.GetActionName(), raw)
	if err != nil {
		if IsExpectedErrors(err, []string{"DomainNotExist"}) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
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
