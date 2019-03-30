<?php
// Sleepless Domain
// Another ComicPress one, still updating regularly.
// Mary Cagle is amazing and I have no idea how she maintains her separate
// comics with such a detailed art style. She's sorta Willis-y.

$html = file_get_contents('http://www.sleeplessdomain.com/');
preg_match('@/comics/(.+\\.(jpg|png|gif))@', $html, $matches);

if (empty($matches[1])) {
	echo "No comic found on home page! :(\n";
	return;
}

if (!is_dir('sleeplessdomain')) {
	mkdir('sleeplessdomain');
}

// Download current page's comic, load previous comic webpage, repeat
while (true) {
	if (is_file('sleeplessdomain/' . $matches[1])) {
		return;
	}

	echo "Downloading {$matches[1]}\n";
	$url = "http://www.sleeplessdomain.com/comics/{$matches[1]}";
	$data = @file_get_contents($url);
	if ($data) {
		file_put_contents("sleeplessdomain/{$matches[1]}", $data);
	}

    // Find previous page link
	$regex = '@rel="prev" href="(http://www.sleeplessdomain.com/comic/[0-9a-zA-Z-]+)"@';
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

	preg_match('@/comics/(.+\\.(jpg|png|gif))@', $html, $matches);
	if (empty($matches[1])) {
		echo "No image found on page!\n";
		return;
	}

	usleep(5e5);
}
