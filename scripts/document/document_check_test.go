package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func resourceDocument(resourceName string) []string {
	return []string{
		"---",
		`layout: "alicloud"`,
		`page_title: "Alicloud: ` + resourceName + `"`,
		"description: |-",
		"  Provides a Alicloud resource.",
		"---",
		"",
		"# " + resourceName,
		"",
		"-> **NOTE:** Available since v1.0.0.",
		"",
		"## Example Usage",
		"",
		"```terraform",
		`resource "alicloud_demo" "example" {`,
		`  name = "example"`,
		"}",
		"```",
		"",
		"## Argument Reference",
		"",
		"The following arguments are supported:",
		"",
		"* `name` - (Required) The name.",
		"",
		"## Attributes Reference",
		"",
		"* `id` - The ID of the resource.",
		"",
		"## Timeouts",
		"",
		"## Import",
		"",
		"```shell",
		"terraform import alicloud_demo.example <id>",
		"```",
	}
}

func validEphemeralDocument() []string {
	document := validResourceDocument()
	for i, line := range document {
		if line == "# alicloud_demo" {
			document[i] = "# Ephemeral: alicloud_demo"
			break
		}
	}
	return document
}

func validFunctionDocument() []string {
	return []string{
		"---",
		`subcategory: ""`,
		`layout: "alicloud"`,
		`page_title: "Alibaba Cloud: arn_build"`,
		"description: |-",
		"  Builds an ARN from its constituent parts.",
		"---",
		"",
		"# Function: arn_build",
		"",
		"-> **NOTE:** Available since v2.0.0-beta3.",
		"",
		"Builds an ARN from its constituent parts.",
		"",
		"## Example Usage",
		"",
		"```terraform",
		`provider "alicloud" {}`,
		"",
		`output "arn" {`,
		`  value = provider::alicloud::arn_build("ecs", "cn-hangzhou", "1234567890", "instance/i-abc")`,
		"}",
		"```",
		"",
		"## Signature",
		"",
		"```text",
		"arn_build(service string, region string, account_id string, resource string) string",
		"```",
		"",
		"## Arguments",
		"",
		"The following arguments are supported:",
		"",
		"* `service` (String) The cloud service that owns the resource.",
		"* `region` (String) The region ID of the resource.",
		"",
		"## Return Type",
		"",
		"This function returns a string.",
	}
}

func validGuideDocument() []string {
	return []string{
		"---",
		`layout: "alicloud"`,
		`page_title: "Alibaba Cloud Account Guide"`,
		`sidebar_current: "docs-alicloud-guide-alicloud-account-guide"`,
		"description: |-",
		"  Sign up Alibaba Cloud and distinguish its type.",
		"---",
		"",
		"# Getting an Alibaba Cloud Account",
		"",
		"The Alibaba Cloud has two accounts: International-Site Account and China-Site Account.",
		"",
		"## Sign Up for an Alibaba Cloud International-Site Account",
		"",
		"To sign up for an International-Site account, visit [Alibaba Cloud International-Site Website](https://www.alibabacloud.com/).",
	}
}

func validResourceDocument() []string {
	return resourceDocument("alicloud_demo")
}

func writeDocument(t *testing.T, relativePath string, lines []string) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), filepath.FromSlash(relativePath))
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatalf("create document directory: %v", err)
	}
	if err := os.WriteFile(path, []byte(strings.Join(lines, "\n")+"\n"), 0o644); err != nil {
		t.Fatalf("write document: %v", err)
	}
	return path
}

func TestUnitDocumentCheckAcceptsValidResourcePage(t *testing.T) {
	if code := checkDocument(writeDocument(t, "website/docs/r/demo.html.markdown", validResourceDocument())); code != 0 {
		t.Fatalf("valid resource page rejected with exit code %d", code)
	}
}

func TestUnitDocumentCheckAcceptsValidDataSourcePage(t *testing.T) {
	if code := checkDocument(writeDocument(t, "website/docs/d/demo.html.markdown", validResourceDocument())); code != 0 {
		t.Fatalf("valid data source page rejected with exit code %d", code)
	}
}

