package handler

import (
	"context"
	"errors"
	"ogm-group/model"
	"strings"

	"github.com/asim/go-micro/v3/logger"
	proto "github.com/xtech-cloud/ogm-msp-group/proto/group"
)

type Element struct{}

func (this *Element) Add(_ctx context.Context, _req *proto.ElementAddRequest, _rsp *proto.UuidResponse) error {
	logger.Infof("Received Element.Add, req is %v", _req)
	_rsp.Status = &proto.Status{}

	if "" == _req.Collection {
		_rsp.Status.Code = 1
		_rsp.Status.Message = "collection is required"
		return nil
	}

	if "" == _req.Key {
		_rsp.Status.Code = 1
		_rsp.Status.Message = "key is required"
		return nil
	}

	//collection是否存在
	daoCollection := model.NewCollectionDAO(nil)
	query := model.CollectionQuery{
		UUID: _req.Collection,
	}
	collection, err := daoCollection.QueryOne(&query)
	if errors.Is(err, model.ErrCollectionNotFound) {
		_rsp.Status.Code = 2
		_rsp.Status.Message = err.Error()
		return nil
	}

	daoElement := model.NewElementDAO(nil)
	count, err := daoElement.CountOfCollection(_req.Collection)
	if nil != err {
		return err
	}

	// 容量是否足够，0为无限制
	if collection.Capacity > 0 {
		if uint64(count) >= collection.Capacity {
			_rsp.Status.Code = 3
			_rsp.Status.Message = "out of capacity"
			return nil
		}
	}

	labelStr := ""
	for _, label := range _req.Label {
		labelStr += label + ","
	}
	if "" != labelStr {
		labelStr = labelStr[:len(labelStr)-1]
	}

	uuid := model.ToUUID(_req.Collection + _req.Key)
	element := &model.Element{
		UUID:       uuid,
		Collection: _req.Collection,
		Key:        _req.Key,
		Alias:      _req.Alias,
		Label:      labelStr,
	}

	err = daoElement.Insert(element)
	if errors.Is(err, model.ErrElementExists) {
		_rsp.Status.Code = 4
		_rsp.Status.Message = err.Error()
		return nil
	}
	if nil != err {
		return err
	}

	_rsp.Uuid = uuid
	return nil
}

func (this *Element) Update(_ctx context.Context, _req *proto.ElementUpdateRequest, _rsp *proto.UuidResponse) error {
	logger.Infof("Received Element.Update, req is %v", _req)
	_rsp.Status = &proto.Status{}

	if "" == _req.Uuid {
		_rsp.Status.Code = 1
		_rsp.Status.Message = "uuid is required"
		return nil
	}

	if "" == _req.Key {
		_rsp.Status.Code = 1
		_rsp.Status.Message = "key is required"
		return nil
	}

	if "" == _req.Alias {
		_rsp.Status.Code = 1
		_rsp.Status.Message = "alias is required"
		return nil
	}

	labelStr := ""
	for _, label := range _req.Label {
		labelStr += label + ","
	}
	if "" != labelStr {
		labelStr = labelStr[:len(labelStr)-1]
	}

	dao := model.NewElementDAO(nil)
	element := model.Element{
		UUID:  _req.Uuid,
		Key:   _req.Key,
		Alias: _req.Alias,
		Label: labelStr,
	}
	err := dao.Update(&element)
	if nil != err {
		_rsp.Status.Code = -1
		_rsp.Status.Message = err.Error()
		return nil
	}
	return nil
}

func (this *Element) Get(_ctx context.Context, _req *proto.ElementGetRequest, _rsp *proto.ElementGetResponse) error {
	logger.Infof("Received Element.Get, req is %v", _req)
	_rsp.Status = &proto.Status{}

	if "" == _req.Uuid {
		_rsp.Status.Code = 1
		_rsp.Status.Message = "uuid is required"
		return nil
	}

	dao := model.NewElementDAO(nil)
	query := model.ElementQuery{
		UUID: _req.Uuid,
	}
	element, err := dao.QueryOne(&query)
	if errors.Is(err, model.ErrElementNotFound) {
		_rsp.Status.Code = 2
		_rsp.Status.Message = err.Error()
		return nil
	}

	labelAry := make([]string, 0)
	for _, label := range strings.Split(element.Label, ",") {
		if "" != label {
			labelAry = append(labelAry, label)
		}
	}
	_rsp.Entity = &proto.ElementEntity{
		Uuid:       element.UUID,
		Collection: element.Collection,
		Key:        element.Key,
		Alias:      element.Alias,
		Label:      labelAry,
	}

	return nil
}

