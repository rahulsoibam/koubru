--
-- PostgreSQL database dump
--

-- Dumped from database version 9.6.8
-- Dumped by pg_dump version 10.6 (Ubuntu 10.6-0ubuntu0.18.04.1)

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET client_min_messages = warning;
SET row_security = off;

--
-- Name: apgcc; Type: SCHEMA; Schema: -; Owner: rdsadmin
--

CREATE SCHEMA apgcc;


ALTER SCHEMA apgcc OWNER TO rdsadmin;

--
-- Name: plpgsql; Type: EXTENSION; Schema: -; Owner: 
--

CREATE EXTENSION IF NOT EXISTS plpgsql WITH SCHEMA pg_catalog;


--
-- Name: EXTENSION plpgsql; Type: COMMENT; Schema: -; Owner: 
--

COMMENT ON EXTENSION plpgsql IS 'PL/pgSQL procedural language';


--
-- Name: citext; Type: EXTENSION; Schema: -; Owner: 
--

CREATE EXTENSION IF NOT EXISTS citext WITH SCHEMA public;


--
-- Name: EXTENSION citext; Type: COMMENT; Schema: -; Owner: 
--

COMMENT ON EXTENSION citext IS 'data type for case-insensitive character strings';


--
-- Name: hll; Type: EXTENSION; Schema: -; Owner: 
--

CREATE EXTENSION IF NOT EXISTS hll WITH SCHEMA public;


--
-- Name: EXTENSION hll; Type: COMMENT; Schema: -; Owner: 
--

COMMENT ON EXTENSION hll IS 'type for storing hyperloglog data';


--
-- Name: reaction_type; Type: TYPE; Schema: public; Owner: rahulsoibam
--

CREATE TYPE public.reaction_type AS ENUM (
    'happy',
    'angry',
    'sad',
    'wow',
    'love'
);


ALTER TYPE public.reaction_type OWNER TO rahulsoibam;

--
-- Name: id_generator(); Type: FUNCTION; Schema: public; Owner: rahulsoibam
--

CREATE FUNCTION public.id_generator(OUT result bigint) RETURNS bigint
    LANGUAGE plpgsql
    AS $$
DECLARE
    our_epoch bigint := 1314220021721;
    seq_id bigint;
    now_millis bigint;
    -- the id of this DB shard, must be set for each
    -- schema shard you have - you could pass this as a parameter too
    shard_id int := 1;
BEGIN
	    SELECT nextval('global_id_sequence') % 1024 INTO seq_id;

	    SELECT FLOOR(EXTRACT(EPOCH FROM clock_timestamp()) * 1000) INTO now_millis;
	    result := (now_millis - our_epoch) << 23;
	    result := result | (shard_id << 10);
	    result := result | (seq_id);
END;
$$;


ALTER FUNCTION public.id_generator(OUT result bigint) OWNER TO rahulsoibam;

SET default_tablespace = '';

SET default_with_oids = false;

--
-- Name: category; Type: TABLE; Schema: public; Owner: rahulsoibam
--

CREATE TABLE public.category (
    category_id bigint DEFAULT public.id_generator() NOT NULL,
    name public.citext NOT NULL,
    created_by bigint NOT NULL,
    created_on timestamp with time zone DEFAULT timezone('utc'::text, now())
);


ALTER TABLE public.category OWNER TO rahulsoibam;

--
-- Name: category_follower; Type: TABLE; Schema: public; Owner: rahulsoibam
--

CREATE TABLE public.category_follower (
    category_id bigint NOT NULL,
    user_id bigint NOT NULL,
    followed_on timestamp with time zone DEFAULT timezone('utc'::text, now())
);


ALTER TABLE public.category_follower OWNER TO rahulsoibam;

--
-- Name: global_id_sequence; Type: SEQUENCE; Schema: public; Owner: rahulsoibam
--

CREATE SEQUENCE public.global_id_sequence
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.global_id_sequence OWNER TO rahulsoibam;

--
-- Name: kuser; Type: TABLE; Schema: public; Owner: rahulsoibam
--

CREATE TABLE public.kuser (
    user_id bigint DEFAULT public.id_generator() NOT NULL,
    username public.citext NOT NULL,
    full_name text NOT NULL,
    bio text DEFAULT ''::text,
    email text,
    phone text,
    email_verified boolean DEFAULT false NOT NULL,
    phone_verified boolean DEFAULT false NOT NULL,
    photo_url text DEFAULT ''::text NOT NULL,
    facebook text,
    google text,
    created_on timestamp with time zone DEFAULT timezone('utc'::text, now())
);


