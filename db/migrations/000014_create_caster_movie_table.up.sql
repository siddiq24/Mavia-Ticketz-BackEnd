-- public.caster_movie definition

-- Drop table

-- DROP TABLE public.caster_movie;

CREATE TABLE public.caster_movie (
	movie_id int4 NOT NULL,
	caster_id int4 NOT NULL,
	CONSTRAINT caster_movie_pkey PRIMARY KEY (movie_id, caster_id),
	CONSTRAINT fk_castermovie_caster FOREIGN KEY (caster_id) REFERENCES public.casters(id) ON DELETE CASCADE,
	CONSTRAINT fk_castermovie_movie FOREIGN KEY (movie_id) REFERENCES public.movies(id) ON DELETE CASCADE
);