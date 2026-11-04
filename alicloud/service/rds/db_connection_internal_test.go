package rds

import (
	"encoding/json"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// TestUnitDBConnectionOptionalParams pins the connection-string parameters that
// are optional on both AllocateInstancePublicConnection and
// ModifyDBInstanceConnectionString: pgbouncer_port and general_group_name are
// only sent when they are configured, so an omitted value never changes the
// remote defaults.
func TestUnitDBConnectionOptionalParams(t *testing.T) {
	testCases := []struct {
		name    string
		raw     map[string]interface{}
		want    map[string]interface{}
		notWant []string
	}{
		{
			name: "both configured",
			raw: map[string]interface{}{
				"instance_id":        "rm-test123",
				"pgbouncer_port":     "6432",
				"general_group_name": "rgc-test123",
			},
			want: map[string]interface{}{
				"PGBouncerPort":    "6432",
				"GeneralGroupName": "rgc-test123",
			},
		},
		{
			name: "both omitted",
			raw: map[string]interface{}{
				"instance_id": "rm-test123",
			},
			notWant: []string{"PGBouncerPort", "GeneralGroupName"},
		},
		{
			name: "only pgbouncer port",
			raw: map[string]interface{}{
				"instance_id":    "rm-test123",
				"pgbouncer_port": "6433",
			},
			want: map[string]interface{}{
				"PGBouncerPort": "6433",
			},
			notWant: []string{"GeneralGroupName"},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			d := schema.TestResourceDataRaw(t, resourceAlicloudDBConnection().Schema, tc.raw)
			request := map[string]interface{}{}
			setDBConnectionOptionalParams(request, d)
			for k, v := range tc.want {
				if got := request[k]; got != v {
					t.Fatalf("request[%s] = %v, want %v", k, got, v)
				}
			}
			for _, k := range tc.notWant {
				if got, ok := request[k]; ok {
					t.Fatalf("request[%s] = %v, want the key to be absent", k, got)
				}
			}
		})
	}
}

func dbConnectionMockResourceData(t *testing.T) *schema.ResourceData {
	t.Helper()
	d := schema.TestResourceDataRaw(t, resourceAlicloudDBConnection().Schema, map[string]interface{}{
		"instance_id":       "rm-test",
		"connection_prefix": "abc",
		"pgbouncer_port":    "6432",
	})
	d.SetId("rm-test:abc")
	return d
}

func dbConnectionMockHandler(t *testing.T, instance map[string]interface{}, netInfo map[string]interface{}) http.HandlerFunc {
	t.Helper()
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch r.Form.Get("Action") {
		case "DescribeDBInstanceNetInfo":
			_ = json.NewEncoder(w).Encode(map[string]interface{}{
				"DBInstanceNetInfos": map[string]interface{}{
					"DBInstanceNetInfo": []interface{}{netInfo},
				},
			})
		case "DescribeDBInstanceAttribute":
			_ = json.NewEncoder(w).Encode(map[string]interface{}{
				"Items": map[string]interface{}{
					"DBInstanceAttribute": []interface{}{instance},
				},
			})
		default:
			t.Errorf("unexpected action %q", r.Form.Get("Action"))
		}
	}
}

func dbConnectionMockNetInfo(pgbouncerPort interface{}) map[string]interface{} {
	netInfo := map[string]interface{}{
		"ConnectionString": "abc.mysql.rds.aliyuncs.com",
		"Port":             "3306",
		"IPAddress":        "10.0.0.1",
	}
	if pgbouncerPort != nil {
		netInfo["PGBouncerPort"] = pgbouncerPort
	}
	return netInfo
}