ALTER TABLE public.kuser OWNER TO rahulsoibam;

--
-- Name: opinion; Type: TABLE; Schema: public; Owner: rahulsoibam
--

CREATE TABLE public.opinion (
    opinion_id bigint DEFAULT public.id_generator() NOT NULL,
    parent_id bigint,
    created_by bigint NOT NULL,
    topic_id bigint NOT NULL,
    is_anonymous boolean DEFAULT false NOT NULL,
    thumb_url text DEFAULT ''::text NOT NULL,
    hls_url text DEFAULT ''::text NOT NULL,
    dash_url text DEFAULT ''::text NOT NULL,
    created_on timestamp with time zone DEFAULT timezone('utc'::text, now()),
    reaction public.reaction_type NOT NULL
);


ALTER TABLE public.opinion OWNER TO rahulsoibam;

--
-- Name: opinion_follower; Type: TABLE; Schema: public; Owner: rahulsoibam
--

CREATE TABLE public.opinion_follower (
    opinion_id bigint NOT NULL,
    user_id bigint NOT NULL,
    followed_on timestamp with time zone DEFAULT timezone('utc'::text, now())
);


ALTER TABLE public.opinion_follower OWNER TO rahulsoibam;

--
-- Name: opinion_vote; Type: TABLE; Schema: public; Owner: rahulsoibam
--

CREATE TABLE public.opinion_vote (
    opinion_id bigint NOT NULL,
    user_id bigint NOT NULL,
    vote boolean,
    voted_on timestamp with time zone DEFAULT timezone('utc'::text, now())
);


ALTER TABLE public.opinion_vote OWNER TO rahulsoibam;

--
-- Name: topic; Type: TABLE; Schema: public; Owner: rahulsoibam
--

CREATE TABLE public.topic (
    topic_id bigint DEFAULT public.id_generator() NOT NULL,
    created_by bigint NOT NULL,
    title text NOT NULL,
    details text DEFAULT ''::text NOT NULL,
    created_on timestamp with time zone DEFAULT timezone('utc'::text, now())
);


ALTER TABLE public.topic OWNER TO rahulsoibam;

--
-- Name: topic_category; Type: TABLE; Schema: public; Owner: rahulsoibam
--

CREATE TABLE public.topic_category (
    topic_id bigint NOT NULL,
    category_id bigint NOT NULL,
    created_on timestamp with time zone DEFAULT timezone('utc'::text, now())
);


ALTER TABLE public.topic_category OWNER TO rahulsoibam;

--
-- Name: topic_follower; Type: TABLE; Schema: public; Owner: rahulsoibam
--

CREATE TABLE public.topic_follower (
    topic_id bigint NOT NULL,
    followed_by bigint NOT NULL,
    followed_on timestamp with time zone DEFAULT timezone('utc'::text, now())
);


ALTER TABLE public.topic_follower OWNER TO rahulsoibam;

--
-- Name: usermap; Type: TABLE; Schema: public; Owner: rahulsoibam
--

CREATE TABLE public.usermap (
    user_id bigint NOT NULL,
    follower_id bigint NOT NULL,
    followed_on timestamp with time zone DEFAULT timezone('utc'::text, now())
);


ALTER TABLE public.usermap OWNER TO rahulsoibam;

--
-- Data for Name: category; Type: TABLE DATA; Schema: public; Owner: rahulsoibam
--

COPY public.category (category_id, name, created_by, created_on) FROM stdin;
1968424375409443862	Politics	1967600534613394434	2019-01-30 18:54:53.047921+00
1968424375409443863	Science	1967600534613394434	2019-01-30 18:54:53.047921+00
1968424375417832472	Famine	1967600534613394434	2019-01-30 18:54:53.047921+00
1968424375417832473	India	1967600534613394434	2019-01-30 18:54:53.047921+00
1968469569395754011	Facebook	1967600534613394434	2019-01-30 20:24:40.590181+00
1968469911558685725	google	1967600534613394434	2019-01-30 20:25:21.379915+00
1968471428244177951	internet	1967600534613394434	2019-01-30 20:28:22.18322+00
\.


--
-- Data for Name: category_follower; Type: TABLE DATA; Schema: public; Owner: rahulsoibam
--

