/*

	MIT License

	Copyright (c) 2025 Evandro

	Permission is hereby granted, free of charge, to any person obtaining a copy
	of this software and associated documentation files (the "Software"), to deal
	in the Software without restriction, including without limitation the rights
	to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
	copies of the Software, and to permit persons to whom the Software is
	furnished to do so, subject to the following conditions:

	The above copyright notice and this permission notice shall be included in all
	copies or substantial portions of the Software.

	THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
	IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
	FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
	AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
	LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
	OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
	SOFTWARE.

*/

package sync

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/elaurentium/burrow/cmd/prompt"
	"github.com/elaurentium/burrow/internal/helper"
)

type GithubRelease struct {
	TagName string `json:"tag_name"`
	Name    string `json:"name"`
	Body    string `json:"body"`
	Assets  []struct {
		Name               string `json:"name"`
		BrowserDownloadUrl string `json:"browser_download_url"`
	} `json:"assets"`
}

func CheckForUpdates() (*GithubRelease, bool, error) {
	resp, err := http.Get(helper.GithubApi + "/releases/latest")
	if err != nil {
		return nil, false, fmt.Errorf("error checking for updates: %v", err)
	}
	defer func() {
		if closeErr := resp.Body.Close(); err != nil {
			fmt.Fprintf(os.Stderr, "Failed to close response body: %v", closeErr)
		}
	}()

	if resp.StatusCode != http.StatusOK {
		return nil, false, fmt.Errorf("error checking for updates: %s", resp.Status)
	}

	var release GithubRelease

	decoder := json.NewDecoder(resp.Body)
	if err := decoder.Decode(&release); err != nil {
		return nil, false, fmt.Errorf("error checking for updates: %v", err)
	}

	currentVersion := helper.Version
	latestVersion := strings.TrimPrefix(release.TagName, "v")
	if release.TagName != currentVersion {
		return &release, true, nil
	}

	isNewer, err := helper.IsVersionNewer(latestVersion, currentVersion)
	if err != nil {
		return nil, false, fmt.Errorf("failed to compare versions: %w", err)
	}

	return &release, isNewer, nil
}

func PerformUpdate(release *GithubRelease) error {
	// Find the appropriate asset for current platform
	assetURL, assetName, err := findAssetForPlatform(release)
	if err != nil {
		return err
	}

	fmt.Printf("Downloading update: %s\n", assetName)

	// Download the asset
	tempFile, err := downloadAsset(assetURL, assetName)
	if err != nil {
		return fmt.Errorf("failed to download update: %w", err)
	}
	defer func() {
		if removeErr := os.Remove(tempFile); removeErr != nil {
			fmt.Fprintf(os.Stderr, "Failed to remove temp file: %v\n", removeErr)
		}
	}()

	// Get current executable path
	currentExe, err := os.Executable()
	if err != nil {
		return fmt.Errorf("failed to get current executable path: %w", err)
	}

	// Create backup of current executable
	backupFile := currentExe + ".backup"
	if err := copyFile(currentExe, backupFile); err != nil {
		return fmt.Errorf("failed to create backup: %w", err)
	}
	defer func() {
		if removeErr := os.Remove(backupFile); removeErr != nil {
			fmt.Fprintf(os.Stderr, "Failed to remove backup file: %v\n", removeErr)
		}
	}()

	// Replace current executable
	if err := replaceExecutable(currentExe, tempFile); err != nil {
		// Try to restore backup if replacement fails
		if _, backupErr := os.Stat(backupFile); backupErr == nil {
			if restoreErr := copyFile(backupFile, currentExe); restoreErr != nil {
				fmt.Fprintf(os.Stderr, "Failed to restore backup: %v\n", restoreErr)
			} else {
				fmt.Println("Successfully restored backup after failed update")
			}
		}

		// Check if this is the Windows deferred update case
		if strings.Contains(err.Error(), "update script created") {
			// This is not actually an error - the update will complete after restart
			fmt.Println("Update will complete when you restart the application")
			return nil
		}

		return fmt.Errorf("failed to replace executable: %w", err)
	}

	fmt.Println("Update completed successfully! Please restart the application.")
	return nil

}

// PromptForUpdate shows an interactive prompt asking user if they want to update
func PromptForUpdate(release *GithubRelease) (bool, error) {
	fmt.Printf("Update Available\n"+
		"A new version of %s is available!\n\n"+
		"Current version: v%s\n"+
		"Latest version: %s\n\n"+
		"Release notes:\n%s\n\n",
		helper.Name, helper.Version, release.TagName, release.Body)

	input := prompt.NewPipe(os.Stdout, os.Stdin)
	confirmed, err := input.Confirm("Do you want to update? (y/n): ", true)
	if err != nil {
		return false, fmt.Errorf("failed to read user input: %w", err)
	}

	return confirmed, nil
}

