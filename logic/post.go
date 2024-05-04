package logic

import (
	"goPro/dao"
	"goPro/model"
)

func Pages(page int) ([]model.Post, model.Pagination) {
	var pagination model.Pagination
	posts, _ := dao.SearchPost("")
	// 每页显示的文章数量
	itemsPerPage := 5
	// 使用dao搜索相关文章
	// 总文章数量
	pagination.TotalItems = len(posts)
	// 计算总页数
	pagination.TotalPages = (pagination.TotalItems + itemsPerPage - 1) / itemsPerPage
	// 计算当前页的起始和结束位置
	startIndex := (page - 1) * itemsPerPage
	endIndex := page * itemsPerPage
	if endIndex > pagination.TotalItems {
		endIndex = pagination.TotalItems
	}
	// 获取当前页的文章
	currentPagePosts := posts[startIndex:endIndex]
	pagination.CurrentPage = page
	pagination.HasNext = page < pagination.TotalPages
	pagination.HasPrevious = page > 1
	return currentPagePosts, pagination
}
