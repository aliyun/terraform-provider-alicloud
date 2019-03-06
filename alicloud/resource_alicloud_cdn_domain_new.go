package alicloud

import (
	"encoding/json"
	"strconv"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/cdn"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func resourceAlicloudCdnDomainNew() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudCdnDomainCreateNew,
		Read:   resourceAlicloudCdnDomainReadNew,
		Update: resourceAlicloudCdnDomainUpdateNew,
		Delete: resourceAlicloudCdnDomainDeleteNew,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"domain_name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateDomainName,
			},
			"cdn_type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validateCdnType,
			},
			"sources": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"content": {
							Type:     schema.TypeString,
							Required: true,
						},
						"type": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validateCdnSourceType,
						},
						"port": {
							Type:         schema.TypeInt,
							Optional:     true,
							Default:      80,
							ValidateFunc: validateCdnSourcePort,
						},
						"priority": {
							Type:         schema.TypeInt,
							Optional:     true,
							Default:      20,
							ValidateFunc: validateIntegerInRange(0, 100),
						},
						"weight": {
							Type:     schema.TypeInt,
							Optional: true,
							Default:  10,
							DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
								if v, ok := d.GetOk("sources"); ok && len(v.([]interface{})) > 0 {
									sources := make([]map[string]interface{}, 1)
									byteSources, _ := json.Marshal(v)
									json.Unmarshal(byteSources, &sources)
									if sources[0]["type"].(string) == "ipaddr" && int(sources[0]["weight"].(float64)) != 10 {
										return true
									}
								}
								return false
							},
							ValidateFunc: validateIntegerInRange(0, 100),
						},
					},
				},
				MaxItems: 1,
			},
			"scope": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validateCdnScope,
			},
		},
	}
}

func resourceAlicloudCdnDomainCreateNew(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	request := cdn.CreateAddCdnDomainRequest()
	request.DomainName = d.Get("domain_name").(string)
	request.CdnType = d.Get("cdn_type").(string)
	if v, ok := d.GetOk("scope"); ok {
		request.Scope = v.(string)
	}

	if v, ok := d.GetOk("sources"); ok && len(v.([]interface{})) > 0 {
		sources := make([]map[string]interface{}, 1)
		byteSources, _ := json.Marshal(v)
		err := json.Unmarshal(byteSources, &sources)
		if err != nil {
			return WrapError(err)
		}

		sources[0]["weight"] = strconv.Itoa(int(sources[0]["weight"].(float64)))
		sources[0]["priority"] = strconv.Itoa(int(sources[0]["priority"].(float64)))

		byteSources, _ = json.Marshal(sources)
		request.Sources = string(byteSources)
	}

	raw, err := client.WithCdnClient_new(func(cdnClient *cdn.Client) (interface{}, error) {
		return cdnClient.AddCdnDomain(request)
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "cdn_domain", request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw)

	d.SetId(request.DomainName)

	cdnservice := &CdnService{client: client}
	err = cdnservice.WaitForDomainStatusNew(d.Id(), Online, DefaultTimeoutMedium)
	if err != nil {
		return WrapError(err)
	}

	return resourceAlicloudCdnDomainReadNew(d, meta)
}

func resourceAlicloudCdnDomainUpdateNew(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	request := cdn.CreateModifyCdnDomainRequest()
	request.DomainName = d.Id()

	if d.HasChange("sources") {
		sources := make([]map[string]interface{}, 1)
		v, _ := d.GetOk("sources")
		byteSources, _ := json.Marshal(v)
		err := json.Unmarshal(byteSources, &sources)
		if err != nil {
			return WrapError(err)
		}

		sources[0]["priority"] = strconv.Itoa(int(sources[0]["priority"].(float64)))
		sources[0]["weight"] = strconv.Itoa(int(sources[0]["weight"].(float64)))
		byteSources, _ = json.Marshal(sources)
		request.Sources = string(byteSources)
	}
	raw, err := client.WithCdnClient_new(func(cdnClient *cdn.Client) (interface{}, error) {
		return cdnClient.ModifyCdnDomain(request)
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw)

	cdnservice := &CdnService{client: client}
	err = cdnservice.WaitForDomainStatusNew(d.Id(), Online, DefaultTimeoutMedium)
	if err != nil {
		return WrapError(err)
	}

	return resourceAlicloudCdnDomainReadNew(d, meta)
}

func resourceAlicloudCdnDomainReadNew(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	cdnservice := &CdnService{client: client}
	domain, err := cdnservice.DescribeCdnDomain(d.Id())
	if err != nil {
		if NotFoundError(err) {
			return nil
		}
		return WrapError(err)
	}
	if len(domain.SourceModels.SourceModel) > 0 {
		model := domain.SourceModels.SourceModel[0]
		priority, _ := strconv.Atoi(model.Priority)
		weight, _ := strconv.Atoi(model.Weight)
		sources := make([]map[string]interface{}, 1)
		sources[0] = map[string]interface{}{
			"content":  model.Content,
			"type":     model.Type,
			"port":     model.Port,
			"priority": priority,
			"weight":   weight,
		}
		err := d.Set("sources", sources)
		if err != nil {
			return WrapError(err)
		}
	}

	d.Set("domain_name", domain.DomainName)
	d.Set("cdn_type", domain.CdnType)
	d.Set("scope", domain.Scope)

	return nil
}

func resourceAlicloudCdnDomainDeleteNew(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	request := cdn.CreateDeleteCdnDomainRequest()
	request.DomainName = d.Id()
	return resource.Retry(5*time.Minute, func() *resource.RetryError {
		raw, err := client.WithCdnClient_new(func(cdnClient *cdn.Client) (interface{}, error) {
			return cdnClient.DeleteCdnDomain(request)
		})
		if err != nil {
			if IsExceptedError(err, ServiceBusy) {
				return resource.RetryableError(WrapErrorf(err, DeleteTimeoutMsg, d.Id(), request.GetActionName(), ProviderERROR))
			}
			if IsExceptedError(err, InvalidDomainNotFound) {
				return nil
			}
			return resource.NonRetryableError(WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR))
		}
		addDebug(request.GetActionName(), raw)

		cdnservice := &CdnService{client: client}
		_, err = cdnservice.DescribeCdnDomain(d.Id())
		if err != nil {
			if NotFoundError(err) || IsExceptedError(err, InvalidDomainNotFound) {
				return nil
			}

			return resource.NonRetryableError(WrapError(err))
		}

		return resource.RetryableError(WrapErrorf(err, DeleteTimeoutMsg, d.Id(), request.GetActionName(), ProviderERROR))
	})
}
