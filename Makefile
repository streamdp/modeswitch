default:
	fyne-cross android -image fyneio/fyne-cross-images:v1.2.0-android -no-cache

clean:
	rm -rf *.apk
	rm -rf fyne-cross/
	rm -rf tmp-pkg/