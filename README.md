# ergo

[![CircleCI](https://circleci.com/gh/aludvik/ergo.svg?style=svg)](https://circleci.com/gh/aludvik/ergo) [![GoDoc](https://godoc.org/github.com/aludvik/ergo?status.svg)](https://godoc.org/github.com/aludvik/ergo)

Package ergo provides more ergonomic error handling for Go. The standard error
handling patterns in Go are verbose and unwieldy. The ergo package provides
more ergonomic error handling methods at the cost of type safety. At its core,
it defines a Result type that is based on the Rust error handling primitive
with the same name: https://doc.rust-lang.org/std/result/enum.Result.html.
