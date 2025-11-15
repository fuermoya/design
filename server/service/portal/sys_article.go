package portal

import (
	"errors"
	"strconv"
	"time"

	"github.com/fuermoya/design/server/global"
	"github.com/fuermoya/design/server/model/portal"
	"github.com/fuermoya/design/server/model/portal/request"
	"gorm.io/gorm"
)

type SysArticleService struct{}

// CreateSysArticle 创建文章
func (sysArticleService *SysArticleService) CreateSysArticle(sysArticle request.SysArticleCreate) (err error) {
	var article portal.SysArticle
	article.Title = sysArticle.Title
	article.Content = sysArticle.Content
	article.Summary = sysArticle.Summary
	article.CoverImage = sysArticle.CoverImage
	article.CategoryID = sysArticle.CategoryID
	article.Status = sysArticle.Status
	article.Sort = sysArticle.Sort

	err = global.DB.Transaction(func(tx *gorm.DB) error {
		// 创建文章
		if err := tx.Create(&article).Error; err != nil {
			return err
		}

		// 关联标签
		if len(sysArticle.TagIDs) > 0 {
			var tags []portal.SysTag
			if err := tx.Where("id IN ?", sysArticle.TagIDs).Find(&tags).Error; err != nil {
				return err
			}
			if err := tx.Model(&article).Association("Tags").Replace(tags); err != nil {
				return err
			}
		}

		return nil
	})

	return err
}

// DeleteSysArticle 删除文章
func (sysArticleService *SysArticleService) DeleteSysArticle(id string) (err error) {
	err = global.DB.Delete(&portal.SysArticle{}, "id = ?", id).Error
	return err
}

// DeleteSysArticleByIds 批量删除文章
func (sysArticleService *SysArticleService) DeleteSysArticleByIds(ids []uint) (err error) {
	err = global.DB.Delete(&[]portal.SysArticle{}, "id in ?", ids).Error
	return err
}

// UpdateSysArticle 更新文章
func (sysArticleService *SysArticleService) UpdateSysArticle(sysArticle request.SysArticleUpdate) (err error) {
	var article portal.SysArticle
	if err := global.DB.Where("id = ?", sysArticle.ID).First(&article).Error; err != nil {
		return err
	}

	article.Title = sysArticle.Title
	article.Content = sysArticle.Content
	article.Summary = sysArticle.Summary
	article.CoverImage = sysArticle.CoverImage
	article.CategoryID = sysArticle.CategoryID
	article.Status = sysArticle.Status
	article.Sort = sysArticle.Sort

	err = global.DB.Transaction(func(tx *gorm.DB) error {
		// 更新文章
		if err := tx.Save(&article).Error; err != nil {
			return err
		}

		// 更新标签关联
		if len(sysArticle.TagIDs) > 0 {
			var tags []portal.SysTag
			if err := tx.Where("id IN ?", sysArticle.TagIDs).Find(&tags).Error; err != nil {
				return err
			}
			if err := tx.Model(&article).Association("Tags").Replace(tags); err != nil {
				return err
			}
		} else {
			// 清空标签关联
			if err := tx.Model(&article).Association("Tags").Clear(); err != nil {
				return err
			}
		}

		return nil
	})

	return err
}

// GetSysArticle 根据id获取文章
func (sysArticleService *SysArticleService) GetSysArticle(id string) (sysArticle portal.SysArticle, err error) {
	err = global.DB.Preload("Category").Preload("Tags").Preload("Author").Where("id = ?", id).First(&sysArticle).Error
	return
}

// GetSysArticleInfoList 分页获取文章列表
func (sysArticleService *SysArticleService) GetSysArticleInfoList(info request.SysArticleSearch) (list []portal.SysArticle, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	db := global.DB.Model(&portal.SysArticle{})
	var sysArticles []portal.SysArticle

	// 添加查询条件
	if info.CategoryID != nil {
		db = db.Where("category_id = ?", *info.CategoryID)
	}
	if info.Status != nil {
		db = db.Where("status = ?", *info.Status)
	}
	if info.Keyword != "" {
		db = db.Where("title LIKE ? OR content LIKE ?", "%"+info.Keyword+"%", "%"+info.Keyword+"%")
	}

	err = db.Count(&total).Error
	if err != nil {
		return
	}

	err = db.Preload("Category").Preload("Tags").Preload("Author").
		Limit(limit).Offset(offset).Order("sort DESC, created_at DESC").Find(&sysArticles).Error
	if err != nil {
		return
	}

	// 为每个文章计算浏览量和点赞数
	for i := range sysArticles {
		articleID := strconv.FormatUint(uint64(sysArticles[i].ID), 10)
		viewCount, _ := sysArticleService.GetArticleViewCount(articleID)
		likeCount, _ := sysArticleService.GetArticleLikeCount(articleID)
		sysArticles[i].ViewCount = viewCount
		sysArticles[i].LikeCount = likeCount
	}

	return sysArticles, total, err
}

