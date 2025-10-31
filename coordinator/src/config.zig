const ziggy = @import("ziggy");
const std = @import("std");

pub const Config = struct { heartbeat: []const u8 };
// var buffer: [1024]u8 = undefined;
const gpa_allocator = std.heap.GeneralPurposeAllocator(.{});
var gpa = gpa_allocator.init();
defer gpa_allocator.deinit(&gpa)
// const config = ziggy.parseLeaky();
