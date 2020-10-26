FROM alpine:latest
RUN mkdir -p /data/community-sig/bin
RUN mkdir -p /data/community-sig/log
#ADD ./bin/IMDA /data/bin
ADD  bin/community-sig  /data/community-sig/bin
COPY conf /data/community-sig/conf
EXPOSE 9699
EXPOSE 8092
ENV GRPC_ADDR  ""
# Redis数据库的地址: single_redis_host or cluster_redis_host，格式为：redisIMDA:7379
ENV REDIS_HOST ""
# mongodb数据库的地址：hosts，格式为：mongoIMDA:60000
ENV MONGODB_HOST ""
RUN echo http://mirrors.aliyun.com/alpine/v3.10/main/ > /etc/apk/repositories && \
    echo http://mirrors.aliyun.com/alpine/v3.10/community/ >> /etc/apk/repositories
RUN apk update && apk upgrade
RUN apk add --no-cache tzdata \
    && ln -snf /usr/share/zoneinfo/Asia/Shanghai /etc/localtime \
    && echo "Asia/Shanghai" > /etc/timezone
ENV TZ Asia/Shanghai
CMD ["/data/community-sig/bin/community-sig", "--conf","/data/community-sig/conf/communitysig_docker.toml"]