// GetPublishedArticles 获取已发布的文章列表（前台使用）
func (sysArticleService *SysArticleService) GetPublishedArticles(info request.SysArticleSearch) (list []portal.SysArticle, total int64, err error) {

	db := global.DB.Model(&portal.SysArticle{}).Where("status = 1") // 1表示已发布
	if info.CategoryID != nil {
		db = db.Where("category_id = ?", *info.CategoryID)
	}
	err = db.Count(&total).Error
	if err != nil {
		return
	}

	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)

	err = db.Preload("Category").Preload("Tags").Preload("Author").
		Limit(limit).Offset(offset).Order("publish_time DESC, created_at DESC").Find(&list).Error
	if err != nil {
		return
	}

	// 为每个文章计算浏览量和点赞数
	for i := range list {
		articleID := strconv.FormatUint(uint64(list[i].ID), 10)
		viewCount, _ := sysArticleService.GetArticleViewCount(articleID)
		likeCount, _ := sysArticleService.GetArticleLikeCount(articleID)
		list[i].ViewCount = viewCount
		list[i].LikeCount = likeCount
	}

	return
}

// IncrementViewCount 增加浏览次数
func (sysArticleService *SysArticleService) IncrementViewCount(id string, ip string, userAgent string) error {
	articleID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		return err
	}

	// 检查是否已经记录过浏览（按IP地址检查）
	var existingView portal.SysArticleView
	err = global.DB.Where("article_id = ? AND ip = ?", articleID, ip).First(&existingView).Error
	if err == nil {
		// 已经记录过，更新访问时间
		return global.DB.Model(&existingView).Update("updated_at", time.Now()).Error
	}

	// 创建新的浏览记录
	viewRecord := portal.SysArticleView{
		ArticleID: uint(articleID),
		IP:        ip,
		UserAgent: userAgent,
	}
	return global.DB.Create(&viewRecord).Error
}

// LikeArticle 点赞文章
func (sysArticleService *SysArticleService) LikeArticle(id string, ip string) error {
	articleID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		return err
	}

	// 检查是否已经点赞（按IP地址检查）
	var existingLike portal.SysArticleLike
	err = global.DB.Where("article_id = ? AND ip = ?", articleID, ip).First(&existingLike).Error
	if err == nil {
		// 已经点赞过，返回错误
		return errors.New("已经点赞过该文章")
	}

	// 创建新的点赞记录
	likeRecord := portal.SysArticleLike{
		ArticleID: uint(articleID),
		IP:        ip,
	}
	return global.DB.Create(&likeRecord).Error
}

// GetArticleViewCount 获取文章浏览次数
func (sysArticleService *SysArticleService) GetArticleViewCount(articleID string) (int64, error) {
	var count int64
	err := global.DB.Model(&portal.SysArticleView{}).Where("article_id = ?", articleID).Count(&count).Error
	return count, err
}

// GetArticleLikeCount 获取文章点赞次数
func (sysArticleService *SysArticleService) GetArticleLikeCount(articleID string) (int64, error) {
	var count int64
	err := global.DB.Model(&portal.SysArticleLike{}).Where("article_id = ?", articleID).Count(&count).Error
	return count, err
}

// GetTotalViewCount 获取总浏览量
func (sysArticleService *SysArticleService) GetTotalViewCount() (int64, error) {
	var count int64
	err := global.DB.Model(&portal.SysArticleView{}).Count(&count).Error
	return count, err
}

// GetTotalLikeCount 获取总点赞数
func (sysArticleService *SysArticleService) GetTotalLikeCount() (int64, error) {
	var count int64
	err := global.DB.Model(&portal.SysArticleLike{}).Count(&count).Error
	return count, err
}

// GetTodayViewCount 获取今天的浏览量
func (sysArticleService *SysArticleService) GetTodayViewCount() (int64, error) {
	today := time.Now().Truncate(24 * time.Hour)
	tomorrow := today.Add(24 * time.Hour)

	var count int64
	err := global.DB.Model(&portal.SysArticleView{}).
		Where("created_at >= ? AND created_at < ?", today, tomorrow).
		Count(&count).Error
	return count, err
}

// GetTodayLikeCount 获取今天的点赞数
func (sysArticleService *SysArticleService) GetTodayLikeCount() (int64, error) {
	today := time.Now().Truncate(24 * time.Hour)
	tomorrow := today.Add(24 * time.Hour)

	var count int64
	err := global.DB.Model(&portal.SysArticleLike{}).
		Where("created_at >= ? AND created_at < ?", today, tomorrow).
		Count(&count).Error
	return count, err
}

// GetYesterdayViewCount 获取昨天的浏览量
func (sysArticleService *SysArticleService) GetYesterdayViewCount() (int64, error) {
	yesterday := time.Now().AddDate(0, 0, -1).Truncate(24 * time.Hour)
	today := yesterday.Add(24 * time.Hour)

	var count int64
	err := global.DB.Model(&portal.SysArticleView{}).
		Where("created_at >= ? AND created_at < ?", yesterday, today).
		Count(&count).Error
	return count, err
}

// GetYesterdayLikeCount 获取昨天的点赞数
func (sysArticleService *SysArticleService) GetYesterdayLikeCount() (int64, error) {
	yesterday := time.Now().AddDate(0, 0, -1).Truncate(24 * time.Hour)
	today := yesterday.Add(24 * time.Hour)

	var count int64
	err := global.DB.Model(&portal.SysArticleLike{}).
		Where("created_at >= ? AND created_at < ?", yesterday, today).
		Count(&count).Error
	return count, err
}
