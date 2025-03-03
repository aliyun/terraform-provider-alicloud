package alicloud

import (
	"encoding/json"
	"fmt"
	"github.com/PaesslerAG/jsonpath"
	"github.com/tidwall/sjson"
	"log"
	"strings"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAliCloudOssBucketCname() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudOssBucketCnameCreate,
		Read:   resourceAliCloudOssBucketCnameRead,
		Update: resourceAliCloudOssBucketCnameUpdate,
		Delete: resourceAliCloudOssBucketCnameDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"bucket": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"certificate": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"creation_date": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"valid_end_date": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"fingerprint": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"valid_start_date": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"private_key": {
							Type:      schema.TypeString,
							Optional:  true,
							Sensitive: true,
						},
						"cert_id": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"certificate": {
							Type:      schema.TypeString,
							Optional:  true,
							Sensitive: true,
						},
					},
				},
			},
			"delete_certificate": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"domain": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"force": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"previous_cert_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceAliCloudOssBucketCnameCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := fmt.Sprintf("/?cname&comp=add")
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]*string)
	body := make(map[string]interface{})
	hostMap := make(map[string]*string)
	var err error
	request = make(map[string]interface{})
	hostMap["bucket"] = StringPointer(d.Get("bucket").(string))
	jsonString := "{}"
	jsonString, _ = sjson.Set(jsonString, "BucketCnameConfiguration.Cname.Domain", d.Get("domain"))
	err = json.Unmarshal([]byte(jsonString), &request)
	if err != nil {
		return WrapError(err)
	}

	objectDataLocalMap := make(map[string]interface{})

	if v := d.Get("domain"); !IsNil(v) {
		cname := make(map[string]interface{})
		cname["Domain"] = v
		certificateConfiguration := make(map[string]interface{})
		certId1, _ := jsonpath.Get("$[0].cert_id", d.Get("certificate"))
		if certId1 != nil && certId1 != "" {
			certificateConfiguration["CertId"] = certId1
		}
		certificate1, _ := jsonpath.Get("$[0].certificate", d.Get("certificate"))
		if certificate1 != nil && certificate1 != "" {
			certificateConfiguration["Certificate"] = certificate1
		}
		privateKey1, _ := jsonpath.Get("$[0].private_key", d.Get("certificate"))
		if privateKey1 != nil && privateKey1 != "" {
			certificateConfiguration["PrivateKey"] = privateKey1
		}

		cname["CertificateConfiguration"] = certificateConfiguration

		objectDataLocalMap["Cname"] = cname

		request["BucketCnameConfiguration"] = objectDataLocalMap
	}

	body = request
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.Do("Oss", genXmlParam("POST", "2019-05-17", "PutCname", action), query, body, nil, hostMap, false)
		if err != nil {
			if NeedRetry(err) || IsExpectedErrors(err, []string{"NeedVerifyDomainOwnership"}) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_oss_bucket_cname", action, AlibabaCloudSdkGoERROR)
	}

	BucketCnameConfigurationCnameDomainVar, _ := jsonpath.Get("BucketCnameConfiguration.Cname.Domain", request)
	d.SetId(fmt.Sprintf("%v:%v", *hostMap["bucket"], BucketCnameConfigurationCnameDomainVar))

	ossServiceV2 := OssServiceV2{client}
	stateConf := BuildStateConf([]string{}, []string{fmt.Sprint(d.Get("domain"))}, d.Timeout(schema.TimeoutCreate), 5*time.Second, ossServiceV2.OssBucketCnameStateRefreshFunc(d.Id(), "$.Domain", []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return resourceAliCloudOssBucketCnameRead(d, meta)
}

func resourceAliCloudOssBucketCnameRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	ossServiceV2 := OssServiceV2{client}

	objectRaw, err := ossServiceV2.DescribeOssBucketCname(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_oss_bucket_cname DescribeOssBucketCname Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	if objectRaw["Status"] != nil {
		d.Set("status", objectRaw["Status"])
	}
	if objectRaw["Domain"] != nil {
		d.Set("domain", objectRaw["Domain"])
	}

	certificateMaps := make([]map[string]interface{}, 0)
	certificateMap := make(map[string]interface{})
	certificate1RawObj, _ := jsonpath.Get("$.Certificate", objectRaw)
	certificate1Raw := make(map[string]interface{})
	if certificate1RawObj != nil {
		certificate1Raw = certificate1RawObj.(map[string]interface{})
	}
	if len(certificate1Raw) > 0 {
		certificateMap["cert_id"] = certificate1Raw["CertId"]
		certificateMap["creation_date"] = certificate1Raw["CreationDate"]
		certificateMap["fingerprint"] = certificate1Raw["Fingerprint"]
		certificateMap["status"] = certificate1Raw["Status"]
		certificateMap["type"] = certificate1Raw["Type"]
		certificateMap["valid_end_date"] = certificate1Raw["ValidEndDate"]
		certificateMap["valid_start_date"] = certificate1Raw["ValidStartDate"]

		if _, ok := d.GetOk("certificate"); ok {
			ov := d.Get("certificate")
			oldConfig := ov.([]interface{})
			if oldConfig != nil && len(oldConfig) > 0 {
				val := oldConfig[0].(map[string]interface{})
				certificateMap["private_key"] = val["private_key"]
				certificateMap["certificate"] = val["certificate"]
			}
		}

		certificateMaps = append(certificateMaps, certificateMap)
	}
	if certificate1RawObj != nil {
		if err := d.Set("certificate", certificateMaps); err != nil {
			return err
		}
	}

	parts := strings.Split(d.Id(), ":")
	d.Set("bucket", parts[0])

	return nil
}

func resourceAliCloudOssBucketCnameUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]*string
	var body map[string]interface{}
	update := false
	parts := strings.Split(d.Id(), ":")
	action := fmt.Sprintf("/?cname&comp=add")
	var err error
	request = make(map[string]interface{})
	query = make(map[string]*string)
	body = make(map[string]interface{})
	hostMap := make(map[string]*string)
	hostMap["bucket"] = StringPointer(parts[0])
	jsonString := "{}"
	jsonString, _ = sjson.Set(jsonString, "BucketCnameConfiguration.Cname.Domain", parts[1])
	err = json.Unmarshal([]byte(jsonString), &request)
	if err != nil {
		return WrapError(err)
	}

	objectDataLocalMap := make(map[string]interface{})

	if v := d.Get("domain"); !IsNil(v) {
		cname := make(map[string]interface{})
		cname["Domain"] = d.Get("domain")
		certificateConfiguration := make(map[string]interface{})
		certificateConfiguration["PreviousCertId"] = d.Get("previous_cert_id")
		certificateConfiguration["Force"] = d.Get("force")
		certificateConfiguration["DeleteCertificate"] = d.Get("delete_certificate")
		certId1, _ := jsonpath.Get("$[0].cert_id", v)
		if certId1 != nil && (d.HasChange("certificate.0.cert_id") || certId1 != "") {
			certificateConfiguration["CertId"] = certId1
		}
		certificate1, _ := jsonpath.Get("$[0].certificate", v)
		if certificate1 != nil && (d.HasChange("certificate.0.certificate") || certificate1 != "") {
			certificateConfiguration["Certificate"] = certificate1
		}
		privateKey1, _ := jsonpath.Get("$[0].private_key", v)
		if privateKey1 != nil && (d.HasChange("certificate.0.private_key") || privateKey1 != "") {
			certificateConfiguration["PrivateKey"] = privateKey1
		}

		cname["CertificateConfiguration"] = certificateConfiguration

		objectDataLocalMap["Cname"] = cname

		request["BucketCnameConfiguration"] = objectDataLocalMap
	}
	if d.HasChanges("certificate", "previous_cert_id") {
		update = true
	}

	body = request
	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.Do("Oss", genXmlParam("POST", "2019-05-17", "PutCname", action), query, body, nil, hostMap, false)
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
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
		ossServiceV2 := OssServiceV2{client}
		stateConf := BuildStateConf([]string{}, []string{fmt.Sprint(d.Get("domain"))}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, ossServiceV2.OssBucketCnameStateRefreshFunc(d.Id(), "$.Domain", []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
	}

	return resourceAliCloudOssBucketCnameRead(d, meta)
}

func resourceAliCloudOssBucketCnameDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	parts := strings.Split(d.Id(), ":")
	action := fmt.Sprintf("/?cname&comp=delete")
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]*string)
	body := make(map[string]interface{})
	hostMap := make(map[string]*string)
	var err error
	request = make(map[string]interface{})
	hostMap["bucket"] = StringPointer(parts[0])
	jsonString := "{}"
	jsonString, _ = sjson.Set(jsonString, "BucketCnameConfiguration.Cname.Domain", parts[1])
	err = json.Unmarshal([]byte(jsonString), &request)
	if err != nil {
		return WrapError(err)
	}

	body = request
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = client.Do("Oss", genXmlParam("POST", "2019-05-17", "DeleteCname", action), query, body, nil, hostMap, false)
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
		if IsExpectedErrors(err, []string{"NoSuchBucket", "NoSuchCname"}) || NotFoundError(err) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	ossServiceV2 := OssServiceV2{client}
	stateConf := BuildStateConf([]string{}, []string{}, d.Timeout(schema.TimeoutDelete), 5*time.Second, ossServiceV2.OssBucketCnameStateRefreshFunc(d.Id(), "$.Domain", []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return nil
}
