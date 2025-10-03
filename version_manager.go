package gobrew

import (
	"fmt"
	"regexp"
	"sort"
	"strings"

	"github.com/Masterminds/semver"
)

// VersionManager handles all version-related operations
type VersionManager struct {
	gb *GoBrew
}

// NewVersionManager creates a new VersionManager instance
func NewVersionManager(gb *GoBrew) *VersionManager {
	return &VersionManager{gb: gb}
}

// ParseAndValidateVersion parses a version string and validates it
func (vm *VersionManager) ParseAndValidateVersion(version string) (*semver.Version, error) {
	if version == "" || version == NoneVersion {
		return nil, fmt.Errorf("no version provided")
	}

	// Clean version string by removing suffixes
	cleanVersion := vm.cleanVersionString(version)
	if cleanVersion == "" {
		return nil, fmt.Errorf("invalid version format: %s", version)
	}

	// Parse using semver library
	semverVersion, err := semver.NewVersion(cleanVersion)
	if err != nil {
		return nil, fmt.Errorf("failed to parse version %s: %w", version, err)
	}

	return semverVersion, nil
}

// ResolveVersion resolves a version string to a concrete version
func (vm *VersionManager) ResolveVersion(version string) (string, error) {
	if version == "" || version == NoneVersion {
		return "", fmt.Errorf("no version provided")
	}

	// Handle special version keywords
	switch version {
	case "latest":
		return vm.getLatestStableVersion(), nil
	case "dev-latest":
		return vm.getLatestDevVersion(), nil
	case "mod":
		return vm.resolveModVersion()
	}

	// Handle version patterns
	if strings.HasSuffix(version, "x") || strings.HasSuffix(version, ".x") {
		return vm.resolveVersionPattern(version)
	}

	// Handle @latest and @dev-latest suffixes
	if strings.HasSuffix(version, "@latest") {
		baseVersion := strings.TrimSuffix(version, "@latest")
		return vm.getLatestInMajor(baseVersion), nil
	}

	if strings.HasSuffix(version, "@dev-latest") {
		baseVersion := strings.TrimSuffix(version, "@dev-latest")
		return vm.getLatestDevInMajor(baseVersion), nil
	}

	// Direct version - validate it exists
	cleanVersion := vm.cleanVersionString(version)
	if !vm.versionExists(cleanVersion) {
		return "", fmt.Errorf("version %s does not exist", cleanVersion)
	}

	return cleanVersion, nil
}

// ExtractMajorVersion extracts the major version from a version string
func (vm *VersionManager) ExtractMajorVersion(version string) string {
	if version == "" || version == NoneVersion {
		return ""
	}

	semverVersion, err := semver.NewVersion(version)
	if err != nil {
		// Fallback to string manipulation for non-standard versions
		parts := strings.Split(version, ".")
		if len(parts) < 2 {
			return ""
		}
		// Remove rc and beta suffixes
		parts[1] = strings.Split(parts[1], "rc")[0]
		parts[1] = strings.Split(parts[1], "beta")[0]
		return strings.Join(parts[:2], ".")
	}

	return fmt.Sprintf("%d.%d", semverVersion.Major(), semverVersion.Minor())
}

// Helper methods

func (vm *VersionManager) cleanVersionString(version string) string {
	// Remove @latest and @dev-latest suffixes
	version = strings.TrimSuffix(version, "@latest")
	version = strings.TrimSuffix(version, "@dev-latest")

	// Remove .x and x suffixes
	version = strings.TrimSuffix(version, ".x")
	version = strings.TrimSuffix(version, "x")

	return version
}

func (vm *VersionManager) versionExists(version string) bool {
	remoteVersions := vm.gb.ListRemoteVersions(false)

	// Check if version exists in any major version group
	for _, versions := range remoteVersions {
		for _, v := range versions {
			if v == version {
				return true
			}
		}
	}
	return false
}

func (vm *VersionManager) getLatestStableVersion() string {
	remoteVersions := vm.gb.ListRemoteVersions(false)
	var stableVersions []*semver.Version

	// Collect all stable versions
	for _, versions := range remoteVersions {
		for _, version := range versions {
			if !vm.isBetaOrRC(version) {
				if v, err := semver.NewVersion(version); err == nil {
					stableVersions = append(stableVersions, v)
				}
			}
		}
	}

	if len(stableVersions) == 0 {
		return NoneVersion
	}

	sort.Sort(semver.Collection(stableVersions))
	return stableVersions[len(stableVersions)-1].String()
}

