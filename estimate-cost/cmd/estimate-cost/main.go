// estimate-cost turns a Terraform plan JSON into one or more CC API
// GetApiPrice calls and prints an estimated cost table.
//
// Usage:
//
//	estimate-cost [--mappings <dir>] [--dry-run] <plan.json>
//
// Authentication: AK/SK is read from ALIBABA_CLOUD_ACCESS_KEY_ID/SECRET or
// from ALICLOUD_ACCESS_KEY/SECRET (see internal/ccapi.NewFromEnv).
//
// Strategy:
//
//	create resources: call the matching create target once.
//	update resources:
//	  · If the mapping has a matching update target (e.g. ModifyPrepayInstance-
//	    Spec for PrePaid), call it directly. CC API returns delta mode with
//	    calculatedAmount = the actual upgrade-order amount.
//	  · Otherwise (typically PostPaid spec changes), call the create target
//	    twice — once against plan.before, once against plan.after — and take
//	    the difference. CC API has no native delta mode for these APIs, so
//	    the tool does the subtraction itself.
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/aliyun/terraform-provider-alicloud/estimate-cost/internal/ccapi"
	"github.com/aliyun/terraform-provider-alicloud/estimate-cost/internal/embedmappings"
	"github.com/aliyun/terraform-provider-alicloud/estimate-cost/internal/mapping"
)

type planJSON struct {
	ResourceChanges []resourceChange `json:"resource_changes"`
	Configuration   struct {
		ProviderConfig map[string]struct {
			Expressions map[string]struct {
				ConstantValue interface{} `json:"constant_value"`
			} `json:"expressions"`
		} `json:"provider_config"`
	} `json:"configuration"`
}

type resourceChange struct {
	Address      string `json:"address"`
	Type         string `json:"type"`
	ProviderName string `json:"provider_name"`
	Change       struct {
		Actions      []string               `json:"actions"`
		Before       map[string]interface{} `json:"before"`
		After        map[string]interface{} `json:"after"`
		AfterUnknown map[string]interface{} `json:"after_unknown"`
	} `json:"change"`
}

// row is a single output row. amount already reflects one of three semantics:
// "create price", "delta price" or "diff price".
type row struct {
	addr     string
	amount   float64
	currency string
	mode     string // "create" | "delta" | "diff(after-before)"
	note     string
	skip     bool   // when true, amount is not added to the total
	skipNote string // displayed reason when skip=true
}

const headerFmt = "%-44s %-10s %-9s %-22s %s\n"
const rowFmt = "%-44s %-10.2f %-9s %-22s %s\n"
const skipFmt = "%-44s %-10s %-9s %-22s %s\n"

func main() {
	mappingDir := flag.String("mappings", "", "(dev only) override embedded mappings with a local directory")
	dryRun := flag.Bool("dry-run", false, "build request bodies but do not call CC API")
	flag.Parse()
	if flag.NArg() < 1 {
		fmt.Fprintln(os.Stderr, "Usage: estimate-cost [--mappings dir] [--dry-run] <plan.json>")
		os.Exit(2)
	}

	plan, err := loadPlan(flag.Arg(0))
	if err != nil {
		fail("load plan: %v", err)
	}
	// Default to the embedded mappings; --mappings is for development only.
	var mappings map[string]*mapping.File
	if *mappingDir != "" {
		mappings, err = mapping.Load(*mappingDir)
	} else {
		mappings, err = mapping.LoadFS(embedmappings.FS, embedmappings.Dir)
	}
	if err != nil {
		fail("load mappings: %v", err)
	}
	providerCtx := extractProviderContext(plan)

	var client *ccapi.Client
	if !*dryRun {
		if client, err = ccapi.NewFromEnv(); err != nil {
			fail("init CC API client: %v", err)
		}
	}

	fmt.Printf(headerFmt, "RESOURCE", "AMOUNT", "CURRENCY", "MODE", "NOTE")
	var total float64
	var currency string

	for _, rc := range plan.ResourceChanges {
		if !hasAny(rc.Change.Actions, "create", "update") {
			continue
		}
		r := priceResource(rc, mappings, providerCtx, client, *dryRun)
		if r.skip {
			fmt.Printf(skipFmt, rc.Address, "-", "-", "-", r.skipNote)
			continue
		}
		fmt.Printf(rowFmt, rc.Address, r.amount, r.currency, r.mode, r.note)
		total += r.amount
		if r.currency != "" {
			currency = r.currency
		}
	}
	fmt.Printf(rowFmt, "── TOTAL ──", total, currency, "", "")
}