// CheckAndPromptUpdate is a convenience function that checks for updates and prompts user
func CheckAndPromptUpdate() error {
	fmt.Println("Checking for updates...")

	release, hasUpdate, err := CheckForUpdates()
	if err != nil {
		return fmt.Errorf("failed to check for updates: %w", err)
	}

	if !hasUpdate {
		fmt.Println("You are running the latest version!")
		return nil
	}

	shouldUpdate, err := PromptForUpdate(release)
	if err != nil {
		return err
	}

	if shouldUpdate {
		return PerformUpdate(release)
	}

	if _, err := helper.UpdateVersion(release.TagName); err != nil {
		return fmt.Errorf("failed to update version: %w", err)
	}

	fmt.Println("Update cancelled by user")
	return nil
}

// CheckForUpdatesQuietly checks for updates without user interaction
func CheckForUpdatesQuietly() {
	release, hasUpdate, err := CheckForUpdates()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to check for updates: %v\n", err)
		return
	}

	if hasUpdate {
		fmt.Printf("New version available: %s (current: %s)\n",
			release.TagName, helper.Version)
		fmt.Println("Run with --update flag to update")
	}
}

// PlatformInfo holds platform-specific information
type PlatformInfo struct {
	OS   string
	Arch string
}

// GetCurrentPlatform returns the current platform information
func GetCurrentPlatform() PlatformInfo {
	return PlatformInfo{
		OS:   runtime.GOOS,
		Arch: runtime.GOARCH,
	}
}

func findAssetForPlatform(release *GithubRelease) (string, string, error) {
	return findAssetForPlatformWithInfo(release, GetCurrentPlatform())
}

func findAssetForPlatformWithInfo(release *GithubRelease, platform PlatformInfo) (string, string, error) {
	// Map platform names to expected asset names
	var expectedNames []string
	switch platform.OS {
	case "windows":
		expectedNames = []string{
			fmt.Sprintf("burrow-windows-%s.exe", platform.Arch),
			"burrow-windows.exe",
		}
	case "darwin":
		expectedNames = []string{
			fmt.Sprintf("burrow-darwin-%s", platform.Arch),
			fmt.Sprintf("burrow-macos-%s", platform.Arch),
			"burrow-darwin-universal", // Universal binary (explicit)
			"burrow-darwin",           // Universal binary (generic) or fallback
			"burrow-macos",            // Alternative generic name
		}
	case "linux":
		expectedNames = []string{
			fmt.Sprintf("burrow-linux-%s", platform.Arch),
			"burrow-linux",
		}
	default:
		return "", "", fmt.Errorf("unsupported platform: %s", platform.OS)
	}

	// Find matching asset
	for _, asset := range release.Assets {
		for _, expectedName := range expectedNames {
			if strings.EqualFold(asset.Name, expectedName) {
				return asset.BrowserDownloadUrl, asset.Name, nil
			}
		}
	}

	return "", "", fmt.Errorf("no compatible asset found for %s/%s", platform.OS, platform.Arch)
}

// validateGitHubURLWithTestFlag validates URLs with optional test mode support
func validateGitHubURLWithTestFlag(urlStr string, allowTestMode bool) error {
	parsedURL, err := url.Parse(urlStr)
	if err != nil {
		return fmt.Errorf("invalid URL: %w", err)
	}

	// Allow localhost and 127.0.0.1 in test mode
	if allowTestMode && (parsedURL.Host == "localhost" ||
		strings.HasPrefix(parsedURL.Host, "127.0.0.1") ||
		strings.HasPrefix(parsedURL.Host, "localhost:")) {
		return nil
	}

	// Only allow GitHub domains
	allowedHosts := []string{
		"github.com",
		"api.github.com",
		"objects.githubusercontent.com",
		"github-releases.githubusercontent.com",
	}

	hostAllowed := false
	for _, allowedHost := range allowedHosts {
		if parsedURL.Host == allowedHost {
			hostAllowed = true
			break
		}
	}

	if !hostAllowed {
		return fmt.Errorf("URL host %s not allowed", parsedURL.Host)
	}

	// Ensure HTTPS (except in test mode)
	if !allowTestMode && parsedURL.Scheme != "https" {
		return fmt.Errorf("only HTTPS URLs are allowed")
	}

	return nil
}

