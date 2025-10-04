-- public.schedule_seats definition

-- Drop table

-- DROP TABLE public.schedule_seats;

CREATE TABLE public.schedule_seats (
	id int4 GENERATED ALWAYS AS IDENTITY( INCREMENT BY 1 MINVALUE 1 MAXVALUE 2147483647 START 1 CACHE 1 NO CYCLE) NOT NULL,
	schedule_id int4 NOT NULL,
	seat_id int4 NOT NULL,
	status varchar(20) DEFAULT 'available'::character varying NULL,
	CONSTRAINT schedule_seats_pkey PRIMARY KEY (id),
	CONSTRAINT schedule_seats_schedule_id_seat_id_key UNIQUE (schedule_id, seat_id),
	CONSTRAINT schedule_seats_status_check CHECK (((status)::text = ANY ((ARRAY['available'::character varying, 'selected'::character varying, 'love_nest'::character varying, 'sold'::character varying])::text[]))),
	CONSTRAINT schedule_seats_seat_id_fkey FOREIGN KEY (seat_id) REFERENCES public.seats(id) ON DELETE CASCADE
);