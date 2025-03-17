/*
 Navicat Premium Dump SQL

 Source Server         : local-postgresql
 Source Server Type    : PostgreSQL
 Source Server Version : 170004 (170004)
 Source Host           : 10.0.0.61:5432
 Source Catalog        : sonex-chain-scanner
 Source Schema         : public

 Target Server Type    : PostgreSQL
 Target Server Version : 170004 (170004)
 File Encoding         : 65001

 Date: 15/03/2025 21:59:56
*/


-- ----------------------------
-- Table structure for common_event_weth_deposit
-- ----------------------------
DROP TABLE IF EXISTS "public"."common_event_weth_deposit";
CREATE TABLE "public"."common_event_weth_deposit" (
  "id" int8 NOT NULL,
  "chain_id" int8 NOT NULL,
  "block_number" int8 NOT NULL,
  "tx_hash" varchar(255) COLLATE "pg_catalog"."default" NOT NULL,
  "tx_at" timestamptz(6) NOT NULL,
  "tx_from" varchar(255) COLLATE "pg_catalog"."default" NOT NULL,
  "tx_to" varchar(255) COLLATE "pg_catalog"."default" NOT NULL,
  "tx_value" numeric(78,0) NOT NULL,
  "log_index" int8 NOT NULL,
  "log_address" varchar(255) COLLATE "pg_catalog"."default" NOT NULL,
  "event_dst" varchar(255) COLLATE "pg_catalog"."default" NOT NULL,
  "event_wad" numeric NOT NULL
)
;
ALTER TABLE "public"."common_event_weth_deposit" OWNER TO "minato";

-- ----------------------------
-- Table structure for common_event_weth_withdrawal
-- ----------------------------
DROP TABLE IF EXISTS "public"."common_event_weth_withdrawal";
CREATE TABLE "public"."common_event_weth_withdrawal" (
  "id" int8 NOT NULL,
  "chain_id" int8 NOT NULL,
  "block_number" int8 NOT NULL,
  "tx_hash" varchar(255) COLLATE "pg_catalog"."default" NOT NULL,
  "tx_at" timestamptz(6) NOT NULL,
  "tx_from" varchar(255) COLLATE "pg_catalog"."default" NOT NULL,
  "tx_to" varchar(255) COLLATE "pg_catalog"."default" NOT NULL,
  "tx_value" numeric(78,0) NOT NULL,
  "log_index" int8 NOT NULL,
  "log_address" varchar(255) COLLATE "pg_catalog"."default" NOT NULL,
  "event_src" varchar(255) COLLATE "pg_catalog"."default" NOT NULL,
  "event_wad" numeric NOT NULL
)
;
ALTER TABLE "public"."common_event_weth_withdrawal" OWNER TO "minato";

-- ----------------------------
-- Table structure for scan_prograss
-- ----------------------------
DROP TABLE IF EXISTS "public"."scan_prograss";
CREATE TABLE "public"."scan_prograss" (
  "id" int8 NOT NULL DEFAULT nextval('scan_prograss_id_seq'::regclass),
  "created_at" timestamptz(6) NOT NULL,
  "updated_at" timestamptz(6) NOT NULL,
  "chain_id" int8 NOT NULL,
  "height" int8 NOT NULL
)
;
ALTER TABLE "public"."scan_prograss" OWNER TO "minato";

-- ----------------------------
-- Table structure for sonex_event_factory_pool_created
-- ----------------------------
DROP TABLE IF EXISTS "public"."sonex_event_factory_pool_created";
CREATE TABLE "public"."sonex_event_factory_pool_created" (
  "id" int8 NOT NULL,
  "chain_id" int8 NOT NULL,
  "block_number" int8 NOT NULL,
  "tx_hash" varchar(255) COLLATE "pg_catalog"."default" NOT NULL,
  "tx_at" timestamptz(6) NOT NULL,
  "tx_from" varchar(255) COLLATE "pg_catalog"."default" NOT NULL,
  "tx_to" varchar(255) COLLATE "pg_catalog"."default" NOT NULL,
  "tx_value" numeric(78,0) NOT NULL,
  "log_index" int8 NOT NULL,
  "log_address" varchar(255) COLLATE "pg_catalog"."default" NOT NULL,
  "event_token0" varchar(255) COLLATE "pg_catalog"."default" NOT NULL,
  "event_token1" varchar(255) COLLATE "pg_catalog"."default" NOT NULL,
  "event_fee" int8 NOT NULL,
  "event_tick_spacing" int8 NOT NULL,
  "event_pool" varchar(255) COLLATE "pg_catalog"."default" NOT NULL
)
;
ALTER TABLE "public"."sonex_event_factory_pool_created" OWNER TO "minato";

