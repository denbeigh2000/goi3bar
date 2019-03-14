load("@bazel_tools//tools/build_defs/repo:http.bzl", "http_archive")

http_archive(
    name = "io_bazel_rules_go",
    urls = ["https://github.com/bazelbuild/rules_go/releases/download/0.17.0/rules_go-0.17.0.tar.gz"],
    sha256 = "492c3ac68ed9dcf527a07e6a1b2dcbf199c6bf8b35517951467ac32e421c06c1",
)

http_archive(
    name = "bazel_gazelle",
    urls = ["https://github.com/bazelbuild/bazel-gazelle/releases/download/0.16.0/bazel-gazelle-0.16.0.tar.gz"],
    sha256 = "7949fc6cc17b5b191103e97481cf8889217263acf52e00b560683413af204fcb",
)

load("@io_bazel_rules_go//go:deps.bzl", "go_rules_dependencies", "go_register_toolchains")

go_rules_dependencies()

go_register_toolchains()

load("@bazel_gazelle//:deps.bzl", "gazelle_dependencies", "go_repository")

gazelle_dependencies()

go_repository(
    name = "com_github_shirou_gopsutil",
    commit = "ebc97eefea9b062f9f1624c042c98f196fc90248",
    importpath = "github.com/shirou/gopsutil",
)

go_repository(
    name = "com_github_alyu_configparser",
    commit = "c505e6011694d3c8c1accccea3c9f57eef22afb1",
    importpath = "github.com/alyu/configparser",
)

go_repository(
    name = "com_github_jehiah_go_strftime",
    commit = "1d33003b386959af197ba96475f198c114627b5e",
    importpath = "github.com/jehiah/go-strftime",
)

go_repository(
    name = "org_cloudfoundry_code_bytefmt",
    commit = "2aa6f33b730c79971cfc3c742f279195b0abc627",
    importpath = "code.cloudfoundry.org/bytefmt",
)
