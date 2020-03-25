package alicloud

import (
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	waf_openapi "github.com/aliyun/alibaba-cloud-sdk-go/services/waf-openapi"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
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
				Type:     schema.TypeInt,
				Optional: true,
			},
			"cname": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"connection_time": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"domain": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"http2_port": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"http_port": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"http_to_user_ip": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"https_port": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"https_redirect": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"is_access_product": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"load_balancing": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"log_headers": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"read_time": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"resource_group_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"source_ips": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"write_time": {
				Type:     schema.TypeInt,
				Optional: true,
			},
		},
	}
}

func resourceAlicloudWafDomainCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	request := waf_openapi.CreateCreateDomainRequest()
	if v, ok := d.GetOk("cluster_type"); ok {
		request.ClusterType = requests.NewInteger(v.(int))
	}
	if v, ok := d.GetOk("connection_time"); ok {
		request.ConnectionTime = requests.NewInteger(v.(int))
	}
	request.Domain = d.Get("domain").(string)
	if v, ok := d.GetOk("http2_port"); ok {
		request.Http2Port = v.(string)
	}
	if v, ok := d.GetOk("http_port"); ok {
		request.HttpPort = v.(string)
	}
	if v, ok := d.GetOk("http_to_user_ip"); ok {
		request.HttpToUserIp = requests.NewInteger(v.(int))
	}
	if v, ok := d.GetOk("https_port"); ok {
		request.HttpsPort = v.(string)
	}
	if v, ok := d.GetOk("https_redirect"); ok {
		request.HttpsRedirect = requests.NewInteger(v.(int))
	}
	request.InstanceId = d.Get("instance_id").(string)
	request.IsAccessProduct = requests.NewInteger(d.Get("is_access_product").(int))
	if v, ok := d.GetOk("load_balancing"); ok {
		request.LoadBalancing = requests.NewInteger(v.(int))
	}
	if v, ok := d.GetOk("log_headers"); ok {
		request.LogHeaders = v.(string)
	}
	if v, ok := d.GetOk("read_time"); ok {
		request.ReadTime = requests.NewInteger(v.(int))
	}
	if v, ok := d.GetOk("resource_group_id"); ok {
		request.ResourceGroupId = v.(string)
	}
	if v, ok := d.GetOk("source_ips"); ok {
		request.SourceIps = v.(string)
	}
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
	d.Set("cluster_type", object.Domain.ClusterType)
	d.Set("cname", object.Domain.Cname)
	d.Set("http2_port", object.Domain.Http2Port)
	d.Set("http_port", object.Domain.HttpPort)
	d.Set("http_to_user_ip", object.Domain.HttpToUserIp)
	d.Set("https_port", object.Domain.HttpsPort)
	d.Set("https_redirect", object.Domain.HttpsRedirect)
	d.Set("is_access_product", object.Domain.IsAccessProduct)
	d.Set("load_balancing", object.Domain.LoadBalancing)
	d.Set("log_headers", object.Domain.LogHeaders)
	d.Set("source_ips", object.Domain.SourceIps)
	return nil
}
func resourceAlicloudWafDomainUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}
	update := false
	request := waf_openapi.CreateModifyDomainRequest()
	request.Domain = parts[1]
	request.InstanceId = parts[0]
	request.IsAccessProduct = requests.NewInteger(d.Get("is_access_product").(int))
	request.SourceIps = d.Get("source_ips").(string)
	if d.HasChange("cluster_type") {
		update = true
		request.ClusterType = requests.NewInteger(d.Get("cluster_type").(int))
	}
	if d.HasChange("connection_time") {
		update = true
		request.ConnectionTime = requests.NewInteger(d.Get("connection_time").(int))
	}
	if d.HasChange("http2_port") {
		update = true
		request.Http2Port = d.Get("http2_port").(string)
	}
	if d.HasChange("http_port") {
		update = true
		request.HttpPort = d.Get("http_port").(string)
	}
	if d.HasChange("http_to_user_ip") {
		update = true
		request.HttpToUserIp = requests.NewInteger(d.Get("http_to_user_ip").(int))
	}
	if d.HasChange("https_port") {
		update = true
		request.HttpsPort = d.Get("https_port").(string)
	}
	if d.HasChange("https_redirect") {
		update = true
		request.HttpsRedirect = requests.NewInteger(d.Get("https_redirect").(int))
	}
	if d.HasChange("load_balancing") {
		update = true
		request.LoadBalancing = requests.NewInteger(d.Get("load_balancing").(int))
	}
	if d.HasChange("log_headers") {
		update = true
		request.LogHeaders = d.Get("log_headers").(string)
	}
	if d.HasChange("read_time") {
		update = true
		request.ReadTime = requests.NewInteger(d.Get("read_time").(int))
	}
	if d.HasChange("write_time") {
		update = true
		request.WriteTime = requests.NewInteger(d.Get("write_time").(int))
	}
	if update {
		raw, err := client.WithWafOpenapiClient(func(waf_openapiClient *waf_openapi.Client) (interface{}, error) {
			return waf_openapiClient.ModifyDomain(request)
		})
		addDebug(request.GetActionName(), raw)
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
	}
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
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	return nil
}
