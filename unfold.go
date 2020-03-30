package idoc2txt

import "bufio"
import "code.google.com/p/go-sqlite/go1/sqlite3"
import "fmt"
import "io/ioutil"
import "log"
import "io"
import "os"
import "strconv"
import "strings"

// Data type to convert IDoc classic hierarchical format to flat text file format
type Dunf_tp struct {
  Cnnsq, Cnnst        string
  Dbonm, Dbodr        string
  Objnm, Inpdr, Outdr string
  Ifilt, Ifnam, Ofnam string
  Cntrl, Clien, Rcvpf string
  Idocx               string
  Idocn, Idocb        string
  Sectn, Secnb        string
  Sgnum, Sgnam, Sgdsc string
  Sgnbk, Sghnb, Sglvl string
  Serie               string
  Nsegm               int
  Dirty               bool
  Lctrl          [524]byte
  Lsegm         [1063]byte
  Lstat          [562]byte
  Ifile, Ofile       *os.File
  Db                 *sqlite3.Conn
  Parnt             []Hstruc_tp
  L                   int
}
type Hstruc_tp struct {
  Sgnum, Sgnam, Sglvl string
}

// Constructor of object Dunf: Define input/output file and database location folders, database full connection string as well
func NewDunf(parm Param_tp, s Settings_tp) *Dunf_tp {
  var u Dunf_tp
  s.SetRunVars(parm, s)
  u.Cnnsq, u.Cnnst = s.Cnnsq, s.Cnnst
  u.Dbonm, u.Dbodr = s.Dbonm, s.Dbodr
  u.Objnm          = s.Objnm
  u.Inpdr, u.Outdr = s.Inpdr, s.Outdr
  u.Ifilt          = s.Ifilt
  u.Ifnam, u.Ofnam = s.Ifnam, s.Ofnam
  u.Cntrl, u.Clien = s.Cntrl, s.Clien
  u.Rcvpf          = s.Rcvpf
  u.Idocx          = strings.ToUpper(s.Objnm)
  return &u
}

// Public option UNF: Unfold data IDocs based on specific IDoc-type. Produces system readeable flat text files
func (u *Dunf_tp) UnfoldIdocs(s Settings_tp) {
  files, _ := ioutil.ReadDir(u.Inpdr)
  for _, f := range files {
    if len(u.Ifilt) == 0 || (len(u.Ifilt) > 0 && PassFilter(s, f)) {
      u.ProcDataLines(f)
    }
  }
}

// Function to process IDoc data files, reading line by line and determining measures for format conversion
func (u *Dunf_tp) ProcDataLines(f os.FileInfo) {
  u.OpenProgStreams(f).DetermIdocProps()
  u.Idocn, u.Nsegm, u.L  = "", 0, -1
  u.Parnt = u.Parnt[:u.L+1]
  rdr := bufio.NewReader(u.Ifile)
  wtr := bufio.NewWriter(u.Ofile)
  for l, err := rdr.ReadString(byte('\n')); err != io.EOF; l, err = rdr.ReadString(byte('\n')) {
    l = strings.TrimSpace(l)
    t := strings.Split(l, "\t")
    if len(l) == 0 { // ignores lines in blank
      continue
    }
    if len(u.Idocn) == 0 && len(t) == 1 && l[0:11] == "IDoc Number" { // gets IDoc number
      i := strings.Split(l, " : ")
      u.Idocn = strings.TrimSpace(i[1])
      continue
    }
    if len(t) <= 1 { // ignores lines no containing tabulators (after to have gotten IDoc number)
      continue
    }
    if t[0] == "EDIDC" || t[0] == "EDIDD" || t[0] == "EDIDS" { // determines data section to analyze
      u.SetupSection(t, wtr)
      continue
    }
    if t[0] == "SEGNUM" && len(t) == 3 { // check in segment number to analize
      u.Sgnbk = u.Sgnum
      u.Sgnum = t[2]
      continue
    }
    if t[0] == "SEGNAM" && len(t) == 3 { // check in segment name to analize
      u.SetupSegment(t, wtr)
      continue
    }
    // process fields of each data section
    if u.Sectn == "EDIDC" {
      u.procEdidc(t)
    } else if u.Sectn == "EDIDD" {
      u.procEdidd(u.Sgnum, u.Sgnam, t)
    } else if u.Sectn == "EDIDS" {
      u.procEdids(u.Secnb, t)
    }
  }
  u.CloseProgStreams(f)
}

