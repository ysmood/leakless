const std = @import("std");
const builtin = @import("builtin");
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

    _ = try std.Thread.spawn(.{}, guard, .{ args[1], p.pid });

    std.os.exit(switch (try p.wait()) {
        .Exited => |v| v,
        .Signal => |v| @truncate(u8, v),
        .Stopped => |v| @truncate(u8, v),
        .Unknown => |v| @truncate(u8, v),
    });
}

fn run(args: []const []const u8) !*std.ChildProcess {
    var p = try std.ChildProcess.init(args, allocator);
    try p.spawn();

    return p;
}

fn guard(port: []u8, pid: std.os.system.pid_t) !void {
    defer kill(pid);

    const pi = try std.fmt.parseInt(u16, port, 10);
    const addr = try std.net.Address.parseIp("127.0.0.1", pi);

    if (std.net.tcpConnectToAddress(addr)) |sock| {
        var buf: [1]u8 = undefined;
        _ = sock.reader().read(&buf) catch 0;
    } else |err| {
        _ = err catch null;
    }
}

fn kill(pid: std.os.system.pid_t) void {
    if (builtin.os.tag == .windows) {
        const id = std.fmt.allocPrint(allocator, "{}", .{pid}) catch "";
        if (std.ChildProcess.init([][]u8{"taskkill", "/t", "/f", "/pid", id}, allocator)) |p| {
            _ = p.spawnAndWait() catch .Unknown;
        } else |err| {
            _ = err catch null;
        }
    } else {
        _ = std.os.kill(pid, std.os.SIG.KILL) catch null;
    }
}
