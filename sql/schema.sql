/*
 Navicat Premium Data Transfer

 Source Server         : mysql-tencent
 Source Server Type    : MySQL
 Source Server Version : 80032 (8.0.32)
 Source Host           : 193.112.178.249:3306
 Source Schema         : countdown

 Target Server Type    : MySQL
 Target Server Version : 80032 (8.0.32)
 File Encoding         : 65001

 Date: 26/03/2024 10:05:05
*/

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for countdowns
-- ----------------------------
DROP TABLE IF EXISTS `countdowns`;
CREATE TABLE `countdowns`  (
  `id` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '雪花算法生成id',
  `query_code` varchar(12) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '查询码',
  `expire_at` datetime NOT NULL COMMENT '到期时间',
  `created_at` datetime NOT NULL COMMENT '创建时间',
  `updated_at` datetime NOT NULL ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `message` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT '消息（倒计时结束返回）',
  `remark` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT '备注信息',
  `deleted` tinyint NOT NULL DEFAULT 0 COMMENT '逻辑删除',
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `query_code_index`(`query_code` ASC) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci COMMENT = '倒计时信息表' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Table structure for mail_template
-- ----------------------------
DROP TABLE IF EXISTS `mail_template`;
CREATE TABLE `mail_template`  (
  `id` int NOT NULL,
  `topic` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '邮件主题',
  `format` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '预定义格式',
  `created_at` datetime NOT NULL COMMENT '创建时间',
  `updated_at` datetime NOT NULL ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `deleted` tinyint NOT NULL DEFAULT 0 COMMENT '逻辑删除',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci COMMENT = '邮件模板表' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Table structure for notify_tasks
-- ----------------------------
DROP TABLE IF EXISTS `notify_tasks`;
CREATE TABLE `notify_tasks`  (
  `id` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '雪花算法生成id',
  `cid` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '倒计时id',
  `tid` int NOT NULL COMMENT '模板id',
  `notify_type` int NOT NULL DEFAULT 1 COMMENT '通知类型，1提前2准时3阶段性',
  `status` tinyint NOT NULL DEFAULT 0 COMMENT '发送状态，0未发送，1发送成功',
  `sender` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '发送者邮箱',
  `receiver` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '收信人邮箱',
  `notify_at` datetime NOT NULL COMMENT '通知时间',
  `created_at` datetime NOT NULL COMMENT '创建时间',
  `updated_at` datetime NOT NULL ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `deleted` tinyint NOT NULL DEFAULT 0 COMMENT '逻辑删除',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci COMMENT = '定时通知任务表' ROW_FORMAT = Dynamic;

SET FOREIGN_KEY_CHECKS = 1;