// Function to setup measures to take for each data section. Each new section causes dumping data from previous one
func (u *Dunf_tp) SetupSection(t []string, wtr *bufio.Writer) {
  u.Sectn = t[0]
  if u.Sectn == "EDIDC" {
    for i := 0; i < len(u.Lctrl); i++ {
      u.Lctrl[i] = ' '
    }
  }
  if u.Sectn == "EDIDD" {
    u.DumpControlLine(wtr)
  }
  if u.Sectn == "EDIDS" {
    u.Sgnbk = u.Sgnum
    u.DumpSegmentLine(wtr)
    for i := 0; i < len(u.Lstat); i++ {
      u.Lstat[i] = ' '
    }
    if len(t) == 3 {
      u.Secnb = t[2]
    }
  }
}

// Function to setup measures to take for each data segment in Data Idoc being converted
func (u *Dunf_tp) SetupSegment(t []string, wtr *bufio.Writer) {
  u.Nsegm++
  if u.Nsegm > 1 {
    u.DumpSegmentLine(wtr)
  }
  u.Sgnam = t[2]
  for i := 0; i < len(u.Lsegm); i++ {
    u.Lsegm[i] = ' '
  }
  rdb, err := u.Db.Query(`SELECT dname, level FROM items WHERE idocn=? and rname=? and dtype=?;`, u.Idocx, "SEGMENT", u.Sgnam)
  if err != nil {
    log.Fatalf("Error during searching segment description: %s %s\r\n", u.Sgnam, err)
  }
  var level int
  rdb.Scan(&u.Sgdsc, &level)
  u.Sglvl = fmt.Sprintf("%02d", level)

  if u.Nsegm == 1 {
    u.Parnt = append(u.Parnt, Hstruc_tp{u.Sgnum, u.Sgnam, u.Sglvl})
    u.L++
    u.Sghnb = "000000"
  } else {
    if u.Sglvl > u.Parnt[u.L].Sglvl {
      u.Parnt = append(u.Parnt, Hstruc_tp{u.Sgnum, u.Sgnam, u.Sglvl})
      u.L++
      u.Sghnb = u.Parnt[u.L-1].Sgnum
    } else if u.Sglvl == u.Parnt[u.L].Sglvl {
      u.Parnt[u.L].Sgnum, u.Parnt[u.L].Sgnam, u.Parnt[u.L].Sglvl = u.Sgnum, u.Sgnam, u.Sglvl
      u.Sghnb = u.Parnt[u.L-1].Sgnum
    } else {
      prvlv, _ := strconv.Atoi(u.Parnt[u.L].Sglvl)
      curlv, _ := strconv.Atoi(u.Sglvl)
      nstep := prvlv - curlv
      for i := 1; i <= nstep; i++ {
        u.L--
        u.Parnt = u.Parnt[:u.L+1]
      }
      u.Parnt[u.L].Sgnum, u.Parnt[u.L].Sgnam, u.Parnt[u.L].Sglvl = u.Sgnum, u.Sgnam, u.Sglvl
      u.Sghnb = u.Parnt[u.L-1].Sgnum
    }
  }
  rdb.Close()
}

