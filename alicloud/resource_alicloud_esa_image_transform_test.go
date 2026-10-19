package alicloud

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"sync/atomic"
	"testing"

	credentials "github.com/aliyun/credentials-go/credentials"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
)

// Test ESA ImageTransform. >>> Resource test cases, automatically generated.
// Case resource_ImageTransform_test
func TestAccAliCloudESAImageTransformresource_ImageTransform_test(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_esa_image_transform.default"
	ra := resourceAttrInit(resourceId, AliCloudESAImageTransformresource_ImageTransform_testMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EsaServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeEsaImageTransform")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sESAImageTransform%d", defaultRegionToTest, rand)

	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudESAImageTransformresource_ImageTransform_testBasicDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
		},
		IDRefreshName:     resourceId,
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"site_id":      "${alicloud_esa_site.default.id}",
					"rule_enable":  "off",
					"enable":       "off",
					"rule":         "http.host eq \\\"video.example.com\\\"",
					"site_version": "0",
					"rule_name":    "rule_example",
					"sequence":     "1",
					"auto_avif":    "off",
					"auto_webp":    "off",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"rule_enable": "on",
					"enable":      "on",
					"rule":        "http.host eq \\\"videoo.example.com\\\"",
					"rule_name":   "rule_viedoo",
					"sequence":    "2",
					"auto_avif":   "on",
					"auto_webp":   "on",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"rule_enable": "on",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"rule": "http.host eq \\\"image.example.com\\\"",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"rule_name": "rule_image",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"enable":    "on",
					"auto_avif": "on",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"rule_enable": "off",
					"enable":      "off",
					"rule":        "http.host eq \\\"video.example.com\\\"",
					"rule_name":   "rule_example",
					"auto_webp":   "off",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{},
			},
		},
	})
}

var AliCloudESAImageTransformresource_ImageTransform_testMap = map[string]string{
	"id": CHECKSET,
}

func AliCloudESAImageTransformresource_ImageTransform_testBasicDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

data "alicloud_esa_sites" "default" {
  plan_subscribe_type = "enterpriseplan"
}

resource "alicloud_esa_site" "default" {
  site_name   = "${lower(var.name)}.com"
  instance_id = data.alicloud_esa_sites.default.sites.0.instance_id
  coverage    = "domestic"
  access_type = "NS"
}

