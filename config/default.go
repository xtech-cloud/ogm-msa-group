package config

const defaultYAML string = `
service:
  name: omo.msa.group
  address: :9609
  ttl: 15
  interval: 10
logger:
  level: info
  dir: /var/log/msa/
database:
  lite: true
  timeout: 10
  mysql:
    address: 127.0.0.1:3306
    user: root
    password: mysql@OMO
    db: msa_group
  sqlite:
    path: /tmp/msa-group.db
publisher:
- /collection/make
- /collection/list
- /collection/remove
- /collection/get
- /member/add
- /member/remove
- /member/list
- /member/get
`
