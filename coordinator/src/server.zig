const std = @import("std");
const print = std.debug.print;
const net_address = std.Io.net.IpAddress;
const std_id = std.Io;
const std_thread = std.Io.Threaded;

pub const server = struct {
    const Self = @This();
    address_value: []const u8 = "::",
    port: u16 = 9999,
    udp_port: i64 = 321000,
    std_io: std.Io = undefined,
    server_options: std.Io.net.IpAddress.ListenOptions = undefined,
    var net_server: std.Io.net.Server = undefined;
    var net_stream: std.Io.net.Stream = undefined;
    pub fn init(self: Self, std_io: *std.Io, listen_options: std.Io.net.IpAddress.ListenOptions) !server {
        const server_address = try net_address.parse(self.address_value, self.port);
        net_server = try net_address.listen(server_address, std_io.*, listen_options);
        net_stream = net_server.accept(std_io.*) catch |e| {
            std.debug.panic("error accepting server:{any}\n", .{e});
            return;
        };
        return server{ .server_options = listen_options, .std_io = std_io.* };
    }
    var buffer: [1096]u8 = undefined;
    pub fn stream_data(self: Self) void {
        var stream_reader = net_stream.reader(self.std_io, &buffer);
        var reader = &stream_reader.interface;
        const data = reader.buffered();
        print("data received:{s}\n", .{data});
    }
    pub fn udp_init(self: Self, address: []const u8, port: i64, std_io: *std.Io, listen_options: std.Io.net.IpAddress.ListenOptions) Self {
        self.udp_port = port;
        self.address_value = address;
        self.std_io = std_io;
        self.server_options = listen_options;
        return Self;
    }
};

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
    const nt_server: server = .{};
    const net_server = try server.init(nt_server, &io_thread, opts);
    try std.testing.expect(net_server.address_value.len > 0);
    try std.testing.expect(net_server.port > 0);
    print("address-{s}\n", .{net_server.address_value});
    net_server.stream_data();
}

test "init_with_address" {
    const opts = std.Io.net.IpAddress.ListenOptions{
        .reuse_address = true,
        .mode = .dgram,
        .protocol = .udp,
    };
    var arena = std.heap.ArenaAllocator.init(std.heap.page_allocator);
    defer arena.deinit();
    const allocator = arena.allocator();
    var thread = std_thread.init(allocator);
    defer thread.deinit();
    var io_thread = thread.io();
    const u_server: server = .{};
    _ = try u_server.udp_init("224.0.0.1", 321000, &io_thread, opts);
}