// validateFilePath validates that a file path is safe and doesn't contain directory traversal
func validateFilePath(path string) error {
	// Clean the path to resolve any .. or . components
	cleanPath := filepath.Clean(path)

	// Check for directory traversal attempts
	if strings.Contains(cleanPath, "..") {
		return fmt.Errorf("path contains directory traversal: %s", path)
	}

	// Ensure the path is absolute or within expected directories
	if !filepath.IsAbs(cleanPath) {
		// For relative paths, ensure they don't start with ..
		if strings.HasPrefix(cleanPath, "..") {
			return fmt.Errorf("relative path traversal detected: %s", path)
		}
	}

	return nil
}

// safeTempFile creates a validated temporary file path
func safeTempFile(filename string) (string, error) {
	// Sanitize filename
	filename = filepath.Base(filename) // Remove any path components
	if filename == "" || filename == "." || filename == ".." {
		return "", fmt.Errorf("invalid filename: %s", filename)
	}

	tempDir := os.TempDir()
	tempFile := filepath.Join(tempDir, "burrow-update-"+filename)

	// Validate the resulting path
	if err := validateFilePath(tempFile); err != nil {
		return "", fmt.Errorf("temp file path validation failed: %w", err)
	}

	return tempFile, nil
}

func downloadAsset(url, filename string) (string, error) {
	return downloadAssetWithTestFlag(url, filename, false)
}

// downloadAssetWithTestFlag downloads an asset with optional test mode support
func downloadAssetWithTestFlag(url, filename string, allowTestMode bool) (string, error) {
	var length int64
	// Validate URL before making request
	if err := validateGitHubURLWithTestFlag(url, allowTestMode); err != nil {
		return "", fmt.Errorf("URL validation failed: %w", err)
	}

	// #nosec G107 - URL is validated above to ensure it's from trusted GitHub domains
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer func() {
		if closeErr := resp.Body.Close(); closeErr != nil {
			fmt.Fprintf(os.Stderr, "Failed to close response body: %v\n", closeErr)
		}
	}()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("download failed with status %d", resp.StatusCode)
	}

	// Create temporary file with validation
	tempFile, err := safeTempFile(filename)
	if err != nil {
		return "", fmt.Errorf("failed to create safe temp file: %w", err)
	}

	// #nosec G304 - tempFile is validated by safeTempFile function above
	out, err := os.Create(tempFile)
	if err != nil {
		return "", err
	}
	defer func() {
		if closeErr := out.Close(); closeErr != nil {
			fmt.Fprintf(os.Stderr, "Failed to close temp file: %v\n", closeErr)
		}
	}()

	// Copy with progress indication
	var reader io.Reader = resp.Body
	if length > 0 {
		reader = &helper.DownloadProgress{
			Reader: reader,
			Length: length,
		}
		fmt.Fprintf(os.Stderr, "Downloading %s...\n", filename)
	}
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		if removeErr := os.Remove(tempFile); removeErr != nil {
			fmt.Fprintf(os.Stderr, "Failed to remove temp file after error: %v\n", removeErr)
		}
		return "", err
	}

	return tempFile, nil
}

func copyFile(src, dst string) error {
	// Validate source and destination paths
	if err := validateFilePath(src); err != nil {
		return fmt.Errorf("source path validation failed: %w", err)
	}
	if err := validateFilePath(dst); err != nil {
		return fmt.Errorf("destination path validation failed: %w", err)
	}

	// #nosec G304 - src path is validated above
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer func() {
		if closeErr := sourceFile.Close(); closeErr != nil {
			fmt.Fprintf(os.Stderr, "Failed to close source file: %v\n", closeErr)
		}
	}()

	// #nosec G304 - dst path is validated above
	destFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer func() {
		if closeErr := destFile.Close(); closeErr != nil {
			fmt.Fprintf(os.Stderr, "Failed to close destination file: %v\n", closeErr)
		}
	}()

	_, err = io.Copy(destFile, sourceFile)
	if err != nil {
		return err
	}

	// Copy file permissions
	sourceInfo, err := os.Stat(src)
	if err != nil {
		return err
	}

	size := sourceInfo.Size()

	var reader io.Reader = sourceFile
	if size > 0 {
		fmt.Fprintln(os.Stderr, "Copying...")
		reader = &helper.DownloadProgress{
			Reader: reader,
			Length: int64(size),
		}
	}

	if _, err = io.Copy(destFile, reader); err != nil {
		return err
	}

	return os.Chmod(dst, sourceInfo.Mode())
}

