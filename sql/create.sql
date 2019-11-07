\c testdb;

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE project_data (
    id uuid DEFAULT uuid_generate_v4(), 
    project varchar, 
    data jsonb, 
    meta jsonb,
    PRIMARY KEY (id)
);

CREATE INDEX project_data_idx ON project_data (project);
CREATE INDEX ON project_data USING gin ( to_tsvector('english', data) );

CREATE OR REPLACE FUNCTION create_partition_and_insert() RETURNS trigger AS
  $BODY$
    DECLARE
      partition TEXT;
    BEGIN
      partition := TG_RELNAME || '_' || MD5(NEW.project);
      IF NOT EXISTS(SELECT relname FROM pg_class WHERE relname=partition) THEN
        RAISE NOTICE 'A partition has been created %',partition;
        EXECUTE 'CREATE TABLE ' || partition || ' (check (project = ''' || NEW.project || ''')) INHERITS (' || TG_RELNAME || ');';
      END IF;
      EXECUTE 'INSERT INTO ' || partition || ' SELECT(' || TG_RELNAME || ' ' || quote_literal(NEW) || ').* RETURNING id;';
      RETURN NULL;
    END;
  $BODY$
LANGUAGE plpgsql VOLATILE
COST 100;

CREATE TRIGGER testing_partition_insert_trigger
BEFORE INSERT ON project_data
FOR EACH ROW EXECUTE PROCEDURE create_partition_and_insert();