// priceResource is the per-resource pricing routine. It returns a row to
// print; skip=true means the resource could not be priced.
func priceResource(
	rc resourceChange,
	mappings map[string]*mapping.File,
	providerCtx map[string]map[string]interface{},
	client *ccapi.Client,
	dryRun bool,
) row {
	m, ok := mappings[rc.Type]
	if !ok {
		return skipRow("no mapping, skipped")
	}
	if !m.Billable {
		return skipRow("marked non-billable")
	}

	ctx := providerCtx[shortProviderName(rc.ProviderName)]
	afterCfg := mergeContext(ctx, rc.Change.After)
	beforeCfg := mergeContext(ctx, rc.Change.Before)
	isUpdate := hasAny(rc.Change.Actions, "update")

	if isUpdate {
		// ① A single update may trigger several OpenAPI calls — e.g. a
		//    PrePaid update that changes both instance_type and
		//    system_disk_size will call ModifyPrepayInstanceSpec and
		//    ResizeDisk separately, producing two billing orders. Price
		//    every matching update target and aggregate.
		if ts := m.PickTargets(rc.Change.Actions, afterCfg, beforeCfg); len(ts) > 0 {
			return quoteAndAggregate(ts, afterCfg, rc.Change.Before, rc.Change.AfterUnknown, client, dryRun)
		}
		// ② No update target matched → fall back to the create target and
		//    call it twice (before / after) for a diff (typical for
		//    PostPaid spec changes that CC API does not yet model as delta).
		ct := m.PickTargetByAction("create", afterCfg, beforeCfg)
		if ct == nil {
			return skipRow("update: no modify target and no create-target fallback")
		}
		return quoteDiff(ct, beforeCfg, afterCfg, rc.Change.Before, rc.Change.AfterUnknown, client, dryRun)
	}

	// create path
	t := m.PickTarget(rc.Change.Actions, afterCfg, beforeCfg)
	if t == nil {
		return skipRow("no target matched the action")
	}
	return quoteOnce(t, afterCfg, rc.Change.Before, rc.Change.AfterUnknown, client, dryRun, "create", "create")
}

// quoteAndAggregate prices every update target and sums them. Used for
// "one update triggers several independent OpenAPI calls → several billing
// orders" — e.g. PrePaid changing both instance_type and system_disk_size
// at once. For each target, CC API returns that change's real delta amount;
// summing them yields the total delta.
//
// A single failed or zero target is not fatal — we record the reason in the
// note column and continue.
func quoteAndAggregate(
	targets []*mapping.PricingTarget,
	after, before, afterUnknown map[string]interface{},
	client *ccapi.Client,
	dryRun bool,
) row {
	var total float64
	var currency string
	var notes []string

	for _, t := range targets {
		// Skip targets whose required inputs are unknown (e.g. disk resize
		// when system_disk_id is missing from the plan).
		if u := requiredUnknown(t, afterUnknown); u != "" {
			notes = append(notes, fmt.Sprintf("%s:skipped(unknown %s)", t.Name, u))
			continue
		}
		params, err := mapping.Build(t, after, before)
		if err != nil {
			// Required field missing (e.g. disk resize but disk_id not in
			// plan) — skip this target only; other targets continue.
			notes = append(notes, fmt.Sprintf("%s:skipped(%v)", t.Name, err))
			continue
		}
		req := &ccapi.Request{
			PopCode: t.OpenAPI.PopCode, PopVersion: t.OpenAPI.PopVersion,
			APIName: t.OpenAPI.APIName, Parameters: params,
		}
		if dryRun {
			body, _ := json.MarshalIndent(req, "", "  ")
			fmt.Printf("---- dry-run %s request body ----\n", t.Name)
			fmt.Println(string(body))
			continue
		}
		resp, err := client.Quote(req)
		if err != nil {
			notes = append(notes, fmt.Sprintf("%s:failed(%v)", t.Name, err))
			continue
		}
		if resp.Price == nil || !resp.Price.Success {
			ec, em := "", ""
			if resp.Price != nil {
				ec, em = resp.Price.ErrorCode, resp.Price.ErrorMessage
			}
			notes = append(notes, fmt.Sprintf("%s:%s %s", t.Name, ec, em))
			continue
		}
		amount, cur := resp.Price.FinalAmount()
		total += amount
		if cur != "" {
			currency = cur
		}
		notes = append(notes, fmt.Sprintf("%s=%.2f", t.Name, amount))
	}

	if dryRun {
		return skipRow("dry-run")
	}
	return row{
		amount:   total,
		currency: currency,
		mode:     "delta×N",
		note:     strings.Join(notes, ", "),
	}
}

// quoteOnce makes a single CC API call and parses the amount. mode is the
// label shown in the MODE column.
func quoteOnce(
	t *mapping.PricingTarget,
	after, before, afterUnknown map[string]interface{},
	client *ccapi.Client,
	dryRun bool,
	mode, note string,
) row {
	if u := requiredUnknown(t, afterUnknown); u != "" {
		return skipRow("input " + u + " unknown, priced after apply")
	}
	params, err := mapping.Build(t, after, before)
	if err != nil {
		return skipRow("mapping failed: " + err.Error())
	}
	req := &ccapi.Request{
		PopCode:    t.OpenAPI.PopCode,
		PopVersion: t.OpenAPI.PopVersion,
		APIName:    t.OpenAPI.APIName,
		Parameters: params,
	}
	if dryRun {
		body, _ := json.MarshalIndent(req, "", "  ")
		fmt.Println("---- dry-run request body ----")
		fmt.Println(string(body))
		return skipRow("dry-run")
	}
	resp, err := client.Quote(req)
	if err != nil {
		return skipRow("CC API call failed: " + err.Error())
	}
	if resp.Price == nil || !resp.Price.Success {
		msg := "upstream business error"
		if resp.Price != nil {
			msg = resp.Price.ErrorCode + ": " + resp.Price.ErrorMessage
		}
		return skipRow(msg)
	}
	amount, cur := resp.Price.FinalAmount()
	// If CC API's pricingMode disagrees with our caller's mode label, prefer
	// CC API's — it carries a more precise semantics (for example a PrePaid
	// update naturally returns "delta", which we want to surface as-is).
	actualMode := resp.Price.PricingMode
	if mode == "delta" && actualMode != "" {
		mode = actualMode
	}
	return row{amount: amount, currency: cur, mode: mode, note: note}
}

