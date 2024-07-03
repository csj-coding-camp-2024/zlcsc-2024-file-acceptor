package main

import (
	"fmt"
	"mime/multipart"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	r.POST("/", func(c *gin.Context) {
		// Set memory limit
		r.MaxMultipartMemory = 8 << 20 // 8 MiB

		// Bind Request
		type ModpackUploadRequest struct {
			File *multipart.FileHeader `binding:"required" form:"file"`
		}
		var req ModpackUploadRequest
		if err := c.ShouldBind(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Save uploaded file
		err := c.SaveUploadedFile(req.File, "uploaded/"+req.File.Filename)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err})
			return
		}

		c.JSON(http.StatusOK, gin.H{"name": req.File.Filename})
	})

	r.GET("/:filename", func(c *gin.Context) {
		// Get id from url
		filename := c.Param("filename")

		fmt.Println(filename)

		c.Header("Content-Disposition", "attachment; filename="+filename)
		http.ServeFile(c.Writer, c.Request, "uploaded/"+filename)
	})
	r.Run(":80")
}
