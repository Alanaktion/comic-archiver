<?php
// xkcd
// Not too bad to archive, get the current ID, work in reverse using the API.

$current = json_decode(file_get_contents('https://xkcd.com/info.0.json'));

if (!is_dir('xkcd')) {
	mkdir('xkcd');
}

for ($i = $current->num; $i > 0; $i--) {
	if (glob('xkcd/' . $i . '-*')) {
		continue;
	}
	$meta = json_decode(file_get_contents("https://xkcd.com/$i/info.0.json"));

	echo "Downloading #$i - " . basename($meta->img), "\n";
	$data = @file_get_contents($meta->img);
	if ($data) {
		file_put_contents("xkcd/$i-" . basename($meta->img), $data);
	}
	if ($i > 1084) {
		// Download @2x img
		$data = @file_get_contents(str_replace('.png', '_2x.png', $meta->img));
		if ($data) {
			file_put_contents("xkcd/$i-" . str_replace('.png', '_2x.png', basename($meta->img)), $data);
		}
	}

	usleep(500000);
}
