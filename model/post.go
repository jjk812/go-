package model

import (
	"github.com/jinzhu/gorm"
)

type Post struct {
	gorm.Model
	Title     string `json:"title"`
	Content   string `json:"content"`
	Username  string `json:"username"`
	ViewCount int    `json:"viewcount"`
}
type Pagination struct {
	TotalItems   int
	ItemsPerPage int
	CurrentPage  int
	TotalPages   int
	HasPrevious  bool
	HasNext      bool
}

//        c.HTML(http.StatusOK, "index.html", gin.H{
//            "posts":      currentPagePosts,
//            "pagination": gin.H{
//                "TotalItems":   totalItems,
//                "ItemsPerPage": itemsPerPage,
//                "CurrentPage":  page,
//                "TotalPages":   totalPages,
//                "HasPrevious":  page > 1,
//                "HasNext":      page < totalPages,
//            },
//        })
