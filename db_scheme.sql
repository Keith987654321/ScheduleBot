--
-- PostgreSQL database dump
--

\restrict z5d50fPNx13Ens1suK5wwP9sOZGoEMmF0Q01EQSis2FD3rxF0Z9eKbYuhsAVK9I

-- Dumped from database version 15.14 (Debian 15.14-1.pgdg13+1)
-- Dumped by pg_dump version 15.14 (Debian 15.14-1.pgdg13+1)

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET xmloption = content;
SET client_min_messages = warning;
SET row_security = off;

--
-- Name: db_scheme; Type: SCHEMA; Schema: -; Owner: -
--

CREATE SCHEMA db_scheme;


SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- Name: demo; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.demo (
    cmd_output text
);


--
-- Name: schedule; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.schedule (
    id integer NOT NULL,
    day_of_week integer NOT NULL,
    pair_number integer NOT NULL,
    subject character varying(255) NOT NULL,
    classroom integer DEFAULT 0,
    subgroup integer DEFAULT 0 NOT NULL,
    CONSTRAINT schedule_day_of_week_check CHECK (((day_of_week >= 1) AND (day_of_week <= 7))),
    CONSTRAINT schedule_pair_number_check CHECK (((pair_number >= 1) AND (pair_number <= 8))),
    CONSTRAINT schedule_subgroup_check CHECK (((subgroup >= 0) AND (subgroup <= 2)))
);


--
-- Name: schedule_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.schedule_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: schedule_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.schedule_id_seq OWNED BY public.schedule.id;


--
-- Name: scheduledmessages; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.scheduledmessages (
    id integer NOT NULL,
    cron_spec character varying(24) NOT NULL,
    message character varying(255) NOT NULL,
    subgroup integer DEFAULT 0
);


--
-- Name: scheduledmessages_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.scheduledmessages_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: scheduledmessages_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.scheduledmessages_id_seq OWNED BY public.scheduledmessages.id;


--
-- Name: suggestions; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.suggestions (
    id integer NOT NULL,
    user_id integer,
    day_of_week integer NOT NULL,
    pair_number integer NOT NULL,
    new_subject character varying(255) NOT NULL,
    status character varying(10) DEFAULT 'pending'::character varying,
    classroom integer DEFAULT 0,
    subgroup integer DEFAULT 0 NOT NULL,
    CONSTRAINT suggestions_day_of_week_check CHECK (((day_of_week >= 1) AND (day_of_week <= 7))),
    CONSTRAINT suggestions_pair_number_check CHECK (((pair_number >= 1) AND (pair_number <= 8))),
    CONSTRAINT suggestions_subgroup_check CHECK (((subgroup >= 0) AND (subgroup <= 2)))
);


--
-- Name: suggestions_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.suggestions_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: suggestions_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.suggestions_id_seq OWNED BY public.suggestions.id;


--
-- Name: teachers; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.teachers (
    id integer NOT NULL,
    first_name character varying(64) DEFAULT '-'::character varying,
    middle_name character varying(64) DEFAULT '-'::character varying,
    second_name character varying(64) DEFAULT '-'::character varying,
    subject character varying(64) NOT NULL,
    subgroup integer DEFAULT 0,
    CONSTRAINT teachers_subgroup_check CHECK (((subgroup >= 0) AND (subgroup <= 2)))
);


--
-- Name: teachers_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.teachers_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: teachers_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.teachers_id_seq OWNED BY public.teachers.id;


--
-- Name: users; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.users (
    id integer NOT NULL,
    telegram_id bigint NOT NULL,
    role character varying(10) DEFAULT 'user'::character varying,
    first_name character varying(128) DEFAULT 'default_value'::character varying,
    username character varying(128) DEFAULT 'default_value'::character varying,
    subgroup integer DEFAULT 0,
    CONSTRAINT users_subgroup_check CHECK (((subgroup >= 0) AND (subgroup <= 2)))
);


--
-- Name: users_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.users_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: users_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.users_id_seq OWNED BY public.users.id;


--
-- Name: schedule id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.schedule ALTER COLUMN id SET DEFAULT nextval('public.schedule_id_seq'::regclass);


--
-- Name: scheduledmessages id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.scheduledmessages ALTER COLUMN id SET DEFAULT nextval('public.scheduledmessages_id_seq'::regclass);


--
-- Name: suggestions id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.suggestions ALTER COLUMN id SET DEFAULT nextval('public.suggestions_id_seq'::regclass);


--
-- Name: teachers id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.teachers ALTER COLUMN id SET DEFAULT nextval('public.teachers_id_seq'::regclass);


--
-- Name: users id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.users ALTER COLUMN id SET DEFAULT nextval('public.users_id_seq'::regclass);


--
-- Name: schedule schedule_day_pair_subgroup_key; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.schedule
    ADD CONSTRAINT schedule_day_pair_subgroup_key UNIQUE (day_of_week, pair_number, subgroup);


--
-- Name: schedule schedule_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.schedule
    ADD CONSTRAINT schedule_pkey PRIMARY KEY (id);


--
-- Name: scheduledmessages scheduledmessages_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.scheduledmessages
    ADD CONSTRAINT scheduledmessages_pkey PRIMARY KEY (id);


--
-- Name: suggestions suggestions_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.suggestions
    ADD CONSTRAINT suggestions_pkey PRIMARY KEY (id);


--
-- Name: teachers teachers_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.teachers
    ADD CONSTRAINT teachers_pkey PRIMARY KEY (id);


--
-- Name: users users_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_pkey PRIMARY KEY (id);


--
-- Name: users users_telegram_id_key; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_telegram_id_key UNIQUE (telegram_id);


--
-- Name: suggestions suggestions_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.suggestions
    ADD CONSTRAINT suggestions_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id);


--
-- PostgreSQL database dump complete
--

\unrestrict z5d50fPNx13Ens1suK5wwP9sOZGoEMmF0Q01EQSis2FD3rxF0Z9eKbYuhsAVK9I

