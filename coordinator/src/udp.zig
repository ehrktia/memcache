const std = @import("std");
const print = std.debug.print;
const panic = std.debug.panic;
const config = @import("./config.zig");
const net_address = std.Io.net.IpAddress;
const std_thread = std.Io.Threaded;

const message = "9999";
pub fn read_config_from_file(io: std.Io, allocator: std.mem.Allocator) !u64 {
    const file_location = "config.zgy";
    try config.read_zgy(io, file_location, allocator);
    // TODO: make the udp emit message for this interval using Io.Sleep
    return config.heart_beat;
}

var net_server: std.Io.net.Server = undefined;
pub const udp_server = struct {
    const Self = @This();
    address: []const u8 = "224.0.0.1",
    port: u16 = 3210,
    server_options: std.Io.net.IpAddress.ListenOptions = undefined,
};

pub fn init(self: udp_server, std_io: std.Io, opts: std.Io.net.IpAddress.BindOptions) !void {
    const server_address = try net_address.parse(self.address, self.port);
    const socket = try server_address.bind(std_io, opts);
    var udp_buffer: [1096]u8 = undefined;
    const incoming_message = try socket.receive(std_io, &udp_buffer);
    print("{s}\n", .{incoming_message.data});
}

// =============================================================
// =================== unit test ===============================
// =============================================================

var arena = std.heap.ArenaAllocator.init(std.heap.page_allocator);
test "init" {
    var thread: std.Io.Threaded = std_thread.init(arena.allocator());
    defer arena.deinit();
    defer thread.deinit();
    const opts: std.Io.net.IpAddress.BindOptions = .{ .ip6_only = false, .protocol = std.Io.net.Protocol.udp, .mode = std.Io.net.Socket.Mode.dgram };
    const udp_listen_server: udp_server = .{};
    try init(udp_listen_server, thread.io(), opts);
}

test "read_config" {
    var thread: std.Io.Threaded = std_thread.init(arena.allocator());
    defer arena.deinit();
    defer thread.deinit();
    const hbeat = try read_config_from_file(thread.ioBasic(), arena.allocator());
    try std.testing.expect(hbeat.time_increment_interval > 0);
}
