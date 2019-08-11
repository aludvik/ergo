// Package ergo provides more ergonomic error handling for Go.
//
// The standard error handling patterns in Go are verbose and unwieldy. The
// ergo package provides more ergonomic error handling methods at the cost of
// type safety. At its core, it defines a Result type that is based on the Rust
// error handling primitive with the same name:
// https://doc.rust-lang.org/std/result/enum.Result.html.
package ergo

import (
  "fmt"
  "reflect"
)

// Result represents a the result of an operation that may succeed with a value
// or fail with an error. Its behavior depends on whether the value it wraps
// implements the error interface. It is the core type defined by the ergo
// package.
type Result struct {
  v interface{}
}

// New wraps a value in a Result.
func New(v interface{}) Result {
  return Result{v}
}

// FromPair converts a (value, error) pair into a Result.
func FromPair(ok interface{}, err error) Result {
  if err != nil {
    return Result{err}
  }
  return Result{ok}
}

// Unwrap returns the value wrapped by the Result.
func (r Result) Unwrap() interface{} {
  return r.v
}

// IntoPair converts a Result into a (value, error) pair.
func (r Result) IntoPair() (interface{}, error) {
  switch v := r.Unwrap().(type) {
  case error:
    return nil, v
  default:
    return v, nil
  }
}

// IsErr returns true if the value wrapped by Result implements error.
func (r Result) IsErr() bool {
  _, ok := r.v.(error)
  return ok
}

// Err returns the value wrapped by Result if it implements error or nil.
func (r Result) Err() error {
  if err, ok := r.v.(error); ok {
    return err
  }
  return nil
}

// IsOk returns true if the value wrapped by Result does not implement error.
func (r Result) IsOk() bool {
  return !r.IsErr()
}

// Ok returns the value wrapped by Result if it doesn't implement error or nil.
func (r Result) Ok() interface{} {
  if r.IsErr() {
    return nil
  }
  return r.v
}

// WrapErr wraps the inner value with a new message if it's an error.
func (r Result) WrapErr(msg string) Result {
  return r.MapErr(func(err error) error {
    return fmt.Errorf("%v: %v", msg, err)
  })
}

// MapErr calls the provided function with the inner value if it's an error.
func (r Result) MapErr(fn func(err error) error) Result {
  if err, ok := r.v.(error); ok {
    return New(fn(err))
  }
  return r
}

// Map calls the provided function with the inner value if it's not an error.
// The provided function must take exactly one argument whose type matches the
// type of the inner value when it's Ok and return exactly one value. This is
// not enforced by the compiler.
func (r Result) Map(fn interface{}) Result {
  if r.IsErr() {
    return r
  }

  ret := reflect.ValueOf(fn).Call([]reflect.Value{reflect.ValueOf(r.v)})
  return New(ret[0].Interface())
}

// And returns this Result if it's an error, otherwise it returns the other
// Result.
func (r Result) And(other Result) Result {
  if r.IsOk() {
    return other
  }
  return r
}

// AndThen calls the provided function with the inner value if the inner value
// is not an error. The provided function must take exactly one argument and
// return a Result. This is not enforced by the compiler.
func (r Result) AndThen(fn interface{}) Result {
  if r.IsErr() {
    return r
  }

  ret := reflect.ValueOf(fn).Call([]reflect.Value{reflect.ValueOf(r.v)})
  return ret[0].Interface().(Result)
}

// Or returns this Result if it's Ok, otherwise it returns the other Result.
func (r Result) Or(other Result) Result {
  if r.IsErr() {
    return other
  }
  return r
}

// OrElse calls the provided function with the inner value if it's an error.
func (r Result) OrElse(fn func(error) Result) Result {
  if r.IsOk() {
    return r
  }

  return fn(r.Err())
}
