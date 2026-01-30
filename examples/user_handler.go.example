// Package handler HTTP请求处理器
package handler

import (
	"strconv"
	"sync"

	"github.com/gin-gonic/gin"
	"ruleback/internal/model"
	"ruleback/internal/service"
	"ruleback/pkg/errors"
	"ruleback/pkg/logger"
	"ruleback/pkg/response"
)

var (
	userHandlerInstance *UserHandler
	userHandlerOnce     sync.Once
)

// UserHandler 用户HTTP处理器
type UserHandler struct {
	service *service.UserService
}

// NewUserHandler 创建UserHandler实例（用于依赖注入）
func NewUserHandler(svc *service.UserService) *UserHandler {
	return &UserHandler{service: svc}
}

// GetUserHandler 获取用户Handler单例（保留向后兼容）
func GetUserHandler() *UserHandler {
	userHandlerOnce.Do(func() {
		userHandlerInstance = &UserHandler{
			service: service.GetUserService(),
		}
	})
	return userHandlerInstance
}

// Create 创建用户
func (h *UserHandler) Create(c *gin.Context) {
	var req model.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Warn("创建用户参数错误", logger.Err(err), logger.String("ip", c.ClientIP()))
		response.Fail(c, errors.CodeInvalidParams, "参数错误: "+err.Error())
		return
	}

	user, err := h.service.Create(&req)
	if err != nil {
		h.handleError(c, err)
		return
	}

	response.SuccessWithData(c, user)
}

// GetByID 根据ID获取用户
func (h *UserHandler) GetByID(c *gin.Context) {
	id, err := h.parseIDParam(c)
	if err != nil {
		response.Fail(c, errors.CodeInvalidParams, "无效的用户ID")
		return
	}

	user, err := h.service.GetByID(id)
	if err != nil {
		h.handleError(c, err)
		return
	}

	response.SuccessWithData(c, user)
}

// List 获取用户列表
func (h *UserHandler) List(c *gin.Context) {
	var query model.UserListQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		logger.Warn("获取用户列表参数错误", logger.Err(err))
		response.Fail(c, errors.CodeInvalidParams, "参数错误: "+err.Error())
		return
	}

	users, total, err := h.service.List(&query)
	if err != nil {
		h.handleError(c, err)
		return
	}

	response.SuccessWithPage(c, users, total, query.Page, query.PageSize)
}

// Update 更新用户
func (h *UserHandler) Update(c *gin.Context) {
	id, err := h.parseIDParam(c)
	if err != nil {
		response.Fail(c, errors.CodeInvalidParams, "无效的用户ID")
		return
	}

	var req model.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Warn("更新用户参数错误", logger.Err(err), logger.Uint("user_id", id))
		response.Fail(c, errors.CodeInvalidParams, "参数错误: "+err.Error())
		return
	}

	user, err := h.service.Update(id, &req)
	if err != nil {
		h.handleError(c, err)
		return
	}

	response.SuccessWithData(c, user)
}

// Delete 删除用户
func (h *UserHandler) Delete(c *gin.Context) {
	id, err := h.parseIDParam(c)
	if err != nil {
		response.Fail(c, errors.CodeInvalidParams, "无效的用户ID")
		return
	}

	if err := h.service.Delete(id); err != nil {
		h.handleError(c, err)
		return
	}

	response.SuccessWithMessage(c, "删除成功")
}

// parseIDParam 解析路径中的ID参数
func (h *UserHandler) parseIDParam(c *gin.Context) (uint, error) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		return 0, err
	}
	return uint(id), nil
}

// handleError 统一处理错误响应
func (h *UserHandler) handleError(c *gin.Context, err error) {
	appErr := errors.GetAppError(err)
	if appErr != nil {
		response.Fail(c, appErr.Code, appErr.Message)
		return
	}
	logger.Error("处理请求时发生未知错误", logger.Err(err))
	response.Fail(c, errors.CodeInternalError, "服务器内部错误")
}
