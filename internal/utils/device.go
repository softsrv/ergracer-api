package utils

import (
	"strings"
)

func DetectDeviceType(userAgent string) string {
	ua := strings.ToLower(userAgent)
	
	if strings.Contains(ua, "mobile") || strings.Contains(ua, "android") {
		if strings.Contains(ua, "android") {
			return "android"
		}
		return "mobile"
	}
	
	if strings.Contains(ua, "iphone") || strings.Contains(ua, "ipad") || strings.Contains(ua, "ipod") {
		return "ios"
	}
	
	if strings.Contains(ua, "windows") {
		return "windows"
	}
	
	if strings.Contains(ua, "macintosh") || strings.Contains(ua, "mac os") {
		return "macos"
	}
	
	if strings.Contains(ua, "linux") {
		return "linux"
	}
	
	if strings.Contains(ua, "firefox") {
		return "firefox"
	}
	
	if strings.Contains(ua, "chrome") {
		return "chrome"
	}
	
	if strings.Contains(ua, "safari") {
		return "safari"
	}
	
	if strings.Contains(ua, "edge") {
		return "edge"
	}
	
	return "unknown"
}