const std = @import("std");
const print = std.debug.print;
const net_address = std.Io.net.IpAddress;
const std_thread = std.Io.Threaded;

var net_server: std.Io.net.Server = undefined;
var net_stream: std.Io.net.Stream = undefined;

pub const tcp_server = struct {
    const Self = @This();
    address_value: []const u8 = "::",
    port: u16 = 9999,
    std_io: std.Io = undefined,
    server_options: std.Io.net.IpAddress.ListenOptions = undefined,
    pub fn init(self: Self, std_io: *std.Io, listen_options: std.Io.net.IpAddress.ListenOptions) !tcp_server {
        const server_address = try net_address.parse(self.address_value, self.port);
        net_server = try net_address.listen(server_address, std_io.*, listen_options);
        net_stream = net_server.accept(std_io.*) catch |e| {
            std.debug.panic("error accepting server:{any}\n", .{e});
            return;
        };
        return tcp_server{ .server_options = listen_options, .std_io = std_io.* };
    }
};
var buffer: [1096]u8 = undefined;
pub fn stream_data(self: tcp_server) void {
    var stream_reader = net_stream.reader(self.std_io, &buffer);
    var reader = &stream_reader.interface;
    const data = reader.buffered();
    print("data received:{s}\n", .{data});
}

// =======================================================
// ================== unit test ==========================
// =======================================================

test "init" {
    const opts = std.Io.net.IpAddress.ListenOptions{
        .reuse_address = true,
        .mode = .stream,
        .protocol = .tcp,
    };
    var arena = std.heap.ArenaAllocator.init(std.heap.page_allocator);
    defer arena.deinit();
    const allocator = arena.allocator();
    var thread = std_thread.init(allocator);
    defer thread.deinit();
    var io_thread = thread.io();
    const nt_server: tcp_server = .{};
    const tcp_stream_server = try tcp_server.init(nt_server, &io_thread, opts);
    try std.testing.expect(tcp_stream_server.address_value.len > 0);
    try std.testing.expect(tcp_stream_server.port > 0);
    print("address-{s}\n", .{tcp_stream_server.address_value});
    const tcp_server_thread = try std.Thread.spawn(.{}, stream_data, .{tcp_stream_server});
    tcp_server_thread.join();
}
