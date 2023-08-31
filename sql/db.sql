--
-- PostgreSQL database dump
--

-- Dumped from database version 15.4 (Debian 15.4-1.pgdg120+1)
-- Dumped by pg_dump version 15.3

-- Started on 2023-08-30 19:41:16

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
-- TOC entry 6 (class 2615 OID 16505)
-- Name: operations; Type: SCHEMA; Schema: -; Owner: -
--

CREATE SCHEMA operations;


--
-- TOC entry 7 (class 2615 OID 16506)
-- Name: segments; Type: SCHEMA; Schema: -; Owner: -
--

CREATE SCHEMA segments;


--
-- TOC entry 224 (class 1255 OID 16507)
-- Name: add_segment(character varying); Type: FUNCTION; Schema: operations; Owner: -
--

CREATE FUNCTION operations.add_segment(pslug character varying) RETURNS bigint
    LANGUAGE plpgsql
    AS $$
declare
    vResult bigint;
    vCount integer;
begin
    insert
    into segments.list
    (slug)
    values (pSlug)
        returning
            id
    into
        vResult;

    return vResult;
end
$$;


--
-- TOC entry 220 (class 1255 OID 16508)
-- Name: add_user_link(bigint, character varying); Type: FUNCTION; Schema: operations; Owner: -
--

CREATE FUNCTION operations.add_user_link(puser_id bigint, pslug character varying) RETURNS bigint
    LANGUAGE plpgsql
    AS $$
declare
    vResult bigint;
    vSegment_id bigint;
    vAction_id bigint;
begin
	select id into vSegment_id
    from segments.list
    where slug = pSlug;

    select id into vResult
    from segments.user_links
    where user_id = puser_id and segment_id = vSegment_id;

    if vResult <> 0 then
        return vResult;
    end if;

	insert
    into segments.user_links
    (user_id, segment_id)
    values (pUser_id, vSegment_id)
        returning
            id
    into
        vResult;

    select a.action_id into vAction_id
    from segments.actions as a
    where action_name = 'add';

    insert
    into segments.history
    (user_id, segment_id, action_id, created_at)
    values (pUser_id, vSegment_id, vAction_id, now());

    return vResult;
end
$$;


--
-- TOC entry 221 (class 1255 OID 16509)
-- Name: delete_segment(character varying); Type: FUNCTION; Schema: operations; Owner: -
--

CREATE FUNCTION operations.delete_segment(pslug character varying) RETURNS smallint
    LANGUAGE plpgsql
    AS $$
declare
	vSegment_id bigint;
    vResult smallint;
begin
	select id into vSegment_id
    from segments.list
    where slug = pSlug;

    delete
    from segments.history
    where segment_id = vSegment_id;

	delete
	from segments.user_links
	where segment_id = vSegment_id;

	with deleted as (
		delete
	    from segments.list
	    where slug = pslug
	    returning *
	)
	SELECT count(*) FROM deleted into vResult;

    return vResult;
end
$$;


--
-- TOC entry 222 (class 1255 OID 16510)
-- Name: delete_user_link(bigint, character varying); Type: FUNCTION; Schema: operations; Owner: -
--

CREATE FUNCTION operations.delete_user_link(puser_id bigint, pslug character varying) RETURNS smallint
    LANGUAGE plpgsql
    AS $$
declare
    vSegment_id bigint;
    vResult smallint;
    vAction_id bigint;
begin
	select id into vSegment_id
    from segments.list
    where slug = pSlug;

	with deleted as (
		delete
	    from segments.user_links
	    where user_id = pUser_id and segment_id  = vSegment_id
	    returning *
	)
	SELECT count(*) FROM deleted into vResult;

    if vResult <> 0 then
        select a.action_id into vAction_id
        from segments.actions as a
        where action_name = 'delete';

        insert
        into segments.history
        (user_id, segment_id, action_id, created_at)
        values (pUser_id, vSegment_id, vAction_id, now());
    end if;

    return vResult;
end
$$;


--
-- TOC entry 225 (class 1255 OID 24576)
-- Name: get_segment(character varying); Type: FUNCTION; Schema: operations; Owner: -
--

CREATE FUNCTION operations.get_segment(pslug character varying) RETURNS TABLE(id bigint, slug character varying)
    LANGUAGE plpgsql
    AS $$
begin
	return query
	select sl.id, sl.slug
    from segments.list as sl
    where sl.slug = pSlug;
end
$$;


--
-- TOC entry 223 (class 1255 OID 16511)
-- Name: get_segments_by_user_id(bigint); Type: FUNCTION; Schema: operations; Owner: -
--

CREATE FUNCTION operations.get_segments_by_user_id(puser_id bigint) RETURNS TABLE(slug character varying)
    LANGUAGE plpgsql
    AS $$
