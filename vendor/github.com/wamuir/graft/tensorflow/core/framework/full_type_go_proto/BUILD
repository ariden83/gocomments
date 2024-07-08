load("@rules_go//go:def.bzl", "go_library")

go_library(
    name = "full_type_go_proto",
    srcs = ["full_type.pb.go"],
    importpath = "github.com/wamuir/graft/tensorflow/core/framework/full_type_go_proto",
    visibility = ["//visibility:public"],
    deps = [
        "@org_golang_google_protobuf//reflect/protoreflect",
        "@org_golang_google_protobuf//runtime/protoimpl",
    ],
)

alias(
    name = "go_default_library",
    actual = ":full_type_go_proto",
    visibility = ["//visibility:public"],
)
