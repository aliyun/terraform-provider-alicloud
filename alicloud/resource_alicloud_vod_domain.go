package alicloud

import (
	"fmt"
	"log"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAliCloudVodDomain() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudVodDomainCreate,
		Read:   resourceAliCloudVodDomainRead,
		Update: resourceAliCloudVodDomainUpdate,
		Delete: resourceAliCloudVodDomainDelete,
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
				ValidateFunc: StringInSlice([]string{"domestic", "overseas", "global"}, false),
			},
			"ssl_protocol": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: StringInSlice([]string{"on", "off"}, false),
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
				ValidateFunc: StringInSlice([]string{"cn-hangzhou", "ap-southeast-1"}, false),
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
			"tags": tagsSchemaWithElements(),
			"sources": {
				Type:     schema.TypeSet,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"source_type": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: StringInSlice([]string{"ipaddr", "domain", "oss"}, false),
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
							ValidateFunc: StringInSlice([]string{"20", "30"}, false),
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

func resourceAliCloudVodDomainCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	vodService := VodService{client}
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
	for _, sourcesList := range convertToInterfaceArray(sources) {
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

	sourcesJson, err := convertArrayObjectToJsonString(sourcesMaps)
	if err != nil {
		return WrapError(err)
	}

	request["Sources"] = sourcesJson

	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
		response, err = client.RpcPost("vod", "2017-03-21", action, nil, request, true)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_vod_domain", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(request["DomainName"]))

	stateConf := BuildStateConf([]string{"configuring"}, []string{"online"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, vodService.VodStateRefreshFunc(d.Id(), []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return resourceAliCloudVodDomainUpdate(d, meta)
}

func resourceAliCloudVodDomainRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	vodService := VodService{client}

	object, err := vodService.DescribeVodDomain(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_vod_domain vodService.DescribeVodDomain Failed!!! %s", err)
			d.SetId("")
			return nil
		}

		return WrapError(err)
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
			for _, sourceList := range convertToInterfaceArray(source) {
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

	certInfo, err := vodService.DescribeVodDomainCertificateInfo(d.Id())
	if err != nil {
		return WrapError(err)
	}

	d.Set("cert_id", certInfo["CertId"])
	d.Set("cert_name", certInfo["CertName"])
	d.Set("cert_type", certInfo["CertType"])
	d.Set("cert_region", certInfo["CertRegion"])
	d.Set("ssl_pub", certInfo["ServerCertificate"])

	return nil
}

func resourceAliCloudVodDomainUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	vodService := VodService{client}
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
		for _, sourcesList := range convertToInterfaceArray(sources) {
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

		sourcesJson, err := convertArrayObjectToJsonString(sourcesMaps)
		if err != nil {
			return WrapError(err)
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

		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
			response, err = client.RpcPost("vod", "2017-03-21", action, nil, request, true)
			if err != nil {
				if IsExpectedErrors(err, []string{"3", "2001"}) || NeedRetry(err) {
					wait()
					return retry.RetryableError(err)
				}
				return retry.NonRetryableError(err)
			}
			return nil
		})
		addDebug(action, response, request)

		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}

		stateConf := BuildStateConf([]string{"configuring"}, []string{"online"}, d.Timeout(schema.TimeoutDelete), 5*time.Second, vodService.VodStateRefreshFunc(d.Id(), []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
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

		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
			response, err = client.RpcPost("vod", "2017-03-21", action, nil, setVodDomainSSLCertificateReq, true)
			if err != nil {
				if NeedRetry(err) {
					wait()
					return retry.RetryableError(err)
				}
				return retry.NonRetryableError(err)
			}
			return nil
		})
		addDebug(action, response, setVodDomainSSLCertificateReq)

		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}

		stateConf := BuildStateConf([]string{"configuring"}, []string{"online"}, d.Timeout(schema.TimeoutDelete), 5*time.Second, vodService.VodStateRefreshFunc(d.Id(), []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
	}

	if err := vodService.SetResourceTags(d, "DOMAIN"); err != nil {
		return WrapError(err)
	}

	return resourceAliCloudVodDomainRead(d, meta)
}

func resourceAliCloudVodDomainDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	vodService := VodService{client}

	action := "DeleteVodDomain"
	var response map[string]interface{}
	var err error

	request := map[string]interface{}{
		"DomainName": d.Id(),
	}

	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutDelete), func() *retry.RetryError {
		response, err = client.RpcPost("vod", "2017-03-21", action, nil, request, true)
		if err != nil {
			if IsExpectedErrors(err, []string{"2001"}) || NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	stateConf := BuildStateConf([]string{"deleting"}, []string{}, d.Timeout(schema.TimeoutDelete), 5*time.Second, vodService.VodStateRefreshFunc(d.Id(), []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return nil
}
