# 物联网视频监控服务

服务整体介绍,请点击 [链接](./docs/project_introduce.md)

## 此项目使用前准备

* 新建yml格式配置文件,比如在 conf文件夹下新建 local.yml,配置如下

```yaml
App:
  RunMode: debug
  EnablePProf: false
#  EnablePProf: true

# 注意新建 log文件夹
LogConf:
  FilePath: "./log/video_server.log"

MqttConf:
  UserName: "admin"
  Password: "gRiQ8DwyPGBL7MN"
  Host: "127.0.0.1"
  Port: 1883

UDPConf:
  Port: 9090

```

* 设置环境变量
  ```LOC_CFG=/绝对路径/video_server/conf/local.yml```
