package alicloud

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/cloudapi"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAliyunApigatewayDomain() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliyunApigatewayDomainCreate,
		Read:   resourceAliyunApigatewayDomainRead,
		Update: resourceAliyunApigatewayDomainUpdate,
		Delete: resourceAliyunApigatewayDomainDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"group_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"domain_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			// The following attributes are update-only (@rac operatePrivateType:["update"], @readonly):
			// they are sent to the backing API only during Update, and are NOT persisted back by Read
			// (DescribeDomain does not return them). Their state values are preserved across refresh.
			"ssl_ocsp_enable": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"ssl_ocsp_cache_enable": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"client_cert_s_dn_pass_through": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"wss_enable": {
				Type:     schema.TypeString,
				Optional: true,
			},
			// Computed values populated from DescribeDomain response
			"sub_domain": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"domain_binding_status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"domain_remark": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"domain_web_socket_status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"domain_legal_status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"domain_cname_status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"certificate_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"certificate_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"certificate_valid_start": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"certificate_valid_end": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func resourceAliyunApigatewayDomainCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	groupId := d.Get("group_id").(string)
	domainName := d.Get("domain_name").(string)

	request := cloudapi.CreateSetDomainRequest()
	request.RegionId = client.RegionId
	request.GroupId = groupId
	request.DomainName = domainName

	if err := resource.Retry(5*time.Minute, func() *resource.RetryError {
		raw, err := client.WithCloudApiClient(func(cloudApiClient *cloudapi.Client) (interface{}, error) {
			return cloudApiClient.SetDomain(request)
		})
		if err != nil {
			if IsExpectedErrors(err, []string{"RepeatedCommit"}) {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		return nil
	}); err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_api_gateway_domain", request.GetActionName(), AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprintf("%s%s%s", groupId, COLON_SEPARATED, domainName))

	return resourceAliyunApigatewayDomainUpdate(d, meta)
}

func resourceAliyunApigatewayDomainRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	cloudApiService := CloudApiService{client}
	domain, err := cloudApiService.DescribeApiGatewayDomain(d.Id())
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("group_id", domain.GroupId)
	d.Set("domain_name", domain.DomainName)
	d.Set("sub_domain", domain.SubDomain)
	d.Set("domain_binding_status", domain.DomainBindingStatus)
	d.Set("domain_remark", domain.DomainRemark)
	d.Set("domain_web_socket_status", domain.DomainWebSocketStatus)
	d.Set("domain_legal_status", domain.DomainLegalStatus)
	d.Set("domain_cname_status", domain.DomainCNAMEStatus)
	d.Set("certificate_id", domain.CertificateId)
	d.Set("certificate_name", domain.CertificateName)
	d.Set("certificate_valid_start", domain.CertificateValidStart)
	d.Set("certificate_valid_end", domain.CertificateValidEnd)
	// ssl_ocsp_enable, ssl_ocsp_cache_enable, client_cert_s_dn_pass_through and wss_enable
	// are update-only (@readonly): DescribeDomain does not return them, so we intentionally
	// do NOT call d.Set here — their existing state values are preserved across refresh.
	return nil
}

func resourceAliyunApigatewayDomainUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	parts := strings.Split(d.Id(), COLON_SEPARATED)
	if len(parts) < 2 {
		return WrapError(fmt.Errorf("invalid resource id format: %s, expected <group_id>:<domain_name>", d.Id()))
	}
	groupId := parts[0]
	domainName := parts[1]

	// SslOcspEnable, SslOcspCacheEnable and ClientCertSDnPassThrough are delivered via
	// SetDomainCertificate. The Go SDK request struct does not carry these fields yet, so
	// they are appended as query parameters (same convention as resource_alicloud_api_gateway_api).
	if d.HasChanges("ssl_ocsp_enable", "ssl_ocsp_cache_enable", "client_cert_s_dn_pass_through") {
		request := cloudapi.CreateSetDomainCertificateRequest()
		request.RegionId = client.RegionId
		request.GroupId = groupId
		request.DomainName = domainName
		request.QueryParams["SslOcspEnable"] = strconv.FormatBool(d.Get("ssl_ocsp_enable").(bool))
		request.QueryParams["SslOcspCacheEnable"] = strconv.FormatBool(d.Get("ssl_ocsp_cache_enable").(bool))
		request.QueryParams["ClientCertSDnPassThrough"] = strconv.FormatBool(d.Get("client_cert_s_dn_pass_through").(bool))
		raw, err := client.WithCloudApiClient(func(cloudApiClient *cloudapi.Client) (interface{}, error) {
			return cloudApiClient.SetDomainCertificate(request)
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	}

	// WSSEnable is delivered via SetDomainWebSocketStatus.
	if d.HasChange("wss_enable") {
		request := cloudapi.CreateSetDomainWebSocketStatusRequest()
		request.RegionId = client.RegionId
		request.GroupId = groupId
		request.DomainName = domainName
		request.WSSEnable = d.Get("wss_enable").(string)
		raw, err := client.WithCloudApiClient(func(cloudApiClient *cloudapi.Client) (interface{}, error) {
			return cloudApiClient.SetDomainWebSocketStatus(request)
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	}

	return resourceAliyunApigatewayDomainRead(d, meta)
}

func resourceAliyunApigatewayDomainDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	parts := strings.Split(d.Id(), COLON_SEPARATED)
	if len(parts) < 2 {
		return WrapError(fmt.Errorf("invalid resource id format: %s, expected <group_id>:<domain_name>", d.Id()))
	}
	request := cloudapi.CreateDeleteDomainRequest()
	request.RegionId = client.RegionId
	request.GroupId = parts[0]
	request.DomainName = parts[1]

	if err := resource.Retry(5*time.Minute, func() *resource.RetryError {
		raw, err := client.WithCloudApiClient(func(cloudApiClient *cloudapi.Client) (interface{}, error) {
			return cloudApiClient.DeleteDomain(request)
		})
		if err != nil {
			if NotFoundError(err) || IsExpectedErrors(err, []string{"NotFoundDomain", "NotFound"}) {
				return nil
			}
			return resource.RetryableError(err)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		return nil
	}); err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	return nil
}
