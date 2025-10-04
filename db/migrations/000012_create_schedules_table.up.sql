

-- public.schedules definition

-- Drop table

-- DROP TABLE public.schedules;

CREATE TABLE public.schedules (
	id int4 GENERATED ALWAYS AS IDENTITY( INCREMENT BY 1 MINVALUE 1 MAXVALUE 2147483647 START 1 CACHE 1 NO CYCLE) NOT NULL,
	movie_id int4 NOT NULL,
	cinema_id int4 NOT NULL,
	time_id int4 NOT NULL,
	"date" date NOT NULL,
	city_id int4 NULL,
	price int4 DEFAULT 10 NULL,
	CONSTRAINT schedules_pkey PRIMARY KEY (id),
	CONSTRAINT unique_schedule_rule UNIQUE (movie_id, cinema_id, time_id, date, city_id),
	CONSTRAINT fk_schedule_city FOREIGN KEY (city_id) REFERENCES public.cities(id),
	CONSTRAINT fk_schedules_cinema FOREIGN KEY (cinema_id) REFERENCES public.cinemas(id),
	CONSTRAINT fk_schedules_movie FOREIGN KEY (movie_id) REFERENCES public.movies(id),
	CONSTRAINT fk_schedules_time FOREIGN KEY (time_id) REFERENCES public.times(id)
);