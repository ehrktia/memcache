const std = @import("std");
const print = std.debug.print;

pub fn tcp_server() !void {
    const address = std.net.Address.initIp4([4]u8{ 0, 0, 0, 0 }, 9999);
    // std.debug.print("starting tcp server\n", .{});
    var server = try address.listen(.{ .reuse_address = true });
    defer server.deinit();
    const conn = try server.accept();
    defer conn.stream.close();
    // var read_buffer: [1024]u8 = undefined;
    var writer_buffer: [1024]u8 = undefined;
    // var conn_writer = conn.stream.writer(&.{});
    const net_writer = std.net.Stream.Writer.init(conn.stream, &writer_buffer);
    var writer: std.Io.Writer = net_writer.interface;
    std.debug.print("accepting conn from:{f}\n", .{conn.address});
    var net_conn_reader = conn.stream.reader(&.{});
    // var net_conn_reader = std.net.Stream.Reader.init(conn.stream, &read_buffer);
    var reader = &net_conn_reader.interface().*;
    // try reader.readSliceAll(&read_buffer);
    std.debug.print("reader:{d}\n", .{reader.buffer.len});

    // try reader.readSliceAll(&buffer);
    // _ = try reader.streamDelimiter(&writer, '\n');
    // const data = reader.buffered();
    while (true) {
        const data = reader.takeDelimiterInclusive(&read_buffer) catch |e| {
            print("error reading:{any}\n", .{e});
            return;
        };
        try writer.writeAll(data);
    }
}
