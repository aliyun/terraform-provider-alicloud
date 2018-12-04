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
			"load_balancer_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"frontend_port": &schema.Schema{
				Type:         schema.TypeInt,
				ValidateFunc: validateIntegerInRange(1, 65535),
				Required:     true,
				ForceNew:     true,
			},

			"name": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: false,
				Default:  "tf-slb-rule",
			},
			"domain": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"url": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"server_group_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func resourceAliyunSlbRuleCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	slbService := SlbService{client}
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
		return fmt.Errorf("At least one 'domain' or 'url' must be set.")
	} else if domain == "" {
		rule = fmt.Sprintf("[{'RuleName':'%s','Url':'%s','VServerGroupId':'%s'}]", name, url, group_id)
	} else if url == "" {
		rule = fmt.Sprintf("[{'RuleName':'%s','Domain':'%s','VServerGroupId':'%s'}]", name, domain, group_id)
	} else {
		rule = fmt.Sprintf("[{'RuleName':'%s','Domain':'%s','Url':'%s','VServerGroupId':'%s'}]", name, domain, url, group_id)
	}

	req := slb.CreateCreateRulesRequest()
	req.LoadBalancerId = slb_id
	req.ListenerPort = requests.NewInteger(port)
	req.RuleList = rule
	if err := resource.Retry(3*time.Minute, func() *resource.RetryError {
		_, err := client.WithSlbClient(func(slbClient *slb.Client) (interface{}, error) {
			return slbClient.CreateRules(req)
		})
		if err != nil {
			if IsExceptedErrors(err, []string{BackendServerConfiguring}) {
				return resource.RetryableError(fmt.Errorf("CreateRule got an error: %#v", err))
			}
			if IsExceptedError(err, RuleDomainExist) {
				if ruleId, err := slbService.DescribeLoadBalancerRuleId(slb_id, port, domain, url); err != nil {
					return resource.NonRetryableError(err)
				} else {
					return resource.NonRetryableError(fmt.Errorf("The rule with same domain and url already exists. "+
						"Please import it using ID '%s' to import it or specify a different 'domain' or 'url' and then try again.", ruleId))
				}
			}
			return resource.NonRetryableError(fmt.Errorf("CreateRule got an error: %#v", err))
		}
		return nil
	}); err != nil {
		return err
	}

	ruleId, err := slbService.DescribeLoadBalancerRuleId(slb_id, port, domain, url)
	if err != nil {
		return err
	}

	if ruleId == "" {
		return fmt.Errorf("There is not found any rules in the load balancer %s and listener port %d.", slb_id, port)
	}

	d.SetId(ruleId)

	return resourceAliyunSlbRuleRead(d, meta)
}

func resourceAliyunSlbRuleRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	slbService := SlbService{client}
	rule, err := slbService.DescribeLoadBalancerRuleAttribute(d.Id())

	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return err
	}

	d.Set("name", rule.RuleName)
	d.Set("load_balancer_id", rule.LoadBalancerId)
	if port, err := strconv.Atoi(rule.ListenerPort); err != nil {
		return fmt.Errorf("Convertting listener port from string to int got an error: %#v.", err)
	} else {
		d.Set("frontend_port", port)
	}
	d.Set("domain", rule.Domain)
	d.Set("url", rule.Url)
	d.Set("server_group_id", rule.VServerGroupId)

	return nil
}

func resourceAliyunSlbRuleUpdate(d *schema.ResourceData, meta interface{}) error {

	d.Partial(true)

	if d.HasChange("server_group_id") {
		req := slb.CreateSetRuleRequest()
		req.RuleId = d.Id()
		req.VServerGroupId = d.Get("server_group_id").(string)
		client := meta.(*connectivity.AliyunClient)
		_, err := client.WithSlbClient(func(slbClient *slb.Client) (interface{}, error) {
			return slbClient.SetRule(req)
		})
		if err != nil {
			return fmt.Errorf("Modify rule %s server group got an error: %#v", d.Id(), err)
		}
		d.SetPartial("server_group_id")
	}

	d.Partial(false)

	return resourceAliyunSlbRuleRead(d, meta)
}

func resourceAliyunSlbRuleDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	req := slb.CreateDeleteRulesRequest()
	req.RuleIds = fmt.Sprintf("['%s']", d.Id())
	return resource.Retry(5*time.Minute, func() *resource.RetryError {
		_, err := client.WithSlbClient(func(slbClient *slb.Client) (interface{}, error) {
			return slbClient.DeleteRules(req)
		})
		if err != nil {
			if IsExceptedErrors(err, []string{InvalidRuleIdNotFound}) {
				return nil
			}
			return resource.NonRetryableError(err)
		}

		client := meta.(*connectivity.AliyunClient)
		slbService := SlbService{client}
		if _, err := slbService.DescribeLoadBalancerRuleAttribute(d.Id()); err != nil {
			if NotFoundError(err) {
				return nil
			}
			return resource.NonRetryableError(fmt.Errorf("While deleting rule, DescribeRuleAttribute got an error: %#v", err))
		}

		return resource.RetryableError(fmt.Errorf("DeleteRule %s timeout.", d.Id()))
	})
}
