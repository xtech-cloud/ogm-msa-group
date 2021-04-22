package handler

import (
	"context"
	"errors"
	"ogm-msa-group/config"
	"ogm-msa-group/model"

	"github.com/micro/go-micro/v2/logger"
	proto "github.com/xtech-cloud/ogm-msp-group/proto/group"
)

type Collection struct{}

func (this *Collection) Make(_ctx context.Context, _req *proto.CollectionMakeRequest, _rsp *proto.BlankResponse) error {
	logger.Infof("Received Collection.Make, req is %v", _req)
	_rsp.Status = &proto.Status{}

	if "" == _req.Name {
		_rsp.Status.Code = 1
		_rsp.Status.Message = "name is required"
		return nil
	}

	// 本地数据库使用存储桶名生成UUID，方便测试和开发
	uuid := model.NewUUID()
	if config.Schema.Database.Lite {
		uuid = model.ToUUID(_req.Name)
	}

	Collection := &model.Collection{
		UUID:     uuid,
		Name:     _req.Name,
		Capacity: _req.Capacity,
	}

	dao := model.NewCollectionDAO(nil)
	err := dao.Insert(Collection)
	if errors.Is(err, model.ErrCollectionExists) {
		_rsp.Status.Code = 2
		_rsp.Status.Message = err.Error()
		return nil
	}

	return err
}

func (this *Collection) List(_ctx context.Context, _req *proto.CollectionListRequest, _rsp *proto.CollectionListResponse) error {
	logger.Infof("Received Collection.List, req is %v", _req)
	_rsp.Status = &proto.Status{}

	offset := int64(0)
	count := int64(100)

	if _req.Offset > 0 {
		offset = _req.Offset
	}

	if _req.Count > 0 {
		count = _req.Count
	}

	dao := model.NewCollectionDAO(nil)

	total, err := dao.Count()
	if nil != err {
		return nil
	}
	Collections, err := dao.List(offset, count)
	if nil != err {
		return nil
	}

	_rsp.Total = uint64(total)
	_rsp.Entity = make([]*proto.CollectionEntity, len(Collections))
	for i, collection := range Collections {
		_rsp.Entity[i] = &proto.CollectionEntity{
			Uuid:     collection.UUID,
			Name:     collection.Name,
			Capacity: collection.Capacity,
		}
	}

	return nil
}

func (this *Collection) Remove(_ctx context.Context, _req *proto.CollectionRemoveRequest, _rsp *proto.BlankResponse) error {
	logger.Infof("Received Collection.Remove, req is %v", _req)
	_rsp.Status = &proto.Status{}

	if "" == _req.Uuid {
		_rsp.Status.Code = 1
		_rsp.Status.Message = "uuid is required"
		return nil
	}

	dao := model.NewCollectionDAO(nil)
	err := dao.Delete(_req.Uuid)
	if errors.Is(err, model.ErrCollectionNotFound) {
		_rsp.Status.Code = 2
		_rsp.Status.Message = err.Error()
		return nil
	}
	return err
}

func (this *Collection) Get(_ctx context.Context, _req *proto.CollectionGetRequest, _rsp *proto.CollectionGetResponse) error {
	logger.Infof("Received Collection.Get, req is %v", _req)
	_rsp.Status = &proto.Status{}

	if "" == _req.Uuid {
		_rsp.Status.Code = 1
		_rsp.Status.Message = "uuid is required"
		return nil
	}

	dao := model.NewCollectionDAO(nil)
	query := model.CollectionQuery{
		UUID: _req.Uuid,
	}
	collection, err := dao.QueryOne(&query)
	if errors.Is(err, model.ErrCollectionNotFound) {
		_rsp.Status.Code = 2
		_rsp.Status.Message = err.Error()
		return nil
	}
    if collection.UUID == "" {
		_rsp.Status.Code = 2
		_rsp.Status.Message = "not found"
		return nil
    }
	_rsp.Entity = &proto.CollectionEntity{
		Uuid:     collection.UUID,
		Name:     collection.Name,
		Capacity: collection.Capacity,
	}

	return nil
}
