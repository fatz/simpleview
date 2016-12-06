# SimpleView
a small go project for having a the current state of your icinga2 monitoring on a screen.

It needs an user/password to the icinga api.

## State
this is a really small and fast project the current state works for us but might need some improvement. Let me know if you've ideas or just send pull requests.

```
target/simpleview.darwin.amd64 -help
Usage of target/simpleview.darwin.amd64:
  -icingaendpoint="http://api-icinga2.mycompany.example.com:5665": the url to the icinga2 api
  -icingapassword="password": password
  -icingausername="simpleview": username to authenticate against icinga
  -projects="": filter on this comma-sep. list of projects
```

__projects:__ the projects options is currently not implemented

All commandline options can also be set through environment variables `SIMPLEVIEW_$var`

# License
```
----------------------------------------------------------------------------
"THE BEER-WARE LICENSE" (Revision 42):
<janu@cpan.org> wrote this file.  As long as you retain this notice you
can do whatever you want with this stuff. If we meet some day, and you think
this stuff is worth it, you can buy me a beer in return.   Jan Ulferts
----------------------------------------------------------------------------
```