func replaceExecutable(currentExe, newExe string) error {
	// On Unix systems (Linux, macOS), we can move the running executable
	// and place the new one in its location. The running process continues
	// to use the moved executable until it exits.
	if runtime.GOOS != "windows" {
		// Generate a unique temporary name for the old executable
		tempName := currentExe + ".old." + fmt.Sprintf("%d", os.Getpid())

		// Move current executable to temporary location
		// This works even while the executable is running on Unix systems
		if err := os.Rename(currentExe, tempName); err != nil {
			return fmt.Errorf("failed to move current executable: %w", err)
		}

		// Schedule cleanup of the old executable
		defer func() {
			if removeErr := os.Remove(tempName); removeErr != nil {
				fmt.Fprintf(os.Stderr, "Failed to remove old executable: %v\n", removeErr)
			}
		}()

		// Copy new executable to the original location
		if err := copyFile(newExe, currentExe); err != nil {
			// Try to restore the original executable if copy fails
			if restoreErr := os.Rename(tempName, currentExe); restoreErr != nil {
				fmt.Fprintf(os.Stderr, "Failed to restore original executable: %v\n", restoreErr)
			}
			return fmt.Errorf("failed to copy new executable: %w", err)
		}

		// Make sure the new executable has proper permissions
		// #nosec G302 - 0755 is appropriate for executable files
		if err := os.Chmod(currentExe, 0755); err != nil {
			fmt.Fprintf(os.Stderr, "Failed to set executable permissions: %v\n", err)
			// Don't fail the update for permission issues
		}

		return nil
	}

	// Windows handling - more complex due to file locking
	// We need to use a different approach for Windows
	tempName := currentExe + ".old"

	// Try to rename the current executable
	// This may fail if the executable is in use
	if err := os.Rename(currentExe, tempName); err != nil {
		// If rename fails, try a different approach
		// Create a batch script that will replace the executable after this process exits
		return createWindowsUpdateScript(currentExe, newExe)
	}

	defer func() {
		if removeErr := os.Remove(tempName); removeErr != nil {
			fmt.Fprintf(os.Stderr, "Failed to remove old executable: %v\n", removeErr)
		}
	}()

	// Copy new executable to current location
	if err := copyFile(newExe, currentExe); err != nil {
		// Try to restore the original executable if copy fails
		if restoreErr := os.Rename(tempName, currentExe); restoreErr != nil {
			fmt.Fprintf(os.Stderr, "Failed to restore original executable: %v\n", restoreErr)
		}
		return fmt.Errorf("failed to copy new executable: %w", err)
	}

	return nil
}

// createWindowsUpdateScript creates a batch script to replace the executable
// after the current process exits. This is used as a fallback when the
// executable cannot be moved while running.
func createWindowsUpdateScript(currentExe, newExe string) error {
	if runtime.GOOS != "windows" {
		return fmt.Errorf("createWindowsUpdateScript called on non-Windows platform")
	}

	// Create a batch script that will:
	// 1. Wait for the current process to exit
	// 2. Replace the executable
	// 3. Clean up the temporary files
	// 4. Restart the application (optional)

	scriptPath := filepath.Join(filepath.Dir(currentExe), "update_burrow.bat")

	// Create the batch script content
	scriptContent := fmt.Sprintf(`@echo off
echo Updating burrow...
timeout /t 2 /nobreak > nul
:WAIT
tasklist /FI "PID eq %d" 2>NUL | find /I /N "%d">NUL
if "%%ERRORLEVEL%%"=="0" (
    timeout /t 1 /nobreak > nul
    goto WAIT
)
echo Replacing executable...
copy /Y "%s" "%s" > nul
if exist "%s" del "%s"
echo Update completed!
del "%%~f0"
`, os.Getpid(), os.Getpid(), newExe, currentExe, newExe, newExe)

	// Write the script to file
	if err := os.WriteFile(scriptPath, []byte(scriptContent), 0600); err != nil {
		return fmt.Errorf("failed to create update script: %w", err)
	}

	fmt.Println("Created update script. The application will be updated after exit.")
	fmt.Println("Please close the application to complete the update.")

	// Validate script path before executing
	if err := validateFilePath(scriptPath); err != nil {
		// Clean up script if path is invalid
		if removeErr := os.Remove(scriptPath); removeErr != nil {
			fmt.Fprintf(os.Stderr, "Failed to remove invalid script: %v\n", removeErr)
		}
		return fmt.Errorf("script path validation failed: %w", err)
	}

	// Execute the script in the background
	// Note: We don't wait for it to complete as it needs to run after this process exits
	// #nosec G204 - scriptPath is validated above for safety
	if err := exec.Command("cmd", "/C", "start", "/B", scriptPath).Start(); err != nil {
		// Clean up script if we can't execute it
		if removeErr := os.Remove(scriptPath); removeErr != nil {
			fmt.Fprintf(os.Stderr, "Failed to remove script after execution failure: %v\n", removeErr)
		}
		return fmt.Errorf("failed to execute update script: %w", err)
	}

	return fmt.Errorf("update script created - please restart the application to complete the update")
}
