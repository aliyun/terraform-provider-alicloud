package alicloud

import (
	"fmt"
	"strings"
	"time"

	"strconv"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/slb"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func resourceAliyunSlbRule() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliyunSlbRuleCreate,
		Read:   resourceAliyunSlbRuleRead,
		Update: resourceAliyunSlbRuleUpdate,
		Delete: resourceAliyunSlbRuleDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"load_balancer_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"frontend_port": {
				Type:         schema.TypeInt,
				ValidateFunc: validateIntegerInRange(1, 65535),
				Required:     true,
				ForceNew:     true,
			},

			"name": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Default:  "tf-slb-rule",
			},
			"domain": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"url": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"server_group_id": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func resourceAliyunSlbRuleCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	slb_id := d.Get("load_balancer_id").(string)
	port := d.Get("frontend_port").(int)
	name := strings.Trim(d.Get("name").(string), " ")
	group_id := strings.Trim(d.Get("server_group_id").(string), " ")

	var domain, url, rule string
	if v, ok := d.GetOk("domain"); ok {
		domain = v.(string)
	}
	if v, ok := d.GetOk("url"); ok {
		url = v.(string)
	}

	if domain == "" && url == "" {
		return WrapError(Error("At least one 'domain' or 'url' must be set."))
	} else if domain == "" {
		rule = fmt.Sprintf("[{'RuleName':'%s','Url':'%s','VServerGroupId':'%s'}]", name, url, group_id)
	} else if url == "" {
		rule = fmt.Sprintf("[{'RuleName':'%s','Domain':'%s','VServerGroupId':'%s'}]", name, domain, group_id)
	} else {
		rule = fmt.Sprintf("[{'RuleName':'%s','Domain':'%s','Url':'%s','VServerGroupId':'%s'}]", name, domain, url, group_id)
	}

	request := slb.CreateCreateRulesRequest()
	request.LoadBalancerId = slb_id
	request.ListenerPort = requests.NewInteger(port)
	request.RuleList = rule
	var raw interface{}
	var err error
	if err = resource.Retry(3*time.Minute, func() *resource.RetryError {
		raw, err = client.WithSlbClient(func(slbClient *slb.Client) (interface{}, error) {
			return slbClient.CreateRules(request)
		})
		if err != nil {
			if IsExceptedErrors(err, []string{BackendServerConfiguring}) {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(request.GetActionName(), raw)
		return nil
	}); err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_slb_rule", request.GetActionName(), AlibabaCloudSdkGoERROR)
	}

	response, _ := raw.(*slb.CreateRulesResponse)
	d.SetId(response.Rules.Rule[0].RuleId)

	return resourceAliyunSlbRuleRead(d, meta)
}

func resourceAliyunSlbRuleRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	slbService := SlbService{client}
	object, err := slbService.DescribeSlbRule(d.Id())

	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("name", object.RuleName)
	d.Set("load_balancer_id", object.LoadBalancerId)
	if port, err := strconv.Atoi(object.ListenerPort); err != nil {
		return WrapError(err)
	} else {
		d.Set("frontend_port", port)
	}
	d.Set("domain", object.Domain)
	d.Set("url", object.Url)
	d.Set("server_group_id", object.VServerGroupId)

	return nil
}

func resourceAliyunSlbRuleUpdate(d *schema.ResourceData, meta interface{}) error {
	if d.HasChange("server_group_id") {
		request := slb.CreateSetRuleRequest()
		request.RuleId = d.Id()
		request.VServerGroupId = d.Get("server_group_id").(string)
		client := meta.(*connectivity.AliyunClient)
		raw, err := client.WithSlbClient(func(slbClient *slb.Client) (interface{}, error) {
			return slbClient.SetRule(request)
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		addDebug(request.GetActionName(), raw)
	}

	return resourceAliyunSlbRuleRead(d, meta)
}

func resourceAliyunSlbRuleDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	slbService := SlbService{client}
	request := slb.CreateDeleteRulesRequest()
	request.RuleIds = fmt.Sprintf("['%s']", d.Id())

	raw, err := client.WithSlbClient(func(slbClient *slb.Client) (interface{}, error) {
		return slbClient.DeleteRules(request)
	})
	if err != nil {
		if IsExceptedError(err, InvalidRuleIdNotFound) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw)

	return WrapError(slbService.WaitForSlbRule(d.Id(), Deleted, DefaultTimeoutMedium))

}
