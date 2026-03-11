package order

import (
	"encoding/json"
	"net/http"
	"time"

	authmw "github.com/Gab-Mello/service-finder/internal/http/middleware/auth"
	"github.com/Gab-Mello/service-finder/internal/http/response"
	domain "github.com/Gab-Mello/service-finder/internal/order"
)

const basePath = "/api/v1/orders/"

type Handler struct{ svc *domain.Service }

func NewHandler(s *domain.Service) *Handler { return &Handler{svc: s} }

type requestOrder struct {
	PostingID  string `json:"postingId"`
	ProviderID string `json:"providerId"`
}
type acceptReq struct {
	ScheduledAt string `json:"scheduledAt"`
}

func (h *Handler) Request(w http.ResponseWriter, r *http.Request) {
	uid, ok := authmw.UserIDFromContext(r)
	if !ok {
		response.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}
	var req requestOrder
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.PostingID == "" || req.ProviderID == "" {
		response.Error(w, http.StatusBadRequest, "invalid json")
		return
	}
	o, err := h.svc.Request(uid, req.PostingID, req.ProviderID)
	if err != nil {
		response.Error(w, statusFor(err), err.Error())
		return
	}
	response.JSON(w, http.StatusCreated, o)
}

func (h *Handler) Accept(w http.ResponseWriter, r *http.Request) {
	uid, ok := authmw.UserIDFromContext(r)
	if !ok {
		response.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}
	id := response.PathParam(r.URL.Path, basePath, "/accept")

	var req acceptReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.ScheduledAt == "" {
		response.Error(w, http.StatusBadRequest, "scheduledAt required (RFC3339)")
		return
	}
	when, err := time.Parse(time.RFC3339, req.ScheduledAt)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "invalid scheduledAt")
		return
	}

	o, err := h.svc.Accept(uid, id, when)
	if err != nil {
		response.Error(w, statusFor(err), err.Error())
		return
	}
	response.JSON(w, http.StatusOK, o)
}

func (h *Handler) Start(w http.ResponseWriter, r *http.Request) {
	uid, ok := authmw.UserIDFromContext(r)
	if !ok {
		response.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}
	id := response.PathParam(r.URL.Path, basePath, "/start")

	o, err := h.svc.Start(uid, id)
	if err != nil {
		response.Error(w, statusFor(err), err.Error())
		return
	}
	response.JSON(w, http.StatusOK, o)
}

func (h *Handler) Complete(w http.ResponseWriter, r *http.Request) {
	uid, ok := authmw.UserIDFromContext(r)
	if !ok {
		response.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}
	id := response.PathParam(r.URL.Path, basePath, "/complete")

	o, err := h.svc.Complete(uid, id)
	if err != nil {
		response.Error(w, statusFor(err), err.Error())
		return
	}
	response.JSON(w, http.StatusOK, o)
}

func (h *Handler) Cancel(w http.ResponseWriter, r *http.Request) {
	uid, ok := authmw.UserIDFromContext(r)
	if !ok {
		response.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}
	id := response.PathParam(r.URL.Path, basePath, "/cancel")

	o, err := h.svc.Cancel(uid, id)
	if err != nil {
		response.Error(w, statusFor(err), err.Error())
		return
	}
	response.JSON(w, http.StatusOK, o)
}

func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	uid, ok := authmw.UserIDFromContext(r)
	if !ok {
		response.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}
	id := response.PathParam(r.URL.Path, basePath, "")

	o, err := h.svc.GetForUser(uid, id)
	if err != nil {
		response.Error(w, statusFor(err), err.Error())
		return
	}
	response.JSON(w, http.StatusOK, o)
}

func (h *Handler) ListMine(w http.ResponseWriter, r *http.Request) {
	uid, ok := authmw.UserIDFromContext(r)
	if !ok {
		response.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}
	list, err := h.svc.ListMine(uid)
	if err != nil {
		response.InternalError(w, err)
		return
	}
	response.JSON(w, http.StatusOK, list)
}

func statusFor(err error) int {
	switch err {
	case domain.ErrForbidden:
		return http.StatusForbidden
	case domain.ErrInvalidState, domain.ErrInvalidFields:
		return http.StatusBadRequest
	case domain.ErrNotFound:
		return http.StatusNotFound
	default:
		return http.StatusInternalServerError
	}
}
