## Run mvdparser on all demos in a relative directory

```shell
cd quake
# copy mvdparser, template.dat, fragfile.dat to dir
find . -type f -iregex ".*\.mvd" -exec ./mvdparser_linux_amd64 "{}" \;
```