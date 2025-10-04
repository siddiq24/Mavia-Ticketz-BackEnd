-- public.times definition

-- Drop table

-- DROP TABLE public.times;

CREATE TABLE public.times (
	id int4 GENERATED ALWAYS AS IDENTITY( INCREMENT BY 1 MINVALUE 1 MAXVALUE 2147483647 START 1 CACHE 1 NO CYCLE) NOT NULL,
	"time" time NOT NULL,
	CONSTRAINT times_pkey PRIMARY KEY (id)
);