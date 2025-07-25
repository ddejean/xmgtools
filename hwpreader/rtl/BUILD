# Copyright (C) 2025 - Damien Dejean <dam.dejean@gmail.com>
load("@rules_go//go:def.bzl", "go_binary", "go_library")

go_library(
    name = "rtl_lib",
    srcs = [
        "consts.go",
        "chipid.go",
        "leds.go",
        "phy.go",
        "ports.go",
        "serdes.go",
        "switch.go",
    ],
    importpath = "xioxoz.fr/hwpreader/rtl",
    visibility = ["//hwpreader:__pkg__"],
)