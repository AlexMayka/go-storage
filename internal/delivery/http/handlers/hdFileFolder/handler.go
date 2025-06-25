package hdFileFolder

import (
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"go-storage/internal/domain"
	"go-storage/pkg/errors"
	"go-storage/pkg/logger"
)

type HandlerFileFolder struct {
	userCase UseCaseFileFolder
}

func NewHandlerFileFolder(useCase UseCaseFileFolder) *HandlerFileFolder {
	return &HandlerFileFolder{
		userCase: useCase,
	}
}

// CreateFolder
// @Summary      Create new folder
// @Description  Creates a new folder in the file system
// @Tags         folders
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        folder  body      RequestCreateFolder  true  "Folder details"
// @Success      201     {object}  ResponseFolder
// @Failure      400,500 {object}  errors.ErrorResponse
// @Failure      401,403 {object}  errors.ErrorResponse
// @Router       /folders [post]
func (h *HandlerFileFolder) CreateFolder(ctx *gin.Context) {
	log := logger.FromContext(ctx)
	companyID := ctx.GetString("company_id")

	if companyID == "" {
		log.Error("func createFolder: Company ID is required", "func", "createFolder", "err", "empty companyId from JWT")
		errors.HandleError(ctx, errors.BadRequest("Company ID is required"))
		return
	}

	var inputData RequestCreateFolder
	if err := ctx.ShouldBind(&inputData); err != nil {
		log.Error("func createFolder: Error in parse input param", "func", "createFolder", "err", err.Error())
		errors.HandleError(ctx, errors.BadRequest("Invalid JSON"))
		return
	}

	domainObj, errValid := ToDomainCreateFolder(&inputData)
	if errValid != nil {
		log.Error("func createFolder: Error in valid param", "func", "createFolder", "err", errValid.Error())
		errors.HandleError(ctx, errors.BadRequest("Invalid JSON"))
		return
	}

	domainObj.CompanyId = companyID

	folder, errUc := h.userCase.CreateFolder(ctx, domainObj)
	if errUc != nil {
		log.Error("func createFolder: Error work UseCase/Repository", "func", "createFolder", "err", errUc.Error())
		errors.HandleError(ctx, errUc)
		return
	}

	ctx.JSON(http.StatusCreated, ToResponseFolder(folder))
}

// GetFolderContents
// @Summary      Get folder contents
// @Description  Returns list of files and folders in specified path
// @Tags         folders
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        request  body      RequestGetFolder  true  "Path and filter options"
// @Success      200      {object}  ResponseGetFolder
// @Failure      400,500  {object}  errors.ErrorResponse
// @Failure      401,403  {object}  errors.ErrorResponse
// @Router       /folders/contents [post]
func (h *HandlerFileFolder) GetFolderContents(ctx *gin.Context) {
	log := logger.FromContext(ctx)
	companyID := ctx.GetString("company_id")

	if companyID == "" {
		log.Error("func getFolderContents: Company ID is required", "func", "getFolderContents", "err", "empty companyId from JWT")
		errors.HandleError(ctx, errors.BadRequest("Company ID is required"))
		return
	}

	var inputData RequestGetFolder
	if err := ctx.ShouldBind(&inputData); err != nil {
		log.Error("func getFolderContents: Error in parse input param", "func", "getFolderContents", "err", err.Error())
		errors.HandleError(ctx, errors.BadRequest("Invalid JSON"))
		return
	}

	path, fileType, err := ToDomainGetFolder(&inputData)
	if err != nil {
		log.Error("func getFolderContents: Error in valid param", "func", "getFolderContents", "err", err.Error())
		errors.HandleError(ctx, errors.BadRequest("Invalid JSON"))
		return
	}

	files, errUc := h.userCase.GetFolderContents(ctx, companyID, path, fileType)
	if errUc != nil {
		log.Error("func getFolderContents: Error work UseCase/Repository", "func", "getFolderContents", "err", errUc.Error())
		errors.HandleError(ctx, errUc)
		return
	}

	ctx.JSON(http.StatusOK, ToResponseFolders(files))
}

