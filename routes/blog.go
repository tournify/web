package routes

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func (controller Controller) BlogView(c *gin.Context) {
	c.HTML(http.StatusOK, "blog.html", controller.defaultPageData(c))
}

func (controller Controller) BlogCreate(c *gin.Context) {
	c.HTML(http.StatusOK, "blog.html", controller.defaultPageData(c))
}

func (controller Controller) BlogCreatePost(c *gin.Context) {
	c.HTML(http.StatusOK, "blog.html", controller.defaultPageData(c))
}

func (controller Controller) BlogViewPage(c *gin.Context) {
	c.HTML(http.StatusOK, "blog.html", controller.defaultPageData(c))
}
