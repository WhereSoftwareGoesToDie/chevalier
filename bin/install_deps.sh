# This is here so we can build with Travis; if you just want to use the
# library then you would probably be better off linking to a packaged
# version of libmarquise.

mkdir -p deps/
cd deps

wget http://download.zeromq.org/zeromq-4.0.4.tar.gz
tar -xf zeromq-4.0.4.tar.gz
cd zeromq-4.0.4
./configure 
make
sudo make install
sudo su -c "echo '/usr/local/lib' > /etc/ld.so.conf.d/local.conf"
sudo /sbin/ldconfig
cd ..

sudo apt-get install -y autoconf libtool automake build-essential libglib2.0-dev libprotobuf-c0-dev protobuf-c-compiler protobuf-compiler
go get code.google.com/p/goprotobuf/{proto,protoc-gen-go}