func TestUnitDocumentCheckAcceptsValidEphemeralPage(t *testing.T) {
	if code := checkDocument(writeDocument(t, "website/docs/ephemeral-resources/demo.html.markdown", validEphemeralDocument())); code != 0 {
		t.Fatalf("valid ephemeral page rejected with exit code %d", code)
	}
}

func TestUnitDocumentCheckRejectsEphemeralPageWithPlainResourceTitle(t *testing.T) {
	if code := checkDocument(writeDocument(t, "website/docs/ephemeral-resources/demo.html.markdown", validResourceDocument())); code != 1 {
		t.Fatalf("ephemeral page with a plain resource title returned exit code %d, want 1", code)
	}
}

func TestUnitDocumentCheckAcceptsValidFunctionPage(t *testing.T) {
	if code := checkDocument(writeDocument(t, "website/docs/functions/arn_build.html.markdown", validFunctionDocument())); code != 0 {
		t.Fatalf("valid function page rejected with exit code %d", code)
	}
}

func TestUnitDocumentCheckRejectsFunctionPageWithResourcePageTitle(t *testing.T) {
	document := validFunctionDocument()
	for i, line := range document {
		if line == `page_title: "Alibaba Cloud: arn_build"` {
			document[i] = `page_title: "Alicloud: alicloud_arn_build"`
			break
		}
	}
	if code := checkDocument(writeDocument(t, "website/docs/functions/arn_build.html.markdown", document)); code != 1 {
		t.Fatalf("function page with a resource page_title returned exit code %d, want 1", code)
	}
}

func TestUnitDocumentCheckRejectsNonTerraformExampleFenceOnFunctionPage(t *testing.T) {
	document := validFunctionDocument()
	for i, line := range document {
		if line == "```terraform" {
			document[i] = "```hcl"
			break
		}
	}
	if code := checkDocument(writeDocument(t, "website/docs/functions/arn_build.html.markdown", document)); code != 1 {
		t.Fatalf("function page with a non-terraform example fence returned exit code %d, want 1", code)
	}
}

func TestUnitDocumentCheckRejectsUnclosedExampleFenceBeforeFunctionSections(t *testing.T) {
	// The example section must still report an unclosed fence when the
	// function page sections terminate it.
	document := validFunctionDocument()
	for i, line := range document {
		if line == "```" {
			document = append(document[:i:i], document[i+1:]...)
			break
		}
	}
	if code := checkDocument(writeDocument(t, "website/docs/functions/arn_build.html.markdown", document)); code != 1 {
		t.Fatalf("function page with an unclosed example fence returned exit code %d, want 1", code)
	}
}

func TestUnitDocumentCheckAcceptsBareFileNameAsGuide(t *testing.T) {
	// A bare file name has no owning directory and falls back to the guide
	// classification instead of a broken alicloud_ resource-name template.
	t.Chdir(t.TempDir())
	if err := os.WriteFile("getting-account.html.md", []byte(strings.Join(validGuideDocument(), "\n")+"\n"), 0o644); err != nil {
		t.Fatalf("write document: %v", err)
	}
	if code := checkDocument("getting-account.html.md"); code != 0 {
		t.Fatalf("bare guide page name rejected with exit code %d", code)
	}
}

func TestUnitDocumentCheckAcceptsGuidePageWithBusinessTitle(t *testing.T) {
	for _, name := range []string{"getting-account.html.md", "getting-account.html.markdown"} {
		if code := checkDocument(writeDocument(t, "website/docs/guides/"+name, validGuideDocument())); code != 0 {
			t.Fatalf("valid guide page %s rejected with exit code %d", name, code)
		}
	}
}

func TestUnitDocumentCheckRejectsGuidePageWithUnquotedPageTitle(t *testing.T) {
	document := validGuideDocument()
	for i, line := range document {
		if line == `page_title: "Alibaba Cloud Account Guide"` {
			document[i] = `page_title: Alibaba Cloud Account Guide`
			break
		}
	}
	if code := checkDocument(writeDocument(t, "website/docs/guides/getting-account.html.md", document)); code != 1 {
		t.Fatalf("guide page with an unquoted page_title returned exit code %d, want 1", code)
	}
}

