package script

import (
	"regexp"
	"strings"

	"ezpkg.io/colorz"
)

func ProcessUsageText(s string) string {
	{
		re := regexp.MustCompile(`\{[^}]+}`)
		s = re.ReplaceAllStringFunc(s, func(s string) string {
			s = strings.TrimPrefix(s, "{")
			s = strings.TrimSuffix(s, "}")
			return colorz.Yellow.Wrap(s)
		})
	}
	{
		re := regexp.MustCompile(`(^|\n)#\s+([^\n]+\n)`)
		s = re.ReplaceAllStringFunc(s, func(s string) string {
			matches := re.FindStringSubmatch(s)
			return matches[1] + colorz.Yellow.Wrap(matches[2])
		})
	}
	return s
}
