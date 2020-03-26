package idoc2txt

import "encoding/xml"
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

// Data type to get an xml structure from configuration file to define potential program-level and run-level parameters
type Config_tp struct {
  XMLName xml.Name   `xml:"config"`
  Progm   Program_tp `xml:"program"`
  Runlv   []Run_tp   `xml:"run"`
}
type Program_tp struct {
  Dbnam string `xml:"dbName,attr"`
  Dbodr string `xml:"dboDir,attr"`
  Inpdr string `xml:"inpDir,attr"`
  Outdr string `xml:"outDir,attr"`
}
type Run_tp struct {
  Optcd string `xml:"option,attr"`
  Objnm string `xml:"objName,attr"`
  Dbodr string `xml:"dboDir,attr"`
  Inpdr string `xml:"inpDir,attr"`
  Outdr string `xml:"outDir,attr"`
  Ifilt string `xml:"inpFilter,attr"`
  Ifnam string `xml:"inpName,attr"`
  Ofnam string `xml:"outName,attr"`
}
func (c *Config_tp) NewConfig(fname string) {
  f, err := os.Open(fname)
  if err != nil {
    log.Fatalf("File config.xml opening error: %s\n", err)
  } else {
    defer f.Close()
    xmlv, _ := ioutil.ReadAll(f)
    err = xml.Unmarshal(xmlv, &c)
    if err != nil {
      log.Fatalf("File config.xml reading error: %s\n", err)
    }
  }
}

// Data type to define global environment variables
type Envmnt_tp struct {
  Cnnsq, Cnnst, Idocf string    // SQlite database connection string
  Dtsys, Dtcur, Dtnul time.Time // System, Current and Null datetime
}
func (e *Envmnt_tp) NewEnvmnt(progm Program_tp) {
  e.Cnnsq = "file:@?file:locked.sqlite?cache=shared&mode=rwc"
  e.Cnnst = strings.Replace(e.Cnnsq, "@", progm.Dbodr+progm.Dbnam, 1)
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
  s.NewEnvmnt(s.Progm)
  return s
}
