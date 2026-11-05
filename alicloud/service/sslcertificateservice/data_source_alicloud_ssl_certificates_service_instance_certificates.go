// Package sslcertificateservice. This file is generated automatically. Please do not modify it manually, thank you!
package sslcertificateservice

import (
	"fmt"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/convert"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAliCloudSslCertificatesServiceInstanceCertificates() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAliCloudSslCertificatesServiceInstanceCertificateRead,
		Schema: map[string]*schema.Schema{
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"certificate_source": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"certificate_status": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"instance_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"certificates": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"algorithm": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"cert_identifier": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"certificate_id": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"certificate_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"certificate_source": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"certificate_status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"common_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"domain": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"exist_private_key": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"finger_print": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"instance_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"issuer": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"key_size": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"not_after": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"not_before": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"id": {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func dataSourceAliCloudSslCertificatesServiceInstanceCertificateRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

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

	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	action := "ListCertificates"
	var err error
	request = make(map[string]interface{})
	query = make(map[string]interface{})

	if v, ok := d.GetOk("certificate_source"); ok {
		request["CertificateSource"] = v
	}
	if v, ok := d.GetOk("certificate_status"); ok {
		request["CertificateStatus"] = v
	}
	if v, ok := d.GetOk("instance_id"); ok {
		request["InstanceId"] = v
	}
	request["ShowSize"] = helper.PageSizeLarge
	request["CurrentPage"] = 1
	for {
		wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
		err = retry.Retry(d.Timeout(schema.TimeoutRead), func() *retry.RetryError {
			response, err = client.RpcPost("cas", "2020-04-07", action, query, request, true)

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

		resp, _ := jsonpath.Get("$.CertificateList[*]", response)

		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})
			if len(idsMap) > 0 {
				if _, ok := idsMap[fmt.Sprint(item["CertificateId"])]; !ok {
					continue
				}
			}
			objects = append(objects, item)
		}

		if len(result) < helper.PageSizeLarge {
			break
		}
		request["CurrentPage"] = request["CurrentPage"].(int) + 1
	}

	ids := make([]string, 0)
	s := make([]map[string]interface{}, 0)
	for _, objectRaw := range objects {
		mapping := map[string]interface{}{}

		// ListCertificates types CertificateId as a string while every other API in this
		// product returns an integer; normalise it so the declared type holds.
		actualId, err := convert.FormatInt(objectRaw["CertificateId"])
		if err != nil {
			return err
		}
		mapping["id"] = actualId

		mapping["algorithm"] = objectRaw["Algorithm"]
		mapping["cert_identifier"] = objectRaw["CertIdentifier"]
		certificateId, err := convert.FormatInt(objectRaw["CertificateId"])
		if err != nil {
			return err
		}
		mapping["certificate_id"] = certificateId
		mapping["certificate_name"] = objectRaw["CertificateName"]
		mapping["certificate_source"] = objectRaw["CertificateSource"]
		mapping["certificate_status"] = objectRaw["CertificateStatus"]
		mapping["common_name"] = objectRaw["CommonName"]
		mapping["domain"] = objectRaw["Domain"]
		mapping["exist_private_key"] = objectRaw["ExistPrivateKey"]
		mapping["finger_print"] = objectRaw["FingerPrint"]
		mapping["instance_id"] = objectRaw["InstanceId"]
		mapping["issuer"] = objectRaw["Issuer"]
		mapping["key_size"] = objectRaw["KeySize"]
		mapping["not_after"] = objectRaw["NotAfter"]
		mapping["not_before"] = objectRaw["NotBefore"]
		// ListCertificates does not return the certificate serial; it is deliberately not declared
		// here rather than declared and always empty. The instance_certificate resource reads it
		// through GetCertificateDetail, which does return it.

		ids = append(ids, fmt.Sprint(mapping["id"]))
		s = append(s, mapping)
	}

	d.SetId(helper.DataResourceIdHash(ids))
	if err := d.Set("ids", ids); err != nil {
		return errs.WrapError(err)
	}

	if err := d.Set("certificates", s); err != nil {
		return errs.WrapError(err)
	}

	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		helper.WriteToFile(output.(string), s)
	}
	return nil
}
