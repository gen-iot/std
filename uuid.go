package std

import (
	guuid "github.com/satori/go.uuid"
	"strings"
)

//创建随机uuid id
func GenRandomUUID() string {
	return strings.Replace(guuid.NewV4().String(), "-", "", -1)
}
