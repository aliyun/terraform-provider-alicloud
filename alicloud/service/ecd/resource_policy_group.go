package ecd

import (
	"context"
	"fmt"
	"log"
	"sort"
	"strings"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/convert"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func ResourcePolicyGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourcePolicyGroupCreate,
		Read:   resourcePolicyGroupRead,
		Update: resourcePolicyGroupUpdate,
		Delete: resourcePolicyGroupDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		CustomizeDiff: func(ctx context.Context, d *schema.ResourceDiff, meta interface{}) error {
			clientTypes, _ := d.Get("client_types").(*schema.Set)
			if clientTypes != nil {
				if err := ValidatePolicyGroupClientTypes(clientTypes.List()); err != nil {
					return err
				}
			}
			// domain_list (legacy) and domain_resolve_rule[_type] (structured) are
			// mutually exclusive through ConflictsWith. The structured fields also
			// need a cross-field check that the schema cannot express: rules are only
			// effective while the switch is on, and ModifyPolicyGroup keeps omitted
			// repeatList parameters, so an "on" switch with no rule blocks would
			// leave server-side rules behind and keep planning a diff.
			if rc := d.GetRawConfig(); !rc.IsNull() && rc.IsKnown() {
				if err := validatePolicyGroupDomainConfig(rc.GetAttr("domain_resolve_rule_type"), rc.GetAttr("domain_resolve_rule")); err != nil {
					return err
				}
				priorList, _ := d.GetChange("domain_list")
				priorRules, _ := d.GetChange("domain_resolve_rule")
				priorListValue, _ := priorList.(string)
				priorRuleCount := 0
				if rules, ok := priorRules.(*schema.Set); ok {
					priorRuleCount = rules.Len()
				}
				return validatePolicyGroupDomainFormSwitch(priorListValue, priorRuleCount, expandPolicyGroupDomainForm(rc))
			}
			return nil
		},
		Schema: map[string]*schema.Schema{
			"authorize_access_policy_rules": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"cidr_ip": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"description": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"authorize_security_policy_rules": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"cidr_ip": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"description": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"ip_protocol": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"policy": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"port_range": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"priority": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"type": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"clipboard": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringInSlice([]string{"off", "read", "readwrite"}, false),
			},
			"domain_list": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"domain_resolve_rule", "domain_resolve_rule_type"},
			},
			"domain_resolve_rule_type": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringInSlice([]string{"off", "on"}, false),
			},
			"domain_resolve_rule": {
				Type:          schema.TypeSet,
				Optional:      true,
				ConflictsWith: []string{"domain_list"},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"domain": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"policy": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringInSlice([]string{"allow", "block"}, false),
						},
						"description": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"html_access": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringInSlice([]string{"off", "on"}, false),
			},
			"html_file_transfer": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringInSlice([]string{"all", "download", "off", "upload"}, false),
			},
			"local_drive": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringInSlice([]string{"readwrite", "off", "read"}, false),
			},
			"policy_group_name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"usb_redirect": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringInSlice([]string{"off", "on"}, false),
			},
			"visual_quality": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringInSlice([]string{"high", "lossless", "low", "medium"}, false),
			},
			"watermark": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringInSlice([]string{"off", "on"}, false),
			},
			"watermark_transparency": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringInSlice([]string{"DARK", "LIGHT", "MIDDLE"}, false),
			},
			"watermark_type": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringInSlice([]string{"EndUserId", "HostName"}, false),
			},
			"recording": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringInSlice([]string{"off", "alltime", "period"}, false),
			},
			"recording_start_time": {
				Type:     schema.TypeString,
				Optional: true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					if v, ok := d.GetOk("recording"); ok && v.(string) == "period" {
						return false
					}
					return true
				},
			},
			"recording_end_time": {
				Type:     schema.TypeString,
				Optional: true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					if v, ok := d.GetOk("recording"); ok && v.(string) == "period" {
						return false
					}
					return true
				},
			},
			"recording_fps": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.IntInSlice([]int{2, 5, 10, 15}),
			},
			"camera_redirect": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringInSlice([]string{"on", "off"}, false),
			},
			"recording_expires": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.IntBetween(15, 180),
			},
			"internet_communication_protocol": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringInSlice([]string{"BOTH", "TCP"}, false),
			},
			"recording_duration": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.IntInSlice([]int{10, 20, 30, 60}),
			},
			"client_types": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"client_type": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringInSlice([]string{"android", "html5", "ios", "macos", "windows"}, false),
						},
						"status": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringInSlice([]string{"off", "on"}, false),
						},
					},
				},
			},
			"end_user_apply_admin_coordinate": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringInSlice([]string{"off", "on"}, false),
			},
			"end_user_group_coordinate": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringInSlice([]string{"off", "on"}, false),
			},
		},
	}
}

func resourcePolicyGroupCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "CreatePolicyGroup"
	request := make(map[string]interface{})
	var err error
	if v, ok := d.GetOk("authorize_access_policy_rules"); ok {
		for authorizeAccessPolicyRulesPtr, authorizeAccessPolicyRules := range v.(*schema.Set).List() {
			authorizeAccessPolicyRulesArg := authorizeAccessPolicyRules.(map[string]interface{})
			request["AuthorizeAccessPolicyRule."+fmt.Sprint(authorizeAccessPolicyRulesPtr+1)+".CidrIp"] = authorizeAccessPolicyRulesArg["cidr_ip"]
			request["AuthorizeAccessPolicyRule."+fmt.Sprint(authorizeAccessPolicyRulesPtr+1)+".Description"] = authorizeAccessPolicyRulesArg["description"]
		}
	}
	if v, ok := d.GetOk("authorize_security_policy_rules"); ok {
		for authorizeSecurityPolicyRulesPtr, authorizeSecurityPolicyRules := range v.(*schema.Set).List() {
			authorizeSecurityPolicyRulesArg := authorizeSecurityPolicyRules.(map[string]interface{})
			request["AuthorizeSecurityPolicyRule."+fmt.Sprint(authorizeSecurityPolicyRulesPtr+1)+".CidrIp"] = authorizeSecurityPolicyRulesArg["cidr_ip"]
			request["AuthorizeSecurityPolicyRule."+fmt.Sprint(authorizeSecurityPolicyRulesPtr+1)+".Description"] = authorizeSecurityPolicyRulesArg["description"]
			request["AuthorizeSecurityPolicyRule."+fmt.Sprint(authorizeSecurityPolicyRulesPtr+1)+".IpProtocol"] = authorizeSecurityPolicyRulesArg["ip_protocol"]
			request["AuthorizeSecurityPolicyRule."+fmt.Sprint(authorizeSecurityPolicyRulesPtr+1)+".Policy"] = authorizeSecurityPolicyRulesArg["policy"]
			request["AuthorizeSecurityPolicyRule."+fmt.Sprint(authorizeSecurityPolicyRulesPtr+1)+".PortRange"] = authorizeSecurityPolicyRulesArg["port_range"]
			request["AuthorizeSecurityPolicyRule."+fmt.Sprint(authorizeSecurityPolicyRulesPtr+1)+".Priority"] = authorizeSecurityPolicyRulesArg["priority"]
			request["AuthorizeSecurityPolicyRule."+fmt.Sprint(authorizeSecurityPolicyRulesPtr+1)+".Type"] = authorizeSecurityPolicyRulesArg["type"]
		}
	}
	if v, ok := d.GetOk("clipboard"); ok {
		request["Clipboard"] = v
	}
	if v, ok := d.GetOk("domain_list"); ok {
		request["DomainList"] = v
	}
	if v, ok := d.GetOk("domain_resolve_rule_type"); ok {
		// The wire enum is uppercase OFF/ON; the schema and state use lowercase.
		request["DomainResolveRuleType"] = strings.ToUpper(v.(string))
	}
	if v, ok := d.GetOk("domain_resolve_rule"); ok {
		for domainResolveRulePtr, domainResolveRule := range v.(*schema.Set).List() {
			domainResolveRuleArg, ok := domainResolveRule.(map[string]interface{})
			if !ok {
				continue
			}
			request["DomainResolveRule."+fmt.Sprint(domainResolveRulePtr+1)+".Domain"] = domainResolveRuleArg["domain"]
			request["DomainResolveRule."+fmt.Sprint(domainResolveRulePtr+1)+".Policy"] = domainResolveRuleArg["policy"]
			request["DomainResolveRule."+fmt.Sprint(domainResolveRulePtr+1)+".Description"] = domainResolveRuleArg["description"]
		}
	}
	if v, ok := d.GetOk("html_access"); ok {
		request["Html5Access"] = v
	}
	if v, ok := d.GetOk("html_file_transfer"); ok {
		request["Html5FileTransfer"] = v
	}
	if v, ok := d.GetOk("local_drive"); ok {
		request["LocalDrive"] = v
	}
	if v, ok := d.GetOk("policy_group_name"); ok {
		request["Name"] = v
	}
	request["RegionId"] = client.RegionId
	if v, ok := d.GetOk("usb_redirect"); ok {
		request["UsbRedirect"] = v
	}
	if v, ok := d.GetOk("visual_quality"); ok {
		request["VisualQuality"] = v
	}
	if v, ok := d.GetOk("watermark"); ok {
		request["Watermark"] = v
	}
	if v, ok := d.GetOk("watermark_transparency"); ok {
		request["WatermarkTransparency"] = v
	}
	if v, ok := d.GetOk("watermark_type"); ok {
		request["WatermarkType"] = v
	}
	if v, ok := d.GetOk("recording"); ok {
		request["Recording"] = v
	}
	if v, ok := d.GetOk("recording_start_time"); ok {
		request["RecordingStartTime"] = v
	}
	if v, ok := d.GetOk("recording_end_time"); ok {
		request["RecordingEndTime"] = v
	}
	if v, ok := d.GetOk("recording_fps"); ok {
		request["RecordingFps"] = v
	}
	if v, ok := d.GetOk("camera_redirect"); ok {
		request["CameraRedirect"] = v
	}
	if v, ok := d.GetOk("recording_expires"); ok {
		request["RecordingExpires"] = v
	}
	if v, ok := d.GetOk("internet_communication_protocol"); ok {
		request["InternetCommunicationProtocol"] = v
	}
	if v, ok := d.GetOk("recording_duration"); ok {
		request["RecordingDuration"] = v
	}
	if v, ok := d.GetOk("client_types"); ok {
		for clientTypesPtr, clientTypes := range v.(*schema.Set).List() {
			clientTypesArg, ok := clientTypes.(map[string]interface{})
			if !ok {
				continue
			}
			request["ClientType."+fmt.Sprint(clientTypesPtr+1)+".ClientType"] = clientTypesArg["client_type"]
			// Status is optional: send it only when configured, otherwise the API
			// keeps the type's default. The wire enum is uppercase OFF/ON.
			if status, ok := clientTypesArg["status"].(string); ok && status != "" {
				request["ClientType."+fmt.Sprint(clientTypesPtr+1)+".Status"] = strings.ToUpper(status)
			}
		}
	}
	if v, ok := d.GetOk("end_user_apply_admin_coordinate"); ok {
		request["EndUserApplyAdminCoordinate"] = v
	}
	if v, ok := d.GetOk("end_user_group_coordinate"); ok {
		request["EndUserGroupCoordinate"] = v
	}
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
		response, err = client.RpcPost("ecd", "2020-09-30", action, nil, request, false)
		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	logs.AddDebug(action, response, request)
	if err != nil {
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, "alicloud_ecd_policy_group", action, errs.AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(response["PolicyGroupId"]))

	// The domain access control this create asked for can be echoed stale by
	// DescribePolicyGroups right after CreatePolicyGroup returns, and the
	// trailing Read would record that into state. Wait for the requested
	// payload to be observed twice first, best-effort like the Update
	// convergence wait: a timeout logs and leaves Read to record whatever the
	// server reports.
	requestedType, _ := request["DomainResolveRuleType"].(string)
	target := strings.ToLower(requestedType)
	requestedList, _ := request["DomainList"].(string)
	var requestedRules []map[string]interface{}
	if v, ok := d.GetOk("domain_resolve_rule"); ok {
		for _, rule := range v.(*schema.Set).List() {
			if ruleMap, ok := rule.(map[string]interface{}); ok {
				requestedRules = append(requestedRules, ruleMap)
			}
		}
	}
	if target != "" || requestedList != "" || len(requestedRules) > 0 {
		if err := waitPolicyGroupDomainConvergence(NewService(client), d.Id(), d.Timeout(schema.TimeoutCreate), target, requestedList, requestedRules); err != nil {
			log.Printf("[WARN] ecd_policy_group: domain access control did not converge after create: %s", err)
		}
	}

	return resourcePolicyGroupRead(d, meta)
}
func resourcePolicyGroupRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	ecdService := NewService(client)
	object, err := ecdService.DescribeEcdPolicyGroup(d.Id())
	if err != nil {
		if errs.NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_ecd_policy_group ecdService.DescribeEcdPolicyGroup Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return errs.WrapError(err)
	}
	if authorizeAccessPolicyRulesList, ok := object["AuthorizeAccessPolicyRules"]; ok && authorizeAccessPolicyRulesList != nil {
		authorizeAccessPolicyRulesMaps := make([]map[string]interface{}, 0)
		for _, authorizeAccessPolicyRulesListItem := range authorizeAccessPolicyRulesList.([]interface{}) {
			if authorizeAccessPolicyRulesListItemMap, ok := authorizeAccessPolicyRulesListItem.(map[string]interface{}); ok {
				authorizeAccessPolicyRulesMap := make(map[string]interface{})
				authorizeAccessPolicyRulesMap["cidr_ip"] = authorizeAccessPolicyRulesListItemMap["CidrIp"]
				authorizeAccessPolicyRulesMap["description"] = authorizeAccessPolicyRulesListItemMap["Description"]
				authorizeAccessPolicyRulesMaps = append(authorizeAccessPolicyRulesMaps, authorizeAccessPolicyRulesMap)
			}
		}
		d.Set("authorize_access_policy_rules", authorizeAccessPolicyRulesMaps)
	}
	if v, ok := object["AuthorizeSecurityPolicyRules"].([]interface{}); ok {
		authorizeSecurityPolicyRules := make([]map[string]interface{}, 0)
		for _, val := range v {
			item := val.(map[string]interface{})
			temp := map[string]interface{}{
				"cidr_ip":     item["CidrIp"],
				"description": item["Description"],
				"ip_protocol": item["IpProtocol"],
				"policy":      item["Policy"],
				"port_range":  item["PortRange"],
				"priority":    item["Priority"],
				"type":        item["Type"],
			}

			authorizeSecurityPolicyRules = append(authorizeSecurityPolicyRules, temp)
		}
		if err := d.Set("authorize_security_policy_rules", authorizeSecurityPolicyRules); err != nil {
			return errs.WrapError(err)
		}
	}
	d.Set("clipboard", object["Clipboard"])
	if err := flattenPolicyGroupDomainConfig(d, object); err != nil {
		return err
	}
	d.Set("html_access", object["Html5Access"])
	d.Set("html_file_transfer", object["Html5FileTransfer"])
	d.Set("local_drive", object["LocalDrive"])
	d.Set("policy_group_name", object["Name"])
	d.Set("status", object["PolicyStatus"])
	d.Set("usb_redirect", object["UsbRedirect"])
	d.Set("visual_quality", object["VisualQuality"])
	d.Set("watermark", object["Watermark"])
	d.Set("watermark_transparency", object["WatermarkTransparency"])
	d.Set("watermark_type", object["WatermarkType"])
	d.Set("recording", object["Recording"])
	recording_fpsVal, err := convert.FormatInt(object["RecordingFps"])
	if err != nil {
		return err
	}
	d.Set("recording_fps", recording_fpsVal)
	d.Set("camera_redirect", object["CameraRedirect"])
	recording_expiresVal, err := convert.FormatInt(object["RecordingExpires"])
	if err != nil {
		return err
	}
	d.Set("recording_expires", recording_expiresVal)
	if protocol, ok := object["InternetCommunicationProtocol"].(string); ok && protocol != "" {
		d.Set("internet_communication_protocol", strings.ToUpper(protocol))
	}
	recordingDurationVal, err := convert.FormatInt(object["RecordingDuration"])
	if err != nil {
		return err
	}
	d.Set("recording_duration", recordingDurationVal)
	// DescribePolicyGroups returns client types under the plural "ClientTypes"
	// key, consistent with AuthorizeAccessPolicyRules / AuthorizeSecurityPolicyRules
	// in this same Read: the request uses the singular repeatList prefix
	// ClientType.{N}.ClientType/.Status, and the response echoes it as a plural
	// array. The singular "ClientType" key is empty (verified via ACC). The
	// Status enum is uppercase OFF/ON on the wire, so it is lowercased for state.
	if err := setPolicyGroupClientTypes(d, object); err != nil {
		return err
	}
	d.Set("end_user_apply_admin_coordinate", object["EndUserApplyAdminCoordinate"])
	d.Set("end_user_group_coordinate", object["EndUserGroupCoordinate"])
	return nil
}

