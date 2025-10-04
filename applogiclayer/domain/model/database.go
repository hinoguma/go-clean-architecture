package model

type UpdateValueType string

const (
	UpdateValueTypeSet    UpdateValueType = "SET"
	UpdateValueTypeAdd    UpdateValueType = "ADD"
	UpdateValueTypeRemove UpdateValueType = "REMOVE"
	UpdateValueTypeDelete UpdateValueType = "DELETE"
)

type UpdateValueRequest struct {
	Type  UpdateValueType
	Field string
	Value any
}

func (model UpdateValueRequest) IsSet() bool {
	return model.Type == UpdateValueTypeSet
}
func (model UpdateValueRequest) IsAdd() bool {
	return model.Type == UpdateValueTypeAdd
}
func (model UpdateValueRequest) IsRemove() bool {
	return model.Type == UpdateValueTypeRemove
}
func (model UpdateValueRequest) IsDelete() bool {
	return model.Type == UpdateValueTypeDelete
}

type UpdateItemRequest struct {
	KeyCondition ItemCondition
	Values       []UpdateValueRequest
	Condition    ItemCondition
}

func (model UpdateItemRequest) HasCondition() bool {
	if model.Condition == nil {
		return false
	}
	return true
}

type ComparisonOperator string

const (
	ComparisonOperatorEqual ComparisonOperator = "EQUAL"
)

type ItemConditionType string

const (
	ItemConditionTypeSingle ItemConditionType = "SINGLE"
	ItemConditionTypeAnd    ItemConditionType = "AND"
	ItemConditionTypeOr     ItemConditionType = "OR"
)

type ItemCondition interface {
	GetType() ItemConditionType
}

func (model ItemConditionSingle) GetType() ItemConditionType {
	return ItemConditionTypeSingle
}

func (model ItemConditionAnd) GetType() ItemConditionType {
	return ItemConditionTypeAnd
}

func (model ItemConditionOr) GetType() ItemConditionType {
	return ItemConditionTypeOr
}

type ItemConditionSingle struct {
	Operator ComparisonOperator
	Field    string
	Value    any
}

type ItemConditionAnd struct {
	BaseItemConditionGroup
}

type ItemConditionOr struct {
	BaseItemConditionGroup
}

type BaseItemConditionGroup struct {
	conditions []ItemCondition
}

func (model BaseItemConditionGroup) Conditions() []ItemCondition {
	if model.conditions == nil {
		return make([]ItemCondition, 0)
	}
	return model.conditions
}

func (model *BaseItemConditionGroup) AddCondition(cond ...ItemCondition) {
	model.conditions = append(model.conditions, cond...)
}

/**
Transaction
*/

type TxOptions struct {
	TransactionID  string
	IsolationLevel string
	ReadOnly       bool
}
