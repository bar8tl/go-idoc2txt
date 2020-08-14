package idoc2txt

import "encoding/json"
import "io/ioutil"
import "log"
import "os"
import "strings"
import "time"

// Data type to get a list of run-level parameters from command execution line
// Structured as: -Option:Parameter1:Parameter2
type Params_tp struct {
  Cmdpr []Param_tp
}

func (p *Params_tp) NewParams() {
  if len(os.Args) == 0 {
    log.Printf("Run option missing\r\n")
    return
  }
  for _, curarg := range os.Args {
    if curarg[0:1] == "-" || curarg[0:1] == "/" {
      optn := strings.ToLower(curarg[1:len(curarg)])
      prm1 := ""
      prm2 := ""
      if optn != "" {
        if strings.Index(optn, ":") != -1 {
          prm1 = optn[strings.Index(optn, ":")+1:len(optn)]
          optn = strings.TrimSpace(optn[0:strings.Index(optn, ":")])
          if strings.Index(prm1, ":") != -1 {
            prm2 = strings.TrimSpace(prm1[strings.Index(prm1, ":")+1:len(prm1)])
            prm1 = strings.TrimSpace(prm1[0:strings.Index(prm1, ":")])
          }
        }
        p.Cmdpr = append(p.Cmdpr, Param_tp{optn, prm1, prm2})
      } else {
        log.Printf("Run option missing\r\n")
      }
    }
  }
  return
}

// Data type to get an json structure from configuration file to define
// potential program-level and run-level parameters
type Config_tp struct {
  Const Constant_tp `json:"constants"`
  Progm Program_tp  `json:"program"`
  Runlv []Run_tp    `json:"run"`
}

func (c *Config_tp) NewConfig(fname string) {
  f, err := os.Open(fname)
  if err != nil {
    log.Fatalf("File %s opening error: %s\n", fname, err)
  }
  defer f.Close()
  jsonv, _ := ioutil.ReadAll(f)
  err = json.Unmarshal(jsonv, &c)
  if err != nil {
    log.Fatalf("File %s reading error: %s\n", fname, err)
  }
  c.Const.Cntrl = strings.TrimSpace(c.Const.Cntrl)
  c.Const.Clien = strings.TrimSpace(c.Const.Clien)
  c.Progm.Dbonm = strings.TrimSpace(c.Progm.Dbonm)
  c.Progm.Dbodr = strings.TrimSpace(c.Progm.Dbodr)
  c.Progm.Inpdr = strings.TrimSpace(c.Progm.Inpdr)
  c.Progm.Outdr = strings.TrimSpace(c.Progm.Outdr)
  c.Progm.Ifilt = strings.TrimSpace(c.Progm.Ifilt)
  c.Progm.Ifnam = strings.TrimSpace(c.Progm.Ifnam)
  c.Progm.Ofnam = strings.TrimSpace(c.Progm.Ofnam)
  for i := 0; i < len(c.Runlv); i++ {
    c.Runlv[i].Optcd = strings.TrimSpace(c.Runlv[i].Optcd)
    c.Runlv[i].Objnm = strings.TrimSpace(c.Runlv[i].Objnm)
    c.Runlv[i].Dbonm = strings.TrimSpace(c.Runlv[i].Dbonm)
    c.Runlv[i].Dbodr = strings.TrimSpace(c.Runlv[i].Dbodr)
    c.Runlv[i].Inpdr = strings.TrimSpace(c.Runlv[i].Inpdr)
    c.Runlv[i].Outdr = strings.TrimSpace(c.Runlv[i].Outdr)
    c.Runlv[i].Ifilt = strings.TrimSpace(c.Runlv[i].Ifilt)
    c.Runlv[i].Ifnam = strings.TrimSpace(c.Runlv[i].Ifnam)
    c.Runlv[i].Ofnam = strings.TrimSpace(c.Runlv[i].Ofnam)
    c.Runlv[i].Rcvpf = strings.TrimSpace(c.Runlv[i].Rcvpf)
  }
}

// Data type to define global environment variables
type Envmnt_tp struct {
  Cnnsq string
  Cnnst string
  Cntrl string
  Clien string
  Dbonm string
  Dbodr string
  Inpdr string
  Outdr string
  Ifilt string
  Ifnam string
  Ofnam string
  Objnm string
  Rcvpf string
  Found bool
  Mitm  bool
  Sgrp  bool
  Ssgm  bool
  Dtsys time.Time
  Dtcur time.Time
  Dtnul time.Time
}

