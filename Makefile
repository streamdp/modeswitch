PATH := $(GOPATH)/bin:$(PATH)

default:
	fyne-cross android -image fyneio/fyne-cross-images:v1.2.0-android

clean:
	rm -rf *.apk
	rm -rf fyne-cross/
	rm -rf tmp-pkg/