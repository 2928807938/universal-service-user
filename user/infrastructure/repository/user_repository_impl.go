package repository

import (
	"context"

	"gorm.io/gorm"

	"universal-service-user/share/repository"
	basegorm "universal-service-user/share/repository/gorm"
	"universal-service-user/user/domain/entity"
	domainRepo "universal-service-user/user/domain/repository"
	"universal-service-user/user/infrastructure/converter"
	infraEntity "universal-service-user/user/infrastructure/entity"
)

// UserRepositoryImpl 用户仓储实现
type UserRepositoryImpl struct {
	repo      *basegorm.QueryableGormRepository[infraEntity.UserPO, int]
	converter *converter.UserConverter
}

// NewUserRepositoryImpl 创建用户仓储实现
func NewUserRepositoryImpl(db *gorm.DB) domainRepo.UserRepository {
	return &UserRepositoryImpl{
		repo:      basegorm.NewQueryableGormRepository[infraEntity.UserPO, int](db),
		converter: converter.NewUserConverter(),
	}
}

// Create 创建用户（实现 BaseRepository）
func (r *UserRepositoryImpl) Create(ctx context.Context, user *entity.User) error {
	po := r.converter.ToPO(user)
	return r.repo.Create(ctx, po)
}

// CreateBatch 批量创建用户（实现 BaseRepository）
func (r *UserRepositoryImpl) CreateBatch(ctx context.Context, users []*entity.User) error {
	if len(users) == 0 {
		return nil
	}
	pos := make([]*infraEntity.UserPO, len(users))
	for i, user := range users {
		pos[i] = r.converter.ToPO(user)
	}
	return r.repo.CreateBatch(ctx, pos)
}

// GetByID 根据 ID 查找用户（实现 BaseRepository）
func (r *UserRepositoryImpl) GetByID(ctx context.Context, id int) (*entity.User, error) {
	po, err := r.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if po == nil {
		return nil, nil
	}
	return r.converter.ToEntity(po), nil
}

// FindByID 根据 ID 查找用户（别名方法，保持兼容）
func (r *UserRepositoryImpl) FindByID(ctx context.Context, id int) (*entity.User, error) {
	return r.GetByID(ctx, id)
}

// FindByEmail 根据邮箱查找用户
func (r *UserRepositoryImpl) FindByEmail(ctx context.Context, email string) (*entity.User, error) {
	poList, err := r.repo.Where(ctx, repository.Eq("email", email))
	if err != nil {
		return nil, err
	}
	if len(poList) == 0 {
		return nil, nil
	}
	return r.converter.ToEntity(poList[0]), nil
}

// FindByUsername 根据用户名查找用户
func (r *UserRepositoryImpl) FindByUsername(ctx context.Context, username string) (*entity.User, error) {
	poList, err := r.repo.Where(ctx, repository.Eq("username", username))
	if err != nil {
		return nil, err
	}
	if len(poList) == 0 {
		return nil, nil
	}
	return r.converter.ToEntity(poList[0]), nil
}

// Update 更新用户（实现 BaseRepository）
func (r *UserRepositoryImpl) Update(ctx context.Context, user *entity.User) error {
	po := r.converter.ToPO(user)
	return r.repo.Update(ctx, po)
}

// Delete 删除用户（实现 BaseRepository）
func (r *UserRepositoryImpl) Delete(ctx context.Context, id int) error {
	return r.repo.Delete(ctx, id)
}

// List 查询全部用户列表（实现 BaseRepository）
func (r *UserRepositoryImpl) List(ctx context.Context) ([]*entity.User, error) {
	pos, err := r.repo.List(ctx)
	if err != nil {
		return nil, err
	}
	users := make([]*entity.User, len(pos))
	for i, po := range pos {
		users[i] = r.converter.ToEntity(po)
	}
	return users, nil
}

// Page 分页查询（实现 BaseRepository）
func (r *UserRepositoryImpl) Page(ctx context.Context, request *repository.PageRequest) (*repository.PageResult[*entity.User], error) {
	poResult, err := r.repo.Page(ctx, request)
	if err != nil {
		return nil, err
	}
	// 转换为领域实体
	items := make([]*entity.User, len(poResult.Items))
	for i, po := range poResult.Items {
		items[i] = r.converter.ToEntity(po)
	}
	return repository.NewPageResult(items, poResult.Total, poResult.Page, poResult.Size), nil
}

