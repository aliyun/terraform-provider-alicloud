package oss

import (
	"fmt"
	"log"
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

func ResourceAliCloudOssBucketHttpsConfig() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudOssBucketHttpsConfigCreate,
		Read:   resourceAliCloudOssBucketHttpsConfigRead,
		Update: resourceAliCloudOssBucketHttpsConfigUpdate,
		Delete: resourceAliCloudOssBucketHttpsConfigDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"bucket": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"cipher_suit": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"custom_cipher_suite": {
							Type:     schema.TypeSet,
							Optional: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"enable": {
							Type:     schema.TypeBool,
							Optional: true,
							Computed: true,
						},
						"tls13_custom_cipher_suite": {
							Type:     schema.TypeSet,
							Optional: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"strong_cipher_suite": {
							Type:     schema.TypeBool,
							Optional: true,
						},
					},
				},
			},
			"enable": {
				Type:     schema.TypeBool,
				Required: true,
			},
			"tls_versions": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func resourceAliCloudOssBucketHttpsConfigCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := fmt.Sprintf("/?httpsConfig")
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]*string)
	body := make(map[string]interface{})
	hostMap := make(map[string]*string)
	var err error
	request = make(map[string]interface{})
	hostMap["bucket"] = helper.StringPointer(d.Get("bucket").(string))

	httpsConfiguration := make(map[string]interface{})

	tLS := make(map[string]interface{})
	nodeNative1, _ := jsonpath.Get("$", d.Get("tls_versions"))
	tLS["TLSVersion"] = make([]interface{}, 0)
	if nodeNative1 != nil && nodeNative1 != "" {
		tLS["TLSVersion"] = nodeNative1.(*schema.Set).List()
	}
	tLS["Enable"] = d.Get("enable")
	httpsConfiguration["TLS"] = tLS

	if v := d.Get("cipher_suit"); !helper.IsNil(v) {
		cipherSuite := make(map[string]interface{})
		enable3, _ := jsonpath.Get("$[0].enable", d.Get("cipher_suit"))
		if enable3 != nil && enable3 != "" {
			cipherSuite["Enable"] = enable3
		}
		strongCipherSuite1, _ := jsonpath.Get("$[0].strong_cipher_suite", d.Get("cipher_suit"))
		if strongCipherSuite1 != nil && strongCipherSuite1 != "" {
			cipherSuite["StrongCipherSuite"] = strongCipherSuite1
		}
		customCipherSuite1, _ := jsonpath.Get("$[0].custom_cipher_suite", d.Get("cipher_suit"))
		if customCipherSuite1 != nil && customCipherSuite1 != "" {
			cipherSuite["CustomCipherSuite"] = convert.ConvertToInterfaceArray(customCipherSuite1)
		}
		tls13CustomCipherSuite, _ := jsonpath.Get("$[0].tls13_custom_cipher_suite", d.Get("cipher_suit"))
		if tls13CustomCipherSuite != nil && tls13CustomCipherSuite != "" {
			cipherSuite["TLS13CustomCipherSuite"] = convert.ConvertToInterfaceArray(tls13CustomCipherSuite)
		}

		if len(cipherSuite) > 0 {
			httpsConfiguration["CipherSuite"] = cipherSuite
		}
	}

	request["HttpsConfiguration"] = httpsConfiguration

	body = request
	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
		response, err = client.Do("Oss", helper.XmlParam("PUT", "2019-05-17", "PutBucketHttpsConfig", action), query, body, nil, hostMap, false)
		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	logs.AddDebug(action, response, request)

	if err != nil {
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, "alicloud_oss_bucket_https_config", action, errs.AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(*hostMap["bucket"]))

	ossServiceV2 := OssServiceV2{client}
	stateConf := helper.BuildStateConf([]string{}, []string{"#CHECKSET"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, ossServiceV2.OssBucketHttpsConfigStateRefreshFunc(d.Id(), "#TLS.Enable", []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return errs.WrapErrorf(err, errs.IdMsg, d.Id())
	}

	return resourceAliCloudOssBucketHttpsConfigRead(d, meta)
}

func resourceAliCloudOssBucketHttpsConfigRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	ossServiceV2 := OssServiceV2{client}

	objectRaw, err := ossServiceV2.DescribeOssBucketHttpsConfig(d.Id())
	if err != nil {
		if !d.IsNewResource() && errs.NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_oss_bucket_https_config DescribeOssBucketHttpsConfig Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return errs.WrapError(err)
	}

	tLSRawObj, _ := jsonpath.Get("$.TLS", objectRaw)
	tLSRaw := make(map[string]interface{})
	if tLSRawObj != nil {
		tLSRaw = tLSRawObj.(map[string]interface{})
	}
	d.Set("enable", tLSRaw["Enable"])

	cipherSuitMaps := make([]map[string]interface{}, 0)
	cipherSuitMap := make(map[string]interface{})
	cipherSuiteRaw := make(map[string]interface{})
	if objectRaw["CipherSuite"] != nil {
		cipherSuiteRaw = objectRaw["CipherSuite"].(map[string]interface{})
	}
	if len(cipherSuiteRaw) > 0 {
		cipherSuitMap["enable"] = cipherSuiteRaw["Enable"]
		cipherSuitMap["strong_cipher_suite"] = cipherSuiteRaw["StrongCipherSuite"]

		customCipherSuiteRaw := make([]interface{}, 0)
		if cipherSuiteRaw["CustomCipherSuite"] != nil {
			customCipherSuiteRaw = convert.ConvertToInterfaceArray(cipherSuiteRaw["CustomCipherSuite"])
		}

		cipherSuitMap["custom_cipher_suite"] = customCipherSuiteRaw
		tLS13CustomCipherSuiteRaw := make([]interface{}, 0)
		if cipherSuiteRaw["TLS13CustomCipherSuite"] != nil {
			tLS13CustomCipherSuiteRaw = convert.ConvertToInterfaceArray(cipherSuiteRaw["TLS13CustomCipherSuite"])
		}

		cipherSuitMap["tls13_custom_cipher_suite"] = tLS13CustomCipherSuiteRaw
		cipherSuitMaps = append(cipherSuitMaps, cipherSuitMap)
	}
	if err := d.Set("cipher_suit", cipherSuitMaps); err != nil {
		return err
	}
	tLSVersionRaw, _ := jsonpath.Get("$.TLS.TLSVersion", objectRaw)
	d.Set("tls_versions", tLSVersionRaw)

	d.Set("bucket", d.Id())

	return nil
}

