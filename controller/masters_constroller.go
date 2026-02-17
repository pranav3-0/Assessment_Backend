package controller

import (
	"dhl/constant"
	"dhl/models"
	"dhl/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type MastersController struct {
	svc                        services.DHLBusinessPartnerService
	centerServices             services.DHLCenterService
	resCompanyServices         services.DHLResCompanyService
	resPartnerIndustryServices services.DHLResPartnerIndustryService
	dhlServiceServices         services.DHLServiceService
	dhlServiceGroupServices    services.DHLServiceGroupService
	dhlServiceLineSerService   services.DHLServiceLineService
	subBusinessPartnerSvc      services.DHLSubBusinessPartnerService
	subServiceSvc              services.DHLSubServiceService
}

func NewMastersController(svc services.DHLBusinessPartnerService, centerServices services.DHLCenterService, resCompanyServices services.DHLResCompanyService,
	resPartnerIndustryServices services.DHLResPartnerIndustryService, dhlServiceServices services.DHLServiceService, dhlServiceGroupServices services.DHLServiceGroupService,
	dhlServiceLineSerService services.DHLServiceLineService, subBusinessPartnerSvc services.DHLSubBusinessPartnerService, subServiceSvc services.DHLSubServiceService,
) *MastersController {
	return &MastersController{svc: svc,
		centerServices:             centerServices,
		resCompanyServices:         resCompanyServices,
		resPartnerIndustryServices: resPartnerIndustryServices,
		dhlServiceServices:         dhlServiceServices,
		dhlServiceGroupServices:    dhlServiceGroupServices,
		dhlServiceLineSerService:   dhlServiceLineSerService,
		subBusinessPartnerSvc:      subBusinessPartnerSvc,
		subServiceSvc:              subServiceSvc,
	}
}

func (ctrl *MastersController) CreateDHLBusinessPartner(c *gin.Context) {
	var req models.DHLBusinessPartner

	if err := c.ShouldBindJSON(&req); err != nil {
		models.ErrorResponse(c, constant.Failure, http.StatusBadRequest, "invalid request", nil, err)
		return
	}

	if err := ctrl.svc.CreatePartner(c, req); err != nil {
		models.ErrorResponse(c, constant.Failure, http.StatusBadRequest, "failed to add partner", nil, err)
		return
	}

	models.SuccessResponse(c, constant.Success, http.StatusOK, "business partner created", nil, nil, nil)
	return
}

func (ctrl *MastersController) ListDHLBusinessPartners(c *gin.Context) {
	data, err := ctrl.svc.ListPartners(c)
	if err != nil {
		models.ErrorResponse(c, constant.Failure, http.StatusBadRequest, "failed to get partner", nil, err)
		return
	}

	models.SuccessResponse(c, constant.Success, http.StatusOK, "Partners fetched", data, nil, nil)
}

func (ctrl *MastersController) UpdateDHLBusinessPartner(c *gin.Context) {
	var req models.DHLBusinessPartner

	if err := c.ShouldBindJSON(&req); err != nil {
		models.ErrorResponse(c, constant.Failure, http.StatusBadRequest, "invalid request", nil, err)
		return
	}

	if err := ctrl.svc.UpdatePartner(c, req); err != nil {
		models.ErrorResponse(c, constant.Failure, http.StatusBadRequest, "failed to update partner", nil, err)
		return
	}

	models.SuccessResponse(c, constant.Success, http.StatusOK, "business partner updated", nil, nil, nil)
	return
}

func (ctrl *MastersController) DeleteDHLBusinessPartner(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		models.ErrorResponse(c, constant.Failure, http.StatusBadRequest, "invalid id", nil, err)
		return
	}

	if err := ctrl.svc.DeletePartner(c, id); err != nil {
		models.ErrorResponse(c, constant.Failure, http.StatusBadRequest, "failed to delete", nil, err)
		return
	}

	models.SuccessResponse(c, constant.Success, http.StatusOK, "Partner deleted", nil, nil, nil)
}

func (ctrl *MastersController) CreateDHLCenter(c *gin.Context) {
	var req models.DHLCenter

	if err := c.ShouldBindJSON(&req); err != nil {
		models.ErrorResponse(c, constant.Failure, http.StatusBadRequest, "invalid request", nil, err)
		return
	}

	if err := ctrl.centerServices.CreateCenter(c, req); err != nil {
		models.ErrorResponse(c, constant.Failure, http.StatusBadRequest, "failed to add center", nil, err)
		return
	}

	models.SuccessResponse(c, constant.Success, http.StatusOK, "center created", nil, nil, nil)
	return
}

