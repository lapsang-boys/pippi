all: fetch

fetch: src

VER=1.4.6

src: v${VER}.tar.gz
	tar -x ace-builds-${VER}/src -z -v -f $<
	mv ace-builds-${VER}/src src
	rmdir ace-builds-${VER}

v${VER}.tar.gz:
	wget https://github.com/ajaxorg/ace-builds/archive/v${VER}.tar.gz
