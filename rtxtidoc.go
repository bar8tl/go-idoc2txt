// RTXTIDOC: 2019-11-19 BAR8TL
// E-invoicing Tools: Execution manager for SAP-IDoc definition and conversion
// Version 1.0.0
package main

import rb "bar8tl/p/idoc2txt"

// Selector of functions to execute
func main() {
  s := rb.NewSettings("config.json")
  for _, parm := range s.Cmdpr {
    if parm.Optn == "cdb" {        // Create reference IDoc-definition database
      dbo := rb.NewDdbo()
      dbo.CrtTables(parm, s)
    } else if parm.Optn == "rid" { // Read and upload IDoc-definition files
      rid := rb.NewDrid()
      rid.ProcInput(parm, s)
    } else if parm.Optn == "doc" { // Output of IDoc-definition documentation
      doc := rb.NewDdoc()
      doc.ProcDocument(parm, s)
    } else if parm.Optn == "unf" { // Convert IDOC-data format from SAP to TXT
      unf := rb.NewDunf()
      unf.UnfoldIdocs(parm, s)
    }
  }
}