func TestUnitDocumentCheckRejectsGuidePageWithWrongLayout(t *testing.T) {
	document := validGuideDocument()
	for i, line := range document {
		if line == `layout: "alicloud"` {
			document[i] = `layout: "default"`
			break
		}
	}
	if code := checkDocument(writeDocument(t, "website/docs/guides/getting-account.html.markdown", document)); code != 1 {
		t.Fatalf("guide page with a wrong layout returned exit code %d, want 1", code)
	}
}

func TestUnitDocumentCheckRejectsGuidePageWithChineseSiteLink(t *testing.T) {
	document := validGuideDocument()
	document = append(document, "",
		"For more details, see [Sign up with Alibaba Cloud](https://help.aliyun.com/document_detail/50482.html).")
	if code := checkDocument(writeDocument(t, "website/docs/guides/getting-account.html.md", document)); code != 1 {
		t.Fatalf("guide page with a Chinese site link returned exit code %d, want 1", code)
	}
}

func TestUnitDocumentCheckRejectsResourcePageWithGuideTitle(t *testing.T) {
	document := validResourceDocument()
	for i, line := range document {
		switch line {
		case `page_title: "Alicloud: alicloud_demo"`:
			document[i] = `page_title: "Alibaba Cloud Demo Resource Guide"`
		case "# alicloud_demo":
			document[i] = "# Alibaba Cloud Demo Resource Guide"
		}
	}
	if code := checkDocument(writeDocument(t, "website/docs/r/demo.html.markdown", document)); code != 1 {
		t.Fatalf("resource page with a guide-style title returned exit code %d, want 1", code)
	}
}

func TestUnitDocumentCheckRejectsResourcePageWithWrongObjectName(t *testing.T) {
	document := validResourceDocument()
	for i, line := range document {
		if line == "# alicloud_demo" {
			document[i] = "# alicloud_demo_guide"
			break
		}
	}
	if code := checkDocument(writeDocument(t, "website/docs/r/demo.html.markdown", document)); code != 1 {
		t.Fatalf("resource page with a wrong object name returned exit code %d, want 1", code)
	}
}

func TestUnitDocumentCheckRejectsMissingVersionNote(t *testing.T) {
	var document []string
	for _, line := range validResourceDocument() {
		if strings.HasPrefix(line, "-> **NOTE:** Available since v") {
			continue
		}
		document = append(document, line)
	}
	if code := checkDocument(writeDocument(t, "website/docs/r/demo.html.markdown", document)); code != 1 {
		t.Fatalf("page without a version note returned exit code %d, want 1", code)
	}
}

func TestUnitDocumentCheckRejectsNonTerraformExampleFence(t *testing.T) {
	document := validResourceDocument()
	for i, line := range document {
		if line == "```terraform" {
			document[i] = "```hcl"
			break
		}
	}
	if code := checkDocument(writeDocument(t, "website/docs/r/demo.html.markdown", document)); code != 1 {
		t.Fatalf("page with a non-terraform example fence returned exit code %d, want 1", code)
	}
}

// writeProviderFixture builds a minimal provider source tree holding one
// service package registration, as expected under alicloud/service, and
// returns its root. The checker reads these files as text, so the fixture
// does not need to compile.
func writeProviderFixture(t *testing.T, register, definition string) string {
	t.Helper()
	return writeProviderFixtureFiles(t, register, map[string]string{"resource.go": definition})
}

// writeProviderFixtureFiles is writeProviderFixture for packages defining
// several registered forms, e.g. a resource and a data source sharing one
// TypeName.
func writeProviderFixtureFiles(t *testing.T, register string, definitions map[string]string) string {
	t.Helper()
	root := t.TempDir()
	pkg := filepath.Join(root, "alicloud", "service", "demopkg")
	if err := os.MkdirAll(pkg, 0o755); err != nil {
		t.Fatalf("create fixture package: %v", err)
	}
	if err := os.WriteFile(filepath.Join(pkg, "register.go"), []byte(register), 0o644); err != nil {
		t.Fatalf("write register.go: %v", err)
	}
	for name, definition := range definitions {
		if err := os.WriteFile(filepath.Join(pkg, name), []byte(definition), 0o644); err != nil {
			t.Fatalf("write %s: %v", name, err)
		}
	}
	return root
}