// Where 条件查询（实现 QueryableRepository）
func (r *UserRepositoryImpl) Where(ctx context.Context, conditions ...*repository.Condition) ([]*entity.User, error) {
	poList, err := r.repo.Where(ctx, conditions...)
	if err != nil {
		return nil, err
	}
	users := make([]*entity.User, len(poList))
	for i, po := range poList {
		users[i] = r.converter.ToEntity(po)
	}
	return users, nil
}

// Count 统计数量（实现 QueryableRepository）
func (r *UserRepositoryImpl) Count(ctx context.Context, conditions ...*repository.Condition) (int64, error) {
	return r.repo.Count(ctx, conditions...)
}

// Exists 存在性检查（实现 QueryableRepository）
func (r *UserRepositoryImpl) Exists(ctx context.Context, conditions ...*repository.Condition) (bool, error) {
	return r.repo.Exists(ctx, conditions...)
}

// Query 获取查询构建器（实现 QueryableRepository）
func (r *UserRepositoryImpl) Query() repository.QueryBuilder[entity.User] {
	return NewUserQueryBuilder(r.repo.Query(), r.converter)
}

// ExistsByEmail 检查邮箱是否存在
func (r *UserRepositoryImpl) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	return r.repo.Exists(ctx, repository.Eq("email", email))
}

// ExistsByUsername 检查用户名是否存在
func (r *UserRepositoryImpl) ExistsByUsername(ctx context.Context, username string) (bool, error) {
	return r.repo.Exists(ctx, repository.Eq("username", username))
}

// FindByPhone 根据手机号查找用户
func (r *UserRepositoryImpl) FindByPhone(ctx context.Context, phone string) (*entity.User, error) {
	poList, err := r.repo.Where(ctx, repository.Eq("phone", phone))
	if err != nil {
		return nil, err
	}
	if len(poList) == 0 {
		return nil, nil
	}
	return r.converter.ToEntity(poList[0]), nil
}

// ExistsByPhone 检查手机号是否存在
func (r *UserRepositoryImpl) ExistsByPhone(ctx context.Context, phone string) (bool, error) {
	return r.repo.Exists(ctx, repository.Eq("phone", phone))
}

// UserQueryBuilder 用户查询构建器（包装 PO 构建器，自动转换）
type UserQueryBuilder struct {
	poBuilder repository.QueryBuilder[infraEntity.UserPO]
	converter *converter.UserConverter
}

// NewUserQueryBuilder 创建用户查询构建器
func NewUserQueryBuilder(poBuilder repository.QueryBuilder[infraEntity.UserPO], converter *converter.UserConverter) *UserQueryBuilder {
	return &UserQueryBuilder{
		poBuilder: poBuilder,
		converter: converter,
	}
}

func (b *UserQueryBuilder) Where(condition *repository.Condition) repository.QueryBuilder[entity.User] {
	b.poBuilder.Where(condition)
	return b
}

func (b *UserQueryBuilder) And(conditions ...*repository.Condition) repository.QueryBuilder[entity.User] {
	b.poBuilder.And(conditions...)
	return b
}

func (b *UserQueryBuilder) OrderBy(field string) repository.QueryBuilder[entity.User] {
	b.poBuilder.OrderBy(field)
	return b
}

func (b *UserQueryBuilder) OrderByDesc(field string) repository.QueryBuilder[entity.User] {
	b.poBuilder.OrderByDesc(field)
	return b
}

func (b *UserQueryBuilder) Limit(limit int) repository.QueryBuilder[entity.User] {
	b.poBuilder.Limit(limit)
	return b
}

func (b *UserQueryBuilder) Offset(offset int) repository.QueryBuilder[entity.User] {
	b.poBuilder.Offset(offset)
	return b
}

func (b *UserQueryBuilder) Select(fields ...string) repository.QueryBuilder[entity.User] {
	b.poBuilder.Select(fields...)
	return b
}

func (b *UserQueryBuilder) Find(ctx context.Context) ([]*entity.User, error) {
	pos, err := b.poBuilder.Find(ctx)
	if err != nil {
		return nil, err
	}
	users := make([]*entity.User, len(pos))
	for i, po := range pos {
		users[i] = b.converter.ToEntity(po)
	}
	return users, nil
}

func (b *UserQueryBuilder) First(ctx context.Context) (*entity.User, error) {
	po, err := b.poBuilder.First(ctx)
	if err != nil || po == nil {
		return nil, err
	}
	return b.converter.ToEntity(po), nil
}

func (b *UserQueryBuilder) Count(ctx context.Context) (int64, error) {
	return b.poBuilder.Count(ctx)
}

func (b *UserQueryBuilder) Exists(ctx context.Context) (bool, error) {
	return b.poBuilder.Exists(ctx)
}
