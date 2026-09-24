package alicloud

import (
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"net/http"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/errors"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// newTestServerError builds a synthetic SDK ServerError with the given HTTP status
// and Alibaba Cloud error code, mirroring how the real SDK deserializes an error
// response (the error code is parsed out of the JSON body's "Code" field). It
// keeps message text free of any code substring so IsExpectedErrors' message
// matching cannot produce false positives in the assertions below.
func newTestServerError(t *testing.T, httpStatus int, code, message string) *errors.ServerError {
	t.Helper()
	body := fmt.Sprintf(`{"Code":%q,"Message":%q}`, code, message)
	err := errors.NewServerError(httpStatus, body, "")
	se, ok := err.(*errors.ServerError)
	if !ok {
		t.Fatalf("NewServerError did not return *ServerError: %T", err)
	}
	return se
}

func writeRdsError(t *testing.T, w http.ResponseWriter, status int, code string) {
	t.Helper()
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(map[string]string{"Code": code, "Message": "request could not be completed"}); err != nil {
		t.Error(err)
	}
}

func writeRdsParent(w http.ResponseWriter, status string) {
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]interface{}{"Items": map[string]interface{}{"DBInstanceAttribute": []map[string]interface{}{{"DBInstanceStatus": status, "DBInstanceId": "rm-test"}}}})
}

// Business NotFound responses can use HTTP 400 and messages that do not contain
// "NotFound" or "instance is not found". Keep the wire response independent of
// the generic HTTP 404 and message-based fallbacks.
func writeRdsHTTP400Error(t *testing.T, w http.ResponseWriter, code, message string) {
	t.Helper()
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusBadRequest)
	if err := json.NewEncoder(w).Encode(map[string]string{"Code": code, "Message": message}); err != nil {
		t.Error(err)
	}
}

func TestRdsDatabaseReadHTTP400NotFound(t *testing.T) {
	for _, tc := range []struct {
		code, message string
		gone          bool
	}{
		{"InvalidDBInstanceId.NotFound", "DBInstanceIdentifier does not refer to an existing DB instance.", true},
		{"InvalidDBName.NotFound", "The specified database does not exist.", true},
		{"InvalidAccessKeyId.NotFound", "The specified access key does not exist.", false},
		{"Unknown", "DBInstanceIdentifier does not refer to an existing DB instance.", false},
	} {
		t.Run(tc.code, func(t *testing.T) {
			var calls int32
			client := rdsAccountTestClient(t, func(w http.ResponseWriter, r *http.Request) {
				atomic.AddInt32(&calls, 1)
				if r.Form.Get("Action") != "DescribeDatabases" || r.Form.Get("DBInstanceId") != "rm-test" || r.Form.Get("DBName") != "testdb" {
					t.Errorf("unexpected query: %v", r.Form)
				}
				writeRdsHTTP400Error(t, w, tc.code, tc.message)
			})
			id := "rm-test:testdb"
			d := rdsTestData(t, "alicloud_db_database", id)
			err := resourceAliCloudRdsDatabaseRead(d, client)
			if tc.gone {
				if err != nil || d.Id() != "" {
					t.Fatalf("business absence must clear existing state: id=%q err=%v", d.Id(), err)
				}
			} else if err == nil || d.Id() != id {
				t.Fatalf("query failure must preserve state: id=%q err=%v", d.Id(), err)
			}
			if calls != 1 {
				t.Fatalf("unexpected retry or parent lookup: calls=%d", calls)
			}
		})
	}
}

func TestRdsDatabaseDeleteHTTP400NotFound(t *testing.T) {
	for _, stage := range []string{"before_delete", "after_deleting", "delete_request", "completion_query"} {
		t.Run(stage, func(t *testing.T) {
			var parents, deletes, databases int32
			client := rdsAccountTestClient(t, func(w http.ResponseWriter, r *http.Request) {
				switch r.Form.Get("Action") {
				case "DescribeDBInstanceAttribute":
					n := atomic.AddInt32(&parents, 1)
					if n > 2 {
						// Bound the test even if a regression starts polling an absent parent.
						writeRdsError(t, w, 403, "Forbidden.RAM")
					} else if stage == "after_deleting" && n == 1 {
						writeRdsParent(w, "Deleting")
					} else if stage == "before_delete" || stage == "after_deleting" {
						writeRdsHTTP400Error(t, w, "InvalidDBInstanceName.NotFound", "The specified DB instance name does not exist.")
					} else {
						writeRdsParent(w, "Running")
					}
				case "DeleteDatabase":
					atomic.AddInt32(&deletes, 1)
					if stage == "delete_request" {
						writeRdsHTTP400Error(t, w, "InvalidDBInstanceId.NotFound", "DBInstanceIdentifier does not refer to an existing DB instance.")
					} else {
						w.Header().Set("Content-Type", "application/json")
						_, _ = w.Write([]byte(`{}`))
					}
				case "DescribeDatabases":
					atomic.AddInt32(&databases, 1)
					writeRdsHTTP400Error(t, w, "InvalidDBInstanceId.NotFound", "DBInstanceIdentifier does not refer to an existing DB instance.")
				default:
					t.Errorf("unexpected action: %s", r.Form.Get("Action"))
				}
			})
			res := resourceAliCloudRdsDatabase()
			res.Timeouts.Delete = schema.DefaultTimeout(3 * time.Second)
			d := res.Data(nil)
			d.SetId("rm-test:testdb")
			if err := res.Delete(d, client); err != nil {
				t.Fatalf("confirmed absence must complete deletion: %v", err)
			}
			wantParents, wantDeletes, wantDatabases := int32(1), int32(0), int32(0)
			if stage == "after_deleting" {
				wantParents = 2
			}
			if stage == "delete_request" || stage == "completion_query" {
				wantDeletes = 1
			}
			if stage == "completion_query" {
				wantDatabases = 1
			}
			if parents != wantParents || deletes != wantDeletes || databases != wantDatabases {
				t.Fatalf("unexpected calls: parents=%d deletes=%d databases=%d; want %d/%d/%d", parents, deletes, databases, wantParents, wantDeletes, wantDatabases)
			}
		})
	}
}