// Functions to process format conversion to fields in control record
func (u *Dunf_tp) procEdidc(t []string) {
  flkey := t[0]
  if flkey == "RVCPRN" {
    flkey = "RCVPRN"
  }
  flval := ""
  if len(t) == 3 {
    f := strings.Split(t[2], " :")
    flval = strings.TrimSpace(f[0])
  }
  if flkey == "CREDAT" {
    u.Serie = flval
  }
  if flkey == "CRETIM" {
    u.Serie += flval
  }
  if len(flval) > 0 {
    u.Dirty = true
    u.SetControlField(flkey, flval)
  }
}
func (u *Dunf_tp) SetControlField(flkey, flval string) {
  var strps int
  rdb, err := u.Db.Query(`SELECT strps FROM items WHERE idocn=? and rname=? and dname=?;`, u.Idocx, "CONTROL", flkey)
  if err != nil {
    log.Fatalf("Error during reading database for control data: %v\r\n", err)
  }
  rdb.Scan(&strps)
  rdb.Close()
  if flkey == "IDOCTYP" && flval == "14" {
    flval = u.Idocb
  }
  if flkey == "CIMTYP" && flval == "14" {
    flval = u.Idocx
  }
  k := strps - 1
  for i := 0; i < len(flval); i++ {
    u.Lctrl[k] = flval[i]
    k++
  }
}
func (u *Dunf_tp) DumpControlLine(wtr *bufio.Writer) {
  if u.Dirty {
    u.SetControlField("TABNAM", u.Cntrl)
    u.SetControlField("MANDT",  u.Clien)
    u.SetControlField("DOCNUM", u.Idocn)
    u.SetControlField("RCVPFC", u.Rcvpf)
    u.SetControlField("SERIAL", u.Serie)
    fmt.Fprintf(wtr, "%s\r\n",  u.Lctrl)
    wtr.Flush()
    u.Dirty = false
  }
}

// Functions to process format conversion to fields in data records
func (u *Dunf_tp) procEdidd(sgnum, sgnam string, t []string) {
  flkey := t[0]
  flval := ""
  if len(t) == 3 {
    f := strings.Split(t[2], " :")
    flval = strings.TrimSpace(f[0])
  }
  if len(flval) > 0 {
    u.Dirty = true
    u.SetSegmentField(u.Sgdsc, flkey, flval)
  }
}
func (u *Dunf_tp) SetSegmentField(sgdsc, flkey, flval string) {
  var strps int
  rdb, err := u.Db.Query(`SELECT strps FROM items WHERE idocn=? and rname=? and dname=?;`, u.Idocx, sgdsc, flkey)
  if err != nil {
    log.Fatalf("Error during reading database for segment data: %s %s %s %v\r\n", u.Idocx, u.Sgdsc, flkey, err)
  }
  rdb.Scan(&strps)
  rdb.Close()
  k := strps - 1
  for i := 0; i < len(flval); i++ {
    u.Lsegm[k] = flval[i]
    k++
  }
}
func (u *Dunf_tp) DumpSegmentLine(wtr *bufio.Writer) {
  if u.Dirty {
    u.SetSegmentField("DATA", "SEGNAM", u.Sgdsc)
    u.SetSegmentField("DATA", "MANDT",  u.Clien)
    u.SetSegmentField("DATA", "DOCNUM", u.Idocn)
    u.SetSegmentField("DATA", "SEGNUM", u.Sgnbk)
    u.SetSegmentField("DATA", "PSGNUM", u.Sghnb)
    u.SetSegmentField("DATA", "HLEVEL", u.Sglvl)
    fmt.Fprintf(wtr, "%s\r\n", u.Lsegm)
    wtr.Flush()
    u.Dirty = false
  }
}

func (u *Dunf_tp) procEdids(secnb string, t []string) {}

func (u *Dunf_tp) OpenProgStreams(f os.FileInfo) *Dunf_tp {
  var err error
  u.Db, err = sqlite3.Open(u.Cnnst)
  if err != nil {
    log.Fatalf("Open SQLite database error: %s\n", err)
  }
  u.Ifile, err = os.Open(u.Inpdr + f.Name())
  if err != nil {
    log.Fatalf("Input file not found: %s\r\n", err)
  }
  u.Ofile, err = os.Create(u.Outdr + f.Name())
  if err != nil {
    log.Fatalf("Error during output file creation: %s\r\n", err)
  }
  return u
}
func (u *Dunf_tp) CloseProgStreams(f os.FileInfo) *Dunf_tp {
  u.Ifile.Close()
  u.Ofile.Close()
  u.Db.Close()
  RanameInpFile(u.Inpdr, f)
  RanameOutFile(u.Outdr, f)
  return u
}
func (u *Dunf_tp) DetermIdocProps() *Dunf_tp {
  rdb, err := u.Db.Query("SELECT dname FROM items WHERE idocn=? and rname=?;", u.Idocx, "IDOC")
  if err != nil {
    log.Fatalf("Error during searching Idoc properties: %s %v\r\n", u.Idocx, err)
  }
  rdb.Scan(&u.Idocb)
  rdb.Close()
  return u
}
