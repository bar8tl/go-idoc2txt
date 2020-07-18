package idoc2txt

import "strconv"

type Drmitm_tp struct {
  Icol, Gcol, Scol, Fcol []string
  Stack []Parsl_tp  // levels stack
  Lrecd []Recdf_tp  // lists
  Lidoc []Idcdf_tp
  Lgrup []Grpdf_tp
  Lsegm []Segdf_tp
  Lfild []Flddf_tp
  Didoc   Idocf_tp
  Dgrup   Grupf_tp
  Dsegm   Segmf_tp
  Dfild   Fildf_tp
  Out     Outsqlt_tp
  L       int       // level
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

func (r *Drmitm_tp) NewDrmitm(s Settings_tp) {
  r.Out.NewOutsqlt(s)
  r.Icol = []string{"EXTENSION"}
  r.Gcol = []string{"LEVEL", "STATUS", "LOOPMIN", "LOOPMAX"}
  r.Scol = []string{"SEGMENTTYPE", "LEVEL", "STATUS", "LOOPMIN", "LOOPMAX"}
  r.Fcol = []string{"NAME", "TEXT", "TYPE", "LENGTH", "FIELD_POS", "CHARACTER_FIRST", "CHARACTER_LAST"}
  r.L    = -1
}

// Functions to get IDoc items data (records, groups, segments and fields) and to create corresponding items records in ref database
func (r *Drmitm_tp) GetData(sline Parsl_tp) { // Scan SAP parser file to identify IDoc elements
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

// Functions to upload IDoc data elements into a reference definition database
type Args_tp struct {
  idocn, rname, rtype, rclas, dname, dtype, dtext, extsn, stats string
  gnumb, level, minlp, maxlp, lngth, seqno, strps, endps        int
}
func (r *Drmitm_tp) IsrtData(s Settings_tp) {
  r.Out.ClearItems(r.Lidoc[0].Idoc.Col[1])
  r.UpldRecd().UplDidoc().UplDgrup().UplDsegm().UpldFlds()
}
func (r *Drmitm_tp) UpldRecd() *Drmitm_tp { // Upload IDoc records data
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
    r.Out.IsrtItems(a)
  }
  return r
}
func (r *Drmitm_tp) UplDidoc() *Drmitm_tp { // Upload IDoc idoc data
  var a Args_tp
  for i := 0; i < len(r.Lidoc); i++ {
    a.idocn = r.Lidoc[0].Idoc.Col[1]
    a.rname = r.Lidoc[i].Type
    a.rtype, a.rclas, a.dname, a.extsn = r.Lidoc[i].Type, r.Lidoc[i].Name, r.Lidoc[i].Idoc.Col[0], r.Lidoc[i].Idoc.Col[1]
    a.dtype, a.dtext, a.stats = "", "", ""
    a.gnumb, a.level, a.minlp, a.maxlp, a.lngth, a.seqno, a.strps, a.endps = 0, 0, 0, 0, 0, 0, 0, 0
    r.Out.IsrtItems(a)
  }
  return r
}
func (r *Drmitm_tp) UplDgrup() *Drmitm_tp { // Upload IDoc groups data
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
    r.Out.IsrtItems(a)
  }
  return r
}
func (r *Drmitm_tp) UplDsegm() *Drmitm_tp { // Upload IDoc segments data
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
    r.Out.IsrtItems(a)
  }
  return r
}
func (r *Drmitm_tp) UpldFlds() *Drmitm_tp { // Upload IDoc fields data
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
    r.Out.IsrtItems(a)
  }
  return r
}
