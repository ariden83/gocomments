load("@rules_go//go:def.bzl", "go_library")

go_library(
    name = "graph_debug_info_proto",
    srcs = ["graph_debug_info.pb.go"],
    importpath = "github.com/wamuir/graft/tensorflow/core/framework/graph_debug_info_proto",
    visibility = ["//visibility:public"],
    deps = [
        "@org_golang_google_protobuf//reflect/protoreflect",
        "@org_golang_google_protobuf//runtime/protoimpl",
    ],
)

alias(
    name = "go_default_library",
    actual = ":graph_debug_info_proto",
    visibility = ["//visibility:public"],
)
