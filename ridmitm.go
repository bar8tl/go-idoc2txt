package idoc2txt

import "strconv"

// Data type to get IDoc items data (records, groups, segments and fields)
// and to create corresponding items records in ref database
type Drmitm_tp struct {
  Icol  []string
  Gcol  []string
  Scol  []string
  Fcol  []string
  Stack []Parsl_tp // List of Parsl_tp: Levels stack
  Lidoc []Idcdf_tp // List of Idcdf_tp: Idoc
  Lgrup []Grpdf_tp // List of Grpdf_tp: Grup
  Lsegm []Segdf_tp // List of Segdf_tp: Segm
  Lfild []Flddf_tp // List of Flddf_tp: Fild
  Lrecd []Flddf_tp // List of Flddf_tp: Fild
  Colsi [2]string  // Name, Extn
  Colsg [5]string  // Numb, Levl, Stat, Mnlp, Mxlp
  Colss [6]string  // Name, Type, Levl, Stat, Mnlp, Mxlp
  Colsf [7]string  // Name, Text, Type, Lgth, Seqn, Strp, Endp
  Colsr [7]string  // Name, Text, Type, Lgth, Seqn, Strp, Endp
  Out   Outsqlt_tp
  L     int        // Stack level
}

func (r *Drmitm_tp) NewDrmitm(s Settings_tp) {
  r.Out.NewOutsqlt(s)
  r.Icol = []string{"EXTENSION"}
  r.Gcol = []string{"LEVEL", "STATUS", "LOOPMIN", "LOOPMAX"}
  r.Scol = []string{"SEGMENTTYPE", "LEVEL", "STATUS", "LOOPMIN", "LOOPMAX"}
  r.Fcol = []string{"NAME", "TEXT", "TYPE", "LENGTH", "FIELD_POS",
    "CHARACTER_FIRST", "CHARACTER_LAST"}
  r.L = -1
}

