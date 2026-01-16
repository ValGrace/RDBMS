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


SQL Compiler -> 