func (vm *VersionManager) getLatestDevVersion() string {
	remoteVersions := vm.gb.ListRemoteVersions(false)
	var allVersions []*semver.Version

	// Collect all versions including beta/rc
	for _, versions := range remoteVersions {
		for _, version := range versions {
			if v, err := semver.NewVersion(version); err == nil {
				allVersions = append(allVersions, v)
			}
		}
	}

	if len(allVersions) == 0 {
		return NoneVersion
	}

	sort.Sort(semver.Collection(allVersions))
	return allVersions[len(allVersions)-1].String()
}

func (vm *VersionManager) resolveVersionPattern(pattern string) (string, error) {
	basePattern := strings.TrimSuffix(pattern, "x")
	basePattern = strings.TrimSuffix(basePattern, ".")

	remoteVersions := vm.gb.ListRemoteVersions(false)
	var matchingVersions []*semver.Version

	// Find versions that match the pattern
	for _, versions := range remoteVersions {
		for _, version := range versions {
			if strings.HasPrefix(version, basePattern) {
				if v, err := semver.NewVersion(version); err == nil {
					matchingVersions = append(matchingVersions, v)
				}
			}
		}
	}

	if len(matchingVersions) == 0 {
		return "", fmt.Errorf("no versions found matching pattern: %s", pattern)
	}

	sort.Sort(semver.Collection(matchingVersions))
	return matchingVersions[len(matchingVersions)-1].String(), nil
}

func (vm *VersionManager) getLatestInMajor(majorVersion string) string {
	remoteVersions := vm.gb.ListRemoteVersions(false)
	var versionsInMajor []*semver.Version

	// Find versions in the same major version
	for _, versions := range remoteVersions {
		for _, version := range versions {
			if vm.ExtractMajorVersion(version) == majorVersion {
				if !vm.isBetaOrRC(version) {
					if v, err := semver.NewVersion(version); err == nil {
						versionsInMajor = append(versionsInMajor, v)
					}
				}
			}
		}
	}

	if len(versionsInMajor) == 0 {
		// If no stable versions found, try including beta/RC versions
		for _, versions := range remoteVersions {
			for _, version := range versions {
				if vm.ExtractMajorVersion(version) == majorVersion {
					if v, err := semver.NewVersion(version); err == nil {
						versionsInMajor = append(versionsInMajor, v)
					}
				}
			}
		}
	}

	if len(versionsInMajor) == 0 {
		return NoneVersion
	}

	sort.Sort(semver.Collection(versionsInMajor))
	return versionsInMajor[len(versionsInMajor)-1].String()
}

func (vm *VersionManager) getLatestDevInMajor(majorVersion string) string {
	remoteVersions := vm.gb.ListRemoteVersions(false)
	var versionsInMajor []*semver.Version

	// Find versions in the same major version (including dev)
	for _, versions := range remoteVersions {
		for _, version := range versions {
			if vm.ExtractMajorVersion(version) == majorVersion {
				if v, err := semver.NewVersion(version); err == nil {
					versionsInMajor = append(versionsInMajor, v)
				}
			}
		}
	}

	if len(versionsInMajor) == 0 {
		return NoneVersion
	}

	sort.Sort(semver.Collection(versionsInMajor))
	return versionsInMajor[len(versionsInMajor)-1].String()
}

func (vm *VersionManager) resolveModVersion() (string, error) {
	modVersion := vm.gb.getModVersion()
	if modVersion == "" || modVersion == NoneVersion {
		return "", fmt.Errorf("no go.mod version found")
	}

	// If mod version is already a full version, return it
	if strings.Count(modVersion, ".") >= 2 {
		if vm.versionExists(modVersion) {
			return modVersion, nil
		}
		// Try to resolve it directly (for cases like 1.25.1)
		return modVersion, nil
	}

	// If mod version is like 1.19, 1.20, append @latest
	return vm.ResolveVersion(modVersion + "@latest")
}

func (vm *VersionManager) isBetaOrRC(version string) bool {
	re := regexp.MustCompile("beta.*|rc.*")
	matches := re.FindAllString(version, -1)
	return len(matches) > 0
}