func TestRdsWaitInstanceHTTP400NotFound(t *testing.T) {
	for _, target := range []Status{Running, Deleted} {
		t.Run(string(target), func(t *testing.T) {
			var calls int32
			client := rdsAccountTestClient(t, func(w http.ResponseWriter, r *http.Request) {
				atomic.AddInt32(&calls, 1)
				writeRdsHTTP400Error(t, w, "InvalidDBInstanceName.NotFound", "The specified DB instance name does not exist.")
			})
			s := RdsService{client}
			err := s.WaitForDBInstance("rm-test", target, 0)
			if target == Deleted {
				if err != nil {
					t.Fatalf("absence must satisfy Deleted: %v", err)
				}
			} else if !NotFoundError(err) || strings.Contains(strings.ToLower(err.Error()), "timeout") {
				t.Fatalf("absence must fail Running with NotFound rather than timeout: %v", err)
			}
			if calls != 1 {
				t.Fatalf("must not keep polling an absent parent: calls=%d", calls)
			}
		})
	}
}

func rdsTestData(t *testing.T, kind, id string) *schema.ResourceData {
	t.Helper()
	res := Provider().(*schema.Provider).ResourcesMap[kind]
	d, err := schema.InternalMap(res.Schema).Data(nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	d.SetId(id)
	return d
}

func TestRdsParentErrorsDoNotConfirmChildAbsence(t *testing.T) {
	for _, code := range []string{"OperationDenied.DBInstanceStatus", "OperationDenied.ReadDBInstanceStatus", "InvalidAccessKeyId.NotFound", "Unknown"} {
		t.Run(code, func(t *testing.T) {
			var calls int32
			client := rdsAccountTestClient(t, func(w http.ResponseWriter, r *http.Request) {
				atomic.AddInt32(&calls, 1)
				writeRdsError(t, w, 403, code)
			})
			service := RdsService{client}
			_, err := service.DescribeDBInstance("rm-test")
			if err == nil || NotFoundError(err) || !rdsErrorHasCode(err, code) {
				t.Fatalf("parent error must not confirm child absence: %v", err)
			}
			if calls != 1 {
				t.Fatalf("unexpected retry or mutation: calls=%d", calls)
			}
		})
	}
}

func TestRdsChildReadParentConfirmation(t *testing.T) {
	for _, kind := range []string{"alicloud_db_database", "alicloud_db_account_privilege", "alicloud_rds_account", "alicloud_db_account"} {
		for _, parent := range []string{"Running", "gone", "empty", "auth404", "OperationDenied.DBInstanceStatus", "InvalidAccessKeyId.NotFound", "Unknown"} {
			for _, stateCode := range []string{"OperationDenied.DBInstanceStatus", "OperationDenied.ReadDBInstanceStatus"} {
				t.Run(kind+"/"+parent+"/"+stateCode, func(t *testing.T) {
					var parentCalls, childCalls int32
					client := rdsAccountTestClient(t, func(w http.ResponseWriter, r *http.Request) {
						if r.Form.Get("Action") == "DescribeDBInstanceAttribute" {
							atomic.AddInt32(&parentCalls, 1)
							switch parent {
							case "Running":
								writeRdsParent(w, "Running")
							case "gone":
								writeRdsError(t, w, 404, "InvalidDBInstanceId.NotFound")
							case "empty":
								w.Header().Set("Content-Type", "application/json")
								_, _ = w.Write([]byte(`{"Items":{"DBInstanceAttribute":[]}}`))
							case "auth404":
								writeRdsError(t, w, 404, "InvalidAccessKeyId.NotFound")
							default:
								writeRdsError(t, w, 403, parent)
							}
						} else {
							atomic.AddInt32(&childCalls, 1)
							writeRdsError(t, w, 403, stateCode)
						}
					})
					id := "rm-test:test-account"
					if kind == "alicloud_db_account_privilege" {
						id += ":ReadOnly"
					}
					d := rdsTestData(t, kind, id)
					res := Provider().(*schema.Provider).ResourcesMap[kind]
					err := res.Read(d, client)
					if parent == "gone" || parent == "empty" {
						if err != nil || d.Id() != "" {
							t.Fatalf("confirmed absence: err=%v id=%q", err, d.Id())
						}
					} else if err == nil || d.Id() != id {
						t.Fatalf("uncertain parent: err=%v id=%q", err, d.Id())
					}
					if parentCalls != 1 || childCalls != 1 {
						t.Fatalf("expected one child query and one parent check, got %d/%d", childCalls, parentCalls)
					}
				})
			}
		}
	}
}

func TestRdsDeletionWaitersConfirmParent(t *testing.T) {
	for _, waiter := range []string{"database", "revoke", "privilege"} {
		for _, parent := range []string{"Running", "gone", "auth404", "unknown404", "InvalidAccessKeyId.NotFound", "OperationDenied.DBInstanceStatus"} {
			t.Run(waiter+"/"+parent, func(t *testing.T) {
				var parentCalls int32
				client := rdsAccountTestClient(t, func(w http.ResponseWriter, r *http.Request) {
					if r.Form.Get("Action") == "DescribeDBInstanceAttribute" {
						atomic.AddInt32(&parentCalls, 1)
						if parent == "Running" {
							writeRdsParent(w, parent)
						} else if parent == "gone" {
							writeRdsError(t, w, 404, "InvalidDBInstanceId.NotFound")
						} else if parent == "auth404" {
							writeRdsError(t, w, 404, "InvalidAccessKeyId.NotFound")
						} else if parent == "unknown404" {
							writeRdsError(t, w, 404, "Unknown")
						} else {
							writeRdsError(t, w, 403, parent)
						}
					} else {
						writeGone403(t, w)
					}
				})
				service := RdsService{client}
				var err error
				switch waiter {
				case "database":
					err = service.WaitForDBDatabase("rm-test:testdb", Deleted, 1)
				case "revoke":
					err = service.WaitForAccountPrivilegeRevoked("rm-test:account:ReadOnly", "testdb", 1)
				case "privilege":
					err = service.WaitForAccountPrivilege("rm-test:account:ReadOnly", "testdb", Deleted, 1)
				}
				if (err == nil) != (parent == "gone") || parentCalls != 1 {
					t.Fatalf("err=%v parentCalls=%d", err, parentCalls)
				}
			})
		}
	}
}

func TestRdsPostgresFallbackClearsID(t *testing.T) {
	for _, parent := range []string{"gone", "Running", "auth404", "unknown404", "InvalidAccessKeyId.NotFound"} {
		t.Run(parent, func(t *testing.T) {
			client := rdsAccountTestClient(t, func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				switch r.Form.Get("Action") {
				case "DescribeAccounts":
					_ = json.NewEncoder(w).Encode(map[string]interface{}{"Accounts": map[string]interface{}{"DBInstanceAccount": []map[string]interface{}{{"DBInstanceId": "pgm-test", "AccountName": "account", "DatabasePrivileges": map[string]interface{}{"DatabasePrivilege": []interface{}{}}}}}})
				case "DescribeDatabases":
					writeGone403(t, w)
				case "DescribeDBInstanceAttribute":
					if parent == "gone" {
						writeRdsError(t, w, 404, "InvalidDBInstanceId.NotFound")
					} else if parent == "Running" {
						writeRdsParent(w, parent)
					} else if parent == "auth404" {
						writeRdsError(t, w, 404, "InvalidAccessKeyId.NotFound")
					} else if parent == "unknown404" {
						writeRdsError(t, w, 404, "Unknown")
					} else {
						writeRdsError(t, w, 403, parent)
					}
				default:
					t.Errorf("unexpected action %s", r.Form.Get("Action"))
				}
			})
			d := rdsTestData(t, "alicloud_db_account_privilege", "pgm-test:account:DBOwner")
			_ = d.Set("db_names", []string{"testdb"})
			err := resourceAlicloudDBAccountPrivilegeRead(d, client)
			if parent == "gone" {
				if err != nil || d.Id() != "" {
					t.Fatalf("err=%v id=%q", err, d.Id())
				}
			} else if err == nil || d.Id() == "" {
				t.Fatalf("err=%v id=%q", err, d.Id())
			}
		})
	}
}

