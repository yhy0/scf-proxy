# scf-proxy
SCF代理
代码思路参考：https://freewechat.com/a/MzI0MDI5MTQ3OQ==/2247484068/1

下载源码后执行`make`，会生成两个文件
![image](https://user-images.githubusercontent.com/19854253/117415890-65303c00-af4b-11eb-84e6-e3720b16b513.png)
* server.zip，用于腾讯创建云函数时上传。
* client，用于本地代理的使用。

下面是图文教程，橙色标注的是需要修改的地方。

![image](https://user-images.githubusercontent.com/19854253/112804352-e96ae600-90a6-11eb-90dd-89a812d6dfa2.png)
![image](https://user-images.githubusercontent.com/19854253/112804449-04d5f100-90a7-11eb-88bd-0301c55e2399.png)
点击完成之后。这里有API的URL。
![image](https://user-images.githubusercontent.com/19854253/117416694-38305900-af4c-11eb-910e-2a474df886db.png)

执行
```bash
./client -api <API_URL> -port <PORT>
```
![image](https://user-images.githubusercontent.com/19854253/117416417-ef78a000-af4b-11eb-90bb-c49de8756d31.png)


然后走代理就完事了。
![image](https://user-images.githubusercontent.com/19854253/117416265-c1935b80-af4b-11eb-960d-9c4b73be1311.png)