// TestUnitDBConnectionPgBouncerPortReadOmission covers QA-NEW-03: on PostgreSQL
// an omitted PGBouncerPort means the configured value is unconfirmed and must
// not stay in state, whether the instance reports PgBouncer as enabled or its
// flag is unknown. Only an instance that explicitly reports PgBouncer disabled
// (or a non-PostgreSQL engine, where the parameter does not apply) keeps the
// configured value instead of producing a permanent diff.
func TestUnitDBConnectionPgBouncerPortReadOmission(t *testing.T) {
	cases := []struct {
		name          string
		engine        string
		pgbouncerFlag interface{}
		netInfoPort   interface{}
		wantPort      string
	}{
		{name: "enabled and present", engine: "PostgreSQL", pgbouncerFlag: "true", netInfoPort: "6432", wantPort: "6432"},
		{name: "enabled and omitted", engine: "PostgreSQL", pgbouncerFlag: "true", netInfoPort: nil, wantPort: ""},
		{name: "disabled and omitted", engine: "PostgreSQL", pgbouncerFlag: "false", netInfoPort: nil, wantPort: "6432"},
		{name: "unknown and omitted", engine: "PostgreSQL", pgbouncerFlag: nil, netInfoPort: nil, wantPort: ""},
		{name: "mysql and omitted", engine: "MySQL", pgbouncerFlag: nil, netInfoPort: nil, wantPort: "6432"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			instance := map[string]interface{}{"DBInstanceStatus": "Running", "Engine": tc.engine}
			if tc.pgbouncerFlag != nil {
				instance["PGBouncerEnabled"] = tc.pgbouncerFlag
			}
			client := dbInstanceMockClient(t, dbConnectionMockHandler(t, instance, dbConnectionMockNetInfo(tc.netInfoPort)))
			d := dbConnectionMockResourceData(t)
			if err := resourceAlicloudDBConnectionRead(d, client); err != nil {
				t.Fatalf("read failed: %v", err)
			}
			if got := d.Get("pgbouncer_port").(string); got != tc.wantPort {
				t.Fatalf("pgbouncer_port = %q, want %q", got, tc.wantPort)
			}
		})
	}
}

// TestUnitDBConnectionParentInstanceNotFound covers QA-NEW-09: a confirmed
// missing parent instance clears the connection from state instead of keeping a
// dangling record.
func TestUnitDBConnectionParentInstanceNotFound(t *testing.T) {
	client := dbInstanceMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch r.Form.Get("Action") {
		case "DescribeDBInstanceNetInfo":
			_ = json.NewEncoder(w).Encode(map[string]interface{}{
				"DBInstanceNetInfos": map[string]interface{}{
					"DBInstanceNetInfo": []interface{}{dbConnectionMockNetInfo("6432")},
				},
			})
		case "DescribeDBInstanceAttribute":
			w.WriteHeader(http.StatusNotFound)
			_ = json.NewEncoder(w).Encode(map[string]interface{}{
				"Code": "InvalidDBInstanceId.NotFound", "Message": "instance not found",
			})
		default:
			t.Errorf("unexpected action %q", r.Form.Get("Action"))
		}
	})
	d := dbConnectionMockResourceData(t)
	if err := resourceAlicloudDBConnectionRead(d, client); err != nil {
		t.Fatalf("read failed: %v", err)
	}
	if d.Id() != "" {
		t.Fatalf("connection must be removed from state, id = %q", d.Id())
	}
}

// TestUnitDBConnectionPgBouncerPortConvergence covers the shared Create/Update
// waiter: an omitted port is not convergence, and the waiter succeeds once the
// API reports the requested value.
func TestUnitDBConnectionPgBouncerPortConvergence(t *testing.T) {
	t.Run("eventually converged", func(t *testing.T) {
		var calls int
		client := dbInstanceMockClient(t, func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			if r.Form.Get("Action") != "DescribeDBInstanceNetInfo" {
				t.Errorf("unexpected action %q", r.Form.Get("Action"))
			}
			calls++
			port := interface{}(nil)
			if calls > 1 {
				port = "6433"
			}
			_ = json.NewEncoder(w).Encode(map[string]interface{}{
				"DBInstanceNetInfos": map[string]interface{}{
					"DBInstanceNetInfo": []interface{}{dbConnectionMockNetInfo(port)},
				},
			})
		})
		if err := waitDBConnectionPgBouncerPortConverged(RdsService{client}, "rm-test:abc", "6433", 5*time.Second); err != nil {
			t.Fatalf("waiter failed: %v", err)
		}
	})

	t.Run("never confirmed", func(t *testing.T) {
		client := dbInstanceMockClient(t, func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(map[string]interface{}{
				"DBInstanceNetInfos": map[string]interface{}{
					"DBInstanceNetInfo": []interface{}{dbConnectionMockNetInfo(nil)},
				},
			})
		})
		if err := waitDBConnectionPgBouncerPortConverged(RdsService{client}, "rm-test:abc", "6432", 2*time.Second); err == nil {
			t.Fatal("an omitted port must not be reported as converged")
		}
	})
}

