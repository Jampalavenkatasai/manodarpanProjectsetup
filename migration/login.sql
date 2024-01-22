CREATE TABLE IF NOT EXISTS public.login
(
    id SERIAL PRIMARY KEY,
    first_name character varying(255) COLLATE pg_catalog."default" NOT NULL,
    last_name character varying(255) COLLATE pg_catalog."default" NOT NULL,
    email character varying(255) COLLATE pg_catalog."default" NOT NULL,
    CONSTRAINT login_email_key UNIQUE (email) -- Corrected the spelling of CONSTRAINT
    )
    WITH (
        OIDS = FALSE
        )
    TABLESPACE pg_default;

ALTER TABLE IF EXISTS public.login
    OWNER TO postgres;
