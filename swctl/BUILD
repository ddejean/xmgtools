# Copyright (C) 2025 - Damien Dejean <dam.dejean@gmail.com>

load("@gazelle//:def.bzl", "gazelle")
load("@rules_go//go:def.bzl", "go_binary", "go_library")

# Let gazelle generate build rules.
gazelle(name = "gazelle")

go_library(
    name = "swctl_lib",
    srcs = [
        "automator.go",
        "csm.go",
        "plug.go",
        "scanner.go",
        "swctl.go",
    ],
    importpath = "xioxoz.fr/swctl",
    visibility = ["//visibility:private"],
    deps = [
        "@com_github_azurity_xmodem_go//:xmodem-go",
        "@com_github_machinebox_progress//:progress",
        "@com_github_tarm_serial//:serial",
    ],
)

go_binary(
    name = "swctl",
    embed = [":swctl_lib"],
    visibility = ["//visibility:public"],
)
