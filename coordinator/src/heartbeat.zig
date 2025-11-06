const std = @import("std");
const print = std.debug.print;
const config = @import("./config.zig");
var unit_idx_start: usize = undefined;
var time_lookup_store: std.hash_map.HashMap([]const u8, i128, std.hash_map.StringContext, 80) = undefined;
pub fn split_interval(cfg: config.Config) !void {
    // TODO: infinite loop with increment based on the time_increment

    for (cfg.heartbeat, 0..) |value, i| {
        if (value >= '0' and value <= '9') {
            continue;
        }
        unit_idx_start = i;
        break;
    }
    const time_unit_value = cfg.heartbeat[unit_idx_start..];
    const time_unit = time_lookup_store.get(time_unit_value).?;
    print("start time:{d}\n", .{time_unit});
    print("time increment:{s}\n", .{cfg.heartbeat[0..unit_idx_start]});
    const time_increment = std.fmt.parseInt(i32, cfg.heartbeat[0..unit_idx_start], 10) catch |err| {
        print("parse int from str error:{any}\n", .{err});
        return;
    };
    print("end time:{d}\n", .{std.time.timestamp() + time_increment});
}

pub fn initialize_time_lookup_store(allocator: std.mem.Allocator) !void {
    time_lookup_store = std.StringHashMap(i128).init(allocator);
    try time_lookup_store.put("ms", std.time.milliTimestamp());
    try time_lookup_store.put("mis", std.time.microTimestamp());
    try time_lookup_store.put("ns", std.time.nanoTimestamp());
    try time_lookup_store.put("s", std.time.timestamp());
}