-- ----------------------------
-- Table structure for sonex_event_pool_burn
-- ----------------------------
DROP TABLE IF EXISTS "public"."sonex_event_pool_burn";
CREATE TABLE "public"."sonex_event_pool_burn" (
  "id" int8 NOT NULL,
  "chain_id" int8 NOT NULL,
  "block_number" int8 NOT NULL,
  "tx_hash" varchar(255) COLLATE "pg_catalog"."default" NOT NULL,
  "tx_at" timestamptz(6) NOT NULL,
  "tx_from" varchar(255) COLLATE "pg_catalog"."default" NOT NULL,
  "tx_to" varchar(255) COLLATE "pg_catalog"."default" NOT NULL,
  "tx_value" numeric(78,0) NOT NULL,
  "log_index" int8 NOT NULL,
  "log_address" varchar(255) COLLATE "pg_catalog"."default" NOT NULL,
  "event_owner" varchar(255) COLLATE "pg_catalog"."default" NOT NULL,
  "event_tick_lower" int8 NOT NULL,
  "event_tick_upper" int8 NOT NULL,
  "event_amount" numeric(78,0) NOT NULL,
  "event_amount0" numeric(78,0) NOT NULL,
  "event_amount1" numeric(78,0) NOT NULL
)
;
ALTER TABLE "public"."sonex_event_pool_burn" OWNER TO "minato";

-- ----------------------------
-- Table structure for sonex_event_pool_collect
-- ----------------------------
DROP TABLE IF EXISTS "public"."sonex_event_pool_collect";
CREATE TABLE "public"."sonex_event_pool_collect" (
  "id" int8 NOT NULL,
  "chain_id" int8 NOT NULL,
  "block_number" int8 NOT NULL,
  "tx_hash" varchar(255) COLLATE "pg_catalog"."default" NOT NULL,
  "tx_at" timestamptz(6) NOT NULL,
  "tx_from" varchar(255) COLLATE "pg_catalog"."default" NOT NULL,
  "tx_to" varchar(255) COLLATE "pg_catalog"."default" NOT NULL,
  "tx_value" numeric(78,0) NOT NULL,
  "log_index" int8 NOT NULL,
  "log_address" varchar(255) COLLATE "pg_catalog"."default" NOT NULL,
  "event_owner" varchar(255) COLLATE "pg_catalog"."default" NOT NULL,
  "event_recipient" varchar(255) COLLATE "pg_catalog"."default" NOT NULL,
  "event_tick_lower" int8 NOT NULL,
  "event_tick_upper" int8 NOT NULL,
  "event_amount0" numeric(78,0) NOT NULL,
  "event_amount1" numeric(78,0) NOT NULL
)
;
ALTER TABLE "public"."sonex_event_pool_collect" OWNER TO "minato";

-- ----------------------------
-- Table structure for sonex_event_pool_mint
-- ----------------------------
DROP TABLE IF EXISTS "public"."sonex_event_pool_mint";
CREATE TABLE "public"."sonex_event_pool_mint" (
  "id" int8 NOT NULL,
  "chain_id" int8 NOT NULL,
  "block_number" int8 NOT NULL,
  "tx_hash" varchar(255) COLLATE "pg_catalog"."default" NOT NULL,
  "tx_at" timestamptz(6) NOT NULL,
  "tx_from" varchar(255) COLLATE "pg_catalog"."default" NOT NULL,
  "tx_to" varchar(255) COLLATE "pg_catalog"."default" NOT NULL,
  "tx_value" numeric(78,0) NOT NULL,
  "log_index" int8 NOT NULL,
  "log_address" varchar(255) COLLATE "pg_catalog"."default" NOT NULL,
  "event_sender" varchar(255) COLLATE "pg_catalog"."default" NOT NULL,
  "event_owner" varchar(255) COLLATE "pg_catalog"."default" NOT NULL,
  "event_tick_lower" int8 NOT NULL,
  "event_tick_upper" int8 NOT NULL,
  "event_amount" numeric(78,0) NOT NULL,
  "event_amount0" numeric(78,0) NOT NULL,
  "event_amount1" numeric(78,0) NOT NULL
)
;
ALTER TABLE "public"."sonex_event_pool_mint" OWNER TO "minato";

