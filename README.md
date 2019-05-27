[![GoDoc](https://godoc.org/github.com/dc0d/bufferpool?status.svg)](https://godoc.org/github.com/dc0d/bufferpool)


# bufferpool
a pool of byte slice, without the memory fragmentation

Byte slices in this pool are backed by a big array and all of them have a fixed, equal size so they will not grow as a result of `append` actions. They are safe for concurrent use because the ranges have not overlaps.
