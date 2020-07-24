package config

var ConfigYAMLSampleDefault = []byte(`
server:
  address: :8080
  tls:
    enable: False
    cert: ""
    key: ""
  limits:
    groups: 100
    conns: 1000
  seed: 0
  log:
    enable_json: False
    level: info
  enable_broadcast: False
  enable_web: False
`)

var ConfigYAMLSampleAddressAndTLS = []byte(`
server:
  address: :9999
  tls:
    enable: True
    cert: "path/to/cert"
    key: "path/to/key"
`)

var ConfigYAMLSampleBullshitSyntax = []byte(`
server:
   address:
 :9999
  tls:
    enable: True
     cert: "path/to/cert"
    key: "path/to/key"
`)
