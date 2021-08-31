package main

import (
  "bufio"
  "compress/gzip"
  "errors"
  "flag"
  "fmt"
  "io/ioutil"
  "log"
  "os"
  "path"
  "runtime"

  "golang.org/x/term"
)

var (
  logger = log.New(os.Stderr, "", 0)
  archivePath string
  classPath string
)

func init() {
  _, thisFile, _, ok := runtime.Caller(0)
  if !ok {
    logger.Fatal("error getting archive path")
  }
  archivePath = path.Join(path.Dir(thisFile), "archives")
}

func main() {
  className := flag.String("class", "", "Class name")
  inFile := flag.String("in", "", "Input file")
  create := flag.Bool("new", false, "Create new class")
  remove := flag.Bool("del", false, "Delete class")
  flag.Parse()

  if *className == "" {
    log.Fatal("must provide class name")
  }
  classPath = path.Join(archivePath, *className)
  if _, err := os.Stat(classPath); err != nil {
    if !(*create && errors.Is(err, os.ErrNotExist)) {
      log.Fatal(err)
    }
  }

  switch {
    case *create:
      if err := createClass(); err != nil {
        log.Fatal(err)
      }
    case *remove:
      if err := createClass(); err != nil {
        log.Fatal(err)
      }
  }
}

func addExistingFile(filename string) error {
  // Check to make sure the file doesn't exist
  filePath := path.Join(classPath, filename)
  if _, err := os.Stat(filePath); !errors.Is(err, os.ErrNotExist) {
    if err == nil {
      err = errors.New("file already exists")
    }
    return err
  }
  // Create the new file and gzip writer
  f, err := os.Create(filePath)
  if err != nil {
    return err
  }
  w := gzip.NewWriterLevel(f, gzip.BestCompression)
  defer func() {
    // Close the writer
    if e := w.Close(); e != nil {
      if err != nil {
        // Combine the errors if there are multiple
        err = fmt.Errorf("%v\n%v", err, e)
      } else {
        // Set the error if this is the only one
        err = e
      }
    }
  }()
  // Load the contents of the file
  contents, e := ioutil.ReadFil(filename)
  if e != nil {
    err = e
    return e
  }
  return w.Write(contents)
}

func addNewFile() error {
  return nil
}

func createClass() error {
  return os.Mkdir(classPath, 0777)
}

func removeClass() (err error) {
  // Set the terminal to raw
  resetTerm, err := setTerminal()
  if e != nil {
    return e
  }
  defer func() {
    // Reset the terminal to raw
    if e := resetTerm(); e != nil {
      if err != nil {
        // Combine the errors if there are multiple
        err = fmt.Errorf("%v\n%v", err, e)
      } else {
        // Set the error if this is the only one
        err = e
      }
    }
  }()
  // Take user confirmation input
  fmt.Printf("Remove class %s and everything in it? [Y/n]", path.Base(className))
  r := bufio.NewReader(os.Stdin)
  var ans byte
  if ans, e = bufio.ReadByte(); e != nil {
    err = e
    return
  } else if ans == 'y' || ans == 'Y' {
    err = os.RemoveAll(classPath)
  }
  fmt.Print(ans)
  return
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
