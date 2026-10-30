package vpngateway

// Local status type and constants used by the vpngateway resources, previously
// defined in the root package alicloud (common.go). They are reproduced here
// 1:1 with their original values so that the migrated package does not import
// the root package.
type Status string

const (
	Pending         = Status("Pending")
	Creating        = Status("Creating")
	Running         = Status("Running")
	Available       = Status("Available")
	Unavailable     = Status("Unavailable")
	Modifying       = Status("Modifying")
	Deleting        = Status("Deleting")
	Starting        = Status("Starting")
	Stopping        = Status("Stopping")
	Stopped         = Status("Stopped")
	Normal          = Status("Normal")
	Changing        = Status("Changing")
	Online          = Status("online")
	Configuring     = Status("configuring")
	Active          = Status("Active")
	Deleted         = Status("Deleted")
	Null            = Status("Null")
	Enable          = Status("Enable")
	Init            = Status("Init")
	Provisioning    = Status("Provisioning")
	Updating        = Status("Updating")
	FinancialLocked = Status("FinancialLocked")
)

const DefaultTimeout = 120
