load("@rules_go//go:def.bzl", "go_library")

go_library(
    name = "op_def_go_proto",
    srcs = ["op_def.pb.go"],
    importpath = "github.com/wamuir/graft/tensorflow/core/framework/op_def_go_proto",
    visibility = ["//visibility:public"],
    deps = [
        "//tensorflow/core/framework/attr_value_go_proto",
        "//tensorflow/core/framework/full_type_go_proto",
        "//tensorflow/core/framework/resource_handle_go_proto",
        "//tensorflow/core/framework/types_go_proto",
        "@org_golang_google_protobuf//reflect/protoreflect",
        "@org_golang_google_protobuf//runtime/protoimpl",
    ],
)

alias(
    name = "go_default_library",
    actual = ":op_def_go_proto",
    visibility = ["//visibility:public"],
)
