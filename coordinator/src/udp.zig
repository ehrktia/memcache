const std = @import("std");
const print = std.debug.print;
const net_address = std.Io.net.IpAddress;
const net_server = @import("./server.zig");

pub fn udp_server(io: *std.Io) !void {
    const opts: net_address.ListenOptions = .{
        .reuse_address = true,
        .mode = .dgram,
        .protocol = .udp,
    };
    const srv: net_server.server = .{};
    const net_serv = try srv.init(io, opts);
    return net_serv.listen_data();
}
