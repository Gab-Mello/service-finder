package review

import (
	"encoding/json"
	"net/http"

	authmw "github.com/Gab-Mello/service-finder/internal/http/middleware/auth"
	"github.com/Gab-Mello/service-finder/internal/http/response"
	domain "github.com/Gab-Mello/service-finder/internal/review"
)

const basePath = "/api/v1/reviews/"

type Handler struct{ svc *domain.Service }

func NewHandler(s *domain.Service) *Handler { return &Handler{svc: s} }

type createReq struct {
	OrderID string `json:"orderId"`
	Stars   int    `json:"stars"`
	Comment string `json:"comment"`
}
type editReq struct {
	Stars   int    `json:"stars"`
	Comment string `json:"comment"`
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	uid, ok := authmw.UserIDFromContext(r)
	if !ok {
		response.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	var req createReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.OrderID == "" {
		response.Error(w, http.StatusBadRequest, "invalid json")
		return
	}
	rv, err := h.svc.Create(uid, req.OrderID, req.Stars, req.Comment)
	if err != nil {
		response.Error(w, statusFor(err), err.Error())
		return
	}
	response.JSON(w, http.StatusCreated, rv)
}

func (h *Handler) Edit(w http.ResponseWriter, r *http.Request) {
	uid, ok := authmw.UserIDFromContext(r)
	if !ok {
		response.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	orderID := response.PathParam(r.URL.Path, basePath, "")
	var req editReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid json")
		return
	}
	rv, err := h.svc.Edit(uid, orderID, req.Stars, req.Comment)
	if err != nil {
		response.Error(w, statusFor(err), err.Error())
		return
	}
	response.JSON(w, http.StatusOK, rv)
}

func statusFor(err error) int {
	switch err {
	case domain.ErrForbidden:
		return http.StatusForbidden
	case domain.ErrInvalidFields, domain.ErrAlreadyExists, domain.ErrEditWindowOver, domain.ErrOrderNotDone:
		return http.StatusBadRequest
	case domain.ErrNotFound:
		return http.StatusNotFound
	default:
		return http.StatusInternalServerError
	}
}
