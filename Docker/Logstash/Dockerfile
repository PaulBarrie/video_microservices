FROM logstash:7.9.3

USER root
RUN rm -f /usr/share/logstash/pipeline/logstash.conf
# RUN curl -L -o /usr/share/logstash/lib/mysql-connector-java-5.1.47.jar https://repo1.maven.org/maven2/mysql/mysql-connector-java/5.1.47/mysql-connector-java-5.1.47.jar
ADD ./pipeline/ /usr/share/logstash/pipeline/
ADD ./config /usr/share/logstash/config/
ADD mysql-connector-java-8.0.21.jar /usr/share/logstash/logstash-core/lib/jars/mysql/mysql-connector-java-8.0.21.jar

RUN logstash-plugin install logstash-input-jdbc