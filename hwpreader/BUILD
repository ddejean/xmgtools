# Copyright (C) 2025 - Damien Dejean <dam.dejean@gmail.com>

load("@gazelle//:def.bzl", "gazelle")
load("@rules_go//go:def.bzl", "go_binary", "go_library")

# Let gazelle generate build rules.
gazelle(name = "gazelle")

go_library(
    name = "hwpreader_lib",
    srcs = [
        "hwpreader.go",
    ],
    importpath = "xioxoz.fr/hwpreader",
    deps = ["//hwpreader/rtl:rtl_lib"],
    visibility = ["//visibility:private"],
)

go_binary(
    name = "hwpreader",
    embed = [":hwpreader_lib"],
    visibility = ["//visibility:public"],
)