func withProviderRoot(t *testing.T, root string) {
	t.Helper()
	original := providerRoot
	providerRoot = root
	t.Cleanup(func() { providerRoot = original })
}

func writeNamedDocument(t *testing.T, name string, lines []string) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), name)
	if err := os.WriteFile(path, []byte(strings.Join(lines, "\n")+"\n"), 0o644); err != nil {
		t.Fatalf("write document: %v", err)
	}
	return path
}

// writeKindDocument writes the page under a docs directory form (r, d or
// ephemeral-resources) so the checker derives the registration kind from the
// path, as production paths do.
func writeKindDocument(t *testing.T, dir, name string, lines []string) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), dir, name)
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatalf("create docs directory: %v", err)
	}
	if err := os.WriteFile(path, []byte(strings.Join(lines, "\n")+"\n"), 0o644); err != nil {
		t.Fatalf("write document: %v", err)
	}
	return path
}

func deprecatedResourceDocument() []string {
	document := resourceDocument("alicloud_deprecated")
	for i, line := range document {
		if strings.HasPrefix(line, "-> **NOTE:** Available since v") {
			document[i] = "-> **DEPRECATED:** This resource has been deprecated from v1.2.0."
		}
	}
	return document
}

const sdkRegister = `package demopkg

func init() {
	conns.ServicePackages = append(conns.ServicePackages, &conns.ServicePackage{
		SDKResources: []*conns.SDKResource{
			{
				TypeName: "alicloud_deprecated",
				Factory:  NewDeprecatedResource,
			},
		},
	})
}
`

const sdkDeprecatedDefinition = `package demopkg

func NewDeprecatedResource() *schema.Resource {
	return &schema.Resource{
		SchemaVersion: 1,
		DeprecationMessage: "This resource has been deprecated since v1.2.0.",
	}
}
`

func TestUnitDocumentCheckRejectsMissingDeprecatedNoteForDeprecatedType(t *testing.T) {
	root := writeProviderFixture(t, sdkRegister, sdkDeprecatedDefinition)
	withProviderRoot(t, root)
	// The page only states an Available note although the schema declares a
	// type-level DeprecationMessage, as in alicloud_cms_service.
	if code := checkDocument(writeNamedDocument(t, "deprecated.html.markdown", resourceDocument("alicloud_deprecated"))); code != 1 {
		t.Fatalf("page without a DEPRECATED note for a deprecated type returned exit code %d, want 1", code)
	}
}

func TestUnitDocumentCheckAcceptsDeprecatedNoteForDeprecatedType(t *testing.T) {
	root := writeProviderFixture(t, sdkRegister, sdkDeprecatedDefinition)
	withProviderRoot(t, root)
	if code := checkDocument(writeNamedDocument(t, "deprecated.html.markdown", deprecatedResourceDocument())); code != 0 {
		t.Fatalf("page with a DEPRECATED note for a deprecated type returned exit code %d, want 0", code)
	}
}

const frameworkRegister = `package demopkg

func init() {
	conns.ServicePackages = append(conns.ServicePackages, &conns.ServicePackage{
		Resources: []*conns.Resource{
			{
				TypeName: "alicloud_deprecated",
				Factory:  NewDeprecatedResource,
			},
		},
	})
}
`

const frameworkDeprecatedDefinition = `package demopkg

func NewDeprecatedResource() func() resource.Resource {
	return func() resource.Resource {
		return &deprecatedResource{}
	}
}

func (r *deprecatedResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		DeprecationMessage: "This resource has been deprecated since v1.2.0.",
	}
}
`

