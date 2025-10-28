const std = @import("std");

pub fn main() !void {
    const address = try std.net.Address.parseIp("224.0.0.1", 32100);
    const sock = try std.posix.socket(std.posix.AF.INET, std.os.linux.SOCK.DGRAM, std.posix.IPPROTO.UDP);
    errdefer std.posix.close(sock);
    try std.posix.connect(sock, &address.any, address.getOsSockLen());
    const message: []const u8 = "ready";
    const send_bytes = try std.posix.send(sock, message, 0);
    std.debug.print("{d}\n", .{send_bytes});
    // nc -l -u -s 224.0.0.1 -p 32100
}
