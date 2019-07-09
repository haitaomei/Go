

Run cassandra using docker
----
Run 

`docker run --rm --name cassandra-node1 -d -p 9042:9042 cassandra:3.11`


Run

`docker exec -it $(docker ps | grep cassandra-node1 | awk '{print $1}') cqlsh`

to Connect to Cassandra from **cqlsh**


Prepare the database
----
Create a key space using

`CREATE KEYSPACE exampleks WITH REPLICATION = { 'class' : 'SimpleStrategy', 'replication_factor' : 1 };`

create a table

    use exampleks;
    CREATE TABLE users (
      usrid text PRIMARY KEY,
      first_name text,
      last_name text,
      age int
    );

After the experiment, using `Drop keyspace exampleks;` to delete the database.

Ref
------
https://docs.docker.com/samples/library/cassandra/