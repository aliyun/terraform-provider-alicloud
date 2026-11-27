package ecd

import (
	"testing"

	"github.com/hashicorp/go-cty/cty"
)

func policyGroupDomainRules(domains ...string) cty.Value {
	ruleType := cty.Object(map[string]cty.Type{
		"domain":      cty.String,
		"policy":      cty.String,
		"description": cty.String,
	})
	if len(domains) == 0 {
		return cty.SetValEmpty(ruleType)
	}
	rules := make([]cty.Value, 0, len(domains))
	for _, domain := range domains {
		rules = append(rules, cty.ObjectVal(map[string]cty.Value{
			"domain":      cty.StringVal(domain),
			"policy":      cty.StringVal("block"),
			"description": cty.StringVal(""),
		}))
	}
	return cty.SetVal(rules)
}

func TestUnitValidatePolicyGroupDomainConfig(t *testing.T) {
	nullRules := cty.NullVal(cty.Set(cty.Object(map[string]cty.Type{
		"domain":      cty.String,
		"policy":      cty.String,
		"description": cty.String,
	})))
	unknownRules := cty.UnknownVal(cty.Set(cty.Object(map[string]cty.Type{
		"domain":      cty.String,
		"policy":      cty.String,
		"description": cty.String,
	})))
	cases := []struct {
		name     string
		typeVal  cty.Value
		rulesVal cty.Value
		wantErr  bool
	}{
		{
			name:     "type on with a rule",
			typeVal:  cty.StringVal("on"),
			rulesVal: policyGroupDomainRules("*.example.com"),
		},
		{
			name:     "type on without rules is rejected",
			typeVal:  cty.StringVal("on"),
			rulesVal: nullRules,
			wantErr:  true,
		},
		{
			name:     "type on without rules but empty set is rejected",
			typeVal:  cty.StringVal("on"),
			rulesVal: policyGroupDomainRules(),
			wantErr:  true,
		},
		{
			name:     "rule without a domain is rejected",
			typeVal:  cty.StringVal("on"),
			rulesVal: policyGroupDomainRules(""),
			wantErr:  true,
		},
		{
			name:     "type off without rules",
			typeVal:  cty.StringVal("off"),
			rulesVal: nullRules,
		},
		{
			name:     "type off with rules is rejected",
			typeVal:  cty.StringVal("off"),
			rulesVal: policyGroupDomainRules("*.example.com"),
			wantErr:  true,
		},
		{
			name:     "rules without a type are rejected",
			typeVal:  cty.NullVal(cty.String),
			rulesVal: policyGroupDomainRules("*.example.com"),
			wantErr:  true,
		},
		{
			name:     "no domain configuration at all",
			typeVal:  cty.NullVal(cty.String),
			rulesVal: nullRules,
		},
		{
			name:     "unknown type is not validated at plan time",
			typeVal:  cty.UnknownVal(cty.String),
			rulesVal: policyGroupDomainRules("*.example.com"),
		},
		{
			name:     "unknown rules are not validated at plan time",
			typeVal:  cty.StringVal("on"),
			rulesVal: unknownRules,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			err := validatePolicyGroupDomainConfig(c.typeVal, c.rulesVal)
			switch {
			case c.wantErr && err == nil:
				t.Fatalf("expected an error, got nil")
			case !c.wantErr && err != nil:
				t.Fatalf("expected no error, got: %v", err)
			}
		})
	}
}

func TestUnitShouldDisableDomainAccessControl(t *testing.T) {
	cases := []struct {
		name            string
		listChanged     bool
		listConfigured  bool
		rulesChanged    bool
		rulesConfigured bool
		want            bool
	}{
		{
			name: "unmanaged or imported policy group on an unrelated update",
			want: false,
		},
		{
			name:           "legacy domain_list removed",
			listChanged:    true,
			listConfigured: false,
			want:           true,
		},
		{
			name:           "legacy domain_list changed but still configured",
			listChanged:    true,
			listConfigured: true,
			want:           false,
		},
		{
			name:         "structured rules removed",
			rulesChanged: true,
			want:         true,
		},
		{
			name:            "structured rules changed but still configured",
			rulesChanged:    true,
			rulesConfigured: true,
			want:            false,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := shouldDisableDomainAccessControl(c.listChanged, c.listConfigured, c.rulesChanged, c.rulesConfigured)
			if got != c.want {
				t.Fatalf("expected %v, got %v", c.want, got)
			}
		})
	}
}

