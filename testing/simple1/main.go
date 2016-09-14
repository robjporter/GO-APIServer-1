package myPackage

import (
    "sync/atomic"
)

var count uint64

func Count() uint64 {
    atomic.AddUint64(&count,1)
    return count
}
