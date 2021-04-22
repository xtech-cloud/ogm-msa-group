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
	rm -rf /tmp/msa-group.db

.PHONY: call
call:
	MICRO_REGISTRY=consul micro call xtc.api.ogm.group Healthy.Echo '{"msg":"hello"}'
	# -------------------------------------------------------------------------
	# 创建集合, 缺少参数
	MICRO_REGISTRY=consul micro call xtc.api.ogm.group Collection.Make
	MICRO_REGISTRY=consul micro call xtc.api.ogm.group Collection.Make '{"name":"test"}'
	MICRO_REGISTRY=consul micro call xtc.api.ogm.group Collection.Make '{"name":"test1"}'
	# 创建集合，已存在
	MICRO_REGISTRY=consul micro call xtc.api.ogm.group Collection.Make '{"name":"test"}'
	# 列举集合，无参数
	MICRO_REGISTRY=consul micro call xtc.api.ogm.group Collection.List
	# 列举集合
	MICRO_REGISTRY=consul micro call xtc.api.ogm.group Collection.List '{"offset":1, "count":1}'
	# 获取集合，无参数
	MICRO_REGISTRY=consul micro call xtc.api.ogm.group Collection.Get
	# 获取集合，不存在
	MICRO_REGISTRY=consul micro call xtc.api.ogm.group Collection.Get '{"uuid":"00000000"}'
	# 获取集合
	MICRO_REGISTRY=consul micro call xtc.api.ogm.group Collection.Get '{"uuid":"098f6bcd4621d373cade4e832627b4f6"}'
	# 加入成员,无参数
	MICRO_REGISTRY=consul micro call xtc.api.ogm.group Member.Add 
	# 加入成员,集合不存在
	MICRO_REGISTRY=consul micro call xtc.api.ogm.group Member.Add '{"collection":"00000000"}'
	# 加入成员
	MICRO_REGISTRY=consul micro call xtc.api.ogm.group Member.Add '{"collection":"098f6bcd4621d373cade4e832627b4f6", "element":"0001"}'
	MICRO_REGISTRY=consul micro call xtc.api.ogm.group Member.Add '{"collection":"098f6bcd4621d373cade4e832627b4f6", "element":"0002"}'
	MICRO_REGISTRY=consul micro call xtc.api.ogm.group Member.Add '{"collection":"098f6bcd4621d373cade4e832627b4f6", "element":"0003"}'
	MICRO_REGISTRY=consul micro call xtc.api.ogm.group Member.Add '{"collection":"5a105e8b9d40e1329780d62ea2265d8a", "element":"0001"}'
	MICRO_REGISTRY=consul micro call xtc.api.ogm.group Member.Add '{"collection":"5a105e8b9d40e1329780d62ea2265d8a", "element":"0002"}'
	# 定位成员
	MICRO_REGISTRY=consul micro call xtc.api.ogm.group Member.Where '{"element":"0002"}'
	# 列举成员,集合不存在
	MICRO_REGISTRY=consul micro call xtc.api.ogm.group Member.List '{"collection":"00000000"}'
	# 列举成员
	MICRO_REGISTRY=consul micro call xtc.api.ogm.group Member.List 
	MICRO_REGISTRY=consul micro call xtc.api.ogm.group Member.List '{"collection":"098f6bcd4621d373cade4e832627b4f6"}'
	# 列举成员
	MICRO_REGISTRY=consul micro call xtc.api.ogm.group Member.List '{"collection":"098f6bcd4621d373cade4e832627b4f6", "offset":1, "count":1}'
	# 获取成员, 无参数
	MICRO_REGISTRY=consul micro call xtc.api.ogm.group Member.Get 
	# 获取成员, 不存在 
	MICRO_REGISTRY=consul micro call xtc.api.ogm.group Member.Get '{"uuid":"0000000"}'
	# 获取成员
	MICRO_REGISTRY=consul micro call xtc.api.ogm.group Member.Get '{"uuid":"0000000"}'
	# 删除成员，无参数
	MICRO_REGISTRY=consul micro call xtc.api.ogm.group Member.Remove 
	# 删除成员，不存在
	MICRO_REGISTRY=consul micro call xtc.api.ogm.group Member.Remove '{"uuid":"0000000"}'
	# 删除成员
	MICRO_REGISTRY=consul micro call xtc.api.ogm.group Member.Remove '{"uuid":"7a07a596be5f45821274550975675323"}'
	MICRO_REGISTRY=consul micro call xtc.api.ogm.group Member.Remove '{"uuid":"fe141627a1cadad438b1203fb086b0a4"}'
	MICRO_REGISTRY=consul micro call xtc.api.ogm.group Member.Remove '{"uuid":"ddd826052e478dcef4eca0a3b30d3be0"}'
	MICRO_REGISTRY=consul micro call xtc.api.ogm.group Member.Remove '{"uuid":"ccd232d5256626b766d65633c9b68ac8"}'
	MICRO_REGISTRY=consul micro call xtc.api.ogm.group Member.Remove '{"uuid":"77b6798fe86a2e4838ed1f0462fb1ce6"}'
	# 删除集合，缺少参数
	MICRO_REGISTRY=consul micro call xtc.api.ogm.group Collection.Remove
	# 删除集合
	MICRO_REGISTRY=consul micro call xtc.api.ogm.group Collection.Remove '{"uuid":"098f6bcd4621d373cade4e832627b4f6"}'
	MICRO_REGISTRY=consul micro call xtc.api.ogm.group Collection.Remove '{"uuid":"5a105e8b9d40e1329780d62ea2265d8a"}'

.PHONY: post
post:
	curl -X POST -d '{"msg":"hello"}' 127.0.0.1:8080/ogm/group/Healthy/Echo

.PHONY: bm
bm:
	python3 benchmark.py

.PHONY: dist
dist:
	mkdir dist
	tar -zcf dist/${APP_NAME}-${BUILD_VERSION}.tar.gz ./bin/${APP_NAME}
