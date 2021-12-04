package model

import (
	"errors"
	"time"
)

type Collection struct {
	UUID      string `gorm:"column:uuid;type:char(32);not null;unique;primaryKey"`
	Name      string `gorm:"column:name;type:varchar(256);not null;unique"`
	Capacity  uint64 `gorm:"column:capacity;not null;default:0"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

var ErrCollectionExists = errors.New("collection exists")
var ErrCollectionNotFound = errors.New("collection not found")

func (Collection) TableName() string {
	return "ogm_group_collection"
}

type CollectionQuery struct {
	UUID string
}

type CollectionDAO struct {
	conn *Conn
}

func NewCollectionDAO(_conn *Conn) *CollectionDAO {
	conn := DefaultConn
	if nil != _conn {
		conn = _conn
	}
	return &CollectionDAO{
		conn: conn,
	}
}

func (this *CollectionDAO) Count() (int64, error) {
	var count int64
	err := this.conn.DB.Model(&Collection{}).Count(&count).Error
	return count, err
}

func (this *CollectionDAO) Insert(_Collection *Collection) error {
	var count int64
	err := this.conn.DB.Model(&Collection{}).Where("uuid = ? OR name = ?", _Collection.UUID, _Collection.Name).Count(&count).Error
	if nil != err {
		return err
	}

	if count > 0 {
		return ErrCollectionExists
	}

	return this.conn.DB.Create(_Collection).Error
}

func (this *CollectionDAO) Update(_Collection *Collection) error {
	return this.conn.DB.Updates(_Collection).Error
}

func (this *CollectionDAO) Delete(_uuid string) error {
	var count int64
	err := this.conn.DB.Model(&Collection{}).Where("uuid = ?", _uuid).Count(&count).Error
	if nil != err {
		return err
	}

	if 0 == count {
		return ErrCollectionNotFound
	}

	return this.conn.DB.Where("uuid = ?", _uuid).Delete(&Collection{}).Error
}

func (this *CollectionDAO) List(_offset int64, _count int64) (int64, []*Collection, error) {
	db := this.conn.DB.Model(&Collection{})
	var count int64
	res := db.Count(&count)
	if nil != res.Error {
		return 0, nil, res.Error
	}
	var Collections []*Collection
	res = db.Offset(int(_offset)).Limit(int(_count)).Order("created_at desc").Find(&Collections)
	return count, Collections, res.Error
}

func (this *CollectionDAO) Search(_offset int64, _count int64, _name string) (int64, []*Collection, error) {
	db := this.conn.DB.Model(&Collection{}).Where("name LIKE ?", "%"+_name+"%")
	var count int64
	res := db.Count(&count)
	if nil != res.Error {
		return 0, nil, res.Error
	}
	var Collections []*Collection
	res = db.Offset(int(_offset)).Limit(int(_count)).Order("created_at desc").Find(&Collections)
	return count, Collections, res.Error
}

func (this *CollectionDAO) QueryOne(_query *CollectionQuery) (*Collection, error) {
	db := this.conn.DB.Model(&Collection{})
	hasWhere := false
	if "" != _query.UUID {
		db = db.Where("uuid = ?", _query.UUID)
		hasWhere = true
	}
	// 没有where子句时，返回未找到错误
	if !hasWhere {
		return nil, ErrCollectionNotFound
	}

	var collection Collection
	err := db.Limit(1).Find(&collection).Error
	if collection.UUID == "" {
		return nil, ErrCollectionNotFound
	}
	return &collection, err
}
