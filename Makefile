APP_NAME := ogm-group
BUILD_VERSION   := $(shell git tag --contains)
BUILD_TIME      := $(shell date "+%F %T")
COMMIT_SHA1     := $(shell git rev-parse HEAD )

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
	gomu --registry=etcd --client=grpc call xtc.ogm.group Collection.Get '{"uuid":"098f6bcd4621d373cade4e832627b4f6"}'
	# 搜索集合
	gomu --registry=etcd --client=grpc call xtc.ogm.group Collection.Search '{"name":"1122"}'
	gomu --registry=etcd --client=grpc call xtc.ogm.group Collection.Search '{"name":"es"}'
	gomu --registry=etcd --client=grpc call xtc.ogm.group Collection.Search '{"name":"t1"}'
	# 加入成员,无参数
	gomu --registry=etcd --client=grpc call xtc.ogm.group Member.Add 
	# 加入成员,集合不存在
	gomu --registry=etcd --client=grpc call xtc.ogm.group Member.Add '{"collection":"00000000"}'
	# 加入成员
	gomu --registry=etcd --client=grpc call xtc.ogm.group Member.Add '{"collection":"098f6bcd4621d373cade4e832627b4f6", "element":"0001"}'
	gomu --registry=etcd --client=grpc call xtc.ogm.group Member.Add '{"collection":"098f6bcd4621d373cade4e832627b4f6", "element":"0002"}'
	gomu --registry=etcd --client=grpc call xtc.ogm.group Member.Add '{"collection":"098f6bcd4621d373cade4e832627b4f6", "element":"0003"}'
	gomu --registry=etcd --client=grpc call xtc.ogm.group Member.Add '{"collection":"5a105e8b9d40e1329780d62ea2265d8a", "element":"0001"}'
	gomu --registry=etcd --client=grpc call xtc.ogm.group Member.Add '{"collection":"5a105e8b9d40e1329780d62ea2265d8a", "element":"0002"}'
	# 定位成员
	gomu --registry=etcd --client=grpc call xtc.ogm.group Member.Where '{"element":"0002"}'
	# 列举成员,集合不存在
	gomu --registry=etcd --client=grpc call xtc.ogm.group Member.List '{"collection":"00000000"}'
	# 列举成员
	gomu --registry=etcd --client=grpc call xtc.ogm.group Member.List 
	gomu --registry=etcd --client=grpc call xtc.ogm.group Member.List '{"collection":"098f6bcd4621d373cade4e832627b4f6"}'
	# 列举成员
	gomu --registry=etcd --client=grpc call xtc.ogm.group Member.List '{"collection":"098f6bcd4621d373cade4e832627b4f6", "offset":1, "count":1}'
	# 获取成员, 无参数
	gomu --registry=etcd --client=grpc call xtc.ogm.group Member.Get 
	# 获取成员, 不存在 
	gomu --registry=etcd --client=grpc call xtc.ogm.group Member.Get '{"uuid":"0000000"}'
	# 获取成员
	gomu --registry=etcd --client=grpc call xtc.ogm.group Member.Get '{"uuid":"0000000"}'
	# 删除成员，无参数
	gomu --registry=etcd --client=grpc call xtc.ogm.group Member.Remove 
	# 删除成员，不存在
	gomu --registry=etcd --client=grpc call xtc.ogm.group Member.Remove '{"uuid":"0000000"}'
	# 删除成员
	gomu --registry=etcd --client=grpc call xtc.ogm.group Member.Remove '{"uuid":"7a07a596be5f45821274550975675323"}'
	gomu --registry=etcd --client=grpc call xtc.ogm.group Member.Remove '{"uuid":"fe141627a1cadad438b1203fb086b0a4"}'

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
