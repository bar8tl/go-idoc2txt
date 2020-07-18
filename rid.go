package idoc2txt

import "bufio"
import "io"
import "log"
import "os"
import "strings"

// Data type to read SAP IDoc parser file and to upload IDoc definition detail and structure into an internal reference database
type Drid_tp struct {
  ri Drmitm_tp; rg Drsgrp_tp; rs Drssgm_tp
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
  for line, _, err := rdr.ReadLine(); err != io.EOF; line, _, err = rdr.ReadLine() {
    if l := strings.Trim(string(line), " "); len(l) > 0 {
      sline := r.ScanTextIdocLine(l)
      r.ProcLinesOfFile(s, sline)
    }
  }
  r.ProcEndOfFile(s)
}

func (r *Drid_tp) ProcStartOfFile(s Settings_tp) {
  if s.Mitm { r.ri.NewDrmitm(s)      }
  if s.Sgrp { r.rg.NewDrsgrp(s, GRP) }
  if s.Ssgm { r.rs.NewDrssgm(s, SGM) }
}

func (r *Drid_tp) ProcLinesOfFile(s Settings_tp, sline Parsl_tp) {
  if s.Mitm { r.ri.GetData(sline) }
  if s.Sgrp { r.rg.GetData(sline) }
  if s.Ssgm { r.rs.GetData(sline) }
}

func (r *Drid_tp) ProcEndOfFile(s Settings_tp) {
  if s.Mitm { r.ri.IsrtData(s)    }
}

// Function to identify individual tokens in SAP IDoc parser file
type Parsl_tp struct {
  Label Reclb_tp
  Value string
}
type Reclb_tp struct {
  Ident, Recnm, Rectp string
}

func (r *Drid_tp) ScanTextIdocLine(s string) (p Parsl_tp) {
  var key, val string
  flds := strings.Fields(s)
  if len(flds) > 0 {
    key = flds[0]
    if (len(key) >= 6 && key[0:6] == "BEGIN_") || (len(key) >= 4 && key[0:4] == "END_") {
      tokn := strings.Split(key, "_")
      if len(tokn) == 2 {
        p.Label.Ident, p.Label.Recnm, p.Label.Rectp = tokn[0], tokn[1], ""
      } else if len(tokn) == 3 {
        p.Label.Ident, p.Label.Recnm, p.Label.Rectp = tokn[0], tokn[1], tokn[2]
      }
    } else {
      p.Label.Ident, p.Label.Recnm, p.Label.Rectp = key, "", ""
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
