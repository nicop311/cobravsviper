// MIT License
//
// Copyright (c) 2025 nicop311. All rights reserved.
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT.

package version

import (
	"encoding/json"
	"fmt"

	go_version "github.com/hashicorp/go-version"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

// populated by the Go LDFLAGS at build
var (
	RawGitDescribe     string
	GitDirtyStr        string // "true" or "false" but as strings as they are retrieved from git bash
	GitCommitIdShort   string
	GitCommitIdLong    string
	GitCommitTimestamp string
	GoVersion          string
	BuildPlatform      string
	BuildDate          string
)

// VersionOutput represents the JSON output structure.
type VersionOutput struct {
	VersionData VersionData `json:"cobravsviper-cli"`
}

// VersionData holds structured versioning details.
type VersionData struct {
	Major              uint64 `json:"major"`
	Minor              uint64 `json:"minor"`
	Patch              uint64 `json:"patch"`
	Version            string `json:"version"` // raw git describe
	IsGitDirty         bool   `json:"isGitDirty"`
	GitCommitIdLong    string `json:"gitCommitIdLong"`
	GitCommitIdShort   string `json:"gitCommitIdShort"`
	GitCommitTimestamp string `json:"gitCommitTimestamp"`
	GoVersion          string `json:"goVersion"`
	BuildDate          string `json:"buildDate"`
	BuildPlatform      string `json:"buildPlatform"`
}

// IsPopulated checks if the global variables for version information are populated.
// Returns true if at least RawGitDescribe is not empty, false otherwise.
// If false, most probably this is an issue with LDFLAGS.
func IsPopulated() bool {
	return RawGitDescribe != ""
}

// IsDirty takes a string from a build flag and returns a boolean indicating whether
// the build is from a dirty git tree.
func IsDirty(isDirtyStr string) (bool, error) {
	switch isDirtyStr {
	case "true":
		return true, nil
	case "false":
		return false, nil
	default:
		logrus.WithField("GitDirtyStr", isDirtyStr).Warn("Unexpected Git dirty string, assuming clean")
		return false, fmt.Errorf("invalid dirty information: %s", GitDirtyStr)
	}
}

// unset (zero v0.0.0).
func NewVersionData() (VersionData, error) {
	// this is a minimal content of the VersionData information
	versionData := VersionData{
		Version:            RawGitDescribe,
		GitCommitIdLong:    GitCommitIdLong,
		GitCommitIdShort:   GitCommitIdShort,
		GitCommitTimestamp: GitCommitTimestamp,
		GoVersion:          GoVersion,
		BuildDate:          BuildDate,
		BuildPlatform:      BuildPlatform,
	}

	// add the git state dirty true or false
	isDirty, err := IsDirty(GitDirtyStr)
	if err != nil {
		// only do a warning, do not return an error
		logrus.WithError(err).Warning("Failed to parse Git dirty status")
	}
	versionData.IsGitDirty = isDirty

	// Check if RawGitDescribe is a valid semantic version or a commit hash
	version, err := go_version.NewSemver(RawGitDescribe)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"raw_git_describe": RawGitDescribe,
			"error":            err,
		}).Debug("Invalid semantic versioning, falling back to snapshot version")
		return versionData, nil
	}

	// If version parsing is successful, populate the major, minor, and patch fields
	versionSegments := version.Segments()
	if len(versionSegments) < 3 {
		err = fmt.Errorf("raw git describe --tags --always version %s is not parsable as "+
			"semantic versioning. Expected 3 segments (major, minor, patch) but got %d",
			RawGitDescribe, len(versionSegments))
		return versionData, err
	}

	// set major, minor and patch to values that have been parsed by go_version
	versionData.Major = uint64(versionSegments[0])
	versionData.Minor = uint64(versionSegments[1])
	versionData.Patch = uint64(versionSegments[2])

	return versionData, nil
}

// returnJsonVersion returns the version as a JSON object.
func returnJsonVersion(prettyPrint bool) ([]byte, error) {
	versionOutput, err := NewVersionData()
	if err != nil {
		return nil, err
	}

	if prettyPrint {
		return json.MarshalIndent(versionOutput, "", "  ")
	}
	return json.Marshal(versionOutput)
}

// returnYamlVersion returns the version as a YAML object.
func returnYamlVersion() ([]byte, error) {
	jsonData, err := returnJsonVersion(false)
	if err != nil {
		logrus.WithError(err).Error("Failed to marshal JSON for YAML conversion")
		return nil, err
	}

	var data map[string]interface{}
	if err := json.Unmarshal(jsonData, &data); err != nil {
		logrus.WithError(err).Error("Failed to unmarshal JSON for YAML conversion")
		return nil, err
	}

	yamlData, err := yaml.Marshal(data)
	if err != nil {
		logrus.WithError(err).Error("Failed to marshal YAML")
		return nil, err
	}
	return yamlData, nil
}

// LogrusOutputVersion logs the version details at server startup. For server logging.
func LogrusOutputVersion() {
	versionData, err := NewVersionData()
	if err != nil {
		logrus.WithError(err).Error("Failed to fetch version data")
		return
	}

	logrus.Infof("cobravsviper version: %s", versionData.Version)
	logrus.WithFields(logrus.Fields{
		"build-date":       versionData.BuildDate,
		"build-platform":   versionData.BuildPlatform,
		"commit":           versionData.GitCommitIdLong,
		"go-version":       versionData.GoVersion,
		"raw-git-describe": versionData.Version,
		"is-git-dirty":     versionData.IsGitDirty,
		"short-commit":     versionData.GitCommitIdShort,
	}).Debug("cobravsviper version details")
}

// VersionOutputToString returns the version as a formatted string.
func VersionOutputToString(outputFormat string, prettyPrint bool) string {
	switch outputFormat {
	case "json":
		data, err := returnJsonVersion(prettyPrint)
		if err != nil {
			logrus.WithError(err).Error("Failed to generate JSON version output")
			return "Error generating JSON output"
		}
		return string(data)
	case "yaml":
		data, err := returnYamlVersion()
		if err != nil {
			logrus.WithError(err).Error("Failed to generate YAML version output")
			return "Error generating YAML output"
		}
		return string(data)
	default:
		version, err := go_version.NewSemver(RawGitDescribe)
		if err != nil {
			logrus.WithError(err).Debug("Invalid semantic versioning, falling back to snapshot version")
			return fmt.Sprintf("cobravsviper: (snapshot) %s", RawGitDescribe)
		}

		return fmt.Sprintf("cobravsviper: %s", version.String())
	}
}
