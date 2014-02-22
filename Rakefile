JARFILE = 'kafka/core/target/scala-2.8.0/kafka_2.8.0-0.8.1.jar'

def sh_cd(cmd)
  sh "cd kafka && #{cmd}"
end

directory 'kafka' do
  sh "git clone https://github.com/apache/kafka.git"
end

file JARFILE => [ :kafka ]do
  %w(update package assembly-package-dependency).each do |cmd|
    sh_cd "./sbt #{cmd}"
  end
end

task :build => [ JARFILE ]

namespace :topic do
  task :create do
    # NOTE this reflects changes currently in kafka trunk, 0.8 uses different commands
    sh_cd 'bin/kafka-topics.sh --create --zookeeper localhost:2181 --replication-factor 1 --partition 1 --topic test'
  end
end

namespace :agent do
  task :producer do
    sh_cd 'bin/kafka-console-producer.sh --broker-list localhost:9092 --topic test'
  end

  task :consumer do
    sh_cd 'bin/kafka-console-consumer.sh --zookeeper localhost:2181 --topic test --from-beginning'
  end
end

namespace :run do
  task :zk do
    sh_cd 'bin/zookeeper-server-start.sh config/zookeeper.properties'
  end

  task :kafka do
    sh_cd 'bin/kafka-server-start.sh config/server.properties'
  end
end
