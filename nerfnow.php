<?php
// Nerf NOW!!
// For some reason the image numbers don't match the comic numbers, but we'll
// work around that by renaming the files I guess. I didn't want to do that, but
// we don't have a good alternative.

// This should be safe to re-run multiple times to catch missed images from a
// previous run, though that hasn't actually been thoroughly tested.

$html = file_get_contents('http://www.nerfnow.com/');
preg_match('@/img/([0-9]+)/([0-9]+\\.[a-z]{3,4})@', $html, $matches);

if (!$matches[1]) {
	echo "No comic found on home page! :(\n";
	return;
}

if (!is_dir('nerfnow')) {
	mkdir('nerfnow');
}

$start = $matches[1];

for ($i = $start; $i > 0; $i--) {
	if (glob("nerfnow/$i-*")) {
		continue;
	}

	echo "Downloading #$i\n";
	$html = file_get_contents('http://www.nerfnow.com/comic/' . $i);
	preg_match('@/img/([0-9]+)/([0-9]+\\.[a-z]{3,4})@', $html, $matches);

	$url = "http://www.nerfnow.com/img/$i/{$matches[2]}";
	$data = @file_get_contents($url);
	if ($data) {
		file_put_contents("nerfnow/$i-{$matches[2]}", $data);
	}

	usleep(500000);
}
