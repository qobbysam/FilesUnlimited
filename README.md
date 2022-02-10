

This is a file server written in GO.



This application uses Minio for the the data location.

You can serve files from data location by rest server

You can write files to data location  by rpc 




The Rpc and Fileservers can be started on different machines.

This can scale horizontally really well for both rpc and restserver.



###Basic Start

Create a copy of config.Json


Start with "./main -st="all" -c=path_to_config" 


-st
options
    rest
        restserver

options 
    rpc 
        rpc server

-c config is optional . 
