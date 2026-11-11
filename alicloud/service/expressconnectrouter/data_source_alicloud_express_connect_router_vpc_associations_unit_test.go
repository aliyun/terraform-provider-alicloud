package expressconnectrouter

import (
	"encoding/json"
	"net/http"
	"sync/atomic"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// TestUnitVpcAssociationsReadPaginatesToTargetID proves that the data source
// keeps paging past a non-matching association and only emits the requested id.
func TestUnitVpcAssociationsReadPaginatesToTargetID(t *testing.T) {
	var calls int32
	client := newVpcAssociationUnitClient(t, func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			t.Errorf("parse request: %v", err)
		}
		if got := r.Form.Get("Action"); got != "DescribeExpressConnectRouterAssociation" {
			t.Errorf("unexpected action %q", got)
		}
		page := atomic.AddInt32(&calls, 1)
		w.Header().Set("Content-Type", "application/json")
		switch page {
		case 1:
			if err := json.NewEncoder(w).Encode(map[string]interface{}{
				"RequestId": "unit-test",
				"NextToken": "page-2",
				"AssociationList": []map[string]interface{}{{
					"EcrId":               "ecr-1",
					"AssociationId":       "assoc-a",
					"VpcId":               "vpc-a",
					"AssociationNodeType": "VPC",
					"Status":              "ACTIVE",
					"Description":         "desc-a",
				}},
			}); err != nil {
				t.Errorf("encode page 1: %v", err)
			}
		case 2:
			if got := r.Form.Get("NextToken"); got != "page-2" {
				t.Errorf("second page NextToken = %q, want page-2", got)
			}
			if err := json.NewEncoder(w).Encode(map[string]interface{}{
				"RequestId": "unit-test",
				"AssociationList": []map[string]interface{}{{
					"EcrId":               "ecr-1",
					"AssociationId":       "assoc-b",
					"VpcId":               "vpc-b",
					"AssociationNodeType": "VPC",
					"Status":              "ACTIVE",
					"Description":         "desc-b",
				}},
			}); err != nil {
				t.Errorf("encode page 2: %v", err)
			}
		default:
			t.Errorf("unexpected extra request page %d", page)
			w.WriteHeader(http.StatusBadRequest)
			_ = json.NewEncoder(w).Encode(map[string]string{"Code": "InvalidParameter", "Message": "unexpected extra page"})
		}
	})

	d := schema.TestResourceDataRaw(t, DataSourceVpcAssociations().Schema, map[string]interface{}{
		"ecr_id": "ecr-1",
		"ids":    []interface{}{"ecr-1:assoc-b:vpc-b"},
	})
	if err := dataSourceVpcAssociationsRead(d, client); err != nil {
		t.Fatalf("read data source: %v", err)
	}
	if got := atomic.LoadInt32(&calls); got != 2 {
		t.Fatalf("DescribeExpressConnectRouterAssociation calls = %d, want 2", got)
	}
	associations, ok := d.Get("associations").([]interface{})
	if !ok || len(associations) != 1 {
		t.Fatalf("associations = %#v, want exactly the target association", d.Get("associations"))
	}
	association, ok := associations[0].(map[string]interface{})
	if !ok {
		t.Fatalf("association entry = %#v, want a map", associations[0])
	}
	if got := association["id"]; got != "ecr-1:assoc-b:vpc-b" {
		t.Errorf("association id = %v, want ecr-1:assoc-b:vpc-b", got)
	}
	if got := association["description"]; got != "desc-b" {
		t.Errorf("association description = %v, want desc-b", got)
	}
}
