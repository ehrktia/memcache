const std = @import("std");
const print = std.debug.print;
const net_address = std.Io.net.IpAddress;
const std_thread = std.Io.Threaded;
const panic = std.debug.panic;

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
        net_stream = try net_server.accept(std_io.*);
        return tcp_server{ .server_options = listen_options, .std_io = std_io.* };
    }
};
pub fn stream_data(server: tcp_server) !void {
    print("ready to receive data....\n", .{});
    {
        const handle_thread = try std.Thread.spawn(.{}, read_data_from_stream, .{server});
        defer handle_thread.join();
    }
}

fn read_data_from_stream(server: tcp_server) !void {
    var buffer: [1096]u8 = undefined;
    var reader = net_stream.reader(server.std_io, &buffer);
    var count: usize = 0;
    while (count <= 5) : (count += 1) {
        // while (true) {
        const read_size = reader.interface.takeDelimiterInclusive('\n') catch |e| {
            if (e == error.EndOfStream) {
                print("client closed connection \n", .{});
                continue;
            }
            return e;
        };
        print("read:{s}\n", .{read_size});
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
    var tcp_stream_server = try tcp_server.init(nt_server, &io_thread, opts);
    try std.testing.expect(tcp_stream_server.address_value.len > 0);
    try std.testing.expect(tcp_stream_server.port > 0);
    const address = try net_address.parse(tcp_stream_server.address_value, tcp_stream_server.port);
    {
        const tcp_server_thread = try std.Thread.spawn(.{}, stream_data, .{tcp_stream_server});
        defer tcp_server_thread.join();
        const write_data = try std.Thread.spawn(.{}, write_data_to_stream, .{ io_thread, address });
        defer write_data.join();
    }
}

fn write_data_to_stream(io: std.Io, address: std.Io.net.IpAddress) !void {
    const opts = std.Io.net.IpAddress.ConnectOptions{
        .mode = .stream,
        .protocol = .tcp,
        .timeout = .none,
    };
    const msg = try std.Io.net.IpAddress.connect(address, io, opts);
    defer msg.close(io);
    var buf: [100]u8 = undefined;
    var writer = msg.writer(io, &buf);
    try writer.interface.writeAll("data from client");
    try writer.interface.flush();
}