// FolderRename
// @Summary      Rename folder
// @Description  Renames a folder at the specified path
// @Tags         folders
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        path     path      string               true  "Folder path"
// @Param        request  body      RequestRenameFolder  true  "New folder name"
// @Success      200      {object}  ResponsePath
// @Failure      400,500  {object}  errors.ErrorResponse
// @Failure      401,403  {object}  errors.ErrorResponse
// @Router       /folders/{path}/rename [put]
func (h *HandlerFileFolder) FolderRename(ctx *gin.Context) {
	log := logger.FromContext(ctx)
	companyID := ctx.GetString("company_id")

	if companyID == "" {
		log.Error("func FolderRename: Company ID is required", "func", "FolderRename", "err", "empty companyId from JWT")
		errors.HandleError(ctx, errors.BadRequest("Company ID is required"))
		return
	}

	var paramPath = ctx.Param("path")
	path, errPath := domain.NewPath(paramPath)

	if errPath != nil {
		log.Error("func FolderRename: Error in parse input param", "func", "FolderRename", "err", errPath.Error())
		errors.HandleError(ctx, errors.BadRequest("Invalid path"))
		return
	}

	var inputData RequestRenameFolder
	if err := ctx.ShouldBind(&inputData); err != nil {
		log.Error("func FolderRename: Error in parse input param", "func", "FolderRename", "err", err.Error())
		errors.HandleError(ctx, errors.BadRequest("Invalid JSON"))
		return
	}

	newName := path.GetParent().Join(inputData.Name)

	newPath, errUc := h.userCase.MoveFolder(ctx, companyID, &path, &newName)
	if errUc != nil {
		log.Error("func FolderRename: Error work UseCase/Repository", "func", "FolderRename", "err", errUc.Error())
		errors.HandleError(ctx, errUc)
		return
	}

	ctx.JSON(http.StatusOK, ToResponsePath(newPath))
}

// MoveFolder
// @Summary      Move folder
// @Description  Moves a folder to a new location
// @Tags         folders
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        path     path      string             true  "Folder path"
// @Param        request  body      RequestMoveFolder  true  "New parent path"
// @Success      200      {object}  ResponsePath
// @Failure      400,500  {object}  errors.ErrorResponse
// @Failure      401,403  {object}  errors.ErrorResponse
// @Router       /folders/{path}/move [put]
func (h *HandlerFileFolder) MoveFolder(ctx *gin.Context) {
	log := logger.FromContext(ctx)
	companyID := ctx.GetString("company_id")

	if companyID == "" {
		log.Error("func MoveFolder: Company ID is required", "func", "MoveFolder", "err", "empty companyId from JWT")
		errors.HandleError(ctx, errors.BadRequest("Company ID is required"))
		return
	}

	var paramPath = ctx.Param("path")
	path, errPath := domain.NewPath(paramPath)

	if errPath != nil {
		log.Error("func MoveFolder: Error in parse input param", "func", "MoveFolder", "err", errPath.Error())
		errors.HandleError(ctx, errors.BadRequest("Invalid path"))
		return
	}

	var inputData RequestMoveFolder
	if err := ctx.ShouldBind(&inputData); err != nil {
		log.Error("func MoveFolder: Error in parse input param", "func", "MoveFolder", "err", err.Error())
		errors.HandleError(ctx, errors.BadRequest("Invalid JSON"))
		return
	}

	newParentPath, errName := domain.NewPath(inputData.ParentPath)
	if errName != nil {
		log.Error("func MoveFolder: Error in parse input newName", "func", "MoveFolder", "err", errName.Error())
		errors.HandleError(ctx, errors.BadRequest("Invalid JSON"))
		return
	}

	newPath := newParentPath.Join(path.GetName())
	_, errUc := h.userCase.MoveFolder(ctx, companyID, &path, &newPath)
	if errUc != nil {
		log.Error("func MoveFolder: Error work UseCase/Repository", "func", "MoveFolder", "err", errUc.Error())
		errors.HandleError(ctx, errUc)
		return
	}

	ctx.JSON(http.StatusOK, ToResponsePath(&newPath))
}

