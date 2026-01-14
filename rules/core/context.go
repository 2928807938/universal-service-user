package core

import "reflect"

// RuleContext 封装验证所需的输入数据
type RuleContext interface {
	// FieldName 返回当前验证的字段名
	FieldName() string

	// Value 返回当前字段的值
	Value() any

	// OriginalValue 返回原始数据对象
	OriginalValue() any

	// GetFieldValue 获取指定字段的值（用于跨字段验证）
	GetFieldValue(fieldName string) (any, bool)

	// Metadata 返回附加的元数据
	Metadata() map[string]any

	// WithMetadata 设置元数据并返回新的上下文
	WithMetadata(key string, value any) RuleContext
}

// DefaultRuleContext 是 RuleContext 的默认实现
type DefaultRuleContext struct {
	fieldName     string
	value         any
	originalValue any
	metadata      map[string]any
}

// NewRuleContext 创建一个新的规则上下文
func NewRuleContext(fieldName string, value any, originalValue any) *DefaultRuleContext {
	return &DefaultRuleContext{
		fieldName:     fieldName,
		value:         value,
		originalValue: originalValue,
		metadata:      make(map[string]any),
	}
}

func (c *DefaultRuleContext) FieldName() string {
	return c.fieldName
}

func (c *DefaultRuleContext) Value() any {
	return c.value
}

func (c *DefaultRuleContext) OriginalValue() any {
	return c.originalValue
}

func (c *DefaultRuleContext) GetFieldValue(fieldName string) (any, bool) {
	if c.originalValue == nil {
		return nil, false
	}

	val := reflect.ValueOf(c.originalValue)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	if val.Kind() != reflect.Struct {
		return nil, false
	}

	field := val.FieldByName(fieldName)
	if !field.IsValid() {
		return nil, false
	}

	return field.Interface(), true
}

func (c *DefaultRuleContext) Metadata() map[string]any {
	return c.metadata
}

func (c *DefaultRuleContext) WithMetadata(key string, value any) RuleContext {
	newCtx := &DefaultRuleContext{
		fieldName:     c.fieldName,
		value:         c.value,
		originalValue: c.originalValue,
		metadata:      make(map[string]any),
	}
	for k, v := range c.metadata {
		newCtx.metadata[k] = v
	}
	newCtx.metadata[key] = value
	return newCtx
}
