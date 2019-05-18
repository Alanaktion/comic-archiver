<?php
// Gaia - The fantasy webcomic
// Another ComicPress one, still updating regularly, but ending soon (2019)
// This is hosted on the same domain as Sandra and Woo, in a sub-directory

$html = file_get_contents('http://www.sandraandwoo.com/gaia/');
preg_match('@/gaia/comics/(.+\\.(jpg|png|gif))@', $html, $matches);

if (empty($matches[1])) {
	echo "No comic found on home page! :(\n";
	return;
}

if (!is_dir('gaia')) {
	mkdir('gaia');
}

// Download current page's comic, load previous comic webpage, repeat
while (true) {
	if (is_file('gaia/' . $matches[1])) {
		return;
	}

	echo "Downloading {$matches[1]}\n";
	$url = "http://www.sandraandwoo.com/gaia/comics/{$matches[1]}";
	$data = @file_get_contents($url);
	if ($data) {
		file_put_contents("gaia/{$matches[1]}", $data);
	}

    // Find previous page link
	$regex = '@href="(http://www.sandraandwoo.com/gaia/[0-9]{4}/[0-9]+/[0-9]+/[0-9a-zA-Z-]+/?)" rel="prev"@';
	preg_match($regex, $html, $prevMatch);

	if (empty($prevMatch[1])) {
		echo "No previous URL found!\n";
		return;
	}

	$html = @file_get_contents($prevMatch[1]);
	if (!$html) {
		echo "Failed to load previous page!\n";
		return;
	}

	preg_match('@/gaia/comics/(.+\\.(jpg|png|gif))@', $html, $matches);
	if (empty($matches[1])) {
		echo "No image found on page!\n";
		return;
	}

	usleep(5e5);
}
