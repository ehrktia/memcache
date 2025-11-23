const std = @import("std");
const print = std.debug.print;

pub fn tcp_server() !void {
    const address = std.net.Address.initIp4([4]u8{ 0, 0, 0, 0 }, 9999);
    const server = try address.listen(.{});
    defer server.deinit();
    const conn = try server.accept();
    var buffer: [512]u8 = undefined;
    const fsStdOut = std.fs.File.stdout().writer(&buffer);
    const stdout: *std.Io.Writer = &fsStdOut.interface;
    try conn.address.format(stdout);
}
