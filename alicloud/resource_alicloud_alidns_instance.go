package alicloud

import (
	"fmt"
	"log"
	"strconv"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/bssopenapi"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceAlicloudAlidnsInstance() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudAlidnsInstanceCreate,
		Read:   resourceAlicloudAlidnsInstanceRead,
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
			"domain_numbers": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"payment_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"Subscription"}, false),
				Default:      "Subscription",
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
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"version_enterprise_advanced", "version_enterprise_basic", "version_personal"}, false),
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

	request := bssopenapi.CreateCreateInstanceRequest()
	request.SubscriptionType = d.Get("payment_type").(string)
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

	request.Parameter = &[]bssopenapi.CreateInstanceParameter{
		{
			Code:  "DNSSecurity",
			Value: d.Get("dns_security").(string),
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
	if !response.Success {
		return WrapErrorf(fmt.Errorf("%v", response), DefaultErrorMsg, "alicloud_alidns_instance", request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	d.SetId(fmt.Sprintf("%v", response.Data.InstanceId))

	return resourceAlicloudAlidnsInstanceRead(d, meta)
}
func resourceAlicloudAlidnsInstanceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	alidnsService := AlidnsService{client}
	object, err := alidnsService.DescribeAlidnsInstance(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_alidns_instance alidnsService.DescribeAlidnsInstance Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("dns_security", convertDnsSecurityResponse(object.DnsSecurity))
	d.Set("domain_numbers", strconv.FormatInt(object.BindDomainCount, 10))
	d.Set("version_code", object.VersionCode)
	d.Set("version_name", object.VersionName)

	return nil
}
func resourceAlicloudAlidnsInstanceDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[WARN] Cannot destroy resourceAlicloudAlidnsInstance. Terraform will remove this resource from the state file, however resources may remain.")
	return nil
}
func convertDnsSecurityResponse(source string) string {
	switch source {
	case "DNS Anti-DDoS Advanced":
		return "advanced"
	case "DNS Anti-DDoS Basic":
		return "basic"
	case "Not Required":
		return "no"
	}
	return source
}