func TestUnitValidatePolicyGroupDomainFormSwitch(t *testing.T) {
	cases := []struct {
		name            string
		priorDomainList string
		priorRuleCount  int
		nextForm        string
		wantErr         bool
	}{
		{
			name:     "first configuration",
			nextForm: "structured",
		},
		{
			name:            "legacy value updated",
			priorDomainList: "blocked.example.com",
			nextForm:        "legacy",
		},
		{
			name:           "structured rules updated",
			priorRuleCount: 2,
			nextForm:       "structured",
		},
		{
			name:            "legacy removed",
			priorDomainList: "blocked.example.com",
			nextForm:        "",
		},
		{
			name:           "structured rules removed",
			priorRuleCount: 1,
			nextForm:       "",
		},
		{
			name:            "legacy to structured is rejected",
			priorDomainList: "blocked.example.com",
			nextForm:        "structured",
			wantErr:         true,
		},
		{
			name:           "structured to legacy is rejected",
			priorRuleCount: 1,
			nextForm:       "legacy",
			wantErr:        true,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			err := validatePolicyGroupDomainFormSwitch(c.priorDomainList, c.priorRuleCount, c.nextForm)
			switch {
			case c.wantErr && err == nil:
				t.Fatalf("expected an error, got nil")
			case !c.wantErr && err != nil:
				t.Fatalf("expected no error, got: %v", err)
			}
		})
	}
}

func policyGroupDomainRawConfig(domainList, ruleType string, rulesVal cty.Value) cty.Value {
	return cty.ObjectVal(map[string]cty.Value{
		"domain_list":              cty.StringVal(domainList),
		"domain_resolve_rule_type": cty.StringVal(ruleType),
		"domain_resolve_rule":      rulesVal,
	})
}

func TestUnitExpandPolicyGroupDomainForm(t *testing.T) {
	nullRules := cty.NullVal(cty.Set(cty.Object(map[string]cty.Type{
		"domain":      cty.String,
		"policy":      cty.String,
		"description": cty.String,
	})))
	cases := []struct {
		name      string
		rawConfig cty.Value
		want      string
	}{
		{
			name:      "explicit off is the cleanup path, not a form switch",
			rawConfig: policyGroupDomainRawConfig("", "off", nullRules),
			want:      "",
		},
		{
			name:      "explicit on manages the structured form",
			rawConfig: policyGroupDomainRawConfig("", "on", nullRules),
			want:      "structured",
		},
		{
			name:      "rules without a switch manage the structured form",
			rawConfig: policyGroupDomainRawConfig("", "", policyGroupDomainRules("*.example.com")),
			want:      "structured",
		},
		{
			name:      "domain_list manages the legacy form",
			rawConfig: policyGroupDomainRawConfig("[black:],blocked.example.com", "", nullRules),
			want:      "legacy",
		},
		{
			name:      "no domain configuration",
			rawConfig: policyGroupDomainRawConfig("", "", nullRules),
			want:      "",
		},
		{
			name:      "null configuration",
			rawConfig: cty.NullVal(cty.EmptyObject),
			want:      "",
		},
		{
			name:      "unknown configuration",
			rawConfig: cty.UnknownVal(cty.EmptyObject),
			want:      "",
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := expandPolicyGroupDomainForm(c.rawConfig); got != c.want {
				t.Fatalf("expected %q, got %q", c.want, got)
			}
		})
	}
}

func TestUnitResolvePolicyGroupDomainRuleType(t *testing.T) {
	cases := []struct {
		name       string
		priorType  string
		serverType string
		want       string
	}{
		{
			name:       "the server echo is recorded, so an off that has not converged yet stays a visible drift",
			priorType:  "off",
			serverType: "on",
			want:       "on",
		},
		{
			name:       "prior off is sticky while the server stays off",
			priorType:  "off",
			serverType: "off",
			want:       "off",
		},
		{
			name:       "prior off is sticky while the server echo is absent",
			priorType:  "off",
			serverType: "",
			want:       "off",
		},
		{
			name:       "server off is recorded for drift detection against a configured on",
			priorType:  "on",
			serverType: "off",
			want:       "off",
		},
		{
			name:       "a domain_list re-enable is not masked by the prior off",
			priorType:  "off",
			serverType: "on",
			want:       "on",
		},
		{
			name:       "server on for the first read or import",
			serverType: "on",
			want:       "on",
		},
		{
			name:       "no server value and no prior state",
			serverType: "",
			want:       "",
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := resolvePolicyGroupDomainRuleType(c.priorType, c.serverType); got != c.want {
				t.Fatalf("expected %q, got %q", c.want, got)
			}
		})
	}
}
