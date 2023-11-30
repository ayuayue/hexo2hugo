- [x] golang


#### 说明

使用 `golang` 原生编写,将 `hexo` 文章的 `front matter` 标签改为 `hugo` 可以解析的格式.
只过滤了常用的 `date tags categories` 标签, 如有其他定义的标签,可以自己在源码添加

#### 使用方法
下载`release` 或者`clone` 源码进行编译,在命令行运行即可

示例
```bash
# 复制几个文件来运行脚本测试
cp ~/wolfdan.cn/content/posts/hello-next test/
go run main.go test
# 针对正式目录执行
go run main.go ~/wolfdan.cn/content/posts/
# 在临时目录 newPost1701319905 中查看结果
# NOTICE: 强烈建议先生成到一个测试目录，修改hugo配置来测试hugo能生成再覆盖
mv newPost1701319905/* ~/wolfdan.cn/content/posts-test/
# NOTICE: 完全确认结果满意 覆盖 原来的目录
mv newPost1701319905/* ~/wolfdan.cn/content/posts/
# cd ~/wolfdan.cn && hugo 生成
```


#### 注意事项
没有成命令行接受参数的形式,为了防止文章修改失败而导致文件内容改变或丢失,请将要修改的`markdown` 文件复制到该应用的目录之下尽行操作

#### 支持的格式
目前只支持单个标签的更改并标签后要有空格,
支持的格式1:

    tags: hexo,blog
    
支持的格式2:

    tags: 
        - hexo
        - blog

----
优化
1. 写入到独立文件中,不更改元数据内容
