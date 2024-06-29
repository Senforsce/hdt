const std = @import("std");
const Builder = std.build.Builder;
const builtin = std.builtin;

pub fn build(b: *std.build.Builder) void {
    // link C++ hdtwrapper bridge to zig
    const target = b.standardTargetOptions(.{});
    const optimize = b.standardOptimizeOption(.{});

    const exe = b.addExecutable(.{
        .name = "zighdt",
        .root_source_file = .{ .path = "./main.zig" },
        .target = target,
        .optimize = optimize,
    });

    exe.linkLibC();
    exe.linkLibCpp();
    exe.addIncludePath("./");
    exe.addIncludePath("./libhdt");
    exe.addCSourceFile("./hdtwrapper.cpp", &.{});

    b.installArtifact(exe);
    // end link C++ hdtwrapper bridge to zig

    // link golang bridge to zig

    const mode = b.standardReleaseOptions();
    const lib = b.addStaticLibrary("zgo", "zgo.zig");
    lib.bundle_compiler_rt = true;
    lib.use_stage1 = false;
    lib.emit_h = true;
    lib.emit_bin = .{ .emit_to = "libzgo.a" };
    lib.setBuildMode(mode);
    lib.install();

    const go = build_go(b);
    const make_step = b.step("go", "Make go executable");
    make_step.dependOn(&go.step);
}

fn build_go(b: *std.build.Builder) *std.build.RunStep {
    const go = b.addSystemCommand(
        &[_][]const u8{
            "go",
            "build",
            "-ldflags",
            "-linkmode external -extldflags -static",
            "bridge.go",
        },
    );
    return go;
}
