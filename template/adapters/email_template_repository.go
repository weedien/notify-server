package adapters

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/weedien/notify-server/template/app/query"
	t "github.com/weedien/notify-server/template/domain/template"
	"gorm.io/gorm"
)

// EmailTemplateModel 实体类
type EmailTemplateModel struct {
	ID        int       `gorm:"primaryKey;autoIncrement"`
	Type      int       `gorm:"not null"`
	Topic     string    `gorm:"size:64;not null"`
	Slots     string    `gorm:"not null"`
	Content   string    `gorm:"not null"`
	CreatedAt time.Time `gorm:"not null"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
	Deleted   bool      `gorm:"not null"`
}

// TableName 指定表名
func (EmailTemplateModel) TableName() string {
	return "email_templates"
}

type EmailTemplateRepository struct {
	db *gorm.DB
}

func NewEmailTemplateRepository(db *gorm.DB) *EmailTemplateRepository {
	return &EmailTemplateRepository{db: db}
}

// Create 插入模板
func (r *EmailTemplateRepository) Create(ctx context.Context, template *t.EmailTemplate) error {
	model := EmailTemplateModel{
		Type:    template.Type(),
		Topic:   template.Topic(),
		Slots:   strings.Join(template.Slots(), ","),
		Content: template.Content().String,
	}
	return r.db.WithContext(ctx).Create(&model).Error
}

// Update 更新模板
func (r *EmailTemplateRepository) Update(ctx context.Context, template *t.EmailTemplate) error {
	// 查询模板是否存在
	var model EmailTemplateModel
	if err := r.db.WithContext(ctx).Where("deleted = ?", 0).First(&model, template.ID()).Error; err != nil {
		return errors.New("template not found")
	}

	q := r.db.WithContext(ctx).Model(&model).Where("id = ?", template.ID())
	if template.Topic() != "" {
		q = q.Update("topic", template.Topic())
	}
	if template.Content().Valid {
		q = q.Update("content", template.Content().String)
		q = q.Update("type", template.Type())
	}
	if len(template.Slots()) > 0 {
		q = q.Update("slots", strings.Join(template.Slots(), ","))
	}
	return q.Error
}

// 删除模板
func (r *EmailTemplateRepository) Delete(ctx context.Context, id int) error {
	// 逻辑删除
	return r.db.WithContext(ctx).Where("id = ?", id).Update("deleted", 1).Error
}

// // 通过ID获取模板
// func (r *EmailTemplateRepository) GetByID(ctx context.Context, id int) (*t.EmailTemplate, error) {
// 	var model EmailTemplateModel
// 	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&model).Error; err != nil {
// 		return nil, errors.New("template not found")
// 	}

// 	return t.NewEmailTemplateWithID(model.ID, model.Type, model.Topic, model.Content)
// }

// // 通过类型获取模板
// func (r *EmailTemplateRepository) GetByType(ctx context.Context, _type int) ([]t.EmailTemplate, error) {
// 	var models []EmailTemplateModel
// 	if err := r.db.WithContext(ctx).Where("type = ?", _type).Find(&models).Error; err != nil {
// 		return nil, errors.New("template not found")
// 	}

// 	var templates []t.EmailTemplate
// 	for _, model := range models {
// 		template, err := t.NewEmailTemplateWithID(model.ID, model.Type, model.Topic, model.Content)
// 		if err != nil {
// 			return nil, err
// 		}
// 		templates = append(templates, *template)
// 	}
// 	return templates, nil
// }

// // 获取所有模板
// func (r *EmailTemplateRepository) GetAll(ctx context.Context) ([]t.EmailTemplate, error) {
// 	var models []EmailTemplateModel
// 	if err := r.db.WithContext(ctx).Find(&models).Error; err != nil {
// 		return nil, errors.New("template not found")
// 	}

// 	var templates []t.EmailTemplate
// 	for _, model := range models {
// 		template, err := t.NewEmailTemplateWithID(model.ID, model.Type, model.Topic, model.Content)
// 		if err != nil {
// 			return nil, err
// 		}
// 		templates = append(templates, *template)
// 	}
// 	return templates, nil
// }

func (r *EmailTemplateRepository) FindTemplateByID(ctx context.Context, id int) (query.EmailTemplate, error) {
	var model EmailTemplateModel
	if err := r.db.WithContext(ctx).Where("deleted = ?", 0).Where("id = ?", id).First(&model).Error; err != nil {
		return query.EmailTemplate{}, err
	}

	return query.EmailTemplate{
		ID:      model.ID,
		Type:    model.Type,
		Topic:   model.Topic,
		Content: model.Content,
		Slots:   strings.Split(model.Slots, ","),
	}, nil
}

func (r *EmailTemplateRepository) AllTemplates(ctx context.Context) ([]query.EmailTemplate, error) {
	var models []EmailTemplateModel
	if err := r.db.WithContext(ctx).Where("delete = ?", 0).Find(&models).Error; err != nil {
		return nil, err
	}

	var templates []query.EmailTemplate
	for _, model := range models {
		template := query.EmailTemplate{
			ID:      model.ID,
			Type:    model.Type,
			Topic:   model.Topic,
			Content: model.Content,
			Slots:   strings.Split(model.Slots, ","),
		}
		templates = append(templates, template)
	}
	return templates, nil
}

func (r *EmailTemplateRepository) FindTemplatesByType(ctx context.Context, _type int) ([]query.EmailTemplate, error) {
	var models []EmailTemplateModel
	if err := r.db.WithContext(ctx).Where("type = ? and deleted = ?", _type, 0).Find(&models).Error; err != nil {
		return nil, err
	}

	var templates []query.EmailTemplate
	for _, model := range models {
		template := query.EmailTemplate{
			ID:      model.ID,
			Type:    model.Type,
			Topic:   model.Topic,
			Content: model.Content,
			Slots:   strings.Split(model.Slots, ","),
		}
		templates = append(templates, template)
	}
	return templates, nil
}
