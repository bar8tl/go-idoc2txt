package idoc2txt

import "code.google.com/p/go-sqlite/go1/sqlite3"
import "log"

const SELSTRUC =
 `SELECT * FROM struc WHERE idocn=? and strtp=? and prnam=? and pdnam=?;`

type Ddoc_tp struct {
  Db    *sqlite3.Conn
  Rdb   [9]*sqlite3.Stmt
  Cnnst string
  Idocn string
  Strtp string
  Prnam string
  Pdnam string
}

func NewDdoc() *Ddoc_tp {
  var d Ddoc_tp
  d.Idocn = "/RB04/YP3_DELVRY_RBNA"
  d.Strtp = "SGM"
  d.Prnam = "IDOC"
  d.Pdnam = "DELVRY07"
  return &d
}

func (d *Ddoc_tp) ProcDocument(parm Param_tp, s Settings_tp) {
  s.SetRunVars(parm)
  d.Cnnst = s.Cnnst
  var err error
  d.Db, err = sqlite3.Open(d.Cnnst)
  if err != nil {
    log.Fatalf("Open SQLite database error: %v\n", err)
  }
  d.BwseStruc(d.Idocn, d.Strtp, d.Prnam, d.Pdnam, -1)
}

func (d *Ddoc_tp) BwseStruc(Idocn, Strtp, Prnam, Pdnam string, i int) {
  var err error
  i++
  for d.Rdb[i], err = d.Db.Query(SELSTRUC, Idocn, Strtp, Prnam, Pdnam);
    err == nil; err = d.Rdb[i].Next() {
    var idocn string
    var strtp string
    var prnam string
    var pdnam string
    var seqno string
    var crnam string
    var cdnam string
    d.Rdb[i].Scan(&idocn, &strtp, &prnam, &pdnam, &seqno, &crnam, &cdnam)
    log.Printf("%d %s %s %s %s %s\r\n", i, prnam, pdnam, seqno, crnam, cdnam)
    d.BwseStruc(idocn, strtp, crnam, cdnam, i)
  }
}
