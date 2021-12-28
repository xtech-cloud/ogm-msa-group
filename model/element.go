package model

import (
	"errors"
	"time"
)

type Element struct {
	UUID       string `gorm:"column:uuid;type:char(32);not null;unique;primaryKey"`
	Collection string `gorm:"column:collection;type:char(32);not null;"`
	Key        string `gorm:"column:key;type:char(32);not null;"`
	Alias      string `gorm:"column:alias;type:varchar(128);not null;default:''"`
	Label      string `gorm:"column:label;type:varchar(1024);not null;default:''"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

var ErrElementExists = errors.New("element exists")
var ErrElementNotFound = errors.New("element not found")

func (Element) TableName() string {
	return "ogm_group_element"
}

type ElementDAO struct {
	conn *Conn
}

type ElementQuery struct {
	UUID       string
	Collection string
	Key        string
}

func NewElementDAO(_conn *Conn) *ElementDAO {
	conn := DefaultConn
	if nil != _conn {
		conn = _conn
	}
	return &ElementDAO{
		conn: conn,
	}
}

func (this *ElementDAO) Count() (int64, error) {
	var count int64
	err := this.conn.DB.Model(&Element{}).Count(&count).Error
	return count, err
}

func (this *ElementDAO) CountOfCollection(_collection string) (int64, error) {
	var count int64
	db := this.conn.DB.Model(&Element{})
	if "" != _collection {
		db = db.Where("collection = ?", _collection)
	}
	err := db.Count(&count).Error
	return count, err
}

func (this *ElementDAO) Insert(_element *Element) error {
	var count int64
	err := this.conn.DB.Model(&Element{}).Where("uuid = ?", _element.UUID).Count(&count).Error
	if nil != err {
		return err
	}

	if count > 0 {
		return ErrElementExists
	}

	return this.conn.DB.Create(_element).Error
}

func (this *ElementDAO) Update(_element *Element) error {
	return this.conn.DB.Updates(_element).Error
}

func (this *ElementDAO) Delete(_uuid string) error {
	var count int64
	err := this.conn.DB.Model(&Element{}).Where("uuid = ?", _uuid).Count(&count).Error
	if nil != err {
		return err
	}

	if 0 == count {
		return ErrElementNotFound
	}

	return this.conn.DB.Where("uuid = ?", _uuid).Delete(&Element{}).Error
}

func (this *ElementDAO) List(_offset int64, _count int64, _collection string) (int64, []*Element, error) {
	var elements []*Element
	db := this.conn.DB.Model(&Element{})
	if "" != _collection {
		db = db.Where("collection = ?", _collection)
	}
	var count int64
	err := db.Count(&count).Error
	if nil != err {
		return 0, nil, err
	}
	res := db.Offset(int(_offset)).Limit(int(_count)).Order("created_at desc").Find(&elements)
	return count, elements, res.Error
}

func (this *ElementDAO) Search(_offset int64, _count int64, _collection string, _key string, _alias string) (int64, []*Element, error) {
	var elements []*Element
	db := this.conn.DB.Model(&Element{})
	if "" != _collection {
		db = db.Where("`collection` = ?", _collection)
	}
	if "" != _key {
		db = db.Where("`key` LIKE ?", "%"+_key+"%")
	}
	if "" != _alias {
		db = db.Where("`alias` LIKE ?", "%"+_alias+"%")
	}
	var count int64
	err := db.Count(&count).Error
	if nil != err {
		return 0, nil, err
	}
	res := db.Offset(int(_offset)).Limit(int(_count)).Order("created_at desc").Find(&elements)
	return count, elements, res.Error
}

func (this *ElementDAO) QueryOne(_query *ElementQuery) (*Element, error) {
	db := this.conn.DB.Model(&Element{})
	hasWhere := false
	if "" != _query.UUID {
		db = db.Where("`uuid` = ?", _query.UUID)
		hasWhere = true
	}
	if "" != _query.Collection {
		db = db.Where("`collection` = ?", _query.Collection)
		hasWhere = true
	}
	if "" != _query.Key {
		db = db.Where("`key` = ?", _query.Key)
		hasWhere = true
	}
	if !hasWhere {
		return nil, ErrElementNotFound
	}

	var element Element
	err := db.Limit(1).Find(&element).Error
	if element.UUID == "" {
		return nil, ErrElementNotFound
	}
	return &element, err
}

func (this *ElementDAO) QueryMany(_query *ElementQuery) ([]*Element, error) {
	db := this.conn.DB.Model(&Element{})
	hasWhere := false
	if "" != _query.UUID {
		db = db.Where("uuid = ?", _query.UUID)
		hasWhere = true
	}
	if "" != _query.Collection {
		db = db.Where("collection = ?", _query.Collection)
		hasWhere = true
	}
	if "" != _query.Key {
		db = db.Where("`key` = ?", _query.Key)
		hasWhere = true
	}
	if !hasWhere {
		return nil, ErrElementNotFound
	}

	var elements []*Element
	err := db.Find(&elements).Error
	return elements, err
}
