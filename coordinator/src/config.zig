const ziggy = @import("ziggy");
const std = @import("std");
const read = @import("./read.zig");
const print = std.debug.print;

pub const Config = struct {
    heartbeat: []const u8,
};

pub fn read_config(io: std.Io, file_name: []const u8, allocator: std.mem.Allocator) !Config {
    const data = try read.read_file(io, file_name, allocator);
    var meta: ziggy.Deserializer.Meta = .init;
    const file_data: [:0]u8 = @ptrCast(data);
    print("data from file: {s}\n", .{file_data});
    // TODO: fix leaky unexpected error
    const cfg = try ziggy.deserializeLeaky(Config, allocator, file_data, &meta, .{});
    return cfg;
}

// ============================================================
// ======================= unit test ==========================
// ============================================================

test "read_config" {
    var arena = std.heap.ArenaAllocator.init(std.heap.page_allocator);
    defer arena.deinit();
    var threaded_io = std.Io.Threaded.init(arena.allocator());
    defer threaded_io.deinit();
    const result = try read_config(threaded_io.ioBasic(), "config.zgy", arena.allocator());
    print("config:{s}\n", .{result.heartbeat});
}
