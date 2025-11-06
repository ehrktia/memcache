const std = @import("std");
const print = std.debug.print;
const config = @import("./config.zig");
const heartbeat = @import("./heartbeat.zig");

pub fn main() !void {
    var arena = std.heap.ArenaAllocator.init(std.heap.page_allocator);
    const allocator = arena.allocator();
    const file_location: []const u8 = "config.zgy";
    try heartbeat.initialize_time_lookup_store(allocator);
    const config_value = config.read_config(file_location, allocator) catch |err| {
        std.debug.print("error reading config from ziggy:{any}\n", .{err});
        return;
    };
    std.debug.print("{s}\n", .{config_value.heartbeat});
    const address = try std.net.Address.parseIp("224.0.0.1", 32100);
    const sock = try std.posix.socket(std.posix.AF.INET, std.os.linux.SOCK.DGRAM, std.posix.IPPROTO.UDP);
    errdefer std.posix.close(sock);
    try std.posix.connect(sock, &address.any, address.getOsSockLen());
    const message: []const u8 = "9999";
    // TODO: get heartbeat.interval struct use start time and time_increment
    // to run a infinite loop and emit one beat with sleep set to time_increment
    var count: i8 = 10;
    while (count > 0) : (count -= 1) {
        const send_bytes = try std.posix.send(sock, message, 0);
        std.debug.print("{d}\n", .{send_bytes});
    }
    // nc -l -u -s 224.0.0.1 -p 32100
    try heartbeat.split_interval(config_value);
}
