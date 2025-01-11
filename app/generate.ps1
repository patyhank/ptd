# Temp workaround for go workspace

Rename-Item .\go.work .\go1work
go generate ./ent
Rename-Item .\go1work .\go.work