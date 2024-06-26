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
	prefix := fmt.Sprintf("%x", md5.Sum([]byte(namespace)))
	if strings.HasPrefix(match, prefix) { // idempotent
		return match
	}
	return fmt.Sprintf("%x.%s", md5.Sum([]byte(namespace)), match)
}

func GenerateNamespacedMatchRegExpr(namespace string, matchRegex string) string {
	matchRegex = strings.TrimPrefix(matchRegex, "^")
	prefix := fmt.Sprintf("^%x", md5.Sum([]byte(namespace)))
	if strings.HasPrefix(matchRegex, prefix) { // idempotent
		return matchRegex
	}
	return fmt.Sprintf("^%x\\.%s", md5.Sum([]byte(namespace)), matchRegex)
}

func YamlIndent(depth int) string {
	return strings.Repeat("  ", depth)
}

func AdjustYamlIndent(yamlStr string, depth int) string {
	lines := strings.Split(yamlStr, "\n")
	nameIdx := -1
	for i, line := range lines {
		if line == "" {
			continue
		}
		if line == "logs:" {
			lines[i] = fmt.Sprintf("%s%s", YamlIndent(depth-1), line)
		} else if strings.HasPrefix(line, "name:") {
			lines[i] = fmt.Sprintf("%s- %s", YamlIndent(depth-1), line)
			nameIdx = i
		} else {
			lines[i] = fmt.Sprintf("%s%s", YamlIndent(depth), line)
		}
	}
	if nameIdx != -1 {
		// shuffle the name line  to the first
		lines[0], lines[nameIdx] = lines[nameIdx], lines[0]
	}
	return strings.Join(lines, "\n")
}
