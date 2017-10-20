<?php
// JL8
// This one totally cheats by using someone else's easy-to-mirror archive
// Thanks, Axel! Way easier for me :)

// If the limbero.org archive one day falls out of sync or goes offline, I'll
// update this to pull from the actual Tumblr blog instead.

if (!is_dir('jl8')) {
	mkdir('jl8');
}

$url = 'http://limbero.org/jl8/comics/';
$html = file_get_contents($url);
preg_match_all('@href="([^"]+\\.jpe?g)"@', $html, $matches);
foreach ($matches[1] as $img) {
	if (!is_file('jl8/' . $img)) {
		echo "Downloading $img\n";
		$data = file_get_contents($url . $img);
		file_put_contents('jl8/' . $img, $data);
	}
}
