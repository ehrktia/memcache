const std = @import("std");
const print = std.debug.print;
const net_address = std.Io.net.IpAddress;

pub fn tcp_server(io: std.Io) !void {
    const address = try net_address.parse("::", 9999);
    const opts: net_address.ListenOptions = .{
        .reuse_address = true,
        .mode = .stream,
        .protocol = .tcp,
    };
    var server = try net_address.listen(address, io, opts);
    defer server.deinit(io);
    var buffer: [1096]u8 = undefined;
    const stream = try server.accept(io);
    defer stream.close(io);
    const stream_reader = stream.reader(io, &buffer);
    var reader = stream_reader.interface;
    const data_received = reader.buffered();
    print("data received:{s}\n", .{data_received});
}