-- ----------------------------
-- Table structure for sonex_event_pool_swap
-- ----------------------------
DROP TABLE IF EXISTS "public"."sonex_event_pool_swap";
CREATE TABLE "public"."sonex_event_pool_swap" (
  "id" int8 NOT NULL,
  "chain_id" int8 NOT NULL,
  "block_number" int8 NOT NULL,
  "tx_hash" varchar(255) COLLATE "pg_catalog"."default" NOT NULL,
  "tx_at" timestamptz(6) NOT NULL,
  "tx_from" varchar(255) COLLATE "pg_catalog"."default" NOT NULL,
  "tx_to" varchar(255) COLLATE "pg_catalog"."default" NOT NULL,
  "tx_value" numeric(78,0) NOT NULL,
  "log_index" int8 NOT NULL,
  "log_address" varchar(255) COLLATE "pg_catalog"."default" NOT NULL,
  "event_sender" varchar(255) COLLATE "pg_catalog"."default" NOT NULL,
  "event_recipient" varchar(255) COLLATE "pg_catalog"."default" NOT NULL,
  "event_amount0" numeric(78,0) NOT NULL,
  "event_amount1" numeric(78,0) NOT NULL,
  "event_sqrt_price_x96" numeric NOT NULL,
  "event_liquidity" numeric NOT NULL,
  "event_tick" int8 NOT NULL
)
;
ALTER TABLE "public"."sonex_event_pool_swap" OWNER TO "minato";

-- ----------------------------
-- Table structure for sonex_event_position_manager_collect
-- ----------------------------
DROP TABLE IF EXISTS "public"."sonex_event_position_manager_collect";
CREATE TABLE "public"."sonex_event_position_manager_collect" (
  "id" int8 NOT NULL,
  "chain_id" int8 NOT NULL,
  "block_number" int8 NOT NULL,
  "tx_hash" varchar(255) COLLATE "pg_catalog"."default" NOT NULL,
  "tx_at" timestamptz(6) NOT NULL,
  "tx_from" varchar(255) COLLATE "pg_catalog"."default" NOT NULL,
  "tx_to" varchar(255) COLLATE "pg_catalog"."default" NOT NULL,
  "tx_value" numeric(78,0) NOT NULL,
  "log_index" int8 NOT NULL,
  "log_address" varchar(255) COLLATE "pg_catalog"."default" NOT NULL,
  "event_token_id" numeric NOT NULL,
  "event_recipient" varchar(255) COLLATE "pg_catalog"."default" NOT NULL,
  "event_amount0" numeric(78,0) NOT NULL,
  "event_amount1" numeric(78,0) NOT NULL
)
;
ALTER TABLE "public"."sonex_event_position_manager_collect" OWNER TO "minato";

-- ----------------------------
-- Table structure for sonex_event_position_manager_decrease_liquidity
-- ----------------------------
DROP TABLE IF EXISTS "public"."sonex_event_position_manager_decrease_liquidity";
CREATE TABLE "public"."sonex_event_position_manager_decrease_liquidity" (
  "id" int8 NOT NULL,
  "chain_id" int8 NOT NULL,
  "block_number" int8 NOT NULL,
  "tx_hash" varchar(255) COLLATE "pg_catalog"."default" NOT NULL,
  "tx_at" timestamptz(6) NOT NULL,
  "tx_from" varchar(255) COLLATE "pg_catalog"."default" NOT NULL,
  "tx_to" varchar(255) COLLATE "pg_catalog"."default" NOT NULL,
  "tx_value" numeric(78,0) NOT NULL,
  "log_index" int8 NOT NULL,
  "log_address" varchar(255) COLLATE "pg_catalog"."default" NOT NULL,
  "event_token_id" numeric NOT NULL,
  "event_liquidity" numeric NOT NULL,
  "event_amount0" numeric NOT NULL,
  "event_amount1" numeric NOT NULL
)
;
ALTER TABLE "public"."sonex_event_position_manager_decrease_liquidity" OWNER TO "minato";