func TestUnitDocumentCheckRejectsMissingDeprecatedNoteForDeprecatedFrameworkType(t *testing.T) {
	root := writeProviderFixture(t, frameworkRegister, frameworkDeprecatedDefinition)
	withProviderRoot(t, root)
	if code := checkDocument(writeNamedDocument(t, "deprecated.html.markdown", resourceDocument("alicloud_deprecated"))); code != 1 {
		t.Fatalf("page without a DEPRECATED note for a deprecated framework type returned exit code %d, want 1", code)
	}
}

// frameworkAttributeDefinition is the Framework shape with only an
// attribute-level deprecation: the DeprecationMessage field belongs to the
// schema.StringAttribute inside the attribute map, not to the schema.Schema
// the Schema method assigns.
const frameworkAttributeDefinition = `package demopkg

func NewDeprecatedResource() func() resource.Resource {
	return func() resource.Resource {
		return &deprecatedResource{}
	}
}

func (r *deprecatedResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"status": schema.StringAttribute{
				Required:           true,
				DeprecationMessage: "Use is_enabled instead.",
			},
		},
	}
}
`

func TestUnitDocumentCheckAcceptsAvailableNoteForAttributeLevelFrameworkDeprecation(t *testing.T) {
	root := writeProviderFixture(t, frameworkRegister, frameworkAttributeDefinition)
	withProviderRoot(t, root)
	// The page states an ordinary Available note: only an attribute is
	// deprecated, so no type-level DEPRECATED statement is required.
	if code := checkDocument(writeNamedDocument(t, "deprecated.html.markdown", resourceDocument("alicloud_deprecated"))); code != 0 {
		t.Fatalf("page with an Available note for an attribute-level framework deprecation returned exit code %d, want 0", code)
	}
}

// frameworkMixedDefinition declares both levels at once: the type-level
// DeprecationMessage on schema.Schema and an attribute-level one on
// schema.StringAttribute. Only the type-level one demands a DEPRECATED
// statement, but the attribute-level token must not mask it either.
const frameworkMixedDefinition = `package demopkg

func NewDeprecatedResource() func() resource.Resource {
	return func() resource.Resource {
		return &deprecatedResource{}
	}
}

func (r *deprecatedResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		DeprecationMessage: "This resource has been deprecated since v1.2.0.",
		Attributes: map[string]schema.Attribute{
			"status": schema.StringAttribute{
				DeprecationMessage: "Use is_enabled instead.",
			},
		},
	}
}
`

func TestUnitDocumentCheckRejectsMissingDeprecatedNoteWithTypeLevelBesideAttributeLevel(t *testing.T) {
	root := writeProviderFixture(t, frameworkRegister, frameworkMixedDefinition)
	withProviderRoot(t, root)
	if code := checkDocument(writeNamedDocument(t, "deprecated.html.markdown", resourceDocument("alicloud_deprecated"))); code != 1 {
		t.Fatalf("page without a DEPRECATED note for a type-level deprecation beside an attribute-level one returned exit code %d, want 1", code)
	}
}

const sdkPlainDefinition = `package demopkg

func NewDeprecatedResource() *schema.Resource {
	return &schema.Resource{
		SchemaVersion: 1,
		Schema: map[string]*schema.Schema{
			"status": {
				Type:     schema.TypeString,
				Deprecated: "Use is_default instead.",
			},
		},
	}
}
`

func TestUnitDocumentCheckAcceptsAvailableNoteForNonDeprecatedType(t *testing.T) {
	root := writeProviderFixture(t, sdkRegister, sdkPlainDefinition)
	withProviderRoot(t, root)
	// An ordinary Available note stays valid when the type is not deprecated;
	// the attribute-level Deprecated: field must not trigger the check.
	if code := checkDocument(writeNamedDocument(t, "deprecated.html.markdown", resourceDocument("alicloud_deprecated"))); code != 0 {
		t.Fatalf("page with an Available note for a non-deprecated type returned exit code %d, want 0", code)
	}
}

const multiLineRegister = `package demopkg

func init() {
	conns.ServicePackages = append(conns.ServicePackages, &conns.ServicePackage{
		SDKDataSources: []*conns.SDKResource{
			{TypeName: "alicloud_other", Factory: NewOtherDataSource},
			{TypeName: "alicloud_deprecated", Factory: NewDeprecatedResource},
		},
	})
}
`