func TestRdsAccountDeleteBusinessNotFound(t *testing.T) {
	for _, code := range []string{"InvalidDBInstanceId.NotFound", "InvalidAccountName.NotFound", "InvalidAccessKeyId.NotFound", "Unknown"} {
		t.Run(code, func(t *testing.T) {
			var deletes int32
			client := rdsAccountTestClient(t, func(w http.ResponseWriter, r *http.Request) {
				switch r.Form.Get("Action") {
				case "DescribeDBInstanceAttribute":
					writeRdsParent(w, "Running")
				case "DeleteAccount":
					atomic.AddInt32(&deletes, 1)
					writeRdsError(t, w, 404, code)
				default:
					t.Errorf("unexpected action %s", r.Form.Get("Action"))
					w.WriteHeader(400)
				}
			})
			err := resourceAliCloudRdsAccountDelete(rdsAccountResourceData(t), client)
			wantGone := code == "InvalidDBInstanceId.NotFound" || code == "InvalidAccountName.NotFound"
			if (err == nil) != wantGone || deletes != 1 {
				t.Fatalf("err=%v deletes=%d", err, deletes)
			}
		})
	}
}

func TestRdsAccountDeletePostcheckAuthError(t *testing.T) {
	for _, status := range []int{403, 404} {
		t.Run(fmt.Sprint(status), func(t *testing.T) {
			var parents int32
			client := rdsAccountTestClient(t, func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				switch r.Form.Get("Action") {
				case "DescribeDBInstanceAttribute":
					if atomic.AddInt32(&parents, 1) == 1 {
						writeRdsParent(w, "Running")
					} else {
						writeRdsError(t, w, status, "InvalidAccessKeyId.NotFound")
					}
				case "DeleteAccount":
					_, _ = w.Write([]byte(`{}`))
				case "DescribeAccounts":
					_, _ = w.Write([]byte(`{"Accounts":{"DBInstanceAccount":[{"AccountStatus":"Deleting"}]}}`))
				default:
					t.Errorf("unexpected action %s", r.Form.Get("Action"))
				}
			})
			d := rdsAccountResourceData(t)
			if err := resourceAliCloudRdsAccountDelete(d, client); err == nil || NotFoundError(err) || d.Id() == "" || parents != 2 {
				t.Fatalf("post-delete auth error must not be swallowed: err=%v id=%q parents=%d", err, d.Id(), parents)
			}
		})
	}
}

