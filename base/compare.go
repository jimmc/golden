package base

import (
  "bytes"
  "fmt"
  "io/ioutil"
)

func CompareOutToGolden(outfilepath, goldenfilepath string) error {
  outcontent, err := ioutil.ReadFile(outfilepath)
  if err != nil {
    return fmt.Errorf("error reading back output file %s: %v", outfilepath, err)
  }
  goldencontent, err := ioutil.ReadFile(goldenfilepath)
  if err != nil {
    return fmt.Errorf("error reading golden file %s: %v", goldenfilepath, err)
  }
  if !bytes.Equal(outcontent, goldencontent) {
    return fmt.Errorf("outfile %s does not match golden file %s", outfilepath, goldenfilepath)
  }
  return nil
}
