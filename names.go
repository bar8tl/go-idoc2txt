package idoc2txt

import "log"
import "os"
import "path/filepath"
import "strings"

// Functions in this file need to be further developped, status is in progress
func PassFilter(s Settings_tp, f os.FileInfo) bool {
  if strings.Contains(f.Name(), "processed") {
    return false
  }
  return true
}

func RanameInpFile(Inpdr string, f os.FileInfo) {
  filtp := filepath.Ext(f.Name())
  filnm := strings.TrimRight(f.Name(), filtp)
  oldName := Inpdr + f.Name()
  newName := Inpdr + "inp_" + filnm + "_processed" + filtp
  err := os.Rename(oldName, newName)
  if err != nil {
    log.Fatalf("Input file %s renaming error: %s\r\n", oldName, err)
  }
}

func RanameOutFile(Outdr string, f os.FileInfo) {
  filtp := filepath.Ext(f.Name())
  filnm := strings.TrimRight(f.Name(), filtp)
  oldName := Outdr + f.Name()
  newName := Outdr + "out_" + filnm + filtp
  err := os.Rename(oldName, newName)
  if err != nil {
    log.Fatalf("Output file %s renaming error: %s\r\n", oldName, err)
  }
}
