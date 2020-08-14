package idoc2txt

import "bufio"
import "io"
import "log"
import "os"
import "strings"

// Data type to read SAP IDoc parser file and to upload IDoc definition detail
// and structure into an internal reference database
type Drid_tp struct {
  ri Drmitm_tp
  rg Drsgrp_tp
  rs Drssgm_tp
}

func NewDrid() *Drid_tp {
  var r Drid_tp
  return &r
}

func (r *Drid_tp) ProcInput(parm Param_tp, s Settings_tp) {
  s.SetRunVars(parm)
  ifile, err := os.Open(s.Inpdr + s.Objnm)
  if err != nil {
    log.Fatalf("Input file %s not found: %s\r\n", s.Inpdr+s.Objnm, err)
  }
  defer ifile.Close()
  r.ProcStartOfFile(s)
  rdr := bufio.NewReader(ifile)
  for l, _, err := rdr.ReadLine(); err != io.EOF; l, _, err = rdr.ReadLine() {
    if line := strings.TrimSpace(string(l)); len(line) > 0 {
      sline := r.ScanTextIdocLine(line)
      r.ProcLinesOfFile(s, sline)
    }
  }
  r.ProcEndOfFile(s)
}

func (r *Drid_tp) ProcStartOfFile(s Settings_tp) {
  if s.Mitm {
    r.ri.NewDrmitm(s)
  }
  if s.Sgrp {
    r.rg.NewDrsgrp(s, GRP)
  }
  if s.Ssgm {
    r.rs.NewDrssgm(s, SGM)
  }
}

func (r *Drid_tp) ProcLinesOfFile(s Settings_tp, sline Parsl_tp) {
  if s.Mitm {
    r.ri.GetData(sline)
  }
  if s.Sgrp {
    r.rg.GetData(sline)
  }
  if s.Ssgm {
    r.rs.GetData(sline)
  }
}

func (r *Drid_tp) ProcEndOfFile(s Settings_tp) {
  if s.Mitm {
    r.ri.IsrtData(s)
  }
}

// Function to identify individual tokens in SAP IDoc parser file
func (r *Drid_tp) ScanTextIdocLine(s string) (p Parsl_tp) {
  var key string
  var val string
  flds := strings.Fields(s)
  if len(flds) > 0 {
    key = flds[0]
    if (len(key) >= 6 && key[0:6] == "BEGIN_") ||
      (len(key) >= 4 && key[0:4] == "END_") {
      tokn := strings.Split(key, "_")
      if len(tokn) == 2 {
        p.Label.Ident = tokn[0]
        p.Label.Recnm = tokn[1]
        p.Label.Rectp = ""
      } else if len(tokn) == 3 {
        p.Label.Ident = tokn[0]
        p.Label.Recnm = tokn[1]
        p.Label.Rectp = tokn[2]
      }
    } else {
      p.Label.Ident = key
      p.Label.Recnm = ""
      p.Label.Rectp = ""
    }
  }
  if len(flds) > 1 {
    val = flds[1]
    for i := 2; i < len(flds); i++ {
      val += " " + flds[i]
    }
    p.Value = val
  }
  return p
}
