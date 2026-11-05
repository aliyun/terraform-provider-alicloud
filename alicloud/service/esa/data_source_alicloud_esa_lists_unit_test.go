package esa

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"testing"

	credentials "github.com/aliyun/credentials-go/credentials"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func TestUnitEsaListsPagination(t *testing.T) {
	records := make([]map[string]interface{}, 100)
	for i := range records {
		records[i] = map[string]interface{}{
			"Id":          1000 + i,
			"Name":        fmt.Sprintf("fixture-list-%03d", i),
			"Kind":        "ip",
			"Description": fmt.Sprintf("fixture-description-%03d", i),
			"Length":      2,
			"UpdateTime":  "2025-01-01T00:00:00Z",
		}
	}

	for _, tc := range []struct {
		name          string
		total         int
		filters       map[string]interface{}
		wantStart     int
		wantCount     int
		wantPages     int
		wantQueryArgs map[string]interface{}
	}{
		{
			name: "item_like_passthrough", total: 51, wantCount: 51, wantPages: 2,
			filters: map[string]interface{}{
				"query_args": []interface{}{map[string]interface{}{"item_like": "192.0.2.", "desc": false}},
			},
			wantQueryArgs: map[string]interface{}{"ItemLike": "192.0.2.", "NameItemLike": "", "Desc": false},
		},
		{name: "50_plus_1", total: 51, wantCount: 51, wantPages: 2},
		{name: "ids_second_page", total: 51, filters: map[string]interface{}{"ids": []interface{}{"1050"}}, wantStart: 50, wantCount: 1, wantPages: 2},
		{name: "name_regex_second_page", total: 51, filters: map[string]interface{}{"name_regex": "^fixture-list-050$"}, wantStart: 50, wantCount: 1, wantPages: 2},
		{name: "ids_no_match", total: 51, filters: map[string]interface{}{"ids": []interface{}{"9999"}}, wantPages: 2},
		{name: "name_regex_no_match", total: 51, filters: map[string]interface{}{"name_regex": "^missing$"}, wantPages: 2},
		{name: "50_plus_50_plus_0", total: 100, wantCount: 100, wantPages: 3},
	} {
		t.Run(tc.name, func(t *testing.T) {
			var requests int32
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				requestNumber := int(atomic.AddInt32(&requests, 1))
				valid := true
				if err := r.ParseForm(); err != nil {
					t.Errorf("parse request: %v", err)
					valid = false
				}
				for key, want := range map[string]string{"Action": "ListLists", "PageSize": "50"} {
					if got := r.Form.Get(key); got != want {
						t.Errorf("request %d: %s = %q, want %q", requestNumber, key, got, want)
						valid = false
					}
				}
				page, err := strconv.Atoi(r.Form.Get("PageNumber"))
				if err != nil || page != requestNumber || page < 1 || page > tc.wantPages {
					t.Errorf("request %d: PageNumber = %q, want %d within 1..%d", requestNumber, r.Form.Get("PageNumber"), requestNumber, tc.wantPages)
					valid = false
				}
				if tc.wantQueryArgs != nil {
					var queryArgs map[string]interface{}
					if err := json.Unmarshal([]byte(r.Form.Get("QueryArgs")), &queryArgs); err != nil {
						t.Errorf("request %d: parse QueryArgs: %v", requestNumber, err)
						valid = false
					} else {
						for key, want := range tc.wantQueryArgs {
							if got := queryArgs[key]; got != want {
								t.Errorf("request %d: QueryArgs.%s = %#v, want %#v", requestNumber, key, got, want)
								valid = false
							}
						}
					}
				}
				// An invalid request gets an empty page so a pagination regression cannot loop forever.
				items := make([]map[string]interface{}, 0)
				if valid {
					start, end := (page-1)*50, page*50
					if end > tc.total {
						end = tc.total
					}
					items = records[start:end]
				}
				w.Header().Set("Content-Type", "application/json")
				if err := json.NewEncoder(w).Encode(map[string]interface{}{"Lists": items}); err != nil {
					t.Errorf("encode response: %v", err)
				}
			}))
			t.Cleanup(server.Close)

			credential, err := credentials.NewCredential(new(credentials.Config).
				SetType("access_key").SetAccessKeyId("test-key").SetAccessKeySecret("test-secret"))
			if err != nil {
				t.Fatal(err)
			}
			endpoints := new(sync.Map)
			endpoint := strings.TrimPrefix(server.URL, "http://")
			t.Setenv("NO_PROXY", endpoint)
			config := &connectivity.Config{
				AccessKey: "test-key", SecretKey: "test-secret", Credential: credential,
				RegionId: "cn-hangzhou", AccountType: "test", Protocol: "http",
				Endpoints: endpoints, SignVersion: new(sync.Map), SkipRegionValidation: true,
			}
			client, err := config.Client()
			if err != nil {
				t.Fatal(err)
			}
			endpoints.Store("esa", endpoint)

			ds := dataSourceAliCloudEsaLists()
			d := schema.TestResourceDataRaw(t, ds.Schema, tc.filters)
			if err := ds.Read(d, client); err != nil {
				t.Fatal(err)
			}
			if got := atomic.LoadInt32(&requests); int(got) != tc.wantPages {
				t.Fatalf("requests = %d, want %d", got, tc.wantPages)
			}

			wantIDs := make([]interface{}, tc.wantCount)
			wantNames := make([]interface{}, tc.wantCount)
			wantLists := make([]interface{}, tc.wantCount)
			for i := 0; i < tc.wantCount; i++ {
				record := records[tc.wantStart+i]
				wantIDs[i] = fmt.Sprint(record["Id"])
				wantNames[i] = record["Name"]
				wantLists[i] = map[string]interface{}{
					"id":          wantIDs[i],
					"list_id":     wantIDs[i],
					"name":        record["Name"],
					"kind":        record["Kind"],
					"description": record["Description"],
					"length":      record["Length"],
					"update_time": record["UpdateTime"],
				}
			}
			for field, want := range map[string][]interface{}{"ids": wantIDs, "names": wantNames, "lists": wantLists} {
				got := d.Get(field).([]interface{})
				if len(got) != len(want) {
					t.Fatalf("%s count = %d, want %d", field, len(got), len(want))
				}
				for i := range want {
					if !reflect.DeepEqual(got[i], want[i]) {
						t.Errorf("%s[%d] = %#v, want %#v", field, i, got[i], want[i])
					}
				}
			}
		})
	}
}