func TestRdsDatabaseDeleteWaitsForRelease(t *testing.T) {
	for _, gone := range []bool{false, true} {
		t.Run(fmt.Sprint(gone), func(t *testing.T) {
			var parents, deletes int32
			client := rdsAccountTestClient(t, func(w http.ResponseWriter, r *http.Request) {
				switch r.Form.Get("Action") {
				case "DescribeDBInstanceAttribute":
					n := atomic.AddInt32(&parents, 1)
					if n == 1 {
						writeRdsParent(w, "Running")
					} else if n == 2 || !gone {
						writeRdsParent(w, "Deleting")
					} else {
						writeRdsError(t, w, 404, "InvalidDBInstanceId.NotFound")
					}
				case "DeleteDatabase":
					atomic.AddInt32(&deletes, 1)
					writeGone403(t, w)
				default:
					t.Errorf("unexpected action %s", r.Form.Get("Action"))
				}
			})
			res := resourceAliCloudRdsDatabase()
			budget := 50 * time.Millisecond
			if gone {
				budget = 6 * time.Second
			}
			res.Timeouts.Delete = schema.DefaultTimeout(budget)
			d := res.Data(nil)
			d.SetId("rm-test:testdb")
			start := time.Now()
			err := res.Delete(d, client)
			if (err == nil) != gone {
				t.Fatalf("gone=%v err=%v", gone, err)
			}
			if time.Since(start) > budget+time.Second {
				t.Fatalf("delete exceeded shared budget: %s", time.Since(start))
			}
			if deletes != 1 {
				t.Fatalf("delete called while parent Deleting: %d", deletes)
			}
		})
	}
}

func TestRdsWaitForRunningFailsOnDeleting(t *testing.T) {
	client := rdsAccountTestClient(t, func(w http.ResponseWriter, r *http.Request) { writeRdsParent(w, "Deleting") })
	s := RdsService{client}
	err := s.WaitForDBInstance("rm-test", Running, 1)
	if err == nil || !strings.Contains(err.Error(), "Deleting") || strings.Contains(strings.ToLower(err.Error()), "timeout") {
		t.Fatalf("expected explicit in-progress error, got %v", err)
	}
}

func TestRdsDatabaseDeletionRequiresAbsence(t *testing.T) {
	var calls int32
	client := rdsAccountTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt32(&calls, 1)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"Databases":{"Database":[{"DBName":"testdb"}]}}`))
	})
	s := RdsService{client}
	if err := s.WaitForDBDatabase("rm-test:testdb", Deleted, 0); err == nil {
		t.Fatal("a present database is not deleted")
	}
	if calls != 1 {
		t.Fatalf("zero budget must not poll repeatedly: %d", calls)
	}
}

func TestRdsDatabaseInitiallyDeletingPreservesState(t *testing.T) {
	var calls int32
	client := rdsAccountTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Form.Get("Action") != "DescribeDBInstanceAttribute" {
			t.Errorf("unexpected mutation %s", r.Form.Get("Action"))
		}
		atomic.AddInt32(&calls, 1)
		writeRdsParent(w, "Deleting")
	})
	res := resourceAliCloudRdsDatabase()
	res.Timeouts.Delete = schema.DefaultTimeout(50 * time.Millisecond)
	d := res.Data(nil)
	d.SetId("rm-test:testdb")
	start := time.Now()
	err := res.Delete(d, client)
	if err == nil || d.Id() == "" || !strings.Contains(err.Error(), "Deleting") {
		t.Fatalf("err=%v id=%q", err, d.Id())
	}
	if time.Since(start) > time.Second || calls != 1 {
		t.Fatalf("deadline exceeded or repeated queries: calls=%d elapsed=%s", calls, time.Since(start))
	}
}

func TestRdsChildNetworkErrorPreservesState(t *testing.T) {
	client := rdsAccountTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Form.Get("Action") == "DescribeDBInstanceAttribute" {
			conn, _, err := w.(http.Hijacker).Hijack()
			if err != nil {
				t.Error(err)
				return
			}
			_ = conn.Close()
			return
		}
		writeGone403(t, w)
	})
	d := rdsAccountResourceData(t)
	if err := resourceAliCloudRdsAccountRead(d, client); err == nil || d.Id() == "" {
		t.Fatalf("network error must preserve ID: err=%v id=%q", err, d.Id())
	}
}

func TestRdsAccountWaiterPreservesAuthError(t *testing.T) {
	client := rdsAccountTestClient(t, func(w http.ResponseWriter, r *http.Request) { writeRdsError(t, w, 403, "InvalidAccessKeyId.NotFound") })
	s := RdsService{client}
	_, _, err := s.RdsAccountStateRefreshFunc("rm-test:account", nil)()
	if err == nil {
		t.Fatal("account waiter swallowed auth error")
	}
}

func TestRdsRevokeBusinessNotFoundSkipsWait(t *testing.T) {
	for _, code := range []string{"InvalidDBInstanceId.NotFound", "InvalidDB.NotFound", "InvalidAccountName.NotFound", "InvalidAccessKeyId.NotFound"} {
		t.Run(code, func(t *testing.T) {
			var calls int32
			client := rdsAccountTestClient(t, func(w http.ResponseWriter, r *http.Request) {
				atomic.AddInt32(&calls, 1)
				if r.Form.Get("Action") != "RevokeAccountPrivilege" {
					t.Errorf("unexpected call %s", r.Form.Get("Action"))
				}
				writeRdsError(t, w, 404, code)
			})
			s := RdsService{client}
			err := s.RevokeAccountPrivilege("rm-test:account:ReadOnly", "testdb")
			if (err != nil) != (code == "InvalidAccessKeyId.NotFound") || calls != 1 {
				t.Fatalf("err=%v calls=%d", err, calls)
			}
		})
	}
}

func TestRdsWaitInstanceRetriesTransient(t *testing.T) {
	var calls int32
	client := rdsAccountTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if atomic.AddInt32(&calls, 1) == 1 {
			writeRdsError(t, w, 400, "InvalidParameter")
			return
		}
		writeRdsParent(w, "Running")
	})
	s := RdsService{client}
	err := s.WaitForDBInstance("rm-test", Running, 6)
	if err != nil || calls != 2 {
		t.Fatalf("want recovery after transient error: calls=%d err=%v", calls, err)
	}
}

func TestRdsWaitDatabaseRetriesTransient(t *testing.T) {
	var calls int32
	client := rdsAccountTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if atomic.AddInt32(&calls, 1) == 1 {
			writeRdsError(t, w, 400, "InternalError")
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"Databases":{"Database":[]}}`))
	})
	s := RdsService{client}
	err := s.WaitForDBDatabase("rm-test:testdb", Deleted, 6)
	if err != nil || calls != 2 {
		t.Fatalf("want confirmed disappearance after transient error: calls=%d err=%v", calls, err)
	}
}

