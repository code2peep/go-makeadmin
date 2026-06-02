package adapter

import (
	"context"
	"errors"
	"mime/multipart"
	"path"

	"gorm.io/gorm"

	"go-makeadmin/admin/schemas/req"
	"go-makeadmin/admin/schemas/resp"
	"go-makeadmin/config"
	"go-makeadmin/core"
	"go-makeadmin/core/request"
	"go-makeadmin/core/response"
	"go-makeadmin/makeadmin/repository"
	makeadminsvc "go-makeadmin/makeadmin/service"
	"go-makeadmin/model/makeadmin"
	"go-makeadmin/plugin"
	"go-makeadmin/util"
)

type FileAdapter interface {
	Available(ctx context.Context) bool
	AlbumList(ctx context.Context, page request.PageReq, listReq req.CommonAlbumListReq) (response.PageResp, error)
	AlbumRename(ctx context.Context, id uint, name string) error
	AlbumMove(ctx context.Context, ids []uint, cid int) error
	AlbumAdd(ctx context.Context, addReq req.CommonAlbumAddReq) (uint, error)
	AlbumDel(ctx context.Context, ids []uint) error
	CateList(ctx context.Context, listReq req.CommonCateListReq) ([]interface{}, error)
	CateAdd(ctx context.Context, addReq req.CommonCateAddReq) error
	CateRename(ctx context.Context, id uint, name string) error
	CateDel(ctx context.Context, id uint) error
	UploadImage(ctx context.Context, file *multipart.FileHeader, cid uint, adminID uint) (resp.CommonUploadFileResp, error)
	UploadVideo(ctx context.Context, file *multipart.FileHeader, cid uint, adminID uint) (resp.CommonUploadFileResp, error)
}

type fileAdapter struct {
	db *gorm.DB
}

func NewFileAdapter(db *gorm.DB) FileAdapter {
	return fileAdapter{db: db}
}

func (adapter fileAdapter) Available(ctx context.Context) bool {
	if adapter.db == nil ||
		!adapter.db.Migrator().HasTable(&makeadmin.FileCategory{}) ||
		!adapter.db.Migrator().HasTable(&makeadmin.File{}) {
		return false
	}
	var count int64
	err := adapter.db.WithContext(ctx).
		Model(&makeadmin.FileCategory{}).
		Where("tenant_id = ? AND delete_time = ?", tenantIDFromContext(ctx), 0).
		Count(&count).
		Error
	return err == nil && count > 0
}

func (adapter fileAdapter) AlbumList(ctx context.Context, page request.PageReq, listReq req.CommonAlbumListReq) (response.PageResp, error) {
	page = normalizePage(page)
	fileType, err := fileTypeFromLegacy(listReq.Type)
	if err != nil {
		return response.PageResp{}, mapFileError(err)
	}
	categoryID := uint64(0)
	if listReq.Cid > 0 {
		categoryID = uint64(listReq.Cid)
	}
	tenantID := tenantIDFromContext(ctx)
	result, err := adapter.fileService().ListFiles(ctx, tenantID, makeadminsvc.FileFilter{
		CategoryID: categoryID,
		FileType:   fileType,
		Name:       listReq.Name,
	}, makeadminsvc.FilePageInput{PageNo: page.PageNo, PageSize: page.PageSize})
	if err != nil {
		return response.PageResp{}, mapFileError(err)
	}
	return response.PageResp{
		PageNo:   page.PageNo,
		PageSize: page.PageSize,
		Count:    result.Count,
		Lists:    albumListResponses(result.Items),
	}, nil
}

func (adapter fileAdapter) AlbumRename(ctx context.Context, id uint, name string) error {
	return mapFileError(adapter.fileService().RenameFile(ctx, tenantIDFromContext(ctx), uint64(id), name))
}

func (adapter fileAdapter) AlbumMove(ctx context.Context, ids []uint, cid int) error {
	categoryID := uint64(0)
	if cid > 0 {
		categoryID = uint64(cid)
	}
	return mapFileError(adapter.fileService().MoveFiles(ctx, tenantIDFromContext(ctx), uintSliceToUint64(ids), categoryID))
}

func (adapter fileAdapter) AlbumAdd(ctx context.Context, addReq req.CommonAlbumAddReq) (uint, error) {
	fileType, err := fileTypeFromLegacy(addReq.Type)
	if err != nil {
		return 0, mapFileError(err)
	}
	tenantID := tenantIDFromContext(ctx)
	item, err := adapter.fileService().AddFile(ctx, makeadminsvc.FileInput{
		TenantID:   tenantID,
		CategoryID: uint64(addReq.Cid),
		AdminID:    uint64(addReq.Aid),
		FileType:   fileType,
		Name:       addReq.Name,
		URI:        addReq.Uri,
		URL:        util.UrlUtil.ToAbsoluteUrl(addReq.Uri),
		Ext:        addReq.Ext,
		Size:       addReq.Size,
	})
	if err != nil {
		return 0, mapFileError(err)
	}
	return uint(item.ID), nil
}

func (adapter fileAdapter) AlbumDel(ctx context.Context, ids []uint) error {
	return mapFileError(adapter.fileService().DeleteFiles(ctx, tenantIDFromContext(ctx), uintSliceToUint64(ids)))
}

func (adapter fileAdapter) CateList(ctx context.Context, listReq req.CommonCateListReq) ([]interface{}, error) {
	fileType, err := fileTypeFromLegacy(listReq.Type)
	if err != nil {
		return nil, mapFileError(err)
	}
	categories, err := adapter.fileService().ListCategories(ctx, tenantIDFromContext(ctx), fileType, listReq.Name)
	if err != nil {
		return nil, mapFileError(err)
	}
	return util.ArrayUtil.ListToTree(
		util.ConvertUtil.StructsToMaps(cateListResponses(categories)), "id", "pid", "children"), nil
}

