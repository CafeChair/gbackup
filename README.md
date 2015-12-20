##gobackup
原理

                            backup server
                                |
                                |
                            redis cluster
                                |
                                |
                            backup client
    1 script脚本备份指定数据库(备份格式:ipaddress-date-type.tar.gz)
    2 gobackup连接redis集群存入格式(key:rsyncmod,value：备份格式)
    3 Server端连接redis集群拉取格式(key:rsyncmod,value：备份格式)
    4 Server端通过rsync方式，拉取备份文件(rsync -azcv rsyncmod::备份格式 .)

配置文件(client)

    {
        "script":{
            //当前目录script下放置脚本文件,脚本文件格式:60_script.sh
            //60表示script目录下的脚本文件每隔多长时间循环执行,单位默认是分钟
            "dir":"./script",
            //当前目录logs下是脚本文件执行的日志记录
            "logdir": "./logs"
        }
    }

配置文件(server)