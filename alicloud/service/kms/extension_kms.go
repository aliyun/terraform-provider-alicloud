package kms

type KeyState string

const (
	Enabled         = KeyState("Enabled")
	Disabled        = KeyState("Disabled")
	PendingDeletion = KeyState("PendingDeletion")
)