// DeleteFolder
// @Summary      Delete folder
// @Description  Deletes a folder and all its contents
// @Tags         folders
// @Security     BearerAuth
// @Produce      json
// @Param        path  path      string  true  "Folder path"
// @Success      200   {object}  ResponsePath
// @Failure      400,500  {object}  errors.ErrorResponse
// @Failure      401,403  {object}  errors.ErrorResponse
// @Router       /folders/{path} [delete]
func (h *HandlerFileFolder) DeleteFolder(ctx *gin.Context) {
	log := logger.FromContext(ctx)
	companyID := ctx.GetString("company_id")

	if companyID == "" {
		log.Error("func DeleteFolder: Company ID is required", "func", "DeleteFolder", "err", "empty companyId from JWT")
		errors.HandleError(ctx, errors.BadRequest("Company ID is required"))
		return
	}

	var paramPath = ctx.Param("path")
	path, errPath := domain.NewPath(paramPath)
	if errPath != nil {
		log.Error("func DeleteFolder: Error in parse input param", "func", "DeleteFolder", "err", errPath.Error())
		errors.HandleError(ctx, errors.BadRequest("Invalid path"))
		return
	}

	errUc := h.userCase.DeleteFolder(ctx, companyID, &path)
	if errUc != nil {
		log.Error("func DeleteFolder: Error work UseCase/Repository", "func", "DeleteFolder", "err", errUc.Error())
		errors.HandleError(ctx, errUc)
		return
	}

	ctx.JSON(http.StatusOK, ToResponsePath(&path))
}

// UploadFile
// @Summary      Upload file
// @Description  Uploads a file to the specified folder
// @Tags         files
// @Security     BearerAuth
// @Accept       multipart/form-data
// @Produce      json
// @Param        parentPath  formData  string  true  "Parent folder path"
// @Param        file        formData  file    true  "File to upload"
// @Success      201         {object}  ResponseFile
// @Failure      400,500     {object}  errors.ErrorResponse
// @Failure      401,403     {object}  errors.ErrorResponse
// @Router       /files/upload [post]
func (h *HandlerFileFolder) UploadFile(ctx *gin.Context) {
	log := logger.FromContext(ctx)
	companyID := ctx.GetString("company_id")
	userID := ctx.GetString("user_id")

	if companyID == "" {
		log.Error("func UploadFile: Company ID is required", "func", "UploadFile", "err", "empty companyId from JWT")
		errors.HandleError(ctx, errors.BadRequest("Company ID is required"))
		return
	}

	if userID == "" {
		log.Error("func UploadFile: User ID is required", "func", "UploadFile", "err", "empty userId from JWT")
		errors.HandleError(ctx, errors.BadRequest("User ID is required"))
		return
	}

	var inputData RequestUploadFile
	if err := ctx.ShouldBind(&inputData); err != nil {
		log.Error("func UploadFile: Error in parse input param", "func", "UploadFile", "err", err.Error())
		errors.HandleError(ctx, errors.BadRequest("Invalid form data"))
		return
	}

	parentPath, err := domain.NewPath(inputData.ParentPath)
	if err != nil {
		log.Error("func UploadFile: Error in parse parent path", "func", "UploadFile", "err", err.Error())
		errors.HandleError(ctx, errors.BadRequest("Invalid parent path"))
		return
	}

	fileHeader, err := ctx.FormFile("file")
	if err != nil {
		log.Error("func UploadFile: Error getting file from form", "func", "UploadFile", "err", err.Error())
		errors.HandleError(ctx, errors.BadRequest("File is required"))
		return
	}

	file, err := fileHeader.Open()
	if err != nil {
		log.Error("func UploadFile: Error opening file", "func", "UploadFile", "err", err.Error())
		errors.HandleError(ctx, errors.InternalServer("Failed to open file"))
		return
	}
	defer file.Close()

	uploadedFile, errUc := h.userCase.UploadFile(ctx, companyID, userID, &parentPath, fileHeader.Filename, fileHeader.Size, file)
	if errUc != nil {
		log.Error("func UploadFile: Error work UseCase/Repository", "func", "UploadFile", "err", errUc.Error())
		errors.HandleError(ctx, errUc)
		return
	}

	ctx.JSON(http.StatusCreated, ToResponseFile(uploadedFile))
}

