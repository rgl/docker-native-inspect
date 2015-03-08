Inspects the native Docker driver container state.

Normally, the driver stores its data at:

	/var/lib/docker/execdriver/native/<container id>/state.json
	/var/lib/docker/execdriver/native/<container id>/container.json


# Example usages

Show the whole container state:

	sudo ./docker-native-inspect state container-name

Show just the `veth_host` property:

	sudo ./docker-native-inspect -format '{{.network_state.veth_host}}' state container-name

Show the whole container configuration:

	sudo ./docker-native-inspect container container-name

Show just the mounts:

	sudo ./docker-native-inspect -format '{{range .mount_config.mounts}}{{printf "%v -> %v\n" .source .destination}}{{end}}' container container-name


**NB** The syntax of the template language used on the `-format` argument is described at the Go [text/template package documentation page](http://golang.org/pkg/text/template/).
