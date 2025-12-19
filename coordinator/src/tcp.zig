const std = @import("std");
const print = std.debug.print;

pub fn tcp_server(io: std.Io) !void {
    const address = try std.Io.net.Ip4Address.parse("0.0.0.0", 9999);
    const stream = try std.Io.net.IpAddress.connect(address, io, .{ .mode = .stream, .protocol = .tcp, .timeout = .{ .duration = .{ .raw = .{ .nanoseconds = 3 * 1000000000 }, .clock = .awake } } });
    defer stream.close(io);
    var read_buffer: [1096]u8 = undefined;
    const reader = stream.reader(io, &read_buffer);

    // var server = try address.listen(.{ .reuse_address = true });
    // defer server.deinit();
    // const conn = try server.accept();
    // defer conn.stream.close();
    // var read_buffer: [1024]u8 = undefined;
    // var writer_buffer: [1024]u8 = undefined;
    // var conn_writer = conn.stream.writer(&.{});
    // const net_writer = std.net.Stream.Writer.init(conn.stream, &writer_buffer);
    // var writer: std.Io.Writer = net_writer.interface;
    // std.debug.print("accepting conn from:{f}\n", .{conn.address});
    // var net_conn_reader = conn.stream.reader(&.{});
    // var net_conn_reader = std.net.Stream.Reader.init(conn.stream, &read_buffer);
    // var reader = &net_conn_reader.interface().*;
    // try reader.readSliceAll(&read_buffer);
    // std.debug.print("reader:{d}\n", .{reader.buffer.len});
    // print("reader:{d}\n", .{reader.buffer.len});

    // try reader.readSliceAll(&buffer);
    // _ = try reader.streamDelimiter(&writer, '\n');
    // const data = reader.buffered();
    // while (true) {
    //     const data = reader.takeDelimiterInclusive(&read_buffer) catch |e| {
    //         print("error reading:{any}\n", .{e});
    //         return;
    //     };
    //     try writer.writeAll(data);
    // }
}
