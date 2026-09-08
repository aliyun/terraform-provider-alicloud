package alicloud

import (
	"regexp"
	"strconv"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/cloudapi"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func dataSourceAlicloudApiGatewaySignatures() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudApiGatewaySignaturesRead,

		Schema: map[string]*schema.Schema{
			"name_regex": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsValidRegExp,
			},
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			// Computed values
			"names": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"signatures": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"signature_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"signature_key": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"signature_secret": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"created_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"modified_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"region_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceAlicloudApiGatewaySignaturesRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	request := cloudapi.CreateDescribeSignaturesRequest()
	request.RegionId = client.RegionId
	request.PageSize = requests.NewInteger(PageSizeLarge)
	request.PageNumber = requests.NewInteger(1)

	var signatures []cloudapi.SignatureInfo

	for {
		raw, err := client.WithCloudApiClient(func(cloudApiClient *cloudapi.Client) (interface{}, error) {
			return cloudApiClient.DescribeSignatures(request)
		})
		if err != nil {
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_api_gateway_signatures", request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		response, _ := raw.(*cloudapi.DescribeSignaturesResponse)

		signatures = append(signatures, response.SignatureInfos.SignatureInfo...)

		if len(signatures) >= response.TotalCount {
			break
		}

		page, err := getNextpageNumber(request.PageNumber)
		if err != nil {
			return WrapError(err)
		}
		request.PageNumber = page
	}

	var (
		filteredSignaturesTemp []cloudapi.SignatureInfo
		nameRegex              *regexp.Regexp
	)
	if v, ok := d.GetOk("name_regex"); ok && v.(string) != "" {
		r, err := regexp.Compile(v.(string))
		if err != nil {
			return WrapError(err)
		}
		nameRegex = r
	}

	var ids []string
	if v, ok := d.GetOk("ids"); ok {
		for _, vv := range v.([]interface{}) {
			ids = append(ids, vv.(string))
		}
	}

	for _, sig := range signatures {
		if nameRegex != nil && !nameRegex.MatchString(sig.SignatureName) {
			continue
		}
		if len(ids) > 0 {
			match := false
			for _, id := range ids {
				if sig.SignatureId == id {
					match = true
					break
				}
			}
			if !match {
				continue
			}
		}
		filteredSignaturesTemp = append(filteredSignaturesTemp, sig)
	}

	var names []string
	var sids []string
	var sMaps []map[string]interface{}
	for _, sig := range filteredSignaturesTemp {
		sMaps = append(sMaps, map[string]interface{}{
			"id":               sig.SignatureId,
			"signature_name":   sig.SignatureName,
			"signature_key":    sig.SignatureKey,
			"signature_secret": sig.SignatureSecret,
			"created_time":     sig.CreatedTime,
			"modified_time":    sig.ModifiedTime,
			"region_id":        sig.RegionId,
		})
		names = append(names, sig.SignatureName)
		sids = append(sids, sig.SignatureId)
	}

	d.SetId(strconv.Itoa(len(sMaps)))
	if err := d.Set("signatures", sMaps); err != nil {
		return WrapError(err)
	}
	if err := d.Set("names", names); err != nil {
		return WrapError(err)
	}
	if err := d.Set("ids", sids); err != nil {
		return WrapError(err)
	}

	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		if err := writeToFile(output.(string), sMaps); err != nil {
			return WrapError(err)
		}
	}
	return nil
}
