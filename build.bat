@echo off



echo == Preparation > compil
ibt -stats preps.ibt >> compil
echo = >> compil
echo == Pink Samurai >> compil
ibt -stats ps.ibt >> compil



ibt -begin preps.ibt

rem go get

minify home.css -o min/home.css
minify home.js -o min/home.js
rem mess up my html
rem minify home.html -o home.min.html
rem minify login.html -o min/login.html

minify bonus/zalgo.js -o bonus/zalgo.min.js
minify bonus/konami.js -o bonus/konami.min.js
minify bonus/bonus.css -o bonus/bonus.min.css

minify fonts/LiberationMono.css -o fonts/LiberationMono.min.css

minify compatibility/ie.js -o compatibility/ie.min.js

go-bindata favicon.ico compil min info.html 404.html 404 home.html login.html PFUDOR.mp3 bonus bonus/glitch bonus/gay vortex-1.png vortex-2.png vortex-3.png background.jpg ChangeLog tentacles.png fonts compatibility emojione emojione/svg
set LastError=%ERRORLEVEL%
ibt -end preps.ibt %LastError%




ibt -begin ps.ibt
set /p ver=<build
set /a ver=ver
set /a ver=ver+1
echo %ver% > build
echo Build: %ver%
go build -ldflags "-X main.build=%ver% -X main.version=0.2"
set LastError=%ERRORLEVEL%
ibt -end ps.ibt %LastError%
