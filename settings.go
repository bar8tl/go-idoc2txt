package idoc2txt

import "encoding/json"
import "io/ioutil"
import "log"
import "os"
import "strings"
import "time"

// Data type to get a list of run-level parameters from command execution line (Structured like: -Option:Parameter1:Parameter2)
type Params_tp struct {
  Cmdpr []Param_tp
}
type Param_tp struct {
  Optn, Prm1, Prm2 string
}
func (p *Params_tp) NewParams() {
  if len(os.Args) == 0 {
    log.Printf("Run option missing\r\n")
    return
  }
  for i := 0; i < len(os.Args); i++ {
    curarg := os.Args[i]
    if curarg[0:1] == "-" || curarg[0:1] == "/" {
      optn := strings.ToLower(curarg[1:len(curarg)])
      prm1, prm2 := "", ""
      if optn != "" {
        if strings.Index(optn, ":") != -1 {
          prm1 = optn[strings.Index(optn, ":")+1 : len(optn)]
          optn = optn[0:strings.Index(optn, ":")]
          if strings.Index(prm1, ":") != -1 {
            prm2 = prm1[strings.Index(prm1, ":")+1 : len(prm1)]
            prm1 = prm1[0:strings.Index(prm1, ":")]
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

// Data type to get an json structure from configuration file to define potential program-level and run-level parameters
type Config_tp struct {
  Const   Constant_tp `json:"constants"`
  Progm   Program_tp  `json:"program"`
  Runlv   []Run_tp    `json:"run"`
}
type Constant_tp struct {
  Cntrl string `json:"contrl"`
  Clien string `json:"client"`
}
type Program_tp struct {
  Dbonm string `json:"dboNam"`
  Dbodr string `json:"dboDir"`
  Inpdr string `json:"inpDir"`
  Outdr string `json:"outDir"`
  Ifilt string `json:"inFilt"`
  Ifnam string `json:"inName"`
  Ofnam string `json:"ouName"`
}
type Run_tp struct {
  Optcd string `json:"option"`
  Objnm string `json:"objNam"`
  Dbonm string `json:"dboNam"`
  Dbodr string `json:"dboDir"`
  Inpdr string `json:"inpDir"`
  Outdr string `json:"outDir"`
  Ifilt string `json:"inFilt"`
  Ifnam string `json:"inName"`
  Ofnam string `json:"ouName"`
  Rcvpf string `json:"rcPrnF"`
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
}

// Data type to define global environment variables
type Envmnt_tp struct {
  Cnnsq, Cnnst        string
  Dbonm, Dbodr        string
  Objnm, Inpdr, Outdr string
  Ifilt, Ifnam, Ofnam string
  Cntrl, Clien, Rcvpf string
  Found               bool
  Dtsys, Dtcur, Dtnul time.Time
}
func (e *Envmnt_tp) NewEnvmnt(s Settings_tp) {
  e.Cnnsq = CNNS_SQLIT3
  e.Dtsys, e.Dtcur, e.Dtnul = time.Now(), time.Now(), time.Date(1901, 1, 1, 0, 0, 0, 0, time.UTC)
  e.Cntrl = ternary_op(len(s.Const.Cntrl) > 0, strings.TrimSpace(s.Const.Cntrl), CONTROL_CODE)
  e.Clien = ternary_op(len(s.Const.Clien) > 0, strings.TrimSpace(s.Const.Clien), CLIENT_CODE)
  e.Dbonm = ternary_op(len(s.Progm.Dbonm) > 0, strings.TrimSpace(s.Progm.Dbonm), DB_NAME)
  e.Dbodr = ternary_op(len(s.Progm.Dbodr) > 0, strings.TrimSpace(s.Progm.Dbodr), DB_DIR)
  e.Inpdr = ternary_op(len(s.Progm.Inpdr) > 0, strings.TrimSpace(s.Progm.Inpdr), INPUTS_DIR)
  e.Outdr = ternary_op(len(s.Progm.Outdr) > 0, strings.TrimSpace(s.Progm.Outdr), OUTPUTS_DIR)
  e.Ifilt = ternary_op(len(s.Progm.Ifilt) > 0, strings.TrimSpace(s.Progm.Ifilt), INPUTS_FILTER)
  e.Ifnam = ternary_op(len(s.Progm.Ifnam) > 0, strings.TrimSpace(s.Progm.Ifnam), INPUTS_NAMING)
  e.Ofnam = ternary_op(len(s.Progm.Ofnam) > 0, strings.TrimSpace(s.Progm.Ofnam), OUTPUTS_NAMING)
}

// Data type to be used as container of program-level and run-level settings. Pseudo-inheritance is used to simplify names of settings
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

//----------------------------------------------------------------------------------------------------------------------------------
func (e *Envmnt_tp) SetRunVars(p Param_tp, s Settings_tp) {
  if len(p.Prm1) > 0 {
    e.Objnm = strings.TrimSpace(p.Prm1)
  } else {
    log.Fatalf("Error: Not possible to determine IDOC-Type name.\r\n")
  }
  found := false
  for i := 0; i < len(s.Runlv) && !found; i++ {
    if p.Optn == s.Runlv[i].Optcd && p.Prm1 == s.Runlv[i].Objnm {
      e.Found = true
      if p.Optn == "cdb" || p.Optn == "rid" || p.Optn == "unf" {
        e.Objnm = ternary_op(len(s.Runlv[i].Objnm) > 0, strings.TrimSpace(s.Runlv[i].Objnm), e.Objnm)
        e.Dbonm = ternary_op(len(s.Runlv[i].Dbonm) > 0, strings.TrimSpace(s.Runlv[i].Dbonm), e.Dbonm)
        e.Dbodr = ternary_op(len(s.Runlv[i].Dbodr) > 0, strings.TrimSpace(s.Runlv[i].Dbodr), e.Dbodr)
      }
      if p.Optn == "rid" || p.Optn == "unf" {
        e.Inpdr = ternary_op(len(s.Runlv[i].Inpdr) > 0, strings.TrimSpace(s.Runlv[i].Inpdr), e.Inpdr)
        e.Outdr = ternary_op(len(s.Runlv[i].Outdr) > 0, strings.TrimSpace(s.Runlv[i].Outdr), e.Outdr)
      }
      if p.Optn == "unf" {
        e.Ifilt = ternary_op(len(s.Runlv[i].Ifilt) > 0, strings.TrimSpace(s.Runlv[i].Ifilt), e.Ifilt)
        e.Ifnam = ternary_op(len(s.Runlv[i].Ifnam) > 0, strings.TrimSpace(s.Runlv[i].Ifnam), e.Ifnam)
        e.Ofnam = ternary_op(len(s.Runlv[i].Ofnam) > 0, strings.TrimSpace(s.Runlv[i].Ofnam), e.Ofnam)
        e.Rcvpf = ternary_op(len(s.Runlv[i].Rcvpf) > 0, strings.TrimSpace(s.Runlv[i].Rcvpf), e.Rcvpf)
      }
    }
  }
  e.Cnnst = strings.Replace(e.Cnnsq, "@", e.Dbodr+e.Dbonm, 1)
}

func ternary_op(statement bool, tcond, fcond string) string {
  if statement {
    return tcond
  }
  return fcond
}