func (ctrl *MastersController) ListDHLCenters(c *gin.Context) {
	data, err := ctrl.centerServices.ListCenters(c)
	if err != nil {
		models.ErrorResponse(c, constant.Failure, http.StatusBadRequest, "failed to get partner", nil, err)
		return
	}

	models.SuccessResponse(c, constant.Success, http.StatusOK, "Centers fetched", data, nil, nil)
}

func (ctrl *MastersController) UpdateDHLCenter(c *gin.Context) {
	var req models.DHLCenter

	if err := c.ShouldBindJSON(&req); err != nil {
		models.ErrorResponse(c, constant.Failure, http.StatusBadRequest, "invalid request", nil, err)
		return
	}

	if err := ctrl.centerServices.UpdateCenter(c, req); err != nil {
		models.ErrorResponse(c, constant.Failure, http.StatusBadRequest, "failed to update center", nil, err)
		return
	}

	models.SuccessResponse(c, constant.Success, http.StatusOK, "center updated", nil, nil, nil)
	return
}

func (ctrl *MastersController) DeleteDHLCenter(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		models.ErrorResponse(c, constant.Failure, http.StatusBadRequest, "invalid id", nil, err)
		return
	}

	if err := ctrl.centerServices.DeleteCenter(c, id); err != nil {
		models.ErrorResponse(c, constant.Failure, http.StatusBadRequest, "failed to delete", nil, err)
		return
	}

	models.SuccessResponse(c, constant.Success, http.StatusOK, "Center deleted", nil, nil, nil)
}

func (ctrl *MastersController) CreateDHLResCompany(c *gin.Context) {
	var req models.DHLResCompany

	if err := c.ShouldBindJSON(&req); err != nil {
		models.ErrorResponse(c, constant.Failure, http.StatusBadRequest, "invalid request", nil, err)
		return
	}

	if err := ctrl.resCompanyServices.CreateCompany(c, req); err != nil {
		models.ErrorResponse(c, constant.Failure, http.StatusBadRequest, "failed to add center", nil, err)
		return
	}

	models.SuccessResponse(c, constant.Success, http.StatusOK, "center created", nil, nil, nil)
	return
}

func (ctrl *MastersController) ListDHLResCompanys(c *gin.Context) {
	data, err := ctrl.resCompanyServices.ListCompanies(c)
	if err != nil {
		models.ErrorResponse(c, constant.Failure, http.StatusBadRequest, "failed to get partner", nil, err)
		return
	}

	models.SuccessResponse(c, constant.Success, http.StatusOK, "Companys fetched", data, nil, nil)
}

func (ctrl *MastersController) UpdateDHLResCompany(c *gin.Context) {
	var req models.DHLResCompany

	if err := c.ShouldBindJSON(&req); err != nil {
		models.ErrorResponse(c, constant.Failure, http.StatusBadRequest, "invalid request", nil, err)
		return
	}

	if err := ctrl.resCompanyServices.UpdateCompany(c, req); err != nil {
		models.ErrorResponse(c, constant.Failure, http.StatusBadRequest, "failed to update partner", nil, err)
		return
	}

	models.SuccessResponse(c, constant.Success, http.StatusOK, "business partner updated", nil, nil, nil)
	return
}

func (ctrl *MastersController) DeleteDHLResCompany(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		models.ErrorResponse(c, constant.Failure, http.StatusBadRequest, "invalid id", nil, err)
		return
	}

	if err := ctrl.resCompanyServices.DeleteCompany(c, id); err != nil {
		models.ErrorResponse(c, constant.Failure, http.StatusBadRequest, "failed to delete", nil, err)
		return
	}

	models.SuccessResponse(c, constant.Success, http.StatusOK, "Company deleted", nil, nil, nil)
}

func (ctrl *MastersController) CreateDHLResPartnerIndustry(c *gin.Context) {
	var req models.DHLResPartnerIndustry

	if err := c.ShouldBindJSON(&req); err != nil {
		models.ErrorResponse(c, constant.Failure, http.StatusBadRequest, "invalid request", nil, err)
		return
	}

	if err := ctrl.resPartnerIndustryServices.Create(c, req); err != nil {
		models.ErrorResponse(c, constant.Failure, http.StatusBadRequest, "failed to create industry", nil, err)
		return
	}

	models.SuccessResponse(c, constant.Success, http.StatusOK, "Industry created", nil, nil, nil)
}

