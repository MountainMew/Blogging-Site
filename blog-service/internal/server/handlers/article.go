package handlers

import (
	"blog-service/internal/db/mongo"
	"blog-service/internal/db/postgres"
	"blog-service/internal/server/models"
	"encoding/json"
	"errors"
	"net/http"
	"regexp"
	"strconv"
)

type ArticleHandler struct {
	MongoDB    *mongo.Client
	PostgresDB *postgres.Client
}

const commentsPerPage = 10

var (
	ArticleIDRe        = regexp.MustCompile(`/article/[a-f0-9]{24}/$`)
	ArticleIDNoSlashRe = regexp.MustCompile(`/article/[a-f0-9]{24}$`)
)

func (h *ArticleHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost && (r.URL.Path == "/article" || r.URL.Path == "/article/") {
		h.ArticleCreate(w, r)
		return
	}

	if ArticleIDRe.MatchString(r.URL.Path) || ArticleIDNoSlashRe.MatchString(r.URL.Path) {
		switch r.Method {
		case http.MethodGet:
			h.ArticleGet(w, r)
			return
		case http.MethodDelete:
			h.ArticleDelete(w, r)
			return
		case http.MethodPut:
			h.ArticleUpdate(w, r)
			return
		}
	}
	http.NotFound(w, r)
}

func (h *ArticleHandler) ArticleGet(w http.ResponseWriter, r *http.Request) {
	articleID := r.PathValue("id")
	pageParam := r.URL.Query().Get("page")
	var page int
	var err error

	if pageParam == "" {
		page = 1
	} else {
		page, err = strconv.Atoi(pageParam)
		if err != nil || page <= 0 {
			page = 1
		}
	}

	article, err := models.GetArticleByID(r.Context(), h.MongoDB, articleID)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	comments, err := models.GetCommentsByArticleID(r.Context(), h.PostgresDB, articleID, commentsPerPage, page)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	article.Comments = comments

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(article)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (h *ArticleHandler) ArticleUpdate(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
}

func (h *ArticleHandler) ArticleDelete(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
}

func (h *ArticleHandler) ArticleCreate(w http.ResponseWriter, r *http.Request) {
	var article models.ArticleCreateDTO
	var paramError *models.ParamError
	var UnauthorizedErr *models.UnauthorizedError

	err := json.NewDecoder(r.Body).Decode(&article)
	if err != nil {
		http.Error(w, paramError.Error(), http.StatusBadRequest)
		return
	}

	id, err := models.CreateArticle(r.Context(), h.MongoDB, &article)
	switch {
	case errors.As(err, &paramError):
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	case errors.As(err, &UnauthorizedErr):
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	case err != nil:
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	err = json.NewEncoder(w).Encode(models.ArticleCreateResponse{ID: id})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}