`, name)
}

// Test ESA ImageTransform. <<< Resource test cases, automatically generated.

// setEsaImageTransformTestAttrs sets the shared unit-test attribute set
// (including the new auto_avif / auto_webp enum fields) onto d using literal
// schema keys so tflint R001 (ResourceData.Set key must be a string literal)
// stays satisfied. The attributes map is retained for newInstanceDiff.
func setEsaImageTransformTestAttrs(t *testing.T, d *schema.ResourceData) {
	assert.Nil(t, d.Set("site_id", "12345"))
	assert.Nil(t, d.Set("site_version", 0))
	assert.Nil(t, d.Set("sequence", 1))
	assert.Nil(t, d.Set("rule_enable", "off"))
	assert.Nil(t, d.Set("rule_name", "rule_example"))
	assert.Nil(t, d.Set("rule", `http.host eq "video.example.com"`))
	assert.Nil(t, d.Set("enable", "off"))
	assert.Nil(t, d.Set("auto_avif", "off"))
	assert.Nil(t, d.Set("auto_webp", "off"))
}

// esaMockConfig configures the canned behavior of newEsaMockServer per subcase.
// It is the offline counterpart of the previous gomonkey patches: instead of
// monkey-patching RpcPost / DescribeEsaImageTransform (which silently fails
// under the ACube FC container's seccomp profile because ptrace/mprotect is
// denied), the real RPC client is wired to an httptest server, so the CRUD
// path exercises the actual request building and response parsing code.
type esaMockConfig struct {
	readResponse    map[string]interface{} // GetImageTransform success body
	createResponse  map[string]interface{} // CreateImageTransform success body (default {ConfigId:67890})
	updateResponse  map[string]interface{} // UpdateImageTransform success body (default {})
	deleteResponse  map[string]interface{} // DeleteImageTransform success body (default {})
	getNotFound     bool                   // GetImageTransform returns 404 NotFound
	deleteNotFound  bool                   // DeleteImageTransform returns 404 NotFound
	createError     string                 // if non-empty, CreateImageTransform returns this error code
	createFailCount int                    // number of Create attempts that fail with createError before success
}

// newEsaMockServer builds an httptest server dispatching ESA RPC actions and a
// client wired to it. The server captures request form params per action in
// the returned *sync.Map (action -> map[string]string). RPC params may arrive
// as URL query (RpcGet) or POST body form (RpcPost); ParseForm merges both, and
// a JSON-body fallback covers the case where the tea SDK encodes the body as
// JSON rather than form-encoded.
func newEsaMockServer(t *testing.T, cfg esaMockConfig) (*httptest.Server, *sync.Map, *connectivity.AliyunClient) {
	captured := &sync.Map{}
	var createAttempts int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		bodyBytes, _ := io.ReadAll(r.Body)
		_ = r.Body.Close()
		r.Body = io.NopCloser(bytes.NewReader(bodyBytes))
		if err := r.ParseForm(); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		action := r.Form.Get("Action")
		params := map[string]string{}
		for k, vs := range r.Form {
			if len(vs) > 0 {
				params[k] = vs[0]
			}
		}
		// Fallback: tea SDK may encode the POST body as JSON rather than form.
		if action == "" {
			var j map[string]interface{}
			if json.Unmarshal(bodyBytes, &j) == nil {
				if a, ok := j["Action"].(string); ok {
					action = a
				}
				for k, v := range j {
					params[k] = fmt.Sprint(v)
				}
			}
		}
		if action != "" {
			captured.Store(action, params)
		}

		w.Header().Set("Content-Type", "application/json")
		encode := func(v map[string]interface{}) {
			if err := json.NewEncoder(w).Encode(v); err != nil {
				t.Errorf("encode response: %v", err)
			}
		}
		// errBody returns an RPC soft-error: HTTP 200 with a body carrying a
		// non-success Code. formatError (connectivity/client.go) inspects the
		// response body regardless of HTTP status and rebuilds a tea.SDKError
		// from Code/Message, which is what NeedRetry / IsExpectedErrors /
		// NotFoundError match on. Using 200 + body avoids depending on how the
		// tea SDK maps non-2xx HTTP status to errors.
		errBody := func(code string) {
			encode(map[string]interface{}{
				"Code":      code,
				"Message":   code,
				"RequestId": "test-req-id",
			})
		}
		switch action {
		case "CreateImageTransform":
			if cfg.createError != "" {
				n := int(atomic.AddInt32(&createAttempts, 1))
				if n <= cfg.createFailCount {
					errBody(cfg.createError)
					return
				}
			}
			resp := cfg.createResponse
			if resp == nil {
				resp = map[string]interface{}{"ConfigId": 67890}
			}
			encode(resp)
		case "UpdateImageTransform":
			resp := cfg.updateResponse
			if resp == nil {
				resp = map[string]interface{}{}
			}
			encode(resp)
		case "GetImageTransform":
			if cfg.getNotFound {
				errBody("NotFound")
				return
			}
			encode(cfg.readResponse)
		case "DeleteImageTransform":
			if cfg.deleteNotFound {
				errBody("NotFound")
				return
			}
			resp := cfg.deleteResponse
			if resp == nil {
				resp = map[string]interface{}{}
			}
			encode(resp)
		default:
			errBody("UnknownAction")
		}
	}))
	t.Cleanup(server.Close)

	endpoint := strings.TrimPrefix(server.URL, "http://")
	t.Setenv("NO_PROXY", endpoint)
	credential, err := credentials.NewCredential(new(credentials.Config).
		SetType("access_key").SetAccessKeyId("test-key").SetAccessKeySecret("test-secret"))
	if err != nil {
		t.Fatal(err)
	}
	endpoints := new(sync.Map)
	endpoints.Store("esa", endpoint)
	config := &connectivity.Config{
		AccessKey: "test-key", SecretKey: "test-secret", Credential: credential,
		RegionId: "cn-hangzhou", AccountType: "test", Protocol: "http",
		Endpoints: endpoints, SignVersion: new(sync.Map), SkipRegionValidation: true,
	}
	// config.Client() returns (*AliyunClient, error) directly; no interface
	// assertion is needed (asserting a concrete type triggers a vet error).
	client, err := config.Client()
	if err != nil {
		t.Fatal(err)
	}
	return server, captured, client
}

// TestUnitAliCloudEsaImageTransform covers the ImageTransform resource CRUD
// logic without touching the cloud: an httptest server mocks the ESA RPC
// endpoints (Create/Update/Get/Delete ImageTransform) and the real RpcPost /
// RpcGet client talks to it. It exercises the auto_avif / auto_webp fields
// added for the property_completeness remediation, including retry and failure
// paths, and asserts the new params are sent on Create / Update and omitted
// when unchanged.
func TestUnitAliCloudEsaImageTransform(t *testing.T) {
	p := Provider().ResourcesMap
	resourceSchema := p["alicloud_esa_image_transform"].Schema

	// ReadMockResponse is the GetImageTransform response shape; includes the
	// new AutoAvif / AutoWebp response keys so the Read path is exercised.
	readMock := map[string]interface{}{
		"ConfigId": 67890, "SiteId": 12345, "SiteVersion": 0, "Sequence": 1,
		"RuleEnable": "off", "RuleName": "rule_example",
		"Rule":   `http.host eq "video.example.com"`,
		"Enable": "off", "AutoAvif": "off", "AutoWebp": "off",
	}
	readMockUpdated := map[string]interface{}{
		"ConfigId": 67890, "SiteId": 12345, "SiteVersion": 0, "Sequence": 1,
		"RuleEnable": "off", "RuleName": "rule_example",
		"Rule":   `http.host eq "video.example.com"`,
		"Enable": "on", "AutoAvif": "on", "AutoWebp": "on",
	}
	attributes := map[string]interface{}{
		"site_id":      "12345",
		"site_version": 0,
		"sequence":     1,
		"rule_enable":  "off",
		"rule_name":    "rule_example",
		"rule":         `http.host eq "video.example.com"`,
		"enable":       "off",
		"auto_avif":    "off",
		"auto_webp":    "off",
	}

	// newRD builds a fresh new-resource ResourceData populated with the shared
	// attributes (including auto_avif / auto_webp) for Create-path subcases.
	newRD := func() *schema.ResourceData {
		d, _ := schema.InternalMap(resourceSchema).Data(nil, nil)
		d.MarkNewResource()
		setEsaImageTransformTestAttrs(t, d)
		return d
	}

	// Create: success path + verify request contains AutoAvif/AutoWebp.
	t.Run("Create", func(t *testing.T) {
		_, captured, client := newEsaMockServer(t, esaMockConfig{readResponse: readMock})
		dInit := newRD()

		err := resourceAliCloudEsaImageTransformCreate(dInit, client)
		assert.Nil(t, err)
		// The new params must be present on the create request body.
		createParams, ok := captured.Load("CreateImageTransform")
		assert.True(t, ok, "CreateImageTransform request should have been captured")
		assert.Equal(t, "off", createParams.(map[string]string)["AutoAvif"], "AutoAvif should be sent on create")
		assert.Equal(t, "off", createParams.(map[string]string)["AutoWebp"], "AutoWebp should be sent on create")
		// Read populated the computed/returned values including the new fields.
		assert.Equal(t, "off", dInit.State().Attributes["auto_avif"])
		assert.Equal(t, "off", dInit.State().Attributes["auto_webp"])
		assert.Equal(t, "67890", dInit.State().Attributes["config_id"])
		assert.Equal(t, "12345:67890", dInit.State().ID)
	})

	// Update: change auto_avif / auto_webp / enable and verify the request carries them.
	t.Run("Update", func(t *testing.T) {
		_, captured, client := newEsaMockServer(t, esaMockConfig{readResponse: readMockUpdated})
		dExisted := newRD()
		dExisted.SetId("12345:67890")

		attributesDiff := map[string]interface{}{
			"enable":    "on",
			"auto_avif": "on",
			"auto_webp": "on",
		}
		diff, err := newInstanceDiff("alicloud_esa_image_transform", attributes, attributesDiff, dExisted.State())
		assert.Nil(t, err)
		dExisted, _ = schema.InternalMap(resourceSchema).Data(dExisted.State(), diff)

		err = resourceAliCloudEsaImageTransformUpdate(dExisted, client)
		assert.Nil(t, err)
		updateParams, ok := captured.Load("UpdateImageTransform")
		assert.True(t, ok, "UpdateImageTransform request should have been captured")
		assert.Equal(t, "on", updateParams.(map[string]string)["AutoAvif"], "AutoAvif should be sent on update when changed")
		assert.Equal(t, "on", updateParams.(map[string]string)["AutoWebp"], "AutoWebp should be sent on update when changed")
		assert.Equal(t, "on", updateParams.(map[string]string)["Enable"], "Enable should be sent on update when changed")
		assert.Equal(t, "on", dExisted.State().Attributes["auto_avif"])
		assert.Equal(t, "on", dExisted.State().Attributes["auto_webp"])
	})

	// Update with no change to auto_avif / auto_webp: those params must NOT be sent.
	t.Run("UpdateNoChangeOmitsAutoParams", func(t *testing.T) {
		_, captured, client := newEsaMockServer(t, esaMockConfig{readResponse: readMock})
		dInit := newRD()
		dInit.SetId("12345:67890")
		attributesDiff := map[string]interface{}{
			"sequence": 2,
		}
		diff, err := newInstanceDiff("alicloud_esa_image_transform", attributes, attributesDiff, dInit.State())
		assert.Nil(t, err)
		dNoChange, _ := schema.InternalMap(resourceSchema).Data(dInit.State(), diff)

		err = resourceAliCloudEsaImageTransformUpdate(dNoChange, client)
		assert.Nil(t, err)
		updateParams, ok := captured.Load("UpdateImageTransform")
		assert.True(t, ok, "UpdateImageTransform request should have been captured")
		_, hasAvif := updateParams.(map[string]string)["AutoAvif"]
		assert.False(t, hasAvif, "AutoAvif should not be sent when unchanged")
		_, hasWebp := updateParams.(map[string]string)["AutoWebp"]
		assert.False(t, hasWebp, "AutoWebp should not be sent when unchanged")
	})

	// Delete: normal path returns nil error.
	t.Run("Delete", func(t *testing.T) {
		_, _, client := newEsaMockServer(t, esaMockConfig{})
		dDelete := newRD()
		dDelete.SetId("12345:67890")

		err := resourceAliCloudEsaImageTransformDelete(dDelete, client)
		assert.Nil(t, err)
	})

	// Delete: NotFound is treated as success (idempotent delete).
	t.Run("DeleteNotFound", func(t *testing.T) {
		_, _, client := newEsaMockServer(t, esaMockConfig{deleteNotFound: true})
		dDelete := newRD()
		dDelete.SetId("12345:67890")

		err := resourceAliCloudEsaImageTransformDelete(dDelete, client)
		assert.Nil(t, err)
	})

	// Create failure path: non-retryable error surfaces; retryable errors eventually succeed.
	t.Run("CreateFailureAndRetry", func(t *testing.T) {
		// NonRetryableError -> CreateImageTransform surfaces an error.
		_, _, client1 := newEsaMockServer(t, esaMockConfig{
			readResponse:    readMock,
			createError:     "NonRetryableError",
			createFailCount: 100, // always fail (attempts never exceed 100)
		})
		dFail := newRD()
		err := resourceAliCloudEsaImageTransformCreate(dFail, client1)
		assert.NotNil(t, err)

		// Retryable (Throttling) twice then success.
		_, captured, client2 := newEsaMockServer(t, esaMockConfig{
			readResponse:    readMock,
			createError:     "Throttling",
			createFailCount: 2, // first 2 attempts fail with Throttling, 3rd succeeds
		})
		dFail2 := newRD()
		err = resourceAliCloudEsaImageTransformCreate(dFail2, client2)
		assert.Nil(t, err)
		createParams, ok := captured.Load("CreateImageTransform")
		assert.True(t, ok, "successful create request should be captured")
		assert.Equal(t, "off", createParams.(map[string]string)["AutoAvif"], "AutoAvif should be sent on the successful retry create")
	})

	// Read: NotFound on an existing (non-new) resource drops the resource cleanly.
	t.Run("ReadNotFoundDropsResource", func(t *testing.T) {
		_, _, client := newEsaMockServer(t, esaMockConfig{getNotFound: true})
		dRead, _ := schema.InternalMap(resourceSchema).Data(nil, nil)
		dRead.SetId("12345:67890")
		_ = dRead.Set("site_id", "12345")

		err := resourceAliCloudEsaImageTransformRead(dRead, client)
		assert.Nil(t, err)
		assert.Equal(t, "", dRead.Id(), "resource should be removed from state on NotFound")
	})
}
