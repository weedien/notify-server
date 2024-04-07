// Package template provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen/v2 version v2.0.0-20240331212514-80f0b978ef16 DO NOT EDIT.
package ports

import (
	"time"
)

// Template defines model for template.
type Template struct {
	Content string `json:"content"`

	// Slots 模板中需要填充的数据
	Slots []string `json:"slots"`
	Topic string   `json:"topic"`

	// Type 1不完整需填充数据2完整
	Type       int       `json:"type"`
	UpdateTime time.Time `json:"updateTime"`
}

// GetTemplatesParams defines parameters for GetTemplates.
type GetTemplatesParams struct {
	// Type 模板类型
	Type *int `form:"type,omitempty" json:"type,omitempty"`

	// ContentLen 界面中显示的content长度
	ContentLen *int `form:"contentLen,omitempty" json:"contentLen,omitempty"`
}

// CreateTemplateJSONBody defines parameters for CreateTemplate.
type CreateTemplateJSONBody struct {
	Content string `json:"content"`
	Topic   string `json:"topic"`
}

// UpdateTemplateJSONBody defines parameters for UpdateTemplate.
type UpdateTemplateJSONBody struct {
	Content *string `json:"content"`
	Topic   *string `json:"topic"`
}

// CreateTemplateJSONRequestBody defines body for CreateTemplate for application/json ContentType.
type CreateTemplateJSONRequestBody CreateTemplateJSONBody

// UpdateTemplateJSONRequestBody defines body for UpdateTemplate for application/json ContentType.
type UpdateTemplateJSONRequestBody UpdateTemplateJSONBody