func TestUnitDocumentCheckRejectsMissingDeprecatedNoteForRegisteredDataSource(t *testing.T) {
	root := writeProviderFixture(t, multiLineRegister, sdkDeprecatedDefinition)
	withProviderRoot(t, root)
	if code := checkDocument(writeNamedDocument(t, "deprecated.html.markdown", resourceDocument("alicloud_deprecated"))); code != 1 {
		t.Fatalf("page without a DEPRECATED note for a deprecated data source returned exit code %d, want 1", code)
	}
}

// sharedNameRegister registers one TypeName as both an SDK resource and an
// SDK data source, as about twenty production names do (e.g.
// alicloud_fc_service or alicloud_log_alert_resource). The factories differ
// in whether their schema declares a type-level DeprecationMessage.
const sharedNameRegister = `package demopkg

func init() {
	conns.ServicePackages = append(conns.ServicePackages, &conns.ServicePackage{
		SDKResources: []*conns.SDKResource{
			{
				TypeName: "alicloud_shared",
				Factory:  NewSharedResource,
			},
		},
		SDKDataSources: []*conns.SDKDataSource{
			{
				TypeName: "alicloud_shared",
				Factory:  NewSharedDataSource,
			},
		},
	})
}
`

const sharedResourceDeprecated = `package demopkg

func NewSharedResource() *schema.Resource {
	return &schema.Resource{
		SchemaVersion: 1,
		DeprecationMessage: "This resource has been deprecated since v1.2.0.",
	}
}
`

const sharedDataSourceDeprecated = `package demopkg

func NewSharedDataSource() *schema.Resource {
	return &schema.Resource{
		SchemaVersion: 1,
		DeprecationMessage: "This data source has been deprecated since v1.2.0.",
	}
}
`

const plainSchemaDefinition = `package demopkg

func %s() *schema.Resource {
	return &schema.Resource{
		SchemaVersion: 1,
		Schema: map[string]*schema.Schema{
			"status": {
				Type:     schema.TypeString,
			},
		},
	}
}
`

func sharedPlainDefinition(factory string) string {
	return fmt.Sprintf(plainSchemaDefinition, factory)
}

// TestUnitDocumentCheckAcceptsNonDeprecatedDataSourceSharingNameWithDeprecatedResource
// reproduces the fc_service shape: the resource is deprecated while the data
// source of the same TypeName is not, so a d/ page with an ordinary
// Available note must not be asked for a DEPRECATED statement.
func TestUnitDocumentCheckAcceptsNonDeprecatedDataSourceSharingNameWithDeprecatedResource(t *testing.T) {
	root := writeProviderFixtureFiles(t, sharedNameRegister, map[string]string{
		"resource.go":    sharedResourceDeprecated,
		"data_source.go": sharedPlainDefinition("NewSharedDataSource"),
	})
	withProviderRoot(t, root)
	if code := checkDocument(writeKindDocument(t, "d", "shared.html.markdown", resourceDocument("alicloud_shared"))); code != 0 {
		t.Fatalf("non-deprecated data source page sharing a TypeName with a deprecated resource returned exit code %d, want 0", code)
	}
}

// TestUnitDocumentCheckRejectsMissingDeprecatedNoteForDataSourceShadowedByResource
// reproduces the log_alert_resource shape: the data source of the shared
// TypeName is deprecated while the resource is not; the resource entry must
// not mask the data source deprecation for a d/ page.
func TestUnitDocumentCheckRejectsMissingDeprecatedNoteForDataSourceShadowedByResource(t *testing.T) {
	root := writeProviderFixtureFiles(t, sharedNameRegister, map[string]string{
		"resource.go":    sharedPlainDefinition("NewSharedResource"),
		"data_source.go": sharedDataSourceDeprecated,
	})
	withProviderRoot(t, root)
	if code := checkDocument(writeKindDocument(t, "d", "shared.html.markdown", resourceDocument("alicloud_shared"))); code != 1 {
		t.Fatalf("page without a DEPRECATED note for a deprecated data source shadowed by a resource returned exit code %d, want 1", code)
	}
}

