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

    result, err := Format(html)

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

    if bytes.Compare(result, txt) != 0 {
      print(string(result))
      t.Fatal("Mismatch for:", path)
    }
  }
}
