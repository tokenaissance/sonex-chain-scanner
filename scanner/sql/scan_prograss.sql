/*
 Navicat Premium Dump SQL

 Source Server         : postgres
 Source Server Type    : PostgreSQL
 Source Server Version : 160003 (160003)
 Source Host           : 10.0.0.150:5432
 Source Catalog        : chain_scanner
 Source Schema         : public

 Target Server Type    : PostgreSQL
 Target Server Version : 160003 (160003)
 File Encoding         : 65001

 Date: 05/02/2025 19:37:07
*/

-- ----------------------------
-- Sequence structure for scan_prograss_id_seq
-- ----------------------------
DROP SEQUENCE IF EXISTS "scan_prograss_id_seq"; 
CREATE SEQUENCE "scan_prograss_id_seq" 
INCREMENT 1
MINVALUE  1
MAXVALUE 9223372036854775807
START 1
CACHE 1;

-- ----------------------------
-- Table structure for scan_prograss
-- ----------------------------
DROP TABLE IF EXISTS "scan_prograss";
CREATE TABLE "scan_prograss" (
  "id" int8 NOT NULL DEFAULT nextval('scan_prograss_id_seq'::regclass),
  "created_at" timestamptz(6) NOT NULL,
  "updated_at" timestamptz(6) NOT NULL,
  "chain_id" int8 NOT NULL,
  "height" int8 NOT NULL
)
;

-- ----------------------------
-- Indexes structure for table scan_prograss
-- ----------------------------
CREATE UNIQUE INDEX "idx_scan_prograss_chain_id" ON "scan_prograss" USING btree (
  "chain_id" "pg_catalog"."int8_ops" ASC NULLS LAST
);

-- ----------------------------
-- Primary Key structure for table scan_prograss
-- ----------------------------
ALTER TABLE "scan_prograss" ADD CONSTRAINT "scan_prograss_pkey" PRIMARY KEY ("id");
