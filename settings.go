package idoc2txt

import "encoding/json"
import "io/ioutil"
import "log"
import "os"
import "strings"
import "time"

const CNNSTR = "file:@?file:locked.sqlite?cache=shared&mode=rwc"

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
  e.Cnnsq = CNNSTR
  e.Cntrl = s.Const.Cntrl
  e.Clien = s.Const.Clien
  e.Dtsys, e.Dtcur, e.Dtnul = time.Now(), time.Now(), time.Date(1901, 1, 1, 0, 0, 0, 0, time.UTC)
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

func (e *Envmnt_tp) SetRunVars(p Param_tp, s Settings_tp) {
  for i := 0; i < len(s.Runlv) && !e.Found; i++ {
    if p.Optn == s.Runlv[i].Optcd {
      e.Found = true
      if p.Optn == "cdb" || p.Optn == "rid" || p.Optn == "unf" {
        e.seekRunObjnm(i,p,s).seekRunDbonm(i,s).seekRunDbodr(i,s)
      }
      if p.Optn == "rid" || p.Optn == "unf" {
        e.seekRunInpdr(i,s).seekRunOutdr(i,s)
      }
      if p.Optn == "unf" {
        e.seekRunIfilt(i,s).seekRunIfnam(i,s).seekRunOfnam(i,s).seekRunRcvpf(i,s)
      }
    }
  }
  if !e.Found {
    if p.Optn == "cdb" || p.Optn == "rid" || p.Optn == "unf" {
      e.seekPgmObjnm(p).seekPgmDbonm(s).seekPgmDbodr(s)
    }
    if p.Optn == "rid" || p.Optn == "unf" {
      e.seekPgmInpdr(s).seekPgmOutdr(s)
    }
  }
  e.seekConstCntrl(s).seekConstClien(s)
  e.Cnnst = strings.Replace(e.Cnnsq, "@", e.Dbodr+e.Dbonm, 1)
}

func (e *Envmnt_tp) seekRunObjnm(i int, p Param_tp, s Settings_tp) *Envmnt_tp {
  if len(s.Runlv[i].Objnm) > 0 {
    e.Objnm = strings.TrimSpace(s.Runlv[i].Objnm)
  } else {
    e.seekPgmObjnm(p)
  }
  return e
}
func (e *Envmnt_tp) seekRunDbonm(i int, s Settings_tp) *Envmnt_tp {
  if len(s.Runlv[i].Dbonm) > 0 {
    e.Dbonm = strings.TrimSpace(s.Runlv[i].Dbonm)
  } else {
    e.seekPgmDbonm(s)
  }
  return e
}
func (e *Envmnt_tp) seekRunDbodr(i int, s Settings_tp) *Envmnt_tp {
  if len(s.Runlv[i].Dbodr) > 0 {
    e.Dbodr = strings.TrimSpace(s.Runlv[i].Dbodr)
  } else {
    e.seekPgmDbodr(s)
  }
  return e
}
func (e *Envmnt_tp) seekRunInpdr(i int, s Settings_tp) *Envmnt_tp {
  if len(s.Runlv[i].Inpdr) > 0 {
    e.Inpdr = strings.TrimSpace(s.Runlv[i].Inpdr)
  } else {
    e.seekPgmInpdr(s)
  }
  return e
}
func (e *Envmnt_tp) seekRunOutdr(i int, s Settings_tp) *Envmnt_tp {
  if len(s.Runlv[i].Outdr) > 0 {
    e.Outdr = strings.TrimSpace(s.Runlv[i].Outdr)
  } else {
    e.seekPgmOutdr(s)
  }
  return e
}
func (e *Envmnt_tp) seekRunIfilt(i int, s Settings_tp) *Envmnt_tp {
  if len(s.Runlv[i].Ifilt) > 0 {
    e.Ifilt = strings.TrimSpace(s.Runlv[i].Ifilt)
  } else {
    e.seekPgmIfilt(s)
  }
  return e
}
func (e *Envmnt_tp) seekRunIfnam(i int, s Settings_tp) *Envmnt_tp {
  if len(s.Runlv[i].Ifnam) > 0 {
    e.Ifnam = strings.TrimSpace(s.Runlv[i].Ifnam)
  } else {
    e.seekPgmIfnam(s)
  }
  return e
}
func (e *Envmnt_tp) seekRunOfnam(i int, s Settings_tp) *Envmnt_tp {
  if len(s.Runlv[i].Ofnam) > 0 {
    e.Ofnam = strings.TrimSpace(s.Runlv[i].Ofnam)
  } else {
    e.seekPgmOfnam(s)
  }
  return e
}
func (e *Envmnt_tp) seekRunRcvpf(i int, s Settings_tp) *Envmnt_tp {
  if len(s.Runlv[i].Rcvpf) > 0 {
    e.Rcvpf = strings.TrimSpace(s.Runlv[i].Rcvpf)
  } else {
    log.Fatalf("Error: Not possible to determine Receiver Partner Function.\r\n")
  }
  return e
}

