package user

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/opaulochaves/myserver/apperrors"
	"github.com/opaulochaves/myserver/pkg/pagination"
)

func RegisterHandlers(service Service) *chi.Mux {
	res := resource{service}
	r := chi.NewRouter()

	r.Get("/", res.list)    // GET /users - read a list of users
	r.Post("/", res.create) // POST /users - create a new user and persist it

	r.Route("/{id}", func(r chi.Router) {
		r.Use(res.userContext) // lets have a users map, and lets actually load/manipulate
		r.Get("/", res.get)    // GET /users/{id} - read a single user by :id
		// r.Put("/", res.Update)    // PUT /users/{id} - update a single user by :id
		// r.Delete("/", res.Delete) // DELETE /users/{id} - delete a single user by :id
	})

	return r
}

type resource struct {
	service Service
}

func (r resource) list(w http.ResponseWriter, req *http.Request) {
	count, err := r.service.Count()

	if err != nil {
		render.Render(w, req, apperrors.ErrInternalError(err))
		return
	}

	pages := pagination.NewFromRequest(req, count)
	users, err := r.service.Query(pages.Offset(), pages.Limit())

	if err != nil {
		render.Render(w, req, apperrors.ErrInternalError(err))
		return
	}

	pages.Items = users

	if err := render.Render(w, req, pages); err != nil {
		render.Render(w, req, apperrors.ErrRender(err))
		return
	}
}

func (c resource) create(w http.ResponseWriter, r *http.Request) {
	input := CreateUserRequest{}

	if err := render.Bind(r, &input); err != nil {
		render.Render(w, r, apperrors.ErrInvalidRequest(err))
		return
	}

	user, err := c.service.Create(input)

	if err != nil {
		render.Render(w, r, apperrors.ErrInternalError(err))
		return
	}

	render.Status(r, http.StatusCreated)
	render.Render(w, r, &UserResponse{User: user})
}

func (c resource) get(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user").(User)

	fmt.Println(">>>>>>", user)

	if err := render.Render(w, r, &UserResponse{User: user}); err != nil {
		render.Render(w, r, apperrors.ErrRender(err))
		return
	}
}

func (c resource) userContext(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID, err := strconv.Atoi(chi.URLParam(r, "id"))

		if err != nil {
			// TODO use logger
			log.Printf("err: %v", err)
			// TODO improve app errors definitions
			render.Render(w, r, apperrors.ErrInvalidRequest(err))
			return
		}

		user, err := c.service.Get(int64(userID))

		if err != nil {
			render.Render(w, r, apperrors.ErrNotFound)
			return
		}

		ctx := context.WithValue(r.Context(), "user", user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