func (ctrl *MastersController) ListDHLResPartnerIndustry(c *gin.Context) {
	data, err := ctrl.resPartnerIndustryServices.List(c)
	if err != nil {
		models.ErrorResponse(c, constant.Failure, http.StatusBadRequest, "failed to fetch list", nil, err)
		return
	}
	models.SuccessResponse(c, constant.Success, http.StatusOK, "Industry list fetched", data, nil, nil)
}

func (ctrl *MastersController) UpdateDHLResPartnerIndustry(c *gin.Context) {
	var req models.DHLResPartnerIndustry

	if err := c.ShouldBindJSON(&req); err != nil {
		models.ErrorResponse(c, constant.Failure, http.StatusBadRequest, "invalid request", nil, err)
		return
	}

	if req.PartnerIndustryID == 0 {
		models.ErrorResponse(c, constant.Failure, http.StatusBadRequest, "partner_industry_id required", nil, nil)
		return
	}

	if err := ctrl.resPartnerIndustryServices.Update(c, req); err != nil {
		models.ErrorResponse(c, constant.Failure, http.StatusBadRequest, "failed to update", nil, err)
		return
	}

	models.SuccessResponse(c, constant.Success, http.StatusOK, "Industry updated", nil, nil, nil)
}

func (ctrl *MastersController) DeleteDHLResPartnerIndustry(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		models.ErrorResponse(c, constant.Failure, http.StatusBadRequest, "invalid id", nil, err)
		return
	}

	if err := ctrl.resPartnerIndustryServices.Delete(c, id); err != nil {
		models.ErrorResponse(c, constant.Failure, http.StatusBadRequest, "failed to delete", nil, err)
		return
	}

	models.SuccessResponse(c, constant.Success, http.StatusOK, "Industry deleted", nil, nil, nil)
}

func (ctrl *MastersController) CreateDHLService(c *gin.Context) {
	var req models.DHLService

	if err := c.ShouldBindJSON(&req); err != nil {
		models.ErrorResponse(c, constant.Failure, http.StatusBadRequest, "invalid request", nil, err)
		return
	}

	if err := ctrl.dhlServiceServices.Create(c, req); err != nil {
		models.ErrorResponse(c, constant.Failure, http.StatusBadRequest, "failed to create service", nil, err)
		return
	}

	models.SuccessResponse(c, constant.Success, http.StatusOK, "Service created", nil, nil, nil)
}

func (ctrl *MastersController) ListDHLServices(c *gin.Context) {
	data, err := ctrl.dhlServiceServices.List(c)
	if err != nil {
		models.ErrorResponse(c, constant.Failure, http.StatusBadRequest, "failed to fetch services", nil, err)
		return
	}
	models.SuccessResponse(c, constant.Success, http.StatusOK, "Services fetched", data, nil, nil)
}

func (ctrl *MastersController) UpdateDHLService(c *gin.Context) {
	var req models.DHLService

	if err := c.ShouldBindJSON(&req); err != nil {
		models.ErrorResponse(c, constant.Failure, http.StatusBadRequest, "invalid request", nil, err)
		return
	}

	if req.ServiceID == 0 {
		models.ErrorResponse(c, constant.Failure, http.StatusBadRequest, "service_id required", nil, nil)
		return
	}

	if err := ctrl.dhlServiceServices.Update(c, req); err != nil {
		models.ErrorResponse(c, constant.Failure, http.StatusBadRequest, "failed to update service", nil, err)
		return
	}

	models.SuccessResponse(c, constant.Success, http.StatusOK, "Service updated", nil, nil, nil)
}

func (ctrl *MastersController) DeleteDHLService(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		models.ErrorResponse(c, constant.Failure, http.StatusBadRequest, "invalid id", nil, err)
		return
	}

	if err := ctrl.dhlServiceServices.Delete(c, id); err != nil {
		models.ErrorResponse(c, constant.Failure, http.StatusBadRequest, "failed to delete service", nil, err)
		return
	}

	models.SuccessResponse(c, constant.Success, http.StatusOK, "Service deleted", nil, nil, nil)
}

func (ctrl *MastersController) CreateDHLServiceGroup(c *gin.Context) {
	var req models.DHLServiceGroup

	if err := c.ShouldBindJSON(&req); err != nil {
		models.ErrorResponse(c, constant.Failure, http.StatusBadRequest, "invalid request", nil, err)
		return
	}

	if err := ctrl.dhlServiceGroupServices.Create(c, req); err != nil {
		models.ErrorResponse(c, constant.Failure, http.StatusBadRequest, "failed to create service group", nil, err)
		return
	}

	models.SuccessResponse(c, constant.Success, http.StatusOK, "Service group created", nil, nil, nil)
}