func TestRdsDatabaseDeleteRetriesParent(t *testing.T) {
	var parents, deletes int32
	client := rdsAccountTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch r.Form.Get("Action") {
		case "DescribeDBInstanceAttribute":
			if atomic.AddInt32(&parents, 1) == 1 {
				writeRdsError(t, w, 400, "InvalidParameter")
				return
			}
			writeRdsParent(w, "Running")
		case "DeleteDatabase":
			atomic.AddInt32(&deletes, 1)
			_, _ = w.Write([]byte(`{}`))
		case "DescribeDatabases":
			_, _ = w.Write([]byte(`{"Databases":{"Database":[]}}`))
		default:
			t.Errorf("unexpected action: %s", r.Form.Get("Action"))
		}
	})
	res := resourceAliCloudRdsDatabase()
	res.Timeouts.Delete = schema.DefaultTimeout(6 * time.Second)
	d := res.Data(nil)
	d.SetId("rm-test:testdb")
	err := res.Delete(d, client)
	if err != nil || parents != 2 || deletes != 1 {
		t.Fatalf("want delete after parent recovery: parents=%d deletes=%d err=%v", parents, deletes, err)
	}
}

func TestRdsAccountGoneSkipsParentRecheck(t *testing.T) {
	var parents, accounts int32
	client := rdsAccountTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch r.Form.Get("Action") {
		case "DescribeDBInstanceAttribute":
			if atomic.AddInt32(&parents, 1) > 1 {
				writeRdsError(t, w, 403, "OperationDenied.DBInstanceStatus")
				return
			}
			writeRdsParent(w, "Running")
		case "DeleteAccount":
			_, _ = w.Write([]byte(`{}`))
		case "DescribeAccounts":
			atomic.AddInt32(&accounts, 1)
			_, _ = w.Write([]byte(`{"Accounts":{"DBInstanceAccount":[]}}`))
		default:
			t.Errorf("unexpected action: %s", r.Form.Get("Action"))
		}
	})
	err := resourceAliCloudRdsAccountDelete(rdsAccountResourceData(t), client)
	if err != nil || parents != 1 || accounts != 1 {
		t.Fatalf("confirmed absent account must skip parent recheck: parents=%d accounts=%d err=%v", parents, accounts, err)
	}
}

func TestRdsNeedRetryWrappedError(t *testing.T) {
	raw := newTestServerError(t, 503, "ServiceUnavailable", "temporarily unavailable")
	if !NeedRetry(raw) {
		t.Fatal("raw error should be retryable")
	}
	if NeedRetry(WrapError(raw)) {
		t.Fatal("probe expected existing NeedRetry not to unwrap ComplexError")
	}
}

func TestRdsSingleQueryPreservesRetryDecision(t *testing.T) {
	for _, tc := range []struct {
		query, code string
		retryable   bool
	}{
		{"parent", "InvalidParameter", true},
		{"parent", "ServiceUnavailable", true},
		{"parent", "InternalError", false},
		{"parent", "Forbidden.RAM", false},
		{"database", "InternalError", true},
		{"database", "InvalidParameter", false},
		{"database", "Forbidden.RAM", false},
	} {
		t.Run(tc.query+"/"+tc.code, func(t *testing.T) {
			client := rdsAccountTestClient(t, func(w http.ResponseWriter, r *http.Request) { writeRdsError(t, w, 400, tc.code) })
			s := RdsService{client}
			var err error
			if tc.query == "parent" {
				_, err = s.describeDBInstance("rm-test", 0)
			} else {
				_, err = s.describeDBDatabase("rm-test:testdb", 0)
			}
			err = WrapError(WrapError(err))
			if isRdsRetryableQueryError(err) != tc.retryable || !rdsErrorHasCode(err, tc.code) {
				t.Fatalf("retryable=%v expected=%v code=%s err=%v", isRdsRetryableQueryError(err), tc.retryable, tc.code, err)
			}
		})
	}
}

