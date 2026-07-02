// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"log"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAliCloudApigDomain() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudApigDomainCreate,
		Read:   resourceAliCloudApigDomainRead,
		Update: resourceAliCloudApigDomainUpdate,
		Delete: resourceAliCloudApigDomainDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"ca_cert_identifier": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"cert_identifier": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"client_ca_cert": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"domain_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"domain_scope": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"force_https": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"gateway_type": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"http2_option": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"m_tls_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"protocol": {
				Type:     schema.TypeString,
				Required: true,
			},
			"resource_group_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"tls_cipher_suites_config": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"tls_cipher_suite": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"support_versions": {
										Type:     schema.TypeList,
										Optional: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
									"name": {
										Type:     schema.TypeString,
										Optional: true,
									},
								},
							},
						},
						"config_type": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"tls_max": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"tls_min": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceAliCloudApigDomainCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := fmt.Sprintf("/v1/domains")
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]*string)
	body := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})

	if v, ok := d.GetOk("resource_group_id"); ok {
		request["resourceGroupId"] = v
	}
	if v, ok := d.GetOkExists("force_https"); ok {
		request["forceHttps"] = v
	}
	request["name"] = d.Get("domain_name")
	request["protocol"] = d.Get("protocol")
	if v, ok := d.GetOk("tls_min"); ok {
		request["tlsMin"] = v
	}
	if v, ok := d.GetOk("client_ca_cert"); ok {
		request["clientCACert"] = v
	}
	tlsCipherSuitesConfig := make(map[string]interface{})

	if v := d.Get("tls_cipher_suites_config"); !IsNil(v) {
		if v, ok := d.GetOk("tls_cipher_suites_config"); ok {
			localData, err := jsonpath.Get("$[0].tls_cipher_suite", v)
			if err != nil {
				localData = make([]interface{}, 0)
			}
			localMaps := make([]interface{}, 0)
			for _, dataLoop := range convertToInterfaceArray(localData) {
				dataLoopTmp := make(map[string]interface{})
				if dataLoop != nil {
					dataLoopTmp = dataLoop.(map[string]interface{})
				}
				dataLoopMap := make(map[string]interface{})
				dataLoopMap["supportVersions"] = dataLoopTmp["support_versions"]
				dataLoopMap["name"] = dataLoopTmp["name"]
				localMaps = append(localMaps, dataLoopMap)
			}
			tlsCipherSuitesConfig["tlsCipherSuite"] = localMaps
		}

		configType1, _ := jsonpath.Get("$[0].config_type", v)
		if configType1 != nil && configType1 != "" {
			tlsCipherSuitesConfig["configType"] = configType1
		}

		request["tlsCipherSuitesConfig"] = tlsCipherSuitesConfig
	}

	if v, ok := d.GetOk("http2_option"); ok {
		request["http2Option"] = v
	}
	if v, ok := d.GetOk("domain_scope"); ok {
		request["domainScope"] = v
	}
	if v, ok := d.GetOk("cert_identifier"); ok {
		request["certIdentifier"] = v
	}
	if v, ok := d.GetOk("tls_max"); ok {
		request["tlsMax"] = v
	}
	if v, ok := d.GetOk("gateway_type"); ok {
		request["gatewayType"] = v
	}
	if v, ok := d.GetOk("ca_cert_identifier"); ok {
		request["caCertIdentifier"] = v
	}
	if v, ok := d.GetOkExists("m_tls_enabled"); ok {
		request["mTLSEnabled"] = v
	}
	body = request
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.RoaPost("APIG", "2024-03-27", action, query, nil, body, true)
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_apig_domain", action, AlibabaCloudSdkGoERROR)
	}

	id, _ := jsonpath.Get("$.data.domainId", response)
	d.SetId(fmt.Sprint(id))

	return resourceAliCloudApigDomainRead(d, meta)
}

func resourceAliCloudApigDomainRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	apigServiceV2 := ApigServiceV2{client}

	objectRaw, err := apigServiceV2.DescribeApigDomain(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_apig_domain DescribeApigDomain Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("ca_cert_identifier", objectRaw["caCertIdentifier"])
	d.Set("cert_identifier", objectRaw["certIdentifier"])
	d.Set("client_ca_cert", objectRaw["clientCACert"])
	d.Set("domain_name", objectRaw["name"])
	d.Set("domain_scope", objectRaw["domainScope"])
	d.Set("force_https", objectRaw["forceHttps"])
	d.Set("http2_option", objectRaw["http2Option"])
	d.Set("m_tls_enabled", objectRaw["mTLSEnabled"])
	d.Set("protocol", objectRaw["protocol"])
	d.Set("resource_group_id", objectRaw["resourceGroupId"])
	d.Set("tls_max", objectRaw["tlsMax"])
	d.Set("tls_min", objectRaw["tlsMin"])

	tlsCipherSuitesConfigMaps := make([]map[string]interface{}, 0)
	tlsCipherSuitesConfigMap := make(map[string]interface{})
	tlsCipherSuitesConfigRaw := make(map[string]interface{})
	if objectRaw["tlsCipherSuitesConfig"] != nil {
		tlsCipherSuitesConfigRaw = objectRaw["tlsCipherSuitesConfig"].(map[string]interface{})
	}
	if len(tlsCipherSuitesConfigRaw) > 0 {
		tlsCipherSuitesConfigMap["config_type"] = tlsCipherSuitesConfigRaw["configType"]

		tlsCipherSuiteRaw := tlsCipherSuitesConfigRaw["tlsCipherSuite"]
		tlsCipherSuiteMaps := make([]map[string]interface{}, 0)
		if tlsCipherSuiteRaw != nil {
			for _, tlsCipherSuiteChildRaw := range convertToInterfaceArray(tlsCipherSuiteRaw) {
				tlsCipherSuiteMap := make(map[string]interface{})
				tlsCipherSuiteChildRaw := tlsCipherSuiteChildRaw.(map[string]interface{})
				tlsCipherSuiteMap["name"] = tlsCipherSuiteChildRaw["name"]

				supportVersionsRaw := make([]interface{}, 0)
				if tlsCipherSuiteChildRaw["supportVersions"] != nil {
					supportVersionsRaw = convertToInterfaceArray(tlsCipherSuiteChildRaw["supportVersions"])
				}

				tlsCipherSuiteMap["support_versions"] = supportVersionsRaw
				tlsCipherSuiteMaps = append(tlsCipherSuiteMaps, tlsCipherSuiteMap)
			}
		}
		tlsCipherSuitesConfigMap["tls_cipher_suite"] = tlsCipherSuiteMaps
		tlsCipherSuitesConfigMaps = append(tlsCipherSuitesConfigMaps, tlsCipherSuitesConfigMap)
	}
	if err := d.Set("tls_cipher_suites_config", tlsCipherSuitesConfigMaps); err != nil {
		return err
	}

	return nil
}

func resourceAliCloudApigDomainUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]*string
	var body map[string]interface{}
	update := false

	var err error
	domainId := d.Id()
	action := fmt.Sprintf("/v1/domains/%s", domainId)
	request = make(map[string]interface{})
	query = make(map[string]*string)
	body = make(map[string]interface{})

	if d.HasChange("force_https") {
		update = true
	}
	if v, ok := d.GetOkExists("force_https"); ok || d.HasChange("force_https") {
		request["forceHttps"] = v
	}
	if d.HasChange("protocol") {
		update = true
	}
	request["protocol"] = d.Get("protocol")
	if d.HasChange("tls_min") {
		update = true
	}
	if v, ok := d.GetOk("tls_min"); ok || d.HasChange("tls_min") {
		request["tlsMin"] = v
	}
	if d.HasChange("client_ca_cert") {
		update = true
	}
	if v, ok := d.GetOk("client_ca_cert"); ok || d.HasChange("client_ca_cert") {
		request["clientCACert"] = v
	}
	if d.HasChange("tls_cipher_suites_config") {
		update = true
	}
	tlsCipherSuitesConfig := make(map[string]interface{})

	if v := d.Get("tls_cipher_suites_config"); !IsNil(v) || d.HasChange("tls_cipher_suites_config") {
		if v, ok := d.GetOk("tls_cipher_suites_config"); ok {
			localData, err := jsonpath.Get("$[0].tls_cipher_suite", v)
			if err != nil {
				localData = make([]interface{}, 0)
			}
			localMaps := make([]interface{}, 0)
			for _, dataLoop := range convertToInterfaceArray(localData) {
				dataLoopTmp := make(map[string]interface{})
				if dataLoop != nil {
					dataLoopTmp = dataLoop.(map[string]interface{})
				}
				dataLoopMap := make(map[string]interface{})
				dataLoopMap["supportVersions"] = dataLoopTmp["support_versions"]
				dataLoopMap["name"] = dataLoopTmp["name"]
				localMaps = append(localMaps, dataLoopMap)
			}
			tlsCipherSuitesConfig["tlsCipherSuite"] = localMaps
		}

		configType1, _ := jsonpath.Get("$[0].config_type", v)
		if configType1 != nil && configType1 != "" {
			tlsCipherSuitesConfig["configType"] = configType1
		}

		request["tlsCipherSuitesConfig"] = tlsCipherSuitesConfig
	}

	if d.HasChange("http2_option") {
		update = true
	}
	if v, ok := d.GetOk("http2_option"); ok || d.HasChange("http2_option") {
		request["http2Option"] = v
	}
	if d.HasChange("domain_scope") {
		update = true
	}
	if v, ok := d.GetOk("domain_scope"); ok || d.HasChange("domain_scope") {
		request["domainScope"] = v
	}
	if d.HasChange("cert_identifier") {
		update = true
	}
	if v, ok := d.GetOk("cert_identifier"); ok || d.HasChange("cert_identifier") {
		request["certIdentifier"] = v
	}
	if d.HasChange("tls_max") {
		update = true
	}
	if v, ok := d.GetOk("tls_max"); ok || d.HasChange("tls_max") {
		request["tlsMax"] = v
	}
	if d.HasChange("ca_cert_identifier") {
		update = true
	}
	if v, ok := d.GetOk("ca_cert_identifier"); ok || d.HasChange("ca_cert_identifier") {
		request["caCertIdentifier"] = v
	}
	if d.HasChange("m_tls_enabled") {
		update = true
	}
	if v, ok := d.GetOkExists("m_tls_enabled"); ok || d.HasChange("m_tls_enabled") {
		request["mTLSEnabled"] = v
	}
	body = request
	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RoaPut("APIG", "2024-03-27", action, query, nil, body, true)
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
	}
	update = false
	action = fmt.Sprintf("/move-resource-group")
	request = make(map[string]interface{})
	query = make(map[string]*string)
	body = make(map[string]interface{})
	query["ResourceId"] = StringPointer(d.Id())
	query["RegionId"] = StringPointer(client.RegionId)
	if d.HasChange("resource_group_id") {
		update = true
	}
	if v, ok := d.GetOk("resource_group_id"); ok {
		query["ResourceGroupId"] = StringPointer(v.(string))
	}

	query["Service"] = StringPointer("APIG")
	query["ResourceType"] = StringPointer("Domain")
	body = request
	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RoaPost("APIG", "2024-03-27", action, query, nil, body, true)
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
	}

	return resourceAliCloudApigDomainRead(d, meta)
}

func resourceAliCloudApigDomainDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	domainId := d.Id()
	action := fmt.Sprintf("/v1/domains/%s", domainId)
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]*string)
	var err error
	request = make(map[string]interface{})

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = client.RoaDelete("APIG", "2024-03-27", action, query, nil, nil, true)
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
		if NotFoundError(err) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	return nil
}
