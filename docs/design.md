Commands
 
command
 
App is a route
 
Route.flags
 
C
 
App.Commands
 
Have documentation folder to auto load text files and spit on screen




Command


Route = zeal install name


Strings with no flags are considered command path


Flags parsed will be added in flags var


Sample path: install {{--packageName}}
When “zeal install service1” is entered. Service1 is auto binded to flag --packageName
