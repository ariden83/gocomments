load("@rules_go//go:def.bzl", "go_library")

go_library(
    name = "cost_graph_go_proto",
    srcs = ["cost_graph.pb.go"],
    importpath = "github.com/wamuir/graft/tensorflow/core/framework/cost_graph_go_proto",
    visibility = ["//visibility:public"],
    deps = [
        "//tensorflow/core/framework/tensor_shape_go_proto",
        "//tensorflow/core/framework/types_go_proto",
        "@org_golang_google_protobuf//reflect/protoreflect",
        "@org_golang_google_protobuf//runtime/protoimpl",
    ],
)

alias(
    name = "go_default_library",
    actual = ":cost_graph_go_proto",
    visibility = ["//visibility:public"],
)