// ValidatePolicyGroupClientTypes rejects client_types entries that could never
// converge: an empty client_type cannot be matched back during Read, and a
// duplicated client_type is collapsed into one state element.
func ValidatePolicyGroupClientTypes(clientTypes []interface{}) error {
	seen := make(map[string]bool, len(clientTypes))
	for _, rawEntry := range clientTypes {
		entryMap, ok := rawEntry.(map[string]interface{})
		if !ok {
			return fmt.Errorf("client_types: client_type is required when a client_types block is configured")
		}
		clientType, _ := entryMap["client_type"].(string)
		if clientType == "" {
			return fmt.Errorf("client_types: client_type is required when a client_types block is configured")
		}
		if seen[clientType] {
			return fmt.Errorf("client_types: duplicate client_type %q", clientType)
		}
		seen[clientType] = true
	}
	return nil
}

// validatePolicyGroupDomainConfig enforces the cross-field constraints of the
// structured domain access control fields. Values that are not known at plan
// time (for example a rule set derived from a resource that does not exist yet)
// are skipped here; there is no provider-side re-evaluation at apply time, so
// such values rely on the server rejecting the request.
func validatePolicyGroupDomainConfig(typeVal, rulesVal cty.Value) error {
	if !typeVal.IsKnown() {
		return nil
	}
	typeConfigured := !typeVal.IsNull() && typeVal.AsString() != ""
	domains := make([]string, 0, 1)
	if rulesVal.IsKnown() && !rulesVal.IsNull() {
		for it := rulesVal.ElementIterator(); it.Next(); {
			_, rule := it.Element()
			if !rule.IsKnown() {
				return nil
			}
			if rule.IsNull() {
				return fmt.Errorf("each domain_resolve_rule block must set a domain")
			}
			domain := rule.GetAttr("domain")
			if !domain.IsKnown() {
				return nil
			}
			if domain.IsNull() || strings.TrimSpace(domain.AsString()) == "" {
				return fmt.Errorf("each domain_resolve_rule block requires a non-empty domain")
			}
			domains = append(domains, domain.AsString())
		}
	}
	if typeConfigured {
		switch strings.ToLower(typeVal.AsString()) {
		case "on":
			if rulesVal.IsKnown() && len(domains) == 0 {
				return fmt.Errorf("domain_resolve_rule_type=on requires at least one domain_resolve_rule block; to disable domain access control, set domain_resolve_rule_type to off and drop the domain_resolve_rule blocks")
			}
		case "off":
			if len(domains) > 0 {
				return fmt.Errorf("domain_resolve_rule blocks must not be configured together with domain_resolve_rule_type=off")
			}
		}
		return nil
	}
	if len(domains) > 0 {
		return fmt.Errorf("domain_resolve_rule blocks require domain_resolve_rule_type=on")
	}
	return nil
}

// policyGroupClientTypeSelection is one client type the user configured, plus
// whether an explicit status was configured for it.
type policyGroupClientTypeSelection struct {
	status    string
	hasStatus bool
}

