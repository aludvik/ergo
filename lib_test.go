package ergo_test

import (
  "errors"
  "fmt"
  "testing"

  . "github.com/aludvik/ergo"
)

func address(i interface{}) string {
  return fmt.Sprintf("%p", i)
}

func TestBasics(t *testing.T) {
  t.Run("Ok", func(t *testing.T) {
    t.Run("Scalar", func(t *testing.T) {
      r := New(0)
      if !r.IsOk() { t.Error("Should be Ok") }
      v, ok := r.Ok().(int)
      if !ok { t.Error("Should be Ok") }
      if v != 0 { t.Error("Should be 0") }

      if r.IsErr() { t.Error("Should not be Err") }
      if r.Err() != nil { t.Error("Should be nil") }
    })

    t.Run("Pointer", func(t *testing.T) {
      b := []byte{0x0}
      r := New(b)
      if !r.IsOk() { t.Error("Should be Ok") }

      v, ok := r.Ok().([]byte)
      if !ok { t.Error("Should be Ok") }
      if address(v) != address(b) { t.Error("Should be same") }

      b[0] = 0x1
      if 0x1 != v[0] { t.Error("Should be same") }

      if r.IsErr() { t.Error("Should not be Err") }
      if r.Err() != nil { t.Error("Should be nil") }
    })

    t.Run("nil", func(t *testing.T) {
      r := New(nil)
      if !r.IsOk() { t.Error("Should be Ok") }
      if r.IsErr() { t.Error("Should not be error") }
    })
  })

  t.Run("Err", func(t *testing.T) {
    err := errors.New("error")
    r := New(err)
    if !r.IsErr() { t.Error("Should be Err") }
    frr := r.Err()
    if err != frr { t.Error("Should be same") }

    if r.IsOk() { t.Error("Should not be Ok") }
    if r.Ok() != nil { t.Error("Should be nil") }
  })
}

func existingFunc(i interface{}) (int, error) {
  if err, ok := i.(error); ok {
    return 0, err
  }
  return i.(int), nil
}

func TestFromPair(t *testing.T) {
  r1 := FromPair(existingFunc(1))
  if !r1.IsOk() { t.Error("Should be Ok") }
  if r1.Ok().(int) != 1 { t.Error("Should be 1") }

  r2 := FromPair(existingFunc(errors.New("error")))
  if !r2.IsErr() { t.Error("Should be Err") }
  if r2.Err().Error() != "error" { t.Error("Should be error") }
}

func TestIntoPair(t *testing.T) {
  r1 := New(1)
  i1, e1 := r1.IntoPair()
  if e1 != nil { t.Error("Should be nil") }
  if i1.(int) != 1 { t.Error("Should be 1") }

  r2 := New(errors.New("error"))
  i2, e2 := r2.IntoPair()
  if e2 == nil { t.Error("Should not be nil") }
  if i2 != nil { t.Error("Should be nil") }
}

func TestMapErr(t *testing.T) {
  t.Run("Ok", func(t *testing.T) {
    v := new(int)
    *v = 5
    r := New(v)
    s := r.MapErr(func(err error) error {
      return errors.New("error")
    })

    if !s.IsOk() { t.Error("Should be Ok") }
    w, ok := s.Ok().(*int)
    if !ok { t.Error("Should be Ok") }
    if *w != *v { t.Error("Should be same") }
    if address(w) != address(v) { t.Error("Should be same") }
  })

  t.Run("Err", func(t *testing.T) {
    er1 := errors.New("1")
    er2 := errors.New("2")
    r := New(er1)
    s := r.MapErr(func(err error) error {
      return er2
    })

    if !s.IsErr() { t.Error("Should be Err") }
    if s.Err() != er2 { t.Error("Should be same") }
  })

  t.Run("Wrap", func(t *testing.T) {
    inner := errors.New("inner")
    r := New(inner)
    s := r.WrapErr("outer")

    if !s.IsErr() { t.Error("Should be Err") }
    if s.Err().Error() != "outer: inner" {
      t.Errorf("Incorrect wrap: %v", r)
    }
  })
}

