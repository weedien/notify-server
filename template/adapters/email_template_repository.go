package adapters

import (
	"context"
	"errors"
	"strings"
	"time"

	q "github.com/weedien/notify-server/template/app/query"
	t "github.com/weedien/notify-server/template/domain/template"
	"gorm.io/gorm"
)

// EmailTemplateModel 实体类
type EmailTemplateModel struct {
	ID        int64     `gorm:"size:64;primaryKey;autoIncrement"`
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
	return "notify_templates"
}

type EmailTemplateRepository struct {
	db *gorm.DB
}

func NewEmailTemplateRepository(db *gorm.DB) *EmailTemplateRepository {
	return &EmailTemplateRepository{db: db}
}

// Create 插入模板
func (r *EmailTemplateRepository) Create(ctx context.Context, template *t.EmailTemplate) error {
	// 查询数据库中是否存在同名的topic
	var count int64
	if err := r.db.WithContext(ctx).Model(&EmailTemplateModel{}).Where("topic = ?", template.Topic()).Count(&count).Error; err != nil {
		return err
	}
	if count > 0 {
		return errors.New("topic already exists")
	}

	model := EmailTemplateModel{
		Type:    template.Type(),
		Topic:   template.Topic(),
		Slots:   strings.Join(template.Content().Slots(), ","),
		Content: template.Content().String(),
	}
	return r.db.WithContext(ctx).Create(&model).Error
}

// Update 更新模板
func (r *EmailTemplateRepository) Update(ctx context.Context, template *t.EmailTemplate) error {
	// 查询模板是否存在
	var count int64
	err := r.db.WithContext(ctx).Model(&EmailTemplateModel{}).Where("id = ? AND deleted = ?", template.ID(), 0).Count(&count).Error
	if err != nil {
		return err
	}

	if count == 0 {
		return errors.New("template not found")
	}
	// 查询数据库中是否存在同名的topic
	if template.Topic() != "" {
		if err := r.db.WithContext(ctx).Model(&EmailTemplateModel{}).Where("topic = ?", template.Topic()).Count(&count).Error; err != nil {
			return err
		}
		if count > 0 {
			return errors.New("topic already exists")
		}
	}

	chain := r.db.WithContext(ctx).Model(&EmailTemplateModel{}).Where("id = ?", template.ID())
	if template.Topic() != "" {
		chain = chain.Update("topic", template.Topic())
	}
	if template.Content() != nil {
		chain = chain.Update("content", template.Content().String())
		chain = chain.Update("type", template.Type())
		chain = chain.Update("slots", strings.Join(template.Content().Slots(), ","))
	}
	return chain.Error
}

// Delete 删除模板
func (r *EmailTemplateRepository) Delete(ctx context.Context, id int64) error {
	// 逻辑删除
	return r.db.WithContext(ctx).Where("id = ?", id).Update("deleted", 1).Error
}

// // 通过ID获取模板
// func (r *EmailTemplateRepository) GetByID(ctx context.Context, id int64) (*t.EmailTemplate, error) {
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

func (r *EmailTemplateRepository) FindTemplateByID(ctx context.Context, id int64) (q.EmailTemplate, error) {
	var model EmailTemplateModel
	if err := r.db.WithContext(ctx).Where("deleted = ?", 0).Where("id = ?", id).First(&model).Error; err != nil {
		return q.EmailTemplate{}, err
	}

	return q.EmailTemplate{
		ID:         model.ID,
		Type:       model.Type,
		Topic:      model.Topic,
		Content:    model.Content,
		Slots:      strings.Split(model.Slots, ","),
		UpdateTime: model.UpdatedAt,
	}, nil
}

func (r *EmailTemplateRepository) Templates(ctx context.Context, query q.TemplatesQuery) ([]q.EmailTemplate, error) {
	var models []EmailTemplateModel

	chain := r.db.WithContext(ctx).Where("deleted = ?", 0)

	if query.Type != nil {
		chain = chain.Where("type = ?", *query.Type)
	}

	if err := chain.Find(&models).Error; err != nil {
		return nil, err
	}

	var templates []q.EmailTemplate
	for _, model := range models {
		// 对内容进行截取
		trimContent := model.Content
		if query.ContentLen != nil && *query.ContentLen > 0 {
			if len(model.Content) > *query.ContentLen {
				trimContent = model.Content[:*query.ContentLen]
			}
		}
		template := q.EmailTemplate{
			ID:         model.ID,
			Type:       model.Type,
			Topic:      model.Topic,
			Content:    trimContent,
			Slots:      strings.Split(model.Slots, ","),
			UpdateTime: model.UpdatedAt,
		}
		templates = append(templates, template)
	}
	return templates, nil
}
