$exename = "focusrite-autoclock-windows.exe"
$startuppath = Join-Path $env:APPDATA 'Microsoft\Windows\Start Menu\Programs\Startup'
taskkill /IM $exename
go build -ldflags "-H=windowsgui" -o .\$exename .\src\.
Copy-Item $exename $startuppath
$exepath = Join-Path $startuppath $exename
Start-Process $exepath