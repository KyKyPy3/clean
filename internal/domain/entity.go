package domain

import (
	"github.com/KyKyPy3/clean/internal/domain/common"
)

type Entity interface {
	GetID() common.UID
	SetID(id common.UID)
}

type BaseEntity struct {
	id common.UID
}

func NewBaseEntity() *BaseEntity {
	return &BaseEntity{}
}

func (e *BaseEntity) GetID() common.UID {
	return e.id
}

func (e *BaseEntity) SetID(id common.UID) {
	e.id = id
}

var _ Entity = (*BaseEntity)(nil)
