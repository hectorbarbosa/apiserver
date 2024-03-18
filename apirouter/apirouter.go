package apirouter

import (
	"database/sql"
	"encoding/json"
	"filmoteka/database"
	"filmoteka/model"
	"fmt"
	"log"
	"mime"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	_ "github.com/lib/pq"
)

const (
	host      = "localhost"
	port      = 5432
	adminUser = "adm"
	userUser  = "u"
	password  = "pass"
	dbname    = "filmoteka"
)

type ApiRouter struct {
	adminConn *sql.DB
	userConn  *sql.DB
}

func NewApiRouter() *ApiRouter {
	return &ApiRouter{}
}

func (s *ApiRouter) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var ActorRe = regexp.MustCompile(`^/actors/*$`)
	var ActorReWithID = regexp.MustCompile(`^/actors/([0-9]+(?:-[0-9]+)+)$`)
	var FilmRe = regexp.MustCompile(`^/films/*$`)
	var FilmReWithID = regexp.MustCompile(`^/films/([0-9]+(?:-[0-9]+)+)$`)

	switch {
	case r.Method == http.MethodPost && ActorRe.MatchString(r.URL.Path):
		s.CreateActor(w, r)
	case r.Method == http.MethodGet && ActorRe.MatchString(r.URL.Path):
		s.GetActorsList(w, r)
	case r.Method == http.MethodGet && ActorReWithID.MatchString(r.URL.Path):
		s.GetActorById(w, r)
	case r.Method == http.MethodPut && ActorReWithID.MatchString(r.URL.Path):
		s.UpdateActor(w, r)
	case r.Method == http.MethodDelete && ActorReWithID.MatchString(r.URL.Path):
		s.DeleteActor(w, r)

	case r.Method == http.MethodPost && FilmRe.MatchString(r.URL.Path):
		s.CreateFilm(w, r)
	case r.Method == http.MethodGet && FilmRe.MatchString(r.URL.Path):
		s.GetFilmsList(w, r)
	case r.Method == http.MethodGet && FilmReWithID.MatchString(r.URL.Path):
		s.GetFilmById(w, r)
	// case r.Method == http.MethodPut && FilmReWithID.MatchString(r.URL.Path):
	// s.UpdateFilm(w, r)
	case r.Method == http.MethodDelete && FilmReWithID.MatchString(r.URL.Path):
		s.DeleteFilm(w, r)
	}

}

func (s *ApiRouter) Start() error {
	var err error
	psqlAdmin := fmt.Sprintf("host=%s port=%d user=%s password=%s "+
		"dbname=%s sslmode=disable", host, port, adminUser, password, dbname)
	s.adminConn, err = sql.Open("postgres", psqlAdmin)
	if err != nil {
		return err
	}
	psqlUser := fmt.Sprintf("host=%s port=%d user=%s password=%s "+
		"dbname=%s sslmode=disable", host, port, userUser, password, dbname)
	s.userConn, err = sql.Open("postgres", psqlUser)
	if err != nil {
		return err
	}

	log.Println("Router started...")
	return nil
}

func (s *ApiRouter) Stop() error {
	s.adminConn.Close()
	s.userConn.Close()

	log.Println("Router stopped.")
	return nil
}

