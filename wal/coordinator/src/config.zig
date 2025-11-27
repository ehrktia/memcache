const ziggy = @import("ziggy");
const std = @import("std");
const print = std.debug.print;

pub const Config = struct {
    heartbeat: []const u8,
};

pub fn read_config(file_name: []const u8, allocator: std.mem.Allocator) !Config {
    const file = std.fs.cwd().openFile(file_name, .{ .mode = .read_only }) catch |err| {
        print("file open error:{any}\n", .{err});
        return err;
    };
    defer file.close();
    const file_stat = (try file.stat()).size;
    const file_size: usize = @as(usize, file_stat);
    const buffer = allocator.alloc(u8, file_size) catch |err| {
        print("buffer allocate error:{any}\n", .{err});
        return err;
    };
    _ = file.read(buffer) catch |err| {
        print("file read error:{any}\n", .{err});
        return err;
    };
    const file_data_str: [:0]u8 = @ptrCast(buffer);
    const cfg = ziggy.parseLeaky(Config, allocator, file_data_str, .{}) catch |err| {
        print("parse ziggy error:{any}\n", .{err});
        return err;
    };
    return cfg;
}
