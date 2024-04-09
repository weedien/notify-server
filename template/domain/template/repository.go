package template

import "context"

type Repository interface {
	// Create 插入模板
	Create(ctx context.Context, template *EmailTemplate) error
	// Update 更新模板
	Update(ctx context.Context, template *EmailTemplate) error
	// Delete 删除模板
	Delete(ctx context.Context, id int64) error
}