func (s *ApiRouter) CreateActor(w http.ResponseWriter, req *http.Request) {
	log.Printf("handling actor create at %s\n", req.URL.Path)

	type ResponseId struct {
		Id int `json:"id"`
	}

	// Enforce a JSON Content-Type.
	contentType := req.Header.Get("Content-Type")
	mediatype, _, err := mime.ParseMediaType(contentType)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if mediatype != "application/json" {
		http.Error(w, "expect application/json Content-Type", http.StatusUnsupportedMediaType)
		return
	}

	dec := json.NewDecoder(req.Body)
	dec.DisallowUnknownFields()
	var ra model.RequestActor
	if err := dec.Decode(&ra); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if ra.Name == "" || ra.Gender == "" || ra.BirthDate == "" {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	id, err := database.CreateActor(ra, s.adminConn)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	js, err := json.Marshal(ResponseId{Id: id})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func (s *ApiRouter) GetActorsList(w http.ResponseWriter, req *http.Request) {
	log.Printf("handling get all actors at %s\n", req.URL.Path)
	log.Printf("method: %s\n", req.Method)

	actors, err := database.GetActorsList(s.userConn)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	js, err := json.Marshal(actors)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func (s *ApiRouter) GetActorById(w http.ResponseWriter, req *http.Request) {
	log.Printf("handling get actor by ID at %s\n", req.URL.Path)
	log.Printf("method: %s\n", req.Method)

	id, err := strconv.Atoi(strings.TrimPrefix(req.URL.Path, "/actor/"))
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	actor, err := database.GetActor(id, s.userConn)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	js, err := json.Marshal(actor)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func (s *ApiRouter) UpdateActor(w http.ResponseWriter, req *http.Request) {
	log.Printf("handling update actor at %s\n", req.URL.Path)

	id, err := strconv.Atoi(strings.TrimPrefix(req.URL.Path, "/actor/"))
	if err != nil {
		http.Error(w, "actor update invalid id", http.StatusBadRequest)
		return
	}

	type ResponseId struct {
		Id int `json:"id"`
	}

	// Enforce a JSON Content-Type.
	contentType := req.Header.Get("Content-Type")
	mediatype, _, err := mime.ParseMediaType(contentType)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if mediatype != "application/json" {
		http.Error(w, "expect application/json Content-Type", http.StatusUnsupportedMediaType)
		return
	}

	dec := json.NewDecoder(req.Body)
	dec.DisallowUnknownFields()
	var ra model.RequestActor
	if err := dec.Decode(&ra); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = database.UpdateActor(id, ra, s.adminConn)
	js, err := json.Marshal(ResponseId{Id: id})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func (s *ApiRouter) DeleteActor(w http.ResponseWriter, req *http.Request) {
	log.Printf("handling delete actor at %s\n", req.URL.Path)

	id, err := strconv.Atoi(strings.TrimPrefix(req.URL.Path, "/actor/"))
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	err = database.DeleteActor(id, s.adminConn)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
	}
}

func (s *ApiRouter) CreateFilm(w http.ResponseWriter, req *http.Request) {
	log.Printf("handling film create at %s\n", req.URL.Path)

	type ResponseId struct {
		Id int `json:"id"`
	}

	// Enforce a JSON Content-Type.
	contentType := req.Header.Get("Content-Type")
	mediatype, _, err := mime.ParseMediaType(contentType)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if mediatype != "application/json" {
		http.Error(w, "expect application/json Content-Type", http.StatusUnsupportedMediaType)
		return
	}

	dec := json.NewDecoder(req.Body)
	dec.DisallowUnknownFields()
	var rf model.RequestFilm
	if err := dec.Decode(&rf); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	id, err := database.CreateFilm(rf, s.adminConn)
	js, err := json.Marshal(ResponseId{Id: id})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func (s *ApiRouter) GetFilmsList(w http.ResponseWriter, req *http.Request) {
	log.Printf("handling get all films at %s\n", req.URL.Path)

	films, err := database.GetAllFilms(s.userConn)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	js, err := json.Marshal(films)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func (s *ApiRouter) GetFilmById(w http.ResponseWriter, req *http.Request) {
	log.Printf("handling get film at %s\n", req.URL.Path)

	id, err := strconv.Atoi(strings.TrimPrefix(req.URL.Path, "/film/"))
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	task, err := database.GetFilm(id, s.userConn)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	js, err := json.Marshal(task)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func (s *ApiRouter) DeleteFilm(w http.ResponseWriter, req *http.Request) {
	log.Printf("handling delete film at %s\n", req.URL.Path)

	id, err := strconv.Atoi(strings.TrimPrefix(req.URL.Path, "/film/"))
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	err = database.DeleteFilm(id, s.adminConn)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
	}
}
