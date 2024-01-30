package utils

import (
	"crypto/md5"
	"fmt"
	"strings"
)

func HashCode(msg string) string {
	var h = md5.New()
	h.Write([]byte(msg))
	return string(h.Sum(nil))
}

func GenerateNamespacedMatchExpr(namespace string, match string) string {
	return fmt.Sprintf("%x.%s", md5.Sum([]byte(namespace)), match)
}

func GenerateNamespacedMatchRegExpr(namespace string, matchRegex string) string {
	matchRegex = strings.TrimPrefix(matchRegex, "^")
	return fmt.Sprintf("^%x\\.%s", md5.Sum([]byte(namespace)), matchRegex)
}
