package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// linkTestIndex builds a small website/docs tree in a temp directory and
// returns the resulting docIndex plus the repo root so callers can also write
// documents that checkDocument will resolve against it.
func linkTestIndex(t *testing.T) (*docIndex, string) {
	t.Helper()
	root := t.TempDir()
	for _, dir := range []string{"r", "d", "functions", "guides", "actions", "ephemeral-resources"} {
		if err := os.MkdirAll(filepath.Join(root, "website", "docs", dir), 0o755); err != nil {
			t.Fatalf("mkdir %s: %v", dir, err)
		}
	}
	write := func(dir, name string) {
		if err := os.WriteFile(filepath.Join(root, "website", "docs", dir, name+".html.markdown"), []byte("# stub\n"), 0o644); err != nil {
			t.Fatalf("write %s/%s: %v", dir, name, err)
		}
	}
	write("r", "foo")
	write("r", "bar")
	write("r", "ecs_disk")
	write("d", "foos")
	write("d", "bars")
	write("functions", "arn_build")
	write("guides", "getting-account")
	return buildDocIndex(root), root
}

func assertLinkError(t *testing.T, text string, idx *docIndex) {
	t.Helper()
	if msgs := checkCrossReferenceLinks(text, 1, idx); len(msgs) == 0 {
		t.Fatalf("expected a link error for %q, got none", text)
	}
}

func assertNoLinkError(t *testing.T, text string, idx *docIndex) {
	t.Helper()
	if msgs := checkCrossReferenceLinks(text, 1, idx); len(msgs) != 0 {
		t.Fatalf("expected no link error for %q, got: %s", text, strings.Join(msgs, "; "))
	}
}

func TestLinkCheckAcceptsValidRegistryResourceLink(t *testing.T) {
	idx, _ := linkTestIndex(t)
	assertNoLinkError(t, "see [bar](https://registry.terraform.io/providers/aliyun/alicloud/latest/docs/resources/bar)", idx)
}

func TestLinkCheckAcceptsValidRegistryDataSourceLinkWithAnchor(t *testing.T) {
	idx, _ := linkTestIndex(t)
	assertNoLinkError(t, "see [foos](https://registry.terraform.io/providers/aliyun/alicloud/latest/docs/data-sources/foos#argument-reference)", idx)
}

func TestLinkCheckAcceptsValidFunctionAndGuideLinks(t *testing.T) {
	idx, _ := linkTestIndex(t)
	assertNoLinkError(t, "uses [arn_build](https://registry.terraform.io/providers/aliyun/alicloud/latest/docs/functions/arn_build)", idx)
	assertNoLinkError(t, "see [getting-account](https://registry.terraform.io/providers/aliyun/alicloud/latest/docs/guides/getting-account)", idx)
}

func TestLinkCheckIgnoresThirdPartyProviderLink(t *testing.T) {
	idx, _ := linkTestIndex(t)
	assertNoLinkError(t, "see [vault](https://registry.terraform.io/providers/hashicorp/vault/latest/docs/resources/database)", idx)
}

func TestLinkCheckAcceptsProviderIndexAnchor(t *testing.T) {
	idx, _ := linkTestIndex(t)
	assertNoLinkError(t, "see [account_id](https://registry.terraform.io/providers/aliyun/alicloud/latest/docs#assume-role)", idx)
}

func TestLinkCheckAcceptsSectionIndexPage(t *testing.T) {
	idx, _ := linkTestIndex(t)
	assertNoLinkError(t, "browse [resources](https://registry.terraform.io/providers/aliyun/alicloud/latest/docs/resources)", idx)
}

func TestLinkCheckRejectsNonExistentResourceTarget(t *testing.T) {
	idx, _ := linkTestIndex(t)
	assertLinkError(t, "see [nope](https://registry.terraform.io/providers/aliyun/alicloud/latest/docs/resources/nope)", idx)
}

func TestLinkCheckRejectsSectionMismatch(t *testing.T) {
	idx, _ := linkTestIndex(t)
	// foo is a resource; linking it under data-sources is a mismatch.
	assertLinkError(t, "see [foo](https://registry.terraform.io/providers/aliyun/alicloud/latest/docs/data-sources/foo)", idx)
}