// Scan SAP parser file to identify IDoc elements
func (r *Drmitm_tp) GetData(sline Parsl_tp) {
  if sline.Label.Ident == "BEGIN" {
    r.L++
    r.Stack = append(r.Stack, Parsl_tp{
      Reclb_tp{sline.Label.Ident, sline.Label.Recnm, sline.Label.Rectp},
      sline.Value})
    if sline.Value != "" {
      if sline.Label.Recnm == "IDOC" {
        r.Colsi[0] = sline.Value
        r.Colsi[1] = sline.Value
        r.Lidoc = append(r.Lidoc, Idcdf_tp{
          r.Colsi[0], r.Stack[r.L].Label.Recnm, r.Colsi})
      } else if sline.Label.Recnm == "GROUP" {
        r.Colsg[0] = sline.Value
      } else if sline.Label.Recnm == "SEGMENT" {
        r.Colss[0] = sline.Value
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
    for i := 0; i < len(r.Icol); i++ {
      if sline.Label.Ident == r.Icol[i] {
        r.Colsi[i+1] = sline.Value
        if i == (len(r.Icol) - 1) {
          r.Lidoc[0].Cols[1] = r.Colsi[i+1]
        }
        break
      }
    }
  }

  if r.Stack[r.L].Label.Recnm == "GROUP" {
    for i := 0; i < len(r.Gcol); i++ {
      if sline.Label.Ident == r.Gcol[i] {
        r.Colsg[i+1] = sline.Value
        if i == (len(r.Gcol) - 1) {
          r.Lgrup = append(r.Lgrup, Grpdf_tp{
            r.Colsg[0], r.Stack[r.L].Label.Recnm, r.Colsg})
        }
        break
      }
    }
  }

  if r.Stack[r.L].Label.Recnm == "SEGMENT" {
    for i := 0; i < len(r.Scol); i++ {
      if sline.Label.Ident == r.Scol[i] {
        r.Colss[i+1] = sline.Value
        if i == (len(r.Scol) - 1) {
          r.Lsegm = append(r.Lsegm, Segdf_tp{
            r.Colss[0], r.Stack[r.L].Label.Recnm, r.Colss})
        }
        break
      }
    }
  }

  if r.Stack[r.L].Label.Recnm == "FIELDS" {
    match := false
    for i := 0; i < len(r.Fcol) && !match; i++ {
      if sline.Label.Ident == r.Fcol[i] {
        r.Colsf[i] = sline.Value
        match = true
      }
      if i == (len(r.Fcol) - 1) {
        if r.Stack[r.L-1].Label.Rectp == "RECORD" {
          r.Lrecd = append(r.Lrecd, Flddf_tp{
            r.Stack[r.L-1].Label.Recnm, r.Stack[r.L].Label.Recnm,
            r.Stack[r.L-1].Label.Rectp, r.Colsf})
        } else if r.Stack[r.L-1].Label.Recnm == "SEGMENT" {
          r.Lfild = append(r.Lfild, Flddf_tp{
            r.Colss[0], r.Stack[r.L].Label.Recnm, r.Stack[r.L-1].Label.Recnm,
            r.Colsf})
        }
      }
    }
  }
}

// Functions to upload IDoc data elements into a reference definition database
func (r *Drmitm_tp) IsrtData(s Settings_tp) {
  r.Out.ClearItems(r.Lidoc[0].Cols[1])
  r.UpldRecd()
  r.UplDidoc()
  r.UplDgrup()
  r.UplDsegm()
  r.UpldFlds()
}

// /RB04/YP3_DELVRY_RBNA|IDOC|DELVRY07|DELVRY07|IDOC|||/RB04/YP3_DELVRY_RBNA|0|
// 0||0|0|0|0|0|0
func (r *Drmitm_tp) UplDidoc() *Drmitm_tp { // Upload IDoc idoc data
  var a Args_tp
  for _, lidoc := range r.Lidoc {
    a.idocn = r.Lidoc[0].Cols[1] // EXTENSION       /RB04/YP3_DELVRY_RBNA
    a.rname = lidoc.Type         // B…_IDOC         IDOC
    a.dname = lidoc.Cols[0]      // BEGIN_IDOC      DELVRY07
    a.rclas = lidoc.Name         // BEGIN_IDOC      DELVRY07
    a.rtype = lidoc.Type         // B…_IDOC         IDOC
    a.dtype = ""
    a.dtext = ""
    a.extsn = lidoc.Cols[1]      // EXTENSION       /RB04/YP3_DELVRY_RBNA
    a.gnumb = 0
    a.level = 0
    a.stats = ""
    a.minlp = 0
    a.maxlp = 0
    a.lngth = 0
    a.seqno = 0
    a.strps = 0
    a.endps = 0
    r.Out.IsrtItems(a)
  }
  return r
}

// /RB04/YP3_DELVRY_RBNA|GROUP|1|1|GROUP||||1|2|MANDATORY|1|9999|0|0|0|0
func (r *Drmitm_tp) UplDgrup() *Drmitm_tp { // Upload IDoc groups data
  var a Args_tp
  for _, lgrup := range r.Lgrup {
    a.idocn = r.Lidoc[0].Cols[1] // EXTENSION       /RB04/YP3_DELVRY_RBNA
    a.rname = lgrup.Type         // B…_GROUP        GROUP
    a.dname = lgrup.Cols[0]      // BEGIN_GROUP     1
    a.rclas = lgrup.Name         // BEGIN_GROUP     1
    a.rtype = lgrup.Type         // B…_GROUP        GROUP
    a.dtype = ""
    a.dtext = ""
    a.extsn = ""
    a.gnumb, _ = strconv.Atoi(lgrup.Cols[0]) // BEGIN_GROUP     1
    a.level, _ = strconv.Atoi(lgrup.Cols[1]) // LEVEL           02
    a.stats = lgrup.Cols[2]                  // STATUS          MANDATORY
    a.minlp, _ = strconv.Atoi(lgrup.Cols[3]) // LOOPMIN         0000000001
    a.maxlp, _ = strconv.Atoi(lgrup.Cols[4]) // LOOPMAX         0000009999
    a.lngth = 0
    a.seqno = 0
    a.strps = 0
    a.endps = 0
    r.Out.IsrtItems(a)
  }
  return r
}

// /RB04/YP3_DELVRY_RBNA|SEGMENT|E2EDL20004|E2EDL20004|SEGMENT|E1EDL20|||0|2|
// MANDATORY|1|1|0|0|0|0
func (r *Drmitm_tp) UplDsegm() *Drmitm_tp { // Upload IDoc segments data
  var a Args_tp
  for _, lsegm := range r.Lsegm {
    a.idocn = r.Lidoc[0].Cols[1] // EXTENSION       /RB04/YP3_DELVRY_RBNA
    a.rname = lsegm.Type         // B…_SEGMENT      SEGMENT
    a.dname = lsegm.Cols[0]      // BEGIN_SEGMENT   E2EDL20004
    a.rclas = lsegm.Name         // BEGIN_SEGMENT   E2EDL20004
    a.rtype = lsegm.Type         // B…_SEGMENT      SEGMENT
    a.dtype = lsegm.Cols[1]      // SEGMENTTYPE     E1EDL20
    a.dtext = ""
    a.extsn = ""
    a.gnumb = 0
    a.level, _ = strconv.Atoi(lsegm.Cols[2]) // LEVEL           02
    a.stats = lsegm.Cols[3]                  // STATUS          MANDATORY
    a.minlp, _ = strconv.Atoi(lsegm.Cols[4]) // LOOPMIN         0000000001
    a.maxlp, _ = strconv.Atoi(lsegm.Cols[5]) // LOOPMAX         0000000001
    a.lngth = 0
    a.seqno = 0
    a.strps = 0
    a.endps = 0
    r.Out.IsrtItems(a)
  }
  return r
}

// /RB04/YP3_DELVRY_RBNA|E2EDL20004|VKBUR|SEGMENT|FIELDS|CHARACTER|Sales Office|
// |0|0||0|0|4|5|84|87
func (r *Drmitm_tp) UpldFlds() *Drmitm_tp { // Upload IDoc fields data
  var a Args_tp
  for _, lfild := range r.Lfild {
    a.idocn = r.Lidoc[0].Cols[1] // EXTENSION       /RB04/YP3_DELVRY_RBNA
    a.rname = lfild.Name         // BEGIN_SEGMENT   E2EDL20004
    a.dname = lfild.Cols[0]      // NAME            VKBUR
    a.rclas = lfild.Clas         // B…_SEGMENT      SEGMENT
    a.rtype = lfild.Type         // B…_FIELDS       FIELDS
    a.dtype = lfild.Cols[2]      // TYPE            CHARACTER
    a.dtext = lfild.Cols[1]      // TEXT            Sales Office
    a.extsn = ""
    a.gnumb = 0
    a.level = 0
    a.stats = ""
    a.minlp = 0
    a.maxlp = 0
    a.lngth, _ = strconv.Atoi(lfild.Cols[3]) // LENGTH          000004
    a.seqno, _ = strconv.Atoi(lfild.Cols[4]) // FIELD_POS       0005
    a.strps, _ = strconv.Atoi(lfild.Cols[5]) // CHARACTER_FIRST 000084
    a.endps, _ = strconv.Atoi(lfild.Cols[6]) // CHARACTER_LAST  000087
    r.Out.IsrtItems(a)
  }
  return r
}

// /RB04/YP3_DELVRY_RBNA|CONTROL|TABNAM|RECORD|FIELDS|CHARACTER|
// Name of Table Structure||0|0||0|0|10|1|1|10
func (r *Drmitm_tp) UpldRecd() *Drmitm_tp { // Upload IDoc records data
  var a Args_tp
  for _, lrecd := range r.Lrecd {
    a.idocn = r.Lidoc[0].Cols[1] // EXTENSION       /RB04/YP3_DELVRY_RBNA
    a.rname = lrecd.Name         // B…_CONTROL_R…   CONTROL
    a.dname = lrecd.Cols[0]      // NAME            TABNAM
    a.rclas = lrecd.Clas         // B…_C…_RECORD    RECORD
    a.rtype = lrecd.Type         // B…_FIELDS       FIELDS
    a.dtype = lrecd.Cols[2]      // TYPE            CHARACTER
    a.dtext = lrecd.Cols[1]      // TEXT            Name of Table Stru...
    a.extsn = ""
    a.gnumb = 0
    a.level = 0
    a.stats = ""
    a.minlp = 0
    a.maxlp = 0
    a.lngth, _ = strconv.Atoi(lrecd.Cols[3]) // LENGTH          000010
    a.seqno, _ = strconv.Atoi(lrecd.Cols[4]) // FIELD_POS       0001
    a.strps, _ = strconv.Atoi(lrecd.Cols[5]) // CHARACTER_FIRST 000001
    a.endps, _ = strconv.Atoi(lrecd.Cols[6]) // CHARACTER_LAST  000010
    r.Out.IsrtItems(a)
  }
  return r
}