-- ----------------------------
-- Table structure for sonex_event_position_manager_increase_liquidity
-- ----------------------------
DROP TABLE IF EXISTS "public"."sonex_event_position_manager_increase_liquidity";
CREATE TABLE "public"."sonex_event_position_manager_increase_liquidity" (
  "id" int8 NOT NULL,
  "chain_id" int8 NOT NULL,
  "block_number" int8 NOT NULL,
  "tx_hash" varchar(255) COLLATE "pg_catalog"."default" NOT NULL,
  "tx_at" timestamptz(6) NOT NULL,
  "tx_from" varchar(255) COLLATE "pg_catalog"."default" NOT NULL,
  "tx_to" varchar(255) COLLATE "pg_catalog"."default" NOT NULL,
  "tx_value" numeric(78,0) NOT NULL,
  "log_index" int8 NOT NULL,
  "log_address" varchar(255) COLLATE "pg_catalog"."default" NOT NULL,
  "event_token_id" numeric NOT NULL,
  "event_liquidity" numeric NOT NULL,
  "event_amount0" numeric NOT NULL,
  "event_amount1" numeric NOT NULL
)
;
ALTER TABLE "public"."sonex_event_position_manager_increase_liquidity" OWNER TO "minato";

-- ----------------------------
-- Table structure for sonex_event_position_manager_transfer
-- ----------------------------
DROP TABLE IF EXISTS "public"."sonex_event_position_manager_transfer";
CREATE TABLE "public"."sonex_event_position_manager_transfer" (
  "id" int8 NOT NULL,
  "chain_id" int8 NOT NULL,
  "block_number" int8 NOT NULL,
  "tx_hash" varchar(255) COLLATE "pg_catalog"."default" NOT NULL,
  "tx_at" timestamptz(6) NOT NULL,
  "tx_from" varchar(255) COLLATE "pg_catalog"."default" NOT NULL,
  "tx_to" varchar(255) COLLATE "pg_catalog"."default" NOT NULL,
  "tx_value" numeric(78,0) NOT NULL,
  "log_index" int8 NOT NULL,
  "log_address" varchar(255) COLLATE "pg_catalog"."default" NOT NULL,
  "event_from" varchar(255) COLLATE "pg_catalog"."default" NOT NULL,
  "event_to" varchar(255) COLLATE "pg_catalog"."default" NOT NULL,
  "event_token_id" numeric NOT NULL
)
;
ALTER TABLE "public"."sonex_event_position_manager_transfer" OWNER TO "minato";

-- ----------------------------
-- Indexes structure for table common_event_weth_deposit
-- ----------------------------
CREATE UNIQUE INDEX "common_event_w_d_idx_chain_hash_index" ON "public"."common_event_weth_deposit" USING btree (
  "chain_id" "pg_catalog"."int8_ops" ASC NULLS LAST,
  "tx_hash" COLLATE "pg_catalog"."default" "pg_catalog"."text_ops" ASC NULLS LAST,
  "log_index" "pg_catalog"."int8_ops" ASC NULLS LAST
);

-- ----------------------------
-- Primary Key structure for table common_event_weth_deposit
-- ----------------------------
ALTER TABLE "public"."common_event_weth_deposit" ADD CONSTRAINT "common_event_weth_deposit_pkey" PRIMARY KEY ("id");

-- ----------------------------
-- Indexes structure for table common_event_weth_withdrawal
-- ----------------------------
CREATE UNIQUE INDEX "common_event_w_w_idx_chain_hash_index" ON "public"."common_event_weth_withdrawal" USING btree (
  "chain_id" "pg_catalog"."int8_ops" ASC NULLS LAST,
  "tx_hash" COLLATE "pg_catalog"."default" "pg_catalog"."text_ops" ASC NULLS LAST,
  "log_index" "pg_catalog"."int8_ops" ASC NULLS LAST
);

