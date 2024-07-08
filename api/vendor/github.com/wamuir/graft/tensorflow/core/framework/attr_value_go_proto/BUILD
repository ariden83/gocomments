load("@rules_go//go:def.bzl", "go_library")

go_library(
    name = "attr_value_go_proto",
    srcs = ["attr_value.pb.go"],
    importpath = "github.com/wamuir/graft/tensorflow/core/framework/attr_value_go_proto",
    visibility = ["//visibility:public"],
    deps = [
        "//tensorflow/core/framework/tensor_go_proto",
        "//tensorflow/core/framework/tensor_shape_go_proto",
        "//tensorflow/core/framework/types_go_proto",
        "@org_golang_google_protobuf//reflect/protoreflect",
        "@org_golang_google_protobuf//runtime/protoimpl",
    ],
)

alias(
    name = "go_default_library",
    actual = ":attr_value_go_proto",
    visibility = ["//visibility:public"],
)
