const std = @import("std");
const print = std.debug.print;
const net_address = std.Io.net.IpAddress;
const std_thread = std.Io.Threaded;
const panic = std.debug.panic;

var net_server: std.Io.net.Server = undefined;
var net_stream: std.Io.net.Stream = undefined;
var group: std.Io.Group = .init;

pub const tcp_server = struct {
    const Self = @This();
    address_value: []const u8 = "0.0.0.0",
    port: u16 = 9999,
    std_io: std.Io = undefined,
    server_options: std.Io.net.IpAddress.ListenOptions = undefined,
    pub fn init(self: Self, std_io: std.Io, listen_options: std.Io.net.IpAddress.ListenOptions) tcp_server {
        _ = self;
        return tcp_server{ .server_options = listen_options, .std_io = std_io };
    }
};

pub fn start_server(server: tcp_server, io: std.Io, opts: std.Io.net.IpAddress.ListenOptions) !void {
    const address = try std.Io.net.IpAddress.parse(server.address_value, server.port);
    var stream = try net_address.listen(address, io, opts);
    var client: std.Io.net.Stream = undefined;
    defer {
        client.close(io);
        stream.deinit(io);
    }
    std.debug.print("starting server...\n", .{});
    while (true) {
        client = try stream.accept(io);
        try group.concurrent(io, stream_data, .{ client, io });
    }
}

fn stream_data(client: std.Io.net.Stream, io: std.Io) void {
    std.debug.print("reading data from server\n", .{});
    var buf: [127]u8 = undefined;
    var client_reader = client.reader(io, &buf);
    const data = client_reader.interface.takeDelimiterInclusive('\n') catch unreachable;
    if (data.len > 0) {
        std.debug.print("data:{s}\n", .{data});
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
    const connect_opts = std.Io.net.IpAddress.ConnectOptions{
        .mode = .stream,
        .protocol = .tcp,
    };
    const nt_server: tcp_server = .{};
    var tcp_stream_server = tcp_server.init(nt_server, std.testing.io, opts);
    try std.testing.expect(tcp_stream_server.address_value.len > 0);
    try std.testing.expect(tcp_stream_server.port > 0);
    const address = try net_address.parse(tcp_stream_server.address_value, tcp_stream_server.port);
    _ = std.Io.async(std.testing.io, start_server, .{ tcp_stream_server, std.testing.io, opts });
    try std.Io.sleep(std.testing.io, std.Io.Duration{ .nanoseconds = 100000 }, std.Io.Clock.awake);
    var client_async = std.Io.async(std.testing.io, write_data_to_stream, .{ address, std.testing.io, connect_opts });
    try client_async.await(std.testing.io);
}

fn write_data_to_stream(address: std.Io.net.IpAddress, io: std.Io, opts: std.Io.net.IpAddress.ConnectOptions) !void {
    std.debug.print("starting client\n", .{});
    const client_stream = address.connect(io, opts) catch |e| {
        std.debug.print("error connecting to server:{s}\n", .{@errorName(e)});
        return;
    };
    defer client_stream.close(io);
    var buffer: [124]u8 = undefined;
    var client_writer = client_stream.writer(io, &buffer);
    try client_writer.interface.writeAll("some data from client\n");
    try client_writer.interface.flush();
    std.debug.print("sent data...\n", .{});
}
