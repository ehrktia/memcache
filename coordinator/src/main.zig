const std = @import("std");
const config = @import("./config.zig");

pub fn main() !void {
    var arena = std.heap.ArenaAllocator.init(std.heap.page_allocator);
    const allocator = arena.allocator();
    // var cfg: config.Config = .{ .heartbeat = "default", .allocator = allocator };
    const config_value = config.read_config(allocator) catch |err| {
        std.debug.print("error reading config from ziggy:{any}\n", .{err});
        return;
    };
    std.debug.print("{s}\n", .{config_value.heartbeat});
    const address = try std.net.Address.parseIp("224.0.0.1", 32100);
    const sock = try std.posix.socket(std.posix.AF.INET, std.os.linux.SOCK.DGRAM, std.posix.IPPROTO.UDP);
    errdefer std.posix.close(sock);
    try std.posix.connect(sock, &address.any, address.getOsSockLen());
    const message: []const u8 = "9999";
    var count: i8 = 10;
    while (count > 0) : (count -= 1) {
        const send_bytes = try std.posix.send(sock, message, 0);
        std.debug.print("{d}\n", .{send_bytes});
    }
    // nc -l -u -s 224.0.0.1 -p 32100
}
