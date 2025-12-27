const ziggy = @import("ziggy");
const std = @import("std");
const read = @import("./read.zig");
const print = std.debug.print;
pub var heart_beat: u64 = undefined;

pub const Config = struct {
    heart_beat: []const u8,
};

// reads config file in zgy format
// populate the heart beat interval
pub fn read_zgy(io: std.Io, file_name: []const u8, allocator: std.mem.Allocator) !void {
    const config_dir = std.Io.Dir.cwd();
    const file = try config_dir.openFile(io, file_name, .{});
    const stat = try file.stat(io);
    var buf = try allocator.alloc(u8, stat.size + 1);
    const data = try config_dir.readFile(io, file_name, buf);
    buf[buf.len - 1] = 0;
    const file_data: [:0]u8 = @ptrCast(data);
    var meta: ziggy.Deserializer.Meta = .init;
    const cfg = try ziggy.deserializeLeaky(Config, allocator, file_data, &meta, .{});
    // populate heart_beat value
    try split_interval(cfg.heart_beat);
}

// read json config as back when zgy is not working
pub fn read_config(io: std.Io, file_name: []const u8, allocator: std.mem.Allocator) !Config {
    const data = try read.read_file(io, file_name, allocator);
    const file_data: [:0]u8 = @ptrCast(data);
    const cfg: Config = try std.json.parseFromSliceLeaky(Config, allocator, file_data, .{});
    // populate heart_beat value
    try split_interval(cfg.heart_beat);
    return cfg;
}

pub fn read_config_from_file(io: std.Io, allocator: std.mem.Allocator) !u64 {
    const file_location = "config.zgy";
    try read_zgy(io, file_location, allocator);
    // TODO: make the udp emit message for this interval using Io.Sleep
    return heart_beat;
}

var unit_idx_start: usize = undefined;
fn split_interval(heartbeat: []const u8) !void {
    for (heartbeat, 0..) |value, i| {
        if (value >= '0' and value <= '9') {
            continue;
        }
        unit_idx_start = i;
        break;
    }
    heart_beat = try std.fmt.parseInt(u64, heartbeat[0..unit_idx_start], 10);
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

test "create_zgy" {
    var arena = std.heap.ArenaAllocator.init(std.heap.page_allocator);
    defer arena.deinit();
    var threaded_io = std.Io.Threaded.init(arena.allocator());
    defer threaded_io.deinit();
    const result = try read_zgy(threaded_io.ioBasic(), "config.zgy", arena.allocator());
    print("result:{s}\n", .{result.heart_beat});
}