-- ----------------------------
-- Primary Key structure for table common_event_weth_withdrawal
-- ----------------------------
ALTER TABLE "public"."common_event_weth_withdrawal" ADD CONSTRAINT "common_event_weth_withdrawal_pkey" PRIMARY KEY ("id");

-- ----------------------------
-- Indexes structure for table scan_prograss
-- ----------------------------
CREATE UNIQUE INDEX "idx_scan_prograss_chain_id" ON "public"."scan_prograss" USING btree (
  "chain_id" "pg_catalog"."int8_ops" ASC NULLS LAST
);

-- ----------------------------
-- Primary Key structure for table scan_prograss
-- ----------------------------
ALTER TABLE "public"."scan_prograss" ADD CONSTRAINT "scan_prograss_pkey" PRIMARY KEY ("id");

-- ----------------------------
-- Indexes structure for table sonex_event_factory_pool_created
-- ----------------------------
CREATE UNIQUE INDEX "sonex_event_f_p_c_idx_chain_hash_index" ON "public"."sonex_event_factory_pool_created" USING btree (
  "chain_id" "pg_catalog"."int8_ops" ASC NULLS LAST,
  "tx_hash" COLLATE "pg_catalog"."default" "pg_catalog"."text_ops" ASC NULLS LAST,
  "log_index" "pg_catalog"."int8_ops" ASC NULLS LAST
);

-- ----------------------------
-- Primary Key structure for table sonex_event_factory_pool_created
-- ----------------------------
ALTER TABLE "public"."sonex_event_factory_pool_created" ADD CONSTRAINT "sonex_event_factory_pool_created_pkey" PRIMARY KEY ("id");

-- ----------------------------
-- Indexes structure for table sonex_event_pool_burn
-- ----------------------------
CREATE UNIQUE INDEX "sonex_event_p_b_idx_chain_hash_index" ON "public"."sonex_event_pool_burn" USING btree (
  "chain_id" "pg_catalog"."int8_ops" ASC NULLS LAST,
  "tx_hash" COLLATE "pg_catalog"."default" "pg_catalog"."text_ops" ASC NULLS LAST,
  "log_index" "pg_catalog"."int8_ops" ASC NULLS LAST
);

-- ----------------------------
-- Primary Key structure for table sonex_event_pool_burn
-- ----------------------------
ALTER TABLE "public"."sonex_event_pool_burn" ADD CONSTRAINT "sonex_event_pool_burn_pkey" PRIMARY KEY ("id");

-- ----------------------------
-- Indexes structure for table sonex_event_pool_collect
-- ----------------------------
CREATE UNIQUE INDEX "sonex_event_p_c_idx_chain_hash_index" ON "public"."sonex_event_pool_collect" USING btree (
  "chain_id" "pg_catalog"."int8_ops" ASC NULLS LAST,
  "tx_hash" COLLATE "pg_catalog"."default" "pg_catalog"."text_ops" ASC NULLS LAST,
  "log_index" "pg_catalog"."int8_ops" ASC NULLS LAST
);

-- ----------------------------
-- Primary Key structure for table sonex_event_pool_collect
-- ----------------------------
ALTER TABLE "public"."sonex_event_pool_collect" ADD CONSTRAINT "sonex_event_pool_collect_pkey" PRIMARY KEY ("id");

-- ----------------------------
-- Indexes structure for table sonex_event_pool_mint
-- ----------------------------
CREATE UNIQUE INDEX "sonex_event_p_m_idx_chain_hash_index" ON "public"."sonex_event_pool_mint" USING btree (
  "chain_id" "pg_catalog"."int8_ops" ASC NULLS LAST,
  "tx_hash" COLLATE "pg_catalog"."default" "pg_catalog"."text_ops" ASC NULLS LAST,
  "log_index" "pg_catalog"."int8_ops" ASC NULLS LAST
);

-- ----------------------------
-- Primary Key structure for table sonex_event_pool_mint
-- ----------------------------
ALTER TABLE "public"."sonex_event_pool_mint" ADD CONSTRAINT "sonex_event_pool_mint_pkey" PRIMARY KEY ("id");

