#!/bin/bash

# 修改项 start 
remoteuser="forest"
remotedir="/home/forest/forestAdminBin/server"
binname="forest-admin"
# 修改项 end


env="$1"
remotehost=""
if [ "$env" = "dev" ]; then
  remotehost="192.168.0.57"
elif [ "$env" = "test" ]; then
  remotehost="47.114.72.210"
elif [ "$env" = "prod" ]; then 
  remotehost="13.212.153.142"
else
  echo "没传入环境参数"
  exit 1
fi



echo "准备打包..."
GOOS=linux GOARCH=amd64 go build -o tempapp

echo "准备上传..."
scp tempapp $remoteuser@$remotehost:$remotedir/$binname

# 杀掉旧的，启动新的
echo '准备重启服务...'
ssh $remoteuser@$remotehost -t  << EOF
cd $remotedir 
ps -ef | grep $binname | grep -v grep | awk '{print \$2}' | xargs kill -9
nohup ./$binname server -c config/settings.$env.yml >> /dev/null 2>&1 &
mv $binname $binname.bak
exit
EOF

echo 'Done'
