package idoc2txt

import "code.google.com/p/go-sqlite/go1/sqlite3"
import "log"
import "strings"

// Data type Ddbo for creation of reference IDoc definition database
type Ddbo_tp struct {
  Dbodr, Cnnst string
}

// Constructor of object Ddbo: Define database location folder and SQlite3 database full connection string
func NewDdbo(parm Param_tp, s Settings_tp) *Ddbo_tp {
  var d Ddbo_tp
  for _, run := range s.Runlv {
    if parm.Optn == run.Optcd {
      if len(run.Dbodr) == 0 {
        d.Dbodr = s.Progm.Dbodr
      } else {
        d.Dbodr = run.Dbodr
      }
    }
  }
  if len(parm.Prm1) == 0 {
    d.Cnnst = s.Cnnst
  } else {
    d.Cnnst = strings.Replace(s.Cnnsq, "@", d.Dbodr+parm.Prm1, 1)
  }
  return &d
}

// Public option CDB: Creation of tables in database
func (d *Ddbo_tp) CrtTables() {
  d.CrtItems().CrtStruc()
}

// Function to create table ITEMS: Which table contains specifications for IDOC control/data/status records
func (d *Ddbo_tp) CrtItems() *Ddbo_tp {
  db, _ := sqlite3.Open(d.Cnnst)
  defer db.Close()
  db.Exec(`
    DROP TABLE IF EXISTS items;
  `)
  err := db.Exec(`
    CREATE TABLE items (
      idocn TEXT,
      rname TEXT,
      dname TEXT,
      rclas TEXT,
      rtype TEXT,
      dtype TEXT,
      dtext TEXT,
      extsn TEXT,
      gnumb INTEGER,
      level INTEGER,
      stats TEXT,
      minlp INTEGER,
      maxlp INTEGER,
      lngth INTEGER,
      seqno INTEGER,
      strps INTEGER,
      endps INTEGER,
      PRIMARY KEY (idocn, rname, dname)
    );`,
  )
  if err != nil {
    log.Fatalf("Table items creation error: %s\n", err)
  }
  return d
}

// Function to create table STRUC: Which table contains specifications for structure of IDOC elements
func (d *Ddbo_tp) CrtStruc() *Ddbo_tp {
  db, _ := sqlite3.Open(d.Cnnst)
  defer db.Close()
  db.Exec(`
    DROP TABLE IF EXISTS struc;
  `)
  err := db.Exec(`
    CREATE TABLE struc (
      parnt TEXT,
      child TEXT,
      PRIMARY KEY (parnt, child)
    );
  `)
  if err != nil {
    log.Fatalf("Table struc creation error: %s\n", err)
  }
  return d
}
