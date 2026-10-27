package vod

import (
	"fmt"
	"log"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/convert"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/tags"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/validators"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceDomain() *schema.Resource {
	return &schema.Resource{
		Create: resourceDomainCreate,
		Read:   resourceDomainRead,
		Update: resourceDomainUpdate,
		Delete: resourceDomainDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(20 * time.Minute),
			Update: schema.DefaultTimeout(20 * time.Minute),
			Delete: schema.DefaultTimeout(20 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"domain_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"scope": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Default:      "domestic",
				ValidateFunc: validators.StringInSlice([]string{"domestic", "overseas", "global"}, false),
			},
			"ssl_protocol": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validators.StringInSlice([]string{"on", "off"}, false),
			},
			"cert_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"cert_name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"cert_type": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"cert_region": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validators.StringInSlice([]string{"cn-hangzhou", "ap-southeast-1"}, false),
			},
			"ssl_pub": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"ssl_pri": {
				Type:      schema.TypeString,
				Optional:  true,
				Sensitive: true,
			},
			"check_url": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"top_level_domain": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"tags": tags.TagsSchemaWithElements(),
			"sources": {
				Type:     schema.TypeSet,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"source_type": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validators.StringInSlice([]string{"ipaddr", "domain", "oss"}, false),
						},
						"source_content": {
							Type:     schema.TypeString,
							Required: true,
						},
						"source_port": {
							Type:     schema.TypeString,
							Required: true,
						},
						"source_priority": {
							Type:         schema.TypeString,
							Optional:     true,
							Computed:     true,
							ValidateFunc: validators.StringInSlice([]string{"20", "30"}, false),
						},
					},
				},
			},
			"cname": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"weight": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"gmt_created": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"gmt_modified": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceDomainCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	vodService := &Service{client: client}
	var response map[string]interface{}
	action := "AddVodDomain"
	request := make(map[string]interface{})
	var err error

	request["DomainName"] = d.Get("domain_name")

	if v, ok := d.GetOk("scope"); ok {
		request["Scope"] = v
	}

	if v, ok := d.GetOk("check_url"); ok {
		request["CheckUrl"] = v
	}

	if v, ok := d.GetOk("top_level_domain"); ok {
		request["TopLevelDomain"] = v
	}

	sources := d.Get("sources")
	sourcesMaps := make([]map[string]interface{}, 0)
	for _, sourcesList := range convert.ConvertToInterfaceArray(sources) {
		sourcesMap := map[string]interface{}{}
		sourcesArg := sourcesList.(map[string]interface{})

		sourcesMap["type"] = sourcesArg["source_type"]
		sourcesMap["content"] = sourcesArg["source_content"]
		sourcesMap["port"] = sourcesArg["source_port"]

		if priority, ok := sourcesArg["source_priority"]; ok {
			sourcesMap["priority"] = priority
		}

		sourcesMaps = append(sourcesMaps, sourcesMap)
	}

	sourcesJson, err := convert.ConvertArrayObjectToJsonString(sourcesMaps)
	if err != nil {
		return errs.WrapError(err)
	}

	request["Sources"] = sourcesJson

	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
		response, err = client.RpcPost("vod", "2017-03-21", action, nil, request, true)
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
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, "alicloud_vod_domain", action, errs.AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(request["DomainName"]))

	stateConf := helper.BuildStateConf([]string{"configuring"}, []string{"online"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, vodService.DomainStateRefreshFunc(d.Id(), []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return errs.WrapErrorf(err, errs.IdMsg, d.Id())
	}

	return resourceDomainUpdate(d, meta)
}

func resourceDomainRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	vodService := &Service{client: client}

	object, err := vodService.DescribeDomain(d.Id())
	if err != nil {
		if errs.NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_vod_domain vodService.DescribeDomain Failed!!! %s", err)
			d.SetId("")
			return nil
		}

		return errs.WrapError(err)
	}

	d.Set("domain_name", object["DomainName"])
	d.Set("scope", object["Scope"])
	d.Set("ssl_protocol", object["SSLProtocol"])
	d.Set("cname", object["Cname"])
	d.Set("weight", object["Weight"])
	d.Set("description", object["Description"])
	d.Set("status", object["DomainStatus"])
	d.Set("gmt_created", object["GmtCreated"])
	d.Set("gmt_modified", object["GmtModified"])

	if sources, ok := object["Sources"]; ok {
		if source, ok := sources.(map[string]interface{})["Source"]; ok {
			sourcesMaps := make([]map[string]interface{}, 0)
			for _, sourceList := range convert.ConvertToInterfaceArray(source) {
				sourcesMap := map[string]interface{}{}
				sourcesArg := sourceList.(map[string]interface{})

				sourcesMap["source_type"] = sourcesArg["Type"]
				sourcesMap["source_content"] = sourcesArg["Content"]
				sourcesMap["source_port"] = sourcesArg["Port"]

				if priority, ok := sourcesArg["Priority"]; ok {
					sourcesMap["source_priority"] = priority
				}

				sourcesMaps = append(sourcesMaps, sourcesMap)
			}

			d.Set("sources", sourcesMaps)
		}
	}

	certInfo, err := vodService.DescribeDomainCertificateInfo(d.Id())
	if err != nil {
		return errs.WrapError(err)
	}

	d.Set("cert_id", certInfo["CertId"])
	d.Set("cert_name", certInfo["CertName"])
	d.Set("cert_type", certInfo["CertType"])
	d.Set("cert_region", certInfo["CertRegion"])
	d.Set("ssl_pub", certInfo["ServerCertificate"])

	return nil
}

func resourceDomainUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	vodService := &Service{client: client}
	var response map[string]interface{}
	var err error
	update := false

	request := map[string]interface{}{
		"DomainName": d.Id(),
	}

	if !d.IsNewResource() && d.HasChange("sources") {
		update = true

		sources := d.Get("sources")
		sourcesMaps := make([]map[string]interface{}, 0)
		for _, sourcesList := range convert.ConvertToInterfaceArray(sources) {
			sourcesMap := map[string]interface{}{}
			sourcesArg := sourcesList.(map[string]interface{})

			sourcesMap["type"] = sourcesArg["source_type"]
			sourcesMap["content"] = sourcesArg["source_content"]
			sourcesMap["port"] = sourcesArg["source_port"]

			if priority, ok := sourcesArg["source_priority"]; ok {
				sourcesMap["priority"] = priority
			}

			sourcesMaps = append(sourcesMaps, sourcesMap)
		}

		sourcesJson, err := convert.ConvertArrayObjectToJsonString(sourcesMaps)
		if err != nil {
			return errs.WrapError(err)
		}

		request["Sources"] = sourcesJson
	}

	if !d.IsNewResource() && d.HasChange("top_level_domain") {
		update = true

		if v, ok := d.GetOk("top_level_domain"); ok {
			request["TopLevelDomain"] = v
		}
	}

	if update {
		action := "UpdateVodDomain"

		wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
		err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
			response, err = client.RpcPost("vod", "2017-03-21", action, nil, request, true)
			if err != nil {
				if errs.IsExpectedErrors(err, []string{"3", "2001"}) || helper.NeedRetry(err) {
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

		stateConf := helper.BuildStateConf([]string{"configuring"}, []string{"online"}, d.Timeout(schema.TimeoutDelete), 5*time.Second, vodService.DomainStateRefreshFunc(d.Id(), []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return errs.WrapErrorf(err, errs.IdMsg, d.Id())
		}
	}

	update = false
	setVodDomainSSLCertificateReq := map[string]interface{}{
		"DomainName": d.Id(),
	}

	if d.HasChange("ssl_protocol") {
		update = true
	}
	if v, ok := d.GetOk("ssl_protocol"); ok {
		setVodDomainSSLCertificateReq["SSLProtocol"] = v
	}

	if d.HasChange("cert_type") {
		update = true
	}
	if v, ok := d.GetOk("cert_type"); ok {
		setVodDomainSSLCertificateReq["CertType"] = v
	}

	if d.HasChange("cert_id") {
		update = true

		if v, ok := d.GetOk("cert_id"); ok {
			setVodDomainSSLCertificateReq["CertId"] = v
		}
	}

	if d.HasChange("cert_name") {
		update = true

		if v, ok := d.GetOk("cert_name"); ok {
			setVodDomainSSLCertificateReq["CertName"] = v
		}
	}

	if d.HasChange("cert_region") {
		update = true

		if v, ok := d.GetOk("cert_region"); ok {
			setVodDomainSSLCertificateReq["CertRegion"] = v
		}
	}

	if d.HasChange("ssl_pub") {
		update = true

		if v, ok := d.GetOk("ssl_pub"); ok {
			setVodDomainSSLCertificateReq["SSLPub"] = v
		}
	}
	if d.HasChange("ssl_pri") {
		update = true

		if v, ok := d.GetOk("ssl_pri"); ok {
			setVodDomainSSLCertificateReq["SSLPri"] = v
		}
	}

	if update {
		action := "SetVodDomainSSLCertificate"

		wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
		err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
			response, err = client.RpcPost("vod", "2017-03-21", action, nil, setVodDomainSSLCertificateReq, true)
			if err != nil {
				if helper.NeedRetry(err) {
					wait()
					return retry.RetryableError(err)
				}
				return retry.NonRetryableError(err)
			}
			return nil
		})
		logs.AddDebug(action, response, setVodDomainSSLCertificateReq)

		if err != nil {
			return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
		}

		stateConf := helper.BuildStateConf([]string{"configuring"}, []string{"online"}, d.Timeout(schema.TimeoutDelete), 5*time.Second, vodService.DomainStateRefreshFunc(d.Id(), []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return errs.WrapErrorf(err, errs.IdMsg, d.Id())
		}
	}

	if err := vodService.SetResourceTags(d, "DOMAIN"); err != nil {
		return errs.WrapError(err)
	}

	return resourceDomainRead(d, meta)
}

func resourceDomainDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	vodService := &Service{client: client}

	action := "DeleteVodDomain"
	var response map[string]interface{}
	var err error

	request := map[string]interface{}{
		"DomainName": d.Id(),
	}

	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutDelete), func() *retry.RetryError {
		response, err = client.RpcPost("vod", "2017-03-21", action, nil, request, true)
		if err != nil {
			if errs.IsExpectedErrors(err, []string{"2001"}) || helper.NeedRetry(err) {
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

	stateConf := helper.BuildStateConf([]string{"deleting"}, []string{}, d.Timeout(schema.TimeoutDelete), 5*time.Second, vodService.DomainStateRefreshFunc(d.Id(), []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return errs.WrapErrorf(err, errs.IdMsg, d.Id())
	}

	return nil
}