func resourceAliCloudOssBucketHttpsConfigUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]*string
	var body map[string]interface{}
	update := false

	var err error
	action := fmt.Sprintf("/?httpsConfig")
	request = make(map[string]interface{})
	query = make(map[string]*string)
	body = make(map[string]interface{})
	hostMap := make(map[string]*string)
	hostMap["bucket"] = helper.StringPointer(d.Id())

	httpsConfiguration := make(map[string]interface{})

	if d.HasChanges("enable", "tls_versions") {
		update = true
	}
	if v := d.Get("enable"); !helper.IsNil(v) || d.HasChange("enable") {
		tLS := make(map[string]interface{})
		if v, ok := d.GetOkExists("enable"); ok {
			tLS["Enable"] = v
		}
		tlsVersions, _ := jsonpath.Get("$", d.Get("tls_versions"))
		if tlsVersions != nil && tlsVersions != "" {
			tLS["TLSVersion"] = convert.ConvertToInterfaceArray(tlsVersions)
		}

		if len(tLS) > 0 {
			httpsConfiguration["TLS"] = tLS
		}
	}

	if d.HasChange("cipher_suit") {
		update = true
	}
	if v := d.Get("cipher_suit"); !helper.IsNil(v) || d.HasChange("cipher_suit") {
		cipherSuite := make(map[string]interface{})
		enable3, _ := jsonpath.Get("$[0].enable", d.Get("cipher_suit"))
		if enable3 != nil && enable3 != "" {
			cipherSuite["Enable"] = enable3
		}
		strongCipherSuite1, _ := jsonpath.Get("$[0].strong_cipher_suite", d.Get("cipher_suit"))
		if strongCipherSuite1 != nil && strongCipherSuite1 != "" {
			cipherSuite["StrongCipherSuite"] = strongCipherSuite1
		}
		customCipherSuite1, _ := jsonpath.Get("$[0].custom_cipher_suite", d.Get("cipher_suit"))
		if customCipherSuite1 != nil && customCipherSuite1 != "" {
			cipherSuite["CustomCipherSuite"] = convert.ConvertToInterfaceArray(customCipherSuite1)
		}
		tls13CustomCipherSuite, _ := jsonpath.Get("$[0].tls13_custom_cipher_suite", d.Get("cipher_suit"))
		if tls13CustomCipherSuite != nil && tls13CustomCipherSuite != "" {
			cipherSuite["TLS13CustomCipherSuite"] = convert.ConvertToInterfaceArray(tls13CustomCipherSuite)
		}

		if len(cipherSuite) > 0 {
			httpsConfiguration["CipherSuite"] = cipherSuite
		}
	}

	request["HttpsConfiguration"] = httpsConfiguration

	body = request
	if update {
		wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
		err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
			response, err = client.Do("Oss", helper.XmlParam("PUT", "2019-05-17", "PutBucketHttpsConfig", action), query, body, nil, hostMap, false)
			if err != nil {
				if helper.NeedRetry(err) {
					wait()
					return retry.RetryableError(err)
				}
				return retry.NonRetryableError(err)
			}
			return nil
		})
		logs.AddDebug(action, response, request)
		if err != nil {
			return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
		}
		ossServiceV2 := OssServiceV2{client}
		stateConf := helper.BuildStateConf([]string{}, []string{"#CHECKSET"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, ossServiceV2.OssBucketHttpsConfigStateRefreshFunc(d.Id(), "#TLS.Enable", []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return errs.WrapErrorf(err, errs.IdMsg, d.Id())
		}
	}

	return resourceAliCloudOssBucketHttpsConfigRead(d, meta)
}

func resourceAliCloudOssBucketHttpsConfigDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[WARN] Cannot destroy resource AliCloud Resource Bucket Https Config. Terraform will remove this resource from the state file, however resources may remain.")
	return nil
}