func (ctrl *MastersController) ListDHLServiceGroups(c *gin.Context) {
	data, err := ctrl.dhlServiceGroupServices.List(c)
	if err != nil {
		models.ErrorResponse(c, constant.Failure, http.StatusBadRequest, "failed to fetch groups", nil, err)
		return
	}
	models.SuccessResponse(c, constant.Success, http.StatusOK, "Service groups fetched", data, nil, nil)
}

func (ctrl *MastersController) UpdateDHLServiceGroup(c *gin.Context) {
	var req models.DHLServiceGroup

	if err := c.ShouldBindJSON(&req); err != nil {
		models.ErrorResponse(c, constant.Failure, http.StatusBadRequest, "invalid request", nil, err)
		return
	}

	if req.ServiceGrpID == 0 {
		models.ErrorResponse(c, constant.Failure, http.StatusBadRequest, "service_grp_id required", nil, nil)
		return
	}

	if err := ctrl.dhlServiceGroupServices.Update(c, req); err != nil {
		models.ErrorResponse(c, constant.Failure, http.StatusBadRequest, "failed to update group", nil, err)
		return
	}

	models.SuccessResponse(c, constant.Success, http.StatusOK, "Service group updated", nil, nil, nil)
}

func (ctrl *MastersController) DeleteDHLServiceGroup(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		models.ErrorResponse(c, constant.Failure, http.StatusBadRequest, "invalid id", nil, err)
		return
	}

	if err := ctrl.dhlServiceGroupServices.Delete(c, id); err != nil {
		models.ErrorResponse(c, constant.Failure, http.StatusBadRequest, "failed to delete group", nil, err)
		return
	}

	models.SuccessResponse(c, constant.Success, http.StatusOK, "Service group deleted", nil, nil, nil)
}

func (ctrl *MastersController) CreateDHLServiceLine(c *gin.Context) {
	var req models.DHLServiceLine

	if err := c.ShouldBindJSON(&req); err != nil {
		models.ErrorResponse(c, constant.Failure, http.StatusBadRequest, "invalid request", nil, err)
		return
	}

	if err := ctrl.dhlServiceLineSerService.Create(c, req); err != nil {
		models.ErrorResponse(c, constant.Failure, http.StatusBadRequest, "failed to create service line", nil, err)
		return
	}

	models.SuccessResponse(c, constant.Success, http.StatusOK, "Service line created", nil, nil, nil)
}

func (ctrl *MastersController) ListDHLServiceLine(c *gin.Context) {
	data, err := ctrl.dhlServiceLineSerService.List(c)
	if err != nil {
		models.ErrorResponse(c, constant.Failure, http.StatusBadRequest, "failed to fetch service lines", nil, err)
		return
	}

	models.SuccessResponse(c, constant.Success, http.StatusOK, "Service lines fetched", data, nil, nil)
}

func (ctrl *MastersController) UpdateDHLServiceLine(c *gin.Context) {
	var req models.DHLServiceLine

	if err := c.ShouldBindJSON(&req); err != nil {
		models.ErrorResponse(c, constant.Failure, http.StatusBadRequest, "invalid request", nil, err)
		return
	}

	if req.ServiceLineID == 0 {
		models.ErrorResponse(c, constant.Failure, http.StatusBadRequest, "service_line_id is required", nil, nil)
		return
	}

	if err := ctrl.dhlServiceLineSerService.Update(c, req); err != nil {
		models.ErrorResponse(c, constant.Failure, http.StatusBadRequest, "failed to update service line", nil, err)
		return
	}

	models.SuccessResponse(c, constant.Success, http.StatusOK, "Service line updated", nil, nil, nil)
}

func (ctrl *MastersController) DeleteDHLServiceLine(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		models.ErrorResponse(c, constant.Failure, http.StatusBadRequest, "invalid id", nil, err)
		return
	}

	if err := ctrl.dhlServiceLineSerService.Delete(c, id); err != nil {
		models.ErrorResponse(c, constant.Failure, http.StatusBadRequest, "failed to delete", nil, err)
		return
	}

	models.SuccessResponse(c, constant.Success, http.StatusOK, "Service line deleted", nil, nil, nil)
}

