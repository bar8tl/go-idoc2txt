package idoc2txt

import "code.google.com/p/go-sqlite/go1/sqlite3"
import "log"

// Data type Ddbo for creation of reference IDoc definition database
type Ddbo_tp struct {
  Cnnsq, Cnnst string
  Dbonm, Dbodr string
  Inpdr, Outdr string
}

// Constructor of object Ddbo: Define database name, location folder and SQlite3 database full connection string
func NewDdbo(parm Param_tp, s Settings_tp) *Ddbo_tp {
  var d Ddbo_tp
  s.SetRunVars(parm, s)
  d.Cnnsq, d.Cnnst = s.Cnnsq, s.Cnnst
  d.Dbonm, d.Dbodr = s.Dbonm, s.Dbodr
  d.Inpdr          = s.Inpdr
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
