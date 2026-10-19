package paiplugin_test

import (
	"net/http"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// TestUnitPaiPluginAlgorithmDataSourceList covers the collection path: ListAlgorithms
// returns two algorithms in a fixed order, and the data source surfaces them as
// parallel ids/names lists whose index order matches the gateway response.
func TestUnitPaiPluginAlgorithmDataSourceList(t *testing.T) {
	const dataId = "data.alicloud_pai_plugin_algorithm.default"

	paiPluginTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == http.MethodGet && r.URL.Path == "/api/v2/algorithms":
			writeJSON(t, w, http.StatusOK, paiPluginEnvelope(map[string]interface{}{
				"Algorithms": []interface{}{
					map[string]interface{}{"Id": "alg-1", "Name": "AlgorithmOne"},
					map[string]interface{}{"Id": "alg-2", "Name": "AlgorithmTwo"},
				},
				"TotalCount": 2,
				"PageNumber": 1,
				"PageSize":   10,
			}))
		default:
			paiPluginNotFound(w)
		}
	})

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() {},
		ProtoV5ProviderFactories: factories(),
		IsUnitTest:               true,
		Steps: []resource.TestStep{
			{
				Config: paiPluginProviderConfig() + `
data "alicloud_pai_plugin_algorithm" "default" {}
`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dataId, "ids.#", "2"),
					// Collection order: index 0 is the first gateway element.
					resource.TestCheckResourceAttr(dataId, "ids.0", "alg-1"),
					resource.TestCheckResourceAttr(dataId, "ids.1", "alg-2"),
					resource.TestCheckResourceAttr(dataId, "names.0", "AlgorithmOne"),
					resource.TestCheckResourceAttr(dataId, "names.1", "AlgorithmTwo"),
					resource.TestCheckResourceAttr(dataId, "algorithm_id_out", "alg-1"),
					resource.TestCheckResourceAttr(dataId, "algorithm_name_out", "AlgorithmOne"),
				),
			},
		},
	})
}

// TestUnitPaiPluginAlgorithmDataSourceById covers the single-algorithm detail
// path: with algorithm_id set, GetAlgorithm returns the full object including
// description and the user config maps.
func TestUnitPaiPluginAlgorithmDataSourceById(t *testing.T) {
	const dataId = "data.alicloud_pai_plugin_algorithm.default"

	paiPluginTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == http.MethodGet && r.URL.Path == "/api/v2/algorithms/alg-1":
			writeJSON(t, w, http.StatusOK, paiPluginEnvelope(map[string]interface{}{
				"Id":                 "alg-1",
				"Name":               "AlgorithmOne",
				"Description":        "First algorithm",
				"InferUserConfigMap": `{"key":"infer"}`,
				"TrainUserConfigMap": `{"key":"train"}`,
			}))
		default:
			paiPluginNotFound(w)
		}
	})

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() {},
		ProtoV5ProviderFactories: factories(),
		IsUnitTest:               true,
		Steps: []resource.TestStep{
			{
				Config: paiPluginProviderConfig() + `
data "alicloud_pai_plugin_algorithm" "default" {
  algorithm_id = "alg-1"
}
`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dataId, "algorithm_id", "alg-1"),
					resource.TestCheckResourceAttr(dataId, "algorithm_id_out", "alg-1"),
					resource.TestCheckResourceAttr(dataId, "algorithm_name_out", "AlgorithmOne"),
					resource.TestCheckResourceAttr(dataId, "description", "First algorithm"),
					resource.TestCheckResourceAttr(dataId, "ids.#", "1"),
					resource.TestCheckResourceAttr(dataId, "ids.0", "alg-1"),
				),
			},
		},
	})
}

