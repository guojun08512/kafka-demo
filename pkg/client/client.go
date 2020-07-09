package client

// GetDocDefs views preset
func GetDocDefs(types []string) DocDefs {
	defs := make(DocDefs)
	for _, t := range types {
		defs[t] = getDocDef(t)
	}
	return defs
}

func getDocDef(docType string) *DocDef {
	switch docType {
	}
	return &DocDef{}
}
