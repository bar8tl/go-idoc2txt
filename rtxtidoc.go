/*** rtxtidoc.go : 2019-11-19 BAR8TL - E-invoicing Tools: Execution manager for SAP-IDoc definition and conversions - Version 1.0.0 ***/
package main

import rb "bar8tl/p/idoc2txt"

func main() {
  s := rb.NewSettings("config.json")
  for _, parm := range s.Cmdpr {   // Browse declared parameters in command execution line and process accordingly:
           if parm.Optn == "cdb" { //   CDB Option to create reference IDoc-definition database
      dbo := rb.NewDdbo()
      dbo.CrtTables(parm, s)
    } else if parm.Optn == "rid" { //   RID Option to read and upload IDoc-definition files
      rid := rb.NewDrid()
      rid.ProcInput(parm, s)
    } else if parm.Optn == "unf" { //   UNF Option to read data IDocs and convert the format to flat positional text file
      unf := rb.NewDunf()
      unf.UnfoldIdocs(parm, s)
    }
  }
}
