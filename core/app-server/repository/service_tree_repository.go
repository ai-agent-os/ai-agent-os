package repository

import (
	"fmt"

	"github.com/ai-agent-os/ai-agent-os/core/app-server/model"
	"gorm.io/gorm"
)

type ServiceTreeRepository struct {
	db *gorm.DB
}

func NewServiceTreeRepository(db *gorm.DB) *ServiceTreeRepository {
	return &ServiceTreeRepository{db: db}
}

// CreateServiceTreeWithParentPath 创建服务目录
func (r *ServiceTreeRepository) CreateServiceTreeWithParentPath(serviceTree *model.ServiceTree, parentFullIDPath string) error {
	// 直接创建，不再计算FullIDPath
	return r.db.Create(serviceTree).Error
}

// CreateServiceTreeWithAppPrefix 创建带有用户应用前缀的服务目录
func (r *ServiceTreeRepository) CreateServiceTreeWithAppPrefix(serviceTree *model.ServiceTree, userAppPrefix string) error {
	// 先保存到数据库获取ID
	if err := r.db.Create(serviceTree).Error; err != nil {
		return err
	}

	// 然后计算包含用户应用前缀的路径信息
	if err := r.calculatePathsWithAppPrefix(serviceTree, userAppPrefix); err != nil {
		return fmt.Errorf("failed to calculate paths with app prefix: %w", err)
	}

	// 更新路径信息到数据库
	return r.db.Save(serviceTree).Error
}

// GetServiceTreeByID 根据ID获取服务目录
func (r *ServiceTreeRepository) GetServiceTreeByID(id int64) (*model.ServiceTree, error) {
	var serviceTree model.ServiceTree
	err := r.db.Where("id = ?", id).First(&serviceTree).Error
	if err != nil {
		return nil, err
	}
	return &serviceTree, nil
}

// GetServiceTreesByAppID 根据应用ID获取所有服务目录
func (r *ServiceTreeRepository) GetServiceTreesByAppID(appID int64) ([]*model.ServiceTree, error) {
	var serviceTrees []*model.ServiceTree
	err := r.db.Where("app_id = ?", appID).Order("created_at ASC").Find(&serviceTrees).Error
	if err != nil {
		return nil, err
	}
	return serviceTrees, nil
}

// GetServiceTreeChildren 获取子服务目录
func (r *ServiceTreeRepository) GetServiceTreeChildren(parentID int64) ([]*model.ServiceTree, error) {
	var children []*model.ServiceTree
	err := r.db.Where("parent_id = ?", parentID).Order("created_at ASC").Find(&children).Error
	if err != nil {
		return nil, err
	}
	return children, nil
}

// BuildServiceTree 构建树形结构
func (r *ServiceTreeRepository) BuildServiceTree(appID int64) ([]*model.ServiceTree, error) {
	// 获取所有服务目录
	allTrees, err := r.GetServiceTreesByAppID(appID)
	if err != nil {
		return nil, err
	}

	// 构建树形结构
	treeMap := make(map[int64]*model.ServiceTree)
	var rootNodes []*model.ServiceTree

	// 先创建所有节点的映射
	for _, tree := range allTrees {
		treeMap[tree.ID] = tree
	}

	// 构建父子关系
	for _, tree := range allTrees {
		if tree.ParentID == 0 {
			// 根节点
			rootNodes = append(rootNodes, tree)
		} else {
			// 子节点
			if parent, exists := treeMap[tree.ParentID]; exists {
				parent.Children = append(parent.Children, tree)
			}
		}
	}

	return rootNodes, nil
}

// UpdateServiceTree 更新服务目录
func (r *ServiceTreeRepository) UpdateServiceTree(serviceTree *model.ServiceTree) error {
	return r.db.Save(serviceTree).Error
}

// DeleteServiceTree 删除服务目录（级联删除子目录）
func (r *ServiceTreeRepository) DeleteServiceTree(id int64) error {
	// 先删除所有子目录
	children, err := r.GetServiceTreeChildren(id)
	if err != nil {
		return err
	}

	for _, child := range children {
		if err := r.DeleteServiceTree(child.ID); err != nil {
			return err
		}
	}

	// 删除当前目录
	return r.db.Delete(&model.ServiceTree{}, id).Error
}

// calculatePathsWithAppPrefix 计算带有用户应用前缀的路径信息
func (r *ServiceTreeRepository) calculatePathsWithAppPrefix(serviceTree *model.ServiceTree, userAppPrefix string) error {
	// FullCodePath使用预加载的app信息计算
	if serviceTree.App != nil {
		// 使用预加载的App对象构建路径
		appPrefix := fmt.Sprintf("/%s/%s", serviceTree.App.User, serviceTree.App.Code)
		serviceTree.FullCodePath = fmt.Sprintf("%s/%s", appPrefix, serviceTree.Code)
	} else {
		// 回退到传入的用户应用前缀
		serviceTree.FullCodePath = fmt.Sprintf("%s/%s", userAppPrefix, serviceTree.Code)
	}

	return nil
}

// GetServiceTreeByFullPath 根据完整路径获取服务目录（full_code_path全局唯一）
func (r *ServiceTreeRepository) GetServiceTreeByFullPath(fullPath string) (*model.ServiceTree, error) {
	var serviceTree model.ServiceTree
	err := r.db.Where("full_code_path = ?", fullPath).First(&serviceTree).Error
	if err != nil {
		return nil, err
	}
	return &serviceTree, nil
}

// GetServiceTreeByFullPaths 批量根据完整路径获取服务目录
func (r *ServiceTreeRepository) GetServiceTreeByFullPaths(fullPaths []string) (map[string]*model.ServiceTree, error) {
	if len(fullPaths) == 0 {
		return make(map[string]*model.ServiceTree), nil
	}

	var serviceTrees []*model.ServiceTree
	err := r.db.Where("full_code_path IN ?", fullPaths).Find(&serviceTrees).Error
	if err != nil {
		return nil, err
	}

	result := make(map[string]*model.ServiceTree)
	for _, tree := range serviceTrees {
		result[tree.FullCodePath] = tree
	}
	return result, nil
}

// CheckNameExists 检查名称是否已存在（在同一父目录下和同一应用内）
func (r *ServiceTreeRepository) CheckNameExists(parentID int64, code string, appID int64) (bool, error) {
	var count int64
	query := r.db.Model(&model.ServiceTree{}).Where("parent_id = ? AND code = ? AND app_id = ?", parentID, code, appID)

	err := query.Count(&count).Error
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func (r *ServiceTreeRepository) GetByID(parentId int64) (*model.ServiceTree, error) {
	var tree model.ServiceTree
	err := r.db.Model(&model.ServiceTree{}).Where("id = ?", parentId).First(&tree).Error
	if err != nil {
		return nil, err
	}
	return &tree, nil
}