func (adapter fileAdapter) CateAdd(ctx context.Context, addReq req.CommonCateAddReq) error {
	fileType, err := fileTypeFromLegacy(addReq.Type)
	if err != nil {
		return mapFileError(err)
	}
	tenantID := tenantIDFromContext(ctx)
	return mapFileError(adapter.fileService().AddCategory(ctx, makeadminsvc.FileCategoryInput{
		TenantID: tenantID,
		ParentID: uint64(addReq.Pid),
		FileType: fileType,
		Name:     addReq.Name,
	}))
}

func (adapter fileAdapter) CateRename(ctx context.Context, id uint, name string) error {
	return mapFileError(adapter.fileService().RenameCategory(ctx, tenantIDFromContext(ctx), uint64(id), name))
}

func (adapter fileAdapter) CateDel(ctx context.Context, id uint) error {
	return mapFileError(adapter.fileService().DeleteCategory(ctx, tenantIDFromContext(ctx), uint64(id)))
}

func (adapter fileAdapter) UploadImage(ctx context.Context, file *multipart.FileHeader, cid uint, adminID uint) (resp.CommonUploadFileResp, error) {
	return adapter.uploadFile(ctx, file, "image", 10, cid, adminID)
}

func (adapter fileAdapter) UploadVideo(ctx context.Context, file *multipart.FileHeader, cid uint, adminID uint) (resp.CommonUploadFileResp, error) {
	return adapter.uploadFile(ctx, file, "video", 20, cid, adminID)
}

func (adapter fileAdapter) uploadFile(ctx context.Context, file *multipart.FileHeader, folder string, legacyType int, cid uint, adminID uint) (resp.CommonUploadFileResp, error) {
	upRes, err := plugin.StorageDriver.Upload(file, folder, legacyType)
	if err != nil {
		return resp.CommonUploadFileResp{}, err
	}
	addReq := req.CommonAlbumAddReq{
		Cid:  cid,
		Aid:  adminID,
		Type: legacyType,
		Name: upRes.Name,
		Uri:  upRes.Uri,
		Ext:  upRes.Ext,
		Size: upRes.Size,
	}
	albumID, err := adapter.AlbumAdd(ctx, addReq)
	if err != nil {
		return resp.CommonUploadFileResp{}, err
	}
	return resp.CommonUploadFileResp{
		ID:   albumID,
		Cid:  cid,
		Aid:  adminID,
		Type: legacyType,
		Name: upRes.Name,
		Uri:  upRes.Uri,
		Path: upRes.Path,
		Ext:  upRes.Ext,
		Size: upRes.Size,
	}, nil
}

func (adapter fileAdapter) fileService() makeadminsvc.FileService {
	return makeadminsvc.NewFileService(repository.NewFileRepository(adapter.db))
}

func albumListResponses(items []makeadminsvc.FileItem) []resp.CommonAlbumListResp {
	result := make([]resp.CommonAlbumListResp, 0, len(items))
	for _, item := range items {
		result = append(result, resp.CommonAlbumListResp{
			ID:         uint(item.ID),
			Cid:        uint(item.CategoryID),
			Name:       item.Name,
			Path:       path.Join(config.Config.PublicPrefix, item.URI),
			Uri:        util.UrlUtil.ToAbsoluteUrl(item.URI),
			Ext:        item.Ext,
			Size:       util.ServerUtil.GetFmtSize(uint64(item.Size)),
			CreateTime: core.TsTime(item.CreateTime),
			UpdateTime: core.TsTime(item.UpdateTime),
		})
	}
	return result
}

func cateListResponses(items []makeadminsvc.FileCategoryItem) []resp.CommonCateListResp {
	result := make([]resp.CommonCateListResp, 0, len(items))
	for _, item := range items {
		result = append(result, resp.CommonCateListResp{
			ID:         uint(item.ID),
			Pid:        uint(item.ParentID),
			Name:       item.Name,
			CreateTime: core.TsTime(item.CreateTime),
			UpdateTime: core.TsTime(item.UpdateTime),
		})
	}
	return result
}

func mapFileError(err error) error {
	switch {
	case err == nil:
		return nil
	case errors.Is(err, makeadminsvc.ErrFileNotFound):
		return response.AssertArgumentError.Make("文件丢失！")
	case errors.Is(err, makeadminsvc.ErrFileCategoryNotFound):
		return response.AssertArgumentError.Make("分类已不存在！")
	case errors.Is(err, makeadminsvc.ErrFileCategoryInUse):
		return response.AssertArgumentError.Make("当前分类正被使用中,不能删除！")
	case errors.Is(err, makeadminsvc.ErrFileCategoryHasChildren):
		return response.AssertArgumentError.Make("当前分类存在子分类,不能删除！")
	case errors.Is(err, makeadminsvc.ErrUnsupportedFileType):
		return response.AssertArgumentError.Make("文件类型不支持！")
	default:
		return err
	}
}

func fileTypeFromLegacy(fileType int) (string, error) {
	switch fileType {
	case 0:
		return "", nil
	case 10:
		return makeadmin.FileTypeImage, nil
	case 20:
		return makeadmin.FileTypeVideo, nil
	case 30:
		return makeadminsvc.FileTypeOther, nil
	default:
		return "", makeadminsvc.ErrUnsupportedFileType
	}
}
