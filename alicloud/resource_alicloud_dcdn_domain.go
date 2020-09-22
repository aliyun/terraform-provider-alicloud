package alicloud

import (
	"fmt"
	"log"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/dcdn"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceAlicloudDcdnDomain() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudDcdnDomainCreate,
		Read:   resourceAlicloudDcdnDomainRead,
		Update: resourceAlicloudDcdnDomainUpdate,
		Delete: resourceAlicloudDcdnDomainDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"cert_name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"cert_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"free", "cas", "upload"}, false),
			},
			"check_url": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"domain_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"force_set": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"resource_group_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"ssl_pri": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"ssl_protocol": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"off", "on"}, false),
				Default:      "off",
			},
			"ssl_pub": {
				Type:     schema.TypeString,
				Optional: true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return d.Get("ssl_protocol").(string) != "on"
				},
			},
			"scope": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"domestic", "global", "overseas"}, false),
				Default:      "domestic",
			},
			"security_token": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"sources": {
				Type:     schema.TypeSet,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"content": {
							Type:     schema.TypeString,
							Required: true,
						},
						"port": {
							Type:         schema.TypeInt,
							Optional:     true,
							Default:      80,
							ValidateFunc: validation.IntInSlice([]int{443, 80}),
						},
						"priority": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  "20",
						},
						"type": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringInSlice([]string{"ipaddr", "domain", "oss"}, false),
						},
						"weight": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  "10",
						},
					},
				},
			},
			"status": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"offline", "online"}, false),
				Default:      "online",
			},
			"top_level_domain": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceAlicloudDcdnDomainCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	dcdnService := DcdnService{client}

	request := dcdn.CreateAddDcdnDomainRequest()
	if v, ok := d.GetOk("check_url"); ok {
		request.CheckUrl = v.(string)
	}

	request.DomainName = d.Get("domain_name").(string)
	if v, ok := d.GetOk("resource_group_id"); ok {
		request.ResourceGroupId = v.(string)
	}

	if v, ok := d.GetOk("scope"); ok {
		request.Scope = v.(string)
	}

	if v, ok := d.GetOk("security_token"); ok {
		request.SecurityToken = v.(string)
	}

	sources, err := dcdnService.convertSourcesToString(d.Get("sources").(*schema.Set).List())
	if err != nil {
		return WrapError(err)
	}
	request.Sources = sources
	if v, ok := d.GetOk("top_level_domain"); ok {
		request.TopLevelDomain = v.(string)
	}

	raw, err := client.WithDcdnClient(func(dcdnClient *dcdn.Client) (interface{}, error) {
		return dcdnClient.AddDcdnDomain(request)
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_dcdn_domain", request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw)
	d.SetId(fmt.Sprintf("%v", request.DomainName))
	stateConf := BuildStateConf([]string{}, []string{"online"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, dcdnService.DcdnDomainStateRefreshFunc(d.Id(), []string{"check_failed", "configure_failed"}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return resourceAlicloudDcdnDomainUpdate(d, meta)
}
func resourceAlicloudDcdnDomainRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	dcdnService := DcdnService{client}
	object, err := dcdnService.DescribeDcdnDomain(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_dcdn_domain dcdnService.DescribeDcdnDomain Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("domain_name", d.Id())
	d.Set("resource_group_id", object.ResourceGroupId)
	d.Set("ssl_protocol", convertSSLProtocolResponse(object.SSLProtocol))
	d.Set("scope", object.Scope)
	source := make([]map[string]interface{}, len(object.Sources.Source))
	for i, v := range object.Sources.Source {
		source[i] = map[string]interface{}{
			"content":  v.Content,
			"port":     v.Port,
			"priority": v.Priority,
			"type":     v.Type,
			"weight":   v.Weight,
		}
	}
	if err := d.Set("sources", source); err != nil {
		return WrapError(err)
	}
	d.Set("status", object.DomainStatus)

	describeDcdnDomainCertificateInfoObject, err := dcdnService.DescribeDcdnDomainCertificateInfo(d.Id())
	if err != nil {
		return WrapError(err)
	}
	d.Set("cert_name", describeDcdnDomainCertificateInfoObject.CertName)
	d.Set("ssl_pub", describeDcdnDomainCertificateInfoObject.SSLPub)
	return nil
}
func resourceAlicloudDcdnDomainUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	dcdnService := DcdnService{client}
	d.Partial(true)

	if !d.IsNewResource() && d.HasChange("scope") {
		request := dcdn.CreateModifyDCdnDomainSchdmByPropertyRequest()
		request.DomainName = d.Id()
		request.Property = fmt.Sprintf(`{"coverage":"%s"}`, d.Get("scope").(string))
		raw, err := client.WithDcdnClient(func(dcdnClient *dcdn.Client) (interface{}, error) {
			return dcdnClient.ModifyDCdnDomainSchdmByProperty(request)
		})
		addDebug(request.GetActionName(), raw)
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		stateConf := BuildStateConf([]string{}, []string{"online"}, d.Timeout(schema.TimeoutUpdate), 3*time.Second, dcdnService.DcdnDomainStateRefreshFunc(d.Id(), []string{"configure_failed"}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
		d.SetPartial("scope")
	}
	update := false
	request := dcdn.CreateUpdateDcdnDomainRequest()
	request.DomainName = d.Id()
	if !d.IsNewResource() && d.HasChange("resource_group_id") {
		update = true
		request.ResourceGroupId = d.Get("resource_group_id").(string)
	}
	if !d.IsNewResource() && d.HasChange("security_token") {
		update = true
		request.SecurityToken = d.Get("security_token").(string)
	}
	if !d.IsNewResource() && d.HasChange("sources") {
		update = true
		sources, err := dcdnService.convertSourcesToString(d.Get("sources").(*schema.Set).List())
		if err != nil {
			return WrapError(err)
		}
		request.Sources = sources
	}
	if !d.IsNewResource() && d.HasChange("top_level_domain") {
		update = true
		request.TopLevelDomain = d.Get("top_level_domain").(string)
	}
	if update {
		raw, err := client.WithDcdnClient(func(dcdnClient *dcdn.Client) (interface{}, error) {
			return dcdnClient.UpdateDcdnDomain(request)
		})
		addDebug(request.GetActionName(), raw)
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		stateConf := BuildStateConf([]string{}, []string{"online"}, d.Timeout(schema.TimeoutUpdate), 3*time.Second, dcdnService.DcdnDomainStateRefreshFunc(d.Id(), []string{"configure_failed"}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
		d.SetPartial("resource_group_id")
		d.SetPartial("security_token")
		d.SetPartial("sources")
		d.SetPartial("top_level_domain")
	}
	update = false
	setDcdnDomainCertificateReq := dcdn.CreateSetDcdnDomainCertificateRequest()
	setDcdnDomainCertificateReq.DomainName = d.Id()
	if d.HasChange("ssl_protocol") {
		update = true
	}
	setDcdnDomainCertificateReq.SSLProtocol = d.Get("ssl_protocol").(string)
	if d.HasChange("cert_name") {
		update = true
		setDcdnDomainCertificateReq.CertName = d.Get("cert_name").(string)
	}
	if d.HasChange("cert_type") {
		update = true
		setDcdnDomainCertificateReq.CertType = d.Get("cert_type").(string)
	}
	if d.HasChange("force_set") {
		update = true
		setDcdnDomainCertificateReq.ForceSet = d.Get("force_set").(string)
	}
	setDcdnDomainCertificateReq.Region = client.RegionId
	if d.HasChange("ssl_pri") {
		update = true
		setDcdnDomainCertificateReq.SSLPri = d.Get("ssl_pri").(string)
	}
	if d.HasChange("ssl_pub") {
		update = true
		setDcdnDomainCertificateReq.SSLPub = d.Get("ssl_pub").(string)
	}
	if !d.IsNewResource() && d.HasChange("security_token") {
		update = true
		setDcdnDomainCertificateReq.SecurityToken = d.Get("security_token").(string)
	}
	if update {
		raw, err := client.WithDcdnClient(func(dcdnClient *dcdn.Client) (interface{}, error) {
			return dcdnClient.SetDcdnDomainCertificate(setDcdnDomainCertificateReq)
		})
		addDebug(setDcdnDomainCertificateReq.GetActionName(), raw)
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), setDcdnDomainCertificateReq.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		d.SetPartial("ssl_protocol")
		d.SetPartial("cert_name")
		d.SetPartial("cert_type")
		d.SetPartial("force_set")
		d.SetPartial("ssl_pri")
		d.SetPartial("ssl_pub")
		d.SetPartial("security_token")
	}
	if d.HasChange("status") {
		object, err := dcdnService.DescribeDcdnDomain(d.Id())
		if err != nil {
			return WrapError(err)
		}
		target := d.Get("status").(string)
		if object.DomainStatus != target {
			if target == "offline" {
				request := dcdn.CreateStopDcdnDomainRequest()
				request.DomainName = d.Id()
				if v, ok := d.GetOk("security_token"); ok {
					request.SecurityToken = v.(string)
				}
				raw, err := client.WithDcdnClient(func(dcdnClient *dcdn.Client) (interface{}, error) {
					return dcdnClient.StopDcdnDomain(request)
				})
				addDebug(request.GetActionName(), raw)
				if err != nil {
					return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
				}
				stateConf := BuildStateConf([]string{}, []string{"offline"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, dcdnService.DcdnDomainStateRefreshFunc(d.Id(), []string{"configure_failed", "check_failed"}))
				if _, err := stateConf.WaitForState(); err != nil {
					return WrapErrorf(err, IdMsg, d.Id())
				}
			}
			if target == "online" {
				request := dcdn.CreateStartDcdnDomainRequest()
				request.DomainName = d.Id()
				if v, ok := d.GetOk("security_token"); ok {
					request.SecurityToken = v.(string)
				}
				raw, err := client.WithDcdnClient(func(dcdnClient *dcdn.Client) (interface{}, error) {
					return dcdnClient.StartDcdnDomain(request)
				})
				addDebug(request.GetActionName(), raw)
				if err != nil {
					return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
				}
				stateConf := BuildStateConf([]string{}, []string{"online"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, dcdnService.DcdnDomainStateRefreshFunc(d.Id(), []string{"configure_failed", "check_failed"}))
				if _, err := stateConf.WaitForState(); err != nil {
					return WrapErrorf(err, IdMsg, d.Id())
				}
			}
			d.SetPartial("status")
		}
	}
	d.Partial(false)
	return resourceAlicloudDcdnDomainRead(d, meta)
}
func resourceAlicloudDcdnDomainDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	dcdnService := DcdnService{client}
	request := dcdn.CreateDeleteDcdnDomainRequest()
	request.DomainName = d.Id()
	if v, ok := d.GetOk("resource_group_id"); ok {
		request.ResourceGroupId = v.(string)
	}
	if v, ok := d.GetOk("security_token"); ok {
		request.SecurityToken = v.(string)
	}
	raw, err := client.WithDcdnClient(func(dcdnClient *dcdn.Client) (interface{}, error) {
		return dcdnClient.DeleteDcdnDomain(request)
	})
	addDebug(request.GetActionName(), raw)
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	stateConf := BuildStateConf([]string{}, []string{}, d.Timeout(schema.TimeoutDelete), 3*time.Second, dcdnService.DcdnDomainStateRefreshFunc(d.Id(), []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}
	return nil
}
func convertSSLProtocolResponse(source string) string {
	switch source {
	case "":
		return "off"
	}
	return source
}
