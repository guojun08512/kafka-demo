package client

import "net/http"

type (
	DocDefs map[string]*DocDef
)

// ServiceManifest model
type ServiceManifest struct {
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Version     string                 `json:"version"`
	Repository  string                 `json:"repository"`
	Categories  []string               `json:"categories"`
	SupportMIME []string               `json:"support_mime"`
	Scope       []string               `json:"scope"`
	Params      []Param                `json:"params"`
	DocTypes    DocDefs                `json:"doctypes,omitempty"`
	Services    map[string]interface{} `json:"services,omitempty"`
}

// Param model
type Param struct {
	Name         string                                                    `json:"name"`
	Type         string                                                    `json:"type"`
	Default      interface{}                                               `json:"default,omitempty"`
	TypeChecker  func(value map[string]interface{}) bool                   `json:"-"`
	ParamUpdater func(value map[string]interface{}) map[string]interface{} `json:"-"`
	Array        bool                                                      `json:"is_array,omitempty"`
	Description  string                                                    `json:"description"`
}

type ServTrigger struct {
	Debounce       string `json:"debounce"`
	Type           string `json:"type"`
	TriggerOptions string `json:"trigger"`
}

// DocDef service doctype define
type DocDef struct {
	Index  map[string][]string `json:"index,omitempty"`
	Unique string              `json:"unique,omitempty"`
}

type SealClient struct {
	Domain            string
	Scheme            string
	Authorizer        Authorizer
	HTTPClient        *http.Client
	RefreshAuthorizer func() (Authorizer, http.CookieJar, error)
}
