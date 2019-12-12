# proro

a personal streaming server, still developing...

usage
```
ffmpeg -re -i demo.flv -c copy -f flv rtmp://localhost:1935/live/movie

ffmpeg -re -i ".\example.OverwatchReplay.mp4" -c copy -f flv rtmp://localhost:1935/live/movie
```
## 设计功能

- 申请串流地址/播放地址
- 超时未使用自动关闭串流地址
- 播放地址不存在跳转至首页
- 停止串流（网页端）