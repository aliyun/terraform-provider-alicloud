package dataworks

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

// listItem builds one ListResources item.
func listItem(id, name string) map[string]interface{} {
	return map[string]interface{}{
		"Id":         id,
		"Name":       name,
		"ProjectId":  "proj-1",
		"Owner":      "owner-1",
		"Type":       "file",
		"CreateTime": 1700000000,
		"ModifyTime": 1700000001,
		"DataSource": map[string]interface{}{"Name": "ds-1", "Type": "mysql"},
		"Script":     map[string]interface{}{"Path": "/s", "Id": "sc-1"},
	}
}

// listResponse wraps items in the PagingInfo.Resources envelope the ListResources
// API returns.
func listResponse(items []map[string]interface{}) map[string]interface{} {
	return map[string]interface{}{
		"PagingInfo": map[string]interface{}{
			"Resources": items,
		},
	}
}

// TestUnitDataWorksResourceDataSourceLists: a single page lists all resources and
// maps each to a flat entry with a composite id.
func TestUnitDataWorksResourceDataSourceLists(t *testing.T) {
	client := dataWorksTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		switch r.Form.Get("Action") {
		case "ListResources":
			assert.Equal(t, "proj-1", r.Form.Get("ProjectId"))
			writeJSON(t, w, listResponse([]map[string]interface{}{
				listItem("r1", "n1"),
				listItem("r2", "n2"),
			}))
		default:
			t.Errorf("unexpected action %q", r.Form.Get("Action"))
		}
	})

	entries, err := resourceDataSourceRead(client, "proj-1", "", "", nil)
	assert.NoError(t, err)
	assert.Len(t, entries, 2)
	assert.Equal(t, "proj-1:r1", entries[0].Id.ValueString())
	assert.Equal(t, "r1", entries[0].ResourceId.ValueString())
	assert.Equal(t, "n1", entries[0].ResourceName.ValueString())
	assert.Equal(t, "ds-1", entries[0].DataSourceName.ValueString())
	assert.Equal(t, "sc-1", entries[0].ScriptId.ValueString())
	assert.Equal(t, int64(1700000000), entries[0].CreateTime.ValueInt64())
}

// TestUnitDataWorksResourceDataSourcePaginates: a full first page (PageSize
// items) drives a second request, and both pages are concatenated.
func TestUnitDataWorksResourceDataSourcePaginates(t *testing.T) {
	full := make([]map[string]interface{}, pageSize)
	for i := range full {
		full[i] = listItem(fmt.Sprintf("r%d", i), fmt.Sprintf("n%d", i))
	}
	tail := []map[string]interface{}{listItem("rTail", "nTail")}

	client := dataWorksTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		switch r.Form.Get("Action") {
		case "ListResources":
			page := r.Form.Get("PageNumber")
			switch page {
			case "1":
				assert.Equal(t, fmt.Sprint(pageSize), r.Form.Get("PageSize"))
				writeJSON(t, w, listResponse(full))
			case "2":
				writeJSON(t, w, listResponse(tail))
			default:
				t.Errorf("unexpected PageNumber %q", page)
			}
		default:
			t.Errorf("unexpected action %q", r.Form.Get("Action"))
		}
	})

	entries, err := resourceDataSourceRead(client, "proj-1", "", "", nil)
	assert.NoError(t, err)
	assert.Len(t, entries, pageSize+1, "both pages must be concatenated")
	assert.Equal(t, "r0", entries[0].ResourceId.ValueString())
	assert.Equal(t, "rTail", entries[pageSize].ResourceId.ValueString(), "the second page's item must appear after the first page's")
}

// TestUnitDataWorksResourceDataSourceIdsFilter: an ids filter narrows the result
// to the requested resource id.
func TestUnitDataWorksResourceDataSourceIdsFilter(t *testing.T) {
	client := dataWorksTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		switch r.Form.Get("Action") {
		case "ListResources":
			writeJSON(t, w, listResponse([]map[string]interface{}{
				listItem("r1", "n1"),
				listItem("r2", "n2"),
				listItem("r3", "n3"),
			}))
		default:
			t.Errorf("unexpected action %q", r.Form.Get("Action"))
		}
	})

	entries, err := resourceDataSourceRead(client, "proj-1", "", "", []string{"r2"})
	assert.NoError(t, err)
	assert.Len(t, entries, 1, "only the requested id survives the filter")
	assert.Equal(t, "r2", entries[0].ResourceId.ValueString())
}
