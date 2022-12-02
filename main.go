package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
	"github.com/opaulochaves/myserver/config"
	"github.com/opaulochaves/myserver/users"
)

func main() {
	log.Println("Starting server...")

	// TODO: Only use dotenv if not production
	err := godotenv.Load()
	if err != nil {
		log.Fatalln("Error loading .env file")
	}

	ctx := context.Background()
	cfg, err := config.LoadConfig(ctx)

	if err != nil {
		log.Fatalf("Could not load the config: %v\n", err)
	}

	// initialize data sources
	ds, err := initDS(ctx, cfg)

	defer ds.DBPool.Close()

	if err != nil {
		log.Fatalf("Unable to initialize data sources: %v\n", err)
	}

	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)

	userRepository := users.NewUserRepository(ds.DBPool)

	userService := users.NewUserService(&users.USConfig{
		UserRepository: userRepository,
	})

	userController := users.NewUserController(&users.UCConfig{
		UserService: userService,
	})

	router.Mount("/api/users", userController.Routes())

	server := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: router,
	}

	log.Printf("Listening on port %v\n", server.Addr)

	// Server run context
	serverCtx, serverStopCtx := context.WithCancel(context.Background())

	// Listen for syscall signals for process to interrupt/quit
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	go func() {
		<-sig

		// Shutdown signal with grace period of 30 seconds
		shutdownCtx, shutdownCancel := context.WithTimeout(serverCtx, 30*time.Second)
		defer shutdownCancel()

		go func() {
			<-shutdownCtx.Done()
			if shutdownCtx.Err() == context.DeadlineExceeded {
				log.Fatal("graceful shutdown timed out.. forcing exit.")
			}
		}()

		// Trigger graceful shutdown
		err := server.Shutdown(shutdownCtx)
		if err != nil {
			log.Fatal(err)
		}
		serverStopCtx()
	}()

	err = server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}

	log.Println("Shutting down server...")

	// Wait for server context to be stopped
	<-serverCtx.Done()
}