COPY public.category_follower (category_id, user_id, followed_on) FROM stdin;
1968424375409443863	1967600534613394434	2019-01-31 22:13:22.120287+00
1968424375409443862	1967600534613394434	2019-02-02 00:12:15.74799+00
1968424375417832472	1967600534613394434	2019-02-02 00:12:38.371699+00
1968469569395754011	1967601543494501379	2019-02-02 00:13:43.691412+00
1968424375417832472	1967646297590596612	2019-02-02 00:13:43.691412+00
1968424375417832473	1967600534613394434	2019-02-02 00:13:43.691412+00
1968424375417832472	1967601543494501379	2019-02-02 13:24:21.674695+00
1968424375417832472	1967816205448250383	2019-02-02 13:24:21.674695+00
1968471428244177951	1969131975155385386	2019-02-02 14:27:05.943741+00
\.


--
-- Data for Name: kuser; Type: TABLE DATA; Schema: public; Owner: rahulsoibam
--

COPY public.kuser (user_id, username, full_name, bio, email, phone, email_verified, phone_verified, photo_url, facebook, google, created_on) FROM stdin;
1967601543494501379	rahulsoibam2	Rahul Soibam		rahulsoibam2@gmail.com	\N	f	f	https://lh5.googleusercontent.com/-pd3dU_GSnLs/AAAAAAAAAAI/AAAAAAAAAAc/benV9ng5i4U/s96-c/photo.jpg	\N	\N	2019-01-29 15:40:03.83993+00
1967600534613394434	rahulsoibam	Rahul Soibam		rahulsoibam@gmail.com	\N	f	f	https://lh5.googleusercontent.com/-pd3dU_GSnLs/AAAAAAAAAAI/AAAAAAAAAAc/benV9ng5i4U/s96-c/photo.jpg	\N	114532928067913638606	2019-01-29 15:38:03.570338+00
1967646297590596612	nengkhoiba	Nengkhoiba Chungkham		nengkhoiba@mobimp.com	\N	f	f	https://lh5.googleusercontent.com/-pd3dU_GSnLs/AAAAAAAAAAI/AAAAAAAAAAc/benV9ng5i4U/s96-c/photo.jpg	\N	\N	2019-01-29 17:08:58.944716+00
1967648996113515525	nenenendnne	nnenenenennen		neneenen@gmail.com	\N	f	f	https://lh5.googleusercontent.com/-pd3dU_GSnLs/AAAAAAAAAAI/AAAAAAAAAAc/benV9ng5i4U/s96-c/photo.jpg	\N	\N	2019-01-29 17:14:20.633075+00
1967651834281591814	ndndndndn	nsmsnns		ndndndn@hotmail.com	\N	f	f	https://lh5.googleusercontent.com/-pd3dU_GSnLs/AAAAAAAAAAI/AAAAAAAAAAc/benV9ng5i4U/s96-c/photo.jpg	\N	\N	2019-01-29 17:19:58.969156+00
1967654962175411207	kvkvkv	kvvk		kvkvk@gmail.com	\N	f	f	https://lh5.googleusercontent.com/-pd3dU_GSnLs/AAAAAAAAAAI/AAAAAAAAAAc/benV9ng5i4U/s96-c/photo.jpg	\N	\N	2019-01-29 17:26:11.843413+00
1967657599654429704	dhane	Dhaneshori		hijam.dhane@gmail.com	\N	f	f	https://lh5.googleusercontent.com/-pd3dU_GSnLs/AAAAAAAAAAI/AAAAAAAAAAc/benV9ng5i4U/s96-c/photo.jpg	\N	\N	2019-01-29 17:31:26.255667+00
1967811162141623310	lohenyumnam	Lohen Yumnam		\N	\N	f	f	https://lh5.googleusercontent.com/-pd3dU_GSnLs/AAAAAAAAAAI/AAAAAAAAAAc/benV9ng5i4U/s96-c/photo.jpg	10212773160579104	\N	2019-01-29 22:36:32.329202+00
1967816205448250383	undyinglegend	Undying Legend		\N	\N	f	f	https://lh5.googleusercontent.com/-pd3dU_GSnLs/AAAAAAAAAAI/AAAAAAAAAAc/benV9ng5i4U/s96-c/photo.jpg	\N	\N	2019-01-29 22:46:33.539396+00
1967826376593507344	rahulsoibam3	Rahul Soibam		\N	\N	f	f	https://lh5.googleusercontent.com/-pd3dU_GSnLs/AAAAAAAAAAI/AAAAAAAAAAc/benV9ng5i4U/s96-c/photo.jpg	\N	101597058952368758781	2019-01-29 23:06:46.033+00
1968006000665429009	nengkhoibachungkham	nengkhoiba chungkham		\N	\N	f	f	https://lh5.googleusercontent.com/-pd3dU_GSnLs/AAAAAAAAAAI/AAAAAAAAAAc/benV9ng5i4U/s96-c/photo.jpg	\N	113291100633847616878	2019-01-30 05:03:38.892503+00
1968011192307811346	nengkhoibachungkham2	Nengkhoiba Chungkham		\N	\N	f	f	https://lh5.googleusercontent.com/-pd3dU_GSnLs/AAAAAAAAAAI/AAAAAAAAAAc/benV9ng5i4U/s96-c/photo.jpg	2309205209113306	\N	2019-01-30 05:13:57.784807+00
1968030117368169491	nengkhoibach	nengkhoiba Chungkham		neng.ch@gmail.com	\N	f	f	https://lh5.googleusercontent.com/-pd3dU_GSnLs/AAAAAAAAAAI/AAAAAAAAAAc/benV9ng5i4U/s96-c/photo.jpg	\N	\N	2019-01-30 05:51:33.827134+00
1968102053884462100	142pm	Babu Musai		142pm@a.com	\N	f	f	https://lh5.googleusercontent.com/-pd3dU_GSnLs/AAAAAAAAAAI/AAAAAAAAAAc/benV9ng5i4U/s96-c/photo.jpg	\N	\N	2019-01-30 08:14:29.32769+00
1968264739200959509	lohenyumnam2	Lohen Yumnam		\N	\N	f	f	https://lh5.googleusercontent.com/-pd3dU_GSnLs/AAAAAAAAAAI/AAAAAAAAAAc/benV9ng5i4U/s96-c/photo.jpg	\N	114864775164927326365	2019-01-30 13:37:42.92782+00
1968988760016159778	torobabu	Toro		t@a.com	\N	f	f		\N	\N	2019-01-31 13:36:12.936727+00
1968993577316910115	torobi	Toro		tt@a.com	\N	f	f		\N	\N	2019-01-31 13:45:47.203767+00
1969009052411233320	pupu	pupu		pupu@a.com	\N	f	f		\N	\N	2019-01-31 14:16:31.976454+00
1969130633657582633	popo	Popo Tum		popo@a.com	\N	f	f		\N	\N	2019-01-31 18:18:05.590805+00
1969131975155385386	pop	popoBu		pop@a.com	\N	f	f		\N	\N	2019-01-31 18:20:45.511022+00
1969134295687627819	popopo	popo Yum		popopo@a.com	\N	f	f		\N	\N	2019-01-31 18:25:22.140483+00
\.


