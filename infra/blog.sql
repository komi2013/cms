--
-- PostgreSQL database dump
--

-- Dumped from database version 12.9 (Ubuntu 12.9-0ubuntu0.20.04.1)
-- Dumped by pg_dump version 12.9 (Ubuntu 12.9-0ubuntu0.20.04.1)

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

SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- Name: m_category_tree; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.m_category_tree (
    leaf_id integer DEFAULT 0 NOT NULL,
    level_1 integer DEFAULT 0 NOT NULL,
    level_2 integer DEFAULT 0 NOT NULL,
    level_3 integer DEFAULT 0 NOT NULL,
    level_4 integer DEFAULT 0 NOT NULL,
    level_5 integer DEFAULT 0 NOT NULL,
    level_6 integer DEFAULT 0 NOT NULL,
    level_7 integer DEFAULT 0 NOT NULL,
    level_8 integer DEFAULT 0 NOT NULL,
    updated_at timestamp without time zone DEFAULT now() NOT NULL
);


ALTER TABLE public.m_category_tree OWNER TO postgres;

--
-- Name: category_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.category_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.category_id_seq OWNER TO postgres;

--
-- Name: category_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.category_id_seq OWNED BY public.m_category_tree.leaf_id;


--
-- Name: m_category_name; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.m_category_name (
    category_id integer DEFAULT nextval('public.category_id_seq'::regclass) NOT NULL,
    category_name text DEFAULT ''::text NOT NULL,
    updated_at timestamp without time zone DEFAULT now() NOT NULL,
    category_description text DEFAULT ''::text NOT NULL
);


ALTER TABLE public.m_category_name OWNER TO postgres;

--
-- Name: nologin_usr_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.nologin_usr_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1
    CYCLE;


ALTER TABLE public.nologin_usr_id_seq OWNER TO postgres;

--
-- Name: t_comment; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.t_comment (
    comment_id integer NOT NULL,
    comment_txt text DEFAULT ''::text NOT NULL,
    usr_id integer DEFAULT 0 NOT NULL,
    note_id integer DEFAULT 0 NOT NULL,
    created_at timestamp without time zone DEFAULT now() NOT NULL
);


ALTER TABLE public.t_comment OWNER TO postgres;

--
-- Name: t_comment_comment_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.t_comment_comment_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.t_comment_comment_id_seq OWNER TO postgres;

--
-- Name: t_comment_comment_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.t_comment_comment_id_seq OWNED BY public.t_comment.comment_id;


--
-- Name: t_note_note_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.t_note_note_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.t_note_note_id_seq OWNER TO postgres;

--
-- Name: t_note; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.t_note (
    note_id integer DEFAULT nextval('public.t_note_note_id_seq'::regclass) NOT NULL,
    note_txt text DEFAULT ''::text NOT NULL,
    updated_at timestamp without time zone DEFAULT now() NOT NULL,
    category_id integer DEFAULT 0 NOT NULL,
    note_img text DEFAULT ''::text NOT NULL,
    sequence integer DEFAULT 0 NOT NULL,
    note_title text DEFAULT ''::text NOT NULL,
    list_category_id integer DEFAULT 0 NOT NULL
);


ALTER TABLE public.t_note OWNER TO postgres;

--
-- Name: COLUMN t_note.sequence; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.t_note.sequence IS 'within leaf category';


--
-- Name: usr_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.usr_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.usr_id_seq OWNER TO postgres;

--
-- Name: t_usr; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.t_usr (
    usr_id integer DEFAULT nextval('public.usr_id_seq'::regclass) NOT NULL,
    pv_u_id text DEFAULT ''::text NOT NULL,
    provider integer DEFAULT 0 NOT NULL,
    usr_img text DEFAULT ''::text NOT NULL,
    updated_at timestamp without time zone DEFAULT now() NOT NULL,
    introduce text DEFAULT ''::text NOT NULL,
    push_tokens json
);


ALTER TABLE public.t_usr OWNER TO postgres;

--
-- Name: COLUMN t_usr.provider; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.t_usr.provider IS '1=FB, 2=TW, 3=G+, 4=LN';


--
-- Name: t_comment comment_id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.t_comment ALTER COLUMN comment_id SET DEFAULT nextval('public.t_comment_comment_id_seq'::regclass);


--
-- Name: m_category_name m_category_name_category_id; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.m_category_name
    ADD CONSTRAINT m_category_name_category_id PRIMARY KEY (category_id);


--
-- Name: m_category_tree m_category_tree_leaf_id; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.m_category_tree
    ADD CONSTRAINT m_category_tree_leaf_id PRIMARY KEY (leaf_id);


--
-- Name: t_comment t_comment_comment_id; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.t_comment
    ADD CONSTRAINT t_comment_comment_id PRIMARY KEY (comment_id);


--
-- Name: t_note t_note_note_id; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.t_note
    ADD CONSTRAINT t_note_note_id PRIMARY KEY (note_id);


--
-- Name: t_usr t_usr_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.t_usr
    ADD CONSTRAINT t_usr_pkey PRIMARY KEY (usr_id);


--
-- Name: t_comment_note_id; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX t_comment_note_id ON public.t_comment USING btree (note_id);


--
-- Name: SCHEMA public; Type: ACL; Schema: -; Owner: postgres
--

GRANT USAGE ON SCHEMA public TO exam_8099;


--
-- Name: TABLE m_category_tree; Type: ACL; Schema: public; Owner: postgres
--

GRANT ALL ON TABLE public.m_category_tree TO exam_8099;


--
-- Name: SEQUENCE category_id_seq; Type: ACL; Schema: public; Owner: postgres
--

GRANT SELECT,USAGE ON SEQUENCE public.category_id_seq TO exam_8099;


--
-- Name: TABLE m_category_name; Type: ACL; Schema: public; Owner: postgres
--

GRANT ALL ON TABLE public.m_category_name TO exam_8099;


--
-- Name: SEQUENCE nologin_usr_id_seq; Type: ACL; Schema: public; Owner: postgres
--

GRANT SELECT,USAGE ON SEQUENCE public.nologin_usr_id_seq TO exam_8099;


--
-- Name: TABLE t_comment; Type: ACL; Schema: public; Owner: postgres
--

GRANT ALL ON TABLE public.t_comment TO exam_8099;


--
-- Name: SEQUENCE t_comment_comment_id_seq; Type: ACL; Schema: public; Owner: postgres
--

GRANT SELECT,USAGE ON SEQUENCE public.t_comment_comment_id_seq TO exam_8099;


--
-- Name: SEQUENCE t_note_note_id_seq; Type: ACL; Schema: public; Owner: postgres
--

GRANT SELECT,USAGE ON SEQUENCE public.t_note_note_id_seq TO exam_8099;


--
-- Name: TABLE t_note; Type: ACL; Schema: public; Owner: postgres
--

GRANT ALL ON TABLE public.t_note TO exam_8099;


--
-- Name: SEQUENCE usr_id_seq; Type: ACL; Schema: public; Owner: postgres
--

GRANT SELECT,USAGE ON SEQUENCE public.usr_id_seq TO exam_8099;


--
-- Name: TABLE t_usr; Type: ACL; Schema: public; Owner: postgres
--

GRANT ALL ON TABLE public.t_usr TO exam_8099;


--
-- PostgreSQL database dump complete
--

