package idoc2txt

import "strconv"
import "strings"

type Drssgm_tp struct {
  Out          Outsqlt_tp
  Stack        []Keyst_tp
  Tnode, Fnode Keyst_tp
  Snode        Keyst_tp
  Idocn, Strtp string
  L            int
}
func (r *Drssgm_tp) NewDrssgm(s Settings_tp, strtp string) {
  r.Strtp, r.L = strings.ToUpper(strtp), -1
  r.Out.NewOutsqlt(s)
}

func (r *Drssgm_tp) GetData(sline Parsl_tp) {
  if sline.Label.Ident == "BEGIN" {
    if sline.Label.Recnm == "IDOC" {
      r.Stack = append(r.Stack, Keyst_tp{sline.Label.Recnm, sline.Value, 0, 0}); r.L++
      r.Tnode.Rname, r.Tnode.Dname = sline.Label.Recnm, sline.Value
      r.Idocn = sline.Value
      r.Out.ClearStruc(r.Idocn, r.Strtp)
    } else if sline.Label.Recnm == "SEGMENT" && len(sline.Label.Rectp) == 0 {
      r.Tnode.Rname, r.Tnode.Dname = sline.Label.Recnm, sline.Value
    } else if sline.Label.Recnm == "FIELDS" && r.L >= 0 {
      r.Fnode.Rname, r.Fnode.Dname = sline.Label.Recnm, sline.Value
    }
    return
  }
  if sline.Label.Ident == "END" && r.L >= 0 {
    if sline.Label.Recnm == "IDOC" {
      r.Stack = r.Stack[:r.L]; r.L--
    } else if sline.Label.Recnm == "SEGMENT" && len(sline.Label.Rectp) == 0 {
      if r.L == 0 {
        r.Stack[r.L].Seqno += 1
        r.Stack = append(r.Stack, Keyst_tp{r.Tnode.Rname, r.Tnode.Dname, r.Tnode.Level, 0}); r.L++
      } else if r.Tnode.Level <= r.Stack[r.L].Level {
        for r.Tnode.Level <= r.Stack[r.L].Level {
          r.Out.IsrtStruc(r.Idocn, r.Strtp, r.Stack[r.L-1], r.Stack[r.L])
          r.Stack = r.Stack[:r.L]; r.L--
        }
        r.Stack[r.L].Seqno += 1
        r.Stack = append(r.Stack, Keyst_tp{r.Tnode.Rname, r.Tnode.Dname, r.Tnode.Level, 0}); r.L++
      } else if r.Tnode.Level > r.Stack[r.L].Level {
        r.Stack[r.L].Seqno += 1
        r.Stack = append(r.Stack, Keyst_tp{r.Tnode.Rname, r.Tnode.Dname, r.Tnode.Level, 0}); r.L++
      }
    } else if sline.Label.Recnm == "FIELDS" && r.L >= 0 {
      r.Fnode.Rname, r.Fnode.Dname = "", "" 
    }
    return
  }
  if r.Fnode.Rname == "FIELDS" && r.Tnode.Rname == "SEGMENT" {
    if sline.Label.Ident == "NAME" {
      r.Snode.Rname, r.Snode.Dname, r.Snode.Level, r.Snode.Seqno = r.Tnode.Dname, sline.Value, 0, 0
    } else if sline.Label.Ident == "FIELD_POS" {
      pos, _ := strconv.Atoi(sline.Value)
      r.Tnode.Seqno = pos
      r.Out.IsrtStruc(r.Idocn, r.Strtp, r.Tnode, r.Snode)
    }
    return
  }
  if r.Tnode.Rname == "SEGMENT" && len(r.Tnode.Dname) > 0 {
    if sline.Label.Ident == "LEVEL" {
      l, _ := strconv.Atoi(sline.Value)
      r.Tnode.Level = l
    }
    return
  }
  if r.Tnode.Rname == "IDOC" {
    if sline.Label.Ident == "EXTENSION" {
      r.Idocn = sline.Value
      r.Out.ClearStruc(r.Idocn, r.Strtp)
    }
    return
  }
}
