const std = @import("std");

pub export fn x(y: c_int) c_int {
    return y + 2;
}
