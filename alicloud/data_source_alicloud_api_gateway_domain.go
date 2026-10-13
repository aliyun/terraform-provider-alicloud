package alicloud

import (
	"fmt"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceAlicloudApiGatewayDomain() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudApiGatewayDomainRead,

		Schema: map[string]*schema.Schema{
			"group_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"domain_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"output_file": {
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

func dataSourceAlicloudApiGatewayDomainRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	cloudApiService := CloudApiService{client}
	groupId := d.Get("group_id").(string)
	domainName := d.Get("domain_name").(string)

	id := fmt.Sprintf("%s%s%s", groupId, COLON_SEPARATED, domainName)
	domain, err := cloudApiService.DescribeApiGatewayDomain(id)
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.SetId(id)
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

	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		d.Set("output_file", nil)
	}

	return nil
}
