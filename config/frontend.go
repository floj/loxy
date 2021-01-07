package config

type Frontend struct {
	Name        string               `hcl:"name,label"`
	Bind        string               `hcl:"bind,optional"`
	Port        int                  `hcl:"port"`
	TLS         *FrontendTLS         `hcl:"tls,block"`
	Routes      []FrontendRoute      `hcl:"route,block"`
	Middlewares *FrontendMiddlewares `hcl:"middleware,block"`
}

type FrontendTLS struct {
	Cert string `hcl:"cert"`
	Key  string `hcl:"key"`
}

type FrontendMiddlewares struct {
	ProxyHeaders *MiddlewareProxyHeaders `hcl:"proxy_headers,block"`
	Logger       *MiddlewareLogger       `hcl:"logger,block"`
}

type FrontendRoute struct {
	Name         string                     `hcl:"name,label"`
	Backend      string                     `hcl:"backend"`
	Condition    *FrontendRouteCondition    `hcl:"match,block"`
	Modification *FrontendRouteModification `hcl:"modify,block"`
}

type FrontendRouteCondition struct {
	Headers []FieldMatcher  `hcl:"header,block"`
	Paths   []StringMatcher `hcl:"path,block"`
}

type StringMatcher struct {
	Eq     string `hcl:"is,optional"`
	Prefix string `hcl:"has_prefix,optional"`
	Suffix string `hcl:"has_suffix,optional"`
	Regexp string `hcl:"matches,optional"`

	EqA     []string `hcl:"is_any,optional"`
	PrefixA []string `hcl:"has_any_prefix,optional"`
	SuffixA []string `hcl:"has_any_suffix,optional"`
	RegexpA []string `hcl:"matches_any,optional"`
}

type FieldMatcher struct {
	Field string `hcl:"field,label"`

	Eq     string `hcl:"is,optional"`
	Prefix string `hcl:"has_prefix,optional"`
	Suffix string `hcl:"has_suffix,optional"`
	Regexp string `hcl:"matches,optional"`

	EqA     []string `hcl:"is_any,optional"`
	PrefixA []string `hcl:"has_any_prefix,optional"`
	SuffixA []string `hcl:"has_any_suffix,optional"`
	RegexpA []string `hcl:"matches_any,optional"`
}

type FrontendRouteModification struct {
	Headers []FieldModifier  `hcl:"header,block"`
	Paths   []StringModifier `hcl:"path,block"`
}

type FieldModifier struct {
	Field string `hcl:"field,label"`

	Remove   bool   `hcl:"remove,optional"`
	SetValue string `hcl:"set,optional"`
	AddValue string `hcl:"add,optional"`
}

type StringModifier struct {
	StripPrefix []string `hcl:"strip_prefix,optional"`
	StripSuffix []string `hcl:"strip_suffix,optional"`
}
