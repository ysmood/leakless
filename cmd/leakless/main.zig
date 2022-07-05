const std = @import("std");
const stdout = std.io.getStdOut().writer();
var arena = std.heap.ArenaAllocator.init(std.heap.page_allocator);
const allocator = arena.allocator();

pub fn main() anyerror!void {
    const args = try std.process.argsAlloc(allocator);

    if (args.len < 3) {
        _ = try stdout.write("usage: leakless <port> <cmd...>\n");
        std.os.exit(1);
    }

    var p = try run(args[2..]);

    _ = try std.Thread.spawn(.{}, guard, .{args[1], p});

    std.os.exit(switch (try p.wait()) {
        std.ChildProcess.Term.Exited => |v| v,
        std.ChildProcess.Term.Signal => |v| @truncate(u8, v),
        std.ChildProcess.Term.Stopped => |v| @truncate(u8, v),
        std.ChildProcess.Term.Unknown => |v| @truncate(u8, v),
    });
}

fn run(args: []const []const u8) !*std.ChildProcess {
    var p = try std.ChildProcess.init(args, allocator);
    try p.spawn();

    return p;
}

fn guard(port: []u8, p: *std.ChildProcess) !void {
    defer _ = p.kill() catch std.ChildProcess.Term.Unknown;

    const addr = try std.net.Address.parseIp("127.0.0.1", try std.fmt.parseInt(u16, port, 10));
    const sock = try std.net.tcpConnectToAddress(addr);

    var buf: [1]u8 = undefined;
    _ = sock.reader().read(&buf) catch 0;
}