// quoteDiff prices the create target once against `before` and once against
// `after` and returns the difference. Typically used for PostPaid spec
// changes, where CC API does not (yet) expose a delta mode, so the tool
// performs the subtraction.
func quoteDiff(
	ct *mapping.PricingTarget,
	beforeCfg, afterCfg map[string]interface{},
	beforeRaw, afterUnknown map[string]interface{},
	client *ccapi.Client,
	dryRun bool,
) row {
	if u := requiredUnknown(ct, afterUnknown); u != "" {
		return skipRow("input " + u + " unknown, priced after apply")
	}
	// quote with "before" config
	beforeAmount, currency, err := quoteFor(ct, beforeCfg, beforeRaw, client, dryRun, "before")
	if err != "" {
		return skipRow("before quote failed: " + err)
	}
	// quote with "after" config
	afterAmount, cur2, err := quoteFor(ct, afterCfg, beforeRaw, client, dryRun, "after")
	if err != "" {
		return skipRow("after quote failed: " + err)
	}
	if cur2 != "" {
		currency = cur2
	}
	return row{
		amount:   afterAmount - beforeAmount,
		currency: currency,
		mode:     "diff(after-before)",
		note:     fmt.Sprintf("before=%.2f after=%.2f", beforeAmount, afterAmount),
	}
}

// quoteFor is an internal helper for quoteDiff: it runs the create target
// once against a given config. Returns (amount, currency, errMsg); a
// non-empty errMsg signals failure.
func quoteFor(
	ct *mapping.PricingTarget,
	cfg, before map[string]interface{},
	client *ccapi.Client,
	dryRun bool,
	tag string,
) (float64, string, string) {
	params, err := mapping.Build(ct, cfg, before)
	if err != nil {
		return 0, "", err.Error()
	}
	req := &ccapi.Request{
		PopCode:    ct.OpenAPI.PopCode,
		PopVersion: ct.OpenAPI.PopVersion,
		APIName:    ct.OpenAPI.APIName,
		Parameters: params,
	}
	if dryRun {
		body, _ := json.MarshalIndent(req, "", "  ")
		fmt.Printf("---- dry-run %s request body ----\n", tag)
		fmt.Println(string(body))
		return 0, "", ""
	}
	resp, err := client.Quote(req)
	if err != nil {
		return 0, "", err.Error()
	}
	if resp.Price == nil || !resp.Price.Success {
		if resp.Price != nil {
			return 0, "", resp.Price.ErrorCode + ": " + resp.Price.ErrorMessage
		}
		return 0, "", "upstream business error"
	}
	amount, cur := resp.Price.FinalAmount()
	return amount, cur, ""
}

func skipRow(note string) row {
	return row{skip: true, skipNote: note}
}

func loadPlan(path string) (*planJSON, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var p planJSON
	return &p, json.Unmarshal(raw, &p)
}

func extractProviderContext(plan *planJSON) map[string]map[string]interface{} {
	out := map[string]map[string]interface{}{}
	for name, pc := range plan.Configuration.ProviderConfig {
		ctx := map[string]interface{}{}
		for k, v := range pc.Expressions {
			if v.ConstantValue != nil {
				ctx[k] = v.ConstantValue
			}
		}
		out[name] = ctx
	}
	return out
}

func mergeContext(ctx, m map[string]interface{}) map[string]interface{} {
	out := map[string]interface{}{}
	for k, v := range ctx {
		out[k] = v
	}
	for k, v := range m {
		out[k] = v
	}
	return out
}

func requiredUnknown(t *mapping.PricingTarget, au map[string]interface{}) string {
	for name, def := range t.Params {
		if !def.Required || def.From == "" || len(def.From) < 3 {
			continue
		}
		key := def.From[2:]
		if v, ok := au[key]; ok {
			if b, _ := v.(bool); b {
				return name
			}
		}
	}
	return ""
}

func hasAny(actions []string, want ...string) bool {
	for _, a := range actions {
		for _, w := range want {
			if a == w {
				return true
			}
		}
	}
	return false
}

func shortProviderName(full string) string {
	if i := strings.LastIndex(full, "/"); i >= 0 {
		return full[i+1:]
	}
	return full
}

func fail(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, format+"\n", args...)
	os.Exit(1)
}
