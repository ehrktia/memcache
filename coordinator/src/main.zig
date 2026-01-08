const std = @import("std");
const print = std.debug.print;
const tcp = @import("./tcp.zig");
const config = @import("./config.zig");
const std_thread = std.Io.Threaded;
var heart_beat: u64 = undefined;
pub fn main() !void {
    var arena_allocator = std.heap.ArenaAllocator.init(std.heap.page_allocator);
    defer arena_allocator.deinit();
    var arena = std.heap.ArenaAllocator.init(std.heap.page_allocator);
    var threaded = std.Io.Threaded.init(arena.allocator(), .{});
    defer threaded.deinit();
    heart_beat = try config.read_config_from_file(threaded.io(), arena_allocator.allocator());
    print("heart beat:{d}\n", .{heart_beat});
    const tcp_stream_server: tcp.tcp_server = .{};
    const tcp_opts = std.Io.net.IpAddress.ListenOptions{
        .reuse_address = true,
        .mode = .stream,
        .protocol = .tcp,
    };
    const tcp_stream = tcp_stream_server.init(threaded.io(), tcp_opts);
    var udp_address = try std.Io.net.IpAddress.parse("224.0.0.1", 32100);
    const udp_opts: std.Io.net.IpAddress.BindOptions = .{ .ip6_only = false, .protocol = std.Io.net.Protocol.udp, .mode = std.Io.net.Socket.Mode.dgram };
    const udp_socket = try udp_address.bind(threaded.io(), udp_opts);
    defer udp_socket.close(threaded.io());
    _ = std.Io.async(threaded.io(), tcp.start_server, .{ tcp_stream, threaded.io(), tcp_opts });
    var udp_heart_beat = std.Io.async(threaded.io(), send_heart_beat, .{ threaded.io(), udp_socket, &udp_address, "9999" });
    try udp_heart_beat.await(threaded.io());
}

fn convert_to_sec() void {
    heart_beat = heart_beat * 1000000000;
}

fn send_heart_beat(io: std.Io, socket: std.Io.net.Socket, address: *std.Io.net.IpAddress, message: []const u8) !void {
    while (true) {
        try std.Io.sleep(io, std.Io.Duration{ .nanoseconds = heart_beat * 1000000000 }, std.Io.Clock.awake);
        try socket.send(io, address, message);
    }
}

// ======================================================================
// ============================== test unit =============================
// ======================================================================

test "convert" {
    heart_beat = 9;
    convert_to_sec();
    std.debug.print("heart_beat:{d}\n", .{heart_beat});
}
