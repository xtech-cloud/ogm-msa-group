APP_NAME := ogm-group
BUILD_VERSION   := $(shell git tag --contains)
BUILD_TIME      := $(shell date "+%F %T")
COMMIT_SHA1     := $(shell git rev-parse HEAD )
collection := $(shell cat /tmp/ogm-group-collection)

.PHONY: build
build:
	go build -ldflags \
		"\
		-X 'main.BuildVersion=${BUILD_VERSION}' \
		-X 'main.BuildTime=${BUILD_TIME}' \
		-X 'main.CommitID=${COMMIT_SHA1}' \
		"\
		-o ./bin/${APP_NAME}

.PHONY: run
run:
	./bin/${APP_NAME}

.PHONY: install
install:
	go install

.PHONY: clean
clean:
	rm -rf /tmp/ogm-group.db

.PHONY: call
call:
	gomu --registry=etcd --client=grpc call xtc.ogm.group Healthy.Echo '{"msg":"hello"}'
	# -------------------------------------------------------------------------
	# 创建集合, 缺少参数
	gomu --registry=etcd --client=grpc call xtc.ogm.group Collection.Make
	gomu --registry=etcd --client=grpc call xtc.ogm.group Collection.Make '{"name":"test"}'
	gomu --registry=etcd --client=grpc call xtc.ogm.group Collection.Make '{"name":"test1"}'
	# 创建集合，已存在
	gomu --registry=etcd --client=grpc call xtc.ogm.group Collection.Make '{"name":"test"}'
	# 列举集合，无参数
	gomu --registry=etcd --client=grpc call xtc.ogm.group Collection.List
	# 列举集合
	gomu --registry=etcd --client=grpc call xtc.ogm.group Collection.List '{"offset":1, "count":1}'
	# 获取集合，无参数
	gomu --registry=etcd --client=grpc call xtc.ogm.group Collection.Get
	# 获取集合，不存在
	gomu --registry=etcd --client=grpc call xtc.ogm.group Collection.Get '{"uuid":"00000000"}'
	# 获取集合
	gomu --registry=etcd --client=grpc call xtc.ogm.group Collection.Get '{"uuid":"${collection}"}'
	# 搜索集合
	gomu --registry=etcd --client=grpc call xtc.ogm.group Collection.Search '{"name":"1122"}'
	gomu --registry=etcd --client=grpc call xtc.ogm.group Collection.Search '{"name":"es"}'
	gomu --registry=etcd --client=grpc call xtc.ogm.group Collection.Search '{"name":"t1"}'
	# 加入成员,无参数
	gomu --registry=etcd --client=grpc call xtc.ogm.group Element.Add 
	# 加入成员,集合不存在
	gomu --registry=etcd --client=grpc call xtc.ogm.group Element.Add '{"collection":"00000000"}'
	# 加入成员
	gomu --registry=etcd --client=grpc call xtc.ogm.group Element.Add '{"collection":"${collection}", "key":"0001", "label":["a1","a2"]}'
	gomu --registry=etcd --client=grpc call xtc.ogm.group Element.Add '{"collection":"${collection}", "key":"0002"}'
	gomu --registry=etcd --client=grpc call xtc.ogm.group Element.Add '{"collection":"${collection}", "key":"0003"}'
	gomu --registry=etcd --client=grpc call xtc.ogm.group Element.Add '{"collection":"${collection}", "key":"0001"}'
	gomu --registry=etcd --client=grpc call xtc.ogm.group Element.Add '{"collection":"${collection}", "key":"0002"}'
	# 定位成员
	gomu --registry=etcd --client=grpc call xtc.ogm.group Element.Where '{"key":"0002"}'
	# 列举成员,集合不存在
	gomu --registry=etcd --client=grpc call xtc.ogm.group Element.List '{"collection":"00000000"}'
	# 列举成员
	gomu --registry=etcd --client=grpc call xtc.ogm.group Element.List 
	gomu --registry=etcd --client=grpc call xtc.ogm.group Element.List '{"collection":"${collection}"}'
	# 列举成员
	gomu --registry=etcd --client=grpc call xtc.ogm.group Element.List '{"collection":"${collection}", "offset":1, "count":1}'
	# 获取成员, 无参数
	gomu --registry=etcd --client=grpc call xtc.ogm.group Element.Get 
	# 获取成员, 不存在 
	gomu --registry=etcd --client=grpc call xtc.ogm.group Element.Get '{"uuid":"0000000"}'
	# 获取成员
	gomu --registry=etcd --client=grpc call xtc.ogm.group Element.Get '{"uuid":"0000000"}'
	# 删除成员，无参数
	gomu --registry=etcd --client=grpc call xtc.ogm.group Element.Remove 
	# 删除成员，不存在
	gomu --registry=etcd --client=grpc call xtc.ogm.group Element.Remove '{"uuid":"0000000"}'
	# 删除成员
	gomu --registry=etcd --client=grpc call xtc.ogm.group Element.Remove '{"uuid":"7a07a596be5f45821274550975675323"}'
	gomu --registry=etcd --client=grpc call xtc.ogm.group Element.Remove '{"uuid":"fe141627a1cadad438b1203fb086b0a4"}'

.PHONY: post
post:
	curl -X POST -d '{"msg":"hello"}' localhost/ogm/group/Healthy/Echo

.PHONY: dist
dist:
	mkdir dist
	tar -zcf dist/${APP_NAME}-${BUILD_VERSION}.tar.gz ./bin/${APP_NAME}

.PHONY: docker
docker:
	docker build -t xtechcloud/${APP_NAME}:${BUILD_VERSION} .
	docker rm -f ${APP_NAME}
	docker run --restart=always --name=${APP_NAME} --net=host -v /data/${APP_NAME}:/ogm -e MSA_REGISTRY_ADDRESS='localhost:2379' -e MSA_CONFIG_DEFINE='{"source":"file","prefix":"/ogm/config","key":"${APP_NAME}.yaml"}' -d xtechcloud/${APP_NAME}:${BUILD_VERSION}
	docker logs -f ${APP_NAME}
