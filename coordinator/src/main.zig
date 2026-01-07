const std = @import("std");
const print = std.debug.print;
const tcp = @import("./tcp.zig");
const config = @import("./config.zig");
// const udp = @import("./udp.zig");
const std_thread = std.Io.Threaded;

pub fn main() !void {
    var arena_allocator = std.heap.ArenaAllocator.init(std.heap.page_allocator);
    defer arena_allocator.deinit();
    var thread = std_thread.init(arena_allocator.allocator());
    var thread_io = thread.io();
    defer thread.deinit();
    // config file heartbeat interval
    const h_beat = try config.read_config_from_file(thread_io, arena_allocator.allocator());
    print("heart beat:{d}\n", .{h_beat});
    const tcp_stream_server: tcp.tcp_server = .{};
    const tcp_opts = std.Io.net.IpAddress.ListenOptions{
        .reuse_address = true,
        .mode = .stream,
        .protocol = .tcp,
    };
    const tcp_stream = try tcp_stream_server.init(&thread_io, tcp_opts);
    var udp_address = try std.Io.net.IpAddress.parse("224.0.0.1", 32100);
    const udp_opts: std.Io.net.IpAddress.BindOptions = .{ .ip6_only = false, .protocol = std.Io.net.Protocol.udp, .mode = std.Io.net.Socket.Mode.dgram };
    const udp_socket = try udp_address.bind(thread_io, udp_opts);
    defer udp_socket.close(thread_io);
    while (true) {
        _ = std.Io.async(thread_io, tcp.start_server, .{ tcp_stream, thread_io, tcp_opts });
        // const tcp_thread = try std.Thread.spawn(.{}, tcp.stream_data, .{tcp_stream});
        _ = try std.Io.async(thread_io, send_heart_beat, .{ thread_io, udp_socket, &udp_address, "9999" });
        // tcp_thread.join();
        // udp_thread.join();
    }
}

fn send_heart_beat(io: std.Io, socket: std.Io.net.Socket, address: *std.Io.net.IpAddress, message: []const u8) !void {
    print("sending message-{s}\n", .{message});
    var count: usize = 0;
    // TODO: make it infinite loop with sleep timer matching heartbeat
    while (count < 25) : (count += 1) {
        print("count={d}\n", .{count});
        try socket.send(io, address, message);
    }
}
