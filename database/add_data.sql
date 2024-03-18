\c filmoteka;

insert into f.films (film_id, film_name, description, release_date, rating) 
    values(DEFAULT, 'День сурка', 'Description1', '1995-01-06', 9);

insert into f.actors (actor_id, actor_name, gender, birth_date) 
    values(DEFAULT, 'Билл Мюррей', 'M', '1950-09-21');
    
-- insert into f.roles (role_id, actor_id,	film_id) values(DEFAULT, 1, 3);

\c postgres;