--
-- Data for Name: opinion; Type: TABLE DATA; Schema: public; Owner: rahulsoibam
--

COPY public.opinion (opinion_id, parent_id, created_by, topic_id, is_anonymous, thumb_url, hls_url, dash_url, created_on, reaction) FROM stdin;
\.


--
-- Data for Name: opinion_follower; Type: TABLE DATA; Schema: public; Owner: rahulsoibam
--

COPY public.opinion_follower (opinion_id, user_id, followed_on) FROM stdin;
\.


--
-- Data for Name: opinion_vote; Type: TABLE DATA; Schema: public; Owner: rahulsoibam
--

COPY public.opinion_vote (opinion_id, user_id, vote, voted_on) FROM stdin;
\.


--
-- Data for Name: topic; Type: TABLE DATA; Schema: public; Owner: rahulsoibam
--

COPY public.topic (topic_id, created_by, title, details, created_on) FROM stdin;
1969001063495238692	1967600534613394434	First example topic		2019-01-31 14:00:39.62561+00
1969001138715886629	1967600534613394434	Second example topic		2019-01-31 14:00:48.591991+00
1969001195338990630	1967600534613394434	Third example topic		2019-01-31 14:00:55.342028+00
1969001300381140007	1967600534613394434	Fourth example topic		2019-01-31 14:01:07.864258+00
1969256566880207917	1967600534613394434	This is created from API		2019-01-31 22:28:18.000012+00
1969257777029186606	1967600534613394434	This is created from API		2019-01-31 22:30:42.261383+00
1969286968462279728	1967600534613394434	This is created from API with Categories		2019-01-31 23:28:42.151447+00
1969288191764595761	1967600534613394434	This is created from API with Categories		2019-01-31 23:31:07.980426+00
1969293378566751282	1967600534613394434	This is from postman		2019-01-31 23:41:26.296922+00
1969295909837603891	1967600534613394434	This is a topic created with categories from Paw		2019-01-31 23:46:28.048499+00
1969457782558032948	1967600534613394434	This is a topic created with categories from Neng		2019-02-01 05:08:04.778396+00
1969721059473097781	1967600534613394434	This is a topic created with categories from Neng		2019-02-01 13:51:09.832169+00
1969843750498731062	1968006000665429009	Testing from Android device		2019-02-01 17:54:55.744406+00
1969924293324178487	1968006000665429009	Testing topic post android		2019-02-01 20:34:57.195639+00
1970258447886713912	1967600534613394434	what is your view towards koubru		2019-02-02 07:38:51.52363+00
1970534525792420921	1968006000665429009	Testing new Android post		2019-02-02 16:47:22.576135+00
\.


