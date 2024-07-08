load("@rules_go//go:def.bzl", "go_library")

go_library(
    name = "function_go_proto",
    srcs = ["function.pb.go"],
    importpath = "github.com/wamuir/graft/tensorflow/core/framework/function_go_proto",
    visibility = ["//visibility:public"],
    deps = [
        "//tensorflow/core/framework/attr_value_go_proto",
        "//tensorflow/core/framework/node_def_go_proto",
        "//tensorflow/core/framework/op_def_go_proto",
        "@org_golang_google_protobuf//reflect/protoreflect",
        "@org_golang_google_protobuf//runtime/protoimpl",
    ],
)

alias(
    name = "go_default_library",
    actual = ":function_go_proto",
    visibility = ["//visibility:public"],
)
