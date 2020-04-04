DIR=`pwd`
mv $DIR/initial_bood/build.bood $DIR
bood
mv $DIR/build.bood $DIR/initial_bood/
mv $DIR/custom_bood/build.bood $DIR
cd $DIR/out/bin/
mv bood $GOPATH/bin
