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
        print("error reading config from ziggy:{any}\n", .{err});
        return;
    };
    const address = try std.net.Address.parseIp("224.0.0.1", 32100);
    const sock = try std.posix.socket(std.posix.AF.INET, std.os.linux.SOCK.DGRAM, std.posix.IPPROTO.UDP);
    errdefer std.posix.close(sock);
    try std.posix.connect(sock, &address.any, address.getOsSockLen());
    const message: []const u8 = "9999";
    // nc -l -u -s 224.0.0.1 -p 32100
    const heartbeat_config = try heartbeat.split_interval(config_value);
    print("time_increment_interval:{d}\n", .{heartbeat_config.time_increment_interval});
    print("increment_unit:{s}\n", .{heartbeat_config.increment_unit});
    while (true) {
        const bytes = try std.posix.send(sock, message, 0);
        print("sent heart_beat with size:{d}\tmessage with data:{s}\n", .{ bytes, message });
        std.posix.nanosleep(heartbeat_config.time_increment_interval, std.time.ms_per_s);
    }
}
