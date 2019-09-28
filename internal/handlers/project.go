package handlers

import (
	"github.com/Nerufa/go-shared/logger"
	"github.com/Nerufa/go-shared/provider"
	"github.com/davecgh/go-spew/spew"
	"github.com/labstack/echo/v4"
	"github.com/paysuper/paysuper-billing-server/pkg"
	"github.com/paysuper/paysuper-billing-server/pkg/proto/billing"
	"github.com/paysuper/paysuper-billing-server/pkg/proto/grpc"
	"github.com/paysuper/paysuper-management-api/internal/dispatcher/common"
	"net/http"
)

const (
	projectsPath   = "/projects"
	projectsIdPath = "/projects/:id"
)

type ProjectRoute struct {
	dispatch common.HandlerSet
	cfg      common.Config
	provider.LMT
}

func NewProjectRoute(set common.HandlerSet, cfg *common.Config) *ProjectRoute {
	set.AwareSet.Logger = set.AwareSet.Logger.WithFields(logger.Fields{"router": "ProjectRoute"})
	return &ProjectRoute{
		dispatch: set,
		LMT:      &set.AwareSet,
		cfg:      *cfg,
	}
}

func (h *ProjectRoute) Route(groups *common.Groups) {
	groups.AuthUser.GET(projectsPath, h.listProjects)
	groups.AuthUser.GET(projectsIdPath, h.getProject)
	groups.AuthUser.POST(projectsPath, h.createProject)
	groups.AuthUser.PATCH(projectsIdPath, h.updateProject)
	groups.AuthUser.DELETE(projectsIdPath, h.deleteProject)
}

func (h *ProjectRoute) createProject(ctx echo.Context) error {
	req := &billing.Project{}
	err := ctx.Bind(req)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest /*ErrorRequestParamsIncorrect*/, err)
	}

	err = h.dispatch.Validate.Struct(req)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, common.GetValidationError(err))
	}

	res, err := h.dispatch.Services.Billing.ChangeProject(ctx.Request().Context(), req)

	if err != nil {
		h.L().Error(common.InternalErrorTemplate, logger.WithFields(logger.Fields{"err": err.Error()}))
		return echo.NewHTTPError(http.StatusInternalServerError, common.ErrorUnknown)
	}

	if res.Status != pkg.ResponseStatusOk {
		return echo.NewHTTPError(int(res.Status), res.Message)
	}

	return ctx.JSON(http.StatusCreated, res)
}

func (h *ProjectRoute) updateProject(ctx echo.Context) error {
	req := &billing.Project{}
	binder := common.NewChangeProjectRequestBinder(h.dispatch, h.cfg)
	err := binder.Bind(req, ctx)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest /*ErrorRequestParamsIncorrect*/, err)
	}

	err = h.dispatch.Validate.Struct(req)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, common.GetValidationError(err))
	}

	res, err := h.dispatch.Services.Billing.ChangeProject(ctx.Request().Context(), req)

	if err != nil {
		h.L().Error(common.InternalErrorTemplate, logger.WithFields(logger.Fields{"err": err.Error()}))
		return echo.NewHTTPError(http.StatusInternalServerError, common.ErrorUnknown)
	}

	if res.Status != pkg.ResponseStatusOk {
		return echo.NewHTTPError(int(res.Status), res.Message)
	}

	return ctx.JSON(http.StatusOK, res.Item)
}

func (h *ProjectRoute) getProject(ctx echo.Context) error {
	req := &grpc.GetProjectRequest{
		ProjectId: ctx.Param(common.RequestParameterId),
	}

	err := h.dispatch.Validate.Struct(req)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, common.GetValidationError(err))
	}

	res, err := h.dispatch.Services.Billing.GetProject(ctx.Request().Context(), req)

	if err != nil {
		h.L().Error(common.InternalErrorTemplate, logger.WithFields(logger.Fields{"err": err.Error()}))
		return echo.NewHTTPError(http.StatusInternalServerError, common.ErrorUnknown)
	}

	if res.Status != pkg.ResponseStatusOk {
		return echo.NewHTTPError(int(res.Status), res.Message)
	}

	return ctx.JSON(http.StatusOK, res)
}

func (h *ProjectRoute) listProjects(ctx echo.Context) error {
	req := &grpc.ListProjectsRequest{}
	err := ctx.Bind(req)

	spew.Dump(req)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, common.ErrorRequestParamsIncorrect)
	}

	if req.Limit <= 0 {
		req.Limit = h.cfg.LimitDefault
	}

	err = h.dispatch.Validate.Struct(req)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, common.GetValidationError(err))
	}

	res, err := h.dispatch.Services.Billing.ListProjects(ctx.Request().Context(), req)

	if err != nil {
		h.L().Error(common.InternalErrorTemplate, logger.WithFields(logger.Fields{"err": err.Error()}))
		return echo.NewHTTPError(http.StatusInternalServerError, common.ErrorUnknown)
	}

	return ctx.JSON(http.StatusOK, res)
}

func (h *ProjectRoute) deleteProject(ctx echo.Context) error {
	req := &grpc.GetProjectRequest{
		ProjectId: ctx.Param(common.RequestParameterId),
	}

	err := h.dispatch.Validate.Struct(req)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, common.GetValidationError(err))
	}

	res, err := h.dispatch.Services.Billing.DeleteProject(ctx.Request().Context(), req)

	if err != nil {
		h.L().Error(common.InternalErrorTemplate, logger.WithFields(logger.Fields{"err": err.Error()}))
		return echo.NewHTTPError(http.StatusInternalServerError, common.ErrorUnknown)
	}

	if res.Status != pkg.ResponseStatusOk {
		return echo.NewHTTPError(int(res.Status), res.Message)
	}

	return ctx.JSON(http.StatusOK, res)
}
