load("@rules_go//go:def.bzl", "go_library")

go_library(
    name = "step_stats_go_proto",
    srcs = ["step_stats.pb.go"],
    importpath = "github.com/wamuir/graft/tensorflow/core/framework/step_stats_go_proto",
    visibility = ["//visibility:public"],
    deps = [
        "//tensorflow/core/framework/allocation_description_go_proto",
        "//tensorflow/core/framework/tensor_description_go_proto",
        "@org_golang_google_protobuf//reflect/protoreflect",
        "@org_golang_google_protobuf//runtime/protoimpl",
    ],
)

alias(
    name = "go_default_library",
    actual = ":step_stats_go_proto",
    visibility = ["//visibility:public"],
)