func TestLinkCheckRejectsResidualHtmlSuffix(t *testing.T) {
	idx, _ := linkTestIndex(t)
	assertLinkError(t, "see [foo](https://registry.terraform.io/providers/aliyun/alicloud/latest/docs/resources/foo.html)", idx)
}

func TestLinkCheckRejectsOldFormatLink(t *testing.T) {
	idx, _ := linkTestIndex(t)
	assertLinkError(t, "see [foo](https://www.terraform.io/docs/providers/alicloud/r/foo)", idx)
}

func TestLinkCheckRejectsOldFormatProviderIndex(t *testing.T) {
	idx, _ := linkTestIndex(t)
	assertLinkError(t, "see [account_id](https://www.terraform.io/docs/providers/alicloud/index.html#account_id)", idx)
}

func TestLinkCheckReportsFormatAndExistenceForOldBrokenLink(t *testing.T) {
	idx, _ := linkTestIndex(t)
	msgs := checkCrossReferenceLinks("see [nope](https://www.terraform.io/docs/providers/alicloud/r/nope)", 1, idx)
	if len(msgs) < 2 {
		t.Fatalf("expected format + existence errors, got %d: %v", len(msgs), msgs)
	}
}

func TestLinkCheckHandlesBareUrlWithTrailingPunctuation(t *testing.T) {
	idx, _ := linkTestIndex(t)
	// A bare URL followed by sentence punctuation must not see the period as
	// part of the target name; foo exists, so no error is expected.
	assertNoLinkError(t, "see https://registry.terraform.io/providers/aliyun/alicloud/latest/docs/resources/foo.", idx)
}

func TestLinkCheckSkipsWhenIndexIsNil(t *testing.T) {
	// Paths outside website/docs (unit-test fixtures) must skip the check so
	// existing structure tests remain unaffected by link content.
	msgs := checkCrossReferenceLinks("see [nope](https://registry.terraform.io/providers/aliyun/alicloud/latest/docs/resources/nope)", 1, nil)
	if len(msgs) != 0 {
		t.Fatalf("expected no link error with nil index, got: %v", msgs)
	}
}

func TestDocumentCheckIntegratesValidLink(t *testing.T) {
	_, root := linkTestIndex(t)
	doc := validResourceDocument()
	withLink := append(append([]string{}, doc[:8]...),
		"  See [bar](https://registry.terraform.io/providers/aliyun/alicloud/latest/docs/resources/bar).")
	withLink = append(withLink, doc[8:]...)
	path := filepath.Join(root, "website", "docs", "r", "demo.html.markdown")
	if err := os.WriteFile(path, []byte(strings.Join(withLink, "\n")+"\n"), 0o644); err != nil {
		t.Fatalf("write: %v", err)
	}
	if code := checkDocument(path); code != 0 {
		t.Fatalf("valid doc with a valid cross-reference link rejected: exit %d", code)
	}
}

func TestDocumentCheckIntegratesBrokenLink(t *testing.T) {
	_, root := linkTestIndex(t)
	doc := validResourceDocument()
	broken := append(append([]string{}, doc[:8]...),
		"  See [nope](https://registry.terraform.io/providers/aliyun/alicloud/latest/docs/resources/nope).")
	broken = append(broken, doc[8:]...)
	path := filepath.Join(root, "website", "docs", "r", "demo.html.markdown")
	if err := os.WriteFile(path, []byte(strings.Join(broken, "\n")+"\n"), 0o644); err != nil {
		t.Fatalf("write: %v", err)
	}
	if code := checkDocument(path); code != 1 {
		t.Fatalf("doc with a broken cross-reference link should fail, got exit %d", code)
	}
}

func TestDocumentCheckIntegratesOldFormatLink(t *testing.T) {
	_, root := linkTestIndex(t)
	doc := validResourceDocument()
	old := append(append([]string{}, doc[:8]...),
		"  See [bar](https://www.terraform.io/docs/providers/alicloud/r/bar).")
	old = append(old, doc[8:]...)
	path := filepath.Join(root, "website", "docs", "r", "demo.html.markdown")
	if err := os.WriteFile(path, []byte(strings.Join(old, "\n")+"\n"), 0o644); err != nil {
		t.Fatalf("write: %v", err)
	}
	if code := checkDocument(path); code != 1 {
		t.Fatalf("doc with a deprecated www.terraform.io link should fail, got exit %d", code)
	}
}
