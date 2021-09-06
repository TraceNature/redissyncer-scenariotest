# redissyncer-scenariotest

## 项目目的
本项目用于redissyncer 涉及场景的自动化测试，branch distribution用于分布式场景的自动化测试

## 项目架构

* caseyaml目录中的yaml文件用于描述场景用例，包括生成数据的批次，同步程序的位置等信息
* tasks目录用于存放任务创建json，用来创建各种同步任务,供case调用创建对应的任务