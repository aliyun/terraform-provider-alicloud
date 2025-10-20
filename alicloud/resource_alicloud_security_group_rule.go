package alicloud

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAliyunSecurityGroupRule() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliyunSecurityGroupRuleCreate,
		Read:   resourceAliyunSecurityGroupRuleRead,
		Update: resourceAliyunSecurityGroupRuleUpdate,
		Delete: resourceAliyunSecurityGroupRuleDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"security_group_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"ingress", "egress"}, false),
				Description:  "Type of rule, ingress (inbound) or egress (outbound).",
			},
			"ip_protocol": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"tcp", "udp", "icmp", "icmpv6", "gre", "all"}, false),
			},
			"policy": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Default:      GroupRulePolicyAccept,
				ValidateFunc: validation.StringInSlice([]string{"accept", "drop"}, false),
			},
			"priority": {
				Type:         schema.TypeInt,
				Optional:     true,
				ForceNew:     true,
				Default:      1,
				ValidateFunc: validation.IntBetween(1, 100),
			},
			"cidr_ip": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				AtLeastOneOf: []string{"cidr_ip", "ipv6_cidr_ip", "source_security_group_id", "prefix_list_id"},
			},
			"ipv6_cidr_ip": {
				Type:          schema.TypeString,
				Optional:      true,
				ForceNew:      true,
				ConflictsWith: []string{"cidr_ip"},
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					v, _ := compressIPv6OrCIDR(new)
					return v == old
				},
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
			"prefix_list_id": {
				Type:             schema.TypeString,
				Optional:         true,
				Computed:         true,
				ForceNew:         true,
				DiffSuppressFunc: ecsSecurityGroupRulePreFixListIdDiffSuppressFunc,
			},
			"port_range": {
				Type:             schema.TypeString,
				Optional:         true,
				ForceNew:         true,
				Default:          AllPortRange,
				DiffSuppressFunc: ecsSecurityGroupRulePortRangeDiffSuppressFunc,
			},
			"nic_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Computed:     true,
				ValidateFunc: validation.StringInSlice([]string{"internet", "intranet"}, false),
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"security_group_rule_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceAliyunSecurityGroupRuleCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	ecsService := EcsService{client}
	var response map[string]interface{}
	var cidrIp string
	var sourceSecurityGroupId string
	var prefixListId string
	request := make(map[string]interface{})
	var err error

	request["RegionId"] = client.RegionId

	securityGroupId := d.Get("security_group_id").(string)
	request["SecurityGroupId"] = securityGroupId

	direction := d.Get("type").(string)

	permissionsMaps := make([]map[string]interface{}, 0)
	permissionsMap := map[string]interface{}{}
	permissionsMap["IpProtocol"] = d.Get("ip_protocol")

	if v, ok := d.GetOk("policy"); ok {
		permissionsMap["Policy"] = v
	}

	if v, ok := d.GetOk("priority"); ok {
		permissionsMap["Priority"] = strconv.Itoa(v.(int))
	}

	if v, ok := d.GetOk("cidr_ip"); ok {
		cidrIp = v.(string)
		if direction == string(DirectionIngress) {
			permissionsMap["SourceCidrIp"] = v
		} else {
			permissionsMap["DestCidrIp"] = v
		}
	}

	if v, ok := d.GetOk("ipv6_cidr_ip"); ok {
		cidrIp = strings.Replace(v.(string), ":", "_", -1)
		if direction == string(DirectionIngress) {
			permissionsMap["Ipv6SourceCidrIp"] = v
		} else {
			permissionsMap["Ipv6DestCidrIp"] = v
		}
	}

	if v, ok := d.GetOk("source_security_group_id"); ok {
		cidrIp = v.(string)
		sourceSecurityGroupId = v.(string)
		if direction == string(DirectionIngress) {
			permissionsMap["SourceGroupId"] = v
		} else {
			permissionsMap["DestGroupId"] = v
		}
	}

	if v, ok := d.GetOk("source_group_owner_account"); ok {
		if direction == string(DirectionIngress) {
			permissionsMap["SourceGroupOwnerAccount"] = v
		} else {
			permissionsMap["DestGroupOwnerAccount"] = v
		}
	}

	if v, ok := d.GetOk("prefix_list_id"); ok {
		prefixListId = v.(string)
		if direction == string(DirectionIngress) {
			permissionsMap["SourcePrefixListId"] = v
		} else {
			permissionsMap["DestPrefixListId"] = v
		}
	}

	if v, ok := d.GetOk("port_range"); ok {
		permissionsMap["PortRange"] = v

		if permissionsMap["IpProtocol"].(string) == string(Tcp) || permissionsMap["IpProtocol"].(string) == string(Udp) {
			if v.(string) == AllPortRange {
				return fmt.Errorf(" 'tcp' and 'udp' can support port range: [1, 65535]. Please correct it and try again.")
			}
		} else if v.(string) != AllPortRange {
			return fmt.Errorf(" 'icmp', 'gre' and 'all' only support port range '-1/-1'. Please correct it and try again.")
		}
	}

	securityGroup, err := ecsService.DescribeSecurityGroup(securityGroupId)
	if err != nil {
		return WrapError(err)
	}

	if v, ok := d.GetOk("nic_type"); ok {
		if securityGroup.VpcId != "" || sourceSecurityGroupId != "" {
			if GroupRuleNicType(v.(string)) != GroupRuleIntranet {
				return fmt.Errorf(" When security group in the vpc or authorizing permission for source/destination security group, the nic_type must be 'intranet'.")
			}
		}

		permissionsMap["NicType"] = v
	}

	if v, ok := d.GetOk("description"); ok {
		permissionsMap["Description"] = v
	}

	permissionsMaps = append(permissionsMaps, permissionsMap)
	request["Permissions"] = permissionsMaps

	if direction == string(DirectionIngress) {
		action := "AuthorizeSecurityGroup"
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutCreate)), func() *resource.RetryError {
			response, err = client.RpcPost("Ecs", "2014-05-26", action, nil, request, true)
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})
		addDebug(action, response, request)

		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, "alicloud_security_group_rule", action, AlibabaCloudSdkGoERROR)
		}
	} else {
		action := "AuthorizeSecurityGroupEgress"

		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutCreate)), func() *resource.RetryError {
			response, err = client.RpcPost("Ecs", "2014-05-26", action, nil, request, true)
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})
		addDebug(action, response, request)

		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, "alicloud_security_group_rule", action, AlibabaCloudSdkGoERROR)
		}
	}

	if len(cidrIp) != 0 {
		d.SetId(fmt.Sprintf("%v:%v:%v:%v:%v:%v:%v:%v", request["SecurityGroupId"], direction, permissionsMap["IpProtocol"], permissionsMap["PortRange"], permissionsMap["NicType"], cidrIp, permissionsMap["Policy"], permissionsMap["Priority"]))
	} else {
		d.SetId(fmt.Sprintf("%v:%v:%v:%v:%v:%v:%v:%v", request["SecurityGroupId"], direction, permissionsMap["IpProtocol"], permissionsMap["PortRange"], permissionsMap["NicType"], prefixListId, permissionsMap["Policy"], permissionsMap["Priority"]))
	}

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

	object, err := ecsService.DescribeSecurityGroupRule(d.Id())
	if err != nil {
		if NotFoundError(err) && !d.IsNewResource() {
			log.Printf("[DEBUG] Resource alicloud_security_group_rule ecsService.DescribeSecurityGroupRule Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("type", object.Direction)
	d.Set("ip_protocol", strings.ToLower(string(object.IpProtocol)))
	d.Set("nic_type", object.NicType)
	d.Set("policy", strings.ToLower(string(object.Policy)))
	d.Set("port_range", object.PortRange)
	d.Set("description", object.Description)
	d.Set("security_group_rule_id", object.SecurityGroupRuleId)

	if pri, err := strconv.Atoi(object.Priority); err != nil {
		return WrapError(err)
	} else {
		d.Set("priority", pri)
	}
	d.Set("security_group_id", sgId)

	//support source and desc by type
	if direction == string(DirectionIngress) {
		d.Set("cidr_ip", object.SourceCidrIp)
		d.Set("ipv6_cidr_ip", object.Ipv6SourceCidrIp)
		d.Set("source_security_group_id", object.SourceGroupId)
		d.Set("source_group_owner_account", object.SourceGroupOwnerAccount)
		d.Set("prefix_list_id", object.SourcePrefixListId)
	} else {
		d.Set("cidr_ip", object.DestCidrIp)
		d.Set("ipv6_cidr_ip", object.Ipv6DestCidrIp)
		d.Set("source_security_group_id", object.DestGroupId)
		d.Set("source_group_owner_account", object.DestGroupOwnerAccount)
		d.Set("prefix_list_id", object.DestPrefixListId)
	}

	return nil
}

func resourceAliyunSecurityGroupRuleUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

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

	request, err := buildAliyunSGRuleRequest(d, meta)
	if err != nil {
		return WrapError(err)
	}

	direction := d.Get("type").(string)

	if direction == string(DirectionIngress) {
		request.ApiName = "ModifySecurityGroupRule"
		_, err = client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
			return ecsClient.ProcessCommonRequest(request)
		})
	} else {
		request.ApiName = "ModifySecurityGroupEgressRule"
		_, err = client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
			return ecsClient.ProcessCommonRequest(request)
		})
	}

	raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
		return ecsClient.ProcessCommonRequest(request)
	})

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
	}

	addDebug(request.GetActionName(), raw, request.Headers, request)

	return resourceAliyunSecurityGroupRuleRead(d, meta)
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

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	return nil
}

