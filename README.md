portrdr
=======

portrdr is a simple tcp/udp port redirection program.

Example
-------

A sample configuration (portrdr.json) that redirects port 80 of 192.168.0.1 to
www.example.com port 8080.

~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
[
	{
		"info": "Redirect to www.example.com",
		"protocol": "tcp",
		"localAddr": "192.168.0.1:80",
		"remoteAddr": "www.example.com:8080"
	}
]
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

ToDo
----

The following features will be implemented in the next updates ...

-   SSL redirects

-   HTTP(S) support
