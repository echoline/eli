# what is this?

meant to be a talking bash shell. helps the user a little bit to learn the basics and to read the manual pages. remembers responses to use them later.

first put the .rive files in ~/.config/rivescript/

then run the server and put the client in ~/.bashrc

	function command_not_found_handle {
		echo $* | client
	}
