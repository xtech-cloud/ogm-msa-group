package model

import (
	"errors"
	"time"
)

type Member struct {
	UUID       string `gorm:"column:uuid;type:char(32);not null;unique;primaryKey"`
	Collection string `gorm:"column:collection;type:char(32);not null;"`
	Element    string `gorm:"column:element;type:char(32);not null;"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

var ErrMemberExists = errors.New("member exists")
var ErrMemberNotFound = errors.New("member not found")

func (Member) TableName() string {
	return "ogm_group_member"
}

type MemberDAO struct {
	conn *Conn
}

type MemberQuery struct {
	UUID       string
	Collection string
	Element    string
}

func NewMemberDAO(_conn *Conn) *MemberDAO {
	conn := DefaultConn
	if nil != _conn {
		conn = _conn
	}
	return &MemberDAO{
		conn: conn,
	}
}

func (this *MemberDAO) Count() (int64, error) {
	var count int64
	err := this.conn.DB.Model(&Member{}).Count(&count).Error
	return count, err
}

func (this *MemberDAO) CountOfCollection(_collection string) (int64, error) {
	var count int64
	db := this.conn.DB.Model(&Member{})
	if "" != _collection {
		db = db.Where("collection = ?", _collection)
	}
	err := db.Count(&count).Error
	return count, err
}

func (this *MemberDAO) Insert(_member *Member) error {
	var count int64
	err := this.conn.DB.Model(&Member{}).Where("uuid = ?", _member.UUID).Count(&count).Error
	if nil != err {
		return err
	}

	if count > 0 {
		return ErrMemberExists
	}

	return this.conn.DB.Create(_member).Error
}

func (this *MemberDAO) Update(_member *Member) error {
	var count int64
	err := this.conn.DB.Model(&Member{}).Where("uuid = ?", _member.UUID).Count(&count).Error
	if nil != err {
		return err
	}

	if 0 == count {
		return ErrMemberNotFound
	}

	return this.conn.DB.Updates(_member).Error
}

func (this *MemberDAO) Delete(_uuid string) error {
	var count int64
	err := this.conn.DB.Model(&Member{}).Where("uuid = ?", _uuid).Count(&count).Error
	if nil != err {
		return err
	}

	if 0 == count {
		return ErrMemberNotFound
	}

	return this.conn.DB.Where("uuid = ?", _uuid).Delete(&Member{}).Error
}

func (this *MemberDAO) List(_offset int64, _count int64, _collection string) ([]*Member, error) {
	var members []*Member
	db := this.conn.DB
	if "" != _collection {
		db = db.Where("collection = ?", _collection)
	}
	res := db.Offset(int(_offset)).Limit(int(_count)).Order("created_at desc").Find(&members)
	return members, res.Error
}

func (this *MemberDAO) QueryOne(_query *MemberQuery) (*Member, error) {
	db := this.conn.DB.Model(&Member{})
	hasWhere := false
	if "" != _query.UUID {
		db = db.Where("uuid = ?", _query.UUID)
		hasWhere = true
	}
	if "" != _query.Collection {
		db = db.Where("collection = ?", _query.Collection)
		hasWhere = true
	}
	if "" != _query.Element {
		db = db.Where("element = ?", _query.Element)
		hasWhere = true
	}
	if !hasWhere {
		return nil, ErrMemberNotFound
	}

	var member Member
	err := db.Limit(1).Find(&member).Error
	if member.UUID == "" {
		return nil, ErrMemberNotFound
	}
	return &member, err
}

func (this *MemberDAO) QueryMany(_query *MemberQuery) ([]*Member, error) {
	db := this.conn.DB.Model(&Member{})
	hasWhere := false
	if "" != _query.UUID {
		db = db.Where("uuid = ?", _query.UUID)
		hasWhere = true
	}
	if "" != _query.Collection {
		db = db.Where("collection = ?", _query.Collection)
		hasWhere = true
	}
	if "" != _query.Element {
		db = db.Where("element = ?", _query.Element)
		hasWhere = true
	}
	if !hasWhere {
		return nil, ErrMemberNotFound
	}

	var members []*Member
	err := db.Find(&members).Error
	return members, err
}
