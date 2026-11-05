package vpc

// COLON_SEPARATED is the ":" separator used by composite-ID resources in the
// VPC service package (HaVip attachment, Network Acl attachment/entries). The
// root alicloud package cannot be imported from a service package, so this
// local constant mirrors the root-package value verbatim.
const COLON_SEPARATED = ":"
