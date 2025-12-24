const std = @import("std");
const print = std.debug.print;
const net_address = std.Io.net.IpAddress;
const std_thread = std.Io.Threaded;

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

test "init" {
    var arena = std.heap.ArenaAllocator.init(std.heap.page_allocator);
    defer arena.deinit();
    const opts: std.Io.net.IpAddress.BindOptions = .{ .ip6_only = false, .protocol = std.Io.net.Protocol.udp, .mode = std.Io.net.Socket.Mode.dgram };
    const allocator = arena.allocator();
    var thread = std_thread.init(allocator);
    defer thread.deinit();
    const udp_listen_server: udp_server = .{};
    try init(udp_listen_server, thread.io(), opts);
}
