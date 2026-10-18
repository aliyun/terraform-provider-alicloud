package dataworks

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/alibabacloud-go/tea/tea"
)

// recordedCall captures a single RpcPost invocation so tests can assert the
// request shape and the call sequence (pagination, action selection).
type recordedCall struct {
	product   string
	version   string
	action    string
	query     map[string]interface{}
	body      map[string]interface{}
	autoRetry bool
}

// fakeCaller is a scriptable rpcPostCaller. It records every call and dispatches
// to a user-provided handler keyed by call index, so a test can script a full
// pagination run or a failure path without touching the network.
type fakeCaller struct {
	calls   []recordedCall
	handler func(idx int, c recordedCall) (map[string]interface{}, error)
}

func (f *fakeCaller) RpcPost(apiProductCode string, apiVersion string, apiName string, query map[string]interface{}, body map[string]interface{}, autoRetry bool) (map[string]interface{}, error) {
	c := recordedCall{
		product:   apiProductCode,
		version:   apiVersion,
		action:    apiName,
		query:     query,
		body:      body,
		autoRetry: autoRetry,
	}
	idx := len(f.calls)
	f.calls = append(f.calls, c)
	return f.handler(idx, c)
}

func runChecksResponse(items ...map[string]interface{}) map[string]interface{} {
	list := make([]interface{}, len(items))
	for i, item := range items {
		list[i] = item
	}
	return map[string]interface{}{
		"Data": map[string]interface{}{
			"RuleChecks": list,
		},
	}
}

func sampleQualityCheckResult(id, entityID, ruleID string) map[string]interface{} {
	return map[string]interface{}{
		"Id":                id,
		"EntityId":          entityID,
		"RuleId":            ruleID,
		"ProjectName":       "dw-project",
		"TableName":         "orders",
		"Op":                "gt",
		"TaskId":            "task-1",
		"CheckerName":       "row-count",
		"CheckerId":         "checker-1",
		"CheckerType":       "0",
		"BeginTime":         "2024-01-01 00:00:00",
		"EndTime":           "2024-01-01 00:10:00",
		"ExpectValue":       "100",
		"UpperValue":        "200",
		"LowerValue":        "10",
		"WarningThreshold":  "80",
		"CriticalThreshold": "50",
		"Trend":             "0",
		"CheckResult":       "1",
		"FixedCheck":        "0",
		"IsPrediction":      "0",
		"BlockType":         "0",
		"MatchExpression":   "ds=20240101",
		"ActualExpression":  "ds=20240101",
		"WhereCondition":    "id > 0",
		"Property":          "row_count",
		"MethodName":        "table",
		"DateType":          "DAY",
		"RuleName":          "row-count-rule",
		"TemplateName":      "tpl",
		"TemplateId":        "tpl-1",
		"Comment":           "",
		"ExternalId":        "ext-1",
		"ExternalType":      "node",
	}
}

