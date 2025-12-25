// const ziggy = @import("ziggy");
const std = @import("std");
const read = @import("./read.zig");
const print = std.debug.print;

pub const Config = struct {
    heart_beat: []const u8,
};

pub fn read_config(io: std.Io, file_name: []const u8, allocator: std.mem.Allocator) !Config {
    const data = try read.read_file(io, file_name, allocator);
    const file_data: [:0]u8 = @ptrCast(data);
    const cfg: Config = try std.json.parseFromSliceLeaky(Config, allocator, file_data, .{});
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
    print("config:{s}\n", .{result.heart_beat});
}
