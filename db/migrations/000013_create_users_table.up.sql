

-- public.users definition

-- Drop table

-- DROP TABLE public.users;

CREATE TABLE public.users (
	id int4 GENERATED ALWAYS AS IDENTITY( INCREMENT BY 1 MINVALUE 1 MAXVALUE 2147483647 START 1 CACHE 1 NO CYCLE) NOT NULL,
	email varchar NOT NULL,
	"password" varchar NOT NULL,
	point int4 DEFAULT 0 NOT NULL,
	created_at timestamp DEFAULT CURRENT_TIMESTAMP NULL,
	"role" varchar(100) DEFAULT 'user'::character varying NULL,
	username varchar(100) NOT NULL,
	updated_at timestamptz DEFAULT CURRENT_TIMESTAMP NULL,
	city_id int4 NULL,
	CONSTRAINT users_email_key UNIQUE (email),
	CONSTRAINT users_pkey PRIMARY KEY (id),
	CONSTRAINT users_username_key UNIQUE (username),
	CONSTRAINT users_city_id_fkey FOREIGN KEY (city_id) REFERENCES public.cities(id)
);