// TestUnitDBConnectionPgBouncerPortSkipWhenDisabled covers the gate in front of
// the port waiter: when the instance reports PgBouncer disabled the parameter is
// not applicable, so the waiter must not block on a value the API never returns.
func TestUnitDBConnectionPgBouncerPortSkipWhenDisabled(t *testing.T) {
	client := dbInstanceMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch r.Form.Get("Action") {
		case "DescribeDBInstanceAttribute":
			_ = json.NewEncoder(w).Encode(map[string]interface{}{
				"Items": map[string]interface{}{
					"DBInstanceAttribute": []interface{}{map[string]interface{}{"Engine": "PostgreSQL", "PGBouncerEnabled": "false"}},
				},
			})
		case "DescribeDBInstanceNetInfo":
			t.Error("the port waiter must not query the connection when PgBouncer is disabled")
		default:
			t.Errorf("unexpected action %q", r.Form.Get("Action"))
		}
	})
	if err := waitDBConnectionPgBouncerPortIfEnabled(RdsService{client}, "rm-test", "rm-test:abc", "6432", 2*time.Second); err != nil {
		t.Fatalf("waiter must skip instances with PgBouncer disabled, got %v", err)
	}
}

// TestUnitDBConnectionPgBouncerPortUnknownFlag covers QA-NEW-03's remaining
// hole: an omitted PGBouncerEnabled is unknown, not "disabled". The waiter must
// resolve the flag before deciding, and must fail instead of accepting an
// unconfirmed port when the flag never becomes known.
func TestUnitDBConnectionPgBouncerPortUnknownFlag(t *testing.T) {
	t.Run("unknown then enabled", func(t *testing.T) {
		var attributeCalls int
		client := dbInstanceMockClient(t, func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			switch r.Form.Get("Action") {
			case "DescribeDBInstanceAttribute":
				attributeCalls++
				attribute := map[string]interface{}{"Engine": "PostgreSQL", "DBInstanceStatus": "CONFIG_SWITCHING"}
				if attributeCalls > 1 {
					attribute["PGBouncerEnabled"] = "true"
				}
				_ = json.NewEncoder(w).Encode(map[string]interface{}{
					"Items": map[string]interface{}{"DBInstanceAttribute": []interface{}{attribute}},
				})
			case "DescribeDBInstanceNetInfo":
				_ = json.NewEncoder(w).Encode(map[string]interface{}{
					"DBInstanceNetInfos": map[string]interface{}{
						"DBInstanceNetInfo": []interface{}{dbConnectionMockNetInfo("6432")},
					},
				})
			default:
				t.Errorf("unexpected action %q", r.Form.Get("Action"))
			}
		})
		if err := waitDBConnectionPgBouncerPortIfEnabled(RdsService{client}, "rm-test", "rm-test:abc", "6432", 5*time.Second); err != nil {
			t.Fatalf("waiter must converge once the flag becomes known, got %v", err)
		}
		if attributeCalls < 2 {
			t.Fatalf("expected the flag to be re-read, got %d calls", attributeCalls)
		}
	})

	t.Run("never known", func(t *testing.T) {
		client := dbInstanceMockClient(t, func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			if r.Form.Get("Action") != "DescribeDBInstanceAttribute" {
				t.Errorf("unexpected action %q", r.Form.Get("Action"))
			}
			_ = json.NewEncoder(w).Encode(map[string]interface{}{
				"Items": map[string]interface{}{
					"DBInstanceAttribute": []interface{}{map[string]interface{}{"Engine": "PostgreSQL", "DBInstanceStatus": "CONFIG_SWITCHING"}},
				},
			})
		})
		if err := waitDBConnectionPgBouncerPortIfEnabled(RdsService{client}, "rm-test", "rm-test:abc", "6432", 2*time.Second); err == nil {
			t.Fatal("an unconfirmed PgBouncer state must not be accepted")
		}
	})

	t.Run("non postgresql", func(t *testing.T) {
		client := dbInstanceMockClient(t, func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			if r.Form.Get("Action") != "DescribeDBInstanceAttribute" {
				t.Errorf("unexpected action %q", r.Form.Get("Action"))
			}
			_ = json.NewEncoder(w).Encode(map[string]interface{}{
				"Items": map[string]interface{}{
					"DBInstanceAttribute": []interface{}{map[string]interface{}{"Engine": "MySQL"}},
				},
			})
		})
		if err := waitDBConnectionPgBouncerPortIfEnabled(RdsService{client}, "rm-test", "rm-test:abc", "6432", 2*time.Second); err != nil {
			t.Fatalf("the parameter does not apply to other engines, got %v", err)
		}
	})
}

