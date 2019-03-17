package alicloud

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func resourceAliyunSecurityGroupRule() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliyunSecurityGroupRuleCreate,
		Read:   resourceAliyunSecurityGroupRuleRead,
		Delete: resourceAliyunSecurityGroupRuleDelete,

		Schema: map[string]*schema.Schema{
			"type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validateSecurityRuleType,
				Description:  "Type of rule, ingress (inbound) or egress (outbound).",
			},

			"ip_protocol": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validateSecurityRuleIpProtocol,
			},

			"nic_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Computed:     true,
				ValidateFunc: validateSecurityRuleNicType,
			},

			"policy": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Default:      GroupRulePolicyAccept,
				ValidateFunc: validateSecurityRulePolicy,
			},

			"port_range": {
				Type:             schema.TypeString,
				Optional:         true,
				ForceNew:         true,
				Default:          AllPortRange,
				DiffSuppressFunc: ecsSecurityGroupRulePortRangeDiffSuppressFunc,
			},

			"priority": {
				Type:         schema.TypeInt,
				Optional:     true,
				ForceNew:     true,
				Default:      1,
				ValidateFunc: validateSecurityPriority,
			},

			"security_group_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"cidr_ip": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},

			"source_security_group_id": {
				Type:          schema.TypeString,
				Optional:      true,
				ForceNew:      true,
				ConflictsWith: []string{"cidr_ip"},
			},

			"source_group_owner_account": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
		},
	}
}

func resourceAliyunSecurityGroupRuleCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	direction := d.Get("type").(string)
	sgId := d.Get("security_group_id").(string)
	ptl := d.Get("ip_protocol").(string)
	port := d.Get("port_range").(string)
	if port == "" {
		return WrapError(fmt.Errorf("'port_range': required field is not set or invalid."))
	}
	nicType := d.Get("nic_type").(string)
	policy := d.Get("policy").(string)
	priority := d.Get("priority").(int)

	if _, ok := d.GetOk("cidr_ip"); !ok {
		if _, ok := d.GetOk("source_security_group_id"); !ok {
			return WrapError(fmt.Errorf("Either 'cidr_ip' or 'source_security_group_id' must be specified."))
		}
	}

	request, err := buildAliyunSGRuleRequest(d, meta)
	if err != nil {
		return WrapError(err)
	}

	if direction == string(DirectionIngress) {
		request.ApiName = "AuthorizeSecurityGroup"
		_, err = client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
			return ecsClient.ProcessCommonRequest(request)
		})
	} else {
		request.ApiName = "AuthorizeSecurityGroupEgress"
		_, err = client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
			return ecsClient.ProcessCommonRequest(request)
		})
	}
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "security_group_rule", request.GetActionName(), AlibabaCloudSdkGoERROR)
	}

	var cidr_ip string
	if ip, ok := d.GetOk("cidr_ip"); ok {
		cidr_ip = ip.(string)
	} else {
		cidr_ip = d.Get("source_security_group_id").(string)
	}
	d.SetId(sgId + ":" + direction + ":" + ptl + ":" + port + ":" + nicType + ":" + cidr_ip + ":" + policy + ":" + strconv.Itoa(priority))

	return resourceAliyunSecurityGroupRuleRead(d, meta)
}

func resourceAliyunSecurityGroupRuleRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	ecsService := EcsService{client}
	parts := strings.Split(d.Id(), ":")
	policy := parseSecurityRuleId(d, meta, 6)
	strPriority := parseSecurityRuleId(d, meta, 7)
	var priority int
	if policy == "" || strPriority == "" {
		policy = d.Get("policy").(string)
		priority = d.Get("priority").(int)
		d.SetId(d.Id() + ":" + policy + ":" + strconv.Itoa(priority))
	} else {
		prior, err := strconv.Atoi(strPriority)
		if err != nil {
			return WrapError(err)
		}
		priority = prior
	}
	sgId := parts[0]
	direction := parts[1]

	rule, err := ecsService.DescribeSecurityGroupRule(sgId, direction, parts[2], parts[3], parts[4], parts[5], policy, priority)
	if err != nil {
		if NotFoundError(err) || IsExceptedError(err, InvalidSecurityGroupIdNotFound) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("type", rule.Direction)
	d.Set("ip_protocol", strings.ToLower(string(rule.IpProtocol)))
	d.Set("nic_type", rule.NicType)
	d.Set("policy", strings.ToLower(string(rule.Policy)))
	d.Set("port_range", rule.PortRange)
	if pri, err := strconv.Atoi(rule.Priority); err != nil {
		return WrapError(err)
	} else {
		d.Set("priority", pri)
	}
	d.Set("security_group_id", sgId)
	//support source and desc by type
	if direction == string(DirectionIngress) {
		d.Set("cidr_ip", rule.SourceCidrIp)
		d.Set("source_security_group_id", rule.SourceGroupId)
		d.Set("source_group_owner_account", rule.SourceGroupOwnerAccount)
	} else {
		d.Set("cidr_ip", rule.DestCidrIp)
		d.Set("source_security_group_id", rule.DestGroupId)
		d.Set("source_group_owner_account", rule.DestGroupOwnerAccount)
	}
	return nil
}