// DownloadFile
// @Summary      Download file
// @Description  Downloads a file by ID
// @Tags         files
// @Security     BearerAuth
// @Produce      application/octet-stream
// @Param        id      path   string  true   "File ID"
// @Param        inline  query  bool    false  "Display inline instead of attachment"
// @Success      200     {file}  binary
// @Failure      400,404,500  {object}  errors.ErrorResponse
// @Failure      401,403      {object}  errors.ErrorResponse
// @Router       /files/{id}/download [get]
func (h *HandlerFileFolder) DownloadFile(ctx *gin.Context) {
	log := logger.FromContext(ctx)
	companyID := ctx.GetString("company_id")

	if companyID == "" {
		log.Error("func DownloadFile: Company ID is required", "func", "DownloadFile", "err", "empty companyId from JWT")
		errors.HandleError(ctx, errors.BadRequest("Company ID is required"))
		return
	}

	var inputData RequestDownloadFile
	if err := ctx.ShouldBindUri(&inputData); err != nil {
		log.Error("func DownloadFile: Error in parse URI param", "func", "DownloadFile", "err", err.Error())
		errors.HandleError(ctx, errors.BadRequest("Invalid file ID"))
		return
	}

	if err := ctx.ShouldBindQuery(&inputData); err != nil {
		log.Error("func DownloadFile: Error in parse query param", "func", "DownloadFile", "err", err.Error())
		errors.HandleError(ctx, errors.BadRequest("Invalid query parameters"))
		return
	}

	reader, fileInfo, errUc := h.userCase.DownloadFile(ctx, companyID, inputData.ID)
	if errUc != nil {
		log.Error("func DownloadFile: Error work UseCase/Repository", "func", "DownloadFile", "err", errUc.Error())
		errors.HandleError(ctx, errUc)
		return
	}
	defer reader.Close()

	disposition := "attachment"
	if inputData.Inline {
		disposition = "inline"
	}

	ctx.Header("Content-Disposition", fmt.Sprintf(`%s; filename="%s"`, disposition, fileInfo.Name))
	ctx.Header("Content-Type", *fileInfo.MimeType)
	ctx.Header("Content-Length", strconv.FormatInt(*fileInfo.Size, 10))

	if _, err := io.Copy(ctx.Writer, reader); err != nil {
		log.Error("func DownloadFile: Error streaming file", "func", "DownloadFile", "err", err.Error())
		return
	}
}

// GetFileInfo
// @Summary      Get file information
// @Description  Returns metadata about a file
// @Tags         files
// @Security     BearerAuth
// @Produce      json
// @Param        id  path      string  true  "File ID"
// @Success      200 {object}  ResponseFile
// @Failure      400,404,500  {object}  errors.ErrorResponse
// @Failure      401,403      {object}  errors.ErrorResponse
// @Router       /files/{id} [get]
func (h *HandlerFileFolder) GetFileInfo(ctx *gin.Context) {
	log := logger.FromContext(ctx)
	companyID := ctx.GetString("company_id")

	if companyID == "" {
		log.Error("func GetFileInfo: Company ID is required", "func", "GetFileInfo", "err", "empty companyId from JWT")
		errors.HandleError(ctx, errors.BadRequest("Company ID is required"))
		return
	}

	var inputData RequestGetFileInfo
	if err := ctx.ShouldBindUri(&inputData); err != nil {
		log.Error("func GetFileInfo: Error in parse URI param", "func", "GetFileInfo", "err", err.Error())
		errors.HandleError(ctx, errors.BadRequest("Invalid file ID"))
		return
	}

	fileInfo, errUc := h.userCase.GetFileInfo(ctx, companyID, inputData.ID)
	if errUc != nil {
		log.Error("func GetFileInfo: Error work UseCase/Repository", "func", "GetFileInfo", "err", errUc.Error())
		errors.HandleError(ctx, errUc)
		return
	}

	ctx.JSON(http.StatusOK, ToResponseFile(fileInfo))
}

// RenameFile
// @Summary      Rename file
// @Description  Renames a file
// @Tags         files
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        id       path      string             true  "File ID"
// @Param        request  body      RequestRenameFile  true  "New file name"
// @Success      200      {object}  ResponseFile
// @Failure      400,404,500  {object}  errors.ErrorResponse
// @Failure      401,403      {object}  errors.ErrorResponse
// @Router       /files/{id}/rename [put]
func (h *HandlerFileFolder) RenameFile(ctx *gin.Context) {
	log := logger.FromContext(ctx)
	companyID := ctx.GetString("company_id")

	if companyID == "" {
		log.Error("func RenameFile: Company ID is required", "func", "RenameFile", "err", "empty companyId from JWT")
		errors.HandleError(ctx, errors.BadRequest("Company ID is required"))
		return
	}

	fileID := ctx.Param("id")
	if fileID == "" {
		log.Error("func RenameFile: File ID is required", "func", "RenameFile", "err", "empty file ID")
		errors.HandleError(ctx, errors.BadRequest("File ID is required"))
		return
	}

	var inputData RequestRenameFile
	if err := ctx.ShouldBind(&inputData); err != nil {
		log.Error("func RenameFile: Error in parse input param", "func", "RenameFile", "err", err.Error())
		errors.HandleError(ctx, errors.BadRequest("Invalid JSON"))
		return
	}

	renamedFile, errUc := h.userCase.RenameFile(ctx, companyID, fileID, inputData.Name)
	if errUc != nil {
		log.Error("func RenameFile: Error work UseCase/Repository", "func", "RenameFile", "err", errUc.Error())
		errors.HandleError(ctx, errUc)
		return
	}

	ctx.JSON(http.StatusOK, ToResponseFile(renamedFile))
}

