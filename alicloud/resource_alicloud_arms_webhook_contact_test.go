package alicloud

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"

	credentials "github.com/aliyun/credentials-go/credentials"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

// Test Arms WebhookContact. >>> Resource test cases, automatically generated.
func TestAccAliCloudArmsWebhookContact_basic(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_arms_webhook_contact.default"
	ra := resourceAttrInit(resourceId, AlicloudArmsWebhookContactMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ArmsServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeArmsWebhookContact")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sarmswebhookcontact%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudArmsWebhookContactBasicDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"webhook_contact_name": name,
					"method":               "Post",
					"url":                  "https://example.com/webhook",
					"body":                 "alert fired",
					"recover_body":         "alert recovered",
					"biz_headers":          map[string]interface{}{"Content-Type": "application/json"},
					"biz_params":           map[string]interface{}{"foo": "bar"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"webhook_contact_name":     name,
						"webhook_contact_id":       CHECKSET,
						"method":                   "Post",
						"url":                      "https://example.com/webhook",
						"body":                     "alert fired",
						"recover_body":             "alert recovered",
						"biz_headers.%":            "1",
						"biz_headers.Content-Type": "application/json",
						"biz_params.%":             "1",
						"biz_params.foo":           "bar",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"webhook_contact_name": name,
					"method":               "Post",
					"url":                  "https://example.com/webhook",
					"body":                 "alert fired v2",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"webhook_contact_name": name,
						"method":               "Post",
						"body":                 "alert fired v2",
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

var AlicloudArmsWebhookContactMap = map[string]string{}

func AlicloudArmsWebhookContactBasicDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}
`, name)
}

func armsWebhookContactTestClient(t *testing.T, handler http.HandlerFunc) *connectivity.AliyunClient {
	t.Helper()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			t.Error(err)
		}
		handler(w, r)
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
	endpoints.Store("arms", endpoint)
	return client
}

func TestUnitArmsWebhookContactCrud(t *testing.T) {
	t.Run("create_read_delete", func(t *testing.T) {
		var deleted string
		client := armsWebhookContactTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			var response interface{}
			switch r.Form.Get("Action") {
			case "CreateOrUpdateWebhookContact":
				if r.Form.Get("WebhookName") != "test-webhook" {
					t.Errorf("unexpected WebhookName: %q", r.Form.Get("WebhookName"))
				}
				if r.Form.Get("Method") != "Post" || r.Form.Get("Url") != "https://example.com/hook" {
					t.Errorf("unexpected webhook request params: Method=%q Url=%q", r.Form.Get("Method"), r.Form.Get("Url"))
				}
				if r.Form.Get("Body") != `{"message":"alert fired"}` || r.Form.Get("RecoverBody") != `{"message":"alert recovered"}` {
					t.Errorf("unexpected template params: Body=%q RecoverBody=%q", r.Form.Get("Body"), r.Form.Get("RecoverBody"))
				}
				if r.Form.Get("BizHeaders") != `[{"Content-Type":"application/json"}]` {
					t.Errorf("unexpected BizHeaders: %q", r.Form.Get("BizHeaders"))
				}
				if r.Form.Get("BizParams") != `[{"foo":"bar"}]` {
					t.Errorf("unexpected BizParams: %q", r.Form.Get("BizParams"))
				}
				response = map[string]interface{}{
					"WebhookContact": map[string]interface{}{"WebhookId": "wc-test-123"},
				}
			case "DescribeWebhookContacts":
				if r.Form.Get("Page") == "" || r.Form.Get("Size") == "" {
					w.WriteHeader(http.StatusBadRequest)
					response = map[string]interface{}{
						"Code":    "MissingPage",
						"Message": "Page is mandatory for this action.",
					}
					break
				}
				if r.Form.Get("ContactIds") != "wc-test-123" {
					t.Errorf("unexpected ContactIds: %q", r.Form.Get("ContactIds"))
				}
				if r.Form.Get("Page") != "1" {
					t.Errorf("unexpected Page: %q", r.Form.Get("Page"))
				}
				response = map[string]interface{}{
					"PageBean": map[string]interface{}{
						"WebhookContacts": []interface{}{
							map[string]interface{}{
								"WebhookId":   "wc-test-123",
								"WebhookName": "test-webhook",
								"Webhook": map[string]interface{}{
									"Method":      "Post",
									"Url":         "https://example.com/hook",
									"Body":        `{"message":"alert fired"}`,
									"RecoverBody": `{"message":"alert recovered"}`,
									"BizHeaders":  []interface{}{map[string]interface{}{"Content-Type": "application/json"}},
									"BizParams":   `[{"foo":"bar"}]`,
								},
							},
						},
					},
				}
			case "DeleteWebhookContact":
				deleted = r.Form.Get("WebhookId")
				response = map[string]interface{}{"Code": 200}
			default:
				t.Errorf("unexpected action: %q", r.Form.Get("Action"))
			}
			if err := json.NewEncoder(w).Encode(response); err != nil {
				t.Error(err)
			}
		})

		d := schema.TestResourceDataRaw(t, resourceAliCloudArmsWebhookContact().Schema, map[string]interface{}{
			"webhook_contact_name": "test-webhook",
			"method":               "Post",
			"url":                  "https://example.com/hook",
			"body":                 `{"message":"alert fired"}`,
			"recover_body":         `{"message":"alert recovered"}`,
			"biz_headers":          map[string]interface{}{"Content-Type": "application/json"},
			"biz_params":           map[string]interface{}{"foo": "bar"},
		})
		if err := resourceAliCloudArmsWebhookContactCreate(d, client); err != nil {
			t.Fatalf("create failed: %v", err)
		}
		if d.Id() != "wc-test-123" {
			t.Fatalf("expected id wc-test-123, got %q", d.Id())
		}
		if got := d.Get("webhook_contact_name"); got != "test-webhook" {
			t.Fatalf("read back name %q", got)
		}
		if got := d.Get("method"); got != "Post" {
			t.Fatalf("read back method %q", got)
		}
		if got := d.Get("url"); got != "https://example.com/hook" {
			t.Fatalf("read back url %q", got)
		}
		if got := d.Get("body"); got != `{"message":"alert fired"}` {
			t.Fatalf("read back body %q", got)
		}
		if got := d.Get("recover_body"); got != `{"message":"alert recovered"}` {
			t.Fatalf("read back recover_body %q", got)
		}
		headers, ok := d.Get("biz_headers").(map[string]interface{})
		if !ok || headers["Content-Type"] != "application/json" {
			t.Fatalf("unexpected biz_headers read back: %v", d.Get("biz_headers"))
		}
		params, ok := d.Get("biz_params").(map[string]interface{})
		if !ok || params["foo"] != "bar" {
			t.Fatalf("unexpected biz_params read back: %v", d.Get("biz_params"))
		}
		if got := d.Get("webhook_contact_id"); got != "wc-test-123" {
			t.Fatalf("read back webhook_contact_id %q", got)
		}
		if err := resourceAliCloudArmsWebhookContactDelete(d, client); err != nil {
			t.Fatalf("delete failed: %v", err)
		}
		if deleted != "wc-test-123" {
			t.Fatalf("expected DeleteWebhookContact WebhookId=wc-test-123, got %q", deleted)
		}
	})

	t.Run("read_not_found_clears_id", func(t *testing.T) {
		client := armsWebhookContactTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			response := map[string]interface{}{
				"PageBean": map[string]interface{}{"WebhookContacts": []interface{}{}},
			}
			if err := json.NewEncoder(w).Encode(response); err != nil {
				t.Error(err)
			}
		})
		d := schema.TestResourceDataRaw(t, resourceAliCloudArmsWebhookContact().Schema, map[string]interface{}{
			"webhook_contact_name": "test-webhook",
		})
		d.SetId("wc-test-gone")
		if err := resourceAliCloudArmsWebhookContactRead(d, client); err != nil {
			t.Fatalf("read must clear the id on not found: %v", err)
		}
		if d.Id() != "" {
			t.Fatalf("expected id to be cleared, got %q", d.Id())
		}
	})
}
