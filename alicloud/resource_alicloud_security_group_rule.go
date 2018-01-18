package alicloud

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/denverdino/aliyungo/ecs"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceAliyunSecurityGroupRule() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliyunSecurityGroupRuleCreate,
		Read:   resourceAliyunSecurityGroupRuleRead,
		Delete: resourceAliyunSecurityGroupRuleDelete,

		Schema: map[string]*schema.Schema{
			"type": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validateSecurityRuleType,
				Description:  "Type of rule, ingress (inbound) or egress (outbound).",
			},

			"ip_protocol": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validateSecurityRuleIpProtocol,
			},

			"nic_type": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Computed:     true,
				ValidateFunc: validateSecurityRuleNicType,
			},

			"policy": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Default:      GroupRulePolicyAccept,
				ValidateFunc: validateSecurityRulePolicy,
			},

			"port_range": &schema.Schema{
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				DiffSuppressFunc: ecsSecurityGroupRulePortRangeDiffSuppressFunc,
			},

			"priority": &schema.Schema{
				Type:         schema.TypeInt,
				Optional:     true,
				ForceNew:     true,
				Default:      1,
				ValidateFunc: validateSecurityPriority,
			},

			"security_group_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"cidr_ip": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},

			"source_security_group_id": &schema.Schema{
				Type:          schema.TypeString,
				Optional:      true,
				ForceNew:      true,
				ConflictsWith: []string{"cidr_ip"},
			},

			"source_group_owner_account": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
		},
	}
}

func resourceAliyunSecurityGroupRuleCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*AliyunClient)
	conn := client.ecsconn

	direction := d.Get("type").(string)
	sgId := d.Get("security_group_id").(string)
	ptl := d.Get("ip_protocol").(string)
	port := d.Get("port_range").(string)
	if ecs.IpProtocol(ptl) != ecs.IpProtocolTCP && ecs.IpProtocol(ptl) != ecs.IpProtocolUDP {
		port = AllPortRange
	}
	nicType := d.Get("nic_type").(string)
	policy := d.Get("policy").(string)
	priority := d.Get("priority").(int)

	if _, ok := d.GetOk("cidr_ip"); !ok {
		if _, ok := d.GetOk("source_security_group_id"); !ok {
			return fmt.Errorf("Either 'cidr_ip' or 'source_security_group_id' must be specified.")
		}
	}

	var autherr error
	switch ecs.Direction(direction) {
	case ecs.DirectionIngress:
		args, err := buildAliyunSecurityIngressArgs(d, meta)
		if err != nil {
			return err
		}
		autherr = conn.AuthorizeSecurityGroup(args)
	case ecs.DirectionEgress:
		args, err := buildAliyunSecurityEgressArgs(d, meta)
		if err != nil {
			return err
		}
		autherr = conn.AuthorizeSecurityGroupEgress(args)
	default:
		return fmt.Errorf("Security Group Rule must be type 'ingress' or type 'egress'")
	}

	if autherr != nil {
		return fmt.Errorf(
			"Error authorizing security group rule type %s: %s",
			direction, autherr)
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
	client := meta.(*AliyunClient)
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
			return fmt.Errorf("SecrityGroupRuleRead parse rule id gets an error: %#v", err)
		}
		priority = prior
	}
	sgId := parts[0]
	direction := parts[1]

	rule, err := client.DescribeSecurityGroupRule(sgId, direction, parts[2], parts[3], parts[4], parts[5], policy, priority)
	if err != nil {
		if NotFoundError(err) || IsExceptedError(err, InvalidSecurityGroupIdNotFound) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error describing security group rule: %#v", err)
	}

	d.Set("type", rule.Direction)
	d.Set("ip_protocol", strings.ToLower(string(rule.IpProtocol)))
	d.Set("nic_type", rule.NicType)
	d.Set("policy", strings.ToLower(string(rule.Policy)))
	d.Set("port_range", rule.PortRange)
	d.Set("priority", rule.Priority)
	d.Set("security_group_id", sgId)
	//support source and desc by type
	if ecs.Direction(direction) == ecs.DirectionIngress {
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
	client := meta.(*AliyunClient)
	ruleType := d.Get("type").(string)

	if ecs.Direction(ruleType) == ecs.DirectionIngress {
		args, err := buildAliyunSecurityIngressArgs(d, meta)
		if err != nil {
			return err
		}
		revokeArgs := &ecs.RevokeSecurityGroupArgs{
			AuthorizeSecurityGroupArgs: *args,
		}
		return client.RevokeSecurityGroup(revokeArgs)
	}

	args, err := buildAliyunSecurityEgressArgs(d, meta)

	if err != nil {
		return err
	}
	revokeArgs := &ecs.RevokeSecurityGroupEgressArgs{
		AuthorizeSecurityGroupEgressArgs: *args,
	}
	return client.RevokeSecurityGroupEgress(revokeArgs)
}

func resourceAliyunSecurityGroupRuleDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*AliyunClient)
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
			return fmt.Errorf("SecrityGroupRuleRead parse rule id gets an error: %#v", err)
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

		_, err = client.DescribeSecurityGroupRule(parts[0], parts[1], parts[2], parts[3], parts[4], parts[5], policy, priority)
		if err != nil {
			if NotFoundError(err) || IsExceptedError(err, InvalidSecurityGroupIdNotFound) {
				return nil
			}
			return resource.NonRetryableError(err)
		}

		return resource.RetryableError(fmt.Errorf("Delete security group rule timeout and got an error: %#v", err))
	})

}

func buildAliyunSecurityIngressArgs(d *schema.ResourceData, meta interface{}) (*ecs.AuthorizeSecurityGroupArgs, error) {
	conn := meta.(*AliyunClient).ecsconn

	args := &ecs.AuthorizeSecurityGroupArgs{
		RegionId: getRegion(d, meta),
	}
	if v, ok := d.GetOk("ip_protocol"); ok {
		args.IpProtocol = ecs.IpProtocol(v.(string))
	}

	if args.IpProtocol == ecs.IpProtocolTCP || args.IpProtocol == ecs.IpProtocolUDP {
		if v, ok := d.GetOk("port_range"); ok {
			args.PortRange = v.(string)
		}
	} else {
		args.PortRange = AllPortRange
	}

	if v, ok := d.GetOk("policy"); ok {
		args.Policy = ecs.PermissionPolicy(v.(string))
	}

	if v, ok := d.GetOk("priority"); ok {
		args.Priority = v.(int)
	}

	if v, ok := d.GetOk("cidr_ip"); ok {
		args.SourceCidrIp = v.(string)
	}

	if v, ok := d.GetOk("source_security_group_id"); ok {
		args.SourceGroupId = v.(string)
	}

	if v, ok := d.GetOk("source_group_owner_account"); ok {
		args.SourceGroupOwnerAccount = v.(string)
	}

	sgId := d.Get("security_group_id").(string)

	sgArgs := &ecs.DescribeSecurityGroupAttributeArgs{
		SecurityGroupId: sgId,
		RegionId:        getRegion(d, meta),
	}

	group, err := conn.DescribeSecurityGroupAttribute(sgArgs)
	if err != nil {
		return nil, fmt.Errorf("Error get security group %s error: %#v", sgId, err)
	}

	if v, ok := d.GetOk("nic_type"); ok {
		if (group != nil && group.VpcId != "") || args.SourceGroupId != "" {
			if GroupRuleNicType(v.(string)) != GroupRuleIntranet {
				return nil, fmt.Errorf("When security group in the vpc or authorizing permission for source security group, " +
					"the nic_type must be 'intranet'.")
			}
		}
		args.NicType = ecs.NicType(v.(string))
	}

	args.SecurityGroupId = sgId

	return args, nil
}

func buildAliyunSecurityEgressArgs(d *schema.ResourceData, meta interface{}) (*ecs.AuthorizeSecurityGroupEgressArgs, error) {
	conn := meta.(*AliyunClient).ecsconn

	args := &ecs.AuthorizeSecurityGroupEgressArgs{
		RegionId: getRegion(d, meta),
	}

	if v, ok := d.GetOk("ip_protocol"); ok {
		args.IpProtocol = ecs.IpProtocol(v.(string))
	}

	if v, ok := d.GetOk("port_range"); ok {
		args.PortRange = v.(string)
	}

	if v, ok := d.GetOk("policy"); ok {
		args.Policy = ecs.PermissionPolicy(v.(string))
	}

	if v, ok := d.GetOk("priority"); ok {
		args.Priority = v.(int)
	}

	if v, ok := d.GetOk("cidr_ip"); ok {
		args.DestCidrIp = v.(string)
	}

	if v, ok := d.GetOk("source_security_group_id"); ok {
		args.DestGroupId = v.(string)
	}

	if v, ok := d.GetOk("source_group_owner_account"); ok {
		args.DestGroupOwnerAccount = v.(string)
	}

	sgId := d.Get("security_group_id").(string)

	sgArgs := &ecs.DescribeSecurityGroupAttributeArgs{
		SecurityGroupId: sgId,
		RegionId:        getRegion(d, meta),
	}

	group, err := conn.DescribeSecurityGroupAttribute(sgArgs)
	if err != nil {
		return nil, fmt.Errorf("Error get security group %s error: %#v", sgId, err)
	}

	if v, ok := d.GetOk("nic_type"); ok {
		if (group != nil && group.VpcId != "") || args.DestGroupId != "" {
			if GroupRuleNicType(v.(string)) != GroupRuleIntranet {
				return nil, fmt.Errorf("When security group in the vpc or authorizing permission for destination security group, " +
					"the nic_type must be 'intranet'.")
			}
		}
		args.NicType = ecs.NicType(v.(string))
	}

	args.SecurityGroupId = sgId

	return args, nil
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