// TestUnitQualityCheckResultsRead exercises the no-cloud Read path: request
// shape, action selection, pagination, ids filtering and the 403 failure path.
func TestUnitQualityCheckResultsRead(t *testing.T) {
	ctx := context.Background()
	d := &qualityCheckResultsDataSource{}

	t.Run("success by entity maps fields and stops on partial page", func(t *testing.T) {
		fc := &fakeCaller{
			handler: func(idx int, c recordedCall) (map[string]interface{}, error) {
				if c.action != "ListQualityResultsByEntity" {
					return nil, fmt.Errorf("unexpected action %q", c.action)
				}
				if c.product != "dataworks-public" || c.version != "2020-05-18" {
					return nil, fmt.Errorf("unexpected product/version %q/%q", c.product, c.version)
				}
				if c.body["ProjectName"] != "dw-project" || c.body["StartDate"] != "2024-01-01 00:00:00" || c.body["EndDate"] != "2024-01-02 00:00:00" {
					return nil, fmt.Errorf("unexpected request body: %v", c.body)
				}
				if c.body["EntityId"] != "entity-1" {
					return nil, fmt.Errorf("EntityId not forwarded: %v", c.body["EntityId"])
				}
				if c.body["RuleId"] != nil {
					return nil, fmt.Errorf("RuleId leaked into entity request: %v", c.body["RuleId"])
				}
				if c.body["PageSize"] != 10 || c.body["PageNumber"] != 1 {
					return nil, fmt.Errorf("unexpected page sizing: %v/%v", c.body["PageSize"], c.body["PageNumber"])
				}
				return runChecksResponse(
					sampleQualityCheckResult("1001", "entity-1", "rule-9"),
					sampleQualityCheckResult("1002", "entity-1", "rule-9"),
				), nil
			},
		}

		items, err := d.listQualityCheckResults(ctx, fc, qualityCheckParams{
			ProjectName: "dw-project",
			StartDate:   "2024-01-01 00:00:00",
			EndDate:     "2024-01-02 00:00:00",
			EntityId:    "entity-1",
			PageSize:    10,
			PageNumber:  1,
		})
		if err != nil {
			t.Fatalf("expected success, got error: %v", err)
		}
		if len(items) != 2 {
			t.Fatalf("expected 2 results, got %d", len(items))
		}
		if len(fc.calls) != 1 {
			t.Fatalf("expected 1 call (partial page stops), got %d", len(fc.calls))
		}
		model := mapQualityCheckResult(items[0])
		if model.Id.ValueString() != "1001" {
			t.Errorf("expected Id 1001, got %q", model.Id.ValueString())
		}
		if model.QualityCheckResultId.ValueString() != "1001" {
			t.Errorf("expected quality_check_result_id 1001, got %q", model.QualityCheckResultId.ValueString())
		}
		if model.EntityId.ValueString() != "entity-1" {
			t.Errorf("expected EntityId entity-1, got %q", model.EntityId.ValueString())
		}
		if model.ProjectName.ValueString() != "dw-project" {
			t.Errorf("expected ProjectName dw-project, got %q", model.ProjectName.ValueString())
		}
	})

	t.Run("pagination walks pages until a partial page", func(t *testing.T) {
		fc := &fakeCaller{
			handler: func(idx int, c recordedCall) (map[string]interface{}, error) {
				page := c.body["PageNumber"].(int)
				switch page {
				case 1:
					return runChecksResponse(
						sampleQualityCheckResult("1", "e", "r"),
						sampleQualityCheckResult("2", "e", "r"),
					), nil
				case 2:
					return runChecksResponse(
						sampleQualityCheckResult("3", "e", "r"),
					), nil
				default:
					return nil, fmt.Errorf("unexpected page %d", page)
				}
			},
		}

		items, err := d.listQualityCheckResults(ctx, fc, qualityCheckParams{
			ProjectName: "dw-project",
			StartDate:   "2024-01-01 00:00:00",
			EndDate:     "2024-01-02 00:00:00",
			RuleId:      "rule-9",
			PageSize:    2,
			PageNumber:  1,
		})
		if err != nil {
			t.Fatalf("expected success, got error: %v", err)
		}
		if len(items) != 3 {
			t.Fatalf("expected 3 results across pages, got %d", len(items))
		}
		if len(fc.calls) != 2 {
			t.Fatalf("expected 2 paginated calls, got %d", len(fc.calls))
		}
		if fc.calls[0].action != "ListQualityResultsByRule" {
			t.Errorf("expected rule action when RuleId set, got %q", fc.calls[0].action)
		}
		if fc.calls[0].body["RuleId"] != "rule-9" {
			t.Errorf("expected RuleId forwarded, got %v", fc.calls[0].body["RuleId"])
		}
		if fc.calls[0].body["EntityId"] != nil {
			t.Errorf("EntityId leaked into rule request: %v", fc.calls[0].body["EntityId"])
		}
		if fc.calls[1].body["PageNumber"] != 2 {
			t.Errorf("expected second call PageNumber=2, got %v", fc.calls[1].body["PageNumber"])
		}
	})

	t.Run("ids filter keeps only requested result ids", func(t *testing.T) {
		fc := &fakeCaller{
			handler: func(idx int, c recordedCall) (map[string]interface{}, error) {
				return runChecksResponse(
					sampleQualityCheckResult("1", "e", "r"),
					sampleQualityCheckResult("2", "e", "r"),
					sampleQualityCheckResult("3", "e", "r"),
				), nil
			},
		}

		items, err := d.listQualityCheckResults(ctx, fc, qualityCheckParams{
			ProjectName: "dw-project",
			StartDate:   "2024-01-01 00:00:00",
			EndDate:     "2024-01-02 00:00:00",
			EntityId:    "e",
			PageSize:    10,
			PageNumber:  1,
			IDs:         []string{"2"},
		})
		if err != nil {
			t.Fatalf("expected success, got error: %v", err)
		}
		if len(items) != 1 {
			t.Fatalf("expected 1 filtered result, got %d", len(items))
		}
		if got := fmt.Sprint(items[0]["Id"]); got != "2" {
			t.Errorf("expected filtered id 2, got %q", got)
		}
	})

	t.Run("403 permission error propagates as a failure path", func(t *testing.T) {
		fc := &fakeCaller{
			handler: func(idx int, c recordedCall) (map[string]interface{}, error) {
				return nil, &tea.SDKError{
					Code:       tea.String("104000401"),
					Message:    tea.String("projectName is not relative with your account id"),
					StatusCode: tea.Int(403),
				}
			},
		}

		_, err := d.listQualityCheckResults(ctx, fc, qualityCheckParams{
			ProjectName: "dw-project",
			StartDate:   "2024-01-01 00:00:00",
			EndDate:     "2024-01-02 00:00:00",
			EntityId:    "entity-1",
			PageSize:    10,
			PageNumber:  1,
		})
		if err == nil {
			t.Fatalf("expected a 403 error, got nil")
		}
		if len(fc.calls) != 1 {
			t.Fatalf("expected exactly 1 call before failure, got %d", len(fc.calls))
		}
		if fc.calls[0].action != "ListQualityResultsByEntity" {
			t.Errorf("expected entity action, got %q", fc.calls[0].action)
		}
		msg := err.Error()
		if !strings.Contains(msg, "ListQualityResultsByEntity") || !strings.Contains(msg, "104000401") {
			t.Errorf("expected error to surface action and code, got %q", msg)
		}
	})
}
