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

package helper

import (
	"fmt"
	"strconv"
	"strings"
)

const (
	Name      = "Burrow"
	Usage     = "Burrow"
	Owner     = "elaurentium"
	GithubApi = "https://api.github.com/repos/" + Owner + "/" + Name

	// Build tools
	Makefile   = "Makefile"
	CMakeLists = "CMakeLists"
	Rakefile   = "Rakefile"
	Jakefile   = "Jakefile"
	Gruntfile  = "Gruntfile"
	Gulpfile   = "Gulpfile"

	// Container/VM
	Dockerfile    = "Dockerfile"
	Containerfile = "Containerfile"
	Vagrantfile   = "Vagrantfile"

	// CI/CD
	Jenkinsfile = "Jenkinsfile"
	Procfile    = "Procfile"
	Buildfile   = "Buildfile"

	// Package managers
	Gemfile  = "Gemfile"
	Podfile  = "Podfile"
	Cartfile = "Cartfile"
	Brewfile = "Brewfile"

	// Documentation
	README       = "README"
	LICENSE      = "LICENSE"
	CHANGELOG    = "CHANGELOG"
	CONTRIBUTING = "CONTRIBUTING"
	AUTHORS      = "AUTHORS"
	CONTRIBUTORS = "CONTRIBUTORS"
	COPYING      = "COPYING"
	INSTALL      = "INSTALL"
	NEWS         = "NEWS"
	TODO         = "TODO"
	HISTORY      = "HISTORY"
	NOTICE       = "NOTICE"

	// Config files
	Cakefile  = "Cakefile"
	Capfile   = "Capfile"
	Guardfile = "Guardfile"

	// Version control
	CODEOWNERS = "CODEOWNERS"
)

var (
	FilesWithoutExtension = []string{
		// Build tools
		Makefile, CMakeLists, Rakefile, Jakefile, Gruntfile, Gulpfile,

		// Container/VM
		Dockerfile, Containerfile, Vagrantfile,

		// CI/CD
		Jenkinsfile, Procfile, Buildfile,

		// Package managers
		Gemfile, Podfile, Cartfile, Brewfile,

		// Documentation
		README, LICENSE, CHANGELOG, CONTRIBUTING, AUTHORS, CONTRIBUTORS,
		COPYING, INSTALL, NEWS, TODO, HISTORY, NOTICE,

		// Config files
		Cakefile, Capfile, Guardfile,

		// Version control
		CODEOWNERS,
	}
)

func IsVersionNewer(latest, current string) (bool, error) {
	latestParts := strings.Split(latest, ".")
	currentParts := strings.Split(current, ".")

	// Pad shorter version with zeros
	maxLen := len(latestParts)
	if len(currentParts) > maxLen {
		maxLen = len(currentParts)
	}

	for len(latestParts) < maxLen {
		latestParts = append(latestParts, "0")
	}
	for len(currentParts) < maxLen {
		currentParts = append(currentParts, "0")
	}

	// Compare each part
	for i := 0; i < maxLen; i++ {
		latestNum, err := strconv.Atoi(latestParts[i])
		if err != nil {
			return false, fmt.Errorf("invalid version format in latest: %s", latest)
		}

		currentNum, err := strconv.Atoi(currentParts[i])
		if err != nil {
			return false, fmt.Errorf("invalid version format in current: %s", current)
		}

		if latestNum > currentNum {
			return true, nil
		} else if latestNum < currentNum {
			return false, nil
		}
	}

	return false, nil // Versions are equal
}