--
-- Data for Name: topic_category; Type: TABLE DATA; Schema: public; Owner: rahulsoibam
--

COPY public.topic_category (topic_id, category_id, created_on) FROM stdin;
1969001063495238692	1968424375409443862	2019-01-31 14:03:06.266794+00
1969001063495238692	1968424375409443863	2019-01-31 14:03:06.266794+00
1969001138715886629	1968424375417832473	2019-01-31 14:03:06.266794+00
1969286968462279728	1968424375409443863	2019-01-31 23:28:42.151447+00
1969288191764595761	1968424375409443863	2019-01-31 23:31:07.980426+00
1969293378566751282	1968424375417832473	2019-01-31 23:41:26.296922+00
1969293378566751282	1968424375417832472	2019-01-31 23:41:26.296922+00
1969295909837603891	1968424375409443863	2019-01-31 23:46:28.048499+00
1969295909837603891	1968469911558685725	2019-01-31 23:46:28.048499+00
1969295909837603891	1968471428244177951	2019-01-31 23:46:28.048499+00
1969457782558032948	1968424375409443863	2019-02-01 05:08:04.778396+00
1969457782558032948	1968469911558685725	2019-02-01 05:08:04.778396+00
1969457782558032948	1968471428244177951	2019-02-01 05:08:04.778396+00
1969721059473097781	1968424375409443863	2019-02-01 13:51:09.832169+00
1969721059473097781	1968469911558685725	2019-02-01 13:51:09.832169+00
1969721059473097781	1968471428244177951	2019-02-01 13:51:09.832169+00
1969843750498731062	1968424375409443863	2019-02-01 17:54:55.744406+00
1969843750498731062	1968469911558685725	2019-02-01 17:54:55.744406+00
1969843750498731062	1968424375417832473	2019-02-01 17:54:55.744406+00
1969924293324178487	1968424375409443862	2019-02-01 20:34:57.195639+00
1969924293324178487	1968424375409443863	2019-02-01 20:34:57.195639+00
1969924293324178487	1968424375417832472	2019-02-01 20:34:57.195639+00
1970258447886713912	1968424375409443863	2019-02-02 07:38:51.52363+00
1970258447886713912	1968424375417832473	2019-02-02 07:38:51.52363+00
1970258447886713912	1968469911558685725	2019-02-02 07:38:51.52363+00
1970534525792420921	1968424375409443862	2019-02-02 16:47:22.576135+00
1970534525792420921	1968424375409443863	2019-02-02 16:47:22.576135+00
1970534525792420921	1968424375417832472	2019-02-02 16:47:22.576135+00
\.


--
-- Data for Name: topic_follower; Type: TABLE DATA; Schema: public; Owner: rahulsoibam
--

