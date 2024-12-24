// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"encoding/json"
	"fmt"
	"log"
	"regexp"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAliCloudDcdnDomain() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudDcdnDomainCreate,
		Read:   resourceAliCloudDcdnDomainRead,
		Update: resourceAliCloudDcdnDomainUpdate,
		Delete: resourceAliCloudDcdnDomainDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
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
			"cert_region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"cert_type": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: StringInSlice([]string{"upload", "cas", "free"}, false),
			},
			"check_url": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"domain_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"env": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"function_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: StringInSlice([]string{"routine", "image", "cloudFunction"}, false),
			},
			"resource_group_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"scene": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"scope": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: StringInSlice([]string{"domestic", "overseas", "global"}, false),
			},
			"sources": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: StringInSlice([]string{"domain", "ipaddr", "oss", "common"}, false),
						},
						"content": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: StringMatch(regexp.MustCompile(".*"), "The address of the source station."),
						},
						"priority": {
							Type:         schema.TypeString,
							Optional:     true,
							Computed:     true,
							ValidateFunc: StringInSlice([]string{"20", "30"}, false),
						},
						"port": {
							Type:         schema.TypeInt,
							Optional:     true,
							Computed:     true,
							ValidateFunc: IntBetween(0, 65535),
						},
						"weight": {
							Type:         schema.TypeString,
							Optional:     true,
							Computed:     true,
							ValidateFunc: StringMatch(regexp.MustCompile("\\d+"), "The weight of the origin if multiple origins are specified. Default to `10`."),
						},
					},
				},
			},
			"ssl_pri": {
				Type:      schema.TypeString,
				Optional:  true,
				Sensitive: true,
			},
			"ssl_protocol": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: StringInSlice([]string{"off", "on"}, false),
				Default:      "off",
			},
			"ssl_pub": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"status": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: StringInSlice([]string{"online", "offline", "configuring", "configure_failed", "checking", "check_failed"}, false),
			},
			"tags": tagsSchema(),
			"top_level_domain": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"cname": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceAliCloudDcdnDomainCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "AddDcdnDomain"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	query["DomainName"] = d.Get("domain_name")

	if v, ok := d.GetOk("resource_group_id"); ok {
		request["ResourceGroupId"] = v
	}
	if v, ok := d.GetOk("tags"); ok {
		tagsMap := ConvertTags(v.(map[string]interface{}))
		request = expandTagsToMap(request, tagsMap)
	}

	if v, ok := d.GetOk("scope"); ok {
		request["Scope"] = v
	}
	if v, ok := d.GetOk("function_type"); ok {
		request["FunctionType"] = v
	}
	if v, ok := d.GetOk("scene"); ok {
		request["Scene"] = v
	}
	if v, ok := d.GetOk("check_url"); ok {
		request["CheckUrl"] = v
	}
	if v, ok := d.GetOk("top_level_domain"); ok {
		request["TopLevelDomain"] = v
	}
	if v, ok := d.GetOk("sources"); ok {
		sourcesMaps := make([]interface{}, 0)
		for _, dataLoop1 := range v.(*schema.Set).List() {
			dataLoop1Tmp := dataLoop1.(map[string]interface{})
			dataLoop1Map := make(map[string]interface{})
			dataLoop1Map["type"] = dataLoop1Tmp["type"]
			dataLoop1Map["content"] = dataLoop1Tmp["content"]
			dataLoop1Map["priority"] = dataLoop1Tmp["priority"]
			if dataLoop1Tmp["port"].(int) > 0 {
				dataLoop1Map["port"] = dataLoop1Tmp["port"]
			}
			dataLoop1Map["weight"] = dataLoop1Tmp["weight"]
			sourcesMaps = append(sourcesMaps, dataLoop1Map)
		}
		sourcesMapsJson, err := json.Marshal(sourcesMaps)
		if err != nil {
			return WrapError(err)
		}
		request["Sources"] = string(sourcesMapsJson)
	}

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.RpcPost("dcdn", "2018-01-15", action, query, request, false)
		if err != nil {
			if IsExpectedErrors(err, []string{"ServiceBusy"}) || NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(action, response, request)
		return nil
	})

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_dcdn_domain", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(query["DomainName"]))

	dcdnServiceV2 := DcdnServiceV2{client}
	stateConf := BuildStateConf([]string{}, []string{"online"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, dcdnServiceV2.DcdnDomainStateRefreshFunc(d.Id(), "DomainStatus", []string{"check_failed", "configure_failed"}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return resourceAliCloudDcdnDomainUpdate(d, meta)
}

func resourceAliCloudDcdnDomainRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	dcdnServiceV2 := DcdnServiceV2{client}

	objectRaw, err := dcdnServiceV2.DescribeDcdnDomain(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_dcdn_domain DescribeDcdnDomain Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	if objectRaw["GmtCreated"] != nil {
		d.Set("create_time", objectRaw["GmtCreated"])
	}
	if objectRaw["FunctionType"] != nil {
		d.Set("function_type", objectRaw["FunctionType"])
	}
	if objectRaw["ResourceGroupId"] != nil {
		d.Set("resource_group_id", objectRaw["ResourceGroupId"])
	}
	if objectRaw["Scene"] != nil {
		d.Set("scene", objectRaw["Scene"])
	}
	if objectRaw["Scope"] != nil {
		d.Set("scope", objectRaw["Scope"])
	}
	if objectRaw["SSLProtocol"] != nil {
		d.Set("ssl_protocol", objectRaw["SSLProtocol"])
	}
	if objectRaw["SSLPub"] != nil {
		d.Set("ssl_pub", objectRaw["SSLPub"])
	}
	if objectRaw["DomainStatus"] != nil {
		d.Set("status", objectRaw["DomainStatus"])
	}
	if objectRaw["DomainName"] != nil {
		d.Set("domain_name", objectRaw["DomainName"])
	}

	source1Raw, _ := jsonpath.Get("$.Sources.Source", objectRaw)
	sourcesMaps := make([]map[string]interface{}, 0)
	if source1Raw != nil {
		for _, sourceChild1Raw := range source1Raw.([]interface{}) {
			sourcesMap := make(map[string]interface{})
			sourceChild1Raw := sourceChild1Raw.(map[string]interface{})
			sourcesMap["content"] = sourceChild1Raw["Content"]
			sourcesMap["port"] = sourceChild1Raw["Port"]
			sourcesMap["priority"] = sourceChild1Raw["Priority"]
			sourcesMap["type"] = sourceChild1Raw["Type"]
			sourcesMap["weight"] = sourceChild1Raw["Weight"]

			sourcesMaps = append(sourcesMaps, sourcesMap)
		}
	}
	if source1Raw != nil {
		d.Set("sources", sourcesMaps)
	}
	d.Set("cname", objectRaw["Cname"])

	objectRaw, err = dcdnServiceV2.DescribeDescribeDcdnDomainCertificateInfo(d.Id())
	if err != nil {
		return WrapError(err)
	}

	if objectRaw["CertId"] != nil {
		d.Set("cert_id", objectRaw["CertId"])
	}
	if objectRaw["CertName"] != nil {
		d.Set("cert_name", objectRaw["CertName"])
	}
	if objectRaw["CertRegion"] != nil {
		d.Set("cert_region", objectRaw["CertRegion"])
	}
	if objectRaw["CertType"] != nil {
		d.Set("cert_type", objectRaw["CertType"])
	}
	if objectRaw["SSLProtocol"] != nil {
		d.Set("ssl_protocol", objectRaw["SSLProtocol"])
	}
	if objectRaw["SSLPub"] != nil {
		d.Set("ssl_pub", objectRaw["SSLPub"])
	}
	if objectRaw["DomainName"] != nil {
		d.Set("domain_name", objectRaw["DomainName"])
	}

	d.Set("domain_name", d.Id())

	dcdnService := DcdnService{client}
	listTagResourcesObject, err := dcdnService.ListTagResources(d.Id(), "DOMAIN")
	if err != nil {
		return WrapError(err)
	}

	d.Set("tags", tagsToMap(listTagResourcesObject))

	return nil
}

func resourceAliCloudDcdnDomainUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	update := false
	d.Partial(true)
	action := "ModifyDCdnDomainSchdmByProperty"
	var err error
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	query["DomainName"] = d.Id()

	if !d.IsNewResource() && d.HasChange("scope") {
		update = true
		request["Property"] = d.Get("scope")
	}

	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("dcdn", "2018-01-15", action, query, request, false)
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
		dcdnServiceV2 := DcdnServiceV2{client}
		stateConf := BuildStateConf([]string{}, []string{"online"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, dcdnServiceV2.DcdnDomainStateRefreshFunc(d.Id(), "DomainStatus", []string{"configure_failed"}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
	}
	update = false
	action = "UpdateDcdnDomain"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	query["DomainName"] = d.Id()

	if _, ok := d.GetOk("resource_group_id"); ok && !d.IsNewResource() && d.HasChange("resource_group_id") {
		update = true
		request["ResourceGroupId"] = d.Get("resource_group_id")
	}

	if !d.IsNewResource() && d.HasChange("sources") {
		update = true
		if v, ok := d.GetOk("sources"); ok {
			sourcesMaps := make([]interface{}, 0)
			for _, dataLoop := range v.(*schema.Set).List() {
				dataLoopTmp := dataLoop.(map[string]interface{})
				dataLoopMap := make(map[string]interface{})
				dataLoopMap["Type"] = dataLoopTmp["type"]
				dataLoopMap["Content"] = dataLoopTmp["content"]
				dataLoopMap["Priority"] = dataLoopTmp["priority"]
				if dataLoopTmp["port"].(int) > 0 {
					dataLoopMap["Port"] = dataLoopTmp["port"]
				}
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

	if v, ok := d.GetOk("top_level_domain"); ok {
		request["TopLevelDomain"] = v
	}
	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("dcdn", "2018-01-15", action, query, request, false)
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
		dcdnServiceV2 := DcdnServiceV2{client}
		stateConf := BuildStateConf([]string{}, []string{"online"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, dcdnServiceV2.DcdnDomainStateRefreshFunc(d.Id(), "DomainStatus", []string{"configure_failed", "check_failed"}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
	}

	dcdnService := DcdnService{client}
	if !d.IsNewResource() && d.HasChange("tags") {
		if err := dcdnService.SetResourceTags(d, "DOMAIN"); err != nil {
			return WrapError(err)
		}
	}

	update = false
	action = "SetDcdnDomainSSLCertificate"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	query["DomainName"] = d.Id()

	if d.HasChange("cert_name") {
		update = true
		request["CertName"] = d.Get("cert_name")
	}

	if d.HasChange("cert_type") {
		update = true
		request["CertType"] = d.Get("cert_type")
	}

	if d.HasChange("ssl_protocol") {
		update = true
	}
	request["SSLProtocol"] = d.Get("ssl_protocol")

	if d.HasChange("ssl_pub") {
		update = true
		request["SSLPub"] = d.Get("ssl_pub")
	}

	if d.HasChange("ssl_pri") {
		update = true
		request["SSLPri"] = d.Get("ssl_pri")
	}

	if d.HasChange("cert_id") {
		update = true
		request["CertId"] = d.Get("cert_id")
	}

	if d.HasChange("cert_region") {
		update = true
		request["CertRegion"] = d.Get("cert_region")
	}

	if v, ok := d.GetOk("env"); ok {
		request["Env"] = v
	}
	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("dcdn", "2018-01-15", action, query, request, false)
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
		dcdnServiceV2 := DcdnServiceV2{client}
		stateConf := BuildStateConf([]string{}, []string{"online"}, d.Timeout(schema.TimeoutUpdate), 20*time.Second, dcdnServiceV2.DcdnDomainStateRefreshFunc(d.Id(), "DomainStatus", []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
	}

	if d.HasChange("status") {
		client := meta.(*connectivity.AliyunClient)
		dcdnServiceV2 := DcdnServiceV2{client}
		object, err := dcdnServiceV2.DescribeDcdnDomain(d.Id())
		if err != nil {
			return WrapError(err)
		}

		target := d.Get("status").(string)
		if object["DomainStatus"].(string) != target {
			if target == "online" {
				action = "StartDcdnDomain"
				request = make(map[string]interface{})
				query = make(map[string]interface{})
				query["DomainName"] = d.Id()

				wait := incrementalWait(3*time.Second, 5*time.Second)
				err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
					response, err = client.RpcPost("dcdn", "2018-01-15", action, query, request, false)
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
				dcdnServiceV2 := DcdnServiceV2{client}
				stateConf := BuildStateConf([]string{}, []string{"online"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, dcdnServiceV2.DcdnDomainStateRefreshFunc(d.Id(), "DomainStatus", []string{"configure_failed", "check_failed"}))
				if _, err := stateConf.WaitForState(); err != nil {
					return WrapErrorf(err, IdMsg, d.Id())
				}

			}
			if target == "offline" {
				action = "StopDcdnDomain"
				request = make(map[string]interface{})
				query = make(map[string]interface{})
				query["DomainName"] = d.Id()

				wait := incrementalWait(3*time.Second, 5*time.Second)
				err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
					response, err = client.RpcPost("dcdn", "2018-01-15", action, query, request, false)
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
				dcdnServiceV2 := DcdnServiceV2{client}
				stateConf := BuildStateConf([]string{}, []string{"offline"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, dcdnServiceV2.DcdnDomainStateRefreshFunc(d.Id(), "DomainStatus", []string{"configure_failed", "check_failed"}))
				if _, err := stateConf.WaitForState(); err != nil {
					return WrapErrorf(err, IdMsg, d.Id())
				}

			}
		}
	}

	d.Partial(false)
	return resourceAliCloudDcdnDomainRead(d, meta)
}

func resourceAliCloudDcdnDomainDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	action := "DeleteDcdnDomain"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	query["DomainName"] = d.Id()

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = client.RpcPost("dcdn", "2018-01-15", action, query, request, false)

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

	dcdnServiceV2 := DcdnServiceV2{client}
	stateConf := BuildStateConf([]string{}, []string{}, d.Timeout(schema.TimeoutDelete), 3*time.Minute, dcdnServiceV2.DcdnDomainStateRefreshFunc(d.Id(), "DomainStatus", []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return nil
}
