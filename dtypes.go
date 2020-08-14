package idoc2txt

type Param_tp struct {
  Optn string
  Prm1 string
  Prm2 string
}

type Constant_tp struct {
  Cntrl string `json:"contrl"`
  Clien string `json:"client"`
}

type Program_tp struct {
  Dbonm string `json:"dboNam"`
  Dbodr string `json:"dboDir"`
  Inpdr string `json:"inpDir"`
  Outdr string `json:"outDir"`
  Ifilt string `json:"inFilt"`
  Ifnam string `json:"inName"`
  Ofnam string `json:"ouName"`
}

type Run_tp struct {
  Optcd string `json:"option"`
  Objnm string `json:"objNam"`
  Dbonm string `json:"dboNam"`
  Dbodr string `json:"dboDir"`
  Inpdr string `json:"inpDir"`
  Outdr string `json:"outDir"`
  Ifilt string `json:"inFilt"`
  Ifnam string `json:"inName"`
  Ofnam string `json:"ouName"`
  Rcvpf string `json:"rcPrnF"`
}

type Parsl_tp struct {
  Label Reclb_tp
  Value string
}

type Reclb_tp struct {
  Ident string
  Recnm string
  Rectp string
}

type Idcdf_tp struct {
  Name string
  Type string
  Cols [2]string // Name, Extn
}

type Grpdf_tp struct {
  Name string
  Type string
  Cols [5]string // Numb, Levl, Stat, Mnlp, Mxlp
}

type Segdf_tp struct {
  Name string
  Type string
  Cols [6]string // Name, Type, Levl, Stat, Mnlp, Mxlp
}

type Flddf_tp struct {
  Name string
  Type string
  Clas string
  Cols [7]string // Name, Text, Type, Lgth, Seqn, Strp, Endp
}

type Args_tp struct {
  idocn string
  rname string
  rtype string
  rclas string
  dname string
  dtype string
  dtext string
  extsn string
  stats string
  gnumb int
  level int
  minlp int
  maxlp int
  lngth int
  seqno int
  strps int
  endps int
}

type Keyst_tp struct {
  Rname string
  Dname string
  Level int
  Seqno int
}

type Hstruc_tp struct {
  Sgnum string
  Sgnam string
  Sglvl string
}
