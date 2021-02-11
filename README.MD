# 12306 后端微服务结构

## 目录结构和功能

12306A/&emsp;&emsp;&emsp;&emsp;&emsp;&emsp;&emsp;12306后端A小组</br>
|------rpc&emsp;&emsp;&emsp;&emsp;&emsp;&emsp;&emsp;grpc相关的接口和协议文件</br>
|&emsp;&emsp; |------pay&emsp;&emsp; &emsp;&emsp;pay服务器的rpc代码, 同理如果是user服务应该在该文件夹下建立user文件夹</br>
|&emsp;&emsp;&emsp;&emsp; |------proto&emsp;&emsp;.proto文件存放</br>
|&emsp;&emsp;&emsp;&emsp; |------client&emsp;&emsp;grpc客户端, grpc服务再server中自己实现</br>
|------server&emsp;&emsp;&emsp;&emsp;&emsp;&emsp;每个微服务项目</br>
|&emsp;&emsp; |------candidate&emsp;&emsp;候补服务器</br>
|&emsp;&emsp; &emsp;&emsp; |------controller&emsp;&emsp;控制层,数据的接受的校验</br>
|&emsp;&emsp; &emsp;&emsp; |------service&emsp;&emsp;服务层,业务逻辑</br>
|&emsp;&emsp; &emsp;&emsp; |------model&emsp;&emsp;模型层,与数据库连接</br>
|&emsp;&emsp; &emsp;&emsp; |------redis&emsp;&emsp;缓存连接</br>
|&emsp;&emsp; &emsp;&emsp; |------setting&emsp;&emsp;配置服务</br>
|&emsp;&emsp; &emsp;&emsp; |------config&emsp;&emsp;配置文件存放</br>
|&emsp;&emsp; |------pay&emsp;&emsp;&emsp;&emsp; 支付服务器</br>
|&emsp;&emsp; |------reticket&emsp;&emsp;&emsp;退票服务器</br>
|&emsp;&emsp; |------search&emsp;&emsp;&emsp;搜索</br>
|&emsp;&emsp;&emsp;&emsp; |------dynamic&emsp;&emsp;动态搜索</br>
|&emsp;&emsp;&emsp;&emsp; |------static&emsp;&emsp;静态搜索</br>
|&emsp;&emsp; |------ticket&emsp;&emsp;&emsp;&emsp;购票服务器</br>
|&emsp;&emsp; |------user&emsp;&emsp;&emsp;&emsp; 用户服务器</br>
|------ticketPool&emsp;&emsp;&emsp;&emsp; 线程池服务,主要是对内提供服务</br>
|------interface&emsp;&emsp;&emsp;&emsp; 每个服务能够提供的接口</br>
|------implement&emsp;&emsp;&emsp;&emsp; 每个服务的接口实现</br>


this is dev brench

## 说明
1. 项目使用consul作为服务发现
2. 使用nginx将https请求转为内部http请求
3. 使用自建gin反向代理作为token鉴权并进行为url加入用户信息
4. 使用fabio作为gin反向代理之后的三级代理,目的是做到负载均衡和按照服务名转发