func (this *Element) Remove(_ctx context.Context, _req *proto.ElementRemoveRequest, _rsp *proto.UuidResponse) error {
	logger.Infof("Received Element.Remove, req is %v", _req)
	_rsp.Status = &proto.Status{}

	if "" == _req.Uuid {
		_rsp.Status.Code = 1
		_rsp.Status.Message = "uuid is required"
		return nil
	}

	dao := model.NewElementDAO(nil)
	err := dao.Delete(_req.Uuid)
	if errors.Is(err, model.ErrElementNotFound) {
		_rsp.Status.Code = 2
		_rsp.Status.Message = err.Error()
		return nil
	}
	_rsp.Uuid = _req.Uuid
	return nil
}

func (this *Element) List(_ctx context.Context, _req *proto.ElementListRequest, _rsp *proto.ElementListResponse) error {
	logger.Infof("Received Element.List, req is %v", _req)
	_rsp.Status = &proto.Status{}

	offset := int64(0)
	count := int64(100)

	if _req.Offset > 0 {
		offset = _req.Offset
	}

	if _req.Count > 0 {
		count = _req.Count
	}

	dao := model.NewElementDAO(nil)

	total, elements, err := dao.List(offset, count, _req.Collection)
	if nil != err {
		return nil
	}

	_rsp.Total = uint64(total)
	_rsp.Entity = make([]*proto.ElementEntity, len(elements))
	for i, element := range elements {
		labelAry := make([]string, 0)
		for _, label := range strings.Split(element.Label, ",") {
			if "" != label {
				labelAry = append(labelAry, label)
			}
		}
		_rsp.Entity[i] = &proto.ElementEntity{
			Uuid:       element.UUID,
			Collection: element.Collection,
			Key:        element.Key,
			Alias:      element.Alias,
			Label:      labelAry,
		}
	}

	return nil
}

func (this *Element) Search(_ctx context.Context, _req *proto.ElementSearchRequest, _rsp *proto.ElementListResponse) error {
	logger.Infof("Received Element.Search, req is %v", _req)
	_rsp.Status = &proto.Status{}

	offset := int64(0)
	count := int64(100)

	if _req.Offset > 0 {
		offset = _req.Offset
	}

	if _req.Count > 0 {
		count = _req.Count
	}

	dao := model.NewElementDAO(nil)

	total, elements, err := dao.Search(offset, count, _req.Collection, _req.Key, _req.Alias)
	if nil != err {
		return nil
	}

	_rsp.Total = uint64(total)
	_rsp.Entity = make([]*proto.ElementEntity, len(elements))
	for i, element := range elements {
		labelAry := make([]string, 0)
		for _, label := range strings.Split(element.Label, ",") {
			if "" != label {
				labelAry = append(labelAry, label)
			}
		}
		_rsp.Entity[i] = &proto.ElementEntity{
			Uuid:       element.UUID,
			Collection: element.Collection,
			Key:        element.Key,
			Alias:      element.Alias,
			Label:      labelAry,
		}
	}

	return nil
}
func (this *Element) Where(_ctx context.Context, _req *proto.ElementWhereRequest, _rsp *proto.ElementWhereResponse) error {
	logger.Infof("Received Element.Where, req is %v", _req)
	_rsp.Status = &proto.Status{}

	if "" == _req.Key {
		_rsp.Status.Code = 1
		_rsp.Status.Message = "key is required"
		return nil
	}

	dao := model.NewElementDAO(nil)
	query := model.ElementQuery{
		Key: _req.Key,
	}
	elements, err := dao.QueryMany(&query)
	if nil != err {
		return nil
	}

	_rsp.Collection = make([]string, len(elements))
	for i, element := range elements {
		_rsp.Collection[i] = element.Collection
	}

	return nil
}
