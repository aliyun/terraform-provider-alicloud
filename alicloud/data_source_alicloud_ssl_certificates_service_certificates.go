package alicloud

import (
	"fmt"
	"regexp"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func dataSourceAliCloudSslCertificatesServiceCertificates() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAliCloudSslCertificatesServiceCertificatesRead,
		Schema: map[string]*schema.Schema{
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"name_regex": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.ValidateRegexp,
			},
			"keyword": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"enable_details": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"names": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"certificates": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"cert_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"certificate_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"fingerprint": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"common": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"sans": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"org_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"issuer": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"country": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"province": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"city": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"expired": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"start_date": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"end_date": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"cert": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"key": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"buy_in_aliyun": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"name": {
							Type:       schema.TypeString,
							Computed:   true,
							Deprecated: "Field `name` has been deprecated from provider version 1.129.0. New field `certificate_name` instead.",
						},
					},
				},
			},
			"lang": {
				Type:       schema.TypeString,
				Optional:   true,
				Deprecated: "Field `lang` has been deprecated from provider version 1.232.0.",
			},
		},
	}
}

func dataSourceAliCloudSslCertificatesServiceCertificatesRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	action := "ListUserCertificateOrder"
	request := make(map[string]interface{})
	request["OrderType"] = "CERT"
	request["ShowSize"] = PageSizeLarge
	request["CurrentPage"] = 1

	if v, ok := d.GetOk("keyword"); ok {
		request["Keyword"] = v
	}

	var objects []map[string]interface{}
	idsMap := make(map[string]string)
	if v, ok := d.GetOk("ids"); ok {
		for _, vv := range v.([]interface{}) {
			if vv == nil {
				continue
			}
			idsMap[vv.(string)] = vv.(string)
		}
	}

	var certificateNameRegex *regexp.Regexp
	if v, ok := d.GetOk("name_regex"); ok {
		r, err := regexp.Compile(v.(string))
		if err != nil {
			return WrapError(err)
		}
		certificateNameRegex = r
	}

	var response map[string]interface{}
	var err error
	for {
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			response, err = client.RpcPost("cas", "2020-04-07", action, nil, request, true)
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
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_ssl_certificates_service_certificates", action, AlibabaCloudSdkGoERROR)
		}

		resp, err := jsonpath.Get("$.CertificateOrderList", response)
		if err != nil {
			return WrapErrorf(err, FailedGetAttributeMsg, action, "$.CertificateOrderList", response)
		}

		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})
			if len(idsMap) > 0 {
				if _, ok := idsMap[fmt.Sprint(item["CertificateId"])]; !ok {
					continue
				}
			}

			if certificateNameRegex != nil && !certificateNameRegex.MatchString(fmt.Sprint(item["Name"])) {
				continue
			}

			objects = append(objects, item)
		}

		if len(result) < request["ShowSize"].(int) {
			break
		}

		request["CurrentPage"] = request["CurrentPage"].(int) + 1
	}

	ids := make([]string, 0)
	names := make([]interface{}, 0)
	s := make([]map[string]interface{}, 0)
	for _, object := range objects {
		mapping := map[string]interface{}{
			"id":               fmt.Sprint(object["CertificateId"]),
			"cert_id":          fmt.Sprint(object["CertificateId"]),
			"certificate_name": object["Name"],
			"fingerprint":      object["Fingerprint"],
			"common":           object["CommonName"],
			"sans":             object["Sans"],
			"org_name":         object["OrgName"],
			"issuer":           object["Issuer"],
			"country":          object["Country"],
			"province":         object["Province"],
			"city":             object["City"],
			"expired":          object["Expired"],
			"start_date":       object["StartDate"],
			"end_date":         object["EndDate"],
			"name":             object["Name"],
		}

		ids = append(ids, fmt.Sprint(mapping["id"]))
		names = append(names, object["Name"])

		if detailedEnabled := d.Get("enable_details"); !detailedEnabled.(bool) {
			s = append(s, mapping)
			continue
		}

		id := fmt.Sprint(object["CertificateId"])
		casService := CasService{client}

		describeUserCertificateDetail, err := casService.DescribeSslCertificatesServiceCertificate(id)
		if err != nil {
			return WrapError(err)
		}

		mapping["cert"] = describeUserCertificateDetail["Cert"]
		mapping["key"] = describeUserCertificateDetail["Key"]
		mapping["buy_in_aliyun"] = describeUserCertificateDetail["BuyInAliyun"]

		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))

	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}

	if err := d.Set("names", names); err != nil {
		return WrapError(err)
	}

	if err := d.Set("certificates", s); err != nil {
		return WrapError(err)
	}

	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}

	return nil
}
