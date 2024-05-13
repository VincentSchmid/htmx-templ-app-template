package main

import (
	"fmt"
	"regexp"
	"strings"
)

// permissionMapping defines mappings including patterns for dynamic paths and wildcards.
var permissionMapping = map[string]string{
	"^/v1/user/clients[^0-9]*$":                        "client",
	"^/v1/user/contracts[^0-9]*$":                      "contract",
	"^/v1/user/contracts/([a-fA-F0-9\\-]{36})[^0-9]*$": "contract:{uuid}",
	"^/v1/user/invites[^0-9]*$":                        "invite",
}

var permissionMapping2 = map[string]string{
	"/v1/user/clients*":          "client",
	"/v1/user/contracts*":        "contract",
	"/v1/user/contracts/{uuid}*": "contract:{uuid}",
	"/v1/user/invites*":          "invite",
}

var testPaths = []string{
	"/v1/user/clients/view",
	"/v1/user/contracts/new",
	"/v1/user/contracts",
	"/v1/user/contracts/f0c66c6d-e340-4776-ac45-20c34521e7c5/edit",
	"/v1/user/invites/send",
}

var regexes = map[*regexp.Regexp]string{
	regexp.MustCompile("^/v1/user/clients[^0-9]*$"):                      "client",
	regexp.MustCompile("^/v1/user/contracts[^0-9]*$"):                    "contract",
	regexp.MustCompile(`^/v1/user/contracts/([a-fA-F0-9-]{36})[^0-9]*$`): "contract:{uuid}",
	regexp.MustCompile("^/v1/user/invites[^0-9]*$"):                      "invite",
}

// findPermissionForPath attempts to find a permission for the given path,
// accounting for dynamic segments and wildcards.
func findPermissionForPath(path string) (string, bool) {
	for pattern, permission := range permissionMapping {
		re := regexp.MustCompile(pattern)
		if matches := re.FindStringSubmatch(path); matches != nil {
			if strings.Contains(permission, "{uuid}") && len(matches) > 1 {
				return strings.Replace(permission, "{uuid}", matches[1], 1), true
			}
			return permission, true
		}
	}

	return "", false
}

func findPermissionForPath2(path string) (string, bool) {
	for pattern, permission := range permissionMapping2 {
		re := regexp.MustCompile(pattern)
		if matches := re.FindStringSubmatch(path); matches != nil {
			if strings.Contains(permission, "{uuid}") && len(matches) > 1 {
				return strings.Replace(permission, "{uuid}", matches[1], 1), true
			}
			return permission, true
		}
	}

	return "", false
}

func findPermissionForPath3(path string) (string, bool) {
	for re, permission := range regexes {
		if matches := re.FindStringSubmatch(path); matches != nil {
			if strings.Contains(permission, "{uuid}") && len(matches) > 1 {
				return strings.Replace(permission, "{uuid}", matches[1], 1), true
			}
			return permission, true
		}
	}

	return "", false
}

func main() {
	for _, path := range testPaths {
		if permission, found := findPermissionForPath3(path); found {
			fmt.Printf("Permission for '%s': %s\n", path, permission)
		} else {
			fmt.Printf("No permission found for '%s'\n", path)
		}
	}
}
