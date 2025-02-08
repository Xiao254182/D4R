## D4R - 使用 Docker CLI 来管理您的集群

D4R 以k9s项目为启发，取 docker 的首字母、尾字母和取中间4个字母的数量命名

D4R 提供了一个终端 UI 来与您的 docker 集群进行交互。该项目的目的是让您更轻松地在 docker 集群中发现服务、观察和管理您的应用程序。D4R 会持续监视 `docker ps` 的变化并集成了`logs`、`exec`、`rm`等命令来与您观察到的资源进行交互

展示：

进入系统回自动检测是否有docker-compose服务，没有则不自动展示docker-compose部分

![](https://github.com/Xiao254182/D4R/blob/master/%E5%B1%95%E7%A4%BA/%E5%9F%BA%E6%9C%AC%E4%BD%BF%E7%94%A8.gif)

docker-compose部分

![](https://github.com/Xiao254182/D4R/blob/master/%E5%B1%95%E7%A4%BA/docker-compose%E5%B1%95%E7%A4%BA.gif)

部署：

```shell
bash <(wget -qO - https://raw.githubusercontent.com/Xiao254182/D4R/refs/heads/master/d4r_setup.sh)
```
(PS:如果网络不好，可以先将源代码下载到服务器中，然后执行 `local_setup.sh` 脚本即可)
```shell
https://github.com/Xiao254182/D4R/releases
```