// func GetNote(w http.ResponseWriter, r *http.Request) {
// 	note := r.Context().Value("note").(*Note)
//
// 	if err := render.Render(w, r, NewNoteResponse(note)); err != nil {
// 		render.Render(w, r, ErrRender(err))
// 		return
// 	}
// }
//
// // NoteRequest is the request payload for Note data model.
// //
// // NOTE: It's good practice to have well defined request and response payloads
// // so you can manage the specific inputs and outputs for clients, and also gives
// // you the opportunity to transform data on input or output, for example
// // on request, we'd like to protect certain fields and on output perhaps
// // we'd like to include a computed field based on other values that aren't
// // in the data model. Also, check out this awesome blog post on struct composition:
// // http://attilaolah.eu/2014/09/10/json-and-struct-composition-in-go/
// type NoteRequest struct {
// 	*Note
//
// 	User *UserPayload `json:"user,omitempty"`
//
// 	ProtectedID string `json:"id"` // override 'id' json to have more control
// }
//
// func CreateNote(w http.ResponseWriter, r *http.Request) {
// 	data := &NoteRequest{}
// 	if err := render.Bind(r, data); err != nil {
// 		render.Render(w, r, ErrInvalidRequest(err))
// 		return
// 	}
//
// 	note := data.Note
// 	dbNewNote(note)
//
// 	render.Status(r, http.StatusCreated)
// 	render.Render(w, r, NewNoteResponse(note))
// }
//
// func (n *NoteRequest) Bind(r *http.Request) error {
// 	// n.Note is nil if no Note fields are sent in the request. Return an
// 	// error to avoid a nil pointer dereference.
// 	if n.Note == nil {
// 		return errors.New("missing required Note fields.")
// 	}
//
// 	// just a post-process after a decode..
// 	n.ProtectedID = ""
// 	n.Note.Title = strings.ToLower(n.Note.Title)
// 	return nil
// }
//
// // NoteCtx middleware is used to load an Note object from
// // the URL parameters passed through as the request. In case
// // the Note could not be found, we stop here and return a 404.
// func NoteCtx(next http.Handler) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		var note *Note
// 		var err error
//
// 		if noteID := chi.URLParam(r, "noteID"); noteID != "" {
// 			note, err = dbGetNote(noteID)
// 		} else {
// 			render.Render(w, r, ErrNotFound)
// 			return
// 		}
// 		if err != nil {
// 			render.Render(w, r, ErrNotFound)
// 			return
// 		}
//
// 		ctx := context.WithValue(r.Context(), "note", note)
// 		next.ServeHTTP(w, r.WithContext(ctx))
// 	})
// }
//
// func ListNotes(w http.ResponseWriter, r *http.Request) {
// 	if err := render.RenderList(w, r, NewNoteListResponse(notes)); err != nil {
// 		render.Render(w, r, ErrRender(err))
// 		return
// 	}
// }
//
// type UserPayload struct {
// 	*User
//
// 	Role string `json:"role"`
// }
//
// func NewUserPayloadResponse(user *User) *UserPayload {
// 	return &UserPayload{User: user}
// }
//
// type NoteResponse struct {
// 	*Note
//
// 	User *UserPayload `json:"user,omitempty"`
//
// 	Elapsed int64 `json:"elapsed"`
// }
//
// func NewNoteResponse(note *Note) *NoteResponse {
// 	resp := &NoteResponse{Note: note}
//
// 	if resp.User == nil {
// 		if user, _ := dbGetUser(resp.UserID); user != nil {
// 			resp.User = NewUserPayloadResponse(user)
// 		}
// 	}
//
// 	return resp
// }
//
// func (rd *NoteResponse) Render(w http.ResponseWriter, r *http.Request) error {
// 	// Pre-processing before a response is marshalled and sent across the wire
// 	rd.Elapsed = 10
// 	return nil
// }
//
// func NewNoteListResponse(notes []*Note) []render.Renderer {
// 	list := []render.Renderer{}
// 	for _, note := range notes {
// 		list = append(list, NewNoteResponse(note))
// 	}
// 	return list
// }
//
// //--
// // Error response payloads & renderers
// //--
//
// // ErrResponse renderer type for handling all sorts of errors.
// //
// // In the best case scenario, the excellent github.com/pkg/errors package
// // helps reveal information on the error, setting it on Err, and in the Render()
// // method, using it to set the application-specific error code in AppCode.
// type ErrResponse struct {
// 	Err            error `json:"-"` // low-level runtime error
// 	HTTPStatusCode int   `json:"-"` // http response status code
//
// 	StatusText string `json:"status"`          // user-level status message
// 	AppCode    int64  `json:"code,omitempty"`  // application-specific error code
// 	ErrorText  string `json:"error,omitempty"` // application-level error message, for debugging
// }
//
// func (e *ErrResponse) Render(w http.ResponseWriter, r *http.Request) error {
// 	render.Status(r, e.HTTPStatusCode)
// 	return nil
// }
//
// func ErrInvalidRequest(err error) render.Renderer {
// 	return &ErrResponse{
// 		Err:            err,
// 		HTTPStatusCode: 400,
// 		StatusText:     "Invalid request.",
// 		ErrorText:      err.Error(),
// 	}
// }
//
// func ErrRender(err error) render.Renderer {
// 	return &ErrResponse{
// 		Err:            err,
// 		HTTPStatusCode: 422,
// 		StatusText:     "Error rendering response.",
// 		ErrorText:      err.Error(),
// 	}
// }
//
// var ErrNotFound = &ErrResponse{HTTPStatusCode: 404, StatusText: "Resource not found."}
//
// //--
// // Data model objects and persistence mocks:
// //--
//
// type User struct {
// 	ID   int64  `json:"id"`
// 	Name string `json:"name"`
// }
//
// // Note data model. I suggest looking at https://upper.io for an easy
// // and powerful data persistence adapter.
// type Note struct {
// 	ID      string `json:"id"`
// 	UserID  int64  `json:"user_id"` // the author
// 	Title   string `json:"title"`
// 	Content string `json:"content"`
// }
//
// // Note fixture data
// var notes = []*Note{
// 	{ID: "1", UserID: 100, Title: "Go", Content: "Programming Language"},
// 	{ID: "2", UserID: 200, Title: "JavaScript", Content: "Programming Language"},
// }
//
// // User fixture data
// var users = []*User{
// 	{ID: 100, Name: "Peter"},
// 	{ID: 200, Name: "Julia"},
// }
//
// func dbNewNote(note *Note) (string, error) {
// 	note.ID = fmt.Sprintf("%d", rand.Intn(100)+10)
// 	notes = append(notes, note)
// 	return note.ID, nil
// }
//
// func dbGetNote(id string) (*Note, error) {
// 	for _, note := range notes {
// 		if note.ID == id {
// 			return note, nil
// 		}
// 	}
// 	return nil, errors.New("note not found.")
// }
//
// func dbGetUser(id int64) (*User, error) {
// 	for _, u := range users {
// 		if u.ID == id {
// 			return u, nil
// 		}
// 	}
// 	return nil, errors.New("user not found.")
// }
