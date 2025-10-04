-- public.order_seats definition

-- Drop table

-- DROP TABLE public.order_seats;

CREATE TABLE public.order_seats (
	order_id int4 NOT NULL,
	seat_id int4 NOT NULL,
	CONSTRAINT order_seats_pkey PRIMARY KEY (order_id, seat_id),
	CONSTRAINT unique_rule UNIQUE (order_id, seat_id),
	CONSTRAINT fk_orderseats_orders FOREIGN KEY (order_id) REFERENCES public.order_tickets(id) ON DELETE CASCADE,
	CONSTRAINT fk_orderseats_seats FOREIGN KEY (seat_id) REFERENCES public.seats(id) ON DELETE CASCADE
);