// TestUnitDocumentCheckReconcilesResourcePageAgainstResourceRegistration
// keeps the resource side of a shared TypeName honest: an r/ page must
// reconcile against the resource registration, not the data source entry.
func TestUnitDocumentCheckReconcilesResourcePageAgainstResourceRegistration(t *testing.T) {
	deprecatedResource := writeProviderFixtureFiles(t, sharedNameRegister, map[string]string{
		"resource.go":    sharedResourceDeprecated,
		"data_source.go": sharedPlainDefinition("NewSharedDataSource"),
	})
	withProviderRoot(t, deprecatedResource)
	if code := checkDocument(writeKindDocument(t, "r", "shared.html.markdown", resourceDocument("alicloud_shared"))); code != 1 {
		t.Fatalf("page without a DEPRECATED note for a deprecated resource returned exit code %d, want 1", code)
	}

	plainResource := writeProviderFixtureFiles(t, sharedNameRegister, map[string]string{
		"resource.go":    sharedPlainDefinition("NewSharedResource"),
		"data_source.go": sharedDataSourceDeprecated,
	})
	withProviderRoot(t, plainResource)
	if code := checkDocument(writeKindDocument(t, "r", "shared.html.markdown", resourceDocument("alicloud_shared"))); code != 0 {
		t.Fatalf("page with an Available note for a non-deprecated resource returned exit code %d, want 0", code)
	}
}

const ephemeralRegister = `package demopkg

func init() {
	conns.ServicePackages = append(conns.ServicePackages, &conns.ServicePackage{
		EphemeralResources: []conns.EphemeralResource{
			{
				TypeName: "alicloud_ephemeral",
				Factory:  NewSharedEphemeral,
			},
		},
	})
}
`

const ephemeralDeprecatedDefinition = `package demopkg

func NewSharedEphemeral() func() ephemeral.Resource {
	return func() ephemeral.Resource {
		return &sharedEphemeral{}
	}
}

func (r *sharedEphemeral) Schema(ctx context.Context, req ephemeral.SchemaRequest, resp *ephemeral.SchemaResponse) {
	resp.Schema = schema.Schema{
		DeprecationMessage: "This ephemeral resource has been deprecated since v1.2.0.",
	}
}
`

func ephemeralDocument(resourceName string) []string {
	document := resourceDocument(resourceName)
	for i, line := range document {
		if line == "# "+resourceName {
			document[i] = "# Ephemeral: " + resourceName
		}
	}
	return document
}

func deprecatedEphemeralDocument(resourceName string) []string {
	document := ephemeralDocument(resourceName)
	for i, line := range document {
		if strings.HasPrefix(line, "-> **NOTE:** Available since v") {
			document[i] = "-> **DEPRECATED:** This ephemeral resource has been deprecated from v1.2.0."
		}
	}
	return document
}

// TestUnitDocumentCheckReconcilesEphemeralPageAgainstEphemeralRegistration
// covers the third docs form: ephemeral-resources/ pages reconcile against
// the EphemeralResources registration section.
func TestUnitDocumentCheckReconcilesEphemeralPageAgainstEphemeralRegistration(t *testing.T) {
	root := writeProviderFixture(t, ephemeralRegister, ephemeralDeprecatedDefinition)
	withProviderRoot(t, root)
	if code := checkDocument(writeKindDocument(t, "ephemeral-resources", "ephemeral.html.markdown", ephemeralDocument("alicloud_ephemeral"))); code != 1 {
		t.Fatalf("page without a DEPRECATED note for a deprecated ephemeral resource returned exit code %d, want 1", code)
	}
	if code := checkDocument(writeKindDocument(t, "ephemeral-resources", "ephemeral.html.markdown", deprecatedEphemeralDocument("alicloud_ephemeral"))); code != 0 {
		t.Fatalf("page with a DEPRECATED note for a deprecated ephemeral resource returned exit code %d, want 0", code)
	}
}
