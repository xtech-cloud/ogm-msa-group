package handler

import (
	"context"
	"errors"
	"fmt"
	"omo-msa-group/model"
	"omo-msa-group/publisher"

	"github.com/micro/go-micro/v2/logger"
	proto "github.com/xtech-cloud/omo-msp-group/proto/group"
)

type Member struct{}

func (this *Member) Add(_ctx context.Context, _req *proto.MemberAddRequest, _rsp *proto.BlankResponse) error {
	logger.Infof("Received Member.Add, req is %v", _req)
	_rsp.Status = &proto.Status{}

	if "" == _req.Collection {
		_rsp.Status.Code = 1
		_rsp.Status.Message = "collection is required"
		return nil
	}

	if "" == _req.Element {
		_rsp.Status.Code = 1
		_rsp.Status.Message = "element is required"
		return nil
	}

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

	daoMember := model.NewMemberDAO(nil)
	count, err := daoMember.CountOfCollection(_req.Collection)
	if nil != err {
		return err
	}

	// 0为无限制
	if collection.Capacity > 0 {
		if uint64(count) >= collection.Capacity {
			_rsp.Status.Code = 3
			_rsp.Status.Message = "out of capacity"
			return nil
		}
	}

	uuid := model.ToUUID(_req.Collection + _req.Element)
	member := &model.Member{
		UUID:       uuid,
		Collection: _req.Collection,
		Element:    _req.Element,
	}

	err = daoMember.Insert(member)
	if errors.Is(err, model.ErrMemberExists) {
		_rsp.Status.Code = 4
		_rsp.Status.Message = err.Error()
		return nil
	}
	if nil != err {
		return err
	}

	ctx := buildNotifyContext(_ctx, "root")
	publisher.Publish(ctx, "member/add", "", fmt.Sprintf("%v", _req))
	return nil
}

func (this *Member) Get(_ctx context.Context, _req *proto.MemberGetRequest, _rsp *proto.MemberGetResponse) error {
	logger.Infof("Received Member.Get, req is %v", _req)
	_rsp.Status = &proto.Status{}

	if "" == _req.Uuid {
		_rsp.Status.Code = 1
		_rsp.Status.Message = "uuid is required"
		return nil
	}

	dao := model.NewMemberDAO(nil)
	query := model.MemberQuery{
		UUID: _req.Uuid,
	}
	member, err := dao.QueryOne(&query)
	if errors.Is(err, model.ErrMemberNotFound) {
		_rsp.Status.Code = 2
		_rsp.Status.Message = err.Error()
		return nil
	}
	_rsp.Entity = &proto.MemberEntity{
		Uuid:       member.UUID,
		Collection: member.Collection,
		Element:    member.Element,
	}

	ctx := buildNotifyContext(_ctx, "root")
	publisher.Publish(ctx, "member/get", "", fmt.Sprintf("%v", _req))
	return nil
}

func (this *Member) Remove(_ctx context.Context, _req *proto.MemberRemoveRequest, _rsp *proto.BlankResponse) error {
	logger.Infof("Received Member.Remove, req is %v", _req)
	_rsp.Status = &proto.Status{}

	if "" == _req.Uuid {
		_rsp.Status.Code = 1
		_rsp.Status.Message = "uuid is required"
		return nil
	}

	dao := model.NewMemberDAO(nil)
	err := dao.Delete(_req.Uuid)
	if errors.Is(err, model.ErrMemberNotFound) {
		_rsp.Status.Code = 2
		_rsp.Status.Message = err.Error()
		return nil
	}
	ctx := buildNotifyContext(_ctx, "root")
	publisher.Publish(ctx, "member/remove", "", fmt.Sprintf("%v", _req))

	return nil
}

func (this *Member) List(_ctx context.Context, _req *proto.MemberListRequest, _rsp *proto.MemberListResponse) error {
	logger.Infof("Received Member.List, req is %v", _req)
	_rsp.Status = &proto.Status{}

	offset := int64(0)
	count := int64(100)

	if _req.Offset > 0 {
		offset = _req.Offset
	}

	if _req.Count > 0 {
		count = _req.Count
	}

	dao := model.NewMemberDAO(nil)

	total, err := dao.CountOfCollection(_req.Collection)
	if nil != err {
		return nil
	}
	members, err := dao.List(offset, count, _req.Collection)
	if nil != err {
		return nil
	}

	_rsp.Total = uint64(total)
	_rsp.Entity = make([]*proto.MemberEntity, len(members))
	for i, member := range members {
		_rsp.Entity[i] = &proto.MemberEntity{
			Uuid:       member.UUID,
			Collection: member.Collection,
			Element:    member.Element,
		}
	}

	ctx := buildNotifyContext(_ctx, "root")
	publisher.Publish(ctx, "member/list", "", fmt.Sprintf("%v", _req))
	return nil
}

func (this *Member) Where(_ctx context.Context, _req *proto.MemberWhereRequest, _rsp *proto.MemberWhereResponse) error {
	logger.Infof("Received Member.Where, req is %v", _req)
	_rsp.Status = &proto.Status{}

	if "" == _req.Element{
		_rsp.Status.Code = 1
		_rsp.Status.Message = "element is required"
		return nil
	}

	dao := model.NewMemberDAO(nil)
    query := model.MemberQuery {
        Element: _req.Element,
    }
	members, err := dao.QueryMany(&query)
	if nil != err {
		return nil
	}

	_rsp.Collection = make([]string, len(members))
	for i, member := range members {
		_rsp.Collection[i] = member.Collection;
	}

	ctx := buildNotifyContext(_ctx, "root")
	publisher.Publish(ctx, "member/where", "", fmt.Sprintf("%v", _req))
	return nil
}