func deleteSecurityGroupRule(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	ruleType := d.Get("type").(string)
	request, err := buildAliyunSGRuleRequest(d, meta)
	if err != nil {
		return WrapError(err)
	}

	if ruleType == string(DirectionIngress) {
		request.ApiName = "RevokeSecurityGroup"
		_, err = client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
			return ecsClient.ProcessCommonRequest(request)
		})
	} else {
		request.ApiName = "RevokeSecurityGroupEgress"
		_, err = client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
			return ecsClient.ProcessCommonRequest(request)
		})
	}

	return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)

}

func resourceAliyunSecurityGroupRuleDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	ecsService := EcsService{client}
	parts := strings.Split(d.Id(), ":")
	policy := parseSecurityRuleId(d, meta, 6)
	strPriority := parseSecurityRuleId(d, meta, 7)
	var priority int
	if policy == "" || strPriority == "" {
		policy = d.Get("policy").(string)
		priority = d.Get("priority").(int)
		d.SetId(d.Id() + ":" + policy + ":" + strconv.Itoa(priority))
	} else {
		prior, err := strconv.Atoi(strPriority)
		if err != nil {
			return WrapError(err)
		}
		priority = prior
	}

	return resource.Retry(5*time.Minute, func() *resource.RetryError {
		err := deleteSecurityGroupRule(d, meta)

		if err != nil {
			if NotFoundError(err) || IsExceptedError(err, InvalidSecurityGroupIdNotFound) {
				return nil
			}
			resource.RetryableError(fmt.Errorf("Delete security group rule timeout and got an error: %#v", err))
		}

		_, err = ecsService.DescribeSecurityGroupRule(parts[0], parts[1], parts[2], parts[3], parts[4], parts[5], policy, priority)
		if err != nil {
			if NotFoundError(err) || IsExceptedError(err, InvalidSecurityGroupIdNotFound) {
				return nil
			}
			return resource.NonRetryableError(err)
		}

		return resource.RetryableError(fmt.Errorf("Delete security group rule timeout and got an error: %#v", err))
	})

}

func buildAliyunSGRuleRequest(d *schema.ResourceData, meta interface{}) (*requests.CommonRequest, error) {
	client := meta.(*connectivity.AliyunClient)
	ecsService := EcsService{client}
	// Get product code from the built request
	ruleReq := ecs.CreateModifySecurityGroupRuleRequest()
	request, err := client.NewCommonRequest(ruleReq.GetProduct(), ruleReq.GetLocationServiceCode(), strings.ToUpper(string(Https)), connectivity.ApiVersion20140526)
	if err != nil {
		return request, WrapError(err)
	}

	direction := d.Get("type").(string)

	port_range := d.Get("port_range").(string)
	request.QueryParams["PortRange"] = port_range

	if v, ok := d.GetOk("ip_protocol"); ok {
		request.QueryParams["IpProtocol"] = v.(string)
		if v.(string) == string(Tcp) || v.(string) == string(Udp) {
			if port_range == AllPortRange {
				return nil, fmt.Errorf("'tcp' and 'udp' can support port range: [1, 65535]. Please correct it and try again.")
			}
		} else if port_range != AllPortRange {
			return nil, fmt.Errorf("'icmp', 'gre' and 'all' only support port range '-1/-1'. Please correct it and try again.")
		}
	}

	if v, ok := d.GetOk("policy"); ok {
		request.QueryParams["Policy"] = v.(string)
	}

	if v, ok := d.GetOk("priority"); ok {
		request.QueryParams["Priority"] = strconv.Itoa(v.(int))
	}

	if v, ok := d.GetOk("cidr_ip"); ok {
		if direction == string(DirectionIngress) {
			request.QueryParams["SourceCidrIp"] = v.(string)
		} else {
			request.QueryParams["DestCidrIp"] = v.(string)
		}
	}

	var targetGroupId string
	if v, ok := d.GetOk("source_security_group_id"); ok {
		targetGroupId = v.(string)
		if direction == string(DirectionIngress) {
			request.QueryParams["SourceGroupId"] = targetGroupId
		} else {
			request.QueryParams["DestGroupId"] = targetGroupId
		}
	}

	if v, ok := d.GetOk("source_group_owner_account"); ok {
		if direction == string(DirectionIngress) {
			request.QueryParams["SourceGroupOwnerAccount"] = v.(string)
		} else {
			request.QueryParams["DestGroupOwnerAccount"] = v.(string)
		}
	}

	sgId := d.Get("security_group_id").(string)

	group, err := ecsService.DescribeSecurityGroupAttribute(sgId)
	if err != nil {
		return nil, WrapError(err)
	}

	if v, ok := d.GetOk("nic_type"); ok {
		if group.VpcId != "" || targetGroupId != "" {
			if GroupRuleNicType(v.(string)) != GroupRuleIntranet {
				return nil, fmt.Errorf("When security group in the vpc or authorizing permission for source/destination security group, " +
					"the nic_type must be 'intranet'.")
			}
		}
		request.QueryParams["NicType"] = v.(string)
	}

	request.QueryParams["SecurityGroupId"] = sgId

	return request, nil
}

func parseSecurityRuleId(d *schema.ResourceData, meta interface{}, index int) (result string) {
	parts := strings.Split(d.Id(), ":")
	defer func() {
		if e := recover(); e != nil {
			fmt.Printf("Panicing %s\r\n", e)
			result = ""
		}
	}()
	return parts[index]
}
