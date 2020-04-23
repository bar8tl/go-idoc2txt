package idoc2txt

import "bufio"
import "code.google.com/p/go-sqlite/go1/sqlite3"
import "io"
import "log"
import "os"
import "path/filepath"
import "strings"
import "strconv"

// Data type to read SAP IDoc parser file and to upload IDoc definition detail and structure into an internal reference database
type Drid_tp struct {
  Idocf, Idocn, Idoct         string
  Icol,  Gcol,  Scol,  Fcol []string
  Stack []Parsl_tp
  Lrecd []Recdf_tp
  Lidoc []Idcdf_tp
  Lgrup []Grpdf_tp
  Lsegm []Segdf_tp
  Lfild []Flddf_tp
  Didoc   Idocf_tp
  Dgrup   Grupf_tp
  Dsegm   Segmf_tp
  Dfild   Fildf_tp
  L       int
}
type Parsl_tp struct {
  Label Reclb_tp
  Value string
}
type Reclb_tp struct {
  Ident, Recnm, Rectp string
}
type Recdf_tp struct {
  Name, Type, Clas string
  Flds             Fildf_tp
}
type Idcdf_tp struct {
  Name, Type string
  Idoc       Idocf_tp
}
type Grpdf_tp struct {
  Name, Type string
  Grup       Grupf_tp
}
type Segdf_tp struct {
  Name, Type string
  Segm       Segmf_tp
}
type Flddf_tp struct {
  Name, Type, Clas string
  Flds             Fildf_tp
}
type Idocf_tp struct {
  Col [2]string // Name, Extsn
}
type Grupf_tp struct {
  Col [5]string // Numbr, Level, Stats, Minlp, Maxlp
}
type Segmf_tp struct {
  Col [6]string // Name, Type, Level, Stats, Minlp, Maxlp
}
type Fildf_tp struct {
  Col [7]string // Name, Text, Type, Lngth, Seqno, Strps, Endps
}

// Constructor of object Drid: Define input file and database name and location folders, SQlite3 database full connection string as well
func NewDrid() *Drid_tp {
  var r Drid_tp
  return &r
}

// Public option RID: Read and upload IDoc-definition text file (commonly known as IDoc Parser-File in SAP Transaction Code WE60)
func (r *Drid_tp) ProcInput(parm Param_tp, s Settings_tp) {
  s.SetRunVars(parm, s)
  r.Idocf          = s.Objnm
  r.Idoct, r.Idocn = filepath.Ext(r.Idocf), strings.TrimRight(r.Idocf, r.Idoct)
  r.Idocf = s.Inpdr + r.Idocf
  r.Icol  = []string{"EXTENSION"}
  r.Gcol  = []string{"LEVEL", "STATUS", "LOOPMIN", "LOOPMAX"}
  r.Scol  = []string{"SEGMENTTYPE", "LEVEL", "STATUS", "LOOPMIN", "LOOPMAX"}
  r.Fcol  = []string{"NAME", "TEXT", "TYPE", "LENGTH", "FIELD_POS", "CHARACTER_FIRST", "CHARACTER_LAST"}
  r.L     = -1
  ifile, err := os.Open(r.Idocf)
  if err != nil {
    log.Fatalf("Input file %s not found: %s\r\n", r.Idocf, err)
  }
  defer ifile.Close()
  rdr := bufio.NewReader(ifile)
  r.ProcStartOfFile()
  for line, _, err := rdr.ReadLine(); err != io.EOF; line, _, err = rdr.ReadLine() {
    if l := strings.Trim(string(line), " "); len(l) > 0 {
      r.ProcLines(string(l))
    }
  }
  r.ProcEndOfFile(s)
}

func (r *Drid_tp) ProcStartOfFile() {}