// MoveFile
// @Summary      Move file
// @Description  Moves a file to a new location
// @Tags         files
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        id       path      string           true  "File ID"
// @Param        request  body      RequestMoveFile  true  "New parent path"
// @Success      200      {object}  ResponseFile
// @Failure      400,404,500  {object}  errors.ErrorResponse
// @Failure      401,403      {object}  errors.ErrorResponse
// @Router       /files/{id}/move [put]
func (h *HandlerFileFolder) MoveFile(ctx *gin.Context) {
	log := logger.FromContext(ctx)
	companyID := ctx.GetString("company_id")

	if companyID == "" {
		log.Error("func MoveFile: Company ID is required", "func", "MoveFile", "err", "empty companyId from JWT")
		errors.HandleError(ctx, errors.BadRequest("Company ID is required"))
		return
	}

	fileID := ctx.Param("id")
	if fileID == "" {
		log.Error("func MoveFile: File ID is required", "func", "MoveFile", "err", "empty file ID")
		errors.HandleError(ctx, errors.BadRequest("File ID is required"))
		return
	}

	var inputData RequestMoveFile
	if err := ctx.ShouldBind(&inputData); err != nil {
		log.Error("func MoveFile: Error in parse input param", "func", "MoveFile", "err", err.Error())
		errors.HandleError(ctx, errors.BadRequest("Invalid JSON"))
		return
	}

	newParentPath, err := domain.NewPath(inputData.ParentPath)
	if err != nil {
		log.Error("func MoveFile: Error in parse parent path", "func", "MoveFile", "err", err.Error())
		errors.HandleError(ctx, errors.BadRequest("Invalid parent path"))
		return
	}

	movedFile, errUc := h.userCase.MoveFile(ctx, companyID, fileID, &newParentPath)
	if errUc != nil {
		log.Error("func MoveFile: Error work UseCase/Repository", "func", "MoveFile", "err", errUc.Error())
		errors.HandleError(ctx, errUc)
		return
	}

	ctx.JSON(http.StatusOK, ToResponseFile(movedFile))
}

// DeleteFile
// @Summary      Delete file
// @Description  Deletes a file from the system
// @Tags         files
// @Security     BearerAuth
// @Produce      json
// @Param        id  path      string  true  "File ID"
// @Success      200 {object}  ResponseSuccess
// @Failure      400,404,500  {object}  errors.ErrorResponse
// @Failure      401,403      {object}  errors.ErrorResponse
// @Router       /files/{id} [delete]
func (h *HandlerFileFolder) DeleteFile(ctx *gin.Context) {
	log := logger.FromContext(ctx)
	companyID := ctx.GetString("company_id")

	if companyID == "" {
		log.Error("func DeleteFile: Company ID is required", "func", "DeleteFile", "err", "empty companyId from JWT")
		errors.HandleError(ctx, errors.BadRequest("Company ID is required"))
		return
	}

	fileID := ctx.Param("id")
	if fileID == "" {
		log.Error("func DeleteFile: File ID is required", "func", "DeleteFile", "err", "empty file ID")
		errors.HandleError(ctx, errors.BadRequest("File ID is required"))
		return
	}

	errUc := h.userCase.DeleteFile(ctx, companyID, fileID)
	if errUc != nil {
		log.Error("func DeleteFile: Error work UseCase/Repository", "func", "DeleteFile", "err", errUc.Error())
		errors.HandleError(ctx, errUc)
		return
	}

	ctx.JSON(http.StatusOK, ToResponseSuccess("File deleted successfully"))
}

