/*
 * SPDX-License-Identifier: GPL-3.0
 * TestCord Installer, a cross platform gui/cli app for installing TestCord
 * Copyright (c) 2025 TestcordDev and TestCord contributors
 */

package main

import (
	"equilotl/buildinfo"
	"image/color"
)

const ReleaseUrl = "https://api.github.com/repos/TestcordDev/Testcord/releases/tags/latest"
const ReleaseUrlFallback = "https://github.com/TestcordDev/Testcord/releases/tag/latest"
const InstallerReleaseUrl = "https://api.github.com/repos/TestcordDev/Testcordinstaller/releases/latest"
const InstallerReleaseUrlFallback = "https://testcord.org/releases/xcinstaller"

// ReleaseRepoApi / InstallerRepoApi are the GitHub API roots of the mod and the
// installer, used to resolve release tags to the commit they point at.
const ReleaseRepoApi = "https://api.github.com/repos/TestcordDev/Testcord"
const InstallerRepoApi = "https://api.github.com/repos/TestcordDev/Testcordinstaller"

var UserAgent = "TestCordInstaller/" + buildinfo.InstallerGitHash + " (https://github.com/TestcordDev/Testcordinstaller)"

const SupportUrl = "https://github.com/TestcordDev/Testcordinstaller/issues"

var (
	DiscordGreen        = color.RGBA{R: 0x2D, G: 0x7C, B: 0x46, A: 0xFF}
	DiscordGreenHovered = color.RGBA{R: 0x25, G: 0x64, B: 0x39, A: 0xFF}
	DiscordRed          = color.RGBA{R: 0xEC, G: 0x41, B: 0x44, A: 0xFF}
	DiscordRedHovered   = color.RGBA{R: 0xBE, G: 0x34, B: 0x37, A: 0xFF}
	DiscordBlue         = color.RGBA{R: 0x58, G: 0x65, B: 0xF2, A: 0xFF}
	DiscordBlueHovered  = color.RGBA{R: 0x45, G: 0x4F, B: 0xBD, A: 0xFF}
	DiscordYellow       = color.RGBA{R: 0xfe, G: 0xe7, B: 0x5c, A: 0xff}
)

var LinuxDiscordNames = []string{
	"Discord",
	"DiscordPTB",
	"DiscordCanary",
	"DiscordDevelopment",
	"discord",
	"discordptb",
	"discordcanary",
	"discorddevelopment",
	"discord-ptb",
	"discord-canary",
	"discord-development",
	// Flatpak
	"com.discordapp.Discord",
	"com.discordapp.DiscordPTB",
	"com.discordapp.DiscordCanary",
	"com.discordapp.DiscordDevelopment",
}
