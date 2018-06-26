// Config is put into a different package to prevent cyclic imports in case
// it is needed in several locations

package config

import "time"

type Config struct {
	Period        time.Duration `config:"period"`
	UrlWS         string        `config:"url_ws"`
	Username      string        `config:"username"`
	Password      string        `config:"password"`
	Group         string        `config:"group"`
	RootConfigXML string        `config:"rootConfigXML"`
}

var DefaultConfig = Config{
	Period:        2 * time.Second,
	UrlWS:         "http://172.16.253.60:8080/",
	Username:      "98765432100",
	Password:      "33C3109AAA028CCB",
	Group:         "SERPRO",
	RootConfigXML: "/home/youre/workspaceGo/src/github.com/yourepena/qwcfp-client-go/xml-conf/",
}