// GetUploadStrategy
// @Summary      Get upload strategy
// @Description  Returns recommended upload strategy based on file size
// @Tags         files
// @Security     BearerAuth
// @Produce      json
// @Param        fileSize  query     int  true  "File size in bytes"
// @Success      200       {object}  ResponseUploadStrategy
// @Failure      400,500   {object}  errors.ErrorResponse
// @Failure      401,403   {object}  errors.ErrorResponse
// @Router       /files/upload-strategy [get]
func (h *HandlerFileFolder) GetUploadStrategy(ctx *gin.Context) {
	log := logger.FromContext(ctx)

	fileSizeStr := ctx.Query("fileSize")
	if fileSizeStr == "" {
		log.Error("func GetUploadStrategy: File size is required", "func", "GetUploadStrategy", "err", "missing fileSize query parameter")
		errors.HandleError(ctx, errors.BadRequest("File size is required"))
		return
	}

	fileSize, err := strconv.ParseInt(fileSizeStr, 10, 64)
	if err != nil {
		log.Error("func GetUploadStrategy: Invalid file size", "func", "GetUploadStrategy", "err", err.Error())
		errors.HandleError(ctx, errors.BadRequest("Invalid file size"))
		return
	}

	strategy, errUc := h.userCase.GetUploadStrategy(ctx, fileSize)
	if errUc != nil {
		log.Error("func GetUploadStrategy: Error work UseCase", "func", "GetUploadStrategy", "err", errUc.Error())
		errors.HandleError(ctx, errUc)
		return
	}

	ctx.JSON(http.StatusOK, ToResponseUploadStrategy(strategy))
}

// InitChunkedUpload
// @Summary      Initialize chunked upload
// @Description  Initializes a chunked upload session for large files
// @Tags         chunked-upload
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        request  body      RequestInitChunkedUpload  true  "Upload initialization details"
// @Success      201      {object}  ResponseInitChunkedUpload
// @Failure      400,500  {object}  errors.ErrorResponse
// @Failure      401,403  {object}  errors.ErrorResponse
// @Router       /files/chunked/init [post]
func (h *HandlerFileFolder) InitChunkedUpload(ctx *gin.Context) {
	log := logger.FromContext(ctx)
	companyID := ctx.GetString("company_id")
	userID := ctx.GetString("user_id")

	if companyID == "" {
		log.Error("func InitChunkedUpload: Company ID is required", "func", "InitChunkedUpload", "err", "empty companyId from JWT")
		errors.HandleError(ctx, errors.BadRequest("Company ID is required"))
		return
	}

	if userID == "" {
		log.Error("func InitChunkedUpload: User ID is required", "func", "InitChunkedUpload", "err", "empty userId from JWT")
		errors.HandleError(ctx, errors.BadRequest("User ID is required"))
		return
	}

	var inputData RequestInitChunkedUpload
	if err := ctx.ShouldBind(&inputData); err != nil {
		log.Error("func InitChunkedUpload: Error in parse input param", "func", "InitChunkedUpload", "err", err.Error())
		errors.HandleError(ctx, errors.BadRequest("Invalid JSON"))
		return
	}

	parentPath, err := domain.NewPath(inputData.ParentPath)
	if err != nil {
		log.Error("func InitChunkedUpload: Error in parse parent path", "func", "InitChunkedUpload", "err", err.Error())
		errors.HandleError(ctx, errors.BadRequest("Invalid parent path"))
		return
	}

	chunkedUpload, errUc := h.userCase.InitChunkedUpload(ctx, companyID, userID, inputData.FileName, inputData.FileSize, &parentPath, inputData.MimeType)
	if errUc != nil {
		log.Error("func InitChunkedUpload: Error work UseCase/Repository", "func", "InitChunkedUpload", "err", errUc.Error())
		errors.HandleError(ctx, errUc)
		return
	}

	strategy, _ := h.userCase.GetUploadStrategy(ctx, inputData.FileSize)

	ctx.JSON(http.StatusCreated, ToResponseInitChunkedUpload(chunkedUpload, strategy))
}

