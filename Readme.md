# Important features
- The Database system has to be relational ~ (supports tables, joins, SQL)
- Object relational Extensibility ~ (support for custom types, inheritance, stored procedures)
- ACID properties
- CRUD operations
- SQL parsing

# Why I built an Embedded Database
~ Sits on the application layer reducing the need for a separate database server. This means I can store data in a local file
# Logging 
Zerolog is minimalistic and the best logger for Embedded DBMS

# REPL
~ Added a simple CLI engine to easily interface with REPL
~ Create a cross-platform CLI tool with live validation using prompts

# Table Implementation

~ use of B-Tree which can allow each node to have multiple children
    ~ B-Tree supports Operations like Search, Insert, Delete, Traverse making it highly efficient for database systems.
      Rules 
      1. The leaves should be at the same level of the tree for balance
      2. Every node has a maximum and minimum number of keys ( min = max / 2)
      3. The root node can have fewer keys
      4. Bottom up creation process

Every key has a record pointer

SQL Compiler -> 

** Considerations
1. Use of indexes and caches
2. The order of table joins
3. Concurrency control
4. Transaction management

Parse the SQL statement -> transform the SQL into a relational representation -> create an execution plan that utilizes index info -> return results