func TestRdsQueryRetriesRespectDeadline(t *testing.T) {
	for _, path := range []string{"instance", "database", "delete"} {
		for _, transient := range []bool{false, true} {
			t.Run(fmt.Sprintf("%s/%t", path, transient), func(t *testing.T) {
				var calls int32
				code := "Forbidden.RAM"
				if transient {
					code = "InvalidParameter"
					if path == "database" {
						code = "InternalError"
					}
				}
				client := rdsAccountTestClient(t, func(w http.ResponseWriter, r *http.Request) {
					atomic.AddInt32(&calls, 1)
					writeRdsError(t, w, 400, code)
				})
				s := RdsService{client}
				start := time.Now()
				var err error
				switch path {
				case "instance":
					err = s.WaitForDBInstance("rm-test", Running, 1)
				case "database":
					err = s.waitForDBDatabaseUntil("rm-test:testdb", Deleted, start.Add(time.Second))
				case "delete":
					res := resourceAliCloudRdsDatabase()
					res.Timeouts.Delete = schema.DefaultTimeout(time.Second)
					d := res.Data(nil)
					d.SetId("rm-test:testdb")
					err = res.Delete(d, client)
				}
				if err == nil || !rdsErrorHasCode(err, code) {
					t.Fatalf("must retain original error: %v", err)
				}
				if strings.Contains(strings.ToLower(err.Error()), "timeout") != transient {
					t.Fatalf("unexpected timeout classification: %v", err)
				}
				maxCalls := int32(1)
				if path == "delete" && transient {
					// resource.Retry uses a 500 ms minimum backoff within this budget.
					maxCalls = 2
				}
				if calls < 1 || calls > maxCalls || time.Since(start) > 2*time.Second {
					t.Fatalf("budget exceeded or unexpected retry: calls=%d elapsed=%s", calls, time.Since(start))
				}
			})
		}
	}
}

func TestRdsDatabaseWaitRetriesParentConfirmation(t *testing.T) {
	var parents, databases int32
	client := rdsAccountTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		switch r.Form.Get("Action") {
		case "DescribeDatabases":
			atomic.AddInt32(&databases, 1)
			writeGone403(t, w)
		case "DescribeDBInstanceAttribute":
			if atomic.AddInt32(&parents, 1) == 1 {
				writeRdsError(t, w, 400, "InvalidParameter")
			} else {
				writeRdsError(t, w, 404, "InvalidDBInstanceId.NotFound")
			}
		default:
			t.Errorf("unexpected action %s", r.Form.Get("Action"))
		}
	})
	s := RdsService{client}
	err := s.WaitForDBDatabase("rm-test:testdb", Deleted, 6)
	if err != nil || parents != 2 || databases != 2 {
		t.Fatalf("parent retry decision was lost: parents=%d databases=%d err=%v", parents, databases, err)
	}
}

func TestRdsInstanceRetryTimeoutRetainsStatus(t *testing.T) {
	var calls int32
	client := rdsAccountTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if atomic.AddInt32(&calls, 1) == 1 {
			writeRdsParent(w, "Creating")
		} else {
			writeRdsError(t, w, 400, "ServiceUnavailable")
		}
	})
	s := RdsService{client}
	start := time.Now()
	err := s.WaitForDBInstance("rm-test", Running, 6)
	if err == nil || !rdsErrorHasCode(err, "ServiceUnavailable") || !strings.Contains(err.Error(), "Creating") || !strings.Contains(err.Error(), "timeout") {
		t.Fatalf("timeout must retain status and cause: %v", err)
	}
	if calls != 2 || time.Since(start) > 7*time.Second {
		t.Fatalf("unexpected calls or elapsed: %d %s", calls, time.Since(start))
	}
}

func TestRdsDatabaseDeleteSharesRetryAndWaiterBudget(t *testing.T) {
	var parents, deletes, databases int32
	client := rdsAccountTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch r.Form.Get("Action") {
		case "DescribeDBInstanceAttribute":
			if atomic.AddInt32(&parents, 1) == 1 {
				writeRdsParent(w, "Creating")
				return
			}
			writeRdsParent(w, "Running")
		case "DeleteDatabase":
			atomic.AddInt32(&deletes, 1)
			_, _ = w.Write([]byte(`{}`))
		case "DescribeDatabases":
			atomic.AddInt32(&databases, 1)
			_, _ = w.Write([]byte(`{"Databases":{"Database":[{"DBName":"testdb","DBStatus":"Deleting"}]}}`))
		default:
			t.Errorf("unexpected action: %s", r.Form.Get("Action"))
		}
	})
	res := resourceAliCloudRdsDatabase()
	budget := 1200 * time.Millisecond
	res.Timeouts.Delete = schema.DefaultTimeout(budget)
	d := res.Data(nil)
	d.SetId("rm-test:testdb")
	start := time.Now()
	err := res.Delete(d, client)
	elapsed := time.Since(start)
	if err == nil || !strings.Contains(err.Error(), "timeout") || d.Id() == "" {
		t.Fatalf("database is still present: err=%v id=%q", err, d.Id())
	}
	if parents != 2 || deletes != 1 || databases != 1 {
		t.Fatalf("unexpected request counts: parents=%d deletes=%d databases=%d", parents, deletes, databases)
	}
	// The first retry consumes at least 500 ms; the waiter must not get a fresh budget.
	if elapsed < budget || elapsed > budget+350*time.Millisecond {
		t.Fatalf("delete and waiter must share one budget: elapsed=%s budget=%s", elapsed, budget)
	}
}