func TestMap(t *testing.T) {
  t.Run("Ok", func(t *testing.T) {
    r := New(5).Map(func(i int) int {
      return i + 1
    })

    if !r.IsOk() { t.Error("Should be Ok") }
    if r.Ok().(int) != 6 {
      t.Errorf("Incorrect int: %v", r.Ok().(int))
    }
  })

  t.Run("Err", func(t *testing.T) {
    r := New(errors.New("error")).Map(func(i int) int {
      return i + 1
    })

    if r.Err().Error() != "error" {
      t.Error("Expected error")
    }
  })
}

func TestAnd(t *testing.T) {
  t.Run("Ok", func(t *testing.T) {
    o1 := New(1)
    o2 := New(2)

    a1 := o1.And(o2)
    if !a1.IsOk() { t.Error("Should be Ok") }
    if a1.Ok().(int) != 2 { t.Error("Should be 2") }

    e1 := New(errors.New("error"))
    a2 := o1.And(e1)
    if !a2.IsErr() { t.Error("Should be Err") }
  })

  t.Run("Err", func(t *testing.T) {
    e1 := New(errors.New("e1"))
    o1 := New(1)

    a1 := e1.And(o1)
    if !a1.IsErr() { t.Error("Should be Err") }
    if a1.Err().Error() != "e1" { t.Error("Should be e1") }

    e2 := New(errors.New("e2"))
    a2 := e1.And(e2)
    if !a2.IsErr() { t.Error("Should be Err") }
    if a2.Err().Error() != "e1" { t.Error("Should be e1") }
  })
}

func TestAndThen(t *testing.T) {
  t.Run("Ok", func(t *testing.T) {
    b := []byte{0x0}
    o := New(b)
    r1 := o.AndThen(func(c []byte) Result {
      return New(append(c, 0x1))
    })
    if r1.Ok().([]byte)[1] != 0x1 { t.Error("Should be 0x1") }
    r2 := o.AndThen(func(c []byte) Result {
      return New(errors.New("error"))
    })
    if r2.Err().Error() != "error" { t.Error("Should be error") }
  })

  t.Run("Err", func(t *testing.T) {
    e := New(errors.New("error"))
    r1 := e.AndThen(func(c []byte) Result {
      return New(1)
    })
    if r1.Err().Error() != "error" { t.Error("Should be error") }
    r2 := e.AndThen(func(c []byte) Result {
      return New(errors.New("nope"))
    })
    if r2.Err().Error() != "error" { t.Error("Should be error") }
  })
}

func TestOr(t *testing.T) {
  t.Run("Ok", func(t *testing.T) {
    o1 := New(1)
    o2 := New(2)

    a1 := o1.Or(o2)
    if a1.Ok().(int) != 1 { t.Error("Should be 1") }

    e1 := New(errors.New("error"))
    a2 := o1.Or(e1)
    if a2.Ok().(int) != 1 { t.Error("Should be 1") }
  })

  t.Run("Err", func(t *testing.T) {
    e1 := New(errors.New("e1"))
    o1 := New(1)

    a1 := e1.Or(o1)
    if a1.Ok().(int) != 1 { t.Error("Should be 1") }

    e2 := New(errors.New("e2"))
    a2 := e1.Or(e2)
    if a2.Err().Error() != "e2" { t.Error("Should be e2") }
  })
}

func TestOrElse(t *testing.T) {
  t.Run("Ok", func(t *testing.T) {
    o := New(1)
    r := o.OrElse(func(err error) Result {
      return New(errors.New("error"))
    })
    if r.Ok().(int) != 1 { t.Error("Should be 1") }
  })

  t.Run("Err", func(t *testing.T) {
    e1 := New(errors.New("e1"))
    r1 := e1.OrElse(func(err error) Result {
      return New(errors.New("e2"))
    })
    if r1.Err().Error() != "e2" { t.Error("Should be e2") }

    r2 := e1.OrElse(func(err error) Result {
      return New(1)
    })
    if r2.Ok().(int) != 1 { t.Error("Should be 1") }
  })
}
