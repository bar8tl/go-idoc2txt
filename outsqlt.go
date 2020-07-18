package idoc2txt

import "code.google.com/p/go-sqlite/go1/sqlite3"
import "fmt"
import "log"
import "strconv"

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
  args := sqlite3.NamedArgs{
    ":01": a.idocn, ":02": a.rname, ":03": a.dname, ":04": a.rclas, ":05": a.rtype, ":06": a.dtype, ":07": a.dtext, ":08": a.extsn,
    ":09": a.gnumb, ":10": a.level, ":11": a.stats, ":12": a.minlp, ":13": a.maxlp, ":14": a.lngth, ":15": a.seqno, ":16": a.strps,
    ":17": a.endps,
  }
  err := o.Db.Exec(`INSERT INTO items VALUES(:01,:02,:03,:04,:05,:06,:07,:08,:09,:10,:11,:12,:13,:14,:15,:16,:17)`, args)
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
  args := sqlite3.NamedArgs{
    ":01": idocn, ":02": strtp, ":03": pnode.Rname, ":04": pnode.Dname, ":05": seqno, ":06": cnode.Rname, ":07": cnode.Dname,
  }
  err := o.Db.Exec(`INSERT INTO struc VALUES(:01,:02,:03,:04,:05,:06,:07)`, args)
  if err != nil {
    log.Fatalf("Insert struc sql table error: %v\n", err)
  }
}
