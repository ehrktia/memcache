const std = @import("std");
const config = @import("./config.zig");
const print = std.debug.print;
var unit_idx_start: usize = undefined;
pub const Heartbeat_Config = struct {
    time_increment_interval: u64,
};

pub fn split_interval(cfg: config.Config) !Heartbeat_Config {
    for (cfg.heart_beat, 0..) |value, i| {
        if (value >= '0' and value <= '9') {
            continue;
        }
        unit_idx_start = i;
        break;
    }
    const time_increment = std.fmt.parseInt(u64, cfg.heart_beat[0..unit_idx_start], 10) catch |err| {
        print("parse int from str error:{any}\n", .{err});
        return err;
    };
    return .{
        .time_increment_interval = time_increment,
    };
}
