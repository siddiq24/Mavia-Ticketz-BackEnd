-- public.profiles definition

-- Drop table

-- DROP TABLE public.profiles;

CREATE TABLE public.profiles (
	user_id int4 NOT NULL,
	avatar varchar NULL,
	phone varchar(20) NULL,
	address text NULL,
	birthdate date NULL,
	CONSTRAINT profiles_pkey PRIMARY KEY (user_id),
	CONSTRAINT fk_profile_user FOREIGN KEY (user_id) REFERENCES public.users(id) ON DELETE CASCADE
);