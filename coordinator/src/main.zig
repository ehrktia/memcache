const std = @import("std");
const print = std.debug.print;
// const config = @import("./config.zig");
const heartbeat = @import("./heartbeat.zig");
const tcp = @import("./tcp.zig");
const udp = @import("./udp.zig");
const std_thread = std.Io.Threaded;

pub fn main() !void {
    var arena_allocator = std.heap.ArenaAllocator.init(std.heap.page_allocator);
    defer arena_allocator.deinit();
    var thread = std_thread.init(arena_allocator.allocator());
    var thread_io = thread.io();
    defer thread.deinit();
    const tcp_stream_server: tcp.tcp_server = .{};
    const tcp_opts = std.Io.net.IpAddress.ListenOptions{
        .reuse_address = true,
        .mode = .stream,
        .protocol = .tcp,
    };
    const tcp_stream = try tcp_stream_server.init(&thread_io, tcp_opts);
    const udp_listen_server: udp.udp_server = .{};
    const udp_opts: std.Io.net.IpAddress.BindOptions = .{ .ip6_only = false, .protocol = std.Io.net.Protocol.udp, .mode = std.Io.net.Socket.Mode.dgram };
    while (true) {
        const tcp_thread = try std.Thread.spawn(.{}, tcp.stream_data, .{tcp_stream});
        const udp_thread = try std.Thread.spawn(.{}, udp.init, .{ udp_listen_server, thread_io, udp_opts });
        tcp_thread.join();
        udp_thread.join();
    }
}