begin
    return query
        select segments.list.slug
        from segments.user_links
             join segments.list on segments.list.id = segment_id
        where user_id = puser_id;
end
$$;


--
-- TOC entry 216 (class 1259 OID 16512)
-- Name: segments_seq; Type: SEQUENCE; Schema: segments; Owner: -
--

CREATE SEQUENCE segments.segments_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- TOC entry 217 (class 1259 OID 16513)
-- Name: list; Type: TABLE; Schema: segments; Owner: -
--

CREATE TABLE segments.list (
    id bigint DEFAULT nextval('segments.segments_seq'::regclass) NOT NULL,
    slug character varying(100) NOT NULL
);


CREATE SEQUENCE segments.actions_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

CREATE TABLE segments.actions
(
    action_id   bigint DEFAULT nextval('segments.actions_seq'::regclass) NOT NULL,
    action_name varchar(50)                                              NOT NULL,

    CONSTRAINT pk_action PRIMARY KEY (action_id)
);

CREATE UNIQUE INDEX idx1_action ON segments.actions USING btree (action_name);

INSERT INTO segments.actions (action_name) VALUES ('add');
INSERT INTO segments.actions (action_name) VALUES ('delete');


--
-- TOC entry 218 (class 1259 OID 16517)
-- Name: user_links_seq; Type: SEQUENCE; Schema: segments; Owner: -
--

CREATE SEQUENCE segments.user_links_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- TOC entry 219 (class 1259 OID 16518)
-- Name: user_links; Type: TABLE; Schema: segments; Owner: -
--

CREATE TABLE segments.user_links (
    id bigint DEFAULT nextval('segments.user_links_seq'::regclass) NOT NULL,
    user_id bigint NOT NULL,
    segment_id bigint NOT NULL
);


--
-- TOC entry 3216 (class 2606 OID 16523)
-- Name: list pk_segment; Type: CONSTRAINT; Schema: segments; Owner: -
--

ALTER TABLE ONLY segments.list
    ADD CONSTRAINT pk_segment PRIMARY KEY (id);


--
-- TOC entry 3220 (class 2606 OID 16525)
-- Name: user_links pk_user; Type: CONSTRAINT; Schema: segments; Owner: -
--

ALTER TABLE ONLY segments.user_links
    ADD CONSTRAINT pk_user PRIMARY KEY (id);


--
-- TOC entry 3214 (class 1259 OID 16526)
-- Name: idx1_segments; Type: INDEX; Schema: segments; Owner: -
--

CREATE UNIQUE INDEX idx1_segments ON segments.list USING btree (slug);


--
-- TOC entry 3217 (class 1259 OID 16527)
-- Name: idx1_users; Type: INDEX; Schema: segments; Owner: -
--

CREATE INDEX idx1_users ON segments.user_links USING btree (user_id);


--
-- TOC entry 3218 (class 1259 OID 16528)
-- Name: idx2_segments; Type: INDEX; Schema: segments; Owner: -
--

CREATE INDEX idx2_segments ON segments.user_links USING btree (segment_id);


--
-- TOC entry 3221 (class 2606 OID 16529)
-- Name: user_links fk1_segments; Type: FK CONSTRAINT; Schema: segments; Owner: -
--

ALTER TABLE ONLY segments.user_links
    ADD CONSTRAINT fk1_segments FOREIGN KEY (segment_id) REFERENCES segments.list(id);


-- Completed on 2023-08-30 19:41:16

--
-- PostgreSQL database dump complete
--

CREATE SEQUENCE segments.history_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

CREATE TABLE segments.history
(
    history_id         bigint DEFAULT nextval('segments.history_seq'::regclass) NOT NULL,
    user_id    bigint                                                      NOT NULL,
    segment_id bigint NOT NULL,
    action_id bigint                                                      NOT NULL,
    created_at timestamp NOT NULL,

    CONSTRAINT pk_history PRIMARY KEY (history_id),
    CONSTRAINT fk1_segments FOREIGN KEY (segment_id) REFERENCES segments.list (id),
    CONSTRAINT fk2_actions FOREIGN KEY (action_id) REFERENCES segments.actions
);

CREATE INDEX idx1_user ON segments.history USING btree (user_id);
CREATE INDEX idx2_segment ON segments.history USING btree (segment_id);


CREATE FUNCTION operations.get_history(pUser_id bigint, pDateFrom timestamp, pDateTo timestamp) RETURNS
    TABLE(user_id bigint, slug character varying, action character varying, date timestamp)
    LANGUAGE plpgsql
AS $$
begin
    return query
        select sh.user_id, sl.slug, sa.action_name, sh.created_at
        from segments.history as sh
                 join segments.actions as sa on sa.action_id = sh.action_id
                 join segments.list as sl on sl.id = segment_id
        where sh.user_id = pUser_id and sh.created_at >= pDateFrom and sh.created_at < pDateTo;
end
$$;