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
        // net_stream = net_server.accept(std_io.*) catch |e| {
        //     std.debug.panic("error accepting server:{any}\n", .{e});
        //     return;
        // };
        return tcp_server{ .server_options = listen_options, .std_io = std_io.* };
    }
};
var buffer: [1096]u8 = undefined;
pub fn stream_data(server: tcp_server) void {
    print("ready to receive data....\n", .{});
    var stream_reader = net_stream.reader(server.std_io, &buffer);
    var reader = &stream_reader.interface;
    const data = reader.buffered();
    if (data.len > 0) {
        print("data received:{s}\n", .{data});
    }
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
    const address = try std.Io.net.IpAddress.parse(tcp_stream_server.address_value, tcp_stream_server.port);
    const tcp_server_thread = try std.Thread.spawn(.{}, stream_data, .{tcp_stream_server});
    const tcp_read = try std.Thread.spawn(.{}, read_data, .{ io_thread, address });
    tcp_server_thread.join();
    tcp_read.join();
}

fn read_data(io: std.Io, address: std.Io.net.IpAddress) !void {
    const opts = std.Io.net.IpAddress.ConnectOptions{
        .mode = .stream,
        .protocol = .tcp,
        .timeout = .none,
    };
    const msg = try std.Io.net.IpAddress.connect(address, io, opts);
    defer msg.close(io);
    const file = std.Io.File.stdin();
    var buf: [1096]u8 = undefined;
    buf[0] = "some test data";
    var data_to_send: [][]u8 = buf;
    std.Io.File.writeStreaming(file, io, &data_to_send);
    var fs_reader = std.Io.File.reader(file, io, &buf);
    const reader = &fs_reader.interface;
    const data = reader.buffered();
    print("data:{s}\n", .{data});
}
