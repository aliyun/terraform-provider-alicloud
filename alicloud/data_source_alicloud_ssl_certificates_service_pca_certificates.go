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

func dataSourceAliCloudSslCertificatesServicePcaCertificates() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAliCloudSslCertificatesServicePcaCertificatesRead,
		Schema: map[string]*schema.Schema{
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"name_regex": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsValidRegExp,
			},
			"ca_status": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"issue", "forbidden", "revoke"}, false),
			},
			"cert_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"root", "subRoot", "externalCa"}, false),
			},
			"issuer_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"local", "iTrusChina", "external"}, false),
			},
			"valid_status": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"valid", "notValid"}, false),
			},
			"resource_group_id": {
				Type:     schema.TypeString,
				Optional: true,
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
						"identifier": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"serial_number": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"x509_certificate": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"certificate_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"algorithm": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"sign_algorithm": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"sha2": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"md5": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"locality": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"organization": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"organization_unit": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"common_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"country_code": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"state": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"parent_identifier": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"years": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"before_date": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"after_date": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"resource_group_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceAliCloudSslCertificatesServicePcaCertificatesRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	action := "DescribeCACertificateList"
	request := make(map[string]interface{})
	request["ShowSize"] = PageSizeLarge
	request["CurrentPage"] = 1

	if v, ok := d.GetOk("ca_status"); ok {
		request["CaStatus"] = v
	}
	if v, ok := d.GetOk("cert_type"); ok {
		request["CertType"] = v
	}
	if v, ok := d.GetOk("issuer_type"); ok {
		request["IssuerType"] = v
	}
	if v, ok := d.GetOk("valid_status"); ok {
		request["ValidStatus"] = v
	}
	if v, ok := d.GetOk("resource_group_id"); ok {
		request["ResourceGroupId"] = v
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

	var commonNameRegex *regexp.Regexp
	if v, ok := d.GetOk("name_regex"); ok {
		r, err := regexp.Compile(v.(string))
		if err != nil {
			return WrapError(err)
		}
		commonNameRegex = r
	}

	var expectedCertType string
	if v, ok := d.GetOk("cert_type"); ok {
		certTypeMap := map[string]string{
			"root":       "ROOT",
			"subRoot":    "SUB_ROOT",
			"externalCa": "EXTERNAL_CA",
		}
		if expected, ok := certTypeMap[v.(string)]; ok {
			expectedCertType = expected
		}
	}

	var response map[string]interface{}
	var err error
	for {
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			response, err = client.RpcPost("cas", "2020-06-30", action, nil, request, true)
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
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_ssl_certificates_service_pca_certificates", action, AlibabaCloudSdkGoERROR)
		}

		resp, err := jsonpath.Get("$.CertificateList", response)
		if err != nil {
			return WrapErrorf(err, FailedGetAttributeMsg, action, "$.CertificateList", response)
		}

		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})
			if len(idsMap) > 0 {
				if _, ok := idsMap[fmt.Sprint(item["Identifier"])]; !ok {
					continue
				}
			}

			if commonNameRegex != nil && !commonNameRegex.MatchString(fmt.Sprint(item["CommonName"])) {
				continue
			}

			if expectedCertType != "" && fmt.Sprint(item["CertificateType"]) != expectedCertType {
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
			"id":                fmt.Sprint(object["Identifier"]),
			"identifier":        fmt.Sprint(object["Identifier"]),
			"serial_number":     object["SerialNumber"],
			"x509_certificate":  object["X509Certificate"],
			"certificate_type":  object["CertificateType"],
			"algorithm":         object["Algorithm"],
			"sign_algorithm":    object["SignAlgorithm"],
			"sha2":              object["Sha2"],
			"md5":               object["Md5"],
			"locality":          object["Locality"],
			"organization":      object["Organization"],
			"organization_unit": object["OrganizationUnit"],
			"common_name":       object["CommonName"],
			"country_code":      object["CountryCode"],
			"state":             object["State"],
			"parent_identifier": object["ParentIdentifier"],
			"status":            object["Status"],
			"years":             formatInt(object["Years"]),
			"before_date":       object["BeforeDate"],
			"after_date":        object["AfterDate"],
			"resource_group_id": object["ResourceGroupId"],
		}

		ids = append(ids, fmt.Sprint(mapping["id"]))
		names = append(names, object["CommonName"])
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
