package models

// ApplicationVersion represents the application version information
type ApplicationVersion struct {
	MobileAppVersion      string `json:"mobileAppVersion"`      // Mobile app version (e.g., "1.0.0")
	MobileAppBuildNumber  string `json:"mobileAppBuildNumber"`  // Mobile app build number
	WebVersion            string `json:"webVersion"`            // Web version
}
