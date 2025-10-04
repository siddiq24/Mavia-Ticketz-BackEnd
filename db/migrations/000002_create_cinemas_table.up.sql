


-- public.cinemas definition

-- Drop table

-- DROP TABLE public.cinemas;

CREATE TABLE public.cinemas (
	id int4 GENERATED ALWAYS AS IDENTITY( INCREMENT BY 1 MINVALUE 1 MAXVALUE 2147483647 START 1 CACHE 1 NO CYCLE) NOT NULL,
	"name" varchar(50) NOT NULL,
	image varchar NULL,
	CONSTRAINT cinemas_pkey PRIMARY KEY (id)
);
CREATE INDEX idx_cinemas_name ON public.cinemas USING btree (name);