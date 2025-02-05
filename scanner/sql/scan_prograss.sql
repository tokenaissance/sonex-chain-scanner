DROP SEQUENCE IF EXISTS "scan_prograss_id_seq"; 
CREATE SEQUENCE "scan_prograss_id_seq" 
INCREMENT 1
MINVALUE  1
MAXVALUE 9223372036854775807
START 1
CACHE 1;

DROP TABLE IF EXISTS "scan_prograss";
CREATE TABLE "scan_prograss" (
  "id" int8 NOT NULL DEFAULT nextval('scan_prograss_id_seq'::regclass),
  "created_at" timestamptz(6) NOT NULL,
  "updated_at" timestamptz(6) NOT NULL,
  "chain_id" int8 NOT NULL,
  "height" int8 NOT NULL
)
;

ALTER TABLE "scan_prograss" ADD CONSTRAINT "scan_prograss_pkey" PRIMARY KEY ("id");

CREATE UNIQUE INDEX "idx_scan_prograss_chain_id" ON "scan_prograss" USING btree (
  "chain_id" "pg_catalog"."int8_ops" ASC NULLS LAST
);