func (ctrl *MastersController) CreateDHLSubBusinessPartner(c *gin.Context) {
	var req models.DHLSubBusinessPartner

	if err := c.ShouldBindJSON(&req); err != nil {
		models.ErrorResponse(c, constant.Failure, http.StatusBadRequest, "invalid request", nil, err)
		return
	}

	if err := ctrl.subBusinessPartnerSvc.Create(c, req); err != nil {
		models.ErrorResponse(c, constant.Failure, http.StatusBadRequest, "failed to add sub business partner", nil, err)
		return
	}

	models.SuccessResponse(c, constant.Success, http.StatusOK, "sub business partner created", nil, nil, nil)
}

func (ctrl *MastersController) ListDHLSubBusinessPartner(c *gin.Context) {
	data, err := ctrl.subBusinessPartnerSvc.List(c)
	if err != nil {
		models.ErrorResponse(c, constant.Failure, http.StatusBadRequest, "failed to fetch list", nil, err)
		return
	}

	models.SuccessResponse(c, constant.Success, http.StatusOK, "sub business partner list fetched", data, nil, nil)
}

func (ctrl *MastersController) UpdateDHLSubBusinessPartner(c *gin.Context) {
	var req models.DHLSubBusinessPartner

	if err := c.ShouldBindJSON(&req); err != nil {
		models.ErrorResponse(c, constant.Failure, http.StatusBadRequest, "invalid request", nil, err)
		return
	}

	if req.SubBusinessPartnerID == 0 {
		models.ErrorResponse(c, constant.Failure, http.StatusBadRequest, "sub_business_partner_id required", nil, nil)
		return
	}

	if err := ctrl.subBusinessPartnerSvc.Update(c, req); err != nil {
		models.ErrorResponse(c, constant.Failure, http.StatusBadRequest, "failed to update", nil, err)
		return
	}

	models.SuccessResponse(c, constant.Success, http.StatusOK, "sub business partner updated", nil, nil, nil)
}

func (ctrl *MastersController) DeleteDHLSubBusinessPartner(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		models.ErrorResponse(c, constant.Failure, http.StatusBadRequest, "invalid id", nil, err)
		return
	}

	if err := ctrl.subBusinessPartnerSvc.Delete(c, id); err != nil {
		models.ErrorResponse(c, constant.Failure, http.StatusBadRequest, "failed to delete", nil, err)
		return
	}

	models.SuccessResponse(c, constant.Success, http.StatusOK, "sub business partner deleted", nil, nil, nil)
}

func (ctrl *MastersController) CreateDHLSubService(c *gin.Context) {
	var req models.DHLSubService

	if err := c.ShouldBindJSON(&req); err != nil {
		models.ErrorResponse(c, constant.Failure, http.StatusBadRequest, "invalid request", nil, err)
		return
	}

	if err := ctrl.subServiceSvc.Create(c, req); err != nil {
		models.ErrorResponse(c, constant.Failure, http.StatusBadRequest, "failed to create sub service", nil, err)
		return
	}

	models.SuccessResponse(c, constant.Success, http.StatusOK, "sub service created", nil, nil, nil)
}

func (ctrl *MastersController) ListDHLSubService(c *gin.Context) {
	data, err := ctrl.subServiceSvc.List(c)
	if err != nil {
		models.ErrorResponse(c, constant.Failure, http.StatusBadRequest, "failed to fetch list", nil, err)
		return
	}

	models.SuccessResponse(c, constant.Success, http.StatusOK, "sub service list fetched", data, nil, nil)
}

func (ctrl *MastersController) UpdateDHLSubService(c *gin.Context) {
	var req models.DHLSubService

	if err := c.ShouldBindJSON(&req); err != nil {
		models.ErrorResponse(c, constant.Failure, http.StatusBadRequest, "invalid request", nil, err)
		return
	}

	if req.SubServiceID == 0 {
		models.ErrorResponse(c, constant.Failure, http.StatusBadRequest, "sub_service_id required", nil, nil)
		return
	}

	if err := ctrl.subServiceSvc.Update(c, req); err != nil {
		models.ErrorResponse(c, constant.Failure, http.StatusBadRequest, "failed to update sub service", nil, err)
		return
	}

	models.SuccessResponse(c, constant.Success, http.StatusOK, "sub service updated", nil, nil, nil)
}

func (ctrl *MastersController) DeleteDHLSubService(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		models.ErrorResponse(c, constant.Failure, http.StatusBadRequest, "invalid id", nil, err)
		return
	}

	if err := ctrl.subServiceSvc.Delete(c, id); err != nil {
		models.ErrorResponse(c, constant.Failure, http.StatusBadRequest, "failed to delete sub service", nil, err)
		return
	}

	models.SuccessResponse(c, constant.Success, http.StatusOK, "sub service deleted", nil, nil, nil)
}
