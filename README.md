# WIP symmetrical-pancake

Simple web document viewer. Allows to navigate a library of files, page by page (like comics). From a local source

### Perequisites

    go
    docker


### Setup   

Start the database 

    docker-compose up

Build viewer 

    go build .\viewer.go

Execute viewer

     .\viewer.exe  

Create library folder at root

    mkdir library

Put subfolder and file structure under the library folder (fileserver location)

Populate database with title -> path mapping

### Usage

Access specific file by opening localhost:8081 in a browser (with mappings)

    localhost:8081/reader/<title>/<page>

Or just access the fileserver (no mapping)

    localhost:8081/library

### Next steps :

-Adding navigation between pages
-Add support for not just jpg
-Adding a proper visualization too
-Adding tagging and searches
-Adding scripts for populating database