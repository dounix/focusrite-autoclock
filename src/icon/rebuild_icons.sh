unixhead='//+build linux darwin'
winhead='//+build windows'

for file in *.png ; do
  filesuf="${file%.*}"
  convert $file -define icon:auto-resize=256,64,48,32,16 ${filesuf}.ico

  echo $unixhead > ../${filesuf}unix.go
  #cat $file | 2goarray Data ${filesuf}icon >> ${filesuf}unix.go
  cat $file | 2goarray Data$filesuf main >> ../${filesuf}unix.go

  echo $winhead > ../${filesuf}win.go
  #cat ${filesuf}.ico| 2goarray Data ${filesuf}icon >> ${filesuf}win.go
  cat ${filesuf}.ico| 2goarray Data$filesuf main >> ../${filesuf}win.go

done