COPY public.topic_follower (topic_id, followed_by, followed_on) FROM stdin;
1969256566880207917	1967600534613394434	2019-01-31 22:28:18.000012+00
1969293378566751282	1967600534613394434	2019-01-31 23:41:26.296922+00
1969295909837603891	1967600534613394434	2019-01-31 23:46:28.048499+00
1969001300381140007	1967600534613394434	2019-02-01 13:25:26.547001+00
1969001195338990630	1967600534613394434	2019-02-01 13:26:06.126815+00
1969001138715886629	1967600534613394434	2019-02-01 13:26:07.666878+00
1969721059473097781	1968006000665429009	2019-02-01 15:28:04.2111+00
1969721059473097781	1967600534613394434	2019-02-01 15:40:51.20486+00
1969001063495238692	1967600534613394434	2019-02-01 15:41:07.945444+00
1969843750498731062	1968006000665429009	2019-02-01 17:54:55.744406+00
1969457782558032948	1968006000665429009	2019-02-01 18:20:21.29891+00
1969295909837603891	1968006000665429009	2019-02-01 18:20:23.969013+00
1969293378566751282	1968006000665429009	2019-02-01 20:32:22.956769+00
1969288191764595761	1968006000665429009	2019-02-01 20:32:24.796442+00
1969286968462279728	1968006000665429009	2019-02-01 20:32:26.456525+00
1969257777029186606	1968006000665429009	2019-02-01 20:32:27.896487+00
1969256566880207917	1968006000665429009	2019-02-01 20:32:28.726588+00
1969001063495238692	1968006000665429009	2019-02-01 20:32:31.136502+00
1969001138715886629	1968006000665429009	2019-02-01 20:32:32.021559+00
1969001195338990630	1968006000665429009	2019-02-02 02:43:26.399635+00
1969924293324178487	1968006000665429009	2019-02-02 06:18:22.734594+00
1969457782558032948	1967600534613394434	2019-02-02 07:36:37.551618+00
1969924293324178487	1967600534613394434	2019-02-02 07:36:40.952693+00
1969843750498731062	1967600534613394434	2019-02-02 07:36:41.997763+00
1969288191764595761	1967600534613394434	2019-02-02 07:36:46.735102+00
1969286968462279728	1967600534613394434	2019-02-02 07:37:22.961079+00
1970258447886713912	1967600534613394434	2019-02-02 07:38:51.52363+00
1970534525792420921	1968006000665429009	2019-02-02 16:47:22.576135+00
1969001138715886629	1968102053884462100	2019-02-02 20:55:13.098824+00
1969257777029186606	1967600534613394434	2019-02-04 07:29:31.275704+00
1970258447886713912	1968264739200959509	2019-02-04 08:50:10.667522+00
1970534525792420921	1968264739200959509	2019-02-04 09:49:22.233819+00
\.


--
-- Data for Name: usermap; Type: TABLE DATA; Schema: public; Owner: rahulsoibam
--

COPY public.usermap (user_id, follower_id, followed_on) FROM stdin;
1967600534613394434	1967816205448250383	2019-01-30 07:34:33.467922+00
1967600534613394434	1968011192307811346	2019-01-30 07:35:04.03225+00
1967600534613394434	1967657599654429704	2019-01-30 07:35:14.613542+00
1967657599654429704	1967600534613394434	2019-01-30 07:36:05.706672+00
1967601543494501379	1967600534613394434	2019-01-30 07:36:49.435108+00
1968011192307811346	1967600534613394434	2019-02-01 14:37:48.663945+00
1967816205448250383	1967600534613394434	2019-02-01 14:37:49.60878+00
1967601543494501379	1968006000665429009	2019-02-01 19:04:36.663392+00
1967657599654429704	1969130633657582633	2019-02-01 19:04:36.663392+00
1967657599654429704	1967811162141623310	2019-02-01 19:04:36.663392+00
1967600534613394434	1968006000665429009	2019-02-01 19:04:36.663392+00
1967600534613394434	1968264739200959509	2019-02-01 19:04:36.663392+00
1967811162141623310	1967600534613394434	2019-02-01 19:04:36.663392+00
1968102053884462100	1967600534613394434	2019-02-01 19:04:36.663392+00
1967826376593507344	1967600534613394434	2019-02-01 19:04:36.663392+00
1967816205448250383	1967657599654429704	2019-02-01 20:01:15.869606+00
1968264739200959509	1967657599654429704	2019-02-01 20:01:15.869606+00
1968006000665429009	1967657599654429704	2019-02-01 20:01:15.869606+00
1967657599654429704	1968006000665429009	2019-02-01 20:33:05.140987+00
1967826376593507344	1968006000665429009	2019-02-01 20:33:30.239561+00
1968102053884462100	1968006000665429009	2019-02-01 20:33:31.901241+00
1967811162141623310	1968006000665429009	2019-02-01 20:33:33.455297+00
1967816205448250383	1968006000665429009	2019-02-01 20:33:38.484028+00
1968011192307811346	1968006000665429009	2019-02-02 06:16:04.443482+00
1968264739200959509	1967600534613394434	2019-02-02 07:38:00.922207+00
1968006000665429009	1967600534613394434	2019-02-02 07:38:01.982938+00
\.


--
-- Name: global_id_sequence; Type: SEQUENCE SET; Schema: public; Owner: rahulsoibam
--

SELECT pg_catalog.setval('public.global_id_sequence', 57, true);


--
-- Name: category_follower category_follower_pkey; Type: CONSTRAINT; Schema: public; Owner: rahulsoibam
--

