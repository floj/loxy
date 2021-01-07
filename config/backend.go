package config

type Backend struct {
	Name         string               `hcl:"name,label"`
	ReverseProxy *BackendReverseProxy `hcl:"reverse_proxy,block"`
	FileServer   *BackendFileServer   `hcl:"file_server,block"`
	Prometheus   *BackendPrometheus   `hcl:"prometheus,block"`
}

type BackendReverseProxy struct {
	Targets []string `hcl:"targets"`
}

type BackendFileServer struct {
	Root string `hcl:"root"`
}

type BackendPrometheus struct{}
