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

func YamlIndent(depth int) string {
	return strings.Repeat("  ", depth)
}

func AdjustYamlIndent(yamlStr string, depth int) string {
	lines := strings.Split(yamlStr, "\n")
	for i, line := range lines {
		if line == "" {
			continue
		}
		if line == "logs:" {
			lines[i] = fmt.Sprintf("%s%s", YamlIndent(depth-1), line)
		} else {
			lines[i] = fmt.Sprintf("%s%s", YamlIndent(depth), line)
		}
	}
	return strings.Join(lines, "\n")
}