func TestRdsChildReadDistinguishesParentHTTPNotFound(t *testing.T) {
	for _, kind := range []string{"alicloud_db_database", "alicloud_db_account_privilege", "alicloud_rds_account", "alicloud_db_account"} {
		for _, parentCheck := range []bool{false, true} {
			t.Run(fmt.Sprintf("%s/parent=%t", kind, parentCheck), func(t *testing.T) {
				var parentCalls, childCalls int32
				client := rdsAccountTestClient(t, func(w http.ResponseWriter, r *http.Request) {
					if r.Form.Get("Action") == "DescribeDBInstanceAttribute" {
						atomic.AddInt32(&parentCalls, 1)
						writeRdsError(t, w, 404, "Unknown")
						return
					}
					atomic.AddInt32(&childCalls, 1)
					if parentCheck {
						writeGone403(t, w)
					} else {
						writeRdsError(t, w, 404, "Unknown")
					}
				})
				id := "rm-test:account"
				if kind == "alicloud_db_account_privilege" {
					id += ":ReadOnly"
				}
				d := rdsTestData(t, kind, id)
				err := Provider().(*schema.Provider).ResourcesMap[kind].Read(d, client)
				if parentCheck {
					if err == nil || NotFoundError(err) || d.Id() != id {
						t.Fatalf("failed parent query must preserve state: err=%v id=%q", err, d.Id())
					}
				} else if err != nil || d.Id() != "" {
					t.Fatalf("direct child HTTP 404 must retain existing absence handling: err=%v id=%q", err, d.Id())
				}
				wantParents := int32(0)
				if parentCheck {
					wantParents = 1
				}
				if childCalls != 1 || parentCalls != wantParents {
					t.Fatalf("unexpected query counts: child=%d parent=%d", childCalls, parentCalls)
				}
			})
		}
	}
}

func TestRdsDatabaseDeleteDistinguishesParentHTTPNotFound(t *testing.T) {
	for _, missingAt := range []string{"DescribeDBInstanceAttribute", "DeleteDatabase", "DescribeDatabases"} {
		t.Run(missingAt, func(t *testing.T) {
			var calls int32
			client := rdsAccountTestClient(t, func(w http.ResponseWriter, r *http.Request) {
				atomic.AddInt32(&calls, 1)
				action := r.Form.Get("Action")
				if action == missingAt {
					writeRdsError(t, w, 404, "Unknown")
					return
				}
				if action == "DescribeDBInstanceAttribute" {
					writeRdsParent(w, "Running")
				} else if action == "DeleteDatabase" {
					w.Header().Set("Content-Type", "application/json")
					_, _ = w.Write([]byte(`{}`))
				} else {
					t.Errorf("unexpected action %s", action)
				}
			})
			res := resourceAliCloudRdsDatabase()
			d := res.Data(nil)
			d.SetId("rm-test:testdb")
			err := res.Delete(d, client)
			if missingAt == "DescribeDBInstanceAttribute" {
				if err == nil || NotFoundError(err) || d.Id() == "" {
					t.Fatalf("failed parent query must not complete deletion: err=%v id=%q", err, d.Id())
				}
			} else if err != nil {
				t.Fatalf("direct child HTTP 404 at %s must complete deletion: %v", missingAt, err)
			}
			wantCalls := map[string]int32{"DescribeDBInstanceAttribute": 1, "DeleteDatabase": 2, "DescribeDatabases": 3}[missingAt]
			if calls != wantCalls {
				t.Fatalf("unexpected query count: got=%d want=%d", calls, wantCalls)
			}
		})
	}
}

func TestRdsParentLookupRequiresAbsenceEvidence(t *testing.T) {
	for _, tc := range []struct {
		code string
		gone bool
	}{
		{"Unknown", false},
		{"InvalidAccessKeyId.NotFound", false},
		{"InvalidDBInstanceId.NotFound", true},
		{"InvalidDBInstanceName.NotFound", true},
		{"InvalidDBInstanceId.NotFoundError", true},
		{"Forbidden.InstanceNotFound", true},
	} {
		t.Run(tc.code, func(t *testing.T) {
			var calls int32
			client := rdsAccountTestClient(t, func(w http.ResponseWriter, r *http.Request) {
				atomic.AddInt32(&calls, 1)
				writeRdsError(t, w, 404, tc.code)
			})
			s := RdsService{client}
			_, err := s.describeRdsParentInstance("rm-test", 0)
			err = WrapError(WrapError(err))
			if err == nil || NotFoundError(err) != tc.gone || !strings.Contains(err.Error(), tc.code) || calls != 1 {
				t.Fatalf("unexpected parent classification: gone=%t calls=%d err=%v", tc.gone, calls, err)
			}
			// Existing instance-resource callers retain their original HTTP 404 behavior.
			_, legacyErr := s.DescribeDBInstance("rm-test")
			if !NotFoundError(legacyErr) || calls != 2 {
				t.Fatalf("legacy instance query changed: calls=%d err=%v", calls, legacyErr)
			}
		})
	}
}