ALTER TABLE ONLY public.category_follower
    ADD CONSTRAINT category_follower_pkey PRIMARY KEY (category_id, user_id);


--
-- Name: category category_name_key; Type: CONSTRAINT; Schema: public; Owner: rahulsoibam
--

ALTER TABLE ONLY public.category
    ADD CONSTRAINT category_name_key UNIQUE (name);


--
-- Name: category category_pkey; Type: CONSTRAINT; Schema: public; Owner: rahulsoibam
--

ALTER TABLE ONLY public.category
    ADD CONSTRAINT category_pkey PRIMARY KEY (category_id);


--
-- Name: kuser kuser_email_key; Type: CONSTRAINT; Schema: public; Owner: rahulsoibam
--

ALTER TABLE ONLY public.kuser
    ADD CONSTRAINT kuser_email_key UNIQUE (email);


--
-- Name: kuser kuser_facebook_key; Type: CONSTRAINT; Schema: public; Owner: rahulsoibam
--

ALTER TABLE ONLY public.kuser
    ADD CONSTRAINT kuser_facebook_key UNIQUE (facebook);


--
-- Name: kuser kuser_google_key; Type: CONSTRAINT; Schema: public; Owner: rahulsoibam
--

ALTER TABLE ONLY public.kuser
    ADD CONSTRAINT kuser_google_key UNIQUE (google);


--
-- Name: kuser kuser_phone_key; Type: CONSTRAINT; Schema: public; Owner: rahulsoibam
--

ALTER TABLE ONLY public.kuser
    ADD CONSTRAINT kuser_phone_key UNIQUE (phone);


--
-- Name: kuser kuser_pkey; Type: CONSTRAINT; Schema: public; Owner: rahulsoibam
--

ALTER TABLE ONLY public.kuser
    ADD CONSTRAINT kuser_pkey PRIMARY KEY (user_id);


--
-- Name: kuser kuser_username_key; Type: CONSTRAINT; Schema: public; Owner: rahulsoibam
--

ALTER TABLE ONLY public.kuser
    ADD CONSTRAINT kuser_username_key UNIQUE (username);


--
-- Name: opinion_follower opinion_follower_pkey; Type: CONSTRAINT; Schema: public; Owner: rahulsoibam
--

ALTER TABLE ONLY public.opinion_follower
    ADD CONSTRAINT opinion_follower_pkey PRIMARY KEY (opinion_id, user_id);


--
-- Name: opinion opinion_pkey; Type: CONSTRAINT; Schema: public; Owner: rahulsoibam
--

ALTER TABLE ONLY public.opinion
    ADD CONSTRAINT opinion_pkey PRIMARY KEY (opinion_id);


--
-- Name: opinion_vote opinion_vote_pkey; Type: CONSTRAINT; Schema: public; Owner: rahulsoibam
--

ALTER TABLE ONLY public.opinion_vote
    ADD CONSTRAINT opinion_vote_pkey PRIMARY KEY (opinion_id, user_id);


--
-- Name: topic_category topic_category_pkey; Type: CONSTRAINT; Schema: public; Owner: rahulsoibam
--

ALTER TABLE ONLY public.topic_category
    ADD CONSTRAINT topic_category_pkey PRIMARY KEY (topic_id, category_id);


--
-- Name: topic_follower topic_follower_pkey; Type: CONSTRAINT; Schema: public; Owner: rahulsoibam
--

ALTER TABLE ONLY public.topic_follower
    ADD CONSTRAINT topic_follower_pkey PRIMARY KEY (topic_id, followed_by);


--
-- Name: topic topic_pkey; Type: CONSTRAINT; Schema: public; Owner: rahulsoibam
--

ALTER TABLE ONLY public.topic
    ADD CONSTRAINT topic_pkey PRIMARY KEY (topic_id);


--
-- Name: usermap usermap_pkey; Type: CONSTRAINT; Schema: public; Owner: rahulsoibam
--

ALTER TABLE ONLY public.usermap
    ADD CONSTRAINT usermap_pkey PRIMARY KEY (user_id, follower_id);


--
-- Name: category category_created_by_fkey; Type: FK CONSTRAINT; Schema: public; Owner: rahulsoibam
--

ALTER TABLE ONLY public.category
    ADD CONSTRAINT category_created_by_fkey FOREIGN KEY (created_by) REFERENCES public.kuser(user_id);


--
-- Name: category_follower category_follower_category_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: rahulsoibam
--

