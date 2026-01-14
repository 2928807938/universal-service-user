package repository

import (
	"context"

	"gorm.io/gorm"

	"universal-service-user/share/repository"
	basegorm "universal-service-user/share/repository/gorm"
	"universal-service-user/oauth/domain/entity"
	domainRepo "universal-service-user/oauth/domain/repository"
	"universal-service-user/oauth/infrastructure/converter"
	infraEntity "universal-service-user/oauth/infrastructure/entity"
)

// OAuthBindingRepositoryImpl OAuth 绑定仓储实现
type OAuthBindingRepositoryImpl struct {
	repo      *basegorm.QueryableGormRepository[infraEntity.OAuthBindingPO, int]
	converter *converter.OAuthBindingConverter
}

// NewOAuthBindingRepositoryImpl 创建 OAuth 绑定仓储实现
func NewOAuthBindingRepositoryImpl(db *gorm.DB) domainRepo.OAuthBindingRepository {
	return &OAuthBindingRepositoryImpl{
		repo:      basegorm.NewQueryableGormRepository[infraEntity.OAuthBindingPO, int](db),
		converter: converter.NewOAuthBindingConverter(),
	}
}

// Create 创建绑定 (实现 BaseRepository)
func (r *OAuthBindingRepositoryImpl) Create(ctx context.Context, binding *entity.OAuthBinding) error {
	po := r.converter.ToPO(binding)
	return r.repo.Create(ctx, po)
}

// CreateBatch 批量创建绑定 (实现 BaseRepository)
func (r *OAuthBindingRepositoryImpl) CreateBatch(ctx context.Context, bindings []*entity.OAuthBinding) error {
	if len(bindings) == 0 {
		return nil
	}
	pos := r.converter.ToPOs(bindings)
	return r.repo.CreateBatch(ctx, pos)
}

// GetByID 根据 ID 查找绑定 (实现 BaseRepository)
func (r *OAuthBindingRepositoryImpl) GetByID(ctx context.Context, id int) (*entity.OAuthBinding, error) {
	po, err := r.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if po == nil {
		return nil, nil
	}
	return r.converter.ToEntity(po), nil
}

// Update 更新绑定 (实现 BaseRepository)
func (r *OAuthBindingRepositoryImpl) Update(ctx context.Context, binding *entity.OAuthBinding) error {
	po := r.converter.ToPO(binding)
	return r.repo.Update(ctx, po)
}

// Delete 删除绑定 (实现 BaseRepository)
func (r *OAuthBindingRepositoryImpl) Delete(ctx context.Context, id int) error {
	return r.repo.Delete(ctx, id)
}

// List 查询全部绑定列表 (实现 BaseRepository)
func (r *OAuthBindingRepositoryImpl) List(ctx context.Context) ([]*entity.OAuthBinding, error) {
	pos, err := r.repo.List(ctx)
	if err != nil {
		return nil, err
	}
	return r.converter.ToEntities(pos), nil
}

// Page 分页查询 (实现 BaseRepository)
func (r *OAuthBindingRepositoryImpl) Page(ctx context.Context, request *repository.PageRequest) (*repository.PageResult[*entity.OAuthBinding], error) {
	poResult, err := r.repo.Page(ctx, request)
	if err != nil {
		return nil, err
	}
	// 转换为领域实体
	items := r.converter.ToEntities(poResult.Items)
	return repository.NewPageResult(items, poResult.Total, poResult.Page, poResult.Size), nil
}

// Where 条件查询 (实现 QueryableRepository)
func (r *OAuthBindingRepositoryImpl) Where(ctx context.Context, conditions ...*repository.Condition) ([]*entity.OAuthBinding, error) {
	poList, err := r.repo.Where(ctx, conditions...)
	if err != nil {
		return nil, err
	}
	return r.converter.ToEntities(poList), nil
}

// Count 统计数量 (实现 QueryableRepository)
func (r *OAuthBindingRepositoryImpl) Count(ctx context.Context, conditions ...*repository.Condition) (int64, error) {
	return r.repo.Count(ctx, conditions...)
}

// Exists 存在性检查 (实现 QueryableRepository)
func (r *OAuthBindingRepositoryImpl) Exists(ctx context.Context, conditions ...*repository.Condition) (bool, error) {
	return r.repo.Exists(ctx, conditions...)
}

// Query 获取查询构建器 (实现 QueryableRepository)
func (r *OAuthBindingRepositoryImpl) Query() repository.QueryBuilder[entity.OAuthBinding] {
	return NewOAuthBindingQueryBuilder(r.repo.Query(), r.converter)
}

// FindByUserIDAndProvider 根据用户 ID 和平台查询绑定
func (r *OAuthBindingRepositoryImpl) FindByUserIDAndProvider(ctx context.Context, userID int, provider string) (*entity.OAuthBinding, error) {
	poList, err := r.repo.Where(ctx,
		repository.Eq("user_id", userID),
		repository.Eq("provider", provider),
	)
	if err != nil {
		return nil, err
	}
	if len(poList) == 0 {
		return nil, nil
	}
	return r.converter.ToEntity(poList[0]), nil
}

