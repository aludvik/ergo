package ergo_test

import (
  "errors"

  "github.com/aludvik/ergo"
)


func Example_1() {
  ergo.New(1).
    Map(func(i int) int { return i + 1 }).
    AndThen(func(i int) ergo.Result { return ergo.New(errors.New("error")) }).
    WrapErr("wrapper")
}

func Example_2() {
  v, err := ergo.FromPair(nil, errors.New("error")).
    MapErr(func(error) error { return errors.New("new error") }).
    IntoPair()
  if err != nil {
    // ...
  } else {
    _ = v.(int)
  }
}