ALTER TABLE ONLY public.category_follower
    ADD CONSTRAINT category_follower_category_id_fkey FOREIGN KEY (category_id) REFERENCES public.category(category_id);


--
-- Name: category_follower category_follower_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: rahulsoibam
--

ALTER TABLE ONLY public.category_follower
    ADD CONSTRAINT category_follower_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.kuser(user_id);


--
-- Name: opinion opinion_created_by_fkey; Type: FK CONSTRAINT; Schema: public; Owner: rahulsoibam
--

ALTER TABLE ONLY public.opinion
    ADD CONSTRAINT opinion_created_by_fkey FOREIGN KEY (created_by) REFERENCES public.kuser(user_id);


--
-- Name: opinion_follower opinion_follower_opinion_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: rahulsoibam
--

ALTER TABLE ONLY public.opinion_follower
    ADD CONSTRAINT opinion_follower_opinion_id_fkey FOREIGN KEY (opinion_id) REFERENCES public.opinion(opinion_id);


--
-- Name: opinion_follower opinion_follower_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: rahulsoibam
--

ALTER TABLE ONLY public.opinion_follower
    ADD CONSTRAINT opinion_follower_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.kuser(user_id);


--
-- Name: opinion opinion_parent_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: rahulsoibam
--

ALTER TABLE ONLY public.opinion
    ADD CONSTRAINT opinion_parent_id_fkey FOREIGN KEY (parent_id) REFERENCES public.opinion(opinion_id);


--
-- Name: opinion opinion_topic_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: rahulsoibam
--

ALTER TABLE ONLY public.opinion
    ADD CONSTRAINT opinion_topic_id_fkey FOREIGN KEY (topic_id) REFERENCES public.topic(topic_id);


--
-- Name: opinion_vote opinion_vote_opinion_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: rahulsoibam
--

ALTER TABLE ONLY public.opinion_vote
    ADD CONSTRAINT opinion_vote_opinion_id_fkey FOREIGN KEY (opinion_id) REFERENCES public.opinion(opinion_id);


--
-- Name: opinion_vote opinion_vote_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: rahulsoibam
--

ALTER TABLE ONLY public.opinion_vote
    ADD CONSTRAINT opinion_vote_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.kuser(user_id);


--
-- Name: topic_category topic_category_category_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: rahulsoibam
--

ALTER TABLE ONLY public.topic_category
    ADD CONSTRAINT topic_category_category_id_fkey FOREIGN KEY (category_id) REFERENCES public.category(category_id);


--
-- Name: topic_category topic_category_topic_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: rahulsoibam
--

ALTER TABLE ONLY public.topic_category
    ADD CONSTRAINT topic_category_topic_id_fkey FOREIGN KEY (topic_id) REFERENCES public.topic(topic_id);


--
-- Name: topic topic_created_by_fkey; Type: FK CONSTRAINT; Schema: public; Owner: rahulsoibam
--

ALTER TABLE ONLY public.topic
    ADD CONSTRAINT topic_created_by_fkey FOREIGN KEY (created_by) REFERENCES public.kuser(user_id);


--
-- Name: topic_follower topic_follower_followed_by_fkey; Type: FK CONSTRAINT; Schema: public; Owner: rahulsoibam
--

ALTER TABLE ONLY public.topic_follower
    ADD CONSTRAINT topic_follower_followed_by_fkey FOREIGN KEY (followed_by) REFERENCES public.kuser(user_id);


--
-- Name: topic_follower topic_follower_topic_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: rahulsoibam
--

ALTER TABLE ONLY public.topic_follower
    ADD CONSTRAINT topic_follower_topic_id_fkey FOREIGN KEY (topic_id) REFERENCES public.topic(topic_id);


--
-- Name: usermap usermap_follower_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: rahulsoibam
--

ALTER TABLE ONLY public.usermap
    ADD CONSTRAINT usermap_follower_id_fkey FOREIGN KEY (follower_id) REFERENCES public.kuser(user_id);


--
-- Name: usermap usermap_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: rahulsoibam
--

ALTER TABLE ONLY public.usermap
    ADD CONSTRAINT usermap_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.kuser(user_id);


--
-- Name: SCHEMA public; Type: ACL; Schema: -; Owner: rahulsoibam
--

REVOKE ALL ON SCHEMA public FROM rdsadmin;
REVOKE ALL ON SCHEMA public FROM PUBLIC;
GRANT ALL ON SCHEMA public TO rahulsoibam;
GRANT ALL ON SCHEMA public TO PUBLIC;


--
-- PostgreSQL database dump complete
--

