package database

import (
	"database/sql"
	"filmoteka/model"
	"log"
)

func CreateActor(a model.RequestActor, db *sql.DB) (int, error) {
	_, err := db.Exec("insert into f.actors (actor_id, actor_name,"+
		" gender, birth_date) values (DEFAULT, $1, $2, $3);",
		a.Name, a.Gender, a.BirthDate)
	if err != nil {
		log.Fatal("duplicate key value violates unique constraint")
		return 0, err
	}

	// driver is not supported
	// id, err := result.LastInsertId()
	// if err != nil {
	// 	return 0, err
	// }

	return 0, nil
}

func GetActorsList(db *sql.DB) ([]model.Actor, error) {
	rows, err := db.Query("SELECT * FROM f.actors;")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var actors []model.Actor
	for rows.Next() {
		var actor model.Actor
		err := rows.Scan(&actor.Id, &actor.Name, &actor.Gender, &actor.BirthDate)
		if err != nil {
			return actors, err
		}
		actors = append(actors, actor)
	}
	if err := rows.Err(); err != nil {
		return actors, err
	}
	return actors, nil
}

func GetActor(id int, db *sql.DB) (model.Actor, error) {
	log.Printf("Actor id: %d", id)
	var actor model.Actor
	row := db.QueryRow("SELECT * FROM f.actors WHERE actor_id=?;", id)

	err := row.Scan(&actor.Id, &actor.Name, &actor.Gender, &actor.BirthDate)
	if err != nil {
		return actor, err
	}
	return actor, nil
}

func UpdateActor(id int, a model.RequestActor, db *sql.DB) error {
	_, err := db.Exec("UPDATE f.actors "+
		" SET actor_name='?', gender='?', birth_date='?' "+
		"where f.actor_id=?",
		a.Name, a.Gender, a.BirthDate)
	if err != nil {
		return err
	}

	return nil
}

func DeleteActor(id int, db *sql.DB) error {
	_, err := db.Exec("DELETE FROM f.actors WHERE id=?;", id)
	if err != nil {
		return err
	}

	return nil
}

func CreateFilm(f model.RequestFilm, db *sql.DB) (int, error) {
	_, err := db.Exec("insert into f.films (film_id, film_name, description, "+
		"release, rating) values (NULL, $1, $2, $3, $4)",
		f.Name, f.Description, f.Release, f.Rating)
	if err != nil {
		return 0, err
	}
	// id, err := result.LastInsertId()
	// if err != nil {
	// 	return 0, err
	// }

	return 0, nil
}

func GetAllFilms(db *sql.DB) ([]model.Film, error) {
	rows, err := db.Query("SELECT * FROM f.films;")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var films []model.Film
	for rows.Next() {
		var film model.Film
		err := rows.Scan(&film.Id, &film.Name, &film.Description, &film.Release, &film.Rating)
		if err != nil {
			return films, err
		}
		films = append(films, film)
	}
	if err := rows.Err(); err != nil {
		return films, err
	}
	return films, nil
}

func GetFilm(id int, db *sql.DB) (model.Film, error) {
	var film model.Film
	row := db.QueryRow("SELECT * FROM f.films WHERE id = ?;", id)

	err := row.Scan(&film.Id, &film.Name, &film.Description, &film.Rating)
	if err != nil {
		return film, err
	}
	return film, nil
}

func DeleteFilm(id int, db *sql.DB) error {
	_, err := db.Exec("DELETE FROM f.films WHERE id=?;", id)
	if err != nil {
		return err
	}

	return nil
}
