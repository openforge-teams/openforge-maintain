package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/openforge-maintain/openforge-maintain/agent/service/file"
	"github.com/openforge-maintain/openforge-maintain/pkg/response"
)

// FileHandler handles file management API requests.
type FileHandler struct {
	fileService *file.FileManagerService
}

// NewFileHandler creates a new FileHandler.
func NewFileHandler(fileService *file.FileManagerService) *FileHandler {
	return &FileHandler{fileService: fileService}
}

// ListDir lists the contents of a directory.
func (h *FileHandler) ListDir(c *gin.Context) {
	path := c.DefaultQuery("path", "/")
	items, err := h.fileService.ListDir(path)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, items)
}

// GetFileContent reads and returns the content of a file.
func (h *FileHandler) GetFileContent(c *gin.Context) {
	path := c.Query("path")
	if path == "" {
		response.BadRequest(c, "path is required")
		return
	}

	content, err := h.fileService.GetFileContent(path)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, gin.H{"content": content})
}

// SaveFileContent writes content to a file.
func (h *FileHandler) SaveFileContent(c *gin.Context) {
	var req struct {
		Path    string `json:"path" binding:"required"`
		Content string `json:"content"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if err := h.fileService.SaveFileContent(req.Path, req.Content); err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, nil)
}

// Upload handles file upload.
func (h *FileHandler) Upload(c *gin.Context) {
	path := c.PostForm("path")
	if path == "" {
		response.BadRequest(c, "path is required")
		return
	}

	overwrite := c.PostForm("overwrite") == "true"
	file, err := c.FormFile("file")
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	src, err := file.Open()
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	defer src.Close()

	if err := h.fileService.Upload(src, path, overwrite); err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, nil)
}

// Download handles file download.
func (h *FileHandler) Download(c *gin.Context) {
	path := c.Query("path")
	if path == "" {
		response.BadRequest(c, "path is required")
		return
	}

	info, err := h.fileService.GetInfo(path)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}

	reader, err := h.fileService.Download(path)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	defer reader.Close()

	c.DataFromReader(http.StatusOK, info.Size, "application/octet-stream", reader, nil)
}

// DeleteFile deletes a file or directory.
func (h *FileHandler) DeleteFile(c *gin.Context) {
	var req struct {
		Path string `json:"path" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if err := h.fileService.Delete(req.Path); err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, nil)
}

// RenameFile renames or moves a file.
func (h *FileHandler) RenameFile(c *gin.Context) {
	var req struct {
		OldPath string `json:"old_path" binding:"required"`
		NewPath string `json:"new_path" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if err := h.fileService.Rename(req.OldPath, req.NewPath); err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, nil)
}

// ChmodFile changes file permissions.
func (h *FileHandler) ChmodFile(c *gin.Context) {
	var req struct {
		Path string `json:"path" binding:"required"`
		Mode uint32 `json:"mode" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if err := h.fileService.Chmod(req.Path, req.Mode); err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, nil)
}

// ChownFile changes file ownership.
func (h *FileHandler) ChownFile(c *gin.Context) {
	var req struct {
		Path string `json:"path" binding:"required"`
		UID  int    `json:"uid" binding:"required"`
		GID  int    `json:"gid" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if err := h.fileService.Chown(req.Path, req.UID, req.GID); err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, nil)
}

// CompressFiles creates an archive.
func (h *FileHandler) CompressFiles(c *gin.Context) {
	var req struct {
		Paths  []string `json:"paths" binding:"required"`
		Dest   string   `json:"dest" binding:"required"`
		Format string   `json:"format"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if req.Format == "" {
		req.Format = "tar.gz"
	}

	if err := h.fileService.Compress(req.Paths, req.Dest, req.Format); err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, nil)
}

// ExtractFiles extracts an archive.
func (h *FileHandler) ExtractFiles(c *gin.Context) {
	var req struct {
		Src  string `json:"src" binding:"required"`
		Dest string `json:"dest" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if err := h.fileService.Extract(req.Src, req.Dest); err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, nil)
}

// Mkdir creates a directory.
func (h *FileHandler) Mkdir(c *gin.Context) {
	var req struct {
		Path string `json:"path" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if err := h.fileService.Mkdir(req.Path); err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, nil)
}

// GetFileInfo returns detailed information about a file.
func (h *FileHandler) GetFileInfo(c *gin.Context) {
	path := c.Query("path")
	if path == "" {
		response.BadRequest(c, "path is required")
		return
	}

	info, err := h.fileService.GetInfo(path)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, info)
}