// Functions to get IDoc items data (records, groups, segments and fields) and to create corresponding items records in ref database
func (r *Drid_tp) ProcLines(l string) {
  sline := r.ScanTextIdocLine(l)
  r.IsrtMasterData(sline)
  r.IsrtStructData(sline)
}
func (r *Drid_tp) IsrtMasterData(sline Parsl_tp) { // Scan SAP parser file to identify IDoc elements
  if sline.Label.Ident == "BEGIN" {
    r.L++
    r.Stack = append(r.Stack, Parsl_tp{Reclb_tp{sline.Label.Ident, sline.Label.Recnm, sline.Label.Rectp}, sline.Value})
    if sline.Value != "" {
      if sline.Label.Recnm == "IDOC" {
        r.Didoc.Col[0], r.Didoc.Col[1] = sline.Value, sline.Value
        r.Lidoc = append(r.Lidoc, Idcdf_tp{r.Didoc.Col[0], r.Stack[r.L].Label.Recnm, r.Didoc})
      } else if sline.Label.Recnm == "GROUP" {
        r.Dgrup.Col[0] = sline.Value
      } else if sline.Label.Recnm == "SEGMENT" {
        r.Dsegm.Col[0] = sline.Value
      }
    }
    return
  }
  if sline.Label.Ident == "END" {
    r.L--
    r.Stack = r.Stack[:r.L+1]
    return
  }
  if r.Stack[r.L].Label.Recnm == "IDOC" {
    match := false
    for i := 0; i < len(r.Icol) && !match; i++ {
      if sline.Label.Ident == r.Icol[i] {
        r.Didoc.Col[i+1] = sline.Value
        match = true
        if i == (len(r.Icol) - 1) {
          r.Lidoc[0].Idoc.Col[1] = r.Didoc.Col[i+1]
        }
      }
    }
  }
  if r.Stack[r.L].Label.Recnm == "GROUP" {
    match := false
    for i := 0; i < len(r.Gcol) && !match; i++ {
      if sline.Label.Ident == r.Gcol[i] {
        r.Dgrup.Col[i+1] = sline.Value
        match = true
        if i == (len(r.Gcol) - 1) {
          r.Lgrup = append(r.Lgrup, Grpdf_tp{r.Dgrup.Col[0], r.Stack[r.L].Label.Recnm, r.Dgrup})
        }
      }
    }
  }
  if r.Stack[r.L].Label.Recnm == "SEGMENT" {
    match := false
    for i := 0; i < len(r.Scol) && !match; i++ {
      if sline.Label.Ident == r.Scol[i] {
        r.Dsegm.Col[i+1] = sline.Value
        match = true
        if i == (len(r.Scol) - 1) {
          r.Lsegm = append(r.Lsegm, Segdf_tp{r.Dsegm.Col[0], r.Stack[r.L].Label.Recnm, r.Dsegm})
        }
      }
    }
  }
  if r.Stack[r.L].Label.Recnm == "FIELDS" {
    match := false
    for i := 0; i < len(r.Fcol) && !match; i++ {
      if sline.Label.Ident == r.Fcol[i] {
        r.Dfild.Col[i] = sline.Value
        match = true
      }
      if i == (len(r.Fcol) - 1) {
        if r.Stack[r.L-1].Label.Rectp == "RECORD" {
          r.Lrecd = append(r.Lrecd, Recdf_tp{r.Stack[r.L-1].Label.Recnm, r.Stack[r.L].Label.Recnm, r.Stack[r.L-1].Label.Rectp, r.Dfild})
        } else if r.Stack[r.L-1].Label.Recnm == "SEGMENT" {
          r.Lfild = append(r.Lfild, Flddf_tp{r.Dsegm.Col[0], r.Stack[r.L].Label.Recnm, r.Stack[r.L-1].Label.Recnm, r.Dfild})
        }
      }
    }
  }
}

func (r *Drid_tp) IsrtStructData(sline Parsl_tp) {}

// Functions to upload IDoc data elements into a reference definition database
type Args_tp struct {
  idocn, rname, rtype, rclas, dname, dtype, dtext, extsn, stats string
  gnumb, level, minlp, maxlp, lngth, seqno, strps, endps        int
}

func (r *Drid_tp) ProcEndOfFile(s Settings_tp) {
  db, err := sqlite3.Open(s.Cnnst)
  if err != nil {
    log.Fatalf("Open SQLite database error: %v\n", err)
  }
  db.Exec(`DELETE FROM items WHERE idocn=?;`, r.Lidoc[0].Idoc.Col[1])
  r.UpldRecd(s, db).UplDidoc(s, db).UplDgrup(s, db).UplDsegm(s, db).UpldFlds(s, db)
}

func (r *Drid_tp) UpldRecd(s Settings_tp, db *sqlite3.Conn) *Drid_tp { // Upload IDoc records data
  var a Args_tp
  for i := 0; i < len(r.Lrecd); i++ {
    a.idocn = r.Lidoc[0].Idoc.Col[1]
    a.rname = r.Lrecd[i].Name
    a.rtype, a.rclas, a.dname, a.extsn = r.Lrecd[i].Type, r.Lrecd[i].Clas, r.Lrecd[i].Flds.Col[0], ""
    a.dtype, a.dtext, a.stats = r.Lrecd[i].Flds.Col[2], r.Lrecd[i].Flds.Col[1], ""
    a.gnumb, a.level, a.minlp, a.maxlp = 0, 0, 0, 0
    a.lngth, _ = strconv.Atoi(r.Lrecd[i].Flds.Col[3])
    a.seqno, _ = strconv.Atoi(r.Lrecd[i].Flds.Col[4])
    a.strps, _ = strconv.Atoi(r.Lrecd[i].Flds.Col[5])
    a.endps, _ = strconv.Atoi(r.Lrecd[i].Flds.Col[6])
    r.WriteItems(a, db)
  }
  return r
}

func (r *Drid_tp) UplDidoc(s Settings_tp, db *sqlite3.Conn) *Drid_tp { // Upload IDoc idoc data
  var a Args_tp
  for i := 0; i < len(r.Lidoc); i++ {
    a.idocn = r.Lidoc[0].Idoc.Col[1]
    a.rname = r.Lidoc[i].Type
    a.rtype, a.rclas, a.dname, a.extsn = r.Lidoc[i].Type, r.Lidoc[i].Name, r.Lidoc[i].Idoc.Col[0], r.Lidoc[i].Idoc.Col[1]
    a.dtype, a.dtext, a.stats = "", "", ""
    a.gnumb, a.level, a.minlp, a.maxlp, a.lngth, a.seqno, a.strps, a.endps = 0, 0, 0, 0, 0, 0, 0, 0
    r.WriteItems(a, db)
  }
  return r
}

