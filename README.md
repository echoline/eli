# what is this?

meant to be a talking bash shell. helps the user a little bit to learn the basics and to read the manual pages. remembers responses to use them later.

![https://raw.githubusercontent.com/echoline/eli/refs/heads/main/scrn.png](https://raw.githubusercontent.com/echoline/eli/refs/heads/main/scrn.png)

# build

in client and server directories

	go mod init program
	go mod tidy
	go build

# run

first put the .rive files in ~/.config/rivescript/

then run the server program and put the client program in ~/.bashrc

	function command_not_found_handle {
		echo $* | client
	}