// policyGroupClientTypeFilterFromConfig reads the configured client types. The
// second return is false when client_types is absent from the configuration, in
// which case state follows the prior state.
func policyGroupClientTypeFilterFromConfig(rawConfig cty.Value) (map[string]policyGroupClientTypeSelection, bool) {
	if rawConfig.IsNull() || !rawConfig.IsKnown() {
		return nil, false
	}
	configList := rawConfig.GetAttr("client_types")
	if configList.IsNull() || !configList.IsKnown() {
		return nil, false
	}
	selection := make(map[string]policyGroupClientTypeSelection)
	for it := configList.ElementIterator(); it.Next(); {
		_, rawEntry := it.Element()
		clientType := rawEntry.GetAttr("client_type")
		if clientType.IsNull() || !clientType.IsKnown() || clientType.AsString() == "" {
			continue
		}
		entry := policyGroupClientTypeSelection{}
		status := rawEntry.GetAttr("status")
		if !status.IsNull() && status.IsKnown() && status.AsString() != "" {
			entry.status = strings.ToLower(status.AsString())
			entry.hasStatus = true
		}
		selection[clientType.AsString()] = entry
	}
	return selection, true
}

// policyGroupClientTypeFilterFromState mirrors the last applied subset during
// refresh, where Terraform does not pass the configuration to Read.
func policyGroupClientTypeFilterFromState(clientTypes *schema.Set) (map[string]policyGroupClientTypeSelection, bool) {
	if clientTypes == nil || clientTypes.Len() == 0 {
		return nil, false
	}
	selection := make(map[string]policyGroupClientTypeSelection)
	for _, rawEntry := range clientTypes.List() {
		entryMap, ok := rawEntry.(map[string]interface{})
		if !ok {
			continue
		}
		clientType, _ := entryMap["client_type"].(string)
		if clientType == "" {
			continue
		}
		entry := policyGroupClientTypeSelection{}
		if status, _ := entryMap["status"].(string); status != "" {
			entry.status = strings.ToLower(status)
			entry.hasStatus = true
		}
		selection[clientType] = entry
	}
	return selection, true
}

// client_types is Optional and DescribePolicyGroups always returns every client
// type, so state must mirror exactly what the user configured instead of echoing
// the API set: raw configuration during apply, prior state during refresh, and
// an empty set when the blocks are removed. Import has neither configuration nor
// prior state, so it leaves client_types empty.
func setPolicyGroupClientTypes(d *schema.ResourceData, object map[string]interface{}) error {
	clientTypesList, _ := object["ClientTypes"].([]interface{})
	apiClientTypes := make(map[string]map[string]interface{}, len(clientTypesList))
	for _, rawEntry := range clientTypesList {
		entryMap, ok := rawEntry.(map[string]interface{})
		if !ok {
			continue
		}
		if clientType, ok := entryMap["ClientType"].(string); ok && clientType != "" {
			apiClientTypes[clientType] = entryMap
		}
	}

	selection, filtered := policyGroupClientTypeFilterFromConfig(d.GetRawConfig())
	if !filtered {
		clientTypes, _ := d.Get("client_types").(*schema.Set)
		selection, filtered = policyGroupClientTypeFilterFromState(clientTypes)
	}
	if !filtered {
		selection = map[string]policyGroupClientTypeSelection{}
	}

	clientTypesMaps := make([]map[string]interface{}, 0, len(selection))
	for clientType, entry := range selection {
		clientTypesMap := map[string]interface{}{"client_type": clientType}
		// Only entries with a configured status are echoed back, so an entry whose
		// status is left unset keeps an identical element hash across refresh.
		if entry.hasStatus {
			status := entry.status
			if apiEntry, ok := apiClientTypes[clientType]; ok {
				if apiStatus, ok := apiEntry["Status"].(string); ok && apiStatus != "" {
					status = strings.ToLower(apiStatus)
				}
			}
			clientTypesMap["status"] = status
		}
		clientTypesMaps = append(clientTypesMaps, clientTypesMap)
	}
	if err := d.Set("client_types", clientTypesMaps); err != nil {
		return errs.WrapError(err)
	}
	return nil
}

// expandPolicyGroupDomainForm classifies the configured domain access control
// representation: "structured" when domain_resolve_rule_type or
// domain_resolve_rule manages the structured form, "legacy" when only
// domain_list is set, and "" when neither is configured. The empty string is
// authoritative when the configuration is present: the user manages no domain
// access control. An explicit domain_resolve_rule_type="off" also classifies as
// "": it is the disable operation, not a switch to the structured form, because
// sending DomainResolveRuleType=OFF clears both representations on the server.
// Classifying it as "structured" would reject the documented single-apply
// cleanup (drop the managed domain_list and set domain_resolve_rule_type to
// "off") as a form switch, even though no stale representation can survive.
func expandPolicyGroupDomainForm(rawConfig cty.Value) string {
	if rawConfig.IsNull() || !rawConfig.IsKnown() {
		return ""
	}
	if t := rawConfig.GetAttr("domain_resolve_rule_type"); !t.IsNull() && t.IsKnown() && t.AsString() != "" {
		if strings.ToLower(t.AsString()) == "off" {
			return ""
		}
		return "structured"
	}
	if r := rawConfig.GetAttr("domain_resolve_rule"); !r.IsNull() && r.IsKnown() && r.LengthInt() > 0 {
		return "structured"
	}
	if d := rawConfig.GetAttr("domain_list"); !d.IsNull() && d.IsKnown() && d.AsString() != "" {
		return "legacy"
	}
	return ""
}

// validatePolicyGroupDomainFormSwitch rejects switching between the legacy
// domain_list form and the structured domain_resolve_rule form in one apply.
// ModifyPolicyGroup keeps omitted parameters, so a switch would leave the other
// form's domain rules on the server while Read hides them, producing a silent
// stale access-control configuration. Users disable domain access control first
// (type=off), then configure the other form on a later apply. The disable step
// itself is not a form switch: expandPolicyGroupDomainForm classifies a
// domain_resolve_rule_type=off configuration as "" — sending OFF clears both
// representations — so the documented single-apply cleanup that drops the
// managed form and sets domain_resolve_rule_type=off is accepted.
func validatePolicyGroupDomainFormSwitch(priorDomainList string, priorRuleCount int, nextForm string) error {
	priorForm := ""
	if priorDomainList != "" {
		priorForm = "legacy"
	} else if priorRuleCount > 0 {
		priorForm = "structured"
	}
	if priorForm != "" && nextForm != "" && priorForm != nextForm {
		return fmt.Errorf("cannot switch between domain_list and domain_resolve_rule in a single apply; set domain_resolve_rule_type to off and drop the configured form first, apply, then configure the other form")
	}
	return nil
}

