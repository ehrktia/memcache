const ziggy = @import("ziggy");
const std = @import("std");
const expect = std.testing.expect;

pub const Config = struct {
    heartbeat: []const u8,
};

pub fn read_config(allocator: std.mem.Allocator) !Config {
    return try ziggy.parseLeaky(Config, allocator, "str", .{});
}