// TestUnitDBConnectionParentInstanceLookupError covers N-06: a parent lookup
// that fails for a reason other than NotFound must propagate through the public
// Read and must not drop the connection from state.
func TestUnitDBConnectionParentInstanceLookupError(t *testing.T) {
	client := dbInstanceMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch r.Form.Get("Action") {
		case "DescribeDBInstanceNetInfo":
			_ = json.NewEncoder(w).Encode(map[string]interface{}{
				"DBInstanceNetInfos": map[string]interface{}{
					"DBInstanceNetInfo": []interface{}{dbConnectionMockNetInfo("6432")},
				},
			})
		case "DescribeDBInstanceAttribute":
			w.WriteHeader(http.StatusForbidden)
			_ = json.NewEncoder(w).Encode(map[string]interface{}{
				"Code": "Forbidden.RAM", "Message": "not authorized",
			})
		default:
			t.Errorf("unexpected action %q", r.Form.Get("Action"))
		}
	})
	d := dbConnectionMockResourceData(t)
	if err := resourceAlicloudDBConnectionRead(d, client); err == nil {
		t.Fatal("a failed parent lookup must be returned as an error")
	}
	if d.Id() != "rm-test:abc" {
		t.Fatalf("connection id must be preserved on a failed lookup, got %q", d.Id())
	}
}

// TestUnitDBConnectionPgBouncerPortSharedBudget covers N-04: the flag wait and
// the port wait share one budget, so consuming most of it in the first stage
// must not give the second stage a full budget again.
func TestUnitDBConnectionPgBouncerPortSharedBudget(t *testing.T) {
	// The port stage only receives the budget the flag stage left behind: an
	// expired deadline yields an error instead of a fresh budget.
	limit := 15 * time.Minute
	remaining, err := pgBouncerPortRemainingBudget(time.Now().Add(2*time.Minute), limit)
	if err != nil {
		t.Fatalf("a live deadline must yield a budget, got %v", err)
	}
	if remaining <= 0 || remaining > 2*time.Minute {
		t.Fatalf("remaining budget must derive from the deadline, got %s", remaining)
	}

	_, err = pgBouncerPortRemainingBudget(time.Now().Add(-time.Second), limit)
	if err == nil {
		t.Fatal("an exhausted budget must be reported instead of starting a new wait")
	}
	if !strings.Contains(err.Error(), "budget") {
		t.Fatalf("error must explain the exhausted budget, got %v", err)
	}
}