// FindByProviderAndOpenID 根据平台和 OpenID 查询绑定
func (r *OAuthBindingRepositoryImpl) FindByProviderAndOpenID(ctx context.Context, provider, openID string) (*entity.OAuthBinding, error) {
	poList, err := r.repo.Where(ctx,
		repository.Eq("provider", provider),
		repository.Eq("open_id", openID),
	)
	if err != nil {
		return nil, err
	}
	if len(poList) == 0 {
		return nil, nil
	}
	return r.converter.ToEntity(poList[0]), nil
}

// FindByProviderAndUnionID 根据平台和 UnionID 查询绑定
func (r *OAuthBindingRepositoryImpl) FindByProviderAndUnionID(ctx context.Context, provider, unionID string) (*entity.OAuthBinding, error) {
	if unionID == "" {
		return nil, nil
	}
	poList, err := r.repo.Where(ctx,
		repository.Eq("provider", provider),
		repository.Eq("union_id", unionID),
	)
	if err != nil {
		return nil, err
	}
	if len(poList) == 0 {
		return nil, nil
	}
	return r.converter.ToEntity(poList[0]), nil
}

// FindByUserID 根据用户 ID 查询所有绑定
func (r *OAuthBindingRepositoryImpl) FindByUserID(ctx context.Context, userID int) ([]*entity.OAuthBinding, error) {
	poList, err := r.repo.Where(ctx, repository.Eq("user_id", userID))
	if err != nil {
		return nil, err
	}
	return r.converter.ToEntities(poList), nil
}

// DeleteByUserIDAndProvider 根据用户 ID 和平台删除绑定
func (r *OAuthBindingRepositoryImpl) DeleteByUserIDAndProvider(ctx context.Context, userID int, provider string) error {
	// 先查找绑定
	binding, err := r.FindByUserIDAndProvider(ctx, userID, provider)
	if err != nil {
		return err
	}
	if binding == nil {
		return nil // 绑定不存在,视为删除成功
	}
	// 删除绑定
	return r.repo.Delete(ctx, binding.ID)
}

// ExistsByProviderAndOpenID 判断绑定是否存在
func (r *OAuthBindingRepositoryImpl) ExistsByProviderAndOpenID(ctx context.Context, provider, openID string) (bool, error) {
	return r.repo.Exists(ctx,
		repository.Eq("provider", provider),
		repository.Eq("open_id", openID),
	)
}

// OAuthBindingQueryBuilder OAuth 绑定查询构建器 (包装 PO 构建器,自动转换)
type OAuthBindingQueryBuilder struct {
	poBuilder repository.QueryBuilder[infraEntity.OAuthBindingPO]
	converter *converter.OAuthBindingConverter
}

// NewOAuthBindingQueryBuilder 创建 OAuth 绑定查询构建器
func NewOAuthBindingQueryBuilder(poBuilder repository.QueryBuilder[infraEntity.OAuthBindingPO], converter *converter.OAuthBindingConverter) *OAuthBindingQueryBuilder {
	return &OAuthBindingQueryBuilder{
		poBuilder: poBuilder,
		converter: converter,
	}
}

func (b *OAuthBindingQueryBuilder) Where(condition *repository.Condition) repository.QueryBuilder[entity.OAuthBinding] {
	b.poBuilder.Where(condition)
	return b
}

func (b *OAuthBindingQueryBuilder) And(conditions ...*repository.Condition) repository.QueryBuilder[entity.OAuthBinding] {
	b.poBuilder.And(conditions...)
	return b
}

func (b *OAuthBindingQueryBuilder) OrderBy(field string) repository.QueryBuilder[entity.OAuthBinding] {
	b.poBuilder.OrderBy(field)
	return b
}

func (b *OAuthBindingQueryBuilder) OrderByDesc(field string) repository.QueryBuilder[entity.OAuthBinding] {
	b.poBuilder.OrderByDesc(field)
	return b
}

func (b *OAuthBindingQueryBuilder) Limit(limit int) repository.QueryBuilder[entity.OAuthBinding] {
	b.poBuilder.Limit(limit)
	return b
}

func (b *OAuthBindingQueryBuilder) Offset(offset int) repository.QueryBuilder[entity.OAuthBinding] {
	b.poBuilder.Offset(offset)
	return b
}

func (b *OAuthBindingQueryBuilder) Select(fields ...string) repository.QueryBuilder[entity.OAuthBinding] {
	b.poBuilder.Select(fields...)
	return b
}

func (b *OAuthBindingQueryBuilder) Find(ctx context.Context) ([]*entity.OAuthBinding, error) {
	pos, err := b.poBuilder.Find(ctx)
	if err != nil {
		return nil, err
	}
	return b.converter.ToEntities(pos), nil
}

func (b *OAuthBindingQueryBuilder) First(ctx context.Context) (*entity.OAuthBinding, error) {
	po, err := b.poBuilder.First(ctx)
	if err != nil || po == nil {
		return nil, err
	}
	return b.converter.ToEntity(po), nil
}

func (b *OAuthBindingQueryBuilder) Count(ctx context.Context) (int64, error) {
	return b.poBuilder.Count(ctx)
}

func (b *OAuthBindingQueryBuilder) Exists(ctx context.Context) (bool, error) {
	return b.poBuilder.Exists(ctx)
}
