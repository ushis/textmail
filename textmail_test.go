package textmail

import (
  "bytes"
  "io/ioutil"
  "path/filepath"
  "os"
  "testing"
)

func TestFormat(t *testing.T) {
  paths, err := filepath.Glob("./testdata/*.html")

  if err != nil {
    t.Fatal(err)
  }

  for _, path := range paths {
    file, err := os.Open(path)

    if err != nil {
      t.Fatal(err)
    }
    html, err := ioutil.ReadAll(file)

    if err != nil {
      t.Fatal(err)
    }
    file.Close()

    formatter, err := Format(html)

    if err != nil {
      t.Fatal(err)
    }
    file, err = os.Open(path[:len(path)-4] + "txt")

    if err != nil {
      t.Fatal(err)
    }
    txt, err := ioutil.ReadAll(file)

    if err != nil {
      t.Fatal(err)
    }
    file.Close()

    if res := formatter.Bytes(); bytes.Compare(res, txt) != 0 {
      print(string(res))
      t.Fatal("Mismatch for:", path)
    }
  }
}