func (r *Drid_tp) UplDgrup(s Settings_tp, db *sqlite3.Conn) *Drid_tp { // Upload IDoc groups data
  var a Args_tp
  for i := 0; i < len(r.Lgrup); i++ {
    a.idocn = r.Lidoc[0].Idoc.Col[1]
    a.rname = r.Lgrup[i].Type
    a.rtype, a.rclas, a.dname, a.extsn = r.Lgrup[i].Type, r.Lgrup[i].Name, r.Lgrup[i].Grup.Col[0], ""
    a.dtype, a.dtext, a.stats = "", "", r.Lgrup[i].Grup.Col[2]
    a.gnumb, _ = strconv.Atoi(r.Lgrup[i].Grup.Col[0])
    a.level, _ = strconv.Atoi(r.Lgrup[i].Grup.Col[1])
    a.minlp, _ = strconv.Atoi(r.Lgrup[i].Grup.Col[3])
    a.maxlp, _ = strconv.Atoi(r.Lgrup[i].Grup.Col[4])
    a.lngth, a.seqno, a.strps, a.endps = 0, 0, 0, 0
    r.WriteItems(a, db)
  }
  return r
}

func (r *Drid_tp) UplDsegm(s Settings_tp, db *sqlite3.Conn) *Drid_tp { // Upload IDoc segments data
  var a Args_tp
  for i := 0; i < len(r.Lsegm); i++ {
    a.idocn = r.Lidoc[0].Idoc.Col[1]
    a.rname = r.Lsegm[i].Type
    a.rtype, a.rclas, a.dname, a.extsn = r.Lsegm[i].Type, r.Lsegm[i].Name, r.Lsegm[i].Segm.Col[0], ""
    a.dtype, a.dtext, a.stats = r.Lsegm[i].Segm.Col[1], "", r.Lsegm[i].Segm.Col[3]
    a.level, _ = strconv.Atoi(r.Lsegm[i].Segm.Col[2])
    a.minlp, _ = strconv.Atoi(r.Lsegm[i].Segm.Col[4])
    a.maxlp, _ = strconv.Atoi(r.Lsegm[i].Segm.Col[5])
    a.gnumb, a.lngth, a.seqno, a.strps, a.endps = 0, 0, 0, 0, 0
    r.WriteItems(a, db)
  }
  return r
}

func (r *Drid_tp) UpldFlds(s Settings_tp, db *sqlite3.Conn) *Drid_tp { // Upload IDoc fields data
  var a Args_tp
  for i := 0; i < len(r.Lfild); i++ {
    a.idocn = r.Lidoc[0].Idoc.Col[1]
    a.rname = r.Lfild[i].Name
    a.rtype, a.rclas, a.dname, a.extsn = r.Lfild[i].Type, r.Lfild[i].Clas, r.Lfild[i].Flds.Col[0], ""
    a.dtype, a.dtext, a.stats = r.Lfild[i].Flds.Col[2], r.Lfild[i].Flds.Col[1], ""
    a.gnumb, a.level, a.minlp, a.maxlp = 0, 0, 0, 0
    a.lngth, _ = strconv.Atoi(r.Lfild[i].Flds.Col[3])
    a.seqno, _ = strconv.Atoi(r.Lfild[i].Flds.Col[4])
    a.strps, _ = strconv.Atoi(r.Lfild[i].Flds.Col[5])
    a.endps, _ = strconv.Atoi(r.Lfild[i].Flds.Col[6])
    r.WriteItems(a, db)
  }
  return r
}

func (r *Drid_tp) WriteItems(a Args_tp, db *sqlite3.Conn) { // Fetch data of IDoc elements ino reference database
  args := sqlite3.NamedArgs{
    ":01": a.idocn, ":02": a.rname, ":03": a.dname, ":04": a.rclas, ":05": a.rtype, ":06": a.dtype, ":07": a.dtext, ":08": a.extsn,
    ":09": a.gnumb, ":10": a.level, ":11": a.stats, ":12": a.minlp, ":13": a.maxlp, ":14": a.lngth, ":15": a.seqno, ":16": a.strps,
    ":17": a.endps,
  }
  err := db.Exec(`
    INSERT INTO items VALUES(:01,:02,:03,:04,:05,:06,:07,:08,:09,:10,:11,:12,:13,:14,:15,:16,:17)`, args)
  if err != nil {
    log.Fatalf("Insert items sql table error: %v\n", err)
  }
}

// Function to identify individual tokens in SAP IDoc parser file
func (r *Drid_tp) ScanTextIdocLine(s string) (p Parsl_tp) {
  var key string
  var val string
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
