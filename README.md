PI - Process Investigator
=========

![pi](https://raw.githubusercontent.com/adamar/PI/master/doc/pi.jpg)

Still under heavy development!


About	
--------------
PI is a tool for inspecting linux processes


Components
-------------
Currently there is a main Process Investigator which pulls information
from the /proc directory
 
 -  pi.go           Display information from the /proc/{pid} directory

There are also a number of secondary tools for inspecting different elements of a running process

 - pi-net.go:       Tracks network system calls
 - pi-gile.go:      Tracks filesystem system calls
 - pi-mem-grep.go   Greps through a processes assigned stack for a string