func (e *Envmnt_tp) NewEnvmnt(s Settings_tp) {
  e.Cnnsq = CNNS_SQLIT3
  e.Cntrl = ternary_op(len(s.Const.Cntrl) > 0, s.Const.Cntrl, CONTROL_CODE)
  e.Clien = ternary_op(len(s.Const.Clien) > 0, s.Const.Clien, CLIENT_CODE)
  e.Dbonm = ternary_op(len(s.Progm.Dbonm) > 0, s.Progm.Dbonm, DB_NAME)
  e.Dbodr = ternary_op(len(s.Progm.Dbodr) > 0, s.Progm.Dbodr, DB_DIR)
  e.Inpdr = ternary_op(len(s.Progm.Inpdr) > 0, s.Progm.Inpdr, INPUTS_DIR)
  e.Outdr = ternary_op(len(s.Progm.Outdr) > 0, s.Progm.Outdr, OUTPUTS_DIR)
  e.Ifilt = ternary_op(len(s.Progm.Ifilt) > 0, s.Progm.Ifilt, INPUTS_FILTER)
  e.Ifnam = ternary_op(len(s.Progm.Ifnam) > 0, s.Progm.Ifnam, INPUTS_NAMING)
  e.Ofnam = ternary_op(len(s.Progm.Ofnam) > 0, s.Progm.Ofnam, OUTPUTS_NAMING)
  e.Dtsys = time.Now()
  e.Dtcur = time.Now()
  e.Dtnul = time.Date(1901, 1, 1, 0, 0, 0, 0, time.UTC)
}

// Data type to be used as container of program-level and run-level settings.
// Pseudo-inheritance is used to simplify names of settings
type Settings_tp struct {
  Config_tp
  Params_tp
  Envmnt_tp
}

func NewSettings(fname string) Settings_tp {
  var s Settings_tp
  s.NewParams()
  s.NewConfig(fname)
  s.NewEnvmnt(s)
  return s
}

//------------------------------------------------------------------------------
func (s *Settings_tp) SetRunVars(p Param_tp) {
  if len(p.Prm1) > 0 {
    s.Objnm = p.Prm1
  } else {
    log.Fatalf("Error: Not possible to determine IDOC-Type name.\r\n")
  }
  s.Found = false
  for _, runlv := range s.Runlv {
    if p.Optn == runlv.Optcd && p.Prm1 == runlv.Objnm {
      if p.Optn == "cdb" || p.Optn == "rid" || p.Optn == "unf" {
        s.Objnm = ternary_op(len(runlv.Objnm) > 0, runlv.Objnm, s.Objnm)
        s.Dbonm = ternary_op(len(runlv.Dbonm) > 0, runlv.Dbonm, s.Dbonm)
        s.Dbodr = ternary_op(len(runlv.Dbodr) > 0, runlv.Dbodr, s.Dbodr)
      }
      if p.Optn == "rid" || p.Optn == "unf" {
        s.Inpdr = ternary_op(len(runlv.Inpdr) > 0, runlv.Inpdr, s.Inpdr)
        s.Outdr = ternary_op(len(runlv.Outdr) > 0, runlv.Outdr, s.Outdr)
      }
      if p.Optn == "unf" {
        s.Ifilt = ternary_op(len(runlv.Ifilt) > 0, runlv.Ifilt, s.Ifilt)
        s.Ifnam = ternary_op(len(runlv.Ifnam) > 0, runlv.Ifnam, s.Ifnam)
        s.Ofnam = ternary_op(len(runlv.Ofnam) > 0, runlv.Ofnam, s.Ofnam)
        s.Rcvpf = ternary_op(len(runlv.Rcvpf) > 0, runlv.Rcvpf, s.Rcvpf)
      }
      s.Found = true
      break
    }
  }
  if p.Optn == "rid" {
    s.Mitm = true
    s.Sgrp = false
    s.Ssgm = false
    if len(p.Prm2) > 0 {
      mflds := strings.Split(p.Prm2, ".")
      for i := 0; i < len(mflds); i++ {
        switch strings.ToLower(mflds[i]) {
        case ITM:
          s.Mitm = true
        case GRP:
          s.Sgrp = true
        case SGM:
          s.Ssgm = true
        default:
          s.Mitm = true
          s.Sgrp = false
          s.Ssgm = false
        }
      }
    }
  }
  s.Cnnst = strings.Replace(s.Cnnsq, "@", s.Dbodr+s.Dbonm, 1)
}

func ternary_op(statement bool, tcond, fcond string) string {
  if statement {
    return tcond
  }
  return fcond
}
