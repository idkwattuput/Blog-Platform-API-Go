package post

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/idkwattuput/blogging-platform-api-go/types"
	"github.com/idkwattuput/blogging-platform-api-go/utils"
)

type Handler struct {
	store types.PostStore
}

func NewHandler(store types.PostStore) *Handler {
	return &Handler{store: store}
}

func (h *Handler) RegisterRoutes(router *http.ServeMux) {
	router.HandleFunc("GET /posts", h.handleGetPosts)
	router.HandleFunc("GET /posts/{id}", h.handleGetPost)
	router.HandleFunc("POST /posts", h.handleCreatePost)
	router.HandleFunc("PUT /posts/{id}", h.handleUpdatePost)
	router.HandleFunc("DELETE /posts/{id}", h.handleDeletePost)
}

func (h *Handler) handleGetPosts(w http.ResponseWriter, r *http.Request) {
	posts, err := h.store.GetPosts()
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	response := map[string]interface{}{
		"data": posts,
	}
	utils.WriteJSON(w, http.StatusOK, response)
}

func (h *Handler) handleGetPost(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	post, err := h.store.GetPostById(id)
	if err != nil {
		if err == sql.ErrNoRows {
			utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("Post not found"))
			return
		}
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	response := map[string]interface{}{
		"data": post,
	}
	utils.WriteJSON(w, http.StatusOK, response)
}

func (h *Handler) handleCreatePost(w http.ResponseWriter, r *http.Request) {
	var post types.PostPayload
	if err := utils.ParseJSON(r, &post); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	if err := utils.Validate.Struct(post); err != nil {
		errors := err.(validator.ValidationErrors)
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("Invalid payload: %v", errors))
		return
	}

	newPost, err := h.store.CreatePost(post)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	response := map[string]interface{}{
		"data": newPost,
	}
	utils.WriteJSON(w, http.StatusOK, response)
}

func (h *Handler) handleUpdatePost(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	_, err = h.store.GetPostById(id)
	if err != nil {
		if err == sql.ErrNoRows {
			utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("Post not found"))
			return
		}
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	var post types.PostPayload
	if err := utils.ParseJSON(r, &post); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	if err := utils.Validate.Struct(post); err != nil {
		errors := err.(validator.ValidationErrors)
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("Invalid payload: %v", errors))
		return
	}

	updatedPost, err := h.store.UpdatePost(id, post)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	response := map[string]interface{}{
		"data": updatedPost,
	}
	utils.WriteJSON(w, http.StatusOK, response)
}

func (h *Handler) handleDeletePost(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	_, err = h.store.GetPostById(id)
	if err != nil {
		if err == sql.ErrNoRows {
			utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("Post not found"))
			return
		}
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	err = h.store.DeletePost(id)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	utils.WriteJSON(w, http.StatusNoContent, "OK")
}
