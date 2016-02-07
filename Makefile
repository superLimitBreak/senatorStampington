INSTALL_DIR=/opt/
SYSTEMD_DIR=/etc/systemd/system/
MAKEFILE_PATH=$$(dirname $$(realpath $(MAKEFILE_LIST)))

.PHONY:
help:
	# Make file for the installation of senatorStampington
	# install (requires you to pass GO_DIR):
	# 	install links to seneatorStampington in a linux system
	# 	this assumes it is already compiled
	# 	example: make install GO_DIR=/home/foo/go/bin
	# clean:
	#   remove all links relating to senatorStampington

.PHONY:
install:
ifeq ($(GO_DIR),)
	@echo "You must pass the location of your go folder via GO_DIR. See make help"
else
	ln -s $(GO_DIR)/senatorStampington $(INSTALL_DIR)
	ln -s $(MAKEFILE_PATH)/senatorStampington.service $(SYSTEMD_DIR)
endif

.PHONY:
clean:
	rm $(INSTALL_DIR)senatorStampington $(SYSTEMD_DIR)senatorStampington.service

