package idoc2txt

import "strings"

// Data type to get IDoc groups structure data and to create corresponding
// structure records in ref database
type Drsgrp_tp struct {
  Out   Outsqlt_tp
  Stack []Keyst_tp
  Idocn string
  Strtp string
  L     int
}

func (r *Drsgrp_tp) NewDrsgrp(s Settings_tp, strtp string) {
  r.Strtp = strings.ToUpper(strtp)
  r.L = -1
  r.Out.NewOutsqlt(s)
}

func (r *Drsgrp_tp) GetData(sline Parsl_tp) {
  if sline.Label.Ident == "BEGIN" {
    if sline.Label.Recnm == "IDOC" {
      r.Stack = append(r.Stack, Keyst_tp{sline.Label.Recnm, sline.Value, 0, 0})
      r.L++
      r.Idocn = sline.Value
      r.Out.ClearStruc(r.Idocn, r.Strtp)
    } else if sline.Label.Recnm == "GROUP" {
      r.Stack[r.L].Seqno += 1
      r.Stack = append(r.Stack, Keyst_tp{sline.Label.Recnm, sline.Value, 0, 0})
      r.L++
    }
    return
  }
  if sline.Label.Ident == "END" {
    if sline.Label.Recnm == "IDOC" {
      r.Stack = r.Stack[:r.L]
      r.L--
    } else if sline.Label.Recnm == "GROUP" {
      r.Out.IsrtStruc(r.Idocn, r.Strtp, r.Stack[r.L-1], r.Stack[r.L])
      r.Stack = r.Stack[:r.L]
      r.L--
    }
    return
  }
  if r.L >= 0 && r.Stack[r.L].Rname == "IDOC" {
    if sline.Label.Ident == "EXTENSION" {
      r.Idocn = sline.Value
      r.Out.ClearStruc(r.Idocn, r.Strtp)
    }
    return
  }
}
