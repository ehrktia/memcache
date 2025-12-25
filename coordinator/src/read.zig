const std = @import("std");
const print = std.debug.print;

pub fn read_file(io: std.Io, file_name: []const u8, allocator: std.mem.Allocator) ![]u8 {
    const config_dir = std.Io.Dir.cwd();
    const buf = try allocator.alloc(u8, 1096);
    const data = try std.Io.Dir.readFile(config_dir, io, file_name, buf);
    return data;
}

test "read_file" {
    var arena = std.heap.ArenaAllocator.init(std.heap.page_allocator);
    defer arena.deinit();
    var thread = std.Io.Threaded.init(arena.allocator());
    defer thread.deinit();
    const data = try read_file(thread.ioBasic(), "config.json", arena.allocator());
    print("data:{s}\n", .{data});
    try std.testing.expect(data.len > 0);
}
