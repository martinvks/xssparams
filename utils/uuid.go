package utils

import (
	"github.com/google/uuid"
	"strings"
)

func MiniUuid() string {
	return strings.Split(uuid.NewString(), "-")[0]
}
