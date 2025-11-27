const std = @import("std");
const print = std.debug.print;
const config = @import("./config.zig");
const heartbeat = @import("./heartbeat.zig");
const tcp = @import("./tcp.zig");

pub fn main() !void {
    // var arena = std.heap.ArenaAllocator.init(std.heap.page_allocator);
    // const allocator = arena.allocator();
    // const file_location: []const u8 = "config.zgy";
    // const config_value = config.read_config(file_location, allocator) catch |err| {
    //     print("error reading config from ziggy:{any}\n", .{err});
    //     return;
    // };
    // TODO: add an option to control the udp via config
    // const address = try std.net.Address.parseIp("224.0.0.1", 32100);
    // const sock = try std.posix.socket(std.posix.AF.INET, std.os.linux.SOCK.DGRAM, std.posix.IPPROTO.UDP);
    // errdefer std.posix.close(sock);
    // std.posix.connect(sock, &address.any, address.getOsSockLen()) catch |e| {
    //     print("error connecting to udp network:{any}\n", .{e});
    //     return;
    // };
    // const message: []const u8 = "9999";
    // // nc -l -u -s 224.0.0.1 -p 32100
    // const heartbeat_config = heartbeat.split_interval(config_value) catch |e| {
    //     print("error getting heart beat interval from config:{any}\n", .{e});
    //     return;
    // };
    // print("time increment interval:{d}\n", .{heartbeat_config.time_increment_interval});
    // tcp server to communicate between data layer and control layer
    try tcp.tcp_server();
    // while (true) {
    //     const bytes = try std.posix.send(sock, message, 0);
    //     print("sent heart_beat with size:{d}\tmessage with data:{s}\n", .{ bytes, message });
    //     std.posix.nanosleep(heartbeat_config.time_increment_interval, 0);
    // }
}
