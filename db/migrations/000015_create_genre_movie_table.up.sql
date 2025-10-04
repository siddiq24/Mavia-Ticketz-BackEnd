-- public.genre_movie definition

-- Drop table

-- DROP TABLE public.genre_movie;

CREATE TABLE public.genre_movie (
	movie_id int4 NOT NULL,
	genre_id int4 NOT NULL,
	CONSTRAINT genre_movie_pkey PRIMARY KEY (movie_id, genre_id),
	CONSTRAINT fk_genremovie_genre FOREIGN KEY (genre_id) REFERENCES public.genres(id) ON DELETE CASCADE,
	CONSTRAINT fk_genremovie_movie FOREIGN KEY (movie_id) REFERENCES public.movies(id) ON DELETE CASCADE
);