func TestUnitPaiPluginUserDataSource(t *testing.T) {
	const dataId = "data.alicloud_pai_plugin_user.default"

	paiPluginTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == http.MethodGet && r.URL.Path == "/api/v2/users":
			writeJSON(t, w, http.StatusOK, paiPluginEnvelope(map[string]interface{}{
				"AccountStatus": 1,
			}))
		default:
			paiPluginNotFound(w)
		}
	})

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() {},
		ProtoV5ProviderFactories: factories(),
		IsUnitTest:               true,
		Steps: []resource.TestStep{
			{
				Config: paiPluginProviderConfig() + `
data "alicloud_pai_plugin_user" "default" {}
`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dataId, "id", "alicloud_pai_plugin_user"),
					resource.TestCheckResourceAttr(dataId, "account_status", "1"),
				),
			},
		},
	})
}

func TestUnitPaiPluginMessageConfigDataSource(t *testing.T) {
	const dataId = "data.alicloud_pai_plugin_message_config.default"

	paiPluginTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == http.MethodGet && r.URL.Path == "/api/v2/users/messageConfig":
			writeJSON(t, w, http.StatusOK, paiPluginEnvelope(map[string]interface{}{
				"SmsReportUrl": "https://example.com/report",
				"SmsUpUrl":     "https://example.com/up",
			}))
		default:
			paiPluginNotFound(w)
		}
	})

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() {},
		ProtoV5ProviderFactories: factories(),
		IsUnitTest:               true,
		Steps: []resource.TestStep{
			{
				Config: paiPluginProviderConfig() + `
data "alicloud_pai_plugin_message_config" "default" {}
`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dataId, "id", "alicloud_pai_plugin_message_config"),
					resource.TestCheckResourceAttr(dataId, "sms_report_url", "https://example.com/report"),
					resource.TestCheckResourceAttr(dataId, "sms_up_url", "https://example.com/up"),
				),
			},
		},
	})
}

// TestUnitPaiPluginMessageMetricsDataSource covers the collection path: the
// metrics data source returns parallel arrays (dates, fails, pendings, rates,
// successes, totals) whose index order matches the gateway Metrics array.
func TestUnitPaiPluginMessageMetricsDataSource(t *testing.T) {
	const dataId = "data.alicloud_pai_plugin_message_metrics.default"

	paiPluginTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == http.MethodGet && r.URL.Path == "/api/v2/messages/metrics":
			writeJSON(t, w, http.StatusOK, paiPluginEnvelope(map[string]interface{}{
				"Metrics": []interface{}{
					map[string]interface{}{
						"Date": "20260101", "Fail": 1, "Pending": 2,
						"Rate": 0.5, "Success": 1, "Total": 2,
					},
					map[string]interface{}{
						"Date": "20260102", "Fail": 0, "Pending": 0,
						"Rate": 1.0, "Success": 3, "Total": 3,
					},
				},
				"TotalCount": 2,
				"PageNumber": 1,
				"PageSize":   10,
			}))
		default:
			paiPluginNotFound(w)
		}
	})

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() {},
		ProtoV5ProviderFactories: factories(),
		IsUnitTest:               true,
		Steps: []resource.TestStep{
			{
				Config: paiPluginProviderConfig() + `
data "alicloud_pai_plugin_message_metrics" "default" {
  start_date = "20260101"
  end_date   = "20260102"
}
`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dataId, "dates.#", "2"),
					// Collection order: index 0 is the first metrics element.
					resource.TestCheckResourceAttr(dataId, "dates.0", "20260101"),
					resource.TestCheckResourceAttr(dataId, "dates.1", "20260102"),
					resource.TestCheckResourceAttr(dataId, "totals.0", "2"),
					resource.TestCheckResourceAttr(dataId, "totals.1", "3"),
					resource.TestCheckResourceAttr(dataId, "successes.0", "1"),
					resource.TestCheckResourceAttr(dataId, "successes.1", "3"),
					resource.TestCheckResourceAttr(dataId, "fails.0", "1"),
					resource.TestCheckResourceAttr(dataId, "fails.1", "0"),
					resource.TestCheckResourceAttr(dataId, "total_count", "2"),
				),
			},
		},
	})
}
