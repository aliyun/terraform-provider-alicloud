package slb_test

// Shared acceptance-test sentinel constants for the split slb test package.
// The flat-layout originals live in the root alicloud package
// (alicloud/service_alicloud_common_test.go) and in alicloud/acctest/common.go
// and cannot be referenced unqualified from a service test package. The slb
// tests use these only as map values handed to acctest.ResourceAttrInit, whose
// machinery compares the literal strings below; the values are therefore kept
// identical to the canonical definitions so behaviour is unchanged. Sibling
// service packages (ecs, agentloop, openapiexplorer) keep the same local copy.

const (
	NOSET      = "#NOSET"       // be equivalent to method "TestCheckNoResourceAttrSet"
	CHECKSET   = "#CHECKSET"    // "TestCheckResourceAttrSet"
	REMOVEKEY  = "#REMOVEKEY"   // remove checkMap key
	CLEARLIST  = "#CLEARLIST"   // clear checkMap key and set empty list
	CLEARMAP   = "#CLEARMAP"    // clear checkMap key and set empty map
	REGEXMATCH = "#REGEXMATCH:" // "TestMatchResourceAttr" ,the map name/key like `"attribute" : REGEXMATCH + "attributeString"`
	ForceSleep = "force_sleep"
)
