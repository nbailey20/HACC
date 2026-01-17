package helpers

import (
	"fmt"
	"time"
)

func GeneratePassword() string {
	return fmt.Sprintf("%d", time.Now().Unix())
}
