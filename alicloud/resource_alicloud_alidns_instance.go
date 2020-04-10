package alicloud

import (
	"log"
	"strconv"
	"strings"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/alidns"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/bssopenapi"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func resourceAlicloudAlidnsInstance() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudAlidnsInstanceCreate,
		Read:   resourceAlicloudAlidnsInstanceRead,
		Update: resourceAlicloudAlidnsInstanceUpdate,
		Delete: resourceAlicloudAlidnsInstanceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"dns_security": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"advanced", "basic", "no"}, false),
			},
			"domain_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"domain_numbers": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"lang": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"period": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
			},
			"renew_period": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
			},
			"renewal_status": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"version_code": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"version_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceAlicloudAlidnsInstanceCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	alidnsService := AlidnsService{client}

	request := bssopenapi.CreateCreateInstanceRequest()
	if v, ok := d.GetOk("period"); ok {
		request.Period = requests.NewInteger(v.(int))
	}
	request.ProductCode = "dns"
	request.ProductType = "alidns_pre"
	if v, ok := d.GetOk("renew_period"); ok {
		request.RenewPeriod = requests.NewInteger(v.(int))
	}
	if v, ok := d.GetOk("renewal_status"); ok {
		request.RenewalStatus = v.(string)
	}
	request.SubscriptionType = "Subscription"
	request.Parameter = &[]bssopenapi.CreateInstanceParameter{
		{
			Code:  "DNSSecurity",
			Value: d.Get("dns_security").(string),
		},
		{
			Code:  "Domain",
			Value: d.Get("domain_name").(string),
		},
		{
			Code:  "DomainNumbers",
			Value: d.Get("domain_numbers").(string),
		},
		{
			Code:  "Version",
			Value: d.Get("version_code").(string),
		},
	}
	raw, err := client.WithBssopenapiClient(func(bssopenapiClient *bssopenapi.Client) (interface{}, error) {
		return bssopenapiClient.CreateInstance(request)
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_alidns_instance", request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw)
	response, _ := raw.(*bssopenapi.CreateInstanceResponse)
	d.SetId(response.Data.InstanceId)
	if err := alidnsService.WaitForAlidnsInstance(d.Id(), map[string]interface{}{"Domain": d.Get("domain_name").(string)}, false, DefaultTimeoutMedium); err != nil {
		return WrapError(err)
	}

	return resourceAlicloudAlidnsInstanceRead(d, meta)
}
func resourceAlicloudAlidnsInstanceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	alidnsService := AlidnsService{client}
	object, err := alidnsService.DescribeAlidnsInstance(d.Id())
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("dns_security", convertDnsSecurityResponse(object.DnsSecurity))
	d.Set("domain_numbers", strconv.FormatInt(object.BindDomainCount, 10))
	d.Set("version_code", object.VersionCode)
	d.Set("version_name", object.VersionName)

	describeInstanceDomainsObject, err := alidnsService.DescribeInstanceDomains(d.Id())
	if err != nil {
		return WrapError(err)
	}
	d.Set("domain_name", strings.Join(flatten(describeInstanceDomainsObject), COMMA_SEPARATED))
	return nil
}
func resourceAlicloudAlidnsInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	alidnsService := AlidnsService{client}
	d.Partial(true)

	if d.HasChange("domain_name") {
		if err := alidnsService.setAlidnsInstanceDomains(d); err != nil {
			return WrapError(err)
		}
		d.SetPartial("domain_name")
	}
	d.Partial(false)
	if err := alidnsService.WaitForAlidnsInstance(d.Id(), map[string]interface{}{"Domain": d.Get("domain_name").(string)}, false, DefaultTimeoutMedium); err != nil {
		return WrapError(err)
	}
	return resourceAlicloudAlidnsInstanceRead(d, meta)
}
func resourceAlicloudAlidnsInstanceDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[WARN] Cannot destroy resourceAlicloudAlidnsInstance. Terraform will remove this resource from the state file, however resources may remain.")
	return nil
}
func convertDnsSecurityResponse(input string) string {
	switch input {
	case "DNS Anti-DDoS Advanced":
		return "advanced"
	case "DNS Anti-DDoS Basic":
		return "basic"
	case "Not Required":
		return "no"
	}
	return ""
}
func flatten(input alidns.DescribeInstanceDomainsResponse) []string {
	domainNames := make([]string, 0)
	for _, v := range input.InstanceDomains {
		domainNames = append(domainNames, v.DomainName)
	}
	return domainNames
}
