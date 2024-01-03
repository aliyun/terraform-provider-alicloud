// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"encoding/json"
	"fmt"
	"log"
	"regexp"
	"strconv"
	"time"

	"github.com/PaesslerAG/jsonpath"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAliCloudCdnDomain() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudCdnDomainCreate,
		Read:   resourceAliCloudCdnDomainRead,
		Update: resourceAliCloudCdnDomainUpdate,
		Delete: resourceAliCloudCdnDomainDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(15 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(15 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"cdn_type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: StringInSlice([]string{"unitcell", "web", "download", "liveStream", "video", "app"}, false),
			},
			"certificate_config": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"server_certificate": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"private_key": {
							Type:      schema.TypeString,
							Optional:  true,
							Sensitive: true,
						},
						"force_set": {
							Type:         schema.TypeString,
							Optional:     true,
							Removed:      "Field 'force_set' has been removed from provider version 1.206.0",
							ValidateFunc: StringInSlice([]string{"1", "0"}, false),
						},
						"cert_id": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"cert_region": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"cert_type": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"cert_name": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"server_certificate_status": {
							Type:         schema.TypeString,
							Optional:     true,
							Default:      "on",
							ValidateFunc: StringInSlice([]string{"on", "off"}, false),
						},
					},
				},
			},
			"check_url": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"cname": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"domain_name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: StringMatch(regexp.MustCompile("^([a-z0-9]([a-z0-9\\-]{0,61}[a-z0-9])?\\.)+[a-z][a-z0-9\\-]{0,62}$"), "Name of the accelerated domain. This name without suffix can have a string of 1 to 63 characters, must contain only alphanumeric characters or \"-\", and must not begin or end with \"-\", and \"-\" must not in the 3th and 4th character positions at the same time. Suffix .sh and .tel are not supported."),
			},
			"resource_group_id": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: StringMatch(regexp.MustCompile("rg-\\w+"), "The ID of the resource group."),
			},
			"scope": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				ValidateFunc: StringInSlice([]string{"domestic", "overseas", "global"}, false),
			},
			"sources": {
				Type:     schema.TypeSet,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: StringInSlice([]string{"domain", "ipaddr", "oss", "common"}, false),
						},
						"content": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"priority": {
							Type:         schema.TypeInt,
							Optional:     true,
							Default:      20,
							ValidateFunc: IntBetween(0, 100),
						},
						"port": {
							Type:         schema.TypeInt,
							Optional:     true,
							Default:      80,
							ValidateFunc: IntInSlice([]int{80, 443}),
						},
						"weight": {
							Type:         schema.TypeInt,
							Optional:     true,
							Default:      10,
							ValidateFunc: IntBetween(0, 100),
						},
					},
				},
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"tags": tagsSchema(),
		},
	}
}

func resourceAliCloudCdnDomainCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	action := "AddCdnDomain"
	var request map[string]interface{}
	var response map[string]interface{}
	conn, err := client.NewCdnClient()
	if err != nil {
		return WrapError(err)
	}
	request = make(map[string]interface{})
	request["DomainName"] = d.Get("domain_name")

	request["CdnType"] = d.Get("cdn_type")
	if v, ok := d.GetOk("resource_group_id"); ok {
		request["ResourceGroupId"] = v
	}
	if v, ok := d.GetOk("sources"); ok {
		sourcesMaps := make([]map[string]interface{}, 0)
		for _, dataLoop := range v.(*schema.Set).List() {
			dataLoopTmp := dataLoop.(map[string]interface{})
			dataLoopMap := make(map[string]interface{})
			dataLoopMap["Weight"] = dataLoopTmp["weight"]
			dataLoopMap["type"] = dataLoopTmp["type"]
			dataLoopMap["content"] = dataLoopTmp["content"]
			dataLoopMap["priority"] = dataLoopTmp["priority"]
			dataLoopMap["port"] = dataLoopTmp["port"]
			sourcesMaps = append(sourcesMaps, dataLoopMap)
		}
		sourcesMapsJson, err := json.Marshal(sourcesMaps)
		if err != nil {
			return WrapError(err)
		}
		request["Sources"] = string(sourcesMapsJson)
	}

	if v, ok := d.GetOk("scope"); ok {
		request["Scope"] = v
	}
	if v, ok := d.GetOk("check_url"); ok {
		request["CheckUrl"] = v
	}
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2018-05-10"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})

		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(action, response, request)
		return nil
	})

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_cdn_domain_new", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(request["DomainName"]))

	cdnServiceV2 := CdnServiceV2{client}
	stateConf := BuildStateConf([]string{}, []string{"online"}, d.Timeout(schema.TimeoutCreate), 50*time.Second, cdnServiceV2.CdnDomainStateRefreshFunc(d.Id(), "DomainStatus", []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return resourceAliCloudCdnDomainUpdate(d, meta)
}

func resourceAliCloudCdnDomainRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	cdnServiceV2 := CdnServiceV2{client}

	objectRaw, err := cdnServiceV2.DescribeCdnDomain(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_cdn_domain_new DescribeCdnDomain Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("cdn_type", objectRaw["CdnType"])
	d.Set("cname", objectRaw["Cname"])
	d.Set("resource_group_id", objectRaw["ResourceGroupId"])
	d.Set("scope", objectRaw["Scope"])
	d.Set("status", objectRaw["DomainStatus"])
	d.Set("domain_name", objectRaw["DomainName"])
	sourceModel1Raw, _ := jsonpath.Get("$.SourceModels.SourceModel", objectRaw)
	sourcesMaps := make([]map[string]interface{}, 0)
	if sourceModel1Raw != nil {
		for _, sourceModelChild1Raw := range sourceModel1Raw.([]interface{}) {
			sourcesMap := make(map[string]interface{})
			sourceModelChild1Raw := sourceModelChild1Raw.(map[string]interface{})
			sourcesMap["content"] = sourceModelChild1Raw["Content"]
			sourcesMap["port"], _ = sourceModelChild1Raw["Port"]
			if v, ok := sourceModelChild1Raw["Priority"].(string); ok {
				sourcesMap["priority"], _ = strconv.Atoi(v)
			} else {
				sourcesMap["priority"] = 0
			}
			sourcesMap["type"] = sourceModelChild1Raw["Type"]
			if v, ok := sourceModelChild1Raw["Weight"].(string); ok {
				sourcesMap["weight"], _ = strconv.Atoi(v)
			} else {
				sourcesMap["weight"] = 0
			}
			sourcesMaps = append(sourcesMaps, sourcesMap)
		}
	}
	d.Set("sources", sourcesMaps)

	objectRaw, err = cdnServiceV2.DescribeDescribeDomainCertificateInfo(d.Id())
	if err != nil {
		return WrapError(err)
	}

	d.Set("domain_name", objectRaw["DomainName"])
	certificateConfigMaps := make([]map[string]interface{}, 0)
	certificateConfigMap := make(map[string]interface{})

	certificateConfigMap["cert_id"] = objectRaw["CertId"]
	certificateConfigMap["cert_name"] = objectRaw["CertName"]
	certificateConfigMap["cert_region"] = objectRaw["CertRegion"]
	certificateConfigMap["cert_type"] = objectRaw["CertType"]
	certificateConfigMap["server_certificate"] = objectRaw["ServerCertificate"]
	certificateConfigMap["server_certificate_status"] = objectRaw["ServerCertificateStatus"]

	if v, ok := d.GetOk("certificate_config"); ok {
		oldConfig := v.([]interface{})
		if len(oldConfig) > 0 {
			val := oldConfig[0].(map[string]interface{})
			certificateConfigMap["private_key"] = val["private_key"]
		}
	}
	certificateConfigMaps = append(certificateConfigMaps, certificateConfigMap)
	d.Set("certificate_config", certificateConfigMaps)

	objectRaw, err = cdnServiceV2.DescribeListTagResources(d.Id())
	if err != nil {
		return WrapError(err)
	}

	tagsMaps, _ := jsonpath.Get("$.TagResources.TagResource", objectRaw)
	d.Set("tags", tagsToMap(tagsMaps))

	return nil
}

func resourceAliCloudCdnDomainUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	update := false
	d.Partial(true)
	action := "ModifyCdnDomain"
	conn, err := client.NewCdnClient()
	if err != nil {
		return WrapError(err)
	}
	request = make(map[string]interface{})

	request["DomainName"] = d.Id()
	if !d.IsNewResource() && d.HasChange("resource_group_id") {
		update = true
		request["ResourceGroupId"] = d.Get("resource_group_id")
	}

	if !d.IsNewResource() && d.HasChange("sources") {
		update = true
		if v, ok := d.GetOk("sources"); ok {
			sourcesMaps := make([]map[string]interface{}, 0)
			for _, dataLoop := range v.(*schema.Set).List() {
				dataLoopTmp := dataLoop.(map[string]interface{})
				dataLoopMap := make(map[string]interface{})
				dataLoopMap["Type"] = dataLoopTmp["type"]
				dataLoopMap["Content"] = dataLoopTmp["content"]
				dataLoopMap["Priority"] = dataLoopTmp["priority"]
				dataLoopMap["Port"] = dataLoopTmp["port"]
				dataLoopMap["Weight"] = dataLoopTmp["weight"]
				sourcesMaps = append(sourcesMaps, dataLoopMap)
			}
			sourcesMapsJson, err := json.Marshal(sourcesMaps)
			if err != nil {
				return WrapError(err)
			}
			request["Sources"] = string(sourcesMapsJson)
		}
	}

	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2018-05-10"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})

			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			addDebug(action, response, request)
			return nil
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
		cdnServiceV2 := CdnServiceV2{client}
		stateConf := BuildStateConf([]string{}, []string{"online"}, d.Timeout(schema.TimeoutUpdate), 60*time.Second, cdnServiceV2.CdnDomainStateRefreshFunc(d.Id(), "DomainStatus", []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
		d.SetPartial("resource_group_id")
	}
	update = false
	action = "SetCdnDomainSSLCertificate"
	conn, err = client.NewCdnClient()
	if err != nil {
		return WrapError(err)
	}
	request = make(map[string]interface{})

	request["DomainName"] = d.Id()
	if d.HasChange("certificate_config.0.server_certificate_status") {
		update = true
	}
	if v, ok := d.GetOk("certificate_config.0.server_certificate_status"); ok {
		request["SSLProtocol"] = v
	}

	if d.HasChange("certificate_config.0.cert_name") {
		update = true
		if v, ok := d.GetOk("certificate_config.0.cert_name"); ok {
			request["CertName"] = v
		}
	}

	if d.HasChange("certificate_config.0.cert_id") {
		update = true
		if v, ok := d.GetOk("certificate_config.0.cert_id"); ok {
			request["CertId"] = v
		}
	}

	if d.HasChange("certificate_config.0.cert_type") {
		update = true
		if v, ok := d.GetOk("certificate_config.0.cert_type"); ok {
			request["CertType"] = v
		}
	}

	if d.HasChange("certificate_config.0.server_certificate") {
		update = true
		if v, ok := d.GetOk("certificate_config.0.server_certificate"); ok {
			request["SSLPub"] = v
		}
	}

	if d.HasChange("certificate_config.0.private_key") {
		update = true
		if v, ok := d.GetOk("certificate_config.0.private_key"); ok {
			request["SSLPri"] = v
		}
	}

	if d.HasChange("certificate_config.0.cert_region") {
		update = true
		if v, ok := d.GetOk("certificate_config.0.cert_region"); ok {
			request["CertRegion"] = v
		}
	}

	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2018-05-10"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})

			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			addDebug(action, response, request)
			return nil
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
		d.SetPartial("certificate_config")
	}

	update = false
	if d.HasChange("tags") {
		update = true
		cdnServiceV2 := CdnServiceV2{client}
		if err := cdnServiceV2.SetResourceTags(d, "DOMAIN"); err != nil {
			return WrapError(err)
		}
		d.SetPartial("tags")
	}
	d.Partial(false)
	return resourceAliCloudCdnDomainRead(d, meta)
}

func resourceAliCloudCdnDomainDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "DeleteCdnDomain"
	var request map[string]interface{}
	var response map[string]interface{}
	conn, err := client.NewCdnClient()
	if err != nil {
		return WrapError(err)
	}
	request = make(map[string]interface{})

	request["DomainName"] = d.Id()

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2018-05-10"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})

		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(action, response, request)
		return nil
	})

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	cdnServiceV2 := CdnServiceV2{client}
	stateConf := BuildStateConf([]string{}, []string{}, d.Timeout(schema.TimeoutDelete), 4*time.Minute, cdnServiceV2.CdnDomainStateRefreshFunc(d.Id(), "DomainStatus", []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}
	return nil
}
