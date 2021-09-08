# redissyncer-scenariotest

## 项目目的

本项目用于redissyncer 涉及场景的自动化测试，branch distribution用于分布式场景的自动化测试

## 项目架构

* caseyaml目录中的yaml文件用于描述场景用例，包括生成数据的批次，同步程序的位置等信息
* tasks目录用于存放任务创建json，用来创建各种同步任务,供case调用创建对应的任务

# 场景描述

* 单实例同步
* 单实例带映射关系的同步
* 单实例到原生集群同步
* 原生集群到原生集群同步
* 单实例同步+断点续传
* 单实例带映射关系的同步+断点续传
* 单实例到原生集群同步+断点续传
* 原生集群到原生集群同步+断点续传
* 任务节点shutdown情况下的任务迁移