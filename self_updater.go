/*
 * SPDX-License-Identifier: GPL-3.0
 * TestCord Installer, a cross platform gui/cli app for installing TestCord
 * Copyright (c) 2025 TestcordDev and TestCord contributors
 */

package main

import (
	"equilotl/buildinfo"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"runtime"
	"strings"
	"time"
)

var IsSelfOutdated = false
var SelfUpdateCheckDoneChan = make(chan bool, 1)

func init() {
	//goland:noinspection GoBoolExpressions
	if buildinfo.InstallerTag == buildinfo.VersionUnknown {
		Log.Debug("Disabling self updater as this is not a release build")
		return
	}

	go DeleteOldExecutable()

	go func() {
		Log.Debug("Checking for Installer Updates...")

		res, err := GetGithubRelease(InstallerReleaseUrl, InstallerReleaseUrl)
		if err != nil {
			Log.Warn("Failed to check for self updates:", err)
			SelfUpdateCheckDoneChan <- false
		} else {
			IsSelfOutdated = isInstallerOutdated(res.TagName)
			Log.Debug("Is self outdated?", IsSelfOutdated)
			SelfUpdateCheckDoneChan <- true
		}
	}()
}

// resolveTagToCommit resolves a release tag to the commit it points at.
func resolveTagToCommit(tag string) (string, error) {
	return resolveTagToCommitInRepo(InstallerRepoApi, tag)
}

// isInstallerOutdated reports whether the running binary differs from the
// published one.
//
// We ship releases under a rolling tag ("latest", "rel") that is moved on every
// publish, so the tag alone cannot answer this: comparing the release tag
// against the tag baked in at build time reports an update on every single run,
// and comparing only tags would never report one at all. The git hash compiled
// into the binary is what actually identifies the build, so resolve the release
// tag to a commit and compare that. The tag comparison is only a fallback for
// builds that have no usable hash.
func isInstallerOutdated(releaseTag string) bool {
	local := buildinfo.InstallerGitHash

	if local == "" || local == buildinfo.VersionUnknown {
		Log.Debug("No git hash baked in, falling back to tag comparison")
		return releaseTag != buildinfo.InstallerTag
	}

	commit, err := resolveTagToCommit(releaseTag)
	if err != nil {
		Log.Debug("Failed to resolve release tag to a commit, falling back to tag comparison:", err)
		return releaseTag != buildinfo.InstallerTag
	}

	// the baked in hash is the short (abbreviated) form of the commit
	if strings.HasPrefix(commit, local) || strings.HasPrefix(local, commit) {
		Log.Debug("Installer matches published commit", commit)
		return false
	}

	Log.Debug("Installer commit", local, "differs from published commit", commit)
	return true
}

func GetInstallerDownloadLink() string {
	const BaseUrl = "https://github.com/TestcordDev/Testcordinstaller/releases/latest/download/"
	isCli := buildinfo.UiType == buildinfo.UiTypeCli

	switch runtime.GOOS {
	case "windows":
		// 32bit builds are published without an arch suffix, so only amd64/arm64 exist
		filename := "Windows_Testcord_installer-rel"
		if isCli {
			filename = "Windows_Testcord_installer-rel_cli"
		}
		return BaseUrl + filename + ".exe"
	case "linux":
		if isCli {
			return BaseUrl + "Linux_Testcord_installer-rel_cli"
		}
		return BaseUrl + "Linux_Testcord_installer-rel"
	default:
		// no macos assets are published to this repo
		return ""
	}
}

func CanUpdateSelf() bool {
	//goland:noinspection GoBoolExpressions
	return IsSelfOutdated && runtime.GOOS != "darwin"
}

func UpdateSelf() error {
	if !CanUpdateSelf() {
		return errors.New("Cannot update self. Either no update available or macos")
	}

	url := GetInstallerDownloadLink()
	if url == "" {
		return errors.New("Failed to get installer download link")
	}

	Log.Debug("Updating self from", url)

	ownExePath, err := os.Executable()
	if err != nil {
		return err
	}

	ownExeDir := path.Dir(ownExePath)

	res, err := http.Get(url)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	tmp, err := os.CreateTemp(ownExeDir, "TestcordinstallerUpdate")
	if err != nil {
		return fmt.Errorf("Failed to create tempfile: %w", err)
	}
	defer func() {
		_ = tmp.Close()
		_ = os.Remove(tmp.Name())
	}()
	if err = tmp.Chmod(0o755); err != nil {
		return fmt.Errorf("Failed to chmod 755", tmp.Name()+":", err)
	}

	if _, err = io.Copy(tmp, res.Body); err != nil {
		return err
	}

	if err = tmp.Close(); err != nil {
		return err
	}

	if err = os.Remove(ownExePath); err != nil {
		if err = os.Rename(ownExePath, ownExePath+".old"); err != nil {
			return fmt.Errorf("Failed to remove/rename own executable: %w", err)
		}
	}

	if err = os.Rename(tmp.Name(), ownExePath); err != nil {
		return fmt.Errorf("Failed to replace self with updated executable. Please manually redownload the installer: %w", err)
	}

	return nil
}

func DeleteOldExecutable() {
	ownExePath, err := os.Executable()
	if err != nil {
		return
	}

	for attempts := 0; attempts < 10; attempts += 1 {
		err = os.Remove(ownExePath + ".old")

		if err == nil || errors.Is(err, os.ErrNotExist) {
			break
		}

		Log.Warn("Failed to remove old executable. Retrying in 1 second.", err)
		time.Sleep(1 * time.Second)
	}
}

func RelaunchSelf() error {
	attr := new(os.ProcAttr)
	attr.Files = []*os.File{os.Stdin, os.Stdout, os.Stderr}

	var argv []string
	if len(os.Args) > 1 {
		argv = os.Args[1:]
	} else {
		argv = []string{}
	}

	Log.Debug("Restarting self with exe", os.Args[0], "and args", argv)

	proc, err := os.StartProcess(os.Args[0], argv, attr)
	if err != nil {
		return fmt.Errorf("Failed to start new process: %w", err)
	}

	if err = proc.Release(); err != nil {
		return fmt.Errorf("Failed to release new process: %w", err)
	}

	os.Exit(0)
	return nil
}