func TestRdsChildDeleteRejectsParentHTTPNotFound(t *testing.T) {
	for _, kind := range []string{"alicloud_db_database", "alicloud_db_account_privilege", "alicloud_rds_account", "alicloud_db_account"} {
		for _, code := range []string{"Unknown", "InvalidAccessKeyId.NotFound"} {
			t.Run(kind+"/"+code, func(t *testing.T) {
				var parents, mutations int32
				client := rdsAccountTestClient(t, func(w http.ResponseWriter, r *http.Request) {
					if r.Form.Get("Action") == "DescribeDBInstanceAttribute" {
						atomic.AddInt32(&parents, 1)
						writeRdsError(t, w, 404, code)
					} else {
						atomic.AddInt32(&mutations, 1)
						writeRdsError(t, w, 400, "UnexpectedAction")
					}
				})
				id := "rm-test:account"
				if kind == "alicloud_db_account_privilege" {
					id += ":ReadOnly"
				}
				d := rdsTestData(t, kind, id)
				err := Provider().(*schema.Provider).ResourcesMap[kind].Delete(d, client)
				if err == nil || NotFoundError(err) || d.Id() != id || parents != 1 || mutations != 0 {
					t.Fatalf("parent lookup failed: err=%v id=%q parents=%d mutations=%d", err, d.Id(), parents, mutations)
				}
			})
		}
	}
}

func TestRdsMutationRetriesParentConfirmation(t *testing.T) {
	for _, operation := range []string{"account", "revoke"} {
		for _, recover := range []bool{true, false} {
			t.Run(fmt.Sprintf("%s/recover=%t", operation, recover), func(t *testing.T) {
				var mutations, confirmations int32
				client := rdsAccountTestClient(t, func(w http.ResponseWriter, r *http.Request) {
					w.Header().Set("Content-Type", "application/json")
					switch r.Form.Get("Action") {
					case "DeleteAccount", "RevokeAccountPrivilege":
						n := atomic.AddInt32(&mutations, 1)
						if !recover || n <= 2 {
							writeGone403(t, w)
						} else {
							_, _ = w.Write([]byte(`{}`))
						}
					case "DescribeDBInstanceAttribute":
						if atomic.LoadInt32(&mutations) == 0 {
							writeRdsParent(w, "Running")
							return
						}
						n := atomic.AddInt32(&confirmations, 1)
						if !recover || n == 1 {
							writeRdsError(t, w, 400, "InvalidParameter")
						} else {
							writeRdsParent(w, "Running")
						}
					case "DescribeAccounts":
						_, _ = w.Write([]byte(`{"Accounts":{"DBInstanceAccount":[]}}`))
					case "DescribeDatabases":
						_, _ = w.Write([]byte(`{"Databases":{"Database":[]}}`))
					default:
						t.Errorf("unexpected action %s", r.Form.Get("Action"))
						writeRdsError(t, w, 400, "UnexpectedAction")
					}
				})
				budget := 1100 * time.Millisecond
				if recover {
					budget = 8 * time.Second
				}
				start := time.Now()
				var err error
				if operation == "account" {
					res := resourceAliCloudRdsAccount()
					res.Timeouts.Delete = schema.DefaultTimeout(budget)
					d := res.Data(nil)
					d.SetId("rm-test:account")
					err = res.Delete(d, client)
				} else {
					s := RdsService{client}
					err = s.revokeAccountPrivilege("rm-test:account:ReadOnly", "testdb", budget)
				}
				elapsed := time.Since(start)
				if recover {
					if err != nil || mutations != 3 || confirmations != 2 || elapsed > budget {
						t.Fatalf("parent retry decision lost: err=%v mutations=%d confirmations=%d elapsed=%s", err, mutations, confirmations, elapsed)
					}
				} else if err == nil || !rdsErrorHasCode(err, "InvalidParameter") || NotFoundError(err) || mutations < 2 || confirmations != mutations || elapsed < budget || elapsed > budget+time.Second {
					t.Fatalf("persistent parent failure must exhaust existing retry budget: err=%v mutations=%d confirmations=%d elapsed=%s", err, mutations, confirmations, elapsed)
				}
			})
		}
	}
}

func TestRdsMutationRejectsParentHTTPNotFound(t *testing.T) {
	for _, operation := range []string{"database", "account", "revoke"} {
		t.Run(operation, func(t *testing.T) {
			var mutations, confirmations int32
			client := rdsAccountTestClient(t, func(w http.ResponseWriter, r *http.Request) {
				switch r.Form.Get("Action") {
				case "DescribeDBInstanceAttribute":
					if atomic.LoadInt32(&mutations) == 0 {
						writeRdsParent(w, "Running")
					} else {
						atomic.AddInt32(&confirmations, 1)
						writeRdsError(t, w, 404, "InvalidAccessKeyId.NotFound")
					}
				case "DeleteDatabase", "DeleteAccount", "RevokeAccountPrivilege":
					atomic.AddInt32(&mutations, 1)
					writeGone403(t, w)
				default:
					t.Errorf("unexpected action %s", r.Form.Get("Action"))
					writeRdsError(t, w, 400, "UnexpectedAction")
				}
			})
			var err error
			switch operation {
			case "database":
				err = resourceAliCloudRdsDatabaseDelete(rdsTestData(t, "alicloud_db_database", "rm-test:testdb"), client)
			case "account":
				err = resourceAliCloudRdsAccountDelete(rdsAccountResourceData(t), client)
			case "revoke":
				s := RdsService{client}
				err = s.RevokeAccountPrivilege("rm-test:account:ReadOnly", "testdb")
			}
			if err == nil || NotFoundError(err) || mutations != 1 || confirmations != 1 {
				t.Fatalf("failed parent query must fail deletion without retry: err=%v mutations=%d confirmations=%d", err, mutations, confirmations)
			}
		})
	}
}