// flattenPolicyGroupDomainForm mirrors the last applied representation from the
// state. The legacy domain_list is checked before the structured rules: both
// representations can coexist in state right after an apply that crosses forms, but
// a non-empty domain_list is the authoritative marker of the legacy form.
// domain_resolve_rule_type is intentionally not a discriminator here: the server
// echoes "on" for both the legacy and structured forms and Read stores it for both,
// so keying on it would misclassify a legacy config as structured, repopulate the
// server-derived rules, clear domain_list and produce a perpetual diff on the next
// refresh. Besides refresh it is also the fallback for a Read whose configuration
// manages no domain field, and for the first import Read before any form has been
// recorded.
func flattenPolicyGroupDomainForm(d *schema.ResourceData) string {
	if v, ok := d.GetOk("domain_list"); ok && v.(string) != "" {
		return "legacy"
	}
	if v, ok := d.GetOk("domain_resolve_rule"); ok && v.(*schema.Set).Len() > 0 {
		return "structured"
	}
	return ""
}

// serverPolicyGroupDomainForm picks the representation to record for domain rules
// that exist only on the server, used by the first import Read where neither the
// configuration nor the prior state records a form: "structured" when a server rule
// carries a description only the structured form can represent, otherwise "legacy"
// for the DomainList string the rules round-trip through losslessly, and "" when
// the server reports no domain rules. Recording the rules keeps the control that is
// actually in effect visible in state, and gives the form-switch guard the
// old-form information it needs to reject a later cross-form adoption that would
// strand the server-side rules.
func serverPolicyGroupDomainForm(object map[string]interface{}, serverRules []map[string]interface{}) string {
	for _, rule := range serverRules {
		if description, ok := rule["description"].(string); ok && description != "" {
			return "structured"
		}
	}
	if domainList, ok := object["DomainList"].(string); ok && domainList != "" {
		return "legacy"
	}
	if len(serverRules) > 0 {
		return "structured"
	}
	return ""
}

// stateRecordedPolicyGroupDomainRuleType reports whether the state this Read runs
// against already carries a domain_resolve_rule_type observation. The first import
// Read runs against state that synthesized nothing but the resource ID, so its raw
// state carries no observation and the import still adopts the server's form. A
// refresh whose state recorded the switch — a group this configuration never
// managed, or one whose removal converged — leaves server-only rules unrecorded
// instead of adopting them, so an unrelated update cannot show a removal diff and
// disable control that exists only out of band.
func stateRecordedPolicyGroupDomainRuleType(rawState cty.Value) bool {
	if rawState == cty.NilVal || rawState.IsNull() || !rawState.IsKnown() || !rawState.Type().IsObjectType() {
		return false
	}
	if !rawState.Type().HasAttribute("domain_resolve_rule_type") {
		return false
	}
	t := rawState.GetAttr("domain_resolve_rule_type")
	return t.IsKnown() && !t.IsNull()
}

// flattenDomainResolveRuleType reads the DomainResolveRuleType field from a
// DescribePolicyGroups object and lowercases it, returning "" when absent. The
// wire enum is uppercase OFF/ON while the response echoes lowercase off/on;
// lowercasing both keeps comparisons consistent with the schema enum.
func flattenDomainResolveRuleType(object map[string]interface{}) string {
	if t, ok := object["DomainResolveRuleType"].(string); ok && t != "" {
		return strings.ToLower(t)
	}
	return ""
}

// policyGroupDomainConfig is the domain access control configuration the user
// manages, expanded from the raw configuration.
type policyGroupDomainConfig struct {
	ruleType   string
	domainList string
	rules      []map[string]interface{}
}

// expandPolicyGroupDomainConfig reads the domain access control fields the user
// configures. The raw configuration is authoritative: Optional(+Computed) fields
// otherwise return the computed state value for fields the configuration omits.
func expandPolicyGroupDomainConfig(rawConfig cty.Value, d *schema.ResourceData) policyGroupDomainConfig {
	config := policyGroupDomainConfig{}
	if rawConfig.IsNull() || !rawConfig.IsKnown() {
		return config
	}
	if t := rawConfig.GetAttr("domain_resolve_rule_type"); !t.IsNull() && t.IsKnown() {
		config.ruleType = t.AsString()
	}
	if v := rawConfig.GetAttr("domain_list"); !v.IsNull() && v.IsKnown() {
		config.domainList = v.AsString()
	}
	if r := rawConfig.GetAttr("domain_resolve_rule"); !r.IsNull() && r.IsKnown() && r.LengthInt() > 0 {
		if v, ok := d.GetOk("domain_resolve_rule"); ok {
			for _, rule := range v.(*schema.Set).List() {
				ruleMap, ok := rule.(map[string]interface{})
				if !ok {
					continue
				}
				config.rules = append(config.rules, ruleMap)
			}
		}
	}
	return config
}

// shouldDisableDomainAccessControl reports whether this update removed domain
// access control the state recorded. It must depend on a plan change, not on prior
// state alone: policy groups whose state never recorded a domain form also read
// back "on", and an unrelated update must not pick up and disable their out-of-band
// configuration. State that did record a form — the configuration managed it, or an
// import adopted it from the server — plans a removal diff the user can see, and
// applying that removal sends OFF deliberately.
func shouldDisableDomainAccessControl(listChanged, listConfigured, rulesChanged, rulesConfigured bool) bool {
	return (listChanged && !listConfigured) || (rulesChanged && !rulesConfigured)
}

// flattenDomainResolveRules reads the structured DomainResolveRule list from a
// DescribePolicyGroups object into the shape state records
// (domain/policy/description).
func flattenDomainResolveRules(object map[string]interface{}) []map[string]interface{} {
	v, ok := object["DomainResolveRule"].([]interface{})
	if !ok {
		return nil
	}
	rules := make([]map[string]interface{}, 0, len(v))
	for _, val := range v {
		item, ok := val.(map[string]interface{})
		if !ok {
			continue
		}
		rules = append(rules, map[string]interface{}{
			"domain":      item["Domain"],
			"policy":      item["Policy"],
			"description": item["Description"],
		})
	}
	return rules
}

// domainRuleSetKey renders one domain rule for set comparison.
// domain_resolve_rule is a set, so order never matters; a field the response
// omits (nil) and an unset schema string ("") normalize to the empty string so
// a converged echo still matches.
func domainRuleSetKey(rule map[string]interface{}) string {
	field := func(name string) string {
		if s, ok := rule[name].(string); ok {
			return s
		}
		return ""
	}
	return field("domain") + "|" + field("policy") + "|" + field("description")
}

// sameDomainRuleSet compares two rule collections with set semantics.
func sameDomainRuleSet(a, b []map[string]interface{}) bool {
	if len(a) != len(b) {
		return false
	}
	keys := make([]string, 0, len(a))
	for _, rule := range a {
		keys = append(keys, domainRuleSetKey(rule))
	}
	sort.Strings(keys)
	for _, rule := range b {
		k := domainRuleSetKey(rule)
		i := sort.SearchStrings(keys, k)
		if i == len(keys) || keys[i] != k {
			return false
		}
		keys = append(keys[:i], keys[i+1:]...)
	}
	return true
}