func (e *Envmnt_tp) seekPgmObjnm(p Param_tp) *Envmnt_tp {
  if len(p.Prm1) > 0 {
    e.Objnm = strings.TrimSpace(p.Prm1)
  } else {
    log.Fatalf("Error: Not possible to determine IDOC-Type name.\r\n")
  }
  return e
}
func (e *Envmnt_tp) seekPgmDbonm(s Settings_tp) *Envmnt_tp {
  if len(s.Progm.Dbonm) > 0 {
    e.Dbonm = strings.TrimSpace(s.Progm.Dbonm)
  } else {
    log.Fatalf("Error: Not possible to determine Database name.\r\n")
  }
  return e
}
func (e *Envmnt_tp) seekPgmDbodr(s Settings_tp) *Envmnt_tp {
  if len(s.Progm.Dbodr) > 0 {
    e.Dbodr = strings.TrimSpace(s.Progm.Dbodr)
  } else {
    log.Fatalf("Error: Not possible to determine Database directory.\r\n")
  }
  return e
}
func (e *Envmnt_tp) seekPgmInpdr(s Settings_tp) *Envmnt_tp {
  if len(s.Progm.Inpdr) > 0 {
    e.Inpdr = strings.TrimSpace(s.Progm.Inpdr)
  } else {
    log.Fatalf("Error: Not possible to determine Input files directory.\r\n")
  }
  return e
}
func (e *Envmnt_tp) seekPgmOutdr(s Settings_tp) *Envmnt_tp {
  if len(s.Progm.Outdr) > 0 {
    e.Outdr = strings.TrimSpace(s.Progm.Outdr)
  } else {
    log.Fatalf("Error: Not possible to determine Output files directory.\r\n")
  }
  return e
}
func (e *Envmnt_tp) seekPgmIfilt(s Settings_tp) *Envmnt_tp {
  if len(s.Progm.Ifilt) > 0 {
    e.Ifilt = strings.TrimSpace(s.Progm.Ifilt)
  } else {
    log.Fatalf("Error: Not possible to determine Input filter.\r\n")
  }
  return e
}
func (e *Envmnt_tp) seekPgmIfnam(s Settings_tp) *Envmnt_tp {
  if len(s.Progm.Ifnam) > 0 {
    e.Ifnam = strings.TrimSpace(s.Progm.Ifnam)
  } else {
    log.Fatalf("Error: Not possible to determine Input files Naming Rule.\r\n")
  }
  return e
}
func (e *Envmnt_tp) seekPgmOfnam(s Settings_tp) *Envmnt_tp {
  if len(s.Progm.Ofnam) > 0 {
    e.Ofnam = strings.TrimSpace(s.Progm.Ofnam)
  } else {
    log.Fatalf("Error: Not possible to determine Output files Naming Rule.\r\n")
  }
  return e
}

func (e *Envmnt_tp) seekConstCntrl(s Settings_tp) *Envmnt_tp {
  if len(s.Const.Cntrl) > 0 {
    e.Cntrl = strings.TrimSpace(s.Const.Cntrl)
  } else {
    log.Fatalf("Error: Not possible to determine Control constant.\r\n")
  }
  return e
}
func (e *Envmnt_tp) seekConstClien(s Settings_tp) *Envmnt_tp {
  if len(s.Const.Clien) > 0 {
    e.Clien = strings.TrimSpace(s.Const.Clien)
  } else {
    log.Fatalf("Error: Not possible to determine Client constant.\r\n")
  }
  return e
}
