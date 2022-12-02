package users

import (
	"context"
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/opaulochaves/myserver/apperrors"
)

type userController struct {
	UserService UserService
}

type createUserRequest struct {
	*User

	ID        bool `json:"id,omitempty"`
	CreatedAt bool `json:"createdAt,omitempty"`
	UpdatedAt bool `json:"updatedAt,omitempty"`
} //@name CreateUserRequest

type userResponse struct {
	*User
	// Password bool `json:"password,omitempty"`
} //@name UserResponse

type UCConfig struct {
	UserService UserService
}

type UserController interface {
	Routes() chi.Router
	List(w http.ResponseWriter, r *http.Request)
	Create(w http.ResponseWriter, r *http.Request)
	Get(w http.ResponseWriter, r *http.Request)
	Update(w http.ResponseWriter, r *http.Request)
	Delete(w http.ResponseWriter, r *http.Request)
}

func NewUserController(c *UCConfig) UserController {
	return &userController{
		UserService: c.UserService,
	}
}

func (uc *userController) Routes() chi.Router {
	r := chi.NewRouter()

	r.Get("/", uc.List)    // GET /users - read a list of users
	r.Post("/", uc.Create) // POST /users - create a new user and persist it

	r.Route("/{id}", func(r chi.Router) {
		r.Use(uc.UserContext)    // lets have a users map, and lets actually load/manipulate
		r.Get("/", uc.Get)       // GET /users/{id} - read a single user by :id
		r.Put("/", uc.Update)    // PUT /users/{id} - update a single user by :id
		r.Delete("/", uc.Delete) // DELETE /users/{id} - delete a single user by :id
	})

	return r
}

func (uc userController) List(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("users list of stuff.."))
}

func (uc userController) Create(w http.ResponseWriter, r *http.Request) {
	data := &createUserRequest{}

	if err := render.Bind(r, data); err != nil {
		render.Render(w, r, apperrors.ErrInvalidRequest(err))
		return
	}

	dataUser := data.User
	user, err := uc.UserService.CreateUser(dataUser)

	if err != nil {
		render.Render(w, r, apperrors.ErrInternalError(err))
		return
	}

	render.Status(r, http.StatusCreated)
	render.Render(w, r, &userResponse{User: user})
}

func (uc userController) Get(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user").(*User)

	if err := render.Render(w, r, &userResponse{User: user}); err != nil {
		render.Render(w, r, apperrors.ErrRender(err))
		return
	}
}

func (uc userController) Update(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("user update"))
}

func (uc userController) Delete(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("user delete"))
}

func (uc userController) UserContext(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var user *User
		var err error

		userID, err := strconv.Atoi(chi.URLParam(r, "id"))

		if err != nil {
			log.Printf("err: %v", err)
			render.Render(w, r, apperrors.ErrInvalidRequest(err))
			return
		}

		user, err = uc.UserService.Get(userID)

		if err != nil {
			log.Printf("err 2: %v", err)
			render.Render(w, r, apperrors.ErrNotFound)
			return
		}

		ctx := context.WithValue(r.Context(), "user", user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (u *createUserRequest) Bind(r *http.Request) (err error) {
	if u.User == nil {
		err = errors.New("missing required User fields.")
	}

	return
}

func (rd *userResponse) Render(w http.ResponseWriter, r *http.Request) (err error) {
	// Pre-processing before a response is marshalled and sent across the wire
	return
}
