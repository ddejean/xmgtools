# Copyright (C) 2025 - Damien Dejean <dam.dejean@gmail.com>

load("@gazelle//:def.bzl", "gazelle")
load("@rules_go//go:def.bzl", "go_binary", "go_library")

# Let gazelle generate build rules.
gazelle(name = "gazelle")

go_library(
    name = "swctl_lib",
    srcs = [
        "plug.go",
        "swctl.go",
    ],
    importpath = "xioxoz.fr/swctl",
    visibility = ["//visibility:private"],
    deps = [
        "//swctl/bootext",
        "//swctl/uboot",
        "@com_github_tarm_serial//:serial",
    ],
)

go_binary(
    name = "swctl",
    embed = [":swctl_lib"],
    visibility = ["//visibility:public"],
)