// UploadChunk
// @Summary      Upload file chunk
// @Description  Uploads a single chunk of a file in a chunked upload session
// @Tags         chunked-upload
// @Security     BearerAuth
// @Accept       multipart/form-data
// @Produce      json
// @Param        uploadId    path      string  true  "Upload session ID"
// @Param        chunkIndex  path      string  true  "Chunk index (0-based)"
// @Param        chunk       formData  file    true  "Chunk data"
// @Success      200         {object}  ResponseUploadChunk
// @Failure      400,500     {object}  errors.ErrorResponse
// @Failure      401,403     {object}  errors.ErrorResponse
// @Router       /files/chunked/{uploadId}/chunk/{chunkIndex} [post]
func (h *HandlerFileFolder) UploadChunk(ctx *gin.Context) {
	log := logger.FromContext(ctx)
	companyID := ctx.GetString("company_id")

	if companyID == "" {
		log.Error("func UploadChunk: Company ID is required", "func", "UploadChunk", "err", "empty companyId from JWT")
		errors.HandleError(ctx, errors.BadRequest("Company ID is required"))
		return
	}

	var inputData RequestUploadChunk
	if err := ctx.ShouldBindUri(&inputData); err != nil {
		log.Error("func UploadChunk: Error in parse URI param", "func", "UploadChunk", "err", err.Error())
		errors.HandleError(ctx, errors.BadRequest("Invalid URI parameters"))
		return
	}

	fileHeader, err := ctx.FormFile("chunk")
	if err != nil {
		log.Error("func UploadChunk: Error getting chunk from form", "func", "UploadChunk", "err", err.Error())
		errors.HandleError(ctx, errors.BadRequest("Chunk data is required"))
		return
	}

	file, err := fileHeader.Open()
	if err != nil {
		log.Error("func UploadChunk: Error opening chunk file", "func", "UploadChunk", "err", err.Error())
		errors.HandleError(ctx, errors.InternalServer("Failed to open chunk"))
		return
	}
	defer file.Close()

	chunkIndex, err := strconv.Atoi(inputData.ChunkIndex)
	if err != nil {
		log.Error("func UploadChunk: Invalid chunk index", "func", "UploadChunk", "err", err.Error())
		errors.HandleError(ctx, errors.BadRequest("Invalid chunk index"))
		return
	}

	chunkedUpload, errUc := h.userCase.UploadChunk(ctx, companyID, inputData.UploadID, chunkIndex, file, fileHeader.Size)
	if errUc != nil {
		log.Error("func UploadChunk: Error work UseCase/Repository", "func", "UploadChunk", "err", errUc.Error())
		errors.HandleError(ctx, errUc)
		return
	}

	ctx.JSON(http.StatusOK, ToResponseUploadChunk(chunkedUpload, chunkIndex))
}

// GetChunkedUploadStatus
// @Summary      Get chunked upload status
// @Description  Returns the current status of a chunked upload session
// @Tags         chunked-upload
// @Security     BearerAuth
// @Produce      json
// @Param        uploadId  path      string  true  "Upload session ID"
// @Success      200       {object}  ResponseChunkedUploadStatus
// @Failure      400,404,500  {object}  errors.ErrorResponse
// @Failure      401,403      {object}  errors.ErrorResponse
// @Router       /files/chunked/{uploadId}/status [get]
func (h *HandlerFileFolder) GetChunkedUploadStatus(ctx *gin.Context) {
	log := logger.FromContext(ctx)
	companyID := ctx.GetString("company_id")

	if companyID == "" {
		log.Error("func GetChunkedUploadStatus: Company ID is required", "func", "GetChunkedUploadStatus", "err", "empty companyId from JWT")
		errors.HandleError(ctx, errors.BadRequest("Company ID is required"))
		return
	}

	var inputData RequestGetChunkedUploadStatus
	if err := ctx.ShouldBindUri(&inputData); err != nil {
		log.Error("func GetChunkedUploadStatus: Error in parse URI param", "func", "GetChunkedUploadStatus", "err", err.Error())
		errors.HandleError(ctx, errors.BadRequest("Invalid upload ID"))
		return
	}

	chunkedUpload, errUc := h.userCase.GetChunkedUploadStatus(ctx, companyID, inputData.UploadID)
	if errUc != nil {
		log.Error("func GetChunkedUploadStatus: Error work UseCase/Repository", "func", "GetChunkedUploadStatus", "err", errUc.Error())
		errors.HandleError(ctx, errUc)
		return
	}

	ctx.JSON(http.StatusOK, ToResponseChunkedUploadStatus(chunkedUpload))
}

