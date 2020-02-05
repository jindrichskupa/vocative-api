SET datestyle = "ISO, DMY";

CREATE EXTENSION IF NOT EXISTS "unaccent";
CREATE EXTENSION IF NOT EXISTS "pg_trgm";

CREATE OR REPLACE FUNCTION to_ascii(bytea, name) RETURNS text STRICT AS 'to_ascii_encname' LANGUAGE internal;
CREATE OR REPLACE FUNCTION lower(text) RETURNS text LANGUAGE internal IMMUTABLE STRICT AS $function$lower$function$;
CREATE OR REPLACE FUNCTION lower(anyrange) RETURNS anyelement LANGUAGE internal IMMUTABLE STRICT AS $function$range_lower$function$;

-- vytvoreni tabulky pro adresni mista primo z CSV
CREATE TABLE import_firstnames (
  count VARCHAR(10),
  name VARCHAR(255),
  vocative VARCHAR(255),
  gender VARCHAR(20)
);

-- nahrani dat z CSV
COPY import_firstnames (count,name,vocative) FROM '/data/krestni_muzi.csv' (DELIMITER ',', FORMAT CSV, NULL '', ENCODING 'UTF8');
UPDATE import_firstnames SET gender='male' WHERE gender IS NULL;
COPY import_firstnames (count,name,vocative) FROM '/data/krestni_zeny.csv' (DELIMITER ',', FORMAT CSV, NULL '', ENCODING 'UTF8');
UPDATE import_firstnames SET gender='female' WHERE gender IS NULL;

CREATE TABLE import_surnames (
  count VARCHAR(10),
  name VARCHAR(255),
  vocative VARCHAR(255),
  gender VARCHAR(20)
);

COPY import_surnames (count,name,vocative) FROM '/data/prijmeni_muzi_1.csv' (DELIMITER ',', FORMAT CSV, NULL '', ENCODING 'UTF8');
COPY import_surnames (count,name,vocative) FROM '/data/prijmeni_muzi_2.csv' (DELIMITER ',', FORMAT CSV, NULL '', ENCODING 'UTF8');
UPDATE import_surnames SET gender='male' WHERE gender IS NULL;
COPY import_surnames (count,name,vocative) FROM '/data/prijmeni_zeny_1.csv' (DELIMITER ',', FORMAT CSV, NULL '', ENCODING 'UTF8');
COPY import_surnames (count,name,vocative) FROM '/data/prijmeni_zeny_2.csv' (DELIMITER ',', FORMAT CSV, NULL '', ENCODING 'UTF8');
UPDATE import_surnames SET gender='female' WHERE gender IS NULL;

-- vytvoreni view pro krestni jmena
CREATE MATERIALIZED VIEW first_names AS
  SELECT 
    name as name,
    vocative as vocative,
    regexp_replace(regexp_replace(regexp_replace(lower(to_ascii(convert_to(import_firstnames.name,'latin2'), 'latin2')), '-', ' '),'\.',' ','g'), '\s+', ' ','g') as name_search,
    count::int as count,
    gender
  FROM 
    import_firstnames;

-- vytvoreni view pro krestni jmena
CREATE MATERIALIZED VIEW sur_names AS
  SELECT 
    name as name,
    vocative as vocative,
    regexp_replace(regexp_replace(regexp_replace(lower(to_ascii(convert_to(import_surnames.name,'latin2'), 'latin2')), '-', ' '),'\.',' ','g'), '\s+', ' ','g') as name_search,
    count::int as count,
    gender
  FROM 
    import_surnames;

DROP INDEX IF EXISTS index_first_names_on_gender;
CREATE INDEX index_first_names_on_gender ON first_names USING btree (gender);
DROP INDEX IF EXISTS index_first_names_on_name;
CREATE INDEX index_first_names_on_name ON first_names USING gin (name_search gin_trgm_ops);

DROP INDEX IF EXISTS index_sur_names_on_gender;
CREATE INDEX index_sur_names_on_gender ON sur_names USING btree (gender);
DROP INDEX IF EXISTS index_sur_names_on_name;
CREATE INDEX index_sur_names_on_name ON sur_names USING gin (name_search gin_trgm_ops);

SET datestyle = "ISO, MDY";