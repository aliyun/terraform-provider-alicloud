package alidns

import (
	alidnssdk "github.com/aliyun/alibaba-cloud-sdk-go/services/alidns"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func dataSourceAlicloudDnsDomainTxtGuid() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudDnsDomainTxtGuidRead,

		Schema: map[string]*schema.Schema{
			"domain_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"type": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"ADD_SUB_DOMAIN", "RETRIEVAL"}, false),
			},
			"lang": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"en", "zh", "jp"}, false),
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"rr": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"value": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceAlicloudDnsDomainTxtGuidRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	request := alidnssdk.CreateGetTxtRecordForVerifyRequest()
	request.RegionId = client.RegionId
	domainName := d.Get("domain_name").(string)
	request.DomainName = domainName
	request.Type = d.Get("type").(string)
	if lang, ok := d.GetOk("lang"); ok {
		request.Lang = lang.(string)
	}

	raw, err := client.WithDnsClient(func(dnsclient *alidnssdk.Client) (i interface{}, err error) {
		return dnsclient.GetTxtRecordForVerify(request)
	})
	if err != nil {
		return errs.WrapErrorf(err, errs.DataDefaultErrorMsg, "alicloud_dns_domain_txt_guid", request.GetActionName(), errs.AlibabaCloudSdkGoERROR)
	}
	logs.AddDebug(request.GetActionName(), raw, request.RpcRequest, request)

	response := raw.(*alidnssdk.GetTxtRecordForVerifyResponse)

	rr := response.RR
	value := response.Value

	s := map[string]interface{}{
		"rr":    rr,
		"value": value,
	}

	ids := []string{domainName, rr, value}
	d.SetId(helper.DataResourceIdHash(ids))
	if err := d.Set("rr", rr); err != nil {
		return errs.WrapError(err)
	}

	if err := d.Set("value", value); err != nil {
		return errs.WrapError(err)
	}

	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		helper.WriteToFile(output.(string), s)
	}
	return nil
}
