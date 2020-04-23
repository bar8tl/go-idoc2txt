package idoc2txt

const CNNS_SQLIT3    = "file:@?file:locked.sqlite?cache=shared&mode=rwc"
const CONTROL_CODE   = "EDI_DC40"
const CLIENT_CODE    = "011"

const DB_NAME        = "idoctp.db"
const DB_DIR         = "c:\\c_portab\\01_rb\\ProgramData\\go-idoc2txt\\"
const INPUTS_DIR     = "c:\\c_portab\\01_rb\\_rbprojects\\go-idoc2txt\\idoctypes\\"
const OUTPUTS_DIR    = "c:\\c_portab\\01_rb\\ProgramData\\go-idoc2txt\\"
const INPUTS_FILTER  = "!(*processed*)"
const INPUTS_NAMING  = "rundt+'_'+idocn+docno+docdt+'_inp_processed'"
const OUTPUTS_NAMING = "rundt+'_'+idocn+docno+docdt+'_out'"
