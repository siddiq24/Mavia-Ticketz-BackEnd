-- public.movies definition

-- Drop table

-- DROP TABLE public.movies;

CREATE TABLE public.movies (
	id int4 GENERATED ALWAYS AS IDENTITY( INCREMENT BY 1 MINVALUE 1 MAXVALUE 2147483647 START 1 CACHE 1 NO CYCLE) NOT NULL,
	title varchar(100) NOT NULL,
	description text NOT NULL,
	release_date date DEFAULT now() NOT NULL,
	duration int4 NOT NULL,
	poster_img varchar NULL,
	director_id int4 NOT NULL,
	backdrop_img varchar NULL,
	rating float8 DEFAULT 10 NOT NULL,
	created_at timestamp DEFAULT CURRENT_TIMESTAMP NOT NULL,
	updated_at timestamp NULL,
	is_deleted bool DEFAULT false NULL,
	CONSTRAINT movies_duration_check CHECK ((duration > 0)),
	CONSTRAINT movies_pkey PRIMARY KEY (id),
	CONSTRAINT movies_rating_check CHECK (((rating >= (0)::double precision) AND (rating <= (10)::double precision))),
	CONSTRAINT fk_movies_director FOREIGN KEY (director_id) REFERENCES public.directors(id) ON DELETE CASCADE ON UPDATE CASCADE
);
CREATE INDEX idx_movies_title ON public.movies USING btree (title);