// CompleteChunkedUpload
// @Summary      Complete chunked upload
// @Description  Completes a chunked upload session and creates the final file
// @Tags         chunked-upload
// @Security     BearerAuth
// @Produce      json
// @Param        uploadId  path      string  true  "Upload session ID"
// @Success      200       {object}  ResponseCompleteChunkedUpload
// @Failure      400,404,500  {object}  errors.ErrorResponse
// @Failure      401,403      {object}  errors.ErrorResponse
// @Router       /files/chunked/{uploadId}/complete [post]
func (h *HandlerFileFolder) CompleteChunkedUpload(ctx *gin.Context) {
	log := logger.FromContext(ctx)
	companyID := ctx.GetString("company_id")

	if companyID == "" {
		log.Error("func CompleteChunkedUpload: Company ID is required", "func", "CompleteChunkedUpload", "err", "empty companyId from JWT")
		errors.HandleError(ctx, errors.BadRequest("Company ID is required"))
		return
	}

	var inputData RequestCompleteChunkedUpload
	if err := ctx.ShouldBindUri(&inputData); err != nil {
		log.Error("func CompleteChunkedUpload: Error in parse URI param", "func", "CompleteChunkedUpload", "err", err.Error())
		errors.HandleError(ctx, errors.BadRequest("Invalid upload ID"))
		return
	}

	completedFile, errUc := h.userCase.CompleteChunkedUpload(ctx, companyID, inputData.UploadID)
	if errUc != nil {
		log.Error("func CompleteChunkedUpload: Error work UseCase/Repository", "func", "CompleteChunkedUpload", "err", errUc.Error())
		errors.HandleError(ctx, errUc)
		return
	}

	ctx.JSON(http.StatusOK, ToResponseCompleteChunkedUpload(completedFile))
}

// AbortChunkedUpload
// @Summary      Abort chunked upload
// @Description  Aborts a chunked upload session and cleans up temporary data
// @Tags         chunked-upload
// @Security     BearerAuth
// @Produce      json
// @Param        uploadId  path      string  true  "Upload session ID"
// @Success      200       {object}  ResponseSuccess
// @Failure      400,404,500  {object}  errors.ErrorResponse
// @Failure      401,403      {object}  errors.ErrorResponse
// @Router       /files/chunked/{uploadId}/abort [delete]
func (h *HandlerFileFolder) AbortChunkedUpload(ctx *gin.Context) {
	log := logger.FromContext(ctx)
	companyID := ctx.GetString("company_id")

	if companyID == "" {
		log.Error("func AbortChunkedUpload: Company ID is required", "func", "AbortChunkedUpload", "err", "empty companyId from JWT")
		errors.HandleError(ctx, errors.BadRequest("Company ID is required"))
		return
	}

	uploadID := ctx.Param("uploadId")
	if uploadID == "" {
		log.Error("func AbortChunkedUpload: Upload ID is required", "func", "AbortChunkedUpload", "err", "empty upload ID")
		errors.HandleError(ctx, errors.BadRequest("Upload ID is required"))
		return
	}

	errUc := h.userCase.AbortChunkedUpload(ctx, companyID, uploadID)
	if errUc != nil {
		log.Error("func AbortChunkedUpload: Error work UseCase/Repository", "func", "AbortChunkedUpload", "err", errUc.Error())
		errors.HandleError(ctx, errUc)
		return
	}

	ctx.JSON(http.StatusOK, ToResponseSuccess("Chunked upload aborted successfully"))
}

// GetResourceStats
// @Summary      Get resource statistics
// @Description  Returns current system resource usage and limits
// @Tags         monitoring
// @Security     BearerAuth
// @Produce      json
// @Success      200     {object}  ResponseResourceStats
// @Failure      500     {object}  errors.ErrorResponse
// @Failure      401,403 {object}  errors.ErrorResponse
// @Router       /files/stats [get]
func (h *HandlerFileFolder) GetResourceStats(ctx *gin.Context) {
	log := logger.FromContext(ctx)

	stats, errUc := h.userCase.GetResourceStats(ctx)
	if errUc != nil {
		log.Error("func GetResourceStats: Error work UseCase", "func", "GetResourceStats", "err", errUc.Error())
		errors.HandleError(ctx, errUc)
		return
	}

	ctx.JSON(http.StatusOK, ToResponseResourceStats(stats))
}
