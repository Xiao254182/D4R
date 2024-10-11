## D4R - 使用 Docker CLI 来管理您的集群

D4R 以k9s项目为启发，取 docker 的首字母、尾字母和取中间4个字母的数量命名

D4R 提供了一个终端 UI 来与您的 docker 集群进行交互。该项目的目的是让您更轻松地在 docker 集群中发现服务、观察和管理您的应用程序。D4R 会持续监视 `docker ps` 的变化并集成了`logs`、`exec`、`rm`等命令来与您观察到的资源进行交互，也可以对 `docker-compose` 容器编排中的容器执行上述操作

部署：

```shell
bash <(wget -qO - https://raw.githubusercontent.com/Xiao254182/D4R/refs/heads/master/d4r_setup.sh)
```

使用：部署后直接在终端中执行d4r进入系统

