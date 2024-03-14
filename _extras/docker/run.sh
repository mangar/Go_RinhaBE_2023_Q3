

docker run -d --name api00 -p 3000:3000 \
    -e FLAG_DEBUG="True" \
    -e SERVER_NAME="SERVER_API00" \
    -e LOG_OUTPUT_DIR="./" \
    -e DB_CONNECTION="postgres://postgres:password@192.168.1.110:5432/postgres?sslmode=disable" \
    -e REDIS_CONNECTION="redis://:password@192.168.1.110:6379/0?protocol=3" \
    -e WEB_PORT="3000" \
 mangar/rinhabe_2023_q3_go:0.0.1