func resourceAliyunSecurityGroupRuleDelete(d *schema.ResourceData, meta interface{}) error {
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

	err := resource.Retry(5*time.Minute, func() *resource.RetryError {
		err := deleteSecurityGroupRule(d, meta)
		if err != nil {
			if NotFoundError(err) || IsExpectedErrors(err, []string{"InvalidSecurityGroupId.NotFound"}) {
				return nil
			}
			return resource.RetryableError(err)
		}
		return nil
	})
	if err != nil {
		return WrapError(err)
	}
	return nil
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
	if v, ok := d.GetOk("ipv6_cidr_ip"); ok {
		if direction == string(DirectionIngress) {
			request.QueryParams["Ipv6SourceCidrIp"] = v.(string)
		} else {
			request.QueryParams["Ipv6DestCidrIp"] = v.(string)
		}
	}
	if v, ok := d.GetOk("prefix_list_id"); ok {
		if direction == string(DirectionIngress) {
			request.QueryParams["SourcePrefixListId"] = v.(string)
		} else {
			request.QueryParams["DestPrefixListId"] = v.(string)
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

	group, err := ecsService.DescribeSecurityGroup(sgId)
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

	description := d.Get("description").(string)
	request.QueryParams["Description"] = description

	return request, nil
}

func parseSecurityRuleId(d *schema.ResourceData, meta interface{}, index int) (result string) {
	parts := strings.Split(d.Id(), ":")
	defer func() {
		if e := recover(); e != nil {
			log.Printf("Panicing %s\r\n", e)
			result = ""
		}
	}()
	return parts[index]
}
