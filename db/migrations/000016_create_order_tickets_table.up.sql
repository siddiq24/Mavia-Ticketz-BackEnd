-- public.order_tickets definition

-- Drop table

-- DROP TABLE public.order_tickets;

CREATE TABLE public.order_tickets (
	id int4 GENERATED ALWAYS AS IDENTITY( INCREMENT BY 1 MINVALUE 1 MAXVALUE 2147483647 START 1 CACHE 1 NO CYCLE) NOT NULL,
	user_id int4 NOT NULL,
	total_amount int4 NOT NULL,
	fullname varchar NOT NULL,
	email varchar NOT NULL,
	point_ticket int4 NULL,
	is_paid bool DEFAULT false NULL,
	phone varchar(20) NOT NULL,
	payment_method_id int4 NULL,
	is_active bool DEFAULT true NULL,
	create_at timestamp DEFAULT CURRENT_TIMESTAMP NOT NULL,
	schedule_id int4 NULL,
	CONSTRAINT order_tickets_pkey PRIMARY KEY (id),
	CONSTRAINT order_tickets_payment_method_id_fkey FOREIGN KEY (payment_method_id) REFERENCES public.payment_method(id),
	CONSTRAINT order_tickets_schedule_id_fkey FOREIGN KEY (schedule_id) REFERENCES public.schedules(id),
	CONSTRAINT order_tickets_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id)
);