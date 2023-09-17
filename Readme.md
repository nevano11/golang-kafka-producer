docker compose exec broker kafka-topics --create --topic humans --bootstrap-server localhost:9092 --replication-factor 1 --partitions 1
github.com/confluentinc/confluent-kafka-go/v2/kafka