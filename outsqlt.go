package idoc2txt

import "code.google.com/p/go-sqlite/go1/sqlite3"
import "fmt"
import "log"
import "strconv"

// Addressing outputs to sqlite3 database
type Outsqlt_tp struct {
  Db    *sqlite3.Conn
  Cnnst string
}

func (o *Outsqlt_tp) NewOutsqlt(s Settings_tp) {
  o.Cnnst = s.Cnnst
}

func (o *Outsqlt_tp) ClearItems(idocn string) {
  var err error
  o.Db, err = sqlite3.Open(o.Cnnst)
  if err != nil {
    log.Fatalf("Open SQLite database error: %v\n", err)
  }
  o.Db.Exec(`DELETE FROM items WHERE idocn=?;`, idocn)
}

func (o *Outsqlt_tp) IsrtItems(a Args_tp) {
  err := o.Db.Exec(`INSERT INTO items VALUES(?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)`,
    a.idocn, a.rname, a.dname, a.rclas, a.rtype, a.dtype, a.dtext, a.extsn,
    a.gnumb, a.level, a.stats, a.minlp, a.maxlp, a.lngth, a.seqno, a.strps,
    a.endps)
  if err != nil {
    log.Fatalf("Insert items sql table error: %v\n", err)
  }
}

func (o *Outsqlt_tp) ClearStruc(idocn, strtp string) {
  var err error
  o.Db, err = sqlite3.Open(o.Cnnst)
  if err != nil {
    log.Fatalf("Open SQLite database error: %v\n", err)
  }
  o.Db.Exec(`DELETE FROM struc WHERE idocn=? and strtp=?;`, idocn, strtp)
}

func (o *Outsqlt_tp) IsrtStruc(idocn, strtp string, pnode, cnode Keyst_tp) {
  if strtp == "GRP" {
    pd, _ := strconv.Atoi(pnode.Dname)
    pnode.Dname = fmt.Sprintf("%02d", pd)
    cd, _ := strconv.Atoi(cnode.Dname)
    cnode.Dname = fmt.Sprintf("%02d", cd)
  }
  seqno := fmt.Sprintf("%04d", pnode.Seqno)
  err := o.Db.Exec(`INSERT INTO struc VALUES(?,?,?,?,?,?,?)`, idocn, strtp,
    pnode.Rname, pnode.Dname, seqno, cnode.Rname, cnode.Dname)
  if err != nil {
    log.Fatalf("Insert struc sql table error: %v\n", err)
  }
}
