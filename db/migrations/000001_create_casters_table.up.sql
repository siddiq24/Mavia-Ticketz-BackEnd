-- public.casters definition

-- Drop table

-- DROP TABLE public.casters;

CREATE TABLE public.casters (
	id int4 GENERATED ALWAYS AS IDENTITY( INCREMENT BY 1 MINVALUE 1 MAXVALUE 2147483647 START 1 CACHE 1 NO CYCLE) NOT NULL,
	"name" varchar(50) NULL,
	CONSTRAINT casters_pkey PRIMARY KEY (id)
);