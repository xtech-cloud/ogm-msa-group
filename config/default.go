package config

const defaultYAML string = `
service:
  name: xtc.api.ogm.group
  address: :9609
  ttl: 15
  interval: 10
logger:
  level: trace
  dir: /var/log/ogm/
database:
  lite: true
  mysql:
    address: localhost:3306
    user: root
    password: mysql@OMO
    db: ogm
  sqlite:
    path: /tmp/ogm-group.db
`