// policyGroupDomainEchoConverged reports whether a DescribePolicyGroups object
// reflects the domain access control an apply asked for. The empty string means
// every requested payload converged; otherwise the first mismatch is described
// for the caller to wait on. Payloads the request did not carry are not
// compared: ModifyPolicyGroup keeps omitted parameters, so an update must not
// wait on domain state it never asked to change. The legacy DomainList is
// compared as the exact string and the structured rules as a set — the shapes
// Read records into state, so a converged echo is exactly what the trailing
// Read will hold.
func policyGroupDomainEchoConverged(object map[string]interface{}, targetType, targetDomainList string, targetRules []map[string]interface{}) string {
	if targetType != "" {
		if got := flattenDomainResolveRuleType(object); got != targetType {
			return fmt.Sprintf("DomainResolveRuleType is %q, want %q", got, targetType)
		}
	}
	if targetDomainList != "" {
		got, _ := object["DomainList"].(string)
		if got != targetDomainList {
			return fmt.Sprintf("DomainList is %q, want %q", got, targetDomainList)
		}
	}
	if len(targetRules) > 0 {
		if got := flattenDomainResolveRules(object); !sameDomainRuleSet(got, targetRules) {
			return fmt.Sprintf("DomainResolveRule has %d rules, want %d", len(got), len(targetRules))
		}
	}
	return ""
}

// waitPolicyGroupDomainConvergence polls DescribePolicyGroups until the server
// echoes the requested domain access control — the DomainResolveRuleType
// switch, the legacy DomainList string, and the structured rule set — twice in
// a row, guarding against an echo that reached the target and then reverted
// before the server finished converging. Best-effort within the given timeout:
// the returned error reports a timeout or Describe failure and the caller
// decides whether to act on it.
func waitPolicyGroupDomainConvergence(ecdService *Service, groupId string, timeout time.Duration, targetType, targetDomainList string, targetRules []map[string]interface{}) error {
	convergeWait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	return retry.Retry(timeout, func() *retry.RetryError {
		obj, err := ecdService.DescribeEcdPolicyGroup(groupId)
		if err != nil {
			return retry.NonRetryableError(err)
		}
		if msg := policyGroupDomainEchoConverged(obj, targetType, targetDomainList, targetRules); msg != "" {
			convergeWait()
			return retry.RetryableError(fmt.Errorf("waiting for domain access control to converge: %s", msg))
		}
		confirmed, err := ecdService.DescribeEcdPolicyGroup(groupId)
		if err != nil {
			return retry.NonRetryableError(err)
		}
		if msg := policyGroupDomainEchoConverged(confirmed, targetType, targetDomainList, targetRules); msg != "" {
			convergeWait()
			return retry.RetryableError(fmt.Errorf("domain access control echo not confirmed: %s", msg))
		}
		return nil
	})
}

// resolvePolicyGroupDomainRuleType resolves the state value of the Optional+Computed
// domain_resolve_rule_type field from the switch the server reports. The server echo
// owns the value so refresh records the actual state: a close that failed or was
// reverted out of band surfaces as a plan diff against a configured "off" and the
// next apply re-sends OFF, instead of the configured value permanently masking the
// server switch. The post-update window where DescribePolicyGroups can still echo
// the old value is handled before Read returns: Update polls until the requested
// type is confirmed twice. The prior state's "off" is sticky only while the server
// has not re-enabled: re-enabling the legacy form sends DomainList without
// DomainResolveRuleType, the server turns itself back on, and Read must mirror that
// echo — sticking to the prior "off" there would clear the domain list from state
// while the server keeps enforcing it, leaving a diff that never converges. An
// absent server value carries no switch observation, so the sticky "off" also keeps
// state stable when the response omits the field.
func resolvePolicyGroupDomainRuleType(priorType, serverType string) string {
	if priorType == "off" && serverType != "on" {
		return "off"
	}
	return serverType
}

// flattenPolicyGroupDomainConfig reconciles the dual representations of ECD domain
// access control — the legacy DomainList string and the structured
// DomainResolveRule repeatList — which the server derives from each other and
// echoes together. State mirrors the representation the configuration manages, or
// the one the state already records (adopted from the server at import), so the
// server-echoed alternate form does not produce a perpetual diff. When domain
// access control is disabled (DomainResolveRuleType=off) both are cleared, since
// removal sends OFF and the server has cleared the configuration.
func flattenPolicyGroupDomainConfig(d *schema.ResourceData, object map[string]interface{}) error {
	// domain_resolve_rule_type is Optional+Computed and is resolved from the switch
	// the server reports, so refresh records the actual server state: a close that
	// failed or was reverted out of band surfaces as a plan diff against a
	// configured value and the next apply re-sends the request. The post-update
	// stale-echo window is handled before Read: Update waits for the requested type
	// to be confirmed twice. The prior state's "off" is sticky while the server has
	// not re-enabled, which also keeps state stable when the response omits the
	// switch entirely.
	priorType := ""
	if t, ok := d.GetOk("domain_resolve_rule_type"); ok {
		priorType = strings.ToLower(t.(string))
	}
	domainResolveRuleType := resolvePolicyGroupDomainRuleType(priorType, flattenDomainResolveRuleType(object))
	d.Set("domain_resolve_rule_type", domainResolveRuleType)

	serverRules := flattenDomainResolveRules(object)

	// DomainResolveRuleType is the authoritative switch; when off, the server has
	// cleared the configuration, so normalize both representations to empty even
	// if a stale echo remains.
	if domainResolveRuleType == "off" {
		d.Set("domain_list", "")
		if err := d.Set("domain_resolve_rule", []map[string]interface{}{}); err != nil {
			return errs.WrapError(err)
		}
		return nil
	}

	form := ""
	if rawConfig := d.GetRawConfig(); rawConfig.IsNull() || !rawConfig.IsKnown() {
		// No configuration is passed to Read. The first import Read also has no
		// recorded representation in state, so it adopts one from the server return:
		// the domain rules actually in effect become visible in state, and the
		// form-switch guard gains the old-form information it needs to reject a
		// later cross-form adoption that would strand them. A later refresh whose
		// state already recorded the switch has no form only when the rules exist
		// out of band, so it does not adopt — an unrelated update must not pick
		// them up, plan their removal, and disable them. The removal whose OFF has
		// not converged is re-recorded by Update itself, so it stays visible
		// without adoption here.
		form = flattenPolicyGroupDomainForm(d)
		if form == "" && !stateRecordedPolicyGroupDomainRuleType(d.GetRawState()) {
			form = serverPolicyGroupDomainForm(object, serverRules)
		}
	} else {
		form = expandPolicyGroupDomainForm(rawConfig)
		if form == "" {
			// The configuration manages no domain field. Mirror the server's rules
			// in the representation the state already records instead of clearing
			// them: a removal whose OFF has not converged (or a re-enable out of
			// band) stays visible as a plan diff and the next apply re-sends OFF,
			// rather than state silently dropping rules the server still enforces.
			// State that never recorded a domain form (never managed, never
			// imported) keeps clearing, so an unrelated update does not pick up
			// out-of-band domain configuration and start removing it.
			form = flattenPolicyGroupDomainForm(d)
		}
	}

	switch form {
	case "structured":
		// The user manages the structured form; do not echo the derived legacy
		// DomainList string, which would drift against an unset domain_list.
		d.Set("domain_list", "")
		if err := d.Set("domain_resolve_rule", serverRules); err != nil {
			return errs.WrapError(err)
		}
	case "legacy":
		// The user manages the legacy DomainList string; do not echo the derived
		// structured rules, which would drift against an unset domain_resolve_rule.
		d.Set("domain_list", object["DomainList"])
		if err := d.Set("domain_resolve_rule", []map[string]interface{}{}); err != nil {
			return errs.WrapError(err)
		}
	default:
		d.Set("domain_list", "")
		if err := d.Set("domain_resolve_rule", []map[string]interface{}{}); err != nil {
			return errs.WrapError(err)
		}
	}
	return nil
}

func resourcePolicyGroupUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	ecdService := NewService(client)
	var err error
	var response map[string]interface{}
	update := false
	request := map[string]interface{}{
		"PolicyGroupId": d.Id(),
	}
	request["RegionId"] = client.RegionId
	rawConfig := d.GetRawConfig()
	if d.HasChange("authorize_access_policy_rules") {
		update = true
		err := ecdService.setAuthAccessPolicyRules(d, request, "authorize_access_policy_rules")
		if err != nil {
			return errs.WrapError(err)
		}
	}
	if d.HasChange("authorize_security_policy_rules") {
		update = true
		err := ecdService.setAuthSecurityPolicyRules(d, request, "authorize_security_policy_rules")
		if err != nil {
			return errs.WrapError(err)
		}
	}
	if d.HasChange("cidr_ip") {
		update = true
		if v, ok := d.GetOk("cidr_ip"); ok {
			request["AuthorizeSecurityPolicyRule.*.CidrIp"] = v
		}
	}
	if d.HasChange("description") {
		update = true
		if v, ok := d.GetOk("description"); ok {
			request["AuthorizeSecurityPolicyRule.*.Description"] = v
		}
	}
	if d.HasChange("ip_protocol") {
		update = true
		if v, ok := d.GetOk("ip_protocol"); ok {
			request["AuthorizeSecurityPolicyRule.*.IpProtocol"] = v
		}
	}
	if d.HasChange("policy") {
		update = true
		if v, ok := d.GetOk("policy"); ok {
			request["AuthorizeSecurityPolicyRule.*.Policy"] = v
		}
	}
	if d.HasChange("port_range") {
		update = true
		if v, ok := d.GetOk("port_range"); ok {
			request["AuthorizeSecurityPolicyRule.*.PortRange"] = v
		}
	}
	if d.HasChange("priority") {
		update = true
		if v, ok := d.GetOk("priority"); ok {
			request["AuthorizeSecurityPolicyRule.*.Priority"] = v
		}
	}
	if d.HasChange("type") {
		update = true
		if v, ok := d.GetOk("type"); ok {
			request["AuthorizeSecurityPolicyRule.*.Type"] = v
		}
	}
	if d.HasChange("clipboard") {
		update = true
		if v, ok := d.GetOk("clipboard"); ok {
			request["Clipboard"] = v
		}
	}
	if d.HasChange("domain_resolve_rule_type") || d.HasChange("domain_list") || d.HasChange("domain_resolve_rule") {
		update = true
	}
	// The raw config records which domain access control fields the user actually
	// manages; Optional(+Computed) fields return the computed state value when the
	// configuration omits them.
	domainConfig := expandPolicyGroupDomainConfig(rawConfig, d)
	if domainConfig.ruleType != "" {
		// The wire enum is uppercase OFF/ON; the schema and state use lowercase.
		request["DomainResolveRuleType"] = strings.ToUpper(domainConfig.ruleType)
	}
	if domainConfig.domainList != "" {
		request["DomainList"] = domainConfig.domainList
	}
	// Configured rules are sent on every update: ModifyPolicyGroup keeps omitted
	// parameters, so an update that left the rules out of the request would keep
	// the server-side rules in place; sending the configured set keeps the server
	// aligned with the config.
	for i, rule := range domainConfig.rules {
		request["DomainResolveRule."+fmt.Sprint(i+1)+".Domain"] = rule["domain"]
		request["DomainResolveRule."+fmt.Sprint(i+1)+".Policy"] = rule["policy"]
		request["DomainResolveRule."+fmt.Sprint(i+1)+".Description"] = rule["description"]
	}
	// Removing a field the state recorded leaves a plan diff. ModifyPolicyGroup keeps
	// omitted parameters, so an empty DomainList cannot clear the server-side
	// configuration and DomainResolveRuleType=OFF is the only way to remove it. The
	// change signal is required: a policy group whose state never recorded a domain
	// form (never managed here, never imported) also reads back "on" with no domain
	// configuration, and an unrelated update must not pick up its out-of-band
	// configuration. State that did record a form — managed here or adopted from the
	// server at import — plans a visible removal diff, and applying it sends OFF
	// deliberately.
	domainConfigured := domainConfig.ruleType != "" || domainConfig.domainList != "" || len(domainConfig.rules) > 0
	if !domainConfigured && shouldDisableDomainAccessControl(
		d.HasChange("domain_list"), domainConfig.domainList != "",
		d.HasChange("domain_resolve_rule"), len(domainConfig.rules) > 0,
	) {
		request["DomainResolveRuleType"] = "OFF"
		update = true
	}
	if d.HasChange("html_access") {
		update = true
		if v, ok := d.GetOk("html_access"); ok {
			request["Html5Access"] = v
		}
	}
	if d.HasChange("html_file_transfer") {
		update = true
		if v, ok := d.GetOk("html_file_transfer"); ok {
			request["Html5FileTransfer"] = v
		}
	}
	if d.HasChange("local_drive") {
		update = true
		if v, ok := d.GetOk("local_drive"); ok {
			request["LocalDrive"] = v
		}
	}
	if d.HasChange("policy_group_name") {
		update = true
		if v, ok := d.GetOk("policy_group_name"); ok {
			request["Name"] = v
		}
	}
	if d.HasChange("usb_redirect") {
		update = true
		if v, ok := d.GetOk("usb_redirect"); ok {
			request["UsbRedirect"] = v
		}
	}
	if d.HasChange("visual_quality") {
		update = true
	}
	if v, ok := d.GetOk("visual_quality"); ok {
		request["VisualQuality"] = v
	}
	if d.HasChange("watermark") {
		update = true
	}
	if v, ok := d.GetOk("watermark"); ok {
		request["Watermark"] = v
	}
	if d.HasChange("watermark_transparency") {
		update = true
	}
	if v, ok := d.GetOk("watermark_transparency"); ok {
		request["WatermarkTransparency"] = v
	}
	if d.HasChange("watermark_type") {
		update = true
	}
	if v, ok := d.GetOk("watermark_type"); ok {
		request["WatermarkType"] = v
	}
	if d.HasChange("recording") {
		update = true
	}
	if v, ok := d.GetOk("recording"); ok {
		request["Recording"] = v
	}
	if d.HasChange("recording_start_time") {
		update = true
	}
	if v, ok := d.GetOk("recording_start_time"); ok {
		request["RecordingStartTime"] = v
	}
	if d.HasChange("recording_end_time") {
		update = true
	}
	if v, ok := d.GetOk("recording_end_time"); ok {
		request["RecordingEndTime"] = v
	}
	if d.HasChange("recording_fps") {
		update = true
	}
	if v, ok := d.GetOk("recording_fps"); ok {
		request["RecordingFps"] = v
	}
	if d.HasChange("camera_redirect") {
		update = true
	}
	if v, ok := d.GetOk("camera_redirect"); ok {
		request["CameraRedirect"] = v
	}
	if d.HasChange("recording_expires") {
		update = true
	}
	if v, ok := d.GetOk("recording_expires"); ok {
		request["RecordingExpires"] = v
	}
	if d.HasChange("internet_communication_protocol") {
		update = true
	}
	if v, ok := d.GetOk("internet_communication_protocol"); ok {
		request["InternetCommunicationProtocol"] = v
	}
	if d.HasChange("recording_duration") {
		update = true
	}
	if v, ok := d.GetOk("recording_duration"); ok {
		request["RecordingDuration"] = v
	}
	if d.HasChange("client_types") {
		update = true
	}
	// Send the configured client types on every update, like the sibling
	// attributes, so every request carries the configured restrictions instead
	// of leaving the client types to the API's handling of an omitted list.
	// Removing every block sends none, which is the documented "allow all
	// client types" form.
	if v, ok := d.GetOk("client_types"); ok {
		for clientTypesPtr, clientTypes := range v.(*schema.Set).List() {
			clientTypesArg, ok := clientTypes.(map[string]interface{})
			if !ok {
				continue
			}
			request["ClientType."+fmt.Sprint(clientTypesPtr+1)+".ClientType"] = clientTypesArg["client_type"]
			// Status is optional: send it only when configured, otherwise the API
			// keeps the type's default. The wire enum is uppercase OFF/ON.
			if status, ok := clientTypesArg["status"].(string); ok && status != "" {
				request["ClientType."+fmt.Sprint(clientTypesPtr+1)+".Status"] = strings.ToUpper(status)
			}
		}
	}
	if d.HasChange("end_user_apply_admin_coordinate") {
		update = true
	}
	if v, ok := d.GetOk("end_user_apply_admin_coordinate"); ok {
		request["EndUserApplyAdminCoordinate"] = v
	}
	if d.HasChange("end_user_group_coordinate") {
		update = true
	}
	if v, ok := d.GetOk("end_user_group_coordinate"); ok {
		request["EndUserGroupCoordinate"] = v
	}
	if update {
		action := "ModifyPolicyGroup"
		wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
		err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
			response, err = client.RpcPost("ecd", "2020-09-30", action, nil, request, false)
			if err != nil {
				if helper.NeedRetry(err) || errs.IsExpectedErrors(err, []string{"InvalidPolicyStatus.Modification"}) {
					wait()
					return retry.RetryableError(err)
				}
				return retry.NonRetryableError(err)
			}
			return nil
		})
		logs.AddDebug(action, response, request)
		if err != nil {
			return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
		}
		// ModifyPolicyGroup returns before DescribePolicyGroups reflects the new
		// domain access control, and the echo can flicker: a transient target value
		// that reverts to the pre-modification value while the old domain_list/rules
		// are still echoed. Poll until every payload this apply asked for — the
		// DomainResolveRuleType switch, the legacy domain_list, and the structured
		// rules — is observed twice in a row so the subsequent Read does not
		// repopulate state from a stale echo: waiting on the switch alone would let
		// a content write that is still unconverged through, and Read would then
		// record values the plan does not expect. Payloads the request did not
		// carry are not compared. The poll shares the declared Update timeout with
		// the write above instead of a fixed budget that expires independently of
		// it. The wait is best-effort: a slow-but-eventual convergence does not
		// fail the update, and on timeout for an OFF the removed form is
		// re-recorded below so the change remains a visible plan diff and the next
		// apply re-sends the request.
		requestedType, _ := request["DomainResolveRuleType"].(string)
		target := strings.ToLower(requestedType)
		requestedList, _ := request["DomainList"].(string)
		if target != "" || requestedList != "" || len(domainConfig.rules) > 0 {
			if err := waitPolicyGroupDomainConvergence(ecdService, d.Id(), d.Timeout(schema.TimeoutUpdate), target, requestedList, domainConfig.rules); err != nil {
				log.Printf("[WARN] ecd_policy_group: domain access control did not converge: %s", err)
				if target == "off" {
					// An OFF that has not converged must not look complete in state.
					// Re-record the removed form so the trailing Read mirrors the
					// stale server echo, the removal stays a visible plan diff, and
					// the next apply re-sends the OFF. Refresh no longer adopts
					// server-only rules, so this is what keeps an unconverged removal
					// visible until the server actually converges; an explicit
					// configured "off" keeps its diff through the switch mirror
					// either way.
					oldList, _ := d.GetChange("domain_list")
					oldRules, _ := d.GetChange("domain_resolve_rule")
					if list, ok := oldList.(string); ok {
						d.Set("domain_list", list)
					}
					if rules, ok := oldRules.(*schema.Set); ok && rules.Len() > 0 {
						if err := d.Set("domain_resolve_rule", rules.List()); err != nil {
							return errs.WrapError(err)
						}
					}
				}
			}
		}
	}
	return resourcePolicyGroupRead(d, meta)
}
func resourcePolicyGroupDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	action := "DeletePolicyGroups"
	var response map[string]interface{}
	var err error
	request := map[string]interface{}{
		"PolicyGroupId": []string{d.Id()},
	}

	request["RegionId"] = client.RegionId
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutDelete), func() *retry.RetryError {
		response, err = client.RpcPost("ecd", "2020-09-30", action, nil, request, false)
		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	logs.AddDebug(action, response, request)
	if err != nil {
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
	}
	return nil
}
