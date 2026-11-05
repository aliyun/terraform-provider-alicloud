package helper

import (
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
)

// BuildClientToken generates a unique client token for idempotent API
// calls, capped at 64 characters.
func BuildClientToken(action string) string {
	token := strings.TrimSpace(fmt.Sprintf("TF-%s-%d-%s", action, time.Now().Unix(), strings.Trim(uuid.New().String(), "-")))
	if len(token) > 64 {
		token = token[0:64]
	}
	return token
}
