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
  flags:
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

var ConfigYAMLSampleAddressAndTLSAndLimits = []byte(`
server:
  address: :9999
  tls:
    enable: True
    cert: "path/to/cert"
    key: "path/to/key"
  limits:
    groups: 144
    conns: 4123
  flags:
    enable_broadcast: True
`)

var ConfigYAMLSampleAddressAndTLSAndLimitsAndCORS = []byte(`
server:
  address: :9999
  tls:
    enable: True
    cert: "path/to/cert"
    key: "path/to/key"
  limits:
    groups: 144
    conns: 4123
  flags:
    enable_broadcast: True
    forbid_cors: True
`)

var ConfigYAMLSampleLimitsAndSentry = []byte(`
server:
  limits:
    groups: 144
    conns: 4123
  sentry:
    enable: True
    dsn: https://public@sentry.example.com/1
`)

var ConfigYAMLSampleSentryAndDebug = []byte(`
server:
  sentry:
    enable: True
    dsn: https://public@sentry.example.com/1
  flags:
    debug: True
`)
