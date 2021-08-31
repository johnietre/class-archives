package main

import (
  "bufio"
  "log"
  "os"

  "golang.org/x/term"
)

func main() {
  reset, err := setTerminal()
  if err != nil {
    log.Fatal(err)
  }
  defer func() {
    if err := reset(); err != nil {
      log.Fatal(err)
    }
  }()
  r := bufio.NewReader(os.Stdin)
  for {
    b, err := r.ReadByte()
    if err != nil {
      log.Println(err)
      return
    }
    if b == 'q' {
      break
    }
  }
}

func setTerminal() (resetFunc func() error, err error) {
  oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
  if err != nil {
    return nil, err
  }
  resetFunc = func() error {
    return term.Restore(int(os.Stdin.Fd()), oldState)
  }
  return
}