-- ----------------------------
-- Indexes structure for table sonex_event_pool_swap
-- ----------------------------
CREATE UNIQUE INDEX "sonex_event_p_s_idx_chain_hash_index" ON "public"."sonex_event_pool_swap" USING btree (
  "chain_id" "pg_catalog"."int8_ops" ASC NULLS LAST,
  "tx_hash" COLLATE "pg_catalog"."default" "pg_catalog"."text_ops" ASC NULLS LAST,
  "log_index" "pg_catalog"."int8_ops" ASC NULLS LAST
);

-- ----------------------------
-- Primary Key structure for table sonex_event_pool_swap
-- ----------------------------
ALTER TABLE "public"."sonex_event_pool_swap" ADD CONSTRAINT "sonex_event_pool_swap_pkey" PRIMARY KEY ("id");

-- ----------------------------
-- Indexes structure for table sonex_event_position_manager_collect
-- ----------------------------
CREATE UNIQUE INDEX "sonex_event_p_m_c_idx_chain_hash_index" ON "public"."sonex_event_position_manager_collect" USING btree (
  "chain_id" "pg_catalog"."int8_ops" ASC NULLS LAST,
  "tx_hash" COLLATE "pg_catalog"."default" "pg_catalog"."text_ops" ASC NULLS LAST,
  "log_index" "pg_catalog"."int8_ops" ASC NULLS LAST
);

-- ----------------------------
-- Primary Key structure for table sonex_event_position_manager_collect
-- ----------------------------
ALTER TABLE "public"."sonex_event_position_manager_collect" ADD CONSTRAINT "sonex_event_position_manager_collect_pkey" PRIMARY KEY ("id");

-- ----------------------------
-- Indexes structure for table sonex_event_position_manager_decrease_liquidity
-- ----------------------------
CREATE UNIQUE INDEX "sonex_event_p_m_d_l_idx_chain_hash_index" ON "public"."sonex_event_position_manager_decrease_liquidity" USING btree (
  "chain_id" "pg_catalog"."int8_ops" ASC NULLS LAST,
  "tx_hash" COLLATE "pg_catalog"."default" "pg_catalog"."text_ops" ASC NULLS LAST,
  "log_index" "pg_catalog"."int8_ops" ASC NULLS LAST
);

-- ----------------------------
-- Primary Key structure for table sonex_event_position_manager_decrease_liquidity
-- ----------------------------
ALTER TABLE "public"."sonex_event_position_manager_decrease_liquidity" ADD CONSTRAINT "sonex_event_position_manager_decrease_liquidity_pkey" PRIMARY KEY ("id");

-- ----------------------------
-- Indexes structure for table sonex_event_position_manager_increase_liquidity
-- ----------------------------
CREATE UNIQUE INDEX "sonex_event_p_m_i_l_idx_chain_hash_index" ON "public"."sonex_event_position_manager_increase_liquidity" USING btree (
  "chain_id" "pg_catalog"."int8_ops" ASC NULLS LAST,
  "tx_hash" COLLATE "pg_catalog"."default" "pg_catalog"."text_ops" ASC NULLS LAST,
  "log_index" "pg_catalog"."int8_ops" ASC NULLS LAST
);

-- ----------------------------
-- Primary Key structure for table sonex_event_position_manager_increase_liquidity
-- ----------------------------
ALTER TABLE "public"."sonex_event_position_manager_increase_liquidity" ADD CONSTRAINT "sonex_event_position_manager_increase_liquidity_pkey" PRIMARY KEY ("id");

-- ----------------------------
-- Indexes structure for table sonex_event_position_manager_transfer
-- ----------------------------
CREATE UNIQUE INDEX "sonex_event_p_m_t_idx_chain_hash_index" ON "public"."sonex_event_position_manager_transfer" USING btree (
  "chain_id" "pg_catalog"."int8_ops" ASC NULLS LAST,
  "tx_hash" COLLATE "pg_catalog"."default" "pg_catalog"."text_ops" ASC NULLS LAST,
  "log_index" "pg_catalog"."int8_ops" ASC NULLS LAST
);

-- ----------------------------
-- Primary Key structure for table sonex_event_position_manager_transfer
-- ----------------------------
ALTER TABLE "public"."sonex_event_position_manager_transfer" ADD CONSTRAINT "sonex_event_position_manager_transfer_pkey" PRIMARY KEY ("id");
